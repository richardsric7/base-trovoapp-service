package users

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"trovo-wallet-api/internal/basetxn"
	tPayErrors "trovo-wallet-api/internal/components/payments/errors"
	usersDB "trovo-wallet-api/internal/components/users/db"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/evmkeypair"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func generateMintingXdr(client *ethclient.Client, owner *userModels.User, sourceWallet *userModels.UserWallet, mintingInfo *userModels.MintingInfo, db *gorm.DB, gc *sharedconfig.GlobalConfig) (string, *userModels.User, error) {
	var err error
	mintingInfo, err = ValidateMintingInfo(mintingInfo)
	nativeAssetCode := os.Getenv("NATIVE_ASSET_CODE")
	if err != nil {
		return "", nil, err
	}
	// amountToSend := decimal.RequireFromString(mintingInfo.Amount).InexactFloat64()
	// amountToSendDec := decimal.RequireFromString(mintingInfo.Amount)
	newAmountToSend := decimal.RequireFromString(mintingInfo.Amount).Truncate(7).String()
	asset := basetxn.CreditAsset{Code: mintingInfo.AssetCode, Issuer: mintingInfo.ContractAddress}
	if sourceWallet.ID != asset.GetIssuer() {

		return "", nil, &tErrors.CustomError{
			Param:      "assetCode",
			Err:        "error cannot mint token you are not an issuer of",
			ErrMessage: "You can only mint tokens you are an issuer of",
		}

	}
	destinationInfo, getDestinationError := usersDB.GetUser(mintingInfo.Destination, db, gc)
	destinationWallet, destinationWalletError := usersDB.GetWallet(mintingInfo.Destination, db)

	if (getDestinationError != nil || destinationWalletError != nil) && len(mintingInfo.Destination) != 42 {
		return "", nil, &tPayErrors.ErrorPaymentDestinationDoesNotExist{}
	}
	if destinationInfo.Suspended == 1 {
		return "", nil, &tErrors.ErrorUsernameIsSuspended{}
	}

	mintingInfo.DestinationFirstName = destinationInfo.FirstName
	if destinationInfo.LastName != nil {
		mintingInfo.DestinationLastName = *destinationInfo.LastName
	}

	mintingInfo.DestinationVerified = destinationInfo.Verified

	if destinationInfo.ImageThumbnailURL != nil {
		mintingInfo.DestinationThumbnail = *destinationInfo.ImageThumbnailURL
	}

	// var destinationAddress string

	destinationAddress := destinationWallet.ID

	//perform ths checks of determining messages to be appended. if destination account property is not checked here, information would be returned without messages set.
	_, destinationAccountTrustsAsset, _, _, destinationBlockchainAccount, destinationAccountErr :=
		network.BlockchainAccountProperties(client, destinationAddress, asset)
		//set base charge to be used in all places it is needed

		//check if to set charges messages

	// if !destinationAccountTrustsAsset && !owner.IsEnterpriseProfile(gc) {
	// 	message := fmt.Sprintf("%v has not yet opted in to receive the asset (%v) you are trying to send. %v %v will be deducted from your account to ensure that this transaction goes through. After this, %v will be able to receive %v anytime, without any further charges to you.", destinationWallet.Alias, mintingInfo.AssetCode, charge, nativeAssetCode, destinationWallet.Alias, mintingInfo.AssetCode)

	// 	mintingInfo.Messages = append(mintingInfo.Messages, message)
	// 	// log.Printf("[generatePaymentXdr]message[1]: %v\n", message)

	// }

	chanAccount, releaseChanAccount, errCheckout := sharedconfig.CheckoutChannelAccount(gc)
	if errCheckout != nil {
		return "", nil, errCheckout
	}
	defer releaseChanAccount()
	// paymentInfo.Messages = messages
	_, _, _, _, chanSourceAccount, _ := network.BlockchainAccountProperties(client, chanAccount.Address(), basetxn.NativeAsset{})
	_, sourceAccountTrustsAsset, sourceAccountNativeBalance, sourceAccountCustomBalance, sourceAccount, sourceAccountErr := network.BlockchainAccountProperties(client, sourceWallet.ID, asset)

	if sourceAccountErr != nil {
		return "", nil, sourceAccountErr
	}

	if !sourceAccountTrustsAsset {

		return "", nil, &tErrors.ErrorUnderfundedAccount{}
	}

	log.Printf("[generateMintingXdr]obtained source account balance:\n%v balance is %v\n%v balance is %v\n", nativeAssetCode, sourceAccountNativeBalance, asset.GetCode(), sourceAccountCustomBalance)

	//check if destination account exists
	if destinationAccountErr != nil {
		log.Println("[generateMintingXdr]destination Account error:", destinationAccountErr)
		return "", nil, destinationAccountErr
	}
	var ops []basetxn.Operation = make([]basetxn.Operation, 0)

	var extraAccountKeyPair *evmkeypair.Full = nil

	//custom asset

	if !destinationAccountTrustsAsset {

		//meaning that destinationWallet and destinationUser objects are valid.
		if destinationWallet.WalletType == 1 {
			//asset issuing wallet is forbidden to receive custom assets. only native assets
			err = &tErrors.CustomError{
				Param:      "destination",
				Err:        "error-destination-forbidden-to-receive-asset",
				ErrMessage: fmt.Sprintf("%v, a token minting wallet, is forbidden from receiving %v.", destinationWallet.Alias, asset.GetCode()),
			}
			return "", nil, err
		}

		// Only market-ready tokenized/regulated assets ever reach here
		// (see network.IsWalletAuthorizedForAsset) - every other B20
		// asset is always authorized on Base, no opt-in step needed. A
		// regulated asset needs an explicit compliance approval (POST
		// /v1/compliance/wallet-authorization, signed by the asset's own
		// issuing wallet) before it can be minted to any destination,
		// custodial wallets included - so this rejects rather than
		// minting around the gate.
		return "", nil, &tErrors.CustomError{
			Param:      "destination",
			Err:        "error-destination-cannot-accept-asset",
			ErrMessage: fmt.Sprintf("%v is not authorized to receive the asset %v.", destinationWallet.Alias, asset.GetCode()),
		}

	}
	ops = append(ops, &basetxn.Payment{
		Destination:   destinationAddress,
		Amount:        newAmountToSend,
		Asset:         asset,
		SourceAccount: sourceWallet.ID,
	})

	var tx *basetxn.Transaction
	// Construct the transaction that holds the operations to execute on the network
	if mintingInfo.Multiparty == 1 {
		mintingInfo.TransactionSource = chanSourceAccount.Address
		tx, err = basetxn.NewTransaction(
			basetxn.TransactionParams{
				SourceAccount:        chanSourceAccount.Address,
				IncrementSequenceNum: true,
				Operations:           ops,
				BaseFee:              2000,
				Memo:                 mintingInfo.Memo,
			},
		)
	} else {
		tx, err = basetxn.NewTransaction(
			basetxn.TransactionParams{
				SourceAccount:        sourceAccount.Address,
				IncrementSequenceNum: true,
				Operations:           ops,
				BaseFee:              2000,
				Memo:                 mintingInfo.Memo,
			},
		)
	}

	if err != nil {
		log.Println("[generateMintingXdr] error constructing transaction ", err)
		return "", nil, err
	}

	if mintingInfo.Multiparty == 1 {

		tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), chanAccount)

		if err != nil {
			log.Println("[generateMintingXdr] error signing transaction with channelAccount key ", err)
			return "", nil, &tErrors.ErrorTemporaryServerError{}
		}
	}

	if extraAccountKeyPair != nil {

		tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), extraAccountKeyPair)

		if err != nil {
			log.Println("[generateMintingXdr] error signing transaction with temporary key ", err)
			return "", nil, &tErrors.ErrorTemporaryServerError{}
		}
	}

	var xdrBase64 string

	xdrBase64, err = tx.Base64()
	if err != nil {
		log.Println("[generateMintingXdr] error getting txn base64", err)
		return "", nil, err
	}

	if destinationBlockchainAccount != nil {
		mintingInfo.CallbackURLS = GetBlockchainAccountDataKey(destinationBlockchainAccount, "orderPaymentCallbackUrl")

	}

	{
		owner.InvalidateUserCache(gc)
	}

	return xdrBase64, &destinationInfo, nil

}

// GetBlockchainAccountDataKey looked up values from a Stellar account's
// on-chain manage_data store (e.g. a payment-callback URL). Base/EVM
// accounts have no equivalent store - see assets.go's GetDataKey doc for
// the same simplification applied elsewhere.
func GetBlockchainAccountDataKey(account *network.AccountInfo, keys ...string) (dataValues map[string]string) {
	return make(map[string]string)
}

func MintAsset(signerUser *userModels.User, sourceWallet *userModels.UserWallet, mintingInfo *userModels.MintingInfo, gc *sharedconfig.GlobalConfig) (*userModels.MintingInfo, *userModels.User, error) {
	db := gc.DB
	client := network.GetBlockchainClient()
	var xdrBase64 string
	var destinationUser *userModels.User
	var err error
	mintingInfo.Messages = make([]string, 0)
	walletHasViewOnlyAccess := true

	walletHasViewOnlyAccess = sourceWallet.HasViewOnlyAccess(gc)
	if !walletHasViewOnlyAccess {
		mintingInfo.Multiparty = 1

	}
	if walletHasViewOnlyAccess {
		mintingInfo.SignatureRequired = 1

	}
	if signerUser.IsEnterpriseProfile(gc) && mintingInfo.Multiparty == 0 {
		mintingInfo.SignatureRequired = 1

	}

	if len(mintingInfo.Transaction) > 0 {
		dUser, e := usersDB.GetUser(mintingInfo.Destination, db, gc)
		if e == nil {
			if len(dUser.ID) > 0 {
				destinationUser = &dUser

			}
		}
	}

	xdrBase64, destinationUser, err = generateMintingXdr(client, signerUser, sourceWallet, mintingInfo, db, gc)
	if err != nil {
		log.Printf("[MintAsset] from [%v] to [%v] generateMintingXdr error:[%v] \n", sourceWallet.ID, mintingInfo.Destination, err)
	}

	mintingInfo.Transaction = xdrBase64

	mintingInfo.NetworkPassPhrase = network.GetBlockchainNetworkPassPhrase()

	if len(mintingInfo.TransactionSignature) == 0 && mintingInfo.Commit == 0 {
		return mintingInfo, nil, err
	}

	// if paymentInfo.SHash != algofuncs.SHash(paymentInfo.Transaction) && paymentInfo.Commit == 0 {
	// 	return paymentInfo, nil, &tPayErrors.ErrorTransactionMismatch{}
	// }
	if len(mintingInfo.ChannelAccountSignature) == 0 && len(mintingInfo.ChannelAccount) == 42 {
		return mintingInfo, nil, &tPayErrors.ErrorTransactionMismatch{Detail: "signature for channel account does not validate"}
	}

	if mintingInfo.Multiparty == 0 {
		//shared access disabled. submit to network is possible
		var txnHash string
		if len(mintingInfo.ChannelAccountSignature) > 0 && len(mintingInfo.ChannelAccount) == 42 {
			// txnHash, err = network.SubmitXdrWithSignatureChannelAccounts(client, sourceWallet.Signer, paymentInfo.ChannelAccount, xdrBase64, paymentInfo.TransactionSignature, paymentInfo.ChannelAccountSignature)
			txnHash, err = network.SubmitXdrWithSignatureChannelAccounts(client, sourceWallet.Signer, mintingInfo.ChannelAccount, mintingInfo.Transaction, mintingInfo.TransactionSignature, mintingInfo.ChannelAccountSignature)
			if err != nil {
				log.Println("MintAsset######################submit with channel account throws error:", err)

			}
		} else {
			// txnHash, err = network.SubmitXdrWithSignature(client, sourceWallet.Signer, xdrBase64, paymentInfo.TransactionSignature)
			txnHash, err = network.SubmitXdrWithSignature(client, sourceWallet.Signer, mintingInfo.Transaction, mintingInfo.TransactionSignature)
			if err != nil {
				log.Printf("[MintAsset] from [%v] to [%v] SubmitXdrWithSignature error:[%v] \n", sourceWallet.Alias, mintingInfo.Destination, err)
			}
		}
		mintingInfo.TransactionID = txnHash
		return mintingInfo, destinationUser, err
	}

	if mintingInfo.Commit == 0 {
		return mintingInfo, nil, nil
	}

	mintingInfo.TransactionID = "PENDING_AUTH"
	log.Printf("[MintAsset]shared access with approver permission enabled for %v \n", sourceWallet.Alias)
	id := uuid.NewString()
	assetOfPayment := mintingInfo.AssetCode
	if len(mintingInfo.ContractAddress) == 42 {
		assetOfPayment = fmt.Sprintf("%v:%v...%v", mintingInfo.AssetCode, mintingInfo.ContractAddress[0:3], mintingInfo.ContractAddress[52:55])
	}
	var msgs string
	for i, m := range mintingInfo.Messages {
		msgs = m
		if i < len(mintingInfo.Messages)-1 {
			msgs = fmt.Sprintf("%s\n", msgs)
		}
	}
	description := fmt.Sprintf("Mint %v, \nTo: %v, \nAmount: %v %v", mintingInfo.AssetCode, mintingInfo.Destination, mintingInfo.Amount, assetOfPayment)
	if len(mintingInfo.Memo) > 0 {
		description = fmt.Sprintf("%v \nFor: %v", description, mintingInfo.Memo)

	}
	if len(msgs) > 0 {
		description = fmt.Sprintf("%v \nMessages: %v", description, msgs)
	}
	mintingInfo.ReturnedDescription = description
	transactionByte, _ := json.Marshal(*mintingInfo)
	transactionStr := string(transactionByte)
	pendingAuth := userModels.PendingAuth{
		ID:                     id,
		Initiator:              signerUser.Username,
		InitiatorSignerAddress: signerUser.PrimarySigner,
		WalletAddress:          sourceWallet.ID,
		TransactionType:        "MINT TOKEN",
		Description:            description,
		TransactionSource:      mintingInfo.TransactionSource,
		ApprovalsNeeded:        sourceWallet.NumberOfApprovalsNeeded,
		TransactionXdr:         mintingInfo.Transaction,
		TransactionInfoStr:     &transactionStr,
	}
	//save and commit this to database
	e := db.Omit(clause.Associations).Create(&pendingAuth).Error
	if e != nil {
		log.Printf("[Pay] Error saving payment txn [%+v] transaction on pending auth table: %s\n", pendingAuth, e.Error())
		err = &tErrors.ErrorTemporaryServerError{}
		return mintingInfo, destinationUser, err
	}

	return mintingInfo, destinationUser, nil

}

// func logDiscordFailedPayment(msg string) {
// 	discord.WebhookURL = "https://discord.com/api/webhooks/827986576415129663/wqMKp9wxB_fxs9Q3zlMKCNPGENXmD_ueUnL8hVCu1wmRfD2wkXAjfP85k1Ro_2_wGfiY"
// 	if len(os.Getenv("FAILED_PAYMENT_ERROR_WEBHOOK")) > 50 {
// 		discord.WebhookURL = os.Getenv("FAILED_PAYMENT_ERROR_WEBHOOK")
// 	}
// 	discord.Say(msg)
// }

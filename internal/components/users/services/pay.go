package users

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	algofuncs "trovo-wallet-api/internal/blockchainalgofuncs"
	tPayErrors "trovo-wallet-api/internal/components/payments/errors"
	paymentModels "trovo-wallet-api/internal/components/payments/models"
	userBc "trovo-wallet-api/internal/components/users/blockchain"
	usersDB "trovo-wallet-api/internal/components/users/db"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
)

// Pay function sends a payment from user to another user
func Pay(signerUser *userModels.User, wallet *userModels.UserWallet, paymentInfo *paymentModels.PaymentInfo, gc *sharedconfig.GlobalConfig) (*paymentModels.PaymentInfo, *userModels.User, error) {
	db := gc.DB
	client := network.GetBlockchainClient()
	var xdrBase64 string
	var destinationUser *userModels.User
	var err error

	//check if destination is a wallet with memo
	if len(paymentInfo.Destination) == 56 {
		memoWalletSlices28byte := strings.Split(os.Getenv("WALLETS_REQUIRE_28_BYTE_MEMO"), ",")

		for _, w := range memoWalletSlices28byte {
			if paymentInfo.Destination == w && len(strings.ReplaceAll(paymentInfo.Memo, " ", "")) != 28 {
				return paymentInfo, destinationUser, &tErrors.CustomError{
					Param:      "memo",
					Err:        "error invalid memo",
					ErrMessage: "Bittrex Global Deposits require 28 character memo. If you do not put exact memo, your funds will be lost. Please carefully provide the memo for your Bittrex Global wallet.",
				}

			}
		}

		memoWalletSlices16byte := strings.Split(os.Getenv("WALLETS_REQUIRE_16_BYTE_MEMO"), ",")

		for _, w := range memoWalletSlices16byte {
			if paymentInfo.Destination == w && len(strings.ReplaceAll(paymentInfo.Memo, " ", "")) != 16 {
				return paymentInfo, destinationUser, &tErrors.CustomError{
					Param:      "memo",
					Err:        "error invalid memo",
					ErrMessage: "FMFW/Tradefada Deposits require 16 character memo. If you do not put exact memo, your funds will be lost. Please carefully provide the memo for your exchange's XBN wallet.",
				}

			}

		}

		memoWalletSlicesVariablebyte := strings.Split(os.Getenv("WALLETS_REQUIRE_VARIABLE_BYTE_MEMO"), ",")

		for _, w := range memoWalletSlicesVariablebyte {
			if paymentInfo.Destination == w && len(strings.ReplaceAll(paymentInfo.Memo, " ", "")) < 9 {
				return paymentInfo, destinationUser, &tErrors.CustomError{
					Param:      "memo",
					Err:        "error invalid memo",
					ErrMessage: "You are attempting to send to an exchange that requires memo for all deposits. If you do not put exact memo in the description field, your funds will be lost. Please carefully provide the memo for your exchange's XBN wallet.",
				}

			}
		}

	}

	if len(paymentInfo.ChannelAccount) == 56 {
		//payment is with channel account
		xdrBase64, destinationUser, err = generatePaymentXdrWithChannelAccountPK(client, signerUser, wallet, paymentInfo, db)

		if err != nil {
			log.Printf("[Pay] from [%v] to [%v] generatePaymentXdrWithChannelAccountPK error:[%v]\n", wallet.Alias, paymentInfo.Destination, err)
		}
	} else {

		xdrBase64, destinationUser, err = generatePaymentXdr(client, signerUser, wallet, paymentInfo, db)
		if err != nil {
			log.Printf("[Pay] from [%v] to [%v] generatePaymentXdr error:[%v] \n", wallet.ID, paymentInfo.Destination, err)
		}
	}
	oldTransaction := paymentInfo.Transaction

	paymentInfo.Transaction = xdrBase64
	paymentInfo.NetworkPassPhrase = network.GetBlockchainNetworkPassPhrase()
	if !wallet.HasViewOnlyAccess(gc) && wallet.SharedAccessEnabled == 1 {
		paymentInfo.Multiparty = 1
	}
	if len(paymentInfo.TransactionSignature) == 0 && paymentInfo.Commit == 0 {
		return paymentInfo, nil, err
	}

	if xdrBase64 != oldTransaction {
		log.Printf("[PAY]oldTransaction: %v\nNewTransaction: %v\n", oldTransaction, xdrBase64)
		return paymentInfo, nil, &tPayErrors.ErrorTransactionMismatch{}
	}
	if len(paymentInfo.ChannelAccountSignature) == 0 && len(paymentInfo.ChannelAccount) == 56 {
		return paymentInfo, nil, &tPayErrors.ErrorTransactionMismatch{Detail: "signature for channel account does not validate"}
	}

	if wallet.SharedAccessEnabled == 0 {
		//shared access disabled. submit to network is possible
		var txnHash string
		if len(paymentInfo.ChannelAccountSignature) > 0 && len(paymentInfo.ChannelAccount) == 56 {
			txnHash, err = network.SubmitXdrWithSignatureChannelAccounts(client, wallet.Signer, paymentInfo.ChannelAccount, xdrBase64, paymentInfo.TransactionSignature, paymentInfo.ChannelAccountSignature)
			if err != nil {
				log.Println("#############################submit with channel account throws error:", err)

			}
		} else {
			txnHash, err = network.SubmitXdrWithSignature(client, wallet.Signer, xdrBase64, paymentInfo.TransactionSignature)
			if err != nil {
				log.Printf("[Pay] from [%v] to [%v] SubmitXdrWithSignature error:[%v] \n", wallet.Alias, paymentInfo.Destination, err)
			}
		}
		paymentInfo.TransactionID = txnHash
		return paymentInfo, destinationUser, err
	}
	walletHasViewOnlyAccess := wallet.HasViewOnlyAccess(gc)
	//put routine for shared access > 0
	//check to be sure wallet permission includes approval
	if walletHasViewOnlyAccess {
		//maker checker not enabled. submit to network is possible
		var txnHash string
		if len(paymentInfo.ChannelAccountSignature) > 0 && len(paymentInfo.ChannelAccount) == 56 {
			txnHash, err = network.SubmitXdrWithSignatureChannelAccounts(client, wallet.Signer, paymentInfo.ChannelAccount, xdrBase64, paymentInfo.TransactionSignature, paymentInfo.ChannelAccountSignature)
			if err != nil {
				log.Println("#############################submit with channel account throws error:", err)

			}
		} else {
			txnHash, err = network.SubmitXdrWithSignature(client, wallet.Signer, xdrBase64, paymentInfo.TransactionSignature)
			if err != nil {
				log.Printf("[Pay] from [%v] to [%v] SubmitXdrWithSignature error:[%v] \n", wallet.Alias, paymentInfo.Destination, err)
			}
		}
		paymentInfo.TransactionID = txnHash
		return paymentInfo, destinationUser, err
	}
	if paymentInfo.Commit == 0 {
		return paymentInfo, nil, nil
	}
	paymentInfo.TransactionID = "PENDING_AUTH"
	log.Printf("[Pay]shared access with approver permission enabled for %v \n", wallet.Alias)
	id := uuid.NewString()
	assetOfPayment := os.Getenv("NATIVE_ASSET_CODE")
	if len(paymentInfo.AssetIssuer) == 56 {
		assetOfPayment = fmt.Sprintf("%v:%v...%v", paymentInfo.AssetCode, paymentInfo.AssetIssuer[0:4], paymentInfo.AssetIssuer[51:55])
	}
	description := fmt.Sprintf("Payment from:%v To: %v. Amount: %v %v.\nMemo: %v\nMessages: %v\n", wallet.Alias, paymentInfo.Destination, paymentInfo.Amount, assetOfPayment, paymentInfo.Memo, paymentInfo.Messages)
	transactionByte, _ := json.Marshal(*paymentInfo)
	transactionStr := string(transactionByte)
	pendingAuth := userModels.PendingAuth{
		ID:                       id,
		Initiator:                signerUser.Username,
		InitiatorSignerPublicKey: signerUser.PrimarySigner,
		WalletPublicKey:          wallet.ID,
		TransactionType:          "PAYMENT",
		Description:              description,
		ApprovalsNeeded:          wallet.NumberOfApprovalsNeeded,
		TransactionXdr:           xdrBase64,
		TransactionInfoStr:       &transactionStr,
	}
	//save and commit this to database
	e := db.Create(&pendingAuth).Error
	if e != nil {
		log.Printf("[Pay] Error saving payment txn [%+v] transaction on pending auth table: %s\n", pendingAuth, e.Error())
		err = &tErrors.ErrorTemporaryServerError{}
		return paymentInfo, destinationUser, err
	}

	return paymentInfo, destinationUser, nil

}

func generatePaymentXdr(client *horizonclient.Client, owner *userModels.User, wallet *userModels.UserWallet, paymentInfo *paymentModels.PaymentInfo, db *gorm.DB) (string, *userModels.User, error) {
	baseReserve := network.GetBlockchainBaseReserve()
	// var messages []string
	//check if it is public key payment
	publicKeyPayment := len(paymentInfo.Destination) == 56
	var err error
	paymentInfo, err = ValidatePaymentInfo(paymentInfo)

	if err != nil {
		return "", nil, err
	}

	var amountToSend float64
	if amountToSend, err = strconv.ParseFloat(paymentInfo.Amount, 64); err != nil {
		return "", nil, &tPayErrors.ErrorInvalidPaymentAmount{}
	}

	newAmountToSend := decimal.NewFromFloat(amountToSend).Truncate(7).String()

	var asset txnbuild.Asset = txnbuild.NativeAsset{}

	if len(paymentInfo.AssetCode) != 0 {
		asset = txnbuild.CreditAsset{Code: paymentInfo.AssetCode, Issuer: paymentInfo.AssetIssuer}
	}

	destinationInfo, getDestinationError := usersDB.GetUser(paymentInfo.Destination, db)
	destinationWallet, _, _ := usersDB.GetWallet(paymentInfo.Destination, db)

	charge := baseReserve.Mul(decimal.NewFromInt(3)).Truncate(7).String()
	if getDestinationError != nil && len(paymentInfo.Destination) != 56 {
		return "", nil, &tPayErrors.ErrorPaymentDestinationDoesNotExist{}
	}
	if destinationInfo.Suspended == 1 {
		return "", nil, &tErrors.ErrorUsernameIsSuspended{}
	}

	if !publicKeyPayment {
		paymentInfo.DestinationFirstName = destinationInfo.FirstName
		if destinationInfo.LastName != nil {
			paymentInfo.DestinationLastName = *destinationInfo.LastName
		}

		paymentInfo.DestinationVerified = destinationInfo.Verified

		if destinationInfo.ImageThumbnailURL != nil {
			paymentInfo.DestinationThumbnail = *destinationInfo.ImageThumbnailURL
		}
	} else {
		//parse public key
		_, err := keypair.ParseAddress(paymentInfo.Destination)
		if err != nil {
			return "", nil, &tPayErrors.ErrorInvalidPaymentDestinationPublicKey{}
		}

		message := "You are about to make payment to a public key directly. Please be sure of the address as the payment cannot be retrieved after confirmation."

		paymentInfo.Messages = append(paymentInfo.Messages, message)
		// log.Printf("[generatePaymentXdr]message for public key logged: %v\n", message)
	}
	var destinationPublicKey string
	if publicKeyPayment {
		destinationPublicKey = paymentInfo.Destination
	} else {
		destinationPublicKey = destinationWallet.ID
	}
	//perform ths checks of determining messages to be appended. if destination account property is not checked here, information would be returned without messages set.
	destinationAccountExists, destinationAccountTrustsAsset, _, _, destinationBlockchainAccount, destinationAccountErr :=
		network.BlockchainAccountProperties(client, destinationPublicKey, asset)
	//set base charge to be used in all places it is needed
	if publicKeyPayment {
		if !asset.IsNative() && !destinationAccountTrustsAsset {
			return "", nil, &tPayErrors.ErrorDestinationPublicKeyCannotReceiveAsset{}
		}
	} else {
		//check if to set charges messages
		if !asset.IsNative() {
			//custom asset
			if !destinationAccountExists {

				message := fmt.Sprintf("Important: The wallet %v is unfunded. %v XBN will be deducted from your account to fund %v’s account. You only need to do this once for %v.", destinationWallet.Alias, charge, destinationWallet.Alias, destinationWallet.Alias)

				paymentInfo.Messages = append(paymentInfo.Messages, message)
				// log.Printf("[generatePaymentXdr]message[0]: %v\n", message)

			}

			if !destinationAccountTrustsAsset {
				message := fmt.Sprintf("Important: %v has not yet activated the asset (%v) you are trying to send. %v XBN will be deducted from your account to ensure that this transaction goes through. After this, %v will be able to receive %v anytime, without any further charges to you.", destinationWallet.Alias, paymentInfo.AssetCode, charge, destinationWallet.Alias, paymentInfo.AssetCode)

				paymentInfo.Messages = append(paymentInfo.Messages, message)
				// log.Printf("[generatePaymentXdr]message[1]: %v\n", message)

			}
		}
	}

	// paymentInfo.Messages = messages
	sourceAccountExists, sourceAccountTrustsAsset, sourceAccountNativeBalance, sourceAccountCustomBalance, sourceAccount, sourceAccountErr := network.BlockchainAccountProperties(client, wallet.ID, asset)

	if sourceAccountErr != nil {
		return "", nil, sourceAccountErr
	}

	if !sourceAccountExists {
		return "", nil, &tErrors.ErrorUnderfundedAccount{}
	}

	if !sourceAccountTrustsAsset {

		return "", nil, &tErrors.ErrorUnderfundedAccount{}
	}

	log.Printf("[generatePaymentXdr]obtained sourced account info \n")

	log.Printf("[generatePaymentXdr]obtained source account balance is %v, custom account balance %v\n", sourceAccountNativeBalance, sourceAccountCustomBalance)

	amountToSendDec := decimal.NewFromFloat(amountToSend)

	if wallet.ID != asset.GetIssuer() {

		if asset.IsNative() {
			if sourceAccountNativeBalance.LessThan(amountToSendDec) {
				return "", nil, &tErrors.ErrorUnderfundedAccount{}
			}
		} else {
			if sourceAccountCustomBalance.LessThan(amountToSendDec) {
				return "", nil, &tErrors.ErrorUnderfundedAccount{}
			}
		}
	}

	//check if destination account exists
	if destinationAccountErr != nil {
		log.Println("[generatePaymentXdr]destination Account error:", destinationAccountErr)
		return "", nil, destinationAccountErr
	}
	var ops []txnbuild.Operation = make([]txnbuild.Operation, 0)

	var tempAccountKeyPair *keypair.Full = nil

	if asset.IsNative() {
		//native asset
		if !destinationAccountExists {

			if amountToSendDec.LessThan(baseReserve.Mul(decimal.NewFromInt(3))) {
				log.Printf("[generatePaymentXdr] trying to create account with %v, meanwhile you need %v\n", amountToSendDec, baseReserve.Mul(decimal.NewFromInt(3)))
				return "", nil, &tPayErrors.ErrorInsufficientAmountToFundAccount{}
			}
			ops = append(ops, &txnbuild.CreateAccount{
				Destination:   destinationPublicKey,
				Amount:        newAmountToSend,
				SourceAccount: wallet.ID,
			})
		} else {
			ops = append(ops, &txnbuild.Payment{
				Destination:   destinationPublicKey,
				Amount:        newAmountToSend,
				Asset:         asset,
				SourceAccount: wallet.ID,
			})
		}
	} else {
		//custom asset

		// claimable assets are for trovo wallet users only. it would return error above when destination does not trust asset

		if !destinationAccountExists {
			ops = append(ops, &txnbuild.CreateAccount{
				Destination:   destinationPublicKey,
				Amount:        charge,
				SourceAccount: wallet.ID,
			})
		}

		if !destinationAccountTrustsAsset {
			if !publicKeyPayment {
				ops2, _tempAccountKeyPair, err := processDestinationAssetDoesNotTrustAsset(destinationInfo, client, &destinationWallet, sourceAccount, asset, newAmountToSend, db)

				if err != nil {
					return "", nil, err
				}

				tempAccountKeyPair = _tempAccountKeyPair

				ops = append(ops, ops2...)
			}

		} else {
			ops = append(ops, &txnbuild.Payment{
				Destination:   destinationPublicKey,
				Amount:        newAmountToSend,
				Asset:         asset,
				SourceAccount: wallet.ID,
			})
		}

	}

	// Construct the transaction that holds the operations to execute on the network
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        sourceAccount,
			IncrementSequenceNum: true,
			Operations:           ops,
			BaseFee:              txnbuild.MinBaseFee,
			Preconditions: txnbuild.Preconditions{
				TimeBounds: txnbuild.NewInfiniteTimeout(),
			},
			Memo: txnbuild.MemoText(paymentInfo.Memo),
		},
	)
	if err != nil {
		log.Println("[generatePaymentXdr] error constructing transaction ", err)
		return "", nil, err
	}

	if tempAccountKeyPair != nil {

		tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), tempAccountKeyPair)

		if err != nil {
			log.Println("[generatePaymentXdr] error signing transaction with temporary key ", err)
			return "", nil, &tErrors.ErrorTemporaryServerError{}
		}
	}
	var xdrBase64 string

	xdrBase64, err = tx.Base64()
	if err != nil {
		log.Println("[generatePaymentXdr] error getting txn base64", err)
		return "", nil, err
	}

	if destinationBlockchainAccount != nil {
		paymentInfo.CallbackURLS = GetBlockchainAccountDataKey(*destinationBlockchainAccount, "orderPaymentCallbackUrl")

	}
	if publicKeyPayment {
		return xdrBase64, nil, nil
	}
	return xdrBase64, &destinationInfo, nil

}

func generatePaymentXdrWithChannelAccountPK(client *horizonclient.Client, owner *userModels.User, wallet *userModels.UserWallet, paymentInfo *paymentModels.PaymentInfo, db *gorm.DB) (string, *userModels.User, error) {
	baseReserve := network.GetBlockchainBaseReserve()
	// var messages []string
	//check if it is public key payment
	publicKeyPayment := len(paymentInfo.Destination) == 56
	var err error
	paymentInfo, err = ValidatePaymentInfo(paymentInfo)

	if err != nil {
		return "", nil, err
	}

	var amountToSend float64
	if amountToSend, err = strconv.ParseFloat(paymentInfo.Amount, 64); err != nil {
		return "", nil, &tPayErrors.ErrorInvalidPaymentAmount{}
	}

	newAmountToSend := decimal.NewFromFloat(amountToSend).Truncate(7).String()

	var asset txnbuild.Asset = txnbuild.NativeAsset{}

	if len(paymentInfo.AssetCode) != 0 {
		asset = txnbuild.CreditAsset{Code: paymentInfo.AssetCode, Issuer: paymentInfo.AssetIssuer}
	}
	destinationInfo, getDestinationError := usersDB.GetUser(paymentInfo.Destination, db)
	destinationWallet, _, _ := usersDB.GetWallet(paymentInfo.Destination, db)
	charge := baseReserve.Mul(decimal.NewFromInt(3)).Truncate(7).String()
	if getDestinationError != nil && len(paymentInfo.Destination) != 56 {
		return "", nil, &tPayErrors.ErrorPaymentDestinationDoesNotExist{}
	}

	var destinationPublicKey string
	if publicKeyPayment {
		destinationPublicKey = paymentInfo.Destination
	} else {
		destinationPublicKey = destinationWallet.ID
	}
	if !publicKeyPayment {
		paymentInfo.DestinationFirstName = destinationInfo.FirstName
		if destinationInfo.LastName != nil {
			paymentInfo.DestinationLastName = *destinationInfo.LastName
		}

		paymentInfo.DestinationVerified = destinationInfo.Verified
		if destinationInfo.ImageThumbnailURL != nil {
			paymentInfo.DestinationThumbnail = *destinationInfo.ImageThumbnailURL
		}
	} else {
		//parse public key
		_, err := keypair.ParseAddress(paymentInfo.Destination)
		if err != nil {
			return "", nil, &tPayErrors.ErrorInvalidPaymentDestinationPublicKey{}
		}

		message := "You are about to make payment to a public key directly. Please be sure of the address as the payment cannot be retrieved after confirmation."

		paymentInfo.Messages = append(paymentInfo.Messages, message)
		// log.Printf("[generatePaymentXdr]message for public key logged: %v\n", message)
	}
	// var destinationPublicKey string
	// if publicKeyPayment {
	// 	destinationPublicKey = paymentInfo.Destination
	// } else {
	// 	destinationPublicKey = destinationUser.PublicKey
	// }
	//perform ths checks to determine messages to be appended. if destination account property is not checked here, information would be returned without messages set.
	destinationAccountExists, destinationAccountTrustsAsset, _, _, destinationBlockchainAccount, destinationAccountErr :=
		network.BlockchainAccountProperties(client, destinationPublicKey, asset)
	//set base charge to be used in all places it is needed

	if publicKeyPayment {
		if !asset.IsNative() && !destinationAccountTrustsAsset {
			return "", nil, &tPayErrors.ErrorDestinationPublicKeyCannotReceiveAsset{}
		}
	} else {
		//check if to set charges messages
		if !asset.IsNative() {
			//custom asset
			if !destinationAccountExists {

				message := fmt.Sprintf("Important: The wallet %v is unfunded. %v XBN will be deducted from your account to fund %v’s account. You only need to do this once for %v.", destinationWallet.Alias, charge, destinationWallet.Alias, destinationWallet.Alias)

				paymentInfo.Messages = append(paymentInfo.Messages, message)
				// log.Printf("[generatePaymentXdr]message[0]: %v\n", message)

			}

			if !destinationAccountTrustsAsset {
				message := fmt.Sprintf("Important: %v has not yet activated the asset (%v), that you are trying to send. %v XBN will be deducted from your account to ensure that this transaction goes through. After this, %v will be able to receive %v anytime, without any further charges to you.", destinationWallet.Alias, paymentInfo.AssetCode, charge, destinationWallet.Alias, paymentInfo.AssetCode)

				paymentInfo.Messages = append(paymentInfo.Messages, message)
				// log.Printf("[generatePaymentXdr]message[1]: %v\n", message)

			}
		}
	}

	// paymentInfo.Messages = messages
	sourceAccountExists, sourceAccountTrustsAsset, sourceAccountNativeBalance, sourceAccountCustomBalance, sourceAccount, sourceAccountErr := network.BlockchainAccountProperties(client, wallet.ID, asset)
	channelSourceAccountExists, _, channelSourceAccountNativeBalance, _, channelSourceAccount, channelSourceAccountErr := network.BlockchainAccountProperties(client, paymentInfo.ChannelAccount, txnbuild.NativeAsset{})

	if channelSourceAccountErr != nil {
		return "", nil, channelSourceAccountErr
	}

	if !channelSourceAccountExists {
		return "", nil, &tErrors.ErrorUnderfundedAccount{}
	}

	if channelSourceAccountNativeBalance.LessThan(decimal.NewFromFloat(6.1)) {
		//min balance of 6XBN and fee of
		return "", nil, &tErrors.ErrorUnderfundedAccount{Detail: fmt.Sprintf("Channel account is underfunded[%v XBN usable]. Needs extra %v XBN", channelSourceAccountNativeBalance.String(), decimal.NewFromFloat(6.1).Sub(channelSourceAccountNativeBalance))}
	}

	if sourceAccountErr != nil {
		return "", nil, sourceAccountErr
	}

	if !sourceAccountExists {
		return "", nil, &tErrors.ErrorUnderfundedAccount{}
	}

	if !sourceAccountTrustsAsset {

		return "", nil, &tErrors.ErrorUnderfundedAccount{}
	}

	log.Printf("[generatePaymentXdrWithChannelAccountPK]obtained sourced account info \n")

	log.Printf("[generatePaymentXdrWithChannelAccountPK]obtained source account balance is %v, custom account balance %v\n", sourceAccountNativeBalance, sourceAccountCustomBalance)

	amountToSendDec := decimal.NewFromFloat(amountToSend)

	if wallet.ID != asset.GetIssuer() {

		if asset.IsNative() {
			if sourceAccountNativeBalance.LessThan(amountToSendDec) {
				return "", nil, &tErrors.ErrorUnderfundedAccount{}
			}
		} else {
			if sourceAccountCustomBalance.LessThan(amountToSendDec) {
				return "", nil, &tErrors.ErrorUnderfundedAccount{}
			}
		}
	}

	//check if destination account exists
	if destinationAccountErr != nil {
		log.Println("[generatePaymentXdrWithChannelAccountPK]destination Account error:", destinationAccountErr)
		return "", nil, destinationAccountErr
	}
	var ops []txnbuild.Operation = make([]txnbuild.Operation, 0)

	var tempAccountKeyPair *keypair.Full = nil

	if asset.IsNative() {
		//native asset
		if !destinationAccountExists {

			if amountToSendDec.LessThan(baseReserve.Mul(decimal.NewFromInt(3))) {
				log.Printf("[generatePaymentXdrWithChannelAccountPK] trying to create account with %v, meanwhile you need %v\n", amountToSendDec, baseReserve.Mul(decimal.NewFromInt(3)))
				return "", nil, &tPayErrors.ErrorInsufficientAmountToFundAccount{}
			}
			ops = append(ops, &txnbuild.CreateAccount{
				Destination:   destinationPublicKey,
				Amount:        newAmountToSend,
				SourceAccount: wallet.ID,
			})
		} else {
			ops = append(ops, &txnbuild.Payment{
				Destination:   destinationPublicKey,
				Amount:        newAmountToSend,
				Asset:         asset,
				SourceAccount: wallet.ID,
			})
		}
	} else {
		//custom asset

		// claimable assets are for trovotech customers only. it would return error above when destination does not trust asset

		if !destinationAccountExists {
			ops = append(ops, &txnbuild.CreateAccount{
				Destination:   destinationPublicKey,
				Amount:        charge,
				SourceAccount: wallet.ID,
			})
		}

		if !destinationAccountTrustsAsset {
			if !publicKeyPayment {

				ops2, _tempAccountKeyPair, err := processDestinationAssetDoesNotTrustAsset(destinationInfo, client, &destinationWallet, sourceAccount, asset, newAmountToSend, db)

				if err != nil {
					return "", nil, err
				}

				tempAccountKeyPair = _tempAccountKeyPair

				ops = append(ops, ops2...)
			}

		} else {
			ops = append(ops, &txnbuild.Payment{
				Destination:   destinationPublicKey,
				Amount:        newAmountToSend,
				Asset:         asset,
				SourceAccount: wallet.ID,
			})
		}

	}

	// Construct the transaction that holds the operations to execute on the network
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        channelSourceAccount,
			IncrementSequenceNum: true,
			Operations:           ops,
			BaseFee:              txnbuild.MinBaseFee,
			Preconditions: txnbuild.Preconditions{
				TimeBounds: txnbuild.NewInfiniteTimeout(),
			},
			Memo: txnbuild.MemoText(paymentInfo.Memo),
		},
	)
	if err != nil {
		log.Println("[generatePaymentXdrWithChannelAccountPK] error constructing transaction ", err)
		return "", nil, err
	}

	if tempAccountKeyPair != nil {

		tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), tempAccountKeyPair)

		if err != nil {
			log.Println("[generatePaymentXdrWithChannelAccountPK] error signing transaction with temporary key ", err)
			return "", nil, &tErrors.ErrorTemporaryServerError{}
		}
	}

	var xdrBase64 string

	xdrBase64, err = tx.Base64()
	if err != nil {
		log.Println("[generatePaymentXdrWithChannelAccountPK] error getting txn base64", err)
		return "", nil, err
	}

	if destinationBlockchainAccount != nil {
		paymentInfo.CallbackURLS = GetBlockchainAccountDataKey(*destinationBlockchainAccount, "orderPaymentCallbackUrl")

	}
	if publicKeyPayment {
		return xdrBase64, nil, nil
	}
	return xdrBase64, &destinationInfo, nil
}

func processDestinationAssetDoesNotTrustAsset(destinationUser userModels.User, client *horizonclient.Client, destinationWallet *userModels.UserWallet, sourceAccount *horizon.Account, asset txnbuild.Asset, amountToSend string, db *gorm.DB) ([]txnbuild.Operation, *keypair.Full, error) {

	ops := make([]txnbuild.Operation, 0)

	var signerKeyPairToReturn *keypair.Full = nil

	tempAccountKeypair, tempAccountError := network.TempAccountKeypair(destinationWallet.ID)

	if tempAccountError != nil {
		log.Printf("[processDestinationAssetDoesNotTrustAsset] error generating temporary account %v\n", tempAccountError)
		return ops, tempAccountKeypair, &tErrors.ErrorTemporaryServerError{}
	}

	// var tempAccount txnbuild.Account = &txnbuild.SimpleAccount{AccountID: tempAccountKeypair.Address(), Sequence: 0}

	tempAccountExists, tempAccountTrustsAsset, _, _, tempAccountSource, tempAccountError :=
		network.BlockchainAccountProperties(client, tempAccountKeypair.FromAddress().Address(), asset)

	if tempAccountError != nil {
		log.Printf("[processDestinationAssetDoesNotTrustAsset] error looking up temporary account %v\n", tempAccountError)
		return ops, tempAccountKeypair, &tErrors.ErrorTemporaryServerError{}
	}

	{

		var _wallet userModels.UserWallet

		//check if temp account exists in buds or not.
		err := db.Where("id = ?", destinationWallet.ID).First(&_wallet).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ops, nil, &tErrors.ErrorInvalidPublicKey{}

			}
			log.Println("[processDestinationAssetDoesNotTrustAsset]", err)
			return ops, nil, &tErrors.ErrorTemporaryServerError{}

		}

		// var _walletToUpdate users.UserWallet
		var update bool

		if _wallet.TempPublicKey != nil {
			if *_wallet.TempPublicKey != tempAccountKeypair.Address() {
				v := tempAccountKeypair.Address()
				_wallet.TempPublicKey = &v
				update = true

			}
		} else {
			v := tempAccountKeypair.Address()
			_wallet.TempPublicKey = &v
			update = true
		}

		if update {
			dbSaveError := db.Save(&_wallet).Error

			if dbSaveError != nil {
				log.Printf("[processDestinationAssetDoesNotTrustAsset]db temp save error %v\n", dbSaveError)
				return ops, nil, &tErrors.ErrorTemporaryServerError{}
			}
		}

	}

	//no errors, so do

	baseReserve := network.GetBlockchainBaseReserve()

	if !tempAccountExists {

		//create temp account and add destination public key as signer.

		ops = append(ops, &txnbuild.CreateAccount{
			Destination:   tempAccountKeypair.Address(),
			Amount:        baseReserve.Mul(decimal.NewFromInt(3)).Truncate(7).String(),
			SourceAccount: sourceAccount.AccountID,
		})

		ops = append(ops, &txnbuild.SetOptions{
			Signer: &txnbuild.Signer{
				Address: destinationWallet.Signer,
				Weight:  1,
			},
			SourceAccount: tempAccountKeypair.Address(),
		})

		signerKeyPairToReturn = tempAccountKeypair
	}

	//add the recovery address if enabled and account exists but recovery is not already signer key
	if len(destinationUser.Username) > 0 {
		if destinationUser.AccountRecoveryEnabled == 1 {
			recoveryKeyAddress := algofuncs.GetRecoveryAccountAddress(destinationUser.Username, destinationUser.PublicKey)
			if len(recoveryKeyAddress) > 0 {
				if !userBc.SignerIsValid(tempAccountKeypair.Address(), recoveryKeyAddress) {

					//just make the recovery key a signer

					ops = append(ops, &txnbuild.SetOptions{
						Signer: &txnbuild.Signer{
							Address: recoveryKeyAddress,
							Weight:  1,
						},
						SourceAccount: tempAccountKeypair.Address(),
					})

					signerKeyPairToReturn = tempAccountKeypair
				}
			}

		}
	}

	// destinationWallet.SignerIsValidWA(destinationWallet.Signer, tempAccountSource)
	if tempAccountExists && !destinationWallet.SignerIsValidWA(destinationWallet.Signer, tempAccountSource) {

		//just make the temporary key a signer

		ops = append(ops, &txnbuild.SetOptions{
			Signer: &txnbuild.Signer{
				Address: destinationWallet.Signer,
				Weight:  1,
			},
			SourceAccount: tempAccountKeypair.Address(),
		})

		signerKeyPairToReturn = tempAccountKeypair
	}

	if !tempAccountTrustsAsset {

		ops = append(ops, &txnbuild.Payment{
			Destination:   tempAccountKeypair.FromAddress().Address(),
			Amount:        baseReserve.Mul(decimal.NewFromInt(2)).Truncate(7).String(),
			Asset:         txnbuild.NativeAsset{},
			SourceAccount: sourceAccount.AccountID,
		})

		ops = append(ops, &txnbuild.ChangeTrust{
			Line:          txnbuild.ChangeTrustAssetWrapper{Asset: asset},
			Limit:         "900000000000",
			SourceAccount: tempAccountKeypair.Address(),
		})

		signerKeyPairToReturn = tempAccountKeypair

	}

	ops = append(ops, &txnbuild.Payment{
		Destination:   tempAccountKeypair.Address(),
		Amount:        amountToSend,
		Asset:         asset,
		SourceAccount: sourceAccount.AccountID,
	})

	return ops, signerKeyPairToReturn, nil

}

// GetBlockchainAccountDataKey fetches the bantu account information using public key
func GetBlockchainAccountDataKey(account horizon.Account, keys ...string) (dataValues map[string]string) {
	dataValues = make(map[string]string)

	data, err := GetBlockchainAccountData(account)
	if err != nil {
		return
	}
	for _, key := range keys {
		d, ok := data[key]
		if !ok {
			continue
		}
		decData, err := base64.StdEncoding.DecodeString(d)

		if err != nil {
			continue
		} else {
			dataValues[key] = string(decData)
		}

	}

	return dataValues
}

func GetBlockchainAccountData(clientAccount horizon.Account) (accountData map[string]string, err error) {

	if err != nil {
		return
	}

	return clientAccount.Data, nil
}

// func logDiscordFailedPayment(msg string) {
// 	discord.WebhookURL = "https://discord.com/api/webhooks/827986576415129663/wqMKp9wxB_fxs9Q3zlMKCNPGENXmD_ueUnL8hVCu1wmRfD2wkXAjfP85k1Ro_2_wGfiY"
// 	if len(os.Getenv("FAILED_PAYMENT_ERROR_WEBHOOK")) > 50 {
// 		discord.WebhookURL = os.Getenv("FAILED_PAYMENT_ERROR_WEBHOOK")
// 	}
// 	discord.Say(msg)
// }

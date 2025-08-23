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
	"gorm.io/gorm/clause"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
)

// Pay function sends a payment from user to another user
func Pay(signerUser *userModels.User, sourceWallet *userModels.UserWallet, paymentInfo *paymentModels.PaymentInfo, gc *sharedconfig.GlobalConfig) (*paymentModels.PaymentInfo, *userModels.User, error) {
	db := gc.DB
	client := network.GetBlockchainClient()
	var xdrBase64 string
	var destinationUser *userModels.User
	var err error
	walletHasViewOnlyAccess := true
	publicKeyPayment := len(paymentInfo.Destination) == 56 || len(paymentInfo.Destination) == 69

	//check if destination is a wallet with memo
	if publicKeyPayment {
		paymentInfo.Destination = strings.ToUpper(paymentInfo.Destination)
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

	} else {
		paymentInfo.Destination = strings.ToLower(paymentInfo.Destination)
	}
	paymentInfo.AmountToPay = paymentInfo.Amount
	walletHasViewOnlyAccess = sourceWallet.HasViewOnlyAccess(gc)
	if !walletHasViewOnlyAccess {
		paymentInfo.Multiparty = 1
		{
			//calculate fees
			fee := decimal.RequireFromString(sourceWallet.GetSharedAccessPaymentFee(gc))
			paymentInfo.Fee = fee.String()
			feeAmount := ((decimal.RequireFromString(paymentInfo.Amount).Mul(fee)).Div(decimal.NewFromInt(100))).Truncate(7)
			paymentInfo.FeeAmount = feeAmount.String()
			amountToPay := decimal.RequireFromString(paymentInfo.Amount).Add(feeAmount)
			paymentInfo.AmountToPay = amountToPay.String()

		}

	}
	if walletHasViewOnlyAccess {
		paymentInfo.SignatureRequired = 1

	}

	if (!publicKeyPayment && len(paymentInfo.Transaction) > 0) && len(paymentInfo.SHash) > 1 {
		dUser, e := usersDB.GetUser(paymentInfo.Destination, db, gc)
		if e == nil {
			if len(dUser.ID) > 0 {
				destinationUser = &dUser
			}

		}
	}
	// if len(paymentInfo.SHash) == 0 {
	if len(paymentInfo.ChannelAccount) == 56 {
		//payment is with channel account
		xdrBase64, destinationUser, err = generatePaymentXdrWithChannelAccountPK(signerUser, sourceWallet, paymentInfo, gc)

		if err != nil {
			log.Printf("[Pay] from [%v] to [%v] generatePaymentXdrWithChannelAccountPK error:[%v]\n", sourceWallet.Alias, paymentInfo.Destination, err)
		}
	} else {
		xdrBase64, destinationUser, err = generatePaymentXdr(client, signerUser, sourceWallet, paymentInfo, db, gc)
		if err != nil {
			log.Printf("[Pay] from [%v] to [%v] generatePaymentXdr error:[%v] \n", sourceWallet.ID, paymentInfo.Destination, err)
		}
	}

	paymentInfo.Transaction = xdrBase64
	// if len(paymentInfo.Transaction) > 0 {
	// 	paymentInfo.SHash = algofuncs.SHash(paymentInfo.Transaction)
	// }
	// }

	paymentInfo.NetworkPassPhrase = network.GetBlockchainNetworkPassPhrase()

	if len(paymentInfo.TransactionSignature) == 0 && paymentInfo.Commit == 0 {
		return paymentInfo, nil, err
	}

	// if paymentInfo.SHash != algofuncs.SHash(paymentInfo.Transaction) && paymentInfo.Commit == 0 {
	// 	return paymentInfo, nil, &tPayErrors.ErrorTransactionMismatch{}
	// }
	if len(paymentInfo.ChannelAccountSignature) == 0 && len(paymentInfo.ChannelAccount) == 56 {
		return paymentInfo, nil, &tPayErrors.ErrorTransactionMismatch{Detail: "signature for channel account does not validate"}
	}

	if paymentInfo.Multiparty == 0 {
		//shared access disabled. submit to network is possible
		var txnHash string
		if len(paymentInfo.ChannelAccountSignature) > 0 && len(paymentInfo.ChannelAccount) == 56 {
			// txnHash, err = network.SubmitXdrWithSignatureChannelAccounts(client, sourceWallet.Signer, paymentInfo.ChannelAccount, xdrBase64, paymentInfo.TransactionSignature, paymentInfo.ChannelAccountSignature)
			txnHash, err = network.SubmitXdrWithSignatureChannelAccounts(client, sourceWallet.Signer, paymentInfo.ChannelAccount, paymentInfo.Transaction, paymentInfo.TransactionSignature, paymentInfo.ChannelAccountSignature)
			if err != nil {
				log.Println("#############################submit with channel account throws error:", err)

			}
		} else {
			// txnHash, err = network.SubmitXdrWithSignature(client, sourceWallet.Signer, xdrBase64, paymentInfo.TransactionSignature)
			txnHash, err = network.SubmitXdrWithSignature(client, sourceWallet.Signer, paymentInfo.Transaction, paymentInfo.TransactionSignature)
			if err != nil {
				log.Printf("[Pay] from [%v] to [%v] SubmitXdrWithSignature error:[%v] \n", sourceWallet.Alias, paymentInfo.Destination, err)
			}
		}
		paymentInfo.TransactionID = txnHash
		if err == nil {
			sourceWallet.InvalidateUserCache(gc)
			signerUser.InvalidateUserWalletCache(gc)
			signerUser.InvalidateUserCache(gc)
			if destinationUser != nil {
				if len(destinationUser.Username) > 0 {
					destinationUser.InvalidateUserCache(gc)
					destinationUser.InvalidateUserWalletCache(gc)
				}
			}
		}

		return paymentInfo, destinationUser, err
	}

	if paymentInfo.Commit == 0 {
		return paymentInfo, nil, nil
	}

	paymentInfo.TransactionID = "PENDING_AUTH"
	log.Printf("[Pay]shared access with approver permission enabled for %v \n", sourceWallet.Alias)
	id := uuid.NewString()
	assetOfPayment := os.Getenv("NATIVE_ASSET_CODE")
	if len(paymentInfo.AssetIssuer) == 56 {
		assetOfPayment = fmt.Sprintf("%v:%v...%v", paymentInfo.AssetCode, paymentInfo.AssetIssuer[0:4], paymentInfo.AssetIssuer[51:55])
	}
	var msgs string
	for i, m := range paymentInfo.Messages {
		msgs = m
		if i < len(paymentInfo.Messages)-1 {
			msgs = fmt.Sprintf("%s\n", msgs)
		}
	}
	description := fmt.Sprintf("Payment \nFrom: %v, \nTo: %v, \nAmount: %v %v", sourceWallet.Alias, paymentInfo.Destination, paymentInfo.Amount, assetOfPayment)
	if len(paymentInfo.Memo) > 0 {
		description = fmt.Sprintf("%v \nFor: %v", description, paymentInfo.Memo)

	}
	if len(msgs) > 0 {
		description = fmt.Sprintf("%v \nMessages: %v", description, msgs)
	}

	transactionByte, _ := json.Marshal(*paymentInfo)
	transactionStr := string(transactionByte)
	pendingAuth := userModels.PendingAuth{
		ID:                       id,
		Initiator:                signerUser.Username,
		InitiatorSignerPublicKey: signerUser.PrimarySigner,
		WalletPublicKey:          sourceWallet.ID,
		TransactionType:          "PAYMENT",
		Description:              description,
		TransactionSource:        paymentInfo.TransactionSource,
		ApprovalsNeeded:          sourceWallet.NumberOfApprovalsNeeded,
		TransactionXdr:           paymentInfo.Transaction,
		TransactionInfoStr:       &transactionStr,
	}
	//save and commit this to database
	e := db.Omit(clause.Associations).Create(&pendingAuth).Error
	if e != nil {
		log.Printf("[Pay] Error saving payment txn [%+v] transaction on pending auth table: %s\n", pendingAuth, e.Error())
		err = &tErrors.ErrorTemporaryServerError{}
		return paymentInfo, destinationUser, err
	}

	return paymentInfo, destinationUser, nil

}

func generatePaymentXdr(client *horizonclient.Client, owner *userModels.User, sourceWallet *userModels.UserWallet, paymentInfo *paymentModels.PaymentInfo, db *gorm.DB, gc *sharedconfig.GlobalConfig) (string, *userModels.User, error) {
	baseReserve := network.GetBlockchainBaseReserve()
	var tokenizedAssetIssuerMustSign bool
	charge := baseReserve.Mul(decimal.NewFromInt(3)).Truncate(7).String()
	nativeAssetCode := os.Getenv("NATIVE_ASSET_CODE")
	// var messages []string
	//check if it is public key payment

	publicKeyPayment := len(paymentInfo.Destination) == 56 || len(paymentInfo.Destination) == 69
	if publicKeyPayment {
		paymentInfo.Destination = strings.ToUpper(paymentInfo.Destination)
	} else {
		paymentInfo.Destination = strings.ToLower(paymentInfo.Destination)
	}

	var err error
	paymentInfo, err = ValidatePaymentInfo(paymentInfo)

	if err != nil {
		return "", nil, err
	}

	var amountToSend float64
	if amountToSend, err = strconv.ParseFloat(paymentInfo.AmountToPay, 64); err != nil {
		return "", nil, &tPayErrors.ErrorInvalidPaymentAmount{}
	}

	newAmountToSend := decimal.NewFromFloat(amountToSend).Truncate(7).String()

	var asset txnbuild.Asset = txnbuild.NativeAsset{}

	if len(paymentInfo.AssetIssuer) > 0 {
		asset = txnbuild.CreditAsset{Code: paymentInfo.AssetCode, Issuer: paymentInfo.AssetIssuer}
	}

	destinationInfo, getDestinationError := usersDB.GetUser(paymentInfo.Destination, db, gc)
	if getDestinationError != nil && !publicKeyPayment {
		return "", nil, &tPayErrors.ErrorPaymentDestinationDoesNotExist{}
	}
	if len(destinationInfo.Username) == 0 && !publicKeyPayment {
		log.Printf("[generatePaymentXdr]Could not get destination user for payment destination: %v\n", paymentInfo.Destination)

		return "", nil, &tPayErrors.ErrorPaymentDestinationDoesNotExist{}
	}
	//replace possible email and the rest
	if (!strings.Contains(paymentInfo.Destination, destinationInfo.Username)) && !publicKeyPayment {
		//email or phone or other ID used for payment. replace it.
		paymentInfo.Messages = append(paymentInfo.Messages, fmt.Sprintf("Notice: [%v] belongs to the wallet alias [%v] and will be used as the destination.", paymentInfo.Destination, destinationInfo.Username))
		paymentInfo.Destination = destinationInfo.Username
	}
	destinationWallet, _, _ := usersDB.GetWallet(paymentInfo.Destination, db)

	if len(destinationWallet.ID) == 0 && !publicKeyPayment {
		log.Printf("[generatePaymentXdr]Could not get destination wallet for payment destination: %v\n", paymentInfo.Destination)

		return "", nil, &tPayErrors.ErrorPaymentDestinationDoesNotExist{}
	}

	if destinationInfo.Suspended == 1 {
		return "", nil, &tErrors.ErrorUsernameIsSuspended{}
	}

	if !publicKeyPayment {
		if gc.IsValidTokenizedAsset(asset.GetCode()) && destinationInfo.KYCVerified == 0 && asset.GetIssuer() != destinationInfo.PublicKey {
			//if it is not a token burn also
			log.Printf("[SubscribeToTokenizedAsset] Error Destination Wallet owner %v has not met KYC status for asset %v\n", destinationInfo.Username, asset.GetCode())
			err = &tErrors.CustomError{Param: "destination", Err: "error-invalid-kyc", ErrMessage: fmt.Sprintf("%v has not passed/met the KYC requirement to receive this tokenized asset %v.", destinationInfo.Username, asset.GetCode())}
			return "", nil, err
		}
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
		paymentInfo.Destination = strings.ToUpper(paymentInfo.Destination)
		_, err := keypair.ParseAddress(paymentInfo.Destination)
		if err != nil {
			log.Printf("[generatePaymentXdr] error validating payment address [%v], %v\n", paymentInfo.Destination, err)

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
	destinationAccountExists, destinationAccountTrustsAsset, _, _, destinationBlockchainAccount, destinationAccountErr := network.BlockchainAccountProperties(client, destinationPublicKey, asset)
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

				message := fmt.Sprintf("The wallet %v is underfunded. %v %v will be deducted from your account to fund %v’s account. You only need to do this once for %v.", destinationWallet.Alias, charge, nativeAssetCode, destinationWallet.Alias, destinationWallet.Alias)

				paymentInfo.Messages = append(paymentInfo.Messages, message)
				// log.Printf("[generatePaymentXdr]message[0]: %v\n", message)

			}
			if gc.IsValidTokenizedAsset(asset.GetCode()) && asset.GetIssuer() != destinationInfo.PublicKey {
				//check if destination has done KYC
				if destinationInfo.KYCVerified == 0 {
					return "", nil, &tErrors.CustomError{
						Param:      "destination",
						Err:        "error-no-kyc",
						ErrMessage: fmt.Sprintf("%v does not meet KYC requirement to receive the asset %v", destinationWallet.Alias, asset.GetCode()),
					}
				}

			}
			if !destinationAccountTrustsAsset && !gc.IsValidTokenizedAsset(asset.GetCode()) {

				bantuAsset := userModels.BantuAsset{
					AssetCode:   asset.GetCode(),
					AssetIssuer: asset.GetIssuer(),
				}
				bcAsset, e := bantuAsset.GetBlockchainAssetProperty(gc)
				if e != nil {
					err = &tErrors.ErrorTemporaryServerError{}
					return "", nil, err

				}
				if len(bcAsset.Code) == 0 {
					err = &tErrors.ErrorTemporaryServerError{}
					return "", nil, err

				}

				if bcAsset.Flags.AuthRequired {

					err = &tErrors.CustomError{
						Param:      "destination",
						Err:        "error-destination-forbidden-to-receive-asset",
						ErrMessage: fmt.Sprintf("%v is a regulated asset. %v has not yet opted to receive this asset. Let the reciepient first add the asset to their trusted assets, successfully.", asset.GetCode(), destinationWallet.Alias),
					}
					return "", nil, err
				}

				message := fmt.Sprintf("%v has not yet opted in to receive the asset (%v) you are trying to send. %v %v will be deducted from your account to ensure that this transaction goes through. After this, %v will be able to receive %v anytime, without any further charges to you.", destinationWallet.Alias, paymentInfo.AssetCode, charge, nativeAssetCode, destinationWallet.Alias, paymentInfo.AssetCode)

				paymentInfo.Messages = append(paymentInfo.Messages, message)
				// log.Printf("[generatePaymentXdr]message[1]: %v\n", message)

			}
		}
	}
	if !publicKeyPayment {
		if destinationWallet.WalletType != 0 {
			if destinationWallet.Signer != sourceWallet.Signer {
				if !asset.IsNative() {
					//ensure transfer is from a sibbling account.
					if !destinationWallet.OwnerOfBlockchainAsset(asset.GetCode()) {
						err = &tErrors.CustomError{
							Param:      "destination",
							Err:        "error-action-forbidden",
							ErrMessage: fmt.Sprintf("%v is only allowed to be funded by a sibbling wallet, unless you are burning a token minted by %v. Only wallets belonging to the same account can fund a non-standard wallet.", destinationWallet.Alias, destinationWallet.Alias),
						}
						return "", nil, err
					}
				} else {
					err = &tErrors.CustomError{
						Param:      "destination",
						Err:        "error-action-forbidden",
						ErrMessage: fmt.Sprintf("%v is only allowed to be funded by a sibbling wallet. Only wallets belonging to the same account can fund a non-standard wallet.", destinationWallet.Alias),
					}
					return "", nil, err
				}

			}

		}
	}
	chanAccount := <-gc.ChannelAccounts
	defer func(c *keypair.Full) {
		gc.ChannelAccounts <- c
	}(chanAccount)
	// paymentInfo.Messages = messages
	_, _, _, _, chanSourceAccount, _ := network.BlockchainAccountProperties(client, chanAccount.Address(), txnbuild.NativeAsset{})
	sourceAccountExists, sourceAccountTrustsAsset, sourceAccountNativeBalance, sourceAccountCustomBalance, sourceAccount, sourceAccountErr := network.BlockchainAccountProperties(client, sourceWallet.ID, asset)

	if sourceAccountErr != nil {
		return "", nil, sourceAccountErr
	}

	if !sourceAccountExists {
		return "", nil, &tErrors.ErrorUnderfundedAccount{}
	}

	if !sourceAccountTrustsAsset {

		return "", nil, &tErrors.ErrorUnderfundedAccount{}
	}

	log.Printf("[generatePaymentXdr]obtained source account balance of owner %v:\n%v balance is %v\n%v balance is %v\n", owner.Username, nativeAssetCode, sourceAccountNativeBalance, asset.GetCode(), sourceAccountCustomBalance)

	amountToSendDec := decimal.NewFromFloat(amountToSend)
	//prevent minting of new tokens from this routine
	if !asset.IsNative() {
		if sourceWallet.ID == asset.GetIssuer() {
			err = &tErrors.CustomError{
				Param:      "destination",
				Err:        "error-source-forbidden-to-sending-asset",
				ErrMessage: fmt.Sprintf("%v, a token minting wallet, is forbidden from sending %v.", sourceWallet.Alias, asset.GetCode()),
			}

		}

	}
	if sourceWallet.ID != asset.GetIssuer() {

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

	var extraAccountKeyPair *keypair.Full = nil

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
				SourceAccount: sourceWallet.ID,
			})
		} else {
			ops = append(ops, &txnbuild.Payment{
				Destination:   destinationPublicKey,
				Amount:        newAmountToSend,
				Asset:         asset,
				SourceAccount: sourceWallet.ID,
			})
		}
	} else {
		//custom asset

		// claimable assets are for TrovoApp users only. it would return error above when destination does not trust asset

		if !destinationAccountExists {
			ops = append(ops, &txnbuild.CreateAccount{
				Destination:   destinationPublicKey,
				Amount:        charge,
				SourceAccount: sourceWallet.ID,
			})
		}

		if !destinationAccountTrustsAsset {
			if !publicKeyPayment {
				if gc.IsValidTokenizedAsset(asset.GetCode()) {

					return "", nil, &tErrors.CustomError{
						Param:      "destination",
						Err:        "error-destination-cannot-accept-asset",
						ErrMessage: fmt.Sprintf("%v does not accept the asset %v at this time.", destinationWallet.Alias, asset.GetCode()),
					}

				}
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
				if destinationWallet.WalletType == 0 {
					//standard wallet, create pending asset
					ops2, _tempAccountKeyPair, tokenIssuerMustSign, err := processDestinationWalletDoesNotTrustAsset(&destinationInfo, &destinationWallet, sourceAccount, asset, newAmountToSend, gc)
					if err != nil {
						return "", nil, err
					}

					tokenizedAssetIssuerMustSign = tokenIssuerMustSign

					extraAccountKeyPair = _tempAccountKeyPair

					ops = append(ops, ops2...)
				}

				if destinationWallet.WalletType == 2 || destinationWallet.WalletType == 3 {

					ops2, _dSignerAccountKeyPair, err := processCustodialDestinationWalletDoesNotTrustAsset(&destinationWallet, sourceAccount, asset, newAmountToSend)

					if err != nil {
						return "", nil, err
					}

					extraAccountKeyPair = _dSignerAccountKeyPair

					ops = append(ops, ops2...)
				}

			}

		} else {
			ops = append(ops, &txnbuild.Payment{
				Destination:   destinationPublicKey,
				Amount:        newAmountToSend,
				Asset:         asset,
				SourceAccount: sourceWallet.ID,
			})
		}

	}
	//service fee
	signForFeeTrustLine := 0
	serviceFee, e := decimal.NewFromString(paymentInfo.FeeAmount)
	if e != nil {
		serviceFee = decimal.Zero
	}
	if serviceFee.IsPositive() && os.Getenv("SHARED_ACCESS_FEE_ENABLED") == "1" {
		if paymentInfo.Multiparty == 1 {
			//process service fee
			feeLabel := paymentInfo.Fee + "%"
			// assetCode := os.Getenv("NATIVE_ASSET_CODE")
			// if !asset.IsNative() {
			// 	assetCode = asset.GetCode()
			// }
			feeKeypair := keypair.MustParseFull(os.Getenv("SHARED_ACCESS_FEE_WALLET"))
			feeAddress := feeKeypair.Address()

			if !asset.IsNative() {

				_, feeAccountTrustsAsset, _, _, _, _ := network.BlockchainAccountProperties(gc.BantuExpansionClient, feeAddress, asset)
				if !feeAccountTrustsAsset {
					signForFeeTrustLine = 1
					//establish trustline automatically
					ops = append(ops, &txnbuild.ChangeTrust{
						Line:          txnbuild.ChangeTrustAssetWrapper{Asset: asset},
						Limit:         "900000000000",
						SourceAccount: feeAddress,
					})

					if gc.IsValidTokenizedAsset(asset.GetCode()) {
						//check if it is a tokenized asset
						// allow trust from issuer to destination wallet
						ops = append(ops, &txnbuild.SetTrustLineFlags{
							Trustor:       feeAddress,
							Asset:         txnbuild.CreditAsset{Code: asset.GetCode(), Issuer: asset.GetIssuer()},
							SetFlags:      []txnbuild.TrustLineFlag{txnbuild.TrustLineAuthorized},
							SourceAccount: asset.GetIssuer(),
						})
						tokenizedAssetIssuerMustSign = true
					}

				}
			}
			ops = append(ops, &txnbuild.Payment{
				Destination:   feeAddress,
				Amount:        serviceFee.String(),
				SourceAccount: sourceWallet.ID,
				Asset:         asset,
			})
			// paymentInfo.Messages = append(paymentInfo.Messages, fmt.Sprintf("%v %v will be added from wallet %v as service fee (%v).", serviceFee.String(), assetCode, sourceWallet.Alias, feeLabel))
			paymentInfo.Messages = append(paymentInfo.Messages, fmt.Sprintf("%v will be added from wallet %v as service fee.", feeLabel, sourceWallet.Alias))

		}
	}

	var tx *txnbuild.Transaction
	// Construct the transaction that holds the operations to execute on the network
	if paymentInfo.Multiparty == 1 {
		paymentInfo.TransactionSource = chanSourceAccount.AccountID
		tx, err = txnbuild.NewTransaction(
			txnbuild.TransactionParams{
				SourceAccount:        chanSourceAccount,
				IncrementSequenceNum: true,
				Operations:           ops,
				BaseFee:              txnbuild.MinBaseFee,
				Preconditions: txnbuild.Preconditions{
					TimeBounds: txnbuild.NewInfiniteTimeout(),
				},
				Memo: txnbuild.MemoText(paymentInfo.Memo),
			},
		)
	} else {
		tx, err = txnbuild.NewTransaction(
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
	}

	if err != nil {
		log.Println("[generatePaymentXdr] error constructing transaction ", err)
		return "", nil, err
	}

	if paymentInfo.Multiparty == 1 {

		tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), chanAccount)

		if err != nil {
			log.Println("[generatePaymentXdr] error signing transaction with channelAccount key ", err)
			return "", nil, &tErrors.ErrorTemporaryServerError{}
		}
	}

	if signForFeeTrustLine == 1 && !asset.IsNative() && paymentInfo.Multiparty == 1 {
		feeKeypair := keypair.MustParseFull(os.Getenv("SHARED_ACCESS_FEE_WALLET"))
		tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), feeKeypair)

		if err != nil {
			log.Println("[generatePaymentXdr] error signing transaction with shared access fee key", err)
			return "", nil, &tErrors.ErrorTemporaryServerError{}
		}
	}

	if extraAccountKeyPair != nil {

		tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), extraAccountKeyPair)

		if err != nil {
			log.Println("[generatePaymentXdr] error signing transaction with temporary key ", err)
			return "", nil, &tErrors.ErrorTemporaryServerError{}
		}
	}

	if tokenizedAssetIssuerMustSign {
		log.Println("[generatePaymentXdr] <<<<<<<<<<<<<<<<<<<<<<<<<<<< signing transaction with issuer key>>>>>>>>>>>>>>>>>>>>>>>>")
		//get atprofile
		var tokenizationIssuerProfileWallet string

		if len(os.Getenv("TOKENIZATION_ISSUING_PROFILE_WALLET")) > 1 {
			tokenizationIssuerProfileWallet = strings.TrimSpace(os.Getenv("TOKENIZATION_ISSUING_PROFILE_WALLET"))
		}

		tokenizationIssuerProfileWalletKP := keypair.MustParseFull(tokenizationIssuerProfileWallet)

		tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), tokenizationIssuerProfileWalletKP)
		if err != nil {
			log.Println("[generatePaymentXdr] error signing transaction with issuer key to authorize trustline", err)
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

func generateMintingXdr(client *horizonclient.Client, owner *userModels.User, sourceWallet *userModels.UserWallet, mintingInfo *userModels.MintingInfo, db *gorm.DB, gc *sharedconfig.GlobalConfig) (string, *userModels.User, error) {
	baseReserve := network.GetBlockchainBaseReserve()
	var tokenizedAssetIssuerMustSign bool
	charge := baseReserve.Mul(decimal.NewFromInt(3)).Truncate(7).String()

	var err error
	mintingInfo, err = ValidateMintingInfo(mintingInfo)
	nativeAssetCode := os.Getenv("NATIVE_ASSET_CODE")
	if err != nil {
		return "", nil, err
	}
	// amountToSend := decimal.RequireFromString(mintingInfo.Amount).InexactFloat64()
	// amountToSendDec := decimal.RequireFromString(mintingInfo.Amount)
	newAmountToSend := decimal.RequireFromString(mintingInfo.Amount).Truncate(7).String()
	asset := txnbuild.CreditAsset{Code: mintingInfo.AssetCode, Issuer: mintingInfo.AssetIssuer}
	if sourceWallet.ID != asset.GetIssuer() {

		return "", nil, &tErrors.CustomError{
			Param:      "assetCode",
			Err:        "error cannot mint token you are not an issuer of",
			ErrMessage: "You can only mint tokens you are an issuer of",
		}

	}
	destinationInfo, getDestinationError := usersDB.GetUser(mintingInfo.Destination, db, gc)
	destinationWallet, _, destinationWalletError := usersDB.GetWallet(mintingInfo.Destination, db)

	if (getDestinationError != nil || destinationWalletError != nil) && len(mintingInfo.Destination) != 56 {
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

	// var destinationPublicKey string

	destinationPublicKey := destinationWallet.ID

	//perform ths checks of determining messages to be appended. if destination account property is not checked here, information would be returned without messages set.
	destinationAccountExists, destinationAccountTrustsAsset, _, _, destinationBlockchainAccount, destinationAccountErr :=
		network.BlockchainAccountProperties(client, destinationPublicKey, asset)
		//set base charge to be used in all places it is needed

		//check if to set charges messages

	//custom asset
	if !destinationAccountExists {

		message := fmt.Sprintf("The wallet %v is unfunded. %v %v will be deducted from your account to fund %v’s account. You only need to do this once for %v.", destinationWallet.Alias, charge, nativeAssetCode, destinationWallet.Alias, destinationWallet.Alias)

		mintingInfo.Messages = append(mintingInfo.Messages, message)
		// log.Printf("[generatePaymentXdr]message[0]: %v\n", message)

	}

	if !destinationAccountTrustsAsset {
		message := fmt.Sprintf("%v has not yet opted in to receive the asset (%v) you are trying to send. %v %v will be deducted from your account to ensure that this transaction goes through. After this, %v will be able to receive %v anytime, without any further charges to you.", destinationWallet.Alias, mintingInfo.AssetCode, charge, nativeAssetCode, destinationWallet.Alias, mintingInfo.AssetCode)

		mintingInfo.Messages = append(mintingInfo.Messages, message)
		// log.Printf("[generatePaymentXdr]message[1]: %v\n", message)

	}

	chanAccount := <-gc.ChannelAccounts
	defer func(c *keypair.Full) {
		gc.ChannelAccounts <- c
	}(chanAccount)
	// paymentInfo.Messages = messages
	_, _, _, _, chanSourceAccount, _ := network.BlockchainAccountProperties(client, chanAccount.Address(), txnbuild.NativeAsset{})
	sourceAccountExists, sourceAccountTrustsAsset, sourceAccountNativeBalance, sourceAccountCustomBalance, sourceAccount, sourceAccountErr := network.BlockchainAccountProperties(client, sourceWallet.ID, asset)

	if sourceAccountErr != nil {
		return "", nil, sourceAccountErr
	}

	if !sourceAccountExists {
		return "", nil, &tErrors.ErrorUnderfundedAccount{}
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
	var ops []txnbuild.Operation = make([]txnbuild.Operation, 0)

	var extraAccountKeyPair *keypair.Full = nil

	//custom asset

	// claimable assets are for TrovoApp users only. it would return error above when destination does not trust asset

	if !destinationAccountExists {
		ops = append(ops, &txnbuild.CreateAccount{
			Destination:   destinationPublicKey,
			Amount:        charge,
			SourceAccount: sourceWallet.ID,
		})
	}

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
		if destinationWallet.WalletType == 0 {
			//standard wallet, create pending asset
			ops2, _tempAccountKeyPair, tokenIssuerMustSign, err := processDestinationWalletDoesNotTrustAsset(&destinationInfo, &destinationWallet, sourceAccount, asset, newAmountToSend, gc)
			tokenizedAssetIssuerMustSign = tokenIssuerMustSign
			if err != nil {
				return "", nil, err
			}

			extraAccountKeyPair = _tempAccountKeyPair

			ops = append(ops, ops2...)
		}

		if destinationWallet.WalletType == 2 || destinationWallet.WalletType == 3 {

			ops2, _dSignerAccountKeyPair, err := processCustodialDestinationWalletDoesNotTrustAsset(&destinationWallet, sourceAccount, asset, newAmountToSend)

			if err != nil {
				return "", nil, err
			}

			extraAccountKeyPair = _dSignerAccountKeyPair

			ops = append(ops, ops2...)
		}

	} else {
		ops = append(ops, &txnbuild.Payment{
			Destination:   destinationPublicKey,
			Amount:        newAmountToSend,
			Asset:         asset,
			SourceAccount: sourceWallet.ID,
		})
	}

	var tx *txnbuild.Transaction
	// Construct the transaction that holds the operations to execute on the network
	if mintingInfo.Multiparty == 1 {
		mintingInfo.TransactionSource = chanSourceAccount.AccountID
		tx, err = txnbuild.NewTransaction(
			txnbuild.TransactionParams{
				SourceAccount:        chanSourceAccount,
				IncrementSequenceNum: true,
				Operations:           ops,
				BaseFee:              txnbuild.MinBaseFee,
				Preconditions: txnbuild.Preconditions{
					TimeBounds: txnbuild.NewInfiniteTimeout(),
				},
				Memo: txnbuild.MemoText(mintingInfo.Memo),
			},
		)
	} else {
		tx, err = txnbuild.NewTransaction(
			txnbuild.TransactionParams{
				SourceAccount:        sourceAccount,
				IncrementSequenceNum: true,
				Operations:           ops,
				BaseFee:              txnbuild.MinBaseFee,
				Preconditions: txnbuild.Preconditions{
					TimeBounds: txnbuild.NewInfiniteTimeout(),
				},
				Memo: txnbuild.MemoText(mintingInfo.Memo),
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

	if tokenizedAssetIssuerMustSign {
		log.Println("[generateMintingXdr] <<<<<<<<<<<<<<<<<<<<<<<<<<<< signing transaction with issuer key>>>>>>>>>>>>>>>>>>>>>>>>")
		//get atprofile
		var tokenizationIssuerProfileWallet string

		if len(os.Getenv("TOKENIZATION_ISSUING_PROFILE_WALLET")) > 1 {
			tokenizationIssuerProfileWallet = strings.TrimSpace(os.Getenv("TOKENIZATION_ISSUING_PROFILE_WALLET"))
		}

		tokenizationIssuerProfileWalletKP := keypair.MustParseFull(tokenizationIssuerProfileWallet)

		tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), tokenizationIssuerProfileWalletKP)
		if err != nil {
			log.Println("[generateMintingXdr] error signing transaction with issuer key to authorize trustline", err)
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
		mintingInfo.CallbackURLS = GetBlockchainAccountDataKey(*destinationBlockchainAccount, "orderPaymentCallbackUrl")

	}

	{
		owner.InvalidateUserCache(gc)
	}

	return xdrBase64, &destinationInfo, nil

}

func generatePaymentXdrWithChannelAccountPK(owner *userModels.User, sourceWallet *userModels.UserWallet, paymentInfo *paymentModels.PaymentInfo, gc *sharedconfig.GlobalConfig) (string, *userModels.User, error) {
	baseReserve := network.GetBlockchainBaseReserve()
	var tokenizedAssetIssuerMustSign bool
	// var messages []string
	//check if it is public key payment
	nativeAssetCode := os.Getenv("NATIVE_ASSET_CODE")
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
	destinationInfo, getDestinationError := usersDB.GetUser(paymentInfo.Destination, gc.DB, gc)
	destinationWallet, _, _ := usersDB.GetWallet(paymentInfo.Destination, gc.DB)
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
			log.Println("[generatePaymentXdr] error validating payment address [%v], %v", paymentInfo.Destination, err)

			return "", nil, &tPayErrors.ErrorInvalidPaymentDestinationPublicKey{}
		}

		message := "You are about to make payment to a public key directly. Please be sure of the address as the payment cannot be retrieved after confirmation."

		paymentInfo.Messages = append(paymentInfo.Messages, message)
		// log.Printf("[generatePaymentXdr]message for public key logged: %v\n", message)
	}

	//perform ths checks to determine messages to be appended. if destination account property is not checked here, information would be returned without messages set.
	destinationAccountExists, destinationAccountTrustsAsset, _, _, destinationBlockchainAccount, destinationAccountErr :=
		network.BlockchainAccountProperties(gc.BantuExpansionClient, destinationPublicKey, asset)
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

				message := fmt.Sprintf("The wallet %v is unfunded. %v %v will be deducted from your account to fund %v’s account. You only need to do this once for %v.", destinationWallet.Alias, charge, nativeAssetCode, destinationWallet.Alias, destinationWallet.Alias)

				paymentInfo.Messages = append(paymentInfo.Messages, message)
				// log.Printf("[generatePaymentXdr]message[0]: %v\n", message)

			}

			if !destinationAccountTrustsAsset {
				message := fmt.Sprintf("%v has not yet opted in to receive (%v), that you are trying to send. %v %v will be deducted from your account to ensure that this transaction goes through. After this, %v will be able to receive %v anytime, without any further charges to you.", destinationWallet.Alias, paymentInfo.AssetCode, charge, os.Getenv("NATIVE_ASSET_CODE"), destinationWallet.Alias, paymentInfo.AssetCode)

				paymentInfo.Messages = append(paymentInfo.Messages, message)
				// log.Printf("[generatePaymentXdr]message[1]: %v\n", message)

			}
		}
	}

	if !publicKeyPayment {
		if destinationWallet.WalletType != 0 {
			if destinationWallet.Signer != sourceWallet.Signer {
				if !asset.IsNative() {
					//ensure transfer is from a sibbling account.
					if !destinationWallet.OwnerOfBlockchainAsset(asset.GetCode()) {
						err = &tErrors.CustomError{
							Param:      "destination",
							Err:        "error-action-forbidden",
							ErrMessage: fmt.Sprintf("%v is only allowed to be funded by a sibbling wallet, unless you are burning a token minted by %v. Only wallets belonging to the same account can fund a non-standard wallet.", destinationWallet.Alias, destinationWallet.Alias),
						}
						return "", nil, err
					}
				} else {
					err = &tErrors.CustomError{
						Param:      "destination",
						Err:        "error-action-forbidden",
						ErrMessage: fmt.Sprintf("%v is only allowed to be funded by a sibbling wallet. Only wallets belonging to the same account can fund a non-standard wallet.", destinationWallet.Alias),
					}
					return "", nil, err
				}

			}

		}
	}

	sourceAccountExists, sourceAccountTrustsAsset, sourceAccountNativeBalance, sourceAccountCustomBalance, sourceAccount, sourceAccountErr := network.BlockchainAccountProperties(gc.BantuExpansionClient, sourceWallet.ID, asset)
	channelSourceAccountExists, _, channelSourceAccountNativeBalance, _, channelSourceAccount, channelSourceAccountErr := network.BlockchainAccountProperties(gc.BantuExpansionClient, paymentInfo.ChannelAccount, txnbuild.NativeAsset{})

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

	log.Printf("[generatePaymentXdrChannelAccountPk]obtained source account balance:\n%v balance is %v\n%v balance is %v\n", nativeAssetCode, sourceAccountNativeBalance, asset.GetCode(), sourceAccountCustomBalance)

	amountToSendDec := decimal.NewFromFloat(amountToSend)

	if sourceWallet.ID != asset.GetIssuer() {

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

	var extraAccountKeyPair *keypair.Full = nil

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
				SourceAccount: sourceWallet.ID,
			})
		} else {
			ops = append(ops, &txnbuild.Payment{
				Destination:   destinationPublicKey,
				Amount:        newAmountToSend,
				Asset:         asset,
				SourceAccount: sourceWallet.ID,
			})
		}
	} else {
		//custom asset

		// claimable assets are for trovotech customers only. it would return error above when destination does not trust asset

		if !destinationAccountExists {
			ops = append(ops, &txnbuild.CreateAccount{
				Destination:   destinationPublicKey,
				Amount:        charge,
				SourceAccount: sourceWallet.ID,
			})
		}

		if !destinationAccountTrustsAsset {
			if !publicKeyPayment {
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
				if destinationWallet.WalletType == 0 {
					//standard wallet, create pending asset
					ops2, _tempAccountKeyPair, tokenIssuerMustSign, err := processDestinationWalletDoesNotTrustAsset(&destinationInfo, &destinationWallet, sourceAccount, asset, newAmountToSend, gc)
					tokenizedAssetIssuerMustSign = tokenIssuerMustSign
					if err != nil {
						return "", nil, err
					}

					extraAccountKeyPair = _tempAccountKeyPair

					ops = append(ops, ops2...)
				}

				if destinationWallet.WalletType == 2 || destinationWallet.WalletType == 3 {

					ops2, _dSignerAccountKeyPair, err := processCustodialDestinationWalletDoesNotTrustAsset(&destinationWallet, sourceAccount, asset, newAmountToSend)

					if err != nil {
						return "", nil, err
					}

					extraAccountKeyPair = _dSignerAccountKeyPair

					ops = append(ops, ops2...)
				}

			}

		} else {
			ops = append(ops, &txnbuild.Payment{
				Destination:   destinationPublicKey,
				Amount:        newAmountToSend,
				Asset:         asset,
				SourceAccount: sourceWallet.ID,
			})
		}

	}

	//service fee
	fee := decimal.RequireFromString(os.Getenv("SHARED_ACCESS_FEE_AMOUNT"))
	if !fee.IsZero() {
		if paymentInfo.Multiparty == 1 {
			//process service fee
			if len(os.Getenv("SHARED_ACCESS_FEE_ASSET_ISSUER")) != 56 {
				ops = append(ops, &txnbuild.Payment{
					Destination:   os.Getenv("SHARED_ACCESS_FEE_ADDRESS"),
					Amount:        os.Getenv("SHARED_ACCESS_FEE_AMOUNT"),
					SourceAccount: sourceWallet.ID,
					Asset:         txnbuild.NativeAsset{},
				})
				paymentInfo.Messages = append(paymentInfo.Messages, fmt.Sprintf("%v %v will be deducted from wallet %v as service fee.", os.Getenv("SHARED_ACCESS_FEE_AMOUNT"), nativeAssetCode, sourceWallet.Alias))

			} else {
				ops = append(ops, &txnbuild.Payment{
					Destination:   os.Getenv("SHARED_ACCESS_FEE_ADDRESS"),
					Amount:        os.Getenv("SHARED_ACCESS_FEE_AMOUNT"),
					SourceAccount: sourceWallet.ID,
					Asset:         txnbuild.CreditAsset{Code: os.Getenv("SHARED_ACCESS_FEE_ASSET_CODE"), Issuer: os.Getenv("SHARED_ACCESS_FEE_ASSET_ISSUER")},
				})
				paymentInfo.Messages = append(paymentInfo.Messages, fmt.Sprintf("%v %v will be deducted from wallet %v as service fee.", os.Getenv("SHARED_ACCESS_FEE_AMOUNT"), os.Getenv("SHARED_ACCESS_FEE_ASSET_CODE"), sourceWallet.Alias))

			}

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

	if extraAccountKeyPair != nil {

		tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), extraAccountKeyPair)

		if err != nil {
			log.Println("[generatePaymentXdrWithChannelAccountPK] error signing transaction with temporary key ", err)
			return "", nil, &tErrors.ErrorTemporaryServerError{}
		}
	}
	if tokenizedAssetIssuerMustSign {
		log.Println("[generatePaymentXdrWithChannelAccountPK] <<<<<<<<<<<<<<<<<<<<<<<<<<<< signing transaction with issuer key>>>>>>>>>>>>>>>>>>>>>>>>")
		//get atprofile
		var tokenizationIssuerProfileWallet string

		if len(os.Getenv("TOKENIZATION_ISSUING_PROFILE_WALLET")) > 1 {
			tokenizationIssuerProfileWallet = strings.TrimSpace(os.Getenv("TOKENIZATION_ISSUING_PROFILE_WALLET"))
		}

		tokenizationIssuerProfileWalletKP := keypair.MustParseFull(tokenizationIssuerProfileWallet)

		tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), tokenizationIssuerProfileWalletKP)
		if err != nil {
			log.Println("[generatePaymentXdrWithChannelAccountPK] error signing transaction with issuer key to authorize trustline", err)
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
	{
		owner.InvalidateUserCache(gc)
	}
	return xdrBase64, &destinationInfo, nil
}

func processDestinationWalletDoesNotTrustAsset(destinationUser *userModels.User, destinationWallet *userModels.UserWallet, sourceAccount *horizon.Account, asset txnbuild.Asset, amountToSend string, gc *sharedconfig.GlobalConfig) ([]txnbuild.Operation, *keypair.Full, bool, error) {

	ops := make([]txnbuild.Operation, 0)

	var signerKeyPairToReturn *keypair.Full = nil
	var tokenizedAssetIssuerMustSign bool
	tempAccountKeypair, tempAccountError := network.TempAccountKeypair(destinationWallet.ID)

	if tempAccountError != nil {
		log.Printf("[processDestinationAssetDoesNotTrustAsset] error generating temporary account %v\n", tempAccountError)
		return ops, tempAccountKeypair, false, &tErrors.ErrorTemporaryServerError{}
	}

	// var tempAccount txnbuild.Account = &txnbuild.SimpleAccount{AccountID: tempAccountKeypair.Address(), Sequence: 0}

	tempAccountExists, tempAccountTrustsAsset, _, _, tempAccountSource, tempAccountError :=
		network.BlockchainAccountProperties(gc.BantuExpansionClient, tempAccountKeypair.FromAddress().Address(), asset)

	if tempAccountError != nil {
		log.Printf("[processDestinationAssetDoesNotTrustAsset] error looking up temporary account %v\n", tempAccountError)
		return ops, tempAccountKeypair, false, &tErrors.ErrorTemporaryServerError{}
	}

	{

		var _wallet userModels.UserWallet

		//check if temp account exists in buds or not.
		err := gc.DB.Where("id = ?", destinationWallet.ID).First(&_wallet).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ops, nil, false, &tErrors.ErrorInvalidPublicKey{}

			}
			log.Println("[processDestinationAssetDoesNotTrustAsset]", err)
			return ops, nil, false, &tErrors.ErrorTemporaryServerError{}

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
			dbSaveError := gc.DB.Omit(clause.Associations).Save(&_wallet).Error

			if dbSaveError != nil {
				log.Printf("[processDestinationAssetDoesNotTrustAsset]db temp save error %v\n", dbSaveError)
				return ops, nil, tokenizedAssetIssuerMustSign, &tErrors.ErrorTemporaryServerError{}
			}
		}

	}

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

		if gc.IsValidTokenizedAsset(asset.GetCode()) {
			//check if it is a tokenized asset
			// allow trust from issuer to destination wallet
			ops = append(ops, &txnbuild.SetTrustLineFlags{
				Trustor:       tempAccountKeypair.FromAddress().Address(),
				Asset:         txnbuild.CreditAsset{Code: asset.GetCode(), Issuer: asset.GetIssuer()},
				SetFlags:      []txnbuild.TrustLineFlag{txnbuild.TrustLineAuthorized},
				SourceAccount: asset.GetIssuer(),
			})
			tokenizedAssetIssuerMustSign = true
		}

	}

	ops = append(ops, &txnbuild.Payment{
		Destination:   tempAccountKeypair.Address(),
		Amount:        amountToSend,
		Asset:         asset,
		SourceAccount: sourceAccount.AccountID,
	})

	return ops, signerKeyPairToReturn, tokenizedAssetIssuerMustSign, nil

}

// func processCustodialDestinationWalletDoesNotTrustAsset(destinationUser *userModels.User, destinationWallet *userModels.UserWallet, sourceAccount, destinationAccount *horizon.Account, asset txnbuild.Asset, amountToSend string) (ops []txnbuild.Operation, signerKeyPairToReturn *keypair.Full, err error) {
func processCustodialDestinationWalletDoesNotTrustAsset(destinationWallet *userModels.UserWallet, sourceAccount *horizon.Account, asset txnbuild.Asset, amountToSend string) (ops []txnbuild.Operation, signerKeyPairToReturn *keypair.Full, err error) {

	ops = make([]txnbuild.Operation, 0)

	if destinationWallet.WalletType == 1 && !asset.IsNative() {
		err = &tErrors.CustomError{
			Param:      "destination",
			Err:        "error-destination-forbidden-to-receive-asset",
			ErrMessage: fmt.Sprintf("%v, a token minting wallet, is forbidden from receiving %v", destinationWallet.Alias, asset.GetCode()),
		}
		return
	}

	ops = append(ops, &txnbuild.ChangeTrust{
		Line:          txnbuild.ChangeTrustAssetWrapper{Asset: asset},
		Limit:         "900000000000",
		SourceAccount: destinationWallet.ID,
	})

	ops = append(ops, &txnbuild.Payment{
		Destination:   destinationWallet.ID,
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

	return clientAccount.Data, nil
}

func MintAsset(signerUser *userModels.User, sourceWallet *userModels.UserWallet, mintingInfo *userModels.MintingInfo, gc *sharedconfig.GlobalConfig) (*userModels.MintingInfo, *userModels.User, error) {
	db := gc.DB
	client := network.GetBlockchainClient()
	var xdrBase64 string
	var destinationUser *userModels.User
	var err error
	walletHasViewOnlyAccess := true

	walletHasViewOnlyAccess = sourceWallet.HasViewOnlyAccess(gc)
	if !walletHasViewOnlyAccess {
		mintingInfo.Multiparty = 1

	}
	if walletHasViewOnlyAccess {
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
	if len(mintingInfo.ChannelAccountSignature) == 0 && len(mintingInfo.ChannelAccount) == 56 {
		return mintingInfo, nil, &tPayErrors.ErrorTransactionMismatch{Detail: "signature for channel account does not validate"}
	}

	if mintingInfo.Multiparty == 0 {
		//shared access disabled. submit to network is possible
		var txnHash string
		if len(mintingInfo.ChannelAccountSignature) > 0 && len(mintingInfo.ChannelAccount) == 56 {
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
	if len(mintingInfo.AssetIssuer) == 56 {
		assetOfPayment = fmt.Sprintf("%v:%v...%v", mintingInfo.AssetCode, mintingInfo.AssetIssuer[0:3], mintingInfo.AssetIssuer[52:55])
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
		ID:                       id,
		Initiator:                signerUser.Username,
		InitiatorSignerPublicKey: signerUser.PrimarySigner,
		WalletPublicKey:          sourceWallet.ID,
		TransactionType:          "MINT TOKEN",
		Description:              description,
		TransactionSource:        mintingInfo.TransactionSource,
		ApprovalsNeeded:          sourceWallet.NumberOfApprovalsNeeded,
		TransactionXdr:           mintingInfo.Transaction,
		TransactionInfoStr:       &transactionStr,
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

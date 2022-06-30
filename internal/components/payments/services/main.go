package payments

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"
	tPayErrors "trovo-wallet-api/internal/components/payments/errors"
	payments "trovo-wallet-api/internal/components/payments/models"
	usersdb "trovo-wallet-api/internal/components/users/db"
	users "trovo-wallet-api/internal/components/users/models"

	"log"
	"strconv"
	paymentsDB "trovo-wallet-api/internal/components/payments/db"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/network"

	"github.com/ecnepsnai/discord"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
)

//Pay function sends a payment from user to another user
func Pay(owner *paymentsDB.User, wallet *paymentsDB.UserWallet, paymentInfo *payments.PaymentInfo, db *gorm.DB) (*payments.PaymentInfo, *paymentsDB.User, error) {

	client := network.GetBlockchainClient()
	var xdrBase64 string
	var destinationUser *paymentsDB.User
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
		xdrBase64, destinationUser, err = generatePaymentXdrWithChannelAccountPK(client, owner, destinationUser, wallet, paymentInfo, db)

		if err != nil {
			log.Printf("[Pay] from [%v] to [%v] generatePaymentXdrWithChannelAccountPK error:[%v] \n", wallet.Alias, paymentInfo.Destination, err)
		}
	} else {

		xdrBase64, destinationUser, err = generatePaymentXdr(client, owner, destinationUser, wallet, paymentInfo, db)
		if err != nil {
			log.Printf("[Pay] from [%v] to [%v] generatePaymentXdr error:[%v] \n", wallet.ID, paymentInfo.Destination, err)
		}
	}
	oldTransaction := paymentInfo.Transaction

	paymentInfo.Transaction = xdrBase64
	paymentInfo.NetworkPassPhrase = network.GetBlockchainNetworkPassPhrase()

	if len(paymentInfo.TransactionSignature) == 0 {
		return paymentInfo, nil, err
	}

	if xdrBase64 != oldTransaction {
		return paymentInfo, nil, &tPayErrors.ErrorTransactionMismatch{}
	}
	if len(paymentInfo.ChannelAccountSignature) == 0 && len(paymentInfo.ChannelAccount) == 56 {
		return paymentInfo, nil, &tPayErrors.ErrorTransactionMismatch{Detail: "signature for channel account does not validate"}
	}

	if wallet.ManagedAccessEnabled == 0 {
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
	} else {
		//put routine for managed access > 0
		log.Printf("managed access enabled for %v \n", wallet.Alias)
	}
	return paymentInfo, destinationUser, err

}

//PayWithChannelAccount function sends a payment from user to another user using channel accounts
func PayWithChannelAccount(senderPublicKey string, paymentInfo *payments.PaymentInfo, db *gorm.DB) (*payments.PaymentInfo, error) {
	// db, err := conDB.OpenDb()
	// if err != nil {
	// 	log.Println("-------------- ------DB error in Pay:", err)
	// 	return nil, err
	// }
	channelAccount := keypair.MustParseFull(os.Getenv("CHANNEL_ACCOUNT"))
	client := network.GetBlockchainClient()

	xdrBase64, err := generatePaymentXdrWithChannelAccount(client, senderPublicKey, paymentInfo, channelAccount, db)

	oldTransaction := paymentInfo.Transaction

	paymentInfo.Transaction = xdrBase64
	paymentInfo.NetworkPassPhrase = network.GetBlockchainNetworkPassPhrase()

	if len(paymentInfo.TransactionSignature) == 0 {
		return paymentInfo, err
	}

	if xdrBase64 != oldTransaction {
		return paymentInfo, &tPayErrors.ErrorTransactionMismatch{}
	}

	txnHash, err := network.SubmitXdrWithSignature(client, senderPublicKey, xdrBase64, paymentInfo.TransactionSignature)

	paymentInfo.TransactionID = txnHash
	return paymentInfo, err
}
func generatePaymentXdr(client *horizonclient.Client, owner *paymentsDB.User, destinationUser *paymentsDB.User, wallet *paymentsDB.UserWallet, paymentInfo *payments.PaymentInfo, db *gorm.DB) (string, *paymentsDB.User, error) {
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

	destinationInfo, getDestinationError := usersdb.GetUser(paymentInfo.Destination, db)
	destinationWallet, _, _ := usersdb.GetWallet(paymentInfo.Destination, db)

	charge := baseReserve.Mul(decimal.NewFromInt(3)).Truncate(7).String()
	if getDestinationError != nil && len(paymentInfo.Destination) != 56 {
		return "", nil, &tPayErrors.ErrorPaymentDestinationDoesNotExist{}
	}
	if destinationInfo.Suspended == 1 {
		return "", nil, &tErrors.ErrorUsernameIsSuspended{}
	}
	// if banned, errBanned := usersdb.PublicKeyIsBanned(destinationInfo.PublicKey, db); banned {
	// 	return "", nil, errBanned
	// }
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

		message := "Important: You are about to make payment to a public key directly. Please be sure of the address as the payment cannot be retrieved after confirmation."

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

				message := fmt.Sprintf("Important: The account of %v is unfunded. %v XBN will be deducted from your account to fund %v’s account. You only need to do this once for %v.", destinationInfo.Username, charge, destinationInfo.Username, destinationInfo.Username)

				paymentInfo.Messages = append(paymentInfo.Messages, message)
				// log.Printf("[generatePaymentXdr]message[0]: %v\n", message)

			}

			if !destinationAccountTrustsAsset {
				message := fmt.Sprintf("Important: %v has not yet activated the asset (%v) you are trying to send. %v XBN will be deducted from your account to ensure that this transaction goes through. After this, %v will be able to receive %v anytime, without any further charges to you.", destinationInfo.Username, paymentInfo.AssetCode, charge, destinationInfo.Username, paymentInfo.AssetCode)

				paymentInfo.Messages = append(paymentInfo.Messages, message)
				// log.Printf("[generatePaymentXdr]message[1]: %v\n", message)

			}
		}
	}

	// paymentInfo.Messages = messages
	sourceAccountExists, sourceAccountTrustsAsset, sourceAccountNativeBalance, sourceAccountCustomBalance, sourceAccount, sourceAccountErr := network.BlockchainAccountProperties(client, wallet.Signer, asset)

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
				SourceAccount: sourceAccount.AccountID,
			})
		} else {
			ops = append(ops, &txnbuild.Payment{
				Destination:   destinationPublicKey,
				Amount:        newAmountToSend,
				Asset:         asset,
				SourceAccount: sourceAccount.AccountID,
			})
		}
	} else {
		//custom asset

		// claimable assets are for trovo wallet users only. it would return error above when destination does not trust asset

		if !destinationAccountExists {
			ops = append(ops, &txnbuild.CreateAccount{
				Destination:   destinationPublicKey,
				Amount:        charge,
				SourceAccount: sourceAccount.AccountID,
			})
		}

		if !destinationAccountTrustsAsset {
			if !publicKeyPayment {
				ops2, _tempAccountKeyPair, err := processDestinationAssetDoesNotTrustAsset(client, sourceAccount, destinationPublicKey, asset, newAmountToSend, db)

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
				SourceAccount: sourceAccount.AccountID,
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

	return xdrBase64, destinationUser, nil
}

func generatePaymentXdrWithChannelAccount(client *horizonclient.Client, senderPublicKey string, paymentInfo *payments.PaymentInfo, channelAccount *keypair.Full, db *gorm.DB) (string, error) {
	baseReserve := network.GetBlockchainBaseReserve()
	var messages []string
	//check if it is publc key payment
	publicKeyPayment := len(paymentInfo.Destination) == 56
	var err error
	paymentInfo, err = ValidatePaymentInfo(paymentInfo)

	if err != nil {
		return "", err
	}

	var amountToSend float64
	if amountToSend, err = strconv.ParseFloat(paymentInfo.Amount, 64); err != nil {
		return "", &tPayErrors.ErrorInvalidPaymentAmount{}
	}

	newAmountToSend := decimal.NewFromFloat(amountToSend).Truncate(7).String()

	var asset txnbuild.Asset = txnbuild.NativeAsset{}

	if len(paymentInfo.AssetCode) != 0 {
		asset = txnbuild.CreditAsset{Code: paymentInfo.AssetCode, Issuer: paymentInfo.AssetIssuer}
	}

	destinationInfo, getDestinationError := usersdb.GetUser(paymentInfo.Destination, db)
	destinationWallet, _, _ := usersdb.GetWallet(paymentInfo.Destination, db)
	charge := baseReserve.Mul(decimal.NewFromInt(3)).Truncate(7).String()
	if getDestinationError != nil && len(paymentInfo.Destination) != 56 {
		return "", &tPayErrors.ErrorPaymentDestinationDoesNotExist{}
	}
	if !publicKeyPayment {
		paymentInfo.DestinationFirstName = destinationInfo.FirstName
		if destinationInfo.LastName != nil {
			paymentInfo.DestinationLastName = *destinationInfo.LastName
		}

		if destinationInfo.ImageThumbnailURL != nil {
			paymentInfo.DestinationThumbnail = *destinationInfo.ImageThumbnailURL
		}
	} else {
		//parse public key
		_, err := keypair.ParseAddress(paymentInfo.Destination)
		if err != nil {
			return "", &tPayErrors.ErrorInvalidPaymentDestinationPublicKey{}
		}

		message := "Important: You are about to make payment to a public key directly. Please be sure of the address as the payment cannot be retrieved after confirmation."

		messages = append(messages, message)
		log.Printf("[generatePaymentXdr]message for public key logged: %v\n", message)
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
			return "", &tPayErrors.ErrorDestinationPublicKeyCannotReceiveAsset{}
		}
	} else {
		//check if to set charges messages
		if !asset.IsNative() {
			//custom asset
			if !destinationAccountExists {

				message := fmt.Sprintf("Important: The account of %v is unfunded. %v XBN will be deducted from your account to fund %v’s account. You only need to do this once for %v.", destinationInfo.Username, charge, destinationInfo.Username, destinationInfo.Username)

				messages = append(messages, message)
				log.Printf("[generatePaymentXdr]message[0]: %v\n", message)

			}

			if !destinationAccountTrustsAsset {
				message := fmt.Sprintf("Important: %v has not yet activated the asset (%v) you are trying to send. To reduce spam, %v XBN will be deducted from your account. No extra XBN will be deducted from you after %v claims your asset.", destinationInfo.Username, paymentInfo.AssetCode, charge, destinationInfo.Username)

				messages = append(messages, message)
				log.Printf("[generatePaymentXdr]message[1]: %v\n", message)

			}
		}
	}

	paymentInfo.Messages = messages
	sourceAccountExists, sourceAccountTrustsAsset, sourceAccountNativeBalance, sourceAccountCustomBalance, sourceAccount, sourceAccountErr := network.BlockchainAccountProperties(client, senderPublicKey, asset)
	_, _, _, _, channelSourceAccount, _ := network.BlockchainAccountProperties(client, channelAccount.Address(), txnbuild.NativeAsset{})

	if sourceAccountErr != nil {
		return "", sourceAccountErr
	}

	if !sourceAccountExists {
		return "", &tErrors.ErrorUnderfundedAccount{}
	}

	if !sourceAccountTrustsAsset {

		return "", &tErrors.ErrorUnderfundedAccount{}
	}

	log.Printf("[generatePaymentXdr]obtained sourced account info \n")

	log.Printf("[generatePaymentXdr]obtained source account balance is %v, custom account balance %v\n", sourceAccountNativeBalance, sourceAccountCustomBalance)

	amountToSendDec := decimal.NewFromFloat(amountToSend)

	if senderPublicKey != asset.GetIssuer() {

		if asset.IsNative() {
			if sourceAccountNativeBalance.LessThan(amountToSendDec) {
				return "", &tErrors.ErrorUnderfundedAccount{}
			}
		} else {
			if sourceAccountCustomBalance.LessThan(amountToSendDec) {
				return "", &tErrors.ErrorUnderfundedAccount{}
			}
		}
	}

	//check if destination account exists
	if destinationAccountErr != nil {
		log.Println("[generatePaymentXdr]destination Account error:", destinationAccountErr)
		return "", destinationAccountErr
	}
	var ops []txnbuild.Operation = make([]txnbuild.Operation, 0)

	var tempAccountKeyPair *keypair.Full = nil

	if asset.IsNative() {
		//native asset
		if !destinationAccountExists {

			if amountToSendDec.LessThan(baseReserve.Mul(decimal.NewFromInt(3))) {
				log.Printf("[generatePaymentXdr] trying to create account with %v, meanwhile you need %v\n", amountToSendDec, baseReserve.Mul(decimal.NewFromInt(3)))
				return "", &tPayErrors.ErrorInsufficientAmountToFundAccount{}
			}
			ops = append(ops, &txnbuild.CreateAccount{
				Destination:   destinationPublicKey,
				Amount:        newAmountToSend,
				SourceAccount: senderPublicKey,
			})
		} else {
			ops = append(ops, &txnbuild.Payment{
				Destination:   destinationPublicKey,
				Amount:        newAmountToSend,
				Asset:         asset,
				SourceAccount: senderPublicKey,
			})
		}
	} else {
		//custom asset

		// claimable assets are for trovowallet users only. it would return error above when destination does not trust asset

		if !destinationAccountExists {
			ops = append(ops, &txnbuild.CreateAccount{
				Destination:   destinationPublicKey,
				Amount:        charge,
				SourceAccount: senderPublicKey,
			})
		}

		if !destinationAccountTrustsAsset {
			if !publicKeyPayment {
				ops2, _tempAccountKeyPair, err := processDestinationAssetDoesNotTrustAsset(client, sourceAccount, destinationPublicKey, asset, newAmountToSend, db)

				if err != nil {
					return "", err
				}

				tempAccountKeyPair = _tempAccountKeyPair

				ops = append(ops, ops2...)
			}

		} else {
			ops = append(ops, &txnbuild.Payment{
				Destination:   destinationPublicKey,
				Amount:        newAmountToSend,
				Asset:         asset,
				SourceAccount: senderPublicKey,
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
		log.Println("[generatePaymentXdr] error constructing transaction ", err)
		return "", err
	}

	if tempAccountKeyPair != nil {

		tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), tempAccountKeyPair)

		if err != nil {
			log.Println("[generatePaymentXdr] error signing transaction with temporary key ", err)
			return "", &tErrors.ErrorTemporaryServerError{}
		}
	}

	tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), channelAccount)

	if err != nil {
		log.Println("[generatePaymentXdr] error signing transaction with temporary key ", err)
		return "", &tErrors.ErrorTemporaryServerError{}
	}

	xdrBase64, err := tx.Base64()

	if err != nil {
		log.Println("[generatePaymentXdr] error getting txn base64", err)
		return "", err
	}

	if destinationBlockchainAccount != nil {
		paymentInfo.CallbackURLS = GetBlockchainAccountDataKey(*destinationBlockchainAccount, "orderPaymentCallbackUrl")

	}

	return xdrBase64, nil
}

func generatePaymentXdrWithChannelAccountPK(client *horizonclient.Client, owner, destinationUser *paymentsDB.User, wallet *paymentsDB.UserWallet, paymentInfo *payments.PaymentInfo, db *gorm.DB) (string, *paymentsDB.User, error) {
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

	charge := baseReserve.Mul(decimal.NewFromInt(3)).Truncate(7).String()

	if !publicKeyPayment {
		paymentInfo.DestinationFirstName = destinationUser.FirstName
		if destinationUser.LastName != nil {
			paymentInfo.DestinationLastName = *destinationUser.LastName
		}

		paymentInfo.DestinationVerified = destinationUser.Verified
		if destinationUser.ImageThumbnailURL != nil {
			paymentInfo.DestinationThumbnail = *destinationUser.ImageThumbnailURL
		}
	} else {
		//parse public key
		_, err := keypair.ParseAddress(paymentInfo.Destination)
		if err != nil {
			return "", nil, &tPayErrors.ErrorInvalidPaymentDestinationPublicKey{}
		}

		message := "Important: You are about to make payment to a public key directly. Please be sure of the address as the payment cannot be retrieved after confirmation."

		paymentInfo.Messages = append(paymentInfo.Messages, message)
		// log.Printf("[generatePaymentXdr]message for public key logged: %v\n", message)
	}
	var destinationPublicKey string
	if publicKeyPayment {
		destinationPublicKey = paymentInfo.Destination
	} else {
		destinationPublicKey = destinationUser.PublicKey
	}
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

				message := fmt.Sprintf("Important: The account of %v is unfunded. %v XBN will be deducted from your account to fund %v’s account. You only need to do this once for %v.", destinationUser.Username, charge, destinationUser.Username, destinationUser.Username)

				paymentInfo.Messages = append(paymentInfo.Messages, message)
				// log.Printf("[generatePaymentXdr]message[0]: %v\n", message)

			}

			if !destinationAccountTrustsAsset {
				message := fmt.Sprintf("Important: %v has not yet activated the asset (%v), that you are trying to send. %v XBN will be deducted from your account to ensure that this transaction goes through. After this, %v will be able to receive %v anytime, without any further charges to you.", destinationUser.Username, paymentInfo.AssetCode, charge, destinationUser.Username, paymentInfo.AssetCode)

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

	// rid, e := ProcessBailArrestedAccount(client, owner, sourceAccount, db)
	// if e != nil {
	// 	//  error while releasing from prison
	// 	log.Println("[generatePaymentXdr] error:", e)
	// } else {
	// 	//account freed from prison
	// 	log.Println("[generatePaymentXdr] account freed from prison", senderPublicKey, "hash:", rid)
	// }

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

		// claimable assets are for bantupay customers only. it would return error above when destination does not trust asset

		if !destinationAccountExists {
			ops = append(ops, &txnbuild.CreateAccount{
				Destination:   destinationPublicKey,
				Amount:        charge,
				SourceAccount: wallet.ID,
			})
		}

		if !destinationAccountTrustsAsset {
			if !publicKeyPayment {
				ops2, _tempAccountKeyPair, err := processDestinationAssetDoesNotTrustAsset(client, sourceAccount, destinationPublicKey, asset, newAmountToSend, db)

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
	return xdrBase64, destinationUser, nil
}

func ProcessArrestAccount(client *horizonclient.Client, owner *users.User, sourceAccount *horizon.Account, db *gorm.DB) (xdrBase64 string, goinToPrison bool) {

	prisonSigner := keypair.MustParseFull(os.Getenv("PRISON_SIGNER"))

	//check if account has prison signer key already. checks if account is already arrested.
	if owner.SignerIsValidWA(prisonSigner.Address(), sourceAccount) {
		//account is already in prison
		return "", false
	}

	//check if account is enlisted for prison arrest
	// if !owner.PublicKeyBanned(db, false) {
	// 	return "", false
	// }
	_, _, _, _, channelSourceAccount, _ := network.BlockchainAccountProperties(client, prisonSigner.Address(), txnbuild.NativeAsset{})

	//no errors, so create transaction to in prison
	log.Println("[ProcessArrestAccount]processing arrest XDR")
	arrestRequest := &txnbuild.SetOptions{
		Signer: &txnbuild.Signer{
			Address: prisonSigner.Address(),
			Weight:  2,
		},
		LowThreshold:    txnbuild.NewThreshold(2),
		MediumThreshold: txnbuild.NewThreshold(2),
		HighThreshold:   txnbuild.NewThreshold(2),
		SourceAccount:   sourceAccount.AccountID,
	}
	// sourceAccount.Sequence = (decimal.RequireFromString(sourceAccount.Sequence).Add(decimal.NewFromInt(1))).String()
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        channelSourceAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{arrestRequest},
			BaseFee:              txnbuild.MinBaseFee,
			Preconditions: txnbuild.Preconditions{
				TimeBounds: txnbuild.NewInfiniteTimeout(),
			},
			Memo: txnbuild.MemoText("arrest operation"),
		},
	)
	if err != nil {
		log.Println("[ProcessArrestAccount] error constructing transaction ", err)
		return "", false
	}

	tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), prisonSigner)

	if err != nil {
		log.Println("[ProcessArrestAccount] error signing transaction with prison signer key ", err)
		return "", false
	}

	xdrBase64, err = tx.Base64()

	if err != nil {
		log.Println("[ProcessArrestAccount] error getting txn base64", err)
		return "", false
	}

	return xdrBase64, true

}

func ProcessBailArrestedAccount(client *horizonclient.Client, owner *users.User, sourceAccount *horizon.Account, db *gorm.DB) (txHash string, err error) {
	discord.WebhookURL = "https://discord.com/api/webhooks/824381163367170058/OXSX51RHd9DyLFbFipjdW3yXmyYC8SWwqd6HiXl6UtDzu75RxS1LzWA800hWereJJumw"
	if len(os.Getenv("EXPANSION_NETWORK_ERROR_WEBHOOK")) > 50 {
		discord.WebhookURL = os.Getenv("EXPANSION_NETWORK_ERROR_WEBHOOK")
	}
	prisonSigner := keypair.MustParseFull(os.Getenv("PRISON_SIGNER"))

	//check if account has prison signer key already. checks if account is already arrested.
	if !owner.SignerIsValidWA(prisonSigner.Address(), sourceAccount) {
		//account is not arrested. cannot release account not in prison
		return "", errors.New("account is not arrested")

	}
	// //signer is attached
	// if owner.PublicKeyBanned(db, false) {
	// 	//account is still banned. cannot release account that is banned
	// 	return "", errors.New("account is still banned")
	// }
	// account no longer banned. release from prison
	log.Println("[ProcessBailArrestedAccount]Releasing account from prison", sourceAccount.AccountID)
	_, _, _, _, channelSourceAccount, _ := network.BlockchainAccountProperties(client, prisonSigner.Address(), txnbuild.NativeAsset{})

	freeFromArrestRequest := &txnbuild.SetOptions{
		Signer: &txnbuild.Signer{
			Address: prisonSigner.Address(),
			Weight:  0,
		},
		LowThreshold:    txnbuild.NewThreshold(0),
		MediumThreshold: txnbuild.NewThreshold(0),
		HighThreshold:   txnbuild.NewThreshold(0),
		SourceAccount:   sourceAccount.AccountID,
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        channelSourceAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{freeFromArrestRequest},
			BaseFee:              txnbuild.MinBaseFee,
			Preconditions: txnbuild.Preconditions{
				TimeBounds: txnbuild.NewInfiniteTimeout(),
			},
			Memo: txnbuild.MemoText("free from prison"),
		},
	)
	if err != nil {
		log.Println("[ProcessBailArrestedAccount] error constructing transaction ", err)
		return "", err
	}

	tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), prisonSigner)

	if err != nil {
		log.Println("[ProcessBailArrestedAccount] error signing transaction with prison signer key ", err)
		return "", err
	}
	var xdrBase64 string
	xdrBase64, err = tx.Base64()

	if err != nil {
		log.Println("[ProcessBailArrestedAccount] error getting txn base64", err)
		return "", err
	}

	resp, err := client.SubmitTransactionXDR(xdrBase64)
	if err != nil {
		if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "handshake") || strings.Contains(err.Error(), "read tcp") || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp") || strings.Contains(err.Error(), "no such host") {
			discord.Say(fmt.Sprintf("[ProcessBailArrestedAccount] error connecting to expansion service: %v\nXDR: %v", err, xdrBase64))
		}

		horizonException, ok := err.(*horizonclient.Error)

		if ok {

			extraErrors := horizonException.Problem.Extras

			for key, val := range extraErrors {
				log.Printf("[ProcessBailArrestedAccount] Extras: %v is %v\nOwner publickey: %v\n", key, val, sourceAccount.AccountID)
				logDiscordFailedPayment(fmt.Sprintf("[ProcessBailArrestedAccount] Extras: %v is %v\nOwner publickey: %v\n", key, val, sourceAccount.AccountID))

			}

			resultCodes, _ := horizonException.ResultCodes()

			for key, val := range resultCodes.OperationCodes {
				log.Printf("[ProcessBailArrestedAccount] Result code: %v is %v\nOwner publickey: %v\n", key, val, sourceAccount.AccountID)
				logDiscordFailedPayment(fmt.Sprintf("[ProcessBailArrestedAccount] Result code: %v is %v\nOwner publickey: %v\n", key, val, sourceAccount.AccountID))

			}

		}

		return "", &tErrors.CustomError{
			Param:      "publicKey",
			Err:        "error-failed-to-release-account-from-prison",
			ErrMessage: "Failed to release account from prison",
		}

	}

	return resp.Hash, nil

}

func processDestinationAssetDoesNotTrustAsset(client *horizonclient.Client, sourceAccount *horizon.Account, destinationPublicKey string, asset txnbuild.Asset, amountToSend string, db *gorm.DB) ([]txnbuild.Operation, *keypair.Full, error) {

	ops := make([]txnbuild.Operation, 0)

	var signerKeyPairToReturn *keypair.Full = nil

	tempAccountKeypair, tempAccountError := network.TempAccountKeypair(destinationPublicKey)

	var tempAccount txnbuild.Account = &txnbuild.SimpleAccount{AccountID: tempAccountKeypair.FromAddress().Address(), Sequence: 0}

	if tempAccountError != nil {
		log.Printf("error generating temporary account %v", tempAccountError)
		return ops, tempAccountKeypair, &tErrors.ErrorTemporaryServerError{}
	}

	tempAccountExists, tempAccountTrustsAsset, _, _, _, tempAccountError :=
		network.BlockchainAccountProperties(client, tempAccountKeypair.FromAddress().Address(), asset)

	if tempAccountError != nil {
		log.Printf("error looking up temporary account %v", tempAccountError)
		return ops, tempAccountKeypair, &tErrors.ErrorTemporaryServerError{}
	}

	{

		var _wallets []users.UserWallet

		//check if temp account exists in buds or not.
		err := db.Where("public_key = ?", destinationPublicKey).Find(&_wallets).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ops, nil, &tErrors.ErrorInvalidPublicKey{}

			}
			log.Println("[UserRegistrationDbChecks]", err)
			return ops, nil, &tErrors.ErrorTemporaryServerError{}

		}

		var _walletsToUpdate []users.UserWallet

		for _, _wallet := range _wallets {
			if _wallet.TempPublicKey != nil {
				if *_wallet.TempPublicKey != tempAccountKeypair.FromAddress().Address() {
					v := tempAccountKeypair.FromAddress().Address()
					_wallet.TempPublicKey = &v
					_walletsToUpdate = append(_walletsToUpdate, _wallet)
				}
			} else {
				v := tempAccountKeypair.FromAddress().Address()
				_wallet.TempPublicKey = &v
				_walletsToUpdate = append(_walletsToUpdate, _wallet)
			}

		}

		if len(_walletsToUpdate) > 0 {
			dbSaveError := db.Save(_walletsToUpdate).Error

			if dbSaveError != nil {
				log.Printf("db temp save error %v\n", dbSaveError)
				return ops, nil, &tErrors.ErrorTemporaryServerError{}
			}
		}

	}

	//no errors, so do

	baseReserve := network.GetBlockchainBaseReserve()

	if !tempAccountExists {

		//create temp account and add destination public key as signer.

		ops = append(ops, &txnbuild.CreateAccount{
			Destination:   tempAccountKeypair.FromAddress().Address(),
			Amount:        baseReserve.Mul(decimal.NewFromInt(3)).Truncate(7).String(),
			SourceAccount: sourceAccount.AccountID,
		})

		ops = append(ops, &txnbuild.SetOptions{
			Signer: &txnbuild.Signer{
				Address: destinationPublicKey,
				Weight:  1,
			},
			SourceAccount: tempAccount.GetAccountID(),
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
			SourceAccount: tempAccount.GetAccountID(),
		})

		signerKeyPairToReturn = tempAccountKeypair

	}

	ops = append(ops, &txnbuild.Payment{
		Destination:   tempAccountKeypair.FromAddress().Address(),
		Amount:        amountToSend,
		Asset:         asset,
		SourceAccount: sourceAccount.AccountID,
	})

	return ops, signerKeyPairToReturn, nil

}

//GetBlockchainAccountDataKey fetches the bantu account information using public key
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

func logDiscordFailedPayment(msg string) {
	discord.WebhookURL = "https://discord.com/api/webhooks/827986576415129663/wqMKp9wxB_fxs9Q3zlMKCNPGENXmD_ueUnL8hVCu1wmRfD2wkXAjfP85k1Ro_2_wGfiY"
	if len(os.Getenv("FAILED_PAYMENT_ERROR_WEBHOOK")) > 50 {
		discord.WebhookURL = os.Getenv("FAILED_PAYMENT_ERROR_WEBHOOK")
	}
	discord.Say(msg)
}

package users

import (
	"fmt"
	"log"
	"os"
	"strings"
	userDB "trovo-wallet-api/internal/components/users/db"
	userModels "trovo-wallet-api/internal/components/users/models"

	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/ecnepsnai/discord"
	"github.com/shopspring/decimal"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/txnbuild"
)

func CreateNewSubWallet(user *userModels.User, subWalletInfo *userModels.SubWalletInfo, gc *sharedconfig.GlobalConfig) (*userModels.SubWalletInfo, error) {
	client := network.GetBlockchainClient()
	var err error
	var xdrBase64 string
	var subWalletObj userModels.UserWallet

	subWalletInfo.NetworkPassPhrase = network.GetBlockchainNetworkPassPhrase()

	if len(subWalletInfo.ChannelAccount) == 56 {
		//generate xdr for channel account
		xdrBase64, subWalletObj, err = generateSubWalletXdrWithChannelAccount(user, subWalletInfo, gc, client)
		if err != nil {
			log.Printf("[CreateNewSubWallet] create sub [%v] for [%v] generateSubWalletXdrWithChannelAccount error:[%v] \n", subWalletInfo.PublicKey, user.Username, err)
		}
	} else {
		xdrBase64, subWalletObj, err = generateSubWalletXdr(user, subWalletInfo, gc, client)
		if err != nil {
			log.Printf("[CreateNewSubWallet] create sub [%v] for [%v] generateSubWalletXdr error:[%v] \n", subWalletInfo.PublicKey, user.Username, err)
		}
	}

	oldTransaction := subWalletInfo.Transaction

	subWalletInfo.Transaction = xdrBase64

	if len(subWalletInfo.PrimarySignature) == 0 || len(subWalletInfo.SubWalletSignature) == 0 {
		//it has not been signed before
		return subWalletInfo, nil
	}

	if oldTransaction != subWalletInfo.Transaction {
		err = &tErrors.CustomError{
			Param:      "transaction",
			Err:        "error-transaction-mismatch",
			ErrMessage: "Transaction mismatch",
		}
		return subWalletInfo, err

	}

	if len(subWalletInfo.ChannelAccountSignature) == 0 && len(subWalletInfo.ChannelAccount) == 56 {
		err = &tErrors.CustomError{
			Param:      "transaction",
			Err:        "error-transaction-mismatch",
			ErrMessage: "signature for channel account does not validate",
		}
		return subWalletInfo, err

	}
	//begin a database transaction here
	dbTX := gc.DB.Begin()
	defer dbTX.Rollback()

	//create the data to be sure it goes through
	errDBTX := dbTX.Create(&subWalletObj).Error
	if errDBTX != nil {
		//unable to save sub wallet. abort
		err = &tErrors.CustomError{
			Param:      "publicKey",
			Err:        "error-saving-subwallet",
			ErrMessage: "There is an error saving sub-wallet. Please, try again later.",
		}
		return subWalletInfo, err
	}

	if len(subWalletInfo.ChannelAccountSignature) > 0 && len(subWalletInfo.ChannelAccount) == 56 {
		txnHash, err := SubmitSubWalletXdrForChannelAccountWithSignature(client, user.PublicKey, subWalletInfo.PublicKey, subWalletInfo.ChannelAccount, xdrBase64, subWalletInfo.PrimarySignature, subWalletInfo.SubWalletSignature, subWalletInfo.ChannelAccountSignature)
		if err != nil {
			log.Printf("[CreateNewSubWallet] by [%v] for [%v] SubmitSubWalletXdrForChannelAccountWithSignature error:[%v] \n", user.Username, subWalletInfo.PublicKey, err)
			return subWalletInfo, err
		}
		subWalletInfo.TransactionID = txnHash
		dbTX.Commit()
	} else {
		txnHash, err := SubmitSubWalletXdrWithSignature(client, user.PublicKey, subWalletInfo.PublicKey, xdrBase64, subWalletInfo.PrimarySignature, subWalletInfo.SubWalletSignature)
		if err != nil {
			log.Printf("[CreateNewSubWallet] by [%v] for [%v] SubmitSubwalletXdrWithSignature error:[%v] \n", user.Username, subWalletInfo.PublicKey, err)
			return subWalletInfo, err
		}
		subWalletInfo.TransactionID = txnHash
		dbTX.Commit()

	}

	return subWalletInfo, nil
}

func generateSubWalletXdr(user *userModels.User, subWalletInfo *userModels.SubWalletInfo, gc *sharedconfig.GlobalConfig, client *horizonclient.Client) (xdrbase64 string, subWalletObj userModels.UserWallet, err error) {
	ops := make([]txnbuild.Operation, 0)
	subWalletInfo.Messages = make([]string, 0)
	var activationAmount = decimal.NewFromFloat(3.0)
	var minBalance = decimal.NewFromFloat(3.0)
	if len(os.Getenv("SUB_WALLET_ACTIVATION_AMOUNT")) > 0 {
		activationAmount = decimal.RequireFromString(os.Getenv("SUB_WALLET_ACTIVATION_AMOUNT"))
	}
	if len(os.Getenv("WALLET_MINIMUM_BALANCE")) > 0 {
		minBalance = decimal.RequireFromString(os.Getenv("WALLET_MINIMUM_BALANCE"))
	}
	{
		//check if the sub-wallet passes the validation
		subWalletObj, err = user.BuildNewSubWallet(subWalletInfo.PublicKey, subWalletInfo.WalletTag, subWalletInfo.WalletDescription, gc)
		if err != nil {
			return "", subWalletObj, err
		}

	}

	//check if it is first call to create sub-wallet

	//populate the subwallet Info and generate the transaction

	//check if primary account has native enough native balance
	var nativeAsset txnbuild.Asset = txnbuild.NativeAsset{}
	primaryAccountExists, _, primaryAccountNativeBalance, _, primarySourceAccount, errAct := network.BlockchainAccountProperties(client, user.PublicKey, nativeAsset)
	if errAct != nil {
		return "", subWalletObj, errAct
	}
	if !primaryAccountExists || (primaryAccountNativeBalance.Sub(activationAmount)).LessThan(minBalance) {
		err = &tErrors.CustomError{
			Param:      "publicKey",
			Err:        "error-primary-account-underfunded",
			ErrMessage: "Primary account does not have enough XBN balance to create sub-wallet",
			Code:       404,
		}
		return "", subWalletObj, err
	}

	subWalletAccountExists, _, subWalletAccountNativeBalance, _, subWalletAccountObject, _ := network.BlockchainAccountProperties(client, subWalletInfo.PublicKey, nativeAsset)
	if !subWalletAccountExists {
		//if subwallet is not activated
		//build transaction that will activate the subwallet from the primary wallet

		ops = append(ops, &txnbuild.CreateAccount{
			Destination:   subWalletInfo.PublicKey,
			Amount:        activationAmount.String(),
			SourceAccount: user.PublicKey,
		})

		//after creation, it now exists with enough balance to add primary wallet as signer
		ops = append(ops, &txnbuild.SetOptions{
			Signer: &txnbuild.Signer{
				Address: user.PublicKey,
				Weight:  1,
			},
			SourceAccount: subWalletInfo.PublicKey,
		})

		subWalletInfo.Messages = append(subWalletInfo.Messages, fmt.Sprintf("Important: %v will be deducted from your primary wallet to used to activate the sub-wallet.", activationAmount.String()))

	}

	if subWalletAccountExists && (subWalletAccountNativeBalance.LessThan(minBalance)) {
		//account exists and native balance is less than needed. add 3 native token to the wallet
		ops = append(ops, &txnbuild.Payment{
			Destination:   subWalletInfo.PublicKey,
			Amount:        minBalance.String(),
			Asset:         nativeAsset,
			SourceAccount: user.PublicKey,
		})

		//after topping up, it now has enough balance to add primary wallet as signer if it is not already a signer
		if !user.SignerIsValidWA(user.PublicKey, subWalletAccountObject) {

			ops = append(ops, &txnbuild.SetOptions{
				Signer: &txnbuild.Signer{
					Address: user.PublicKey,
					Weight:  1,
				},
				SourceAccount: subWalletInfo.PublicKey,
			})
		}
		subWalletInfo.Messages = append(subWalletInfo.Messages, fmt.Sprintf("Important: %v will be deducted from your primary wallet to used to complete the sub-wallet process.", activationAmount.String()))

	}
	//TODO: if account exists and subwallet has enough balance, we add the operation to pay TROVO fee from primary Wallet

	// Construct the transaction that holds the operations to execute on the network
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        primarySourceAccount,
			IncrementSequenceNum: true,
			Operations:           ops,
			BaseFee:              txnbuild.MinBaseFee,
			Preconditions: txnbuild.Preconditions{
				TimeBounds: txnbuild.NewInfiniteTimeout(),
			},
			Memo: txnbuild.MemoText("Create Sub Wallet"),
		},
	)
	if err != nil {
		log.Println("[CreateNewSubWallet] error constructing transaction ", err)
		return "", subWalletObj, err
	}

	var xdrBase64 string

	xdrBase64, err = tx.Base64()
	if err != nil {
		log.Println("[CreateNewSubWallet] error getting txn base64", err)
		return "", subWalletObj, err
	}

	return xdrBase64, subWalletObj, nil

}

func generateSubWalletXdrWithChannelAccount(user *userModels.User, subWalletInfo *userModels.SubWalletInfo, gc *sharedconfig.GlobalConfig, client *horizonclient.Client) (xdrbase64 string, subWalletObj userModels.UserWallet, err error) {
	ops := make([]txnbuild.Operation, 0)
	subWalletInfo.Messages = make([]string, 0)
	var activationAmount = decimal.NewFromFloat(3.0)
	var minBalance = decimal.NewFromFloat(3.0)
	if len(os.Getenv("SUB_WALLET_ACTIVATION_AMOUNT")) > 0 {
		activationAmount = decimal.RequireFromString(os.Getenv("SUB_WALLET_ACTIVATION_AMOUNT"))
	}
	if len(os.Getenv("WALLET_MINIMUM_BALANCE")) > 0 {
		minBalance = decimal.RequireFromString(os.Getenv("WALLET_MINIMUM_BALANCE"))
	}
	{
		//check if the sub-wallet passes the validation
		subWalletObj, err = user.BuildNewSubWallet(subWalletInfo.PublicKey, subWalletInfo.WalletTag, subWalletInfo.WalletDescription, gc)
		if err != nil {
			return "", subWalletObj, err
		}

	}

	//check if it is first call to create sub-wallet

	//populate the subwallet Info and generate the transaction

	//check if primary account has native enough native balance
	var nativeAsset txnbuild.Asset = txnbuild.NativeAsset{}
	primaryAccountExists, _, primaryAccountNativeBalance, _, _, errAct := network.BlockchainAccountProperties(client, user.PublicKey, nativeAsset)
	if errAct != nil {
		return "", subWalletObj, errAct
	}
	if !primaryAccountExists || (primaryAccountNativeBalance.Sub(activationAmount)).LessThan(minBalance) {
		err = &tErrors.CustomError{
			Param:      "publicKey",
			Err:        "error-primary-account-underfunded",
			ErrMessage: "Primary account does not have enough XBN balance to create sub-wallet",
			Code:       404,
		}
		return "", subWalletObj, err
	}

	subWalletAccountExists, _, subWalletAccountNativeBalance, _, subWalletAccountObject, _ := network.BlockchainAccountProperties(client, subWalletInfo.PublicKey, nativeAsset)
	if !subWalletAccountExists {
		//if subwallet is not activated
		//build transaction that will activate the subwallet from the primary wallet

		ops = append(ops, &txnbuild.CreateAccount{
			Destination:   subWalletInfo.PublicKey,
			Amount:        activationAmount.String(),
			SourceAccount: user.PublicKey,
		})

		//after creation, it now exists with enough balance to add primary wallet as signer
		ops = append(ops, &txnbuild.SetOptions{
			Signer: &txnbuild.Signer{
				Address: user.PublicKey,
				Weight:  1,
			},
			SourceAccount: subWalletInfo.PublicKey,
		})

		subWalletInfo.Messages = append(subWalletInfo.Messages, fmt.Sprintf("Important: %v will be deducted from your primary wallet to used to activate the sub-wallet.", activationAmount.String()))

	}

	if subWalletAccountExists && (subWalletAccountNativeBalance.LessThan(minBalance)) {
		//account exists and native balance is less than needed. add 3 native token to the wallet
		ops = append(ops, &txnbuild.Payment{
			Destination:   subWalletInfo.PublicKey,
			Amount:        minBalance.String(),
			Asset:         nativeAsset,
			SourceAccount: user.PublicKey,
		})

		//after topping up, it now has enough balance to add primary wallet as signer if it is not already a signer
		if !user.SignerIsValidWA(user.PublicKey, subWalletAccountObject) {

			ops = append(ops, &txnbuild.SetOptions{
				Signer: &txnbuild.Signer{
					Address: user.PublicKey,
					Weight:  1,
				},
				SourceAccount: subWalletInfo.PublicKey,
			})
		}
		subWalletInfo.Messages = append(subWalletInfo.Messages, fmt.Sprintf("Important: %v will be deducted from your primary wallet to used to complete the sub-wallet process.", activationAmount.String()))

	}
	channelSourceAccountExists, _, channelSourceAccountNativeBalance, _, channelSourceAccount, channelSourceAccountErr := network.BlockchainAccountProperties(client, subWalletInfo.ChannelAccount, txnbuild.NativeAsset{})
	if !channelSourceAccountExists || channelSourceAccountErr != nil || (channelSourceAccountNativeBalance.Sub(activationAmount)).LessThan(minBalance) {
		err = &tErrors.CustomError{
			Param:      "channelAccount",
			Err:        "error-channel-account-underfunded",
			ErrMessage: "Channel account does not have minimum XBN balance required to complete this operation",
			Code:       400,
		}
		return "", subWalletObj, err
	}
	//TODO: if account exists and subwallet has enough balance, we add the operation to pay TROVO fee from primary Wallet

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
			Memo: txnbuild.MemoText("Create Sub Wallet"),
		},
	)
	if err != nil {
		log.Println("[CreateNewSubWallet] error constructing transaction ", err)
		return "", subWalletObj, err
	}

	var xdrBase64 string

	xdrBase64, err = tx.Base64()
	if err != nil {
		log.Println("[CreateNewSubWallet] error getting txn base64", err)
		return "", subWalletObj, err
	}

	return xdrBase64, subWalletObj, nil

}

func SubmitSubWalletXdrWithSignature(client *horizonclient.Client, ownerPublicKey, subWalletPublicKey string, xdrBase64 string, primarySignature, subWalletSignature string) (string, error) {
	discord.WebhookURL = "https://discord.com/api/webhooks/824381163367170058/OXSX51RHd9DyLFbFipjdW3yXmyYC8SWwqd6HiXl6UtDzu75RxS1LzWA800hWereJJumw"
	if len(os.Getenv("EXPANSION_NETWORK_ERROR_WEBHOOK")) > 50 {
		discord.WebhookURL = os.Getenv("EXPANSION_NETWORK_ERROR_WEBHOOK")
	}
	gTxn, err := txnbuild.TransactionFromXDR(xdrBase64)

	if err != nil {
		return "", err
	}

	txn, ok := gTxn.Transaction()

	if !ok {
		return "", &tErrors.ErrorInvalidTransaction{}
	}

	{
		//add signature of the primaryWallet to the new transaction Instance
		txn, err = txn.AddSignatureBase64(network.GetBlockchainNetworkPassPhrase(), ownerPublicKey, primarySignature)
		if err != nil {
			log.Println("[SubmitSubwalletXdrWithSignature] Failed to verify primary signature of primaryWallet on [", network.GetBlockchainNetworkPassPhrase(), "] and [", primarySignature, "] for [", xdrBase64, "] and public key ", ownerPublicKey, ", error [", err, "]")

			return "", err
		}

		//add signature of the subWallet to the new transaction Instance
		txn, err = txn.AddSignatureBase64(network.GetBlockchainNetworkPassPhrase(), subWalletPublicKey, subWalletSignature)
		if err != nil {
			log.Println("[SubmitSubwalletXdrWithSignature] Failed to verify signature of subwallet on [", network.GetBlockchainNetworkPassPhrase(), "] and [", subWalletSignature, "] for [", xdrBase64, "] and public key ", subWalletPublicKey, ", error [", err, "]")

			return "", err
		}
	}

	xdrBase64, err = txn.Base64()

	if err != nil {
		log.Printf("[SubmitSubwalletXdrWithSignature] error converting transaction to base64: %v\n", err)
		return "", err
	}

	// log.Println("signed xdr is " + xdrBase64)

	txnResult, err := client.SubmitTransactionXDR(xdrBase64)

	if err != nil {
		if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "handshake") || strings.Contains(err.Error(), "read tcp") || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp") || strings.Contains(err.Error(), "no such host") {
			discord.Say(fmt.Sprintf("[SubmitSubwalletXdrWithSignature] error connecting to expansion service: %v\nXDR: %v", err, xdrBase64))
		}

		horizonException, ok := err.(*horizonclient.Error)

		if ok {

			extraErrors := horizonException.Problem.Extras

			for key, val := range extraErrors {
				log.Printf("[SubmitSubwalletXdrWithSignature] Extras: %v is %v\nOwner publicKey: %v, subwallet: %v\n", key, val, ownerPublicKey, subWalletPublicKey)

			}

			resultCodes, errRes := horizonException.ResultCodes()
			if errRes == nil {
				for key, val := range resultCodes.OperationCodes {
					log.Printf("[SubmitSubwalletXdrWithSignature] Result code: %v is %v\nOwner publicKey: %v\nSubWalletPublicKey: %v\n", key, val, ownerPublicKey, subWalletPublicKey)
					// logDiscordFailedPayment(fmt.Sprintf("[SubmitSubwalletXdrWithSignature]Result code: %v is %v\nOwner publicKey: %v\nSubWalletPublicKey: %v\n", key, val, ownerPublicKey, subWalletPublicKey))

				}
			} else {
				log.Printf("[SubmitSubwalletXdrWithSignature] Error getting result codes: %v\n", errRes)
			}

		} else {
			log.Printf("[SubmitSubwalletXdrWithSignature] not horizon error: %v\n", err)

		}

		return "", &tErrors.CustomError{Param: "publicKey", Err: "error subwallet activation failed", ErrMessage: "SubWallet Failed", Code: 500}

	}

	return txnResult.Hash, nil

}

func SubmitSubWalletXdrForChannelAccountWithSignature(client *horizonclient.Client, ownerPublicKey, subWalletPublicKey, channelPK string, xdrBase64 string, primarySignature, subWalletSignature, channelAccountSignature string) (string, error) {
	discord.WebhookURL = "https://discord.com/api/webhooks/824381163367170058/OXSX51RHd9DyLFbFipjdW3yXmyYC8SWwqd6HiXl6UtDzu75RxS1LzWA800hWereJJumw"
	if len(os.Getenv("EXPANSION_NETWORK_ERROR_WEBHOOK")) > 50 {
		discord.WebhookURL = os.Getenv("EXPANSION_NETWORK_ERROR_WEBHOOK")
	}
	gTxn, err := txnbuild.TransactionFromXDR(xdrBase64)

	if err != nil {
		return "", err
	}

	txn, ok := gTxn.Transaction()

	if !ok {
		return "", &tErrors.ErrorInvalidTransaction{}
	}

	{
		//add signature of the primaryWallet to the new transaction Instance
		txn, err = txn.AddSignatureBase64(network.GetBlockchainNetworkPassPhrase(), ownerPublicKey, primarySignature)
		if err != nil {
			log.Println("[SubmitSubwalletXdrWithSignature] Failed to verify primary signature of primaryWallet on [", network.GetBlockchainNetworkPassPhrase(), "] and [", primarySignature, "] for [", xdrBase64, "] and public key ", ownerPublicKey, ", error [", err, "]")

			return "", err
		}

		//add signature of the subWallet to the new transaction Instance
		txn, err = txn.AddSignatureBase64(network.GetBlockchainNetworkPassPhrase(), subWalletPublicKey, subWalletSignature)
		if err != nil {
			log.Println("[SubmitSubwalletXdrWithSignature] Failed to verify signature of subwallet on [", network.GetBlockchainNetworkPassPhrase(), "] and [", subWalletSignature, "] for [", xdrBase64, "] and public key ", subWalletPublicKey, ", error [", err, "]")

			return "", err
		}
		//add signature of the channelAccount to the new transaction Instance
		txn, err = txn.AddSignatureBase64(network.GetBlockchainNetworkPassPhrase(), subWalletPublicKey, channelAccountSignature)
		if err != nil {
			log.Println("[SubmitSubwalletXdrWithSignature] Failed to verify signature of channelAccount on [", network.GetBlockchainNetworkPassPhrase(), "] and [", channelAccountSignature, "] for [", xdrBase64, "] and public key ", channelPK, ", error [", err, "]")

			return "", err
		}
	}

	xdrBase64, err = txn.Base64()

	if err != nil {
		log.Printf("[SubmitSubwalletXdrWithSignature] error converting transaction to base64: %v\n", err)
		return "", err
	}

	// log.Println("signed xdr is " + xdrBase64)

	txnResult, err := client.SubmitTransactionXDR(xdrBase64)

	if err != nil {
		if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "handshake") || strings.Contains(err.Error(), "read tcp") || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp") || strings.Contains(err.Error(), "no such host") {
			discord.Say(fmt.Sprintf("[SubmitSubwalletXdrWithSignature] error connecting to expansion service: %v\nXDR: %v", err, xdrBase64))
		}

		horizonException, ok := err.(*horizonclient.Error)

		if ok {

			extraErrors := horizonException.Problem.Extras

			for key, val := range extraErrors {
				log.Printf("[SubmitSubwalletXdrWithSignature] Extras: %v is %v\nOwner publicKey: %v, subwallet: %v\n", key, val, ownerPublicKey, subWalletPublicKey)

			}

			resultCodes, errRes := horizonException.ResultCodes()
			if errRes == nil {
				for key, val := range resultCodes.OperationCodes {
					log.Printf("[SubmitSubwalletXdrWithSignature] Result code: %v is %v\nOwner publicKey: %v\nSubWalletPublicKey: %v\n", key, val, ownerPublicKey, subWalletPublicKey)
					// logDiscordFailedPayment(fmt.Sprintf("[SubmitSubwalletXdrWithSignature]Result code: %v is %v\nOwner publicKey: %v\nSubWalletPublicKey: %v\n", key, val, ownerPublicKey, subWalletPublicKey))

				}
			} else {
				log.Printf("[SubmitSubwalletXdrWithSignature] Error getting result codes: %v\n", errRes)
			}

		} else {
			log.Printf("[SubmitSubwalletXdrWithSignature] not horizon error: %v\n", err)

		}

		return "", &tErrors.CustomError{Param: "publicKey", Err: "error subwallet activation failed", ErrMessage: "SubWallet Failed", Code: 500}

	}

	return txnResult.Hash, nil

}

func HasAccessToPublicKey(ownerPublicKey, targetPublicKey string, gc *sharedconfig.GlobalConfig) (hasAccess bool) {
	user, err := userDB.GetUser(ownerPublicKey, gc.DB)

	if err != nil {
		return false
	}
	walletPermissions := user.Fetch3rdPartyWallets(gc)
	if len(walletPermissions) == 0 {
		return false
	}
	for _, walletAccess := range walletPermissions {
		if walletAccess.PublicKey == targetPublicKey {
			return true
		}
	}

	return false
}

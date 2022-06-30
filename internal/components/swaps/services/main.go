package swaps

import (
	"fmt"
	"os"
	"strings"
	swapErrors "trovo-wallet-api/internal/components/swaps/errors"
	swapModels "trovo-wallet-api/internal/components/swaps/models"
	userModels "trovo-wallet-api/internal/components/users/models"
	"trovo-wallet-api/internal/sharedconfig"

	"log"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/network"

	"github.com/ecnepsnai/discord"
	"github.com/shopspring/decimal"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
)

//SwapSend function swaps an asset to another asset
func SwapSend(walletOwner *userModels.User, wallet *userModels.UserWallet, signerPublicKey string, swapInfo *swapModels.SwapSendInfo, gc *sharedconfig.GlobalConfig) error {

	client := gc.BantuExpansionClient
	//transform codes and issuer
	swapInfo.DestinationAssetCode = strings.ToUpper(swapInfo.DestinationAssetCode)
	swapInfo.DestinationAssetIssuer = strings.ToUpper(swapInfo.DestinationAssetIssuer)
	swapInfo.SourceAssetCode = strings.ToUpper(swapInfo.SourceAssetCode)
	swapInfo.SourceAssetIssuer = strings.ToUpper(swapInfo.SourceAssetIssuer)
	if e := ValidateSwapSendInfo(swapInfo); e != nil {
		return e
	}
	xdrBase64, err := generateSwapXdr(signerPublicKey, walletOwner, wallet, swapInfo, gc)

	// oldTransaction := swapInfo.Transaction

	swapInfo.Transaction = xdrBase64
	swapInfo.NetworkPassPhrase = network.GetBlockchainNetworkPassPhrase()

	if len(swapInfo.TransactionSignature) == 0 {
		return err
	}

	//no need to check this since offer can change, therefore changing the transaction
	// if xdrBase64 != oldTransaction {
	// 	return &swapErrors.ErrorTransactionMismatch{}
	// }

	txnHash, err := network.SubmitXdrWithSignature(client, signerPublicKey, xdrBase64, swapInfo.TransactionSignature)
	if err != nil {
		logDiscordFailedSwap(fmt.Sprintf("Error submitting swap [%+v] transaction: %s", swapInfo, err.Error()))
	}
	swapInfo.TransactionID = txnHash
	return err
}

func generateSwapXdr(signerPublicKey string, owner *userModels.User, wallet *userModels.UserWallet, swapInfo *swapModels.SwapSendInfo, gc *sharedconfig.GlobalConfig) (string, error) {
	baseReserve := network.GetBlockchainBaseReserve()
	swapDestMin := network.GetBlockchainSwapDestinationMin()
	client := gc.BantuExpansionClient
	messages := make([]string, 0)

	var err error
	var amountToSwap decimal.Decimal

	if amountToSwap, err = decimal.NewFromString(swapInfo.SourceAmount); err != nil {
		return "", &swapErrors.ErrorInvalidSwapAmount{}
	}

	newAmountToSwap := amountToSwap.Truncate(7).String()

	var sourceAsset, destinationAsset txnbuild.Asset = txnbuild.NativeAsset{}, txnbuild.NativeAsset{}

	if len(swapInfo.DestinationAssetCode) != 0 {

		destinationAsset = txnbuild.CreditAsset{Code: swapInfo.DestinationAssetCode, Issuer: swapInfo.DestinationAssetIssuer}
	}
	if len(swapInfo.SourceAssetCode) != 0 {

		sourceAsset = txnbuild.CreditAsset{Code: swapInfo.SourceAssetCode, Issuer: swapInfo.SourceAssetIssuer}
	}

	// charge := baseReserve.Mul(decimal.NewFromInt(1)).Truncate(7).String()
	appliedCharge := decimal.NewFromFloat(0)
	swapInfo.Messages = messages
	var ops []txnbuild.Operation = make([]txnbuild.Operation, 0)

	sourceAccountExists, _, sourceAccountNativeBalance, sourceAccountCustomBalance, sourceAccount, sourceAccountErr := network.BlockchainAccountProperties(client, wallet.ID, sourceAsset)
	var sourceAccountTrustsDestinationAsset bool
	if !destinationAsset.IsNative() {
		_, sourceAccountTrustsDestinationAsset, _, _, _, _ = network.BlockchainAccountProperties(client, wallet.ID, destinationAsset)

	}

	if sourceAccountErr != nil {
		return "", sourceAccountErr
	}

	if !sourceAccountExists {
		return "", &tErrors.ErrorUnderfundedAccount{}
	}

	if !destinationAsset.IsNative() {

		if !sourceAccountTrustsDestinationAsset {
			message := fmt.Sprintf("Important: you have not yet activated the asset (%v) you are trying to swap to. Some XBN will be deducted from the wallet [%v] to activate the asset.", swapInfo.DestinationAssetCode, wallet.Alias)
			messages = append(messages, message)
			log.Printf("message[0]: %v\n", message)
			appliedCharge = baseReserve.Mul(decimal.NewFromInt(1)).Truncate(7)
			//establish trustline
			ops = append(ops, &txnbuild.ChangeTrust{
				Line:          txnbuild.ChangeTrustAssetWrapper{Asset: destinationAsset},
				Limit:         "900000000000",
				SourceAccount: wallet.ID,
			})

		}
	}

	log.Printf("obtained sourced account info \n")

	log.Printf("obtained source account balance is XBN %v, custom account balance %v\n", sourceAccountNativeBalance, sourceAccountCustomBalance)

	amountToSwapDec := amountToSwap

	if sourceAsset.IsNative() {
		if sourceAccountNativeBalance.LessThan(amountToSwapDec.Add(appliedCharge)) {
			return "", &tErrors.ErrorUnderfundedAccount{}
		}
	} else {
		if sourceAccountCustomBalance.LessThan(amountToSwapDec) {
			return "", &tErrors.ErrorUnderfundedAccount{}
		}
		//if it is custom asset check if the native balance can carry the charge for trusting asset
		if sourceAccountNativeBalance.LessThan(appliedCharge) {
			return "", &tErrors.ErrorUnderfundedAccount{}
		}
	}

	//get sendPath
	//using DestinationAccount will get paths to all assets in the destination account.
	//using destinationAssets gets path to only the asset
	destAsset := ""
	if !destinationAsset.IsNative() {
		destAsset = fmt.Sprintf("%s:%s", swapInfo.DestinationAssetCode, swapInfo.DestinationAssetIssuer)
	}
	pathInput := swapModels.SwapSendPathInput{
		DestinationAssets: destAsset,
		SourceAssetCode:   swapInfo.SourceAssetCode,
		SourceAssetIssuer: swapInfo.SourceAssetIssuer,
		SourceAmount:      newAmountToSwap,
	}
	path, swappedEstimate, err := getStrictSendPaths(pathInput, client)
	if err != nil {
		log.Println("[generateSwapXdr]error fetching valid swap Path ", err)
		return "", err
	}

	//native asset
	ops = append(ops, &txnbuild.PathPaymentStrictSend{
		SendAsset:     sourceAsset,
		SendAmount:    swapInfo.SourceAmount,
		Destination:   wallet.ID,
		DestAsset:     destinationAsset,
		DestMin:       swapDestMin.String(),
		Path:          path,
		SourceAccount: wallet.ID,
	})

	// Construct the transaction that holds the operations to execute on the network
	var memoSAC, memoDAC string
	memoSAC = swapInfo.SourceAssetCode
	memoDAC = swapInfo.DestinationAssetCode
	if swapInfo.SourceAssetIssuer == "" || swapInfo.SourceAssetIssuer == "native" {
		memoSAC = "XBN"
	}
	if swapInfo.DestinationAssetIssuer == "" || swapInfo.DestinationAssetIssuer == "native" {
		memoDAC = "XBN"
	}

	memo := fmt.Sprintf("%v>%v", memoSAC, memoDAC)
	log.Println("[generateSwapXdr] Memo:", memo)

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        sourceAccount,
			IncrementSequenceNum: true,
			Operations:           ops,
			BaseFee:              2000,
			Preconditions:        txnbuild.Preconditions{TimeBounds: txnbuild.NewInfiniteTimeout()},
			Memo:                 txnbuild.MemoText(memo),
		},
	)

	if err != nil {
		log.Println("[generateSwapXdr]error constructing transaction ", err)
		return "", err
	}

	var xdrBase64 string

	xdrBase64, err = tx.Base64()
	if err != nil {
		log.Println("[generateSwapXdr] error getting txn base64", err)
		return "", err
	}
	swapInfo.Memo = memo
	swapInfo.Messages = messages
	swapInfo.SwappedEstimate = swappedEstimate
	return xdrBase64, nil
}

//getStrictSendPaths gets Strict Send Paths for Strict Send Path Payment request
func getStrictSendPaths(pathInput swapModels.SwapSendPathInput, client *horizonclient.Client) (paths []txnbuild.Asset, swappedEstimate string, err error) {
	discord.WebhookURL = "https://discord.com/api/webhooks/824381163367170058/OXSX51RHd9DyLFbFipjdW3yXmyYC8SWwqd6HiXl6UtDzu75RxS1LzWA800hWereJJumw"
	if len(os.Getenv("EXPANSION_NETWORK_ERROR_WEBHOOK")) > 50 {
		discord.WebhookURL = os.Getenv("EXPANSION_NETWORK_ERROR_WEBHOOK")
	}
	var swapPaths horizon.PathsPage
	paths = make([]txnbuild.Asset, 0)
	var sourceAssetType horizonclient.AssetType
	if len(pathInput.SourceAssetCode) == 0 {
		sourceAssetType = horizonclient.AssetTypeNative
		pathInput.SourceAssetCode = ""
		pathInput.SourceAssetIssuer = ""
	} else if len(pathInput.SourceAssetCode) < 5 {
		sourceAssetType = horizonclient.AssetType4
	} else {
		sourceAssetType = horizonclient.AssetType12
	}
	if pathInput.DestinationAccount != "" {
		pathInput.DestinationAssets = ""
	}

	if pathInput.DestinationAssets == "" && pathInput.DestinationAccount == "" {
		pathInput.DestinationAssets = "native"
	}
	if pathInput.DestinationAssets == "native" {
		pathInput.DestinationAccount = ""
	}

	if pathInput.DestinationAssets != "" {
		pathInput.DestinationAccount = ""
	}
	sspr := horizonclient.StrictSendPathsRequest{
		DestinationAccount: pathInput.DestinationAccount,
		DestinationAssets:  pathInput.DestinationAssets,
		SourceAssetType:    sourceAssetType,
		SourceAssetCode:    pathInput.SourceAssetCode,
		SourceAssetIssuer:  pathInput.SourceAssetIssuer,
		SourceAmount:       pathInput.SourceAmount,
	}

	// tryCount := 0
	// for err != nil && tryCount <= 5 {
	swapPaths, err = client.StrictSendPaths(sspr)
	// log.Printf("[client.StrictSendPathsErr] failed at count: %v, err: %v\n", tryCount, err)
	// tryCount++
	// }
	//successful or 5 retries exceeded

	if err != nil {
		if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "handshake") || strings.Contains(err.Error(), "read tcp") || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp") || strings.Contains(err.Error(), "no such host") {
			log.Println("#######################@@@@@@@@@@@@@@@@[client.StrictSendPathsErr] expansion connection problem:", err)
			discord.Say(fmt.Sprintf("[getStrictSendPaths] error connecting to expansion service: %v\nSwapPathRequest: %+v", err, sspr))

			return paths, "", &tErrors.ErrorTemporaryServerError{}
		}
		horizonException, ok := err.(*horizonclient.Error)

		if ok {

			extraErrors := horizonException.Problem.Extras

			for key, val := range extraErrors {
				log.Printf("Extras: %v is %v\n", key, val)
			}

			resultCodes, _ := horizonException.ResultCodes()

			for key, val := range resultCodes.OperationCodes {
				log.Printf("Result code: %v is %v\n", key, val)
			}
			discord.Say(fmt.Sprintf("[getStrictSendPaths] error submitting: %v\nSwapPathRequest: %+v\nResultCodes: %+v", err, sspr, resultCodes))

		}
		log.Println("[client.StrictSendPathsErr] Error submitting:", err)

		return paths, "", &tErrors.ErrorTemporaryServerError{}

	}
	// discord.Say(fmt.Sprintf("[getStrictSendPaths] swapPaths: %+v\nRequestParams: %+v", swapPaths, sspr))
	// log.Printf("[getStrictSendPaths] swapPaths: %+v\n", swapPaths)
	if len(swapPaths.Embedded.Records) == 0 {
		return paths, "", &swapErrors.ErrorSwapOfferNotAvailable{}
	}
	destAmountDec, _ := decimal.NewFromString(swapPaths.Embedded.Records[0].DestinationAmount)
	if destAmountDec.LessThan(network.GetBlockchainSwapDestinationMin()) {
		return paths, "", &swapErrors.ErrorSwapAmountTooSmall{}
	}

	bestPath := swapPaths.Embedded.Records[0]
	//build assets
	swappedEstimate = bestPath.DestinationAmount

	for _, v := range bestPath.Path {
		if len(v.Issuer) == 0 {
			paths = append(paths, txnbuild.NativeAsset{})
		} else {
			paths = append(paths, txnbuild.CreditAsset{Code: v.Code, Issuer: v.Issuer})
		}

	}

	return paths, swappedEstimate, nil
}

func logDiscordFailedSwap(msg string) {
	discord.WebhookURL = "https://discord.com/api/webhooks/827986576415129663/wqMKp9wxB_fxs9Q3zlMKCNPGENXmD_ueUnL8hVCu1wmRfD2wkXAjfP85k1Ro_2_wGfiY"
	if len(os.Getenv("FAILED_PAYMENT_ERROR_WEBHOOK")) > 50 {
		discord.WebhookURL = os.Getenv("FAILED_PAYMENT_ERROR_WEBHOOK")
	}
	discord.Say(msg)
}

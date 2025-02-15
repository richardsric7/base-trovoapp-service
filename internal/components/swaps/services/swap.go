package swaps

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	// algofuncs "trovo-wallet-api/internal/blockchainalgofuncs"
	swapErrors "trovo-wallet-api/internal/components/swaps/errors"
	swapModels "trovo-wallet-api/internal/components/swaps/models"
	userModels "trovo-wallet-api/internal/components/users/models"
	"trovo-wallet-api/internal/sharedconfig"

	"log"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/network"

	"github.com/ecnepsnai/discord"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
)

// SwapSend function swaps an asset to another asset
func SwapSend(signerUser, walletOwner *userModels.User, wallet *userModels.UserWallet, swapInfo *swapModels.SwapSendInfo, gc *sharedconfig.GlobalConfig) error {
	swapInfo.Messages = make([]string, 0)
	if wallet.SharedAccessEnabled == 1 && wallet.NumberOfApprovalsNeeded > 0 {
		swapInfo.Multiparty = 1
	}
	fee := decimal.RequireFromString(wallet.GetSwapFee(gc))

	feeAmount := ((fee.Mul(decimal.RequireFromString(swapInfo.SourceAmount))).Div(decimal.NewFromInt(100))).Truncate(7)

	swapAmount := decimal.RequireFromString(swapInfo.SourceAmount).Sub(feeAmount)
	swapInfo.SwapAmount = swapAmount.String()
	swapInfo.Fee = fee.String()
	swapInfo.FeeAmount = feeAmount.String()

	if wallet.HasViewOnlyAccess(gc) {
		swapInfo.SignatureRequired = 1
	}
	client := gc.BantuExpansionClient
	//transform codes and issuer
	swapInfo.DestinationAssetCode = strings.ToUpper(swapInfo.DestinationAssetCode)
	swapInfo.DestinationAssetIssuer = strings.ToUpper(swapInfo.DestinationAssetIssuer)
	swapInfo.SourceAssetCode = strings.ToUpper(swapInfo.SourceAssetCode)
	swapInfo.SourceAssetIssuer = strings.ToUpper(swapInfo.SourceAssetIssuer)
	if e := ValidateSwapSendInfo(swapInfo); e != nil {
		return e
	}
	if gc.IsValidTokenizedAsset(swapInfo.DestinationAssetCode) {
		//check if user has done KYC
		if walletOwner.KYCVerified == 0 {
			return &tErrors.CustomError{
				Param:      "destinationAssetCode",
				Err:        "error-no-kyc",
				ErrMessage: fmt.Sprintf("%v does not meet KYC requirement to receive the asset %v", walletOwner.Username, swapInfo.DestinationAssetCode),
			}
		}

		//Check if it is still in primary sales
		if gc.IsTokenizedAssetInPrimarySales(swapInfo.DestinationAssetCode) {
			return &tErrors.CustomError{
				Param:      "destinationAssetCode",
				Err:        "error-primary-sales-active",
				ErrMessage: fmt.Sprintf("%v is still in primary sales. Please go to the tokenized asset market place to purchase from there.", swapInfo.DestinationAssetCode),
			}
		}

	}
	if len(swapInfo.TransactionSignature) == 0 {
		xdrBase64, err := generateSwapSendXdr(wallet, swapInfo, gc)
		if err != nil {
			return err
		}

		swapInfo.Transaction = xdrBase64
		// swapInfo.SHash = algofuncs.SHash(swapInfo.Transaction)
	}

	swapInfo.NetworkPassPhrase = network.GetBlockchainNetworkPassPhrase()

	if len(swapInfo.TransactionSignature) == 0 && swapInfo.Commit == 0 {
		return nil
	}
	//no need to check this since offer can change, therefore changing the transaction

	if len(swapInfo.TransactionSignature) > 0 && swapInfo.Commit == 0 {
		txnHash, err := network.SubmitXdrWithSignature(client, signerUser.PrimarySigner, swapInfo.Transaction, swapInfo.TransactionSignature)
		if err != nil {
			logDiscordFailedSwap(fmt.Sprintf("Error submitting swap [%+v] transaction: %s", swapInfo, err.Error()))
			if strings.Contains(err.Error(), "liquid") {
				destAsset := os.Getenv("NATIVE_ASSET_CODE")
				sourceAsset := os.Getenv("NATIVE_ASSET_CODE")
				if len(swapInfo.SourceAssetCode) > 0 {
					sourceAsset = swapInfo.SourceAssetCode
				}
				if len(swapInfo.DestinationAssetCode) > 0 {
					destAsset = swapInfo.DestinationAssetCode
				}
				return &tErrors.CustomError{
					Param:      "destinationAssetCode",
					Err:        "error-low-liquidity",
					ErrMessage: fmt.Sprintf("There is not enough %v market to exchange for your %v at this time. Please try again later or reduce the quantity of %v to try again.", destAsset, sourceAsset, sourceAsset),
				}
			}
		}
		swapInfo.TransactionID = txnHash
		wallet.InvalidateUserCache(gc)
		return err

	}
	//multi Party
	if swapInfo.Multiparty == 1 {
		swapInfo.TransactionID = "PENDING_AUTH"

		id := uuid.NewString()
		destinationAsset := os.Getenv("NATIVE_ASSET_CODE")
		sourceAsset := os.Getenv("NATIVE_ASSET_CODE")
		if len(swapInfo.SourceAssetIssuer) == 56 {
			sourceAsset = fmt.Sprintf("%v:%v...%v", swapInfo.SourceAssetCode, swapInfo.SourceAssetIssuer[0:4], swapInfo.SourceAssetIssuer[51:55])
		}
		if len(swapInfo.DestinationAssetIssuer) == 56 {
			destinationAsset = fmt.Sprintf("%v:%v...%v", swapInfo.DestinationAssetCode, swapInfo.DestinationAssetIssuer[0:4], swapInfo.DestinationAssetIssuer[51:55])
		}
		description := fmt.Sprintf("Swap\n From:%v,\n To:%v,\n Est. Value After: %v", sourceAsset, destinationAsset, swapInfo.SwappedEstimate)
		if len(swapInfo.Memo) > 0 {
			description = fmt.Sprintf("%v\nMemo: %v", description, swapInfo.Memo)

		}
		if len(swapInfo.Messages) > 0 {
			var msgs string
			for i, m := range swapInfo.Messages {
				msgs = m
				if i < len(swapInfo.Messages)-1 {
					msgs = fmt.Sprintf("%s\n", msgs)
				}
			}
			description = fmt.Sprintf("%v\nMessages: %v", description, msgs)

		}
		swapInfo.ReturnedDescription = description
		transactionByte, _ := json.Marshal(*swapInfo)
		transactionStr := string(transactionByte)
		pendingAuth := userModels.PendingAuth{
			ID:                       id,
			Initiator:                signerUser.Username,
			InitiatorSignerPublicKey: signerUser.PrimarySigner,
			WalletPublicKey:          wallet.ID,
			TransactionType:          "SWAP",
			Description:              description,
			TransactionSource:        swapInfo.TransactionSource,
			ApprovalsNeeded:          wallet.NumberOfApprovalsNeeded,
			TransactionXdr:           swapInfo.Transaction,
			TransactionInfoStr:       &transactionStr,
		}
		//save and commit this to database
		e := gc.DB.Create(&pendingAuth).Error
		if e != nil {
			log.Printf("[SwapSend] Error saving swap txn [%+v] transaction on pending auth table: %s\n", pendingAuth, e.Error())
			err := &tErrors.ErrorTemporaryServerError{}
			return err
		}
		return nil
	}
	log.Println("[SwapSend]UNKNOWN OPTION FOR ACTION")
	err := &tErrors.ErrorTemporaryServerError{}
	return err
}

// SwapReceive function swaps an asset to specic amount of another asset
func SwapReceive(signerUser, walletOwner *userModels.User, wallet *userModels.UserWallet, swapInfo *swapModels.SwapReceiveInfo, gc *sharedconfig.GlobalConfig) error {
	swapInfo.Messages = make([]string, 0)
	if wallet.SharedAccessEnabled == 1 && wallet.NumberOfApprovalsNeeded > 0 {
		swapInfo.Multiparty = 1
	}
	fee := decimal.RequireFromString(wallet.GetSwapFee(gc))

	feeAmount := ((fee.Mul(decimal.RequireFromString(swapInfo.DestinationAmount))).Div(decimal.NewFromInt(100))).Truncate(7)

	swapAmount := decimal.RequireFromString(swapInfo.DestinationAmount).Sub(feeAmount)
	swapInfo.SwapAmount = swapAmount.String()
	swapInfo.Fee = fee.String()
	swapInfo.FeeAmount = feeAmount.String()

	if wallet.HasViewOnlyAccess(gc) {
		swapInfo.SignatureRequired = 1
	}
	client := gc.BantuExpansionClient
	//transform codes and issuer
	swapInfo.DestinationAssetCode = strings.ToUpper(swapInfo.DestinationAssetCode)
	swapInfo.DestinationAssetIssuer = strings.ToUpper(swapInfo.DestinationAssetIssuer)
	swapInfo.SourceAssetCode = strings.ToUpper(swapInfo.SourceAssetCode)
	swapInfo.SourceAssetIssuer = strings.ToUpper(swapInfo.SourceAssetIssuer)
	if e := ValidateSwapReceiveInfo(swapInfo); e != nil {
		return e
	}

	if len(swapInfo.TransactionSignature) == 0 {
		xdrBase64, _, err := generateSwapReceiveXdr(wallet, swapInfo, gc)
		if err != nil {
			return err
		}

		swapInfo.Transaction = xdrBase64
		// swapInfo.SHash = algofuncs.SHash(swapInfo.Transaction)
	}

	swapInfo.NetworkPassPhrase = network.GetBlockchainNetworkPassPhrase()

	if len(swapInfo.TransactionSignature) == 0 && swapInfo.Commit == 0 {
		return nil
	}
	//no need to check this since offer can change, therefore changing the transaction
	// if swapInfo.SHash != algofuncs.SHash(swapInfo.Transaction) {
	// 	return &swapErrors.ErrorTransactionMismatch{}
	// }

	if len(swapInfo.TransactionSignature) > 0 && swapInfo.Commit == 0 {
		txnHash, err := network.SubmitXdrWithSignature(client, signerUser.PrimarySigner, swapInfo.Transaction, swapInfo.TransactionSignature)
		if err != nil {
			logDiscordFailedSwap(fmt.Sprintf("Error submitting swap [%+v] transaction: %s", swapInfo, err.Error()))
			if strings.Contains(err.Error(), "liquid") {
				destAsset := os.Getenv("NATIVE_ASSET_CODE")
				sourceAsset := os.Getenv("NATIVE_ASSET_CODE")
				if len(swapInfo.SourceAssetCode) > 0 {
					sourceAsset = swapInfo.SourceAssetCode
				}
				if len(swapInfo.DestinationAssetCode) > 0 {
					destAsset = swapInfo.DestinationAssetCode
				}
				return &tErrors.CustomError{
					Param:      "destinationAssetCode",
					Err:        "error-low-liquidity",
					ErrMessage: fmt.Sprintf("There is not enough %v market to exchange for your %v at this time. Please try again later or reduce the quantity of %v to try again.", destAsset, sourceAsset, sourceAsset),
				}
			}
		}
		swapInfo.TransactionID = txnHash
		wallet.InvalidateUserCache(gc)
		return err

	}
	//multi Party
	if swapInfo.Multiparty == 1 {
		swapInfo.TransactionID = "PENDING_AUTH"

		id := uuid.NewString()
		destinationAsset := os.Getenv("NATIVE_ASSET_CODE")
		sourceAsset := os.Getenv("NATIVE_ASSET_CODE")
		if len(swapInfo.SourceAssetIssuer) == 56 {
			sourceAsset = fmt.Sprintf("%v:%v...%v", swapInfo.SourceAssetCode, swapInfo.SourceAssetIssuer[0:4], swapInfo.SourceAssetIssuer[51:55])
		}
		if len(swapInfo.DestinationAssetIssuer) == 56 {
			destinationAsset = fmt.Sprintf("%v:%v...%v", swapInfo.DestinationAssetCode, swapInfo.DestinationAssetIssuer[0:4], swapInfo.DestinationAssetIssuer[51:55])
		}
		description := fmt.Sprintf("Swap\n From:%v,\n To:%v,\n Est. Value After: %v", sourceAsset, destinationAsset, swapInfo.RequiredEstimate)
		if len(swapInfo.Memo) > 0 {
			description = fmt.Sprintf("%v\nMemo: %v", description, swapInfo.Memo)

		}
		if len(swapInfo.Messages) > 0 {
			var msgs string
			for i, m := range swapInfo.Messages {
				msgs = m
				if i < len(swapInfo.Messages)-1 {
					msgs = fmt.Sprintf("%s\n", msgs)
				}
			}
			description = fmt.Sprintf("%v\nMessages: %v", description, msgs)

		}
		swapInfo.ReturnedDescription = description
		transactionByte, _ := json.Marshal(*swapInfo)
		transactionStr := string(transactionByte)
		pendingAuth := userModels.PendingAuth{
			ID:                       id,
			Initiator:                signerUser.Username,
			InitiatorSignerPublicKey: signerUser.PrimarySigner,
			WalletPublicKey:          wallet.ID,
			TransactionType:          "SWAP",
			Description:              description,
			TransactionSource:        swapInfo.TransactionSource,
			ApprovalsNeeded:          wallet.NumberOfApprovalsNeeded,
			TransactionXdr:           swapInfo.Transaction,
			TransactionInfoStr:       &transactionStr,
		}
		//save and commit this to database
		e := gc.DB.Create(&pendingAuth).Error
		if e != nil {
			log.Printf("[SwapSend] Error saving swap txn [%+v] transaction on pending auth table: %s\n", pendingAuth, e.Error())
			err := &tErrors.ErrorTemporaryServerError{}
			return err
		}
		return nil
	}
	log.Println("[SwapSend]UNKNOWN OPTION FOR ACTION")
	err := &tErrors.ErrorTemporaryServerError{}
	return err
}

func generateSwapSendXdr(wallet *userModels.UserWallet, swapInfo *swapModels.SwapSendInfo, gc *sharedconfig.GlobalConfig) (string, error) {
	baseReserve := network.GetBlockchainBaseReserve()
	// charge := baseReserve.Mul(decimal.NewFromInt(3)).Truncate(7).String()
	swapDestMin := network.GetBlockchainSwapDestinationMin()
	client := gc.BantuExpansionClient
	messages := make([]string, 0)
	nativeAssetCode := os.Getenv("NATIVE_ASSET_CODE")
	var err error
	var amountToSwap, totalFees decimal.Decimal

	if amountToSwap, err = decimal.NewFromString(swapInfo.SourceAmount); err != nil {
		return "", &swapErrors.ErrorInvalidSwapAmount{}
	}

	newAmountToSwap := amountToSwap.Truncate(7).String()

	var sourceAsset txnbuild.Asset = txnbuild.NativeAsset{}
	var destinationAsset txnbuild.Asset = txnbuild.NativeAsset{}

	if len(swapInfo.DestinationAssetCode) != 0 && !strings.EqualFold(swapInfo.DestinationAssetCode, nativeAssetCode) {

		destinationAsset = txnbuild.CreditAsset{Code: swapInfo.DestinationAssetCode, Issuer: swapInfo.DestinationAssetIssuer}
	}
	if len(swapInfo.SourceAssetCode) != 0 && !strings.EqualFold(swapInfo.SourceAssetCode, nativeAssetCode) {

		sourceAsset = txnbuild.CreditAsset{Code: swapInfo.SourceAssetCode, Issuer: swapInfo.SourceAssetIssuer}
	}

	// charge := baseReserve.Mul(decimal.NewFromInt(1)).Truncate(7).String()
	appliedCharge := decimal.NewFromFloat(0)
	swapInfo.Messages = messages
	var ops []txnbuild.Operation = make([]txnbuild.Operation, 0)
	chanAccount := <-gc.ChannelAccounts
	defer func(c *keypair.Full) {
		gc.ChannelAccounts <- c
	}(chanAccount)

	_, _, _, _, chanSourceAccount, _ := network.BlockchainAccountProperties(client, chanAccount.Address(), txnbuild.NativeAsset{})

	sourceAccountExists, _, sourceAccountNativeBalance, sourceAccountCustomBalance, _, sourceAccountErr := network.BlockchainAccountProperties(client, wallet.ID, sourceAsset)
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
			appliedCharge = baseReserve.Mul(decimal.NewFromInt(2)).Truncate(7)
			message := fmt.Sprintf("%v not yet accepted on [%v]. Continuing will activate %v on [%v].", swapInfo.DestinationAssetCode, wallet.Alias, swapInfo.DestinationAssetCode, wallet.Alias)
			messages = append(messages, message)
			log.Printf("message[0]: %v\n", message)

			// appliedCharge = baseReserve.Mul(decimal.RequireFromString(charge)).Truncate(7)
			totalFees = appliedCharge
			log.Println("[generateSwapXdr] total fees:", totalFees)
			//establish trustline
			ops = append(ops, &txnbuild.ChangeTrust{
				Line:          txnbuild.ChangeTrustAssetWrapper{Asset: destinationAsset},
				Limit:         "900000000000",
				SourceAccount: wallet.ID,
			})

		}
	}

	log.Printf("[generateSwapXdr]obtained source account balance:\n%v balance is %v\n%v balance is %v\n", nativeAssetCode, sourceAccountNativeBalance, sourceAsset.GetCode(), sourceAccountCustomBalance)

	amountToSwapDec := amountToSwap

	if sourceAsset.IsNative() {
		if !sourceAccountTrustsDestinationAsset {
			if sourceAccountNativeBalance.LessThan(amountToSwapDec.Add(appliedCharge)) {
				return "", &tErrors.ErrorUnderfundedAccount{Detail: fmt.Sprintf("Not enough funds. Needs Extra %v %v to accommodate the amount needed to opt you into the destination asset or you reduce same from the amount you want to swap.", (amountToSwapDec.Add(appliedCharge)).Sub(sourceAccountNativeBalance), os.Getenv("NATIVE_ASSET_CODE"))}
			}
		} else {
			if sourceAccountNativeBalance.LessThan(amountToSwapDec) {
				return "", &tErrors.ErrorUnderfundedAccount{Detail: fmt.Sprintf("Not enough funds. Needs Extra %v %v or you reduce same from the amount you want to swap.", (amountToSwapDec).Sub(sourceAccountNativeBalance), os.Getenv("NATIVE_ASSET_CODE"))}
			}
		}

	} else {

		if !sourceAccountTrustsDestinationAsset {
			if sourceAccountNativeBalance.LessThan(appliedCharge) {
				return "", &tErrors.ErrorUnderfundedAccount{Detail: fmt.Sprintf("Not enough funds. Needs Extra %v %v to complete this transaction", appliedCharge.Sub(sourceAccountNativeBalance), os.Getenv("NATIVE_ASSET_CODE"))}
			}
		}
		if sourceAccountCustomBalance.LessThan(amountToSwapDec) {
			return "", &tErrors.ErrorUnderfundedAccount{Detail: fmt.Sprintf("Not enough funds. Needs Extra %v %v or you reduce same from the amount you want to swap.", (amountToSwapDec).Sub(sourceAccountCustomBalance), sourceAsset.GetCode())}
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
	path, swappedEstimate, err := GetStrictSendPaths(pathInput, client)
	if err != nil {
		log.Println("[generateSwapXdr]error fetching valid swap Path ", err)
		return "", err
	}

	//native asset
	ops = append(ops, &txnbuild.PathPaymentStrictSend{
		SendAsset:     sourceAsset,
		SendAmount:    swapInfo.SwapAmount,
		Destination:   wallet.ID,
		DestAsset:     destinationAsset,
		DestMin:       swapDestMin.String(),
		Path:          path,
		SourceAccount: wallet.ID,
	})
	serviceFee, e := decimal.NewFromString(swapInfo.FeeAmount)
	if e != nil {
		serviceFee = decimal.Zero
	}
	// totalFees = totalFees.Add(serviceFee)
	// feeLabel := swapInfo.Fee + "%"
	signForFeeTrustLine := 0
	if serviceFee.IsPositive() && os.Getenv("SWAP_FEE_ENABLED") == "1" {
		//process service fee

		//ensure that the fee address is can accept the asset.
		// but bcos  fee address needs to sign, it cannot be done here
		feeKeypair := keypair.MustParseFull(os.Getenv("SWAP_FEE_WALLET"))
		feeAddress := feeKeypair.Address()

		if !sourceAsset.IsNative() {

			_, feeAccountTrustsAsset, _, _, _, _ := network.BlockchainAccountProperties(gc.BantuExpansionClient, feeAddress, sourceAsset)
			if !feeAccountTrustsAsset {
				signForFeeTrustLine = 1
				//establish trustline automatically
				ops = append(ops, &txnbuild.ChangeTrust{
					Line:          txnbuild.ChangeTrustAssetWrapper{Asset: sourceAsset},
					Limit:         "900000000000",
					SourceAccount: feeAddress,
				})

			}
		}

		ops = append(ops, &txnbuild.Payment{
			Destination:   feeAddress,
			Amount:        serviceFee.String(),
			SourceAccount: wallet.ID,
			Asset:         sourceAsset,
		})

		messages = append(messages, "Service fee will apply.")

	}

	// Construct the transaction that holds the operations to execute on the network
	var memoSAC, memoDAC string
	memoSAC = swapInfo.SourceAssetCode
	memoDAC = swapInfo.DestinationAssetCode
	if swapInfo.SourceAssetIssuer == "" || swapInfo.SourceAssetIssuer == "native" {
		memoSAC = os.Getenv("NATIVE_ASSET_CODE")
	}
	if swapInfo.DestinationAssetIssuer == "" || swapInfo.DestinationAssetIssuer == "native" {
		memoDAC = os.Getenv("NATIVE_ASSET_CODE")
	}

	memo := fmt.Sprintf("%v>%v", memoSAC, memoDAC)
	log.Println("[generateSwapXdr] Memo:", memo)
	swapInfo.Memo = memo

	var tx *txnbuild.Transaction
	// Construct the transaction that holds the operations to execute on the network
	// if swapInfo.Multiparty == 1 {
	swapInfo.TransactionSource = chanSourceAccount.AccountID
	tx, err = txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        chanSourceAccount,
			IncrementSequenceNum: true,
			Operations:           ops,
			BaseFee:              2000,
			Preconditions: txnbuild.Preconditions{
				TimeBounds: txnbuild.NewInfiniteTimeout(),
			},
			Memo: txnbuild.MemoText(memo),
		},
	)

	if err != nil {
		log.Println("[generatePaymentXdr] error constructing transaction ", err)
		return "", err
	}

	// if swapInfo.Multiparty == 1 {

	tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), chanAccount)

	if err != nil {
		log.Println("[generateSwapXdr] error signing transaction with channelAccount key ", err)
		return "", &tErrors.ErrorTemporaryServerError{}
	}
	// }

	if signForFeeTrustLine == 1 && !sourceAsset.IsNative() {
		log.Printf("[generateSwapXdr] <<<<<<<<<<<<<<<<<<<<<<<<<<<< signing transaction with swap fee key>>>>>>>>>>>>>>>>>>>>>>>>:[%v]\n\n", sourceAsset)

		feeKeypair := keypair.MustParseFull(os.Getenv("SWAP_FEE_WALLET"))

		tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), feeKeypair)
		if err != nil {
			log.Println("[generateSwapXdr] error signing transaction with swap fee key", err)
			return "", &tErrors.ErrorTemporaryServerError{}
		}
	}

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

// generateSwapReceiveXdr generates xdr for strict receive operation. Returns the base64 xdr transaction string, the operation object, and error
func generateSwapReceiveXdr(wallet *userModels.UserWallet, swapInfo *swapModels.SwapReceiveInfo, gc *sharedconfig.GlobalConfig) (string, []txnbuild.Operation, error) {
	// baseReserve := network.GetBlockchainBaseReserve()
	// charge := baseReserve.Mul(decimal.NewFromInt(3)).Truncate(7).String()
	client := gc.BantuExpansionClient
	messages := make([]string, 0)
	nativeAssetCode := os.Getenv("NATIVE_ASSET_CODE")
	var err error
	var amountToSwap decimal.Decimal

	if amountToSwap, err = decimal.NewFromString(swapInfo.DestinationAmount); err != nil {
		return "", []txnbuild.Operation{}, &swapErrors.ErrorInvalidSwapAmount{}
	}

	newAmountToSwap := amountToSwap.Truncate(7).String()

	var sourceAsset txnbuild.Asset = txnbuild.NativeAsset{}
	var destinationAsset txnbuild.Asset = txnbuild.NativeAsset{}

	if len(swapInfo.DestinationAssetCode) != 0 && !strings.EqualFold(swapInfo.DestinationAssetCode, nativeAssetCode) {

		destinationAsset = txnbuild.CreditAsset{Code: swapInfo.DestinationAssetCode, Issuer: swapInfo.DestinationAssetIssuer}
	}
	if len(swapInfo.SourceAssetCode) != 0 && !strings.EqualFold(swapInfo.SourceAssetCode, nativeAssetCode) {

		sourceAsset = txnbuild.CreditAsset{Code: swapInfo.SourceAssetCode, Issuer: swapInfo.SourceAssetIssuer}
	}

	// charge := baseReserve.Mul(decimal.NewFromInt(1)).Truncate(7).String()
	appliedCharge := decimal.NewFromFloat(0)
	swapInfo.Messages = messages
	var ops []txnbuild.Operation = make([]txnbuild.Operation, 0)
	chanAccount := <-gc.ChannelAccounts
	defer func(c *keypair.Full) {
		gc.ChannelAccounts <- c
	}(chanAccount)

	_, _, _, _, chanSourceAccount, _ := network.BlockchainAccountProperties(client, chanAccount.Address(), txnbuild.NativeAsset{})

	sourceAccountExists, _, sourceAccountNativeBalance, sourceAccountCustomBalance, _, sourceAccountErr := network.BlockchainAccountProperties(client, wallet.ID, sourceAsset)
	var sourceAccountTrustsDestinationAsset bool
	if !destinationAsset.IsNative() {
		_, sourceAccountTrustsDestinationAsset, _, _, _, _ = network.BlockchainAccountProperties(client, wallet.ID, destinationAsset)

	}
	// destAccountExists, destAccountTrustsDestinationAsset, destAccountNativeBalance, destAccountCustomBalance, _, destAccountErr := network.BlockchainAccountProperties(client, swapInfo.DestinationAccount, destinationAsset)
	_, destAccountTrustsDestinationAsset, _, _, _, _ := network.BlockchainAccountProperties(client, swapInfo.DestinationAccount, destinationAsset)

	if sourceAccountErr != nil {
		return "", []txnbuild.Operation{}, sourceAccountErr
	}

	if !sourceAccountExists {
		return "", []txnbuild.Operation{}, &tErrors.ErrorUnderfundedAccount{}
	}

	if !destinationAsset.IsNative() {

		if !destAccountTrustsDestinationAsset {
			return "", []txnbuild.Operation{}, &tErrors.CustomError{
				Param:      "destinationAccount",
				Err:        "error-destination-cannot-accept-asset",
				ErrMessage: "Destination address cannot accept the asset " + swapInfo.DestinationAssetCode,
			}

		}
	}

	amountToSwapDec := amountToSwap

	if sourceAsset.IsNative() {
		if !sourceAccountTrustsDestinationAsset {
			if sourceAccountNativeBalance.LessThan(amountToSwapDec.Add(appliedCharge)) {
				return "", []txnbuild.Operation{}, &tErrors.ErrorUnderfundedAccount{Detail: fmt.Sprintf("Not enough funds. Needs Extra %v %v to accommodate the amount needed to opt you into the destination asset or you reduce same from the amount you want to swap.", (amountToSwapDec.Add(appliedCharge)).Sub(sourceAccountNativeBalance), os.Getenv("NATIVE_ASSET_CODE"))}
			}
		} else {
			if sourceAccountNativeBalance.LessThan(amountToSwapDec) {
				return "", []txnbuild.Operation{}, &tErrors.ErrorUnderfundedAccount{Detail: fmt.Sprintf("Not enough funds. Needs Extra %v %v or you reduce same from the amount you want to swap.", (amountToSwapDec).Sub(sourceAccountNativeBalance), os.Getenv("NATIVE_ASSET_CODE"))}
			}
		}

	} else {

		if !sourceAccountTrustsDestinationAsset {
			if sourceAccountNativeBalance.LessThan(appliedCharge) {
				return "", []txnbuild.Operation{}, &tErrors.ErrorUnderfundedAccount{Detail: fmt.Sprintf("Not enough funds. Needs Extra %v %v to complete this transaction", appliedCharge.Sub(sourceAccountNativeBalance), os.Getenv("NATIVE_ASSET_CODE"))}
			}
		}
		if sourceAccountCustomBalance.LessThan(amountToSwapDec) {
			return "", []txnbuild.Operation{}, &tErrors.ErrorUnderfundedAccount{Detail: fmt.Sprintf("Not enough funds. Needs Extra %v %v or you reduce same from the amount you want to swap.", (amountToSwapDec).Sub(sourceAccountCustomBalance), sourceAsset.GetCode())}
		}

	}

	//get sendPath
	//using DestinationAccount will get paths to all assets in the destination account.
	//using destinationAssets gets path to only the asset
	sourceAssets := ""
	if !sourceAsset.IsNative() {
		sourceAssets = fmt.Sprintf("%s:%s", swapInfo.SourceAssetCode, swapInfo.SourceAssetIssuer)
	}
	pathInput := swapModels.SwapPathInput{
		SourceAssets: sourceAssets,
		// SourceAccount:          swapInfo.SourceAccount,
		DestinationAssetCode:   swapInfo.DestinationAssetCode,
		DestinationAssetIssuer: swapInfo.DestinationAssetIssuer,
		DestinationAmount:      newAmountToSwap,
	}
	path, requiredEstimate, err := GetStrictReceivePaths(pathInput, client)
	if err != nil {
		log.Println("[generateSwapReceiveXdr]error fetching valid swap Path ", err)
		return "", []txnbuild.Operation{}, err
	}

	//check is source account has minimum of that required estimate
	if decimal.RequireFromString(requiredEstimate).GreaterThan(sourceAccountNativeBalance) {
		acode := swapInfo.SourceAssetCode
		if len(acode) == 0 {
			acode = nativeAssetCode
		}
		return "", []txnbuild.Operation{}, &tErrors.ErrorUnderfundedAccount{
			Detail: fmt.Sprintf("You need extra %v %v to complete the transaction at this time.", decimal.RequireFromString(requiredEstimate).Sub(sourceAccountNativeBalance), acode),
		}
	}

	//native asset
	ops = append(ops, &txnbuild.PathPaymentStrictReceive{
		SendAsset: sourceAsset,
		// SendMax:       requiredEstimate,
		Destination:   swapInfo.DestinationAccount,
		DestAsset:     destinationAsset,
		DestAmount:    swapInfo.DestinationAmount,
		Path:          path,
		SourceAccount: wallet.ID,
	})
	// serviceFee, e := decimal.NewFromString(swapInfo.FeeAmount)
	// if e != nil {
	// 	serviceFee = decimal.Zero
	// }
	// totalFees = totalFees.Add(serviceFee)
	// feeLabel := swapInfo.Fee + "%"
	signForFeeTrustLine := 0
	// if serviceFee.IsPositive() && os.Getenv("SWAP_FEE_ENABLED") == "1" {
	// 	//process service fee

	// 	//ensure that the fee address is can accept the asset.
	// 	// but bcos  fee address needs to sign, it cannot be done here
	// 	feeKeypair := keypair.MustParseFull(os.Getenv("SWAP_FEE_WALLET"))
	// 	feeAddress := feeKeypair.Address()

	// 	if !sourceAsset.IsNative() {

	// 		_, feeAccountTrustsAsset, _, _, _, _ := network.BlockchainAccountProperties(gc.BantuExpansionClient, feeAddress, sourceAsset)
	// 		if !feeAccountTrustsAsset {
	// 			signForFeeTrustLine = 1
	// 			//establish trustline automatically
	// 			ops = append(ops, &txnbuild.ChangeTrust{
	// 				Line:          txnbuild.ChangeTrustAssetWrapper{Asset: sourceAsset},
	// 				Limit:         "900000000000",
	// 				SourceAccount: feeAddress,
	// 			})

	// 		}
	// 	}

	// 	ops = append(ops, &txnbuild.Payment{
	// 		Destination:   feeAddress,
	// 		Amount:        serviceFee.String(),
	// 		SourceAccount: wallet.ID,
	// 		Asset:         sourceAsset,
	// 	})

	// 	messages = append(messages, "Service fee will apply.")

	// }

	// Construct the transaction that holds the operations to execute on the network
	var memoSAC, memoDAC string
	memoSAC = swapInfo.SourceAssetCode
	memoDAC = swapInfo.DestinationAssetCode
	if swapInfo.SourceAssetIssuer == "" || swapInfo.SourceAssetIssuer == "native" {
		memoSAC = os.Getenv("NATIVE_ASSET_CODE")
	}
	if swapInfo.DestinationAssetIssuer == "" || swapInfo.DestinationAssetIssuer == "native" {
		memoDAC = os.Getenv("NATIVE_ASSET_CODE")
	}

	memo := fmt.Sprintf("%v>%v", memoSAC, memoDAC)
	log.Println("[generateSwapReceiveXdr] Memo:", memo)
	swapInfo.Memo = memo

	var tx *txnbuild.Transaction
	// Construct the transaction that holds the operations to execute on the network
	// if swapInfo.Multiparty == 1 {
	swapInfo.TransactionSource = chanSourceAccount.AccountID
	tx, err = txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        chanSourceAccount,
			IncrementSequenceNum: true,
			Operations:           ops,
			BaseFee:              2000,
			Preconditions: txnbuild.Preconditions{
				TimeBounds: txnbuild.NewInfiniteTimeout(),
			},
			Memo: txnbuild.MemoText(memo),
		},
	)

	if err != nil {
		log.Println("[generateSwapReceiveXdr] error constructing transaction ", err)
		return "", ops, err
	}

	// if swapInfo.Multiparty == 1 {

	tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), chanAccount)

	if err != nil {
		log.Println("[generateSwapReceiveXdr] error signing transaction with channelAccount key ", err)
		return "", ops, &tErrors.ErrorTemporaryServerError{}
	}
	// }

	if signForFeeTrustLine == 1 && !sourceAsset.IsNative() {
		log.Printf("[generateSwapReceiveXdr] <<<<<<<<<<<<<<<<<<<<<<<<<<<< signing transaction with swap fee key>>>>>>>>>>>>>>>>>>>>>>>>:[%v]\n\n", sourceAsset)

		feeKeypair := keypair.MustParseFull(os.Getenv("SWAP_FEE_WALLET"))

		tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), feeKeypair)
		if err != nil {
			log.Println("[generateSwapReceiveXdr] error signing transaction with swap fee key", err)
			return "", ops, &tErrors.ErrorTemporaryServerError{}
		}
	}

	if err != nil {
		log.Println("[generateSwapReceiveXdr]error constructing transaction ", err)
		return "", ops, err
	}

	var xdrBase64 string

	xdrBase64, err = tx.Base64()
	if err != nil {
		log.Println("[generateSwapXdr] error getting txn base64", err)
		return "", ops, err
	}
	swapInfo.Memo = memo
	swapInfo.Messages = messages
	swapInfo.RequiredEstimate = requiredEstimate
	return xdrBase64, ops, nil
}

// GetStrictSendPaths gets Strict Send Paths for Strict Send Path Payment request
func GetStrictSendPaths(pathInput swapModels.SwapSendPathInput, client *horizonclient.Client) (paths []txnbuild.Asset, swappedEstimate string, err error) {
	discord.WebhookURL = "https://discord.com/api/webhooks/824381163367170058/75RxS1LzWA800hWereJJumw"
	if len(os.Getenv("EXPANSION_NETWORK_ERROR_WEBHOOK")) > 50 {
		discord.WebhookURL = os.Getenv("EXPANSION_NETWORK_ERROR_WEBHOOK")
	}
	var swapPaths horizon.PathsPage
	paths = make([]txnbuild.Asset, 0)
	var sourceAssetType horizonclient.AssetType
	if len(pathInput.SourceAssetIssuer) == 0 {
		sourceAssetType = horizonclient.AssetTypeNative
		pathInput.SourceAssetCode = ""
		pathInput.SourceAssetIssuer = ""
	} else if len(pathInput.SourceAssetCode) < 5 && len(pathInput.SourceAssetIssuer) == 56 {
		sourceAssetType = horizonclient.AssetType4
	} else if len(pathInput.SourceAssetCode) > 4 && len(pathInput.SourceAssetCode) <= 12 && len(pathInput.SourceAssetIssuer) == 56 {
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

	swapPaths, err = client.StrictSendPaths(sspr)

	if err != nil {
		if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "handshake") || strings.Contains(err.Error(), "read tcp") || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp") || strings.Contains(err.Error(), "no such host") {
			log.Println("#######################@@@@@@@@@@@@@@@@[client.StrictSendPathsErr] expansion connection problem:", err)
			discord.Say(fmt.Sprintf("[getStrictSendPaths] error connecting to expansion service: %v\nSwapPathRequest: %+v", err, sspr))

			return paths, "", &tErrors.ErrorTemporaryServerError{}
		}
		if strings.Contains(err.Error(), "liquid") {
			destAsset := os.Getenv("NATIVE_ASSET_CODE")
			sourceAsset := os.Getenv("NATIVE_ASSET_CODE")
			if len(pathInput.SourceAssetCode) > 0 {
				sourceAsset = pathInput.SourceAssetCode
			}
			if pathInput.DestinationAssets != "native" {
				destAsset = strings.Split(pathInput.DestinationAssets, ":")[0]
			}
			return paths, "", &tErrors.CustomError{
				Param:      "destinationAssetCode",
				Err:        "error-low-liquidity",
				ErrMessage: fmt.Sprintf("There is not enough %v market to exchange for your %v at this time. Please try again later or reduce the quantity of %v to try again.", destAsset, sourceAsset, sourceAsset),
			}
		}
		horizonException, ok := err.(*horizonclient.Error)

		if ok {

			extraErrors := horizonException.Problem.Extras

			for key, val := range extraErrors {
				log.Printf("Extras: %v is %v\n", key, val)
			}

			resultCodes, e := horizonException.ResultCodes()
			if e != nil {
				log.Println("[client.StrictSendPathsErr] Error getting result codes:", e)

				destAsset := os.Getenv("NATIVE_ASSET_CODE")
				sourceAsset := os.Getenv("NATIVE_ASSET_CODE")
				if len(pathInput.SourceAssetCode) > 0 {
					sourceAsset = pathInput.SourceAssetCode
				}
				if pathInput.DestinationAssets != "native" {
					destAsset = strings.Split(pathInput.DestinationAssets, ":")[0]
				}
				return paths, "", &tErrors.CustomError{
					Param:      "destinationAssetCode",
					Err:        "error-low-liquidity",
					ErrMessage: fmt.Sprintf("There is not enough %v market to exchange for your %v at this time. Please try again later or reduce the quantity of %v to try again.", destAsset, sourceAsset, sourceAsset),
				}

			}

			for key, val := range resultCodes.OperationCodes {
				log.Printf("Result code: %v is %v\n", key, val)
			}
			discord.Say(fmt.Sprintf("[getStrictSendPaths] error submitting: %v\nSwapPathRequest: %+v\nResultCodes: %+v", err, sspr, resultCodes))

		}
		log.Println("[client.StrictSendPathsErr] Error submitting:", err)
		if strings.Contains(err.Error(), "liquid") {
			destAsset := os.Getenv("NATIVE_ASSET_CODE")
			sourceAsset := os.Getenv("NATIVE_ASSET_CODE")
			if len(pathInput.SourceAssetCode) > 0 {
				sourceAsset = pathInput.SourceAssetCode
			}
			if pathInput.DestinationAssets != "native" {
				destAsset = strings.Split(pathInput.DestinationAssets, ":")[0]
			}
			return paths, "", &tErrors.CustomError{
				Param:      "destinationAssetCode",
				Err:        "error-low-liquidity",
				ErrMessage: fmt.Sprintf("There is not enough %v market to exchange for your %v at this time. Please try again later or reduce the quantity of %v to try again.", destAsset, sourceAsset, sourceAsset),
			}
		}
		return paths, "", &tErrors.ErrorTemporaryServerError{}

	}
	// discord.Say(fmt.Sprintf("[getStrictSendPaths] swapPaths: %+v\nRequestParams: %+v", swapPaths, sspr))
	// log.Printf("[getStrictSendPaths] swapPaths: %+v\n", swapPaths)
	if len(swapPaths.Embedded.Records) == 0 {
		destAsset := os.Getenv("NATIVE_ASSET_CODE")
		sourceAsset := os.Getenv("NATIVE_ASSET_CODE")
		if len(pathInput.SourceAssetCode) > 0 {
			sourceAsset = pathInput.SourceAssetCode
		}
		if pathInput.DestinationAssets != "native" {
			destAsset = strings.Split(pathInput.DestinationAssets, ":")[0]
		}
		return paths, "", &tErrors.CustomError{
			Param:      "destinationAssetCode",
			Err:        "error-low-liquidity",
			ErrMessage: fmt.Sprintf("There is not enough %v market to exchange for your %v at this time. Please try again later or reduce the quantity of %v to try again.", destAsset, sourceAsset, sourceAsset),
		}
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

// getStrictReceivePaths gets Strict Receive Paths for Strict Receive Path Payment request
func GetStrictReceivePaths(pathInput swapModels.SwapPathInput, client *horizonclient.Client) (paths []txnbuild.Asset, requiredEstimate string, err error) {
	discord.WebhookURL = "https://discord.com/api/webhooks/824381163367170058/75RxS1LzWA800hWereJJumw"
	if len(os.Getenv("EXPANSION_NETWORK_ERROR_WEBHOOK")) > 50 {
		discord.WebhookURL = os.Getenv("EXPANSION_NETWORK_ERROR_WEBHOOK")
	}
	var swapPaths horizon.PathsPage
	paths = make([]txnbuild.Asset, 0)
	var destinationAssetType horizonclient.AssetType
	if len(pathInput.DestinationAssetIssuer) == 0 {
		destinationAssetType = horizonclient.AssetTypeNative
		pathInput.DestinationAssetCode = ""
		pathInput.DestinationAssetIssuer = ""
	} else if len(pathInput.DestinationAssetCode) < 5 && len(pathInput.DestinationAssetIssuer) == 56 {
		destinationAssetType = horizonclient.AssetType4
	} else if len(pathInput.DestinationAssetCode) > 4 && len(pathInput.DestinationAssetCode) <= 12 && len(pathInput.DestinationAssetIssuer) == 56 {
		destinationAssetType = horizonclient.AssetType12
	}
	// if pathInput.SourceAccount != "" {
	// 	pathInput.SourceAssets = ""
	// }

	if pathInput.SourceAssets == "" {
		pathInput.SourceAssets = "native"
	}

	sspr := horizonclient.PathsRequest{
		DestinationAssetType:   destinationAssetType,
		DestinationAssetCode:   pathInput.DestinationAssetCode,
		DestinationAssetIssuer: pathInput.DestinationAssetIssuer,
		DestinationAmount:      pathInput.DestinationAmount,
		SourceAssets:           pathInput.SourceAssets,
	}

	swapPaths, err = client.StrictReceivePaths(sspr)

	if err != nil {
		if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "handshake") || strings.Contains(err.Error(), "read tcp") || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "dial tcp") || strings.Contains(err.Error(), "no such host") {
			log.Println("#######################@@@@@@@@@@@@@@@@[client.StrictSendPathsErr] expansion connection problem:", err)
			discord.Say(fmt.Sprintf("[getStrictSendPaths] error connecting to expansion service: %v\nSwapPathRequest: %+v", err, sspr))

			return paths, "", &tErrors.ErrorTemporaryServerError{}
		}
		if strings.Contains(err.Error(), "liquid") {
			destAssetCode := os.Getenv("NATIVE_ASSET_CODE")
			sourceAssetCode := os.Getenv("NATIVE_ASSET_CODE")
			if len(pathInput.DestinationAssetCode) > 0 {
				destAssetCode = pathInput.DestinationAssetCode
			}
			if pathInput.SourceAssets != "native" {
				sourceAssetCode = strings.Split(pathInput.SourceAssets, ":")[0]
			}
			return paths, "", &tErrors.CustomError{
				Param:      "destinationAssetCode",
				Err:        "error-low-liquidity",
				ErrMessage: fmt.Sprintf("There is not enough %v market to exchange for your %v at this time. Please try again later or reduce the quantity of %v to try again.", destAssetCode, sourceAssetCode, destAssetCode),
			}
		}
		horizonException, ok := err.(*horizonclient.Error)

		if ok {

			extraErrors := horizonException.Problem.Extras

			for key, val := range extraErrors {
				log.Printf("Extras: %v is %v\n", key, val)
			}

			resultCodes, e := horizonException.ResultCodes()
			if e != nil {
				log.Println("[client.StrictSendPathsErr] Error getting result codes:", e)

				destAssetCode := os.Getenv("NATIVE_ASSET_CODE")
				sourceAssetCode := os.Getenv("NATIVE_ASSET_CODE")
				if len(pathInput.DestinationAssetCode) > 0 {
					destAssetCode = pathInput.DestinationAssetCode
				}
				if pathInput.SourceAssets != "native" {
					sourceAssetCode = strings.Split(pathInput.SourceAssets, ":")[0]
				}
				return paths, "", &tErrors.CustomError{
					Param:      "destinationAssetCode",
					Err:        "error-low-liquidity",
					ErrMessage: fmt.Sprintf("There is not enough %v market to exchange for your %v at this time. Please try again later or reduce the quantity of %v to try again.", destAssetCode, sourceAssetCode, destAssetCode),
				}

			}

			for key, val := range resultCodes.OperationCodes {
				log.Printf("Result code: %v is %v\n", key, val)
			}
			discord.Say(fmt.Sprintf("[getStrictSendPaths] error submitting: %v\nSwapPathRequest: %+v\nResultCodes: %+v", err, sspr, resultCodes))

		}
		log.Println("[client.StrictSendPathsErr] Error submitting:", err)
		if strings.Contains(err.Error(), "liquid") {
			destAssetCode := os.Getenv("NATIVE_ASSET_CODE")
			sourceAssetCode := os.Getenv("NATIVE_ASSET_CODE")
			if len(pathInput.DestinationAssetCode) > 0 {
				destAssetCode = pathInput.DestinationAssetCode
			}
			if pathInput.SourceAssets != "native" {
				sourceAssetCode = strings.Split(pathInput.SourceAssets, ":")[0]
			}
			return paths, "", &tErrors.CustomError{
				Param:      "destinationAssetCode",
				Err:        "error-low-liquidity",
				ErrMessage: fmt.Sprintf("There is not enough %v market to exchange for your %v at this time. Please try again later or reduce the quantity of %v to try again.", destAssetCode, sourceAssetCode, destAssetCode),
			}
		}
		return paths, "", &tErrors.ErrorTemporaryServerError{}

	}
	// discord.Say(fmt.Sprintf("[getStrictSendPaths] swapPaths: %+v\nRequestParams: %+v", swapPaths, sspr))
	// log.Printf("[getStrictSendPaths] swapPaths: %+v\n", swapPaths)
	if len(swapPaths.Embedded.Records) == 0 {
		destAssetCode := os.Getenv("NATIVE_ASSET_CODE")
		sourceAssetCode := os.Getenv("NATIVE_ASSET_CODE")
		if len(pathInput.DestinationAssetCode) > 0 {
			destAssetCode = pathInput.DestinationAssetCode
		}
		if pathInput.SourceAssets != "native" {
			sourceAssetCode = strings.Split(pathInput.SourceAssets, ":")[0]
		}
		return paths, "", &tErrors.CustomError{
			Param:      "destinationAssetCode",
			Err:        "error-low-liquidity",
			ErrMessage: fmt.Sprintf("There is not enough %v market to exchange for your %v at this time. Please try again later or reduce the quantity of %v to try again.", destAssetCode, sourceAssetCode, destAssetCode),
		}
	}
	// sourceAmountDec, _ := decimal.NewFromString(swapPaths.Embedded.Records[0].SourceAmount)

	bestPath := swapPaths.Embedded.Records[0]
	//build assets
	requiredEstimate = bestPath.SourceAmount

	for _, v := range bestPath.Path {
		if len(v.Issuer) == 0 {
			paths = append(paths, txnbuild.NativeAsset{})
		} else {
			paths = append(paths, txnbuild.CreditAsset{Code: v.Code, Issuer: v.Issuer})
		}

	}

	return paths, requiredEstimate, nil
}

func logDiscordFailedSwap(msg string) {
	discord.WebhookURL = "https://discord.com/api/webhooks/827986576415129663/wqMKp9wxB_fxs9Q3zlMKCNPGENXmD_ueUnL8hVCu1wmRfD2wkXAjfP85k1Ro_2_wGfiY"
	if len(os.Getenv("FAILED_PAYMENT_ERROR_WEBHOOK")) > 50 {
		discord.WebhookURL = os.Getenv("FAILED_PAYMENT_ERROR_WEBHOOK")
	}
	discord.Say(msg)
}

package users

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"trovo-wallet-api/internal/basetxn"
	userBc "trovo-wallet-api/internal/components/users/blockchain"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/evmkeypair"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/sharedconfig"
	"trovo-wallet-api/internal/validators"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm/clause"
)

// ClaimPendingAsset claim pending assets
func ClaimPendingAsset(signerUser *userModels.User, wallet *userModels.UserWallet, pendingAssetToClaim *userModels.PendingAssetToClaim, gc *sharedconfig.GlobalConfig) (*userModels.PendingAssetToClaim, bool, error) {
	if gc.IsValidTokenizedAsset(pendingAssetToClaim.AssetCode) {
		t := gc.GetTokenizedAssetByCode(pendingAssetToClaim.AssetCode)
		if t.AssetTokenizationStatus < 5 {
			return pendingAssetToClaim, false, &tErrors.CustomError{
				Param:      "assetIssuer",
				Err:        "error-asset-not-yet-available-for-sale",
				ErrMessage: "Asset is not yet available for sale.",
				Code:       http.StatusForbidden,
			}
		}
	}

	if len(pendingAssetToClaim.AssetIssuer) == 0 {
		return pendingAssetToClaim, false, &tErrors.CustomError{
			Param:      "assetIssuer",
			Err:        "error-missing-parameter",
			ErrMessage: "Asset issuer is invalid",
			Code:       http.StatusBadRequest,
		}
	}

	if wallet.NumberOfApprovalsNeeded > 0 && wallet.SharedAccessEnabled == 1 {
		pendingAssetToClaim.Multiparty = 1
	}
	if wallet.HasViewOnlyAccess(gc) {
		pendingAssetToClaim.SignatureRequired = 1
	}
	var err error
	if wallet.Signer != signerUser.PrimarySigner && pendingAssetToClaim.Multiparty == 0 {
		return pendingAssetToClaim, false, &tErrors.CustomError{Param: "publicKey", Err: "error-wallet-not-managed-by-user", ErrMessage: "You do not have permission to operate on this wallet", Code: http.StatusBadRequest}

	}
	//validators
	{

		err = validators.ValidateAssetCodeFormat(pendingAssetToClaim.AssetCode)

		if err != nil {
			return pendingAssetToClaim, false, err
		}

		err = validators.ValidateAddressFormat(pendingAssetToClaim.AssetIssuer)

		if err != nil {
			return pendingAssetToClaim, false, err
		}
	}

	err = gc.DB.Where("transaction_type = 'ACCEPT PENDING ASSET' AND transaction_status = 'PENDING' AND wallet_address = ? AND description like ?", wallet.ID, "%"+pendingAssetToClaim.AssetCode+":%").First(&userModels.PendingAuth{}).Error
	if err == nil {
		return pendingAssetToClaim, false, &tErrors.CustomError{
			Param:      "assetCode",
			Err:        "error-duplicate-entry",
			ErrMessage: "There is already a pending asset claim request. Please fulfil that one first. or reject it before continuing.",
			Code:       http.StatusBadRequest,
		}
	}

	horizonClient := network.GetBlockchainClient()

	xdrBase64, err := generateClaimPendingAssetXdr(wallet, pendingAssetToClaim, gc)

	if err != nil {
		return pendingAssetToClaim, false, err
	}

	oldTxn := pendingAssetToClaim.Transaction

	pendingAssetToClaim.NetworkPassPhrase = network.GetBlockchainNetworkPassPhrase()

	pendingAssetToClaim.Transaction = xdrBase64

	if len(pendingAssetToClaim.TransactionSignature) == 0 && pendingAssetToClaim.Commit == 0 {
		//no signature
		return pendingAssetToClaim, false, err

	}

	//there was a signature... let's submit

	if oldTxn != xdrBase64 && pendingAssetToClaim.Commit == 0 && pendingAssetToClaim.Multiparty == 0 {
		return pendingAssetToClaim, false, &tErrors.CustomError{
			Param:      "transaction",
			Err:        "transaction mismatch",
			ErrMessage: "transaction mismatch, please try again",
			Code:       404,
		}
	}
	if (pendingAssetToClaim.Commit == 0 && wallet.SharedAccessEnabled == 1 && wallet.NumberOfApprovalsNeeded == 0) || (wallet.SharedAccessEnabled == 0 && pendingAssetToClaim.Commit == 0) {
		txnID, err := network.SubmitXdrWithSignature(horizonClient, wallet.Signer, xdrBase64, pendingAssetToClaim.TransactionSignature)

		if err != nil {
			log.Printf("[ClaimPendingAsset]error submitting txn: %v\n", err)
			return pendingAssetToClaim, false, &tErrors.ErrorTemporaryServerError{}
		}

		pendingAssetToClaim.TransactionID = txnID

		return pendingAssetToClaim, true, nil
	}

	if pendingAssetToClaim.Multiparty == 1 {
		pendingAssetToClaim.TransactionID = "PENDING_AUTH"
		log.Printf("[ClaimPendingAsset]shared access with approver permission enabled for %v \n", wallet.Alias)
		id := uuid.NewString()
		assetOfPayment := os.Getenv("NATIVE_ASSET_CODE")
		if len(pendingAssetToClaim.AssetIssuer) == 42 {
			assetOfPayment = fmt.Sprintf("%v:%v...%v", pendingAssetToClaim.AssetCode, pendingAssetToClaim.AssetIssuer[0:4], pendingAssetToClaim.AssetIssuer[51:55])
		}
		description := fmt.Sprintf("Accept & claim pending balance for asset %v.\nMessages:%v", assetOfPayment, pendingAssetToClaim.Messages)
		transactionByte, _ := json.Marshal(*pendingAssetToClaim)
		transactionStr := string(transactionByte)
		pendingAuth := userModels.PendingAuth{
			ID:                     id,
			Initiator:              signerUser.Username,
			InitiatorSignerAddress: signerUser.PrimarySigner,
			WalletAddress:          wallet.ID,
			TransactionType:        "ACCEPT PENDING ASSET",
			Description:            description,
			TransactionSource:      pendingAssetToClaim.TransactionSource,
			ApprovalsNeeded:        wallet.NumberOfApprovalsNeeded,
			TransactionXdr:         xdrBase64,
			TransactionInfoStr:     &transactionStr,
		}
		//save and commit this to database
		e := gc.DB.Omit(clause.Associations).Create(&pendingAuth).Error
		if e != nil {
			log.Printf("[ClaimPendingAsset] Error saving payment txn [%+v] transaction on pending auth table: %s\n", pendingAuth, e.Error())
			err = &tErrors.ErrorTemporaryServerError{}
			return pendingAssetToClaim, false, err
		}
		pendingAssetToClaim.ReturnedDescription = description
		return pendingAssetToClaim, true, nil

	}

	return pendingAssetToClaim, false, &tErrors.ErrorTemporaryServerError{}
}

// RejectPendingAsset rejects pending assets
func RejectPendingAsset(signerUser *userModels.User, wallet *userModels.UserWallet, pendingAssetToClaim *userModels.PendingAssetToClaim, gc *sharedconfig.GlobalConfig) (*userModels.PendingAssetToClaim, bool, error) {
	if gc.IsValidTokenizedAsset(pendingAssetToClaim.AssetCode) {
		t := gc.GetTokenizedAssetByCode(pendingAssetToClaim.AssetCode)
		if t.AssetTokenizationStatus < 5 {
			return pendingAssetToClaim, false, &tErrors.CustomError{
				Param:      "assetIssuer",
				Err:        "error-asset-not-yet-available-for-sale",
				ErrMessage: "Asset is not yet available for sale.",
				Code:       http.StatusForbidden,
			}
		}
	}
	if wallet.NumberOfApprovalsNeeded > 0 && wallet.SharedAccessEnabled == 1 {
		pendingAssetToClaim.Multiparty = 1
	}
	if wallet.HasViewOnlyAccess(gc) {
		pendingAssetToClaim.SignatureRequired = 1
	}
	var err error
	if wallet.Signer != signerUser.PrimarySigner && pendingAssetToClaim.Multiparty == 0 {
		return pendingAssetToClaim, false, &tErrors.CustomError{Param: "publicKey", Err: "error-wallet-not-managed-by-user", ErrMessage: "You do not have permission to operate on this wallet", Code: http.StatusBadRequest}

	}
	//validators
	{

		err = validators.ValidateAssetCodeFormat(pendingAssetToClaim.AssetCode)

		if err != nil {
			return pendingAssetToClaim, false, err
		}

		err = validators.ValidateAddressFormat(pendingAssetToClaim.AssetIssuer)

		if err != nil {
			return pendingAssetToClaim, false, err
		}
	}

	horizonClient := network.GetBlockchainClient()

	xdrBase64, err := generateRejectPendingAssetXdr(wallet, pendingAssetToClaim, gc)

	if err != nil {
		return pendingAssetToClaim, false, err
	}

	oldTxn := pendingAssetToClaim.Transaction

	pendingAssetToClaim.NetworkPassPhrase = network.GetBlockchainNetworkPassPhrase()

	pendingAssetToClaim.Transaction = xdrBase64

	if len(pendingAssetToClaim.TransactionSignature) == 0 && pendingAssetToClaim.Commit == 0 {
		//no signature
		return pendingAssetToClaim, false, err

	}

	//there was a signature... let's submit

	if oldTxn != xdrBase64 && pendingAssetToClaim.Commit == 0 {
		return pendingAssetToClaim, false, &tErrors.CustomError{
			Param:      "transaction",
			Err:        "transaction mismatch",
			ErrMessage: "transaction mismatch, please try again",
			Code:       404,
		}
	}
	if (pendingAssetToClaim.Commit == 0 && wallet.SharedAccessEnabled == 1 && wallet.NumberOfApprovalsNeeded == 0) || (wallet.SharedAccessEnabled == 0 && pendingAssetToClaim.Commit == 0) {
		txnID, err := network.SubmitXdrWithSignature(horizonClient, wallet.Signer, xdrBase64, pendingAssetToClaim.TransactionSignature)

		if err != nil {
			log.Printf("[ClaimPendingAsset]error submitting txn: %v\n", err)
			return pendingAssetToClaim, false, &tErrors.ErrorTemporaryServerError{}
		}

		pendingAssetToClaim.TransactionID = txnID

		return pendingAssetToClaim, true, nil
	}

	if pendingAssetToClaim.Multiparty == 1 {
		pendingAssetToClaim.TransactionID = "PENDING_AUTH"
		log.Printf("[ClaimPendingAsset]shared access with approver permission enabled for %v \n", wallet.Alias)
		id := uuid.NewString()
		assetOfPayment := os.Getenv("NATIVE_ASSET_CODE")
		if len(pendingAssetToClaim.AssetIssuer) == 42 {
			assetOfPayment = fmt.Sprintf("%v:%v...%v", pendingAssetToClaim.AssetCode, pendingAssetToClaim.AssetIssuer[0:4], pendingAssetToClaim.AssetIssuer[51:55])
		}
		description := fmt.Sprintf("Reject pending balance for %v.\nMessages: %v", assetOfPayment, pendingAssetToClaim.Messages)
		transactionByte, _ := json.Marshal(*pendingAssetToClaim)
		transactionStr := string(transactionByte)
		pendingAuth := userModels.PendingAuth{
			ID:                     id,
			Initiator:              signerUser.Username,
			InitiatorSignerAddress: signerUser.PrimarySigner,
			WalletAddress:          wallet.ID,
			TransactionType:        "REJECT PENDING ASSET",
			Description:            description,
			TransactionSource:      pendingAssetToClaim.TransactionSource,
			ApprovalsNeeded:        wallet.NumberOfApprovalsNeeded,
			TransactionXdr:         xdrBase64,
			TransactionInfoStr:     &transactionStr,
		}
		//save and commit this to database
		e := gc.DB.Omit(clause.Associations).Create(&pendingAuth).Error
		if e != nil {
			log.Printf("[ClaimPendingAsset] Error saving payment txn [%+v] transaction on pending auth table: %s\n", pendingAuth, e.Error())
			err = &tErrors.ErrorTemporaryServerError{}
			return pendingAssetToClaim, false, err
		}
		pendingAssetToClaim.ReturnedDescription = description
		return pendingAssetToClaim, true, nil

	}

	return pendingAssetToClaim, false, &tErrors.ErrorTemporaryServerError{}
}

func generateClaimPendingAssetXdr(wallet *userModels.UserWallet, pendingAssetToClaim *userModels.PendingAssetToClaim, gc *sharedconfig.GlobalConfig) (string, error) {
	var tokenizedAssetIssuerMustSign bool
	if len(pendingAssetToClaim.AssetIssuer) != 42 {
		return "", &tErrors.CustomError{
			Param:      "assetIssuer",
			Err:        "error-missing-parameter",
			ErrMessage: "Asset issuer is invalid",
			Code:       http.StatusBadRequest,
		}
	}
	tempKeyPair, err := network.TempAccountKeypair(wallet.ID)
	log.Printf("[generatePendingAssetXdr]tempKey: %v, main key: %v, alias: %v\n", tempKeyPair.Address(), wallet.ID, wallet.Alias)

	if err != nil {
		return "", err
	}

	//asset to Claim

	var asset basetxn.Asset = nil

	// asset = basetxn.NativeAsset{}

	if len(pendingAssetToClaim.AssetIssuer) > 0 {
		asset = basetxn.CreditAsset{Code: pendingAssetToClaim.AssetCode, Issuer: pendingAssetToClaim.AssetIssuer}
	}

	tempAccountExist, tempAccountTrustsAsset, nativeAccountBalance, customAccountBalance, tempAccount, err := network.BlockchainAccountProperties(gc.BantuExpansionClient, tempKeyPair.Address(), asset)

	if err != nil {
		return "", err
	}

	if !tempAccountExist {
		return "", &tErrors.ErrorAssetNotClaimable{}
	}

	if !tempAccountTrustsAsset {
		return "", &tErrors.ErrorAssetNotClaimable{}
	}
	chanAccount := <-gc.ChannelAccounts
	defer func(c *evmkeypair.Full) {
		gc.ChannelAccounts <- c
	}(chanAccount)
	// paymentInfo.Messages = messages
	_, _, _, _, chanSourceAccount, _ := network.BlockchainAccountProperties(gc.BantuExpansionClient, chanAccount.Address(), basetxn.NativeAsset{})

	//source account details

	_, sourceAccountTrustsAsset, _, _, sourceAccount, _ := network.BlockchainAccountProperties(gc.BantuExpansionClient, wallet.ID, asset)

	var ops []basetxn.Operation = make([]basetxn.Operation, 0)

	if nativeAccountBalance.GreaterThan(decimal.Zero) {
		ops = append(ops, &basetxn.Payment{
			Destination:   wallet.ID,
			Amount:        nativeAccountBalance.Truncate(7).String(),
			Asset:         basetxn.NativeAsset{},
			SourceAccount: tempAccount.Address,
		})
	}

	if !sourceAccountTrustsAsset {
		ops = append(ops, &basetxn.ChangeTrust{
			Line:          asset,
			Limit:         "900000000000",
			SourceAccount: sourceAccount.Address,
		})
	}

	if gc.IsValidTokenizedAsset(asset.GetCode()) {
		//check if it is a tokenized asset
		// allow trust from issuer to destination wallet
		ops = append(ops, &basetxn.SetTrustLineFlags{
			Trustor:       wallet.ID,
			Asset:         basetxn.CreditAsset{Code: asset.GetCode(), Issuer: asset.GetIssuer()},
			SetFlags:      []basetxn.TrustLineFlag{basetxn.TrustLineAuthorized},
			SourceAccount: asset.GetIssuer(),
		})
		tokenizedAssetIssuerMustSign = true
	}

	if customAccountBalance.GreaterThan(decimal.Zero) {
		ops = append(ops, &basetxn.Payment{
			Destination:   wallet.ID,
			Amount:        customAccountBalance.Truncate(7).String(),
			Asset:         asset,
			SourceAccount: tempAccount.Address,
		})
	}

	if len(ops) == 0 {

		return "", &tErrors.ErrorAssetNotClaimable{}
	}

	var tx *basetxn.Transaction
	// Construct the transaction that holds the operations to execute on the network
	if pendingAssetToClaim.Multiparty == 1 {
		pendingAssetToClaim.TransactionSource = chanSourceAccount.Address
		tx, err = basetxn.NewTransaction(
			basetxn.TransactionParams{
				SourceAccount:        chanSourceAccount.Address,
				IncrementSequenceNum: true,
				Operations:           ops,
				BaseFee:              2000,
				Memo:                 "claim-asset",
			},
		)
	} else {
		tx, err = basetxn.NewTransaction(
			basetxn.TransactionParams{
				SourceAccount:        tempAccount.Address,
				IncrementSequenceNum: true,
				Operations:           ops,
				BaseFee:              2000,
				Memo:                 "claim-asset",
			},
		)
	}

	if err != nil {
		log.Println("[generatePendingAssetXdr]error constructing transaction ", err)
		return "", &tErrors.ErrorTemporaryServerError{}
	}

	if pendingAssetToClaim.Multiparty == 1 {

		tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), chanAccount)

		if err != nil {
			log.Println("[generatePendingAssetXdr] error signing transaction with channelAccount key ", err)
			return "", &tErrors.ErrorTemporaryServerError{}
		}
	}
	if tokenizedAssetIssuerMustSign {
		log.Println("[generateTrustAssetXdr] <<<<<<<<<<<<<<<<<<<<<<<<<<<< signing transaction with issuer key>>>>>>>>>>>>>>>>>>>>>>>>")
		//get atprofile
		var tokenizationIssuerProfileWallet string

		if len(os.Getenv("TOKENIZATION_ISSUING_PROFILE_WALLET")) > 1 {
			tokenizationIssuerProfileWallet = strings.TrimSpace(os.Getenv("TOKENIZATION_ISSUING_PROFILE_WALLET"))
		}

		tokenizationIssuerProfileWalletKP := evmkeypair.MustParseFull(tokenizationIssuerProfileWallet)

		tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), tokenizationIssuerProfileWalletKP)
		if err != nil {
			log.Println("[generateTrustAssetXdr] error signing transaction with issuer key to authorize trustline", err)
			return "", &tErrors.ErrorTemporaryServerError{}
		}
	}
	xdrBase64, err := tx.Base64()

	if err != nil {
		return "", err
	}

	return xdrBase64, nil

}

func generateRejectPendingAssetXdr(wallet *userModels.UserWallet, pendingAssetToClaim *userModels.PendingAssetToClaim, gc *sharedconfig.GlobalConfig) (string, error) {
	if len(pendingAssetToClaim.AssetIssuer) != 42 {
		return "", &tErrors.CustomError{
			Param:      "assetIssuer",
			Err:        "error-missing-parameter",
			ErrMessage: "Asset issuer is invalid",
			Code:       http.StatusBadRequest,
		}
	}
	tempKeyPair, err := network.TempAccountKeypair(wallet.ID)
	log.Printf("[generateRejectPendingAssetXdr]tempKey: %v, main key: %v, alias: %v\n", tempKeyPair.Address(), wallet.ID, wallet.Alias)

	if err != nil {
		return "", err
	}
	originAddress := userBc.BlockchainAssetLastPaymentSource(tempKeyPair.Address(), pendingAssetToClaim.AssetCode, pendingAssetToClaim.AssetIssuer, gc)

	if len(originAddress) == 0 {
		err = &tErrors.CustomError{
			Param:      "assetCode",
			Err:        "error-could not get payment source",
			ErrMessage: "Unable to get payment source for rejection.",
		}
		return "", err
	}

	//asset to Claim

	var asset basetxn.Asset = nil

	asset = basetxn.NativeAsset{}

	if len(pendingAssetToClaim.AssetIssuer) > 0 {
		asset = basetxn.CreditAsset{Code: pendingAssetToClaim.AssetCode, Issuer: pendingAssetToClaim.AssetIssuer}
	}

	tempAccountExist, tempAccountTrustsAsset, _, customAccountBalance, tempAccount, err := network.BlockchainAccountProperties(gc.BantuExpansionClient, tempKeyPair.Address(), asset)

	if err != nil {
		return "", err
	}

	if !tempAccountExist {
		return "", &tErrors.ErrorAssetNotClaimable{}
	}

	if !tempAccountTrustsAsset {
		return "", &tErrors.ErrorAssetNotClaimable{}
	}
	chanAccount := <-gc.ChannelAccounts
	defer func(c *evmkeypair.Full) {
		gc.ChannelAccounts <- c
	}(chanAccount)
	// paymentInfo.Messages = messages
	_, _, _, _, chanSourceAccount, _ := network.BlockchainAccountProperties(gc.BantuExpansionClient, chanAccount.Address(), basetxn.NativeAsset{})

	//source account details

	// sourceAccountExists, sourceAccountTrustsAsset, _, _, sourceAccount, _ := network.BlockchainAccountProperties(gc.BantuExpansionClient, wallet.ID, asset)

	var ops []basetxn.Operation = make([]basetxn.Operation, 0)

	if customAccountBalance.GreaterThan(decimal.Zero) {

		ops = append(ops, &basetxn.Payment{
			Destination:   originAddress,
			Amount:        customAccountBalance.Truncate(7).String(),
			Asset:         asset,
			SourceAccount: tempAccount.Address,
		})
		//remove trustline
		ops = append(ops, &basetxn.ChangeTrust{
			Line:          asset,
			Limit:         "0",
			SourceAccount: tempAccount.Address,
		})
	}

	if len(ops) == 0 {

		return "", &tErrors.ErrorAssetNotClaimable{}
	}

	var tx *basetxn.Transaction
	// Construct the transaction that holds the operations to execute on the network
	if pendingAssetToClaim.Multiparty == 1 {
		pendingAssetToClaim.TransactionSource = chanSourceAccount.Address
		tx, err = basetxn.NewTransaction(
			basetxn.TransactionParams{
				SourceAccount:        chanSourceAccount.Address,
				IncrementSequenceNum: true,
				Operations:           ops,
				BaseFee:              2000,
				Memo:                 "reject-asset",
			},
		)
	} else {
		tx, err = basetxn.NewTransaction(
			basetxn.TransactionParams{
				SourceAccount:        tempAccount.Address,
				IncrementSequenceNum: true,
				Operations:           ops,
				BaseFee:              2000,
				Memo:                 "reject-asset",
			},
		)
	}

	if err != nil {
		log.Println("[generateRejectPendingAssetXdr]error constructing transaction ", err)
		return "", &tErrors.ErrorTemporaryServerError{}
	}

	if pendingAssetToClaim.Multiparty == 1 {

		tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), chanAccount)

		if err != nil {
			log.Println("[generateRejectPendingAssetXdr] error signing transaction with channelAccount key ", err)
			return "", &tErrors.ErrorTemporaryServerError{}
		}
	}

	xdrBase64, err := tx.Base64()

	if err != nil {
		return "", err
	}

	return xdrBase64, nil

}

func generateTrustAssetXdr(wallet *userModels.UserWallet, trustLineInfo *userModels.Trustline, gc *sharedconfig.GlobalConfig) (txnBase64 string, err error) {
	var tokenizedAssetIssuerMustSign bool
	if len(trustLineInfo.AssetIssuer) != 42 {
		return "", &tErrors.CustomError{
			Param:      "assetIssuer",
			Err:        "error-missing-parameter",
			ErrMessage: "Asset issuer is invalid",
			Code:       http.StatusBadRequest,
		}
	}
	assetIssuer := trustLineInfo.AssetIssuer
	assetCode := trustLineInfo.AssetCode
	minBalance := decimal.RequireFromString(os.Getenv("STANDARD_WALLET_MINIMUM_BALANCE"))
	trustLineInfo.Messages = make([]string, 0)
	if len(assetCode) == 0 || len(assetCode) > 12 || len(assetIssuer) != 42 {
		return "", &tErrors.CustomError{Param: "assetCode", Err: "error-invalid-asset", ErrMessage: "Asset Supplied is invalid.", Code: http.StatusBadRequest}
	}

	asset := basetxn.CreditAsset{Code: assetCode, Issuer: assetIssuer}
	chanAccount := <-gc.ChannelAccounts
	defer func(c *evmkeypair.Full) {
		gc.ChannelAccounts <- c
	}(chanAccount)

	_, _, _, _, chanSourceAccount, _ := network.BlockchainAccountProperties(gc.BantuExpansionClient, chanAccount.Address(), basetxn.NativeAsset{})

	_, sourceAccountTrustsAsset, nativeAccountBalance, _, sourceAccount, _ := network.BlockchainAccountProperties(gc.BantuExpansionClient, wallet.ID, asset)

	if nativeAccountBalance.LessThan(minBalance) {
		return "", &tErrors.CustomError{Param: "publicKey", Err: "error-wallet-underfunded", ErrMessage: fmt.Sprintf("The Wallet is currently underfunded. Please maintain min %v %v balance before you can perform this task", minBalance.String(), os.Getenv("NATIVE_ASSET_CODE")), Code: http.StatusBadRequest}

	}

	if sourceAccountTrustsAsset {

		return "", &tErrors.CustomError{Param: "publicKey", Err: "error-asset-already-trusted", ErrMessage: fmt.Sprintf("The wallet already accepted this %v asset before.", assetCode), Code: http.StatusBadRequest}

	}

	var ops []basetxn.Operation = make([]basetxn.Operation, 0)
	ops = append(ops, &basetxn.ChangeTrust{
		Line:          asset,
		Limit:         "900000000000",
		SourceAccount: wallet.ID,
	})

	if gc.IsValidTokenizedAsset(asset.GetCode()) {
		//check if it is a tokenized asset
		// allow trust from issuer to destination wallet
		//check if it is the asset tokenization profile that owns the wallet
		// 		walletOwner, _:=wallet.GetWalletOwner(gc.DB,gc)
		// 		if walletOwner.Username==strings.TrimSpace(os.Getenv("TOKENIZATION_ISSUING_PROFILE")){
		// //request if from a tokenization profile give full rights
		// 		ops = append(ops, &basetxn.SetTrustLineFlags{
		// 			Trustor:       wallet.ID,
		// 			Asset:         basetxn.CreditAsset{Code: asset.GetCode(), Issuer: asset.GetIssuer()},
		// 			SetFlags:      []basetxn.TrustLineFlag{basetxn.TrustLineAuthorized},
		// 			SourceAccount: asset.GetIssuer(),
		// 		})

		// 		}else{
		// 			//it is from third party. remove market maker right.
		// 		ops = append(ops, &basetxn.SetTrustLineFlags{
		// 			Trustor:       wallet.ID,
		// 			Asset:         basetxn.CreditAsset{Code: asset.GetCode(), Issuer: asset.GetIssuer()},
		// 			SetFlags:      []basetxn.TrustLineFlag{basetxn.TrustLineAuthorized, txnbuild.trustlin},
		// 			SourceAccount: asset.GetIssuer(),
		// 		})
		// 		}

		ops = append(ops, &basetxn.SetTrustLineFlags{
			Trustor:       wallet.ID,
			Asset:         basetxn.CreditAsset{Code: asset.GetCode(), Issuer: asset.GetIssuer()},
			SetFlags:      []basetxn.TrustLineFlag{basetxn.TrustLineAuthorized},
			SourceAccount: asset.GetIssuer(),
		})
		tokenizedAssetIssuerMustSign = true
	}
	// Construct the transaction that holds the operations to execute on the network

	var tx *basetxn.Transaction
	// Construct the transaction that holds the operations to execute on the network
	if trustLineInfo.Multiparty == 1 {
		trustLineInfo.TransactionSource = chanSourceAccount.Address
		tx, err = basetxn.NewTransaction(
			basetxn.TransactionParams{
				SourceAccount:        chanSourceAccount.Address,
				IncrementSequenceNum: true,
				Operations:           ops,
				BaseFee:              2000,
				Memo:                 "opt-in-" + assetCode,
			},
		)
	} else {
		tx, err = basetxn.NewTransaction(
			basetxn.TransactionParams{
				SourceAccount:        sourceAccount.Address,
				IncrementSequenceNum: true,
				Operations:           ops,
				BaseFee:              2000,
				Memo:                 "opt-in-" + assetCode,
			},
		)
	}
	if err != nil {
		log.Println("[generateTrustAssetXdr]error constructing transaction", err)
		return "", &tErrors.ErrorTemporaryServerError{}
	}

	if trustLineInfo.Multiparty == 1 {

		tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), chanAccount)

		if err != nil {
			log.Println("[generateTrustAssetXdr] error signing transaction with channelAccount key ", err)
			return "", &tErrors.ErrorTemporaryServerError{}
		}
	}
	if tokenizedAssetIssuerMustSign {
		log.Println("[generateTrustAssetXdr] <<<<<<<<<<<<<<<<<<<<<<<<<<<< signing transaction with issuer key>>>>>>>>>>>>>>>>>>>>>>>>")
		//get atprofile
		var tokenizationIssuerProfileWallet string

		if len(os.Getenv("TOKENIZATION_ISSUING_PROFILE_WALLET")) > 1 {
			tokenizationIssuerProfileWallet = strings.TrimSpace(os.Getenv("TOKENIZATION_ISSUING_PROFILE_WALLET"))
		}

		tokenizationIssuerProfileWalletKP := evmkeypair.MustParseFull(tokenizationIssuerProfileWallet)

		tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), tokenizationIssuerProfileWalletKP)
		if err != nil {
			log.Println("[generateTrustAssetXdr] error signing transaction with issuer key to authorize trustline", err)
			return "", &tErrors.ErrorTemporaryServerError{}
		}
	}
	xdrBase64, err := tx.Base64()

	if err != nil {
		return "", err
	}

	return xdrBase64, nil

}

func generateRemoveTrustAssetXdr(wallet *userModels.UserWallet, trustLineInfo *userModels.Trustline, gc *sharedconfig.GlobalConfig) (txnBase64 string, err error) {

	assetIssuer := trustLineInfo.AssetIssuer
	assetCode := trustLineInfo.AssetCode
	minBalance := decimal.RequireFromString(os.Getenv("STANDARD_WALLET_MINIMUM_BALANCE"))
	if len(assetIssuer) != 42 {
		return "", &tErrors.CustomError{
			Param:      "assetIssuer",
			Err:        "error-missing-parameter",
			ErrMessage: "Asset issuer is invalid",
			Code:       http.StatusBadRequest,
		}
	}
	if len(assetCode) == 0 || len(assetCode) > 12 || len(assetIssuer) != 42 {
		return "", &tErrors.CustomError{Param: "assetCode", Err: "error-invalid-asset", ErrMessage: "Asset Supplied is invalid.", Code: http.StatusBadRequest}
	}

	asset := basetxn.CreditAsset{Code: assetCode, Issuer: assetIssuer}
	chanAccount := <-gc.ChannelAccounts
	defer func(c *evmkeypair.Full) {
		gc.ChannelAccounts <- c
	}(chanAccount)

	_, _, _, _, chanSourceAccount, _ := network.BlockchainAccountProperties(gc.BantuExpansionClient, chanAccount.Address(), basetxn.NativeAsset{})

	_, sourceAccountTrustsAsset, nativeAccountBalance, assetBalance, sourceAccount, _ := network.BlockchainAccountProperties(gc.BantuExpansionClient, wallet.ID, asset)

	if nativeAccountBalance.LessThan(minBalance) {
		return "", &tErrors.CustomError{Param: "publicKey", Err: "error-wallet-underfunded", ErrMessage: fmt.Sprintf("The Wallet is currently underfunded. Please maintain min %v %v balance before you can perform this task", minBalance.String(), os.Getenv("NATIVE_ASSET_CODE")), Code: http.StatusBadRequest}

	}

	if assetBalance.GreaterThan(decimal.Zero) {
		return "", &tErrors.CustomError{Param: "publicKey", Err: "error-wallet-has-asset-balance", ErrMessage: fmt.Sprintf("The Wallet is currently has %v %v balance. Please transfer all of them and maintain 0 %v balance before you can perform this task", assetBalance.Truncate(7), assetCode, assetCode), Code: http.StatusBadRequest}

	}

	if !sourceAccountTrustsAsset {

		return "", &tErrors.CustomError{Param: "publicKey", Err: "error-asset-not-currently-trusted", ErrMessage: fmt.Sprintf("The wallet currently does not accept this %v asset. No need for this operation.", assetCode), Code: http.StatusBadRequest}

	}

	var ops []basetxn.Operation = make([]basetxn.Operation, 0)
	ops = append(ops, &basetxn.ChangeTrust{
		Line:          asset,
		Limit:         "0",
		SourceAccount: wallet.ID,
	})

	// //service fee
	// serviceFee, e := decimal.NewFromString(os.Getenv("SHARED_ACCESS_FEE_AMOUNT"))
	// if e != nil {
	// 	serviceFee = decimal.Zero
	// }
	// if serviceFee.IsPositive() {
	// 	if trustLineInfo.Multiparty == 1 {
	// 		//process service fee
	// 		if len(os.Getenv("SHARED_ACCESS_FEE_ASSET_ISSUER")) != 42 {
	// 			ops = append(ops, &basetxn.Payment{
	// 				Destination:   os.Getenv("SHARED_ACCESS_FEE_ADDRESS"),
	// 				Amount:        os.Getenv("SHARED_ACCESS_FEE_AMOUNT"),
	// 				SourceAccount: wallet.ID,
	// 				Asset:         basetxn.NativeAsset{},
	// 			})
	// 			trustLineInfo.Messages = append(trustLineInfo.Messages, fmt.Sprintf("%v %v will be deducted from wallet %v as service fee.", os.Getenv("SHARED_ACCESS_FEE_AMOUNT"), os.Getenv("NATIVE_ASSET_CODE"), wallet.Alias))

	// 		} else {
	// 			ops = append(ops, &basetxn.Payment{
	// 				Destination:   os.Getenv("SHARED_ACCESS_FEE_ADDRESS"),
	// 				Amount:        os.Getenv("SHARED_ACCESS_FEE_AMOUNT"),
	// 				SourceAccount: wallet.ID,
	// 				Asset:         basetxn.CreditAsset{Code: os.Getenv("SHARED_ACCESS_FEE_ASSET_CODE"), Issuer: os.Getenv("SHARED_ACCESS_FEE_ASSET_ISSUER")},
	// 			})
	// 			trustLineInfo.Messages = append(trustLineInfo.Messages, fmt.Sprintf("%v %v will be deducted from wallet %v as service fee.", os.Getenv("SHARED_ACCESS_FEE_AMOUNT"), os.Getenv("SHARED_ACCESS_FEE_ASSET_CODE"), wallet.Alias))

	// 		}

	// 	}
	// }

	var tx *basetxn.Transaction
	// Construct the transaction that holds the operations to execute on the network
	if trustLineInfo.Multiparty == 1 {
		trustLineInfo.TransactionSource = chanSourceAccount.Address
		tx, err = basetxn.NewTransaction(
			basetxn.TransactionParams{
				SourceAccount:        chanSourceAccount.Address,
				IncrementSequenceNum: true,
				Operations:           ops,
				BaseFee:              2000,
				Memo:                 "opt-out-" + assetCode,
			},
		)
	} else {
		tx, err = basetxn.NewTransaction(
			basetxn.TransactionParams{
				SourceAccount:        sourceAccount.Address,
				IncrementSequenceNum: true,
				Operations:           ops,
				BaseFee:              2000,
				Memo:                 "opt-out-" + assetCode,
			},
		)
	}
	if err != nil {
		log.Println("[generateRemoveTrustAssetXdr]error constructing transaction", err)
		return "", &tErrors.ErrorTemporaryServerError{}
	}

	if trustLineInfo.Multiparty == 1 {

		tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), chanAccount)

		if err != nil {
			log.Println("[generateRemoveTrustAssetXdr] error signing transaction with channelAccount key ", err)
			return "", &tErrors.ErrorTemporaryServerError{}
		}
	}

	xdrBase64, err := tx.Base64()

	if err != nil {
		return "", err
	}

	return xdrBase64, nil

}

func TrustAsset(signerUser *userModels.User, wallet *userModels.UserWallet, trustLineInfo *userModels.Trustline, gc *sharedconfig.GlobalConfig) (*userModels.Trustline, error) {
	if gc.IsValidTokenizedAsset(trustLineInfo.AssetCode) {
		t := gc.GetTokenizedAssetByCode(trustLineInfo.AssetCode)
		if t.AssetTokenizationStatus < 5 {
			return trustLineInfo, &tErrors.CustomError{
				Param:      "assetIssuer",
				Err:        "error-asset-not-yet-available-for-sale",
				ErrMessage: "This tokenized Asset is not yet available for sale. Add/Remove is not allowed at this time.",
				Code:       http.StatusForbidden,
			}
		}
	}

	trustLineInfo.NetworkPassPhrase = gc.BantuNetworkPassphrase
	if wallet.NumberOfApprovalsNeeded > 0 && wallet.SharedAccessEnabled == 1 {
		trustLineInfo.Multiparty = 1
	}
	if wallet.HasViewOnlyAccess(gc) {
		trustLineInfo.SignatureRequired = 1
	}
	if wallet.Signer != signerUser.PrimarySigner && trustLineInfo.Multiparty == 0 {
		return trustLineInfo, &tErrors.CustomError{Param: "publicKey", Err: "error-wallet-not-managed-by-user", ErrMessage: "You do not have permission to operate on this wallet", Code: http.StatusBadRequest}

	}

	xdrBase64, err := generateTrustAssetXdr(wallet, trustLineInfo, gc)
	if err != nil {
		return trustLineInfo, err
	}
	oldTransaction := trustLineInfo.Transaction
	trustLineInfo.Transaction = xdrBase64
	if len(trustLineInfo.TransactionSignature) == 0 && trustLineInfo.Commit == 0 {
		//needs to be signed first
		return trustLineInfo, nil
	}
	if oldTransaction != xdrBase64 && trustLineInfo.Commit == 0 {
		return trustLineInfo, &tErrors.CustomError{
			Param:      "transaction",
			Err:        "transaction mismatch",
			ErrMessage: "transaction mismatch, please try again",
			Code:       http.StatusBadRequest,
		}
	}

	if (trustLineInfo.Commit == 0 && wallet.SharedAccessEnabled == 1 && wallet.NumberOfApprovalsNeeded == 0) || (wallet.SharedAccessEnabled == 0 && trustLineInfo.Commit == 0) {

		//submit to network
		txnHash, err := network.SubmitXdrWithSignature(gc.BantuExpansionClient, signerUser.PrimarySigner, xdrBase64, trustLineInfo.TransactionSignature)
		if err != nil {
			return trustLineInfo, err
		}
		trustLineInfo.TransactionID = txnHash
		return trustLineInfo, nil
	}

	if trustLineInfo.Multiparty == 1 {
		trustLineInfo.TransactionID = "PENDING_AUTH"
		log.Printf("[TrustAsset]shared access with approver permission enabled for %v \n", wallet.Alias)
		id := uuid.NewString()
		assetOfPayment := os.Getenv("NATIVE_ASSET_CODE")
		if len(trustLineInfo.AssetIssuer) == 42 {
			assetOfPayment = fmt.Sprintf("%v:%v...%v", trustLineInfo.AssetCode, trustLineInfo.AssetIssuer[0:4], trustLineInfo.AssetIssuer[51:55])
		}
		description := fmt.Sprintf("Opt in asset %v.\nMessages: %v", assetOfPayment, trustLineInfo.Messages)
		transactionByte, _ := json.Marshal(*trustLineInfo)
		transactionStr := string(transactionByte)
		pendingAuth := userModels.PendingAuth{
			ID:                     id,
			Initiator:              signerUser.Username,
			InitiatorSignerAddress: signerUser.PrimarySigner,
			WalletAddress:          wallet.ID,
			TransactionType:        "OPT IN ASSET",
			Description:            description,
			TransactionSource:      trustLineInfo.TransactionSource,
			ApprovalsNeeded:        wallet.NumberOfApprovalsNeeded,
			TransactionXdr:         xdrBase64,
			TransactionInfoStr:     &transactionStr,
		}
		//save and commit this to database

		//check for duplicate
		if CheckDuplicatePendingApproval(pendingAuth.WalletAddress, pendingAuth.TransactionType, pendingAuth.Description, gc.DB) {
			//duplicate exists. resist the duplicate
			return trustLineInfo, &tErrors.CustomError{
				Param:      "transaction",
				Err:        "error-duplicate-operation-exists",
				ErrMessage: "There is an existing duplicate transaction pending approval. Please Approve or reject that one before continuing with this.",
				Code:       http.StatusBadRequest,
			}
		}

		e := gc.DB.Omit(clause.Associations).Create(&pendingAuth).Error
		if e != nil {
			log.Printf("[TrustAsset] Error saving opt in txn [%+v] transaction on pending auth table: %s\n", pendingAuth, e.Error())
			err = &tErrors.ErrorTemporaryServerError{}
			return trustLineInfo, err
		}
		trustLineInfo.ReturnedDescription = description
		return trustLineInfo, nil

	}
	// did not match any of the conditions
	return trustLineInfo, &tErrors.ErrorTemporaryServerError{}
}

func RemoveAssetTrust(signerUser *userModels.User, wallet *userModels.UserWallet, trustLineInfo *userModels.Trustline, gc *sharedconfig.GlobalConfig) (*userModels.Trustline, error) {
	if gc.IsValidTokenizedAsset(trustLineInfo.AssetCode) {
		t := gc.GetTokenizedAssetByCode(trustLineInfo.AssetCode)
		if t.AssetTokenizationStatus < 5 {
			return trustLineInfo, &tErrors.CustomError{
				Param:      "assetIssuer",
				Err:        "error-asset-not-yet-available-for-sale",
				ErrMessage: "This tokenized Asset is not yet available for sale. Add/Remove is not allowed at this time.",
				Code:       http.StatusForbidden,
			}
		}
	}

	trustLineInfo.NetworkPassPhrase = gc.BantuNetworkPassphrase
	if wallet.NumberOfApprovalsNeeded > 0 && wallet.SharedAccessEnabled == 1 {
		trustLineInfo.Multiparty = 1
	}
	if wallet.HasViewOnlyAccess(gc) {
		trustLineInfo.SignatureRequired = 1
	}
	if wallet.Signer != signerUser.PrimarySigner && trustLineInfo.Multiparty == 0 {
		return trustLineInfo, &tErrors.CustomError{Param: "publicKey", Err: "error-wallet-not-managed-by-user", ErrMessage: "You do not have permission to operate on this wallet", Code: http.StatusBadRequest}

	}

	xdrBase64, err := generateRemoveTrustAssetXdr(wallet, trustLineInfo, gc)
	if err != nil {
		return trustLineInfo, err
	}
	oldTransaction := trustLineInfo.Transaction
	trustLineInfo.Transaction = xdrBase64
	if len(trustLineInfo.TransactionSignature) == 0 && trustLineInfo.Commit == 0 {
		//needs to be signed first
		return trustLineInfo, nil
	}
	if oldTransaction != xdrBase64 && trustLineInfo.Commit == 0 {
		return trustLineInfo, &tErrors.CustomError{
			Param:      "transaction",
			Err:        "transaction mismatch",
			ErrMessage: "transaction mismatch, please try again",
			Code:       http.StatusBadRequest,
		}
	}

	if (trustLineInfo.Commit == 0 && wallet.SharedAccessEnabled == 1 && wallet.NumberOfApprovalsNeeded == 0) || (wallet.SharedAccessEnabled == 0 && trustLineInfo.Commit == 0) {

		//submit to network
		txnHash, err := network.SubmitXdrWithSignature(gc.BantuExpansionClient, signerUser.PrimarySigner, xdrBase64, trustLineInfo.TransactionSignature)
		if err != nil {
			return trustLineInfo, err
		}
		trustLineInfo.TransactionID = txnHash
		return trustLineInfo, nil
	}

	if trustLineInfo.Multiparty == 1 {
		trustLineInfo.TransactionID = "PENDING_AUTH"
		log.Printf("[RemoveAssetTrust]shared access with approver permission enabled for %v \n", wallet.Alias)
		id := uuid.NewString()
		assetOfPayment := os.Getenv("NATIVE_ASSET_CODE")
		if len(trustLineInfo.AssetIssuer) == 42 {
			assetOfPayment = fmt.Sprintf("%v:%v...%v", trustLineInfo.AssetCode, trustLineInfo.AssetIssuer[0:4], trustLineInfo.AssetIssuer[51:55])
		}
		description := fmt.Sprintf("Opt out asset %v.\nMessages: %v", assetOfPayment, trustLineInfo.Messages)
		transactionByte, _ := json.Marshal(*trustLineInfo)
		transactionStr := string(transactionByte)
		pendingAuth := userModels.PendingAuth{
			ID:                     id,
			Initiator:              signerUser.Username,
			InitiatorSignerAddress: signerUser.PrimarySigner,
			WalletAddress:          wallet.ID,
			TransactionType:        "OPT OUT ASSET",
			Description:            description,
			TransactionSource:      trustLineInfo.TransactionSource,
			ApprovalsNeeded:        wallet.NumberOfApprovalsNeeded,
			TransactionXdr:         xdrBase64,
			TransactionInfoStr:     &transactionStr,
		}
		//save and commit this to database
		e := gc.DB.Omit(clause.Associations).Create(&pendingAuth).Error
		if e != nil {
			log.Printf("[RemoveAssetTrust] Error saving opt out txn [%+v] transaction on pending auth table: %s\n", pendingAuth, e.Error())
			err = &tErrors.ErrorTemporaryServerError{}
			return trustLineInfo, err
		}
		trustLineInfo.ReturnedDescription = description
		return trustLineInfo, nil

	}
	// did not match any of the conditions
	return trustLineInfo, &tErrors.ErrorTemporaryServerError{}
}

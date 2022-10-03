package users

import (
	"fmt"
	"log"
	"net/http"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/sharedconfig"
	"trovo-wallet-api/internal/validators"

	"github.com/shopspring/decimal"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/txnbuild"
	"gorm.io/gorm"
)

//ClaimPendingAsset claim pending assets
func ClaimPendingAsset(owner *userModels.User, wallet *userModels.UserWallet, pendingAssetToClaim *userModels.PendingAssetToClaim, db *gorm.DB) (*userModels.PendingAssetToClaim, bool, error) {

	var err error

	//validators
	{

		err = validators.ValidateAssetCodeFormat(pendingAssetToClaim.AssetCode)

		if err != nil {
			return pendingAssetToClaim, false, err
		}

		err = validators.ValidatePublicKeyFormat(pendingAssetToClaim.AssetIssuer)

		if err != nil {
			return pendingAssetToClaim, false, err
		}
	}

	horizonClient := network.GetBlockchainClient()

	xdrBase64, err := generateXdr(horizonClient, *owner, wallet, db, pendingAssetToClaim)

	if err != nil {
		return pendingAssetToClaim, false, err
	}

	oldTxn := pendingAssetToClaim.Transaction

	pendingAssetToClaim.NetworkPassPhrase = network.GetBlockchainNetworkPassPhrase()

	pendingAssetToClaim.Transaction = xdrBase64

	if len(pendingAssetToClaim.TransactionSignature) == 0 {
		//no signature
		return pendingAssetToClaim, false, err

	}

	//there was a signature... let's submit

	if oldTxn != xdrBase64 {
		return pendingAssetToClaim, false, &tErrors.CustomError{
			Param:      "transaction",
			Err:        "transaction mismatch",
			ErrMessage: "transaction mismatch, please try again",
			Code:       404,
		}
	}

	txnID, err := network.SubmitXdrWithSignature(horizonClient, wallet.Signer, xdrBase64, pendingAssetToClaim.TransactionSignature)

	if err != nil {
		log.Printf("[ClaimPendingAsset]error submitting txn: %v\n", err)
		return pendingAssetToClaim, false, &tErrors.ErrorTemporaryServerError{}
	}

	pendingAssetToClaim.TransactionID = txnID

	return pendingAssetToClaim, true, nil

}

func generateXdr(horizonClient *horizonclient.Client, owner userModels.User, wallet *userModels.UserWallet, db *gorm.DB, pendingAssetToClaim *userModels.PendingAssetToClaim) (string, error) {

	tempKeyPair, err := network.TempAccountKeypair(wallet.ID)
	log.Printf("[generateXdr]tempKey: %v, main key: %v, alias: %v\n", tempKeyPair.Address(), wallet.ID, wallet.Alias)

	if err != nil {
		return "", err
	}

	//asset to Claim

	var asset txnbuild.Asset = nil

	asset = txnbuild.NativeAsset{}

	if len(pendingAssetToClaim.AssetCode) > 0 {
		asset = txnbuild.CreditAsset{Code: pendingAssetToClaim.AssetCode, Issuer: pendingAssetToClaim.AssetIssuer}
	}

	tempAccountExist, tempAccountTrustsAsset, nativeAccountBalance, customAccountBalance, tempAccount, err := network.BlockchainAccountProperties(horizonClient, tempKeyPair.Address(), asset)

	if err != nil {
		return "", err
	}

	if !tempAccountExist {
		return "", &tErrors.ErrorAssetNotClaimable{}
	}

	if !tempAccountTrustsAsset {
		return "", &tErrors.ErrorAssetNotClaimable{}
	}

	//source account details

	sourceAccountExists, sourceAccountTrustsAsset, _, _, sourceAccount, _ := network.BlockchainAccountProperties(horizonClient, wallet.ID, asset)

	var ops []txnbuild.Operation = make([]txnbuild.Operation, 0)

	if nativeAccountBalance.GreaterThan(decimal.Zero) {
		if !sourceAccountExists {
			//todo: check if nativeAccount balance > 1
			ops = append(ops, &txnbuild.CreateAccount{
				Destination:   wallet.ID,
				Amount:        nativeAccountBalance.Truncate(7).String(),
				SourceAccount: tempAccount.AccountID,
			})
		} else {
			ops = append(ops, &txnbuild.Payment{
				Destination:   wallet.ID,
				Amount:        nativeAccountBalance.Truncate(7).String(),
				Asset:         txnbuild.NativeAsset{},
				SourceAccount: tempAccount.AccountID,
			})
		}
	}

	if !sourceAccountTrustsAsset {
		ops = append(ops, &txnbuild.ChangeTrust{
			Line:          txnbuild.ChangeTrustAssetWrapper{Asset: asset},
			Limit:         "900000000000",
			SourceAccount: sourceAccount.AccountID,
		})
	}

	if customAccountBalance.GreaterThan(decimal.Zero) {
		ops = append(ops, &txnbuild.Payment{
			Destination:   wallet.ID,
			Amount:        customAccountBalance.Truncate(7).String(),
			Asset:         asset,
			SourceAccount: tempAccount.AccountID,
		})
	}

	if len(ops) == 0 {

		return "", &tErrors.ErrorAssetNotClaimable{}
	}

	// Construct the transaction that holds the operations to execute on the network
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        tempAccount,
			IncrementSequenceNum: true,
			Operations:           ops,
			BaseFee:              txnbuild.MinBaseFee,
			Memo:                 txnbuild.MemoText("claim-asset"),
			Preconditions: txnbuild.Preconditions{
				TimeBounds: txnbuild.NewInfiniteTimeout(),
			},
		},
	)

	if err != nil {
		log.Println("[claim-asset generateXdr]error constructing transaction ", err)
		return "", &tErrors.ErrorTemporaryServerError{}
	}

	xdrBase64, err := tx.Base64()

	if err != nil {
		return "", err
	}

	return xdrBase64, nil

}

func generateTrustAssetXdr(wallet *userModels.UserWallet, assetCode, assetIssuer string, gc *sharedconfig.GlobalConfig) (txnBase64 string, err error) {
	if len(assetCode) == 0 || len(assetCode) > 12 || len(assetIssuer) != 56 {
		return "", &tErrors.CustomError{Param: "assetCode", Err: "error-invalid-asset", ErrMessage: "Asset Supplied is invalid.", Code: http.StatusBadRequest}
	}

	asset := txnbuild.CreditAsset{Code: assetCode, Issuer: assetIssuer}
	sourceAccountExists, sourceAccountTrustsAsset, nativeAccountBalance, _, sourceAccount, _ := network.BlockchainAccountProperties(gc.BantuExpansionClient, wallet.ID, asset)

	if !sourceAccountExists {
		return "", &tErrors.CustomError{Param: "publicKey", Err: "error-account-not-activated-on-blockchain", ErrMessage: "The Wallet public key is currently underfunded. Please send about 3XBN to it to activate it before you can perform this task", Code: http.StatusBadRequest}

	}
	if nativeAccountBalance.LessThan(decimal.NewFromFloat(2.8)) {
		return "", &tErrors.CustomError{Param: "publicKey", Err: "error-wallet-underfunded", ErrMessage: "The Wallet is currently underfunded. Please maintain about 3XBN balance before you can perform this task", Code: http.StatusBadRequest}

	}

	if sourceAccountTrustsAsset {

		return "", &tErrors.CustomError{Param: "publicKey", Err: "error-asset-already-trusted", ErrMessage: fmt.Sprintf("The wallet already accepted this %v asset before.", assetCode), Code: http.StatusBadRequest}

	}

	var ops []txnbuild.Operation = make([]txnbuild.Operation, 0)
	ops = append(ops, &txnbuild.ChangeTrust{
		Line:          txnbuild.ChangeTrustAssetWrapper{Asset: asset},
		Limit:         "900000000000",
		SourceAccount: wallet.ID,
	})

	// Construct the transaction that holds the operations to execute on the network
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        sourceAccount,
			IncrementSequenceNum: true,
			Operations:           ops,
			BaseFee:              txnbuild.MinBaseFee,
			Memo:                 txnbuild.MemoText("trust-" + assetCode),
			Preconditions: txnbuild.Preconditions{
				TimeBounds: txnbuild.NewInfiniteTimeout(),
			},
		},
	)

	if err != nil {
		log.Println("[generateTrustAssetXdr]error constructing transaction", err)
		return "", &tErrors.ErrorTemporaryServerError{}
	}

	xdrBase64, err := tx.Base64()

	if err != nil {
		return "", err
	}

	return xdrBase64, nil

}

func generateRemoveTrustAssetXdr(wallet *userModels.UserWallet, assetCode, assetIssuer string, gc *sharedconfig.GlobalConfig) (txnBase64 string, err error) {
	if len(assetCode) == 0 || len(assetCode) > 12 || len(assetIssuer) != 56 {
		return "", &tErrors.CustomError{Param: "assetCode", Err: "error-invalid-asset", ErrMessage: "Asset Supplied is invalid.", Code: http.StatusBadRequest}
	}

	asset := txnbuild.CreditAsset{Code: assetCode, Issuer: assetIssuer}
	sourceAccountExists, sourceAccountTrustsAsset, nativeAccountBalance, assetBalance, sourceAccount, _ := network.BlockchainAccountProperties(gc.BantuExpansionClient, wallet.ID, asset)

	if !sourceAccountExists {
		return "", &tErrors.CustomError{Param: "publicKey", Err: "error-account-not-activated-on-blockchain", ErrMessage: "The Wallet public key is currently underfunded. Please send about 3XBN to it to activate it before you can perform this task", Code: http.StatusBadRequest}

	}
	if nativeAccountBalance.LessThan(decimal.NewFromFloat(2.8)) {
		return "", &tErrors.CustomError{Param: "publicKey", Err: "error-wallet-underfunded", ErrMessage: "The Wallet is currently underfunded. Please maintain about 3XBN balance before you can perform this task", Code: http.StatusBadRequest}

	}

	if assetBalance.GreaterThan(decimal.Zero) {
		return "", &tErrors.CustomError{Param: "publicKey", Err: "error-wallet-has-asset-balance", ErrMessage: fmt.Sprintf("The Wallet is currently has %v %v balance. Please transfer all of them and maintain 0 %v balance before you can perform this task", assetBalance.Truncate(7), assetCode, assetCode), Code: http.StatusBadRequest}

	}

	if !sourceAccountTrustsAsset {

		return "", &tErrors.CustomError{Param: "publicKey", Err: "error-asset-not-currently-trusted", ErrMessage: fmt.Sprintf("The wallet currently does not accept this %v asset. No need for this operation.", assetCode), Code: http.StatusBadRequest}

	}

	var ops []txnbuild.Operation = make([]txnbuild.Operation, 0)
	ops = append(ops, &txnbuild.ChangeTrust{
		Line:          txnbuild.ChangeTrustAssetWrapper{Asset: asset},
		Limit:         "0",
		SourceAccount: wallet.ID,
	})

	// Construct the transaction that holds the operations to execute on the network
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        sourceAccount,
			IncrementSequenceNum: true,
			Operations:           ops,
			BaseFee:              txnbuild.MinBaseFee,
			Memo:                 txnbuild.MemoText("untrust-" + assetCode),
			Preconditions: txnbuild.Preconditions{
				TimeBounds: txnbuild.NewInfiniteTimeout(),
			},
		},
	)

	if err != nil {
		log.Println("[generateRemoveTrustAssetXdr]error constructing transaction", err)
		return "", &tErrors.ErrorTemporaryServerError{}
	}

	xdrBase64, err := tx.Base64()

	if err != nil {
		return "", err
	}

	return xdrBase64, nil

}

func TrustAsset(signerUser *userModels.User, wallet *userModels.UserWallet, trustLineInfo *userModels.Trustline, gc *sharedconfig.GlobalConfig) (*userModels.Trustline, error) {
	trustLineInfo.NetworkPassPhrase = gc.BantuNetworkPassphrase
	if wallet.SharedAccessEnabled == 0 {
		//managed access not enabled
		if wallet.Signer != signerUser.PrimarySigner {
			return trustLineInfo, &tErrors.CustomError{Param: "publicKey", Err: "error-wallet-not-managed-by-user", ErrMessage: "You do not have permission to operate on this wallet", Code: http.StatusBadRequest}

		}
		//generatexdr
		xdrBase64, err := generateTrustAssetXdr(wallet, trustLineInfo.AssetCode, trustLineInfo.AssetIssuer, gc)
		if err != nil {
			return trustLineInfo, err
		}
		oldTransaction := trustLineInfo.Transaction
		trustLineInfo.Transaction = xdrBase64
		if len(trustLineInfo.TransactionSignature) == 0 {
			//needs to be signed first
			return trustLineInfo, nil
		}
		if oldTransaction != xdrBase64 {
			return trustLineInfo, &tErrors.CustomError{
				Param:      "transaction",
				Err:        "transaction mismatch",
				ErrMessage: "transaction mismatch, please try again",
				Code:       http.StatusBadRequest,
			}
		}
		//submit to network
		txnHash, err := network.SubmitXdrWithSignature(gc.BantuExpansionClient, signerUser.PrimarySigner, xdrBase64, trustLineInfo.TransactionSignature)
		if err != nil {
			return trustLineInfo, err
		}
		trustLineInfo.TransactionID = txnHash
		return trustLineInfo, nil
	} else if wallet.SharedAccessEnabled == 1 {
		//perform managed access operation and save to table
		return trustLineInfo, nil
	}
	//did not match any of the conditions
	return trustLineInfo, &tErrors.ErrorTemporaryServerError{}
}

func RemoveAssetTrust(signerUser *userModels.User, wallet *userModels.UserWallet, trustLineInfo *userModels.Trustline, gc *sharedconfig.GlobalConfig) (*userModels.Trustline, error) {
	trustLineInfo.NetworkPassPhrase = gc.BantuNetworkPassphrase
	if wallet.SharedAccessEnabled == 0 {
		//managed access not enabled
		if wallet.Signer != signerUser.PrimarySigner {
			return trustLineInfo, &tErrors.CustomError{Param: "publicKey", Err: "error-wallet-not-managed-by-user", ErrMessage: "You do not have permission to operate on this wallet", Code: http.StatusBadRequest}

		}
		//generatexdr
		xdrBase64, err := generateRemoveTrustAssetXdr(wallet, trustLineInfo.AssetCode, trustLineInfo.AssetIssuer, gc)
		if err != nil {
			return trustLineInfo, err
		}
		oldTransaction := trustLineInfo.Transaction
		trustLineInfo.Transaction = xdrBase64
		if len(trustLineInfo.TransactionSignature) == 0 {
			//needs to be signed first
			return trustLineInfo, nil
		}
		if oldTransaction != xdrBase64 {
			return trustLineInfo, &tErrors.CustomError{
				Param:      "transaction",
				Err:        "transaction mismatch",
				ErrMessage: "transaction mismatch, please try again",
				Code:       http.StatusBadRequest,
			}
		}
		//submit to network
		txnHash, err := network.SubmitXdrWithSignature(gc.BantuExpansionClient, signerUser.PrimarySigner, xdrBase64, trustLineInfo.TransactionSignature)
		if err != nil {
			return trustLineInfo, err
		}
		trustLineInfo.TransactionID = txnHash
		return trustLineInfo, nil
	} else if wallet.SharedAccessEnabled == 1 {
		//perform managed access operation and save to table
		return trustLineInfo, nil
	}
	//did not match any of the conditions
	return trustLineInfo, &tErrors.ErrorTemporaryServerError{}
}

package users

import (
	"log"
	usersdb "trovo-wallet-api/internal/components/users/db"
	usermodels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/validators"

	"github.com/shopspring/decimal"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/txnbuild"
	"gorm.io/gorm"
)

//ClaimPendingAsset claim pending assets
func ClaimPendingAsset(identifier string, sourcePublicKey string, pendingAssetToClaim *usermodels.PendingAssetToClaim, db *gorm.DB) (*usermodels.PendingAssetToClaim, bool, error) {

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

	xdrBase64, err := generateXdr(horizonClient, identifier, sourcePublicKey, db, pendingAssetToClaim)

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

	txnID, err := network.SubmitXdrWithSignature(horizonClient, sourcePublicKey, xdrBase64, pendingAssetToClaim.TransactionSignature)

	if err != nil {
		log.Printf("error submitting txn: %v", err)
		return pendingAssetToClaim, false, &tErrors.ErrorTemporaryServerError{}
	}

	pendingAssetToClaim.TransactionID = txnID

	return pendingAssetToClaim, true, nil

}

func generateXdr(horizonClient *horizonclient.Client, identifier string, sourcePublicKey string, db *gorm.DB, pendingAssetToClaim *usermodels.PendingAssetToClaim) (string, error) {
	user, err := usersdb.GetUser(identifier, db)

	if err != nil {
		return "", err
	}
	if user.Suspended == 1 {
		return "", &tErrors.ErrorUsernameIsSuspended{}
	}
	// if banned, errBanned := usersdb.PublicKeyIsBanned(user.PublicKey, db); banned {
	// 	return "", errBanned
	// }
	if user.PublicKey != sourcePublicKey {
		return "", &tErrors.ErrorInvalidAuthorization{}
	}

	tempKeyPair, err := network.TempAccountKeypair(user.PublicKey)

	if err != nil {
		return "", err
	}

	//asset to Claim

	var asset txnbuild.Asset = nil

	asset = txnbuild.NativeAsset{}

	if len(pendingAssetToClaim.AssetCode) > 0 {
		asset = txnbuild.CreditAsset{Code: pendingAssetToClaim.AssetCode, Issuer: pendingAssetToClaim.AssetIssuer}
	}

	tempAccountExist, tempAccountTrustsAsset, nativeAccountBalance, customAccountBalance, tempAccount, err := network.BlockchainAccountProperties(horizonClient, tempKeyPair.FromAddress().Address(), asset)

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

	sourceAccountExists, sourceAccountTrustsAsset, _, _, sourceAccount, _ := network.BlockchainAccountProperties(horizonClient, sourcePublicKey, asset)

	var ops []txnbuild.Operation = make([]txnbuild.Operation, 0)

	if nativeAccountBalance.GreaterThan(decimal.Zero) {
		if !sourceAccountExists {
			//todo: check if nativeAccount balance > 1
			ops = append(ops, &txnbuild.CreateAccount{
				Destination:   sourcePublicKey,
				Amount:        nativeAccountBalance.Truncate(7).String(),
				SourceAccount: tempAccount.AccountID,
			})
		} else {
			ops = append(ops, &txnbuild.Payment{
				Destination:   sourcePublicKey,
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
			Destination:   sourcePublicKey,
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
		log.Println("error constructing transaction ", err)
		return "", &tErrors.ErrorTemporaryServerError{}
	}

	xdrBase64, err := tx.Base64()

	if err != nil {
		return "", err
	}

	return xdrBase64, nil

}

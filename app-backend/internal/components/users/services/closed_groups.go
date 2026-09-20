package users

import (
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"strings"
	"trovo-wallet-api/internal/basetxn"
	swapModel "trovo-wallet-api/internal/components/swaps/models"
	swaps "trovo-wallet-api/internal/components/swaps/services"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/evmkeypair"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetUserClosedGroups(db *gorm.DB) (ucgs []userModels.UserClosedGroup) {
	ucgs = make([]userModels.UserClosedGroup, 0)
	db.Preload(clause.Associations).Order("closed_group_Id").Find(&ucgs)

	return
}

func GetClosedGroupByOwner(groupOwner string, db *gorm.DB) (cgs []userModels.ClosedGroup) {
	cgs = make([]userModels.ClosedGroup, 0)
	db.Preload(clause.Associations).Order("group_name").Where("group_owner = ?", groupOwner).Find(&cgs)

	return
}

func CreateClosedGroup(owner *userModels.User, cgInput userModels.ClosedGroupJSONInput, gc *sharedconfig.GlobalConfig) (cg userModels.ClosedGroup, err error) {
	cgInput.NetworkPassPhrase = network.GetBlockchainNetworkPassPhrase()
	cgInput.Messages = make([]string, 0)

	tx := gc.DB.Begin()
	defer tx.Rollback()
	cgInput.GroupName = strings.ToUpper(cgInput.GroupName)
	cg = userModels.ClosedGroup{
		ID:                 strings.ReplaceAll(cgInput.GroupName, " ", ""),
		GroupName:          cgInput.GroupName,
		GroupDescription:   cgInput.GroupDescription,
		RegisteredEntity:   cgInput.RegisteredEntity,
		RegistrationName:   cgInput.RegistrationName,
		RegistrationNumber: cgInput.RegistrationNumber,
	}
	e := tx.Omit(clause.Associations).Create(&cg).Error

	if e != nil {
		log.Printf("[CreateClosedGroup] error creating closed group %v for user [%v], error: %v\n", cgInput, owner.Username, e)

		return cg, &tErrors.CustomError{
			Param:      "Id",
			Err:        "error-unable-to-create-group",
			ErrMessage: "Unable to create closed group",
		}
	}
	var base64xdr, txnHash string
	//get the xdr to sign
	if len(cgInput.TransactionSignature) == 0 {
		base64xdr, err = generateClosedGroupXdr(owner, &cgInput, gc)

		if err != nil {
			log.Println("[CreateClosedGroup] error generating patron subscription xdr", err)
			return cg, err
		}

		cgInput.Transaction = base64xdr

		return

	}

	if len(cgInput.Transaction) > 0 && len(cgInput.TransactionSignature) > 0 {
		//transaction signed
		txnHash, err = network.SubmitXdrWithSignature(gc.BantuExpansionClient, owner.PrimarySigner, cgInput.Transaction, cgInput.TransactionSignature)
		if err != nil {
			logDiscordFailedSubscription(fmt.Sprintf("Error submitting closed group [%+v] transaction: %s", cgInput, err.Error()))

			return cg, err
		}

		//on success commit the subscription
		tx.Commit()
		cgInput.TransactionID = txnHash

	}

	//clear user cache.
	owner.InvalidateUserCache(gc)

	return

}

func UploadClosedGroupRegistrationDocument(groupOwner *userModels.User, file multipart.File, fileNameWithExt string, closedGroup *userModels.ClosedGroup, gc *sharedconfig.GlobalConfig) (string, error) {

	newThumbnail, err := gc.FirebaseStorageUploader.UploadFile(file, fileNameWithExt, "")
	if err != nil {
		return "", err
	}

	//update the thumbnail url
	url := fmt.Sprintf("https://storage.googleapis.com/%v/%v", gc.FirebaseStorageUploader.BucketName, newThumbnail)
	//check if document already saved and then retireve it:
	if len(closedGroup.RegistrationDocumentUrl) > 0 {
		//existing record match, update
		closedGroup.RegistrationDocumentUrl = url
		es := gc.DB.Omit(clause.Associations).Save(closedGroup).Error
		if es != nil {

			log.Printf("[UploadClosedGroupRegistrationDocument]error saving existing document in database  [%v] for %v: %v\n", closedGroup, groupOwner.Username, es)
			return "", fmt.Errorf("error saving document %v", closedGroup.GroupName)

		}
	}

	groupOwner.InvalidateUserCache(gc)
	owner, err := userModels.Username(groupOwner.Username).GetFullUser(gc.DB, gc)
	if err == nil {
		if owner.Username == groupOwner.Username {
			groupOwner = &owner
		}

	}

	return url, nil
}

func generateClosedGroupXdr(owner *userModels.User, closedGroupInput *userModels.ClosedGroupJSONInput, gc *sharedconfig.GlobalConfig) (string, error) {
	nativeAssetCode := os.Getenv("NATIVE_ASSET_CODE")
	CLOSED_GROUP_FEE := owner.UserWallets[0].GetClosedGroupFee(gc)
	cgFeeAmountUSD := CLOSED_GROUP_FEE.FeeFixed
	cgFeeAssetCode := CLOSED_GROUP_FEE.FeeAssetCode
	cgFeeAssetIssuer := CLOSED_GROUP_FEE.FeeAssetIssuer
	if len(cgFeeAssetCode) == 0 {
		cgFeeAmountUSD = 0
	}
	if cgFeeAmountUSD > 0 {
		cgFeeAssetCode = "TROV"
		cgFeeAssetIssuer = "GAXMBPVA2GNG6A3NV6Q664VZASMROS5ZACKSMTPVCRIKPOJIV43A2CTJ"
	}
	var ops []basetxn.Operation = make([]basetxn.Operation, 0)
	cgFeeKP := evmkeypair.MustParseFull(CLOSED_GROUP_FEE.FeeWalletSecretKey)
	sourceAssets := strings.ToUpper(fmt.Sprintf("%v:%v", cgFeeAssetCode, cgFeeAssetIssuer))

	var errGetEstimate error
	var requiredUsdWorth string
	// var path []basetxn.Asset

	// var asset basetxn.Asset
	asset := basetxn.CreditAsset{Code: cgFeeAssetCode, Issuer: cgFeeAssetIssuer}

	sourceAccountExists, _, nativeBalance, customBalance, sourceAccount, sourceAccountErr := network.BlockchainAccountProperties(gc.BantuExpansionClient, owner.Address, asset)

	if sourceAccountErr != nil {
		log.Println("[generateClosedGroupXdr] error checking account properties on blockchain. Error ", sourceAccountErr)

		return "", sourceAccountErr
	}

	if !sourceAccountExists {
		log.Println("[generateClosedGroupXdr] error account does not exist on ledger. Error ")

		return "", &tErrors.ErrorUnderfundedAccount{
			Detail: "You need to activate your wallet first and fund it with TROV token to proceed.",
		}
	}

	//get the trov quantity/equivalent needed for the USD from the market.
	pathInput := swapModel.SwapPathInput{
		SourceAssets:           sourceAssets,
		DestinationAssetCode:   CLOSED_GROUP_FEE.FeeAssetCode,
		DestinationAssetIssuer: CLOSED_GROUP_FEE.FeeAssetIssuer,
		DestinationAmount:      fmt.Sprintf("%v", cgFeeAmountUSD),
	}
	_, requiredUsdWorth, errGetEstimate = swaps.GetStrictReceivePaths(pathInput, gc.BantuExpansionClient)
	// requiredTrovAssetEstimate = requiredUsdEstimate

	log.Printf("requires %v %v to convert to %v %v\n", requiredUsdWorth, cgFeeAssetCode, cgFeeAmountUSD, CLOSED_GROUP_FEE.FeeAssetCode)
	if errGetEstimate != nil && requiredUsdWorth == "" {
		log.Println("[generateClosedGroupXdr] error getting required TROV estimate. Error ", errGetEstimate, requiredUsdWorth)

		return "", errGetEstimate
	}

	if !asset.IsNative() {
		// log.Println("[generateClosedGroupXdr] error account does not exist on ledger. Error ")
		if customBalance.LessThan(decimal.RequireFromString(requiredUsdWorth)) {
			return "", &tErrors.ErrorUnderfundedAccount{
				Detail: fmt.Sprintf("You need to add at least %v %v to make up for the fee.", decimal.RequireFromString(requiredUsdWorth).Sub(customBalance).String(), cgFeeAssetCode),
			}
		}

	} else {
		if nativeBalance.LessThan(decimal.RequireFromString(requiredUsdWorth)) {
			return "", &tErrors.ErrorUnderfundedAccount{
				Detail: fmt.Sprintf("You need to add at least %v %v to make up for the fee.", decimal.RequireFromString(requiredUsdWorth).Sub(nativeBalance).String(), nativeAssetCode),
			}
		}
	}

	ops = append(ops, &basetxn.Payment{
		Destination:   cgFeeKP.Address(),
		Amount:        requiredUsdWorth,
		Asset:         basetxn.CreditAsset{Code: cgFeeAssetCode, Issuer: cgFeeAssetIssuer},
		SourceAccount: owner.Address, //primary wallet
	})
	closedGroupInput.Messages = append(closedGroupInput.Messages, fmt.Sprintf("%v %v will be debited from wallet %v to complete the creation of the closed group.", requiredUsdWorth, cgFeeAssetCode, owner.Username))

	tx, err := basetxn.NewTransaction(
		basetxn.TransactionParams{
			SourceAccount:        sourceAccount.Address,
			IncrementSequenceNum: true,
			Operations:           ops,
			BaseFee:              3000,
			Memo:                 "New ClosedGroup",
		},
	)
	if err != nil {
		log.Println("[generateClosedGroupXdr] error constructing transaction ", err)
		return "", err
	}

	var xdrBase64 string

	xdrBase64, err = tx.Base64()
	if err != nil {
		log.Println("[generateClosedGroupXdr] error getting txn base64", err)
		return "", err
	}

	return xdrBase64, nil

}

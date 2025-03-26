package users

import (
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"strings"
	swapModel "trovo-wallet-api/internal/components/swaps/models"
	swaps "trovo-wallet-api/internal/components/swaps/services"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/shopspring/decimal"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
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
	// var nativeAsset txnbuild.Asset = txnbuild.NativeAsset{} "CLOSED_GROUP_FEE_WALLET","CLOSED_GROUP_FEE_QUOTE_AMOUNT","CLOSED_GROUP_FEE_ASSET_CODE","CLOSED_GROUP_FEE_ASSET_ISSUER"
	nativeAssetCode := os.Getenv("NATIVE_ASSET_CODE")
	cgFeeAmountUSD := strings.TrimSpace(os.Getenv("CLOSED_GROUP_FEE_QUOTE_AMOUNT"))
	cgFeeAssetCode := strings.TrimSpace(os.Getenv("CLOSED_GROUP_FEE_ASSET_CODE"))
	cgFeeAssetIssuer := strings.TrimSpace(os.Getenv("CLOSED_GROUP_FEE_ASSET_ISSUER"))
	if len(cgFeeAssetCode) == 0 {
		cgFeeAmountUSD = "0"
	}
	if len(cgFeeAmountUSD) == 0 {
		cgFeeAssetCode = "TROV"
		cgFeeAssetIssuer = "GAXMBPVA2GNG6A3NV6Q664VZASMROS5ZACKSMTPVCRIKPOJIV43A2CTJ"
	}
	var ops []txnbuild.Operation = make([]txnbuild.Operation, 0)
	cgFeeKP := keypair.MustParseFull(os.Getenv("CLOSED_GROUP_FEE_WALLET"))
	sourceAssets := strings.ToUpper(fmt.Sprintf("%v:%v", cgFeeAssetCode, cgFeeAssetIssuer))

	var errGetEstimate error
	var requiredUsdWorth string
	// var path []txnbuild.Asset

	// var asset txnbuild.Asset
	asset := txnbuild.CreditAsset{Code: cgFeeAssetCode, Issuer: cgFeeAssetIssuer}

	sourceAccountExists, _, nativeBalance, customBalance, sourceAccount, sourceAccountErr := network.BlockchainAccountProperties(gc.BantuExpansionClient, owner.PublicKey, asset)

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
		DestinationAssetCode:   strings.Split(os.Getenv("FEE_QUOTE_DEX_ASSET"), ":")[0],
		DestinationAssetIssuer: strings.Split(os.Getenv("FEE_QUOTE_DEX_ASSET"), ":")[1],
		DestinationAmount:      cgFeeAmountUSD,
	}
	_, requiredUsdWorth, errGetEstimate = swaps.GetStrictReceivePaths(pathInput, gc.BantuExpansionClient)
	// requiredTrovAssetEstimate = requiredUsdEstimate

	log.Printf("requires %v %v to convert to %v %v\n", requiredUsdWorth, cgFeeAssetCode, cgFeeAmountUSD, strings.Split(os.Getenv("FEE_QUOTE_DEX_ASSET"), ":")[0])
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
	// //assume the primary wallet does not have trustline to the trov asset. Build the trustline.
	// _, destAccountTrustsDestinationAsset, _, _, _, _ := network.BlockchainAccountProperties(gc.BantuExpansionClient, owner.PublicKey, txnbuild.CreditAsset{Code: "TROV", Issuer: "GAXMBPVA2GNG6A3NV6Q664VZASMROS5ZACKSMTPVCRIKPOJIV43A2CTJ"})

	// if !destAccountTrustsDestinationAsset {

	// 	ops = append(ops, &txnbuild.ChangeTrust{
	// 		Line:          txnbuild.ChangeTrustAssetWrapper{Asset: txnbuild.CreditAsset{Code: "TROV", Issuer: "GAXMBPVA2GNG6A3NV6Q664VZASMROS5ZACKSMTPVCRIKPOJIV43A2CTJ"}},
	// 		Limit:         "900000000000",
	// 		SourceAccount: owner.PublicKey,
	// 	})
	// }
	// if !strings.EqualFold(patronSubInput.PaymentAssetCode, "TROV") {

	// 	//get swap the asset amount to TROV.
	// 	pathInput := swapModel.SwapSendPathInput{
	// 		DestinationAssets: "TROV:GAXMBPVA2GNG6A3NV6Q664VZASMROS5ZACKSMTPVCRIKPOJIV43A2CTJ",
	// 		SourceAssetCode:   patronSubInput.PaymentAssetCode,
	// 		SourceAssetIssuer: patronSubInput.PaymentAssetIssuer,
	// 		SourceAmount:      requiredUsdWorth,
	// 	}

	// 	path, estimatedTrov, errGetEstimate = swaps.GetStrictSendPaths(pathInput, gc.BantuExpansionClient)
	// 	log.Printf(" %v %v converts to %v %v\n", requiredUsdWorth, patronSubInput.PaymentAssetCode, estimatedTrov, "TROV")
	// 	if errGetEstimate != nil && estimatedTrov == "" {
	// 		log.Printf("[generateClosedGroupXdr] error getting required %v estimate. Error %v", patronSubInput.PaymentAssetCode, errGetEstimate)
	// 		return "", errGetEstimate
	// 	}

	// 	{
	// 		//build a swap operation to swap the non-trov asset to trov so that trov can be debited.
	// 		var sendAsset txnbuild.Asset
	// 		if len(patronSubInput.PaymentAssetIssuer) == 0 {
	// 			sendAsset = txnbuild.NativeAsset{}
	// 		} else {
	// 			sendAsset = txnbuild.CreditAsset{Code: patronSubInput.PaymentAssetCode, Issuer: patronSubInput.PaymentAssetIssuer}
	// 		}

	// 		ops = append(ops, &txnbuild.PathPaymentStrictSend{
	// 			SendAsset:     sendAsset,
	// 			SendAmount:    requiredUsdWorth,
	// 			Destination:   patronFeeKP.Address(),
	// 			DestAsset:     txnbuild.CreditAsset{Code: "TROV", Issuer: "GAXMBPVA2GNG6A3NV6Q664VZASMROS5ZACKSMTPVCRIKPOJIV43A2CTJ"},
	// 			DestMin:       "0.0000001",
	// 			Path:          path,
	// 			SourceAccount: owner.PublicKey, //primary wallet
	// 		})
	// 	}
	// } else {

	// }
	ops = append(ops, &txnbuild.Payment{
		Destination:   cgFeeKP.Address(),
		Amount:        requiredUsdWorth,
		Asset:         txnbuild.CreditAsset{Code: cgFeeAssetCode, Issuer: cgFeeAssetIssuer},
		SourceAccount: owner.PublicKey, //primary wallet
	})
	closedGroupInput.Messages = append(closedGroupInput.Messages, fmt.Sprintf("%v %v will be debited from wallet %v to complete the creation of the closed group.", requiredUsdWorth, cgFeeAssetCode, owner.Username))

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        sourceAccount,
			IncrementSequenceNum: true,
			Operations:           ops,
			BaseFee:              3000,
			Preconditions: txnbuild.Preconditions{
				TimeBounds: txnbuild.NewInfiniteTimeout(),
			},
			Memo: txnbuild.MemoText("New ClosedGroup"),
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

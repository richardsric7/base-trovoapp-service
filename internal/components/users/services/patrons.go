package users

import (
	"log"
	"os"
	"time"
	swapModel "trovo-wallet-api/internal/components/swaps/models"
	swaps "trovo-wallet-api/internal/components/swaps/services"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/network"

	"trovo-wallet-api/internal/sharedconfig"

	"github.com/shopspring/decimal"
	"github.com/stellar/go/txnbuild"
)

func GetPatronPackages(gc *sharedconfig.GlobalConfig) (patronPackages []userModels.PatronPackage) {
	patronPackages = make([]userModels.PatronPackage, 0)
	gc.DB.Order("priority_order ASC").Where("inactive = ?", 0).Find(&patronPackages)
	return
}

func GetPatronTiers(gc *sharedconfig.GlobalConfig) (patronTiers []userModels.PatronTier) {
	patronTiers = make([]userModels.PatronTier, 0)
	gc.DB.Order("priority_order ASC").Where("inactive = ?", 0).Find(&patronTiers)
	return
}

func GetPatronMembershipGrades(gc *sharedconfig.GlobalConfig) (memberships []userModels.PatronMembershipGrade) {
	memberships = make([]userModels.PatronMembershipGrade, 0)
	gc.DB.Find(&memberships)
	return
}
func GetPatronMembershipGradeByID(id uint64, gc *sharedconfig.GlobalConfig) (membership userModels.PatronMembershipGrade, err error) {
	err = gc.DB.First(&membership, id).Error
	return
}

func GetPatronSubscriptionLogs(username string, gc *sharedconfig.GlobalConfig) (patronSubLogs []userModels.UserPatronSubscriptionLog) {
	patronSubLogs = make([]userModels.UserPatronSubscriptionLog, 0)
	gc.DB.Order("createdAt DESC").Where("username = ?", username).Find(&patronSubLogs)
	return
}

func GetPatronSubscription(username string, gc *sharedconfig.GlobalConfig) (patronSub userModels.UserPatronMembership, err error) {

	err = gc.DB.Order("createdAt DESC").Where("username = ?", username).First(&patronSub).Error

	return
}

func SubscribeToPatronPackage(signerUser *userModels.User, patronSubInput *userModels.PatronSubscriptionInput, gc *sharedconfig.GlobalConfig) (err error) {
	priceConfig, err := userModels.PatronMembershipGradeID(patronSubInput.PatronMembershipGradeID).GetPatronMemberShipConfig(gc)
	if err != nil {
		return
	} else {
		log.Println(priceConfig)
	}

	return
}

func generatePatronSubscriptionXdr(owner *userModels.User, primaryWallet *userModels.UserWallet, patronSubInput *userModels.PatronSubscriptionInput, priceConfig *userModels.PatronMembershipGrade, gc *sharedconfig.GlobalConfig) (string, error) {
	var nativeAsset txnbuild.Asset = txnbuild.NativeAsset{}
	// check if it is a new subscription or old
	maxDateTime := time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC)
	log.Println(maxDateTime)
	var subscriptionExists bool
	var subscription userModels.UserPatronMembership
	subscription, errGetSub := GetPatronSubscription(owner.Username, gc)

	if errGetSub == nil {
		log.Println(subscription)
		subscriptionExists = true
	}

	patronMembership, errGetMem := GetPatronMembershipGradeByID(patronSubInput.PatronMembershipGradeID, gc)

	if errGetMem != nil {
		return "", &tErrors.CustomError{
			Param:      "patronpackageId",
			Err:        "error-invalid-membership",
			ErrMessage: "Submitted tier and package are invalid.",
		}
	}

	sourceAccountExists, _, _, _, _, sourceAccountErr := network.BlockchainAccountProperties(gc.BantuExpansionClient, primaryWallet.ID, nativeAsset)

	if sourceAccountErr != nil {
		log.Println("[generatePatronSubscriptionXdr] error checking account properties on blockchain. Error ", sourceAccountErr)

		return "", sourceAccountErr
	}

	if !sourceAccountExists {
		log.Println("[generatePatronSubscriptionXdr] error account does not exist on ledger. Error ")

		return "", &tErrors.ErrorUnderfundedAccount{
			Detail: "You need to activate your wallet first and fund it with TROV token to proceed.",
		}
	}

	//get the trov needed.
	pathInput := swapModel.SwapPathInput{
		SourceAssets:           os.Getenv("DOLLAR_ASSET"),
		DestinationAssetCode:   "TROV",
		DestinationAssetIssuer: "GAXMBPVA2GNG6A3NV6Q664VZASMROS5ZACKSMTPVCRIKPOJIV43A2CTJ",
		DestinationAmount:      decimal.NewFromFloat(priceConfig.Price).Truncate(7).String(),
	}
	_, requiredTrovEstimate, errGetEstimate := swaps.GetStrictReceivePaths(pathInput, gc.BantuExpansionClient)

	if errGetEstimate != nil {
		log.Println("[generatePatronSubscriptionXdr] error getting required TROV estimate. Error ", errGetEstimate, requiredTrovEstimate)
		return "", errGetEstimate
	}
	{
		//routine checks for package subscription ability
		if subscriptionExists {
			//run routine for subscription exists
			//check if it is higest tier
			if subscription.PatronPackageID == "DIAMOND" {

				if subscription.PatronTierID == "LIFETIME" {
					//cannot upgrade anymore
					return "", &tErrors.CustomError{
						Param:      "patronpackageId",
						Err:        "error-already-higest-tier",
						ErrMessage: "Account is already member of the highest available tier and package",
					}
				}

			} else if subscription.PatronPackageID == "PLATINUM" {
				if subscription.PatronTierID == "LIFETIME" && (patronMembership.PatronPackage == "GOLD" || patronMembership.PatronPackage == "PLATINUM") {
					//cannot downgrade
					return "", &tErrors.CustomError{
						Param:      "patronpackageId",
						Err:        "error-already-higest-tier",
						ErrMessage: "Account is already member of the highest available tier in this package",
					}
				}

			} else if subscription.PatronPackageID == "GOLD" {
				if subscription.PatronTierID == "LIFETIME" && patronMembership.PatronPackage == "GOLD" {
					//cannot downgrade
					return "", &tErrors.CustomError{
						Param:      "patronpackageId",
						Err:        "error-already-higest-tier",
						ErrMessage: "Account is already member of the highest available tier in this package",
					}
				}
			}
			// if subscription.ValidTill.Year()==maxDateTime.Year(){
			// 	//already life time
			// }
		} else {
			//run routine for new subscription
		}
	}

	return "", nil

}

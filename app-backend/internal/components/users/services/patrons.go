package users

import (
	"log"
	"os"
	"time"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/network"

	"trovo-wallet-api/internal/sharedconfig"

	"github.com/ecnepsnai/discord"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
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
	gc.DB.Order("created_at DESC").Where("username = ?", username).Find(&patronSubLogs)
	return
}

func GetPatronSubscriptionPaymentAssets(gc *sharedconfig.GlobalConfig) (paymentAssets []userModels.PatronSubscriptionPaymentAsset) {
	paymentAssets = make([]userModels.PatronSubscriptionPaymentAsset, 0)
	gc.DB.Where("inactive = ?", 0).Find(&paymentAssets)
	return
}

func GetPatronSubscription(username string, gc *sharedconfig.GlobalConfig) (patronSub userModels.UserPatronMembership, err error) {

	e := gc.DB.Preload(clause.Associations).Where("username = ?", username).First(&patronSub).Error
	if e != nil {
		log.Printf("[GetPatronSubscription] error : %v\n", e)
		err = &tErrors.ErrorTemporaryServerError{}
		return
	}
	return patronSub, nil
}

func countPendingSubscriptionForUser(username string, gc *sharedconfig.GlobalConfig) (int64, error) {
	var pendingSubscriptionsCount int64
	err := gc.DB.Model(&userModels.UserPatronSubscriptionLog{}).
		Where("username = ? AND CAST(effective_date AS DATE) > CAST(? AS DATE)", username, time.Now()).
		Count(&pendingSubscriptionsCount).Error
	if err != nil {
		log.Printf("[CountPendingSubscription] error : %v\n", err)
		err = &tErrors.ErrorTemporaryServerError{}
		return -1, err
	}

	return pendingSubscriptionsCount, nil
}

func SubscribeToPatronPackage(owner *userModels.User, patronSubInput *userModels.PatronSubscriptionInput, gc *sharedconfig.GlobalConfig) (subscriptionLog userModels.UserPatronSubscriptionLog, err error) {

	patronSubInput.NetworkPassPhrase = network.GetBlockchainNetworkPassPhrase()
	patronSubInput.Messages = make([]string, 0)
	priceConfig, err := userModels.PatronMembershipGradeID(patronSubInput.PatronMembershipGradeID).GetPatronMemberShipConfig(gc)
	if err != nil {
		return
	} else {
		log.Println(priceConfig)
	}

	pendingSubscriptionCount, err := countPendingSubscriptionForUser(owner.Username, gc)
	if err != nil {
		return
	}
	if pendingSubscriptionCount > 0 {
		return subscriptionLog, &tErrors.CustomError{
			Param:      "patronPackageId",
			Err:        "you-have-pending-subscription",
			ErrMessage: "You have a pending subscription",
		}
	}
	// check if it is a new subscription or old
	lifetime := time.Date(9999, 12, 1, 23, 59, 59, 000000000, time.UTC)
	log.Println(lifetime)
	var subscriptionExists, subscriptionRenewal, instantActivation bool
	var subscription userModels.UserPatronMembership
	subscription, errGetSub := GetPatronSubscription(owner.Username, gc)

	if errGetSub == nil {
		log.Println(subscription)
		subscriptionExists = true
	}

	patronMembership, errGetMem := GetPatronMembershipGradeByID(patronSubInput.PatronMembershipGradeID, gc)

	if errGetMem != nil {
		return subscriptionLog, &tErrors.CustomError{
			Param:      "patronPackageId",
			Err:        "error-invalid-membership",
			ErrMessage: "Submitted tier and package are invalid.",
		}
	}

	// primaryWallet, _:=owner.GetWalletByAddress(owner.Address, gc.DB)

	{
		//routine checks for package subscription qualification
		if subscriptionExists {
			//run routine for subscription exists
			//check if it is higest tier
			if subscription.PatronPackageID == "DIAMOND" {

				if subscription.PatronTierID == "LIFETIME" {
					//cannot upgrade anymore
					return subscriptionLog, &tErrors.CustomError{
						Param:      "patronPackageId",
						Err:        "error-already-highest-tier",
						ErrMessage: "Account is already member of the highest available tier and package",
					}
				}

			} else if subscription.PatronPackageID == "PLATINUM" {
				if subscription.PatronTierID == "LIFETIME" && (patronMembership.PatronPackage == "GOLD" || patronMembership.PatronPackage == "PLATINUM") {
					//cannot downgrade
					return subscriptionLog, &tErrors.CustomError{
						Param:      "patronPackageId",
						Err:        "error-already-highest-tier",
						ErrMessage: "Account is already member of the highest available tier in this package",
					}
				}

			} else if subscription.PatronPackageID == "GOLD" {
				if subscription.PatronTierID == "LIFETIME" && patronMembership.PatronPackage == "GOLD" {
					//cannot downgrade
					return subscriptionLog, &tErrors.CustomError{
						Param:      "patronPackageId",
						Err:        "error-already-highest-tier",
						ErrMessage: "Account is already member of the highest available tier in this package",
					}
				}
			}

		}
	}
	// var subscriptionLog userModels.UserPatronSubscriptionLog
	tx := gc.DB.Begin()
	defer tx.Rollback()

	subscriptionLog = userModels.UserPatronSubscriptionLog{
		ID:                    uuid.NewString(),
		Username:              owner.Username,
		PatronPackageID:       patronMembership.PatronPackage,
		PatronTierID:          patronMembership.PatronTierID,
		ActivePatronPackageID: &subscription.PatronPackageID,
		ActivePatronTierID:    &subscription.PatronTierID,
		// EffectiveDate:         subscription.ValidTill.AddDate(0, 0, 1), // starts the next day that the active subsription expires
	}
	if subscriptionExists {
		//can be upgrade or renewal

		// subscription.PatronPackageID = patronMembership.PatronPackage
		// subscription.PatronTierID = patronMembership.PatronTierID
		if patronMembership.PatronTierID == "LIFETIME" {
			// subscription.ValidTill = lifetime
			subscriptionLog.ValidTill = lifetime
			if subscription.ValidTill.UTC().After(time.Now().UTC()) {
				//EXISTING subscription still valid
				if subscription.ValidTill.Year() == 9999 || subscription.PatronTierID == "ANNUAL" || subscription.PatronTierID == "LIFETIME" {
					//existing subscription is lifetime or annual, activate immediately
					subscriptionLog.EffectiveDate = time.Now()
					instantActivation = true

				} else {
					subscriptionLog.EffectiveDate = subscription.ValidTill.AddDate(0, 0, 1) //1 day after the active subscription expires

				}

			} else {
				//subscription expired. start immediately
				subscriptionLog.EffectiveDate = time.Now()
				subscriptionRenewal = true
			}

		}
		if patronMembership.PatronTierID == "ANNUAL" {
			if subscription.ValidTill.UTC().After(time.Now().UTC()) {
				//subscription still valid
				if subscription.ValidTill.Year() == 9999 || subscription.PatronTierID == "ANNUAL" || subscription.PatronTierID == "LIFETIME" {
					//lifetime
					subscriptionLog.EffectiveDate = time.Now()
					subscriptionLog.ValidTill = time.Now().AddDate(1, 0, 0) //1 year after the active subscription expires
					instantActivation = true

				} else {
					subscriptionLog.EffectiveDate = subscription.ValidTill.AddDate(0, 0, 1) //1 day after the active subscription expires
					subscriptionLog.ValidTill = subscription.ValidTill.AddDate(1, 0, 0)     //1 year after the active subscription expires

				}
			} else {
				// subscription already expired
				subscriptionLog.EffectiveDate = time.Now()
				subscriptionLog.ValidTill = time.Now().AddDate(1, 0, 0) //1 year after the active subscription expires
				subscriptionRenewal = true
			}

			// subscription.ValidTill = subscription.ValidTill.AddDate(1, 0, 0) //1 year after the active subscription expires
		}
		if patronMembership.PatronTierID == "MONTHLY" {
			if subscription.ValidTill.UTC().After(time.Now().UTC()) {
				//subscription still valid
				if subscription.ValidTill.Year() == 9999 || subscription.PatronTierID == "ANNUAL" || subscription.PatronTierID == "LIFETIME" {
					//lifetime
					subscriptionLog.EffectiveDate = time.Now()
					subscriptionLog.ValidTill = time.Now().AddDate(0, 1, 0) //1 month after
					instantActivation = true

				} else {
					subscriptionLog.EffectiveDate = subscription.ValidTill.AddDate(0, 0, 1) //1 day after the active subscription expires
					subscriptionLog.ValidTill = subscription.ValidTill.AddDate(0, 1, 0)     //1 month after

				}

			} else {
				//expired subscription. start immediately
				subscriptionLog.ValidTill = time.Now().AddDate(0, 1, 0) //1 month after today
				subscriptionLog.EffectiveDate = time.Now()
				subscriptionRenewal = true
			}
		}

		// do not create or modify subscription until the effective date.
		e := tx.Omit(clause.Associations).Create(&subscriptionLog).Error

		if e != nil {
			log.Printf("[SubscribeToPatronPackage] error creating subscriptionLog for user [%v], error: %v\n", owner.Username, e)

			return subscriptionLog, &tErrors.CustomError{
				Param:      "patronPackageId",
				Err:        "error-unable to subscribe",
				ErrMessage: "Unable to subscribe to this package",
			}
		}

		//if it is a renewal of expired subscription or instantActivation is set then activate immediately.
		if subscriptionRenewal || instantActivation {
			subscription.PatronPackageID = patronMembership.PatronPackage
			subscription.PatronTierID = patronMembership.PatronTierID
			subscription.ValidTill = subscriptionLog.ValidTill
			subscriptionLog.EffectiveDate = time.Now()

			e := tx.Omit(clause.Associations).Save(&subscription).Error

			if e != nil {
				log.Printf("[SubscribeToPatronPackage] error saving subscription for user [%v], [%+v], error: %v\n", owner.Username, subscription, e)

				return subscriptionLog, &tErrors.CustomError{
					Param:      "patronPackageId",
					Err:        "error-unable to subscribe",
					ErrMessage: "Unable to subscribe to this package",
				}
			}
		}

	} else {
		//new subscription

		if patronMembership.PatronTierID == "LIFETIME" {
			// subscription.ValidTill = lifetime
			subscriptionLog.ValidTill = lifetime
			subscriptionLog.EffectiveDate = time.Now()
		}
		if patronMembership.PatronTierID == "ANNUAL" {

			subscriptionLog.EffectiveDate = time.Now()
			subscriptionLog.ValidTill = time.Now().AddDate(1, 0, 0) //1 year after the active subscription expires

			// subscription.ValidTill = subscription.ValidTill.AddDate(1, 0, 0) //1 year after the active subscription expires
		}

		if patronMembership.PatronTierID == "MONTHLY" {

			subscriptionLog.EffectiveDate = time.Now()
			subscriptionLog.ValidTill = time.Now().AddDate(0, 1, 0) //1 month after

		}
		e := tx.Omit(clause.Associations).Create(&subscriptionLog).Error

		if e != nil {
			log.Printf("[SubscribeToPatronPackage] error creating subscriptionLog for user [%v], [%+v], error: %v\n", owner.Username, subscriptionLog, e)

			return subscriptionLog, &tErrors.CustomError{
				Param:      "patronPackageId",
				Err:        "error-unable to subscribe",
				ErrMessage: "Unable to subscribe to this package",
			}
		}
		//if new subscription, activate it immediately
		subscription = userModels.UserPatronMembership{
			Username:        owner.Username,
			PatronPackageID: subscriptionLog.PatronPackageID,
			PatronTierID:    subscriptionLog.PatronTierID,
			ValidTill:       subscriptionLog.ValidTill,
		}
		e = tx.Omit(clause.Associations).Create(&subscription).Error

		if e != nil {
			log.Printf("[SubscribeToPatronPackage] error creating subscription for user [%v], [%+v], error: %v\n", owner.Username, subscription, e)

			return subscriptionLog, &tErrors.CustomError{
				Param:      "patronPackageId",
				Err:        "error-unable-to-subscribe",
				ErrMessage: "Unable to subscribe to this package",
			}
		}

	}

	return subscriptionLog, submitPatronSubscription(owner, patronSubInput, &priceConfig, tx, gc)
}

func logDiscordFailedSubscription(msg string) {
	discord.WebhookURL = "https://discord.com/api/webhooks/827986576415129663/wqMKp9wxB_fxs9Q3zlMKCNPGENXmD_ueUnL8hVCu1wmRfD2wkXAjfP85k1Ro_2_wGfiY"
	if len(os.Getenv("FAILED_PAYMENT_ERROR_WEBHOOK")) > 50 {
		discord.WebhookURL = os.Getenv("FAILED_PAYMENT_ERROR_WEBHOOK")
	}
	discord.Say(msg)
}

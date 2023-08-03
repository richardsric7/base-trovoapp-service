package users

import (
	"fmt"
	"log"
	"os"
	"time"
	swapModel "trovo-wallet-api/internal/components/swaps/models"
	swaps "trovo-wallet-api/internal/components/swaps/services"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/network"

	"trovo-wallet-api/internal/sharedconfig"

	"github.com/ecnepsnai/discord"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stellar/go/keypair"
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

func SubscribeToPatronPackage(owner *userModels.User, patronSubInput *userModels.PatronSubscriptionInput, gc *sharedconfig.GlobalConfig) (subscriptionLog userModels.UserPatronSubscriptionLog, err error) {

	patronSubInput.NetworkPassPhrase = network.GetBlockchainNetworkPassPhrase()
	patronSubInput.Messages = make([]string, 0)
	priceConfig, err := userModels.PatronMembershipGradeID(patronSubInput.PatronMembershipGradeID).GetPatronMemberShipConfig(gc)
	if err != nil {
		return
	} else {
		log.Println(priceConfig)
	}
	// check if it is a new subscription or old
	lifetime := time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC)
	log.Println(lifetime)
	var subscriptionExists, subscriptionRenewal bool
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

	// primaryWallet, _:=owner.GetWalletByPublicKey(owner.PublicKey, gc.DB)

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

		} else {
			//run routine for new subscription
		}
	}
	// var subscriptionLog userModels.UserPatronSubscriptionLog
	tx := gc.DB.Begin()
	defer tx.Rollback()
	if subscriptionExists {
		//can be upgrade or renewal
		subscriptionLog = userModels.UserPatronSubscriptionLog{
			ID:                    uuid.NewString(),
			Username:              owner.Username,
			PatronPackageID:       patronMembership.PatronPackage,
			PatronTierID:          patronMembership.PatronTierID,
			ActivePatronPackageID: &subscription.PatronPackageID,
			ActivePatronTierID:    &subscription.PatronTierID,
			// EffectiveDate:         subscription.ValidTill.AddDate(0, 0, 1), // starts the next day that the active subsription expires
		}
		// subscription.PatronPackageID = patronMembership.PatronPackage
		// subscription.PatronTierID = patronMembership.PatronTierID
		if patronMembership.PatronTierID == "LIFETIME" {
			// subscription.ValidTill = lifetime
			subscriptionLog.ValidTill = lifetime
			if subscription.ValidTill.UTC().After(time.Now().UTC()) {
				//subscription still valid
				if subscription.ValidTill.Year() == 9999 {
					//lifetime
					subscriptionLog.EffectiveDate = time.Now()

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
				if subscription.ValidTill.Year() == 9999 {
					//lifetime
					subscriptionLog.EffectiveDate = time.Now()
					subscriptionLog.ValidTill = time.Now().AddDate(1, 0, 0) //1 year after the active subscription expires

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
				if subscription.ValidTill.Year() == 9999 {
					//lifetime
					subscriptionLog.EffectiveDate = time.Now()
					subscriptionLog.ValidTill = time.Now().AddDate(0, 1, 0) //1 month after

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
		e := tx.Save(&subscriptionLog).Error

		if e != nil {
			log.Printf("[SubscribeToPatronPackage] error creating subscriptionLog for user [%v], error: %v\n", owner.Username, e)

			return subscriptionLog, &tErrors.CustomError{
				Param:      "patronPackageId",
				Err:        "error-unable to subscribe",
				ErrMessage: "Unable to subscribe to this package",
			}
		}

		//if it is a renewal of expired subscription then activate immediately.
		if subscriptionRenewal {
			subscription.PatronPackageID = patronMembership.PatronPackage
			subscription.PatronTierID = patronMembership.PatronTierID
			subscription.ValidTill = subscriptionLog.ValidTill

			e := tx.Save(&subscription).Error

			if e != nil {
				log.Printf("[SubscribeToPatronPackage] error saving subscription for user [%v], error: %v\n", owner.Username, e)

				return subscriptionLog, &tErrors.CustomError{
					Param:      "patronPackageId",
					Err:        "error-unable to subscribe",
					ErrMessage: "Unable to subscribe to this package",
				}
			}
		}

	} else {
		//new subscription

		//if new subscription, activate it immediately
		subscription = userModels.UserPatronMembership{
			Username:        owner.Username,
			PatronPackageID: subscriptionLog.PatronPackageID,
			PatronTierID:    subscriptionLog.PatronTierID,
			ValidTill:       subscriptionLog.ValidTill,
		}
		e := tx.Create(&subscription).Error

		if e != nil {
			log.Printf("[SubscribeToPatronPackage] error creating subscription for user [%v], error: %v\n", owner.Username, e)

			return subscriptionLog, &tErrors.CustomError{
				Param:      "patronPackageId",
				Err:        "error-unable-to-subscribe",
				ErrMessage: "Unable to subscribe to this package",
			}
		}

	}

	var base64xdr, txnHash string
	//get the xdr to sign
	if len(patronSubInput.TransactionSignature) == 0 {
		base64xdr, err = generatePatronSubscriptionXdr(owner, patronSubInput, &priceConfig, gc)

		if err != nil {
			log.Println("[SubscribeToPatronPackage] error generating patron subscription xdr", err)
			return subscriptionLog, err
		}

		patronSubInput.Transaction = base64xdr

		return

	}

	if len(patronSubInput.Transaction) > 0 && len(patronSubInput.TransactionSignature) > 0 {
		//transaction signed
		txnHash, err = network.SubmitXdrWithSignature(gc.BantuExpansionClient, owner.PrimarySigner, patronSubInput.Transaction, patronSubInput.TransactionSignature)
		if err != nil {
			logDiscordFailedSubscription(fmt.Sprintf("Error submitting patron subscription [%+v] transaction: %s", patronSubInput, err.Error()))

			return subscriptionLog, err
		}

		//on success commit the subscription
		tx.Commit()
		patronSubInput.TransactionID = txnHash

	}

	//clear user cache.
	owner.InvalidateUserCache(gc)

	return

}

func generatePatronSubscriptionXdr(owner *userModels.User, patronSubInput *userModels.PatronSubscriptionInput, priceConfig *userModels.PatronMembershipGrade, gc *sharedconfig.GlobalConfig) (string, error) {
	var nativeAsset txnbuild.Asset = txnbuild.NativeAsset{}
	var ops []txnbuild.Operation = make([]txnbuild.Operation, 0)
	chanAccount := <-gc.ChannelAccounts
	defer func(c *keypair.Full) {
		gc.ChannelAccounts <- c
	}(chanAccount)

	_, _, _, _, chanSourceAccount, _ := network.BlockchainAccountProperties(gc.BantuExpansionClient, chanAccount.Address(), nativeAsset)
	sourceAccountExists, _, _, _, _, sourceAccountErr := network.BlockchainAccountProperties(gc.BantuExpansionClient, owner.PublicKey, nativeAsset)

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

	//get the trov quantity/equivalent needed for the USD from the market.
	pathInput := swapModel.SwapPathInput{
		SourceAssets:           os.Getenv("DOLLAR_ASSET"),
		DestinationAssetCode:   "TROV",
		DestinationAssetIssuer: "GAXMBPVA2GNG6A3NV6Q664VZASMROS5ZACKSMTPVCRIKPOJIV43A2CTJ",
		DestinationAmount:      decimal.NewFromFloat(priceConfig.Price).Truncate(7).String(),
	}
	_, requiredTrovEstimate, errGetEstimate := swaps.GetStrictReceivePaths(pathInput, gc.BantuExpansionClient)

	if errGetEstimate != nil && requiredTrovEstimate == "" {
		log.Println("[generatePatronSubscriptionXdr] error getting required TROV estimate. Error ", errGetEstimate, requiredTrovEstimate)
		return "", errGetEstimate
	}
	patronSubInput.Messages = append(patronSubInput.Messages, fmt.Sprintf("%v TROV will be debited from wallet %v to complete the subscription.", requiredTrovEstimate, owner.Username))

	patronFeeKP := keypair.MustParseFull(os.Getenv("PATRON_FEE_WALLET"))
	ops = append(ops, &txnbuild.Payment{
		Destination:   patronFeeKP.Address(),
		Amount:        requiredTrovEstimate,
		Asset:         txnbuild.CreditAsset{Code: "TROV", Issuer: "GAXMBPVA2GNG6A3NV6Q664VZASMROS5ZACKSMTPVCRIKPOJIV43A2CTJ"},
		SourceAccount: owner.PublicKey, //primary wallet
	})

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        chanSourceAccount,
			IncrementSequenceNum: true,
			Operations:           ops,
			BaseFee:              3000,
			Preconditions: txnbuild.Preconditions{
				TimeBounds: txnbuild.NewInfiniteTimeout(),
			},
			Memo: txnbuild.MemoText(fmt.Sprintf("Patron %v sub", priceConfig.PatronPackage)),
		},
	)
	if err != nil {
		log.Println("[generatePatronSubscriptionXdr] error constructing transaction ", err)
		return "", err
	}

	var xdrBase64 string

	xdrBase64, err = tx.Base64()
	if err != nil {
		log.Println("[generatePatronSubscriptionXdr] error getting txn base64", err)
		return "", err
	}

	return xdrBase64, nil

}

func logDiscordFailedSubscription(msg string) {
	discord.WebhookURL = "https://discord.com/api/webhooks/827986576415129663/wqMKp9wxB_fxs9Q3zlMKCNPGENXmD_ueUnL8hVCu1wmRfD2wkXAjfP85k1Ro_2_wGfiY"
	if len(os.Getenv("FAILED_PAYMENT_ERROR_WEBHOOK")) > 50 {
		discord.WebhookURL = os.Getenv("FAILED_PAYMENT_ERROR_WEBHOOK")
	}
	discord.Say(msg)
}

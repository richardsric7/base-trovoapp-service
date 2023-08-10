package users

import (
	"fmt"
	"log"
	"os"
	"strings"
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
		e := tx.Create(&subscriptionLog).Error

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
		e := tx.Create(&subscriptionLog).Error

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
		e = tx.Create(&subscription).Error

		if e != nil {
			log.Printf("[SubscribeToPatronPackage] error creating subscription for user [%v], [%+v], error: %v\n", owner.Username, subscription, e)

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
	// var nativeAsset txnbuild.Asset = txnbuild.NativeAsset{}
	nativeAssetCode := os.Getenv("NATIVE_ASSET_CODE")
	var ops []txnbuild.Operation = make([]txnbuild.Operation, 0)
	patronFeeKP := keypair.MustParseFull(os.Getenv("PATRON_FEE_WALLET"))
	sourceAssets := ""
	var errGetEstimate error
	var requiredUsdWorth, estimatedTrov string
	var path []txnbuild.Asset
	if len(patronSubInput.PaymentAssetIssuer) == 56 {
		sourceAssets = strings.ToUpper(fmt.Sprintf("%v:%v", patronSubInput.PaymentAssetCode, patronSubInput.PaymentAssetIssuer))
	}

	var asset txnbuild.Asset
	if len(patronSubInput.PaymentAssetIssuer) == 0 {
		asset = txnbuild.NativeAsset{}
	} else {
		asset = txnbuild.CreditAsset{Code: patronSubInput.PaymentAssetCode, Issuer: patronSubInput.PaymentAssetIssuer}
	}
	// chanAccount := <-gc.ChannelAccounts
	// defer func(c *keypair.Full) {
	// 	gc.ChannelAccounts <- c
	// }(chanAccount)

	// _, _, _, _, chanSourceAccount, _ := network.BlockchainAccountProperties(gc.BantuExpansionClient, chanAccount.Address(), nativeAsset)

	sourceAccountExists, _, nativeBalance, customBalance, sourceAccount, sourceAccountErr := network.BlockchainAccountProperties(gc.BantuExpansionClient, owner.PublicKey, asset)

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

	// //get the trov quantity/equivalent needed for the USD from the market.
	// pathInput := swapModel.SwapPathInput{
	// 	SourceAssets:           "TROV:GAXMBPVA2GNG6A3NV6Q664VZASMROS5ZACKSMTPVCRIKPOJIV43A2CTJ",
	// 	DestinationAssetCode:   strings.Split(os.Getenv("DOLLAR_ASSET"), ":")[0],
	// 	DestinationAssetIssuer: strings.Split(os.Getenv("DOLLAR_ASSET"), ":")[1],
	// 	DestinationAmount:      decimal.NewFromFloat(priceConfig.Price).Truncate(7).String(),
	// }
	//get the trov quantity/equivalent needed for the USD from the market.
	pathInput := swapModel.SwapPathInput{
		SourceAssets:           sourceAssets,
		DestinationAssetCode:   strings.Split(os.Getenv("DOLLAR_ASSET"), ":")[0],
		DestinationAssetIssuer: strings.Split(os.Getenv("DOLLAR_ASSET"), ":")[1],
		DestinationAmount:      decimal.NewFromFloat(priceConfig.Price).Truncate(7).String(),
	}
	_, requiredUsdWorth, errGetEstimate = swaps.GetStrictReceivePaths(pathInput, gc.BantuExpansionClient)
	// requiredTrovAssetEstimate = requiredUsdEstimate

	log.Printf("requires %v %v to convert to %v %v\n", requiredUsdWorth, patronSubInput.PaymentAssetCode, priceConfig.Price, "USDT")
	if errGetEstimate != nil && requiredUsdWorth == "" {
		log.Println("[generatePatronSubscriptionXdr] error getting required TROV estimate. Error ", errGetEstimate, requiredUsdWorth)

		return "", errGetEstimate
	}

	if !asset.IsNative() {
		// log.Println("[generatePatronSubscriptionXdr] error account does not exist on ledger. Error ")
		if customBalance.LessThan(decimal.RequireFromString(requiredUsdWorth)) {
			return "", &tErrors.ErrorUnderfundedAccount{
				Detail: fmt.Sprintf("You need to add at least %v %v to make up for the subscription fee.", decimal.RequireFromString(requiredUsdWorth).Sub(customBalance).String(), patronSubInput.PaymentAssetCode),
			}
		}

	} else {
		if nativeBalance.LessThan(decimal.RequireFromString(requiredUsdWorth)) {
			return "", &tErrors.ErrorUnderfundedAccount{
				Detail: fmt.Sprintf("You need to add at least %v %v to make up for the subscription fee.", decimal.RequireFromString(requiredUsdWorth).Sub(nativeBalance).String(), nativeAssetCode),
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
	if !strings.EqualFold(patronSubInput.PaymentAssetCode, "TROV") {

		//get swap the asset amount to TROV.
		pathInput := swapModel.SwapSendPathInput{
			DestinationAssets: "TROV:GAXMBPVA2GNG6A3NV6Q664VZASMROS5ZACKSMTPVCRIKPOJIV43A2CTJ",
			SourceAssetCode:   patronSubInput.PaymentAssetCode,
			SourceAssetIssuer: patronSubInput.PaymentAssetIssuer,
			SourceAmount:      requiredUsdWorth,
		}

		path, estimatedTrov, errGetEstimate = swaps.GetStrictSendPaths(pathInput, gc.BantuExpansionClient)
		log.Printf(" %v %v converts to %v %v\n", requiredUsdWorth, patronSubInput.PaymentAssetCode, estimatedTrov, "TROV")
		if errGetEstimate != nil && estimatedTrov == "" {
			log.Printf("[generatePatronSubscriptionXdr] error getting required %v estimate. Error %v", patronSubInput.PaymentAssetCode, errGetEstimate)
			return "", errGetEstimate
		}

		{
			//build a swap operation to swap the non-trov asset to trov so that trov can be debited.
			var sendAsset txnbuild.Asset
			if len(patronSubInput.PaymentAssetIssuer) == 0 {
				sendAsset = txnbuild.NativeAsset{}
			} else {
				sendAsset = txnbuild.CreditAsset{Code: patronSubInput.PaymentAssetCode, Issuer: patronSubInput.PaymentAssetIssuer}
			}

			ops = append(ops, &txnbuild.PathPaymentStrictSend{
				SendAsset:     sendAsset,
				SendAmount:    estimatedTrov,
				Destination:   patronFeeKP.Address(),
				DestAsset:     txnbuild.CreditAsset{Code: "TROV", Issuer: "GAXMBPVA2GNG6A3NV6Q664VZASMROS5ZACKSMTPVCRIKPOJIV43A2CTJ"},
				DestMin:       "0.0000001",
				Path:          path,
				SourceAccount: owner.PublicKey, //primary wallet
			})
		}
	} else {

		ops = append(ops, &txnbuild.Payment{
			Destination:   patronFeeKP.Address(),
			Amount:        requiredUsdWorth,
			Asset:         txnbuild.CreditAsset{Code: "TROV", Issuer: "GAXMBPVA2GNG6A3NV6Q664VZASMROS5ZACKSMTPVCRIKPOJIV43A2CTJ"},
			SourceAccount: owner.PublicKey, //primary wallet
		})
	}
	if len(patronSubInput.PaymentAssetIssuer) == 0 {
		patronSubInput.Messages = append(patronSubInput.Messages, fmt.Sprintf("%v %v will be debited from wallet %v to complete the subscription.", requiredUsdWorth, nativeAssetCode, owner.Username))

	} else {

		patronSubInput.Messages = append(patronSubInput.Messages, fmt.Sprintf("%v %v will be debited from wallet %v to complete the subscription.", requiredUsdWorth, patronSubInput.PaymentAssetCode, owner.Username))
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        sourceAccount,
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

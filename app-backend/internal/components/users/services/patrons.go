package users

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"trovo-wallet-api/internal/basetxn"
	swapModel "trovo-wallet-api/internal/components/swaps/models"
	swaps "trovo-wallet-api/internal/components/swaps/services"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/evmkeypair"
	"trovo-wallet-api/internal/network"

	"trovo-wallet-api/internal/sharedconfig"

	"github.com/ecnepsnai/discord"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
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

		var dbContractAddress *string
		serviceFee := owner.UserWallets[0].GetPatronFee(gc)
		patronFeeKP, e := evmkeypair.ParseFull(serviceFee.FeeWalletSecretKey)
		if e != nil {
			log.Println("[SubscribeToPatronPackage] error fetching account fee wallet for patron fee ", e)
			logDiscordFailedSubscription("[SubscribeToPatronPackage] error fetching account fee wallet for patron fee")
			return subscriptionLog, &tErrors.ErrorTemporaryServerError{}
		}
		assetLabel := os.Getenv("NATIVE_ASSET_CODE")
		if len(patronSubInput.PaymentAssetCode) > 0 {
			assetLabel = patronSubInput.PaymentAssetCode
		}

		dbAssetCode := assetLabel
		if len(patronSubInput.PaymentContractAddress) > 0 {
			dbContractAddress = &patronSubInput.PaymentContractAddress
		}

		vatFeeCollection := sharedconfig.FeeCollection{
			ID:                         gc.GenerateUUIDString(),
			FromUsername:               owner.Username,
			FromWalletAddress:          owner.Address,
			FromWalletAlias:            owner.Username,
			BelongsToEnterpriseProfile: owner.CreatedByServiceLinkID,
			FeeType:                    "VAT",
			Amount: func() float64 {

				d, e := decimal.NewFromString(patronSubInput.VatAmount)
				if e != nil {
					return 0.00
				}

				return d.InexactFloat64()

			}(),
			AssetCode:             dbAssetCode,
			ContractAddress:       dbContractAddress,
			DestinationWallet:     patronFeeKP.Address(),
			SharedAccessOperation: 0,
		}

		//save vat to database
		e = tx.Omit(clause.Associations).Create(&vatFeeCollection).Error
		if e != nil {

			log.Printf("[SubscribeToPatronPackage] Error saving vat [%+v] transaction on fee collections table: %s\n", vatFeeCollection, e.Error())
			gc.LogDiscordFailedRequest(fmt.Sprintf("[Pay] Error saving vat [%+v] transaction on fee collections table: %s\n", vatFeeCollection, e.Error()))
			err = &tErrors.ErrorTemporaryServerError{}
			return subscriptionLog, err
		}

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
	// var nativeAsset basetxn.Asset = basetxn.NativeAsset{}
	nativeAssetCode := os.Getenv("NATIVE_ASSET_CODE")
	nairaAssetSlice := strings.Split(os.Getenv("NAIRA_ASSET"), ":")
	var ops []basetxn.Operation = make([]basetxn.Operation, 0)
	serviceFee := owner.UserWallets[0].GetPatronFee(gc)
	patronFeeKP, e := evmkeypair.ParseFull(serviceFee.FeeWalletSecretKey)
	if e != nil {
		log.Println("[generatePatronSubscriptionXdr] error fetching account fee wallet for patron fee ", e)
		logDiscordFailedSubscription("[generatePatronSubscriptionXdr] error fetching account fee wallet for patron fee ")
		return "", &tErrors.ErrorTemporaryServerError{}
	}
	sourceAssets := ""
	var errGetEstimate error
	var requiredSourceQuantity, estimatedTrov string
	var path []basetxn.Asset
	if len(patronSubInput.PaymentContractAddress) == 42 {
		sourceAssets = strings.ToUpper(fmt.Sprintf("%v:%v", patronSubInput.PaymentAssetCode, patronSubInput.PaymentContractAddress))
	}

	var asset basetxn.Asset
	if len(patronSubInput.PaymentContractAddress) == 0 {
		asset = basetxn.NativeAsset{}
	} else {
		asset = basetxn.CreditAsset{Code: patronSubInput.PaymentAssetCode, Issuer: patronSubInput.PaymentContractAddress}
	}
	// chanAccount := <-gc.ChannelAccounts
	// defer func(c *evmkeypair.Full) {
	// 	gc.ChannelAccounts <- c
	// }(chanAccount)

	// _, _, _, _, chanSourceAccount, _ := network.BlockchainAccountProperties(gc.BantuExpansionClient, chanAccount.Address(), nativeAsset)

	_, _, nativeBalance, customBalance, sourceAccount, sourceAccountErr := network.BlockchainAccountProperties(gc.BantuExpansionClient, owner.Address, asset)

	if sourceAccountErr != nil {
		log.Println("[generatePatronSubscriptionXdr] error checking account properties on blockchain. Error ", sourceAccountErr)

		return "", sourceAccountErr
	}

	// //get the trov quantity/equivalent needed for the USD from the market.
	// pathInput := swapModel.SwapPathInput{
	// 	SourceAssets:           "TROV:GAXMBPVA2GNG6A3NV6Q664VZASMROS5ZACKSMTPVCRIKPOJIV43A2CTJ",
	// 	DestinationAssetCode:   strings.Split(os.Getenv("DOLLAR_ASSET"), ":")[0],
	// 	DestinationContractAddress: strings.Split(os.Getenv("DOLLAR_ASSET"), ":")[1],
	// 	DestinationAmount:      decimal.NewFromFloat(priceConfig.Price).Truncate(7).String(),
	// }

	//no USD market. get the USD conversion to CNGN. then fetch the TROV/CNGN price from DEX.
	cngnAmount := gc.ConvertUsdToCngn(priceConfig.Price)
	if !cngnAmount.IsPositive() {
		log.Println("[generatePatronSubscriptionXdr] error getting USD-CNGN conversion estimate. Error ")

		return "", &tErrors.ErrorTemporaryServerError{}
	}

	//get the source quantity/equivalent needed for the CNGN amount we now have from the market since there is CNGN offer of the asset.
	pathInput := swapModel.SwapPathInput{
		SourceAssets:               sourceAssets,
		DestinationAssetCode:       nairaAssetSlice[0],
		DestinationContractAddress: nairaAssetSlice[1],
		DestinationAmount:          cngnAmount.Truncate(7).String(),
	}
	_, requiredSourceQuantity, errGetEstimate = swaps.GetStrictReceivePaths(pathInput, gc.BantuExpansionClient)

	log.Printf("requires %v %v to convert to %v %v\n", requiredSourceQuantity, patronSubInput.PaymentAssetCode, priceConfig.Price, "USD")
	if errGetEstimate != nil && requiredSourceQuantity == "" {
		log.Println("[generatePatronSubscriptionXdr] error getting required source estimate. Error ", errGetEstimate, requiredSourceQuantity)

		return "", errGetEstimate
	}

	if !asset.IsNative() {
		// log.Println("[generatePatronSubscriptionXdr] error account does not exist on ledger. Error ")
		if customBalance.LessThan(decimal.RequireFromString(requiredSourceQuantity)) {
			return "", &tErrors.ErrorUnderfundedAccount{
				Detail: fmt.Sprintf("You need to add at least %v %v to make up for the subscription fee.", decimal.RequireFromString(requiredSourceQuantity).Sub(customBalance).String(), patronSubInput.PaymentAssetCode),
			}
		}

	} else {
		if nativeBalance.LessThan(decimal.RequireFromString(requiredSourceQuantity)) {
			return "", &tErrors.ErrorUnderfundedAccount{
				Detail: fmt.Sprintf("You need to add at least %v %v to make up for the subscription fee.", decimal.RequireFromString(requiredSourceQuantity).Sub(nativeBalance).String(), nativeAssetCode),
			}
		}
	}
	// //assume the primary wallet does not have trustline to the trov asset. Build the trustline.
	// _, destAccountTrustsDestinationAsset, _, _, _, _ := network.BlockchainAccountProperties(gc.BantuExpansionClient, owner.Address, basetxn.CreditAsset{Code: "TROV", Issuer: "GAXMBPVA2GNG6A3NV6Q664VZASMROS5ZACKSMTPVCRIKPOJIV43A2CTJ"})

	// if !destAccountTrustsDestinationAsset {

	// 	ops = append(ops, &basetxn.ChangeTrust{
	// 		Line:          txnbuild.ChangeTrustAssetWrapper{Asset: basetxn.CreditAsset{Code: "TROV", Issuer: "GAXMBPVA2GNG6A3NV6Q664VZASMROS5ZACKSMTPVCRIKPOJIV43A2CTJ"}},
	// 		Limit:         "900000000000",
	// 		SourceAccount: owner.Address,
	// 	})
	// }
	if !strings.EqualFold(patronSubInput.PaymentAssetCode, "TROV") {

		//get swap the asset amount to TROV.
		pathInput := swapModel.SwapSendPathInput{
			DestinationAssets:     "TROV:" + os.Getenv("TROV_ASSET_CONTRACT_ADDRESS"),
			SourceAssetCode:       patronSubInput.PaymentAssetCode,
			SourceContractAddress: patronSubInput.PaymentContractAddress,
			SourceAmount:          requiredSourceQuantity,
		}

		path, estimatedTrov, errGetEstimate = swaps.GetStrictSendPaths(pathInput, gc)
		log.Printf(" %v %v converts to %v %v\n", requiredSourceQuantity, patronSubInput.PaymentAssetCode, estimatedTrov, "TROV")
		if errGetEstimate != nil && estimatedTrov == "" {
			log.Printf("[generatePatronSubscriptionXdr] error getting required %v estimate. Error %v", patronSubInput.PaymentAssetCode, errGetEstimate)
			return "", errGetEstimate
		}

		{
			//build a swap operation to swap the non-trov asset to trov so that trov can be debited.
			var sendAsset basetxn.Asset
			if len(patronSubInput.PaymentContractAddress) == 0 {
				sendAsset = basetxn.NativeAsset{}
			} else {
				sendAsset = basetxn.CreditAsset{Code: patronSubInput.PaymentAssetCode, Issuer: patronSubInput.PaymentContractAddress}
			}

			ops = append(ops, &basetxn.PathPaymentStrictSend{
				SendAsset:     sendAsset,
				SendAmount:    estimatedTrov,
				Destination:   patronFeeKP.Address(),
				DestAsset:     basetxn.CreditAsset{Code: "TROV", Issuer: os.Getenv("TROV_ASSET_CONTRACT_ADDRESS")},
				DestMin:       "0.0000001",
				Path:          path,
				SourceAccount: owner.Address, //primary wallet
			})
		}
	} else {

		ops = append(ops, &basetxn.Payment{
			Destination:   patronFeeKP.Address(),
			Amount:        requiredSourceQuantity,
			Asset:         basetxn.CreditAsset{Code: "TROV", Issuer: os.Getenv("TROV_ASSET_CONTRACT_ADDRESS")},
			SourceAccount: owner.Address, //primary wallet
		})
	}

	//calculate vat here

	assetLabel := os.Getenv("NATIVE_ASSET_CODE")
	if len(patronSubInput.PaymentAssetCode) > 0 {
		assetLabel = patronSubInput.PaymentAssetCode
	}

	//calculate VAT on the fee amount.
	feeAmount := decimal.RequireFromString(requiredSourceQuantity)
	vatFee := gc.GetVATValue(feeAmount)
	vatRate := decimal.NewFromFloat(gc.GetVATRate()).String()
	patronSubInput.Vat = vatRate
	patronSubInput.VatAmount = decimal.NewFromFloat(vatFee).String()
	amountToPay := feeAmount.Add(decimal.NewFromFloat(vatFee))
	patronSubInput.AmountToPay = amountToPay.String()
	patronSubInput.Messages = append(patronSubInput.Messages, fmt.Sprintf("%v %v will be debited from wallet %v to complete the subscription. This is inclusive of VAT (%v %v)", patronSubInput.AmountToPay, assetLabel, owner.Username, patronSubInput.VatAmount, assetLabel))

	// if len(patronSubInput.PaymentContractAddress) == 0 {
	// 	patronSubInput.Messages = append(patronSubInput.Messages, fmt.Sprintf("%v %v will be debited from wallet %v to complete the subscription.", requiredUsdWorth, nativeAssetCode, owner.Username))

	// } else {

	// 	patronSubInput.Messages = append(patronSubInput.Messages, fmt.Sprintf("%v %v will be debited from wallet %v to complete the subscription.", requiredUsdWorth, patronSubInput.PaymentAssetCode, owner.Username))
	// }

	tx, err := basetxn.NewTransaction(
		basetxn.TransactionParams{
			SourceAccount:        sourceAccount.Address,
			IncrementSequenceNum: true,
			Operations:           ops,
			BaseFee:              3000,
			Memo:                 fmt.Sprintf("Patron %v sub", priceConfig.PatronPackage),
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

package payments

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"trovo-wallet-payment-history-engine/internal/cache"
	paymentModels "trovo-wallet-payment-history-engine/internal/components/payments/models"
	userModels "trovo-wallet-payment-history-engine/internal/components/users/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TrackUserWallet(userWallet userModels.UserWallet, roachDB, db *gorm.DB, trackAddress bool, redisCache *cache.RedisCache) error {
	//get names
	type User struct {
		ID        string  `json:"id"`
		FirstName string  `json:"firstName"`
		LastName  *string `json:"lastName"`
		Corporate uint    `json:"corporate"`
	}
	var user User
	cacheKeyInfo := fmt.Sprintf("userObj %v", userWallet.UserID)

	// search cache for user object
	ok, rawdata := redisCache.GetCachedResultRaw(cacheKeyInfo)

	if ok {

		log.Printf("GetUser[%v], served from cache\n", cacheKeyInfo)
		json.Unmarshal(rawdata, &user)

	} else {
		e := db.Where("id = ?", userWallet.UserID).First(&user).Error
		if e != nil {
			log.Printf("[TrackUserWallet]Error while getting user record from wallet for [%s], %v\n", userWallet.UserID, e)
			return e
		}
	}

	var name, lastName string
	if user.LastName != nil {
		lastName = *user.LastName
	}
	if user.Corporate == 0 {
		name = fmt.Sprintf("%s %s", user.FirstName, lastName)
	} else {
		name = user.FirstName
	}

	// Fetch any existing record for this public key first (rather than deleting it
	// unconditionally before checking), so a wallet that's already tracked gets its
	// alias/name updated in place instead of always being deleted and re-created
	// with a brand new ID on every re-track.
	var existingWallet paymentModels.TrackedWallet
	errMain := roachDB.Where("address = ?", userWallet.ID).First(&existingWallet).Error

	if errMain != nil {
		if errors.Is(errMain, gorm.ErrRecordNotFound) {
			//no user was found. create it now
			newID := uuid.NewString()

			trackedWallet := paymentModels.TrackedWallet{
				ID:          newID,
				Address:     userWallet.ID,
				TempAddress: userWallet.TempAddress,
				Alias:       userWallet.Alias,
				Name:        name,
			}
			errCreate := roachDB.Create(&trackedWallet).Error
			if errCreate != nil {
				log.Printf("[TrackUserWallet]Error while creating Main tracked wallet for %s, %v\n", userWallet.Alias, errCreate)
				return errCreate
			}
			if trackAddress {
				//track public key, first remove it if it exists
				roachDB.Where("address = ?", userWallet.ID).Delete(&paymentModels.TrackedAddress{})

				trackedAddress := paymentModels.TrackedAddress{
					Address: userWallet.ID,
				}
				errAddress := roachDB.Create(&trackedAddress).Error
				if errAddress != nil {
					if !strings.Contains(errAddress.Error(), "constraint") {
						log.Printf("[TrackUserWallet]Error while creating tracked public key for %s, %v\n", userWallet.Alias, errAddress)

					}
				}
			}
			return nil

		} else {

			log.Printf("[TrackUserWallet]Error while creating Main tracked wallet for %s, %v\n", userWallet.Alias, errMain)
			return errMain
		}

	} else {

		// record fetched

		//already exists. record was fetched. update the name and save
		existingWallet.Name = name
		existingWallet.Alias = userWallet.Alias
		existingWallet.TempAddress = userWallet.TempAddress
		eSave := roachDB.Save(&existingWallet).Error
		if eSave != nil {
			log.Printf("[TrackUserWallet]Error while saving existing tracked wallet for %s, %v\n", userWallet.Alias, eSave)
			return eSave
		}

		if trackAddress {
			//track public key
			trackedAddress := paymentModels.TrackedAddress{
				Address: userWallet.ID,
			}
			errAddress := roachDB.Create(&trackedAddress).Error
			if errAddress != nil {
				if !strings.Contains(errAddress.Error(), "constraint") {
					log.Printf("[TrackUserWallet]Error while creating tracked public key after saving existing wallet for %s, %v\n", userWallet.Alias, errAddress)

				}
			}
		}

		return nil
	}

}

func SavePaymentHistory(fromPK, fromAlias, fromName, toPK, toAlias, toName, memo, contractAddress, assetCode, amount, transactionHash, transactionType, PT, ID, sss string, transactionTime time.Time, db *gorm.DB) error {
	var from, to string
	var fromVal, toVal, memoVal, issuerVal *string
	if len(fromAlias) > 1 {
		from = fmt.Sprintf("%s[%s]", fromName, fromAlias)
		fromVal = &from
	}
	if len(toAlias) > 1 {
		to = fmt.Sprintf("%s[%s]", toName, toAlias)
		toVal = &to
	}
	if len(memo) > 0 {
		memoVal = &memo
	}
	if len(contractAddress) > 0 {
		issuerVal = &contractAddress
	}

	if assetCode == "" {
		assetCode = os.Getenv("NATIVE_ASSET_CODE")
	}
	paymentHistory := paymentModels.PaymentHistory{
		ID:                    ID,
		TransactionDate:       transactionTime,
		From:                  fromVal,
		FromAddress:           fromPK,
		To:                    toVal,
		ToAddress:             toPK,
		Memo:                  memoVal,
		ContractAddress:       issuerVal,
		AssetCode:             assetCode,
		Amount:                amount,
		TransactionID:         transactionHash,
		TransactionType:       transactionType,
		PT:                    PT,
		SourceAccountSequence: sss,
	}
	// paymentHistory := paymentModels.PaymentHistory{
	// 	ID:              uuid.NewString(),
	// 	TransactionDate: transactionTime,
	// 	From:            fromVal,
	// 	FromAddress:   fromPK,
	// 	To:              toVal,
	// 	ToAddress:     toPK,
	// 	Memo:            memoVal,
	// 	ContractAddress:     issuerVal,
	// 	AssetCode:       assetCode,
	// 	Amount:          amount,
	// 	TransactionID:   transactionHash,
	// 	TransactionType: transactionType,
	// 	PT:              PT,
	// }

	err := db.Create(&paymentHistory).Error
	if err == nil {
		return nil
	}

	log.Printf("[SavePaymentHistory]Error while creating payment history for [%s][%s] [%+v], %v\n", fromPK, toPK, paymentHistory, err)
	if !strings.Contains(err.Error(), "constraint") {
		//not a duplicate-key situation we know how to recover from (e.g. a transient
		//connection/timeout error) - surface it to the caller instead of crashing the engine.
		return err
	}

	//a unique-constraint violation means a matching record already exists (we've seen this
	//operation before, e.g. after a stream reconnect/replay) - find it and update it in place
	//rather than losing the new values.
	var existing paymentModels.PaymentHistory
	e := db.Where("transaction_id = ? AND id = ?", transactionHash, ID).First(&existing).Error
	if e != nil {
		//couldn't locate the conflicting record via this lookup (e.g. the collision was on a
		//different id/tx hash combination). Nothing more we can safely do here - surface the
		//original creation error rather than silently dropping the payment.
		log.Printf("[SavePaymentHistory]unable to locate the conflicting payment history record for [%s][%s]: %v\n", fromPK, toPK, e)
		return err
	}

	existing.PT = PT
	existing.From = fromVal
	existing.To = toVal
	existing.TransactionType = transactionType
	existing.SourceAccountSequence = sss
	esave := db.Save(&existing).Error
	if esave != nil {
		log.Printf("[SavePaymentHistory]unable to update conflicting payment history record for [%s][%s]: %v\n", fromPK, toPK, esave)
		return esave
	}
	return nil
}

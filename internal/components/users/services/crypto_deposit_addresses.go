package users

import (
	"log"
	"os"
	"strings"
	userModels "trovo-wallet-api/internal/components/users/models"
	dl "trovo-wallet-api/internal/dynamiclinks"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/google/uuid"
)

func GenerateDepositAddresses(wallet *userModels.UserWallet, currency string, gc *sharedconfig.GlobalConfig) (depositAddresses []userModels.CryptoWalletDepositAddress, err error) {
	currency = strings.ToUpper(currency)
	depositAddresses = make([]userModels.CryptoWalletDepositAddress, 0)
	if os.Getenv("ENABLE_CRYPTO_WITHDRAWAL_SERVICE") != "1" {
		return
	}
	sub, e := CreateCryptoSubwalletRequest(wallet, currency, gc)

	if e != nil {
		log.Printf("[GenerateDepositAddresses] error generating deposit address for %v %v error: %v\n[GenerateDepositAddresses] checking if it exits already....\n", wallet.Alias, currency, e)

		//try to get it if it already exists
		sub, e = GetCryptoSubwallet(wallet, currency, gc)
		if e != nil {
			log.Printf("[GenerateDepositAddresses] error fetching deposit addresses for %v %v error: %v\n", wallet.Alias, currency, e)

			err = &tErrors.CustomError{
				Param:      "walletID",
				Err:        "error unable to generate deposit address",
				ErrMessage: "Unable to generate deposit address. Please try again.",
			}
			return
		}

	}
	// var depositAddresses []userModels.CryptoWalletDepositAddress
	dbTX := gc.DB.Begin()
	defer dbTX.Rollback()
	for _, v := range sub.Addresses {

		eCheck := dbTX.Where("trovo_wallet_public_key = ? AND currency = ? AND network = ?", wallet.ID, currency, v.Network).First(&userModels.CryptoWalletDepositAddress{}).Error

		if eCheck == nil {
			//address already exists...skip
			continue
		}
		qrc, _ := dl.GenerateQRCode(v.Address, gc)
		var qrCode *string
		if len(qrc) > 0 {
			qrCode = &qrc
		}
		da := userModels.CryptoWalletDepositAddress{
			ID:                   uuid.NewString(),
			UserID:               wallet.UserID,
			TrovoWalletPublicKey: wallet.ID,
			Currency:             currency,
			DepositAddress:       v.Address,
			Network:              v.Network,
			QRCode:               qrCode,
		}

		depositAddresses = append(depositAddresses, da)
	}
	if len(depositAddresses) > 0 {
		e := dbTX.Create(&depositAddresses).Error
		if e != nil {
			log.Printf("[GenerateDepositAddresses] error saving deposit addresses for %v %v error: %v\n", wallet.Alias, currency, e)
			err = &tErrors.ErrorTemporaryServerError{}
			return
		}
		dbTX.Commit()
		wallet.InvalidateUserCache(gc)
		return depositAddresses, nil
	}
	if len(sub.Addresses) > 0 {
		log.Printf("[GenerateDepositAddresses] No new deposit addresses for %v %v. Retrieved: [%+v] Fetching existing addresses.\n", wallet.Alias, currency, sub.Addresses)
		dbTX.Where("trovo_wallet_public_key = ? AND currency = ?", wallet.ID, currency).First(&depositAddresses)
		wallet.InvalidateUserCache(gc)
		return depositAddresses, nil
	}
	log.Printf("[GenerateDepositAddresses] Unable to get any deposit addresses for %v %v.\n", wallet.Alias, currency)

	err = &tErrors.ErrorTemporaryServerError{}
	return
}

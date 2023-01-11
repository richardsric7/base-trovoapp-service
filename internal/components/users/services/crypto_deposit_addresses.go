package users

import (
	"log"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/google/uuid"
)

func GenerateDepositAddresses(wallet *userModels.UserWallet, currency string, gc *sharedconfig.GlobalConfig) (depositAddresses []userModels.CryptoWalletDepositAddress, err error) {
	depositAddresses = make([]userModels.CryptoWalletDepositAddress, 0)
	sub, e := CreateCryptoSubwalletRequest(wallet, currency, gc)

	if e != nil {
		log.Printf("[GenerateDepositAddresses] error generating deposit address for %v error: %v\n[GenerateDepositAddresses] checking if it exits already....\n", currency, e)

		//try to get it if it already exists
		sub, e = GetCryptoSubwallet(wallet, currency, gc)
		log.Printf("[GenerateDepositAddresses] error fetching deposit addresses for %v error: %v\n", currency, e)
		err = &tErrors.CustomError{
			Param:      "walletID",
			Err:        "error unable to generate deposit address",
			ErrMessage: "Unable to generate deposit address. Please try again.",
		}
		return
	}
	// var depositAddresses []userModels.CryptoWalletDepositAddress
	for _, v := range sub.Addresses {
		da := userModels.CryptoWalletDepositAddress{
			ID:                   uuid.NewString(),
			UserID:               wallet.UserID,
			TrovoWalletPublicKey: wallet.ID,
			Currency:             currency,
			DepositAddress:       v.Address,
			Network:              v.Network,
		}
		depositAddresses = append(depositAddresses, da)
	}
	if len(depositAddresses) > 0 {
		e := gc.DB.Create(&depositAddresses).Error
		if e != nil {
			err = &tErrors.ErrorTemporaryServerError{}
		}
		wallet.InvalidateUserCache(gc)
		return depositAddresses, nil
	}

	err = &tErrors.ErrorTemporaryServerError{}
	return
}

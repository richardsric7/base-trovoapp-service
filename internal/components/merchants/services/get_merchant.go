package merchants

import (
	merchantModels "trovo-wallet-api/internal/components/merchants/models"
	"trovo-wallet-api/internal/errors"

	"gorm.io/gorm"
)

//GetMerchant gets user information
func GetMerchant(ID string, db *gorm.DB) (mInfo merchantModels.Merchant, err error) {

	mInfo, err = GetMerchantInfo(ID, db)

	if err != nil {
		return mInfo, err
	}
	if mInfo.Suspended == 1 {
		return mInfo, &errors.ErrorMerchantIsSuspended{}
	}

	return

}

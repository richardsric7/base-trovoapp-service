package servicelinks

import (
	servicelinkModels "trovo-wallet-api/internal/components/servicelinks/models"
	"trovo-wallet-api/internal/errors"

	"gorm.io/gorm"
)

// GetService gets user information
func GetService(ID, apikey string, db *gorm.DB) (mInfo servicelinkModels.ServiceLink, err error) {

	mInfo, err = GetServiceInfo(ID, apikey, db)

	if err != nil {
		return mInfo, err
	}
	if mInfo.Suspended == 1 {
		return mInfo, &errors.ErrorServiceIsSuspended{}
	}

	return

}

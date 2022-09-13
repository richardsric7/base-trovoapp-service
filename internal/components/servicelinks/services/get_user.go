package servicelinks

import (
	merchantdb "trovo-wallet-api/internal/components/servicelinks/db"
	merchantUserModels "trovo-wallet-api/internal/components/servicelinks/models"

	"gorm.io/gorm"
)

// GetUserFromPrimarySigner fetches the user linked to the primary signer
func GetUserFromPrimarySigner(publicKey string, db *gorm.DB) (user merchantUserModels.User, err error) {

	return merchantdb.GetUserFromPrimarySigner(publicKey, db)

}

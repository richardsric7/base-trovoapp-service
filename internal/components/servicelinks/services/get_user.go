package merchants

import (
	merchantUserModels "trovo-wallet-api/internal/components/servicelinks/models"
	"trovo-wallet-api/internal/sharedconfig"

	"gorm.io/gorm"
)

// GetUser gets user information
func GetUser(ID string, db *gorm.DB, publicKey string, gc *sharedconfig.GlobalConfig) (userInfo merchantUserModels.User, err error) {

	return

}

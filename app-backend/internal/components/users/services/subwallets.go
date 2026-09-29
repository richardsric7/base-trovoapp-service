package users

import (
	userModels "trovo-wallet-api/internal/components/users/models"
	"trovo-wallet-api/internal/sharedconfig"
)

func GetUserWallets(user userModels.User, gc *sharedconfig.GlobalConfig) (wallets []userModels.UserWallet) {
	wallets = make([]userModels.UserWallet, 0)
	gc.DB.Where("user_id = ?", user.ID).Find(&wallets)
	return

}

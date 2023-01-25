package users

import (
	userModels "trovo-wallet-api/internal/components/users/models"

	"trovo-wallet-api/internal/sharedconfig"
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

func GetPatronMembershipPrices(gc *sharedconfig.GlobalConfig) (memberships []userModels.PatronMembershipPrice) {
	memberships = make([]userModels.PatronMembershipPrice, 0)
	gc.DB.Find(&memberships)
	return
}

func GetPatronSubscriptionLogs(username string, gc *sharedconfig.GlobalConfig) (patronSubLogs []userModels.UserPatronSubscriptionLog) {
	patronSubLogs = make([]userModels.UserPatronSubscriptionLog, 0)
	gc.DB.Order("createdAt DESC").Where("username = ?", username).Find(&patronSubLogs)
	return
}

func SubscribeToPatronPackage(signerUser *userModels.User, patronSubInput *userModels.PatronSubscriptionInput, gc *sharedconfig.GlobalConfig) (err error) {

	return
}

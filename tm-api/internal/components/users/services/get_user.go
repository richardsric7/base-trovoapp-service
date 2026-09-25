package users

import (
	users "admin-panel-dashboard/internal/components/users/db"
	conDB "admin-panel-dashboard/internal/db"
	"admin-panel-dashboard/internal/errors"
	models "admin-panel-dashboard/internal/models"
	"os"
	"strings"

	"gorm.io/gorm"
)

// GetUser gets user information
func GetUserRecord(ID string, db *gorm.DB) (*models.User, error) {
	conDB.PrintDBStats("GetUserRecord", db)

	user, err := users.GetUserInfo(ID, db)
	return &user, err

}

// GetUser gets user information
func GetUser(userID string, db *gorm.DB) (userInfo models.UserInfo, err error) {
	conDB.PrintDBStats("GetUser", db)

	user, err := users.GetUserInfo(userID, db)
	if err != nil {
		return userInfo, err
	}
	if user.Suspended == 1 {
		return userInfo, &errors.ErrorUsernameIsSuspended{}
	}

	userInfo = models.UserInfo{
		ID:        user.ID,
		Username:  user.Username,
		Address:   user.Address,
		Email:     user.Email,
		LastName:  user.LastName,
		FirstName: user.FirstName,
		//KYCLevel:              user.KYCLevel,
		//AdminLevel:            user.AdminLevel,
		//Offline:               user.Offline,
		//TelegramNotifications: user.TelegramNotifications,
		//MaxAssetPerOffer:      user.MaxAssetPerOffer,
		//MaxAssetPerOrder:      user.MaxAssetPerOrder,
		Suspended: uint(user.Suspended),
	}

	if user.Mobile != "" {
		userInfo.Mobile = user.Mobile
	}
	// if user.Telegram != nil {
	//	userInfo.Telegram = user.Telegram
	//}
	//if user.ContactPhone != nil {
	//	userInfo.ContactPhone = *user.ContactPhone
	//}

	if user.CountryCode != "" {
		userInfo.CountryCode = user.CountryCode
	}
	// if user.ImageThumbnail != nil {
	//	userInfo.ImageThumbnail = *user.ImageThumbnail
	//}
	if len(strings.ReplaceAll(os.Getenv("TAKER_FEE"), " ", "")) > 0 {
		userInfo.TakerFee = strings.ReplaceAll(os.Getenv("TAKER_FEE"), " ", "")
	} else {
		userInfo.TakerFee = "0"
	}
	userInfo.TradeStats = userInfo.GetMakerStat(db)

	userInfo.TakerReputation = userInfo.GetTakerReputation(db)

	return userInfo, nil

}

// GetUserDetails gets user information
func GetUserDetails(userID string, db *gorm.DB) (user models.User, err error) {
	conDB.PrintDBStats("GetUser", db)

	return users.GetUserInfo(userID, db)
}

// GetUserByTelegramID gets user information
func GetUserByTelegramID(ID int64, db *gorm.DB) (userInfo models.UserInfo, err error) {
	conDB.PrintDBStats("GetUser", db)

	user, err := users.GetUserInfoByTelegramID(ID, db)
	if err != nil {
		return userInfo, err
	}
	if user.Suspended == 1 {
		return userInfo, &errors.ErrorUsernameIsSuspended{}
	}

	userInfo = models.UserInfo{
		ID:        user.ID,
		Username:  user.Username,
		Address:   user.Address,
		Email:     user.Email,
		LastName:  user.LastName,
		FirstName: user.FirstName,
		//KYCLevel:              user.KYCLevel,
		//AdminLevel:            user.AdminLevel,
		//Offline:               user.Offline,
		//TelegramNotifications: user.TelegramNotifications,
		//MaxAssetPerOffer:      user.MaxAssetPerOffer,
		//MaxAssetPerOrder:      user.MaxAssetPerOrder,
		Suspended: uint(user.Suspended),
	}

	if user.Mobile != "" {
		userInfo.Mobile = user.Mobile
	}
	// if user.Telegram != nil {
	//	userInfo.Telegram = user.Telegram
	//}
	//if user.ContactPhone != nil {
	//	userInfo.ContactPhone = *user.ContactPhone
	//}

	if user.CountryCode != "" {
		userInfo.CountryCode = user.CountryCode
	}

	// if user.ImageThumbnail != nil {
	//	userInfo.ImageThumbnail = *user.ImageThumbnail
	//}
	if len(strings.ReplaceAll(os.Getenv("TAKER_FEE"), " ", "")) > 0 {
		userInfo.TakerFee = strings.ReplaceAll(os.Getenv("TAKER_FEE"), " ", "")
	} else {
		userInfo.TakerFee = "0"
	}
	userInfo.TradeStats = userInfo.GetMakerStat(db)

	userInfo.TakerReputation = userInfo.GetTakerReputation(db)

	return

}

// GetUserInoByTelegramID gets user information
func GetUserInfoByTelegramID(ID int64, db *gorm.DB) (userInfo *models.User, err error) {
	conDB.PrintDBStats("GetUserInfoByTelegramID", db)

	userInfo, err = users.GetUserInfoByTelegramID(ID, db)
	if err != nil {
		return nil, err
	}

	return userInfo, nil

}

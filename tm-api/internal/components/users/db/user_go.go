package users

import (
	conDB "admin-panel-dashboard/internal/db"
	trovowalletErrors "admin-panel-dashboard/internal/errors"
	models "admin-panel-dashboard/internal/models"
	"admin-panel-dashboard/internal/observe"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/ecnepsnai/discord"
	"gorm.io/gorm"
)

// GetUserInfo gets user data
func GetUserInf(userID string, db *gorm.DB) (user models.User, err error) {

	conDB.PrintDBStats("GetUserInfo", db)
	// username or ID is supplied
	//userInfo = strings.ToLower(userInfo)
	// GET user by ID
	e := db.Where("id = ?", userID).Error
	if e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			// no user was found
			err = &trovowalletErrors.ErrorUserDoesNotExist{}
			return
		}
		log.Println("[GetUserInfo] error: ", e)
		err = &trovowalletErrors.ErrorTemporaryServerError{}
		return

	}

	// log.Printf("user for %v is %v\n", userInfo, user)
	return user, nil

}
func GetUserInfo(userID string, db *gorm.DB) (user models.User, err error) {

	conDB.PrintDBStats("GetUserInfo", db)

	// Fetch the user by ID
	e := db.Where("id = ?", userID).First(&user).Error
	if e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			// No user was found
			err = &trovowalletErrors.ErrorUserDoesNotExist{}
			return
		}
		log.Println("[GetUserInfo] error: ", e)
		err = &trovowalletErrors.ErrorTemporaryServerError{}
		return
	}

	return user, nil
}

// GetUserInfoByTelegramID gets user data
func GetUserInfoByTelegramID(telegramID int64, db *gorm.DB) (user *models.User, err error) {

	conDB.PrintDBStats("GetUserInfo", db)

	// e returns execution errors

	e := db.Where("telegram = ?", &telegramID).First(&user).Error

	if e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			// no user was found
			err = &trovowalletErrors.ErrorUserDoesNotExist{Username: fmt.Sprintf("%v", telegramID)}
			return nil, err
		}
		log.Println("[GetUserInfo] query error: ", e)
		LogDiscordError(fmt.Sprintf("[GetUserInfo] query error: %v", e))
		err = &trovowalletErrors.ErrorTemporaryServerError{}
		return nil, err

	}

	// log.Printf("user for %v is %v\n", userInfo, user)
	return user, nil

}

// GetUserRecord gets referenced user information
func GetUserRecord(ID string, db *gorm.DB) (*models.User, error) {
	conDB.PrintDBStats("GetUserRecord", db)

	user, err := GetUserInfo(ID, db)
	return &user, err

}

func LogDiscordError(msg string) {
	// Also count it, so the failure shows on the error-rate dashboard
	// instead of only in Discord. Additive: Discord is unchanged.
	observe.RecordHandledFailure(msg)

	discord.WebhookURL = "https://discord.com/api/webhooks/865931042795290636/jObHzZWdnbhX1jomOSZQX8Ip5AXLArh87PI4-ZQ8u6ssnRbZuVdY_iPxz5qoWkHUlZwS"
	if len(os.Getenv("500_ERROR_WEBHOOK")) > 50 {
		discord.WebhookURL = os.Getenv("500_ERROR_WEBHOOK")
	}
	err := discord.Say(msg)
	if err != nil {
		log.Println("[LogDiscordError] error:", err)
		return
	}
}

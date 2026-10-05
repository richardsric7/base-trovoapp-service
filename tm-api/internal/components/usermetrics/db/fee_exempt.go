package usermetrics

import (
	"admin-panel-dashboard/internal/models"
	"errors"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ErrFeeExemptUserNotFound: the username is not a Trovo account.
var ErrFeeExemptUserNotFound = errors.New("user not found")

// ListFeeExemptUsers lists the accounts that pay no platform service fees.
func ListFeeExemptUsers(db *gorm.DB) ([]models.FeeExemptUser, error) {
	var out []models.FeeExemptUser
	err := db.Order("username").Find(&out).Error
	return out, err
}

// AddFeeExemptUser exempts an existing account from platform service
// fees (or updates the reason of one already exempt).
func AddFeeExemptUser(db *gorm.DB, username, reason, addedBy string) (models.FeeExemptUser, error) {
	var user struct{ Username string }
	if err := db.Table("users").Select("username").Where("LOWER(username) = ?", strings.ToLower(strings.TrimSpace(username))).Take(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.FeeExemptUser{}, ErrFeeExemptUserNotFound
		}
		return models.FeeExemptUser{}, err
	}
	row := models.FeeExemptUser{Username: user.Username, Reason: strings.TrimSpace(reason), AddedBy: addedBy}
	err := db.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "username"}}, DoUpdates: clause.AssignmentColumns([]string{"reason", "added_by"})}).Create(&row).Error
	return row, err
}

// RemoveFeeExemptUser makes an account pay platform service fees again.
func RemoveFeeExemptUser(db *gorm.DB, username string) error {
	return db.Where("LOWER(username) = ?", strings.ToLower(strings.TrimSpace(username))).Delete(&models.FeeExemptUser{}).Error
}

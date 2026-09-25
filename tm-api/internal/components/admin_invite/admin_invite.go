package admin_invite

import (
	"admin-panel-dashboard/internal/models"

	"gorm.io/gorm"
)

// FetchAdmins fetches admin users from the database with pagination.
func FetchAdmins(limit, offset int, adminDB *gorm.DB) ([]models.AdminUser, error) {
	var admins []models.AdminUser

	if err := adminDB.Limit(limit).Offset(offset).Find(&admins).Error; err != nil {
		return nil, err
	}

	return admins, nil
}

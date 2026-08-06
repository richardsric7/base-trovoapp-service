package currency

import (
	announcementModels "trovo-wallet-api/internal/components/announcements/models"
	geoModels "trovo-wallet-api/internal/components/users/models"

	"log"

	"gorm.io/gorm"
)

// GetAnnouncements gets announcement based on geo data data
func GetAnnouncements(geoData geoModels.IPAPI, db *gorm.DB) (announcements []announcementModels.Announcement) {
	announcements = make([]announcementModels.Announcement, 0)
	err := db.Order("created_at DESC").Where("expiry::date >= now()::date").
		Where("(level = ? OR level = ? OR level = ? OR level = ?)", "ALL", geoData.CountryCode, geoData.RegionName, geoData.City).
		Find(&announcements).Error
	if err != nil {
		log.Println("[GetAnnouncements] error fetching announcementList:", err)
	}

	return
}

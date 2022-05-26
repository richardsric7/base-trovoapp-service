package currency

import (
	announcementModels "trovo-wallet-api/internal/components/announcements/models"
	geoModels "trovo-wallet-api/internal/components/users/models"

	"log"

	"gorm.io/gorm"
)

//GetAnnouncements gets announcement based on geo data data
func GetAnnouncements(geoData geoModels.IPAPI, db *gorm.DB) (announcements []announcementModels.Announcement, err error) {

	err = db.Where("expiry::date >= now()::date").
		Where(db.Where("level = ?", "ALL").
			Or("level = ?", geoData.CountryCode).
			Or("level = ?", geoData.RegionName).
			Or("level = ?", geoData.City)).
		Order("created_at DESC").Find(&announcements).Error
	if err != nil {
		log.Println("[GetAnnouncements] error fetching announcementList:", err)
	}

	return
}

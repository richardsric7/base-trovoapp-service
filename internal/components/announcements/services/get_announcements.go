package currency

import (
	announcementDB "trovo-wallet-api/internal/components/announcements/db"
	announcementModels "trovo-wallet-api/internal/components/announcements/models"
	geoDB "trovo-wallet-api/internal/components/users/models"

	"gorm.io/gorm"
)

//HandleGetAnnouncement gets the announcements
func HandleGetAnnouncement(ip string, db *gorm.DB) (announcements []announcementModels.Announcement, err error) {
	geoData, _ := geoDB.GetGeoInfo(ip)
	announcements, _ = announcementDB.GetAnnouncements(geoData, db)
	return
}

func GetAppVersion(db *gorm.DB) (appVersion announcementModels.AppVersion) {
	db.First(&appVersion)
	return
}

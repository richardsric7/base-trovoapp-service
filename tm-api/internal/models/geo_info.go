package models

import (
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AppendGeoInfo sets user geoInformation
func (user *User) AppendGeoInfo() *IPAPI {
	fetchedGeoInfo, err := GetGeoInfo(user.PublicIP)
	if err != nil {
		return nil
	}
	user.CountryCode = fetchedGeoInfo.CountryCode
	user.Region = fetchedGeoInfo.Region
	user.RegionName = fetchedGeoInfo.RegionName
	user.City = fetchedGeoInfo.City
	user.ISP = fetchedGeoInfo.Isp
	user.Latitude = fetchedGeoInfo.Lat
	user.Longitude = fetchedGeoInfo.Lon
	return &fetchedGeoInfo
}

// GetGeoInfo sets user geoInformation
func (user *User) GetGeoInfo() (fetchedGeoInfo IPAPI, err error) {
	fetchedGeoInfo, err = GetGeoInfo(user.PublicIP)
	if err != nil {
		return
	}
	return fetchedGeoInfo, nil
}

// AppendGeoInfo sets user geoInformation
func (user *User) LogLastLoginData(db *gorm.DB, fetchedGeoInfo *IPAPI) {
	if fetchedGeoInfo != nil {

		lastLoginData := UserLastLogin{
			ID:          uuid.NewString(),
			Username:    user.Username,
			CreatedAt:   time.Now(),
			CountryCode: &fetchedGeoInfo.CountryCode,
			Region:      &fetchedGeoInfo.Region,
			RegionName:  &fetchedGeoInfo.RegionName,
			City:        &fetchedGeoInfo.City,
			ISP:         &fetchedGeoInfo.Isp,
			Latitude:    &fetchedGeoInfo.Lat,
			Longitude:   &fetchedGeoInfo.Lon,
			PublicIP:    &fetchedGeoInfo.IP,
		}

		err := db.Create(&lastLoginData).Error
		if err != nil {
			log.Println("[LogLastLoginData]error saving user login data,err:", err)
		}
	}

}

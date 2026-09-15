package users

//AppendGeoInfo sets user geoInformation
func (user *User) AppendGeoInfo() {
	fetchedGeoInfo, err := GetGeoInfo(user.PublicIP)
	if err != nil {
		return
	}
	user.CountryCode = &fetchedGeoInfo.CountryCode
	user.Region = &fetchedGeoInfo.Region
	user.RegionName = &fetchedGeoInfo.RegionName
	user.City = &fetchedGeoInfo.City
	user.ISP = &fetchedGeoInfo.Isp
	user.Latitude = &fetchedGeoInfo.Lat
	user.Longitude = &fetchedGeoInfo.Lon
}

//GetGeoInfo sets user geoInformation
func (user *User) GetGeoInfo() (fetchedGeoInfo IPAPI, err error) {
	fetchedGeoInfo, err = GetGeoInfo(user.PublicIP)
	if err != nil {
		return
	}
	return fetchedGeoInfo, nil
}

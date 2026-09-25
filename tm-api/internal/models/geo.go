package models

import (
	bantuerrors "admin-panel-dashboard/internal/errors"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
)

// GetGeoInfo gives the Geo Information
func GetGeoInfo(ip string) (fetchedGeoIP IPAPI, err error) {
	// {"status":"success","country":"Nigeria","countryCode":"NG","region":"LA","regionName":"Lagos","city":"Lagos","zip":"","lat":6.4474,"lon":3.3903,"timezone":"Africa/Lagos","isp":"Emerging Markets Telecommunication Services (EMTS) Limited","org":"EMTS Commercial","as":"AS37076 Emerging Markets Telecommunication Services (EMTS) Limited","query":"41.190.3.70"}
	if ip == "" {
		log.Println("[GetGeoIP] no ip address supplied")
		err = &bantuerrors.ErrorTemporaryServerError{}
		return
	}
	if os.Getenv("IPAPI_KEY") == "" {
		log.Println("[GetGeoIP] no IPAPI_KEY supplied")
		err = errors.New("could not authenticate to geoip service")
		return
	}
	if os.Getenv("IPAPI_HOST") == "" {
		log.Println("[GetGeoIP] no IPAPI_HOST supplied")
		err = &bantuerrors.ErrorTemporaryServerError{}
		return
	}

	resp, err := http.Get(os.Getenv("IPAPI_HOST") + "/" + ip + "?key=" + os.Getenv("IPAPI_KEY"))

	if err == nil {
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if len(string(body)) == 0 {
			log.Println("error getting response for request")
			return fetchedGeoIP, &bantuerrors.ErrorTemporaryServerError{}
		}
		if err == nil {
			// fmt.Println("[PRINTING STRING BODY]", string(body))
			// fetchedGeoIP := models.IPAPI{}
			err = json.Unmarshal(body, &fetchedGeoIP)
			if err != nil {
				// fmt.Println("[PRINTING STRING BODY]:", string(body))
				log.Printf("[GetGeoIP] error unmarshalling fetched geo data for ip: %s: error: %v, body: %s\n\n", ip, err, string(body))
				return fetchedGeoIP, &bantuerrors.ErrorTemporaryServerError{}
			}
			fetchedGeoIP.IP = ip
			return fetchedGeoIP, nil
		}
		log.Printf("[GetGeoIP] read response error for ip: %s, %v\n", ip, err)
		return IPAPI{}, &bantuerrors.ErrorTemporaryServerError{}

	}
	log.Printf("[GetGeoIP] get GeoData from url for %s, error: %v\n", ip, err)

	return IPAPI{}, &bantuerrors.ErrorTemporaryServerError{}
}

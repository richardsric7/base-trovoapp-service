package payments

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	tErrors "trovo-wallet-payment-history-engine/internal/errors"
)

//GetGeoInfo gives the Geo Information
func GetGeoInfo(ip string) (fetchedGeoIP IPAPI, err error) {
	if ip == "" {
		log.Println("[GetGeoIP] no ip address supplied")
		err = &tErrors.ErrorTemporaryServerError{}
		return
	}
	if os.Getenv("IPAPI_KEY") == "" {
		log.Println("[GetGeoIP] no IPAPI_KEY supplied")
		err = errors.New("could not authenticate to geo-ip service")
		return
	}
	if os.Getenv("IPAPI_HOST") == "" {
		log.Println("[GetGeoIP] no IPAPI_HOST supplied")
		err = &tErrors.ErrorTemporaryServerError{}
		return
	}

	resp, err := http.Get(os.Getenv("IPAPI_HOST") + "/" + ip + "?key=" + os.Getenv("IPAPI_KEY"))

	if err == nil {
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if len(string(body)) == 0 {
			log.Println("error getting response for request")
			return fetchedGeoIP, &tErrors.ErrorTemporaryServerError{}
		}
		if err == nil {
			// fmt.Println("[PRINTING STRING BODY]", string(body))
			// fetchedGeoIP := models.IPAPI{}
			err = json.Unmarshal(body, &fetchedGeoIP)
			if err != nil {
				// fmt.Println("[PRINTING STRING BODY]:", string(body))
				log.Printf("[GetGeoIP] error unmarshaling fetched geo data for ip: %s: error: %v, body: %s\n\n", ip, err, string(body))
				return fetchedGeoIP, &tErrors.ErrorTemporaryServerError{}
			}
			fetchedGeoIP.IP = ip
			return fetchedGeoIP, nil
		}
		log.Printf("[GetGeoIP] read response error for ip: %s, %v\n", ip, err)
		return IPAPI{}, &tErrors.ErrorTemporaryServerError{}

	}
	log.Printf("[GetGeoIP] get GeoData from url for %s, error: %v\n", ip, err)

	return IPAPI{}, &tErrors.ErrorTemporaryServerError{}
}

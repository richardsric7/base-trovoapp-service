package trovosdk

import (
	"errors"
	"os"
)

func NewServiceLink() (*ServiceLink, error) {
	serviceLink := new(ServiceLink)
	apiBaseUrl := os.Getenv("TROVO_WALLET_BASE_URL")
	serviceLinkUsername := os.Getenv("SERVICE_LINK_USERNAME")
	apiKey := os.Getenv("SERVICE_LINK_API_KEY")
	if len(apiBaseUrl) == 0 {
		return serviceLink, errors.New("no base url specified")
	}

	if len(serviceLinkUsername) == 0 {
		return serviceLink, errors.New("no service link username specified")
	}
	if len(apiKey) == 0 {
		return serviceLink, errors.New("no service link api key specified")
	}
	serviceLink.ApiBaseUrl = apiBaseUrl
	serviceLink.ApiKey = apiKey
	serviceLink.ServiceUsername = serviceLinkUsername

	return serviceLink, nil
}

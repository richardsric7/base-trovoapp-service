package users

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	userModels "trovo-wallet-api/internal/components/users/models"
	"trovo-wallet-api/internal/sharedconfig"
)

// Function to call the API
func StablerailSaveSupportedBanks(gc *sharedconfig.GlobalConfig) (*userModels.StablerailGetBanksResponse, error) {
	// get config
	var config userModels.StablerailConfig
	gc.DB.First(&config)
	if len(config.ApiKey) == 0 {
		log.Println("[StablerailSaveSupportedBanks] Stablerail configuration not found.")
		return nil, fmt.Errorf("no stablerail config found: %v", 404)
	}
	apiKey, baseUrl := config.ApiKey, config.BaseUrl
	log.Println("[StablerailSaveSupportedBanks] Fetching Stablerail Supported banks")

	// url := "https://beta.stablesrail.io/v1/getbankscode"
	if len(baseUrl) == 0 {
		baseUrl = "https://beta.stablesrail.io/v1"
	}
	url := baseUrl + "/getbankscode"
	// Create HTTP client
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	// Create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add headers
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Decode response
	var result userModels.StablerailGetBanksResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	//save bank codes to DB
	for _, bank := range result.Data {
		gc.DB.Save(&bank)
		log.Printf("[StablerailSaveSupportedBanks] Saved Stablerail Supported bank %+v\n", bank)

	}

	return &result, nil
}

func GetStablerailBanks(gc *sharedconfig.GlobalConfig) (banks []userModels.StablerailBank) {
	banks = make([]userModels.StablerailBank, 0)
	gc.DB.Order("bank_name ASC").Find(&banks)
	return
}

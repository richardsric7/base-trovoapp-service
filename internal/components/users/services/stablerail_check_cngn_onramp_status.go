package users

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	userModels "trovo-wallet-api/internal/components/users/models"
	"trovo-wallet-api/internal/sharedconfig"
)

// Function to call the API
func GetCNGNOnrampStatus(requestID string, gc *sharedconfig.GlobalConfig) (*userModels.CNGNOnrampStatusAPIResponse, error) {
	// url := "https://beta.stablesrail.io/v1/cngnonrampstatus"
	// get config
	var config userModels.StablerailConfig
	gc.DB.First(&config)
	if len(config.ApiKey) == 0 {
		log.Println("[StablerailOnboardUser] Stablerail configuration not found.")
		return nil, fmt.Errorf("no stablerail config found: %v", 404)
	}
	apiKey, baseUrl := config.ApiKey, config.BaseUrl

	if len(baseUrl) == 0 {
		baseUrl = "https://beta.stablesrail.io/v1"
	}
	url := baseUrl + "/cngnonrampstatus"

	// Prepare request body
	payload := userModels.CNGNOnrampStatusRequest{
		RequestID: requestID,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// Create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	// HTTP client
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode response
	var apiResp userModels.CNGNOnrampStatusAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	// Optional: check HTTP status
	if resp.StatusCode != http.StatusOK {
		return &apiResp, fmt.Errorf("non-200 response: %d", resp.StatusCode)
	}

	return &apiResp, nil
}

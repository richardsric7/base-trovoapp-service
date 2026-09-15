package users

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	userModels "trovo-wallet-api/internal/components/users/models"
	"trovo-wallet-api/internal/sharedconfig"
)

// Function to initiate offramp
func StableRailInitiateAssetWithdrawal(reqData *userModels.StablerailAssetWithdrawalRequest, gc *sharedconfig.GlobalConfig) (*userModels.StablerailUserAssetWithdrawalResponse, error) {
	if userModels.IsInternalBalanceAssetCode(reqData.Ticker, gc) {
		return nil, fmt.Errorf("%v cannot be withdrawn externally", reqData.Ticker)
	}

	// get config
	var config userModels.StablerailConfig
	if len(reqData.Status) == 0 {
		reqData.Status = "pending"
	}
	e := gc.DB.Save(reqData).Error
	if e != nil {
		//failed to save withdrawal
		log.Printf("[StableRailInitiateAssetWithdrawal] unable to create deposit request %+v\n", *reqData)
		gc.LogDiscordFailedRequest(fmt.Sprintf("[StableRailInitiateAssetWithdrawal] unable to create deposit request %+v\n", *reqData))
		return nil, fmt.Errorf("%v", "unable to create deposit request")
	}

	gc.DB.First(&config)
	if len(config.ApiKey) == 0 {
		log.Println("[StableRailInitiateOfframp] Stablerail configuration not found.")
		return nil, fmt.Errorf("no stablerail config found: %v", 404)
	}
	if config.EnableStablerail == 0 {
		return nil, fmt.Errorf("stablerail not enabled: %v", 202)
	}
	apiKey, baseUrl := config.ApiKey, config.BaseUrl

	if len(baseUrl) == 0 {
		baseUrl = "https://beta.stablesrail.io/v1"
	}
	url := baseUrl + "/withdrawasset"

	// Convert request struct to JSON
	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	// HTTP client
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	// Parse JSON response
	var wtdResp userModels.StablerailUserAssetWithdrawalResponse
	if err := json.Unmarshal(body, &wtdResp); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}
	//save deposit status
	if wtdResp.ResponseCode == "00" {
		reqData.Status = "funded"
		e = gc.DB.Save(reqData).Error
		if e != nil {
			//failed to update user asset withdrawal
			log.Printf("[StableRailInitiateAssetWithdrawal] unable to save deposit request %+v\n", *reqData)
			gc.LogDiscordFailedRequest(fmt.Sprintf("[StableRailInitiateAssetWithdrawal] unable to save deposit request %+v\n", *reqData))
			return nil, fmt.Errorf("%v", "unable to save deposit request")
		}
	}

	return &wtdResp, nil
}

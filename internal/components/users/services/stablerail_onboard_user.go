package users

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	userModels "trovo-wallet-api/internal/components/users/models"
	"trovo-wallet-api/internal/sharedconfig"
)

func StablerailInitiateOnboardUser(trovoUsername, bvn string, gc *sharedconfig.GlobalConfig) (msg string, err error) {
	//check if username already exists
	var stablerailUser userModels.StablerailUser
	gc.DB.Where("trovo_username = ?", trovoUsername).First(&stablerailUser)
	if len(stablerailUser.ID) > 0 {
		//alrready registered
		return "registration already done", nil
	}
	//initiate onboarding
	// tx := gc.DB.Begin()
	// defer tx.Rollback()
	var stablerailRequest userModels.StablerailRequest
	r, err := StablerailOnboardUser(bvn, gc)
	if err != nil {
		return "", err
	}
	//begin user registration
	if !strings.EqualFold(r.ResponseCode, "00") {
		// was not successful
		return "", fmt.Errorf("could not onboard user %s", trovoUsername)
	}
	stablerailRequest = userModels.StablerailRequest{
		ID:            r.Data.RequestID,
		RequestType:   "Onboarding",
		Status:        r.Data.Status,
		TrovoUsername: trovoUsername,
	}
	e := gc.DB.Save(&stablerailRequest).Error
	if e != nil {
		log.Printf("[StablerailInitiateOnboardUser]Error saving onboarding for username %v: %v\n", trovoUsername, e)
		return "", fmt.Errorf("could not onboard user %s", trovoUsername)
	}
	return r.Data.Message, nil

}

// Function to onboard user
func StablerailOnboardUser(bvn string, gc *sharedconfig.GlobalConfig) (*userModels.StablerailOnboardResponse, error) {
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
	url := baseUrl + "/onboarduser"

	//manage stablerails config from DB

	// Prepare request body
	reqBody := userModels.StablerailOnboardRequest{
		BVN: bvn,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP client
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	// Create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Handle non-200 responses
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse JSON response
	var result userModels.StablerailOnboardResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

func UpdateStablerailUserOnboardingStatus(r userModels.StablerailRequest, gc *sharedconfig.GlobalConfig) (err error) {

	//initiate onboarding
	// tx := gc.DB.Begin()
	// defer tx.Rollback()
	// var stablerailRequest userModels.StablerailRequest
	res, err := StablerailCheckOnboardingStatus(r.ID, gc)
	if err != nil {
		return err
	}
	//begin user registration
	if !strings.EqualFold(res.ResponseCode, "00") {
		// was not successful
		return fmt.Errorf("could not check onboarding status user %s", r.TrovoUsername)
	}

	r.Status = res.Data.Status
	e := gc.DB.Save(&r).Error
	if e != nil {
		log.Printf("[UpdateStablerailUserOnboardingStatus]Error saving onboarding status for username %v: %v\n", r.TrovoUsername, e)
		gc.LogDiscordFailedRequest(fmt.Sprintf("[UpdateStablerailUserOnboardingStatus]Error saving onboarding status for username %v: %v\n", r.TrovoUsername, e))
		return fmt.Errorf("could not onboard user %s", r.TrovoUsername)
	}
	//if the response is successful, then we insert into user table
	if strings.EqualFold(res.Data.Status, "completed") {
		// it has been completed
		su := userModels.StablerailUser{
			ID:            res.Data.UserID,
			TrovoUsername: r.TrovoUsername,
		}
		e := gc.DB.Save(&su).Error
		if e != nil {
			log.Printf("[UpdateStablerailUserOnboardingStatus]Error saving stablerail user record for username %v: %v\n", r.TrovoUsername, e)
			gc.LogDiscordFailedRequest(fmt.Sprintf("[UpdateStablerailUserOnboardingStatus]Error saving stablerail user record for username %v: %v\n", r.TrovoUsername, e))

			return fmt.Errorf("could not onboard user %s", r.TrovoUsername)
		}
		//send PN
		user, _ := userModels.Username(r.TrovoUsername).GetSimpleUser(gc.DB, gc)
		dataPayload := make(map[string]string)
		dataPayload["route"] = ""
		user.SendPushMessage("BNV verified for fiat operations", "your BVN has now been verified for fiat operations. You can now go ahead to deposit or make withdrawals", "", dataPayload, gc)
	}
	return nil

}

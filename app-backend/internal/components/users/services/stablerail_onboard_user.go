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

func StablerailInitiateOnboardUser(trovoUser *userModels.User, bvn string, gc *sharedconfig.GlobalConfig) (msg string, err error) {
	var config userModels.StablerailConfig
	trovoUsername := trovoUser.Username
	gc.DB.First(&config)
	if len(config.ApiKey) == 0 {
		log.Println("[StablerailOnboardUser] Stablerail configuration not found.")
		return "", fmt.Errorf("no stablerail config found: %v", 404)
	}
	if config.EnableStablerail == 0 {
		return "", fmt.Errorf("stablerail not enabled: %v", 202)
	}

	//check if username already exists
	var stablerailUser userModels.StablerailUser
	gc.DB.Where("trovo_username = ?", trovoUsername).First(&stablerailUser)
	if len(stablerailUser.ID) > 0 {
		//alrready registered
		return "registration already done", nil
	}
	// a repeated KYC callback must not start a second onboarding
	var inProgress int64
	gc.DB.Model(&userModels.StablerailRequest{}).Where("trovo_username = ? AND request_type = ? AND status IN ?", trovoUsername, "Onboarding", []string{"processing", "pending", "requested"}).Count(&inProgress)
	if inProgress > 0 {
		return "registration in progress", nil
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
	//send PN
	dataPayload := make(map[string]string)
	dataPayload["route"] = ""
	trovoUser.SendPushMessage("Onboarding your BVN for fiat transactions", "Please wait while we onboard your BVN for fiat operations.", "", dataPayload, gc)
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
	if config.EnableStablerail == 0 {
		return nil, fmt.Errorf("stablerail not enabled: %v", 202)
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
	var config userModels.StablerailConfig
	gc.DB.First(&config)
	if len(config.ApiKey) == 0 {
		log.Println("[StablerailOnboardUser] Stablerail configuration not found.")
		return fmt.Errorf("no stablerail config found: %v", 404)
	}
	if config.EnableStablerail == 0 {
		return fmt.Errorf("stablerail not enabled: %v", 202)
	}
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
		user.SendPushMessage("BVN verified for bank transactions", "Your BVN has been verified for bank transactions. You can now deposit Naira from, and withdraw to, your bank account.", "", dataPayload, gc)
	}
	if strings.EqualFold(res.Data.Status, "failed") {
		gc.LogDiscordFailedRequest(fmt.Sprintf("[UpdateStablerailUserOnboardingStatus] Stablerail onboarding %v of %v failed: %v", r.ID, r.TrovoUsername, res.Data.Message))
		if user, e := userModels.Username(r.TrovoUsername).GetSimpleUser(gc.DB, gc); e == nil {
			user.SendPushMessage("Bank transactions not enabled", "Your BVN could not be verified for bank deposits and withdrawals. Please contact support.", "", map[string]string{"route": ""}, gc)
		}
	}
	return nil

}

func GetStablerailPendingOnboardingRequests(gc *sharedconfig.GlobalConfig) (req []userModels.StablerailRequest) {
	req = make([]userModels.StablerailRequest, 0)
	gc.DB.Where("request_type= ? AND status IN ?", "Onboarding", []string{"processing", "pending", "requested"}).Find(&req)
	return
}
func ProcessUpdateStablerailOnboardingStatus(gc *sharedconfig.GlobalConfig) {
	log.Printf("[ProcessUpdateStablerailOnboardingStatus] <<<<<STARTING PROCESSING STABLERAIL USER ONBOARDING STATUS>>>>>\n")

	reqs := GetStablerailPendingOnboardingRequests(gc)

	for _, req := range reqs {
		perror := UpdateStablerailUserOnboardingStatus(req, gc)
		if perror != nil {
			log.Printf("[ProcessUpdateStablerailOnboardingStatus] error updating stablerail user onboarding status: %v\n", perror)
		}
	}
}

// MaskBVN keeps only the last 3 digits of a BVN, for logs and alerts.
func MaskBVN(bvn string) string {
	if len(bvn) <= 3 {
		return "***"
	}
	return strings.Repeat("*", len(bvn)-3) + bvn[len(bvn)-3:]
}

// maxOnboardingRetries is how often a failed onboarding trigger from the KYC
// callback is retried (about every 5 minutes) before it is left to support.
const maxOnboardingRetries = 24

// ProcessStablerailOnboardingRetries retries the Stablerail onboardings the
// KYC callback could not start (Stablerail down or misconfigured at the
// time). A retry that starts the onboarding, or finds the user already
// onboarded, is removed; the BVN is kept no longer than that.
func ProcessStablerailOnboardingRetries(gc *sharedconfig.GlobalConfig) {
	if !StablerailEnabled(gc) {
		return
	}
	var retries []userModels.StablerailOnboardUserRetry
	gc.DB.Where("attempts < ? AND updated_at < ?", maxOnboardingRetries, time.Now().Add(-5*time.Minute)).Order("id").Limit(50).Find(&retries)
	for _, r := range retries {
		user, err := userModels.Username(r.TrovoUsername).GetSimpleUser(gc.DB, gc)
		if err != nil {
			continue
		}
		if _, err := StablerailInitiateOnboardUser(&user, r.BVN, gc); err == nil {
			gc.DB.Delete(&userModels.StablerailOnboardUserRetry{}, r.ID)
			continue
		} else {
			attempts := r.Attempts + 1
			gc.DB.Model(&userModels.StablerailOnboardUserRetry{}).Where("id = ?", r.ID).Updates(map[string]interface{}{"attempts": attempts, "updated_at": time.Now()})
			if attempts >= maxOnboardingRetries {
				gc.LogDiscordFailedRequest(fmt.Sprintf("[ProcessStablerailOnboardingRetries] Stablerail onboarding of %v (BVN %v) still failing after %d attempts: %v", r.TrovoUsername, MaskBVN(r.BVN), attempts, err))
			}
		}
	}
}

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

// Function to initiate offramp
func StableRailInitiateOfframp(reqData userModels.StablerailOfframpRequest, gc *sharedconfig.GlobalConfig) (*userModels.StableRailOfframpResponse, error) {
	// get config
	var config userModels.StablerailConfig
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
	url := baseUrl + "/cngnofframp"

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
	var offrampResp userModels.StableRailOfframpResponse
	if err := json.Unmarshal(body, &offrampResp); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	// record it, so UpdateOfframpStatuses follows it to the end
	if id := offrampResp.Data.RequestID; id != "" {
		rec := userModels.StablerailOfframp{
			ID: id, BankCode: reqData.BankCode, AccountNumber: reqData.AccountNumber, BaseAmount: float64(reqData.Amount),
			Ticker: reqData.Ticker, Status: nonEmpty(offrampResp.Data.Status, "PENDING"), TrovoUsername: reqData.UserID,
		}
		if e := gc.DB.Create(&rec).Error; e != nil {
			log.Printf("[StableRailInitiateOfframp] recording offramp %v: %v", id, e)
		}
	}
	return &offrampResp, nil
}

// offrampFinal reports whether an offramp status is final.
func offrampFinal(status string) bool {
	switch strings.ToUpper(strings.TrimSpace(status)) {
	case "COMPLETED", "COMPLETE", "SUCCESS", "SUCCESSFUL", "FAILED", "FAILURE", "CANCELLED", "CANCELED", "REJECTED", "REVERSED":
		return true
	}
	return false
}

// UpdateOfframpStatuses asks Stablerail for the status of every offramp not
// yet finished, records it, and tells the user when one finishes.
func UpdateOfframpStatuses(gc *sharedconfig.GlobalConfig) {
	var pending []userModels.StablerailOfframp
	gc.DB.Where("created_at > ?", time.Now().Add(-30*24*time.Hour)).Order("created_at").Limit(100).Find(&pending)
	for _, o := range pending {
		if offrampFinal(o.Status) {
			continue
		}
		res, err := GetOfframpStatus(o.ID, gc)
		if err != nil || res == nil || res.Data.Status == "" {
			continue
		}
		status := res.Data.Status
		if strings.EqualFold(status, o.Status) {
			continue
		}
		gc.DB.Model(&userModels.StablerailOfframp{}).Where("id = ?", o.ID).Update("status", status)
		if !offrampFinal(status) {
			continue
		}
		u, err := userModels.Username(o.TrovoUsername).GetSimpleUser(gc.DB, gc)
		if err != nil {
			continue
		}
		title, body := "Bank withdrawal completed", fmt.Sprintf("Your withdrawal of %v %v to your bank account has been paid.", o.BaseAmount, o.Ticker)
		if f := strings.ToUpper(status); f != "COMPLETED" && f != "COMPLETE" && f != "SUCCESS" && f != "SUCCESSFUL" {
			title, body = "Bank withdrawal failed", fmt.Sprintf("Your withdrawal of %v %v to your bank account could not be completed (%v). Please contact support.", o.BaseAmount, o.Ticker, status)
			gc.LogDiscordFailedRequest(fmt.Sprintf("[UpdateOfframpStatuses] offramp %v of %v ended %v", o.ID, o.TrovoUsername, status))
		}
		u.SendPushMessage(title, body, "", map[string]string{"route": "basicTransactionHistory"}, gc)
	}
}

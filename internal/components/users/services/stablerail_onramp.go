package users

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	userModels "trovo-wallet-api/internal/components/users/models"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/shopspring/decimal"
)

// Function to call the API
func CreateCNGNRequest(payload userModels.CNGNOnrampRequest, gc *sharedconfig.GlobalConfig) (*userModels.CNGNOnrampResponse, error) {
	// url := "https://beta.stablesrail.io/v1/cngnonramp"
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
	url := baseUrl + "/cngnonramp"
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	var result userModels.CNGNOnrampResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return &result, fmt.Errorf("non-200 response: %d", resp.StatusCode)
	}

	return &result, nil
}

func UpdateStablerailCNGNOnrampStatus(r userModels.StablerailRequest, gc *sharedconfig.GlobalConfig) (err error) {

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
	res, err := GetCNGNOnrampStatus(r.ID, gc)
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
		log.Printf("[UpdateStablerailCNGNOnrampStatus]Error saving onboarding status for username %v: %v\n", r.TrovoUsername, e)
		gc.LogDiscordFailedRequest(fmt.Sprintf("[UpdateStablerailCNGNOnrampStatus]Error saving onboarding status for username %v: %v\n", r.TrovoUsername, e))
		return fmt.Errorf("could not onboard user %s", r.TrovoUsername)
	}
	//if the response is successful, then we insert into user table
	if strings.EqualFold(res.Data.Status, "funded") {
		// it has been completed
		su := userModels.StablerailOnramp{
			ID:            res.Data.RequestID,
			WalletAddress: res.Data.Wallet.WalletAddress,
			TotalAmount:   decimal.RequireFromString(res.Data.Wallet.Amount).Truncate(7).InexactFloat64(),
			TargetAsset:   res.Data.Wallet.TokenBuy,
			Status:        res.Data.Status,
			AutoSwapEnabled: func() int {
				if res.Data.Wallet.AutoSwap {
					return 1
				}
				return 0
			}(),
			TrovoUsername: r.TrovoUsername,
		}
		e := gc.DB.Save(&su).Error
		if e != nil {
			log.Printf("[UpdateStablerailCNGNOnrampStatus]Error saving stablerail onramp record for username %v: %v\n", r.TrovoUsername, e)
			gc.LogDiscordFailedRequest(fmt.Sprintf("[UpdateStablerailCNGNOnrampStatus]Error saving stablerail onramp record for username %v: %v\n", r.TrovoUsername, e))

			return fmt.Errorf("could not save onramp record for user %s", r.TrovoUsername)
		}
		//send PN
		user, _ := userModels.Username(r.TrovoUsername).GetSimpleUser(gc.DB, gc)
		dataPayload := make(map[string]string)
		dataPayload["route"] = ""
		user.SendPushMessage(fmt.Sprintf("%v successfully purchased!", func() string {
			if res.Data.Wallet.AutoSwap {
				return res.Data.Wallet.TokenBuy
			}
			return "CNGN"
		}()), "Provider has successfully processed your deposit. Now continuing to transfer your token to your TrovoApp blockchain wallet.", "", dataPayload, gc)

		//start user asset withdrawal action

	}
	return nil

}

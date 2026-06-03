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

func StablerailInitiateCNGNOnrampRequest(trovoUser *userModels.User, wallet *userModels.UserWallet, amount float64 /*cngn amount*/, gc *sharedconfig.GlobalConfig) (rv *userModels.GetVirtualAccountResponse, err error) {
	var config userModels.StablerailConfig
	var payload userModels.CNGNOnrampRequest
	trovoUsername := trovoUser.Username
	gc.DB.First(&config)
	if len(config.ApiKey) == 0 {
		log.Println("[StablerailInitiateCNGNOnrampRequest] Stablerail configuration not found.")
		return nil, fmt.Errorf("no stablerail config found: %v", 404)
	}
	if config.EnableStablerail == 0 {
		return nil, fmt.Errorf("stablerail not enabled: %v", 202)
	}
	// check if pending one exists first.
	_, existsError := StablerailGetPendingCNGNOnrampRequest(trovoUser, gc)
	if existsError == nil {

		return nil, fmt.Errorf("%v", "active or pending fiat deposit exists. wait for it to expire or you proceed with it")

	}
	//check if username already exists
	var stablerailUser userModels.StablerailUser
	e := gc.DB.Where("trovo_username = ?", trovoUsername).First(&stablerailUser).Error
	if e != nil {
		//error getting the stablerail user
		return nil, fmt.Errorf("user not onbaorded for fiat: %v", 404)
	}
	//initiate onboarding
	// tx := gc.DB.Begin()
	// defer tx.Rollback()

	var stablerailRequest userModels.StablerailRequest

	payload = userModels.CNGNOnrampRequest{
		Amount:         amount,
		AssetSwap:      "USDC",
		AutoSwap:       false,
		UserID:         stablerailUser.ID,
		SweepToOfframp: true,
	}
	cr, err := CreateCNGNRequest(payload, gc)
	if err != nil {
		return nil, err
	}
	//begin user registration
	if !strings.EqualFold(cr.ResponseCode, "00") {
		// was not successful
		return nil, fmt.Errorf("could not create funding request for %v: %s", trovoUsername, cr.Message)
	}
	body, _ := json.Marshal(payload)
	jsonstring := string(body)
	stablerailRequest = userModels.StablerailRequest{
		ID:                       cr.Data.RequestID,
		RequestType:              "Onramp",
		Status:                   cr.Data.Status,
		StablerailResponseObject: &jsonstring,
		TrovoUsername:            trovoUsername,
	}
	e = gc.DB.Save(&stablerailRequest).Error
	if e != nil {
		log.Printf("[StablerailInitiateCNGNOnrampRequest]Error saving request record for username %v: %v. Request ID: %v\nOnramp Request Response: %+v", trovoUsername, e, cr.Data.RequestID, cr)
		return nil, fmt.Errorf("could not save funding request for user %s", trovoUsername)
	}
	//stablerail onramp
	stablerailOnramp := userModels.StablerailOnramp{
		ID:              cr.Data.RequestID,
		WalletAddress:   wallet.ID,
		TotalAmount:     rv.Data.VirtualAccount.Amount,
		TargetAsset:     "USDC",
		Status:          cr.Data.Status,
		AutoSwapEnabled: 0,
		TrovoUsername:   trovoUsername,
	}
	e = gc.DB.Save(&stablerailOnramp).Error
	if e != nil {
		log.Printf("[StablerailInitiateCNGNOnrampRequest]Error saving onramp record for username %v: %v. Request ID: %v\nOnramp Record: %+v", trovoUsername, e, cr.Data.RequestID, stablerailOnramp)
		return nil, fmt.Errorf("could not save funding request for user %s", trovoUsername)
	}
	//send PN
	dataPayload := make(map[string]string)
	dataPayload["route"] = ""
	trovoUser.SendPushMessage("fiat funding initiated", fmt.Sprintf("Fiat funding for NGN %v has been initiated. Please use the account to be displayed to make payment.", decimal.NewFromFloat(amount).String()), "", dataPayload, gc)
	time.Sleep(5)
	//get virtual account info

	rv, err = GetCNGNOnrampVirtualAccount(cr.Data.RequestID, gc)
	if err != nil {
		return nil, err
	}
	if rv != nil {
		//send the account information in sms.
		vaMsg := fmt.Sprintf("Please pay NGN %v into:\nAccount Number: %v\nBank Name: %v\nAccount Name: %v\nAmount to pay: NGN %v\n", rv.Data.VirtualAccount.Amount, rv.Data.VirtualAccount.AccountNumber, rv.Data.VirtualAccount.BankName, rv.Data.VirtualAccount.AccountName, rv.Data.VirtualAccount.Amount)

		//send PN
		dataPayload := make(map[string]string)
		dataPayload["route"] = ""
		trovoUser.SendPushMessage("funding payment detail", vaMsg, "", dataPayload, gc)
	}
	return rv, nil

}

func StablerailGetPendingCNGNOnrampRequest(trovoUser *userModels.User, gc *sharedconfig.GlobalConfig) (rv userModels.StablerailOnramp, err error) {
	var config userModels.StablerailConfig
	trovoUsername := trovoUser.Username
	gc.DB.First(&config)
	if len(config.ApiKey) == 0 {
		log.Println("[StablerailInitiateCNGNOnrampRequest] Stablerail configuration not found.")
		return rv, fmt.Errorf("no stablerail config found: %v", 404)
	}
	if config.EnableStablerail == 0 {
		return rv, fmt.Errorf("stablerail not enabled: %v", 202)
	}

	//check if username already exists
	var stablerailUserOnramp userModels.StablerailOnramp
	e := gc.DB.Where("trovo_username = ? AND status != ?", trovoUsername, "completed").First(&stablerailUserOnramp).Error
	if e != nil {
		//error getting the stablerail user onramp
		return rv, fmt.Errorf("no pending user fiat deposit: %v", 404)
	}

	return rv, nil

}

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
	if config.EnableStablerail == 0 {
		return nil, fmt.Errorf("stablerail not enabled: %v", 202)
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
		sruser:=GetStablerailUser(r.TrovoUsername, gc)
		sronramp:=GetStablerailOnrampRequestByID(res.Data.RequestID, gc)
		
		if len(sruser.ID)>0 && len(sronramp.ID)>0{
					assetWdlRq := userModels.StablerailAssetWithdrawalRequest{
			ID:     sronramp.ID,
			UserID: res.Data.RequestID,
			InternalWallet: res.Data.Wallet.WalletAddress,
			DestinationWallet: sronramp.WalletAddress,
			Amount: sronramp.TotalAmount,
			Ticker: "CNGN",
			Network: "xbn",
		}

		StableRailInitiateAssetWithdrawal(&assetWdlRq, gc)
		}


	}
	return nil

}

func GetStablerailUser(username string, gc *sharedconfig.GlobalConfig) (req userModels.StablerailUser) {

	gc.DB.Where("Trovo_Username = ?", username).First(&req)
	return
}

func GetStablerailOnrampRequestByID(id string, gc *sharedconfig.GlobalConfig) (req userModels.StablerailOnramp) {

	gc.DB.Where("id = ?", id).First(&req)
	return
}

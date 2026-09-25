package trovosdk

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"

	"admin-panel-dashboard/internal/evmkeypair"

	"github.com/dghubble/sling"
)

func (s *ServiceLink) JwtTokenVerify(jwtToken string) (jwtresponse *JwtVerifyResponse, err error) {

	if len(jwtToken) == 0 {
		err = errors.New("token is empty")
		return
	}

	fullPath := "/v1/servicelinks/token/verify"

	errorResponse := new(ErrorResponse)
	jwtresponse = new(JwtVerifyResponse)
	_, err = sling.New().Set("User-Agent", "TROVO-WALLET Go SDK").
		Set("Authorization", fmt.Sprintf("Bearer %s", jwtToken)).
		Base(s.ApiBaseUrl).
		Post(fullPath).Receive(jwtresponse, errorResponse)
	//get payload string
	if len(errorResponse.Error) > 0 {
		// log.Println("[SendLoginRequest]server response error:", *errorResponse)
		if len(errorResponse.Message) == 0 {
			return nil, errors.New(errorResponse.Error)

		}
		return nil, errors.New(errorResponse.Message)
	}
	if err != nil {
		// log.Println("[SendLoginRequest]request error:", err)
		return nil, err
	}
	// log.Printf("[SendLoginRequest][%+v]\n", *loginInfo)

	return jwtresponse, nil

}

func (s *ServiceLink) JwtTokenDelete(jwtToken string) (jwtresponse *JwtDeleteResponse, err error) {

	if len(jwtToken) == 0 {
		err = errors.New("token is empty")
		return
	}

	fullPath := "/v1/servicelinks/token"

	errorResponse := new(ErrorResponse)
	jwtresponse = new(JwtDeleteResponse)
	_, err = sling.New().Set("User-Agent", "TROVO-WALLET Go SDK").
		Set("Authorization", fmt.Sprintf("Bearer %s", jwtToken)).
		Base(s.ApiBaseUrl).
		Delete(fullPath).Receive(jwtresponse, errorResponse)
	//get payload string
	if len(errorResponse.Error) > 0 {
		// log.Println("[JwtTokenDelete]server response error:", *errorResponse)
		if len(errorResponse.Message) == 0 {
			return nil, errors.New(errorResponse.Error)

		}
		return nil, errors.New(errorResponse.Message)
	}
	if err != nil {
		// log.Println("[JwtTokenDelete]request error:", err)
		return nil, err
	}

	return jwtresponse, nil

}

func (s *ServiceLink) JwtTokenRefresh(refreshToken string) (jwtRefreshResponse *LoginVerifyResponse, err error) {

	// if len(jwtToken) == 0 {
	// 	err = errors.New("token is empty")
	// 	return
	// }
	if len(refreshToken) == 0 {
		err = errors.New("refresh token is empty")
		return
	}
	jsonBody := JwtRefreshTokenInput{
		RefreshToken: refreshToken,
	}
	fullPath := "/v1/servicelinks/token/refresh"

	errorResponse := new(ErrorResponse)
	jwtRefreshResponse = new(LoginVerifyResponse)
	_, err = sling.New().Set("User-Agent", "TROVO-WALLET Go SDK").
		Base(s.ApiBaseUrl).
		Post(fullPath).BodyJSON(jsonBody).Receive(jwtRefreshResponse, errorResponse)
	//get payload string
	if len(errorResponse.Error) > 0 {
		// log.Println("[JwtTokenRefresh]server response error:", *errorResponse)
		if len(errorResponse.Message) == 0 {
			return nil, errors.New(errorResponse.Error)

		}
		return nil, errors.New(errorResponse.Message)
	}
	if err != nil {
		// log.Println("[JwtTokenRefresh]request error:", err)
		return nil, err
	}
	// log.Printf("[JwtTokenRefresh][%+v]\n", *jwtRefreshResponse)

	return jwtRefreshResponse, nil

}

func (s *ServiceLink) SendLoginRequest(trovoUser, authDescription, deviceInfo, callbackUrl string) (loginInfo *LoginWithTrovoWalletData, err error) {

	if len(trovoUser) == 0 {
		err = errors.New("trovoUser is empty")
		return
	}

	// if len(authDescription) == 0 {
	// 	err = errors.New("authDescription is empty")
	// 	return
	// }
	if len(callbackUrl) == 0 {
		err = errors.New("callbackUrl is empty")
		return
	}

	jsonBody := ServiceLinkRequestInput{
		AuthDescription: authDescription,
		DeviceInfo:      deviceInfo,
		CallbackURL:     callbackUrl,
	}

	fullPath := fmt.Sprintf("/v1/servicelinks/login/request/%v", trovoUser)

	errorResponse := new(ErrorResponse)
	loginInfo = new(LoginWithTrovoWalletData)
	_, err = sling.New().Set("User-Agent", "TROVO-WALLET Go SDK").
		Set("X-TW-SERVICE-LINK-API-KEY", s.ApiKey).
		Base(s.ApiBaseUrl).
		Post(fullPath).BodyJSON(jsonBody).Receive(loginInfo, errorResponse)
	//get payload string
	if len(errorResponse.Error) > 0 {
		// log.Println("[SendLoginRequest]server response error:", *errorResponse)
		if len(errorResponse.Message) == 0 {
			return nil, errors.New(errorResponse.Error)

		}
		return nil, errors.New(errorResponse.Message)
	}
	if err != nil {
		// log.Println("[SendLoginRequest]request error:", err)
		return nil, err
	}
	// log.Printf("[SendLoginRequest][%+v]\n", *loginInfo)

	return loginInfo, nil

}

func (s *ServiceLink) VerifyLoginRequest(trovoUser, loginID string) (verifyResp *LoginVerifyResponse, err error) {

	if len(trovoUser) == 0 {
		err = errors.New("trovoUser is empty")
		return
	}

	// if len(authDescription) == 0 {
	// 	err = errors.New("authDescription is empty")
	// 	return
	// }

	if len(loginID) == 0 {
		err = errors.New("loginID is empty")
		return
	}

	fullPath := fmt.Sprintf("/v1/servicelinks/login/verify/%v/%v/%v", s.ServiceUsername, trovoUser, loginID)

	errorResponse := new(ErrorResponse)
	verifyResp = new(LoginVerifyResponse)
	_, err = sling.New().Set("User-Agent", "TROVO-WALLET Go SDK").
		Set("X-TW-SERVICE-LINK-API-KEY", s.ApiKey).
		Base(s.ApiBaseUrl).
		Get(fullPath).Receive(verifyResp, errorResponse)
	//get payload string
	if len(errorResponse.Error) > 0 {
		// log.Println("[VerifyLoginRequest]server response error:", *errorResponse)
		if len(errorResponse.Message) == 0 {
			return nil, errors.New(errorResponse.Error)

		}
		return nil, errors.New(errorResponse.Message)
	}
	if err != nil {
		// log.Println("[VerifyLoginRequest]request error:", err)
		return nil, err
	}
	// log.Printf("[VerifyLoginRequest][%+v]\n", *loginInfo)

	return verifyResp, nil

}

func (s *ServiceLink) VerifyAuthorizationRequest(trovoUser, authID string) (authResponseData *TrovoWalletAuthorizationRespomseData, err error) {

	if len(trovoUser) == 0 {
		err = errors.New("trovoUser is empty")
		return
	}

	// if len(authDescription) == 0 {
	// 	err = errors.New("authDescription is empty")
	// 	return
	// }

	if len(authID) == 0 {
		err = errors.New("loginID is empty")
		return
	}
	//"/v1/servicelinks/authorize/verify/:targetUser"
	fullPath := fmt.Sprintf("/v1/servicelinks/authorize/verify/%v/%v/%v", s.ServiceUsername, trovoUser, authID)

	log.Printf("[VerifyAuthorizationRequest] DEBUG URL: %s%s | ServiceUsername: %s | ApiKey: %s", s.ApiBaseUrl, fullPath, s.ServiceUsername, s.ApiKey)

	errorResponse := new(ErrorResponse)
	authResponseData = new(TrovoWalletAuthorizationRespomseData)
	resp, err := sling.New().Set("User-Agent", "TROVO-WALLET Go SDK").
		Set("X-TW-SERVICE-LINK-API-KEY", s.ApiKey).
		Base(s.ApiBaseUrl).
		Get(fullPath).Receive(authResponseData, errorResponse)
	if resp != nil {
		log.Printf("[VerifyAuthorizationRequest] DEBUG Response Status: %d", resp.StatusCode)
	}
	//get payload string
	if len(errorResponse.Error) > 0 {
		// log.Println("[VerifyAuthorizationRequest]server response error:", *errorResponse)
		if len(errorResponse.Message) == 0 {
			return nil, errors.New(errorResponse.Error)

		}
		return nil, errors.New(errorResponse.Message)
	}
	if err != nil {
		// log.Println("[VerifyAuthorizationRequest]request error:", err)
		return nil, err
	}
	// log.Printf("[VerifyAuthorizationRequest][%+v]\n", *loginInfo)

	return authResponseData, nil

}

// GetUserInfo returns the target user's flat profile - app-backend's
// GET /v1/servicelinks/:ownerUsername/:targetUser/userinfo returns
// servicelinkServices.GetUserForServiceLink's result directly (a flat
// ServiceLinksUser), not wrapped in the {userData, wallet} shape
// ServiceLinkUserInfo/LoginVerifyResponse use for the login flow.
func (s *ServiceLink) GetUserInfo(trovoUser string) (userInfo *ServiceLinksUser, err error) {

	if len(trovoUser) == 0 {
		err = errors.New("trovoUser is empty")
		return
	}

	// if len(authDescription) == 0 {
	// 	err = errors.New("authDescription is empty")
	// 	return
	// }

	// if len(loginID) == 0 {
	// 	err = errors.New("loginID is empty")
	// 	return
	// }

	fullPath := fmt.Sprintf("/v1/servicelinks/%v/%v/userinfo", s.ServiceUsername, trovoUser)

	errorResponse := new(ErrorResponse)
	userInfo = new(ServiceLinksUser)
	_, err = sling.New().Set("User-Agent", "TROVO-WALLET Go SDK").
		Set("X-TW-SERVICE-LINK-API-KEY", s.ApiKey).
		Base(s.ApiBaseUrl).
		Get(fullPath).Receive(userInfo, errorResponse)
	//get payload string
	if len(errorResponse.Error) > 0 {
		// log.Println("[GetUserInfo]server response error:", *errorResponse)
		if len(errorResponse.Message) == 0 {
			return nil, errors.New(errorResponse.Error)

		}
		return nil, errors.New(errorResponse.Message)
	}
	if err != nil {
		// log.Println("[GetUserInfo]request error:", err)
		return nil, err
	}
	// log.Printf("[GetUserInfo][%+v]\n", *loginInfo)

	return userInfo, nil

}

func (s *ServiceLink) GetTokenizedAssetData(assetCode string) (tokenizedAssetData *TokenizedAssetJSON, err error) {

	if len(assetCode) == 0 {
		err = errors.New("assetCode is empty")
		return
	}

	fullPath := fmt.Sprintf("/v1/servicelinks/tokenized-asset/%v", assetCode)

	errorResponse := new(ErrorResponse)
	tokenizedAssetData = new(TokenizedAssetJSON)
	_, err = sling.New().Set("User-Agent", "TROVO-WALLET Go SDK").
		Set("X-TW-SERVICE-LINK-API-KEY", s.ApiKey).
		Base(s.ApiBaseUrl).
		Get(fullPath).Receive(tokenizedAssetData, errorResponse)
	//get payload string
	if len(errorResponse.Error) > 0 {
		// log.Println("[GetTokenizedAssetData]server response error:", *errorResponse)
		if len(errorResponse.Message) == 0 {
			return nil, errors.New(errorResponse.Error)

		}
		return nil, errors.New(errorResponse.Message)
	}
	if err != nil {
		// log.Println("[GetTokenizedAssetData]request error:", err)
		return nil, err
	}
	// log.Printf("[GetTokenizedAssetData][%+v]\n", *loginInfo)

	return tokenizedAssetData, nil

}

func (s *ServiceLink) SendAuthorizationRequest(trovoUser, authDescription, deviceInfo, callbackUrl string, validityInMinutes int) (authData *TrovoWalletAuthorizationData, err error) {

	if len(trovoUser) == 0 {
		err = errors.New("trovoUser is empty")
		return
	}

	// if len(authDescription) == 0 {
	// 	err = errors.New("authDescription is empty")
	// 	return
	// }
	if len(callbackUrl) == 0 {
		err = errors.New("callbackUrl is empty")
		return
	}

	jsonBody := ServiceLinkRequestInput{
		AuthDescription:   authDescription,
		DeviceInfo:        deviceInfo,
		CallbackURL:       callbackUrl,
		ValidityInMinutes: validityInMinutes,
	}

	fullPath := fmt.Sprintf("/v1/servicelinks/authorize/request/%v", trovoUser)

	errorResponse := new(ErrorResponse)
	authData = new(TrovoWalletAuthorizationData)
	_, err = sling.New().Set("User-Agent", "TROVO-WALLET Go SDK").
		Set("X-TW-SERVICE-LINK-API-KEY", s.ApiKey).
		Base(s.ApiBaseUrl).
		Post(fullPath).BodyJSON(jsonBody).Receive(authData, errorResponse)
	//get payload string
	if len(errorResponse.Error) > 0 {
		// log.Println("[SendAuthorizationRequest]server response error:", *errorResponse)
		if len(errorResponse.Message) == 0 {
			return nil, errors.New(errorResponse.Error)

		}
		return nil, errors.New(errorResponse.Message)
	}
	if err != nil {
		// log.Println("[SendAuthorizationRequest]request error:", err)
		return nil, err
	}
	// log.Printf("[SendAuthorizationRequest][%+v]\n", *loginInfo)

	return authData, nil

}

func (s *ServiceLink) GetTokenizedAssetAuthorizationRequest(assetCode, unsignedTransactionXdr string) (authData *ServiceLinkTokenizedAssetAuthRequest, err error) {

	if len(assetCode) == 0 {
		err = errors.New("assetCode is empty")
		return
	}

	if len(unsignedTransactionXdr) == 0 {
		err = errors.New("unsignedTransactionXdr is empty")
		return
	}

	jsonBody := ServiceLinkTokenizedAssetAuthRequest{
		UnsignedTransaction: unsignedTransactionXdr,
		AssetCode:           assetCode,
	}

	fullPath := "/v1/servicelinks/authorize/tokenized-asset"

	errorResponse := new(ErrorResponse)
	authData = new(ServiceLinkTokenizedAssetAuthRequest)
	_, err = sling.New().Set("User-Agent", "TROVO-WALLET Go SDK").
		Set("X-TW-SERVICE-LINK-API-KEY", s.ApiKey).
		Base(s.ApiBaseUrl).
		Post(fullPath).BodyJSON(jsonBody).Receive(authData, errorResponse)
	//get payload string
	if len(errorResponse.Error) > 0 {
		// log.Println("[SendAuthorizationRequest]server response error:", *errorResponse)
		if len(errorResponse.Message) == 0 {
			return nil, errors.New(errorResponse.Error)

		}
		return nil, errors.New(errorResponse.Message)
	}
	if err != nil {
		// log.Println("[SendAuthorizationRequest]request error:", err)
		return nil, err
	}
	// log.Printf("[SendAuthorizationRequest][%+v]\n", *loginInfo)

	return authData, nil

}

func (s *ServiceLink) SendEventRequest(trovoUser, eventDescription, deviceInfo, callbackUrl string, validityInMinutes int) (eventData *TrovoWalletEventData, err error) {

	if len(trovoUser) == 0 {
		err = errors.New("trovoUser is empty")
		return
	}

	// if len(authDescription) == 0 {
	// 	err = errors.New("authDescription is empty")
	// 	return
	// }
	if len(callbackUrl) == 0 {
		err = errors.New("callbackUrl is empty")
		return
	}

	jsonBody := ServiceLinkEventRequestInput{
		EventDescription:  eventDescription,
		DeviceInfo:        deviceInfo,
		CallbackURL:       callbackUrl,
		ValidityInMinutes: validityInMinutes,
	}

	// app-backend's /v1/servicelinks/events/request no longer takes a
	// :targetUser path segment - events aren't targeted at one user, so
	// there is nothing per-user to put in the URL (trovoUser is still
	// required above as an input-validation gate, matching the other
	// Send* methods' shape, even though it isn't sent to the server).
	fullPath := "/v1/servicelinks/events/request"

	errorResponse := new(ErrorResponse)
	eventData = new(TrovoWalletEventData)
	_, err = sling.New().Set("User-Agent", "TROVO-WALLET Go SDK").
		Set("X-TW-SERVICE-LINK-API-KEY", s.ApiKey).
		Base(s.ApiBaseUrl).
		Post(fullPath).BodyJSON(jsonBody).Receive(eventData, errorResponse)
	//get payload string
	if len(errorResponse.Error) > 0 {
		// log.Println("[SendEventRequest]server response error:", *errorResponse)
		if len(errorResponse.Message) == 0 {
			return nil, errors.New(errorResponse.Error)

		}
		return nil, errors.New(errorResponse.Message)
	}
	if err != nil {
		// log.Println("[SendEventRequest]request error:", err)
		return nil, err
	}
	// log.Printf("[SendEventRequest][%+v]\n", *loginInfo)

	return eventData, nil

}

func (s *ServiceLink) SendPushNotification(trovoUser, title, message, imageUri, action string) (pnr *PNSResponse, err error) {

	if len(trovoUser) == 0 {
		err = errors.New("trovoUser is empty")
		return
	}
	if action == "" {
		action = "none"
	}

	if title == "" {
		err = errors.New("message is empty")
		return
	}
	if message == "" {
		err = errors.New("message is empty")
		return
	}

	jsonBody := ServiceLinkPushNotificationInput{
		Title:    title,
		Message:  message,
		ImageURI: imageUri,
		Action:   action, //login, payment, 2fa, event
	}

	// The real push-notification endpoint - was previously pointed at the
	// events/request path, which doesn't accept a push notification body
	// and doesn't identify a target user at all.
	fullPath := fmt.Sprintf("/v1/servicelinks/%v/%v/push", s.ServiceUsername, trovoUser)

	errorResponse := new(ErrorResponse)
	pnr = new(PNSResponse)
	_, err = sling.New().Set("User-Agent", "TROVO-WALLET Go SDK").
		Set("X-TW-SERVICE-LINK-API-KEY", s.ApiKey).
		Base(s.ApiBaseUrl).
		Post(fullPath).BodyJSON(jsonBody).Receive(pnr, errorResponse)
	//get payload string
	if len(errorResponse.Error) > 0 {
		// log.Println("[SendPushNotification]server response error:", *errorResponse)
		if len(errorResponse.Message) == 0 {
			return nil, errors.New(errorResponse.Error)

		}
		return nil, errors.New(errorResponse.Message)
	}
	if err != nil {
		// log.Println("[SendPushNotification]request error:", err)
		return nil, err
	}
	// log.Printf("[SendPushNotification][%+v]\n", *pnr)

	return pnr, nil

}

func (s *ServiceLink) GetPaymentData(trovoUser, paymentDestination, assetCode, contractAddress, amount, memo string) (paymentData *PayWithTrovoWalletData, err error) {

	if len(trovoUser) == 0 {
		err = errors.New("trovoUser is empty")
		return
	}

	if len(paymentDestination) == 0 {
		err = errors.New("authDescription is empty")
		return
	}

	if len(os.Getenv("NATIVE_ASSET_CODE")) == 0 {
		err = errors.New("native asset code is empty")
		return
	}
	if assetCode == "" && contractAddress != "" {
		return nil, errors.New("assetCode is empty but contractAddress was specified")
	}
	if len(assetCode) > 0 {
		if len(assetCode) < 3 || len(assetCode) > 12 {
			return nil, errors.New("assetCode is invalid")
		}

	}
	if contractAddress == "" && assetCode != "" && assetCode != os.Getenv("NATIVE_ASSET_CODE") {
		return nil, errors.New("contractAddress is empty but assetCode is specified")
	}
	if amount == "" || amount == "0" {
		return nil, errors.New("amount is empty")
	}
	if len(contractAddress) > 0 {
		_, e := evmkeypair.ParseAddress(contractAddress)
		if e != nil {
			return nil, errors.New("invalid contractAddress")
		}
	}
	encPaymentDestination := url.QueryEscape(paymentDestination)
	encAssetCode := url.QueryEscape(assetCode)
	encContractAddress := url.QueryEscape(contractAddress)
	encAmount := url.QueryEscape(amount)
	encMemo := url.QueryEscape(memo)

	fullPath := fmt.Sprintf("/v1/servicelinks/payment/request/%v?amount=%v&assetCode=%v&contractAddress=%v&memo=%v&paymentDestination=%v", trovoUser, encAmount, encAssetCode, encContractAddress, encMemo, encPaymentDestination)

	// fullPath := fmt.Sprintf("/v1/servicelinks/payment/request/%v", trovoUser)

	errorResponse := new(ErrorResponse)
	paymentData = new(PayWithTrovoWalletData)
	_, err = sling.New().Set("User-Agent", "TROVO-WALLET Go SDK").
		Set("X-TW-SERVICE-LINK-API-KEY", s.ApiKey).
		Base(s.ApiBaseUrl).
		Get(fullPath).Receive(paymentData, errorResponse)
	//get payload string
	if len(errorResponse.Error) > 0 {
		// log.Println("[GetPaymentData]server response error:", *errorResponse)
		if len(errorResponse.Message) == 0 {
			return nil, errors.New(errorResponse.Error)

		}
		return nil, errors.New(errorResponse.Message)
	}
	if err != nil {
		// log.Println("[GetPaymentData]request error:", err)
		return nil, err
	}
	// log.Printf("[GetPaymentData][%+v]\n", *loginInfo)

	return paymentData, nil

}

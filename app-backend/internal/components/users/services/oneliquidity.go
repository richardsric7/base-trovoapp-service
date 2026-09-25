package users

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	pUrl "net/url"
	"os"
	"strings"
	"trovo-wallet-api/internal/basetxn"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/evmkeypair"
	"trovo-wallet-api/internal/network"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm/clause"
)

func GetCryptoDepositAddresses(wallet *userModels.UserWallet, currency string, gc *sharedconfig.GlobalConfig) (cryptoAddresses []userModels.CryptoWalletDepositAddress) {
	cryptoAddresses = make([]userModels.CryptoWalletDepositAddress, 0)
	e := gc.DB.Where("trovo_wallet_address = ? AND LOWER(currency) = ?", wallet.ID, strings.ToLower(currency)).Find(&cryptoAddresses).Error
	if e != nil {
		log.Printf("[GetCryptoDepositAddresses] error fetching cryptoAddresses from db %v", e)
	}

	return cryptoAddresses

}

func GetCryptoSubwallet(wallet *userModels.UserWallet, currency string, gc *sharedconfig.GlobalConfig) (subwallet userModels.CryptoSubWallet, err error) {

	var wdlResp userModels.CryptoSubwalletResponse

	client := http.DefaultClient
	url := fmt.Sprintf("%s/%s?currency=%s&uid=%s", os.Getenv("ONELIQUIDITY_BASE_URL"), "wallets/v1/sub", currency, wallet.Alias+"@"+os.Getenv("WALLET_DOMAIN"))

	request, err := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("ONELIQUIDITY_TOKEN")))
	resp, err := client.Do(request)
	if err != nil {
		log.Println("[CreateCryptoSubwalletRequest] error sending request:", err)
		return
	}
	if resp.StatusCode != 200 && resp.StatusCode != 201 {

		type ErrorResponse struct {
			Message string `json:"message"`
		}
		var errorResponse ErrorResponse
		defer resp.Body.Close()
		//Decode the data
		if err = json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			log.Println("[CreateCryptoSubwalletRequest] error decoding response:", err)
			return
		}
		log.Printf("[CreateCryptoSubwalletRequest] error response with code: %v, status: %v,error %v", resp.StatusCode, resp.Status, errorResponse.Message)

		err = &tErrors.ErrorTemporaryServerError{}
		return
	}

	defer resp.Body.Close()
	//Decode the data
	if err = json.NewDecoder(resp.Body).Decode(&wdlResp); err != nil {
		log.Println("[CreateCryptoSubwalletRequest] error decoding response:", err)
		return
	}

	return wdlResp.Data, nil

}

func CreateCryptoSubwalletRequest(wallet *userModels.UserWallet, currency string, gc *sharedconfig.GlobalConfig) (subwallet userModels.CryptoSubWallet, err error) {

	var wdlResp userModels.CryptoSubwalletResponse

	client := http.DefaultClient
	url := fmt.Sprintf("%s/%s", os.Getenv("ONELIQUIDITY_BASE_URL"), "wallets/v1/sub")
	jbody, err := json.Marshal(userModels.OnliquiditySubWalletInput{
		Currency: currency,
		UID:      wallet.Alias + "@" + os.Getenv("WALLET_DOMAIN"),
	})
	if err != nil {
		log.Println("[CreateCryptoSubwalletRequest] error sending request:", err)

		return
	}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jbody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("ONELIQUIDITY_TOKEN")))
	resp, err := client.Do(request)
	if err != nil {
		log.Println("[CreateCryptoSubwalletRequest] error sending request:", err)
		return
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		type ErrorResponse struct {
			Message string `json:"message"`
		}
		var errorResponse ErrorResponse
		defer resp.Body.Close()
		//Decode the data
		if err = json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			log.Println("[CreateCryptoSubwalletRequest] error decoding response:", err)
			return
		}
		log.Printf("[CreateCryptoSubwalletRequest] error response with code: %v, status: %v,error %v", resp.StatusCode, resp.Status, errorResponse.Message)

		err = &tErrors.ErrorTemporaryServerError{}
		return
	}

	defer resp.Body.Close()
	//Decode the data
	if err = json.NewDecoder(resp.Body).Decode(&wdlResp); err != nil {
		log.Println("[CreateCryptoSubwalletRequest] error decoding response:", err)
		return
	}

	return wdlResp.Data, nil

}

func GetCryptoSubwalletRequest(wallet *userModels.UserWallet, currency string, gc *sharedconfig.GlobalConfig) (subwallet userModels.CryptoSubWallet, err error) {

	var wdlResp userModels.CryptoSubwalletResponse
	uidParam := pUrl.QueryEscape(wallet.Alias + "@" + os.Getenv("WALLET_DOMAIN"))
	client := http.DefaultClient
	requestUrl := fmt.Sprintf("%s/%s", os.Getenv("ONELIQUIDITY_BASE_URL"), fmt.Sprintf("wallets/v1/sub?currency=%v&uid=%v", pUrl.QueryEscape(currency), uidParam))

	if err != nil {
		log.Println("[GetCryptoSubwalletRequest] error sending request:", err)

		return
	}
	request, err := http.NewRequest(http.MethodGet, requestUrl, nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("ONELIQUIDITY_TOKEN")))
	resp, err := client.Do(request)
	if err != nil {
		log.Println("[GetCryptoSubwalletRequest] error sending request:", err)
		return
	}
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		type ErrorResponse struct {
			Message string `json:"message"`
		}
		var errorResponse ErrorResponse
		defer resp.Body.Close()
		//Decode the data
		if err = json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			log.Println("[GetCryptoSubwalletRequest] error decoding response:", err)
			return
		}
		log.Printf("[GetCryptoSubwalletRequest] error response with code: %v, status: %v,error %v", resp.StatusCode, resp.Status, errorResponse.Message)

		err = &tErrors.ErrorTemporaryServerError{}
		return
	}

	defer resp.Body.Close()
	//Decode the data
	if err = json.NewDecoder(resp.Body).Decode(&wdlResp); err != nil {
		log.Println("[GetCryptoSubwalletRequest] error decoding response:", err)
		return
	}

	return wdlResp.Data, nil

}

func QueueWithdrawalRequest(signerUser *userModels.User, wallet *userModels.UserWallet, wdlInput *userModels.WithdrawalRequestInput, gc *sharedconfig.GlobalConfig) (err error) {

	if wallet.NumberOfApprovalsNeeded > 0 && wallet.SharedAccessEnabled == 1 {
		wdlInput.Multiparty = 1
	}
	if wallet.HasViewOnlyAccess(gc) || wallet.SharedAccessEnabled == 0 {
		wdlInput.SignatureRequired = 1
	}
	log.Printf("[QueueWithdrawalRequest]>:%+v\n", wdlInput)
	ca, err := userModels.Currency(wdlInput.Currency).GetCurratedAsset(gc)
	if err != nil {
		return err
	}
	wdlInput.AmountSubmitted = decimal.NewFromFloat(wdlInput.AmountSubmitted).Truncate(int32(ca.DecimalPlaces)).InexactFloat64()

	//prepare withdrawal figures
	serviceFee := decimal.RequireFromString(os.Getenv("CRYPTO_WITHDRAWAL_SERVICE_FEE"))
	wdlInput.WithdrawalServiceFee = serviceFee.InexactFloat64()
	serviceFeeAmount := (decimal.NewFromFloat(wdlInput.AmountSubmitted).Mul((serviceFee).Div(decimal.NewFromInt(100)))).Truncate(int32(ca.DecimalPlaces))
	//validate input
	wdlAmount := (decimal.NewFromFloat(wdlInput.AmountSubmitted).Sub(serviceFeeAmount)).Truncate(int32(ca.DecimalPlaces))
	wdlInput.AmountToWithdraw = wdlAmount.InexactFloat64()
	wdlNetworks, err := GetWithdrawalNetworks(wdlInput.Currency, gc)
	if err != nil {
		log.Println("[QueueWithdrawalRequest] error getting network request:", err)

		return err
	}
	validNetwork := false
	var wdn userModels.WithdrawalNetwork
	for _, wdn = range wdlNetworks {
		if validNetwork {
			break
		}
		if strings.EqualFold(wdn.Network, wdlInput.WithdrawalNetwork) {
			validNetwork = true
			wdlInput.WithdrawalNetworkFee = decimal.RequireFromString(wdn.WithdrawFee).InexactFloat64()
			//check amount if valid
			if (decimal.NewFromFloat(wdlInput.AmountSubmitted)).LessThan(decimal.RequireFromString(wdn.WithdrawMin)) {

				err = &tErrors.CustomError{
					Param:      "amount",
					Err:        "error amount less than minimum allowed",
					ErrMessage: fmt.Sprintf("Amount is less than minimum %s allowed", wdn.WithdrawMin),
				}
				break
			}

			//check amount if valid
			if (decimal.NewFromFloat(wdlInput.AmountSubmitted)).GreaterThan(decimal.RequireFromString(wdn.WithdrawMax)) {
				err = &tErrors.CustomError{
					Param:      "amount",
					Err:        "error amount greater than maximum allowed",
					ErrMessage: fmt.Sprintf("Amount is greater than maximum %s allowed", wdn.WithdrawMax),
				}
				break
			}
			//exit loop
			break
		}
	}
	if err != nil {
		log.Println("[QueueWithdrawalRequest] error validating request:", err)

		return err
	}

	if !validNetwork {
		log.Println("[QueueWithdrawalRequest] error invalid network")
		err = &tErrors.CustomError{
			Param:      "network",
			Err:        "error invalid network",
			ErrMessage: "Invalid network",
		}
		return err
	}

	//save the request
	log.Printf("[QueueWithdrawalRequest]>>:%+v\n", wdlInput) //prepare xdr
	xdrBase64, err := generateWithdrawalXdr(wallet, wdlInput, gc)
	if err != nil {

		log.Printf("[QueueWithdrawalRequest] %v withdrawal for %v generateWithdrawalXdr error:[%v] \n", wdlInput.Currency, wallet.Alias, err)
		return err
	}
	log.Printf("[QueueWithdrawalRequest] %v withdrawal for %v  transaction:[%v]\n", wdlInput.Currency, wallet.Alias, xdrBase64)

	oldTxn := wdlInput.Transaction

	wdlInput.NetworkPassPhrase = network.GetBlockchainNetworkPassPhrase()

	wdlInput.Transaction = xdrBase64
	log.Printf("[QueueWithdrawalRequest]>>>:%+v\n", wdlInput)
	if len(wdlInput.TransactionSignature) == 0 && wdlInput.Commit == 0 {
		//no signature
		return nil

	}

	//there was a signature... let's submit

	if oldTxn != xdrBase64 && wdlInput.Commit == 0 {
		return &tErrors.CustomError{
			Param:      "transaction",
			Err:        "transaction mismatch",
			ErrMessage: "transaction mismatch, please try again",
			Code:       404,
		}
	}

	if (wdlInput.Commit == 0 && wallet.SharedAccessEnabled == 1 && wallet.NumberOfApprovalsNeeded == 0) || (wallet.SharedAccessEnabled == 0 && wdlInput.Commit == 0) {
		dbTx := gc.DB.Begin()
		defer dbTx.Rollback()

		//save the withdrawal request
		wdlRequest := userModels.WithdrawalRequest{
			ID:                   uuid.NewString(),
			WalletAddress:        wallet.ID,
			WalletAlias:          wallet.Alias,
			UserID:               wallet.UserID,
			Currency:             wdlInput.Currency,
			AmountSubmitted:      wdlInput.AmountSubmitted,
			AmountToWithdraw:     wdlInput.AmountToWithdraw,
			WithdrawalAddress:    wdlInput.WithdrawalAddress,
			WithdrawalMemo:       wdlInput.WithdrawalMemo,
			WithdrawalNetwork:    wdlInput.WithdrawalNetwork,
			WithdrawalServiceFee: wdlInput.WithdrawalServiceFee,
			WithdrawalNetworkFee: wdlInput.WithdrawalNetworkFee,
		}
		e := dbTx.Omit(clause.Associations).Create(&wdlRequest).Error
		if e != nil {
			log.Printf("[QueueWithdrawalRequest] error creating withdrawal request on db. error: %v\n", e)
			return &tErrors.ErrorTemporaryServerError{}
		}

		txnID, err := network.SubmitXdrWithSignature(gc.BantuExpansionClient, wallet.Signer, xdrBase64, wdlInput.TransactionSignature)

		if err != nil {
			log.Printf("[QueueWithdrawalRequest]error submitting txn: %v\n", err)
			return &tErrors.ErrorTemporaryServerError{}
		}

		wdlInput.TransactionID = txnID
		wdlRequest.TransactionID = wdlInput.TransactionID
		e = dbTx.Omit(clause.Associations).Save(&wdlRequest).Error
		if e != nil {
			log.Printf("[QueueWithdrawalRequest] error saving withdrawal request for transactionID %v on db. error: %v\n", txnID, e)

		}
		dbTx.Commit()
		return nil
	}

	if wdlInput.Multiparty == 1 {
		wdlInput.TransactionID = "PENDING_AUTH"
		log.Printf("[QueueWithdrawalRequest]shared access with approver permission enabled for %v \n", wallet.Alias)
		id := uuid.NewString()

		description := fmt.Sprintf("Withdraw %v (%v),\n Amount: %v,\n Withdrawal Address: %v,\n Service Fee: %v,\n Network Fee: %v %v", wdlInput.Currency, wdlInput.WithdrawalNetwork, wdlInput.AmountSubmitted, wdlInput.WithdrawalAddress, serviceFee.String()+"%", wdlInput.WithdrawalNetworkFee, wdlInput.Currency)
		transactionByte, _ := json.Marshal(*wdlInput)
		transactionStr := string(transactionByte)
		pendingAuth := userModels.PendingAuth{
			ID:                     id,
			Initiator:              signerUser.Username,
			InitiatorSignerAddress: signerUser.PrimarySigner,
			WalletAddress:          wallet.ID,
			TransactionType:        "CRYPTO WITHDRAWAL",
			Description:            description,
			ApprovalsNeeded:        wallet.NumberOfApprovalsNeeded,
			TransactionXdr:         xdrBase64,
			TransactionInfoStr:     &transactionStr,
		}
		//save and commit this to database
		e := gc.DB.Omit(clause.Associations).Create(&pendingAuth).Error
		if e != nil {
			log.Printf("[QueueWithdrawalRequest] Error saving payment txn [%+v] transaction on pending auth table: %s\n", pendingAuth, e.Error())
			err = &tErrors.ErrorTemporaryServerError{}
			return err
		}
		wdlInput.ReturnedDescription = description
		return nil

	}

	return &tErrors.ErrorTemporaryServerError{}
}

func SubmitWithdrawalRequest(wallet *userModels.UserWallet, wdlInput userModels.WithdrawalRequestInput, gc *sharedconfig.GlobalConfig) (wdlItem userModels.CryptoWithdrawal, err error) {

	type WDLResp struct {
		Message string `json:"message"`
		Data    struct {
			WithdrawalID string `json:"withdrawalId"`
			Status       string `json:"status"`
		} `json:"data"`
	}

	var wdlResp WDLResp

	//validate input
	wdlNetworks, _ := GetWithdrawalNetworks(wdlInput.Currency, gc)
	validNetwork := false
	var wdn userModels.WithdrawalNetwork
	for _, wdn = range wdlNetworks {
		if strings.EqualFold(wdn.Network, wdlInput.WithdrawalNetwork) {
			if validNetwork {
				break
			}
			validNetwork = true
			//check amount if valid
			if (decimal.NewFromFloat(wdlInput.AmountSubmitted)).LessThan(decimal.RequireFromString(wdn.WithdrawMin)) {
				err = &tErrors.CustomError{
					Param:      "amount",
					Err:        "error amount less than minimum allowed",
					ErrMessage: fmt.Sprintf("Amount is less than minimum %s allowed", wdn.WithdrawMin),
				}
				break
			}

			//check amount if valid
			if (decimal.NewFromFloat(wdlInput.AmountSubmitted)).GreaterThan(decimal.RequireFromString(wdn.WithdrawMax)) {
				err = &tErrors.CustomError{
					Param:      "amount",
					Err:        "error amount greater than maximum allowed",
					ErrMessage: fmt.Sprintf("Amount is greater than maximum %s allowed", wdn.WithdrawMax),
				}
				break
			}
			//exit loop
			break
		}
	}
	if !validNetwork {
		err = &tErrors.CustomError{
			Param:      "network",
			Err:        "error invalid network",
			ErrMessage: "Invalid network",
		}

	}

	if err != nil {
		log.Println("[SubmitWithdrawalRequest] error validating request:", err)

		return
	}

	//

	client := http.DefaultClient
	cryptoWdlInput := userModels.CryptoWithdrawalRequestInput{
		Currency:  wdlInput.Currency,
		Amount:    wdlInput.AmountToWithdraw,
		ToAddress: wdlInput.WithdrawalAddress,
		Network:   wdlInput.WithdrawalNetwork,
		Memo:      wdlInput.WithdrawalMemo,
	}
	url := fmt.Sprintf("%s/%s", os.Getenv("ONELIQUIDITY_BASE_URL"), "wallets/v1/withdrawal")
	jbody, err := json.Marshal(cryptoWdlInput)
	if err != nil {
		log.Println("[SubmitWithdrawalRequest] error sending request:", err)

		return
	}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jbody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("ONELIQUIDITY_TOKEN")))
	resp, err := client.Do(request)
	if err != nil {
		log.Println("[SubmitWithdrawalRequest] error sending request:", err)
		return
	}
	if resp.StatusCode != 200 && resp.StatusCode != 201 {

		type ErrorResponse struct {
			Message string `json:"message"`
		}
		var errorResponse ErrorResponse
		defer resp.Body.Close()
		//Decode the data
		if err = json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			log.Println("[SubmitWithdrawalRequest] error decoding response:", err)
			return
		}
		log.Printf("[SubmitWithdrawalRequest] error response with code: %v, status: %v,error %v", resp.StatusCode, resp.Status, errorResponse.Message)

		err = &tErrors.ErrorTemporaryServerError{}
		return
	}

	defer resp.Body.Close()
	//Decode the data
	if err = json.NewDecoder(resp.Body).Decode(&wdlResp); err != nil {
		log.Println("[SubmitWithdrawalRequest] error decoding response for wdl networks id:", err)
		return
	}

	wdlItem, e := GetAWithdrawalID(wdlResp.Data.WithdrawalID, gc)
	if e != nil {
		err = &tErrors.CustomError{
			Param:      "withdrwalID",
			Err:        "error withdrawal successful but unable to retrieve status at this time",
			ErrMessage: "Withdrawal is already successful, but the status could not be confirmed at this time. Please refresh withdrwal history after 5mins to confirm status.",
			Code:       http.StatusAccepted,
		}
		return
	}
	wdlItem.TrovoWalletAddress = wallet.ID
	wdlItem.Fees = wdlInput.WithdrawalServiceFee
	e = gc.DB.Omit(clause.Associations).Save(&wdlItem).Error
	if e != nil {
		err = &tErrors.CustomError{
			Param:      "withdrwalID",
			Err:        "error withdrawal successful but unable to save status at this time",
			ErrMessage: "Withdrawal is already successful, but the status could not be saved at this time. Please refresh withdrwal history after 5mins to confirm status.",
			Code:       http.StatusAccepted,
		}
		return
	}

	return wdlItem, nil

}

func GetAWithdrawalID(withdrawalID string, gc *sharedconfig.GlobalConfig) (wdlItem userModels.CryptoWithdrawal, err error) {

	type WDLResp struct {
		Message string                      `json:"message"`
		Data    userModels.CryptoWithdrawal `json:"data"`
	}

	var wdlResp WDLResp
	client := http.DefaultClient

	url := fmt.Sprintf("%s/%s?withdrawalId=%s", os.Getenv("ONELIQUIDITY_BASE_URL"), "wallets/v1/withdrawal", withdrawalID)

	request, err := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("ONELIQUIDITY_TOKEN")))
	resp, err := client.Do(request)
	if err != nil {
		log.Println("[GetAWithdrawalID] error sending request:", err)
		return
	}
	if resp.StatusCode != 200 {

		type ErrorResponse struct {
			Message string `json:"message"`
		}
		var errorResponse ErrorResponse
		defer resp.Body.Close()
		//Decode the data
		if err = json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			log.Println("[GetAWithdrawalID] error decoding response:", err)
			return
		}
		log.Printf("[GetAWithdrawalID] error response with code: %v, status: %v,error %v", resp.StatusCode, resp.Status, errorResponse.Message)

		log.Println("[GetAWithdrawalID] error response with code: ", resp.StatusCode, resp.Status)
		err = &tErrors.ErrorTemporaryServerError{}
		return
	}

	defer resp.Body.Close()
	//Decode the data
	if err = json.NewDecoder(resp.Body).Decode(&wdlResp); err != nil {
		log.Println("[GetAWithdrawalID] error decoding response for wdl networks id:", err)
		return
	}
	wdlItem = wdlResp.Data

	return wdlItem, nil

}

func GetADepositByID(depositID string, gc *sharedconfig.GlobalConfig) (depItem userModels.DepositResponseItem, err error) {

	type DEPResp struct {
		Message string                         `json:"message"`
		Data    userModels.DepositResponseItem `json:"data"`
	}

	var depResp DEPResp
	client := http.DefaultClient

	url := fmt.Sprintf("%s/%s?depositId=%s", os.Getenv("ONELIQUIDITY_BASE_URL"), "wallets/v1/deposit", depositID)

	request, err := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("ONELIQUIDITY_TOKEN")))
	resp, err := client.Do(request)
	if err != nil {
		log.Println("[GetADepositByID] error sending request:", err)
		return
	}
	if resp.StatusCode != 200 {

		type ErrorResponse struct {
			Message string `json:"message"`
		}
		var errorResponse ErrorResponse
		defer resp.Body.Close()
		//Decode the data
		if err = json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			log.Println("[GetADepositByID] error decoding response:", err)
			return
		}
		log.Printf("[GetADepositByID] error response with code: %v, status: %v,error %v", resp.StatusCode, resp.Status, errorResponse.Message)

		log.Println("[GetADepositByID] error response with code: ", resp.StatusCode, resp.Status)
		err = &tErrors.ErrorTemporaryServerError{}
		return
	}

	defer resp.Body.Close()
	//Decode the data
	if err = json.NewDecoder(resp.Body).Decode(&depResp); err != nil {
		log.Println("[GetADepositByID] error decoding response for deposit:", err)
		return
	}
	depItem = depResp.Data

	return depItem, nil

}

func GetAllDeposits(lek, limit string, gc *sharedconfig.GlobalConfig) (depItems []userModels.DepositResponseItem, err error) {
	depItems = make([]userModels.DepositResponseItem, 0)
	type DEPResp struct {
		Message string                           `json:"message"`
		Data    []userModels.DepositResponseItem `json:"data"`
	}

	var depResp DEPResp
	client := http.DefaultClient

	url := fmt.Sprintf("%s/%s?lek=%v&limit=%v", os.Getenv("ONELIQUIDITY_BASE_URL"), "wallets/v1/deposit", lek, limit)

	request, err := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("ONELIQUIDITY_TOKEN")))
	resp, err := client.Do(request)
	if err != nil {
		log.Println("[GetAllDeposits] error sending request:", err)
		return
	}
	if resp.StatusCode != 200 {

		type ErrorResponse struct {
			Message string `json:"message"`
		}
		var errorResponse ErrorResponse
		defer resp.Body.Close()
		//Decode the data
		if err = json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			log.Println("[GetAllDeposits] error decoding response:", err)
			return
		}
		log.Printf("[GetAllDeposits] error response with code: %v, status: %v,error %v", resp.StatusCode, resp.Status, errorResponse.Message)

		log.Println("[GetAllDeposits] error response with code: ", resp.StatusCode, resp.Status)
		err = &tErrors.ErrorTemporaryServerError{}
		return
	}

	defer resp.Body.Close()
	//Decode the data
	if err = json.NewDecoder(resp.Body).Decode(&depResp); err != nil {
		log.Println("[GetAllDeposits] error decoding response for deposits:", err)
		return
	}
	depItems = depResp.Data

	return depItems, nil

}

func GetWithdrawalNetworks(currency string, gc *sharedconfig.GlobalConfig) (wdlNetworks []userModels.WithdrawalNetwork, err error) {
	var wdlNetworksResp userModels.CryptoWithdrawalNetworksResponse
	client := http.DefaultClient
	//get 'https://sandbox-api.oneliquidity.technology/wallets/v1/withdrawal/networks?currency=BTC'
	// cacheKey := fmt.Sprintf("wallets/v1/withdrawal/networks?currency=%s", currency)
	url := fmt.Sprintf("%s/%s?currency=%s", os.Getenv("ONELIQUIDITY_BASE_URL"), "wallets/v1/withdrawal/networks", currency)
	cacheKey := url

	// url := "https://sandbox-api.oneliquidity.technology/wallets/v1/withdrawal/networks?currency=BTC"
	{
		ok, rawData := gc.RedisCache.GetCachedResultRaw(cacheKey)
		if ok {

			json.Unmarshal(rawData, &wdlNetworksResp)
			if len(wdlNetworksResp.Data) > 0 {
				// log.Println("[GetWithdrawalNetworks] served from cache:", cacheKey)
				wdlNetworks = wdlNetworksResp.Data
				for key, v := range wdlNetworks {
					v.Currency = currency
					wdlNetworks[key] = v
				}
				return
			}

		}
	}
	request, err := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("ONELIQUIDITY_TOKEN")))
	resp, err := client.Do(request)
	if err != nil {
		log.Println("[GetWithdrawalNetworks] error sending request:", err)
		return
	}
	if resp.StatusCode != 200 {
		type ErrorResponse struct {
			Message string `json:"message"`
		}
		var errorResponse ErrorResponse
		defer resp.Body.Close()
		//Decode the data
		if err = json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			log.Println("[GetWithdrawalNetworks] error decoding response:", err)
			return
		}
		log.Printf("[GetWithdrawalNetworks] error response with code: %v, status: %v,error %v", resp.StatusCode, resp.Status, errorResponse.Message)

		err = &tErrors.ErrorTemporaryServerError{}
		return
	}

	defer resp.Body.Close()
	//Decode the data
	if err = json.NewDecoder(resp.Body).Decode(&wdlNetworksResp); err != nil {
		log.Println("[GetWithdrawalNetworks] error decoding response for wdl networks id:", err)
		return
	}
	wdlNetworks = wdlNetworksResp.Data
	for key, v := range wdlNetworks {
		v.Currency = currency
		wdlNetworks[key] = v
	}
	if len(wdlNetworks) > 0 {
		gc.RedisCache.StoreResultToCacheRaw(cacheKey, wdlNetworksResp, 1000)
	}
	dbTX := gc.DB.Begin()
	defer dbTX.Rollback()
	eDel := dbTX.Where("currency = ?", currency).Delete(&userModels.WithdrawalNetwork{}).Error
	if eDel != nil {
		log.Printf("[GetWithdrawalNetworks] error deleting from wdlNetworks where currency %v, error: %v\n", currency, eDel)
	}
	e := dbTX.Omit(clause.Associations).Create(&wdlNetworks).Error
	if e != nil {
		log.Printf("[GetWithdrawalNetworks] error creating wdlNetworks for %v error:%v\n", currency, e)

	} else {
		dbTX.Commit()
	}

	return wdlNetworks, nil

}

func ComplianceStartNewVerification(firstName, lastName string) (verificationID string, err error) {

	complianceAccountID := os.Getenv("COMPLIANCE_ACCOUNT_ID")
	vr := userModels.VerificationRequest{
		AccountID: complianceAccountID,
		FirstName: firstName,
		LastName:  lastName,
	}
	jbody, err := json.Marshal(vr)
	if err != nil {
		return "", err
	}
	client := http.DefaultClient
	//get upload credentials
	url := fmt.Sprintf("%s/%s?", os.Getenv("ONELIQUIDITY_BASE_URL"), "compliance/v1/verification")
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jbody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("ONELIQUIDITY_TOKEN")))
	resp, err := client.Do(request)
	if err != nil {
		log.Println("[ComplianceStartNewVerification] error starting new verification:", err)
		return "", err
	}
	log.Println("[ComplianceStartNewVerification] succeeded with code: ", resp.StatusCode)
	defer resp.Body.Close()
	var veriResp userModels.VerificationResponse
	//Decode the data
	if err := json.NewDecoder(resp.Body).Decode(&veriResp); err != nil {
		log.Println("[ComplianceStartNewVerification] error decoding response for verification id:", err)
		return "", err
	}

	return veriResp.Data.VerificationID, nil

}

func GetProofOfresidencyCred(user *userModels.User) (fMCred userModels.ProofOfResidenceCred, verificationID string, err error) {
	if user.Corporate == 1 || user.LastName == nil {
		err = &tErrors.CustomError{
			Param:      "user",
			Err:        "error-not-available-for-corporate",
			ErrMessage: "Service not available for corporate accounts",
		}
		return
	}
	verificationID, err = ComplianceStartNewVerification(user.FirstName, *user.LastName)
	if err != nil {
		return fMCred, "", err
	}

	url := fmt.Sprintf("%s/compliance/v1/kyc/verification/signed-url?verificationId=%s&uploadDocType=proof_of_residency", os.Getenv("ONELIQUIDITY_BASE_URL"), verificationID)

	client := http.DefaultClient
	//get upload credentials

	request, err := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("ONELIQUIDITY_TOKEN")))
	resp, err := client.Do(request)
	if err != nil {
		log.Println("[GetProofOfresidencyCred] error starting new verification:", err)
		return fMCred, verificationID, err
	}
	log.Println("[GetProofOfresidencyCred] succeeded with code: ", resp.StatusCode)
	defer resp.Body.Close()
	// var fmResp userModels.VerificationResponse
	//Decode the data
	if err := json.NewDecoder(resp.Body).Decode(&fMCred); err != nil {
		log.Println("[GetProofOfresidencyCred] error decoding proof of residency response:", err)
		return fMCred, verificationID, err
	}

	return fMCred, verificationID, nil

}
func GetFacematchPassportCred(user *userModels.User) (fMCred userModels.FaceMatchVerificationCred, verificationID string, err error) {
	if user.Corporate == 1 || user.LastName == nil {
		err = &tErrors.CustomError{
			Param:      "user",
			Err:        "error-not-available-for-corporate",
			ErrMessage: "Service not available for corporate accounts",
		}
		return
	}
	verificationID, err = ComplianceStartNewVerification(user.FirstName, *user.LastName)
	if err != nil {
		return fMCred, "", err
	}

	url := fmt.Sprintf("%s/compliance/v1/kyc/verification/signed-url?verificationId=%s&uploadDocType=facematch&sourceDocType=passport", os.Getenv("ONELIQUIDITY_BASE_URL"), verificationID)

	client := http.DefaultClient
	//get upload credentials

	request, err := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("ONELIQUIDITY_TOKEN")))
	resp, err := client.Do(request)
	if err != nil {
		log.Println("[GetFacematchPassportCred] error starting new verification:", err)
		return fMCred, verificationID, err
	}
	log.Println("[GetFacematchPassportCred] succeeded with code: ", resp.StatusCode)
	defer resp.Body.Close()
	// var fmResp userModels.VerificationResponse
	//Decode the data
	if err := json.NewDecoder(resp.Body).Decode(&fMCred); err != nil {
		log.Println("[GetFacematchPassportCred] error decoding facematch response:", err)
		return fMCred, verificationID, err
	}

	return fMCred, verificationID, nil

}

func GetFacematchNationalIDCred(user *userModels.User) (fMCred userModels.FaceMatchVerificationCred, verificationID string, err error) {
	if user.Corporate == 1 || user.LastName == nil {
		err = &tErrors.CustomError{
			Param:      "user",
			Err:        "error-not-available-for-corporate",
			ErrMessage: "Service not available for corporate accounts",
		}
		return
	}
	verificationID, err = ComplianceStartNewVerification(user.FirstName, *user.LastName)
	if err != nil {
		return fMCred, "", err
	}

	url := fmt.Sprintf("%s/compliance/v1/kyc/verification/signed-url?verificationId=%s&uploadDocType=facematch&sourceDocType=national_id", os.Getenv("ONELIQUIDITY_BASE_URL"), verificationID)

	client := http.DefaultClient
	//get upload credentials

	request, err := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("ONELIQUIDITY_TOKEN")))
	resp, err := client.Do(request)
	if err != nil {
		log.Println("[GetFacematchPassportCred] error starting new verification:", err)
		return fMCred, verificationID, err
	}
	log.Println("[GetFacematchPassportCred] succeeded with code: ", resp.StatusCode)
	defer resp.Body.Close()
	// var fmResp userModels.VerificationResponse
	//Decode the data
	if err := json.NewDecoder(resp.Body).Decode(&fMCred); err != nil {
		log.Println("[GetFacematchPassportCred] error decoding facematch response:", err)
		return fMCred, verificationID, err
	}

	return fMCred, verificationID, nil

}

func GetFacematchDrivingLicenseCred(user *userModels.User) (fMCred userModels.FaceMatchVerificationCred, verificationID string, err error) {
	if user.Corporate == 1 || user.LastName == nil {
		err = &tErrors.CustomError{
			Param:      "user",
			Err:        "error-not-available-for-corporate",
			ErrMessage: "Service not available for corporate accounts",
		}
		return
	}
	verificationID, err = ComplianceStartNewVerification(user.FirstName, *user.LastName)
	if err != nil {
		return fMCred, "", err
	}

	url := fmt.Sprintf("%s/compliance/v1/kyc/verification/signed-url?verificationId=%s&uploadDocType=facematch&sourceDocType=driving_license", os.Getenv("ONELIQUIDITY_BASE_URL"), verificationID)

	client := http.DefaultClient
	//get upload credentials

	request, err := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("ONELIQUIDITY_TOKEN")))
	resp, err := client.Do(request)
	if err != nil {
		log.Println("[GetFacematchPassportCred] error starting new verification:", err)
		return fMCred, verificationID, err
	}
	log.Println("[GetFacematchPassportCred] succeeded with code: ", resp.StatusCode)
	defer resp.Body.Close()
	// var fmResp userModels.VerificationResponse
	//Decode the data
	if err := json.NewDecoder(resp.Body).Decode(&fMCred); err != nil {
		log.Println("[GetFacematchPassportCred] error decoding facematch response:", err)
		return fMCred, verificationID, err
	}

	return fMCred, verificationID, nil

}

func StartFacematchForPassport(user *userModels.User, selfieVideo, documentPicture *multipart.FileHeader, gc *sharedconfig.GlobalConfig) (err error) {
	// , gc *sharedconfig.GlobalConfig
	kycdData, err := user.GetKycData(gc)

	if err != nil {
		return errors.New("user country not set")
	}
	fmCred, verificationID, err := GetFacematchPassportCred(user)

	if err != nil {
		return err
	}

	sf := strings.Split(selfieVideo.Filename, ".")
	selfieExtension := fmt.Sprintf(".%s", sf[len(sf)-1])

	pf := strings.Split(documentPicture.Filename, ".")
	documentExtension := fmt.Sprintf(".%s", pf[len(pf)-1])

	// update credentials to db
	doc := userModels.FacematchPassport{
		ID:                uuid.NewString(),
		UserID:            user.ID,
		VerificationID:    verificationID,
		SelfieVideoFormat: mime.TypeByExtension(selfieExtension),
		PictureFormat:     mime.TypeByExtension(documentExtension),
		SelfieKey:         fmCred.Data.SelfieVideo.Fields.Key,
		DocumentKey:       fmCred.Data.Passport.Fields.Key,
	}
	dbTX := gc.DB.Begin()
	defer dbTX.Rollback()
	e := dbTX.Omit(clause.Associations).Create(&doc).Error
	if e != nil {
		return &tErrors.ErrorTemporaryServerError{}
	}
	//uploadVideo
	err = UploadFileToS3(selfieVideo, &fmCred.Data.SelfieVideo)
	if err != nil {
		return err
	}
	//uploadPassport
	err = UploadFileToS3(documentPicture, &fmCred.Data.Passport)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/compliance/v1/kyc/facematch", os.Getenv("ONELIQUIDITY_BASE_URL"))

	client := http.DefaultClient
	// get upload credentials
	fmrequest := userModels.StartDocumentRequest{
		VerificationID: verificationID,
	}
	fmrequest.Inputs = append(fmrequest.Inputs, userModels.DocumentInputs{
		Country:  kycdData.Country,
		DocType:  "selfie_video",
		MimeType: mime.TypeByExtension(selfieExtension),
		Key:      fmCred.Data.SelfieVideo.Fields.Key,
	})
	fmrequest.Inputs = append(fmrequest.Inputs, userModels.DocumentInputs{
		Country:  kycdData.Country,
		DocType:  "passport",
		MimeType: mime.TypeByExtension(documentExtension),
		Key:      fmCred.Data.Passport.Fields.Key,
		Side:     "front",
	})
	jbody, err := json.Marshal(fmrequest)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jbody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("ONELIQUIDITY_TOKEN")))
	resp, err := client.Do(request)
	if err != nil {
		log.Println("[StartFacematchForPassport] error starting new verification:", err)
		return err
	}
	log.Println("[StartFacematchForPassport] succeeded with code: ", resp.StatusCode)
	defer resp.Body.Close()
	var fmResp userModels.OKResponse
	// Decode the data
	if err := json.NewDecoder(resp.Body).Decode(&fmResp); err != nil {
		log.Println("[StartFacematchForPassport] error decoding start facematch for passport response:", err)
		return err
	}

	dbTX.Commit()
	return nil

}

func StartFacematchForDrivingLicense(user *userModels.User, selfieVideo, documentPicture *multipart.FileHeader, gc *sharedconfig.GlobalConfig) (err error) {
	// , gc *sharedconfig.GlobalConfig
	kycdData, err := user.GetKycData(gc)

	if err != nil {
		return errors.New("user country not set")
	}

	fmCred, verificationID, err := GetFacematchDrivingLicenseCred(user)

	if err != nil {
		return err
	}

	sf := strings.Split(selfieVideo.Filename, ".")
	selfieExtension := fmt.Sprintf(".%s", sf[len(sf)-1])

	pf := strings.Split(documentPicture.Filename, ".")
	documentExtension := fmt.Sprintf(".%s", pf[len(pf)-1])

	// update credentials to db
	doc := userModels.FacematchDrivingLicense{
		ID:                uuid.NewString(),
		UserID:            user.ID,
		VerificationID:    verificationID,
		SelfieVideoFormat: mime.TypeByExtension(selfieExtension),
		PictureFormat:     mime.TypeByExtension(documentExtension),
		SelfieKey:         fmCred.Data.SelfieVideo.Fields.Key,
		DocumentKey:       fmCred.Data.Passport.Fields.Key,
	}
	dbTX := gc.DB.Begin()
	defer dbTX.Rollback()
	e := dbTX.Omit(clause.Associations).Create(&doc).Error
	if e != nil {
		return &tErrors.ErrorTemporaryServerError{}
	}
	//uploadVideo
	err = UploadFileToS3(selfieVideo, &fmCred.Data.SelfieVideo)
	if err != nil {
		return err
	}
	//uploadPassport
	err = UploadFileToS3(documentPicture, &fmCred.Data.Passport)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/compliance/v1/kyc/facematch", os.Getenv("ONELIQUIDITY_BASE_URL"))

	client := http.DefaultClient
	// get upload credentials
	fmrequest := userModels.StartDocumentRequest{
		VerificationID: verificationID,
	}
	fmrequest.Inputs = append(fmrequest.Inputs, userModels.DocumentInputs{
		Country:  kycdData.Country,
		DocType:  "selfie_video",
		MimeType: mime.TypeByExtension(selfieExtension),
		Key:      fmCred.Data.SelfieVideo.Fields.Key,
	})
	fmrequest.Inputs = append(fmrequest.Inputs, userModels.DocumentInputs{
		Country:  kycdData.Country,
		DocType:  "driving_license",
		MimeType: mime.TypeByExtension(documentExtension),
		Key:      fmCred.Data.Passport.Fields.Key,
		Side:     "front",
	})
	jbody, err := json.Marshal(fmrequest)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jbody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("ONELIQUIDITY_TOKEN")))
	resp, err := client.Do(request)
	if err != nil {
		log.Println("[StartFacematchForDrivingLicense] error starting new verification:", err)
		return err
	}
	log.Println("[StartFacematchForDrivingLicense] succeeded with code: ", resp.StatusCode)
	defer resp.Body.Close()
	var fmResp userModels.OKResponse
	// Decode the data
	if err := json.NewDecoder(resp.Body).Decode(&fmResp); err != nil {
		log.Println("[StartFacematchForDrivingLicense] error decoding start facematch for passport response:", err)
		return err
	}
	dbTX.Commit()
	return nil

}

func StartFacematchForNationalID(user *userModels.User, selfieVideo, documentPicture *multipart.FileHeader, gc *sharedconfig.GlobalConfig) (err error) {
	// , gc *sharedconfig.GlobalConfig
	kycdData, err := user.GetKycData(gc)

	if err != nil {
		return errors.New("user country not set")
	}
	fmCred, verificationID, err := GetFacematchNationalIDCred(user)

	if err != nil {
		return err
	}

	sf := strings.Split(selfieVideo.Filename, ".")
	selfieExtension := fmt.Sprintf(".%s", sf[len(sf)-1])

	pf := strings.Split(documentPicture.Filename, ".")
	documentExtension := fmt.Sprintf(".%s", pf[len(pf)-1])

	// update credentials to db
	doc := userModels.FacematchNationalID{
		ID:                uuid.NewString(),
		UserID:            user.ID,
		VerificationID:    verificationID,
		SelfieVideoFormat: mime.TypeByExtension(selfieExtension),
		PictureFormat:     mime.TypeByExtension(documentExtension),
		SelfieKey:         fmCred.Data.SelfieVideo.Fields.Key,
		DocumentKey:       fmCred.Data.Passport.Fields.Key,
	}
	dbTX := gc.DB.Begin()
	defer dbTX.Rollback()
	e := dbTX.Omit(clause.Associations).Create(&doc).Error
	if e != nil {
		return &tErrors.ErrorTemporaryServerError{}
	}
	//uploadVideo
	err = UploadFileToS3(selfieVideo, &fmCred.Data.SelfieVideo)
	if err != nil {
		return err
	}
	//uploadPassport
	err = UploadFileToS3(documentPicture, &fmCred.Data.Passport)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/compliance/v1/kyc/facematch", os.Getenv("ONELIQUIDITY_BASE_URL"))

	client := http.DefaultClient
	// get upload credentials
	fmrequest := userModels.StartDocumentRequest{
		VerificationID: verificationID,
	}
	fmrequest.Inputs = append(fmrequest.Inputs, userModels.DocumentInputs{
		Country:  kycdData.Country,
		DocType:  "selfie_video",
		MimeType: mime.TypeByExtension(selfieExtension),
		Key:      fmCred.Data.SelfieVideo.Fields.Key,
	})
	fmrequest.Inputs = append(fmrequest.Inputs, userModels.DocumentInputs{
		Country:  kycdData.Country,
		DocType:  "national_id",
		MimeType: mime.TypeByExtension(documentExtension),
		Key:      fmCred.Data.Passport.Fields.Key,
		Side:     "front",
	})
	jbody, err := json.Marshal(fmrequest)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jbody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("ONELIQUIDITY_TOKEN")))
	resp, err := client.Do(request)
	if err != nil {
		log.Println("[StartFacematchForNationalID] error starting new verification:", err)
		return err
	}
	log.Println("[StartFacematchForNationalID] succeeded with code: ", resp.StatusCode)
	defer resp.Body.Close()
	var fmResp userModels.OKResponse
	// Decode the data
	if err := json.NewDecoder(resp.Body).Decode(&fmResp); err != nil {
		log.Println("[StartFacematchForNationalID] error decoding start facematch for passport response:", err)
		return err
	}
	dbTX.Commit()
	return nil

}

func StartGovernmentIDCheckForProofOfResidency(user *userModels.User, documentPicture *multipart.FileHeader, gc *sharedconfig.GlobalConfig) (err error) {
	// , gc *sharedconfig.GlobalConfig
	kycdData, err := user.GetKycData(gc)

	if err != nil {
		return errors.New("user country not set")
	}
	fmCred, verificationID, err := GetProofOfresidencyCred(user)

	if err != nil {
		return err
	}

	pf := strings.Split(documentPicture.Filename, ".")
	pictureExtension := fmt.Sprintf(".%s", pf[len(pf)-1])

	// update credentials to db
	doc := userModels.GovermentIDProofOfResidency{
		ID:             uuid.NewString(),
		UserID:         user.ID,
		VerificationID: verificationID,
		PictureFormat:  mime.TypeByExtension(pictureExtension),
		DocumentKey:    fmCred.Data.PresignedURL.ProofOfResidency.Fields.Key,
	}
	dbTX := gc.DB.Begin()
	defer dbTX.Rollback()
	e := dbTX.Omit(clause.Associations).Create(&doc).Error
	if e != nil {
		return &tErrors.ErrorTemporaryServerError{}
	}

	//uploadPassport
	err = UploadFileToS3(documentPicture, &fmCred.Data.PresignedURL.ProofOfResidency)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/compliance/v1/kyc/gov-id", os.Getenv("ONELIQUIDITY_BASE_URL"))

	client := http.DefaultClient
	// get upload credentials
	fmrequest := userModels.StartDocumentRequest{
		VerificationID: verificationID,
	}

	fmrequest.Inputs = append(fmrequest.Inputs, userModels.DocumentInputs{
		Country:  kycdData.Country,
		DocType:  "proof_of_residency",
		MimeType: mime.TypeByExtension(pictureExtension),
		Key:      fmCred.Data.PresignedURL.ProofOfResidency.Fields.Key,
		Side:     "front",
	})
	jbody, err := json.Marshal(fmrequest)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jbody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("ONELIQUIDITY_TOKEN")))
	resp, err := client.Do(request)
	if err != nil {
		log.Println("[StartGovernmentIDCheckForProofOfResidency] error starting new verification:", err)
		return err
	}
	log.Println("[StartGovernmentIDCheckForProofOfResidency] succeeded with code: ", resp.StatusCode)
	defer resp.Body.Close()
	var fmResp userModels.OKResponse
	// Decode the data
	if err := json.NewDecoder(resp.Body).Decode(&fmResp); err != nil {
		log.Println("[StartGovernmentIDCheckForProofOfResidency] error decoding start goverment id check response:", err)
		return err
	}
	dbTX.Commit()
	return nil

}

func generateWithdrawalXdr(wallet *userModels.UserWallet, wdlInput *userModels.WithdrawalRequestInput, gc *sharedconfig.GlobalConfig) (base64Xdr string, err error) {

	ca, err := userModels.Currency(wdlInput.Currency).GetCurratedAsset(gc)
	if err != nil {
		log.Printf("[generateWithdrawalXdr] error getting %v property, error: %v\n", wdlInput.Currency, err)
		return "", err
	}

	var asset basetxn.Asset = nil

	// asset = basetxn.NativeAsset{}

	asset = basetxn.CreditAsset{Code: ca.AssetCode, Issuer: ca.ContractAddress}

	_, sourceAccountTrustsAsset, nativeAccountBalance, currencyBalance, sourceAccount, errorSource := network.BlockchainAccountProperties(gc.BantuExpansionClient, wallet.ID, asset)

	if errorSource != nil {
		log.Printf("[generateWithdrawalXdr] error withdrawing %v , error: %v\n", wdlInput.Currency, errorSource)
		return "", errorSource
	}

	if nativeAccountBalance.Equal(decimal.Zero) {
		return "", &tErrors.ErrorUnderfundedAccount{}
	}
	if currencyBalance.LessThan(decimal.NewFromFloat(wdlInput.AmountSubmitted)) {
		return "", &tErrors.CustomError{
			Param:      "submittedAmount",
			Err:        "error insufficient balance for " + wdlInput.Currency,
			ErrMessage: "Wallet does not have enough balance to withraw" + wdlInput.Currency,
		}
	}

	if !sourceAccountTrustsAsset {
		return "", &tErrors.CustomError{
			Param:      "wallet",
			Err:        "error account does not have currency",
			ErrMessage: "Wallet does not have " + wdlInput.Currency,
		}
	}
	chanAccount := <-gc.ChannelAccounts
	defer func(c *evmkeypair.Full) {
		gc.ChannelAccounts <- c
	}(chanAccount)

	_, _, _, _, chanSourceAccount, errorChannel := network.BlockchainAccountProperties(gc.BantuExpansionClient, chanAccount.Address(), basetxn.NativeAsset{})
	if errorChannel != nil {
		log.Printf("[generateWithdrawalXdr] error withdrawing %v , channel account error: %v\n", wdlInput.Currency, errorChannel)
		return "", errorChannel
	}
	var ops []basetxn.Operation = make([]basetxn.Operation, 0)

	ops = append(ops, &basetxn.Payment{
		Destination:   ca.ContractAddress,
		Amount:        fmt.Sprintf("%v", wdlInput.AmountSubmitted),
		Asset:         asset,
		SourceAccount: wallet.ID,
	})

	var tx *basetxn.Transaction
	// Construct the transaction that holds the operations to execute on the network
	if wdlInput.Multiparty == 1 {
		tx, err = basetxn.NewTransaction(
			basetxn.TransactionParams{
				SourceAccount:        chanSourceAccount.Address,
				IncrementSequenceNum: true,
				Operations:           ops,
				BaseFee:              2000,
				Memo:                 "withdraw-" + wdlInput.Currency,
			},
		)
	} else {
		tx, err = basetxn.NewTransaction(
			basetxn.TransactionParams{
				SourceAccount:        sourceAccount.Address,
				IncrementSequenceNum: true,
				Operations:           ops,
				BaseFee:              2000,
				Memo:                 "withdraw-" + wdlInput.Currency,
			},
		)
	}

	if err != nil {
		log.Println("[generateWithdrawalXdr]error constructing transaction ", err)
		return "", &tErrors.ErrorTemporaryServerError{}
	}

	if wdlInput.Multiparty == 1 {

		tx, err = tx.Sign(network.GetBlockchainNetworkPassPhrase(), chanAccount)

		if err != nil {
			log.Println("[generateWithdrawalXdr] error signing transaction with channelAccount key ", err)
			return "", &tErrors.ErrorTemporaryServerError{}
		}
	}

	base64Xdr, err = tx.Base64()

	if err != nil {
		log.Printf("[generateWithdrawalXdr] error withdrawing %v , extracting base 64 xdr error: %v\n", wdlInput.Currency, err)

		return "", err
	}
	if len(base64Xdr) == 0 {
		log.Printf("[generateWithdrawalXdr] error withdrawing %v , transaction is empty\n", wdlInput.Currency)

		return "", &tErrors.ErrorTemporaryServerError{}
	}

	log.Printf("[generateWithdrawalXdr] transaction for withdrawing %v, xdrbase64:[%v]\n", wdlInput.Currency, base64Xdr)

	return base64Xdr, nil

}

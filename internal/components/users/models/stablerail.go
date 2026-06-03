package users

import "time"

type StablerailConfig struct {
	CreatedAt        time.Time
	UpdatedAt        time.Time
	ID               uint64
	ApiKey           string
	FintechID        string
	BaseUrl          string
	EnableStablerail int `gorm:"default:0" json:"enableStablerail"`
}

type StablerailUser struct {
	CreatedAt     time.Time
	UpdatedAt     time.Time
	ID            string
	TrovoUsername string
}
type StablerailOnboardUserRetry struct {
	CreatedAt     time.Time
	UpdatedAt     time.Time
	ID            uint64
	TrovoUsername string
	BVN           string // xx...xxx

}

type StablerailRequest struct {
	CreatedAt                time.Time
	UpdatedAt                time.Time
	ID                       string
	RequestType              string
	Status                   string
	StablerailResponseObject *string
	TrovoUsername            string
}

type StablerailOnramp struct {
	CreatedAt     time.Time
	UpdatedAt     time.Time
	ID            string
	WalletAddress string
	// BaseAmount      float64
	// Fee             float64
	TotalAmount     float64
	TargetAsset     string //USDT
	Status          string
	AutoSwapEnabled int
	TrovoUsername   string
}
type StablerailOfframp struct {
	CreatedAt     time.Time
	UpdatedAt     time.Time
	ID            string
	BankCode      string
	AccountNumber string
	BaseAmount    float64
	Fee           float64
	ActualAmount  float64
	Ticker        string //CNGN
	Status        string
	TrovoUsername string
}

// Request payload
type StablerailOnboardRequest struct {
	BVN string `json:"bvn"`
}

// Request payload
type StablerailOnboardStatus struct {
	RequestId string `json:"requestId"`
}

// Response structure
type StablerailOnboardResponse struct {
	Status       string `json:"status"`
	ResponseCode string `json:"response_code"`
	Message      string `json:"message"`
	Data         struct {
		RequestID               string `json:"requestId"`
		Status                  string `json:"status"`
		UserHash                string `json:"userHash"`
		Message                 string `json:"message"`
		EstimatedCompletionTime string `json:"estimatedCompletionTime"`
	} `json:"data"`
}

// Response structs
type StablerailCheckOnboardingResponse struct {
	Status       string `json:"status"`
	ResponseCode string `json:"response_code"`
	Message      string `json:"message"`
	Data         struct {
		Type      string    `json:"type"`
		UserID    string    `json:"userId"`
		Status    string    `json:"status"`
		CreatedAt time.Time `json:"createdAt"`
		Message   string    `json:"message"`
	} `json:"data"`
}

// Structs to match the API response
type StablerailBank struct {
	BankCode    string `gorm:"primaryKey" json:"bank_code"`
	BankName    string `json:"bank_name"`
	CountryCode string `json:""`
}

type StablerailGetBanksResponse struct {
	Status       string `json:"status"`
	ResponseCode string `json:"response_code"`
	Data         struct {
		CountryCode string `json:"countryCode"`
		Banks       []struct {
			BankCode string `gorm:"primaryKey" json:"bank_code"`
			BankName string `json:"bank_name"`
		} `json:"banks"`
	} `json:"data"`
}

// Request payload struct
type StableRailOfframpRequest struct {
	UserID        string `json:"userId"`
	Amount        int    `json:"amount"`
	AccountNumber string `json:"accountNumber"`
	BankCode      string `json:"bankCode"`
	Ticker        string `json:"ticker"`
}

// Request payload struct
type StableRailWithdrawalRequest struct {
	UserID            string  `json:"userId"`
	InternalWallet    string  `json:"internalWallet"`
	DestinationWallet string  `json:"destinationWallet"`
	Amount            float64 `json:"amount"`
	Ticker            string  `json:"ticker"`
	Network           string  `json:"network"`
}

// Response struct
type StableRailOfframpResponse struct {
	Status       string `json:"status"`
	ResponseCode string `json:"response_code"`
	Message      string `json:"message"`
	Data         struct {
		RequestID string `json:"requestId"`
		Status    string `json:"status"`
		Stage     string `json:"stage"`
	} `json:"data"`
}

// Request payload
type StablerailOfframpStatusRequest struct {
	RequestID string `json:"requestId"`
}

// Response structures
type StablerailOfframpStatusAPIResponse struct {
	Status       string                      `json:"status"`
	ResponseCode string                      `json:"response_code"`
	Message      string                      `json:"message"`
	Data         StablerailOfframpStatusData `json:"data"`
}

type StablerailOfframpStatusData struct {
	RequestID         string                  `json:"requestId"`
	Status            string                  `json:"status"`
	Stage             string                  `json:"stage"`
	CreatedAt         time.Time               `json:"createdAt"`
	UpdatedAt         time.Time               `json:"updatedAt"`
	UserWalletAddress string                  `json:"userWalletAddress"`
	VaultAddress      string                  `json:"vaultAddress"`
	TokenAddress      string                  `json:"tokenAddress"`
	Amount            string                  `json:"amount"`
	FintechUserID     string                  `json:"fintechUserId"`
	TokenTransfer     StablerailTokenTransfer `json:"tokenTransfer"`
	Payout            StablerailPayout        `json:"payout"`
}

type StablerailTokenTransfer struct {
	Status          string    `json:"status"`
	TransactionHash string    `json:"transactionHash"`
	BlockNumber     string    `json:"blockNumber"`
	Timestamp       time.Time `json:"timestamp"`
}

type StablerailPayout struct {
	Status        string    `json:"status"`
	Reference     string    `json:"reference"`
	TransactionID string    `json:"transactionId"`
	Timestamp     time.Time `json:"timestamp"`
	Attempts      int       `json:"attempts"`
}

// Request payload
type CNGNOnrampRequest struct {
	Owner          string  `json:"owner"`     //Owner wallet address if you want to manage the generated transaction wallet. Caveat!!! Doing this, our system won't be able to withdraw funds from such wallet.
	Amount         float64 `json:"amount"`    //Amount to fund in Naira (e.g., 500 = ₦500)
	AssetSwap      string  `json:"assetSwap"` //Asset to swap into after funding (e.g. USDC)
	AutoSwap       bool    `json:"autoSwap"`  //Automatically swap funded amount to assetSwap
	UserID         string  `json:"userId"`
	SweepToOfframp bool    `json:"sweepToOfframp"` //Automatically route funds to the user's default wallet after funding
}

// Response structures
type CNGNOnrampResponse struct {
	Status       string           `json:"status"`
	ResponseCode string           `json:"response_code"`
	Message      string           `json:"message"`
	Data         CNGNResponseData `json:"data"`
}

type CNGNResponseData struct {
	RequestID             string           `json:"requestId"`
	WalletAddress         string           `json:"walletAddress"`
	Status                string           `json:"status"`
	Version               string           `json:"version"`
	Message               string           `json:"message"`
	AutoSwapEnabled       bool             `json:"autoSwapEnabled"`
	SweepToOfframpEnabled bool             `json:"sweepToOfframpEnabled"`
	TargetAsset           string           `json:"targetAsset"`
	FeeBreakdown          CNGNFeeBreakdown `json:"feeBreakdown"`
}

type CNGNFeeBreakdown struct {
	BaseAmount     float64   `json:"baseAmount"`
	FintechFee     float64   `json:"fintechFee"`
	GatewayFee     float64   `json:"gatewayFee"`
	StablesRailFee float64   `json:"stablesRailFee"`
	TotalFee       float64   `json:"totalFee"`
	TotalAmount    float64   `json:"totalAmount"`
	Breakdown      Breakdown `json:"breakdown"`
}

type Breakdown struct {
	UserRequestedAmount      float64 `json:"userRequestedAmount"`
	FintechFeeAmount         float64 `json:"fintechFeeAmount"`
	FintechFeePercentage     float64 `json:"fintechFeePercentage"`
	FintechFeeCapped         bool    `json:"fintechFeeCapped"`
	GatewayFeeAmount         float64 `json:"gatewayFeeAmount"`
	GatewayFeePercentage     float64 `json:"gatewayFeePercentage"`
	StablesRailFeeAmount     float64 `json:"stablesRailFeeAmount"`
	StablesRailFeePercentage float64 `json:"stablesRailFeePercentage"`
	StablesRailFeeCapped     bool    `json:"stablesRailFeeCapped"`
	TotalFeeAmount           float64 `json:"totalFeeAmount"`
	FinalAmount              float64 `json:"finalAmount"`
	AmountToWallet           float64 `json:"amountToWallet"`
}

// Request payload
type GetVirtualAccountRequest struct {
	RequestID string `json:"requestId"`
}

// Response structs
type VirtualAccountFeeBreakdown struct {
	UserRequestedAmount      float64 `json:"userRequestedAmount"`
	FintechFeeAmount         float64 `json:"fintechFeeAmount"`
	FintechFeePercentage     float64 `json:"fintechFeePercentage"`
	FintechFeeCapped         bool    `json:"fintechFeeCapped"`
	GatewayFeeAmount         float64 `json:"gatewayFeeAmount"`
	GatewayFeePercentage     float64 `json:"gatewayFeePercentage"`
	StablesRailFeeAmount     float64 `json:"stablesRailFeeAmount"`
	StablesRailFeePercentage float64 `json:"stablesRailFeePercentage"`
	StablesRailFeeCapped     bool    `json:"stablesRailFeeCapped"`
	TotalFeeAmount           float64 `json:"totalFeeAmount"`
	FinalAmount              float64 `json:"finalAmount"`
	AmountToWallet           float64 `json:"amountToWallet"`
}

type VirtualAccount struct {
	AccountNumber      string                     `json:"accountNumber"`
	BankName           string                     `json:"bankName"`
	AccountName        string                     `json:"accountName"`
	Amount             float64                    `json:"amount"`
	CreatedAt          time.Time                  `json:"createdAt"`
	BaseAmount         float64                    `json:"baseAmount"`
	FeeAmount          float64                    `json:"feeAmount"`
	TotalAmountWithFee float64                    `json:"totalAmountWithFee"`
	FeePercentage      float64                    `json:"feePercentage"`
	FeeBreakdown       VirtualAccountFeeBreakdown `json:"feeBreakdown"`
}

type VirtualAccountData struct {
	RequestID      string         `json:"requestId"`
	VirtualAccount VirtualAccount `json:"virtualAccount"`
	Status         string         `json:"status"`
	WalletAddress  string         `json:"walletAddress"`
	Version        string         `json:"version"`
}

type GetVirtualAccountResponse struct {
	Status       string             `json:"status"`
	ResponseCode string             `json:"response_code"`
	Message      string             `json:"message"`
	Data         VirtualAccountData `json:"data"`
}

// Request payload
type CNGNOnrampStatusRequest struct {
	RequestID string `json:"requestId"`
}

// Response structures
type CNGNOnrampStatusAPIResponse struct {
	Status       string               `json:"status"`
	ResponseCode string               `json:"response_code"`
	Message      string               `json:"message"`
	Data         CNGNOnrampStatusData `json:"data"`
}

type CNGNOnrampStatusData struct {
	Status    string                   `json:"status"`
	RequestID string                   `json:"requestId"`
	Version   string                   `json:"version"`
	Metadata  CNGNOnrampStatusMetadata `json:"metadata"`
	Wallet    CNGNOnrampStatusWallet   `json:"wallet"`
	FXQuote   CNGNOnrmpStatusFXQuote   `json:"fxQuote"`
}

type CNGNOnrampStatusMetadata struct {
	TokenBuy       string `json:"tokenBuy"`
	AutoSwap       bool   `json:"autoSwap"`
	SweepToOfframp bool   `json:"sweepToOfframp"`
}

type CNGNOnrampStatusWallet struct {
	WalletAddress        string                                `json:"walletAddress"`
	Owner                string                                `json:"owner"`
	TokenBuy             string                                `json:"tokenBuy"`
	AutoSwap             bool                                  `json:"autoSwap"`
	SweepToOfframp       bool                                  `json:"sweepToOfframp"`
	Amount               string                                `json:"amount"`
	CreatedAt            time.Time                             `json:"createdAt"`
	FundingRequestedAt   time.Time                             `json:"fundingRequestedAt"`
	FundedAt             time.Time                             `json:"fundedAt"`
	TransactionHash      string                                `json:"transactionHash"`
	VirtualAccount       CNGNOnrampStatusVirtualAccountDetails `json:"virtualAccountDetails"`
	VirtualAccountStatus string                                `json:"virtualAccountStatus"`
}

type CNGNOnrampStatusVirtualAccountDetails struct {
	VAID          string    `json:"vaId"`
	AccountNumber string    `json:"accountNumber"`
	BankName      string    `json:"bankName"`
	AccountName   string    `json:"accountName"`
	Amount        float64   `json:"amount"`
	CreatedAt     time.Time `json:"createdAt"`
}

type CNGNOnrmpStatusFXQuote struct {
	Available     bool      `json:"available"`
	Pair          string    `json:"pair"`
	CNGNAmount    string    `json:"cngnAmount"`
	TokenAmount   string    `json:"tokenAmount"`
	AveragePrice  string    `json:"averagePrice"`
	AverageSpread float64   `json:"averageSpread"`
	ExpiresAt     time.Time `json:"expiresAt"`
	Note          string    `json:"note"`
}

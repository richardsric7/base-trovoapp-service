package users

import "time"

type StablerailConfig struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	ID        uint64
	ApiKey    string
	FintechID string
	BaseUrl   string
}

type StablerailUser struct {
	CreatedAt     time.Time
	UpdatedAt     time.Time
	ID            string
	TrovoUsername string
	// BVN           string // xx...xxx

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
	CreatedAt       time.Time
	UpdatedAt       time.Time
	ID              string
	WalletAddress   string
	BaseAmount      float64
	Fee             float64
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
type StablerailCheckOnboardingResponse  struct {
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

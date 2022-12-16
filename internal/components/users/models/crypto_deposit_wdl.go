package users

type PaginatedCryptoDepositHistory struct {
	Pages        int             `json:"pages"`
	CurrentPage  int             `json:"currentPage"`
	TotalRecords int             `json:"totalRecords"`
	Limit        int             `json:"limit"`
	Records      []CryptoDeposit `json:"records"`
}
type PaginatedCryptoWithdrawalHistory struct {
	Pages        int                `json:"pages"`
	CurrentPage  int                `json:"currentPage"`
	TotalRecords int                `json:"totalRecords"`
	Limit        int                `json:"limit"`
	Records      []CryptoWithdrawal `json:"records"`
}

type CryptoDeposit struct {
	TrovoWalletPublicKey string `gorm:"size:100" json:"TrovoWalletPublicKey"`
	DepositID            string `gorm:"primaryKey" json:"depositId"`
	TxID                 string `gorm:"index:unique_txid,unique" json:"txId"`
	Amount               string `json:"amount"`
	CreatedAt            string `json:"createdAt"`
	UpdatedAt            string `json:"updatedAt"`
	Currency             string `json:"currency"`
	Decimal              int    `json:"decimal"`
	Fees                 string `json:"fees"`
	FromAddress          string `json:"fromAddress"`
	IsCompleted          bool   `json:"isCompleted"`
	IsValid              bool   `json:"isValid"`
	IsVerified           bool   `json:"isVerified"`
	ToAddress            string `gorm:"index:unique_txid,unique" json:"toAddress"`
}

type WithdrawalNetwork struct {
	Network              string `gorm:"primaryKey" json:"network"`
	Name                 string `gorm:"size:100" json:"name"`
	AddressRegex         string `json:"addressRegex"`
	MemoRegex            string `gorm:"size:100" json:"memoRegex"`
	WithdrawFee          string `gorm:"size:100" json:"withdrawFee"`
	WithdrawMin          string `gorm:"size:100" json:"withdrawMin"`
	WithdrawMax          string `gorm:"size:100" json:"withdrawMax"`
	EstimatedArrivalTime int    `json:"estimatedArrivalTime"`
}

type CryptoWithdrawalRequestInput struct {
	Currency  string  `json:"currency"`
	Amount    float64 `json:"amount"`
	ToAddress string  `json:"toAddress"`
	Network   string  `json:"network"`
	Memo      string  `json:"memo"`
}
type WithdrawalRequestInput struct {
	Currency         string  `json:"currency"`
	AmountSubmitted  float64 `json:"amountSubmitted"`
	AmountToWithdraw float64 `json:"amountToWithdraw"`
	ToAddress        string  `json:"toAddress"`
	Network          string  `json:"network"`
	Memo             string  `json:"memo"`
	Fees             float64 `json:"fees"`
}

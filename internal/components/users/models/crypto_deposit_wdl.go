package users

type PaginatedCryptoDepositHistory struct {
	Pages        int             `json:"pages"`
	CurrentPage  int             `json:"currentPage"`
	TotalRecords int             `json:"totalRecords"`
	Limit        int             `json:"limit"`
	Records      []CryptoDeposit `json:"records"`
}

type CryptoDepositResponse struct {
	Message string `json:"message"`
	Data    struct {
		DepositID   string `json:"depositId"`
		TxID        string `json:"txId"`
		Amount      string `json:"amount"`
		CreatedAt   string `json:"createdAt"`
		UpdatedAt   string `json:"updatedAt"`
		Currency    string `json:"currency"`
		Decimal     int    `json:"decimal"`
		Fees        string `json:"fees"`
		FromAddress string `json:"fromAddress"`
		IsCompleted bool   `json:"isCompleted"`
		IsValid     bool   `json:"isValid"`
		IsVerified  bool   `json:"isVerified"`
		ToAddress   string `json:"toAddress"`
	} `json:"data"`
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

type CryptoWithdrawalResponse struct {
	Message string `json:"message"`
	Data    []struct {
		Network              string `json:"network"`
		Name                 string `json:"name"`
		AddressRegex         string `json:"addressRegex"`
		MemoRegex            string `json:"memoRegex"`
		WithdrawFee          string `json:"withdrawFee"`
		WithdrawMin          string `json:"withdrawMin"`
		WithdrawMax          string `json:"withdrawMax"`
		EstimatedArrivalTime int    `json:"estimatedArrivalTime"`
	} `json:"data"`
}

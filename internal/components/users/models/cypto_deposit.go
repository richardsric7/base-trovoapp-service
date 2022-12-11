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
	TrovoWalletPublicKey string `json:"TrovoWalletPublicKey"`
	DepositID            string `gorm:"primaryKey" json:"depositId"`
	TxID                 string `json:"txId"`
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
	ToAddress            string `json:"toAddress"`
}

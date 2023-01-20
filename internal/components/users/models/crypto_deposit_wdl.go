package users

import (
	"time"
	"trovo-wallet-api/internal/sharedconfig"
)

type PaginatedCryptoDepositHistory struct {
	Pages        int                 `json:"pages"`
	CurrentPage  int                 `json:"currentPage"`
	TotalRecords int                 `json:"totalRecords"`
	Limit        int                 `json:"limit"`
	Records      []CryptoDepositJSON `json:"records"`
}
type PaginatedCryptoWithdrawalHistory struct {
	Pages        int                 `json:"pages"`
	CurrentPage  int                 `json:"currentPage"`
	TotalRecords int                 `json:"totalRecords"`
	Limit        int                 `json:"limit"`
	Records      []WithdrawalRequest `json:"records"`
}

type CryptoDeposit struct {
	ID                   uint64
	TrovoWalletPublicKey string `gorm:"size:100" json:"trovoWalletPublicKey"`
	DepositID            string `gorm:"index:unique_depositid,unique" json:"depositId"`
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

func (c *CryptoDeposit) ToJSON(gc *sharedconfig.GlobalConfig) (jsonObj CryptoDepositJSON) {
	jsonObj = CryptoDepositJSON{
		TrovoWalletPublicKey: c.TrovoWalletPublicKey,
		TxID:                 c.TxID,
		Amount:               c.Amount,
		CreatedAt:            c.CreatedAt,
		Currency:             c.Currency,
		FromAddress:          c.FromAddress,
		IsCompleted:          c.IsCompleted,
		IsValid:              c.IsValid,
		IsVerified:           c.IsVerified,
		ToAddress:            c.ToAddress,
	}

	return

}

type CryptoDepositJSON struct {
	TrovoWalletPublicKey string `json:"trovoWalletPublicKey"`
	TxID                 string `json:"txId"`
	Amount               string `json:"amount"`
	CreatedAt            string `json:"createdAt"`
	UpdatedAt            string `json:"updatedAt"`
	Currency             string `json:"currency"`
	FromAddress          string `json:"fromAddress"`
	IsCompleted          bool   `json:"isCompleted"`
	IsValid              bool   `json:"isValid"`
	IsVerified           bool   `json:"isVerified"`
	ToAddress            string `json:"toAddress"`
}

type WithdrawalNetwork struct {
	Currency             string `gorm:"primaryKey" json:"currency"`
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
	Currency             string  `json:"currency"`
	AmountSubmitted      float64 `json:"amountSubmitted"`
	AmountToWithdraw     float64 `json:"amountToWithdraw"` //submitted amount less serviceFee
	WithdrawalAddress    string  `json:"withdrawalAddress"`
	WithdrawalNetwork    string  `json:"withdrawalNetwork"`
	WithdrawalMemo       string  `json:"withdrawalMemo"`
	WithdrawalServiceFee float64 `json:"withdrawalServiceFee"`
	WithdrawalNetworkFee float64 `json:"withdrawalNetworkFee"`
	Transaction          string  `json:"transaction"`
	TransactionSignature string  `json:"transactionSignature"`
	TransactionID        string  `json:"transactionId"`
	NetworkPassPhrase    string  `json:"networkPassPhrase"`
	Multiparty           int     `json:"-"`
	TransactionSource    string  `json:"-"`
	SignatureRequired    int     `json:"signatureRequired"`
	Commit               int     `json:"commit"`
	ReturnedDescription  string  `json:"-"`
}

type WithdrawalRequest struct {
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"-"`
	ID                   string    `gorm:"primaryKey" json:"id"`
	WalletPublicKey      string    `gorm:"size:100" json:"walletPublicKey"`
	WalletAlias          string    `gorm:"size:100" json:"walletAlias"`
	UserID               string    `gorm:"size:100" json:"userId"`
	Currency             string    `json:"currency"`
	AmountSubmitted      float64   `json:"amountSubmitted"`
	AmountToWithdraw     float64   `json:"amountToWithdraw"` //submitted amount less serviceFee
	WithdrawalAddress    string    `json:"withdrawalAddress"`
	WithdrawalNetwork    string    `json:"withdrawalNetwork"`
	WithdrawalMemo       string    `json:"withdrawalMemo"`
	WithdrawalServiceFee float64   `json:"withdrawalServiceFee"`
	WithdrawalNetworkFee float64   `json:"withdrawalNetworkFee"`
	TransactionID        string    `gorm:"size:100" json:"transactionId"`
	WithdrawalStatus     string    `gorm:"size:100;default:'PENDING'" json:"withdrawalStatus"`
}

type CallbackDeposit struct {
	Event   string              `json:"event"`
	Message string              `json:"message"`
	Data    CallbackDepositItem `json:"data"`
}

type CallbackDepositItem struct {
	CreatedAt       time.Time
	DepositID       string `gorm:"primaryKey" json:"depositId"`
	Currency        string `json:"currency"`
	Network         string `json:"network"`
	Txid            string `gorm:"size:100;index:idx_callbacltxid,unique" json:"txid"`
	Amount          string `json:"amount"`
	Fees            string `json:"fees"`
	FromAddress     string `json:"from_address"`
	ToAddress       string `json:"to_address"`
	Decimal         int    `json:"decimal"`
	ProcessingState int    `json:"processing_state"`
	Minted          int    `gorm:"default:0" json:"minted"`
}

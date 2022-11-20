package users

import "time"

type MarketOfferRequest struct {
	OfferType            string   `json:"offerType"`
	AssetCode            string   `json:"assetCode"`
	AssetIssuer          string   `json:"assetIssuer"`
	CurrencyCode         string   `json:"currencyCode"`
	CurrencyIssuer       string   `json:"currencyIssuer"`
	PricePerUnit         string   `json:"pricePerUnit"`
	Quantity             string   `json:"quantity"`
	FeeChargedOnAsset    string   `gorm:"not null;size:100" json:"-"`
	FeeValue             string   `gorm:"not null;size:100" json:"-"`
	NetQuantity          string   `gorm:"not null;size:100" json:"-"`
	Transaction          string   `json:"transaction"`
	TransactionSignature string   `json:"transactionSignature"`
	TransactionID        string   `json:"transactionId"`
	NetworkPassPhrase    string   `json:"networkPassPhrase"`
	Messages             []string `json:"messages"`
	Commit               int      `json:"commit"`
	Multiparty           int      `json:"-"`
	SignatureRequired    int      `json:"signatureRequired"`
	Memo                 string   `json:"memo"`
	ReturnedDescription  string   `json:"-"`
}

type DeleteOfferRequest struct {
	ID                   string   `json:"Id"`
	Transaction          string   `json:"transaction"`
	TransactionSignature string   `json:"transactionSignature"`
	TransactionID        string   `json:"transactionId"`
	NetworkPassPhrase    string   `json:"networkPassPhrase"`
	Messages             []string `json:"messages"`
	Commit               int      `json:"commit"`
	Multiparty           int      `json:"-"`
	Memo                 string   `json:"memo"`
	ReturnedDescription  string   `json:"-"`
}
type MarketOffer struct {
	CreatedAt                   time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt                   time.Time `gorm:"default:now()" json:"-"`
	ID                          string    `gorm:"size:100" json:"id"`
	SourceWalletAlias           string    `gorm:"not null;size:100" json:"sourceWalletAlias"`
	SourcewalletPublicKey       string    `gorm:"vsize:100" json:"sourceWalletPublicKey"`
	MarketMakingWalletPublicKey string    `gorm:"not null;size:100" json:"marketMakingWalletPublicKey"`
	OfferType                   string    `gorm:"not null;size:100" json:"offerType"`
	AssetCode                   string    `gorm:"not null;size:100" json:"assetCode"`
	AssetIssuer                 *string   `gorm:"null;size:100" json:"assetIssuer"`
	CurrencyCode                string    `gorm:"not null;size:100" json:"currencyCode"`
	CurrencyIssuer              *string   `gorm:"null;size:100" json:"currencyIssuer"`
	PricePerUnit                string    `gorm:"not null;size:100" json:"pricePerUnit"`
	Quantity                    string    `gorm:"not null;size:100" json:"quantity"`
	FeeChargedOnAsset           string    `gorm:"not null;size:100" json:"feeChargedOnAsset"`
	FeeValue                    string    `gorm:"not null;size:100" json:"FeeValue"`
	NetQuantity                 string    `gorm:"not null;size:100" json:"netQuantity"`
	TransactionID               *string   `gorm:"not null;size:100" json:"transactionId"`
	BlockchainOfferID           *string   `gorm:"null;size:100" json:"blockchainOfferId"`
	RemainingQuantity           string    `gorm:"not null;size:100" json:"remainingQuantity"`
	RemainingFeeValue           string    `gorm:"not null;size:100" json:"remainingFeeValue"`
	Canceled                    int       `gorm:"type:integer;not null; default:0" json:"canceled"`
}

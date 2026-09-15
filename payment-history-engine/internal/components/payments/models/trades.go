package payments

import "time"

type MarketOffer struct {
	CreatedAt                   time.Time `gorm:"default:now()" json:"-"`
	UpdatedAt                   time.Time `gorm:"default:now()" json:"-"`
	ID                          string    `gorm:"size:100" json:"id"`
	SourceWalletAlias           string    `gorm:"not null;size:100" json:"sourceWalletAlias"`
	SourceWalletPublicKey       string    `gorm:"vsize:100" json:"sourceWalletPublicKey"`
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
	TransactionID               *string   `gorm:"null;size:100" json:"transactionId"`
	BlockchainOfferID           *string   `gorm:"null;size:100" json:"blockchainOfferId"`
	LastProcessedCursor         *string   `gorm:"null;size:100" json:"lastProcessedCursor"`
	RemainingQuantity           string    `gorm:"not null;size:100" json:"remainingQuantity"`
	RemainingFeeValue           string    `gorm:"not null;size:100" json:"remainingFeeValue"`
	Canceled                    int       `gorm:"type:integer;not null; default:0" json:"canceled"`
}

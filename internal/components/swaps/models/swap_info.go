package swaps

import (
	"time"
)

// PaymentInfo holds payment information
type PaymentInfo struct {
	Destination          string   `json:"destination"`
	Memo                 string   `json:"memo"`
	AssetIssuer          string   `json:"assetIssuer"`
	AssetCode            string   `json:"assetCode"`
	Amount               string   `json:"amount"`
	Transaction          string   `json:"transaction"`
	TransactionSignature string   `json:"transactionSignature"`
	TransactionID        string   `json:"transactionId"`
	NetworkPassPhrase    string   `json:"networkPassPhrase"`
	DestinationFirstName string   `json:"destinationFirstName"`
	DestinationLastName  string   `json:"destinationLastName"`
	DestinationThumbnail string   `json:"destinationThumbnail"`
	Messages             []string `json:"messages"`
}

// PaymentLog holds payment information for logging
type PaymentLog struct {
	CreatedAt              time.Time
	TrasanctionType        uint64   `gorm:"size:56;default:0;not null"` //0 = Payment(default), 1 = Swap
	Sender                 string   `gorm:"size:56;not null"`
	Destination            string   `gorm:"size:56;not null"` //sender for type = 0, owner for type = 1
	Memo                   *string  `gorm:"size:60;null"`
	AssetIssuer            *string  `gorm:"size:56;null"`      //source asset
	AssetCode              *string  `gorm:"size:12;null"`      //source asset
	DestinationAssetIssuer *string  `gorm:"size:56;null"`      //for swap
	DestinationAssetCode   *string  `gorm:"size:12;null"`      //for swap
	Amount                 string   `gorm:"size:100;not null"` //source amount
	Transaction            *string  `gorm:"null"`
	TransactionSignature   *string  `gorm:"null"`
	TransactionID          string   `gorm:"size:200;not null"`
	DestinationFirstName   *string  `gorm:"size:50;null"`
	DestinationLastName    *string  `gorm:"size:50;null"`
	PublicIP               string   `gorm:"size:50;not null"`
	CountryCode            *string  `gorm:"size:2;null"`
	Latitude               *float64 `gorm:"null"`
	Longitude              *float64 `gorm:"null"`
	City                   *string  `gorm:"null;size:100"`
	Region                 *string  `gorm:"null;size:100"`
	RegionName             *string  `gorm:"null;size:100"`
	TimeZone               *string  `gorm:"null;size:100"`
	ISP                    *string  `gorm:"null;size:150"`
}

// SwapSendPathInput model struct for user input when importing keys for new user
type SwapSendPathInput struct {
	DestinationAccount string
	DestinationAssets  string
	SourceAssetCode    string
	SourceAssetIssuer  string
	SourceAmount       string
}

// SwapSendPathInput model struct for user input when importing keys for new user
type SwapSendInfo struct {
	DestinationAssetCode   string   `json:"destinationAssetCode"`
	DestinationAssetIssuer string   `json:"destinationAssetIssuer"`
	SwappedEstimate        string   `json:"swappedEstimate"`
	SourceAssetCode        string   `json:"sourceAssetCode"`
	SourceAssetIssuer      string   `json:"sourceAssetIssuer"`
	SourceAmount           string   `json:"sourceAmount" `
	Transaction            string   `json:"transaction"`
	TransactionSignature   string   `json:"transactionSignature"`
	TransactionID          string   `json:"transactionId"`
	NetworkPassPhrase      string   `json:"networkPassPhrase"`
	Messages               []string `json:"messages"`
	Memo                   string   `json:"-"`
	Multiparty             int      `json:"-"`
	SignatureRequired      int      `json:"signatureRequired"`
	Commit                 int      `json:"commit"`
	SHash                  string   `json:"sHash"`
	ReturnedDescription    string   `json:"-"`
}

// SwapReceivePathInput model struct for user input when importing keys for new user
type SwapReceivePathInput struct {
	DestinationAccount     string   `json:"destinationAccount"`
	DestinationAssetType   string   `json:"destinationAssetType"`
	DestinationAssetCode   string   `json:"destinationAssetCode"`
	DestinationAssetIssuer string   `json:"destinationAssetIssuer"`
	DestinationAmount      string   `json:"destinationAmount"`
	SourceAccount          string   `json:"sourceAccount"`
	SourceAssets           string   `json:"sourceAssets"`
	Transaction            string   `json:"transaction"`
	TransactionSignature   string   `json:"transactionSignature"`
	TransactionID          string   `json:"transactionId"`
	NetworkPassPhrase      string   `json:"networkPassPhrase"`
	Messages               []string `json:"messages"`
}

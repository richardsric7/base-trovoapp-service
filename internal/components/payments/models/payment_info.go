package payments

import (
	"time"
)

// PaymentInfo holds payment information
type PaymentInfo struct {
	Destination             string            `json:"destination"`
	Memo                    string            `json:"memo"`
	AssetIssuer             string            `json:"assetIssuer"`
	AssetCode               string            `json:"assetCode"`
	Amount                  string            `json:"amount"`
	Transaction             string            `json:"transaction"`
	TransactionSignature    string            `json:"transactionSignature"`
	TransactionID           string            `json:"transactionId"`
	NetworkPassPhrase       string            `json:"networkPassPhrase"`
	DestinationFirstName    string            `json:"destinationFirstName"`
	DestinationLastName     string            `json:"destinationLastName"`
	DestinationThumbnail    string            `json:"destinationThumbnail"`
	DestinationVerified     int               `json:"destinationVerified"`
	Multiparty              int               `json:"-"`
	TransactionSource       string            `json:"-"`
	SignatureRequired       int               `json:"signatureRequired"`
	Commit                  int               `json:"commit"`
	Vat                     string            `json:"vat"`
	Fee                     string            `json:"fee"`
	FeeAmount               string            `json:"feeAmount"`
	VatAmount               string            `json:"vatAmount"`
	AmountToPay             string            `json:"amountToPay"`
	SHash                   string            `json:"sHash"`
	ChannelAccount          string            `json:"channelAccount"`
	ChannelAccountSignature string            `json:"channelAccountSignature"`
	Messages                []string          `json:"messages"`
	CallbackURLS            map[string]string `json:"-"`
}

// PaymentLog holds payment information for logging
type PaymentLog struct {
	CreatedAt            time.Time
	Sender               string   `gorm:"size:100;not null"`
	Destination          string   `gorm:"size:100;not null"`
	ChannelAccount       *string  `gorm:"size:56;null"`
	Memo                 *string  `gorm:"size:60;null"`
	AssetIssuer          *string  `gorm:"size:56;null"`
	AssetCode            *string  `gorm:"size:12;null"`
	Amount               string   `gorm:"size:100;not null"`
	Transaction          *string  `gorm:"null"`
	TransactionSignature *string  `gorm:"null"`
	TransactionID        string   `gorm:"size:200;not null"`
	DestinationFirstName *string  `gorm:"size:50;null"`
	DestinationLastName  *string  `gorm:"size:50;null"`
	PublicIP             string   `gorm:"size:50;not null"`
	CountryCode          *string  `gorm:"size:2;null"`
	Latitude             *float64 `gorm:"null"`
	Longitude            *float64 `gorm:"null"`
	City                 *string  `gorm:"null;size:100"`
	Region               *string  `gorm:"null;size:100"`
	RegionName           *string  `gorm:"null;size:100"`
	TimeZone             *string  `gorm:"null;size:100"`
	ISP                  *string  `gorm:"null;size:150"`
}

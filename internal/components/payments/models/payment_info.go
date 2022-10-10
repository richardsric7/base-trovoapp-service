package payments

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
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
	ChannelAccount          string            `json:"channelAccount"`
	ChannelAccountSignature string            `json:"channelAccountSignature"`
	Messages                []string          `json:"messages"`
	CallbackURLS            map[string]string `json:"-"`
}

// PaymentLog holds payment information for logging
type PaymentLog struct {
	CreatedAt            time.Time
	Sender               string   `gorm:"size:56;not null"`
	Destination          string   `gorm:"size:56;not null"`
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

// PaymentChanObject holds data for faucet payments
type PaymentChanObject struct {
	Secret         string
	Destination    string `json:"destination"`
	Memo           string `json:"memo"`
	AssetIssuer    string `json:"assetIssuer"`
	AssetCode      string `json:"assetCode"`
	Amount         string `json:"amount"`
	SourceUsername string `json:"sourceUsername"`
}

// DBPaymentChanObject holds data for faucet payments from SQLite Database
type DBPaymentObject struct {
	ID             uint64 `gorm:"primaryKey" json:"-"`
	Secret         string `gorm:"size:56;not null"`
	Destination    string `gorm:"size:56;not null" json:"destination"`
	Memo           string `gorm:"size:28;not null" json:"memo"`
	AssetIssuer    string `gorm:"size:56;null;default:''" json:"assetIssuer"`
	AssetCode      string `gorm:"size:12;null;default:''" json:"assetCode"`
	Amount         string `gorm:"size:50;not null" json:"amount"`
	SourceUsername string `gorm:"size:16;not null" json:"sourceUsername"`
	Paid           int    `gorm:"not null;default:0"`
}

func (po *DBPaymentObject) SaveData(sqlite *gorm.DB) error {
	if po.ID == 0 {
		if po.Secret == "" || po.Destination == "" || po.Memo == "" || po.Amount == "" || po.SourceUsername == "" {
			return errors.New("a required parameter is empty, model: " + fmt.Sprintf("%+v", *po))
		}
		return sqlite.Create(&po).Error
	}
	return sqlite.Save(&po).Error

}

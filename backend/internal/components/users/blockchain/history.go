package users

import (
	"time"
)

//PaymentItem holds blockchain payment item
type PaymentItem struct {
	DestinationUsername       string    `json:"destinationUsername"`
	DestinationFirstName      string    `json:"destinationFirstName"`
	DestinationLastName       string    `json:"destinationLastName"`
	DestinationImageThumbnail string    `json:"-"`
	DestinationVerified       int       `json:"destinationVerified"`
	FunderUsername            string    `json:"senderUsername"`
	FunderFirstName           string    `json:"senderFirstName"`
	FunderLastName            string    `json:"senderLastName"`
	FunderImageThumbnail      string    `json:"-"`
	FunderVerified            int       `json:"funderVerified"`
	Amount                    string    `json:"amount"`
	AssetCode                 string    `json:"assetCode"`
	AssetIssuer               string    `json:"assetIssuer"`
	SourceAssetCode           string    `json:"sourceAssetCode"`        //for path payment
	SourceAssetIssuer         string    `json:"sourceAssetIssuer"`      //for path payment
	DestinationAssetCode      string    `json:"destinationAssetCode"`   //for path payment
	DestinationAssetIssuer    string    `json:"destinationAssetIssuer"` //for path payment
	TransactionID             string    `json:"transactionID"`
	FeeCharged                int64     `json:"-"`
	Ledger                    int32     `json:"-"`
	TransactionMemo           string    `json:"transactionMemo"`
	TransactionTime           time.Time `json:"transactionTime"`
	Cursor                    string    `json:"cursor"` //paging token returned as cursor to resume streaming from this state
}

//PaymentEntities holds entities in the payment history batch
type PaymentEntities struct {
	Username       string `json:"username"`
	ImageThumbnail string `json:"thumbnail"`
}

//PaymentHistory holds blockchain payment history
type PaymentHistory struct {
	PageCursor string                     `json:"pageCursor"`
	Payments   []PaymentItem              `json:"payments"`
	Entities   map[string]PaymentEntities `json:"entities"`
}

//IndexedBantuOperation for holding the index of the operation for sorting later
type IndexedBantuOperation struct {
	Index     int
	Operation PaymentItem
}

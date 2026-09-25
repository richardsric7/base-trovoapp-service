package models

// PendingAssetToClaim holds pensing assets to be claimed
type PendingAssetToClaim struct {
	AssetCode            string `json:"assetCode"`
	AssetIssuer          string `json:"assetIssuer"`
	Transaction          string `json:"transaction"`
	TransactionSignature string `json:"transactionSignature"`
	TransactionID        string `json:"transactionId"`
	NetworkPassPhrase    string `json:"networkPassPhrase"`
}

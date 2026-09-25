package models

// PendingAssetToClaim holds pensing assets to be claimed
type PendingAssetToClaim struct {
	AssetCode            string `json:"assetCode"`
	ContractAddress      string `json:"contractAddress"`
	Transaction          string `json:"transaction"`
	TransactionSignature string `json:"transactionSignature"`
	TransactionID        string `json:"transactionId"`
	NetworkPassPhrase    string `json:"networkPassPhrase"`
}

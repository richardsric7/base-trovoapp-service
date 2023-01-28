package users

// Streams holds model for Stream data object
type Streams struct {
	Stream    string `json:"stream"`
	Cursor    string `json:"cursor"`
	Signature string `json:"signature"`
	Signer    string `json:"signer"`
}

type OrderBookStream struct {
	AssetCode      string `json:"assetCode"`
	AssetIssuer    string `json:"assetIssuer"`
	CurrencyCode   string `json:"currencyCode"`
	CurrencyIssuer string `json:"currencyIssuer"`
}

// Handshake holds model for handshake data object
type Handshake struct {
	Auth    bool   `json:"auth"`
	Message string `json:"message"`
}

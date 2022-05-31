package users

import "github.com/shopspring/decimal"

//Balance model for user
type Balance struct {
	AssetIssuer string          `json:"assetIssuer"`
	AssetCode   string          `json:"assetCode"`
	Amount      decimal.Decimal `json:"amount"`
	QRCode      string          `json:"qrCode"`
}

//Signer model for user
type Signer struct {
	Weight  int    `json:"weight"`
	Key     string `json:"key"`
	Type    string `json:"type"`
	Sponsor string `json:"sponsor"`
}

//Signer model for user
type Thresholds struct {
	LowThreshold    string `json:"low_threshold"`
	MediumThreshold string `json:"medium_threshold"`
	HighThreshold   string `json:"high_threshold"`
}

//AssetBalances holds user balances
type AssetBalances struct {
	Claimed   map[string]Balance `json:"claimed"`
	Unclaimed map[string]Balance `json:"unclaimed"`
}

//NFTBalances holds user NFT balances
type NFTBalances struct {
	NFTs []NFT `json:"nfts"`
}
type NFT struct {
	AssetIssuer string `json:"assetIssuer"`
	AssetCode   string `json:"assetCode"`
}

//UserBalanceForMerchant holds user balances
type UserBalanceForMerchant struct {
	Balances []Balance `json:"balances"`
}

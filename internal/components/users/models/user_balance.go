package users

import "github.com/shopspring/decimal"

//Balance model for user
type Balance struct {
	AssetIssuer string          `json:"assetIssuer"`
	AssetCode   string          `json:"assetCode"`
	Amount      decimal.Decimal `json:"amount"`
	QRCode      string          `json:"qrCode"`
	ImageURL    string          `json:"imageUrl"`
}

type NFT struct {
	AssetIssuer    string `json:"assetIssuer"`
	AssetCode      string `json:"assetCode"`
	NFTName        string `json:"nftName"`
	NFTDescription string `json:"nftDescription"`
	NFTImageURI    string `json:"nftImageURI"`
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
	Claimed   []Balance `json:"claimed"`
	Unclaimed []Balance `json:"unclaimed"`
}

//MappedNFTBalance holds user NFTs per mapped public key
type MappedNFTBalance map[string][]NFT

//UserBalanceForMerchant holds user balances
type UserBalanceForMerchant struct {
	Balances []Balance `json:"balances"`
}

type MappedBalance map[string]Balance

package merchants

//Balance model for user
type Balance struct {
	AssetIssuer string `json:"assetIssuer"`
	AssetCode   string `json:"assetCode"`
	Amount      string `json:"amount"`
	QRCode      string `json:"qrCode"`
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

// UsdPrice     string                    `json:"usdPrice"`
// UsdValue     string                    `json:"usdValue"`
// CuratedAsset bool                      `json:"curatedAsset"`
// AssetInfo    assetsmodels.CuratedAsset `json:"assetInfo"`
// //TempBalance model for user
// type TempBalance struct {
// 	AssetIssuer  string                    `json:"assetIssuer"`
// 	AssetCode    string                    `json:"assetCode"`
// 	Amount       string                    `json:"amount"`
// 	UsdPrice     string                    `json:"usdPrice"`
// 	UsdValue     string                    `json:"usdValue"`
// 	CuratedAsset bool                      `json:"curatedAsset"`
// 	AssetInfo    assetsmodels.CuratedAsset `json:"assetInfo"`
// }

//Balances holds user balances
type Balances struct {
	Claimed   []Balance `json:"claimed"`
	Unclaimed []Balance `json:"unclaimed"`
}

//UserBalanceForMerchant holds user balances
type UserBalanceForMerchant struct {
	Balances []Balance `json:"balances"`
}

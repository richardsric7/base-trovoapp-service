package users

type WalletsSharedWithUser struct {
	Owner             string        `json:"owner"`
	WalletPublicKey   string        `json:"walletPublicKey"`
	Permission        string        `json:"permission"`
	WalletAlias       string        `json:"walletAlias"`
	WalletDescription string        `json:"walletDescription"`
	AssetBalances     AssetBalances `json:"assetBalances"`
}

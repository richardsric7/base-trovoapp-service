package users

type ThirdPartyWalletAccess struct {
	Owner             string `json:"owner"`
	PublicKey         string `json:"publicKey"`
	AccessLevel       string `json:"accessLevel"`
	WalletAlias       string `json:"walletAlias"`
	WalletDescription string `json:"walletDescription"`
}

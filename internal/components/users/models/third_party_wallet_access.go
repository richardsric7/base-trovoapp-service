package users

type ThirdPartyWalletAccess struct {
	Owner             string `json:"owner"`
	WalletPublicKey         string `json:"walletPublicKey"`
	AccessLevel       string `json:"accessLevel"`
	WalletAlias       string `json:"walletAlias"`
	WalletDescription string `json:"walletDescription"`
}

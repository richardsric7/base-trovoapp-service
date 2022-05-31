package users

type UserInfo struct {
	UserData               UserJSON                 `json:"userData"`
	AssetBalances          map[string]AssetBalances `json:"assetBalances"` //map of wallet public key and the asset balances
	NFTBalances            map[string]NFTBalances   `json:"nftBalances"`   //map of nft wallet and nftBalances
	ThirdPartyWalletAccess []ThirdPartyWalletAccess `json:"thirdPartyWalletAccess"`
}

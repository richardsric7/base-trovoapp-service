package users

type UserInfo struct {
	UserData              UserJSON                 `json:"userData"`
	AssetBalances         map[string]AssetBalances `json:"assetBalances"`         //map of wallet public key and the asset balances
	NFTs                  map[string][]NFT         `json:"nfts"`                  //map of nft wallet and nfts
	WalletsSharedWithUser []WalletsSharedWithUser  `json:"walletsSharedWithUser"` //shows all the third party access granted to this user
	DefaultAssets         []DefaultAsset           `json:"defaultAssets"`
}

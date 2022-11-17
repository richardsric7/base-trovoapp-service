package users

type WalletsSharedWithUser struct {
	Owner             string          `json:"owner"`
	WalletPublicKey   string          `json:"walletPublicKey"`
	Permission        string          `json:"permission"`
	WalletAlias       string          `json:"walletAlias"`
	WalletDescription string          `json:"walletDescription"`
	AssetBalances     AssetBalances   `json:"assetBalances"`
	WalletSettings    *WalletSettings `json:"walletSettings"`
}

type WalletSettings struct {
	NumberOfApprovalsNeeded int                    `gorm:"type:integer; default:0" json:"numberOfApprovalsNeeded"`
	WalletType              int                    `gorm:"type:integer; default:0" json:"walletType"` //0=normal, 1= assetIssuing, 2= marketMaking, 3 = bulkPayment
	Permissions             []WalletPermissionJSON `gorm:"foreignKey:WalletPublicKey;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"permissions"`
}

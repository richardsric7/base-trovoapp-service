package users

type WalletsSharedWithUser struct {
	Owner             string          `json:"owner"`
	WalletAddress     string          `json:"walletAddress"`
	Permission        string          `json:"permission"`
	WalletAlias       string          `json:"walletAlias"`
	WalletDescription string          `json:"walletDescription"`
	AssetBalances     AssetBalances   `json:"assetBalances"`
	WalletSettings    *WalletSettings `json:"walletSettings"`
}

type WalletSettings struct {
	NumberOfApprovalsNeeded int                    `gorm:"type:integer; default:0" json:"numberOfApprovalsNeeded"`
	WalletType              int                    `gorm:"type:integer; default:0" json:"walletType"` //0=normal, 1= assetIssuing, 2= marketMaking, 3 = bulkPayment
	WalletThreshold         int                    `json:"walletThreshold"`                           //0=no shared access, 1 = view-Only shared access, 2 = approver is present
	Permissions             []WalletPermissionJSON `gorm:"foreignKey:WalletAddress;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"permissions"`
}

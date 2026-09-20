package payments

type TrackedWallet struct {
	// ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	ID          string  `gorm:""`
	Address     string  `gorm:"index:idx_tracked_wallet_address,unique"`
	TempAddress *string `gorm:"index:idx_tracked_wallet_temp_address,unique"`
	Alias       string  `gorm:"index:idx_tracked_wallet_alias"`
	Name        string  `gorm:"index:idx_tracked_wallet_name"`
}

type TrackedAddress struct {
	Address string `gorm:"primaryKey"`
}

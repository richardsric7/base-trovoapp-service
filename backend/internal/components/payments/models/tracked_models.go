package payments

type TrackedWallet struct {
	// ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	ID            string  `gorm:""`
	PublicKey     string  `gorm:"index:idx_tracked_wallet_public_key,unique"`
	TempPublicKey *string `gorm:"index:idx_tracked_wallet_temp_public_key,unique"`
	Alias         string  `gorm:"index:idx_tracked_wallet_alias"`
	Name          string  `gorm:"index:idx_tracked_wallet_name"`
}

type TrackedPublicKey struct {
	PublicKey string `gorm:"primaryKey"`
}

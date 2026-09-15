package users

import "time"

//ReservedName holds model struct for ReservedName table
type ReservedName struct {
	ID           uint64 `gorm:"primaryKey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	ReservedName *string `gorm:"size:50;not null;index:unique_reserved_name, unique;index:idx_reserved_status"`
	Status       *uint64 `gorm:"default:0;index:idx_reserved_status"`
}

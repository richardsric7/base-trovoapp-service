package currency

import (
	"time"
)

//Announcement model
type Announcement struct {
	ID             uint64    `gorm:"primaryKey" json:"id"`
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"createdAt"`
	Expiry         time.Time `json:"expiry"`
	Message        string    `gorm:"not null" json:"message"`
	Title          string    `gorm:"size:100;null" json:"title"`
	BroadcastLevel string    `gorm:"size:65;not null;default:'ALL'" json:"broadcastLevel"`
	Level          string    `gorm:"size:58;null;default:'ALL'" json:"Level"`
	CreatedBy      uint64    `json:"-"`
}

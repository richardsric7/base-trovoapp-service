package models

import (
	p2pErrors "admin-panel-dashboard/internal/errors"
	"log"

	"gorm.io/gorm"
)

// UserInfo model for bantu user directory info
type UserInfo struct {
	ID                    string          `json:"-"`
	Username              string          `json:"username"`
	Address               string          `json:"-"`
	Email                 string          `json:"email"`
	PublicIP              string          `json:"-"`
	LastName              string          `json:"lastName"`
	FirstName             string          `json:"firstName"`
	Mobile                string          `json:"mobile"`
	Telegram              *int64          `json:"-"`
	TelegramNotifications uint            `json:"-"`
	ContactPhone          string          `json:"contactPhone"`
	ImageThumbnail        string          `json:"imageThumbnail"`
	CountryCode           string          `json:"CountryCode"`
	Suspended             uint            `json:"suspended"`
	KYCLevel              uint            `json:"kycLevel"`
	MaxAssetPerOffer      float64         `json:"maxAssetPerOffer"`
	MaxAssetPerOrder      float64         `json:"maxAssetPerOrder"`
	TradeStats            MakerStats      `json:"tradeStats"`
	TakerReputation       TakerReputation `json:"takerReputation"`
	Offline               uint            `json:"offline"`
	AdminLevel            uint            `json:"adminLevel"`
	TakerFee              string          `json:"takerFee"`
}

// MakerStats/TakerReputation surface this user's real P2P trading
// performance from the new P2P module's MerchantPerformance/
// CustomerPerformance tables (see internal/models/p2p.go), replacing the
// dead legacy P2P platform's maker_stats/taker_reputations tables these
// used to read from (which always returned zero-value structs).
type MakerStats struct {
	CompletedTrades int64  `json:"completedTrades"`
	CompletionRate  string `json:"completionRate"`
	DisputesOpened  int64  `json:"disputesOpened"`
}

type TakerReputation struct {
	CompletedTrades int64  `json:"completedTrades"`
	CompletionRate  string `json:"completionRate"`
	DisputesOpened  int64  `json:"disputesOpened"`
}

func (u *UserInfo) GetMakerStat(db *gorm.DB) (rating MakerStats) {
	var p MerchantPerformance
	if err := db.Where("merchant_id = ?", u.ID).First(&p).Error; err != nil {
		return MakerStats{}
	}
	return MakerStats{CompletedTrades: p.CompletedTrades, CompletionRate: p.CompletionRate, DisputesOpened: p.DisputesOpened}
}

func (u *UserInfo) GetTakerReputation(db *gorm.DB) (rep TakerReputation) {
	var p CustomerPerformance
	if err := db.Where("customer_id = ?", u.ID).First(&p).Error; err != nil {
		return TakerReputation{}
	}
	return TakerReputation{CompletedTrades: p.CompletedTrades, CompletionRate: p.CompletionRate, DisputesOpened: p.DisputesOpened}
}

func (u *UserInfo) ToggleOffline(db *gorm.DB) (updatedState uint, err error) {

	type User struct {
		ID       string `gorm:"size:100"`
		Username string `gorm:"size:100"`
		Offline  uint   `gorm:"not null;default:0" json:"offline"`
	}
	var userState User
	username := u.Username
	log.Printf("[ToggleOffline] User:[%v], Current Offline State:[%v]\n", u.Username, u.Offline)
	if u.Offline == 1 {

		log.Printf("[ToggleOffline] User:[%v], Switching ONLINE\n", u.Username)

		err = db.Raw(`WITH o as (update offers set offline = 0 where maker = ?),
	u as (update users set offline = 0, updated_at = now() where username = ? returning id, username, offline)
	select * from u`, u.Username, username).Scan(&userState).Error

	} else {

		log.Printf("[ToggleOffline] User:[%v], Switching OFFLINE\n", u.Username)

		err = db.Raw(`WITH o as (update offers set offline = 1 where maker = ?),
		u as (update users set offline = 1, updated_at = now() where username = ? returning id, username, offline)
		select * from u`, u.Username, username).Scan(&userState).Error

	}

	if err != nil {
		log.Printf("[ToggleOffline] error fetching user data: [%v]\n", err)
		return 0, &p2pErrors.CustomError{
			Param:      "offline",
			Err:        "error toggling offline state",
			ErrMessage: "Unable to switch user offline state",
		}
	}

	// u.Offline = budsState
	updatedState = userState.Offline
	log.Printf("[ToggleOffline] User:[%v], State of Offline after execution:[%v]\n", u.Username, userState.Offline)

	return updatedState, nil
}

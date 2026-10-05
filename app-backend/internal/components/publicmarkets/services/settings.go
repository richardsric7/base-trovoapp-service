// Package publicmarkets is the Public Markets engine: creation and
// redemption from every channel (the Trovo App and exchange partners),
// session net batches to Dealing Members and Custodians, settlement,
// reconciliation, prices, exchange webhooks and dividends.
package publicmarkets

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	pm "trovo-wallet-api/internal/components/publicmarkets/models"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Lagos is West Africa Time, the markets' (and the record dates') clock.
var Lagos = func() *time.Location {
	if l, err := time.LoadLocation("Africa/Lagos"); err == nil {
		return l
	}
	return time.FixedZone("WAT", 3600)
}()

// Error is a refusal with an HTTP status and a stable code for the apps
// and exchanges.
type Error struct {
	Status  int    `json:"-"`
	Code    string `json:"error"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

func (e *Error) Error() string { return e.Message }

func refuse(status int, code, field, format string, args ...interface{}) *Error {
	return &Error{Status: status, Code: code, Field: field, Message: fmt.Sprintf(format, args...)}
}

// AsError returns err as an *Error when it is one.
func AsError(err error) (*Error, bool) {
	var e *Error
	ok := errors.As(err, &e)
	return e, ok
}

// d parses a decimal string, zero when empty or invalid.
func d(s string) decimal.Decimal {
	v, err := decimal.NewFromString(strings.TrimSpace(s))
	if err != nil {
		return decimal.Zero
	}
	return v
}

// LoadSettings returns the settings row, creating it with the defaults.
func LoadSettings(db *gorm.DB) pm.Settings {
	var s pm.Settings
	if db.First(&s, 1).Error == nil {
		return s
	}
	s = pm.Settings{ID: 1, TradeFeePercent: "0.25", FundingAssetCode: "CNGN", NetCreationThreshold: "15000000", ApprovalsRequired: 2,
		PriceStaleMinutes: 15, InstructionMaxAttempts: 5, SettlementSLAHours: 72, WebhookMaxAttempts: 9, BatchIntervalMinutes: 15,
		ReconciliationHour: 6, WHTResidentPercent: "10", WHTNonResidentPercent: "10", WHTMissingTaxIDPercent: "10",
		SubstantialHoldingPercent: "5", AUMFeePercent: "0.75", FXSpreadPercent: "0.50",
		RateLimitTiers: `{"Tier 1":1200,"Tier 2":3000}`, RevenueShareTiers: `{"Tier A":"20%","Tier B":"30%"}`,
		NGXOpen: "10:00", NGXClose: "14:30", FMDQOpen: "09:00", FMDQClose: "16:00"}
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&s)
	db.First(&s, 1)
	return s
}

// FeePercent is the trade fee of an asset (its own, or the settings').
func FeePercent(a *pm.Asset, s pm.Settings) decimal.Decimal {
	if strings.TrimSpace(a.FeePercent) != "" {
		return d(a.FeePercent)
	}
	return d(s.TradeFeePercent)
}

// CSV splits a comma-separated list, trimmed and lower-cased.
func CSV(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		if p = strings.ToLower(strings.TrimSpace(p)); p != "" {
			out = append(out, p)
		}
	}
	return out
}

// RateLimitForTier returns the requests per minute of a rate-limit tier
// (0 when the tier is unknown: the routes' defaults apply).
func RateLimitForTier(s pm.Settings, tier string) int {
	var m map[string]int
	if json.Unmarshal([]byte(s.RateLimitTiers), &m) != nil {
		return 0
	}
	return m[tier]
}

// Session is a market's trading window on a day.
type Session struct {
	Market string
	Open   bool
	OpenAt time.Time // today's (or the next) opening
	Close  time.Time
}

func clockOn(day time.Time, hhmm string) time.Time {
	t, err := time.Parse("15:04", hhmm)
	if err != nil {
		t, _ = time.Parse("15:04", "10:00")
	}
	return time.Date(day.Year(), day.Month(), day.Day(), t.Hour(), t.Minute(), 0, 0, Lagos)
}

func tradingDay(day time.Time, s pm.Settings) bool {
	if day.Weekday() == time.Saturday || day.Weekday() == time.Sunday {
		return false
	}
	for _, h := range strings.Split(s.MarketHolidays, ",") {
		if strings.TrimSpace(h) == day.Format("2006-01-02") {
			return false
		}
	}
	return true
}

// MarketSession returns the market's session around now: whether it is
// open, and when it closes or next opens. Holidays come from the settings.
func MarketSession(market string, now time.Time, s pm.Settings) Session {
	openS, closeS := s.NGXOpen, s.NGXClose
	if market == pm.MarketFMDQ {
		openS, closeS = s.FMDQOpen, s.FMDQClose
	}
	now = now.In(Lagos)
	open, close := clockOn(now, openS), clockOn(now, closeS)
	if tradingDay(now, s) && !now.Before(open) && now.Before(close) {
		return Session{Market: market, Open: true, OpenAt: open, Close: close}
	}
	day := now
	if !tradingDay(day, s) || !now.Before(open) {
		day = day.AddDate(0, 0, 1)
	}
	for i := 0; i < 14 && !tradingDay(day, s); i++ {
		day = day.AddDate(0, 0, 1)
	}
	return Session{Market: market, Open: false, OpenAt: clockOn(day, openS), Close: clockOn(day, closeS)}
}

// SessionDate is the trading day an order or batch belongs to.
func SessionDate(t time.Time) string { return t.In(Lagos).Format("2006-01-02") }

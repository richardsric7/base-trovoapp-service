package publicmarkets

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Lagos is the markets' (and the engine's) time zone.
var Lagos = func() *time.Location {
	if l, err := time.LoadLocation("Africa/Lagos"); err == nil {
		return l
	}
	return time.FixedZone("WAT", 3600)
}()

// Service works on app-backend's database (DB), where the engine runs.
type Service struct {
	DB    *gorm.DB
	Chain ContractReader // optional: verifies token contracts on registration
	Now   func() time.Time
}

func (s *Service) now() time.Time {
	if s.Now != nil {
		return s.Now().UTC()
	}
	return time.Now().UTC()
}

// Error is a refusal with its HTTP status.
type Error struct {
	Status  int    `json:"-"`
	Message string `json:"error"`
	Field   string `json:"field,omitempty"`
}

func (e *Error) Error() string { return e.Message }

func refuse(status int, format string, args ...interface{}) error {
	return &Error{Status: status, Message: fmt.Sprintf(format, args...)}
}

func invalid(field, format string, args ...interface{}) error {
	return &Error{Status: http.StatusBadRequest, Field: field, Message: fmt.Sprintf(format, args...)}
}

func notFound(what string, err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return refuse(http.StatusNotFound, "%s not found", what)
	}
	return err
}

func d(s string) decimal.Decimal {
	v, err := decimal.NewFromString(strings.TrimSpace(s))
	if err != nil {
		return decimal.Zero
	}
	return v
}

func randomHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func trimTo(s string, n int) string {
	if len(s) > n {
		return s[:n]
	}
	return s
}

// CSV splits a comma-separated setting into lower-case entries.
func CSV(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		if p = strings.ToLower(strings.TrimSpace(p)); p != "" {
			out = append(out, p)
		}
	}
	return out
}

func contains(list []string, s string) bool {
	s = strings.ToLower(strings.TrimSpace(s))
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}

// Page is a page request.
type Page struct {
	Page  int
	Limit int
}

func (p Page) apply(q *gorm.DB) *gorm.DB {
	if p.Limit <= 0 || p.Limit > 200 {
		p.Limit = 50
	}
	if p.Page < 1 {
		p.Page = 1
	}
	return q.Offset((p.Page - 1) * p.Limit).Limit(p.Limit)
}

// LoadSettings returns the settings row (the engine creates it; the
// defaults below match the engine's when it has not run yet).
func (s *Service) LoadSettings() Settings {
	var st Settings
	if s.DB.First(&st, "id = 1").Error != nil {
		st = Settings{ID: 1, TradeFeePercent: "0.25", FundingAssetCode: "CNGN", NetCreationThreshold: "15000000", ApprovalsRequired: 2,
			PriceStaleMinutes: 15, InstructionMaxAttempts: 5, SettlementSLAHours: 72, WebhookMaxAttempts: 9, BatchIntervalMinutes: 15,
			ReconciliationHour: 6, WHTResidentPercent: "10", WHTNonResidentPercent: "10", WHTMissingTaxIDPercent: "10",
			SubstantialHoldingPercent: "5", AUMFeePercent: "0.75", FXSpreadPercent: "0.50", NGXOpen: "10:00", NGXClose: "14:30",
			FMDQOpen: "09:00", FMDQClose: "16:00", UpdatedAt: s.now()}
		s.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&st)
	}
	return st
}

// rateLimitForTier is a tier's requests per minute (0: no override).
func rateLimitForTier(st Settings, tier string) int {
	var m map[string]int
	if json.Unmarshal([]byte(st.RateLimitTiers), &m) != nil {
		return 0
	}
	return m[tier]
}

// marketOpen says whether a market is in session at t (weekdays, the
// configured hours, outside the holidays), as the engine decides it.
func marketOpen(st Settings, market string, t time.Time) bool {
	t = t.In(Lagos)
	if t.Weekday() == time.Saturday || t.Weekday() == time.Sunday {
		return false
	}
	if contains(CSV(st.MarketHolidays), t.Format("2006-01-02")) {
		return false
	}
	open, close := st.NGXOpen, st.NGXClose
	if strings.EqualFold(market, MarketFMDQ) {
		open, close = st.FMDQOpen, st.FMDQClose
	}
	hm := t.Format("15:04")
	return hm >= open && hm < close
}

var (
	codePattern  = regexp.MustCompile(`^[A-Z0-9][A-Z0-9.-]{0,15}$`)
	isinPattern  = regexp.MustCompile(`^[A-Z]{2}[A-Z0-9]{9}[0-9]$`)
	emailPattern = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
	hhmmPattern  = regexp.MustCompile(`^([01][0-9]|2[0-3]):[0-5][0-9]$`)
)

func validAddress(a string) bool { return common.IsHexAddress(strings.TrimSpace(a)) }

// ---------------------------------------------------------------- jobs

// requestJob asks the engine to run a job now (it polls every minute).
func (s *Service) requestJob(job, target, by string) (*JobRequest, error) {
	r := JobRequest{Job: job, Target: target, RequestedBy: by, CreatedAt: s.now()}
	if err := s.DB.Create(&r).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

// JobRequests are the recent requests and their results.
func (s *Service) JobRequests(limit int) []JobRequest {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	var out []JobRequest
	s.DB.Order("id DESC").Limit(limit).Find(&out)
	return out
}

// recordManual stores an event for the engine to process, as a partner
// would have sent it (confirmations of MANUAL partners, declared
// corporate actions). The engine processes it within seconds.
func (s *Service) recordManual(kind, reference string, payload interface{}, by string) (*PartnerEvent, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	ev := PartnerEvent{Source: trimTo("MANUAL:"+by, 40), Kind: kind, Reference: trimTo(reference, 160), Payload: string(body), CreatedAt: s.now()}
	res := s.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&ev)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, refuse(http.StatusConflict, "this was already recorded")
	}
	return &ev, nil
}

// PartnerEvents lists inbound events (source filter: a prefix such as
// MANUAL, CUSTODIAN or DEALING_MEMBER).
func (s *Service) PartnerEvents(source, kind string, p Page) ([]PartnerEvent, int64) {
	q := s.DB.Model(&PartnerEvent{})
	if source != "" {
		q = q.Where("source LIKE ?", strings.ToUpper(source)+"%")
	}
	if kind != "" {
		q = q.Where("kind = ?", kind)
	}
	var total int64
	q.Count(&total)
	var out []PartnerEvent
	p.apply(q.Order("id DESC")).Find(&out)
	return out, total
}

// Health is System Health's Public Markets section.
type Health struct {
	Jobs            []JobRun     `json:"jobs"`
	Requests        []JobRequest `json:"requests"`
	PendingEvents   int64        `json:"pendingEvents"`
	PendingMock     int64        `json:"pendingMockEvents"`
	PendingWebhooks int64        `json:"pendingWebhooks"`
}

func (s *Service) Health() Health {
	var h Health
	s.DB.Order("job").Find(&h.Jobs)
	h.Requests = s.JobRequests(20)
	s.DB.Model(&PartnerEvent{}).Where("processed_at IS NULL").Count(&h.PendingEvents)
	s.DB.Model(&MockEvent{}).Where("delivered_at IS NULL").Count(&h.PendingMock)
	s.DB.Model(&WebhookDelivery{}).Where("status = ?", DeliveryPending).Count(&h.PendingWebhooks)
	return h
}

package publicmarkets

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// serviceLink is the slice of app-backend's service_links an exchange
// partner is built on.
type serviceLink struct {
	ID                 string
	ShortName          string
	LongName           string
	Verified           int
	Suspended          int
	Inactive           int
	RateLimitPerMinute int
}

func (s *Service) serviceLink(id string) (*serviceLink, error) {
	var sl serviceLink
	err := s.DB.Table("service_links").Select("id, short_name, long_name, verified, suspended, inactive, rate_limit_per_minute").Where("id = ?", id).Take(&sl).Error
	if err != nil {
		return nil, notFound("service link", err)
	}
	return &sl, nil
}

func (sl serviceLink) name() string {
	if sl.LongName != "" {
		return sl.LongName
	}
	return sl.ShortName
}

// ExchangeRow is an exchange partner as the list shows it.
type ExchangeRow struct {
	ExchangePartner
	Name               string     `json:"name"`
	ShortName          string     `json:"shortName"`
	RateLimitPerMinute int        `json:"rateLimitPerMinute"`
	Wallets            int64      `json:"wallets"`
	Orders30d          int64      `json:"orders30d"`
	DeadLetters        int64      `json:"deadLetters"`
	PendingWebhooks    int64      `json:"pendingWebhooks"`
	SecretRotatedUntil *time.Time `json:"previousSecretValidUntil"`
}

func (s *Service) exchangeRow(p ExchangePartner) ExchangeRow {
	r := ExchangeRow{ExchangePartner: p, SecretRotatedUntil: p.PreviousSecretExpiresAt}
	if sl, err := s.serviceLink(p.ServiceLinkID); err == nil {
		r.Name, r.ShortName, r.RateLimitPerMinute = sl.name(), sl.ShortName, sl.RateLimitPerMinute
	}
	s.DB.Model(&PartnerWallet{}).Where("service_link_id = ? AND status = ?", p.ServiceLinkID, "active").Count(&r.Wallets)
	s.DB.Model(&Order{}).Where("service_link_id = ? AND created_at >= ?", p.ServiceLinkID, s.now().AddDate(0, 0, -30)).Count(&r.Orders30d)
	s.DB.Model(&WebhookDelivery{}).Where("service_link_id = ? AND status = ?", p.ServiceLinkID, DeliveryDeadLetter).Count(&r.DeadLetters)
	s.DB.Model(&WebhookDelivery{}).Where("service_link_id = ? AND status = ?", p.ServiceLinkID, DeliveryPending).Count(&r.PendingWebhooks)
	return r
}

// Exchanges lists the onboarded exchanges.
func (s *Service) Exchanges(status, search string) []ExchangeRow {
	q := s.DB.Model(&ExchangePartner{})
	if status != "" {
		q = q.Where("status = ?", strings.ToLower(status))
	}
	var partners []ExchangePartner
	q.Order("created_at").Find(&partners)
	out := make([]ExchangeRow, 0, len(partners))
	for _, p := range partners {
		r := s.exchangeRow(p)
		if v := strings.ToLower(strings.TrimSpace(search)); v != "" && !strings.Contains(strings.ToLower(r.Name+" "+r.ShortName+" "+r.ServiceLinkID), v) {
			continue
		}
		out = append(out, r)
	}
	return out
}

// ExchangeCandidate is a verified service link not yet onboarded.
type ExchangeCandidate struct {
	ID        string `json:"id"`
	ShortName string `json:"shortName"`
	LongName  string `json:"longName"`
}

// ExchangeCandidates are the service links an exchange can be onboarded on.
func (s *Service) ExchangeCandidates() []ExchangeCandidate {
	var out []ExchangeCandidate
	s.DB.Table("service_links").Select("id, short_name, long_name").
		Where("verified = 1 AND suspended = 0 AND id NOT IN (?)", s.DB.Model(&ExchangePartner{}).Select("service_link_id")).Order("short_name").Scan(&out)
	return out
}

// ExchangeDetail is an exchange's page.
type ExchangeDetail struct {
	ExchangeRow
	Deliveries []WebhookDelivery     `json:"deliveries"`
	Ledger     []ExchangeLedgerEntry `json:"ledger"`
	Orders     []Order               `json:"recentOrders"`
	Requests   []JobRequest          `json:"withdrawals"`
}

// ExchangeDetail returns an exchange with its recent deliveries, balance
// movements and orders (delivery filter: a status).
func (s *Service) ExchangeDetail(id, deliveryStatus string) (*ExchangeDetail, error) {
	p, err := s.partner(id)
	if err != nil {
		return nil, err
	}
	out := &ExchangeDetail{ExchangeRow: s.exchangeRow(*p)}
	q := s.DB.Where("service_link_id = ?", p.ServiceLinkID)
	if deliveryStatus != "" {
		q = q.Where("status = ?", strings.ToUpper(deliveryStatus))
	}
	q.Order("created_at DESC").Limit(100).Find(&out.Deliveries)
	s.DB.Where("service_link_id = ?", p.ServiceLinkID).Order("id DESC").Limit(100).Find(&out.Ledger)
	s.DB.Where("service_link_id = ?", p.ServiceLinkID).Order("created_at DESC").Limit(20).Find(&out.Orders)
	s.DB.Where("job = ? AND target LIKE ?", "EXCHANGE_WITHDRAWAL", p.ServiceLinkID+":%").Order("id DESC").Limit(20).Find(&out.Requests)
	return out, nil
}

func (s *Service) partner(id string) (*ExchangePartner, error) {
	var p ExchangePartner
	if err := s.DB.First(&p, "service_link_id = ?", id).Error; err != nil {
		return nil, notFound("exchange partner", err)
	}
	return &p, nil
}

// ExchangeRequest onboards or edits an exchange.
type ExchangeRequest struct {
	ServiceLinkID        string `json:"serviceLinkId"`
	CallbackURL          string `json:"callbackUrl"`
	Environment          string `json:"environment"` // sandbox | production
	RateLimitTier        string `json:"rateLimitTier"`
	RevenueShareTier     string `json:"revenueShareTier"`
	FundingAddress       string `json:"fundingAddress"`
	TechContact          string `json:"techContact"`
	ConfirmationSLAHours *int   `json:"confirmationSlaHours"`
}

func (r ExchangeRequest) check(st Settings) error {
	if v := strings.TrimSpace(r.CallbackURL); v != "" {
		u, err := url.Parse(v)
		if err != nil || u.Host == "" || (u.Scheme != "https" && !(u.Scheme == "http" && strings.EqualFold(r.Environment, "sandbox"))) {
			return invalid("callbackUrl", "callbackUrl must be an https URL (http only in sandbox)")
		}
	}
	if e := strings.ToLower(strings.TrimSpace(r.Environment)); e != "" && e != "sandbox" && e != "production" {
		return invalid("environment", "environment must be sandbox or production")
	}
	if v := strings.TrimSpace(r.FundingAddress); v != "" && !validAddress(v) {
		return invalid("fundingAddress", "fundingAddress must be a 0x address")
	}
	if r.RateLimitTier != "" && st.RateLimitTiers != "" && rateLimitForTier(st, r.RateLimitTier) == 0 {
		return invalid("rateLimitTier", "unknown tier %q (Settings › Thresholds & Limits › Rate-limit tiers)", r.RateLimitTier)
	}
	if r.ConfirmationSLAHours != nil && (*r.ConfirmationSLAHours < 0 || *r.ConfirmationSLAHours > 24*30) {
		return invalid("confirmationSlaHours", "confirmationSlaHours must be between 0 and 720")
	}
	return nil
}

// applyTier sets the service link's rate limit from the exchange's tier.
func (s *Service) applyTier(tx *gorm.DB, serviceLinkID, tier string, st Settings) error {
	if tier == "" {
		return nil
	}
	return tx.Table("service_links").Where("id = ?", serviceLinkID).Update("rate_limit_per_minute", rateLimitForTier(st, tier)).Error
}

func newSecret() string { return "pmsk_" + randomHex(24) }

// Onboard makes a verified service link an exchange partner. The signing
// secret is returned once: hand it to the exchange securely.
func (s *Service) Onboard(r ExchangeRequest, by string) (*ExchangeRow, string, error) {
	sl, err := s.serviceLink(strings.TrimSpace(r.ServiceLinkID))
	if err != nil {
		return nil, "", err
	}
	if sl.Verified == 0 || sl.Suspended != 0 || sl.Inactive != 0 {
		return nil, "", refuse(http.StatusConflict, "the service link must be verified and active")
	}
	if _, err := s.partner(sl.ID); err == nil {
		return nil, "", refuse(http.StatusConflict, "this service link is already an exchange partner")
	}
	st := s.LoadSettings()
	if err := r.check(st); err != nil {
		return nil, "", err
	}
	if strings.TrimSpace(r.CallbackURL) == "" {
		return nil, "", invalid("callbackUrl", "callbackUrl is required")
	}
	if strings.TrimSpace(r.FundingAddress) == "" {
		return nil, "", invalid("fundingAddress", "fundingAddress is required: deposits are accepted only from it")
	}
	env := strings.ToLower(strings.TrimSpace(r.Environment))
	if env == "" {
		env = "sandbox"
	}
	secret := newSecret()
	now := s.now()
	p := ExchangePartner{ServiceLinkID: sl.ID, Status: "active", Environment: env, CallbackURL: strings.TrimSpace(r.CallbackURL), SigningSecret: secret,
		RateLimitTier: r.RateLimitTier, RevenueShareTier: r.RevenueShareTier, FundingAddress: common.HexToAddress(r.FundingAddress).Hex(),
		Balance: "0", Reserved: "0", TechContact: r.TechContact, CreatedAt: now, UpdatedAt: now}
	if r.ConfirmationSLAHours != nil {
		p.ConfirmationSLAHours = *r.ConfirmationSLAHours
	}
	err = s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&p).Error; err != nil {
			return err
		}
		return s.applyTier(tx, sl.ID, p.RateLimitTier, st)
	})
	if err != nil {
		return nil, "", err
	}
	row := s.exchangeRow(p)
	return &row, secret, nil
}

// UpdateExchange edits an exchange's integration settings.
func (s *Service) UpdateExchange(id string, r ExchangeRequest, by string) (*ExchangeRow, error) {
	p, err := s.partner(id)
	if err != nil {
		return nil, err
	}
	st := s.LoadSettings()
	if r.Environment == "" {
		r.Environment = p.Environment
	}
	if err := r.check(st); err != nil {
		return nil, err
	}
	updates := map[string]interface{}{"updated_at": s.now()}
	if v := strings.TrimSpace(r.CallbackURL); v != "" {
		updates["callback_url"] = v
	}
	updates["environment"] = strings.ToLower(r.Environment)
	if r.RateLimitTier != "" {
		updates["rate_limit_tier"] = r.RateLimitTier
	}
	if r.RevenueShareTier != "" {
		updates["revenue_share_tier"] = r.RevenueShareTier
	}
	if v := strings.TrimSpace(r.FundingAddress); v != "" {
		updates["funding_address"] = common.HexToAddress(v).Hex()
	}
	if r.TechContact != "" {
		updates["tech_contact"] = r.TechContact
	}
	if r.ConfirmationSLAHours != nil {
		updates["confirmation_sla_hours"] = *r.ConfirmationSLAHours
	}
	err = s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&ExchangePartner{}).Where("service_link_id = ?", p.ServiceLinkID).Updates(updates).Error; err != nil {
			return err
		}
		return s.applyTier(tx, p.ServiceLinkID, r.RateLimitTier, st)
	})
	if err != nil {
		return nil, err
	}
	p, _ = s.partner(id)
	row := s.exchangeRow(*p)
	return &row, nil
}

// secretOverlap is how long the previous signing secret keeps working
// after a rotation, so the exchange can deploy the new one.
const secretOverlap = 24 * time.Hour

// RotateSecret issues a new signing secret (returned once); the previous
// one stays valid for 24 hours.
func (s *Service) RotateSecret(id, by string) (string, *time.Time, error) {
	p, err := s.partner(id)
	if err != nil {
		return "", nil, err
	}
	secret := newSecret()
	until := s.now().Add(secretOverlap)
	res := s.DB.Model(&ExchangePartner{}).Where("service_link_id = ? AND signing_secret = ?", p.ServiceLinkID, p.SigningSecret).Updates(map[string]interface{}{
		"previous_signing_secret": p.SigningSecret, "previous_secret_expires_at": &until, "signing_secret": secret, "updated_at": s.now()})
	if res.Error != nil {
		return "", nil, res.Error
	}
	if res.RowsAffected == 0 {
		return "", nil, refuse(http.StatusConflict, "the secret was rotated at the same time; reload")
	}
	return secret, &until, nil
}

// SetExchangeStatus suspends (API calls answer 403) or reactivates an
// exchange.
func (s *Service) SetExchangeStatus(id string, active bool, by string) (*ExchangeRow, error) {
	from, to := "active", "suspended"
	if active {
		from, to = to, from
	}
	res := s.DB.Model(&ExchangePartner{}).Where("service_link_id = ? AND status = ?", id, from).Updates(map[string]interface{}{"status": to, "updated_at": s.now()})
	if res.RowsAffected == 0 {
		if _, err := s.partner(id); err != nil {
			return nil, err
		}
		return nil, refuse(http.StatusConflict, "the exchange is already %s", to)
	}
	p, _ := s.partner(id)
	row := s.exchangeRow(*p)
	return &row, nil
}

// RequestWithdrawal asks the engine to pay part of an exchange's balance
// back to its funding wallet from the treasury.
func (s *Service) RequestWithdrawal(id, amount, by string) (*JobRequest, error) {
	p, err := s.partner(id)
	if err != nil {
		return nil, err
	}
	a, err := decimal.NewFromString(strings.TrimSpace(amount))
	if err != nil || !a.IsPositive() {
		return nil, invalid("amount", "amount must be positive")
	}
	if a.GreaterThan(d(p.Balance)) {
		return nil, invalid("amount", "amount is above the exchange's balance (%s)", p.Balance)
	}
	if !validAddress(p.FundingAddress) {
		return nil, refuse(http.StatusConflict, "the exchange has no funding wallet to pay to")
	}
	var pending int64
	s.DB.Model(&JobRequest{}).Where("job = ? AND target LIKE ? AND done_at IS NULL", "EXCHANGE_WITHDRAWAL", p.ServiceLinkID+":%").Count(&pending)
	if pending > 0 {
		return nil, refuse(http.StatusConflict, "a withdrawal for this exchange is already waiting")
	}
	return s.requestJob("EXCHANGE_WITHDRAWAL", p.ServiceLinkID+":"+a.String(), by)
}

// ReplayWebhook sends a dead-lettered delivery again.
func (s *Service) ReplayWebhook(deliveryID, by string) (*WebhookDelivery, error) {
	now := s.now()
	res := s.DB.Model(&WebhookDelivery{}).Where("id = ? AND status = ?", deliveryID, DeliveryDeadLetter).Updates(map[string]interface{}{
		"status": DeliveryPending, "attempts": 0, "next_attempt_at": &now, "last_error": trimTo("Replayed by "+by, 300), "updated_at": now})
	if res.RowsAffected == 0 {
		return nil, refuse(http.StatusConflict, "only a dead-lettered delivery can be replayed")
	}
	var w WebhookDelivery
	s.DB.First(&w, "id = ?", deliveryID)
	return &w, nil
}

// ReplayDeadLetters replays every dead-lettered delivery of an exchange.
func (s *Service) ReplayDeadLetters(id, by string) (int64, error) {
	if _, err := s.partner(id); err != nil {
		return 0, err
	}
	now := s.now()
	res := s.DB.Model(&WebhookDelivery{}).Where("service_link_id = ? AND status = ?", id, DeliveryDeadLetter).Updates(map[string]interface{}{
		"status": DeliveryPending, "attempts": 0, "next_attempt_at": &now, "last_error": trimTo("Replayed by "+by, 300), "updated_at": now})
	return res.RowsAffected, res.Error
}

// ---------------------------------------------------------------- wallets

// WalletStats head the wallet provisioning page.
type WalletStats struct {
	Provisioned     int64  `json:"provisioned"`
	ConsentPercent  string `json:"consentPercent"`
	TaxDataPercent  string `json:"taxDataPercent"`
	Rejected30d     int64  `json:"rejected30d"`
	DeployedWallets int64  `json:"deployedWallets"`
}

// WalletRow is a provisioning request; personal data is masked unless
// the admin may see it.
type WalletRow struct {
	PartnerWallet
	ExchangeName string `json:"exchangeName"`
}

func mask(s string, keepStart, keepEnd int) string {
	r := []rune(strings.TrimSpace(s))
	if len(r) == 0 {
		return ""
	}
	if len(r) <= keepStart+keepEnd {
		return strings.Repeat("*", len(r))
	}
	return string(r[:keepStart]) + strings.Repeat("*", len(r)-keepStart-keepEnd) + string(r[len(r)-keepEnd:])
}

func maskName(name string) string {
	parts := strings.Fields(name)
	for i, p := range parts {
		parts[i] = mask(p, 1, 0)
	}
	return strings.Join(parts, " ")
}

// Wallets lists provisioning requests (status: active | rejected).
func (s *Service) Wallets(serviceLinkID, status, search string, reveal bool, p Page) ([]WalletRow, int64, WalletStats) {
	var st WalletStats
	s.DB.Model(&PartnerWallet{}).Where("status = ?", "active").Count(&st.Provisioned)
	s.DB.Model(&PartnerWallet{}).Where("status = ? AND deployed = ?", "active", true).Count(&st.DeployedWallets)
	var consent, tax int64
	s.DB.Model(&PartnerWallet{}).Where("status = ? AND ndpa_consent = ?", "active", true).Count(&consent)
	s.DB.Model(&PartnerWallet{}).Where("status = ? AND tax_identifier <> ''", "active").Count(&tax)
	st.ConsentPercent, st.TaxDataPercent = pct(consent, st.Provisioned), pct(tax, st.Provisioned)
	s.DB.Model(&PartnerWallet{}).Where("status = ? AND created_at >= ?", "rejected", s.now().AddDate(0, 0, -30)).Count(&st.Rejected30d)

	q := s.DB.Model(&PartnerWallet{})
	if serviceLinkID != "" {
		q = q.Where("service_link_id = ?", serviceLinkID)
	}
	if status != "" {
		q = q.Where("status = ?", strings.ToLower(status))
	}
	if v := strings.TrimSpace(search); v != "" {
		like := "%" + strings.ToLower(v) + "%"
		q = q.Where("LOWER(external_user_ref) LIKE ? OR LOWER(wallet_address) LIKE ? OR LOWER(id) LIKE ?", like, like, like)
	}
	var total int64
	q.Count(&total)
	var wallets []PartnerWallet
	p.apply(q.Order("created_at DESC")).Find(&wallets)
	names := s.exchangeNames()
	out := make([]WalletRow, 0, len(wallets))
	for _, w := range wallets {
		if !reveal {
			w.LegalName, w.TaxIdentifier = maskName(w.LegalName), mask(w.TaxIdentifier, 4, 2)
		}
		out = append(out, WalletRow{PartnerWallet: w, ExchangeName: nameOr(names, w.ServiceLinkID)})
	}
	return out, total, st
}

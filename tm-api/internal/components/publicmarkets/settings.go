package publicmarkets

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// SettingsRequest edits Settings; empty / nil fields are left as they are.
type SettingsRequest struct {
	TradeFeePercent           string  `json:"tradeFeePercent"`
	FeeWallet                 string  `json:"feeWallet"`
	NetCreationThreshold      string  `json:"netCreationThreshold"`
	ApprovalsRequired         *int    `json:"approvalsRequired"`
	NetCreationApprovers      *string `json:"netCreationApprovers"`
	DividendApprovers         *string `json:"dividendApprovers"`
	PriceStaleMinutes         *int    `json:"priceStaleMinutes"`
	InstructionMaxAttempts    *int    `json:"instructionMaxAttempts"`
	SettlementSLAHours        *int    `json:"settlementSlaHours"`
	ConfirmationSLAHours      *int    `json:"confirmationSlaHours"`
	WebhookMaxAttempts        *int    `json:"webhookMaxAttempts"`
	BatchIntervalMinutes      *int    `json:"batchIntervalMinutes"`
	ReconciliationHour        *int    `json:"reconciliationHour"`
	WHTResidentPercent        string  `json:"whtResidentPercent"`
	WHTNonResidentPercent     string  `json:"whtNonResidentPercent"`
	WHTMissingTaxIDPercent    string  `json:"whtMissingTaxIdPercent"`
	WHTWallet                 string  `json:"whtWallet"`
	SubstantialHoldingPercent string  `json:"substantialHoldingPercent"`
	RateLimitTiers            *string `json:"rateLimitTiers"`
	AUMFeePercent             string  `json:"aumFeePercent"`
	FXSpreadPercent           string  `json:"fxSpreadPercent"`
	RevenueShareTiers         *string `json:"revenueShareTiers"`
	NGXOpen                   string  `json:"ngxOpen"`
	NGXClose                  string  `json:"ngxClose"`
	FMDQOpen                  string  `json:"fmdqOpen"`
	FMDQClose                 string  `json:"fmdqClose"`
	MarketHolidays            *string `json:"marketHolidays"`
}

func percent(field, v string, max int64) error {
	if v == "" {
		return nil
	}
	x, err := decimal.NewFromString(strings.TrimSpace(v))
	if err != nil || x.IsNegative() || x.GreaterThan(decimal.NewFromInt(max)) {
		return invalid(field, "%s must be between 0 and %d", field, max)
	}
	return nil
}

func intRange(field string, v *int, min, max int) error {
	if v != nil && (*v < min || *v > max) {
		return invalid(field, "%s must be between %d and %d", field, min, max)
	}
	return nil
}

func emails(field string, v *string) (string, error) {
	if v == nil {
		return "", nil
	}
	list := CSV(*v)
	for _, e := range list {
		if !emailPattern.MatchString(e) {
			return "", invalid(field, "%q is not an email address", e)
		}
	}
	return strings.Join(list, ","), nil
}

// UpdateSettings validates and saves settings; changing the rate-limit
// tiers re-applies them to every exchange's service link.
func (s *Service) UpdateSettings(r SettingsRequest, by string) (*Settings, error) {
	st := s.LoadSettings()
	checks := []error{
		percent("tradeFeePercent", r.TradeFeePercent, 10), percent("whtResidentPercent", r.WHTResidentPercent, 100),
		percent("whtNonResidentPercent", r.WHTNonResidentPercent, 100), percent("whtMissingTaxIdPercent", r.WHTMissingTaxIDPercent, 100),
		percent("substantialHoldingPercent", r.SubstantialHoldingPercent, 100), percent("aumFeePercent", r.AUMFeePercent, 10),
		percent("fxSpreadPercent", r.FXSpreadPercent, 10), nonNegative("netCreationThreshold", r.NetCreationThreshold),
		intRange("approvalsRequired", r.ApprovalsRequired, 1, 10), intRange("priceStaleMinutes", r.PriceStaleMinutes, 1, 24*60),
		intRange("instructionMaxAttempts", r.InstructionMaxAttempts, 1, 50), intRange("settlementSlaHours", r.SettlementSLAHours, 1, 24*30),
		intRange("confirmationSlaHours", r.ConfirmationSLAHours, 0, 24*30), intRange("webhookMaxAttempts", r.WebhookMaxAttempts, 1, 50),
		intRange("batchIntervalMinutes", r.BatchIntervalMinutes, 1, 24*60), intRange("reconciliationHour", r.ReconciliationHour, 0, 23),
	}
	for _, err := range checks {
		if err != nil {
			return nil, err
		}
	}
	for field, v := range map[string]string{"feeWallet": r.FeeWallet, "whtWallet": r.WHTWallet} {
		if v != "" && !validAddress(v) {
			return nil, invalid(field, "%s must be a 0x address", field)
		}
	}
	for field, v := range map[string]string{"ngxOpen": r.NGXOpen, "ngxClose": r.NGXClose, "fmdqOpen": r.FMDQOpen, "fmdqClose": r.FMDQClose} {
		if v != "" && !hhmmPattern.MatchString(v) {
			return nil, invalid(field, "%s must be HH:MM (WAT)", field)
		}
	}
	set := func(dst *string, v string) {
		if v = strings.TrimSpace(v); v != "" {
			*dst = v
		}
	}
	set(&st.TradeFeePercent, r.TradeFeePercent)
	set(&st.FeeWallet, r.FeeWallet)
	set(&st.NetCreationThreshold, r.NetCreationThreshold)
	set(&st.WHTResidentPercent, r.WHTResidentPercent)
	set(&st.WHTNonResidentPercent, r.WHTNonResidentPercent)
	set(&st.WHTMissingTaxIDPercent, r.WHTMissingTaxIDPercent)
	set(&st.WHTWallet, r.WHTWallet)
	set(&st.SubstantialHoldingPercent, r.SubstantialHoldingPercent)
	set(&st.AUMFeePercent, r.AUMFeePercent)
	set(&st.FXSpreadPercent, r.FXSpreadPercent)
	set(&st.NGXOpen, r.NGXOpen)
	set(&st.NGXClose, r.NGXClose)
	set(&st.FMDQOpen, r.FMDQOpen)
	set(&st.FMDQClose, r.FMDQClose)
	if st.NGXOpen >= st.NGXClose || st.FMDQOpen >= st.FMDQClose {
		return nil, invalid("ngxClose", "a market must close after it opens")
	}
	for _, x := range []struct {
		dst *int
		v   *int
	}{{&st.ApprovalsRequired, r.ApprovalsRequired}, {&st.PriceStaleMinutes, r.PriceStaleMinutes}, {&st.InstructionMaxAttempts, r.InstructionMaxAttempts},
		{&st.SettlementSLAHours, r.SettlementSLAHours}, {&st.ConfirmationSLAHours, r.ConfirmationSLAHours}, {&st.WebhookMaxAttempts, r.WebhookMaxAttempts},
		{&st.BatchIntervalMinutes, r.BatchIntervalMinutes}, {&st.ReconciliationHour, r.ReconciliationHour}} {
		if x.v != nil {
			*x.dst = *x.v
		}
	}
	if v, err := emails("netCreationApprovers", r.NetCreationApprovers); err != nil {
		return nil, err
	} else if r.NetCreationApprovers != nil {
		st.NetCreationApprovers = v
	}
	if v, err := emails("dividendApprovers", r.DividendApprovers); err != nil {
		return nil, err
	} else if r.DividendApprovers != nil {
		st.DividendApprovers = v
	}
	if n := len(CSV(st.NetCreationApprovers)); n > 0 && n < st.ApprovalsRequired {
		return nil, invalid("netCreationApprovers", "list at least %d Net Creation Approvers (approvals required)", st.ApprovalsRequired)
	}
	if n := len(CSV(st.DividendApprovers)); n > 0 && n < st.ApprovalsRequired {
		return nil, invalid("dividendApprovers", "list at least %d Dividend Approvers (approvals required)", st.ApprovalsRequired)
	}
	tiersChanged := false
	if r.RateLimitTiers != nil {
		v := strings.TrimSpace(*r.RateLimitTiers)
		if v != "" {
			var m map[string]int
			if json.Unmarshal([]byte(v), &m) != nil {
				return nil, invalid("rateLimitTiers", `rateLimitTiers must be JSON like {"Tier 1": 1200, "Tier 2": 3000} (requests per minute)`)
			}
			for k, n := range m {
				if strings.TrimSpace(k) == "" || n < 1 {
					return nil, invalid("rateLimitTiers", "each tier needs a name and a positive requests-per-minute")
				}
			}
		}
		tiersChanged = v != st.RateLimitTiers
		st.RateLimitTiers = v
	}
	if r.RevenueShareTiers != nil {
		v := strings.TrimSpace(*r.RevenueShareTiers)
		if v != "" {
			var m map[string]json.Number
			if json.Unmarshal([]byte(v), &m) != nil {
				return nil, invalid("revenueShareTiers", `revenueShareTiers must be JSON like {"Tier A": 20, "Tier B": 30} (percent)`)
			}
		}
		st.RevenueShareTiers = v
	}
	if r.MarketHolidays != nil {
		days := CSV(*r.MarketHolidays)
		for _, day := range days {
			if _, err := time.Parse("2006-01-02", day); err != nil {
				return nil, invalid("marketHolidays", "%q is not a YYYY-MM-DD date", day)
			}
		}
		st.MarketHolidays = strings.Join(days, ",")
	}
	st.UpdatedAt, st.UpdatedBy = s.now(), by
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&st).Error; err != nil {
			return err
		}
		if !tiersChanged {
			return nil
		}
		var partners []ExchangePartner
		tx.Where("rate_limit_tier <> ''").Find(&partners)
		for _, p := range partners {
			if err := s.applyTier(tx, p.ServiceLinkID, p.RateLimitTier, st); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &st, nil
}

// ---------------------------------------------------------------- dealing members

// DealingMemberRow is a Dealing Member with its asset count.
type DealingMemberRow struct {
	DealingMember
	Assets int64 `json:"assets"`
}

func (s *Service) DealingMembers() []DealingMemberRow {
	var dms []DealingMember
	s.DB.Order("dealing_member_name").Find(&dms)
	out := make([]DealingMemberRow, 0, len(dms))
	for _, dm := range dms {
		r := DealingMemberRow{DealingMember: dm}
		s.DB.Model(&Asset{}).Where("dealing_member_id = ?", dm.ID).Count(&r.Assets)
		out = append(out, r)
	}
	return out
}

// PartnerRequest creates or edits a Dealing Member, or a Custodian's
// Public Markets integration.
type PartnerRequest struct {
	Name                string   `json:"name"`
	Address             string   `json:"address"`
	Country             string   `json:"country"`
	CSCSMemberCode      string   `json:"cscsMemberCode"`
	RequirementDocument string   `json:"requirementDocument"`
	FeePercent          *float64 `json:"feePercent"`
	FeeFixed            *float64 `json:"feeFixed"`
	Code                string   `json:"code"`
	NomineeName         string   `json:"nomineeName"`
	Mode                string   `json:"mode"` // MOCK | REST | MANUAL
	Transport           string   `json:"transport"`
	AuthScheme          string   `json:"authScheme"` // HMAC | MTLS | NONE
	BaseURL             string   `json:"baseUrl"`
	CredentialsRef      string   `json:"credentialsRef"`
	Active              *bool    `json:"active"`
}

func (r *PartnerRequest) checkIntegration() error {
	r.Code = strings.ToUpper(strings.TrimSpace(r.Code))
	if r.Code != "" && !codePattern.MatchString(r.Code) {
		return invalid("code", "code is up to 16 capital letters or digits (sent as X-Partner-Code)")
	}
	r.Mode = strings.ToUpper(strings.TrimSpace(r.Mode))
	switch r.Mode {
	case "", ModeMock, ModeManual:
	case ModeREST:
		if !strings.HasPrefix(strings.TrimSpace(r.BaseURL), "https://") {
			return invalid("baseUrl", "a REST partner needs its https base URL")
		}
		if strings.TrimSpace(r.CredentialsRef) == "" {
			return invalid("credentialsRef", "a REST partner needs a credentials reference (env:NAME or vault://path#FIELD)")
		}
	default:
		return invalid("mode", "mode must be MOCK, REST or MANUAL")
	}
	r.AuthScheme = strings.ToUpper(strings.TrimSpace(r.AuthScheme))
	if r.AuthScheme != "" && r.AuthScheme != "HMAC" && r.AuthScheme != "MTLS" && r.AuthScheme != "NONE" {
		return invalid("authScheme", "authScheme must be HMAC, MTLS or NONE")
	}
	if v := strings.TrimSpace(r.CredentialsRef); v != "" && !strings.HasPrefix(v, "env:") && !strings.HasPrefix(v, "vault://") {
		return invalid("credentialsRef", "credentialsRef must be env:NAME or vault://path#FIELD (never the secret itself)")
	}
	return nil
}

// SaveDealingMember creates (id 0) or edits a Dealing Member.
func (s *Service) SaveDealingMember(id uint64, r PartnerRequest) (*DealingMember, error) {
	if err := r.checkIntegration(); err != nil {
		return nil, err
	}
	var dm DealingMember
	if id != 0 {
		if err := s.DB.First(&dm, "id = ?", id).Error; err != nil {
			return nil, notFound("dealing member", err)
		}
	} else {
		if strings.TrimSpace(r.Name) == "" {
			return nil, invalid("name", "name is required")
		}
		if r.Code == "" {
			return nil, invalid("code", "code is required")
		}
		dm = DealingMember{Mode: ModeMock, Active: true, DealingMemberCountry: "NGA", CreatedAt: s.now()}
	}
	set := func(dst *string, v string) {
		if v = strings.TrimSpace(v); v != "" {
			*dst = v
		}
	}
	set(&dm.DealingMemberName, r.Name)
	set(&dm.DealingMemberAddress, r.Address)
	set(&dm.DealingMemberCountry, strings.ToUpper(r.Country))
	set(&dm.CSCSMemberCode, r.CSCSMemberCode)
	set(&dm.RequirementDocument, r.RequirementDocument)
	set(&dm.Code, r.Code)
	set(&dm.Mode, r.Mode)
	set(&dm.AuthScheme, r.AuthScheme)
	set(&dm.BaseURL, r.BaseURL)
	set(&dm.CredentialsRef, r.CredentialsRef)
	if r.FeePercent != nil {
		if *r.FeePercent < 0 || *r.FeePercent > 10 {
			return nil, invalid("feePercent", "feePercent must be between 0 and 10")
		}
		dm.FeePercent = *r.FeePercent
	}
	if r.FeeFixed != nil {
		if *r.FeeFixed < 0 {
			return nil, invalid("feeFixed", "feeFixed cannot be negative")
		}
		dm.FeeFixed = *r.FeeFixed
	}
	if r.Active != nil {
		if !*r.Active && dm.ID != 0 {
			var n int64
			s.DB.Model(&Asset{}).Where("dealing_member_id = ? AND status <> ?", dm.ID, AssetSetup).Count(&n)
			if n > 0 {
				return nil, refuse(http.StatusConflict, "%d live asset(s) use this Dealing Member; move them first", n)
			}
		}
		dm.Active = *r.Active
	}
	var dup int64
	s.DB.Model(&DealingMember{}).Where("code = ? AND id <> ?", dm.Code, dm.ID).Count(&dup)
	if dup > 0 {
		return nil, refuse(http.StatusConflict, "another Dealing Member has the code %s", dm.Code)
	}
	dm.UpdatedAt = s.now()
	if err := s.DB.Save(&dm).Error; err != nil {
		return nil, err
	}
	return &dm, nil
}

// ---------------------------------------------------------------- custodians

// CustodianRow is an Approved Asset Custodian with its Public Markets
// integration (nil when not configured for Public Markets).
type CustodianRow struct {
	ID          uint64     `json:"id"`
	Name        string     `json:"name"`
	Country     string     `json:"country"`
	FeePercent  float64    `json:"feePercent"`
	Integration *Custodian `json:"integration"`
	Assets      int64      `json:"assets"`
}

func (s *Service) Custodians() []CustodianRow {
	type approved struct {
		ID                    uint64
		AssetCustodianName    string
		AssetCustodianCountry string
		FeePercent            float64
	}
	var rows []approved
	s.DB.Table("approved_asset_custodians").Select("id, asset_custodian_name, asset_custodian_country, fee_percent").Order("asset_custodian_name").Scan(&rows)
	out := make([]CustodianRow, 0, len(rows))
	for _, a := range rows {
		r := CustodianRow{ID: a.ID, Name: a.AssetCustodianName, Country: a.AssetCustodianCountry, FeePercent: a.FeePercent}
		var c Custodian
		if s.DB.First(&c, "custodian_id = ?", a.ID).Error == nil {
			r.Integration = &c
		}
		s.DB.Model(&Asset{}).Where("custodian_id = ?", a.ID).Count(&r.Assets)
		out = append(out, r)
	}
	return out
}

// ConfigureCustodian sets an Approved Asset Custodian's Public Markets
// integration (created on first save).
func (s *Service) ConfigureCustodian(custodianID uint64, r PartnerRequest) (*Custodian, error) {
	var name string
	s.DB.Table("approved_asset_custodians").Where("id = ?", custodianID).Limit(1).Pluck("asset_custodian_name", &name)
	if name == "" {
		return nil, refuse(http.StatusNotFound, "approved asset custodian not found")
	}
	if err := r.checkIntegration(); err != nil {
		return nil, err
	}
	var c Custodian
	if s.DB.First(&c, "custodian_id = ?", custodianID).Error != nil {
		if r.Code == "" {
			return nil, invalid("code", "code is required")
		}
		c = Custodian{CustodianID: custodianID, Mode: ModeMock, Active: true, CreatedAt: s.now()}
	}
	set := func(dst *string, v string) {
		if v = strings.TrimSpace(v); v != "" {
			*dst = v
		}
	}
	set(&c.Code, r.Code)
	set(&c.NomineeName, r.NomineeName)
	set(&c.Mode, r.Mode)
	set(&c.Transport, r.Transport)
	set(&c.AuthScheme, r.AuthScheme)
	set(&c.BaseURL, r.BaseURL)
	set(&c.CredentialsRef, r.CredentialsRef)
	if r.FeePercent != nil {
		c.FeePercent = decimal.NewFromFloat(*r.FeePercent).String()
	}
	if r.Active != nil {
		if !*r.Active {
			var n int64
			s.DB.Model(&Asset{}).Where("custodian_id = ? AND status <> ?", custodianID, AssetSetup).Count(&n)
			if n > 0 {
				return nil, refuse(http.StatusConflict, "%d live asset(s) are held by this Custodian", n)
			}
		}
		c.Active = *r.Active
	}
	var dup int64
	s.DB.Model(&Custodian{}).Where("code = ? AND custodian_id <> ?", c.Code, custodianID).Count(&dup)
	if dup > 0 {
		return nil, refuse(http.StatusConflict, "another Custodian has the code %s", c.Code)
	}
	c.UpdatedAt = s.now()
	if err := s.DB.Save(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

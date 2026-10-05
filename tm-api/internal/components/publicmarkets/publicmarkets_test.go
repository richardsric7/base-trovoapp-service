package publicmarkets

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	token   = common.HexToAddress("0x00000000000000000000000000000000000000cc")
	issuing = common.HexToAddress("0x00000000000000000000000000000000000000bb")
)

type fakeReader struct{ info TokenInfo }

func (f fakeReader) Token(ctx context.Context, c common.Address) (TokenInfo, error) {
	if c != token {
		return TokenInfo{}, errors.New("no contract at this address")
	}
	return f.info, nil
}

// a Wednesday 11:00 WAT
var now = time.Date(2026, 9, 30, 10, 0, 0, 0, time.UTC)

func testService(t *testing.T) *Service {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.TempDir()+"/pm.db"), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(All()...); err != nil {
		t.Fatal(err)
	}
	db.Exec(`CREATE TABLE approved_asset_custodians (id INTEGER PRIMARY KEY, asset_custodian_name TEXT, asset_custodian_country TEXT, fee_percent REAL)`)
	db.Exec(`INSERT INTO approved_asset_custodians VALUES (1, 'Custodian A Ltd', 'NGA', 0.1), (2, 'Custodian B Ltd', 'NGA', 0)`)
	db.Exec(`CREATE TABLE service_links (id TEXT PRIMARY KEY, short_name TEXT, long_name TEXT, verified INTEGER, suspended INTEGER, inactive INTEGER, rate_limit_per_minute INTEGER)`)
	db.Exec(`INSERT INTO service_links VALUES ('sl-palm', 'palmline', 'Palmline Exchange Ltd', 1, 0, 0, 0), ('sl-new', 'newx', 'New Exchange', 0, 0, 0, 0)`)
	db.Create(&Custodian{CustodianID: 1, Code: "CUSTA", Mode: ModeMock, Active: true})
	db.Create(&DealingMember{ID: 1, DealingMemberName: "Broker Partner 2 Ltd", DealingMemberCountry: "NGA", Code: "BP2", Mode: ModeManual, Active: true})
	s := &Service{DB: db, Chain: fakeReader{TokenInfo{Owner: issuing, Decimals: 4, TotalSupply: new(big.Int)}}, Now: func() time.Time { return now }}
	s.LoadSettings()
	return s
}

func status(err error) int {
	var e *Error
	if errors.As(err, &e) {
		return e.Status
	}
	if err != nil {
		return 500
	}
	return 200
}

func TestAssetSetupToLive(t *testing.T) {
	s := testService(t)
	req := AssetRequest{AssetCode: "mtnn-t", Market: "NGX", AssetType: "EQUITY", ISIN: "NGMTNN000002", InstrumentName: "MTN Nigeria Communications Plc",
		CustodianID: 1, DealingMemberID: 1}
	bad := req
	bad.ISIN = "NG123"
	if _, err := s.CreateAsset(bad, "ops@trovotech.io"); status(err) != 400 {
		t.Fatalf("bad ISIN: %v", err)
	}
	bad = req
	bad.CustodianID = 2 // approved, but not configured for Public Markets
	if _, err := s.CreateAsset(bad, "ops@trovotech.io"); status(err) != 400 {
		t.Fatalf("unconfigured custodian: %v", err)
	}
	a, err := s.CreateAsset(req, "ops@trovotech.io")
	if err != nil {
		t.Fatal(err)
	}
	if a.AssetCode != "MTNN-T" || a.Ticker != "MTNN" || a.Status != AssetSetup || a.OmnibusReference != "POOL-MTNN-01" {
		t.Fatalf("asset %+v", a)
	}
	if _, err := s.CreateAsset(req, "ops@trovotech.io"); status(err) != 409 {
		t.Fatalf("duplicate: %v", err)
	}
	if _, err := s.GoLive(a.ID, "ops"); status(err) != 409 {
		t.Fatalf("went live without a token: %v", err)
	}
	// the token must be owned by the issuing Safe and unminted
	if _, err := s.RegisterContract(context.Background(), a.ID, token.Hex(), common.HexToAddress("0x01").Hex()); status(err) != 400 {
		t.Fatalf("wrong owner accepted: %v", err)
	}
	if a, err = s.RegisterContract(context.Background(), a.ID, token.Hex(), issuing.Hex()); err != nil || a.TokenDecimals != 4 {
		t.Fatalf("register: %v %+v", err, a)
	}
	if _, err := s.SetPrice(a.ID, "-1", "ops"); status(err) != 400 {
		t.Fatal("negative price accepted")
	}
	if _, err := s.SetPrice(a.ID, "221.50", "ops"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.GoLive(a.ID, "ops"); status(err) != 409 || !strings.Contains(err.Error(), "position") {
		t.Fatalf("went live without a position: %v", err)
	}
	if _, err := s.RecordPosition(a.ID, "1000.5", "2026-09-29", "STMT-1", "ops"); status(err) != 400 {
		t.Fatal("fractional units accepted")
	}
	if _, err := s.RecordPosition(a.ID, "1000", "2026-09-29", "STMT-1", "ops"); err != nil {
		t.Fatal(err)
	}
	if a, err = s.GoLive(a.ID, "ops"); err != nil || a.Status != AssetLive {
		t.Fatalf("go live: %v", err)
	}
	if _, err := s.RegisterContract(context.Background(), a.ID, token.Hex(), issuing.Hex()); status(err) != 409 {
		t.Fatal("the token of a live asset changed")
	}
	if _, err := s.UpdateAsset(a.ID, AssetRequest{CustodianID: 1, DealingMemberID: 2}, "ops"); status(err) != 409 {
		t.Fatal("dealing member of a live asset changed")
	}
	if _, err := s.Halt(a.ID, "", "ops"); status(err) != 400 {
		t.Fatal("halted without a reason")
	}
	if a, err = s.Halt(a.ID, "Custodian statement disputed", "ops@trovotech.io"); err != nil || a.Status != AssetHalted {
		t.Fatalf("halt: %v", err)
	}
	j, err := s.Resume(a.ID, "ops@trovotech.io")
	if err != nil || j.Job != "RESUME" || j.Target != "MTNN-T" {
		t.Fatalf("resume: %v %+v", err, j)
	}
	detail, _ := s.AssetDetail("MTNN-T")
	for _, st := range detail.Steps {
		if !st.Done {
			t.Fatalf("step %s not done", st.Label)
		}
	}
}

func seedBatch(s *Service, status string) NetBatch {
	b := NetBatch{ID: "NB-20260930-1100-MTNN-T", AssetID: "pma-mtnn-t", AssetCode: "MTNN-T", SessionDate: "2026-09-30", Side: "BUY", Quantity: "355",
		ReferencePrice: "221", Value: "78455", Status: status, ApprovalsRequired: 2, CreatedAt: now, UpdatedAt: now}
	s.DB.Create(&b)
	s.DB.Create(&Instruction{ID: "DM-1", Kind: InstrDealingOrder, PartnerType: "DEALING_MEMBER", PartnerID: 1, AssetID: b.AssetID, AssetCode: b.AssetCode,
		BatchID: b.ID, Side: "BUY", Quantity: "355", Status: InstrEscalated, CreatedAt: now, UpdatedAt: now})
	s.DB.Create(&Instruction{ID: "CI-1", Kind: InstrCustodianCreation, PartnerType: "CUSTODIAN", PartnerID: 1, AssetID: b.AssetID, AssetCode: b.AssetCode,
		BatchID: b.ID, Side: "BUY", Quantity: "355", Status: InstrAccepted, CreatedAt: now, UpdatedAt: now})
	return b
}

func TestBatchApprovals(t *testing.T) {
	s := testService(t)
	b := seedBatch(s, BatchAwaitingApproval)
	if _, err := s.ApproveBatch(b.ID, "a@trovotech.io"); status(err) != 409 {
		t.Fatalf("approved with no approvers configured: %v", err)
	}
	list := "a@trovotech.io, B@trovotech.io"
	if _, err := s.UpdateSettings(SettingsRequest{NetCreationApprovers: &list}, "admin"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.ApproveBatch(b.ID, "outsider@trovotech.io"); status(err) != 403 {
		t.Fatalf("outsider approved: %v", err)
	}
	v, err := s.ApproveBatch(b.ID, "A@trovotech.io")
	if err != nil || len(v.Approvals) != 1 {
		t.Fatalf("approve: %v", err)
	}
	if _, err := s.ApproveBatch(b.ID, "a@trovotech.io"); status(err) != 409 {
		t.Fatal("approved twice")
	}
	if _, err := s.RejectBatch(b.ID, "too large today", "outsider@trovotech.io"); status(err) != 403 {
		t.Fatal("outsider rejected")
	}
	if r, err := s.RejectBatch(b.ID, "too large today", "b@trovotech.io"); err != nil || r.Status != BatchRejected || !r.UpdatedAt.Equal(now) {
		t.Fatalf("reject: %v", err)
	}
}

func TestEscalationsAndOutcomes(t *testing.T) {
	s := testService(t)
	b := seedBatch(s, BatchReleased)
	in, err := s.RetryInstruction("DM-1", "ops")
	if err != nil || in.Status != InstrPending || in.Attempts != 0 {
		t.Fatalf("retry: %v %+v", err, in)
	}
	if _, err := s.RetryInstruction("DM-1", "ops"); status(err) != 409 {
		t.Fatal("retried a pending instruction")
	}
	// a MANUAL Dealing Member's fill, recorded as the engine expects it
	if _, err := s.RecordOutcome("DM-1", OutcomeRequest{Status: "filled"}, "ops@trovotech.io"); status(err) != 400 {
		t.Fatal("fill without a price")
	}
	ev, err := s.RecordOutcome("DM-1", OutcomeRequest{Status: "filled", ExecutedPrice: "221.10"}, "ops@trovotech.io")
	if err != nil || ev.Kind != "execution" || ev.Source != "MANUAL:ops@trovotech.io" || ev.ProcessedAt != nil {
		t.Fatalf("outcome: %v %+v", err, ev)
	}
	var body map[string]string
	json.Unmarshal([]byte(ev.Payload), &body)
	if body["orderId"] != "DM-1" || body["status"] != "FILLED" || body["executedQuantity"] != "355" {
		t.Fatalf("payload %v", body)
	}
	if _, err := s.RecordOutcome("DM-1", OutcomeRequest{Status: "FILLED", ExecutedPrice: "221.10"}, "ops@trovotech.io"); status(err) != 409 {
		t.Fatal("recorded twice")
	}
	if _, err := s.RecordOutcome("CI-1", OutcomeRequest{Status: "settlement_final"}, "ops"); status(err) != 400 {
		t.Fatal("settlement without the Custodian's reference")
	}
	if ev, err = s.RecordOutcome("CI-1", OutcomeRequest{Status: "settlement_final", CustodianReference: "CSD-9"}, "ops"); err != nil || ev.Kind != "settlement" {
		t.Fatalf("settlement: %v", err)
	}
	// rolling: refused once anything executed
	s.DB.Model(&Instruction{}).Where("id = ?", "DM-1").Update("status", InstrExecuted)
	if _, err := s.RollBatch(b.ID, "market closed", "ops"); status(err) != 409 {
		t.Fatal("rolled an executed batch")
	}
	s.DB.Model(&Instruction{}).Where("id = ?", "DM-1").Update("status", InstrEscalated)
	r, err := s.RollBatch(b.ID, "market closed", "ops")
	if err != nil || r.Status != BatchFailed {
		t.Fatalf("roll: %v", err)
	}
	var open int64
	s.DB.Model(&Instruction{}).Where("batch_id = ? AND status <> ?", b.ID, InstrHandled).Count(&open)
	if open != 0 {
		t.Fatal("instructions of a rolled batch left open")
	}
}

func TestCorporateActions(t *testing.T) {
	s := testService(t)
	s.DB.Create(&Asset{ID: "pma-mtnn-t", AssetCode: "MTNN-T", Ticker: "MTNN", Market: MarketNGX, AssetType: TypeEquity, InstrumentName: "MTN", CustodianID: 1, DealingMemberID: 1, Status: AssetLive})
	if _, err := s.Declare(DeclareRequest{AssetCode: "MTNN-T", EventType: "DIVIDEND", RecordDate: "2026-10-10", PayDate: "2026-10-01", AmountPerUnit: "1"}, "ops"); status(err) != 400 {
		t.Fatal("pay date before record date accepted")
	}
	ev, err := s.Declare(DeclareRequest{AssetCode: "MTNN", EventType: "dividend", RecordDate: "2026-10-10", PayDate: "2026-10-24", AmountPerUnit: "1.00"}, "ops@trovotech.io")
	if err != nil || ev.Kind != "corporate-action" || !strings.HasPrefix(ev.Source, "MANUAL:") {
		t.Fatalf("declare: %v", err)
	}
	ca := CorporateAction{ID: "CA-1", AssetID: "pma-mtnn-t", AssetCode: "MTNN-T", EventType: ActionDividend, RecordDate: "2026-10-10", Status: ActionSnapshotted,
		SnapshotChecksum: "abc", Source: SourceManualCA, CreatedAt: now, UpdatedAt: now}
	s.DB.Create(&ca)
	if _, err := s.Declare(DeclareRequest{AssetCode: "MTNN-T", EventType: "DIVIDEND", RecordDate: "2026-10-10", PayDate: "2026-10-24", AmountPerUnit: "1"}, "ops"); status(err) != 409 {
		t.Fatal("declared the same dividend twice")
	}
	list := "a@trovotech.io,b@trovotech.io"
	s.UpdateSettings(SettingsRequest{DividendApprovers: &list}, "admin")
	detail, err := s.ApproveAction("CA-1", "a@trovotech.io")
	if err != nil || len(detail.Approvals) != 1 || detail.Approvals[0].Checksum != "abc" {
		t.Fatalf("approve: %v", err)
	}
	if _, err := s.ApproveAction("CA-1", "c@trovotech.io"); status(err) != 403 {
		t.Fatal("non-approver approved")
	}
	s.DB.Model(&CorporateAction{}).Where("id = ?", "CA-1").Update("status", ActionPaying)
	if _, err := s.CancelAction("CA-1", "wrong amount", "ops"); status(err) != 409 {
		t.Fatal("cancelled while paying")
	}
}

func TestExchanges(t *testing.T) {
	s := testService(t)
	tiers := `{"Tier 1": 1200, "Tier 2": 3000}`
	if _, err := s.UpdateSettings(SettingsRequest{RateLimitTiers: &tiers}, "admin"); err != nil {
		t.Fatal(err)
	}
	if _, _, err := s.Onboard(ExchangeRequest{ServiceLinkID: "sl-new", CallbackURL: "https://x.io/h", FundingAddress: issuing.Hex()}, "ops"); status(err) != 409 {
		t.Fatal("onboarded an unverified service link")
	}
	if _, _, err := s.Onboard(ExchangeRequest{ServiceLinkID: "sl-palm", CallbackURL: "http://x.io/h", Environment: "production", FundingAddress: issuing.Hex()}, "ops"); status(err) != 400 {
		t.Fatal("http callback in production accepted")
	}
	row, secret, err := s.Onboard(ExchangeRequest{ServiceLinkID: "sl-palm", CallbackURL: "https://palmline.ng/hooks", FundingAddress: issuing.Hex(), RateLimitTier: "Tier 1"}, "ops")
	if err != nil || !strings.HasPrefix(secret, "pmsk_") || row.RateLimitPerMinute != 1200 || row.Name != "Palmline Exchange Ltd" {
		t.Fatalf("onboard: %v %+v", err, row)
	}
	if len(s.ExchangeCandidates()) != 0 {
		t.Fatal("an onboarded link is still a candidate")
	}
	newSecret, until, err := s.RotateSecret("sl-palm", "ops")
	var p ExchangePartner
	s.DB.First(&p, "service_link_id = ?", "sl-palm")
	if err != nil || newSecret == secret || p.SigningSecret != newSecret || p.PreviousSigningSecret != secret || !until.Equal(now.Add(24*time.Hour)) {
		t.Fatalf("rotate: %v", err)
	}
	// a tier's new limit reaches the service link
	tiers = `{"Tier 1": 600}`
	s.UpdateSettings(SettingsRequest{RateLimitTiers: &tiers}, "admin")
	var limit int
	s.DB.Table("service_links").Where("id = ?", "sl-palm").Pluck("rate_limit_per_minute", &limit)
	if limit != 600 {
		t.Fatalf("rate limit %d", limit)
	}
	s.DB.Model(&ExchangePartner{}).Where("service_link_id = ?", "sl-palm").Update("balance", "5000")
	if _, err := s.RequestWithdrawal("sl-palm", "6000", "ops"); status(err) != 400 {
		t.Fatal("withdrew more than the balance")
	}
	j, err := s.RequestWithdrawal("sl-palm", "2500.50", "ops")
	if err != nil || j.Target != "sl-palm:2500.5" {
		t.Fatalf("withdrawal: %v %+v", err, j)
	}
	if _, err := s.RequestWithdrawal("sl-palm", "1", "ops"); status(err) != 409 {
		t.Fatal("two withdrawals waiting")
	}
	s.DB.Create(&WebhookDelivery{ID: "evt_1", ServiceLinkID: "sl-palm", Event: "dividend.paid", Status: DeliveryDeadLetter, Attempts: 9, CreatedAt: now})
	if w, err := s.ReplayWebhook("evt_1", "ops"); err != nil || w.Status != DeliveryPending || w.Attempts != 0 {
		t.Fatalf("replay: %v", err)
	}
	if _, err := s.SetExchangeStatus("sl-palm", false, "ops"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.SetExchangeStatus("sl-palm", false, "ops"); status(err) != 409 {
		t.Fatal("suspended twice")
	}
}

func TestWalletMasking(t *testing.T) {
	s := testService(t)
	s.DB.Create(&PartnerWallet{ID: "wlt_1", ServiceLinkID: "sl-palm", ExternalUserRef: "plm_1", LegalName: "Adaeze Okafor", TaxIdentifier: "22339871",
		ResidencyCountry: "NGA", Nationality: "NGA", NDPAConsent: true, Status: "active", CreatedAt: now})
	s.DB.Create(&PartnerWallet{ID: "wlt_2", ServiceLinkID: "sl-palm", ExternalUserRef: "plm_2", Status: "rejected", RejectionReason: "taxIdentifier is required", CreatedAt: now})
	rows, total, st := s.Wallets("", "", "", false, Page{})
	if total != 2 || st.Provisioned != 1 || st.Rejected30d != 1 || st.ConsentPercent != "100" {
		t.Fatalf("stats %+v total %d", st, total)
	}
	for _, r := range rows {
		if r.ID == "wlt_1" && (r.LegalName != "A***** O*****" || r.TaxIdentifier != "2233**71") {
			t.Fatalf("masked %q %q", r.LegalName, r.TaxIdentifier)
		}
	}
	rows, _, _ = s.Wallets("", "active", "", true, Page{})
	if len(rows) != 1 || rows[0].LegalName != "Adaeze Okafor" {
		t.Fatalf("revealed %+v", rows)
	}
}

func TestSettingsValidation(t *testing.T) {
	s := testService(t)
	hour := 25
	if _, err := s.UpdateSettings(SettingsRequest{ReconciliationHour: &hour}, "admin"); status(err) != 400 {
		t.Fatal("hour 25 accepted")
	}
	if _, err := s.UpdateSettings(SettingsRequest{NGXOpen: "15:00"}, "admin"); status(err) != 400 {
		t.Fatal("opening after the close accepted")
	}
	one := "a@trovotech.io"
	if _, err := s.UpdateSettings(SettingsRequest{NetCreationApprovers: &one}, "admin"); status(err) != 400 {
		t.Fatal("fewer approvers than approvals required")
	}
	days := "2026-10-01, 2026-12-25"
	st, err := s.UpdateSettings(SettingsRequest{TradeFeePercent: "0.30", MarketHolidays: &days}, "admin")
	if err != nil || st.TradeFeePercent != "0.30" || st.MarketHolidays != "2026-10-01,2026-12-25" {
		t.Fatalf("save: %v", err)
	}
	if marketOpen(*st, MarketNGX, time.Date(2026, 10, 1, 11, 0, 0, 0, Lagos)) {
		t.Fatal("open on a holiday")
	}
	if !marketOpen(*st, MarketNGX, now) {
		t.Fatal("closed in session")
	}
}

func TestPartnersAndCustodians(t *testing.T) {
	s := testService(t)
	if _, err := s.SaveDealingMember(0, PartnerRequest{Name: "Broker 3", Code: "BP3", Mode: "REST", BaseURL: "https://b3.ng"}); status(err) != 400 {
		t.Fatal("REST partner without credentials")
	}
	if _, err := s.SaveDealingMember(0, PartnerRequest{Name: "Broker 3", Code: "BP3", Mode: "REST", BaseURL: "https://b3.ng", CredentialsRef: "s3cret"}); status(err) != 400 {
		t.Fatal("a secret stored as the credentials reference")
	}
	dm, err := s.SaveDealingMember(0, PartnerRequest{Name: "Broker 3", Code: "bp3", Mode: "REST", BaseURL: "https://b3.ng", CredentialsRef: "env:BROKER3"})
	if err != nil || dm.Code != "BP3" {
		t.Fatalf("save: %v", err)
	}
	if _, err := s.SaveDealingMember(0, PartnerRequest{Name: "Other", Code: "BP3"}); status(err) != 409 {
		t.Fatal("duplicate code")
	}
	c, err := s.ConfigureCustodian(2, PartnerRequest{Code: "CUSTB", NomineeName: "Custodian B Nominees", Mode: "MANUAL"})
	if err != nil || c.CustodianID != 2 || c.Mode != ModeManual {
		t.Fatalf("configure: %v", err)
	}
	rows := s.Custodians()
	if len(rows) != 2 || rows[1].Integration == nil {
		t.Fatalf("custodians %+v", rows)
	}
	if _, err := s.RequestPositionFeed("CUSTB", "ops"); status(err) != 409 {
		t.Fatal("feed requested from a MANUAL custodian")
	}
	if j, err := s.RequestPositionFeed("custa", "ops"); err != nil || j.Target != "CUSTA" {
		t.Fatalf("feed: %v", err)
	}
}

func TestPermissions(t *testing.T) {
	s := testService(t)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("auth_type", "trovo_admin")
		c.Set("trovo_admin_email", "viewer@trovotech.io")
	})
	h := Handlers{S: s, Can: func(c *gin.Context, p string) bool { return false }}
	r.GET("/assets", h.Assets)
	r.POST("/assets", h.CreateAsset)
	r.GET("/wallets", h.Wallets)
	r.PUT("/settings", h.UpdateSettings)
	call := func(method, path string, body interface{}) (int, map[string]interface{}) {
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(body)
		req := httptest.NewRequest(method, path, &buf)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		var out map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &out)
		return w.Code, out
	}
	if code, _ := call("GET", "/assets", nil); code != http.StatusOK {
		t.Fatalf("read: %d", code)
	}
	if code, _ := call("POST", "/assets", AssetRequest{AssetCode: "X"}); code != http.StatusForbidden {
		t.Fatalf("write without MANAGE_PUBLIC_MARKETS: %d", code)
	}
	if code, _ := call("PUT", "/settings", SettingsRequest{TradeFeePercent: "1"}); code != http.StatusForbidden {
		t.Fatalf("settings without MANAGE_SETTINGS: %d", code)
	}
	_, out := call("GET", "/wallets", nil)
	if data, _ := out["data"].(map[string]interface{}); data == nil || data["personalDataVisible"] != false {
		t.Fatalf("wallets %v", out)
	}
}

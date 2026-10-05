package proceedpayouts

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"admin-panel-dashboard/internal/cache"
	stakeholderModels "admin-panel-dashboard/internal/components/stakeholder/models"
	stakeholderServices "admin-panel-dashboard/internal/components/stakeholder/services"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const cngn = "0x00000000000000000000000000000000000000C1"

func testService(t *testing.T) *Service {
	t.Helper()
	open := func(name string) *gorm.DB {
		db, err := gorm.Open(sqlite.Open("file:"+url.QueryEscape(t.Name()+name)+"?mode=memory&cache=shared"), &gorm.Config{Logger: logger.Discard})
		if err != nil {
			t.Fatal(err)
		}
		return db
	}
	db, admin := open("wallet"), open("admin")
	if err := db.AutoMigrate(&ProceedPayout{}, &PayoutItem{}, &PayoutBatch{}, &EngineState{}, &ServiceFee{}, &tokenizationCurrency{}, &countryConfig{}, &FeeCollection{}); err != nil {
		t.Fatal(err)
	}
	db.Exec(`CREATE TABLE proceed_payout_approvals (id INTEGER PRIMARY KEY AUTOINCREMENT, created_at DATETIME, proceed_payout_id INTEGER, admin_email TEXT, schedule_checksum TEXT, UNIQUE(proceed_payout_id, admin_email))`)
	db.Exec(`CREATE UNIQUE INDEX idx_batch ON proceed_payouts(batch)`)
	db.Exec(`CREATE TABLE tokenized_assets (id TEXT PRIMARY KEY, asset_code TEXT, asset_name TEXT, asset_country_location TEXT, proceed_payout_currency TEXT, contract_address TEXT)`)
	db.Exec(`INSERT INTO tokenized_assets VALUES ('42', 'FARM', 'Farm Fund', 'NG', 'CNGN', '0x00000000000000000000000000000000000000A5')`)
	db.Create(&tokenizationCurrency{AssetCode: "CNGN", ContractAddress: cngn})
	db.Create(&ServiceFee{ID: FeeServiceID, FeeWalletSecretKey: "0x00000000000000000000000000000000000000F1", FeePercent: 5, Remarks: "50"})
	if err := admin.AutoMigrate(&stakeholderModels.Distribution{}); err != nil {
		t.Fatal(err)
	}
	admin.Create(&stakeholderModels.Distribution{ID: "dist-1", AssetID: "42", ProposedByOrgID: "o", ProposedByMemberID: "m", TrusteeOrgID: "t",
		Amount: decimal.NewFromInt(1000), Currency: "NGN", Source: "rent", Status: stakeholderModels.DistributionStatusAuthorized})
	return &Service{DB: db, AdminDB: admin}
}

type client struct {
	t *testing.T
	r *gin.Engine
}

// as calls the API as a Trovo admin (or "member:" an organization member).
func (c client) as(who, method, path string, body interface{}) (int, map[string]interface{}) {
	c.t.Helper()
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Test-Admin", who)
	w := httptest.NewRecorder()
	c.r.ServeHTTP(w, req)
	var out map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &out)
	return w.Code, out
}

func router(s *Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		who := c.GetHeader("X-Test-Admin")
		if len(who) > 7 && who[:7] == "member:" {
			c.Set("auth_type", "organization_member")
			return
		}
		c.Set("auth_type", "trovo_admin")
		c.Set("trovo_admin_email", who)
	})
	h := Handlers{S: s}
	g := r.Group("/proceed-payouts")
	g.GET("", h.List)
	g.GET("/engine", h.Engine)
	g.POST("/engine/halt", h.Halt)
	g.POST("/engine/unhalt", h.Unhalt)
	g.POST("/engine/sweep", h.Sweep)
	g.GET("/fee-config", h.FeeConfig)
	g.PUT("/fee-config", h.SetFeeConfig)
	g.GET("/reports/payouts", h.PayoutsReport)
	g.GET("/reports/fees", h.FeesReport)
	g.GET("/:id", h.Detail)
	g.GET("/:id/items", h.Items)
	g.POST("/:id/prepare", h.Prepare)
	g.PUT("/:id/fee", h.SetFee)
	g.POST("/:id/approve", h.Approve)
	g.POST("/:id/reject", h.Reject)
	g.POST("/:id/confirm-funding", h.ConfirmFunding)
	g.POST("/:id/pause", h.Pause)
	g.POST("/:id/resume", h.Resume)
	g.POST("/:id/cancel", h.Cancel)
	g.POST("/:id/retry-failed", h.RetryFailed)
	g.POST("/:id/items/:itemId/exclude", h.ExcludeItem)
	g.POST("/:id/items/:itemId/include", h.IncludeItem)
	g.POST("/:id/items/:itemId/mark-paid", h.MarkPaid)
	return r
}

func status(t *testing.T, s *Service, id uint64) string {
	t.Helper()
	p, err := s.load(id)
	if err != nil {
		t.Fatal(err)
	}
	return p.Status
}

// lock does what payout-engine does when it locks a schedule.
func lock(t *testing.T, s *Service, id uint64, checksum string) {
	t.Helper()
	s.DB.Where("proceed_payout_id = ?", id).Delete(&PayoutItem{})
	items := []PayoutItem{
		{ID: "h1", ProceedPayoutID: id, BeneficiaryAddress: "0x01", BalanceUnits: "600", AmountUnits: "570000000", AmountToReceive: 570, Status: ItemPending, Kind: KindHolder, Username: "alice"},
		{ID: "h2", ProceedPayoutID: id, BeneficiaryAddress: "0x02", BalanceUnits: "300", AmountUnits: "285000000", AmountToReceive: 285, Status: ItemPending, Kind: KindHolder},
		{ID: "h3", ProceedPayoutID: id, BeneficiaryAddress: "0x03", BalanceUnits: "100", AmountUnits: "95000000", AmountToReceive: 95, Status: ItemExcluded, Reason: "the market-making wallet", Kind: KindHolder},
		{ID: "fee", ProceedPayoutID: id, BeneficiaryAddress: "0xF1", AmountUnits: "50000000", Status: ItemPending, Kind: KindFee},
	}
	if err := s.DB.Create(&items).Error; err != nil {
		t.Fatal(err)
	}
	s.DB.Model(&ProceedPayout{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status": StatusLocked, "schedule_checksum": checksum, "payment_schedule_ready": 1, "payout_decimals": 6, "token_decimals": 0,
		"fee_units": "50000000", "vat_units": "3750000", "holder_payable": "946250000", "paid_units": "0",
	})
}

func TestProceedPayoutLifecycle(t *testing.T) {
	s := testService(t)
	var sub *redis.PubSub
	if addr := os.Getenv("TEST_REDIS_ADDR"); addr != "" {
		rdb := redis.NewClient(&redis.Options{Addr: addr})
		if rdb.Ping(context.Background()).Err() == nil {
			s.Cache = &cache.RedisCache{Enabled: true, Client: rdb, Context: context.Background()}
			sub = rdb.Subscribe(context.Background(), CommandsChannel)
			defer sub.Close()
			sub.Receive(context.Background())
		}
	}
	api := client{t, router(s)}
	ctx := context.Background()

	// the trustee's authorization registers the payout, once
	dc := DistributionClient{Service: s}
	resp, err := dc.Register(ctx, stakeholderServicesRegistration("dist-1"))
	if err != nil || resp == nil || resp.Status != "registered" || resp.Currency != "CNGN" {
		t.Fatalf("register: %+v %v", resp, err)
	}
	if _, err := dc.Register(ctx, stakeholderServicesRegistration("dist-1")); err != nil {
		t.Fatal(err)
	}
	var payouts []ProceedPayout
	s.DB.Find(&payouts)
	if len(payouts) != 1 {
		t.Fatalf("registered %d payouts", len(payouts))
	}
	p := payouts[0]
	if p.Batch != "FARM-dist-1" || p.PayoutContractAddress != cngn || p.FeeType != FeePercent || p.FeeValue != "5" || p.FeeCap != "50" || p.ApprovalsRequired != 2 {
		t.Fatalf("registered payout: %+v", p)
	}
	path := "/proceed-payouts/" + itoa(p.ID)

	if code, _ := api.as("member:x", http.MethodPost, path+"/prepare", nil); code != http.StatusForbidden {
		t.Fatalf("a member prepared: %d", code)
	}
	if code, _ := api.as("ann@trovo", http.MethodPost, path+"/approve", nil); code != http.StatusConflict {
		t.Fatalf("approved before locking: %d", code)
	}
	if code, out := api.as("ann@trovo", http.MethodPost, path+"/prepare", nil); code != http.StatusOK || status(t, s, p.ID) != StatusPrepareRequested {
		t.Fatalf("prepare: %d %v", code, out)
	}
	if sub != nil {
		select {
		case m := <-sub.Channel():
			if m.Payload == "" {
				t.Fatal("empty wake-up")
			}
		case <-time.After(2 * time.Second):
			t.Fatal("the engine was not woken")
		}
	}
	var d stakeholderModels.Distribution
	s.AdminDB.First(&d, "id = ?", "dist-1")
	if d.Status != stakeholderModels.DistributionStatusProcessing {
		t.Fatalf("distribution: %s", d.Status)
	}

	// approvals: distinct admins, not the preparer
	lock(t, s, p.ID, "sum-1")
	if code, _ := api.as("ann@trovo", http.MethodPost, path+"/approve", nil); code != http.StatusForbidden {
		t.Fatalf("the preparer approved: %d", code)
	}
	if code, out := api.as("bob@trovo", http.MethodPost, path+"/approve", nil); code != http.StatusOK || status(t, s, p.ID) != StatusLocked {
		t.Fatalf("first approval: %d %v", code, out)
	}
	if code, _ := api.as("BOB@trovo", http.MethodPost, path+"/approve", nil); code != http.StatusConflict {
		t.Fatalf("approved twice: %d", code)
	}
	if code, _ := api.as("cy@trovo", http.MethodPost, path+"/approve", nil); code != http.StatusOK || status(t, s, p.ID) != StatusApproved {
		t.Fatalf("second approval: %d %s", code, status(t, s, p.ID))
	}

	// changing the fee re-prepares the schedule; its setter cannot approve
	if code, _ := api.as("dee@trovo", http.MethodPut, path+"/fee", FeeRequest{FeeType: "FIXED", FeeValue: "2000"}); code != http.StatusBadRequest {
		t.Fatalf("a fee above the payout: %d", code)
	}
	if code, out := api.as("dee@trovo", http.MethodPut, path+"/fee", FeeRequest{FeeType: "percent", FeeValue: "2", FeeCap: "10"}); code != http.StatusOK {
		t.Fatalf("set fee: %d %v", code, out)
	}
	p2, _ := s.load(p.ID)
	var n int64
	s.DB.Model(&PayoutApproval{}).Where("proceed_payout_id = ?", p.ID).Count(&n)
	if p2.Status != StatusPrepareRequested || p2.FeeType != FeePercent || p2.FeeValue != "2" || p2.FeeCap != "10" || p2.FeeSetBy != "dee@trovo" || n != 0 {
		t.Fatalf("after the fee change: %+v, %d approvals", p2, n)
	}
	lock(t, s, p.ID, "sum-2")
	if code, _ := api.as("dee@trovo", http.MethodPost, path+"/approve", nil); code != http.StatusForbidden {
		t.Fatalf("the fee setter approved: %d", code)
	}
	api.as("bob@trovo", http.MethodPost, path+"/approve", nil)
	api.as("cy@trovo", http.MethodPost, path+"/approve", nil)
	if status(t, s, p.ID) != StatusApproved {
		t.Fatalf("not approved: %s", status(t, s, p.ID))
	}

	// item controls
	if code, _ := api.as("bob@trovo", http.MethodPost, path+"/items/h3/include", nil); code != http.StatusBadRequest {
		t.Fatalf("included a platform wallet: %d", code)
	}
	if code, _ := api.as("bob@trovo", http.MethodPost, path+"/items/fee/exclude", ReasonRequest{Reason: "x"}); code != http.StatusBadRequest {
		t.Fatalf("excluded the fee line: %d", code)
	}
	if code, out := api.as("bob@trovo", http.MethodPost, path+"/items/h2/exclude", ReasonRequest{Reason: "sanctions check"}); code != http.StatusOK {
		t.Fatalf("exclude: %d %v", code, out)
	}
	if p3, _ := s.load(p.ID); p3.ExcludedCount != 2 || p3.Status != StatusApproved {
		t.Fatalf("after exclusion: %+v", p3)
	}
	if code, _ := api.as("cy@trovo", http.MethodPost, path+"/items/h2/include", nil); code != http.StatusOK {
		t.Fatalf("include: %d", code)
	}

	// funding, pause, mark paid, resume
	if code, _ := api.as("bob@trovo", http.MethodPost, path+"/confirm-funding", nil); code != http.StatusOK || status(t, s, p.ID) != StatusFundingCheckRequested {
		t.Fatalf("confirm funding: %d", code)
	}
	s.DB.Model(&ProceedPayout{}).Where("id = ?", p.ID).Update("status", StatusPaying) // the engine found it funded
	s.DB.Model(&PayoutItem{}).Where("id = ?", "h1").Update("status", ItemFailed)
	if code, _ := api.as("bob@trovo", http.MethodPost, path+"/items/h1/mark-paid", ReferenceRequest{Reference: "bank"}); code != http.StatusConflict {
		t.Fatalf("marked paid while paying: %d", code)
	}
	if code, _ := api.as("bob@trovo", http.MethodPost, path+"/pause", ReasonRequest{Reason: "checking"}); code != http.StatusOK || status(t, s, p.ID) != StatusPaused {
		t.Fatalf("pause: %d", code)
	}
	if code, _ := api.as("bob@trovo", http.MethodPost, path+"/items/h1/mark-paid", ReferenceRequest{Reference: "paid by transfer 0xabc"}); code != http.StatusOK {
		t.Fatalf("mark paid: %d", code)
	}
	if p4, _ := s.load(p.ID); p4.PaidCount != 1 || p4.PaidUnits != "570000000" {
		t.Fatalf("after mark paid: %+v", p4)
	}
	if code, _ := api.as("bob@trovo", http.MethodPost, path+"/resume", nil); code != http.StatusOK || status(t, s, p.ID) != StatusFundingCheckRequested {
		t.Fatalf("resume: %d", code)
	}

	// completed with failures: retry, then settle the last one by hand
	s.DB.Model(&PayoutItem{}).Where("id = ?", "h2").Update("status", ItemFailed)
	s.DB.Model(&PayoutItem{}).Where("id = ?", "fee").Update("status", ItemPaid)
	s.DB.Model(&ProceedPayout{}).Where("id = ?", p.ID).Update("status", StatusCompletedWithFailures)
	if code, _ := api.as("bob@trovo", http.MethodPost, path+"/retry-failed", nil); code != http.StatusOK || status(t, s, p.ID) != StatusFundingCheckRequested {
		t.Fatalf("retry: %d", code)
	}
	var h2 PayoutItem
	s.DB.First(&h2, "id = ?", "h2")
	if h2.Status != ItemPending {
		t.Fatalf("retried item: %s", h2.Status)
	}
	s.DB.Model(&PayoutItem{}).Where("id = ?", "h2").Update("status", ItemFailed)
	s.DB.Model(&ProceedPayout{}).Where("id = ?", p.ID).Update("status", StatusCompletedWithFailures)
	if code, _ := api.as("cy@trovo", http.MethodPost, path+"/items/h2/mark-paid", ReferenceRequest{Reference: "cheque"}); code != http.StatusOK || status(t, s, p.ID) != StatusCompleted {
		t.Fatalf("settling the last failure: %d %s", code, status(t, s, p.ID))
	}
	s.AdminDB.First(&d, "id = ?", "dist-1")
	if d.Status != stakeholderModels.DistributionStatusCompleted {
		t.Fatalf("distribution after completion: %s", d.Status)
	}
	if code, _ := api.as("bob@trovo", http.MethodPost, path+"/cancel", ReasonRequest{Reason: "x"}); code != http.StatusConflict {
		t.Fatalf("cancelled a completed payout: %d", code)
	}

	// the trustee's view and the admin reads
	view, err := dc.Get(ctx, "dist-1")
	if err != nil || !view.PayoutCompleted || len(view.Payouts) != 3 || view.Payouts[0].Amount != "570" {
		t.Fatalf("trustee view: %+v %v", view, err)
	}
	if code, out := api.as("bob@trovo", http.MethodGet, path, nil); code != http.StatusOK || out["data"].(map[string]interface{})["approvals"].(float64) != 2 {
		t.Fatalf("detail: %d %v", code, out)
	}
	if code, out := api.as("bob@trovo", http.MethodGet, path+"/items?kind=HOLDER&status=PAID", nil); code != http.StatusOK || out["data"].(map[string]interface{})["total"].(float64) != 2 {
		t.Fatalf("items: %d %v", code, out)
	}
	if code, out := api.as("bob@trovo", http.MethodGet, "/proceed-payouts?asset=FARM", nil); code != http.StatusOK || out["data"].(map[string]interface{})["total"].(float64) != 1 {
		t.Fatalf("list: %d %v", code, out)
	}
}

func TestPayoutCancelAndEngineControls(t *testing.T) {
	s := testService(t)
	api := client{t, router(s)}
	p, err := s.Register(context.Background(), Registration{DistributionID: "dist-1", TokenizedAssetID: "42", Amount: "1000", Currency: "NGN"})
	if err != nil {
		t.Fatal(err)
	}
	path := "/proceed-payouts/" + itoa(p.ID)
	s.DB.Create(&PayoutBatch{ProceedPayoutID: p.ID, Status: batchSubmitted})
	if code, _ := api.as("ann@trovo", http.MethodPost, path+"/prepare", nil); code != http.StatusConflict {
		t.Fatalf("prepared a schedule that started paying: %d", code)
	}
	if code, _ := api.as("ann@trovo", http.MethodPost, path+"/cancel", ReasonRequest{Reason: "wrong amount"}); code != http.StatusConflict {
		t.Fatalf("cancelled with a batch in flight: %d", code)
	}
	s.DB.Model(&PayoutBatch{}).Where("proceed_payout_id = ?", p.ID).Update("status", "MINED")
	if code, _ := api.as("ann@trovo", http.MethodPost, path+"/cancel", nil); code != http.StatusBadRequest {
		t.Fatalf("cancelled without a reason: %d", code)
	}
	if code, _ := api.as("ann@trovo", http.MethodPost, path+"/cancel", ReasonRequest{Reason: "wrong amount"}); code != http.StatusOK || status(t, s, p.ID) != StatusCancelled {
		t.Fatalf("cancel: %d", code)
	}
	var d stakeholderModels.Distribution
	s.AdminDB.First(&d, "id = ?", "dist-1")
	if d.Status != stakeholderModels.DistributionStatusFailed {
		t.Fatalf("distribution after cancel: %s", d.Status)
	}

	// kill switch and sweep
	if code, _ := api.as("ann@trovo", http.MethodPost, "/proceed-payouts/engine/halt", ReasonRequest{Reason: "incident"}); code != http.StatusOK {
		t.Fatalf("halt: %d", code)
	}
	st, _ := s.Engine()
	if !st.Halted || st.HaltedBy != "ann@trovo" || st.Online {
		t.Fatalf("halted state: %+v", st)
	}
	api.as("ann@trovo", http.MethodPost, "/proceed-payouts/engine/unhalt", nil)
	if st, _ := s.Engine(); st.Halted {
		t.Fatal("still halted")
	}
	if code, _ := api.as("ann@trovo", http.MethodPost, "/proceed-payouts/engine/sweep", SweepRequest{Token: "nope"}); code != http.StatusBadRequest {
		t.Fatalf("swept a non-address: %d", code)
	}
	if code, _ := api.as("ann@trovo", http.MethodPost, "/proceed-payouts/engine/sweep", SweepRequest{Token: cngn}); code != http.StatusOK {
		t.Fatalf("sweep: %d", code)
	}
	if code, _ := api.as("bob@trovo", http.MethodPost, "/proceed-payouts/engine/sweep", SweepRequest{Token: cngn}); code != http.StatusConflict {
		t.Fatalf("a second sweep while one is pending: %d", code)
	}
}

func TestPayoutFeeConfigAndReports(t *testing.T) {
	s := testService(t)
	api := client{t, router(s)}
	s.DB.Where("id = ?", FeeServiceID).Delete(&ServiceFee{})
	if cfg, _ := s.FeeConfig(); cfg.FeeType != FeeFixed || cfg.FeeValue != "0" || cfg.FeeCap != "0" || cfg.FeeWallet != "" {
		t.Fatalf("defaults without a configuration: %+v", cfg)
	}
	if code, _ := api.as("ann@trovo", http.MethodPut, "/proceed-payouts/fee-config", FeeConfigRequest{FeeWallet: "x"}); code != http.StatusBadRequest {
		t.Fatalf("a non-address wallet: %d", code)
	}
	if code, out := api.as("ann@trovo", http.MethodPut, "/proceed-payouts/fee-config", FeeConfigRequest{FeeWallet: "0x00000000000000000000000000000000000000f2", FeeType: "PERCENT", FeeValue: "1.5", FeeCap: "25"}); code != http.StatusOK {
		t.Fatalf("set config: %d %v", code, out)
	}
	cfg, _ := s.FeeConfig()
	if cfg.FeeWallet != "0x00000000000000000000000000000000000000F2" || cfg.FeeType != FeePercent || cfg.FeeValue != "1.5" || cfg.FeeCap != "25" || cfg.LastUpdatedBy != "ann@trovo" {
		t.Fatalf("saved config: %+v", cfg)
	}
	api.as("ann@trovo", http.MethodPut, "/proceed-payouts/fee-config", FeeConfigRequest{FeeWallet: "0x00000000000000000000000000000000000000f2", FeeType: "FIXED", FeeValue: "3"})
	if cfg, _ := s.FeeConfig(); cfg.FeeType != FeeFixed || cfg.FeeValue != "3" || cfg.FeeCap != "0" {
		t.Fatalf("fixed config: %+v", cfg)
	}

	// reports
	p, _ := s.Register(context.Background(), Registration{DistributionID: "dist-1", TokenizedAssetID: "42", Amount: "1000", Currency: "NGN"})
	lock(t, s, p.ID, "sum")
	s.DB.Model(&ProceedPayout{}).Where("id = ?", p.ID).Updates(map[string]interface{}{"status": StatusCompleted, "paid_units": "855000000", "paid_count": 2})
	hash := "0xhash"
	s.DB.Create(&[]FeeCollection{
		{ID: "payout-fee", CreatedAt: time.Now(), FromUsername: "FARM", FromWalletAlias: "payout FARM-dist-1", FeeType: FeeTypePayout, Amount: 50, AssetCode: "CNGN", TransactionHash: &hash},
		{ID: "payout-vat", CreatedAt: time.Now(), FromUsername: "FARM", FromWalletAlias: "payout FARM-dist-1", FeeType: FeeTypePayoutVat, Amount: 3.75, AssetCode: "CNGN", TransactionHash: &hash},
		{ID: "swap-fee", CreatedAt: time.Now(), FromUsername: "alice", FeeType: "SWAP", Amount: 9, AssetCode: "CNGN"},
	})
	code, out := api.as("ann@trovo", http.MethodGet, "/proceed-payouts/reports/payouts?asset=42", nil)
	totals := out["data"].(map[string]interface{})["totals"].([]interface{})
	tot := totals[0].(map[string]interface{})
	if code != http.StatusOK || len(totals) != 1 || tot["currency"] != "CNGN" || tot["fees"] != "50" || tot["vat"] != "3.75" || tot["paidToHolders"] != "855" || tot["holderPayable"] != "946.25" {
		t.Fatalf("payouts report: %d %v", code, out)
	}
	code, out = api.as("ann@trovo", http.MethodGet, "/proceed-payouts/reports/fees?asset=FARM&from=2000-01-01", nil)
	data := out["data"].(map[string]interface{})
	if code != http.StatusOK || len(data["collections"].([]interface{})) != 2 || data["totalFees"].(map[string]interface{})["CNGN"].(float64) != 50 || data["totalVat"].(map[string]interface{})["CNGN"].(float64) != 3.75 {
		t.Fatalf("fees report: %d %v", code, out)
	}
	if code, _ := api.as("ann@trovo", http.MethodGet, "/proceed-payouts/reports/fees?from=yesterday", nil); code != http.StatusBadRequest {
		t.Fatalf("a bad date: %d", code)
	}
}

func itoa(n uint64) string { return decimal.NewFromInt(int64(n)).String() }

func stakeholderServicesRegistration(distributionID string) stakeholderServices.DistributionPayoutRegistration {
	return stakeholderServices.DistributionPayoutRegistration{DistributionID: distributionID, TokenizedAssetID: "42", Amount: "1000", Currency: "NGN"}
}

package publicmarkets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	pmModels "trovo-wallet-api/internal/components/publicmarkets/models"
	"trovo-wallet-api/internal/components/publicmarkets/partners"
	pmsvc "trovo-wallet-api/internal/components/publicmarkets/services"
	slModels "trovo-wallet-api/internal/components/servicelinks/models"
	appdb "trovo-wallet-api/internal/db"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type server struct {
	t      *testing.T
	router *gin.Engine
	db     *gorm.DB
}

func newServer(t *testing.T) *server {
	t.Helper()
	gin.SetMode(gin.TestMode)
	db := testDB(t)
	appdb.MigrateDB(db)
	gc := &sharedconfig.GlobalConfig{DB: db}
	e := pmsvc.NewEngine(gc, nil) // no chain: provisioning answers 503
	r := gin.New()
	Register(r, gc, e)
	db.Create(&slModels.ServiceLink{ID: "sl-palmline", OwnerUsername: "palmline", Address: "0x01", ApiKey: "pk_test_palmline", ShortName: "palmline", Verified: 1})
	db.Create(&slModels.ServiceLink{ID: "sl-other", OwnerUsername: "otherx", Address: "0x02", ApiKey: "pk_test_other", ShortName: "otherx", Verified: 1})
	db.Create(&pmModels.ExchangePartner{ServiceLinkID: "sl-palmline", Status: "active", SigningSecret: "s3cret", Balance: "1500"})
	db.Create(&pmModels.Asset{ID: "asset-mtnn", AssetCode: "MTNN-T", Ticker: "MTNN", Market: pmModels.MarketNGX, AssetType: pmModels.TypeEquity,
		ISIN: "NGMTNN000002", CustodianID: 1, Status: pmModels.AssetLive, LastPrice: "221", TokenDecimals: 4})
	db.Create(&pmModels.Custodian{CustodianID: 1, Code: "CUSTA", Mode: pmModels.ModeManual, Active: true, CredentialsRef: "env:PM_TEST_CUSTA"})
	t.Setenv("PM_TEST_CUSTA", "kid:custsecret")
	return &server{t: t, router: r, db: db}
}

func (s *server) do(method, path, apiKey, secret, idem string, body []byte, hmacHeader bool) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	ts := fmt.Sprint(time.Now().Unix())
	sig := partners.Sign(secret, ts, body)
	req.Header.Set("X-Trovotech-Timestamp", ts)
	if hmacHeader {
		req.Header.Set("Authorization", "HMAC "+apiKey+":"+sig)
	} else {
		req.Header.Set("X-TW-SERVICE-LINK-API-KEY", apiKey)
		req.Header.Set("X-Trovotech-Signature", sig)
	}
	if idem != "" {
		req.Header.Set("Idempotency-Key", idem)
	}
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)
	return w
}

func TestExchangeAuthentication(t *testing.T) {
	s := newServer(t)
	if w := s.do("GET", "/v1/trovo-api/public-markets/account", "pk_test_palmline", "s3cret", "", nil, true); w.Code != http.StatusOK {
		t.Fatalf("HMAC header: %d %s", w.Code, w.Body)
	}
	if w := s.do("GET", "/v1/trovo-api/public-markets/account", "pk_test_palmline", "s3cret", "", nil, false); w.Code != http.StatusOK {
		t.Fatalf("API key header: %d %s", w.Code, w.Body)
	}
	var acc map[string]interface{}
	json.Unmarshal(s.do("GET", "/v1/trovo-api/public-markets/account", "pk_test_palmline", "s3cret", "", nil, true).Body.Bytes(), &acc)
	if acc["balance"] != "1500" {
		t.Fatalf("account %v", acc)
	}
	if w := s.do("GET", "/v1/trovo-api/public-markets/account", "pk_test_palmline", "wrong", "", nil, true); w.Code != http.StatusUnauthorized {
		t.Fatalf("bad signature: %d", w.Code)
	}
	if w := s.do("GET", "/v1/trovo-api/public-markets/account", "pk_unknown", "s3cret", "", nil, true); w.Code < 400 {
		t.Fatalf("unknown key: %d", w.Code)
	}
	// a service link that is not onboarded for Public Markets
	if w := s.do("GET", "/v1/trovo-api/public-markets/account", "pk_test_other", "s3cret", "", nil, true); w.Code != http.StatusForbidden {
		t.Fatalf("not onboarded: %d %s", w.Code, w.Body)
	}
	s.db.Model(&slModels.ServiceLink{}).Where("id = ?", "sl-palmline").Update("suspended", 1)
	if w := s.do("GET", "/v1/trovo-api/public-markets/assets", "pk_test_palmline", "s3cret", "", nil, true); w.Code != http.StatusForbidden {
		t.Fatalf("suspended: %d", w.Code)
	}
}

func TestExchangeIdempotency(t *testing.T) {
	s := newServer(t)
	bad := []byte(`{"externalUserRef":"plm_1","legalName":"Adaeze Okafor","residencyCountry":"NGA","nationality":"NGA","ndpaConsent":true}`)
	if w := s.do("POST", "/v1/trovo-api/public-markets/wallets", "pk_test_palmline", "s3cret", "", bad, true); w.Code != http.StatusBadRequest {
		t.Fatalf("missing key: %d", w.Code)
	}
	w1 := s.do("POST", "/v1/trovo-api/public-markets/wallets", "pk_test_palmline", "s3cret", "k-1", bad, true)
	if w1.Code != http.StatusBadRequest || !bytes.Contains(w1.Body.Bytes(), []byte(`"taxIdentifier"`)) {
		t.Fatalf("field error: %d %s", w1.Code, w1.Body)
	}
	w2 := s.do("POST", "/v1/trovo-api/public-markets/wallets", "pk_test_palmline", "s3cret", "k-1", bad, true)
	if w2.Header().Get("Idempotent-Replayed") != "true" || w2.Body.String() != w1.Body.String() || w2.Code != w1.Code {
		t.Fatalf("replay: %d %q %s", w2.Code, w2.Header().Get("Idempotent-Replayed"), w2.Body)
	}
	var n int64
	s.db.Model(&pmModels.PartnerWallet{}).Where("external_user_ref = ?", "plm_1").Count(&n)
	if n != 1 {
		t.Fatalf("%d rejected requests kept, want 1", n)
	}
	other := []byte(`{"externalUserRef":"plm_2"}`)
	if w := s.do("POST", "/v1/trovo-api/public-markets/wallets", "pk_test_palmline", "s3cret", "k-1", other, true); w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("reused key: %d", w.Code)
	}
	// a server-side failure is not recorded: the same key may be retried
	good := []byte(`{"externalUserRef":"plm_3","legalName":"Adaeze Okafor","taxIdentifier":"TIN-1","residencyCountry":"NGA","nationality":"NGA","ndpaConsent":true}`)
	if w := s.do("POST", "/v1/trovo-api/public-markets/wallets", "pk_test_palmline", "s3cret", "k-2", good, true); w.Code != http.StatusServiceUnavailable {
		t.Fatalf("no chain: %d %s", w.Code, w.Body)
	}
	if w := s.do("POST", "/v1/trovo-api/public-markets/wallets", "pk_test_palmline", "s3cret", "k-2", good, true); w.Header().Get("Idempotent-Replayed") == "true" {
		t.Fatal("a 5xx was replayed")
	}
}

func TestPartnerCallback(t *testing.T) {
	s := newServer(t)
	body := []byte(`{"asOf":"2026-09-29","positions":[{"assetCode":"MTNN-T","unitsHeld":"1200"}]}`)
	post := func(secret string) *httptest.ResponseRecorder {
		req := httptest.NewRequest("POST", "/v1/custodian-partners/callbacks/positions", bytes.NewReader(body))
		ts := fmt.Sprint(time.Now().Unix())
		req.Header.Set("X-Partner-Code", "CUSTA")
		req.Header.Set("X-Trovotech-Timestamp", ts)
		req.Header.Set("X-Partner-Signature", partners.Sign(secret, ts, body))
		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)
		return w
	}
	if w := post("wrong"); w.Code != http.StatusUnauthorized {
		t.Fatalf("bad signature: %d", w.Code)
	}
	if w := post("custsecret"); w.Code != http.StatusOK || !bytes.Contains(w.Body.Bytes(), []byte("processed")) {
		t.Fatalf("first: %d %s", w.Code, w.Body)
	}
	if w := post("custsecret"); !bytes.Contains(w.Body.Bytes(), []byte("duplicate")) {
		t.Fatalf("redelivery: %d %s", w.Code, w.Body)
	}
	var p pmModels.CustodianPosition
	s.db.Where("asset_id = ?", "asset-mtnn").Order("id DESC").First(&p)
	if p.RealUnitsHeld != "1200" {
		t.Fatalf("position %+v", p)
	}
}

// testDB is a fresh SQLite database, or a fresh schema of the Postgres
// database in PUBLIC_MARKETS_TEST_POSTGRES (a DSN) when it is set.
func testDB(t *testing.T) *gorm.DB {
	t.Helper()
	cfg := &gorm.Config{Logger: logger.Discard}
	dsn := os.Getenv("PUBLIC_MARKETS_TEST_POSTGRES")
	if dsn == "" {
		db, err := gorm.Open(sqlite.Open("file:"+t.TempDir()+"/pm.db"), cfg)
		if err != nil {
			t.Fatal(err)
		}
		return db
	}
	schema := fmt.Sprintf("pmtest_%d", time.Now().UnixNano())
	admin, err := gorm.Open(postgres.Open(dsn), cfg)
	if err != nil {
		t.Fatal(err)
	}
	if err := admin.Exec("CREATE SCHEMA " + schema).Error; err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		admin.Exec("DROP SCHEMA " + schema + " CASCADE")
		if sqlDB, err := admin.DB(); err == nil {
			sqlDB.Close()
		}
	})
	db, err := gorm.Open(postgres.Open(dsn+" search_path="+schema), cfg)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

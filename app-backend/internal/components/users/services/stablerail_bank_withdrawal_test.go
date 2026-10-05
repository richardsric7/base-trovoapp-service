package users

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"trovo-wallet-api/internal/cache"
	userModels "trovo-wallet-api/internal/components/users/models"
	appdb "trovo-wallet-api/internal/db"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/ethereum/go-ethereum/common"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// fakeStablerail answers the Stablerail endpoints a bank withdrawal uses.
type fakeStablerail struct {
	mu        sync.Mutex
	offramps  []map[string]interface{}
	refuse    bool   // /cngnofframp answers 400
	status    string // /cngnofframpstatus answers this
	userCalls int
	onboards  int
	down      bool // every endpoint answers 503
}

func (f *fakeStablerail) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if r.Header.Get("x-api-key") != "key" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if f.down {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	var body map[string]interface{}
	_ = json.NewDecoder(r.Body).Decode(&body)
	switch r.URL.Path {
	case "/onboarduser":
		f.onboards++
		w.Write([]byte(`{"status":"Success","response_code":"00","message":"User registration initiated","data":{"requestId":"onb-1","status":"processing"}}`))
	case "/getuserdetails":
		f.userCalls++
		w.Write([]byte(`{"status":"Success","response_code":"00","data":{"walletDetails":{"evmWallet":"0x00000000000000000000000000000000000000aa"}}}`))
	case "/cngnofframp":
		if f.refuse {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"status":"Failed","response_code":"01","message":"Insufficient wallet balance"}`))
			return
		}
		f.offramps = append(f.offramps, body)
		w.Write([]byte(`{"status":"Success","response_code":"00","data":{"requestId":"req-1","status":"payout_pending","stage":"payout"}}`))
	case "/cngnofframpstatus":
		w.Write([]byte(`{"status":"Success","response_code":"00","data":{"requestId":"req-1","status":"` + f.status + `"}}`))
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func bankWithdrawalTestGC(t *testing.T) (*sharedconfig.GlobalConfig, *fakeStablerail) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.TempDir()+"/offramp.db"), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	appdb.MigrateDB(db) // the app's full schema, as recovery_e2e_test does
	fake := &fakeStablerail{status: "processing"}
	srv := httptest.NewServer(fake)
	t.Cleanup(srv.Close)
	// failure alerts go to a local sink, not Discord
	sink := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) }))
	t.Cleanup(sink.Close)
	t.Setenv("EXPANSION_NETWORK_ERROR_WEBHOOK", sink.URL+"/api/webhooks/0000000000/"+strings.Repeat("x", 40))
	db.Create(&userModels.StablerailConfig{ID: 1, ApiKey: "key", BaseUrl: srv.URL, EnableStablerail: 1})
	db.Create(&userModels.StablerailUser{ID: "sr-alice", TrovoUsername: "alice"})
	return &sharedconfig.GlobalConfig{DB: db, RedisCache: &cache.RedisCache{}}, fake
}

func TestStablerailEnabledAndDepositWallet(t *testing.T) {
	gc, fake := bankWithdrawalTestGC(t)
	if !StablerailEnabled(gc) {
		t.Fatal("configured and enabled")
	}
	su := GetStablerailUser("alice", gc)
	w, err := stablerailDepositWallet(&su, gc)
	if err != nil || w != common.HexToAddress("0xaa") {
		t.Fatalf("deposit wallet %v %v", w, err)
	}
	// cached afterwards
	su = GetStablerailUser("alice", gc)
	if _, err := stablerailDepositWallet(&su, gc); err != nil || fake.userCalls != 1 {
		t.Fatalf("expected one /getuserdetails call, got %d (%v)", fake.userCalls, err)
	}
	gc.DB.Model(&userModels.StablerailConfig{}).Where("id = 1").Update("enable_stablerail", 0)
	if StablerailEnabled(gc) {
		t.Fatal("disabled")
	}
}

func TestBankWithdrawalLifecycle(t *testing.T) {
	gc, fake := bankWithdrawalTestGC(t)
	hash := "0x" + strings.Repeat("ab", 32)
	draft := userModels.StablerailOfframp{
		ID: "w-1", BankCode: "058", BankName: "GTBank", AccountNumber: "0123456789", BaseAmount: 5000, Ticker: "CNGN",
		TrovoUsername: "alice", WalletAddress: "0x01", DepositAddress: "0x00000000000000000000000000000000000000aa",
	}
	ctxJSON, _ := json.Marshal(bankWithdrawalContext{Withdrawal: draft})
	ctxStr := string(ctxJSON)
	op := userModels.WalletOperation{ID: "op-1", Kind: OperationBankWithdrawal, WalletAddress: "0x01", Context: &ctxStr, UserOpHash: &hash, Status: userModels.WalletOperationSubmitted}

	// step 2 recorded it as depositing
	w := draft
	w.UserOpHash, w.Status = &hash, userModels.OfframpDepositing
	gc.DB.Create(&w)

	// the transfer is mined: deposited, then Stablerail is asked to pay out
	recordBankWithdrawalDeposit(op, common.HexToHash("0x01"), gc)
	requestOfframpPayout("w-1", gc) // the hook's goroutine may race this; only one claims it
	var got userModels.StablerailOfframp
	waitFor(t, func() bool {
		got = loadOfframp(gc, "w-1")
		return got.RequestID == "req-1"
	})
	if got.Status != "payout_pending" || got.Attempts != 1 || got.TxHash == "" {
		t.Fatalf("after payout request: %+v", got)
	}
	fake.mu.Lock()
	if len(fake.offramps) != 1 || fake.offramps[0]["userId"] != "sr-alice" || fake.offramps[0]["amount"] != float64(5000) || fake.offramps[0]["bankCode"] != "058" {
		t.Fatalf("offramp requests %+v", fake.offramps)
	}
	fake.mu.Unlock()

	// followed until it ends
	UpdateOfframpStatuses(gc)
	got = loadOfframp(gc, "w-1")
	if got.Status != "processing" {
		t.Fatalf("status %v", got.Status)
	}
	fake.status = "completed"
	UpdateOfframpStatuses(gc)
	got = loadOfframp(gc, "w-1")
	if got.Status != "completed" {
		t.Fatalf("status %v", got.Status)
	}
	// a second run does nothing more
	UpdateOfframpStatuses(gc)
	if len(fake.offramps) != 1 {
		t.Fatal("paid out twice")
	}
}

func TestBankWithdrawalPayoutRetries(t *testing.T) {
	gc, fake := bankWithdrawalTestGC(t)
	fake.refuse = true
	gc.DB.Create(&userModels.StablerailOfframp{ID: "w-2", BankCode: "058", AccountNumber: "0123456789", BaseAmount: 2000, Ticker: "CNGN", TrovoUsername: "alice", Status: userModels.OfframpDeposited})
	for i := 0; i < maxOfframpRequestAttempts+3; i++ {
		UpdateOfframpStatuses(gc)
	}
	var got userModels.StablerailOfframp
	got = loadOfframp(gc, "w-2")
	if got.Status != userModels.OfframpRequestFailed || got.Attempts != maxOfframpRequestAttempts || !strings.Contains(got.Error, "Insufficient") {
		t.Fatalf("after refusals: %+v", got)
	}
	// support tops up the attempts budget (or Stablerail recovers): it goes through
	fake.refuse = false
	gc.DB.Model(&userModels.StablerailOfframp{}).Where("id = ?", "w-2").Update("attempts", 0)
	UpdateOfframpStatuses(gc)
	got = loadOfframp(gc, "w-2")
	if got.RequestID != "req-1" {
		t.Fatalf("after recovery: %+v", got)
	}

	// an interrupted request is retried after 5 minutes
	gc.DB.Create(&userModels.StablerailOfframp{ID: "w-3", BaseAmount: 1000, TrovoUsername: "alice", Status: userModels.OfframpRequesting})
	gc.DB.Exec("UPDATE stablerail_offramps SET updated_at = datetime('now', '-10 minutes') WHERE id = 'w-3'")
	UpdateOfframpStatuses(gc)
	got = loadOfframp(gc, "w-3")
	if got.Status != userModels.OfframpRequestFailed {
		t.Fatalf("interrupted request: %+v", got)
	}
}

func TestBankWithdrawalDepositFailed(t *testing.T) {
	gc, fake := bankWithdrawalTestGC(t)
	hash, failed := "0x"+strings.Repeat("cd", 32), false
	reason := "the operation's calls reverted: transfer amount exceeds balance"
	gc.DB.Create(&userModels.WalletOperation{ID: "op-2", Kind: OperationBankWithdrawal, WalletAddress: "0x01", UserOpHash: &hash, Status: userModels.WalletOperationIncluded, Success: &failed, Error: &reason})
	gc.DB.Create(&userModels.StablerailOfframp{ID: "w-4", BaseAmount: 1000, TrovoUsername: "alice", Status: userModels.OfframpDepositing, UserOpHash: &hash})
	UpdateOfframpStatuses(gc)
	var got userModels.StablerailOfframp
	got = loadOfframp(gc, "w-4")
	if got.Status != userModels.OfframpDepositFailed || got.Error != reason {
		t.Fatalf("failed deposit: %+v", got)
	}
	if len(fake.offramps) != 0 {
		t.Fatal("no payout for a failed deposit")
	}

	// mined but missed by the hook: picked up and paid out
	ok, hash2 := true, "0x"+strings.Repeat("ef", 32)
	tx := "0x" + strings.Repeat("12", 32)
	gc.DB.Create(&userModels.WalletOperation{ID: "op-3", Kind: OperationBankWithdrawal, WalletAddress: "0x01", UserOpHash: &hash2, Status: userModels.WalletOperationIncluded, Success: &ok, TxHash: &tx})
	gc.DB.Create(&userModels.StablerailOfframp{ID: "w-5", BankCode: "058", AccountNumber: "0123456789", BaseAmount: 1500, Ticker: "CNGN", TrovoUsername: "alice", Status: userModels.OfframpDepositing, UserOpHash: &hash2})
	UpdateOfframpStatuses(gc)
	UpdateOfframpStatuses(gc)
	got = loadOfframp(gc, "w-5")
	if got.TxHash != tx || got.RequestID != "req-1" {
		t.Fatalf("missed deposit: %+v", got)
	}
}

func waitFor(t *testing.T, cond func() bool) {
	t.Helper()
	for i := 0; i < 100; i++ {
		if cond() {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatal("timed out")
}

func loadOfframp(gc *sharedconfig.GlobalConfig, id string) (w userModels.StablerailOfframp) {
	gc.DB.Where("id = ?", id).First(&w)
	return
}

func TestStablerailOnboardingRetries(t *testing.T) {
	gc, fake := bankWithdrawalTestGC(t)
	gc.DB.Create(&userModels.User{Username: "bob"})
	// the KYC callback could not reach Stablerail and saved a retry
	gc.DB.Create(&userModels.StablerailOnboardUserRetry{TrovoUsername: "bob", BVN: "22222222222"})
	gc.DB.Exec("UPDATE stablerail_onboard_user_retries SET updated_at = datetime('now', '-10 minutes')")

	fake.down = true
	ProcessStablerailOnboardingRetries(gc)
	var r userModels.StablerailOnboardUserRetry
	if gc.DB.First(&r).Error != nil || r.Attempts != 1 {
		t.Fatalf("a failed retry is counted: %+v", r)
	}
	ProcessStablerailOnboardingRetries(gc) // not again within 5 minutes
	gc.DB.First(&r)
	if r.Attempts != 1 {
		t.Fatalf("retried too soon: %+v", r)
	}

	fake.down = false
	gc.DB.Exec("UPDATE stablerail_onboard_user_retries SET updated_at = datetime('now', '-10 minutes')")
	ProcessStablerailOnboardingRetries(gc)
	var left int64
	gc.DB.Model(&userModels.StablerailOnboardUserRetry{}).Count(&left)
	if left != 0 || fake.onboards != 1 {
		t.Fatalf("retry should start the onboarding and be removed (left %d, onboards %d)", left, fake.onboards)
	}
	// a repeated KYC callback does not onboard twice while it is in progress
	u, _ := userModels.Username("bob").GetSimpleUser(gc.DB, gc)
	if msg, err := StablerailInitiateOnboardUser(&u, "22222222222", gc); err != nil || msg != "registration in progress" || fake.onboards != 1 {
		t.Fatalf("second onboarding: %q %v (onboards %d)", msg, err, fake.onboards)
	}
	if MaskBVN("22222222222") != "********222" {
		t.Fatal(MaskBVN("22222222222"))
	}
}

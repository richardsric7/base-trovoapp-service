package services

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// newMockDB returns a gorm DB backed by sqlmock so ping behaviour can be
// controlled without a real database.
func newMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	sqlDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	// gorm.Open pings during initialisation; that ping is not the one under
	// test, so satisfy it here and let each test set its own expectation.
	mock.ExpectPing()
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn:                 sqlDB,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("gorm open: %v", err)
	}
	return gormDB, mock
}

func depByName(deps []DependencyStatus, name string) (DependencyStatus, bool) {
	for _, d := range deps {
		if d.Name == name {
			return d, true
		}
	}
	return DependencyStatus{}, false
}

func TestLivenessIgnoresDependencies(t *testing.T) {
	// No databases, no cache, no wallet - liveness must still report up,
	// because the process itself is fine.
	s := New(nil, nil, nil, nil)
	got := s.Liveness()
	if got.Status != StatusUp {
		t.Fatalf("status = %q, want %q", got.Status, StatusUp)
	}
	if got.Service != "admin-api" {
		t.Errorf("service = %q, want admin-api", got.Service)
	}
	if got.Version == "" {
		t.Error("version should never be empty")
	}
	if got.Time == "" {
		t.Error("time should be set")
	}
}

func TestReadinessAllUp(t *testing.T) {
	adminDB, adminMock := newMockDB(t)
	walletDB, walletMock := newMockDB(t)
	p2pDB, p2pMock := newMockDB(t)
	for _, m := range []sqlmock.Sqlmock{adminMock, walletMock, p2pMock} {
		m.ExpectPing()
	}
	t.Setenv("ENABLE_CACHING", "0")

	s := New(adminDB, walletDB, p2pDB, nil)
	report, status := s.Readiness(context.Background())

	if status != http.StatusOK {
		t.Fatalf("status = %d, want 200", status)
	}
	if report.Status != StatusUp {
		t.Fatalf("report status = %q, want %q", report.Status, StatusUp)
	}
	for _, name := range []string{"admin_db", "wallet_db", "p2p_db"} {
		dep, ok := depByName(report.Dependencies, name)
		if !ok {
			t.Fatalf("missing dependency %q", name)
		}
		if dep.Status != StatusUp {
			t.Errorf("%s = %q, want up (err=%q)", name, dep.Status, dep.Error)
		}
	}
}

func TestReadinessRequiredDependencyDownReturns503(t *testing.T) {
	adminDB, adminMock := newMockDB(t)
	adminMock.ExpectPing().WillReturnError(context.DeadlineExceeded)
	walletDB, walletMock := newMockDB(t)
	walletMock.ExpectPing()
	p2pDB, p2pMock := newMockDB(t)
	p2pMock.ExpectPing()
	t.Setenv("ENABLE_CACHING", "0")

	s := New(adminDB, walletDB, p2pDB, nil)
	report, status := s.Readiness(context.Background())

	if status != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want 503", status)
	}
	if report.Status != StatusDown {
		t.Fatalf("report status = %q, want %q", report.Status, StatusDown)
	}
	dep, _ := depByName(report.Dependencies, "admin_db")
	if dep.Status != StatusDown {
		t.Errorf("admin_db = %q, want down", dep.Status)
	}
	if dep.Error == "" {
		t.Error("a down dependency must name the error")
	}
}

func TestReadinessOptionalDependencyDownIsDegradedNot503(t *testing.T) {
	adminDB, adminMock := newMockDB(t)
	walletDB, walletMock := newMockDB(t)
	p2pDB, p2pMock := newMockDB(t)
	for _, m := range []sqlmock.Sqlmock{adminMock, walletMock, p2pMock} {
		m.ExpectPing()
	}
	t.Setenv("ENABLE_CACHING", "0")

	// Point the wallet check at a closed port so it fails to connect.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	closedURL := srv.URL
	srv.Close()

	s := New(adminDB, walletDB, p2pDB, nil)
	s.wallet = serviceLinkAdapter{baseURL: closedURL}

	report, status := s.Readiness(context.Background())

	if status != http.StatusOK {
		t.Fatalf("status = %d, want 200 - an optional dependency must not fail readiness", status)
	}
	if report.Status != StatusDegraded {
		t.Fatalf("report status = %q, want %q", report.Status, StatusDegraded)
	}
	dep, _ := depByName(report.Dependencies, "wallet_api")
	if dep.Status != StatusDown {
		t.Errorf("wallet_api = %q, want down", dep.Status)
	}
	if !dep.Optional {
		t.Error("wallet_api must be marked optional")
	}
}

func TestWalletCheckUpOnAnyHTTPResponse(t *testing.T) {
	// A 401 from the wallet base URL still proves it is listening.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	s := New(nil, nil, nil, nil)
	s.wallet = serviceLinkAdapter{baseURL: srv.URL}

	dep := s.checkWallet(context.Background())
	if dep.Status != StatusUp {
		t.Fatalf("wallet_api = %q, want up on a 401 response", dep.Status)
	}
}

func TestUnconfiguredDependenciesAreSkippedNotFailed(t *testing.T) {
	t.Setenv("ENABLE_CACHING", "0")
	s := New(nil, nil, nil, nil)
	report, status := s.Readiness(context.Background())

	// Nothing is configured, so nothing can be down; readiness must not 503.
	if status != http.StatusOK {
		t.Fatalf("status = %d, want 200", status)
	}
	for _, dep := range report.Dependencies {
		if dep.Status != StatusSkipped {
			t.Errorf("%s = %q, want skipped", dep.Name, dep.Status)
		}
		if dep.SkipReason == "" {
			t.Errorf("%s skipped without a reason", dep.Name)
		}
	}
}

func TestRedisSkippedWhenCachingDisabled(t *testing.T) {
	t.Setenv("ENABLE_CACHING", "0")
	s := New(nil, nil, nil, nil)
	dep := s.checkRedis(context.Background())
	if dep.Status != StatusSkipped {
		t.Fatalf("redis = %q, want skipped", dep.Status)
	}
	if dep.SkipReason != "caching disabled" {
		t.Errorf("skip reason = %q", dep.SkipReason)
	}
}

func TestDependencyLatencyIsRecorded(t *testing.T) {
	// Regression: the latency assignment used to sit in a deferred func on an
	// unnamed return, so the write was lost and every dependency reported
	// latency 0 - which hid a dependency slowly degrading.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	s := New(nil, nil, nil, nil)
	s.wallet = serviceLinkAdapter{baseURL: srv.URL}

	dep := s.checkWallet(context.Background())
	if dep.Status != StatusUp {
		t.Fatalf("wallet_api = %q, want up", dep.Status)
	}
	if dep.LatencyMs < 1 {
		t.Errorf("latency_ms = %d, want >= 1 for a 5ms round trip", dep.LatencyMs)
	}
	if dep.LatencyUs < 1000 {
		t.Errorf("latency_us = %d, want >= 1000 for a 5ms round trip", dep.LatencyUs)
	}
}

func TestSlowDependencyLatencyIsVisibleInMicroseconds(t *testing.T) {
	// A healthy local database pings in well under a millisecond, so
	// latency_ms rounds to 0 - latency_us is what makes a dependency that is
	// slowing down visible before it fails outright.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(300 * time.Microsecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	s := New(nil, nil, nil, nil)
	s.wallet = serviceLinkAdapter{baseURL: srv.URL}

	dep := s.checkWallet(context.Background())
	if dep.Status != StatusUp {
		t.Fatalf("wallet_api = %q, want up", dep.Status)
	}
	if dep.LatencyUs < 300 {
		t.Errorf("latency_us = %d, want >= 300 for a 300us round trip", dep.LatencyUs)
	}
}

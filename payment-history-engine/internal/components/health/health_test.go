package health

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

// resetStreamState clears the package-level counters between tests.
func resetStreamState(t *testing.T) {
	t.Helper()
	atomic.StoreInt64(&lastOperationNanos, 0)
	atomic.StoreUint64(&operationsProcessed, 0)
	cursorMu.Lock()
	lastCursor = ""
	cursorMu.Unlock()
	atomic.StoreInt32(&hasWorkFlag, 0)
	startedAt = time.Now()
}

func TestLivenessIgnoresDependencies(t *testing.T) {
	// The process being alive is a separate question from whether its
	// dependencies are: liveness must not fail on a database outage, because
	// restarting the process would not fix one.
	s := New(nil, nil, nil, nil)
	got := s.Liveness()

	if got.Status != StatusUp {
		t.Fatalf("status = %q, want up", got.Status)
	}
	if got.Service != "payment-history-engine" {
		t.Errorf("service = %q", got.Service)
	}
}

func TestStreamReportsProgress(t *testing.T) {
	resetStreamState(t)
	SetHasWork(true)

	RecordOperation("cursor-1")
	RecordOperation("cursor-2")

	st := stream()
	if st.Status != StatusUp {
		t.Fatalf("status = %q, want up after recent operations", st.Status)
	}
	if st.OperationsProcessed != 2 {
		t.Errorf("processed = %d, want 2", st.OperationsProcessed)
	}
	if st.LastCursor != "cursor-2" {
		t.Errorf("cursor = %q, want the most recent", st.LastCursor)
	}
	if st.LastOperationAt == "" {
		t.Error("last operation time should be set")
	}
}

func TestStreamReportsStallAfterThreshold(t *testing.T) {
	// The failure this whole component exists for: the process is fine, the
	// database is fine, and the stream silently stopped.
	resetStreamState(t)
	SetHasWork(true)
	RecordOperation("cursor-1")
	atomic.StoreInt64(&lastOperationNanos, time.Now().Add(-stallThreshold-time.Minute).UnixNano())

	st := stream()
	if st.Status != StatusStalled {
		t.Fatalf("status = %q, want stalled", st.Status)
	}
	if st.Note == "" {
		t.Error("a stalled stream must explain itself")
	}
	if st.SecondsSinceLastOp < int64(stallThreshold.Seconds()) {
		t.Errorf("seconds since last op = %d, expected at least the threshold", st.SecondsSinceLastOp)
	}
}

func TestStreamStartupIsNotAStall(t *testing.T) {
	// Having processed nothing yet is normal in the first moments after boot,
	// and must not be reported as a failure.
	resetStreamState(t)
	SetHasWork(true)

	st := stream()
	if st.Status != StatusSkipped {
		t.Fatalf("status = %q, want skipped during startup", st.Status)
	}
	if st.Note == "" {
		t.Error("should explain that it is starting up")
	}
}

func TestNeverStartedStreamEventuallyReportsStalled(t *testing.T) {
	// If nothing has been processed long after boot, the stream never started -
	// which is a real failure, not a startup delay.
	resetStreamState(t)
	SetHasWork(true)
	startedAt = time.Now().Add(-stallThreshold - time.Minute)

	st := stream()
	if st.Status != StatusStalled {
		t.Fatalf("status = %q, want stalled", st.Status)
	}
}

func TestStalledStreamDegradesButDoesNotReturn503(t *testing.T) {
	// Deliberate: a restart cannot restart a stream that stopped because the
	// network went quiet, so marking the container unhealthy would produce a
	// restart loop rather than a fix. The status is there to be alerted on.
	resetStreamState(t)
	SetHasWork(true)
	RecordOperation("c1")
	atomic.StoreInt64(&lastOperationNanos, time.Now().Add(-stallThreshold-time.Minute).UnixNano())
	t.Setenv("ENABLE_CACHING", "0")

	report, status := New(nil, nil, nil, nil).Readiness(context.Background())

	if status != http.StatusOK {
		t.Fatalf("status = %d, want 200 - a stall must not trigger a restart loop", status)
	}
	if report.Status != StatusDegraded {
		t.Errorf("report status = %q, want degraded", report.Status)
	}
}

func TestUnconfiguredDependenciesAreSkipped(t *testing.T) {
	resetStreamState(t)
	t.Setenv("ENABLE_CACHING", "0")

	report, status := New(nil, nil, nil, nil).Readiness(context.Background())

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

func TestHorizonUpOnAnyHTTPResponse(t *testing.T) {
	// Horizon's root may answer 404; that still proves it is reachable.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	s := New(nil, nil, nil, nil)
	s.rpcURL = srv.URL

	dep := s.checkBlockchainRPC(context.Background())
	if dep.Status != StatusUp {
		t.Fatalf("base_rpc = %q, want up on a 404", dep.Status)
	}
}

func TestHorizonDownDegradesRatherThanFails(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	closed := srv.URL
	srv.Close()

	resetStreamState(t)
	SetHasWork(true)
	RecordOperation("c1")
	t.Setenv("ENABLE_CACHING", "0")

	s := New(nil, nil, nil, nil)
	s.rpcURL = closed

	report, status := s.Readiness(context.Background())
	if status != http.StatusOK {
		t.Fatalf("status = %d, want 200 - the engine resumes from its cursor once Horizon returns", status)
	}
	if report.Status != StatusDegraded {
		t.Errorf("report status = %q, want degraded", report.Status)
	}
}

func TestLatencyIsRecorded(t *testing.T) {
	// Guards the deferred-write-to-unnamed-return bug found in the admin API,
	// where every latency silently reported 0.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Millisecond)
	}))
	defer srv.Close()

	s := New(nil, nil, nil, nil)
	s.rpcURL = srv.URL

	dep := s.checkBlockchainRPC(context.Background())
	if dep.LatencyMs < 1 || dep.LatencyUs < 1000 {
		t.Errorf("latency not recorded: ms=%d us=%d", dep.LatencyMs, dep.LatencyUs)
	}
}

func TestProbesServeOverHTTP(t *testing.T) {
	resetStreamState(t)
	SetHasWork(true)
	RecordOperation("cursor-abc")
	t.Setenv("ENABLE_CACHING", "0")

	s := New(nil, nil, nil, nil)
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, s.Liveness())
	})
	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		report, status := s.Readiness(r.Context())
		writeJSON(w, status, report)
	})

	for _, path := range []string{"/health", "/ready"} {
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))

		if rec.Code != http.StatusOK {
			t.Fatalf("%s = %d, want 200", path, rec.Code)
		}
		var body map[string]any
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("%s: invalid JSON: %v", path, err)
		}
		if body["service"] != "payment-history-engine" {
			t.Errorf("%s service = %v", path, body["service"])
		}
	}
}

func TestConcurrentRecordOperationIsSafe(t *testing.T) {
	// The stream workers call this on their hot path; a data race here would
	// be a real bug in production.
	resetStreamState(t)

	done := make(chan struct{})
	for i := 0; i < 50; i++ {
		go func(i int) {
			RecordOperation("cursor")
			done <- struct{}{}
		}(i)
	}
	for i := 0; i < 50; i++ {
		<-done
	}

	if got := atomic.LoadUint64(&operationsProcessed); got != 50 {
		t.Errorf("processed = %d, want 50", got)
	}
}

func TestNoWorkIsIdleNotStalled(t *testing.T) {
	// The false alarm this distinction exists to prevent: an engine with no
	// wallets to watch has nothing to process, so zero operations is correct
	// and healthy. Reporting it as stalled makes a quiet environment look
	// broken forever, which teaches people to ignore the indicator.
	resetStreamState(t)
	SetHasWork(false)
	startedAt = time.Now().Add(-2 * stallThreshold)

	st := stream()
	if st.Status != StatusIdle {
		t.Fatalf("status = %q, want idle - there is no work to do", st.Status)
	}
	if st.Note == "" {
		t.Error("idle should explain that nothing is being tracked")
	}
}

func TestIdleDoesNotDegradeReadiness(t *testing.T) {
	// And an idle engine must report itself fully healthy, not degraded.
	resetStreamState(t)
	SetHasWork(false)
	startedAt = time.Now().Add(-2 * stallThreshold)
	t.Setenv("ENABLE_CACHING", "0")

	report, status := New(nil, nil, nil, nil).Readiness(context.Background())

	if status != http.StatusOK {
		t.Fatalf("status = %d, want 200", status)
	}
	if report.Status != StatusUp {
		t.Errorf("report status = %q, want up - having no work is not a fault", report.Status)
	}
}

func TestWorkAppearingRestoresStallDetection(t *testing.T) {
	// Once wallets are tracked, a stream that then does nothing IS a problem,
	// so the check must come back rather than being permanently disabled.
	resetStreamState(t)
	SetHasWork(false)
	if stream().Status != StatusIdle {
		t.Fatal("precondition: should start idle")
	}

	SetHasWork(true)
	startedAt = time.Now().Add(-2 * stallThreshold)
	if got := stream().Status; got != StatusStalled {
		t.Errorf("status = %q, want stalled once there is work but nothing processed", got)
	}
}

package observe

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/testutil"
)

// resetFailureState clears the rate limiter between tests so one test's reports
// do not suppress another's.
func resetFailureState(t *testing.T) {
	t.Helper()
	lastReportMu.Lock()
	lastReport = map[string]time.Time{}
	lastReportMu.Unlock()
	handledFailuresTotal.Reset()
}

func TestOperationFromExtractsBracketPrefix(t *testing.T) {
	cases := []struct {
		name string
		msg  string
		want string
	}{
		{
			name: "the convention the existing call sites use",
			msg:  "[UpdateStablerailCNGNOnrampStatus] error saving onramp record for user bob: boom",
			want: "UpdateStablerailCNGNOnrampStatus",
		},
		{
			name: "leading whitespace is tolerated",
			msg:  "  [StableRailInitiateAssetWithdrawal] unable to save",
			want: "StableRailInitiateAssetWithdrawal",
		},
		{
			name: "digits are allowed after the first letter",
			msg:  "[Swap2Send] failed",
			want: "Swap2Send",
		},
		{
			name: "no prefix falls back rather than inventing a label",
			msg:  "FAILED to trigger stablerail onboarding for user alice",
			want: "unlabelled",
		},
		{
			name: "a bracket later in the message is not a prefix",
			msg:  "something broke for user [bob]",
			want: "unlabelled",
		},
		{
			name: "spaces inside the brackets are not an operation name",
			msg:  "[KYC WEBHOOK ERROR] User ID is invalid",
			want: "unlabelled",
		},
		{
			name: "empty message",
			msg:  "",
			want: "unlabelled",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := operationFrom(tc.msg); got != tc.want {
				t.Fatalf("operationFrom(%q) = %q, want %q", tc.msg, got, tc.want)
			}
		})
	}
}

// An operation label becomes a Prometheus time series, so an unbounded label
// would be a cardinality leak. The regex caps the length; anything longer must
// fall back rather than be truncated into a near-duplicate series.
func TestOperationFromRejectsOverlongPrefix(t *testing.T) {
	long := ""
	for i := 0; i < 100; i++ {
		long += "A"
	}
	if got := operationFrom("[" + long + "] boom"); got != "unlabelled" {
		t.Fatalf("expected overlong prefix to fall back, got %q", got)
	}
}

func TestRecordHandledFailureIncrementsCounterPerOperation(t *testing.T) {
	resetFailureState(t)

	RecordHandledFailure("[AlphaOp] first failure for user bob")
	RecordHandledFailure("[AlphaOp] second failure for user alice")
	RecordHandledFailure("[BetaOp] different operation")

	if got := testutil.ToFloat64(handledFailuresTotal.WithLabelValues(serviceName, "AlphaOp")); got != 2 {
		t.Fatalf("AlphaOp count = %v, want 2", got)
	}
	if got := testutil.ToFloat64(handledFailuresTotal.WithLabelValues(serviceName, "BetaOp")); got != 1 {
		t.Fatalf("BetaOp count = %v, want 1", got)
	}
}

// The counter is the signal the dashboard reads, so it must move even when
// crash reporting is switched off - which is the default in local development
// and in any environment without a DSN.
func TestRecordHandledFailureCountsWhenReportingDisabled(t *testing.T) {
	resetFailureState(t)

	orig := initialised
	initialised = false
	defer func() { initialised = orig }()

	RecordHandledFailure("[DisabledOp] boom")

	if got := testutil.ToFloat64(handledFailuresTotal.WithLabelValues(serviceName, "DisabledOp")); got != 1 {
		t.Fatalf("count = %v, want 1 even with reporting disabled", got)
	}
}

// A poll loop failing every 10s must not file 8,640 crash reports a day, but
// every one of those failures still has to show up in the rate the dashboard
// plots. Throttling applies to reports only.
func TestReportRateLimitDoesNotThrottleTheCounter(t *testing.T) {
	resetFailureState(t)

	const n = 50
	for i := 0; i < n; i++ {
		RecordHandledFailure(fmt.Sprintf("[NoisyOp] failure %d", i))
	}

	if got := testutil.ToFloat64(handledFailuresTotal.WithLabelValues(serviceName, "NoisyOp")); got != n {
		t.Fatalf("count = %v, want %d - the counter must not be rate limited", got, n)
	}
}

func TestShouldReportThrottlesPerOperationAndReopens(t *testing.T) {
	resetFailureState(t)

	base := time.Now()

	if !shouldReport("GammaOp", base) {
		t.Fatal("first report for an operation should be allowed")
	}
	if shouldReport("GammaOp", base.Add(time.Second)) {
		t.Fatal("a second report inside the window should be suppressed")
	}
	// A different operation has its own budget - one noisy worker must not
	// silence the rest.
	if !shouldReport("DeltaOp", base.Add(time.Second)) {
		t.Fatal("a different operation should not share the throttle")
	}
	if !shouldReport("GammaOp", base.Add(failureReportLimit+time.Second)) {
		t.Fatal("the window should reopen once the limit has elapsed")
	}
}

func TestRecordHandledFailureErrIgnoresNil(t *testing.T) {
	resetFailureState(t)

	RecordHandledFailureErr("NilOp", nil)

	if got := testutil.ToFloat64(handledFailuresTotal.WithLabelValues(serviceName, "NilOp")); got != 0 {
		t.Fatalf("count = %v, want 0 for a nil error", got)
	}
}

func TestRecordHandledFailureErrCountsAndLabels(t *testing.T) {
	resetFailureState(t)

	// log.Printf writes to the standard logger's own output (stderr by
	// default), not os.Stdout, so redirect the logger itself rather than
	// reusing captureStdout.
	var buf bytes.Buffer
	origOut := log.Writer()
	origFlags := log.Flags()
	log.SetOutput(&buf)
	log.SetFlags(0)
	defer func() {
		log.SetOutput(origOut)
		log.SetFlags(origFlags)
	}()

	RecordHandledFailureErr("EpsilonOp", errors.New("database is down"))

	if got := testutil.ToFloat64(handledFailuresTotal.WithLabelValues(serviceName, "EpsilonOp")); got != 1 {
		t.Fatalf("count = %v, want 1", got)
	}
	out := buf.String()
	if want := "operation=EpsilonOp"; !strings.Contains(out, want) {
		t.Fatalf("log output %q does not contain %q", out, want)
	}
	if want := "database is down"; !strings.Contains(out, want) {
		t.Fatalf("log output %q does not contain the error text %q", out, want)
	}
}

func TestRecordHandledFailureErrDefaultsEmptyOperation(t *testing.T) {
	resetFailureState(t)

	RecordHandledFailureErr("", errors.New("boom"))

	if got := testutil.ToFloat64(handledFailuresTotal.WithLabelValues(serviceName, "unlabelled")); got != 1 {
		t.Fatalf("count = %v, want 1 under the fallback label", got)
	}
}

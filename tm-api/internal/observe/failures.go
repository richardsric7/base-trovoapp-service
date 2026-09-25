package observe

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/prometheus/client_golang/prometheus"
)

// handledFailuresTotal counts failures the code caught and handled deliberately:
// a database write that came back with an error, an upstream call that returned
// a bad response, a state we could not reconcile. These never reach the HTTP
// error rate - most happen in background workers that serve no request at all,
// and the ones that do happen in a handler are usually reported to the user as a
// clean 4xx or an empty result.
//
// That gap is the point. The onramp sweep failing to record a completed deposit
// produced no panic, no 5xx, and no change to any existing metric, so a money
// path stayed broken for two days behind a green dashboard. This counter is what
// moves in that situation.
var handledFailuresTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "handled_failures_total",
		Help: "Failures the application caught and handled, by originating operation.",
	},
	[]string{"service", "operation"},
)

func init() {
	prometheus.MustRegister(handledFailuresTotal)
}

// operationPrefix matches the "[SomeFunction]" convention the existing failure
// messages already start with, e.g.
//
//	[UpdateStablerailCNGNOnrampStatus] error saving onramp record for ...
//
// Extracting it gives a per-operation label without editing any of the call
// sites. Only letters and digits are accepted so that a stray bracket in a
// message body cannot invent a label, and the length is bounded because every
// distinct label value is a new Prometheus time series.
var operationPrefix = regexp.MustCompile(`^\[([A-Za-z][A-Za-z0-9]{0,63})\]`)

// operationFrom pulls the operation name out of a failure message, falling back
// to a single shared bucket when the message does not follow the convention.
// The fallback is deliberately one fixed string rather than anything derived
// from the message: message bodies carry usernames, amounts and request IDs, and
// using them as labels would both explode cardinality and put customer data into
// metric names.
func operationFrom(msg string) string {
	if m := operationPrefix.FindStringSubmatch(strings.TrimSpace(msg)); m != nil {
		return m[1]
	}
	return "unlabelled"
}

// failureReportLimit is the minimum gap between crash reports for the same
// operation. A poll loop that fails every 10s would otherwise file 8,640 events
// a day on its own and drown everything else in the project. The Prometheus
// counter is NOT rate limited - it is cheap, and the rate is the thing the
// dashboard and any future alert actually read. Only the report is throttled.
const failureReportLimit = 5 * time.Minute

var (
	lastReportMu sync.Mutex
	lastReport   = map[string]time.Time{}
)

// shouldReport reports whether enough time has passed to file another crash
// report for this operation.
func shouldReport(operation string, now time.Time) bool {
	lastReportMu.Lock()
	defer lastReportMu.Unlock()

	if prev, ok := lastReport[operation]; ok && now.Sub(prev) < failureReportLimit {
		return false
	}
	lastReport[operation] = now
	return true
}

// RecordHandledFailure records a failure that the caller has already decided is
// worth someone's attention. It always increments the Prometheus counter and,
// subject to the per-operation rate limit above, files a crash report grouped by
// operation rather than by message text.
//
// Grouping matters here: these messages embed usernames, amounts and request
// IDs, so reporting the raw text would file a separate issue per occurrence -
// five hundred onramp failures as five hundred issues instead of one with a
// count of five hundred. The message is attached as context instead, and passes
// through the same scrubbing as every other event.
//
// Safe to call before Init: the counter still moves and reporting is skipped.
func RecordHandledFailure(msg string) {
	operation := operationFrom(msg)
	handledFailuresTotal.WithLabelValues(serviceName, operation).Inc()

	if !initialised || !shouldReport(operation, time.Now()) {
		return
	}

	hub := sentry.CurrentHub().Clone()
	hub.Scope().SetTag("operation", operation)
	hub.Scope().SetTag("failure_kind", "handled")
	hub.Scope().SetLevel(sentry.LevelError)
	// The title is the operation alone so the issue groups on it; the full
	// message goes in context, where scrubbing still applies.
	hub.Scope().SetExtra("message", msg)
	hub.CaptureException(fmt.Errorf("handled failure in %s", operation))
}

// RecordHandledFailureErr is the same as RecordHandledFailure for callers that
// have a real error value rather than a preformatted string. It is preferred
// where available: CaptureException unwraps %w chains and attaches a stack
// trace, neither of which survives being formatted into a message.
func RecordHandledFailureErr(operation string, err error) {
	if err == nil {
		return
	}
	if operation == "" {
		operation = "unlabelled"
	}
	handledFailuresTotal.WithLabelValues(serviceName, operation).Inc()
	log.Printf("[observe] handled failure operation=%s: %v", operation, err)

	if !initialised || !shouldReport(operation, time.Now()) {
		return
	}

	hub := sentry.CurrentHub().Clone()
	hub.Scope().SetTag("operation", operation)
	hub.Scope().SetTag("failure_kind", "handled")
	hub.Scope().SetLevel(sentry.LevelError)
	hub.CaptureException(err)
}

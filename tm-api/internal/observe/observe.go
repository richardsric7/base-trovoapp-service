// Package observe provides the admin API's instrumentation: request
// correlation IDs, structured request logging, crash reporting to GlitchTip
// and Prometheus metrics.
//
// Everything here degrades quietly. Crash reporting with no DSN configured
// logs locally instead of reporting; metrics are always collected but only
// leave the process when something scrapes /metrics. Instrumentation must
// never be the reason a request fails.
package observe

import (
	"log"
	"os"
	"strings"
	"time"

	healthservices "admin-panel-dashboard/internal/components/health/services"

	"github.com/getsentry/sentry-go"
)

// RequestIDHeader is the response header carrying the correlation ID. The
// same value appears on every log line for the request, so an ID quoted in a
// bug report leads straight to that request's logs across services.
const RequestIDHeader = "X-Request-ID"

// CtxRequestID is the gin context key holding the correlation ID.
const CtxRequestID = "observe_request_id"

// serviceName identifies this service in logs and error reports.
const serviceName = "admin-api"

// initialised records whether crash reporting is actually active, so the
// middleware can skip Sentry work entirely when no DSN was configured.
var initialised bool

// Init starts crash reporting. It is safe (and expected) to call with no DSN
// configured: reporting is then disabled and panics are logged locally.
//
// An empty version falls back to the same build stamp /health reports, so a
// crash report always names the build it came from - without it the release
// shows blank and "did my fix work?" becomes unanswerable.
//
// Returns a flush function to call on shutdown so buffered events are
// delivered before the process exits.
func Init(version string) func() {
	if version == "" {
		version = healthservices.Version()
	}
	dsn := strings.TrimSpace(os.Getenv("SENTRY_DSN"))
	if dsn == "" {
		log.Println("[observe] SENTRY_DSN not set - crash reporting disabled, panics will be logged locally")
		return func() {}
	}

	environment := os.Getenv("APP_ENV")
	if environment == "" {
		environment = "development"
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn:         dsn,
		Environment: environment,
		Release:     version,
		// Errors only: GlitchTip implements Sentry's error API, not its
		// performance/tracing API. Response-time trends come from Prometheus.
		EnableTracing: false,
		// Strip anything that could carry credentials or personal data before
		// it leaves the process.
		BeforeSend: scrubEvent,
	})
	if err != nil {
		// A bad DSN must not stop the API from starting.
		log.Printf("[observe] crash reporting disabled, sentry init failed: %v", err)
		return func() {}
	}

	initialised = true
	log.Printf("[observe] crash reporting enabled (environment=%s release=%s)", environment, version)
	return func() { sentry.Flush(2 * time.Second) }
}

// Enabled reports whether crash reporting is active.
func Enabled() bool { return initialised }

// sensitiveHeaders are removed from error reports outright. Authorization and
// Cookie carry live credentials; the rest are common proxy-forwarded secrets.
var sensitiveHeaders = []string{
	"Authorization",
	"Cookie",
	"Set-Cookie",
	"X-Api-Key",
	"X-Auth-Token",
	"Proxy-Authorization",
}

// scrubEvent removes credentials from an outgoing error report. Error
// reporting must never become a channel for leaking tokens into a third-party
// store, even a self-hosted one.
func scrubEvent(event *sentry.Event, _ *sentry.EventHint) *sentry.Event {
	if event.Request != nil {
		for _, h := range sensitiveHeaders {
			delete(event.Request.Headers, h)
		}
		// A query string can carry tokens or personal data; the path alone is
		// enough to identify the endpoint that failed.
		event.Request.QueryString = ""
		event.Request.Cookies = ""
		event.Request.Data = ""
	}
	return event
}

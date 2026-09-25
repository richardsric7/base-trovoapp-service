package observe

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestID assigns a correlation ID to every request.
//
// An inbound X-Request-ID is honoured so an ID assigned upstream (a gateway,
// or another Trovo service calling this one) stays constant across hops -
// that is what makes a single ID traceable across the platform. Otherwise a
// new one is generated. The ID goes on the context, into every log line for
// the request, and back to the caller in the response header.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := strings.TrimSpace(c.GetHeader(RequestIDHeader))
		// Only accept an upstream ID that looks like one. An arbitrary
		// attacker-supplied string would otherwise end up in every log line.
		if id == "" || len(id) > 64 || strings.ContainsAny(id, "\r\n \t") {
			id = uuid.NewString()
		}
		c.Set(CtxRequestID, id)
		c.Header(RequestIDHeader, id)
		c.Next()
	}
}

// RequestIDFrom returns the correlation ID for the request, or "" if the
// RequestID middleware did not run.
func RequestIDFrom(c *gin.Context) string {
	if v, ok := c.Get(CtxRequestID); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// logLine is one structured request log record. JSON keys are chosen so Loki
// queries filter on fields (`| json | status >= 500`) rather than matching
// substrings of a formatted string.
type logLine struct {
	Time       string `json:"time"`
	Level      string `json:"level"`
	Service    string `json:"service"`
	Msg        string `json:"msg"`
	RequestID  string `json:"request_id,omitempty"`
	Method     string `json:"method"`
	Path       string `json:"path"`
	Status     int    `json:"status"`
	DurationMs int64  `json:"duration_ms"`
	ClientIP   string `json:"client_ip,omitempty"`
	UserAgent  string `json:"user_agent,omitempty"`
	Error      string `json:"error,omitempty"`
}

// Logger writes one structured line per request.
//
// It replaces gin's default formatter rather than adding to it, so each
// request produces exactly one line. Health and metrics endpoints are skipped:
// they are polled every few seconds and would otherwise dominate the logs and
// the retention budget.
func Logger() gin.HandlerFunc {
	skip := map[string]bool{
		"/health": true, "/ready": true, "/metrics": true,
		"/api/v1/health": true, "/api/v1/ready": true,
	}
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		if skip[c.Request.URL.Path] {
			return
		}

		status := c.Writer.Status()
		level := "info"
		switch {
		case status >= 500:
			level = "error"
			// Report it as well as logging it. 4xx are deliberately excluded:
			// they are expected outcomes (validation, permissions, an expired
			// session) and would drown the real failures.
			reportServerError(c, status)
		case status >= 400:
			level = "warn"
		}

		line := logLine{
			Time:       start.UTC().Format(time.RFC3339Nano),
			Level:      level,
			Service:    serviceName,
			Msg:        "request",
			RequestID:  RequestIDFrom(c),
			Method:     c.Request.Method,
			Path:       c.Request.URL.Path,
			Status:     status,
			DurationMs: time.Since(start).Milliseconds(),
			ClientIP:   c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
			Error:      c.Errors.ByType(gin.ErrorTypePrivate).String(),
		}

		encoded, err := json.Marshal(line)
		if err != nil {
			// Never lose the record because it would not encode.
			fmt.Fprintf(os.Stdout, `{"level":"error","service":%q,"msg":"log encode failed","path":%q}`+"\n", serviceName, line.Path)
			return
		}
		// Resolved at write time rather than captured when the middleware is
		// built, so tests can redirect stdout.
		fmt.Fprintln(os.Stdout, string(encoded))
	}
}

// Recovery converts a panic into a 500 and reports it with request context.
//
// It replaces gin.Recovery(): a panic today writes a stack trace to the
// container log where nobody sees it. Here it is reported to GlitchTip
// (grouped and alertable) and always logged locally with its request ID, so
// the crash can be tied back to the request that caused it.
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			r := recover()
			if r == nil {
				return
			}

			requestID := RequestIDFrom(c)
			RecordPanic(c.FullPath())
			log.Printf("[observe] PANIC request_id=%s %s %s: %v", requestID, c.Request.Method, c.Request.URL.Path, r)

			if initialised {
				hub := sentry.CurrentHub().Clone()
				hub.Scope().SetRequest(c.Request)
				hub.Scope().SetTag("request_id", requestID)
				hub.Scope().SetTag("endpoint", c.Request.Method+" "+c.FullPath())
				hub.Recover(r)
				// Flush here rather than relying on shutdown: a panicking
				// process may not get a clean shutdown.
				hub.Flush(2 * time.Second)
			}

			// The request ID is returned so a user reporting the error can
			// quote it and have it found immediately in the logs.
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error":      "internal server error",
				"request_id": requestID,
			})
		}()
		c.Next()
	}
}

// reportServerError files a crash report for a response that left with a 5xx
// without panicking - a handler that hit an error and returned
// c.JSON(500, ...) deliberately.
//
// Recovery() above only fires on panics, so until this existed those 500s
// reached GlitchTip from nowhere: the browser reported them (the frontend SDK
// captures 5xx responses), but the backend - which has the actual error, the
// SQL, and the request ID - stayed silent. A GET /country/list returning 500
// on every call for a schema mismatch was found from a browser breadcrumb,
// not from the service that produced it.
//
// Grouped by method and route template, not by URL or message: one issue per
// broken endpoint with a count, rather than one per request.
func reportServerError(c *gin.Context, status int) {
	if !initialised {
		return
	}
	route := c.FullPath()
	if route == "" {
		// Unmatched routes would otherwise create an issue per distinct URL.
		route = "unmatched"
	}

	hub := sentry.CurrentHub().Clone()
	hub.Scope().SetRequest(c.Request)
	hub.Scope().SetTag("request_id", RequestIDFrom(c))
	hub.Scope().SetTag("endpoint", c.Request.Method+" "+route)
	hub.Scope().SetTag("failure_kind", "server_error")
	hub.Scope().SetLevel(sentry.LevelError)
	// Gin collects handler errors in c.Errors even when the handler writes the
	// response itself; attach them so the issue carries the real cause rather
	// than just a status code.
	if errs := c.Errors.Errors(); len(errs) > 0 {
		hub.Scope().SetExtra("errors", errs)
	}
	hub.CaptureException(fmt.Errorf("%d on %s %s", status, c.Request.Method, route))
}

// CaptureError reports a handled error that deserves attention but did not
// panic. Handlers can call it directly; it no-ops when reporting is disabled.
func CaptureError(c *gin.Context, err error) {
	if err == nil {
		return
	}
	requestID := RequestIDFrom(c)
	log.Printf("[observe] error request_id=%s %s %s: %v", requestID, c.Request.Method, c.Request.URL.Path, err)

	if !initialised {
		return
	}
	hub := sentry.CurrentHub().Clone()
	hub.Scope().SetRequest(c.Request)
	hub.Scope().SetTag("request_id", requestID)
	hub.Scope().SetTag("endpoint", fmt.Sprintf("%s %s", c.Request.Method, c.FullPath()))
	hub.CaptureException(err)
}

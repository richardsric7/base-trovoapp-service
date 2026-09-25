package accesslog

import (
	"bytes"
	"encoding/json"
	"io"

	"admin-panel-dashboard/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ctxAuditBody is the context key under which Audit stashes a snapshot of the
// request body, so a BodyField extractor can read it after the handler has already
// consumed c.Request.Body.
const ctxAuditBody = "accesslog_audit_body"

// Audit returns gin middleware that records one access-log event for the request
// AFTER the handler runs, deriving success/failure from the response status. It is
// applied per privileged route, so only actions worth auditing are recorded.
//
// The actor, IP, user-agent, method and path are read from the context/request;
// the caller supplies the event name and category, and optionally a targetFn that
// extracts the target of the action for the audit detail. If targetFn is a
// BodyField, the request body is buffered before the handler runs and restored, so
// the handler still reads it normally.
//
// Because it keys off the final status code, it correctly records a failed action
// (e.g. a forbidden suspend attempt) as status=failed without any handler changes.
func Audit(db *gorm.DB, event, category string, targetFn func(*gin.Context) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Snapshot the body up front so a BodyField extractor can read it after the
		// handler consumes c.Request.Body. Only done when there is a body and an
		// extractor that might need it; the body is restored for the handler.
		if targetFn != nil && c.Request.Body != nil {
			if raw, err := io.ReadAll(c.Request.Body); err == nil {
				c.Request.Body = io.NopCloser(bytes.NewReader(raw))
				c.Set(ctxAuditBody, raw)
			}
		}

		c.Next()

		status := models.AccessStatusSuccessful
		if c.Writer.Status() >= 400 {
			status = models.AccessStatusFailed
		}

		var target string
		if targetFn != nil {
			// Guard the extractor: a malformed/absent field must never panic the
			// request (the handler has already run and responded at this point).
			func() {
				defer func() { _ = recover() }()
				target = targetFn(c)
			}()
		}

		RecordFromContext(db, c, Event{
			Event:    event,
			Status:   status,
			Category: category,
			Target:   target,
		})
	}
}

// Param builds a targetFn that reads a single URL path parameter.
func Param(name string) func(*gin.Context) string {
	return func(c *gin.Context) string { return c.Param(name) }
}

// BodyField builds a targetFn that reads a top-level string field from the JSON
// request body (buffered by Audit before the handler ran). Missing field or
// non-JSON body yields "".
func BodyField(field string) func(*gin.Context) string {
	return func(c *gin.Context) string {
		v, ok := c.Get(ctxAuditBody)
		if !ok {
			return ""
		}
		raw, ok := v.([]byte)
		if !ok || len(raw) == 0 {
			return ""
		}
		var m map[string]any
		if err := json.Unmarshal(raw, &m); err != nil {
			return ""
		}
		if s, ok := m[field].(string); ok {
			return s
		}
		return ""
	}
}

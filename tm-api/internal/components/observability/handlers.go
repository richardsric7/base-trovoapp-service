package observability

import (
	"net/http"
	"strconv"
	"strings"

	"admin-panel-dashboard/internal/middleware"
	serverModels "admin-panel-dashboard/internal/server/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Guard is the single authorisation point for every observability endpoint.
//
// Today it requires nothing beyond a valid admin session, matching the rest of
// the dashboard. It exists as its own function so that adding the planned
// "Observability: view" / "Observability: manage" permissions is a change to
// this one place, rather than to six route registrations - and so it is
// obvious where that check belongs.
func Guard(db *gorm.DB) gin.HandlerFunc {
	return middleware.JwtTokenAuthMiddleware(db)
}

// HealthHandler godoc
//
//	@Summary		Platform health
//	@Description	Live status of every backend service and the host they run on.
//	@Description	Each service reports its own dependencies (databases, cache,
//	@Description	blockchain node) with latency, so a degraded dependency is
//	@Description	visible before it becomes an outage. Polled by the Trovo
//	@Description	Manager System Health screen.
//	@Tags			Observability
//	@Produce		json
//	@Security		JwtTokenAuth
//	@Param			Authorization	header		string	true	"JWT Token"	default(Bearer <your-token>)
//	@Success		200	{object}	HealthResponse
//	@Failure		401	{object}	object	"Unauthorized"
//	@Router			/observability/health [get]
func HealthHandler(cfg Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		out, _ := cached("health", func() (HealthResponse, error) {
			return cfg.Health(c.Request.Context()), nil
		})
		c.JSON(http.StatusOK, gin.H{"data": out})
	}
}

// IssuesHandler godoc
//
//	@Summary		Application errors
//	@Description	Grouped crash reports from every Trovo application (admin,
//	@Description	wallet and P2P, backend and web). Each entry is one distinct
//	@Description	error with how many times it has occurred and when it was
//	@Description	first and last seen.
//	@Tags			Observability
//	@Produce		json
//	@Security		JwtTokenAuth
//	@Param			Authorization	header		string	true	"JWT Token"	default(Bearer <your-token>)
//	@Param			limit			query		int		false	"Maximum issues to return (default 50, max 100)"
//	@Param			app				query		string	false	"Filter to one application (project slug)"
//	@Param			query			query		string	false	"Crash reporter search, e.g. is:unresolved"
//	@Success		200	{object}	IssuesResponse
//	@Failure		401	{object}	object	"Unauthorized"
//	@Router			/observability/issues [get]
func IssuesHandler(cfg Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit := clampInt(c.Query("limit"), 50, 1, 100)
		app := strings.TrimSpace(c.Query("app"))
		query := strings.TrimSpace(c.Query("query"))

		out, _ := cached("issues|"+app+"|"+query+"|"+strconv.Itoa(limit), func() (IssuesResponse, error) {
			return cfg.Issues(c.Request.Context(), limit, app, query), nil
		})
		c.JSON(http.StatusOK, gin.H{"data": out})
	}
}

// TrendsHandler godoc
//
//	@Summary		Traffic and performance trends
//	@Description	Request rate, error rate and 95th-percentile response time per
//	@Description	service over the selected window. The 95th percentile is
//	@Description	reported rather than an average because an average hides the
//	@Description	slow requests users actually complain about.
//	@Tags			Observability
//	@Produce		json
//	@Security		JwtTokenAuth
//	@Param			Authorization	header		string	true	"JWT Token"	default(Bearer <your-token>)
//	@Param			window			query		string	false	"Time window"	Enums(1h, 6h, 24h)
//	@Success		200	{object}	TrendsResponse
//	@Failure		401	{object}	object	"Unauthorized"
//	@Router			/observability/trends [get]
func TrendsHandler(cfg Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		window := strings.TrimSpace(c.Query("window"))
		if window == "" {
			window = "6h"
		}
		out, _ := cached("trends|"+window, func() (TrendsResponse, error) {
			return cfg.Trends(c.Request.Context(), window), nil
		})
		c.JSON(http.StatusOK, gin.H{"data": out})
	}
}

// TraceHandler godoc
//
//	@Summary		Trace one request across services
//	@Description	Returns every log line carrying the given request identifier,
//	@Description	from every service that handled it, in time order. Each API
//	@Description	response includes its identifier in the X-Request-ID header,
//	@Description	and the web applications show it on their error screens - so a
//	@Description	user's report can be followed end to end.
//	@Tags			Observability
//	@Produce		json
//	@Security		JwtTokenAuth
//	@Param			Authorization	header		string	true	"JWT Token"	default(Bearer <your-token>)
//	@Param			request_id		path		string	true	"Request identifier from a response header or error screen"
//	@Success		200	{object}	TraceResponse
//	@Failure		400	{object}	object	"Malformed request identifier"
//	@Failure		401	{object}	object	"Unauthorized"
//	@Router			/observability/trace/{request_id} [get]
func TraceHandler(cfg Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := strings.TrimSpace(c.Param("request_id"))
		// Rejected rather than escaped: this is the only user-supplied value
		// that reaches the log query, and LogQL treats quotes as syntax.
		if !ValidRequestID(requestID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request identifier"})
			return
		}
		out := cfg.Trace(c.Request.Context(), requestID)
		c.JSON(http.StatusOK, gin.H{"data": out})
	}
}

// RecentErrorsHandler godoc
//
//	@Summary		Recent failed requests
//	@Description	The most recent requests that returned an error, across every
//	@Description	service. Gives the trace screen a starting point for an admin
//	@Description	who does not already have a request identifier to hand.
//	@Tags			Observability
//	@Produce		json
//	@Security		JwtTokenAuth
//	@Param			Authorization	header		string	true	"JWT Token"	default(Bearer <your-token>)
//	@Param			limit			query		int		false	"Maximum lines to return (default 50, max 200)"
//	@Success		200	{object}	RecentErrorsResponse
//	@Failure		401	{object}	object	"Unauthorized"
//	@Router			/observability/recent-errors [get]
func RecentErrorsHandler(cfg Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit := clampInt(c.Query("limit"), 50, 1, 200)
		out, _ := cached("recent-errors|"+strconv.Itoa(limit), func() (RecentErrorsResponse, error) {
			return cfg.RecentErrors(c.Request.Context(), limit), nil
		})
		c.JSON(http.StatusOK, gin.H{"data": out})
	}
}

// CapabilitiesResponse tells the UI which screens have a backing service, so it
// can show a clear "not configured" state instead of an empty chart.
type CapabilitiesResponse struct {
	Health       bool   `json:"health"`
	Issues       bool   `json:"issues"`
	Trends       bool   `json:"trends"`
	Trace        bool   `json:"trace"`
	GrafanaURL   string `json:"grafana_url,omitempty"`
	GlitchTipURL string `json:"glitchtip_url,omitempty"`
}

// CapabilitiesHandler godoc
//
//	@Summary		Available observability features
//	@Description	Which System Health screens this environment can serve. A
//	@Description	deployment without log aggregation, for example, reports
//	@Description	trace as unavailable so the interface can say so plainly
//	@Description	rather than showing an empty result.
//	@Tags			Observability
//	@Produce		json
//	@Security		JwtTokenAuth
//	@Param			Authorization	header		string	true	"JWT Token"	default(Bearer <your-token>)
//	@Success		200	{object}	CapabilitiesResponse
//	@Failure		401	{object}	object	"Unauthorized"
//	@Router			/observability/capabilities [get]
func CapabilitiesHandler(cfg Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": CapabilitiesResponse{
			Health:       len(cfg.Services) > 0,
			Issues:       cfg.HasGlitchTip(),
			Trends:       cfg.HasPrometheus(),
			Trace:        cfg.HasLoki(),
			GrafanaURL:   cfg.GrafanaURL,
			GlitchTipURL: cfg.PublicGlitchTipURL(),
		}})
	}
}

// Init registers the module's routes.
func Init(apiV1 *gin.RouterGroup, s *serverModels.Server) {
	cfg := LoadConfig()
	guard := Guard(s.AdminDB)

	group := apiV1.Group("/observability", guard)
	group.GET("/capabilities", CapabilitiesHandler(cfg))
	group.GET("/health", HealthHandler(cfg))
	group.GET("/issues", IssuesHandler(cfg))
	group.GET("/trends", TrendsHandler(cfg))
	group.GET("/trace/:request_id", TraceHandler(cfg))
	group.GET("/recent-errors", RecentErrorsHandler(cfg))
}

// clampInt parses a query parameter, falling back to a default and bounding the
// result so a caller cannot request an unbounded page.
func clampInt(raw string, def, min, max int) int {
	n, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil {
		return def
	}
	if n < min {
		return min
	}
	if n > max {
		return max
	}
	return n
}

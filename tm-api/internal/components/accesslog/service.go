// Package accesslog writes the admin security / access audit trail
// (models.AdminAccessLog). It is intentionally tiny and side-effect-only: a write
// failure is logged but never propagates, so auditing can never break the action
// being audited -- the same guarantee the vault-signer audit helper makes.
package accesslog

import (
	"log"
	"strings"
	"time"

	"admin-panel-dashboard/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Actor identifies who performed an action. For login events it is built from the
// submitted username (an admin row may not even be resolvable on failure); for
// authenticated actions it is read off the gin context by FromContext.
type Actor struct {
	AdminID  *uint
	Username string
	FullName string
	Email    string
	Role     string
}

// Event is the payload passed to Record. Only Event and Status are required; the
// writer fills OccurredAt and enriches Location from IPAddress.
type Event struct {
	Event       string
	Status      string
	Category    string
	Actor       Actor
	Target      string
	IPAddress   string
	UserAgent   string
	Method      string
	Path        string
	LoginMethod string
	Detail      string
}

// Record persists an access-log event. It returns immediately and does the insert
// (plus the geo lookup, which is a network call) on a background goroutine so it
// adds no latency to the request and can never fail it. db must be the admin DB.
func Record(db *gorm.DB, ev Event) {
	if db == nil || ev.Event == "" {
		return
	}
	if ev.Status == "" {
		ev.Status = models.AccessStatusSuccessful
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[accesslog] recovered from panic while recording %q: %v", ev.Event, r)
			}
		}()
		recordSync(db, ev)
	}()
}

// recordSync builds and inserts the row on the calling goroutine. Record wraps it
// in a goroutine; tests call it directly for determinism.
func recordSync(db *gorm.DB, ev Event) {
	row := models.AdminAccessLog{
		ID:           uuid.NewString(),
		Event:        ev.Event,
		Status:       ev.Status,
		Category:     ev.Category,
		Action:       humanise(ev.Event),
		ActorAdminID: ev.Actor.AdminID,
		Username:     ev.Actor.Username,
		FullName:     ev.Actor.FullName,
		Email:        ev.Actor.Email,
		Role:         ev.Actor.Role,
		Target:       ev.Target,
		IPAddress:    ev.IPAddress,
		Location:     lookupLocation(ev.IPAddress),
		UserAgent:    ev.UserAgent,
		Method:       ev.Method,
		Path:         ev.Path,
		LoginMethod:  ev.LoginMethod,
		Detail:       ev.Detail,
		OccurredAt:   time.Now(),
	}

	if err := db.Create(&row).Error; err != nil {
		// Never fatal -- auditing must not break the audited action.
		log.Printf("[accesslog] failed to record %q for %q: %v", ev.Event, ev.Actor.Username, err)
	}
}

// RecordFromContext is the convenience entry point for the audit middleware and
// authenticated handlers: it derives the actor, IP, UA, method and path from the
// gin context, then records the event.
func RecordFromContext(db *gorm.DB, c *gin.Context, ev Event) {
	if ev.Actor == (Actor{}) {
		ev.Actor = FromContext(c)
	}
	if ev.IPAddress == "" {
		ev.IPAddress = clientIP(c)
	}
	if ev.UserAgent == "" {
		ev.UserAgent = c.Request.UserAgent()
	}
	if ev.Method == "" {
		ev.Method = c.Request.Method
	}
	if ev.Path == "" {
		ev.Path = c.FullPath()
	}
	Record(db, ev)
}

// lookupLocation turns an IP into "City, Country" using the existing IPAPI helper.
// Best-effort: any failure yields an empty string.
func lookupLocation(ip string) string {
	if ip == "" {
		return ""
	}
	geo, err := models.GetGeoInfo(ip)
	if err != nil {
		return ""
	}
	switch {
	case geo.City != "" && geo.Country != "":
		return geo.City + ", " + geo.Country
	case geo.Country != "":
		return geo.Country
	default:
		return ""
	}
}

// clientIP mirrors the auth flow's convention of substituting a routable IP when
// the request originates from localhost, so dev/local logins still geo-resolve.
func clientIP(c *gin.Context) string {
	ip := c.ClientIP()
	if ip == "127.0.0.1" || ip == "::1" {
		return "41.190.2.194"
	}
	return ip
}

// humanise turns an event key like "admin.suspend" into the label the UI renders
// in its "Action" column, e.g. "Admin suspend". "login.success" -> "Login".
func humanise(event string) string {
	switch event {
	case models.EventLoginSuccess:
		return "Login"
	case models.EventLoginFailure, models.EventLoginSuspended:
		return "Login failed"
	case models.EventLogout:
		return "Logout"
	case models.EventTokenRefresh:
		return "Token refresh"
	}
	// Generic: "admin.status_change" -> "Admin status change".
	s := strings.NewReplacer(".", " ", "_", " ").Replace(event)
	s = strings.TrimSpace(s)
	if s == "" {
		return event
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

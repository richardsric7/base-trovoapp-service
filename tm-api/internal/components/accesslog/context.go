package accesslog

import "github.com/gin-gonic/gin"

// Context keys under which the authenticating middleware stores the resolved
// admin actor, so downstream audit hooks can record who acted. Defined here (not
// in the middleware package) so accesslog has no dependency on middleware and no
// import cycle can form.
const (
	CtxAdminID   = "accesslog_admin_id"
	CtxUsername  = "accesslog_admin_username"
	CtxFullName  = "accesslog_admin_fullname"
	CtxEmail     = "accesslog_admin_email"
	CtxRole      = "accesslog_admin_role"
)

// SetActor is called by the authenticating middleware after it resolves the admin,
// recording the actor identity on the context for later audit hooks.
func SetActor(c *gin.Context, adminID uint, username, fullName, email, role string) {
	c.Set(CtxAdminID, adminID)
	c.Set(CtxUsername, username)
	c.Set(CtxFullName, fullName)
	c.Set(CtxEmail, email)
	c.Set(CtxRole, role)
}

// FromContext reconstructs the Actor from whatever the middleware set. It is
// tolerant of a missing actor (e.g. an unauthenticated route) -- the returned
// Actor is simply sparse in that case.
func FromContext(c *gin.Context) Actor {
	var a Actor
	if v, ok := c.Get(CtxAdminID); ok {
		if id, ok := v.(uint); ok && id != 0 {
			a.AdminID = &id
		}
	}
	a.Username = getString(c, CtxUsername)
	a.FullName = getString(c, CtxFullName)
	a.Email = getString(c, CtxEmail)
	a.Role = getString(c, CtxRole)
	return a
}

func getString(c *gin.Context, key string) string {
	if v, ok := c.Get(key); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

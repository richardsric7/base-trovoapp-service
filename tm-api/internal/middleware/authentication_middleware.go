package middleware

import (
	"admin-panel-dashboard/internal/components/accesslog"
	"admin-panel-dashboard/internal/trovosdk"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

type AdminUser struct {
	ID uint `gorm:"primaryKey"`

	Username  string `gorm:"unique;not null"`
	Email     string `gorm:"unique;not null"`
	FirstName string
	LastName  string

	Status       string `gorm:"not null"`
	Role         Role   `gorm:"not null"`
	IsAdmin      bool   `gorm:"not null"`
	WalletUserID string `gorm:"size:100"`
}

// fullName is the actor display name recorded on audit events.
func (a AdminUser) fullName() string {
	name := a.FirstName
	if a.LastName != "" {
		if name != "" {
			name += " "
		}
		name += a.LastName
	}
	return name
}
type Role string
type AdminStatus string

const (
	SuperAdmin     Role = "SUPER_ADMIN"
	EditLevelAdmin Role = "EDIT_LEVEL_ADMIN"
	ViewOnlyAdmin  Role = "VIEW_ONLY_ADMIN"
	Suspended      Role = "SUSPENDED"
)

const (
	AdminContextUserID = "admin_user_id"
	AdminContextRole   = "admin_role"
)

// JwtTokenAuthMiddleware is middleware that checks if Token is valid //general middleware
func JwtTokenAuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sl, err := trovosdk.NewServiceLink()
		if err != nil {
			log.Println("[JwtTokenAuthMiddleware] error:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		jwtResponse, err := sl.JwtTokenVerify(ExtractToken(c.Request))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		var existingAdmin AdminUser
		if err = db.First(&existingAdmin, "wallet_user_id = ?", jwtResponse.Token.Claims.UserID).Error; err != nil {
			log.Println("wallet_user_id: ", jwtResponse.Token.Claims.UserID)
			msg := fmt.Sprintf("unauthorized user %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
			c.Abort()
			return
		}

		// Record the resolved admin so downstream audit hooks can attribute actions.
		accesslog.SetActor(c, existingAdmin.ID, existingAdmin.Username, existingAdmin.fullName(), existingAdmin.Email, string(existingAdmin.Role))

		c.Next()

	}
}

func AuthenticateSuperAdmin(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Initialize the Trovo SDK ServiceLink
		sl, err := trovosdk.NewServiceLink()
		if err != nil {
			log.Println("[JwtTokenAuthMiddleware] error:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Verify the JWT token
		jwtResponse, err := sl.JwtTokenVerify(ExtractToken(c.Request))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		var existingAdmin AdminUser
		if err = db.First(&existingAdmin, "wallet_user_id = ?", jwtResponse.Token.Claims.UserID).Error; err != nil {
			msg := "unauthorized user"
			c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
			c.Abort()
			return
		}

		// Check if the user's role is superAdmin
		if existingAdmin.Role != SuperAdmin {
			msg := "forbidden: insufficient permissions"
			c.JSON(http.StatusForbidden, gin.H{"error": msg})
			c.Abort()
			return
		}

		c.Set(AdminContextUserID, fmt.Sprint(existingAdmin.ID))
		c.Set(AdminContextRole, string(existingAdmin.Role))

		// Record the resolved admin so downstream audit hooks can attribute actions.
		accesslog.SetActor(c, existingAdmin.ID, existingAdmin.Username, existingAdmin.fullName(), existingAdmin.Email, string(existingAdmin.Role))

		// Proceed with the request if the user is a superAdmin
		c.Next()
	}
}

// AllowOrgOrTrovoAdmin is a composing middleware that grants access if EITHER
// OrganizationAuthMiddleware OR JwtTokenAuthMiddleware succeeds.
func AllowOrgOrTrovoAdmin(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}
		// ExtractToken strips a "Bearer " prefix if present, so both tryOrganizationAuth
		// and tryTrovoAdminAuth below always see the bare token — matching what
		// JwtTokenVerify expects (it re-adds its own "Bearer " prefix before forwarding to
		// the external Trovo service, so a token that already carries one would otherwise
		// go out doubled and fail verification for every caller, regardless of role).
		token := ExtractToken(c.Request)

		// --- Attempt 1: Try to authenticate as Organization Member ---
		// Try to parse the token as OrganizationClaims
		if tryOrganizationAuth(c, db, token) {
			c.Set("auth_type", "organization_member")
			c.Next()
			return
		}

		// --- Attempt 2: Try to authenticate as Trovo Admin ---
		if tryTrovoAdminAuth(c, db, token) {
			c.Set("auth_type", "trovo_admin")
			c.Next()
			return
		}

		// --- Both authentication attempts failed ---
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token or insufficient permissions"})
		c.Abort()
	}
}

// tryOrganizationAuth attempts to authenticate as organization member. token must already
// have any "Bearer " prefix stripped (see ExtractToken) — both call sites below pass the
// extracted token, never the raw header.
func tryOrganizationAuth(c *gin.Context, db *gorm.DB, token string) bool {
	// Define OrganizationClaims locally to avoid import cycles
	type OrganizationClaims struct {
		OrganizationID    string `json:"organization_id"`
		OrganizationEmail string `json:"organization_email"`
		OrganizationType  string `json:"organization_type"`
		MemberID          string `json:"member_id"`
		MemberEmail       string `json:"member_email"`
		MemberRole        string `json:"member_role"`
		OrganizationName  string `json:"organization_name"`
		SessionVersion    uint64 `json:"session_version"`
		jwt.RegisteredClaims
	}

	// Try to parse as OrganizationClaims first
	parsedToken, err := jwt.ParseWithClaims(token, &OrganizationClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err == nil {
		if claims, ok := parsedToken.Claims.(*OrganizationClaims); ok && parsedToken.Valid && claims.MemberID != "" && claims.OrganizationID != "" {
			if validateOrganizationSession(db, claims.MemberID, claims.OrganizationID, claims.MemberEmail, claims.MemberRole, claims.SessionVersion) != nil {
				return false
			}
			// Set organization and member info in context
			c.Set("organization_id", claims.OrganizationID)
			c.Set("organization_email", claims.OrganizationEmail)
			c.Set("organization_type", claims.OrganizationType)
			c.Set("member_id", claims.MemberID)
			c.Set("member_email", claims.MemberEmail)
			c.Set("member_role", claims.MemberRole)
			c.Set("organization_name", claims.OrganizationName)
			return true
		}
	}

	// If OrganizationClaims failed, try to parse as MapClaims (legacy format)
	parsedToken, err = jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err == nil {
		if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
			// Extract fields from legacy token format
			orgID, _ := claims["organization_id"].(string)
			email, _ := claims["email"].(string)
			role, _ := claims["role"].(string)
			userID, _ := claims["user_id"].(string)

			if orgID != "" && email != "" && role != "" && userID != "" {
				if validateOrganizationSession(db, userID, orgID, email, role, 0) != nil {
					return false
				}
				// Set organization and member info in context (mapping legacy fields)
				c.Set("organization_id", orgID)
				c.Set("organization_email", email)
				c.Set("organization_type", "ASSET_MANAGER") // Default for legacy tokens
				c.Set("member_id", userID)
				c.Set("member_email", email)
				c.Set("member_role", role)
				return true
			}
		}
	}

	return false
}

// tryTrovoAdminAuth attempts to authenticate as Trovo admin. token must already have any
// "Bearer " prefix stripped (see ExtractToken) — JwtTokenVerify re-adds its own "Bearer "
// before forwarding to the external Trovo service, so a token that still carried one here
// would go out doubled and fail verification.
func tryTrovoAdminAuth(c *gin.Context, db *gorm.DB, token string) bool {
	// Initialize the Trovo SDK ServiceLink
	sl, err := trovosdk.NewServiceLink()
	if err != nil {
		log.Println("[tryTrovoAdminAuth] Trovo SDK error:", err)
		return false
	}

	// Verify the JWT token
	jwtResponse, err := sl.JwtTokenVerify(token)
	if err != nil {
		log.Println("[tryTrovoAdminAuth] JWT verification error:", err)
		return false
	}

	// Check if user exists as Trovo admin
	var existingAdmin AdminUser
	if err := db.First(&existingAdmin, "wallet_user_id = ?", jwtResponse.Token.Claims.UserID).Error; err != nil {
		log.Println("[tryTrovoAdminAuth] Admin not found for wallet_user_id:", jwtResponse.Token.Claims.UserID)
		return false
	}

	// SUCCESS: Valid Trovo Admin token
	c.Set("trovo_admin_id", existingAdmin.ID)
	c.Set("trovo_admin_email", existingAdmin.Email)
	// Record the resolved admin so downstream audit hooks can attribute actions.
	accesslog.SetActor(c, existingAdmin.ID, existingAdmin.Username, existingAdmin.fullName(), existingAdmin.Email, string(existingAdmin.Role))
	return true
}

// AllowOrgOrTrovoAdminNormalized provides normalized context keys to eliminate handler switch statements
func AllowOrgOrTrovoAdminNormalized(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}
		// ExtractToken strips a "Bearer " prefix if present — see the note on
		// AllowOrgOrTrovoAdmin above for why both helpers below need the bare token.
		token := ExtractToken(c.Request)

		// Try organization member authentication first
		if tryOrganizationAuth(c, db, token) {
			// Normalize context for organization member
			c.Set("user_type", "org_member")
			c.Set("current_user_id", c.GetString("member_id"))
			c.Set("current_user_email", c.GetString("member_email"))
			c.Set("requires_org_param", false)
			c.Set("user_organization_role", c.GetString("member_role"))
			// organization_id, organization_name, organization_type already set by tryOrganizationAuth
			c.Next()
			return
		}

		// Try Trovo admin authentication
		if tryTrovoAdminAuth(c, db, token) {
			// Normalize context for Trovo admin
			c.Set("user_type", "trovo_admin")
			adminID, _ := c.Get("trovo_admin_id")
			c.Set("current_user_id", fmt.Sprintf("%v", adminID))
			c.Set("current_user_email", c.GetString("trovo_admin_email"))
			c.Set("requires_org_param", true)
			c.Set("user_organization_role", "") // No org role for admins
			c.Next()
			return
		}

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token or insufficient permissions"})
		c.Abort()
	}
}

// GetOrganizationIDFromContext gets organization ID from context or request based on user type
func GetOrganizationIDFromContext(c *gin.Context) (string, error) {
	requiresParam := c.GetBool("requires_org_param")

	if requiresParam {
		// Trovo admin - get from query parameter
		orgID := c.Query("organization_id")
		if orgID == "" {
			// Try from request body for POST requests
			type OrgIDRequest struct {
				OrganizationID string `json:"organization_id"`
			}
			var req OrgIDRequest
			if err := c.ShouldBindJSON(&req); err == nil && req.OrganizationID != "" {
				return req.OrganizationID, nil
			}
			return "", fmt.Errorf("organization_id parameter is required for this action")
		}
		return orgID, nil
	} else {
		// Organization member - get from token context
		orgID := c.GetString("organization_id")
		if orgID == "" {
			return "", fmt.Errorf("organization context not found")
		}
		return orgID, nil
	}
}

// CheckOrganizationPermission checks if the current user has required organization role
func CheckOrganizationPermission(c *gin.Context, requiredRole string) error {
	userType := c.GetString("user_type")

	if userType == "trovo_admin" {
		// Trovo admins bypass organization role checks
		return nil
	}

	if userType == "org_member" {
		userRole := c.GetString("user_organization_role")
		if userRole != requiredRole {
			return fmt.Errorf("forbidden: only users with role %s can perform this action", requiredRole)
		}
		return nil
	}

	return fmt.Errorf("invalid user type")
}

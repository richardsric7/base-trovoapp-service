package middleware

import (
	"admin-panel-dashboard/internal/db"
	"admin-panel-dashboard/internal/models"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

// OrganizationClaims represents the claims in the JWT token for organization members
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

// GenerateOrganizationToken generates a JWT token for an organization member
func GenerateOrganizationToken(member *models.OrganizationMember) (string, error) {
	adminDB, err := db.AdminDB()
	if err != nil {
		return "", err
	}
	return GenerateOrganizationTokenWithDB(adminDB, member)
}

func GenerateOrganizationTokenWithDB(adminDB *gorm.DB, member *models.OrganizationMember) (string, error) {
	// Get organization details
	var org models.Organization
	if err := adminDB.First(&org, "id = ?", member.OrganizationID).Error; err != nil {
		return "", err
	}

	claims := &OrganizationClaims{
		OrganizationID:    org.ID,
		OrganizationEmail: org.Email,
		OrganizationType:  org.Type,
		MemberID:          member.ID,
		MemberEmail:       member.Email,
		MemberRole:        member.Role,
		OrganizationName:  org.Name,
		SessionVersion:    member.SessionVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// OrganizationAuthMiddleware is a middleware to authenticate organization members
func OrganizationAuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Authorization header is required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimSpace(authHeader)
		parts := strings.Fields(tokenString)
		if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
			tokenString = parts[1]
		}

		// Try to parse as OrganizationClaims first
		token, err := jwt.ParseWithClaims(tokenString, &OrganizationClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err == nil {
			if claims, ok := token.Claims.(*OrganizationClaims); ok && token.Valid && claims.MemberID != "" && claims.OrganizationID != "" {
				if err := validateOrganizationSession(db, claims.MemberID, claims.OrganizationID, claims.MemberEmail, claims.MemberRole, claims.SessionVersion); err != nil {
					c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Organization session is no longer valid"})
					c.Abort()
					return
				}

				// Set organization and member info in context
				c.Set("organization_id", claims.OrganizationID)
				c.Set("organization_email", claims.OrganizationEmail)
				c.Set("organization_type", claims.OrganizationType)
				c.Set("member_id", claims.MemberID)
				c.Set("member_email", claims.MemberEmail)
				c.Set("member_role", claims.MemberRole)
				c.Set("organization_name", claims.OrganizationName)
				c.Next()
				return
			}
		}

		// If OrganizationClaims failed, try to parse as MapClaims (legacy format)
		token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Extract fields from legacy token format
			orgID, _ := claims["organization_id"].(string)
			email, _ := claims["email"].(string)
			role, _ := claims["role"].(string)
			userID, _ := claims["user_id"].(string)

			if orgID == "" || email == "" || role == "" || userID == "" {
				c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Invalid token claims"})
				c.Abort()
				return
			}
			if err := validateOrganizationSession(db, userID, orgID, email, role, 0); err != nil {
				c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Organization session is no longer valid"})
				c.Abort()
				return
			}

			// Set organization and member info in context (mapping legacy fields)
			c.Set("organization_id", orgID)
			c.Set("organization_email", email)
			c.Set("organization_type", "ASSET_MANAGER") // Default for legacy tokens
			c.Set("member_id", userID)
			c.Set("member_email", email)
			c.Set("member_role", role)

			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Invalid token claims"})
			c.Abort()
			return
		}
	}
}

// RequireActiveOrganization is composed after OrganizationAuthMiddleware for
// organization-owned resources that must not remain accessible while the
// organization is pending, inactive, suspended, or deleted.
func RequireActiveOrganization(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		organizationID := strings.TrimSpace(c.GetString("organization_id"))
		if organizationID == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Organization context is missing"})
			c.Abort()
			return
		}
		var organization models.Organization
		if err := db.Select("id", "status").First(&organization, "id = ?", organizationID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Organization not found"})
			} else {
				c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Database error"})
			}
			c.Abort()
			return
		}
		if organization.Status != models.OrganizationStatusActive {
			c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Organization is not active"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// OrganizationMemberAuthMiddleware is a middleware for organization member authentication
func OrganizationMemberAuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Authorization header is required"})
			c.Abort()
			return
		}

		// Extract token from Bearer
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims := &OrganizationClaims{}

		// Parse and validate the token
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Invalid token"})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Invalid token"})
			c.Abort()
			return
		}

		// Verify member exists and is active
		var member models.OrganizationMember
		if err := db.First(&member, "id = ?", claims.MemberID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Member not found"})
			} else {
				c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Database error"})
			}
			c.Abort()
			return
		}

		if member.Status != "ACTIVE" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Member account is not active"})
			c.Abort()
			return
		}
		if member.SessionVersion != claims.SessionVersion {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Organization session is no longer valid"})
			c.Abort()
			return
		}

		// Verify organization exists and is active
		var org models.Organization
		if err := db.First(&org, "id = ?", claims.OrganizationID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Organization not found"})
			} else {
				c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Database error"})
			}
			c.Abort()
			return
		}

		if org.Status != "ACTIVE" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Organization is not active"})
			c.Abort()
			return
		}

		// Verify token claims match member and organization
		if member.Email != claims.MemberEmail || member.Role != claims.MemberRole || member.OrganizationID != claims.OrganizationID {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Token claims do not match member"})
			c.Abort()
			return
		}

		// Set member and organization info in context
		c.Set("member_id", claims.MemberID)
		c.Set("member_email", claims.MemberEmail)
		c.Set("member_role", claims.MemberRole)
		c.Set("organization_id", claims.OrganizationID)

		c.Next()
	}
}

func validateOrganizationSession(db *gorm.DB, memberID, organizationID, email, role string, sessionVersion uint64) error {
	if db == nil || memberID == "" || organizationID == "" {
		return errors.New("invalid organization session")
	}
	var member models.OrganizationMember
	if err := db.Select("id", "organization_id", "email", "role", "status", "session_version").First(&member, "id = ?", memberID).Error; err != nil {
		return err
	}
	if member.OrganizationID != organizationID || member.Email != email || member.Role != role ||
		member.Status != models.OrganizationStatusActive || member.SessionVersion != sessionVersion {
		return errors.New("invalid organization session")
	}
	return nil
}

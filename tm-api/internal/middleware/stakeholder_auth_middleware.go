package middleware

import (
	stakeholderModels "admin-panel-dashboard/internal/components/stakeholder/models"
	"admin-panel-dashboard/internal/models"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

const (
	StakeholderContextMemberID            = "member_id"
	StakeholderContextMemberEmail         = "member_email"
	StakeholderContextMemberRole          = "member_role"
	StakeholderContextOrganizationID      = "organization_id"
	StakeholderContextOrganizationName    = "organization_name"
	StakeholderContextOrganizationEmail   = "organization_email"
	StakeholderContextOrganizationType    = "organization_type"
	StakeholderContextStakeholderID       = "stakeholder_id"
	StakeholderContextStakeholderType     = "stakeholder_type"
	StakeholderContextDashboardRole       = "dashboard_role"
	StakeholderContextTrovoWalletUsername = "trovo_wallet_username"
	StakeholderContextIsWalletLinked      = "is_wallet_linked"
)

func StakeholderAuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, ok := stakeholderTokenFromHeader(c.GetHeader("Authorization"))
		if !ok {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Authorization header is required"})
			c.Abort()
			return
		}

		claims := &OrganizationClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Invalid organization member token"})
			c.Abort()
			return
		}
		if claims.MemberID == "" || claims.OrganizationID == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Invalid organization member token claims"})
			c.Abort()
			return
		}

		var member models.OrganizationMember
		if err := db.First(&member, "id = ?", claims.MemberID).Error; err != nil {
			writeStakeholderAuthLookupError(c, err, "Member not found")
			return
		}
		if member.Status != models.OrganizationStatusActive {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Member account is not active"})
			c.Abort()
			return
		}
		if member.SessionVersion != claims.SessionVersion {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Organization session is no longer valid"})
			c.Abort()
			return
		}

		var org models.Organization
		if err := db.First(&org, "id = ?", claims.OrganizationID).Error; err != nil {
			writeStakeholderAuthLookupError(c, err, "Organization not found")
			return
		}
		if org.Status != models.OrganizationStatusActive {
			c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Organization is not active"})
			c.Abort()
			return
		}
		if member.OrganizationID != org.ID ||
			member.OrganizationID != claims.OrganizationID ||
			member.Email != claims.MemberEmail ||
			member.Role != claims.MemberRole ||
			(claims.OrganizationEmail != "" && org.Email != claims.OrganizationEmail) {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Token claims do not match active organization member"})
			c.Abort()
			return
		}

		if org.StakeholderID == nil || org.StakeholderType == nil || strings.TrimSpace(*org.StakeholderType) == "" {
			c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Organization is not linked to a supported portal stakeholder"})
			c.Abort()
			return
		}
		dashboardRole, supported := stakeholderModels.DashboardRoleForStakeholderType(*org.StakeholderType)
		if !supported {
			c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Organization stakeholder type is not supported by the portal"})
			c.Abort()
			return
		}

		c.Set(StakeholderContextMemberID, member.ID)
		c.Set(StakeholderContextMemberEmail, member.Email)
		c.Set(StakeholderContextMemberRole, member.Role)
		c.Set(StakeholderContextOrganizationID, org.ID)
		c.Set(StakeholderContextOrganizationName, org.Name)
		c.Set(StakeholderContextOrganizationEmail, org.Email)
		c.Set(StakeholderContextOrganizationType, org.Type)
		c.Set(StakeholderContextStakeholderID, *org.StakeholderID)
		c.Set(StakeholderContextStakeholderType, *org.StakeholderType)
		c.Set(StakeholderContextDashboardRole, dashboardRole)
		c.Set(StakeholderContextTrovoWalletUsername, member.TrovoWalletUsername)
		c.Set(StakeholderContextIsWalletLinked, member.IsWalletLinked)
		c.Next()
	}
}

func RequireStakeholderRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[stakeholderModels.Normalize(role)] = struct{}{}
	}
	return func(c *gin.Context) {
		role := stakeholderModels.Normalize(c.GetString(StakeholderContextDashboardRole))
		if role == "" {
			c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Stakeholder role is required"})
			c.Abort()
			return
		}
		if _, ok := allowed[role]; !ok {
			c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Stakeholder role is not allowed for this route"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func RequireOrganizationAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString(StakeholderContextMemberRole) != string(models.OrganizationSuperAdmin) {
			c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Organization admin role is required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func stakeholderTokenFromHeader(authHeader string) (string, bool) {
	authHeader = strings.TrimSpace(authHeader)
	if authHeader == "" {
		return "", false
	}
	parts := strings.Fields(authHeader)
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return parts[1], true
	}
	if len(parts) == 1 {
		return parts[0], true
	}
	return "", false
}

func writeStakeholderAuthLookupError(c *gin.Context, err error, notFoundMessage string) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: notFoundMessage})
	} else {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Database error"})
	}
	c.Abort()
}

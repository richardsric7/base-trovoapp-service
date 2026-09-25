package services

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"admin-panel-dashboard/internal/components/stakeholder/models"
	"admin-panel-dashboard/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type HTTPError struct {
	Code    int
	Message string
}

func (e *HTTPError) Error() string {
	return e.Message
}

func NewHTTPError(code int, message string) *HTTPError {
	return &HTTPError{Code: code, Message: message}
}

func ErrorStatus(err error) int {
	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		return httpErr.Code
	}
	return http.StatusInternalServerError
}

func ErrorMessage(err error) string {
	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		return httpErr.Message
	}
	if err == nil {
		return ""
	}
	return "internal server error"
}

type AuthContext struct {
	MemberID            string
	MemberEmail         string
	MemberRole          string
	OrganizationID      string
	OrganizationName    string
	OrganizationEmail   string
	OrganizationType    string
	StakeholderID       uint64
	StakeholderType     string
	DashboardRole       string
	TrovoWalletUsername *string
	IsWalletLinked      bool
}

func AuthContextFromGin(c *gin.Context) (AuthContext, error) {
	stakeholderID, ok := c.Get(middleware.StakeholderContextStakeholderID)
	if !ok {
		return AuthContext{}, NewHTTPError(http.StatusUnauthorized, "stakeholder context is missing")
	}
	id, ok := stakeholderID.(uint64)
	if !ok {
		return AuthContext{}, NewHTTPError(http.StatusUnauthorized, "stakeholder context is invalid")
	}
	var walletUsername *string
	if value, exists := c.Get(middleware.StakeholderContextTrovoWalletUsername); exists {
		if typed, ok := value.(*string); ok {
			walletUsername = typed
		}
	}
	return AuthContext{
		MemberID:            c.GetString(middleware.StakeholderContextMemberID),
		MemberEmail:         c.GetString(middleware.StakeholderContextMemberEmail),
		MemberRole:          c.GetString(middleware.StakeholderContextMemberRole),
		OrganizationID:      c.GetString(middleware.StakeholderContextOrganizationID),
		OrganizationName:    c.GetString(middleware.StakeholderContextOrganizationName),
		OrganizationEmail:   c.GetString(middleware.StakeholderContextOrganizationEmail),
		OrganizationType:    c.GetString(middleware.StakeholderContextOrganizationType),
		StakeholderID:       id,
		StakeholderType:     c.GetString(middleware.StakeholderContextStakeholderType),
		DashboardRole:       c.GetString(middleware.StakeholderContextDashboardRole),
		TrovoWalletUsername: walletUsername,
		IsWalletLinked:      c.GetBool(middleware.StakeholderContextIsWalletLinked),
	}, nil
}

// OrganizationAuthContextFromGin builds the subset of stakeholder auth context
// needed by organization-wide compliance routes. Those routes intentionally do
// not require a portal stakeholder mapping.
func OrganizationAuthContextFromGin(c *gin.Context) (AuthContext, error) {
	memberID := strings.TrimSpace(c.GetString("member_id"))
	organizationID := strings.TrimSpace(c.GetString("organization_id"))
	if memberID == "" || organizationID == "" {
		return AuthContext{}, NewHTTPError(http.StatusUnauthorized, "organization context is missing")
	}
	return AuthContext{
		MemberID:          memberID,
		MemberEmail:       c.GetString("member_email"),
		MemberRole:        c.GetString("member_role"),
		OrganizationID:    organizationID,
		OrganizationName:  c.GetString("organization_name"),
		OrganizationEmail: c.GetString("organization_email"),
		OrganizationType:  c.GetString("organization_type"),
	}, nil
}

func ParsePagination(c *gin.Context) models.PaginationRequest {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", c.DefaultQuery("page_size", "20")))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return models.PaginationRequest{Page: page, Limit: limit}
}

func PaginationMeta(page, limit int, total int64) models.PaginationMeta {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	totalPages := 0
	if total > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}
	return models.PaginationMeta{Page: page, Limit: limit, Total: total, TotalPages: totalPages}
}

func ParseDate(value string) (*time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	for _, layout := range []string{time.RFC3339, "2006-01-02"} {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			return &parsed, nil
		}
	}
	return nil, NewHTTPError(http.StatusBadRequest, "invalid date format")
}

func ParseRequiredDecimal(value, field string) (decimal.Decimal, error) {
	amount, err := decimal.NewFromString(strings.TrimSpace(value))
	if err != nil {
		return decimal.Zero, NewHTTPError(http.StatusBadRequest, field+" must be a decimal number")
	}
	if !amount.GreaterThan(decimal.Zero) {
		return decimal.Zero, NewHTTPError(http.StatusUnprocessableEntity, field+" must be greater than zero")
	}
	return amount, nil
}

func NormalizeCurrency(value string) (string, error) {
	value = strings.ToUpper(strings.TrimSpace(value))
	if len(value) < 3 || len(value) > 12 {
		return "", NewHTTPError(http.StatusBadRequest, "currency must be between 3 and 12 characters")
	}
	return value, nil
}

func ContainsString(values []string, value string) bool {
	value = strings.TrimSpace(value)
	for _, item := range values {
		if strings.TrimSpace(item) == value {
			return true
		}
	}
	return false
}

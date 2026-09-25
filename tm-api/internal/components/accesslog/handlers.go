package accesslog

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"admin-panel-dashboard/internal/models"
	serverModels "admin-panel-dashboard/internal/server/models"

	"github.com/gin-gonic/gin"
)

const (
	defaultPageSize = 10
	maxPageSize     = 100
	dateLayout      = "2006-01-02"
)

// listRow is the access-log row shaped for the admin Audit Trail UI. It carries
// both a machine-readable occurred_at and a pre-formatted `date` display string
// (the table prints `date` verbatim), so the screen works untouched while a later
// client-side formatter can switch to occurred_at.
type listRow struct {
	models.AdminAccessLog
	Date string `json:"date"` // "DD Mon, YYYY hh:mm AM/PM"
}

// paginatedData matches the frontend's IPaginatedResponse<T>.data shape.
type paginatedData struct {
	Data            []listRow `json:"data"`
	Page            int       `json:"page"`
	PageSize        int       `json:"pageSize"`
	Total           int64     `json:"total"`
	HasNextPage     bool      `json:"hasNextPage"`
	HasPreviousPage bool      `json:"hasPreviousPage"`
}

// AuditTrailPageData is the paginated envelope returned by ListHandler. Declared
// for Swagger so the response body renders with real fields.
type AuditTrailPageData struct {
	Status  int           `json:"status" example:"200"`
	Message string        `json:"message" example:"access logs fetched successfully"`
	Data    paginatedData `json:"data"`
}

// AuditTrailRecordData is the single-record envelope returned by DetailHandler.
type AuditTrailRecordData struct {
	Status  int     `json:"status" example:"200"`
	Message string  `json:"message" example:"access log fetched successfully"`
	Data    listRow `json:"data"`
}

// ListHandler serves GET /api/v1/audit-trail — the paginated, filterable access log.
//
// @Summary      List admin access / audit-trail events
// @Description  Paginated, filterable security log of admin logins and privileged actions. Super-admin only. The account-state filter (user=Active|Suspended) is accepted but is a property of the actor's account rather than of an access event, so it is currently a no-op.
// @ID           ListAdminAuditTrail
// @Tags         Audit Trail
// @Produce      json
// @Security     JwtTokenAuth
// @Param        Authorization header   string true  "JWT Token" default(Bearer <your-token>)
// @Param        page          query    int    false "Page number (default 1)"
// @Param        pageSize      query    int    false "Page size (default 10, max 100)"
// @Param        search        query    string false "Free text over username, full name or IP address"
// @Param        status        query    string false "Event outcome" Enums(successful, failed, pending)
// @Param        date          query    string false "Single-day filter, format YYYY-MM-DD"
// @Success      200 {object} AuditTrailPageData "Paginated access-log events"
// @Failure      400 {object} object "Invalid pagination parameters"
// @Failure      401 {object} object "Unauthorized"
// @Failure      403 {object} object "Forbidden: super-admin only"
// @Failure      500 {object} object "Internal server error"
// @Router       /audit-trail [get]
func ListHandler(s *serverModels.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, pageSize, ok := parsePaging(c)
		if !ok {
			return
		}

		q := s.AdminDB.Model(&models.AdminAccessLog{})

		if search := strings.TrimSpace(c.Query("search")); search != "" {
			like := "%" + strings.ToLower(search) + "%"
			q = q.Where(
				"LOWER(username) LIKE ? OR LOWER(full_name) LIKE ? OR ip_address LIKE ?",
				like, like, like,
			)
		}
		if status := canonicalStatus(c.Query("status")); status != "" {
			q = q.Where("status = ?", status)
		}
		if date := strings.TrimSpace(c.Query("date")); date != "" {
			if day, err := time.Parse(dateLayout, date); err == nil {
				next := day.Add(24 * time.Hour)
				q = q.Where("occurred_at >= ? AND occurred_at < ?", day, next)
			}
		}

		var total int64
		if err := q.Count(&total).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count access logs"})
			return
		}

		var rows []models.AdminAccessLog
		if err := q.Order("occurred_at DESC").
			Limit(pageSize).Offset((page - 1) * pageSize).
			Find(&rows).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch access logs"})
			return
		}

		out := make([]listRow, 0, len(rows))
		for _, r := range rows {
			out = append(out, toListRow(r))
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "access logs fetched successfully",
			"data": paginatedData{
				Data:            out,
				Page:            page,
				PageSize:        pageSize,
				Total:           total,
				HasNextPage:     int64(page*pageSize) < total,
				HasPreviousPage: page > 1,
			},
		})
	}
}

// DetailHandler serves GET /api/v1/audit-trail/:id — one access event, for the
// activity detail modal.
//
// @Summary      Get one admin access / audit-trail event
// @Description  Returns a single access-log event by id, for the activity detail view. Super-admin only.
// @ID           GetAdminAuditTrailEvent
// @Tags         Audit Trail
// @Produce      json
// @Security     JwtTokenAuth
// @Param        Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param        id            path   string true "Access-log event id"
// @Success      200 {object} AuditTrailRecordData "The access-log event"
// @Failure      401 {object} object "Unauthorized"
// @Failure      403 {object} object "Forbidden: super-admin only"
// @Failure      404 {object} object "Access log not found"
// @Router       /audit-trail/{id} [get]
func DetailHandler(s *serverModels.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var row models.AdminAccessLog
		if err := s.AdminDB.First(&row, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "access log not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "access log fetched successfully",
			"data":    toListRow(row),
		})
	}
}

func parsePaging(c *gin.Context) (page, pageSize int, ok bool) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page number"})
		return 0, 0, false
	}
	pageSize, err = strconv.Atoi(c.DefaultQuery("pageSize", strconv.Itoa(defaultPageSize)))
	if err != nil || pageSize < 1 || pageSize > maxPageSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page size"})
		return 0, 0, false
	}
	return page, pageSize, true
}

func toListRow(r models.AdminAccessLog) listRow {
	return listRow{AdminAccessLog: r, Date: r.OccurredAt.Format("02 Jan, 2006 03:04 PM")}
}

// canonicalStatus reconciles the UI filter's title-case (Successful/Pending/Failed)
// with the stored lowercase values. Unknown input yields "" (no filter).
func canonicalStatus(v string) string {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case models.AccessStatusSuccessful:
		return models.AccessStatusSuccessful
	case models.AccessStatusFailed:
		return models.AccessStatusFailed
	case models.AccessStatusPending:
		return models.AccessStatusPending
	default:
		return ""
	}
}

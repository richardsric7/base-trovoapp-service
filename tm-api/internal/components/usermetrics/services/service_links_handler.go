package usermetrics

import (
	usermetricsDB "admin-panel-dashboard/internal/components/usermetrics/db"
	"admin-panel-dashboard/internal/models"
	serverResponse "admin-panel-dashboard/internal/server/response"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Each handler below reuses checkAdminAuth (defined in
// user_metrics_handler.go, same package).

// @Summary List service links
// @Description Paginated, filterable list of white-label partner integration (service link) accounts, each enriched with a summary of its owner's user record.
// @Tags ServiceLinks
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Items per page" default(20)
// @Param ownerUsername query string false "Filter by owner username (partial match)"
// @Param shortName query string false "Filter by short name (partial match)"
// @Param inactive query string false "Filter by active status (true/false)"
// @Param verified query string false "Filter by verified status (true/false)"
// @Param suspended query string false "Filter by suspended status (true/false)"
// @Success 200 {object} response.Data
// @Failure 401,500 {object} object
// @Router /service-links [get]
func GetServiceLinksHandler(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !checkAdminAuth(c, walletDB) {
			return
		}
		page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		if err != nil || page < 1 {
			page = 1
		}
		pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
		if err != nil || pageSize <= 0 {
			pageSize = 20
		}
		req := models.ServiceLinkListRequest{
			Page:          page,
			PageSize:      pageSize,
			OwnerUsername: c.Query("ownerUsername"),
			ShortName:     c.Query("shortName"),
		}
		if v := c.Query("inactive"); v != "" {
			inactive := v == "true"
			req.Inactive = &inactive
		}
		if v := c.Query("verified"); v != "" {
			verified := v == "true"
			req.Verified = &verified
		}
		if v := c.Query("suspended"); v != "" {
			suspended := v == "true"
			req.Suspended = &suspended
		}
		links, total, err := usermetricsDB.ListServiceLinks(walletDB, req)
		if err != nil {
			log.Println("[SERVICE_LINKS] error listing:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		withOwners, err := usermetricsDB.AttachServiceLinkOwners(walletDB, links)
		if err != nil {
			log.Println("[SERVICE_LINKS] error attaching owners:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		serverResponse.JSON(c, http.StatusOK, "Service links fetched successfully", gin.H{
			"data":     withOwners,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		}, nil)
	}
}

// @Summary Get service link by ID
// @Tags ServiceLinks
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param id path string true "Service link ID"
// @Success 200 {object} response.Data
// @Failure 400,401,404,500 {object} object
// @Router /service-links/{id} [get]
func GetServiceLinkByIDHandler(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !checkAdminAuth(c, walletDB) {
			return
		}
		id := c.Param("id")
		link, err := usermetricsDB.GetServiceLinkByID(walletDB, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		owner, err := usermetricsDB.GetServiceLinkOwnerSummary(walletDB, link.OwnerUsername)
		if err != nil {
			log.Println("[SERVICE_LINKS] error resolving owner:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		serverResponse.JSON(c, http.StatusOK, "Service link fetched successfully", models.ServiceLinkWithOwner{ServiceLink: link, Owner: owner}, nil)
	}
}

// @Summary Create or update a service link
// @Description Action must be one of: `create`, `update`. ID should only be provided for `update`. On create, the owner username must already exist in the users table - its wallet address is copied onto the new row, and the ID/API key are always server-generated GUIDs, never accepted from the client.
// @Tags ServiceLinks
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param data body models.ServiceLinkRequest true "Service link payload"
// @Success 200 {object} response.Data
// @Failure 400,401,500 {object} object
// @Router /service-links [post]
func SaveServiceLinkHandler(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !checkAdminAuth(c, walletDB) {
			return
		}
		var req models.ServiceLinkRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if req.Action == "update" && req.ID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is required for update"})
			return
		}
		link, err := usermetricsDB.SaveServiceLink(walletDB, req)
		if err != nil {
			log.Println("[SERVICE_LINKS] error saving:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		message := "Service link updated successfully"
		if req.Action == "create" {
			message = "Service link created successfully"
		}
		serverResponse.JSON(c, http.StatusOK, message, link, nil)
	}
}

type setServiceLinkInactiveRequest struct {
	Inactive bool `json:"inactive"`
}

// @Summary Set a service link's active/inactive status
// @Description Retires or restores a service link without deleting it.
// @Tags ServiceLinks
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param id path string true "Service link ID"
// @Param data body setServiceLinkInactiveRequest true "Inactive flag"
// @Success 200 {object} response.Data
// @Failure 400,401,404,500 {object} object
// @Router /service-links/{id}/inactive [put]
func SetServiceLinkInactiveHandler(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !checkAdminAuth(c, walletDB) {
			return
		}
		id := c.Param("id")
		var req setServiceLinkInactiveRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		link, err := usermetricsDB.SetServiceLinkInactive(walletDB, id, req.Inactive)
		if err != nil {
			log.Println("[SERVICE_LINKS] error setting inactive:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		serverResponse.JSON(c, http.StatusOK, "Active status updated successfully", link, nil)
	}
}

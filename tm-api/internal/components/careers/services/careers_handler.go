package careers

import (
	careersDB "admin-panel-dashboard/internal/components/careers/db"
	userServices "admin-panel-dashboard/internal/components/users/services"
	"admin-panel-dashboard/internal/middleware"
	"admin-panel-dashboard/internal/models"
	serverResponse "admin-panel-dashboard/internal/server/response"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// validateLocation checks if the provided location is valid
func validateLocation(location string) bool {
	validLocations := []string{models.LocationRemote, models.LocationHybrid, models.LocationOnsite}
	for _, v := range validLocations {
		if location == v {
			return true
		}
	}
	return false
}

// validateWorkType checks if the provided work type is valid
func validateWorkType(workType string) bool {
	validWorkTypes := []string{models.WorkTypeFulltime, models.WorkTypeParttime, models.WorkTypeContract}
	for _, v := range validWorkTypes {
		if workType == v {
			return true
		}
	}
	return false
}

// validateYearsOfExperience checks if years of experience is within length limits
func validateYearsOfExperience(years string) bool {
	return len(years) <= 100
}

// CreateCareerRoleHandler handles creation of a new career role
// @Summary [Admin] Create career role
// @Description Create a new career role. Role will be created in draft status. Admin endpoint, authentication required.
// @Tags Careers
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param data body models.CareerRoleRequest true "Career Role payload"
// @Success 200 {object} models.CareerRoleResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/careers/roles [post]
func CreateCareerRoleHandler(adminDB *gorm.DB, walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID

		_, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[CAREERS] error for userID:", userID, "error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var req models.CareerRoleRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Println("[CAREERS] validation error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format", "details": err.Error()})
			return
		}

		// Manual validation for map types
		if req.Heading == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "heading is required"})
			return
		}
		if len(req.Sections) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "sections is required"})
			return
		}
		if len(req.Application) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "application is required"})
			return
		}

		// Validate new fields
		if !validateLocation(req.Location) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid location, must be one of: remote, hybrid, onsite"})
			return
		}
		if !validateWorkType(req.WorkType) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid work_type, must be one of: fulltime, parttime, contract"})
			return
		}
		if req.YearsOfExperience != "" && !validateYearsOfExperience(req.YearsOfExperience) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "years_of_experience must be 100 characters or less"})
			return
		}

		roleID, err := careersDB.CreateCareerRole(adminDB, req, userID)
		if err != nil {
			log.Println("[CAREERS] error creating role:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create career role"})
			return
		}

		serverResponse.JSON(c, http.StatusOK, "career role created successfully", gin.H{"id": roleID}, nil)
	}
}

// GetPublicCareerRolesHandler retrieves published career roles (public endpoint)
// @Summary [Public] Get published career roles
// @Description Retrieves only published career roles with optional search. Public endpoint, no authentication required.
// @Tags Careers
// @Produce json
// @Param search query string false "Search in heading and company overview"
// @Success 200 {array} models.CareerRoleResponse
// @Failure 500 {object} map[string]string
// @Router /careers/roles [get]
func GetPublicCareerRolesHandler(adminDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		search := c.Query("search")

		roles, err := careersDB.GetPublishedCareerRoles(adminDB, search)
		if err != nil {
			log.Println("[CAREERS] error fetching roles:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch career roles"})
			return
		}

		serverResponse.JSON(c, http.StatusOK, "career roles fetched successfully", roles, nil)
	}
}

// GetPublicCareerRoleByIDHandler retrieves a published career role by ID (public endpoint)
// @Summary [Public] Get published career role by ID
// @Description Retrieves a single published career role by ID. Public endpoint, no authentication required.
// @Tags Careers
// @Produce json
// @Param id path string true "Career Role ID"
// @Success 200 {object} models.CareerRoleResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /careers/roles/{id} [get]
func GetPublicCareerRoleByIDHandler(adminDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if strings.TrimSpace(id) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
			return
		}

		role, err := careersDB.GetPublishedCareerRoleByID(adminDB, id)
		if err != nil {
			if err.Error() == "career role not found" {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			log.Println("[CAREERS] error fetching role:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch career role"})
			return
		}

		serverResponse.JSON(c, http.StatusOK, "career role fetched successfully", role, nil)
	}
}

// GetAdminCareerRolesHandler retrieves all career roles with filters (admin endpoint)
// @Summary [Admin] Get all career roles with filters
// @Description Retrieves all career roles (including drafts) with optional filters. Admin endpoint, authentication required.
// @Tags Careers
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param status query string false "Filter by status" Enums(draft, published, disabled)
// @Param location query string false "Filter by location" Enums(remote, hybrid, onsite)
// @Param work_type query string false "Filter by work type" Enums(fulltime, parttime, contract)
// @Param years_of_experience query string false "Filter by years of experience"
// @Param search query string false "Search in heading and company overview"
// @Success 200 {array} models.CareerRoleResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/careers/roles [get]
func GetAdminCareerRolesHandler(adminDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filters models.CareerRoleFilters
		if err := c.ShouldBindQuery(&filters); err != nil {
			log.Println("[CAREERS] error binding query params:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters"})
			return
		}

		// Validate filter values
		if filters.Status != "" && filters.Status != models.CareerRoleStatusDraft &&
			filters.Status != models.CareerRoleStatusPublished && filters.Status != models.CareerRoleStatusDisabled {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status filter, must be one of: draft, published, disabled"})
			return
		}
		if filters.Location != "" && !validateLocation(filters.Location) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid location filter, must be one of: remote, hybrid, onsite"})
			return
		}
		if filters.WorkType != "" && !validateWorkType(filters.WorkType) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid work_type filter, must be one of: fulltime, parttime, contract"})
			return
		}

		roles, err := careersDB.GetAllCareerRolesWithFilters(adminDB, filters)
		if err != nil {
			log.Println("[CAREERS] error fetching roles:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch career roles"})
			return
		}

		serverResponse.JSON(c, http.StatusOK, "career roles fetched successfully", roles, nil)
	}
}

// GetAdminCareerRoleByIDHandler retrieves any career role by ID (admin endpoint)
// @Summary [Admin] Get career role by ID
// @Description Retrieves a career role by ID regardless of status. Admin endpoint, authentication required.
// @Tags Careers
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param id path string true "Career Role ID"
// @Success 200 {object} models.CareerRoleResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/careers/roles/{id} [get]
func GetAdminCareerRoleByIDHandler(adminDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if strings.TrimSpace(id) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
			return
		}

		role, err := careersDB.GetCareerRoleByID(adminDB, id, true)
		if err != nil {
			if err.Error() == "career role not found" {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			log.Println("[CAREERS] error fetching role:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch career role"})
			return
		}

		serverResponse.JSON(c, http.StatusOK, "career role fetched successfully", role, nil)
	}
}

// UpdateCareerRoleHandler updates an existing career role
// @Summary [Admin] Update career role
// @Description Updates an existing career role. Status cannot be changed via this endpoint. Admin endpoint, authentication required.
// @Tags Careers
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param id path string true "Career Role ID"
// @Param data body models.CareerRoleRequest true "Career Role update payload"
// @Success 200 {object} models.CareerRoleResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/careers/roles/{id} [put]
func UpdateCareerRoleHandler(adminDB *gorm.DB, walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID

		_, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[CAREERS] error for userID:", userID, "error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id := c.Param("id")
		if strings.TrimSpace(id) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
			return
		}

		var req models.CareerRoleRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
			return
		}

		// Validate new fields
		if !validateLocation(req.Location) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid location, must be one of: remote, hybrid, onsite"})
			return
		}
		if !validateWorkType(req.WorkType) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid work_type, must be one of: fulltime, parttime, contract"})
			return
		}
		if req.YearsOfExperience != "" && !validateYearsOfExperience(req.YearsOfExperience) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "years_of_experience must be 100 characters or less"})
			return
		}

		if err := careersDB.UpdateCareerRole(adminDB, id, req, userID); err != nil {
			if err.Error() == "career role not found" {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			log.Println("[CAREERS] error updating role:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update career role"})
			return
		}

		serverResponse.JSON(c, http.StatusOK, "career role updated successfully", nil, nil)
	}
}

// DeleteCareerRoleHandler deletes a career role by ID
// @Summary [Admin] Delete career role
// @Description Deletes a career role record by ID. Admin endpoint, authentication required.
// @Tags Careers
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param id path string true "Career Role ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/careers/roles/{id} [delete]
func DeleteCareerRoleHandler(adminDB *gorm.DB, walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID

		_, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[CAREERS] error for userID:", userID, "error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id := c.Param("id")
		if strings.TrimSpace(id) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
			return
		}

		if err := careersDB.DeleteCareerRole(adminDB, id); err != nil {
			if err.Error() == "career role not found" {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			log.Println("[CAREERS] error deleting role:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete career role"})
			return
		}

		serverResponse.JSON(c, http.StatusOK, "career role deleted successfully", nil, nil)
	}
}

// PublishCareerRoleHandler publishes a draft role
// @Summary [Admin] Publish career role
// @Description Publishes a draft career role, making it visible to the public. Only draft roles can be published. Admin endpoint, authentication required.
// @Tags Careers
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param id path string true "Career Role ID"
// @Success 200 {object} models.CareerRoleResponse
// @Failure 400 {object} map[string]string "Bad request - role not found or not in draft status"
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/careers/roles/{id}/publish [put]
func PublishCareerRoleHandler(adminDB *gorm.DB, walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID

		_, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[CAREERS] error for userID:", userID, "error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id := c.Param("id")
		if strings.TrimSpace(id) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
			return
		}

		if err := careersDB.PublishCareerRole(adminDB, id); err != nil {
			statusCode := http.StatusInternalServerError
			if err.Error() == "career role not found" || err.Error() == "only draft roles can be published" {
				statusCode = http.StatusBadRequest
			}
			log.Println("[CAREERS] error publishing role:", err)
			c.JSON(statusCode, gin.H{"error": err.Error()})
			return
		}

		serverResponse.JSON(c, http.StatusOK, "career role published successfully", nil, nil)
	}
}

// UpdateCareerRoleStatusHandler updates the status of a career role (enable/disable)
// @Summary [Admin] Update career role status
// @Description Updates the status of a career role. Can toggle between 'published' (enabled) and 'disabled'. Only published or disabled roles can have their status updated. Admin endpoint, authentication required.
// @Tags Careers
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param id path string true "Career Role ID"
// @Param data body models.UpdateStatusRequest true "Status update payload"
// @Success 200 {object} models.CareerRoleResponse
// @Failure 400 {object} map[string]string "Bad request - invalid status or role not found"
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/careers/roles/{id}/status [put]
func UpdateCareerRoleStatusHandler(adminDB *gorm.DB, walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID

		_, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[CAREERS] error for userID:", userID, "error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id := c.Param("id")
		if strings.TrimSpace(id) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
			return
		}

		var req models.UpdateStatusRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
			return
		}

		if err := careersDB.UpdateCareerRoleStatus(adminDB, id, req.Status); err != nil {
			statusCode := http.StatusInternalServerError
			if err.Error() == "career role not found" ||
				err.Error() == "draft roles must be published first" ||
				err.Error() == "invalid status: must be 'published' or 'disabled'" {
				statusCode = http.StatusBadRequest
			}
			log.Println("[CAREERS] error updating role status:", err)
			c.JSON(statusCode, gin.H{"error": err.Error()})
			return
		}

		serverResponse.JSON(c, http.StatusOK, "career role status updated successfully", nil, nil)
	}
}

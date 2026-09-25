package server

import (
	"admin-panel-dashboard/internal/middleware"
	"admin-panel-dashboard/internal/models"
	serverModels "admin-panel-dashboard/internal/server/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stellar/go/support/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func getAdminUserDetails(c *gin.Context, s *serverModels.Server) (models.UserInfo, bool) {
	ad, _ := middleware.ExtractTokenMetadata(c.Request)
	adminUserID := ad.UserID
	adminUser, done, _ := GetUserFromContext(adminUserID, s.TrovoWalletDB, c)
	if done {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve admin information"})
		return models.UserInfo{}, true
	}
	// var user models.User
	//if err := s.TrovoWalletDB.First(&user, "id = ?", adminUserID).Error; err != nil {
	//	c.JSON(http.StatusNotFound, gin.H{"error": "User does not exist"})
	//	return models.UserInfo{}, models.User{}, true
	//}
	return adminUser, false
}

// BulkCreateOrUpdateConfigurations handles the bulk creation or updating of configuration items
// @Summary Bulk create or update configurations
// @Description Creates or updates multiple configuration items of specified types
// @Tags Configurations
// @Accept json
// @Produce json
// @Param payload body []ConfigurationItem true "List of configuration items with type and data"
// @Success 200 {object} BulkCreateOrUpdateResponse
// @Failure 400 {object} ErrorResponse "Invalid request payload"
// @Failure 500 {object} ErrorResponse "Failed to process configurations"
// @Router /configurations/bulk [post]
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Example
// Example payload for different configuration types:
// [
//
//	{
//	    "type": "role",
//	    "data": {
//	        "subject": "Admin",
//	        "description": "Administrator role with full permissions"
//	    }
//	},
//	{
//	    "type": "suspensionReason",
//	    "data": {
//	        "subject": "Violation of terms",
//	        "description": "User violated terms and conditions",
//	        "target": "user"
//	    }
//	},
//	{
//	    "type": "currency",
//	    "data": {
//	        "subject": "USD",
//	        "description": "United States Dollar"
//	    }
//	}
//
// ]
func BulkCreateOrUpdateConfigurations(s *serverModels.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		userInfo, _ := getAdminUserDetails(c, s)

		var requests []ConfigurationItem
		if err := c.ShouldBindJSON(&requests); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request payload"})
			return
		}

		for _, item := range requests {
			var err error
			switch item.Type {
			case "role":
				err = processRoleConfiguration(s, item.Data, userInfo)
				log.Debug("Processing role configuration", "data", item.Data, "error", err)
			case "suspensionReason":
				err = processSuspensionReasonConfiguration(s, item.Data, userInfo)
				log.Debug("Processing suspension configuration", "data", item.Data, "error", err)
			case "currency":
				err = processCurrencyConfiguration(s, item.Data, userInfo)
				log.Debug("Processing currency configuration", "data", item.Data, "error", err)
			default:
				c.JSON(http.StatusBadRequest, ErrorResponse{Error: fmt.Sprintf("Unknown configuration type: %s", item.Type)})
				return
			}

			if err != nil {
				c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
				return
			}
		}

		c.JSON(http.StatusOK, BulkCreateOrUpdateResponse{Message: "Configurations processed successfully"})
	}
}

// ReadConfigurations handles retrieving configurations based on type and optional filters
// @Summary Retrieve configurations
// @Description Retrieve configurations based on specified type and optional filters
// @Tags Configurations
// @Accept json
// @Produce json
// @Param type query string true "Configuration type (role, suspensionReason, currency)"
// @Param subject query string false "Filter by specific subject"
// @Success 200 {array} map[string]interface{} "List of configurations"
// @Failure 400 {object} ErrorResponse "Invalid request parameters"
// @Failure 500 {object} ErrorResponse "Failed to retrieve configurations"
// @Router /configurations [get]
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
func ReadConfigurations(s *serverModels.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		configType := c.Query("type")
		subject := c.Query("subject")

		var results []map[string]interface{}
		var err error

		switch configType {
		case "role":
			results, err = fetchRoles(s, subject)
		case "suspensionReason":
			results, err = fetchSuspensionReasons(s, subject)
		case "currency":
			results, err = fetchCurrencies(s, subject)
		default:
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid or missing configuration type"})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, results)
	}
}

// DeleteConfiguration handles deleting a configuration item by type and subject
// @Summary Delete configuration
// @Description Deletes a configuration item of specified type and subject
// @Tags Configurations
// @Accept json
// @Produce json
// @Param type query string true "Configuration type (role, suspensionReason, currency)"
// @Param subject query string true "The unique subject of the item to delete"
// @Success 200 {object} BulkCreateOrUpdateResponse "Deletion success message"
// @Failure 400 {object} ErrorResponse "Invalid request parameters"
// @Failure 404 {object} ErrorResponse "Configuration item not found"
// @Failure 500 {object} ErrorResponse "Failed to delete configuration"
// @Router /configurations [delete]
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
func DeleteConfiguration(s *serverModels.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		configType := c.Query("type")
		subject := c.Query("subject")

		if subject == "" {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Subject is required"})
			return
		}

		var err error
		switch configType {
		case "role":
			err = deleteRole(s, subject)
		case "suspensionReason":
			err = deleteSuspensionReason(s, subject)
		case "currency":
			err = deleteCurrency(s, subject)
		default:
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid or missing configuration type"})
			return
		}

		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, ErrorResponse{Error: "Configuration item not found"})
			} else {
				c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			}
			return
		}

		c.JSON(http.StatusOK, BulkCreateOrUpdateResponse{Message: "Configuration deleted successfully"})
	}
}

// GetAllConfigurations retrieves all configurations across all types
// @Summary Retrieve all configurations
// @Description Retrieve all configurations across all types without filters
// @Tags Configurations
// @Accept json
// @Produce json
// @Success 200 {object} map[string][]map[string]interface{} "All configurations by type"
// @Failure 500 {object} ErrorResponse "Failed to retrieve configurations"
// @Router /configurations/all [get]
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
func GetAllConfigurations(s *serverModels.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		allConfigs := make(map[string][]map[string]interface{})

		// Fetch all roles
		allConfigs["roles"], err = fetchRoles(s, "")
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve roles"})
			return
		}

		// Fetch all suspension reasons
		allConfigs["suspensionReasons"], err = fetchSuspensionReasons(s, "")
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve suspension reasons"})
			return
		}

		// Fetch all currencies
		allConfigs["currencies"], err = fetchCurrencies(s, "")
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve currencies"})
			return
		}

		// Return all configurations grouped by type
		c.JSON(http.StatusOK, allConfigs)
	}
}

// Delete functions for each configuration type
func deleteRole(s *serverModels.Server, subject string) error {
	return s.AdminDB.Where("role_name = ?", subject).Delete(&models.RoleConfig{}).Error
}

func deleteSuspensionReason(s *serverModels.Server, subject string) error {
	return s.AdminDB.Where("reason = ?", subject).Delete(&models.SuspensionReason{}).Error
}

func deleteCurrency(s *serverModels.Server, subject string) error {
	return s.AdminDB.Where("code = ?", subject).Delete(&models.CurrencyConfig{}).Error
}

// Fetch functions for each configuration type
func fetchRoles(s *serverModels.Server, subject string) ([]map[string]interface{}, error) {
	var roles []models.RoleConfig
	query := s.AdminDB
	if subject != "" {
		query = query.Where("role_name = ?", subject)
	}
	if err := query.Find(&roles).Error; err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for _, role := range roles {
		results = append(results, map[string]interface{}{
			"subject":     role.RoleName,
			"description": role.Description,
		})
	}
	return results, nil
}

func fetchSuspensionReasons(s *serverModels.Server, subject string) ([]map[string]interface{}, error) {
	var reasons []models.SuspensionReason
	query := s.AdminDB
	if subject != "" {
		query = query.Where("reason = ?", subject)
	}
	if err := query.Find(&reasons).Error; err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for _, reason := range reasons {
		results = append(results, map[string]interface{}{
			"subject":     reason.Reason,
			"description": reason.Description,
			"target":      reason.Target,
		})
	}
	return results, nil
}

func fetchCurrencies(s *serverModels.Server, subject string) ([]map[string]interface{}, error) {
	var currencies []models.CurrencyConfig
	query := s.AdminDB
	if subject != "" {
		query = query.Where("code = ?", subject)
	}
	if err := query.Find(&currencies).Error; err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for _, currency := range currencies {
		results = append(results, map[string]interface{}{
			"subject":     currency.Code,
			"description": currency.Description,
		})
	}
	return results, nil
}

func processRoleConfiguration(s *serverModels.Server, data map[string]interface{}, userInfo models.UserInfo) error {
	roleName, _ := data["subject"].(string)
	description, _ := data["description"].(string)

	role := models.RoleConfig{
		RoleName:    roleName,
		Description: description,
	}

	// Upsert (insert or update on conflict)
	return s.AdminDB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "role_name"}},              // Unique key for conflict resolution
		DoUpdates: clause.AssignmentColumns([]string{"description"}), // Fields to update on conflict
	}).Create(&role).Error
}

func processSuspensionReasonConfiguration(s *serverModels.Server, data map[string]interface{}, userInfo models.UserInfo) error {
	subject, _ := data["subject"].(string)
	description, _ := data["description"].(string)
	target, _ := data["target"].(string) // "admin" or "user"

	// Validate target
	if target != "admin" && target != "user" {
		return fmt.Errorf("invalid target: %s", target)
	}

	reason := models.SuspensionReason{
		Reason:      subject,
		Description: description,
		Target:      target,
	}

	// Perform upsert using OnConflict
	return s.AdminDB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "reason"}, {Name: "target"}}, // Unique on reason and target
		DoUpdates: clause.AssignmentColumns([]string{"description"}),   // Update only description on conflict
	}).Create(&reason).Error
}

func processCurrencyConfiguration(s *serverModels.Server, data map[string]interface{}, userInfo models.UserInfo) error {
	subject, _ := data["subject"].(string)
	description, _ := data["description"].(string)

	currency := models.CurrencyConfig{
		Code:        subject,
		Description: description,
	}

	return s.AdminDB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "code"}},                   // Unique key for conflict resolution
		DoUpdates: clause.AssignmentColumns([]string{"description"}), // Fields to update on conflict
	}).Create(&currency).Error
}

// ConfigurationItem represents a generic configuration item in the payload
type ConfigurationItem struct {
	Type string                 `json:"type" binding:"required"` // Configuration type (e.g., "role", "suspensionReason", "currency")
	Data map[string]interface{} `json:"data" binding:"required"` // Key-value pairs for configuration details
	// RoleName string                 `json:"role_name"`
}

// UserSuspensionReason represents the model for a suspension reason in the database
type UserSuspensionReason struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Reason string `gorm:"uniqueIndex;not null" json:"reason"`
}

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// BulkCreateOrUpdateResponse represents a successful bulk creation/updating response
type BulkCreateOrUpdateResponse struct {
	Message string `json:"message"`
}

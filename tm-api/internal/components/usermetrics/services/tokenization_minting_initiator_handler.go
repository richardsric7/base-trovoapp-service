package usermetrics

import (
	usermetrics "admin-panel-dashboard/internal/components/usermetrics/db"
	userServices "admin-panel-dashboard/internal/components/users/services"
	"admin-panel-dashboard/internal/middleware"
	"admin-panel-dashboard/internal/models"
	serverResponse "admin-panel-dashboard/internal/server/response"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

///created_at and soft_delete are handled by
// trustees table

// @Summary Create a new minting user
// @Description Creates a new minting user (initiator or approver)
// @ID create-minting-user
// @Tags Tokenization-initiator-approver
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param request body models.CreateMintingUserRequest true "Create Minting User Request"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tokenization/minting/users [post]
func CreateMintingUser(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract user info from token
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID

		// Get user info
		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[MINTING] error for user:", userInfo.Username, "error: ", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Parse request
		var req models.CreateMintingUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate role
		if req.Role != models.MintingRoleInitiator && req.Role != models.MintingRoleApprover {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role. Must be either 'initiator' or 'approver'"})
			return
		}

		// Create minting user
		user, err := usermetrics.CreateMintingUser(walletDB, req.Username, req.Role)
		if err != nil {
			log.Println("[MINTING] error creating user:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create minting user"})
			return
		}

		serverResponse.JSON(c, http.StatusOK, "minting user created successfully", user, nil)
	}
}

// @Summary Delete a minting user
// @Description Deletes a minting user (initiator or approver)
// @ID delete-minting-user
// @Tags Tokenization-initiator-approver
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param request body models.DeleteMintingUserRequest true "Delete Minting User Request"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tokenization/minting/users [delete]
func DeleteMintingUser(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract user info from token
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID

		// Get user info
		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[MINTING] error for user:", userInfo.Username, "error: ", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Parse request
		var req models.DeleteMintingUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate role
		if req.Role != models.MintingRoleInitiator && req.Role != models.MintingRoleApprover {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role. Must be either 'initiator' or 'approver'"})
			return
		}

		// Delete minting user
		if err := usermetrics.DeleteMintingUser(walletDB, req.ID, req.Role); err != nil {
			log.Println("[MINTING] error deleting user:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete minting user"})
			return
		}

		serverResponse.JSON(c, http.StatusOK, "minting user deleted successfully", nil, nil)
	}
}

// @Summary Get all minting users
// @Description Retrieves all minting users (initiators and approvers) with pagination
// @ID get-all-minting-users
// @Tags Tokenization-initiator-approver
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param role query string false "Filter by role (initiator or approver)"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} models.PaginatedMintingUserResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tokenization/minting/users [get]
func GetAllMintingUsers(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract user info from token
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID

		// Get user info
		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[MINTING] error for user:", userInfo.Username, "error: ", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Get pagination parameters
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

		// Get role filter if provided
		role := c.Query("role")
		if role != "" && role != models.MintingRoleInitiator && role != models.MintingRoleApprover {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role. Must be either 'initiator' or 'approver'"})
			return
		}

		// Get all minting users
		users, err := usermetrics.GetAllMintingUsers(walletDB, role, page, pageSize)
		if err != nil {
			log.Println("[MINTING] error getting users:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get minting users"})
			return
		}

		serverResponse.JSON(c, http.StatusOK, "minting users fetched successfully", users, nil)
	}
}

// @Summary Search minting users
// @Description Searches for minting users (initiators and approvers) with pagination
// @ID search-minting-users
// @Tags Tokenization-initiator-approver
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param query query string true "Search query"
// @Param role query string false "Filter by role (initiator or approver)"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} models.PaginatedMintingUserResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tokenization/minting/users/search [get]
func SearchMintingUsers(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract user info from token
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID

		// Get user info
		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[MINTING] error for user:", userInfo.Username, "error: ", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Get search parameters
		query := c.Query("query")
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

		// Get role filter if provided
		role := c.Query("role")
		if role != "" && role != models.MintingRoleInitiator && role != models.MintingRoleApprover {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role. Must be either 'initiator' or 'approver'"})
			return
		}

		// Search minting users
		users, err := usermetrics.SearchMintingUsers(walletDB, query, role, page, pageSize)
		if err != nil {
			log.Println("[MINTING] error searching users:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search minting users"})
			return
		}

		serverResponse.JSON(c, http.StatusOK, "minting users searched successfully", users, nil)
	}
}

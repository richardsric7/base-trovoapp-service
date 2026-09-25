package usermetrics

import (
	usermetricsDB "admin-panel-dashboard/internal/components/usermetrics/db"
	userServices "admin-panel-dashboard/internal/components/users/services"
	"admin-panel-dashboard/internal/middleware"
	"admin-panel-dashboard/internal/models"
	serverResponse "admin-panel-dashboard/internal/server/response"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SaveFaucetConfigHandler handles create/update faucet config
// @Summary Save faucet config
// @Description Create or update a faucet config. Action must be one of: `create`, `update`. ID should only be provided for `update`.
// @Tags FaucetConfig
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param data body models.FaucetConfigRequest true "Faucet Config payload"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /faucet/configs [post]
func SaveFaucetConfigHandler(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID

		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[FAUCET_CONFIG] error for user:", userInfo.Username, "error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var req models.FaucetConfigRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
			return
		}

		if req.Action != "create" && req.Action != "update" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid action: must be 'create' or 'update'"})
			return
		}

		if req.Action == "update" && req.ID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required for update"})
			return
		}

		if err := usermetricsDB.SaveFaucetConfig(walletDB, req); err != nil {
			log.Println("[FAUCET_CONFIG] error saving config:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		message := "Faucet config updated successfully"
		if req.Action == "create" {
			message = "Faucet config created successfully"
		}

		serverResponse.JSON(c, http.StatusOK, message, nil, nil)
	}
}

// GetAllFaucetConfigsHandler retrieves all faucet configs
// @Summary Get all faucet configs
// @Description Retrieves all faucet configuration records
// @Tags FaucetConfig
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Success 200 {array} models.FaucetConfigResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /faucet/configs [get]
func GetAllFaucetConfigsHandler(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID

		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[FAUCET_CONFIG] error for user:", userInfo.Username, "error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		configs, err := usermetricsDB.GetAllFaucetConfigs(walletDB)
		if err != nil {
			log.Println("[FAUCET_CONFIG] error fetching configs:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch faucet configs"})
			return
		}

		serverResponse.JSON(c, http.StatusOK, "faucet configs fetched successfully", configs, nil)
	}
}

// GetFaucetConfigByIDHandler retrieves a faucet config by ID
// @Summary Get faucet config by ID
// @Description Retrieves a single faucet configuration record by ID
// @Tags FaucetConfig
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param id path int true "Faucet Config ID"
// @Success 200 {object} models.FaucetConfigResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /faucet/configs/{id} [get]
func GetFaucetConfigByIDHandler(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID

		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[FAUCET_CONFIG] error for user:", userInfo.Username, "error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
			return
		}

		config, err := usermetricsDB.GetFaucetConfigByID(walletDB, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		serverResponse.JSON(c, http.StatusOK, "faucet config fetched successfully", config, nil)
	}
}

// DeleteFaucetConfigHandler deletes a faucet config by ID
// @Summary Delete faucet config
// @Description Deletes a faucet configuration record by ID
// @Tags FaucetConfig
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param id path int true "Faucet Config ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /faucet/configs/{id} [delete]
func DeleteFaucetConfigHandler(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID

		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[FAUCET_CONFIG] error for user:", userInfo.Username, "error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
			return
		}

		if err := usermetricsDB.DeleteFaucetConfig(walletDB, id); err != nil {
			log.Println("[FAUCET_CONFIG] error deleting config:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete faucet config"})
			return
		}

		serverResponse.JSON(c, http.StatusOK, "faucet config deleted successfully", nil, nil)
	}
}

// SaveKycConfigHandler handles create/update KYC config
// @Summary Save KYC config
// @Description Create or update a KYC config. Action must be one of: `create`, `update`. ID should only be provided for `update`.
// @Tags KycConfig
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param data body models.KycConfigRequest true "KYC Config payload"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /kyc/configs [post]
func SaveKycConfigHandler(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID

		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[KYC_CONFIG] error for user:", userInfo.Username, "error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var req models.KycConfigRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
			return
		}

		if req.Action != "create" && req.Action != "update" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid action: must be 'create' or 'update'"})
			return
		}

		if req.Action == "update" && req.ID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required for update"})
			return
		}

		if err := usermetricsDB.SaveKycConfig(walletDB, req); err != nil {
			log.Println("[KYC_CONFIG] error saving config:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		message := "KYC config updated successfully"
		if req.Action == "create" {
			message = "KYC config created successfully"
		}

		serverResponse.JSON(c, http.StatusOK, message, nil, nil)
	}
}

// GetAllKycConfigsHandler retrieves all KYC configs
// @Summary Get all KYC configs
// @Description Retrieves all KYC configuration records
// @Tags KycConfig
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Success 200 {array} models.KycConfigResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /kyc/configs [get]
func GetAllKycConfigsHandler(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID

		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[KYC_CONFIG] error for user:", userInfo.Username, "error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		configs, err := usermetricsDB.GetAllKycConfigs(walletDB)
		if err != nil {
			log.Println("[KYC_CONFIG] error fetching configs:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch KYC configs"})
			return
		}

		serverResponse.JSON(c, http.StatusOK, "kyc configs fetched successfully", configs, nil)
	}
}

// GetKycConfigByIDHandler retrieves a KYC config by ID
// @Summary Get KYC config by ID
// @Description Retrieves a single KYC configuration record by ID
// @Tags KycConfig
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param id path int true "KYC Config ID"
// @Success 200 {object} models.KycConfigResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /kyc/configs/{id} [get]
func GetKycConfigByIDHandler(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID

		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[KYC_CONFIG] error for user:", userInfo.Username, "error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
			return
		}

		config, err := usermetricsDB.GetKycConfigByID(walletDB, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		serverResponse.JSON(c, http.StatusOK, "kyc config fetched successfully", config, nil)
	}
}

// DeleteKycConfigHandler deletes a KYC config by ID
// @Summary Delete KYC config
// @Description Deletes a KYC configuration record by ID
// @Tags KycConfig
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param id path int true "KYC Config ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /kyc/configs/{id} [delete]
func DeleteKycConfigHandler(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID

		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[KYC_CONFIG] error for user:", userInfo.Username, "error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
			return
		}

		if err := usermetricsDB.DeleteKycConfig(walletDB, id); err != nil {
			log.Println("[KYC_CONFIG] error deleting config:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete KYC config"})
			return
		}

		serverResponse.JSON(c, http.StatusOK, "kyc config deleted successfully", nil, nil)
	}
}

// SaveDojaWidgetHandler handles create/update doja widget
// @Summary Save doja widget
// @Description Create or update a doja widget. Action must be one of: `create`, `update`. ID should only be provided for `update`.
// @Tags DojaWidget
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param data body models.DojaWidgetRequest true "Doja Widget payload"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /doja/widgets [post]
func SaveDojaWidgetHandler(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID

		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[DOJA_WIDGET] error for user:", userInfo.Username, "error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var req models.DojaWidgetRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
			return
		}

		if req.Action != "create" && req.Action != "update" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid action: must be 'create' or 'update'"})
			return
		}

		if req.Action == "update" && req.ID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required for update"})
			return
		}

		if err := usermetricsDB.SaveDojaWidget(walletDB, req); err != nil {
			log.Println("[DOJA_WIDGET] error saving widget:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		message := "Doja widget updated successfully"
		if req.Action == "create" {
			message = "Doja widget created successfully"
		}

		serverResponse.JSON(c, http.StatusOK, message, nil, nil)
	}
}

// GetAllDojaWidgetsHandler retrieves all doja widgets
// @Summary Get all doja widgets
// @Description Retrieves all doja widget records
// @Tags DojaWidget
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Success 200 {array} models.DojaWidgetResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /doja/widgets [get]
func GetAllDojaWidgetsHandler(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID

		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[DOJA_WIDGET] error for user:", userInfo.Username, "error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		widgets, err := usermetricsDB.GetAllDojaWidgets(walletDB)
		if err != nil {
			log.Println("[DOJA_WIDGET] error fetching widgets:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch doja widgets"})
			return
		}

		serverResponse.JSON(c, http.StatusOK, "doja widgets fetched successfully", widgets, nil)
	}
}

// GetDojaWidgetByIDHandler retrieves a doja widget by ID
// @Summary Get doja widget by ID
// @Description Retrieves a single doja widget record by ID
// @Tags DojaWidget
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param id path string true "Doja Widget ID"
// @Success 200 {object} models.DojaWidgetResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /doja/widgets/{id} [get]
func GetDojaWidgetByIDHandler(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID

		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[DOJA_WIDGET] error for user:", userInfo.Username, "error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
			return
		}

		widget, err := usermetricsDB.GetDojaWidgetByID(walletDB, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		serverResponse.JSON(c, http.StatusOK, "doja widget fetched successfully", widget, nil)
	}
}

// DeleteDojaWidgetHandler deletes a doja widget by ID
// @Summary Delete doja widget
// @Description Deletes a doja widget record by ID
// @Tags DojaWidget
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param id path string true "Doja Widget ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /doja/widgets/{id} [delete]
func DeleteDojaWidgetHandler(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID

		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[DOJA_WIDGET] error for user:", userInfo.Username, "error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
			return
		}

		if err := usermetricsDB.DeleteDojaWidget(walletDB, id); err != nil {
			log.Println("[DOJA_WIDGET] error deleting widget:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete doja widget"})
			return
		}

		serverResponse.JSON(c, http.StatusOK, "doja widget deleted successfully", nil, nil)
	}
}

// SaveKycLevelHandler handles create/update KYC level
// @Summary Save KYC level
// @Description Create or update a KYC level. Action must be one of: `create`, `update`. For create action, ID should not be provided. For update action, ID is required. User category must be unique (case-insensitive).
// @Tags KycLevel
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param data body models.KycLevelRequest true "KYC Level payload"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} map[string]string "Bad request - invalid action, missing ID for update, ID provided for create, or user category already exists"
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /kyc/levels [post]
func SaveKycLevelHandler(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID

		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[KYC_LEVEL] error for user:", userInfo.Username, "error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var req models.KycLevelRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
			return
		}

		if req.Action != "create" && req.Action != "update" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid action: must be 'create' or 'update'"})
			return
		}

		if req.Action == "create" && strings.TrimSpace(req.ID) != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID should not be provided for create action"})
			return
		}

		if req.Action == "update" && req.ID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required for update"})
			return
		}

		if err := usermetricsDB.SaveKycLevel(walletDB, req); err != nil {
			log.Println("[KYC_LEVEL] error saving level:", err)
			// Return 400 for validation errors, 500 for other errors
			statusCode := http.StatusInternalServerError
			errorMsg := err.Error()
			if errorMsg == "user category already exists" || errorMsg == "kyc level not found" {
				statusCode = http.StatusBadRequest
			}
			c.JSON(statusCode, gin.H{"error": errorMsg})
			return
		}

		message := "KYC level updated successfully"
		if req.Action == "create" {
			message = "KYC level created successfully"
		}

		serverResponse.JSON(c, http.StatusOK, message, nil, nil)
	}
}

// GetAllKycLevelsHandler retrieves all KYC levels
// @Summary Get all KYC levels
// @Description Retrieves all KYC level records
// @Tags KycLevel
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Success 200 {array} models.KycLevelResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /kyc/levels [get]
func GetAllKycLevelsHandler(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID

		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[KYC_LEVEL] error for user:", userInfo.Username, "error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		levels, err := usermetricsDB.GetAllKycLevels(walletDB)
		if err != nil {
			log.Println("[KYC_LEVEL] error fetching levels:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch KYC levels"})
			return
		}

		serverResponse.JSON(c, http.StatusOK, "kyc levels fetched successfully", levels, nil)
	}
}

// GetKycLevelByIDHandler retrieves a KYC level by ID
// @Summary Get KYC level by ID
// @Description Retrieves a single KYC level record by ID
// @Tags KycLevel
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param id path string true "KYC Level ID"
// @Success 200 {object} models.KycLevelResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /kyc/levels/{id} [get]
func GetKycLevelByIDHandler(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID

		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[KYC_LEVEL] error for user:", userInfo.Username, "error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
			return
		}

		level, err := usermetricsDB.GetKycLevelByID(walletDB, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		serverResponse.JSON(c, http.StatusOK, "kyc level fetched successfully", level, nil)
	}
}

// DeleteKycLevelHandler deletes a KYC level by ID
// @Summary Delete KYC level
// @Description Deletes a KYC level record by ID
// @Tags KycLevel
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param id path string true "KYC Level ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /kyc/levels/{id} [delete]
func DeleteKycLevelHandler(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, _ := middleware.ExtractTokenMetadata(c.Request)
		userID := ad.UserID

		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[KYC_LEVEL] error for user:", userInfo.Username, "error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
			return
		}

		if err := usermetricsDB.DeleteKycLevel(walletDB, id); err != nil {
			log.Println("[KYC_LEVEL] error deleting level:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete KYC level"})
			return
		}

		serverResponse.JSON(c, http.StatusOK, "kyc level deleted successfully", nil, nil)
	}
}

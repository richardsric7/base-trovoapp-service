package usermetrics

import (
	usermetricsDB "admin-panel-dashboard/internal/components/usermetrics/db"
	p2pErrors "admin-panel-dashboard/internal/errors"
	"admin-panel-dashboard/internal/middleware"
	"admin-panel-dashboard/internal/models"
	serverResponse "admin-panel-dashboard/internal/server/response" // Import bytes package

	// Import encoding/json package
	"fmt" // Import io package
	"log"
	"net/http"
	"strconv"
	"strings" // Import strings package

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	userServices "admin-panel-dashboard/internal/components/users/services"
)

// TODO: Implement models.TrusteesRequest and usermetricsDB.SaveTrustees for trustees CRUD support

// SavePartnerHandler handles create/update of various partner types based on query params
// @Summary Save partner entity (Manager, Issuing House, Custodian)
// @Description Creates or updates a partner entity based on the 'type' and 'action' query parameters.
// @Description The JSON request body should contain the specific fields for the partner type.
// @Description For 'update' action, the 'id' field MUST be included in the JSON request body.
// @Tags Partners
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param type query string true "Partner type" Enums(asset_manager, asset_issuing_house, approved_asset_custodian, legal_and_professionals, rating_agency, trustees, legal_adviser, financial_adviser)
// @Param action query string true "Action to perform" Enums(create, update)
// @Param data body object true "JSON payload structure depends on 'type' query param. Refer to Schemas section for models: models.AssetManagerRequest, models.AssetIssuingHouseRequest, models.ApprovedAssetCustodianRequest, models.LegalAndProfessionalsRequest, models.RatingAgencyRequest, models.TrusteesRequest, models.LegalAdviserRequest, models.FinancialAdviserRequest. Include 'id' in body for update actions."
// @Success 200 {object} models.StandardSuccessResponse "Success message (e.g., 'Asset Manager created successfully')"
// @Failure 400 {object} map[string]string "Bad Request (e.g., invalid format, missing fields, invalid type/action, missing ID for update)"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Not Found (e.g., ID not found for update)"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /partners/save [post]
func SavePartnerHandler(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// --- Auth and User Info --- Start
		ad, err := middleware.ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized access"})
			return
		}
		userID := ad.UserID

		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[SAVE_PARTNER] error retrieving user info for UserID:", userID, "error:", err)
			var ex p2pErrors.GenericError
			var ok bool
			ex, ok = err.(p2pErrors.GenericError)
			var statusCode int
			var response interface{}
			if ok {
				statusCode = ex.HTTPCode()
				response = ex.JSONError()
			} else {
				statusCode = http.StatusBadRequest
				response = gin.H{"error": "failed to retrieve user info"}
			}
			c.JSON(statusCode, response)
			return
		}
		// --- Auth and User Info --- End

		// Get type and action from query parameters
		partnerType := strings.ToLower(strings.TrimSpace(c.Query("type")))
		action := strings.ToLower(strings.TrimSpace(c.Query("action")))

		partnerRequest := &SavePartnerRequest{
			PartnerType: partnerType,
			Action:      action,
			UserInfo:    userInfo,
			Context:     c,
			WalletDB:    walletDB,
		}
		// Validate type
		successMessage, _, err := SaveAndUpdatePartnerData(partnerRequest)
		if err != nil {
			log.Printf("[SAVE_PARTNER] error saving partner data: %v (User: %s)", err, userInfo.Username)
			var ex p2pErrors.GenericError
			var ok bool
			ex, ok = err.(p2pErrors.GenericError)
			var statusCode int
			var response interface{}
			if ok {
				statusCode = ex.HTTPCode()
				response = ex.JSONError()
			} else {
				// Determine status code based on error content
				if strings.Contains(err.Error(), "invalid") || strings.Contains(err.Error(), "missing") {
					statusCode = http.StatusBadRequest
				} else if strings.Contains(err.Error(), "not found for update") {
					statusCode = http.StatusNotFound
				} else {
					statusCode = http.StatusInternalServerError
				}
				response = gin.H{"error": err.Error()}
			}
			c.JSON(statusCode, response)
			return
		}
		serverResponse.JSON(c, http.StatusOK, successMessage, nil, nil)
	}
}

/*
	req.Address = partnerRequest.Address
		req.Country = partnerRequest.Country
		req.FeePercent = partnerRequest.FeePercent
		req.FeeFixed = partnerRequest.FeeFixed
*/

func SaveAndUpdatePartnerData(partnerRequest *SavePartnerRequest) (string, uint64, error) {
	action := partnerRequest.Action
	partnerType := partnerRequest.PartnerType
	userInfo := partnerRequest.UserInfo
	walletDB := partnerRequest.WalletDB
	c := partnerRequest.Context

	if partnerType != "asset_manager" && partnerType != "asset_issuing_house" && partnerType != "approved_asset_custodian" && partnerType != "legal_and_professionals" && partnerType != "rating_agency" && partnerType != "trustees" && partnerType != "legal_adviser" && partnerType != "financial_adviser" {
		log.Printf("[SAVE_PARTNER] Invalid type specified in query param: %s (User: %s)", partnerType, userInfo.Username)
		return "", 0, fmt.Errorf("invalid type specified in query parameter. Must be one of: asset_manager, asset_issuing_house, approved_asset_custodian, legal_and_professionals, rating_agency, trustees, legal_adviser, financial_adviser")
	}

	// Validate action
	if action != "create" && action != "update" {
		log.Printf("[SAVE_PARTNER] Invalid action specified in query param: %s (User: %s)", action, userInfo.Username)
		return "", 0, fmt.Errorf("invalid action specified in query parameter. Must be 'create' or 'update'")
	}

	var successMessage string
	var savedID uint64
	var errDb error

	// Switch based on the type to bind the BODY to the correct struct and call the correct DB function
	switch partnerType {
	case "asset_manager":
		var req models.AssetManagerRequest

		// Check if we should use direct fields or JSON binding
		if partnerRequest.UseDirectFields {
			// Programmatic population (from inviteUserFlow)
			req.AssetManagerName = partnerRequest.PartnerName
			req.AssetManagerAddress = partnerRequest.Address
			req.AssetManagerCountry = partnerRequest.Country
			// Set ID for update action
			if action == "update" {
				req.ID = int64(partnerRequest.PartnerID)
			}
			// Parse fee fields
			if partnerRequest.FeePercent != "" {
				feePercent, err := strconv.ParseFloat(partnerRequest.FeePercent, 64)
				if err != nil {
					return "", 0, fmt.Errorf("invalid fee percent: %v", err)
				}
				req.FeePercent = feePercent
			}
			if partnerRequest.FeeFixed != "" {
				feeFixed, err := strconv.ParseFloat(partnerRequest.FeeFixed, 64)
				if err != nil {
					return "", 0, fmt.Errorf("invalid fee fixed: %v", err)
				}
				req.FeeFixed = feeFixed
			}
		} else {
			// HTTP JSON binding (from SavePartnerHandler)
			if err := c.ShouldBindJSON(&req); err != nil {
				log.Println("[SAVE_PARTNER] AssetManager bind error:", err, "(User:", userInfo.Username, ")")
				return "", 0, fmt.Errorf("invalid asset manager request body format: %w", err)
			}
		}

		savedID, errDb = usermetricsDB.SaveAssetManager(action, req, walletDB)
		if errDb == nil {
			if action == "create" {
				successMessage = "Asset Manager created successfully"
			} else {
				successMessage = "Asset Manager updated successfully"
			}
		}

	case "asset_issuing_house":
		var req models.AssetIssuingHouseRequest

		// Check if we should use direct fields or JSON binding
		if partnerRequest.UseDirectFields {
			// Programmatic population (from inviteUserFlow)
			req.AssetIssuingHouseName = partnerRequest.PartnerName
			req.AssetIssuingHouseAddress = partnerRequest.Address
			req.AssetIssuingHouseCountry = partnerRequest.Country
			// Set ID for update action
			if action == "update" {
				req.ID = uint(partnerRequest.PartnerID)
			}
			// Parse fee fields
			if partnerRequest.FeePercent != "" {
				feePercent, err := strconv.ParseFloat(partnerRequest.FeePercent, 64)
				if err != nil {
					return "", 0, fmt.Errorf("invalid fee percent: %v", err)
				}
				req.FeePercent = feePercent
			}
			if partnerRequest.FeeFixed != "" {
				feeFixed, err := strconv.ParseFloat(partnerRequest.FeeFixed, 64)
				if err != nil {
					return "", 0, fmt.Errorf("invalid fee fixed: %v", err)
				}
				req.FeeFixed = feeFixed
			}
		} else {
			// HTTP JSON binding (from SavePartnerHandler)
			if err := c.ShouldBindJSON(&req); err != nil {
				log.Println("[SAVE_PARTNER] AssetIssuingHouse bind error:", err, "(User:", userInfo.Username, ")")
				return "", 0, fmt.Errorf("invalid asset issuing house request body format: %w", err)
			}
		}

		savedID, errDb = usermetricsDB.SaveAssetIssuingHouse(action, req, walletDB)
		if errDb == nil {
			if action == "create" {
				successMessage = "Asset Issuing House created successfully"
			} else {
				successMessage = "Asset Issuing House updated successfully"
			}
		}

	case "approved_asset_custodian":
		var req models.ApprovedAssetCustodianRequest

		// Check if we should use direct fields or JSON binding
		if partnerRequest.UseDirectFields {
			// Programmatic population (from inviteUserFlow)
			req.AssetCustodianName = partnerRequest.PartnerName
			req.AssetCustodianAddress = partnerRequest.Address
			req.AssetCustodianCountry = partnerRequest.Country
			// Set ID for update action
			if action == "update" {
				req.ID = uint(partnerRequest.PartnerID)
			}
			req.FeePercent = partnerRequest.FeePercent // Already strings
			req.FeeFixed = partnerRequest.FeeFixed     // Already strings
			// RequirementDocument is optional, leave empty for now
		} else {
			// HTTP JSON binding (from SavePartnerHandler)
			if err := c.ShouldBindJSON(&req); err != nil {
				log.Println("[SAVE_PARTNER] ApprovedAssetCustodian bind error:", err, "(User:", userInfo.Username, ")")
				return "", 0, fmt.Errorf("invalid approved asset custodian request body format: %w", err)
			}
		}

		savedID, errDb = usermetricsDB.SaveApprovedAssetCustodian(action, req, walletDB)
		if errDb == nil {
			if action == "create" {
				successMessage = "Approved Asset Custodian created successfully"
			} else {
				successMessage = "Approved Asset Custodian updated successfully"
			}
		}
	case "legal_and_professionals":
		var req models.LegalAndProfessionalsRequest

		// Check if we should use direct fields or JSON binding
		if partnerRequest.UseDirectFields {
			// Programmatic population (from inviteUserFlow)
			req.PartnerName = partnerRequest.PartnerName
			req.PartnerAddress = partnerRequest.Address
			req.PartnerCountry = partnerRequest.Country
			// Set ID for update action
			if action == "update" {
				req.ID = uint(partnerRequest.PartnerID)
			}
			req.FeePercent = partnerRequest.FeePercent // Already strings
			req.FeeFixed = partnerRequest.FeeFixed     // Already strings
		} else {
			// HTTP JSON binding (from SavePartnerHandler)
			if err := c.ShouldBindJSON(&req); err != nil {
				log.Println("[SAVE_PARTNER] LegalAndProfessionals bind error:", err, "(User:", userInfo.Username, ")")
				return "", 0, fmt.Errorf("invalid legal and professionals request body format: %w", err)
			}
		}

		savedID, errDb = usermetricsDB.SaveLegalAndProfessionals(action, req, walletDB)
		if errDb == nil {
			if action == "create" {
				successMessage = "Legal and Professionals created successfully"
			} else {
				successMessage = "Legal and Professionals updated successfully"
			}
		}
	case "rating_agency":
		var req models.RatingAgencyRequest

		// Check if we should use direct fields or JSON binding
		if partnerRequest.UseDirectFields {
			// Programmatic population (from inviteUserFlow)
			req.AgencyName = partnerRequest.PartnerName
			req.AgencyAddress = partnerRequest.Address
			req.AgencyCountry = partnerRequest.Country
			// Set ID for update action
			if action == "update" {
				req.ID = uint(partnerRequest.PartnerID)
			}
			req.FeePercent = partnerRequest.FeePercent // Already strings
			req.FeeFixed = partnerRequest.FeeFixed     // Already strings
		} else {
			// HTTP JSON binding (from SavePartnerHandler)
			if err := c.ShouldBindJSON(&req); err != nil {
				log.Println("[SAVE_PARTNER] RatingAgency bind error:", err, "(User:", userInfo.Username, ")")
				return "", 0, fmt.Errorf("invalid rating agency request body format: %w", err)
			}
		}

		savedID, errDb = usermetricsDB.SaveRatingAgency(action, req, walletDB)
		if errDb == nil {
			if action == "create" {
				successMessage = "Rating Agency created successfully"
			} else {
				successMessage = "Rating Agency updated successfully"
			}
		}
	case "trustees":
		var req models.TrusteesRequest

		// Check if we should use direct fields or JSON binding
		if partnerRequest.UseDirectFields {
			// Programmatic population (from inviteUserFlow)
			req.TrusteeName = partnerRequest.PartnerName
			req.TrusteeAddress = partnerRequest.Address
			req.TrusteeCountry = partnerRequest.Country
			// Set ID for update action
			if action == "update" {
				req.ID = uint(partnerRequest.PartnerID)
			}
			req.FeePercent = partnerRequest.FeePercent // Already strings
			req.FeeFixed = partnerRequest.FeeFixed     // Already strings
		} else {
			// HTTP JSON binding (from SavePartnerHandler)
			if err := c.ShouldBindJSON(&req); err != nil {
				log.Println("[SAVE_PARTNER] Trustees bind error:", err, "(User:", userInfo.Username, ")")
				return "", 0, fmt.Errorf("invalid trustees request body format: %w", err)
			}
		}

		savedID, errDb = usermetricsDB.SaveTrustees(action, req, walletDB)
		if errDb == nil {
			if action == "create" {
				successMessage = "Trustees created successfully"
			} else {
				successMessage = "Trustees updated successfully"
			}
		}
	case "legal_adviser":
		var req models.LegalAdviserRequest

		// Check if we should use direct fields or JSON binding
		if partnerRequest.UseDirectFields {
			// Programmatic population (from inviteUserFlow)
			req.AdviserName = partnerRequest.PartnerName
			req.AdviserAddress = partnerRequest.Address
			req.AdviserCountry = partnerRequest.Country
			// Set ID for update action
			if action == "update" {
				req.ID = uint(partnerRequest.PartnerID)
			}
			req.FeePercent = partnerRequest.FeePercent // Already strings
			req.FeeFixed = partnerRequest.FeeFixed     // Already strings
		} else {
			// HTTP JSON binding (from SavePartnerHandler)
			if err := c.ShouldBindJSON(&req); err != nil {
				log.Println("[SAVE_PARTNER] LegalAdviser bind error:", err, "(User:", userInfo.Username, ")")
				return "", 0, fmt.Errorf("invalid legal adviser request body format: %w", err)
			}
		}

		savedID, errDb = usermetricsDB.SaveLegalAdviser(action, req, walletDB)
		if errDb == nil {
			if action == "create" {
				successMessage = "Legal Adviser created successfully"
			} else {
				successMessage = "Legal Adviser updated successfully"
			}
		}
	case "financial_adviser":
		var req models.FinancialAdviserRequest

		// Check if we should use direct fields or JSON binding
		if partnerRequest.UseDirectFields {
			// Programmatic population (from inviteUserFlow)
			req.AdviserName = partnerRequest.PartnerName
			req.AdviserAddress = partnerRequest.Address
			req.AdviserCountry = partnerRequest.Country
			// Set ID for update action
			if action == "update" {
				req.ID = uint(partnerRequest.PartnerID)
			}
			req.FeePercent = partnerRequest.FeePercent // Already strings
			req.FeeFixed = partnerRequest.FeeFixed     // Already strings
		} else {
			// HTTP JSON binding (from SavePartnerHandler)
			if err := c.ShouldBindJSON(&req); err != nil {
				log.Println("[SAVE_PARTNER] FinancialAdviser bind error:", err, "(User:", userInfo.Username, ")")
				return "", 0, fmt.Errorf("invalid financial adviser request body format: %w", err)
			}
		}

		savedID, errDb = usermetricsDB.SaveFinancialAdviser(action, req, walletDB)
		if errDb == nil {
			if action == "create" {
				successMessage = "Financial Adviser created successfully"
			} else {
				successMessage = "Financial Adviser updated successfully"
			}
		}
	}

	// Handle DB errors
	if errDb != nil {
		log.Printf("[SAVE_PARTNER] DB error for type '%s', action '%s', ID '%d': %v (User: %s)\n", partnerType, action, savedID, errDb, userInfo.Username)
		// Check for specific errors (like missing ID or not found)
		if strings.Contains(errDb.Error(), "missing id") {
			return "", 0, errDb
		} else if strings.Contains(errDb.Error(), "not found for update") {
			return "", 0, errDb
		} else {
			return "", 0, fmt.Errorf("failed to save partner data: %w", errDb)
		}
	}

	// Return success response
	log.Printf("[SAVE_PARTNER] %s (Type: %s, Action: %s, ID: %d, User: %s)\n", successMessage, partnerType, action, savedID, userInfo.Username)
	return successMessage, savedID, nil
}

type SavePartnerRequest struct {
	PartnerType string
	Action      string
	UserInfo    models.UserInfo
	Context     *gin.Context
	WalletDB    *gorm.DB

	// Optional fields for programmatic creation (inviteUserFlow)
	UseDirectFields bool   // Flag to indicate manual population instead of JSON binding
	PartnerName     string // Partner/Organization name
	Address         string `json:"address" example:"123 Main St"`
	Country         string `json:"country" example:"USA"`
	FeePercent      string `json:"fee_percent" example:"0.5"`
	FeeFixed        string `json:"fee_fixed" example:"10.00"`
	PartnerID       uint64 // ID of the partner to update (required for update action)
}

// DeletePartnerHandler handles delete operations
// @Summary Delete partner entity
// @Description Delete partner entity by type and id
// @Tags Partners
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param type query string true "Partner type" Enums(asset_manager, asset_issuing_house, approved_asset_custodian, legal_and_professionals, rating_agency, trustees, legal_adviser, financial_adviser)
// @Param id query int true "Partner ID"
// @Success 200 {object} models.StandardSuccessResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /partners/delete [delete]
func DeletePartnerHandler(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, err := middleware.ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized access"})
			return
		}
		userID := ad.UserID

		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[METRICS] error for user:", userInfo.Username, "error: ", err)

			var ex p2pErrors.GenericError
			var ok bool

			ex, ok = err.(p2pErrors.GenericError)
			var statusCode int
			var response interface{}

			if ok {
				statusCode = ex.HTTPCode()
				response = ex.JSONError()
			} else {
				statusCode = http.StatusBadRequest
				response = gin.H{"error": err.Error()}
			}

			c.JSON(statusCode, response)
			return
		}

		typeStr := c.Query("type")
		id := c.Query("id")
		if typeStr == "" || id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing type or id"})
			return
		}

		err = usermetricsDB.DeletePartnerEntity(typeStr, id, walletDB)
		if err != nil {
			log.Println("[PARTNER_DELETE] error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		serverResponse.JSON(c, http.StatusOK, "Partner deleted successfully", nil, nil)
	}
}

// ListPartnersHandler handles listing partners
// @Summary List partner entities
// @Description List partners. If 'type' query parameter is provided, lists only that type. If 'type' is omitted, lists all types.
// @Tags Partners
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param type query string false "Optional: Partner type to filter by" Enums(asset_manager, asset_issuing_house, approved_asset_custodian, legal_and_professionals, rating_agency, trustees, legal_adviser, financial_adviser)
// @Success 200 {object} interface{} "If type is provided, returns array of that type. If type is omitted, returns map[string]array with keys 'asset_manager', 'asset_issuing_house', 'approved_asset_custodian'."
// @Failure 400 {object} map[string]string "Bad Request (e.g., invalid type)"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /partners/list [get]
func ListPartnersHandler(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// --- Auth and User Info --- Start
		ad, err := middleware.ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized access"})
			return
		}
		userID := ad.UserID

		userInfo, err := userServices.GetUser(userID, walletDB)
		if err != nil {
			log.Println("[LIST_PARTNERS] error retrieving user info for UserID:", userID, "error:", err)

			var ex p2pErrors.GenericError
			var ok bool
			ex, ok = err.(p2pErrors.GenericError)
			var statusCode int
			var response interface{}

			if ok {
				statusCode = ex.HTTPCode()
				response = ex.JSONError()
			} else {
				statusCode = http.StatusBadRequest
				response = gin.H{"error": "failed to retrieve user info"}
			}

			c.JSON(statusCode, response)
			return
		}
		// --- Auth and User Info --- End

		typeStr := c.Query("type")
		var result interface{}

		if typeStr == "" {
			// Type parameter is missing, fetch all partner types
			log.Println("[LIST_PARTNERS] Type parameter missing, fetching all partners.")
			result, err = usermetricsDB.ListAllPartnerEntities(walletDB)
			if err != nil {
				log.Printf("[LIST_PARTNERS] Error fetching all partners: %v (User: %s)", err, userInfo.Username)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch partner list"})
				return
			}
		} else {
			// Type parameter is provided, fetch specific type
			log.Printf("[LIST_PARTNERS] Fetching partners of type: %s (User: %s)", typeStr, userInfo.Username)
			result, err = usermetricsDB.ListPartnerEntities(typeStr, walletDB)
			if err != nil {
				// Check if the error is due to an invalid type from the DB function
				if strings.Contains(err.Error(), "invalid type specified") || strings.Contains(err.Error(), "invalid partner type") {
					log.Printf("[LIST_PARTNERS] Invalid type requested: %s (User: %s)", typeStr, userInfo.Username)
					c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid partner type specified: %s", typeStr)})
				} else {
					log.Printf("[LIST_PARTNERS] Error fetching partners of type %s: %v (User: %s)", typeStr, err, userInfo.Username)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch partner list"})
				}
				return
			}
		}

		log.Printf("[LIST_PARTNERS] Successfully fetched partners. Type requested: '%s' (User: %s)", typeStr, userInfo.Username)
		c.JSON(http.StatusOK, result)
	}
}

// GetPartnerHandler retrieves a single partner entity by type and ID
// @Summary Get one partner entity
// @Description Retrieve a single partner entity by type and ID
// @Tags Partners
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param type query string true "Partner type" Enums(asset_manager, asset_issuing_house, approved_asset_custodian, legal_and_professionals, rating_agency, trustees, legal_adviser, financial_adviser)
// @Param id query int true "Partner ID"
// @Success 200 {object} interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /partners/one [get]
func GetPartnerHandler(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := middleware.ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized access"})
			return
		}

		typeStr := c.Query("type")
		idStr := c.Query("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		entity, err := usermetricsDB.GetPartnerEntityByID(typeStr, id, walletDB)
		if err != nil {
			log.Println("[GET_PARTNER_ENTITY] error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, entity)
	}
}

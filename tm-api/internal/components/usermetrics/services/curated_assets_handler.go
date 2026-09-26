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
// user_metrics_handler.go, same package) - the calling admin's JWT-derived
// wallet_user_id must resolve to a real trovo-wallet user, exactly like
// every other P2P/reports admin endpoint.

// @Summary List curated assets
// @Description Paginated, filterable list of the platform's curated asset catalog.
// @Tags CuratedAssets
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Items per page" default(20)
// @Param assetCode query string false "Filter by asset code (partial match)"
// @Param assetClassId query int false "Filter by asset class"
// @Param p2pEnabled query string false "Filter by P2P-enabled (true/false)"
// @Param inactive query string false "Filter by active status (true/false)"
// @Success 200 {object} response.Data
// @Failure 401,500 {object} object
// @Router /assets/curated [get]
func GetCuratedAssetsHandler(walletDB *gorm.DB) gin.HandlerFunc {
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
		req := models.CuratedAssetListRequest{
			Page:      page,
			PageSize:  pageSize,
			AssetCode: c.Query("assetCode"),
		}
		if v := c.Query("assetClassId"); v != "" {
			if id, e := strconv.ParseUint(v, 10, 64); e == nil {
				req.AssetClassID = id
			}
		}
		if v := c.Query("p2pEnabled"); v != "" {
			enabled := v == "true"
			req.P2PEnabled = &enabled
		}
		if v := c.Query("inactive"); v != "" {
			inactive := v == "true"
			req.Inactive = &inactive
		}
		assets, total, err := usermetricsDB.ListCuratedAssets(walletDB, req)
		if err != nil {
			log.Println("[CURATED_ASSETS] error listing:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		serverResponse.JSON(c, http.StatusOK, "Curated assets fetched successfully", gin.H{
			"data":     assets,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		}, nil)
	}
}

// @Summary Get curated asset by ID
// @Tags CuratedAssets
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param id path int true "Curated asset ID"
// @Success 200 {object} response.Data
// @Failure 400,401,404,500 {object} object
// @Router /assets/curated/{id} [get]
func GetCuratedAssetByIDHandler(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !checkAdminAuth(c, walletDB) {
			return
		}
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		asset, err := usermetricsDB.GetCuratedAssetByID(walletDB, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		serverResponse.JSON(c, http.StatusOK, "Curated asset fetched successfully", asset, nil)
	}
}

// @Summary Create or update a curated asset
// @Description Action must be one of: `create`, `update`. ID should only be provided for `update`. Setting p2pEnabled makes the asset available for P2P offer creation/marketplace search; unset (or absent on create) keeps it unavailable there.
// @Tags CuratedAssets
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param data body models.CuratedAssetRequest true "Curated asset payload"
// @Success 200 {object} response.Data
// @Failure 400,401,500 {object} object
// @Router /assets/curated [post]
func SaveCuratedAssetHandler(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !checkAdminAuth(c, walletDB) {
			return
		}
		var req models.CuratedAssetRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if req.Action == "update" && req.ID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is required for update"})
			return
		}
		asset, err := usermetricsDB.SaveCuratedAsset(walletDB, req)
		if err != nil {
			log.Println("[CURATED_ASSETS] error saving:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		message := "Curated asset updated successfully"
		if req.Action == "create" {
			message = "Curated asset created successfully"
		}
		serverResponse.JSON(c, http.StatusOK, message, asset, nil)
	}
}

type setP2PEnabledRequest struct {
	P2PEnabled bool `json:"p2pEnabled"`
}

// @Summary Set a curated asset's P2P-enabled flag
// @Description Enables or disables an asset for P2P operations without touching any of its other fields.
// @Tags CuratedAssets
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param id path int true "Curated asset ID"
// @Param data body setP2PEnabledRequest true "P2P-enabled flag"
// @Success 200 {object} response.Data
// @Failure 400,401,404,500 {object} object
// @Router /assets/curated/{id}/p2p-enabled [put]
func SetCuratedAssetP2PEnabledHandler(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !checkAdminAuth(c, walletDB) {
			return
		}
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		var req setP2PEnabledRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		asset, err := usermetricsDB.SetCuratedAssetP2PEnabled(walletDB, id, req.P2PEnabled)
		if err != nil {
			log.Println("[CURATED_ASSETS] error setting p2pEnabled:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		serverResponse.JSON(c, http.StatusOK, "P2P-enabled flag updated successfully", asset, nil)
	}
}

type setInactiveRequest struct {
	Inactive bool `json:"inactive"`
}

// @Summary Set a curated asset's active/inactive status
// @Description Retires or restores a curated asset without deleting it - deactivation is the supported way to remove an asset from active use everywhere it's already referenced (wallets, historical offers/orders).
// @Tags CuratedAssets
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param id path int true "Curated asset ID"
// @Param data body setInactiveRequest true "Inactive flag"
// @Success 200 {object} response.Data
// @Failure 400,401,404,500 {object} object
// @Router /assets/curated/{id}/inactive [put]
func SetCuratedAssetInactiveHandler(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !checkAdminAuth(c, walletDB) {
			return
		}
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		var req setInactiveRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		asset, err := usermetricsDB.SetCuratedAssetInactive(walletDB, id, req.Inactive)
		if err != nil {
			log.Println("[CURATED_ASSETS] error setting inactive:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		serverResponse.JSON(c, http.StatusOK, "Active status updated successfully", asset, nil)
	}
}

// @Summary List asset classes
// @Description Reference list of asset categories (token/stablecoin/sto/nft) for the curated-asset admin form.
// @Tags CuratedAssets
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Success 200 {object} response.Data
// @Failure 401,500 {object} object
// @Router /asset-classes [get]
func GetAssetClassesHandler(walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !checkAdminAuth(c, walletDB) {
			return
		}
		classes, err := usermetricsDB.ListAssetClasses(walletDB)
		if err != nil {
			log.Println("[CURATED_ASSETS] error listing asset classes:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		serverResponse.JSON(c, http.StatusOK, "Asset classes fetched successfully", classes, nil)
	}
}


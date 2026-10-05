package users

import (
	"net/http"

	usersDB "trovo-wallet-api/internal/components/users/db"
	userServices "trovo-wallet-api/internal/components/users/services"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/middleware"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gin-gonic/gin"
)

// GasFeeAssetInput is the body of PUT /v1/users/settings/gas-fee-asset.
type GasFeeAssetInput struct {
	// AssetCode is a curated asset with gasFeeEligible; empty = pay in ETH.
	AssetCode string `json:"assetCode" example:"USDC"`
}

// getUsersGasFeeAssetsHandler godoc
// @Summary GET /v1/users/settings/gas-fee-assets
// @Description The stablecoins the user can pay network fees in, and the one currently chosen (gasFeeAsset; null = ETH).
// @Tags users
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /v1/users/settings/gas-fee-assets [get]
func getUsersGasFeeAssetsHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := usersDB.GetSlimUserFromPrimarySigner(middleware.ExtractSigner(c), gc.DB, gc)
		if err != nil {
			respondError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"gasFeeAsset": user.GasFeeAsset, "assets": userServices.GasFeeAssets(gc)})
	}
}

// putUsersGasFeeAssetHandler godoc
// @Summary PUT /v1/users/settings/gas-fee-asset
// @Description Sets the stablecoin the user's wallets pay network fees in (through the paymaster), or clears it with an empty assetCode. A wallet that holds too little of it pays in ETH.
// @Tags users
// @Accept json
// @Produce json
// @Param body body GasFeeAssetInput true "Gas fee asset"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /v1/users/settings/gas-fee-asset [put]
func putUsersGasFeeAssetHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		var in GasFeeAssetInput
		if err := c.ShouldBindJSON(&in); err != nil {
			var invalidJSON tErrors.ErrorInvalidJSON
			c.JSON(http.StatusBadRequest, invalidJSON.JSONError())
			return
		}
		user, err := usersDB.GetSlimUserFromPrimarySigner(middleware.ExtractSigner(c), gc.DB, gc)
		if err != nil {
			respondError(c, err)
			return
		}
		if err := userServices.SetGasFeeAsset(&user, in.AssetCode, gc); err != nil {
			respondError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"gasFeeAsset": user.GasFeeAsset})
	}
}

func respondError(c *gin.Context, err error) {
	if ex, ok := err.(tErrors.GenericError); ok {
		c.JSON(http.StatusBadRequest, ex.JSONError())
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": err.Error()})
}

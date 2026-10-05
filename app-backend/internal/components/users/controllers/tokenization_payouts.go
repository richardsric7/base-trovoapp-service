package users

import (
	"net/http"
	"strings"

	usersDB "trovo-wallet-api/internal/components/users/db"
	userServices "trovo-wallet-api/internal/components/users/services"
	"trovo-wallet-api/internal/middleware"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gin-gonic/gin"
)

// getTokenizationPayoutsHandler godoc
// @Summary GET /v1/tokenization/payouts
// @Description The proceeds payouts (dividends / interest) to the user's wallets, newest first: paid ones with their transaction, and scheduled ones once a payout's holder schedule is locked. Optional tokenizedAssetId narrows it to one asset.
// @Tags users
// @Produce json
// @Param tokenizedAssetId query string false "Tokenized asset id"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /v1/tokenization/payouts [get]
func getTokenizationPayoutsHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := usersDB.GetSlimUserFromPrimarySigner(middleware.ExtractSigner(c), gc.DB, gc)
		if err != nil {
			writeError(c, err)
			return
		}
		list, err := userServices.ListProceedPayoutReceipts(&user, strings.TrimSpace(c.Query("tokenizedAssetId")), gc)
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"payouts": list})
	}
}

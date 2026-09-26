package p2p

import (
	"encoding/json"
	"io"
	"net/http"
	p2pModels "trovo-wallet-api/internal/components/p2p/models"
	p2pServices "trovo-wallet-api/internal/components/p2p/services"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gin-gonic/gin"
)

type adminResolveDisputeRequest struct {
	Resolution        string `json:"resolution"`
	ResolvedByAdminID string `json:"resolvedByAdminId"`
}

var validAdminResolutions = map[string]bool{
	p2pModels.DisputeResolutionInFavorOfBuyer:  true,
	p2pModels.DisputeResolutionInFavorOfSeller: true,
	p2pModels.DisputeResolutionSplit:           true,
}

// postAdminResolveDisputeHandler is the arbiter-driven path for a dispute
// that neither self-resolution shortcut (Section 63) covers. Gated by
// AuthenticationMiddlewareUsingAPIKey - app-backend trusts that whichever
// admin-permissioned system calls this (tm-api) has already authorized the
// specific admin user; the resolving admin's identity travels in the
// request body, not this endpoint's own auth.
func postAdminResolveDisputeHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, _ := io.ReadAll(c.Request.Body)
		var req adminResolveDisputeRequest
		if err := json.Unmarshal(data, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error-invalid-json"})
			return
		}
		if !validAdminResolutions[req.Resolution] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error-invalid-resolution", "message": "resolution must be IN_FAVOR_OF_BUYER, IN_FAVOR_OF_SELLER, or SPLIT"})
			return
		}
		if req.ResolvedByAdminID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error-missing-resolved-by", "message": "resolvedByAdminId is required"})
			return
		}
		order, err := p2pServices.AdminResolveDispute(gc, c.Param("disputeID"), req.Resolution, req.ResolvedByAdminID)
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusOK, order)
	}
}

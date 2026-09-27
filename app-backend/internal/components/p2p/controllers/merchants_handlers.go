package p2p

import (
	"encoding/json"
	"io"
	"net/http"
	p2pServices "trovo-wallet-api/internal/components/p2p/services"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gin-gonic/gin"
)

// merchantStatusResponse is the shape a client checks before deciding
// whether to show a merchant-only page's own content, or the "only
// merchants can do this - request merchant status" notice.
type merchantStatusResponse struct {
	IsMerchant     bool `json:"isMerchant"`
	MerchantOnline bool `json:"merchantOnline"`
	KYCLevel       int  `json:"kycLevel"`
}

// getMerchantStatusHandler godoc
// @Summary Get my P2P merchant status
// @Description Tells the caller whether they are already a P2P merchant, whether they're currently online, and their KYC level - used by clients to decide whether to show merchant-only screens or a "become a merchant" prompt.
// @Tags P2P
// @Produce json
// @Success 200 {object} merchantStatusResponse
// @Failure 400 {object} map[string]interface{}
// @Security SignatureAuth
// @Router /v1/p2p/merchants/status [get]
func getMerchantStatusHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := currentUser(c, gc)
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusOK, merchantStatusResponse{
			IsMerchant:     user.IsMerchant,
			MerchantOnline: user.MerchantOnline,
			KYCLevel:       user.KYCVerified,
		})
	}
}

// postMerchantRequestHandler godoc
// @Summary Request P2P merchant status
// @Description Requests that the caller be upgraded to a P2P merchant, gated on the caller already holding KYC level 2.
// @Tags P2P
// @Produce json
// @Success 200 {object} merchantStatusResponse
// @Failure 400 {object} map[string]interface{} "KYC level too low, or already a merchant"
// @Security SignatureAuth
// @Router /v1/p2p/merchants/request [post]
func postMerchantRequestHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := currentUser(c, gc)
		if err != nil {
			writeError(c, err)
			return
		}
		updated, err := p2pServices.RequestMerchantStatus(gc, user.ID)
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusOK, merchantStatusResponse{
			IsMerchant:     updated.IsMerchant,
			MerchantOnline: updated.MerchantOnline,
			KYCLevel:       updated.KYCVerified,
		})
	}
}

type setMerchantOnlineStatusRequest struct {
	Online bool `json:"online"`
}

// putMerchantOnlineStatusHandler godoc
// @Summary Set my merchant online/offline status
// @Description Lets a merchant toggle whether they currently appear as available to trade. Offline merchants' offers are typically hidden or marked unavailable on the marketplace.
// @Tags P2P
// @Accept json
// @Produce json
// @Param body body setMerchantOnlineStatusRequest true "online: true or false"
// @Success 200 {object} merchantStatusResponse
// @Failure 400 {object} map[string]interface{}
// @Security SignatureAuth
// @Router /v1/p2p/merchants/online-status [put]
func putMerchantOnlineStatusHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := currentUser(c, gc)
		if err != nil {
			writeError(c, err)
			return
		}
		data, _ := io.ReadAll(c.Request.Body)
		var req setMerchantOnlineStatusRequest
		if err := json.Unmarshal(data, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error-invalid-json"})
			return
		}
		updated, err := p2pServices.SetMerchantOnlineStatus(gc, user.ID, req.Online)
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusOK, merchantStatusResponse{
			IsMerchant:     updated.IsMerchant,
			MerchantOnline: updated.MerchantOnline,
			KYCLevel:       updated.KYCVerified,
		})
	}
}

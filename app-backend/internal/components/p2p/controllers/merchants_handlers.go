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

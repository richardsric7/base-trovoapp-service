package p2p

import (
	"encoding/json"
	"io"
	"net/http"
	p2pServices "trovo-wallet-api/internal/components/p2p/services"
	usersDB "trovo-wallet-api/internal/components/users/db"
	"trovo-wallet-api/internal/middleware"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/gin-gonic/gin"
)

type depositEscrowRequest struct {
	Transaction          string `json:"transaction"`
	TransactionSignature string `json:"transactionSignature"`
	Commit               int    `json:"commit"`
}

// postEscrowDepositHandler godoc
// @Summary Deposit the sold asset into P2P escrow
// @Description Used by the asset depositor (the SELL-side party) to fund escrow for an order. Call it twice: first with an empty body to receive an unsigned transaction to sign locally, then again with the signed transaction and commit=1 to actually submit it. The source wallet must be a standard, non-temp wallet with no approver-based shared access.
// postEscrowDepositHandler mirrors postUsersPaymentHandler's own
// constraints exactly (Plan Section 28/29): source wallet must be a
// standard wallet (WalletType == 0), not temp, and not a shared-access
// wallet with approvers. The first call (no Commit) returns an unsigned
// transaction for the client to sign locally; the second call, carrying
// the signature and Commit=1, actually submits it.
// @Tags P2P
// @Accept json
// @Produce json
// @Param orderID path string true "Order ID"
// @Param body body depositEscrowRequest false "transaction (base64, from a prior unsigned call), transactionSignature, commit (1 to submit)"
// @Success 202 {object} map[string]interface{} "Unsigned transaction to sign, or submission result"
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{} "Not the asset depositor, wrong wallet type, or shared-access wallet with approvers"
// @Failure 404 {object} map[string]interface{} "Order not found"
// @Security SignatureAuth
// @Router /v1/p2p/orders/{orderID}/escrow-deposit [post]
func postEscrowDepositHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := currentUser(c, gc)
		if err != nil {
			writeError(c, err)
			return
		}
		order, err := p2pServices.GetOrderByID(gc.DB, c.Param("orderID"))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "error-order-not-found"})
			return
		}
		if order.AssetDepositor != user.ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "error-forbidden", "message": "You are not the asset depositor on this order"})
			return
		}

		sourceWallet, err := usersDB.GetWallet(middleware.ExtractAddress(c), gc.DB)
		if err != nil {
			writeError(c, err)
			return
		}
		if sourceWallet.WalletType != 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "error-wallet-type-forbidden", "message": "Only standard wallets are allowed to fund escrow directly."})
			return
		}
		if sourceWallet.SharedAccessEnabled == 1 && sourceWallet.WalletCountApproverAccess(gc) > 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "error-wallet-with-shared-access-not-allowed", "message": "This wallet has approver access enabled. Share the escrow deposit link with an initiator-access wallet instead."})
			return
		}

		data, _ := io.ReadAll(c.Request.Body)
		var req depositEscrowRequest
		_ = json.Unmarshal(data, &req) // an empty/absent body is valid for the first (unsigned) call

		result, err := p2pServices.DepositEscrowFromOwnWallet(gc, &user, &sourceWallet, &order, p2pServices.DepositEscrowInput{
			Transaction:          req.Transaction,
			TransactionSignature: req.TransactionSignature,
			Commit:               req.Commit,
		})
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusAccepted, result)
	}
}

// postRegenerateEscrowShortlinkHandler godoc
// @Summary Regenerate the shareable escrow-deposit link for an order
// @Description Retries shortlink generation if the automatic attempt made when the order was accepted failed and left the order without a deposit shortlink.
// @Tags P2P
// @Produce json
// @Param orderID path string true "Order ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Security SignatureAuth
// @Router /v1/p2p/orders/{orderID}/escrow-deposit/regenerate-shortlink [post]
func postRegenerateEscrowShortlinkHandler(gc *sharedconfig.GlobalConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := currentUser(c, gc)
		if err != nil {
			writeError(c, err)
			return
		}
		order, err := p2pServices.RegenerateEscrowShortlink(gc, c.Param("orderID"), user.ID)
		if err != nil {
			writeError(c, err)
			return
		}
		c.JSON(http.StatusOK, order)
	}
}

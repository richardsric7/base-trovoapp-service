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

// postEscrowDepositHandler mirrors postUsersPaymentHandler's own
// constraints exactly (Plan Section 28/29): source wallet must be a
// standard wallet (WalletType == 0), not temp, and not a shared-access
// wallet with approvers. The first call (no Commit) returns an unsigned
// transaction for the client to sign locally; the second call, carrying
// the signature and Commit=1, actually submits it.
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

		sourceWallet, temp, err := usersDB.GetWallet(middleware.ExtractAddress(c), gc.DB)
		if err != nil {
			writeError(c, err)
			return
		}
		if temp {
			c.JSON(http.StatusForbidden, gin.H{"error": "error-account-not-primary-account-alias", "message": "only primary/subwallets are allowed for escrow deposits"})
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

// postRegenerateEscrowShortlinkHandler lets the escrow-deposit screen
// retry shortlink generation if the best-effort attempt made during
// AcceptOrder failed and left EscrowDepositShortlink empty.
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

package p2p

import (
	"os"
	"strings"
	"time"
	p2pModels "trovo-wallet-api/internal/components/p2p/models"
	"trovo-wallet-api/internal/components/p2p/safesigner"
	paymentModels "trovo-wallet-api/internal/components/payments/models"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/dynamiclinks"
	userServices "trovo-wallet-api/internal/components/users/services"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// EscrowWalletAddress returns the fixed Base address escrow deposits are
// paid into. This is the same address the Universal Safe settlement release
// (Section 59-61) signs from - a Gnosis Safe contract address can receive
// tokens directly, so one env var serves both purposes.
func EscrowWalletAddress() string {
	return os.Getenv("P2P_ESCROW_WALLET_ADDRESS")
}

// GenerateEscrowShortlink implements Plan Section 30's decision: call
// dynamiclinks.GeneratePaymentData directly so the resulting long URL uses
// action=payment, which the client app's existing deep-link handler already
// knows how to open. Returns both the shortlink and QR in one call and
// persists them on the Order.
func GenerateEscrowShortlink(gc *sharedconfig.GlobalConfig, order *p2pModels.Order) error {
	escrowAddress := EscrowWalletAddress()
	if escrowAddress == "" {
		return &tErrors.CustomError{Param: "escrow", Err: "error-escrow-wallet-not-configured", ErrMessage: "P2P escrow wallet is not configured"}
	}
	data, err := dynamiclinks.GeneratePaymentData(escrowAddress, order.Asset, order.AssetContractAddress, order.SellerEscrowAssetAmount, order.ID, gc)
	if err != nil {
		return err
	}
	order.EscrowDepositShortlink = data.DynamicLink
	order.EscrowDepositQRCode = data.QRCode
	if err := gc.DB.Model(&p2pModels.Order{}).Where("id = ?", order.ID).
		Updates(map[string]interface{}{
			"escrow_deposit_shortlink": data.DynamicLink,
			"escrow_deposit_qr_code":   data.QRCode,
		}).Error; err != nil {
		return &tErrors.ErrorTemporaryServerError{}
	}
	RecordAuditEvent(gc, order.ID, order.OfferID, p2pModels.EventShortlinkCreated, order.AssetDepositor, map[string]string{"shortlink": data.DynamicLink})
	return nil
}

// RegenerateEscrowShortlink retries shortlink generation for the depositor
// when the best-effort attempt AcceptOrder makes inline failed (Plan Section
// 26's own note that a shortlink failure must not fail acceptance itself) -
// the escrow-deposit screen calls this if it loads an order with no
// shortlink yet.
func RegenerateEscrowShortlink(gc *sharedconfig.GlobalConfig, orderID, callerUserID string) (p2pModels.Order, error) {
	order, err := GetOrderByID(gc.DB, orderID)
	if err != nil {
		return order, &tErrors.CustomError{Param: "orderId", Err: "error-order-not-found", ErrMessage: "Order not found"}
	}
	if order.AssetDepositor != callerUserID {
		return order, &tErrors.CustomError{Param: "orderId", Err: "error-forbidden", ErrMessage: "You are not the escrow depositor on this order", Code: 403}
	}
	if order.OrderStatus != p2pModels.OrderStatusAwaitingEscrowDeposit {
		return order, &tErrors.CustomError{Param: "orderId", Err: "error-invalid-order-state", ErrMessage: "This order is not awaiting an escrow deposit"}
	}
	if order.EscrowDepositShortlink != "" {
		return order, nil
	}
	if err := GenerateEscrowShortlink(gc, &order); err != nil {
		return order, err
	}
	return order, nil
}

// DepositEscrowInput proxies the same two-phase build-then-commit contract
// userServices.Pay/postUsersPaymentHandler already use (Plan Section 28/29):
// the first call (no Transaction/Commit) returns an unsigned transaction for
// the client to sign locally; the second call carries the signature and
// Commit=1 to actually submit.
type DepositEscrowInput struct {
	Transaction          string
	TransactionSignature string
	Commit               int
}

// DepositEscrowFromOwnWallet lets the order's own assetDepositor fund escrow
// directly from their own standard wallet, reusing userServices.Pay
// in-process (Plan Section 28) rather than a second payment system. Only
// callable while the order is AWAITING_ESCROW_DEPOSIT.
func DepositEscrowFromOwnWallet(gc *sharedconfig.GlobalConfig, signerUser *userModels.User, sourceWallet *userModels.UserWallet, order *p2pModels.Order, in DepositEscrowInput) (*paymentModels.PaymentInfo, error) {
	if order.OrderStatus != p2pModels.OrderStatusAwaitingEscrowDeposit {
		return nil, &tErrors.CustomError{Param: "orderId", Err: "error-invalid-order-state", ErrMessage: "This order is not awaiting an escrow deposit"}
	}
	escrowAddress := EscrowWalletAddress()
	if escrowAddress == "" {
		return nil, &tErrors.CustomError{Param: "escrow", Err: "error-escrow-wallet-not-configured", ErrMessage: "P2P escrow wallet is not configured"}
	}

	paymentInfo := &paymentModels.PaymentInfo{
		Destination:          escrowAddress,
		Memo:                 order.ID,
		ContractAddress:      order.AssetContractAddress,
		AssetCode:            order.Asset,
		Amount:               order.SellerEscrowAssetAmount,
		Transaction:          in.Transaction,
		TransactionSignature: in.TransactionSignature,
		Commit:               in.Commit,
		Messages:             make([]string, 0),
	}

	returnedInfo, _, err := userServices.Pay(signerUser, sourceWallet, paymentInfo, gc)
	if err != nil {
		return returnedInfo, err
	}

	// Only the final, committed call has a real TransactionID - the first
	// (unsigned-transaction) call returns it empty and must be passed back
	// to the client to sign, not treated as a deposit.
	if returnedInfo != nil && returnedInfo.TransactionID != "" && returnedInfo.TransactionID != "PENDING_AUTH" {
		if err := recordCanonicalDeposit(gc, order, escrowAddress, sourceWallet.ID, returnedInfo.TransactionID, "P2P_API"); err != nil {
			return returnedInfo, err
		}
	}
	return returnedInfo, nil
}

// recordCanonicalDeposit creates the BlockchainDeposit row and advances the
// order's escrow state per Plan Section 43/44.
func recordCanonicalDeposit(gc *sharedconfig.GlobalConfig, order *p2pModels.Order, escrowAddress, sender, txHash, detectionSource string) error {
	var existing p2pModels.BlockchainDeposit
	if err := gc.DB.Where("transaction_hash = ?", txHash).First(&existing).Error; err == nil {
		// already recorded (e.g. reconciliation sweep re-ran) - no-op
		return nil
	}

	deposit := p2pModels.BlockchainDeposit{
		ID:              gc.GenerateUUIDString(),
		OrderID:         order.ID,
		Sender:          sender,
		Token:           order.Asset,
		ContractAddress: order.AssetContractAddress,
		Amount:          order.SellerEscrowAssetAmount,
		TransactionHash: txHash,
		IsCanonical:     false,
		DetectionSource: detectionSource,
	}
	if err := gc.DB.Omit(clause.Associations).Create(&deposit).Error; err != nil {
		return &tErrors.ErrorTemporaryServerError{}
	}
	RecordAuditEvent(gc, order.ID, order.OfferID, p2pModels.EventEscrowDepositDetected, sender, deposit)

	return applyDepositToOrder(gc, order, deposit, decimal.RequireFromString(order.SellerEscrowAssetAmount))
}

// applyDepositToOrder implements Plan Section 43/44's escrow-state machine.
func applyDepositToOrder(gc *sharedconfig.GlobalConfig, order *p2pModels.Order, deposit p2pModels.BlockchainDeposit, depositedAmount decimal.Decimal) error {
	expected := decimal.RequireFromString(order.ExpectedEscrowAmount)
	runningTotal := decimal.RequireFromString(order.DepositedEscrowAmount).Add(depositedAmount)

	updates := map[string]interface{}{
		"deposited_escrow_amount": runningTotal.String(),
	}

	switch {
	case runningTotal.LessThan(expected):
		updates["escrow_deposit_status"] = p2pModels.EscrowDepositStatusPartial
	case runningTotal.Equal(expected):
		updates["escrow_deposit_status"] = p2pModels.EscrowDepositStatusConfirmed
		updates["order_status"] = p2pModels.OrderStatusAwaitingPayment
		updates["escrow_deposit_transaction_hash"] = deposit.TransactionHash
		updates["expires_at"] = time.Now().UTC().Add(paymentWindow)
		if err := gc.DB.Model(&p2pModels.BlockchainDeposit{}).Where("id = ?", deposit.ID).Update("is_canonical", true).Error; err != nil {
			return &tErrors.ErrorTemporaryServerError{}
		}
	default: // overpaid
		overpaidAmount := runningTotal.Sub(expected)
		updates["escrow_deposit_status"] = p2pModels.EscrowDepositStatusOverpaid
		updates["order_status"] = p2pModels.OrderStatusAwaitingPayment
		updates["escrow_deposit_transaction_hash"] = deposit.TransactionHash
		updates["refundable_amount"] = overpaidAmount.String()
		updates["expires_at"] = time.Now().UTC().Add(paymentWindow)
		if err := gc.DB.Model(&p2pModels.BlockchainDeposit{}).Where("id = ?", deposit.ID).Update("is_canonical", true).Error; err != nil {
			return &tErrors.ErrorTemporaryServerError{}
		}
		refund := p2pModels.Refund{
			ID:              gc.GenerateUUIDString(),
			DepositID:       deposit.ID,
			OrderID:         order.ID,
			Sender:          deposit.Sender,
			Token:           deposit.Token,
			ContractAddress: deposit.ContractAddress,
			Amount:          overpaidAmount.String(),
			Reason:          p2pModels.RefundReasonOverpayment,
		}
		if err := gc.DB.Omit(clause.Associations).Create(&refund).Error; err != nil {
			return &tErrors.ErrorTemporaryServerError{}
		}
		RecordAuditEvent(gc, order.ID, order.OfferID, p2pModels.EventRefundIssued, deposit.Sender, refund)
	}

	if err := gc.DB.Model(&p2pModels.Order{}).Where("id = ?", order.ID).Updates(updates).Error; err != nil {
		return &tErrors.ErrorTemporaryServerError{}
	}
	if statusVal, ok := updates["order_status"]; ok && statusVal == p2pModels.OrderStatusAwaitingPayment {
		RecordAuditEvent(gc, order.ID, order.OfferID, p2pModels.EventEscrowConfirmed, "", updates)
		NotifyUsername(gc, order.CustomerUsername, "Trovo P2P: Escrow Confirmed", "Escrow deposit confirmed. Please proceed to payment.", map[string]string{"orderId": order.ID, "type": "P2P_ESCROW_CONFIRMED"})
		NotifyUsername(gc, order.MerchantUsername, "Trovo P2P: Escrow Confirmed", "Escrow deposit confirmed for your order.", map[string]string{"orderId": order.ID, "type": "P2P_ESCROW_CONFIRMED"})
	}
	return nil
}

// ReconcileEscrowDepositsFromPaymentHistory detects deposits made by a third
// party who opened the escrow shortlink independently (Plan Section 33/40)
// through the existing generic Payment endpoint, which has no knowledge of
// P2PEscrowPaymentRequest/Order at all. There is no on-chain Deposit Router
// indexer yet (Plan Section 42 - greenfield, parallel smart-contract track),
// so this reconciles against the existing PaymentHistory table instead,
// matching on Memo=Order.ID and destination=escrow address - both already
// guaranteed by Section 19/30's memo convention. Intended to be called by
// the order-detail/status endpoint (on-demand) and a periodic sweep
// (background goroutine started from controllers/main.go's Init).
func ReconcileEscrowDepositsFromPaymentHistory(gc *sharedconfig.GlobalConfig, order *p2pModels.Order) error {
	if order.OrderStatus != p2pModels.OrderStatusAwaitingEscrowDeposit {
		return nil
	}
	escrowAddress := EscrowWalletAddress()
	if escrowAddress == "" {
		return nil
	}
	var matches []paymentModels.PaymentHistory
	if err := gc.DB.Where("memo = ? AND LOWER(to_address) = LOWER(?)", order.ID, escrowAddress).Find(&matches).Error; err != nil {
		return err
	}
	for _, m := range matches {
		var existing p2pModels.BlockchainDeposit
		if err := gc.DB.Where("transaction_hash = ?", m.TransactionID).First(&existing).Error; err == nil {
			continue // already recorded
		}
		if err := recordCanonicalDeposit(gc, order, escrowAddress, m.FromAddress, m.TransactionID, "PAYMENT_HISTORY_MATCH"); err != nil {
			return err
		}
		// re-fetch in case a prior iteration advanced order state
		refreshed, err := GetOrderByID(gc.DB, order.ID)
		if err == nil {
			*order = refreshed
		}
	}
	return nil
}

// RunEscrowReconciliationSweep runs ReconcileEscrowDepositsFromPaymentHistory
// against every order currently AWAITING_ESCROW_DEPOSIT. Meant to be called
// periodically (see controllers/main.go's Init).
func RunEscrowReconciliationSweep(gc *sharedconfig.GlobalConfig) {
	var orders []p2pModels.Order
	if err := gc.DB.Where("order_status = ?", p2pModels.OrderStatusAwaitingEscrowDeposit).Find(&orders).Error; err != nil {
		return
	}
	for i := range orders {
		_ = ReconcileEscrowDepositsFromPaymentHistory(gc, &orders[i])
	}
}

// FindCanonicalDepositForOrder returns the canonical BlockchainDeposit for
// an order - the actual on-chain sender/token/contract address of record,
// which is what any refund back to the depositor must use. Order.
// AssetDepositor is a user id, not a wallet address (Plan Section 11's role
// mapping), and the depositor may not even be a Trovo user (Section 33's
// third-party payer) - the canonical deposit row is the only place the
// real paying wallet address is recorded.
func FindCanonicalDepositForOrder(gc *sharedconfig.GlobalConfig, orderID string) (p2pModels.BlockchainDeposit, error) {
	var deposit p2pModels.BlockchainDeposit
	err := gc.DB.Where("order_id = ? AND is_canonical = ?", orderID, true).First(&deposit).Error
	return deposit, err
}

// ClaimRefund implements Plan Section 45's claim sequence: verify caller is
// the refund's sender, verify not already claimed, mark claimed before
// transferring (checks-effects-interactions/reentrancy protection), then
// transfer the exact amount back via the same Safe-signing path settlement
// uses (Plan Section 59) - a refund is a release from the same escrow Safe,
// just to the depositor instead of the asset recipient.
func ClaimRefund(gc *sharedconfig.GlobalConfig, refundID, callerAddress string) (p2pModels.Refund, error) {
	var refund p2pModels.Refund
	if err := gc.DB.Where("id = ?", refundID).First(&refund).Error; err != nil {
		return refund, &tErrors.CustomError{Param: "refundId", Err: "error-refund-not-found", ErrMessage: "Refund not found"}
	}
	if !strings.EqualFold(refund.Sender, callerAddress) {
		return refund, &tErrors.CustomError{Param: "refundId", Err: "error-forbidden", ErrMessage: "You are not the sender of this refund", Code: 403}
	}
	if refund.Claimed {
		return refund, &tErrors.CustomError{Param: "refundId", Err: "error-refund-already-claimed", ErrMessage: "This refund has already been claimed"}
	}

	now := time.Now().UTC()
	if err := gc.DB.Model(&p2pModels.Refund{}).Where("id = ? AND claimed = ?", refund.ID, false).
		Updates(map[string]interface{}{"claimed": true, "claimed_at": now}).Error; err != nil {
		return refund, &tErrors.ErrorTemporaryServerError{}
	}
	refund.Claimed = true
	refund.ClaimedAt = &now

	txHash, err := safesigner.ExecuteTransfers([]safesigner.Transfer{
		{Token: refund.ContractAddress, Recipient: refund.Sender, Amount: toWei(refund.Amount)},
	})
	if err != nil {
		// Claimed already flipped to true to prevent a double-submit race;
		// operationally this needs a retry/replay path if the transfer
		// itself fails after the claim flag is set - flagged here rather
		// than silently reverting the flag, which would reopen the
		// reentrancy window this ordering exists to close.
		RecordAuditEvent(gc, refund.OrderID, "", "REFUND_CLAIM_TRANSFER_FAILED", refund.Sender, map[string]string{"refundId": refund.ID, "error": err.Error()})
		return refund, &tErrors.CustomError{Param: "refund", Err: "error-refund-transfer-failed", ErrMessage: "Refund transfer failed. Please contact support."}
	}

	if err := gc.DB.Model(&p2pModels.Refund{}).Where("id = ?", refund.ID).Update("transaction_hash", txHash).Error; err != nil {
		return refund, &tErrors.ErrorTemporaryServerError{}
	}
	refund.TransactionHash = txHash
	RecordAuditEvent(gc, refund.OrderID, "", p2pModels.EventRefundClaimed, refund.Sender, refund)
	return refund, nil
}

// ListMyRefunds returns every refund owed to a wallet address, matching
// Plan Section 45's model - Refund.Sender is a wallet address (the actual
// on-chain depositor), not a user id, so refunds are looked up by address,
// the same identifier /v1/p2p/escrow-deposit already authenticates by.
func ListMyRefunds(db *gorm.DB, callerAddress string) ([]p2pModels.Refund, error) {
	var refunds []p2pModels.Refund
	err := db.Where("LOWER(sender) = LOWER(?)", callerAddress).Order("created_at desc").Find(&refunds).Error
	return refunds, err
}

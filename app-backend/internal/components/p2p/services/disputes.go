package p2p

import (
	"time"
	p2pModels "trovo-wallet-api/internal/components/p2p/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/sharedconfig"

	"gorm.io/gorm/clause"
)

var validDisputeSubjects = map[string]bool{
	p2pModels.DisputeSubjectNoPayment:                true,
	p2pModels.DisputeSubjectUnderPayment:              true,
	p2pModels.DisputeSubjectOverPayment:               true,
	p2pModels.DisputeSubjectPaymentNotReceived:        true,
	p2pModels.DisputeSubjectPaymentMarkedSentInError:  true,
	p2pModels.DisputeSubjectWrongPaymentAmount:        true,
	p2pModels.DisputeSubjectOther:                     true,
}

// OpenDispute implements Plan Section 62. Opening a dispute is not
// automatically treated as misconduct by either party (Section 62's own
// note) - this only flags the order for review, it does not penalize anyone.
func OpenDispute(gc *sharedconfig.GlobalConfig, orderID, openedByUserID, subject, description string, evidence []string) (p2pModels.Dispute, error) {
	order, err := GetOrderByID(gc.DB, orderID)
	if err != nil {
		return p2pModels.Dispute{}, &tErrors.CustomError{Param: "orderId", Err: "error-order-not-found", ErrMessage: "Order not found"}
	}
	if openedByUserID != order.CustomerUserID && openedByUserID != order.MerchantUserID {
		return p2pModels.Dispute{}, &tErrors.CustomError{Param: "orderId", Err: "error-forbidden", ErrMessage: "You are not a party to this order", Code: 403}
	}
	if order.OrderStatus == p2pModels.OrderStatusAwaitingApproval || order.OrderStatus == p2pModels.OrderStatusAwaitingEscrowDeposit {
		return p2pModels.Dispute{}, &tErrors.CustomError{Param: "orderId", Err: "error-invalid-order-state", ErrMessage: "Disputes can only be raised once escrow has been deposited"}
	}
	if order.OrderStatus == p2pModels.OrderStatusCompleted {
		return p2pModels.Dispute{}, &tErrors.CustomError{Param: "orderId", Err: "error-invalid-order-state", ErrMessage: "This order is already completed and cannot be disputed"}
	}
	if order.IsDisputed {
		return p2pModels.Dispute{}, &tErrors.CustomError{Param: "orderId", Err: "error-dispute-already-open", ErrMessage: "This order already has an open dispute"}
	}
	if !validDisputeSubjects[subject] {
		return p2pModels.Dispute{}, &tErrors.CustomError{Param: "subject", Err: "error-invalid-dispute-subject", ErrMessage: "Invalid dispute subject"}
	}

	now := time.Now().UTC()
	dispute := p2pModels.Dispute{
		ID:          gc.GenerateUUIDString(),
		OrderID:     order.ID,
		OpenedBy:    openedByUserID,
		OpenedAt:    now,
		Subject:     subject,
		Description: description,
		Evidence:    evidence,
		Status:      p2pModels.DisputeStatusOpen,
	}
	dbTX := gc.DB.Begin()
	if err := dbTX.Omit(clause.Associations).Create(&dispute).Error; err != nil {
		dbTX.Rollback()
		return dispute, &tErrors.ErrorTemporaryServerError{}
	}
	if err := dbTX.Model(&p2pModels.Order{}).Where("id = ?", order.ID).Update("is_disputed", true).Error; err != nil {
		dbTX.Rollback()
		return dispute, &tErrors.ErrorTemporaryServerError{}
	}
	if err := dbTX.Commit().Error; err != nil {
		return dispute, &tErrors.ErrorTemporaryServerError{}
	}
	RecordAuditEvent(gc, order.ID, order.OfferID, p2pModels.EventDisputeOpened, openedByUserID, dispute)
	RecordDisputeOpened(gc, order)

	counterparty := order.MerchantUserID
	if openedByUserID == order.MerchantUserID {
		counterparty = order.CustomerUserID
	}
	notifyByUserID(gc, order, counterparty, "Trovo P2P: Dispute Opened", "A dispute has been opened on your order.")
	return dispute, nil
}

// GetOpenDisputeForOrder returns the currently-open dispute for an order, if
// any - lets a client that only knows Order.IsDisputed=true (the API never
// otherwise surfaces a dispute id) resolve which dispute to act on.
func GetOpenDisputeForOrder(gc *sharedconfig.GlobalConfig, orderID string) (p2pModels.Dispute, error) {
	var dispute p2pModels.Dispute
	err := gc.DB.Where("order_id = ? AND status = ?", orderID, p2pModels.DisputeStatusOpen).
		Order("opened_at desc").First(&dispute).Error
	return dispute, err
}

// MerchantConfirmsPayment is Section 63's first self-resolution path: the
// merchant confirms fiat payment was in fact received (e.g. a dispute was
// opened over a delayed notification) - the dispute closes and the order
// continues to release.
func MerchantConfirmsPayment(gc *sharedconfig.GlobalConfig, disputeID, merchantUserID string) (p2pModels.Order, error) {
	dispute, order, err := loadDisputeAndOrderForResolution(gc, disputeID, merchantUserID, true)
	if err != nil {
		return order, err
	}
	if err := closeDispute(gc, &dispute, order, p2pModels.DisputeResolutionInFavorOfSeller, merchantUserID); err != nil {
		return order, err
	}
	return releaseAndCompleteOrder(gc, &order)
}

// BuyerConfirmsPaymentNotMade is Section 63's second self-resolution path:
// the buyer confirms they had not actually made the fiat payment - the
// dispute closes and the order returns to AWAITING_PAYMENT.
func BuyerConfirmsPaymentNotMade(gc *sharedconfig.GlobalConfig, disputeID, buyerUserID string) (p2pModels.Order, error) {
	dispute, order, err := loadDisputeAndOrderForResolution(gc, disputeID, buyerUserID, false)
	if err != nil {
		return order, err
	}
	if err := closeDispute(gc, &dispute, order, p2pModels.DisputeResolutionInFavorOfSeller, buyerUserID); err != nil {
		return order, err
	}
	order.OrderStatus = p2pModels.OrderStatusAwaitingPayment
	if err := gc.DB.Model(&p2pModels.Order{}).Where("id = ?", order.ID).
		Updates(map[string]interface{}{"order_status": p2pModels.OrderStatusAwaitingPayment, "is_disputed": false}).Error; err != nil {
		return order, &tErrors.ErrorTemporaryServerError{}
	}
	notifyByUserID(gc, order, order.FiatRecipient, "Trovo P2P: Dispute Resolved", "The dispute was resolved - the order is back to awaiting payment.")
	return order, nil
}

// AdminResolveDispute is the arbiter-driven path when the two self-
// resolution shortcuts above don't apply. IN_FAVOR_OF_SELLER continues the
// order to release; IN_FAVOR_OF_BUYER refunds the escrow deposit back to
// its depositor and cancels the order; SPLIT is recorded for manual/off-
// system follow-up (the plan does not specify an automated split-settlement
// mechanism).
func AdminResolveDispute(gc *sharedconfig.GlobalConfig, disputeID, resolution, resolvedByAdminID string) (p2pModels.Order, error) {
	var dispute p2pModels.Dispute
	if err := gc.DB.Where("id = ?", disputeID).First(&dispute).Error; err != nil {
		return p2pModels.Order{}, &tErrors.CustomError{Param: "disputeId", Err: "error-dispute-not-found", ErrMessage: "Dispute not found"}
	}
	order, err := GetOrderByID(gc.DB, dispute.OrderID)
	if err != nil {
		return order, &tErrors.CustomError{Param: "orderId", Err: "error-order-not-found", ErrMessage: "Order not found"}
	}
	if err := closeDispute(gc, &dispute, order, resolution, resolvedByAdminID); err != nil {
		return order, err
	}
	switch resolution {
	case p2pModels.DisputeResolutionInFavorOfSeller:
		return releaseAndCompleteOrder(gc, &order)
	case p2pModels.DisputeResolutionInFavorOfBuyer:
		order.OrderStatus = p2pModels.OrderStatusCancelled
		order.RefundableAmount = order.DepositedEscrowAmount
		if err := gc.DB.Save(&order).Error; err != nil {
			return order, &tErrors.ErrorTemporaryServerError{}
		}
		// Order.AssetDepositor is a user id, not a wallet address (Plan
		// Section 11) - the canonical deposit record has the real paying
		// wallet address a refund must go to.
		deposit, err := FindCanonicalDepositForOrder(gc, order.ID)
		if err != nil {
			return order, &tErrors.CustomError{Param: "order", Err: "error-no-canonical-deposit", ErrMessage: "Could not find the escrow deposit to refund"}
		}
		refund := p2pModels.Refund{
			ID:              gc.GenerateUUIDString(),
			DepositID:       deposit.ID,
			OrderID:         order.ID,
			Sender:          deposit.Sender,
			Token:           deposit.Token,
			ContractAddress: deposit.ContractAddress,
			Amount:          order.DepositedEscrowAmount,
			Reason:          p2pModels.RefundReasonOther,
		}
		if err := gc.DB.Omit(clause.Associations).Create(&refund).Error; err != nil {
			return order, &tErrors.ErrorTemporaryServerError{}
		}
		RecordAuditEvent(gc, order.ID, order.OfferID, p2pModels.EventRefundIssued, resolvedByAdminID, refund)
		return order, nil
	default:
		return order, nil
	}
}

func loadDisputeAndOrderForResolution(gc *sharedconfig.GlobalConfig, disputeID, callerUserID string, callerMustBeMerchant bool) (p2pModels.Dispute, p2pModels.Order, error) {
	var dispute p2pModels.Dispute
	if err := gc.DB.Where("id = ?", disputeID).First(&dispute).Error; err != nil {
		return dispute, p2pModels.Order{}, &tErrors.CustomError{Param: "disputeId", Err: "error-dispute-not-found", ErrMessage: "Dispute not found"}
	}
	order, err := GetOrderByID(gc.DB, dispute.OrderID)
	if err != nil {
		return dispute, order, &tErrors.CustomError{Param: "orderId", Err: "error-order-not-found", ErrMessage: "Order not found"}
	}
	if dispute.Status != p2pModels.DisputeStatusOpen {
		return dispute, order, &tErrors.CustomError{Param: "disputeId", Err: "error-dispute-already-resolved", ErrMessage: "This dispute is already resolved"}
	}
	if callerMustBeMerchant && callerUserID != order.MerchantUserID {
		return dispute, order, &tErrors.CustomError{Param: "disputeId", Err: "error-forbidden", ErrMessage: "You are not the merchant on this order", Code: 403}
	}
	if !callerMustBeMerchant && callerUserID != order.CustomerUserID {
		return dispute, order, &tErrors.CustomError{Param: "disputeId", Err: "error-forbidden", ErrMessage: "You are not the customer on this order", Code: 403}
	}
	return dispute, order, nil
}

func closeDispute(gc *sharedconfig.GlobalConfig, dispute *p2pModels.Dispute, order p2pModels.Order, resolution, resolvedBy string) error {
	now := time.Now().UTC()
	resolutionTime := now.Sub(dispute.OpenedAt)
	dispute.Status = p2pModels.DisputeStatusResolved
	dispute.Resolution = resolution
	dispute.ResolvedAt = &now
	dispute.ResolvedBy = resolvedBy
	dispute.ResolutionTimeSeconds = int64(resolutionTime.Seconds())
	if err := gc.DB.Save(dispute).Error; err != nil {
		return &tErrors.ErrorTemporaryServerError{}
	}
	if err := gc.DB.Model(&p2pModels.Order{}).Where("id = ?", dispute.OrderID).Update("is_disputed", false).Error; err != nil {
		return &tErrors.ErrorTemporaryServerError{}
	}
	RecordAuditEvent(gc, dispute.OrderID, "", p2pModels.EventDisputeResolved, resolvedBy, dispute)
	RecordDisputeResolved(gc, order, resolution, resolutionTime)
	return nil
}

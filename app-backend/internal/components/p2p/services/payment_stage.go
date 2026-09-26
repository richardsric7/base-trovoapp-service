package p2p

import (
	"time"
	p2pModels "trovo-wallet-api/internal/components/p2p/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/sharedconfig"
)

// MarkFiatPaymentSent is the buyer's "I've sent payment" action (Plan
// Section 24: AWAITING_PAYMENT -> AWAITING_PAYMENT_CONFIRMATION).
func MarkFiatPaymentSent(gc *sharedconfig.GlobalConfig, orderID, callerUserID string) (p2pModels.Order, error) {
	order, err := GetOrderByID(gc.DB, orderID)
	if err != nil {
		return order, &tErrors.CustomError{Param: "orderId", Err: "error-order-not-found", ErrMessage: "Order not found"}
	}
	if order.FiatPayer != callerUserID {
		return order, &tErrors.CustomError{Param: "orderId", Err: "error-forbidden", ErrMessage: "You are not the fiat payer on this order", Code: 403}
	}
	if order.OrderStatus != p2pModels.OrderStatusAwaitingPayment {
		return order, &tErrors.CustomError{Param: "orderId", Err: "error-invalid-order-state", ErrMessage: "This order is not awaiting fiat payment"}
	}
	confirmationDeadline := time.Now().UTC().Add(paymentConfirmationWindow)
	order.OrderStatus = p2pModels.OrderStatusAwaitingPaymentConfirmation
	order.ExpiresAt = &confirmationDeadline
	if err := gc.DB.Save(&order).Error; err != nil {
		return order, &tErrors.ErrorTemporaryServerError{}
	}
	RecordAuditEvent(gc, order.ID, order.OfferID, p2pModels.EventFiatPaymentMarkedSent, callerUserID, order)
	notifyByUserID(gc, order, order.FiatRecipient, "Trovo P2P: Payment Sent", "The counterparty marked the fiat payment as sent. Please confirm receipt to release the asset.")
	return order, nil
}

// ConfirmFiatPaymentReceived is the seller's confirmation action, which
// triggers the Universal Safe settlement release (Plan Sections 60-61) and
// completes the order.
func ConfirmFiatPaymentReceived(gc *sharedconfig.GlobalConfig, orderID, callerUserID string) (p2pModels.Order, error) {
	order, err := GetOrderByID(gc.DB, orderID)
	if err != nil {
		return order, &tErrors.CustomError{Param: "orderId", Err: "error-order-not-found", ErrMessage: "Order not found"}
	}
	if order.FiatRecipient != callerUserID {
		return order, &tErrors.CustomError{Param: "orderId", Err: "error-forbidden", ErrMessage: "You are not the fiat recipient on this order", Code: 403}
	}
	if order.OrderStatus != p2pModels.OrderStatusAwaitingPaymentConfirmation {
		return order, &tErrors.CustomError{Param: "orderId", Err: "error-invalid-order-state", ErrMessage: "This order is not awaiting payment confirmation"}
	}

	now := time.Now().UTC()
	order.PaymentConfirmedAt = &now
	RecordAuditEvent(gc, order.ID, order.OfferID, p2pModels.EventFiatPaymentConfirmed, callerUserID, order)

	return releaseAndCompleteOrder(gc, &order)
}

// releaseAndCompleteOrder calls the settlement/release service and, on
// success, marks the order COMPLETED and updates performance aggregates
// (Plan Section 61).
func releaseAndCompleteOrder(gc *sharedconfig.GlobalConfig, order *p2pModels.Order) (p2pModels.Order, error) {
	txHash, err := ReleaseEscrowSettlement(gc, order)
	if err != nil {
		return *order, err
	}
	now := time.Now().UTC()
	order.AssetReleaseTransactionHash = txHash
	order.AssetReleasedAt = &now
	order.CompletedAt = &now
	order.OrderStatus = p2pModels.OrderStatusCompleted
	if err := gc.DB.Save(order).Error; err != nil {
		return *order, &tErrors.ErrorTemporaryServerError{}
	}
	RecordAuditEvent(gc, order.ID, order.OfferID, p2pModels.EventAssetReleased, "", map[string]string{"transactionHash": txHash})
	RecordAuditEvent(gc, order.ID, order.OfferID, p2pModels.EventOrderCompleted, "", order)
	UpdatePerformanceOnCompletion(gc, *order)
	NotifyUsername(gc, order.CustomerUsername, "Trovo P2P: Order Completed", "Your P2P order has been completed.", map[string]string{"orderId": order.ID, "type": "P2P_ORDER_COMPLETED"})
	NotifyUsername(gc, order.MerchantUsername, "Trovo P2P: Order Completed", "Your P2P order has been completed.", map[string]string{"orderId": order.ID, "type": "P2P_ORDER_COMPLETED"})
	return *order, nil
}

// notifyByUserID resolves a userID to a username on the order (customer or
// merchant) and notifies them - both roles are always one of these two.
func notifyByUserID(gc *sharedconfig.GlobalConfig, order p2pModels.Order, userID, title, body string) {
	switch userID {
	case order.CustomerUserID:
		NotifyUsername(gc, order.CustomerUsername, title, body, map[string]string{"orderId": order.ID})
	case order.MerchantUserID:
		NotifyUsername(gc, order.MerchantUsername, title, body, map[string]string{"orderId": order.ID})
	}
}

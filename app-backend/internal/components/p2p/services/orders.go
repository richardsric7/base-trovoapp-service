package p2p

import (
	"encoding/json"
	"time"
	p2pModels "trovo-wallet-api/internal/components/p2p/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// merchantActiveOrderCap and customerActiveOrderCap implement Plan Section 16.
const (
	merchantActiveOrderCap = 5
	customerActiveOrderCap = 1
	// defaultApprovalWindow is how long a customer's order waits for
	// merchant accept/reject before it auto-expires (Plan Section 24).
	defaultApprovalWindow = 30 * time.Minute
	// escrowDepositWindow is how long the depositor has to fund escrow
	// once the merchant accepts, before the reservation is released.
	escrowDepositWindow = 30 * time.Minute
	// paymentWindow is how long the fiat payer has to send payment once
	// escrow is confirmed, before the deposited asset is refunded back to
	// its depositor.
	paymentWindow = 60 * time.Minute
	// paymentConfirmationWindow is how long the fiat recipient has to
	// confirm receipt once the payer marks payment sent, before the order
	// is auto-escalated into a dispute for review rather than any funds
	// moving unilaterally.
	paymentConfirmationWindow = 60 * time.Minute
)

// activeOrderStatuses are the statuses that count toward capacity limits -
// anything still "in flight" and not yet terminal.
var activeOrderStatuses = []string{
	p2pModels.OrderStatusAwaitingApproval,
	p2pModels.OrderStatusAwaitingEscrowDeposit,
	p2pModels.OrderStatusAwaitingPayment,
	p2pModels.OrderStatusAwaitingPaymentConfirmation,
}

// CreateOrderInput is the client-supplied shape for creating an order.
type CreateOrderInput struct {
	OfferID              string
	SpecifiedAssetAmount string
}

// QuoteOrderFees runs the same offer/amount validation and fee calculation
// CreateOrder does, without persisting anything - lets the order-creation
// screen show a real itemized fee breakdown before the customer commits
// (Plan Section 97.6), rather than only after Order is actually created.
func QuoteOrderFees(gc *sharedconfig.GlobalConfig, offerID, specifiedAssetAmount string) (OrderFeeBreakdown, p2pModels.Offer, error) {
	offer, err := GetOfferByID(gc.DB, offerID)
	if err != nil {
		return OrderFeeBreakdown{}, offer, &tErrors.CustomError{Param: "offerId", Err: "error-offer-not-found", ErrMessage: "Offer not found"}
	}
	specifiedAmount, e := decimal.NewFromString(specifiedAssetAmount)
	if e != nil || specifiedAmount.LessThanOrEqual(decimal.Zero) {
		return OrderFeeBreakdown{}, offer, &tErrors.CustomError{Param: "specifiedAssetAmount", Err: "error-invalid-amount", ErrMessage: "specifiedAssetAmount must be a positive decimal number"}
	}
	feeCfg, err := GetActiveFeeConfiguration(gc.DB, offer.CountryCode)
	if err != nil {
		return OrderFeeBreakdown{}, offer, &tErrors.CustomError{Param: "offerId", Err: "error-no-fee-configuration", ErrMessage: "No fee configuration is available for this offer's country"}
	}
	price := decimal.RequireFromString(offer.Price)
	return CalculateOrderFees(specifiedAmount, price, feeCfg), offer, nil
}

// CreateOrder implements Plan Section 25's full creation sequence: validate
// offer, validate customer eligibility/self-trade/capacity, load curated
// asset (already snapshotted on the Offer), calculate fees/VAT, snapshot
// payment info, generate Order.id, and create as AWAITING_APPROVAL.
func CreateOrder(gc *sharedconfig.GlobalConfig, customerUserID, customerUsername string, in CreateOrderInput) (p2pModels.Order, error) {
	offer, err := GetOfferByID(gc.DB, in.OfferID)
	if err != nil {
		return p2pModels.Order{}, &tErrors.CustomError{Param: "offerId", Err: "error-offer-not-found", ErrMessage: "Offer not found"}
	}
	if offer.Status != p2pModels.OfferStatusActive || offer.AvailabilityStatus != p2pModels.OfferAvailabilityOnline {
		return p2pModels.Order{}, &tErrors.CustomError{Param: "offerId", Err: "error-offer-not-available", ErrMessage: "This offer is not currently available"}
	}

	// Self-trade prevention (Plan Section 12) - backend-enforced, never
	// trust a client-side check alone.
	if offer.MerchantUserID == customerUserID {
		return p2pModels.Order{}, &tErrors.CustomError{Param: "offerId", Err: "error-self-trade-forbidden", ErrMessage: "You cannot trade against your own offer"}
	}

	specifiedAmount, e := decimal.NewFromString(in.SpecifiedAssetAmount)
	if e != nil || specifiedAmount.LessThanOrEqual(decimal.Zero) {
		return p2pModels.Order{}, &tErrors.CustomError{Param: "specifiedAssetAmount", Err: "error-invalid-amount", ErrMessage: "specifiedAssetAmount must be a positive decimal number"}
	}
	minAmt := decimal.RequireFromString(offer.MinOrderAmount)
	maxAmt := decimal.RequireFromString(offer.MaxOrderAmount)
	if specifiedAmount.LessThan(minAmt) || specifiedAmount.GreaterThan(maxAmt) {
		return p2pModels.Order{}, &tErrors.CustomError{Param: "specifiedAssetAmount", Err: "error-amount-out-of-range", ErrMessage: "Amount is outside the offer's allowed range"}
	}

	availableLiquidity := decimal.RequireFromString(offer.AvailableLiquidity)
	if availableLiquidity.LessThan(specifiedAmount) {
		return p2pModels.Order{}, &tErrors.CustomError{Param: "specifiedAssetAmount", Err: "error-insufficient-liquidity", ErrMessage: "Offer does not have enough available liquidity for this amount"}
	}

	// Capacity checks (Plan Section 16). Disputed orders don't count toward
	// the merchant's 5-order operational cap.
	var merchantActiveCount int64
	if err := gc.DB.Model(&p2pModels.Order{}).
		Where("merchant_user_id = ? AND order_status IN ? AND is_disputed = ?", offer.MerchantUserID, activeOrderStatuses, false).
		Count(&merchantActiveCount).Error; err != nil {
		return p2pModels.Order{}, &tErrors.ErrorTemporaryServerError{}
	}
	if merchantActiveCount >= merchantActiveOrderCap {
		return p2pModels.Order{}, &tErrors.CustomError{Param: "offerId", Err: "error-merchant-at-capacity", ErrMessage: "This merchant has reached their open order limit. Please try again shortly."}
	}
	var customerActiveCount int64
	if err := gc.DB.Model(&p2pModels.Order{}).
		Where("customer_user_id = ? AND order_status IN ?", customerUserID, activeOrderStatuses).
		Count(&customerActiveCount).Error; err != nil {
		return p2pModels.Order{}, &tErrors.ErrorTemporaryServerError{}
	}
	if customerActiveCount >= customerActiveOrderCap {
		return p2pModels.Order{}, &tErrors.CustomError{Param: "customer", Err: "error-customer-at-capacity", ErrMessage: "You already have an open P2P order. Complete or cancel it before starting another."}
	}
	var customerDisputedCount int64
	if err := gc.DB.Model(&p2pModels.Order{}).
		Where("customer_user_id = ? AND is_disputed = ?", customerUserID, true).
		Count(&customerDisputedCount).Error; err != nil {
		return p2pModels.Order{}, &tErrors.ErrorTemporaryServerError{}
	}
	if customerDisputedCount > 0 {
		return p2pModels.Order{}, &tErrors.CustomError{Param: "customer", Err: "error-customer-has-active-dispute", ErrMessage: "You have an active dispute and cannot open a new order until it is resolved."}
	}

	feeCfg, err := GetActiveFeeConfiguration(gc.DB, offer.CountryCode)
	if err != nil {
		return p2pModels.Order{}, &tErrors.CustomError{Param: "offerId", Err: "error-no-fee-configuration", ErrMessage: "No fee configuration is available for this offer's country"}
	}
	feeWalletCfg, err := GetActiveFeeWalletConfiguration(gc.DB, offer.CountryCode)
	if err != nil {
		return p2pModels.Order{}, &tErrors.CustomError{Param: "offerId", Err: "error-no-fee-wallet-configuration", ErrMessage: "No fee wallet configuration is available for this offer's country"}
	}

	price := decimal.RequireFromString(offer.Price)
	breakdown := CalculateOrderFees(specifiedAmount, price, feeCfg)

	// Role mapping (Plan Section 11).
	assetDepositor, assetRecipient, fiatPayer, fiatRecipient := resolveTradeRoles(offer.OfferType, offer.MerchantUserID, customerUserID)

	feeCfgSnapshot, _ := json.Marshal(feeCfg)

	expiresAt := time.Now().UTC().Add(defaultApprovalWindow)

	order := p2pModels.Order{
		ID:                   gc.GenerateUUIDString(),
		OfferID:              offer.ID,
		CustomerUserID:       customerUserID,
		CustomerUsername:     customerUsername,
		MerchantUserID:       offer.MerchantUserID,
		MerchantUsername:     offer.MerchantUsername,
		OfferType:            offer.OfferType,
		Asset:                offer.Asset,
		AssetContractAddress: offer.ContractAddress,
		PaymentMethodSnapshot: offer.PaymentMethod,
		Country:              offer.Country,
		CountryCode:          offer.CountryCode,
		Currency:             offer.Currency,
		Price:                offer.Price,
		SpecifiedAssetAmount: specifiedAmount.Truncate(decimalPlaces).String(),
		PaymentAmount:        breakdown.PaymentAmount.String(),

		BuyerPlatformFee:      breakdown.BuyerPlatformFee.String(),
		BuyerRegulatoryFee:    breakdown.BuyerRegulatoryFee.String(),
		BuyerPlatformFeeVat:   breakdown.BuyerPlatformFeeVat.String(),
		BuyerRegulatoryFeeVat: breakdown.BuyerRegulatoryFeeVat.String(),
		BuyerTotalFees:        breakdown.BuyerTotalFees.String(),
		BuyerTotalVat:         breakdown.BuyerTotalVat.String(),
		BuyerTotalCharges:     breakdown.BuyerTotalCharges.String(),
		BuyerNetAssetAmount:   breakdown.BuyerNetAssetAmount.String(),

		SellerPlatformFee:       breakdown.SellerPlatformFee.String(),
		SellerRegulatoryFee:     breakdown.SellerRegulatoryFee.String(),
		SellerPlatformFeeVat:    breakdown.SellerPlatformFeeVat.String(),
		SellerRegulatoryFeeVat:  breakdown.SellerRegulatoryFeeVat.String(),
		SellerTotalFees:         breakdown.SellerTotalFees.String(),
		SellerTotalVat:          breakdown.SellerTotalVat.String(),
		SellerTotalCharges:      breakdown.SellerTotalCharges.String(),
		SellerEscrowAssetAmount: breakdown.SellerEscrowAssetAmount.String(),

		CombinedPlatformFee:   breakdown.CombinedPlatformFee.String(),
		CombinedRegulatoryFee: breakdown.CombinedRegulatoryFee.String(),
		CombinedVat:           breakdown.CombinedVat.String(),

		PlatformFeeWalletAddress:   feeWalletCfg.PlatformFeeWalletAddress,
		RegulatoryFeeWalletAddress: feeWalletCfg.RegulatoryFeeWalletAddress,
		VatWalletAddress:           feeWalletCfg.VatWalletAddress,
		FeeConfigurationVersion:    feeCfg.Version,
		FeeConfigurationSnapshot:   string(feeCfgSnapshot),

		ExpectedEscrowAmount: breakdown.SellerEscrowAssetAmount.String(),
		EscrowDepositStatus:  p2pModels.EscrowDepositStatusNotDeposited,

		AssetDepositor: assetDepositor,
		AssetRecipient: assetRecipient,
		FiatPayer:      fiatPayer,
		FiatRecipient:  fiatRecipient,

		OrderStatus: p2pModels.OrderStatusAwaitingApproval,
		ExpiresAt:   &expiresAt,
	}

	dbTX := gc.DB.Begin()
	if err := dbTX.Omit(clause.Associations).Create(&order).Error; err != nil {
		dbTX.Rollback()
		return p2pModels.Order{}, &tErrors.ErrorTemporaryServerError{}
	}
	// Reserve liquidity immediately on creation so a burst of concurrent
	// orders against the same offer cannot collectively oversell it while
	// awaiting merchant approval.
	if err := dbTX.Model(&p2pModels.Offer{}).Where("id = ?", offer.ID).
		Updates(map[string]interface{}{
			"available_liquidity": availableLiquidity.Sub(specifiedAmount).String(),
			"reserved_liquidity":   decimal.RequireFromString(offer.ReservedLiquidity).Add(specifiedAmount).String(),
		}).Error; err != nil {
		dbTX.Rollback()
		return p2pModels.Order{}, &tErrors.ErrorTemporaryServerError{}
	}
	if err := dbTX.Commit().Error; err != nil {
		return p2pModels.Order{}, &tErrors.ErrorTemporaryServerError{}
	}

	RecordAuditEvent(gc, order.ID, offer.ID, p2pModels.EventOrderCreated, customerUserID, order)
	NotifyUsername(gc, offer.MerchantUsername, "Trovo P2P: New Order", customerUsername+" is waiting for you to accept/reject an order", map[string]string{"orderId": order.ID, "type": "P2P_ORDER_CREATED"})
	return order, nil
}

// resolveTradeRoles implements Plan Section 11's fixed role mapping.
func resolveTradeRoles(offerType, merchantUserID, customerUserID string) (assetDepositor, assetRecipient, fiatPayer, fiatRecipient string) {
	if offerType == p2pModels.OfferTypeSell {
		// merchant is selling the asset: buyer=customer, seller=merchant
		return merchantUserID, customerUserID, customerUserID, merchantUserID
	}
	// BUY: merchant is buying the asset: buyer=merchant, seller=customer
	return customerUserID, merchantUserID, merchantUserID, customerUserID
}

// GetOrderByID fetches a single order.
func GetOrderByID(db *gorm.DB, orderID string) (p2pModels.Order, error) {
	var order p2pModels.Order
	err := db.Where("id = ?", orderID).First(&order).Error
	return order, err
}

// AcceptOrder implements Plan Section 26's Accept path: locks terms
// (already snapshotted at creation), transitions to AWAITING_ESCROW_DEPOSIT.
// Escrow instruction/shortlink generation is handled separately by the
// escrow service once the order is in this state (Section 95/96 - the client
// calls the escrow-deposit endpoint next).
func AcceptOrder(gc *sharedconfig.GlobalConfig, orderID, merchantUserID string) (p2pModels.Order, error) {
	order, err := GetOrderByID(gc.DB, orderID)
	if err != nil {
		return order, &tErrors.CustomError{Param: "orderId", Err: "error-order-not-found", ErrMessage: "Order not found"}
	}
	if order.MerchantUserID != merchantUserID {
		return order, &tErrors.CustomError{Param: "orderId", Err: "error-forbidden", ErrMessage: "You are not the merchant on this order", Code: 403}
	}
	if order.OrderStatus != p2pModels.OrderStatusAwaitingApproval {
		return order, &tErrors.CustomError{Param: "orderId", Err: "error-invalid-order-state", ErrMessage: "This order is no longer awaiting approval"}
	}
	escrowDeadline := time.Now().UTC().Add(escrowDepositWindow)
	order.OrderStatus = p2pModels.OrderStatusAwaitingEscrowDeposit
	order.ExpiresAt = &escrowDeadline
	if err := gc.DB.Save(&order).Error; err != nil {
		return order, &tErrors.ErrorTemporaryServerError{}
	}
	RecordAuditEvent(gc, order.ID, order.OfferID, p2pModels.EventOrderAccepted, merchantUserID, order)

	// Generate the escrow deposit shortlink+QR now, atomically with
	// acceptance (Plan Section 26) - a failure here must not fail the
	// acceptance itself; the client can retry shortlink generation via the
	// escrow-deposit screen if this best-effort attempt fails.
	if err := GenerateEscrowShortlink(gc, &order); err != nil {
		RecordAuditEvent(gc, order.ID, order.OfferID, "SHORTLINK_GENERATION_FAILED", merchantUserID, map[string]string{"error": err.Error()})
	}

	NotifyUsername(gc, order.CustomerUsername, "Trovo P2P: Order Accepted", "Your order has been accepted. Please proceed to escrow deposit.", map[string]string{"orderId": order.ID, "type": "P2P_ORDER_ACCEPTED"})
	return order, nil
}

// RejectOrder implements Plan Section 26's Reject path: releases reserved
// liquidity and records the outcome.
func RejectOrder(gc *sharedconfig.GlobalConfig, orderID, merchantUserID string) (p2pModels.Order, error) {
	order, err := GetOrderByID(gc.DB, orderID)
	if err != nil {
		return order, &tErrors.CustomError{Param: "orderId", Err: "error-order-not-found", ErrMessage: "Order not found"}
	}
	if order.MerchantUserID != merchantUserID {
		return order, &tErrors.CustomError{Param: "orderId", Err: "error-forbidden", ErrMessage: "You are not the merchant on this order", Code: 403}
	}
	if order.OrderStatus != p2pModels.OrderStatusAwaitingApproval {
		return order, &tErrors.CustomError{Param: "orderId", Err: "error-invalid-order-state", ErrMessage: "This order is no longer awaiting approval"}
	}
	order.OrderStatus = p2pModels.OrderStatusRejected
	if err := releaseReservedLiquidityAndSave(gc, &order); err != nil {
		return order, err
	}
	RecordAuditEvent(gc, order.ID, order.OfferID, p2pModels.EventOrderRejected, merchantUserID, order)
	NotifyUsername(gc, order.CustomerUsername, "Trovo P2P: Order Rejected", "Your order was rejected by the merchant.", map[string]string{"orderId": order.ID, "type": "P2P_ORDER_REJECTED"})
	return order, nil
}

// CancelOrder allows the customer to cancel while still AWAITING_APPROVAL or
// AWAITING_ESCROW_DEPOSIT (before any funds have moved).
func CancelOrder(gc *sharedconfig.GlobalConfig, orderID, customerUserID string) (p2pModels.Order, error) {
	order, err := GetOrderByID(gc.DB, orderID)
	if err != nil {
		return order, &tErrors.CustomError{Param: "orderId", Err: "error-order-not-found", ErrMessage: "Order not found"}
	}
	if order.CustomerUserID != customerUserID {
		return order, &tErrors.CustomError{Param: "orderId", Err: "error-forbidden", ErrMessage: "You are not the customer on this order", Code: 403}
	}
	if order.OrderStatus != p2pModels.OrderStatusAwaitingApproval && order.OrderStatus != p2pModels.OrderStatusAwaitingEscrowDeposit {
		return order, &tErrors.CustomError{Param: "orderId", Err: "error-invalid-order-state", ErrMessage: "This order can no longer be cancelled"}
	}
	order.OrderStatus = p2pModels.OrderStatusCancelled
	if err := releaseReservedLiquidityAndSave(gc, &order); err != nil {
		return order, err
	}
	RecordAuditEvent(gc, order.ID, order.OfferID, p2pModels.EventOrderCancelled, customerUserID, order)
	NotifyUsername(gc, order.MerchantUsername, "Trovo P2P: Order Cancelled", "The customer cancelled an order before escrow deposit.", map[string]string{"orderId": order.ID, "type": "P2P_ORDER_CANCELLED"})
	return order, nil
}

// ExpireStaleOrders transitions any AWAITING_APPROVAL order past its
// ExpiresAt into EXPIRED and releases its reserved liquidity. Intended to be
// called from a periodic background sweep (mirrors Plan Section 24's
// AWAITING_APPROVAL -> EXPIRED transition; app-backend has no scheduler
// primitive to hook into beyond a plain ticker goroutine, started from
// Init - see controllers/main.go).
// timeoutEligibleStatuses are the non-terminal statuses that carry a real
// deadline (Order.ExpiresAt is set/reset on every stage transition - Plan
// Section 24). Each has a distinct, safe consequence when it lapses - see
// ExpireStaleOrders.
var timeoutEligibleStatuses = []string{
	p2pModels.OrderStatusAwaitingApproval,
	p2pModels.OrderStatusAwaitingEscrowDeposit,
	p2pModels.OrderStatusAwaitingPayment,
	p2pModels.OrderStatusAwaitingPaymentConfirmation,
}

// ExpireStaleOrders sweeps every order past its current stage's deadline
// and applies that stage's specific, safe consequence:
//   - AWAITING_APPROVAL / AWAITING_ESCROW_DEPOSIT: no funds have moved yet,
//     so the order simply expires and its reserved liquidity is released.
//   - AWAITING_PAYMENT: the depositor's asset is already in escrow, so
//     expiring must not just abandon it - the order is cancelled and the
//     deposited amount is queued as a Refund back to the actual on-chain
//     depositor (Plan Section 45).
//   - AWAITING_PAYMENT_CONFIRMATION: the buyer has claimed they paid: an
//     automated cancel/refund here could wrongly undo a legitimate
//     payment, so this never moves funds unilaterally - it escalates into
//     a dispute for a human (self-resolution or admin) to resolve instead.
func ExpireStaleOrders(gc *sharedconfig.GlobalConfig) (int, error) {
	var staleOrders []p2pModels.Order
	if err := gc.DB.Where("order_status IN ? AND expires_at IS NOT NULL AND expires_at < ? AND is_disputed = ?",
		timeoutEligibleStatuses, time.Now().UTC(), false).Find(&staleOrders).Error; err != nil {
		return 0, err
	}
	count := 0
	for i := range staleOrders {
		order := staleOrders[i]
		switch order.OrderStatus {
		case p2pModels.OrderStatusAwaitingApproval, p2pModels.OrderStatusAwaitingEscrowDeposit:
			order.OrderStatus = p2pModels.OrderStatusExpired
			if err := releaseReservedLiquidityAndSave(gc, &order); err != nil {
				continue
			}
			RecordAuditEvent(gc, order.ID, order.OfferID, p2pModels.EventOrderExpired, "", order)
			NotifyUsername(gc, order.CustomerUsername, "Trovo P2P: Order Expired", "Your order expired before the merchant responded in time.", map[string]string{"orderId": order.ID, "type": "P2P_ORDER_EXPIRED"})
			NotifyUsername(gc, order.MerchantUsername, "Trovo P2P: Order Expired", "An order expired before it reached escrow deposit.", map[string]string{"orderId": order.ID, "type": "P2P_ORDER_EXPIRED"})

		case p2pModels.OrderStatusAwaitingPayment:
			if err := expireAwaitingPaymentOrder(gc, &order); err != nil {
				continue
			}

		case p2pModels.OrderStatusAwaitingPaymentConfirmation:
			if err := escalateStalePaymentConfirmationToDispute(gc, &order); err != nil {
				continue
			}
		}
		count++
	}
	return count, nil
}

// expireAwaitingPaymentOrder cancels an order whose fiat-payment window
// lapsed and refunds the already-escrowed deposit back to its real
// depositor (found via the canonical BlockchainDeposit, not
// Order.AssetDepositor - Plan Section 11's role is a user id, not a wallet
// address).
func expireAwaitingPaymentOrder(gc *sharedconfig.GlobalConfig, order *p2pModels.Order) error {
	deposit, err := FindCanonicalDepositForOrder(gc, order.ID)
	if err != nil {
		// No canonical deposit on record despite reaching AWAITING_PAYMENT
		// should not happen, but expiring without a refund target would
		// silently strand funds - skip and let the next sweep retry rather
		// than guess.
		return err
	}
	order.OrderStatus = p2pModels.OrderStatusCancelled
	order.RefundableAmount = order.DepositedEscrowAmount
	if err := gc.DB.Save(order).Error; err != nil {
		return &tErrors.ErrorTemporaryServerError{}
	}
	refund := p2pModels.Refund{
		ID:              gc.GenerateUUIDString(),
		DepositID:       deposit.ID,
		OrderID:         order.ID,
		Sender:          deposit.Sender,
		Token:           deposit.Token,
		ContractAddress: deposit.ContractAddress,
		Amount:          order.DepositedEscrowAmount,
		Reason:          p2pModels.RefundReasonOrderExpired,
	}
	if err := gc.DB.Omit(clause.Associations).Create(&refund).Error; err != nil {
		return &tErrors.ErrorTemporaryServerError{}
	}
	RecordAuditEvent(gc, order.ID, order.OfferID, p2pModels.EventOrderExpired, "", order)
	RecordAuditEvent(gc, order.ID, order.OfferID, p2pModels.EventRefundIssued, deposit.Sender, refund)
	NotifyUsername(gc, order.CustomerUsername, "Trovo P2P: Order Expired", "Your order was cancelled because payment was not completed in time. Any escrowed deposit is refundable.", map[string]string{"orderId": order.ID, "type": "P2P_ORDER_EXPIRED"})
	NotifyUsername(gc, order.MerchantUsername, "Trovo P2P: Order Expired", "An order was cancelled because payment was not completed in time.", map[string]string{"orderId": order.ID, "type": "P2P_ORDER_EXPIRED"})
	return nil
}

// escalateStalePaymentConfirmationToDispute opens a system-initiated
// dispute rather than moving any funds - see ExpireStaleOrders' doc comment
// for why this stage never auto-cancels/refunds.
func escalateStalePaymentConfirmationToDispute(gc *sharedconfig.GlobalConfig, order *p2pModels.Order) error {
	dispute := p2pModels.Dispute{
		ID:          gc.GenerateUUIDString(),
		OrderID:     order.ID,
		OpenedBy:    "SYSTEM",
		OpenedAt:    time.Now().UTC(),
		Subject:     p2pModels.DisputeSubjectOther,
		Description: "Automatically opened: the payment confirmation window expired before the seller confirmed receipt.",
		Status:      p2pModels.DisputeStatusOpen,
	}
	dbTX := gc.DB.Begin()
	if err := dbTX.Omit(clause.Associations).Create(&dispute).Error; err != nil {
		dbTX.Rollback()
		return &tErrors.ErrorTemporaryServerError{}
	}
	if err := dbTX.Model(&p2pModels.Order{}).Where("id = ?", order.ID).Update("is_disputed", true).Error; err != nil {
		dbTX.Rollback()
		return &tErrors.ErrorTemporaryServerError{}
	}
	if err := dbTX.Commit().Error; err != nil {
		return &tErrors.ErrorTemporaryServerError{}
	}
	RecordAuditEvent(gc, order.ID, order.OfferID, p2pModels.EventDisputeOpened, "SYSTEM", dispute)
	NotifyUsername(gc, order.CustomerUsername, "Trovo P2P: Dispute Opened", "The payment confirmation window expired, so this order was flagged for review.", map[string]string{"orderId": order.ID, "type": "P2P_DISPUTE_OPENED"})
	NotifyUsername(gc, order.MerchantUsername, "Trovo P2P: Dispute Opened", "The payment confirmation window expired, so this order was flagged for review.", map[string]string{"orderId": order.ID, "type": "P2P_DISPUTE_OPENED"})
	return nil
}

func releaseReservedLiquidityAndSave(gc *sharedconfig.GlobalConfig, order *p2pModels.Order) error {
	dbTX := gc.DB.Begin()
	if err := dbTX.Save(order).Error; err != nil {
		dbTX.Rollback()
		return &tErrors.ErrorTemporaryServerError{}
	}
	var offer p2pModels.Offer
	if err := dbTX.Where("id = ?", order.OfferID).First(&offer).Error; err == nil {
		specifiedAmount := decimal.RequireFromString(order.SpecifiedAssetAmount)
		available := decimal.RequireFromString(offer.AvailableLiquidity).Add(specifiedAmount)
		reserved := decimal.RequireFromString(offer.ReservedLiquidity).Sub(specifiedAmount)
		if reserved.IsNegative() {
			reserved = decimal.Zero
		}
		if err := dbTX.Model(&p2pModels.Offer{}).Where("id = ?", offer.ID).
			Updates(map[string]interface{}{
				"available_liquidity": available.String(),
				"reserved_liquidity":   reserved.String(),
			}).Error; err != nil {
			dbTX.Rollback()
			return &tErrors.ErrorTemporaryServerError{}
		}
	}
	if err := dbTX.Commit().Error; err != nil {
		return &tErrors.ErrorTemporaryServerError{}
	}
	return nil
}

// ListOrdersFilter narrows My Orders/My Trades (Plan Section 87).
type ListOrdersFilter struct {
	UserID   string
	Role     string // "customer", "merchant", or "" for either
	Status   string
	Page     int
	PageSize int
}

// ListOrders returns paginated orders for a user, filtered by role/status.
func ListOrders(db *gorm.DB, f ListOrdersFilter) ([]p2pModels.Order, int64, error) {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize < 1 || f.PageSize > 100 {
		f.PageSize = 20
	}
	q := db.Model(&p2pModels.Order{})
	switch f.Role {
	case "customer":
		q = q.Where("customer_user_id = ?", f.UserID)
	case "merchant":
		q = q.Where("merchant_user_id = ?", f.UserID)
	default:
		q = q.Where("customer_user_id = ? OR merchant_user_id = ?", f.UserID, f.UserID)
	}
	if f.Status != "" {
		q = q.Where("order_status = ?", f.Status)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var orders []p2pModels.Order
	err := q.Order("created_at desc").Offset((f.Page - 1) * f.PageSize).Limit(f.PageSize).Find(&orders).Error
	return orders, total, err
}

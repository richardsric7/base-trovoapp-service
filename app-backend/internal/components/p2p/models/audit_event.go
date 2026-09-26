package p2p

import "time"

// P2P audit event names (Plan Section 68). SHORTLINK_CREATED is the sole
// event for shortlink/payment-link generation - PAYMENT_LINK_CREATED was
// dropped as redundant.
const (
	EventOfferCreated          = "OFFER_CREATED"
	EventOfferUpdated          = "OFFER_UPDATED"
	EventOrderCreated          = "ORDER_CREATED"
	EventOrderAccepted         = "ORDER_ACCEPTED"
	EventOrderRejected         = "ORDER_REJECTED"
	EventOrderCancelled        = "ORDER_CANCELLED"
	EventOrderExpired          = "ORDER_EXPIRED"
	EventShortlinkCreated      = "SHORTLINK_CREATED"
	EventEscrowDepositDetected = "ESCROW_DEPOSIT_DETECTED"
	EventEscrowConfirmed       = "ESCROW_CONFIRMED"
	EventFiatPaymentMarkedSent = "FIAT_PAYMENT_MARKED_SENT"
	EventFiatPaymentConfirmed  = "FIAT_PAYMENT_CONFIRMED"
	EventAssetReleased         = "ASSET_RELEASED"
	EventOrderCompleted        = "ORDER_COMPLETED"
	EventDisputeOpened         = "DISPUTE_OPENED"
	EventDisputeResolved       = "DISPUTE_RESOLVED"
	EventRefundIssued          = "REFUND_ISSUED"
	EventRefundClaimed         = "REFUND_CLAIMED"
)

// P2PAuditEvent is an append-only audit trail row. app-backend has no
// general-purpose audit-log system today (Plan Section 68's finding) - this
// is P2P's own, scoped to this module.
type P2PAuditEvent struct {
	ID        string    `json:"id" gorm:"primaryKey;size:36"`
	OrderID   string    `json:"orderId" gorm:"size:36;index:idx_p2p_audit_order_id"`
	OfferID   string    `json:"offerId" gorm:"size:36;index:idx_p2p_audit_offer_id"`
	EventName string    `json:"eventName" gorm:"size:40;not null;index:idx_p2p_audit_event_name"`
	ActorID   string    `json:"actorId" gorm:"size:100;default:''"`
	Detail    string    `json:"detail" gorm:"type:text"` // JSON blob
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TableName overrides GORM's default naming, which mangles "P2P" (a
// capital letter following a digit) into "p2_p_audit_events".
func (P2PAuditEvent) TableName() string {
	return "p2p_audit_events"
}

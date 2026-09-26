package p2p

import "time"

// Dispute statuses
const (
	DisputeStatusOpen     = "OPEN"
	DisputeStatusResolved = "RESOLVED"
)

// Dispute subjects (Plan Section 62)
const (
	DisputeSubjectNoPayment                = "NO_PAYMENT"
	DisputeSubjectUnderPayment             = "UNDER_PAYMENT"
	DisputeSubjectOverPayment              = "OVER_PAYMENT"
	DisputeSubjectPaymentNotReceived       = "PAYMENT_NOT_RECEIVED"
	DisputeSubjectPaymentMarkedSentInError = "PAYMENT_MARKED_SENT_IN_ERROR"
	DisputeSubjectWrongPaymentAmount       = "WRONG_PAYMENT_AMOUNT"
	DisputeSubjectOther                    = "OTHER"
)

// Dispute resolutions
const (
	DisputeResolutionInFavorOfBuyer  = "IN_FAVOR_OF_BUYER"
	DisputeResolutionInFavorOfSeller = "IN_FAVOR_OF_SELLER"
	DisputeResolutionSplit           = "SPLIT"
)

// Dispute.ID is the sole canonical dispute identifier - no separate disputeId
// field, matching Offer.ID/Order.ID's already-established convention.
type Dispute struct {
	ID              string     `json:"id" gorm:"primaryKey;size:36"`
	OrderID         string     `json:"orderId" gorm:"size:36;not null;index:idx_p2p_dispute_order_id"`
	OpenedBy        string     `json:"openedBy" gorm:"size:100;not null"`
	OpenedAt        time.Time  `json:"openedAt"`
	Subject         string     `json:"subject" gorm:"size:40;not null"`
	Description     string     `json:"description" gorm:"type:text"`
	Evidence        string     `json:"evidence" gorm:"type:text"` // JSON array of evidence URLs/attachment references
	Status          string     `json:"status" gorm:"size:20;not null;default:'OPEN';index:idx_p2p_dispute_status"`
	Resolution      string     `json:"resolution" gorm:"size:30;default:''"`
	ResolvedAt      *time.Time `json:"resolvedAt"`
	ResolvedBy      string     `json:"resolvedBy" gorm:"size:100;default:''"`
	ResolutionTimeSeconds int64 `json:"resolutionTime" gorm:"default:0"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}

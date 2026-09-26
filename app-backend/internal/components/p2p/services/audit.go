package p2p

import (
	"encoding/json"
	"log"
	p2pModels "trovo-wallet-api/internal/components/p2p/models"
	"trovo-wallet-api/internal/sharedconfig"

	"gorm.io/gorm/clause"
)

// RecordAuditEvent writes an append-only P2PAuditEvent row. detail may be any
// JSON-marshalable value or nil.
func RecordAuditEvent(gc *sharedconfig.GlobalConfig, orderID, offerID, eventName, actorID string, detail interface{}) {
	var detailStr string
	if detail != nil {
		b, err := json.Marshal(detail)
		if err == nil {
			detailStr = string(b)
		}
	}
	event := p2pModels.P2PAuditEvent{
		ID:        gc.GenerateUUIDString(),
		OrderID:   orderID,
		OfferID:   offerID,
		EventName: eventName,
		ActorID:   actorID,
		Detail:    detailStr,
	}
	if err := gc.DB.Omit(clause.Associations).Create(&event).Error; err != nil {
		log.Printf("[p2p:RecordAuditEvent] failed to write audit event %v for order %v: %v\n", eventName, orderID, err)
	}
}

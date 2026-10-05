package publicmarkets

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	pm "trovo-wallet-api/internal/components/publicmarkets/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Inbound event kinds.
const (
	EventSettlement      = "settlement"       // Custodian §4.4.1
	EventPositionFeed    = "position-feed"    // Custodian §4.4.2
	EventCorporateAction = "corporate-action" // Custodian §4.4.3
	EventExecution       = "execution"        // Dealing Member §5.3
)

func (e *Engine) logf(format string, args ...interface{}) {
	log.Printf("[publicmarkets] "+format, args...)
}

// ReceivePartnerEvent stores an inbound partner message once per partner
// reference and processes it. A redelivered message (same reference) is
// not processed again; duplicate reports that.
func (e *Engine) ReceivePartnerEvent(source, kind, reference string, payload []byte) (duplicate bool, err error) {
	reference = strings.TrimSpace(reference)
	if reference == "" {
		return false, refuse(http.StatusBadRequest, "error-missing-reference", "", "The message has no reference.")
	}
	ev := pm.PartnerEvent{Source: source, Kind: kind, Reference: trimTo(reference, 160), Payload: string(payload), CreatedAt: e.now()}
	res := e.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&ev)
	if res.Error != nil {
		return false, res.Error
	}
	if res.RowsAffected == 0 {
		return true, nil
	}
	return false, e.ProcessPartnerEvent(&ev)
}

// ProcessPendingEvents processes stored partner events not yet processed
// (Trovo Manager records manual confirmations this way; failures retry).
func (e *Engine) ProcessPendingEvents() {
	var evs []pm.PartnerEvent
	e.DB.Where("processed_at IS NULL").Order("id").Limit(100).Find(&evs)
	for i := range evs {
		if err := e.ProcessPartnerEvent(&evs[i]); err != nil {
			e.logf("event %d (%s %s): %v", evs[i].ID, evs[i].Source, evs[i].Kind, err)
		}
	}
}

// ProcessPartnerEvent applies one partner event. It is safe to run again:
// each handler moves state only from where the event expects it.
func (e *Engine) ProcessPartnerEvent(ev *pm.PartnerEvent) error {
	var result string
	var err error
	switch ev.Kind {
	case EventSettlement:
		result, err = e.onSettlement(ev)
	case EventExecution:
		result, err = e.onExecution(ev)
	case EventPositionFeed:
		result, err = e.onPositionFeed(ev)
	case EventCorporateAction:
		result, err = e.onCorporateAction(ev)
	default:
		result = "unknown event kind"
	}
	if err != nil {
		var pe *Error
		if !errors.As(err, &pe) { // transient: leave unprocessed so it is retried
			return err
		}
		result = "refused: " + pe.Message
	}
	now := e.now()
	e.DB.Model(&pm.PartnerEvent{}).Where("id = ?", ev.ID).Updates(map[string]interface{}{"processed_at": &now, "result": trimTo(result, 500)})
	return err
}

type settlementBody struct {
	InstructionID      string `json:"instructionId"`
	Status             string `json:"status"`
	SettledQuantity    string `json:"settledQuantity"`
	SettlementDate     string `json:"settlementDate"`
	CustodianReference string `json:"custodianReference"`
}

func (e *Engine) onSettlement(ev *pm.PartnerEvent) (string, error) {
	var b settlementBody
	if err := json.Unmarshal([]byte(ev.Payload), &b); err != nil || b.InstructionID == "" {
		return "", refuse(http.StatusBadRequest, "error-invalid-payload", "instructionId", "instructionId is required.")
	}
	var in pm.Instruction
	if err := e.DB.First(&in, "id = ?", b.InstructionID).Error; err != nil {
		return "", refuse(http.StatusNotFound, "error-unknown-instruction", "instructionId", "Unknown instruction %s.", b.InstructionID)
	}
	if !strings.HasPrefix(ev.Source, "MANUAL") && ev.Source != "CUSTODIAN:"+e.custodianCode(in.PartnerID) {
		return "", refuse(http.StatusForbidden, "error-wrong-partner", "instructionId", "This instruction was not sent to you.")
	}
	switch strings.ToLower(b.Status) {
	case "settlement_final":
	case "failed", "rejected", "cancelled":
		e.escalate(&in, "Custodian reported "+b.Status)
		return "instruction " + b.Status, nil
	default: // received, submitted, chain_final: progress only
		e.DB.Model(&pm.Instruction{}).Where("id = ?", in.ID).Update("partner_ref", b.CustodianReference)
		return "status " + b.Status + " noted", nil
	}
	if in.Status == pm.InstrSettled {
		return "already settled", nil
	}
	now := e.now()
	settled := d(b.SettledQuantity)
	if !settled.IsPositive() {
		settled = d(in.Quantity)
	}
	err := e.DB.Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&pm.Instruction{}).Where("id = ? AND status <> ?", in.ID, pm.InstrSettled).
			Updates(map[string]interface{}{"status": pm.InstrSettled, "settled_at": &now, "partner_ref": b.CustodianReference})
		if res.Error != nil || res.RowsAffected == 0 {
			return res.Error
		}
		// the Custodian's position moves by what settled
		var a pm.Asset
		if err := tx.First(&a, "id = ?", in.AssetID).Error; err != nil {
			return err
		}
		prev := "0"
		var last pm.CustodianPosition
		if tx.Where("asset_id = ?", a.ID).Order("id DESC").First(&last).Error == nil {
			prev = last.RealUnitsHeld
		}
		next := d(prev).Add(settled)
		if in.Side == "SELL" {
			next = d(prev).Sub(settled)
		}
		return tx.Create(&pm.CustodianPosition{AssetID: a.ID, RealUnitsHeld: next.String(), AsOf: now, Source: pm.PositionSettlement,
			Reference: fmt.Sprintf("%s %s %s (%s)", in.ID, in.Side, settled, b.CustodianReference), CreatedAt: now}).Error
	})
	if err != nil {
		return "", err
	}
	if in.BatchID != "" {
		e.batchSettled(in.BatchID, settled.String(), now)
	}
	return "settled " + settled.String(), nil
}

func (e *Engine) custodianCode(id uint64) string {
	var c pm.Custodian
	e.DB.First(&c, "custodian_id = ?", id)
	return c.Code
}

func (e *Engine) dealingMemberCode(id uint64) string {
	var dm pm.DealingMember
	e.DB.First(&dm, id)
	return dm.Code
}

type executionBody struct {
	OrderID          string `json:"orderId"`
	Status           string `json:"status"`
	ExecutedQuantity string `json:"executedQuantity"`
	ExecutedPrice    string `json:"executedPrice"`
	ExecutionTime    string `json:"executionTime"`
}

func (e *Engine) onExecution(ev *pm.PartnerEvent) (string, error) {
	var b executionBody
	if err := json.Unmarshal([]byte(ev.Payload), &b); err != nil || b.OrderID == "" {
		return "", refuse(http.StatusBadRequest, "error-invalid-payload", "orderId", "orderId is required.")
	}
	var in pm.Instruction
	if err := e.DB.First(&in, "id = ? AND kind = ?", b.OrderID, pm.InstrDealingOrder).Error; err != nil {
		return "", refuse(http.StatusNotFound, "error-unknown-order", "orderId", "Unknown order %s.", b.OrderID)
	}
	if !strings.HasPrefix(ev.Source, "MANUAL") && ev.Source != "DEALING_MEMBER:"+e.dealingMemberCode(in.PartnerID) {
		return "", refuse(http.StatusForbidden, "error-wrong-partner", "orderId", "This order was not sent to you.")
	}
	status := strings.ToUpper(b.Status)
	switch status {
	case "FILLED":
	case "REJECTED", "CANCELLED", "EXPIRED":
		e.escalate(&in, "Dealing Member reported "+status)
		return "order " + status, nil
	case "PARTIALLY_FILLED":
		e.escalate(&in, fmt.Sprintf("Partially filled: %s of %s at %s - allocate manually", b.ExecutedQuantity, in.Quantity, b.ExecutedPrice))
		return "partial fill escalated", nil
	default:
		return "status " + status + " noted", nil
	}
	price := d(b.ExecutedPrice)
	if !price.IsPositive() {
		return "", refuse(http.StatusBadRequest, "error-invalid-payload", "executedPrice", "executedPrice is required for a fill.")
	}
	now := e.now()
	res := e.DB.Model(&pm.Instruction{}).Where("id = ? AND status IN ?", in.ID, []string{pm.InstrPending, pm.InstrAccepted, pm.InstrEscalated}).
		Updates(map[string]interface{}{"status": pm.InstrExecuted, "executed_at": &now})
	if res.RowsAffected == 0 {
		return "already executed", nil
	}
	qty := b.ExecutedQuantity
	if !d(qty).IsPositive() {
		qty = in.Quantity
	}
	if in.BatchID != "" {
		e.batchExecuted(in.BatchID, qty, price.String(), now)
	}
	return fmt.Sprintf("filled %s at %s", qty, price), nil
}

type positionFeedBody struct {
	AsOf      string `json:"asOf"`
	Positions []struct {
		AssetCode string `json:"assetCode"`
		UnitsHeld string `json:"unitsHeld"`
	} `json:"positions"`
}

func (e *Engine) onPositionFeed(ev *pm.PartnerEvent) (string, error) {
	var b positionFeedBody
	if err := json.Unmarshal([]byte(ev.Payload), &b); err != nil {
		return "", refuse(http.StatusBadRequest, "error-invalid-payload", "", "Unreadable position feed.")
	}
	code := strings.TrimPrefix(ev.Source, "CUSTODIAN:")
	var c pm.Custodian
	if !strings.HasPrefix(ev.Source, "MANUAL") {
		if err := e.DB.First(&c, "code = ?", code).Error; err != nil {
			return "", refuse(http.StatusForbidden, "error-unknown-partner", "", "Unknown Custodian.")
		}
	}
	asOf, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(b.AsOf), Lagos)
	if err != nil {
		if t, e2 := time.Parse(time.RFC3339, b.AsOf); e2 == nil {
			asOf = t
		} else {
			asOf = e.now()
		}
	}
	source := pm.PositionFeed
	if strings.HasPrefix(ev.Source, "MANUAL") { // recorded by Operations in Trovo Manager
		source = pm.PositionManual
	}
	n := 0
	for _, p := range b.Positions {
		var a pm.Asset
		if e.DB.First(&a, "asset_code = ?", strings.ToUpper(p.AssetCode)).Error != nil {
			continue
		}
		if c.CustodianID != 0 && a.CustodianID != c.CustodianID {
			continue // a Custodian reports only its own assets
		}
		units := d(p.UnitsHeld)
		if units.IsNegative() {
			continue
		}
		e.DB.Create(&pm.CustodianPosition{AssetID: a.ID, RealUnitsHeld: units.String(), AsOf: asOf, Source: source,
			Reference: ev.Source + " feed " + b.AsOf, CreatedAt: e.now()})
		n++
	}
	return fmt.Sprintf("%d position(s) recorded", n), nil
}

// ---------------------------------------------------------------- outbound instructions

// DispatchInstructions sends due instructions. A refusal or the retry cap
// escalates to Operations (Trovo Manager › Orders › Escalations); a lost
// instruction is never silently dropped.
func (e *Engine) DispatchInstructions() {
	var due []pm.Instruction
	e.DB.Where("status = ? AND (next_attempt_at IS NULL OR next_attempt_at <= ?)", pm.InstrPending, e.now()).Order("created_at").Limit(50).Find(&due)
	for i := range due {
		e.dispatch(&due[i])
	}
	e.checkSettlementSLA()
}

func (e *Engine) dispatch(in *pm.Instruction) {
	ctx, cancel := contextTimeout(45 * time.Second)
	defer cancel()
	var err error
	var ack string
	manual := false
	switch in.PartnerType {
	case "CUSTODIAN":
		var c pm.Custodian
		if err = e.DB.First(&c, "custodian_id = ?", in.PartnerID).Error; err != nil {
			break
		}
		client, cerr := e.Partners.Custodian(c)
		if cerr != nil {
			err = cerr
			break
		}
		if client == nil {
			manual = true
			break
		}
		r, sendErr := client.SendInstruction(ctx, in.Side, instructionPayload(in))
		err, ack = sendErr, r.Status
	case "DEALING_MEMBER":
		var dm pm.DealingMember
		if err = e.DB.First(&dm, in.PartnerID).Error; err != nil {
			break
		}
		client, cerr := e.Partners.DealingMember(dm)
		if cerr != nil {
			err = cerr
			break
		}
		if client == nil {
			manual = true
			break
		}
		r, sendErr := client.PlaceOrder(ctx, orderPayload(in))
		err, ack = sendErr, r.Status
	}
	now := e.now()
	if manual {
		e.DB.Model(&pm.Instruction{}).Where("id = ? AND status = ?", in.ID, pm.InstrPending).Updates(map[string]interface{}{
			"status": pm.InstrAccepted, "accepted_at": &now, "last_error": "MANUAL: awaiting Operations on the partner's portal"})
		return
	}
	if err == nil {
		e.DB.Model(&pm.Instruction{}).Where("id = ? AND status = ?", in.ID, pm.InstrPending).Updates(map[string]interface{}{
			"status": pm.InstrAccepted, "accepted_at": &now, "attempts": in.Attempts + 1, "last_error": "", "partner_ref": ack})
		return
	}
	attempts := in.Attempts + 1
	if isPermanent(err) || attempts >= in.MaxAttempts {
		e.DB.Model(&pm.Instruction{}).Where("id = ?", in.ID).Update("attempts", attempts)
		in.Attempts = attempts
		e.escalate(in, err.Error())
		return
	}
	next := now.Add(time.Duration(30<<uint(attempts)) * time.Second) // 1m, 2m, 4m, 8m ...
	e.DB.Model(&pm.Instruction{}).Where("id = ?", in.ID).Updates(map[string]interface{}{"attempts": attempts, "next_attempt_at": &next, "last_error": trimTo(err.Error(), 500)})
}

func (e *Engine) escalate(in *pm.Instruction, reason string) {
	e.DB.Model(&pm.Instruction{}).Where("id = ?", in.ID).Updates(map[string]interface{}{"status": pm.InstrEscalated, "last_error": trimTo(reason, 500)})
	if in.BatchID != "" {
		e.DB.Model(&pm.NetBatch{}).Where("id = ?", in.BatchID).Update("note", trimTo("Instruction "+in.ID+" escalated: "+reason, 500))
	}
	if e.GC != nil {
		e.GC.LogDiscordFailedRequest(fmt.Sprintf("[publicmarkets] instruction %s to %s escalated after %d attempt(s): %s", in.ID, in.PartnerName, in.Attempts, reason))
	}
}

func (e *Engine) checkSettlementSLA() {
	s := LoadSettings(e.DB)
	cutoff := e.now().Add(-time.Duration(s.SettlementSLAHours) * time.Hour)
	var late []pm.Instruction
	e.DB.Where("partner_type = ? AND status IN ? AND accepted_at < ? AND last_error NOT LIKE ?", "CUSTODIAN",
		[]string{pm.InstrAccepted, pm.InstrExecuted}, cutoff, "SLA%").Find(&late)
	for _, in := range late {
		e.DB.Model(&pm.Instruction{}).Where("id = ?", in.ID).Update("last_error", "SLA: settlement not confirmed within the agreed window")
		if e.GC != nil {
			e.GC.LogDiscordFailedRequest(fmt.Sprintf("[publicmarkets] instruction %s (%s %s %s) has no settlement after %dh", in.ID, in.Side, in.Quantity, in.AssetCode, s.SettlementSLAHours))
		}
	}
}

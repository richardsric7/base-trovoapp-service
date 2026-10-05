package publicmarkets

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	pm "trovo-wallet-api/internal/components/publicmarkets/models"
	"trovo-wallet-api/internal/components/publicmarkets/partners"
	userModels "trovo-wallet-api/internal/components/users/models"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func contextTimeout(t time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), t)
}

func isPermanent(err error) bool { return partners.IsPermanent(err) }

func instructionPayload(in *pm.Instruction) partners.CustodianInstruction {
	var ci partners.CustodianInstruction
	_ = json.Unmarshal([]byte(in.Payload), &ci)
	return ci
}

func orderPayload(in *pm.Instruction) partners.DealingOrder {
	var o partners.DealingOrder
	_ = json.Unmarshal([]byte(in.Payload), &o)
	return o
}

// activeBatch is an asset's batch still in flight (one at a time).
func (e *Engine) activeBatch(assetID string) (*pm.NetBatch, bool) {
	var b pm.NetBatch
	err := e.DB.Where("asset_id = ? AND status IN ?", assetID, []string{pm.BatchAwaitingApproval, pm.BatchReleased, pm.BatchExecuted, pm.BatchSettled}).First(&b).Error
	return &b, err == nil
}

// BuildBatches gathers each asset's slow-path orders into the session's
// net instruction, at most every BatchIntervalMinutes while its market is
// open. Creations and redemptions in a batch net against each other and
// against inventory; only the unmatched remainder goes to the market.
func (e *Engine) BuildBatches(ctx context.Context) {
	s := LoadSettings(e.DB)
	var assets []pm.Asset
	e.DB.Where("status = ?", pm.AssetLive).Find(&assets)
	for i := range assets {
		a := &assets[i]
		if _, busy := e.activeBatch(a.ID); busy {
			continue
		}
		p := e.CurrentPrice(a)
		if !p.Fresh {
			continue // waits for the session (and a live price)
		}
		var last pm.NetBatch
		if e.DB.Where("asset_id = ?", a.ID).Order("created_at DESC").First(&last).Error == nil &&
			last.SessionDate == SessionDate(e.now()) && e.now().Sub(last.CreatedAt) < time.Duration(s.BatchIntervalMinutes)*time.Minute {
			continue
		}
		var orders []pm.Order
		e.DB.Where("asset_id = ? AND state = ? AND batch_id = ''", a.ID, pm.StatePendingExecution).Order("created_at").Find(&orders)
		if len(orders) == 0 {
			continue
		}
		if err := e.buildBatch(ctx, a, p.Value, orders, s); err != nil {
			e.logf("batch for %s: %v", a.AssetCode, err)
		}
	}
}

func (e *Engine) buildBatch(ctx context.Context, a *pm.Asset, price decimal.Decimal, orders []pm.Order, s pm.Settings) error {
	creationUnits, redemptionUnits, redemptionCash := decimal.Zero, decimal.Zero, decimal.Zero
	var nc, nr int
	for _, o := range orders {
		if o.Type == pm.OrderCreation {
			creationUnits = creationUnits.Add(d(o.NetAmount).Div(price))
			nc++
		} else {
			redemptionUnits = redemptionUnits.Add(d(o.Quantity))
			gross := d(o.Quantity).Mul(price)
			redemptionCash = redemptionCash.Add(gross.Sub(gross.Mul(d(o.FeePercent)).Div(decimal.NewFromInt(100))))
			nr++
		}
	}
	inventory, err := e.Inventory(ctx, a)
	if err != nil {
		return err
	}
	if inventory.IsNegative() {
		inventory = decimal.Zero
	}
	side, qty := "NONE", decimal.Zero
	// units: redemptions and inventory cover creations first
	if need := creationUnits.Sub(inventory).Sub(redemptionUnits); need.IsPositive() {
		side = "BUY"
		// buying exactly the shortfall leaves no inventory: top it up to
		// the asset's target so later creations fill instantly
		if target := d(a.InventoryTargetUnits); target.IsPositive() {
			need = need.Add(target)
		}
		qty = need.Ceil() // whole shares / bond units
	} else {
		// cash: redemptions are paid from the treasury when it can; the
		// rest of their units are sold
		free, ferr := e.TreasuryFreeCash(ctx)
		if ferr != nil {
			free = decimal.Zero
		}
		appCash := decimal.Zero
		for _, o := range orders {
			if o.Type == pm.OrderRedemption && o.Channel == pm.ChannelApp {
				gross := d(o.Quantity).Mul(price)
				appCash = appCash.Add(gross.Sub(gross.Mul(d(o.FeePercent)).Div(decimal.NewFromInt(100))))
			}
		}
		if appCash.GreaterThan(free) {
			leftover := inventory.Add(redemptionUnits).Sub(creationUnits)
			qty = decimal.Min(leftover, redemptionUnits).Floor()
			if qty.IsPositive() {
				side = "SELL"
			}
		}
	}
	id := fmt.Sprintf("NB-%s-%s", e.now().In(Lagos).Format("20060102-1504"), a.AssetCode)
	b := pm.NetBatch{ID: id, AssetID: a.ID, AssetCode: a.AssetCode, SessionDate: SessionDate(e.now()), Side: side, Quantity: qty.String(),
		ReferencePrice: price.String(), Value: qty.Mul(price).Round(2).String(), CreationOrders: nc, RedemptionOrders: nr,
		ApprovalsRequired: s.ApprovalsRequired, CreatedAt: e.now(), UpdatedAt: e.now()}
	switch side {
	case "NONE":
		b.Status = pm.BatchInternal
		b.Note = "Fully netted against inventory and the treasury: nothing sent to the market"
	case "BUY":
		b.Status = pm.BatchReleased
		var today []pm.NetBatch
		e.DB.Where("asset_id = ? AND session_date = ? AND side = ? AND status NOT IN ?", a.ID, b.SessionDate, "BUY",
			[]string{pm.BatchRejected, pm.BatchFailed}).Find(&today)
		total := d(b.Value)
		for _, t := range today {
			total = total.Add(d(t.Value))
		}
		if total.GreaterThan(d(s.NetCreationThreshold)) {
			b.Status = pm.BatchAwaitingApproval
			b.Note = fmt.Sprintf("Net daily creation %s NGN is above the approval threshold %s NGN", total.StringFixed(2), d(s.NetCreationThreshold).StringFixed(2))
		}
	default:
		b.Status = pm.BatchReleased
	}
	err = e.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&b).Error; err != nil {
			return err
		}
		ids := make([]string, 0, len(orders))
		for _, o := range orders {
			ids = append(ids, o.ID)
		}
		return tx.Model(&pm.Order{}).Where("id IN ? AND batch_id = ''", ids).Update("batch_id", b.ID).Error
	})
	if err != nil {
		return err
	}
	for i := range orders {
		e.note(&orders[i], "Included in net batch "+b.ID+" ("+strings.ToLower(b.Side)+" "+b.Quantity+")")
	}
	switch b.Status {
	case pm.BatchInternal:
		now := e.now()
		e.DB.Model(&pm.NetBatch{}).Where("id = ?", b.ID).Updates(map[string]interface{}{"executed_price": price.String(), "settled_at": &now})
		e.settleBatchOrders(&b, price, false, now)
	case pm.BatchReleased:
		e.releaseBatch(&b)
	}
	return nil
}

// ReleaseApprovedBatches releases batches whose approvals are complete,
// and refunds the orders of rejected ones.
func (e *Engine) ReleaseApprovedBatches() {
	s := LoadSettings(e.DB)
	allowed := CSV(s.NetCreationApprovers)
	var waiting []pm.NetBatch
	e.DB.Where("status = ?", pm.BatchAwaitingApproval).Find(&waiting)
	for i := range waiting {
		b := &waiting[i]
		var approvals []pm.BatchApproval
		e.DB.Where("batch_id = ?", b.ID).Find(&approvals)
		n := 0
		for _, ap := range approvals {
			if len(allowed) == 0 || contains(allowed, strings.ToLower(ap.Approver)) {
				n++
			}
		}
		if n >= b.ApprovalsRequired {
			res := e.DB.Model(&pm.NetBatch{}).Where("id = ? AND status = ?", b.ID, pm.BatchAwaitingApproval).Updates(map[string]interface{}{"status": pm.BatchReleased, "note": fmt.Sprintf("Approved by %d approver(s)", n)})
			if res.RowsAffected == 1 {
				e.releaseBatch(b)
			}
		}
	}
	// rejected or rolled batches: their orders go back or are refunded
	var closed []pm.NetBatch
	e.DB.Where("status IN ?", []string{pm.BatchRejected, pm.BatchFailed}).Where("updated_at > ?", e.now().Add(-72*time.Hour)).Find(&closed)
	for _, b := range closed {
		var orders []pm.Order
		e.DB.Where("batch_id = ? AND state = ?", b.ID, pm.StatePendingExecution).Find(&orders)
		for i := range orders {
			o := &orders[i]
			if b.Status == pm.BatchFailed { // rolled to the next session
				e.DB.Model(&pm.Order{}).Where("id = ?", o.ID).Update("batch_id", "")
				e.note(o, "Batch "+b.ID+" was rolled; the order waits for the next session")
				continue
			}
			var a pm.Asset
			e.DB.First(&a, "id = ?", o.AssetID)
			if err := e.reject(context.Background(), o, &a, "The net creation batch was not approved; the order was refunded."); err != nil {
				e.fail(o, err)
			}
		}
	}
}

func contains(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}

// releaseBatch queues the batch's instructions: the trade to the Dealing
// Member and the matching instruction to the Custodian.
func (e *Engine) releaseBatch(b *pm.NetBatch) {
	var a pm.Asset
	if e.DB.First(&a, "id = ?", b.AssetID).Error != nil {
		return
	}
	s := LoadSettings(e.DB)
	var dm pm.DealingMember
	e.DB.First(&dm, a.DealingMemberID)
	var cust userModels.ApprovedAssetCustodian
	e.DB.First(&cust, a.CustodianID)
	ref := "CREATE-" + b.ID
	kind := pm.InstrCustodianCreation
	if b.Side == "SELL" {
		ref, kind = "REDEEM-"+b.ID, pm.InstrCustodianRedemption
	}
	dmID := "DM-" + uuid.NewString()
	dmBody, _ := json.Marshal(partners.DealingOrder{OrderID: dmID, AssetCode: a.AssetCode, Market: a.Market, Side: b.Side, Quantity: b.Quantity, OrderType: "MARKET", TrovotechReference: ref})
	ciID := "CI-" + uuid.NewString()
	ciBody, _ := json.Marshal(partners.CustodianInstruction{InstructionID: ciID, AssetCode: a.AssetCode, ISIN: a.ISIN, Quantity: b.Quantity, Side: b.Side,
		ReferenceDate: b.SessionDate, TrovotechAccountReference: a.OmnibusReference})
	now := e.now()
	e.DB.Create(&[]pm.Instruction{
		{ID: dmID, Kind: pm.InstrDealingOrder, PartnerType: "DEALING_MEMBER", PartnerID: dm.ID, PartnerName: dm.DealingMemberName, AssetID: a.ID, AssetCode: a.AssetCode,
			BatchID: b.ID, Side: b.Side, Quantity: b.Quantity, Payload: string(dmBody), Status: pm.InstrPending, MaxAttempts: s.InstructionMaxAttempts, CreatedAt: now, UpdatedAt: now},
		{ID: ciID, Kind: kind, PartnerType: "CUSTODIAN", PartnerID: a.CustodianID, PartnerName: cust.AssetCustodianName, AssetID: a.ID, AssetCode: a.AssetCode,
			BatchID: b.ID, Side: b.Side, Quantity: b.Quantity, Payload: string(ciBody), Status: pm.InstrPending, MaxAttempts: s.InstructionMaxAttempts, CreatedAt: now, UpdatedAt: now},
	})
}

// batchExecuted records the Dealing Member's fill.
func (e *Engine) batchExecuted(batchID, qty, price string, at time.Time) {
	res := e.DB.Model(&pm.NetBatch{}).Where("id = ? AND status = ?", batchID, pm.BatchReleased).
		Updates(map[string]interface{}{"status": pm.BatchExecuted, "executed_quantity": qty, "executed_price": price, "executed_at": &at})
	if res.RowsAffected == 0 {
		return
	}
	var b pm.NetBatch
	e.DB.First(&b, "id = ?", batchID)
	var orders []pm.Order
	e.DB.Where("batch_id = ? AND state = ?", batchID, pm.StatePendingExecution).Find(&orders)
	for i := range orders {
		e.transition(&orders[i], []string{pm.StatePendingExecution}, pm.StateExecuted,
			fmt.Sprintf("The broker %s %s at %s; waiting for the Custodian's settlement", map[string]string{"BUY": "bought", "SELL": "sold"}[b.Side], qty, price),
			map[string]interface{}{"executed_price": price})
	}
	e.maybeSettle(&b)
}

// batchSettled records the Custodian's settlement_final.
func (e *Engine) batchSettled(batchID, qty string, at time.Time) {
	var b pm.NetBatch
	if e.DB.First(&b, "id = ?", batchID).Error != nil {
		return
	}
	e.DB.Model(&pm.NetBatch{}).Where("id = ?", batchID).Update("settled_at", &at)
	b.SettledAt = &at
	e.maybeSettle(&b)
}

// maybeSettle settles the batch's orders once both the fill and the
// Custodian's settlement are in (settlement_final gates the mint).
func (e *Engine) maybeSettle(b *pm.NetBatch) {
	if b.SettledAt == nil || b.ExecutedAt == nil {
		var fresh pm.NetBatch
		e.DB.First(&fresh, "id = ?", b.ID)
		*b = fresh
	}
	if b.SettledAt == nil || b.ExecutedAt == nil || b.Status != pm.BatchExecuted {
		return
	}
	res := e.DB.Model(&pm.NetBatch{}).Where("id = ? AND status = ?", b.ID, pm.BatchExecuted).Update("status", pm.BatchSettled)
	if res.RowsAffected == 0 {
		return
	}
	e.settleBatchOrders(b, d(b.ExecutedPrice), true, *b.SettledAt)
}

// settleBatchOrders fixes each order's price (the execution price for the
// side that traded, the reference price for netted orders) and moves it to
// settlement_final; the order loop then mints, burns and pays.
func (e *Engine) settleBatchOrders(b *pm.NetBatch, execPrice decimal.Decimal, traded bool, at time.Time) {
	var a pm.Asset
	e.DB.First(&a, "id = ?", b.AssetID)
	ref := d(b.ReferencePrice)
	var orders []pm.Order
	e.DB.Where("batch_id = ? AND state IN ?", b.ID, []string{pm.StatePendingExecution, pm.StateExecuted}).Find(&orders)
	for i := range orders {
		o := &orders[i]
		price := ref
		if traded && ((b.Side == "BUY" && o.Type == pm.OrderCreation) || (b.Side == "SELL" && o.Type == pm.OrderRedemption)) {
			price = execPrice
		}
		updates := map[string]interface{}{"settled_at": &at, "executed_price": price.String()}
		var note string
		if o.Type == pm.OrderCreation {
			qty := d(o.NetAmount).Div(price).Truncate(int32(a.TokenDecimals))
			updates["quantity"] = qty.String()
			note = fmt.Sprintf("%s confirmed settlement; %s %s at %s", e.custodianName(&a), qty, a.AssetCode, price)
		} else {
			gross := d(o.Quantity).Mul(price).Round(2)
			fee := gross.Mul(d(o.FeePercent)).Div(decimal.NewFromInt(100)).Round(2)
			updates["amount"], updates["fee"], updates["net_amount"] = gross.String(), fee.String(), gross.Sub(fee).String()
			note = fmt.Sprintf("Settled at %s; %s %s to pay", price, gross.Sub(fee).StringFixed(2), o.FundingAssetCode)
		}
		e.transition(o, []string{pm.StatePendingExecution, pm.StateExecuted}, pm.StateSettlementFinal, note, updates)
	}
	e.DB.Model(&pm.NetBatch{}).Where("id = ?", b.ID).Update("status", pm.BatchProcessed)
	if b.Status == pm.BatchInternal {
		e.DB.Model(&pm.NetBatch{}).Where("id = ?", b.ID).Update("status", pm.BatchInternal)
	}
}

package publicmarkets

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// ---------------------------------------------------------------- overview

// AttentionItem is something an operator should look at.
type AttentionItem struct {
	Kind   string `json:"kind"` // approval | escalation | dead-letter | confirmation | disclosure | halt | drift
	Text   string `json:"text"`
	Where  string `json:"where"`
	Target string `json:"target"` // an id the page links to
}

// Overview is the Public Markets landing page.
type Overview struct {
	Assets            int64           `json:"assets"`
	LiveAssets        int64           `json:"liveAssets"`
	HaltedAssets      int64           `json:"haltedAssets"`
	MarketValue       string          `json:"marketValue"`
	SettledToday      int64           `json:"settledToday"`
	SettledTodayValue string          `json:"settledTodayValue"`
	InProgress        int64           `json:"inProgress"`
	InProgressValue   string          `json:"inProgressValue"`
	CreationsToday    int64           `json:"creationsToday"`
	RedemptionsToday  int64           `json:"redemptionsToday"`
	FastPercent       string          `json:"fastPercent"`   // of today's creations
	NettedPercent     string          `json:"nettedPercent"` // of today's redemptions
	Halted            []Asset         `json:"halted"`
	Attention         []AttentionItem `json:"attention"`
	Jobs              []JobRun        `json:"jobs"`
	BatchesToday      []NetBatch      `json:"batchesToday"`
}

func startOfDay(t time.Time) time.Time {
	l := t.In(Lagos)
	return time.Date(l.Year(), l.Month(), l.Day(), 0, 0, 0, 0, Lagos).UTC()
}

func sumAmounts(orders []Order) string {
	total := decimal.Zero
	for _, o := range orders {
		total = total.Add(d(o.Amount))
	}
	return total.StringFixed(2)
}

func pct(part, whole int64) string {
	if whole == 0 {
		return "0"
	}
	return decimal.NewFromInt(part * 100).Div(decimal.NewFromInt(whole)).StringFixed(0)
}

// Overview gathers the landing page.
func (s *Service) Overview() Overview {
	var o Overview
	st := s.LoadSettings()
	rows := s.Assets(AssetFilters{})
	value := decimal.Zero
	for _, r := range rows {
		o.Assets++
		switch r.Status {
		case AssetLive:
			o.LiveAssets++
		case AssetHalted:
			o.HaltedAssets++
			o.Halted = append(o.Halted, r.Asset)
			o.Attention = append(o.Attention, AttentionItem{Kind: "halt", Text: r.AssetCode + " is halted: " + r.HaltReason, Where: "Assets", Target: r.ID})
		}
		value = value.Add(d(r.MarketValue))
		// substantial holdings (CAMA s.120)
		for _, h := range s.holders(&r.Asset, d(r.Supply), d(st.SubstantialHoldingPercent), 0) {
			if h.Substantial {
				o.Attention = append(o.Attention, AttentionItem{Kind: "disclosure", Text: fmt.Sprintf("A holder has %s%% of %s supply: CAC disclosure", h.PercentOfSupply, r.AssetCode),
					Where: "Assets › Beneficial owners", Target: r.ID})
			}
		}
	}
	o.MarketValue = value.StringFixed(2)

	today := startOfDay(s.now())
	var settled, open []Order
	s.DB.Where("state = ? AND completed_at >= ?", StateComplete, today).Find(&settled)
	s.DB.Where("state NOT IN ?", finalStates).Find(&open)
	o.SettledToday, o.SettledTodayValue = int64(len(settled)), sumAmounts(settled)
	o.InProgress, o.InProgressValue = int64(len(open)), sumAmounts(open)
	var fast, netted int64
	s.DB.Model(&Order{}).Where("type = ? AND created_at >= ?", OrderCreation, today).Count(&o.CreationsToday)
	s.DB.Model(&Order{}).Where("type = ? AND created_at >= ?", OrderRedemption, today).Count(&o.RedemptionsToday)
	s.DB.Model(&Order{}).Where("type = ? AND created_at >= ? AND path = ?", OrderCreation, today, PathFast).Count(&fast)
	s.DB.Model(&Order{}).Where("type = ? AND created_at >= ? AND path = ?", OrderRedemption, today, PathNetted).Count(&netted)
	o.FastPercent, o.NettedPercent = pct(fast, o.CreationsToday), pct(netted, o.RedemptionsToday)

	var waiting []NetBatch
	s.DB.Where("status = ?", BatchAwaitingApproval).Order("created_at").Find(&waiting)
	for _, b := range waiting {
		o.Attention = append(o.Attention, AttentionItem{Kind: "approval", Text: fmt.Sprintf("Net batch %s (%s %s, %s NGN) is above the approval threshold", b.ID, b.Side, b.Quantity, b.Value),
			Where: "Orders › Awaiting approval", Target: b.ID})
	}
	var escalated []Instruction
	s.DB.Where("status = ?", InstrEscalated).Order("updated_at").Find(&escalated)
	for _, in := range escalated {
		o.Attention = append(o.Attention, AttentionItem{Kind: "escalation", Text: fmt.Sprintf("%s instruction %s (%s %s) needs handling: %s", in.PartnerName, in.ID, in.Side, in.AssetCode, in.LastError),
			Where: "Orders › Escalations", Target: in.ID})
	}
	names := s.exchangeNames()
	type count struct {
		ServiceLinkID string
		N             int64
	}
	var dlq []count
	s.DB.Model(&WebhookDelivery{}).Select("service_link_id, COUNT(*) AS n").Where("status = ?", DeliveryDeadLetter).Group("service_link_id").Scan(&dlq)
	for _, c := range dlq {
		o.Attention = append(o.Attention, AttentionItem{Kind: "dead-letter", Text: fmt.Sprintf("%d webhook deliveries to %s dead-lettered", c.N, nameOr(names, c.ServiceLinkID)),
			Where: "Exchange Partners", Target: c.ServiceLinkID})
	}
	var late []count
	s.DB.Model(&WebhookDelivery{}).Select("service_link_id, COUNT(*) AS n").Where("escalated_at IS NOT NULL AND confirmed_at IS NULL").Group("service_link_id").Scan(&late)
	for _, c := range late {
		o.Attention = append(o.Attention, AttentionItem{Kind: "confirmation", Text: fmt.Sprintf("%s: %d dividend.paid confirmations past SLA", nameOr(names, c.ServiceLinkID), c.N),
			Where: "Corporate Actions › Exchange confirmations", Target: c.ServiceLinkID})
	}
	var drift []ReconciliationRun
	s.DB.Where("id IN (?)", s.DB.Model(&ReconciliationRun{}).Select("MAX(id)").Group("asset_id")).Where("result <> ?", ReconMatched).Find(&drift)
	for _, r := range drift {
		o.Attention = append(o.Attention, AttentionItem{Kind: "drift", Text: fmt.Sprintf("%s reconciliation: %s (%s)", r.AssetCode, r.Result, r.Detail), Where: "Reconciliation", Target: r.AssetID})
	}
	s.DB.Order("job").Find(&o.Jobs)
	s.DB.Where("created_at >= ?", today).Order("created_at DESC").Find(&o.BatchesToday)
	return o
}

func nameOr(names map[string]string, id string) string {
	if n := names[id]; n != "" {
		return n
	}
	return id
}

// ---------------------------------------------------------------- orders

// OrderFilters narrow the order list.
type OrderFilters struct {
	Type, State, Channel, Asset, Path, Search, BatchID, ServiceLinkID string
	From, To                                                          *time.Time
	Page
}

// Orders lists orders, newest first. State "open" means not final.
func (s *Service) Orders(f OrderFilters) ([]Order, int64) {
	q := s.DB.Model(&Order{})
	if f.Type != "" {
		q = q.Where("type = ?", strings.ToUpper(f.Type))
	}
	switch strings.ToLower(f.State) {
	case "":
	case "open":
		q = q.Where("state NOT IN ?", finalStates)
	default:
		q = q.Where("state = ?", strings.ToLower(f.State))
	}
	if f.Channel != "" {
		q = q.Where("channel = ?", strings.ToUpper(f.Channel))
	}
	if f.Path != "" {
		q = q.Where("path = ?", strings.ToUpper(f.Path))
	}
	if f.Asset != "" {
		q = q.Where("asset_code = ? OR asset_id = ?", strings.ToUpper(f.Asset), f.Asset)
	}
	if f.BatchID != "" {
		q = q.Where("batch_id = ?", f.BatchID)
	}
	if f.ServiceLinkID != "" {
		q = q.Where("service_link_id = ?", f.ServiceLinkID)
	}
	if f.From != nil {
		q = q.Where("created_at >= ?", *f.From)
	}
	if f.To != nil {
		q = q.Where("created_at < ?", *f.To)
	}
	if v := strings.TrimSpace(f.Search); v != "" {
		like := "%" + strings.ToLower(v) + "%"
		q = q.Where("LOWER(id) LIKE ? OR LOWER(wallet_address) LIKE ? OR LOWER(username) LIKE ? OR LOWER(external_order_ref) LIKE ?", like, like, like, like)
	}
	var total int64
	q.Count(&total)
	var out []Order
	f.Page.apply(q.Order("created_at DESC")).Find(&out)
	return out, total
}

// OrderDetail is an order with its timeline and batch.
type OrderDetail struct {
	Order        Order         `json:"order"`
	LastError    string        `json:"lastError"`
	Attempts     int           `json:"attempts"`
	Events       []OrderEvent  `json:"events"`
	Batch        *NetBatch     `json:"batch"`
	Instructions []Instruction `json:"instructions"`
	ExchangeName string        `json:"exchangeName"`
}

func (s *Service) OrderDetail(id string) (*OrderDetail, error) {
	var o Order
	if err := s.DB.First(&o, "id = ?", id).Error; err != nil {
		return nil, notFound("order", err)
	}
	out := &OrderDetail{Order: o, LastError: o.LastError, Attempts: o.Attempts}
	s.DB.Where("order_id = ?", o.ID).Order("id").Find(&out.Events)
	if o.BatchID != "" {
		var b NetBatch
		if s.DB.First(&b, "id = ?", o.BatchID).Error == nil {
			out.Batch = &b
			s.DB.Where("batch_id = ?", b.ID).Order("created_at").Find(&out.Instructions)
		}
	}
	if o.ServiceLinkID != "" {
		out.ExchangeName = nameOr(s.exchangeNames(), o.ServiceLinkID)
	}
	return out, nil
}

// ---------------------------------------------------------------- net batches

// BatchView is a net batch with its approvals.
type BatchView struct {
	NetBatch
	Approvals []BatchApproval `json:"approvals"`
	Threshold string          `json:"threshold"`
}

// Batches lists net batches (status filter, e.g. AWAITING_APPROVAL).
func (s *Service) Batches(status, asset string, p Page) ([]BatchView, int64) {
	q := s.DB.Model(&NetBatch{})
	if status != "" {
		q = q.Where("status = ?", strings.ToUpper(status))
	}
	if asset != "" {
		q = q.Where("asset_code = ?", strings.ToUpper(asset))
	}
	var total int64
	q.Count(&total)
	var batches []NetBatch
	p.apply(q.Order("created_at DESC")).Find(&batches)
	st := s.LoadSettings()
	out := make([]BatchView, 0, len(batches))
	for _, b := range batches {
		v := BatchView{NetBatch: b, Threshold: st.NetCreationThreshold}
		s.DB.Where("batch_id = ?", b.ID).Order("id").Find(&v.Approvals)
		out = append(out, v)
	}
	return out, total
}

// ApproveBatch records an approver's approval; the engine releases the
// batch once it has enough listed approvers.
func (s *Service) ApproveBatch(id, by string) (*BatchView, error) {
	var b NetBatch
	if err := s.DB.First(&b, "id = ?", id).Error; err != nil {
		return nil, notFound("batch", err)
	}
	if b.Status != BatchAwaitingApproval {
		return nil, refuse(http.StatusConflict, "the batch is %s, not awaiting approval", b.Status)
	}
	approvers := CSV(s.LoadSettings().NetCreationApprovers)
	if len(approvers) == 0 {
		return nil, refuse(http.StatusConflict, "no Net Creation Approvers are configured (Settings › Net Creation Approvers)")
	}
	if !contains(approvers, by) {
		return nil, refuse(http.StatusForbidden, "you are not a Net Creation Approver")
	}
	var n int64
	s.DB.Model(&BatchApproval{}).Where("batch_id = ? AND LOWER(approver) = ?", b.ID, strings.ToLower(by)).Count(&n)
	if n > 0 {
		return nil, refuse(http.StatusConflict, "you already approved this batch")
	}
	if err := s.DB.Create(&BatchApproval{BatchID: b.ID, Approver: strings.ToLower(by), CreatedAt: s.now()}).Error; err != nil {
		return nil, err
	}
	return s.batchView(b.ID)
}

func (s *Service) batchView(id string) (*BatchView, error) {
	var b NetBatch
	if err := s.DB.First(&b, "id = ?", id).Error; err != nil {
		return nil, notFound("batch", err)
	}
	v := BatchView{NetBatch: b, Threshold: s.LoadSettings().NetCreationThreshold}
	s.DB.Where("batch_id = ?", b.ID).Order("id").Find(&v.Approvals)
	return &v, nil
}

// RejectBatch refuses a batch above the threshold; the engine refunds its
// orders.
func (s *Service) RejectBatch(id, reason, by string) (*NetBatch, error) {
	if strings.TrimSpace(reason) == "" {
		return nil, invalid("reason", "a reason is required")
	}
	approvers := CSV(s.LoadSettings().NetCreationApprovers)
	if !contains(approvers, by) {
		return nil, refuse(http.StatusForbidden, "only a Net Creation Approver can reject a batch")
	}
	res := s.DB.Model(&NetBatch{}).Where("id = ? AND status = ?", id, BatchAwaitingApproval).Updates(map[string]interface{}{
		"status": BatchRejected, "note": trimTo("Rejected by "+by+": "+strings.TrimSpace(reason), 500), "updated_at": s.now()})
	if res.RowsAffected == 0 {
		return nil, refuse(http.StatusConflict, "only a batch awaiting approval can be rejected")
	}
	var b NetBatch
	s.DB.First(&b, "id = ?", id)
	return &b, nil
}

// ---------------------------------------------------------------- instructions

// Instructions lists outbound instructions (status filter, e.g. ESCALATED).
func (s *Service) Instructions(status, asset, batchID string, p Page) ([]Instruction, int64) {
	q := s.DB.Model(&Instruction{})
	if status != "" {
		q = q.Where("status = ?", strings.ToUpper(status))
	}
	if asset != "" {
		q = q.Where("asset_code = ?", strings.ToUpper(asset))
	}
	if batchID != "" {
		q = q.Where("batch_id = ?", batchID)
	}
	var total int64
	q.Count(&total)
	var out []Instruction
	p.apply(q.Order("updated_at DESC")).Find(&out)
	return out, total
}

func (s *Service) instruction(id string) (*Instruction, error) {
	var in Instruction
	if err := s.DB.First(&in, "id = ?", id).Error; err != nil {
		return nil, notFound("instruction", err)
	}
	return &in, nil
}

// RetryInstruction puts an escalated instruction back in the queue with
// fresh attempts.
func (s *Service) RetryInstruction(id, by string) (*Instruction, error) {
	now := s.now()
	res := s.DB.Model(&Instruction{}).Where("id = ? AND status = ?", id, InstrEscalated).Updates(map[string]interface{}{
		"status": InstrPending, "attempts": 0, "next_attempt_at": &now, "last_error": trimTo("Re-queued by "+by, 500), "updated_at": now})
	if res.RowsAffected == 0 {
		return nil, refuse(http.StatusConflict, "only an escalated instruction can be retried")
	}
	return s.instruction(id)
}

// MarkHandled closes an instruction Operations dealt with outside the
// system (it no longer retries or alerts). Record the partner's outcome
// with RecordOutcome so the batch moves on.
func (s *Service) MarkHandled(id, note, by string) (*Instruction, error) {
	if strings.TrimSpace(note) == "" {
		return nil, invalid("note", "say how it was handled")
	}
	res := s.DB.Model(&Instruction{}).Where("id = ? AND status IN ?", id, []string{InstrEscalated, InstrPending, InstrAccepted}).Updates(map[string]interface{}{
		"status": InstrHandled, "handled_by": by, "last_error": trimTo("Handled by "+by+": "+strings.TrimSpace(note), 500), "updated_at": s.now()})
	if res.RowsAffected == 0 {
		return nil, refuse(http.StatusConflict, "the instruction is already closed")
	}
	return s.instruction(id)
}

// RollBatch gives up on a batch that was not executed: its orders go back
// to wait for the next session's batch, and its instructions are closed.
func (s *Service) RollBatch(batchID, reason, by string) (*NetBatch, error) {
	if strings.TrimSpace(reason) == "" {
		return nil, invalid("reason", "a reason is required")
	}
	var executed int64
	s.DB.Model(&Instruction{}).Where("batch_id = ? AND status IN ?", batchID, []string{InstrExecuted, InstrSettled}).Count(&executed)
	if executed > 0 {
		return nil, refuse(http.StatusConflict, "the batch was already executed or settled; it cannot be rolled")
	}
	now := s.now()
	res := s.DB.Model(&NetBatch{}).Where("id = ? AND status = ?", batchID, BatchReleased).Updates(map[string]interface{}{
		"status": BatchFailed, "note": trimTo("Rolled to the next session by "+by+": "+strings.TrimSpace(reason), 500), "updated_at": now})
	if res.RowsAffected == 0 {
		return nil, refuse(http.StatusConflict, "only a released batch that has not executed can be rolled")
	}
	s.DB.Model(&Instruction{}).Where("batch_id = ? AND status NOT IN ?", batchID, []string{InstrHandled, InstrRejected}).Updates(map[string]interface{}{
		"status": InstrHandled, "handled_by": by, "last_error": "Batch rolled to the next session", "updated_at": now})
	var b NetBatch
	s.DB.First(&b, "id = ?", batchID)
	return &b, nil
}

// OutcomeRequest is a partner's outcome, recorded by Operations for a
// MANUAL partner (or one whose webhook never arrived).
type OutcomeRequest struct {
	Status             string `json:"status"` // DM: FILLED | PARTIALLY_FILLED | REJECTED; Custodian: settlement_final | failed
	ExecutedQuantity   string `json:"executedQuantity"`
	ExecutedPrice      string `json:"executedPrice"`
	SettledQuantity    string `json:"settledQuantity"`
	SettlementDate     string `json:"settlementDate"`
	CustodianReference string `json:"custodianReference"`
}

// RecordOutcome stores the outcome as the partner's callback would have
// delivered it; the engine applies it within seconds.
func (s *Service) RecordOutcome(id string, r OutcomeRequest, by string) (*PartnerEvent, error) {
	in, err := s.instruction(id)
	if err != nil {
		return nil, err
	}
	if in.Status == InstrHandled || in.Status == InstrRejected {
		return nil, refuse(http.StatusConflict, "the instruction is closed")
	}
	status := strings.TrimSpace(r.Status)
	if in.Kind == InstrDealingOrder {
		status = strings.ToUpper(status)
		switch status {
		case "FILLED":
			if !d(r.ExecutedPrice).IsPositive() {
				return nil, invalid("executedPrice", "executedPrice is required for a fill")
			}
			if r.ExecutedQuantity == "" {
				r.ExecutedQuantity = in.Quantity
			}
		case "PARTIALLY_FILLED", "REJECTED":
		default:
			return nil, invalid("status", "status must be FILLED, PARTIALLY_FILLED or REJECTED")
		}
		if in.Status == InstrExecuted {
			return nil, refuse(http.StatusConflict, "the order was already executed")
		}
		return s.recordManual("execution", fmt.Sprintf("tm:%s:%s", in.ID, status), map[string]string{"orderId": in.ID, "status": status,
			"executedQuantity": r.ExecutedQuantity, "executedPrice": r.ExecutedPrice, "executionTime": s.now().Format(time.RFC3339)}, by)
	}
	status = strings.ToLower(status)
	switch status {
	case "settlement_final":
		if r.SettledQuantity == "" {
			r.SettledQuantity = in.Quantity
		}
		if strings.TrimSpace(r.CustodianReference) == "" {
			return nil, invalid("custodianReference", "the Custodian's reference is required")
		}
	case "failed":
	default:
		return nil, invalid("status", "status must be settlement_final or failed")
	}
	if in.Status == InstrSettled {
		return nil, refuse(http.StatusConflict, "the instruction was already settled")
	}
	if r.SettlementDate == "" {
		r.SettlementDate = s.now().In(Lagos).Format("2006-01-02")
	}
	return s.recordManual("settlement", fmt.Sprintf("tm:%s:%s", in.ID, status), map[string]string{"instructionId": in.ID, "status": status,
		"settledQuantity": r.SettledQuantity, "settlementDate": r.SettlementDate, "custodianReference": r.CustodianReference}, by)
}

// ---------------------------------------------------------------- reconciliation

// ReconRow is an asset's latest reconciliation.
type ReconRow struct {
	ReconciliationRun
	CustodianName string  `json:"custodianName"`
	AssetStatus   string  `json:"assetStatus"`
	OrdersSince   []Order `json:"ordersSince"` // settled after the position's as-of: likely causes of a drift
}

// Reconciliation is the reconciliation page.
type Reconciliation struct {
	Checked   int        `json:"checked"`
	Matched   int        `json:"matched"`
	Drift     int        `json:"drift"`
	Stale     int        `json:"stalePrices"`
	LastRunAt *time.Time `json:"lastRunAt"`
	NextRunAt time.Time  `json:"nextRunAt"`
	Rows      []ReconRow `json:"rows"`
}

func (s *Service) Reconciliation() Reconciliation {
	st := s.LoadSettings()
	var out Reconciliation
	var runs []ReconciliationRun
	s.DB.Where("id IN (?)", s.DB.Model(&ReconciliationRun{}).Select("MAX(id)").Group("asset_id")).Order("asset_code").Find(&runs)
	for _, r := range runs {
		row := ReconRow{ReconciliationRun: r}
		var a Asset
		if s.DB.First(&a, "id = ?", r.AssetID).Error == nil {
			row.CustodianName, row.AssetStatus = s.custodianName(a.CustodianID), a.Status
		}
		out.Checked++
		if r.Result == ReconMatched {
			out.Matched++
		} else {
			out.Drift++
			if r.PositionAsOf != nil {
				s.DB.Where("asset_id = ? AND settled_at > ?", r.AssetID, *r.PositionAsOf).Order("settled_at").Limit(20).Find(&row.OrdersSince)
			}
		}
		if r.PriceStale {
			out.Stale++
		}
		if out.LastRunAt == nil || r.CreatedAt.After(*out.LastRunAt) {
			t := r.CreatedAt
			out.LastRunAt = &t
		}
		out.Rows = append(out.Rows, row)
	}
	now := s.now().In(Lagos)
	next := time.Date(now.Year(), now.Month(), now.Day(), st.ReconciliationHour, 0, 0, 0, Lagos)
	if !next.After(now) {
		next = next.AddDate(0, 0, 1)
	}
	out.NextRunAt = next.UTC()
	return out
}

// ReconciliationHistory is an asset's runs, newest first.
func (s *Service) ReconciliationHistory(asset string, p Page) ([]ReconciliationRun, int64) {
	q := s.DB.Model(&ReconciliationRun{})
	if asset != "" {
		q = q.Where("asset_code = ? OR asset_id = ?", strings.ToUpper(asset), asset)
	}
	var total int64
	q.Count(&total)
	var out []ReconciliationRun
	p.apply(q.Order("id DESC")).Find(&out)
	return out, total
}

// RequestReconciliation asks the engine to reconcile an asset (or ALL) now.
func (s *Service) RequestReconciliation(asset, by string) (*JobRequest, error) {
	target := strings.ToUpper(strings.TrimSpace(asset))
	if target == "" {
		target = "ALL"
	}
	if target != "ALL" {
		a, err := s.asset(target)
		if err != nil {
			return nil, err
		}
		target = a.AssetCode
	}
	return s.requestJob("RECONCILE", target, by)
}

// RequestPositionFeed asks a mock Custodian for its position feed now.
func (s *Service) RequestPositionFeed(code, by string) (*JobRequest, error) {
	var c Custodian
	if err := s.DB.First(&c, "code = ?", strings.ToUpper(strings.TrimSpace(code))).Error; err != nil {
		return nil, notFound("custodian", err)
	}
	if c.Mode != ModeMock {
		return nil, refuse(http.StatusConflict, "only a mock Custodian sends its feed on request; record a MANUAL Custodian's position from its statement")
	}
	return s.requestJob("POSITION_FEED", c.Code, by)
}

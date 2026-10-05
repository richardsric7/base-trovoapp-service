package publicmarkets

import (
	"net/http"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// Corporate actions are distributed by app-backend's own dividend engine
// (not the proceeds payout): record-date snapshot, per-holder withholding
// tax, approvals on the snapshot's checksum, then payment.

// CorporateActions lists corporate actions, newest first.
func (s *Service) CorporateActions(status, asset string, p Page) ([]CorporateAction, int64) {
	q := s.DB.Model(&CorporateAction{})
	if status != "" {
		q = q.Where("status = ?", strings.ToUpper(status))
	}
	if asset != "" {
		q = q.Where("asset_code = ?", strings.ToUpper(asset))
	}
	var total int64
	q.Count(&total)
	var out []CorporateAction
	p.apply(q.Order("created_at DESC")).Find(&out)
	return out, total
}

// ChannelTotals sums a distribution per channel.
type ChannelTotals struct {
	Channel    string `json:"channel"`
	Recipients int    `json:"recipients"`
	Gross      string `json:"gross"`
	WHT        string `json:"wht"`
	Net        string `json:"net"`
	Paid       int    `json:"paid"`
	Failed     int    `json:"failed"`
}

// ActionDetail is a corporate action with its distribution.
type ActionDetail struct {
	CorporateAction
	Approvals    []DividendApproval `json:"approvals"`
	Approvers    []string           `json:"approvers"`
	ByChannel    []ChannelTotals    `json:"byChannel"`
	Entitlements []Entitlement      `json:"entitlements"`
	Total        int64              `json:"totalEntitlements"`
}

// ActionDetail returns a corporate action and a page of its entitlements.
func (s *Service) ActionDetail(id, status string, p Page) (*ActionDetail, error) {
	var ca CorporateAction
	if err := s.DB.First(&ca, "id = ?", id).Error; err != nil {
		return nil, notFound("corporate action", err)
	}
	out := &ActionDetail{CorporateAction: ca, Approvers: CSV(s.LoadSettings().DividendApprovers)}
	s.DB.Where("corporate_action_id = ?", ca.ID).Order("id").Find(&out.Approvals)
	var all []Entitlement
	s.DB.Where("corporate_action_id = ?", ca.ID).Find(&all)
	totals := map[string]*ChannelTotals{}
	var order []string
	sums := map[string][3]decimal.Decimal{}
	for _, e := range all {
		t := totals[e.Channel]
		if t == nil {
			t = &ChannelTotals{Channel: e.Channel}
			totals[e.Channel] = t
			order = append(order, e.Channel)
		}
		t.Recipients++
		switch e.Status {
		case EntitlementPaid:
			t.Paid++
		case EntitlementFailed:
			t.Failed++
		}
		v := sums[e.Channel]
		sums[e.Channel] = [3]decimal.Decimal{v[0].Add(d(e.GrossAmount)), v[1].Add(d(e.WHTAmount)), v[2].Add(d(e.NetAmount))}
	}
	for _, ch := range order {
		t, v := totals[ch], sums[ch]
		t.Gross, t.WHT, t.Net = v[0].StringFixed(2), v[1].StringFixed(2), v[2].StringFixed(2)
		out.ByChannel = append(out.ByChannel, *t)
	}
	q := s.DB.Model(&Entitlement{}).Where("corporate_action_id = ?", ca.ID)
	if status != "" {
		q = q.Where("status = ?", strings.ToUpper(status))
	}
	q.Count(&out.Total)
	p.apply(q.Order("id")).Find(&out.Entitlements)
	return out, nil
}

// DeclareRequest is a corporate action entered by Operations.
type DeclareRequest struct {
	AssetCode     string `json:"assetCode"`
	EventType     string `json:"eventType"`
	RecordDate    string `json:"recordDate"`
	PayDate       string `json:"payDate"`
	AmountPerUnit string `json:"amountPerUnit"`
	Currency      string `json:"currency"`
	Description   string `json:"description"`
	Reference     string `json:"reference"`
}

// Declare records a corporate action for the engine, which declares it as
// it would a Custodian's notice (within seconds; the result is on the
// returned event).
func (s *Service) Declare(r DeclareRequest, by string) (*PartnerEvent, error) {
	a, err := s.asset(r.AssetCode)
	if err != nil {
		return nil, err
	}
	typ := strings.ToUpper(strings.TrimSpace(r.EventType))
	record, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(r.RecordDate), Lagos)
	if err != nil {
		return nil, invalid("recordDate", "recordDate must be YYYY-MM-DD")
	}
	switch typ {
	case ActionDividend, ActionCoupon:
		if !d(r.AmountPerUnit).IsPositive() {
			return nil, invalid("amountPerUnit", "amountPerUnit must be positive")
		}
		pay, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(r.PayDate), Lagos)
		if err != nil {
			return nil, invalid("payDate", "payDate must be YYYY-MM-DD")
		}
		if pay.Before(record) {
			return nil, invalid("payDate", "payDate is before the record date")
		}
		if c := strings.ToUpper(strings.TrimSpace(r.Currency)); c != "" && c != "NGN" {
			return nil, invalid("currency", "only NGN distributions are supported")
		}
	case ActionBonus, ActionRights, ActionSplit:
	default:
		return nil, invalid("eventType", "eventType must be DIVIDEND, COUPON, BONUS, RIGHTS or SPLIT")
	}
	var n int64
	s.DB.Model(&CorporateAction{}).Where("asset_id = ? AND event_type = ? AND record_date = ? AND status <> ?", a.ID, typ, r.RecordDate, ActionCancelledCA).Count(&n)
	if n > 0 {
		return nil, refuse(http.StatusConflict, "a %s with this record date is already declared for %s", strings.ToLower(typ), a.AssetCode)
	}
	ref := strings.TrimSpace(r.Reference)
	if ref == "" {
		ref = "tm:" + a.AssetCode + ":" + typ + ":" + r.RecordDate + ":" + randomHex(3)
	}
	return s.recordManual("corporate-action", ref, map[string]string{"assetCode": a.AssetCode, "eventType": typ, "recordDate": r.RecordDate,
		"payDate": r.PayDate, "amountPerUnit": strings.TrimSpace(r.AmountPerUnit), "currency": "NGN", "description": r.Description, "reference": ref}, by)
}

// ApproveAction approves a snapshotted distribution: the approval is on the
// snapshot's checksum, so a re-taken snapshot needs fresh approvals.
func (s *Service) ApproveAction(id, by string) (*ActionDetail, error) {
	var ca CorporateAction
	if err := s.DB.First(&ca, "id = ?", id).Error; err != nil {
		return nil, notFound("corporate action", err)
	}
	if ca.Status != ActionSnapshotted {
		return nil, refuse(http.StatusConflict, "the distribution is %s; only a snapshotted one awaits approval", ca.Status)
	}
	approvers := CSV(s.LoadSettings().DividendApprovers)
	if len(approvers) == 0 {
		return nil, refuse(http.StatusConflict, "no Dividend Approvers are configured (Settings › Thresholds & Limits)")
	}
	if !contains(approvers, by) {
		return nil, refuse(http.StatusForbidden, "you are not a Dividend Approver")
	}
	var n int64
	s.DB.Model(&DividendApproval{}).Where("corporate_action_id = ? AND LOWER(approver) = ?", ca.ID, strings.ToLower(by)).Count(&n)
	if n > 0 {
		return nil, refuse(http.StatusConflict, "you already approved this distribution")
	}
	if err := s.DB.Create(&DividendApproval{CorporateActionID: ca.ID, Approver: strings.ToLower(by), Checksum: ca.SnapshotChecksum, CreatedAt: s.now()}).Error; err != nil {
		return nil, err
	}
	return s.ActionDetail(ca.ID, "", Page{Limit: 50})
}

// CancelAction cancels a corporate action that has not started paying.
func (s *Service) CancelAction(id, reason, by string) (*CorporateAction, error) {
	if strings.TrimSpace(reason) == "" {
		return nil, invalid("reason", "a reason is required")
	}
	res := s.DB.Model(&CorporateAction{}).Where("id = ? AND status IN ?", id, []string{ActionAnnounced, ActionSnapshotted, ActionApproved, ActionNeedsManual}).
		Updates(map[string]interface{}{"status": ActionCancelledCA, "note": trimTo("Cancelled by "+by+": "+strings.TrimSpace(reason), 500), "updated_at": s.now()})
	if res.RowsAffected == 0 {
		return nil, refuse(http.StatusConflict, "only an action that has not started paying can be cancelled")
	}
	var ca CorporateAction
	s.DB.First(&ca, "id = ?", id)
	return &ca, nil
}

// ConfirmationRow is an exchange's dividend.paid confirmations for an asset.
type ConfirmationRow struct {
	ServiceLinkID string     `json:"serviceLinkId"`
	ExchangeName  string     `json:"exchangeName"`
	AssetCode     string     `json:"assetCode"`
	WalletsPaid   int64      `json:"walletsPaid"`
	Confirmed     int64      `json:"confirmed"`
	Outstanding   int64      `json:"outstanding"`
	Escalated     int64      `json:"escalated"`
	OldestDueAt   *time.Time `json:"oldestDueAt"`
}

// Confirmations summarises exchanges' confirmations of dividend.paid.
func (s *Service) Confirmations() []ConfirmationRow {
	var rows []ConfirmationRow
	s.DB.Model(&WebhookDelivery{}).Select(`service_link_id, asset_code, COUNT(*) AS wallets_paid,
		SUM(CASE WHEN confirmed_at IS NOT NULL THEN 1 ELSE 0 END) AS confirmed,
		SUM(CASE WHEN confirmed_at IS NULL THEN 1 ELSE 0 END) AS outstanding,
		SUM(CASE WHEN confirmed_at IS NULL AND escalated_at IS NOT NULL THEN 1 ELSE 0 END) AS escalated`).
		Where("event = ? AND needs_confirmation = ?", "dividend.paid", true).Group("service_link_id, asset_code").Order("service_link_id, asset_code").Scan(&rows)
	names := s.exchangeNames()
	for i := range rows {
		rows[i].ExchangeName = nameOr(names, rows[i].ServiceLinkID)
		var w WebhookDelivery
		if s.DB.Where("service_link_id = ? AND asset_code = ? AND event = ? AND confirmed_at IS NULL AND confirmation_due_at IS NOT NULL",
			rows[i].ServiceLinkID, rows[i].AssetCode, "dividend.paid").Order("confirmation_due_at").First(&w).Error == nil {
			rows[i].OldestDueAt = w.ConfirmationDueAt
		}
	}
	return rows
}

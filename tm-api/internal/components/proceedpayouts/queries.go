package proceedpayouts

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"

	stakeholderModels "admin-panel-dashboard/internal/components/stakeholder/models"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Filters narrow payout lists and reports.
type Filters struct {
	Status  string // one status, or empty
	AssetID string // tokenized asset ID or asset code
	From    *time.Time
	To      *time.Time // exclusive
	Page    int
	Limit   int
}

func (f *Filters) pageLimit() (int, int) {
	page, limit := f.Page, f.Limit
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 200 {
		limit = 50
	}
	return page, limit
}

// PayoutView is a payout with its asset and amounts in human units.
type PayoutView struct {
	ProceedPayout
	AssetCode   string `json:"assetCode"`
	AssetName   string `json:"assetName"`
	HolderShare string `json:"holderPayable"` // TotalAmount less fee and VAT
	Fee         string `json:"fee"`
	Vat         string `json:"vat"`
	Paid        string `json:"paid"`      // to holders
	Payable     string `json:"payable"`   // to eligible holders
	Retained    string `json:"retained"`  // excluded holders' share and rounding dust
	Approvals   int    `json:"approvals"` // for the current schedule
}

type assetInfo struct {
	ID        string
	AssetCode string
	AssetName string
}

func (s *Service) assetsByID(ids []string) map[string]assetInfo {
	out := map[string]assetInfo{}
	if len(ids) == 0 {
		return out
	}
	var rows []assetInfo
	s.DB.Table("tokenized_assets").Select("id, asset_code, asset_name").Where("id IN ?", ids).Scan(&rows)
	for _, r := range rows {
		out[r.ID] = r
	}
	return out
}

func (s *Service) views(payouts []ProceedPayout) []PayoutView {
	ids := make([]string, 0, len(payouts))
	pids := make([]uint64, 0, len(payouts))
	for _, p := range payouts {
		ids = append(ids, p.TokenizedAssetID)
		pids = append(pids, p.ID)
	}
	assets := s.assetsByID(ids)
	type count struct {
		ProceedPayoutID  uint64
		ScheduleChecksum string
		N                int
	}
	var counts []count
	if len(pids) > 0 {
		s.DB.Model(&PayoutApproval{}).Select("proceed_payout_id, schedule_checksum, COUNT(*) AS n").
			Where("proceed_payout_id IN ?", pids).Group("proceed_payout_id, schedule_checksum").Scan(&counts)
	}
	out := make([]PayoutView, len(payouts))
	for i, p := range payouts {
		v := PayoutView{ProceedPayout: p, AssetCode: assets[p.TokenizedAssetID].AssetCode, AssetName: assets[p.TokenizedAssetID].AssetName}
		d := p.PayoutDecimals
		v.HolderShare, v.Fee, v.Vat = human(p.HolderPayable, d), human(p.FeeUnits, d), human(p.VatUnits, d)
		v.Paid, v.Payable, v.Retained = human(p.PaidUnits, d), human(p.PayableUnits, d), human(p.RetainedUnits, d)
		for _, c := range counts {
			if c.ProceedPayoutID == p.ID && c.ScheduleChecksum == p.ScheduleChecksum && p.ScheduleChecksum != "" {
				v.Approvals = c.N
			}
		}
		out[i] = v
	}
	return out
}

func (s *Service) filtered(f Filters) *gorm.DB {
	q := s.DB.Model(&ProceedPayout{})
	if st := strings.ToUpper(strings.TrimSpace(f.Status)); st != "" {
		q = q.Where("status = ?", st)
	}
	if a := strings.TrimSpace(f.AssetID); a != "" {
		q = q.Where("tokenized_asset_id IN (?)", s.DB.Table("tokenized_assets").Select("id").Where("id = ? OR asset_code = ?", a, a))
	}
	if f.From != nil {
		q = q.Where("created_at >= ?", *f.From)
	}
	if f.To != nil {
		q = q.Where("created_at < ?", *f.To)
	}
	return q
}

// List is a page of payouts, newest first.
func (s *Service) List(f Filters) ([]PayoutView, int64, error) {
	var total int64
	if err := s.filtered(f).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	page, limit := f.pageLimit()
	var rows []ProceedPayout
	if err := s.filtered(f).Order("id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	for i := range rows {
		s.syncDistribution(&rows[i])
	}
	return s.views(rows), total, nil
}

// Detail is a payout with its approvals and batches.
type Detail struct {
	PayoutView
	ApprovalList []PayoutApproval `json:"approvalList"`
	Batches      []PayoutBatch    `json:"batches"`
	ItemCounts   map[string]int   `json:"itemCounts"` // holders by status
	Engine       *EngineStatus    `json:"engine"`
}

func (s *Service) Detail(id uint64) (*Detail, error) {
	p, err := s.load(id)
	if err != nil {
		return nil, err
	}
	s.syncDistribution(p)
	d := &Detail{PayoutView: s.views([]ProceedPayout{*p})[0], ItemCounts: map[string]int{}}
	s.DB.Where("proceed_payout_id = ?", id).Order("id").Find(&d.ApprovalList)
	s.DB.Where("proceed_payout_id = ?", id).Order("id DESC").Limit(200).Find(&d.Batches)
	type agg struct {
		Status string
		N      int
	}
	var rows []agg
	s.DB.Model(&PayoutItem{}).Select("status, COUNT(*) AS n").
		Where("proceed_payout_id = ? AND (kind = ? OR kind = '' OR kind IS NULL)", id, KindHolder).Group("status").Scan(&rows)
	for _, r := range rows {
		d.ItemCounts[r.Status] = r.N
	}
	d.Engine, _ = s.Engine()
	return d, nil
}

// ItemFilters narrow a payout's schedule.
type ItemFilters struct {
	Status string
	Kind   string
	Search string // address or username
	Page   int
	Limit  int
}

// Items is a page of a payout's schedule (fee and VAT lines first, then
// holders by size).
func (s *Service) Items(id uint64, f ItemFilters) ([]PayoutItem, int64, error) {
	q := s.DB.Model(&PayoutItem{}).Where("proceed_payout_id = ?", id)
	if st := strings.ToUpper(strings.TrimSpace(f.Status)); st != "" {
		q = q.Where("status = ?", st)
	}
	if k := strings.ToUpper(strings.TrimSpace(f.Kind)); k != "" {
		q = q.Where("kind = ?", k)
	}
	if t := strings.TrimSpace(f.Search); t != "" {
		q = q.Where("LOWER(beneficiary_address) = LOWER(?) OR LOWER(username) LIKE LOWER(?)", t, "%"+t+"%")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	page, limit := (&Filters{Page: f.Page, Limit: f.Limit}).pageLimit()
	var rows []PayoutItem
	err := q.Order("CASE WHEN kind = 'HOLDER' OR kind = '' OR kind IS NULL THEN 1 ELSE 0 END, amount_to_receive DESC, beneficiary_address").
		Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

// --- the trustee's view (stakeholder distribution detail and CSV) ---

// DistributionPayout is a distribution's payout in the stakeholder
// portal's shape, or nil before it is registered.
func (s *Service) DistributionPayout(ctx context.Context, distributionID string) (*stakeholderModels.DistributionPayoutResponse, error) {
	var p ProceedPayout
	if err := s.DB.WithContext(ctx).Where("distribution_id = ?", distributionID).First(&p).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	s.syncDistribution(&p)
	out := &stakeholderModels.DistributionPayoutResponse{
		DistributionID: p.DistributionID, TokenizedAssetID: p.TokenizedAssetID, Batch: p.Batch,
		Amount: p.TotalAmount, Currency: p.PayoutAssetCode, AmountPerToken: p.AmountPerToken,
		Status: strings.ToLower(p.Status), PaymentScheduleReady: p.PaymentScheduleReady == 1,
		PayoutCompleted: p.Status == StatusCompleted, CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
		Payouts: make([]stakeholderModels.DistributionPayoutRecordResponse, 0),
	}
	if p.PaymentScheduleReady != 1 {
		return out, nil
	}
	var items []PayoutItem
	if err := s.DB.WithContext(ctx).Where("proceed_payout_id = ? AND (kind = ? OR kind = '' OR kind IS NULL)", p.ID, KindHolder).
		Order("amount_to_receive DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	for _, it := range items {
		out.Payouts = append(out.Payouts, stakeholderModels.DistributionPayoutRecordResponse{
			ID: it.ID, BeneficiaryAddress: it.BeneficiaryAddress, ConfirmedTokenBalance: human(it.BalanceUnits, p.TokenDecimals),
			Amount: human(it.AmountUnits, p.PayoutDecimals), Currency: p.PayoutAssetCode,
			CannotReceiveAsset: it.CannotReceiveAsset == 1 || it.Status == ItemExcluded, Paid: it.Status == ItemPaid, CreatedAt: it.CreatedAt,
		})
	}
	return out, nil
}

// --- reports ---

// CurrencyTotals sums payouts in one payout currency.
type CurrencyTotals struct {
	Currency      string `json:"currency"`
	Payouts       int    `json:"payouts"`
	Total         string `json:"total"`         // authorized
	HolderPayable string `json:"holderPayable"` // after fee and VAT
	PaidToHolders string `json:"paidToHolders"`
	Fees          string `json:"fees"`
	Vat           string `json:"vat"`
	HoldersPaid   int    `json:"holdersPaid"`
	HoldersFailed int    `json:"holdersFailed"`
}

// PayoutsReport lists the matching payouts (all of them, newest first)
// with totals per payout currency and counts per status.
type PayoutsReport struct {
	Payouts  []PayoutView     `json:"payouts"`
	Totals   []CurrencyTotals `json:"totals"`
	ByStatus map[string]int   `json:"byStatus"`
}

func (s *Service) PayoutsReport(f Filters) (*PayoutsReport, error) {
	var rows []ProceedPayout
	if err := s.filtered(f).Order("id DESC").Limit(5000).Find(&rows).Error; err != nil {
		return nil, err
	}
	r := &PayoutsReport{Payouts: s.views(rows), ByStatus: map[string]int{}}
	type acc struct {
		n                              int
		total, payable, paid, fee, vat decimal.Decimal
		holdersPaid, holdersFailed     int
	}
	byCurrency := map[string]*acc{}
	for _, v := range r.Payouts {
		r.ByStatus[v.Status]++
		a := byCurrency[v.PayoutAssetCode]
		if a == nil {
			a = &acc{}
			byCurrency[v.PayoutAssetCode] = a
		}
		a.n++
		a.total = a.total.Add(dec(v.TotalAmount))
		a.payable = a.payable.Add(dec(v.HolderShare))
		a.paid = a.paid.Add(dec(v.Paid))
		a.fee = a.fee.Add(dec(v.Fee))
		a.vat = a.vat.Add(dec(v.Vat))
		a.holdersPaid += v.PaidCount
		a.holdersFailed += v.FailedCount
	}
	for c, a := range byCurrency {
		r.Totals = append(r.Totals, CurrencyTotals{Currency: c, Payouts: a.n, Total: a.total.String(), HolderPayable: a.payable.String(),
			PaidToHolders: a.paid.String(), Fees: a.fee.String(), Vat: a.vat.String(), HoldersPaid: a.holdersPaid, HoldersFailed: a.holdersFailed})
	}
	sort.Slice(r.Totals, func(i, j int) bool { return r.Totals[i].Currency < r.Totals[j].Currency })
	return r, nil
}

// FeeTotals sums collected payout fees and VAT per asset and currency.
type FeeTotals struct {
	AssetCode       string  `json:"assetCode"`
	PayoutAssetCode string  `json:"payoutAssetCode"`
	Fees            float64 `json:"fees"`
	Vat             float64 `json:"vat"`
	Count           int     `json:"count"`
}

// FeesReport is the payout fees and VAT collected (paid on-chain to the fee
// and VAT wallets), from fee_collections.
type FeesReport struct {
	Collections []FeeCollection    `json:"collections"`
	Totals      []FeeTotals        `json:"totals"`
	TotalFees   map[string]float64 `json:"totalFees"` // per payout currency
	TotalVat    map[string]float64 `json:"totalVat"`
}

func (s *Service) FeesReport(f Filters) (*FeesReport, error) {
	q := s.DB.Model(&FeeCollection{}).Where("fee_type IN ?", []string{FeeTypePayout, FeeTypePayoutVat})
	if a := strings.TrimSpace(f.AssetID); a != "" {
		q = q.Where("from_username IN (?)", s.DB.Table("tokenized_assets").Select("asset_code").Where("id = ? OR asset_code = ?", a, a))
	}
	if f.From != nil {
		q = q.Where("created_at >= ?", *f.From)
	}
	if f.To != nil {
		q = q.Where("created_at < ?", *f.To)
	}
	r := &FeesReport{TotalFees: map[string]float64{}, TotalVat: map[string]float64{}}
	if err := q.Order("created_at DESC").Limit(5000).Find(&r.Collections).Error; err != nil {
		return nil, err
	}
	byKey := map[[2]string]*FeeTotals{}
	for _, c := range r.Collections {
		k := [2]string{c.FromUsername, c.AssetCode}
		t := byKey[k]
		if t == nil {
			t = &FeeTotals{AssetCode: c.FromUsername, PayoutAssetCode: c.AssetCode}
			byKey[k] = t
		}
		t.Count++
		if c.FeeType == FeeTypePayoutVat {
			t.Vat += c.Amount
			r.TotalVat[c.AssetCode] += c.Amount
		} else {
			t.Fees += c.Amount
			r.TotalFees[c.AssetCode] += c.Amount
		}
	}
	for _, t := range byKey {
		r.Totals = append(r.Totals, *t)
	}
	sort.Slice(r.Totals, func(i, j int) bool {
		if r.Totals[i].AssetCode != r.Totals[j].AssetCode {
			return r.Totals[i].AssetCode < r.Totals[j].AssetCode
		}
		return r.Totals[i].PayoutAssetCode < r.Totals[j].PayoutAssetCode
	})
	return r, nil
}

func dec(s string) decimal.Decimal {
	d, err := decimal.NewFromString(strings.TrimSpace(s))
	if err != nil {
		return decimal.Zero
	}
	return d
}

package proceedpayouts

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"admin-panel-dashboard/internal/cache"
	stakeholderDB "admin-panel-dashboard/internal/components/stakeholder/db"
	stakeholderModels "admin-panel-dashboard/internal/components/stakeholder/models"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Service manages payouts in app-backend's database (DB) and keeps the
// stakeholder distributions (AdminDB) in step with them.
type Service struct {
	DB      *gorm.DB
	AdminDB *gorm.DB
	Cache   *cache.RedisCache // optional: wakes the engine
}

// Error is a refusal with its HTTP status.
type Error struct {
	Status  int
	Message string
}

func (e *Error) Error() string { return e.Message }

func refuse(status int, format string, args ...interface{}) error {
	return &Error{Status: status, Message: fmt.Sprintf(format, args...)}
}

func notFound(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return refuse(http.StatusNotFound, "payout not found")
	}
	return err
}

// approvalsRequired is how many distinct admins approve a schedule
// (PROCEED_PAYOUT_APPROVALS_REQUIRED, default 2, at least 1).
func approvalsRequired() int {
	n, err := strconv.Atoi(strings.TrimSpace(os.Getenv("PROCEED_PAYOUT_APPROVALS_REQUIRED")))
	if err != nil {
		return 2
	}
	if n < 1 {
		return 1
	}
	return n
}

func now() time.Time { return time.Now().UTC() }

// wake tells the engine something changed (best effort: it also polls).
func (s *Service) wake(action string, id uint64) {
	if s.Cache == nil || !s.Cache.Enabled || s.Cache.Client == nil {
		return
	}
	raw, _ := json.Marshal(map[string]interface{}{"action": action, "payoutId": id})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := s.Cache.Client.Publish(ctx, CommandsChannel, raw).Err(); err != nil {
		log.Printf("[proceedpayouts] waking payout-engine: %v", err)
	}
}

func (s *Service) load(id uint64) (*ProceedPayout, error) {
	var p ProceedPayout
	if err := s.DB.First(&p, id).Error; err != nil {
		return nil, notFound(err)
	}
	return &p, nil
}

// transition moves a payout from one of from to the updates' status, only
// if it is still there (the engine or another admin may have moved it).
func (s *Service) transition(tx *gorm.DB, id uint64, from []string, updates map[string]interface{}) error {
	res := tx.Model(&ProceedPayout{}).Where("id = ? AND status IN ?", id, from).Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		var p ProceedPayout
		if err := tx.Select("status").First(&p, id).Error; err != nil {
			return notFound(err)
		}
		return refuse(http.StatusConflict, "not possible while the payout is %s", p.Status)
	}
	return nil
}

// --- registration (the trustee's authorization) ---

// Registration is a distribution to pay out.
type Registration struct {
	DistributionID   string
	TokenizedAssetID string
	Amount           string
	Currency         string
}

func (s *Service) asset(id string) (*stakeholderDB.TokenizedAsset, error) {
	var a stakeholderDB.TokenizedAsset
	err := s.DB.Select("id", "asset_code", "asset_name", "asset_country_location", "proceed_payout_currency", "contract_address").
		Where("id = ?", strings.TrimSpace(id)).First(&a).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, refuse(http.StatusUnprocessableEntity, "tokenized asset %s not found", id)
		}
		return nil, err
	}
	return &a, nil
}

// payoutToken picks the token holders are paid in: the first of the
// distribution's currency, the asset's proceed payout currency and its
// country's internal balance token that is a tokenization currency with a
// contract.
func (s *Service) payoutToken(currency string, a *stakeholderDB.TokenizedAsset) (code, contract string) {
	candidates := []string{currency, a.ProceedPayoutCurrency}
	var cc countryConfig
	if s.DB.Where("country_code = ?", countryOf(a)).First(&cc).Error == nil && cc.InternalBalanceTokenCode != nil {
		candidates = append(candidates, *cc.InternalBalanceTokenCode)
	}
	for _, c := range candidates {
		c = strings.ToUpper(strings.TrimSpace(c))
		if c == "" {
			continue
		}
		var tc tokenizationCurrency
		if s.DB.Where("asset_code = ?", c).First(&tc).Error == nil && common.IsHexAddress(tc.ContractAddress) {
			return tc.AssetCode, common.HexToAddress(tc.ContractAddress).Hex()
		}
	}
	return strings.ToUpper(strings.TrimSpace(currency)), ""
}

func countryOf(a *stakeholderDB.TokenizedAsset) string {
	if c := strings.ToUpper(strings.TrimSpace(a.AssetCountryLocation)); c != "" {
		return c
	}
	return "NG"
}

// FeeDefaults are the fee terms new payouts start with (the
// PROCEED_PAYOUT_FEE service fee; none configured: FIXED 0, cap 0).
type FeeDefaults struct {
	FeeType  string `json:"feeType"`
	FeeValue string `json:"feeValue"`
	FeeCap   string `json:"feeCap"`
}

func (s *Service) feeDefaults() FeeDefaults {
	d := FeeDefaults{FeeType: FeeFixed, FeeValue: "0", FeeCap: "0"}
	var row ServiceFee
	if s.DB.Where("id = ? AND inactive = 0", FeeServiceID).First(&row).Error != nil {
		return d
	}
	if row.FeePercent > 0 {
		d.FeeType, d.FeeValue = FeePercent, decimal.NewFromFloat(row.FeePercent).String()
		if c, err := decimal.NewFromString(strings.TrimSpace(row.Remarks)); err == nil && c.IsPositive() {
			d.FeeCap = c.String()
		}
	} else if row.FeeFixed > 0 {
		d.FeeValue = decimal.NewFromFloat(row.FeeFixed).String()
	}
	return d
}

// Register records an authorized distribution's payout (once: a second
// call returns the same payout).
func (s *Service) Register(ctx context.Context, r Registration) (*ProceedPayout, error) {
	amount, err := decimal.NewFromString(strings.TrimSpace(r.Amount))
	if err != nil || !amount.IsPositive() {
		return nil, refuse(http.StatusBadRequest, "invalid distribution amount %q", r.Amount)
	}
	a, err := s.asset(r.TokenizedAssetID)
	if err != nil {
		return nil, err
	}
	var existing ProceedPayout
	if err := s.DB.WithContext(ctx).Where("distribution_id = ?", r.DistributionID).First(&existing).Error; err == nil {
		return &existing, nil
	}
	code, contract := s.payoutToken(r.Currency, a)
	fee := s.feeDefaults()
	label := a.AssetCode
	if label == "" {
		label = a.ID
	}
	p := ProceedPayout{
		CreatedAt: now(), UpdatedAt: now(), TokenizedAssetID: a.ID, Batch: label + "-" + r.DistributionID, DistributionID: r.DistributionID,
		Status: StatusRegistered, PayoutAssetCode: code, PayoutContractAddress: contract, TotalAmount: amount.String(),
		ApprovalsRequired: approvalsRequired(), FeeType: fee.FeeType, FeeValue: fee.FeeValue, FeeCap: fee.FeeCap,
		PaidUnits: "0",
	}
	if contract == "" {
		p.Note = fmt.Sprintf("No tokenization currency with a token contract matches %s; add it before preparing the schedule.", code)
	}
	if err := s.DB.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&p).Error; err != nil {
		return nil, err
	}
	if p.ID == 0 { // registered concurrently
		if err := s.DB.WithContext(ctx).Where("batch = ?", p.Batch).First(&p).Error; err != nil {
			return nil, err
		}
	}
	return &p, nil
}

// --- distribution status ---

// distributionStatus is the stakeholder distribution's status for a payout
// status ("" leaves it as it is).
func distributionStatus(payoutStatus string) string {
	switch payoutStatus {
	case StatusRegistered:
		return ""
	case StatusCompleted:
		return stakeholderModels.DistributionStatusCompleted
	case StatusCancelled:
		return stakeholderModels.DistributionStatusFailed
	default:
		return stakeholderModels.DistributionStatusProcessing
	}
}

// syncDistribution moves the payout's distribution to the status matching
// the payout (from authorized or processing only).
func (s *Service) syncDistribution(p *ProceedPayout) {
	if s.AdminDB == nil || p.DistributionID == "" {
		return
	}
	to := distributionStatus(p.Status)
	if to == "" {
		return
	}
	err := s.AdminDB.Model(&stakeholderModels.Distribution{}).
		Where("id = ? AND status IN ? AND status <> ?", p.DistributionID,
			[]string{stakeholderModels.DistributionStatusAuthorized, stakeholderModels.DistributionStatusProcessing}, to).
		Updates(map[string]interface{}{"status": to, "updated_at": time.Now()}).Error
	if err != nil {
		log.Printf("[proceedpayouts] syncing distribution %s: %v", p.DistributionID, err)
	}
}

// after re-reads a payout after an action, syncs its distribution and
// wakes the engine.
func (s *Service) after(action string, id uint64) (*ProceedPayout, error) {
	p, err := s.load(id)
	if err != nil {
		return nil, err
	}
	s.syncDistribution(p)
	s.wake(action, id)
	return p, nil
}

// --- stage controls ---

// startedPaying reports whether any of the payout's schedule went out (a
// batch was sent, or an item was paid): such a schedule can no longer be
// prepared again or rejected, or holders would be paid twice.
func (s *Service) startedPaying(tx *gorm.DB, id uint64) (bool, error) {
	var n int64
	if err := tx.Model(&PayoutBatch{}).Where("proceed_payout_id = ?", id).Count(&n).Error; err != nil || n > 0 {
		return n > 0, err
	}
	err := tx.Model(&PayoutItem{}).Where("proceed_payout_id = ? AND status IN ?", id, []string{ItemPaid, ItemQueued}).Count(&n).Error
	return n > 0, err
}

func (s *Service) refuseIfPaying(tx *gorm.DB, id uint64) error {
	started, err := s.startedPaying(tx, id)
	if err != nil {
		return err
	}
	if started {
		return refuse(http.StatusConflict, "this payout has started paying; its schedule cannot be prepared again or rejected (pause or cancel it instead)")
	}
	return nil
}

// Prepare asks the engine to snapshot the holders and lock a schedule; on
// a locked or approved payout it re-prepares (approvals are dropped).
func (s *Service) Prepare(id uint64, admin string) (*ProceedPayout, error) {
	p, err := s.load(id)
	if err != nil {
		return nil, err
	}
	updates := map[string]interface{}{
		"status": StatusPrepareRequested, "prepared_by": admin, "preparation_requested_at": now(),
		"approvals_required": approvalsRequired(), "note": "", "updated_at": now(),
	}
	if !common.IsHexAddress(p.PayoutContractAddress) {
		a, err := s.asset(p.TokenizedAssetID)
		if err != nil {
			return nil, err
		}
		code, contract := s.payoutToken(p.PayoutAssetCode, a)
		if contract == "" {
			return nil, refuse(http.StatusUnprocessableEntity, "no tokenization currency with a token contract matches %s", p.PayoutAssetCode)
		}
		updates["payout_asset_code"], updates["payout_contract_address"] = code, contract
	}
	err = s.DB.Transaction(func(tx *gorm.DB) error {
		if err := s.refuseIfPaying(tx, id); err != nil {
			return err
		}
		if err := s.transition(tx, id, []string{StatusRegistered, StatusLocked, StatusApproved}, updates); err != nil {
			return err
		}
		return tx.Where("proceed_payout_id = ?", id).Delete(&PayoutApproval{}).Error
	})
	if err != nil {
		return nil, err
	}
	return s.after("prepare", id)
}

// FeeRequest sets a payout's processing fee.
type FeeRequest struct {
	FeeType  string `json:"feeType" binding:"required"` // FIXED or PERCENT
	FeeValue string `json:"feeValue" binding:"required"`
	FeeCap   string `json:"feeCap"` // PERCENT only; 0 or empty: no cap
}

// SetFee sets the payout's fee terms. On a registered payout they apply
// when it is prepared; on a locked or approved one the schedule is
// re-prepared with them, so they are approved with it.
func (s *Service) SetFee(id uint64, admin string, r FeeRequest) (*ProceedPayout, error) {
	typ := strings.ToUpper(strings.TrimSpace(r.FeeType))
	value, err := decimal.NewFromString(strings.TrimSpace(r.FeeValue))
	if err != nil || value.IsNegative() {
		return nil, refuse(http.StatusBadRequest, "feeValue must be a number of at least 0")
	}
	feeCap := decimal.Zero
	switch typ {
	case FeeFixed:
	case FeePercent:
		if value.GreaterThan(decimal.NewFromInt(100)) {
			return nil, refuse(http.StatusBadRequest, "a percent fee is at most 100")
		}
		if c := strings.TrimSpace(r.FeeCap); c != "" {
			if feeCap, err = decimal.NewFromString(c); err != nil || feeCap.IsNegative() {
				return nil, refuse(http.StatusBadRequest, "feeCap must be a number of at least 0")
			}
		}
	default:
		return nil, refuse(http.StatusBadRequest, "feeType is FIXED or PERCENT")
	}
	p, err := s.load(id)
	if err != nil {
		return nil, err
	}
	if total, err := decimal.NewFromString(p.TotalAmount); err == nil && typ == FeeFixed && value.GreaterThanOrEqual(total) {
		return nil, refuse(http.StatusBadRequest, "the fee must be less than the payout amount (%s %s)", p.TotalAmount, p.PayoutAssetCode)
	}
	updates := map[string]interface{}{"fee_type": typ, "fee_value": value.String(), "fee_cap": feeCap.String(), "fee_set_by": admin, "updated_at": now()}
	from := []string{StatusRegistered}
	if p.Status == StatusLocked || p.Status == StatusApproved {
		// the fee is part of the approved schedule: lock a new one
		from = []string{StatusLocked, StatusApproved}
		updates["status"], updates["prepared_by"], updates["preparation_requested_at"] = StatusPrepareRequested, admin, now()
		updates["approvals_required"], updates["note"] = approvalsRequired(), "Fee changed; the schedule is being prepared again."
	}
	err = s.DB.Transaction(func(tx *gorm.DB) error {
		if err := s.refuseIfPaying(tx, id); err != nil {
			return err
		}
		if err := s.transition(tx, id, from, updates); err != nil {
			return err
		}
		return tx.Where("proceed_payout_id = ?", id).Delete(&PayoutApproval{}).Error
	})
	if err != nil {
		return nil, err
	}
	return s.after("fee", id)
}

// Approve records an admin's approval of the locked schedule; the payout is
// approved once ApprovalsRequired distinct admins (not the one who prepared
// it) approved this schedule.
func (s *Service) Approve(id uint64, admin string) (*ProceedPayout, error) {
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		var p ProceedPayout
		if err := tx.First(&p, id).Error; err != nil {
			return notFound(err)
		}
		if p.Status != StatusLocked || p.ScheduleChecksum == "" {
			return refuse(http.StatusConflict, "only a locked schedule can be approved (the payout is %s)", p.Status)
		}
		if strings.EqualFold(strings.TrimSpace(p.PreparedBy), strings.TrimSpace(admin)) {
			return refuse(http.StatusForbidden, "the admin who prepared the schedule cannot approve it")
		}
		if strings.EqualFold(strings.TrimSpace(p.FeeSetBy), strings.TrimSpace(admin)) && p.FeeSetBy != "" {
			return refuse(http.StatusForbidden, "the admin who set the fee cannot approve the schedule")
		}
		a := PayoutApproval{CreatedAt: now(), ProceedPayoutID: id, AdminEmail: strings.ToLower(strings.TrimSpace(admin)), ScheduleChecksum: p.ScheduleChecksum}
		var mine int64
		if err := tx.Model(&PayoutApproval{}).Where("proceed_payout_id = ? AND admin_email = ?", id, a.AdminEmail).Count(&mine).Error; err != nil {
			return err
		}
		if mine > 0 {
			return refuse(http.StatusConflict, "you already approved this payout")
		}
		// the unique (payout, admin) index settles a concurrent double click
		res := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&a)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return refuse(http.StatusConflict, "you already approved this payout")
		}
		var n int64
		if err := tx.Model(&PayoutApproval{}).Where("proceed_payout_id = ? AND schedule_checksum = ?", id, p.ScheduleChecksum).Count(&n).Error; err != nil {
			return err
		}
		if n >= int64(p.ApprovalsRequired) && p.ApprovalsRequired > 0 {
			return tx.Model(&ProceedPayout{}).Where("id = ? AND status = ? AND schedule_checksum = ?", id, StatusLocked, p.ScheduleChecksum).
				Updates(map[string]interface{}{"status": StatusApproved, "approved_at": now(), "updated_at": now(),
					"note": fmt.Sprintf("Approved by %d admins; confirm the payout Safe is funded to start paying.", n)}).Error
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return s.after("approve", id)
}

// Reject sends a locked or approved schedule back to REGISTERED, dropping
// its approvals.
func (s *Service) Reject(id uint64, admin, reason string) (*ProceedPayout, error) {
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return nil, refuse(http.StatusBadRequest, "a reason is required")
	}
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		if err := s.refuseIfPaying(tx, id); err != nil {
			return err
		}
		if err := s.transition(tx, id, []string{StatusLocked, StatusApproved}, map[string]interface{}{
			"status": StatusRegistered, "payment_schedule_ready": 0, "approved_at": nil, "updated_at": now(),
			"note": trim(fmt.Sprintf("Schedule rejected by %s: %s", admin, reason), 1000),
		}); err != nil {
			return err
		}
		return tx.Where("proceed_payout_id = ?", id).Delete(&PayoutApproval{}).Error
	})
	if err != nil {
		return nil, err
	}
	return s.after("reject", id)
}

// ConfirmFunding asks the engine to check the payout Safe holds the payout
// (and the executor has gas) and, if so, to start paying.
func (s *Service) ConfirmFunding(id uint64, admin string) (*ProceedPayout, error) {
	if err := s.transition(s.DB, id, []string{StatusApproved}, map[string]interface{}{
		"status": StatusFundingCheckRequested, "funding_requested_by": admin, "note": "Checking the payout Safe's balance.", "updated_at": now(),
	}); err != nil {
		return nil, err
	}
	return s.after("fund", id)
}

// Pause stops a payout between batches (a batch already sent is settled).
func (s *Service) Pause(id uint64, admin, reason string) (*ProceedPayout, error) {
	p, err := s.load(id)
	if err != nil {
		return nil, err
	}
	if err := s.transition(s.DB, id, []string{StatusPaying, StatusFundingCheckRequested}, map[string]interface{}{
		"status": StatusPaused, "status_before_pause": p.Status, "updated_at": now(),
		"note": trim(fmt.Sprintf("Paused by %s: %s", admin, strings.TrimSpace(reason)), 1000),
	}); err != nil {
		return nil, err
	}
	return s.after("pause", id)
}

// Resume continues a paused payout. A payout paused while paying goes
// through the funding check again first.
func (s *Service) Resume(id uint64, admin string) (*ProceedPayout, error) {
	if err := s.transition(s.DB, id, []string{StatusPaused}, map[string]interface{}{
		"status": StatusFundingCheckRequested, "status_before_pause": "", "funding_requested_by": admin,
		"note": "Resumed; checking the payout Safe's balance.", "updated_at": now(),
	}); err != nil {
		return nil, err
	}
	return s.after("resume", id)
}

// Cancel stops a payout for good. Holders already paid stay paid.
func (s *Service) Cancel(id uint64, admin, reason string) (*ProceedPayout, error) {
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return nil, refuse(http.StatusBadRequest, "a reason is required")
	}
	var inFlight int64
	s.DB.Model(&PayoutBatch{}).Where("proceed_payout_id = ? AND status = ?", id, batchSubmitted).Count(&inFlight)
	if inFlight > 0 {
		return nil, refuse(http.StatusConflict, "a batch is being mined; pause the payout and cancel once it settles")
	}
	if err := s.transition(s.DB, id, []string{StatusRegistered, StatusPrepareRequested, StatusPreparing, StatusLocked, StatusApproved, StatusFundingCheckRequested, StatusPaused},
		map[string]interface{}{"status": StatusCancelled, "updated_at": now(), "note": trim(fmt.Sprintf("Cancelled by %s: %s", admin, reason), 1000)}); err != nil {
		return nil, err
	}
	return s.after("cancel", id)
}

// RetryFailed puts a payout's failed transfers back in the schedule and
// pays them after a new funding check.
func (s *Service) RetryFailed(id uint64, admin string) (*ProceedPayout, error) {
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		if err := s.transition(tx, id, []string{StatusCompletedWithFailures}, map[string]interface{}{
			"status": StatusFundingCheckRequested, "funding_requested_by": admin, "completed_at": nil, "payout_completed": 0,
			"note": "Retrying failed transfers; checking the payout Safe's balance.", "updated_at": now(),
		}); err != nil {
			return err
		}
		return tx.Model(&PayoutItem{}).Where("proceed_payout_id = ? AND status = ?", id, ItemFailed).
			Updates(map[string]interface{}{"status": ItemPending, "reason": "", "cannot_receive_asset": 0, "updated_at": now()}).Error
	})
	if err != nil {
		return nil, err
	}
	return s.after("retry", id)
}

// --- schedule items ---

// itemStates are the payout statuses in which admins may change items.
var itemStates = []string{StatusLocked, StatusApproved, StatusFundingCheckRequested, StatusPaying, StatusPaused, StatusCompletedWithFailures}

func (s *Service) item(id uint64, itemID string) (*ProceedPayout, *PayoutItem, error) {
	p, err := s.load(id)
	if err != nil {
		return nil, nil, err
	}
	var it PayoutItem
	if err := s.DB.Where("id = ? AND proceed_payout_id = ?", itemID, id).First(&it).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, refuse(http.StatusNotFound, "schedule item not found")
		}
		return nil, nil, err
	}
	if it.Kind != "" && it.Kind != KindHolder {
		return nil, nil, refuse(http.StatusBadRequest, "the fee and VAT lines follow the payout's fee; change the fee instead")
	}
	return p, &it, nil
}

func inStates(status string, states []string) bool {
	for _, s := range states {
		if s == status {
			return true
		}
	}
	return false
}

// ExcludeItem leaves a holder out of the payout (its share stays in the
// payout Safe). Approvals stand: an admin exclusion only pays less.
func (s *Service) ExcludeItem(id uint64, itemID, admin, reason string) (*PayoutItem, error) {
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return nil, refuse(http.StatusBadRequest, "a reason is required")
	}
	p, it, err := s.item(id, itemID)
	if err != nil {
		return nil, err
	}
	if !inStates(p.Status, itemStates) {
		return nil, refuse(http.StatusConflict, "schedule items cannot change while the payout is %s", p.Status)
	}
	res := s.DB.Model(&PayoutItem{}).Where("id = ? AND status IN ?", it.ID, []string{ItemPending, ItemFailed}).
		Updates(map[string]interface{}{"status": ItemExcluded, "reason": trim("Excluded by "+admin+": "+reason, 300), "action_by": admin, "updated_at": now()})
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, refuse(http.StatusConflict, "only a pending or failed holder can be excluded (this one is %s)", it.Status)
	}
	return s.afterItem(p, it.ID)
}

// IncludeItem puts back a holder an admin excluded (the engine's own
// exclusions, platform wallets, cannot be included).
func (s *Service) IncludeItem(id uint64, itemID, admin string) (*PayoutItem, error) {
	p, it, err := s.item(id, itemID)
	if err != nil {
		return nil, err
	}
	if !inStates(p.Status, itemStates) {
		return nil, refuse(http.StatusConflict, "schedule items cannot change while the payout is %s", p.Status)
	}
	if it.Status == ItemExcluded && it.ActionBy == "" {
		return nil, refuse(http.StatusBadRequest, "this address is excluded by the engine (%s)", it.Reason)
	}
	res := s.DB.Model(&PayoutItem{}).Where("id = ? AND status = ? AND action_by <> ''", it.ID, ItemExcluded).
		Updates(map[string]interface{}{"status": ItemPending, "reason": "", "action_by": admin, "updated_at": now()})
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, refuse(http.StatusConflict, "only an excluded holder can be included (this one is %s)", it.Status)
	}
	if p.Status == StatusCompletedWithFailures {
		// it is paid with the failed ones when they are retried
		s.DB.Model(&ProceedPayout{}).Where("id = ?", id).Update("note", "A holder was included again; use retry to pay it.")
	}
	return s.afterItem(p, it.ID)
}

// MarkPaid records that a holder the engine could not pay was paid another
// way (reference: how).
func (s *Service) MarkPaid(id uint64, itemID, admin, reference string) (*PayoutItem, error) {
	reference = strings.TrimSpace(reference)
	if reference == "" {
		return nil, refuse(http.StatusBadRequest, "a reference for the payment is required")
	}
	p, it, err := s.item(id, itemID)
	if err != nil {
		return nil, err
	}
	if p.Status != StatusPaused && p.Status != StatusCompletedWithFailures {
		return nil, refuse(http.StatusConflict, "pause the payout first (it is %s)", p.Status)
	}
	res := s.DB.Model(&PayoutItem{}).Where("id = ? AND status IN ?", it.ID, []string{ItemFailed, ItemPending}).
		Updates(map[string]interface{}{"status": ItemPaid, "reason": trim("Marked paid by "+admin+": "+reference, 300), "action_by": admin,
			"paid_at": now(), "cannot_receive_asset": 0, "updated_at": now()})
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, refuse(http.StatusConflict, "only a failed or pending holder can be marked paid (this one is %s)", it.Status)
	}
	return s.afterItem(p, it.ID)
}

func (s *Service) afterItem(p *ProceedPayout, itemID string) (*PayoutItem, error) {
	s.Recount(p.ID)
	s.completeIfSettled(p.ID)
	if _, err := s.after("item", p.ID); err != nil {
		return nil, err
	}
	var it PayoutItem
	if err := s.DB.First(&it, "id = ?", itemID).Error; err != nil {
		return nil, err
	}
	return &it, nil
}

// Recount refreshes a payout's holder counts and paid total (the engine
// does the same as it pays).
func (s *Service) Recount(id uint64) {
	type agg struct {
		Status string
		N      int
	}
	var rows []agg
	s.DB.Model(&PayoutItem{}).Select("status, COUNT(*) AS n").
		Where("proceed_payout_id = ? AND (kind = ? OR kind = '' OR kind IS NULL)", id, KindHolder).Group("status").Scan(&rows)
	counts := map[string]int{}
	for _, r := range rows {
		counts[r.Status] = r.N
	}
	var paid []PayoutItem
	s.DB.Select("amount_units").Where("proceed_payout_id = ? AND status = ? AND (kind = ? OR kind = '' OR kind IS NULL)", id, ItemPaid, KindHolder).Find(&paid)
	sum := new(big.Int)
	for _, it := range paid {
		sum.Add(sum, bigOf(it.AmountUnits))
	}
	s.DB.Model(&ProceedPayout{}).Where("id = ?", id).Updates(map[string]interface{}{
		"paid_count": counts[ItemPaid], "failed_count": counts[ItemFailed], "excluded_count": counts[ItemExcluded], "paid_units": sum.String(),
	})
}

// completeIfSettled completes a payout with failures once nothing is left
// failed or to pay.
func (s *Service) completeIfSettled(id uint64) {
	var open int64
	s.DB.Model(&PayoutItem{}).Where("proceed_payout_id = ? AND status IN ?", id, []string{ItemFailed, ItemPending, ItemQueued}).Count(&open)
	if open == 0 {
		s.DB.Model(&ProceedPayout{}).Where("id = ? AND status = ?", id, StatusCompletedWithFailures).
			Updates(map[string]interface{}{"status": StatusCompleted, "note": "All holders paid or settled by an admin.", "updated_at": now()})
	}
}

// --- engine controls ---

// EngineStatus is the engine's state and whether it is alive.
type EngineStatus struct {
	EngineState
	Online bool `json:"online"` // heartbeat within the last 90 seconds
}

func (s *Service) Engine() (*EngineStatus, error) {
	var st EngineState
	if err := s.DB.First(&st, 1).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &EngineStatus{EngineState: st, Online: st.HeartbeatAt != nil && time.Since(*st.HeartbeatAt) < 90*time.Second}, nil
}

func (s *Service) ensureEngineRow() {
	s.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&EngineState{ID: 1})
}

// Halt is the kill switch: the engine stops all work (a batch already sent
// is still settled when it resumes).
func (s *Service) Halt(admin, reason string) (*EngineStatus, error) {
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return nil, refuse(http.StatusBadRequest, "a reason is required")
	}
	s.ensureEngineRow()
	t := now()
	if err := s.DB.Model(&EngineState{}).Where("id = 1").Updates(map[string]interface{}{
		"halted": true, "halt_reason": trim(reason, 500), "halted_by": admin, "halted_at": &t,
	}).Error; err != nil {
		return nil, err
	}
	s.wake("halt", 0)
	return s.Engine()
}

// Unhalt lets the engine work again.
func (s *Service) Unhalt(admin string) (*EngineStatus, error) {
	s.ensureEngineRow()
	if err := s.DB.Model(&EngineState{}).Where("id = 1").Updates(map[string]interface{}{
		"halted": false, "halt_reason": "", "halted_by": admin, "halted_at": nil,
	}).Error; err != nil {
		return nil, err
	}
	s.wake("unhalt", 0)
	return s.Engine()
}

// Sweep asks the engine to move the payout Safe's whole balance of token
// to its configured sweep address (it refuses while a payout of that token
// is being funded or paid).
func (s *Service) Sweep(admin, token string) (*EngineStatus, error) {
	if !common.IsHexAddress(strings.TrimSpace(token)) {
		return nil, refuse(http.StatusBadRequest, "token must be a token contract address")
	}
	s.ensureEngineRow()
	t := now()
	res := s.DB.Model(&EngineState{}).Where("id = 1 AND (sweep_token = '' OR sweep_token IS NULL)").Updates(map[string]interface{}{
		"sweep_token": common.HexToAddress(token).Hex(), "sweep_requested_by": admin, "sweep_requested_at": &t, "sweep_result": "",
	})
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, refuse(http.StatusConflict, "a sweep is already pending")
	}
	s.wake("sweep", 0)
	return s.Engine()
}

// --- fee configuration ---

// FeeConfig is the payout fee wallet (set from TM) and the default fee
// terms of new payouts; VatWallet is the VAT service fee's wallet, shown
// for reference (VAT_WALLET on app-backend and the engine overrides it).
type FeeConfig struct {
	FeeWallet string `json:"feeWallet"`
	FeeDefaults
	VatWallet     string `json:"vatWallet"`
	LastUpdatedBy string `json:"lastUpdatedBy"`
}

func (s *Service) FeeConfig() (*FeeConfig, error) {
	c := &FeeConfig{FeeDefaults: s.feeDefaults()}
	var row ServiceFee
	if err := s.DB.Where("id = ?", FeeServiceID).First(&row).Error; err == nil {
		c.FeeWallet, c.LastUpdatedBy = row.FeeWalletSecretKey, row.LastUpdatedBy
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	var vat ServiceFee
	if s.DB.Where("id = ? AND inactive = 0", "VAT").First(&vat).Error == nil {
		c.VatWallet = vat.FeeWalletSecretKey
	}
	return c, nil
}

// FeeConfigRequest sets the fee wallet and the defaults.
type FeeConfigRequest struct {
	FeeWallet string `json:"feeWallet" binding:"required"`
	FeeType   string `json:"feeType"`
	FeeValue  string `json:"feeValue"`
	FeeCap    string `json:"feeCap"`
}

func (s *Service) SetFeeConfig(admin string, r FeeConfigRequest) (*FeeConfig, error) {
	if !common.IsHexAddress(strings.TrimSpace(r.FeeWallet)) {
		return nil, refuse(http.StatusBadRequest, "feeWallet must be an address")
	}
	value, feeCap := decimal.Zero, decimal.Zero
	var err error
	if v := strings.TrimSpace(r.FeeValue); v != "" {
		if value, err = decimal.NewFromString(v); err != nil || value.IsNegative() {
			return nil, refuse(http.StatusBadRequest, "feeValue must be a number of at least 0")
		}
	}
	row := ServiceFee{ID: FeeServiceID, FeeWalletSecretKey: common.HexToAddress(r.FeeWallet).Hex(), LastUpdatedBy: admin, UpdatedAt: now()}
	switch strings.ToUpper(strings.TrimSpace(r.FeeType)) {
	case "", FeeFixed:
		row.FeeFixed, _ = value.Float64()
	case FeePercent:
		if value.GreaterThan(decimal.NewFromInt(100)) {
			return nil, refuse(http.StatusBadRequest, "a percent fee is at most 100")
		}
		if c := strings.TrimSpace(r.FeeCap); c != "" {
			if feeCap, err = decimal.NewFromString(c); err != nil || feeCap.IsNegative() {
				return nil, refuse(http.StatusBadRequest, "feeCap must be a number of at least 0")
			}
		}
		row.FeePercent, _ = value.Float64()
		row.Remarks = feeCap.String()
	default:
		return nil, refuse(http.StatusBadRequest, "feeType is FIXED or PERCENT")
	}
	err = s.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"fee_wallet_secret_key", "fee_fixed", "fee_percent", "remarks", "last_updated_by", "updated_at", "inactive"}),
	}).Create(&row).Error
	if err != nil {
		return nil, err
	}
	return s.FeeConfig()
}

// --- helpers ---

func bigOf(s string) *big.Int {
	n, ok := new(big.Int).SetString(strings.TrimSpace(s), 10)
	if !ok {
		return new(big.Int)
	}
	return n
}

// human is base units as a decimal string.
func human(units string, decimals int) string {
	return decimal.NewFromBigInt(bigOf(units), -int32(decimals)).String()
}

func trim(s string, n int) string {
	s = strings.TrimSpace(s)
	if len(s) > n {
		return s[:n]
	}
	return s
}

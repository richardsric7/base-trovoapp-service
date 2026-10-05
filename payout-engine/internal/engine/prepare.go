package engine

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"os"
	"sort"
	"strings"
	"time"

	"trovo-payout-engine/internal/chain"
	"trovo-payout-engine/internal/store"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// maxChunksPerTick bounds a tick's scanning, so the heartbeat, pauses and
// other payouts are not held up by a long first scan.
const maxChunksPerTick = 40

func key(a common.Address) string { return strings.ToLower(a.Hex()) }

// prepare builds a payout's holder schedule. The token's Transfer events are
// replayed into payout_token_balances (kept between payouts, so only new
// blocks are scanned); scanning follows the chain head, picking up transfers
// made while it runs, and once it is within CatchUpWithin blocks of the
// confirmed head the balances at that block are frozen into the schedule
// and the payout is LOCKED for approval.
func (e *Engine) prepare(ctx context.Context, p *store.ProceedPayout) error {
	if p.Status == store.StatusPrepareRequested {
		if err := e.startPreparation(ctx, p); err != nil {
			return err
		}
	}
	token := common.HexToAddress(p.TokenContractAddress)
	idx, err := e.tokenIndex(ctx, p, token)
	if err != nil {
		return err
	}
	for i := 0; i < maxChunksPerTick; i++ {
		head, err := e.Chain.BlockNumber(ctx)
		if err != nil {
			return err
		}
		if head < e.Cfg.Confirmations {
			return nil
		}
		confirmed := head - e.Cfg.Confirmations
		if confirmed <= idx.ScannedBlock+e.Cfg.CatchUpWithin && idx.ScannedBlock >= idx.StartBlock {
			return e.lockSchedule(ctx, p, idx)
		}
		from := idx.ScannedBlock + 1
		if idx.ScannedBlock < idx.StartBlock {
			from = idx.StartBlock
		}
		to := from + e.Cfg.LogChunk - 1
		if to > confirmed {
			to = confirmed
		}
		if err := e.applyTransfers(ctx, token, from, to); err != nil {
			return err
		}
		idx.ScannedBlock = to
		e.DB.Model(&store.ProceedPayout{}).Where("id = ?", p.ID).Updates(map[string]interface{}{"scanned_block": to, "snapshot_start_block": idx.StartBlock})
		e.activity(fmt.Sprintf("preparing payout %d: scanned to block %d of %d", p.ID, to, confirmed))
		if e.currentStatus(p.ID) != store.StatusPreparing {
			return nil // cancelled meanwhile
		}
	}
	e.publish(ctx, "preparing", p.ID, fmt.Sprintf("scanned to block %d", idx.ScannedBlock))
	return nil
}

// startPreparation fills in what the payout pays with, clears a previous
// schedule (re-preparation) and moves it to PREPARING.
func (e *Engine) startPreparation(ctx context.Context, p *store.ProceedPayout) error {
	var asset store.TokenizedAsset
	if err := e.DB.First(&asset, "id = ?", p.TokenizedAssetID).Error; err != nil {
		return fmt.Errorf("tokenized asset %s: %w", p.TokenizedAssetID, err)
	}
	if asset.ContractAddress == nil || !common.IsHexAddress(*asset.ContractAddress) {
		return errors.New("the asset has no token contract")
	}
	if !common.IsHexAddress(p.PayoutContractAddress) {
		return fmt.Errorf("payout currency %s has no token contract", p.PayoutAssetCode)
	}
	token := common.HexToAddress(*asset.ContractAddress)
	payoutToken := common.HexToAddress(p.PayoutContractAddress)
	tokenDecimals, err := chain.Decimals(ctx, e.Chain, token)
	if err != nil {
		return fmt.Errorf("reading the asset token's decimals: %w", err)
	}
	payoutDecimals, err := chain.Decimals(ctx, e.Chain, payoutToken)
	if err != nil {
		return fmt.Errorf("reading %s decimals: %w", p.PayoutAssetCode, err)
	}
	total, err := decimal.NewFromString(strings.TrimSpace(p.TotalAmount))
	if err != nil || !total.IsPositive() {
		return fmt.Errorf("invalid payout amount %q", p.TotalAmount)
	}
	totalUnits := total.Shift(int32(payoutDecimals)).Truncate(0)
	// a schedule that has started paying is never prepared again (that
	// would drop its paid lines and pay those holders twice): park it
	var sent, paid int64
	e.DB.Model(&store.PayoutBatch{}).Where("proceed_payout_id = ?", p.ID).Count(&sent)
	e.DB.Model(&store.PayoutItem{}).Where("proceed_payout_id = ? AND status IN ?", p.ID, []string{store.ItemPaid, store.ItemQueued}).Count(&paid)
	if sent > 0 || paid > 0 {
		_, err := e.setStatus(p.ID, store.StatusPrepareRequested, store.StatusPaused, map[string]interface{}{
			"note": "This payout has started paying, so its schedule cannot be prepared again; resume or cancel it.",
		})
		e.publish(ctx, "prepare-refused", p.ID, "the payout has started paying")
		return err
	}
	err = e.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("proceed_payout_id = ?", p.ID).Delete(&store.PayoutItem{}).Error; err != nil {
			return err
		}
		if err := tx.Where("proceed_payout_id = ?", p.ID).Delete(&store.PayoutApproval{}).Error; err != nil {
			return err
		}
		res := tx.Model(&store.ProceedPayout{}).Where("id = ? AND status = ?", p.ID, store.StatusPrepareRequested).Updates(map[string]interface{}{
			"status": store.StatusPreparing, "token_contract_address": token.Hex(), "token_decimals": tokenDecimals,
			"payout_decimals": payoutDecimals, "payout_safe_address": e.Cfg.PayoutSafe.Hex(), "total_units": totalUnits.String(),
			"payment_schedule_ready": 0, "schedule_checksum": "", "snapshot_block": 0, "locked_at": nil, "approved_at": nil, "note": "",
		})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errors.New("the payout changed state meanwhile")
		}
		return nil
	})
	if err != nil {
		return err
	}
	p.Status, p.TokenContractAddress, p.TokenDecimals, p.PayoutDecimals = store.StatusPreparing, token.Hex(), tokenDecimals, payoutDecimals
	p.PayoutSafeAddress, p.TotalUnits = e.Cfg.PayoutSafe.Hex(), totalUnits.String()
	e.publish(ctx, "preparing", p.ID, "schedule preparation started")
	return nil
}

// tokenIndex returns the token's replay state, starting a new one at the
// block of the asset's creation (the token cannot be older than its asset).
func (e *Engine) tokenIndex(ctx context.Context, p *store.ProceedPayout, token common.Address) (*store.TokenIndex, error) {
	var idx store.TokenIndex
	err := e.DB.First(&idx, "token = ?", key(token)).Error
	if err == nil {
		return &idx, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	var asset store.TokenizedAsset
	if err := e.DB.First(&asset, "id = ?", p.TokenizedAssetID).Error; err != nil {
		return nil, err
	}
	start, err := chain.BlockAtTime(ctx, e.Chain, asset.CreatedAt.Add(-time.Hour))
	if err != nil {
		return nil, fmt.Errorf("finding the asset's starting block: %w", err)
	}
	if e.Cfg.DefaultStartBlock > start {
		start = e.Cfg.DefaultStartBlock
	}
	idx = store.TokenIndex{Token: key(token), StartBlock: start, UpdatedAt: time.Now().UTC()}
	if start > 0 {
		idx.ScannedBlock = start - 1
	}
	if err := e.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&idx).Error; err != nil {
		return nil, err
	}
	return &idx, nil
}

// applyTransfers replays token's transfers in [from, to] onto the stored
// balances, atomically with the index's progress.
func (e *Engine) applyTransfers(ctx context.Context, token common.Address, from, to uint64) error {
	transfers, err := chain.Transfers(ctx, e.Chain, token, from, to)
	if err != nil {
		return fmt.Errorf("reading transfers %d-%d: %w", from, to, err)
	}
	deltas := map[string]*big.Int{}
	add := func(a common.Address, v *big.Int) {
		if a == (common.Address{}) { // mints come from, burns go to, the zero address
			return
		}
		k := key(a)
		if deltas[k] == nil {
			deltas[k] = new(big.Int)
		}
		deltas[k].Add(deltas[k], v)
	}
	for _, t := range transfers {
		add(t.From, new(big.Int).Neg(t.Amount))
		add(t.To, t.Amount)
	}
	tk := key(token)
	return e.DB.Transaction(func(tx *gorm.DB) error {
		if len(deltas) > 0 {
			holders := make([]string, 0, len(deltas))
			for h := range deltas {
				holders = append(holders, h)
			}
			var rows []store.TokenBalance
			for i := 0; i < len(holders); i += 500 {
				end := i + 500
				if end > len(holders) {
					end = len(holders)
				}
				var part []store.TokenBalance
				if err := tx.Where("token = ? AND holder IN ?", tk, holders[i:end]).Find(&part).Error; err != nil {
					return err
				}
				rows = append(rows, part...)
			}
			current := map[string]*big.Int{}
			for _, r := range rows {
				current[r.Holder] = bigOf(r.Balance)
			}
			updated := make([]store.TokenBalance, 0, len(deltas))
			for h, d := range deltas {
				b := new(big.Int).Add(orZero(current[h]), d)
				if b.Sign() < 0 {
					return fmt.Errorf("replayed balance of %s went negative at blocks %d-%d (the scan must start before the token's first transfer: set PROCEED_PAYOUT_DEFAULT_START_BLOCK or reset its payout_token_indexes row)", h, from, to)
				}
				updated = append(updated, store.TokenBalance{Token: tk, Holder: h, Balance: b.String()})
			}
			if err := tx.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "token"}, {Name: "holder"}}, DoUpdates: clause.AssignmentColumns([]string{"balance"})}).
				CreateInBatches(updated, 500).Error; err != nil {
				return err
			}
		}
		return tx.Model(&store.TokenIndex{}).Where("token = ?", tk).Updates(map[string]interface{}{"scanned_block": to, "updated_at": time.Now().UTC()}).Error
	})
}

func orZero(b *big.Int) *big.Int {
	if b == nil {
		return new(big.Int)
	}
	return b
}

// lockSchedule freezes the balances as of idx.ScannedBlock into the payout's
// schedule.
func (e *Engine) lockSchedule(ctx context.Context, p *store.ProceedPayout, idx *store.TokenIndex) error {
	var rows []store.TokenBalance
	if err := e.DB.Where("token = ? AND balance <> ?", idx.Token, "0").Find(&rows).Error; err != nil {
		return err
	}
	balances := map[string]*big.Int{}
	supply := new(big.Int)
	for _, r := range rows {
		b := bigOf(r.Balance)
		if b.Sign() > 0 {
			balances[r.Holder] = b
			supply.Add(supply, b)
		}
	}
	if supply.Sign() == 0 {
		return errors.New("the token has no holders")
	}
	excluded, err := e.platformWallets(p)
	if err != nil {
		return err
	}
	if err := e.creditOfferBookSellers(balances, idx.Token, excluded); err != nil {
		return err
	}

	fee, err := e.payoutFee(p)
	if err != nil {
		return err
	}
	total := fee.distributable
	holders := make([]string, 0, len(balances))
	for h := range balances {
		holders = append(holders, h)
	}
	sort.Strings(holders)
	usernames, err := e.usernames(holders)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	items := make([]store.PayoutItem, 0, len(holders))
	payable, eligible := new(big.Int), new(big.Int)
	excludedCount := 0
	for _, h := range holders {
		bal := balances[h]
		amount := new(big.Int).Div(new(big.Int).Mul(bal, total), supply)
		it := store.PayoutItem{
			ID: uuid.NewString(), CreatedAt: now, UpdatedAt: now, TokenizedAssetID: p.TokenizedAssetID, Batch: p.Batch,
			PayoutAssetCode: p.PayoutAssetCode, PayoutContractAddress: p.PayoutContractAddress,
			BeneficiaryAddress: common.HexToAddress(h).Hex(), ProceedPayoutID: p.ID, BalanceUnits: bal.String(), AmountUnits: amount.String(),
			ConfirmedTokenizedAssetBalance: human(bal, p.TokenDecimals), AmountToReceive: human(amount, p.PayoutDecimals),
			Status: store.ItemPending, Username: usernames[h], Kind: store.KindHolder,
		}
		switch reason, isExcluded := excluded[h]; {
		case isExcluded:
			it.Status, it.Reason = store.ItemExcluded, reason
			excludedCount++
		case amount.Cmp(e.Cfg.MinPayoutUnits) < 0:
			it.Status, it.Reason = store.ItemSkipped, "share below the minimum payout"
		default:
			payable.Add(payable, amount)
			eligible.Add(eligible, bal)
		}
		items = append(items, it)
	}
	// the fee and its VAT, paid with the holders
	for _, f := range []struct {
		kind, wallet, reason string
		amount               *big.Int
	}{{store.KindFee, fee.feeWallet, "payout processing fee", fee.fee}, {store.KindVat, fee.vatWallet, "VAT on the payout processing fee", fee.vat}} {
		if f.amount.Sign() == 0 {
			continue
		}
		if _, isHolder := balances[key(common.HexToAddress(f.wallet))]; isHolder {
			return fmt.Errorf("the %s wallet %s holds the asset token; use a wallet that does not", strings.ToLower(f.kind), f.wallet)
		}
		items = append(items, store.PayoutItem{
			ID: uuid.NewString(), CreatedAt: now, UpdatedAt: now, TokenizedAssetID: p.TokenizedAssetID, Batch: p.Batch,
			PayoutAssetCode: p.PayoutAssetCode, PayoutContractAddress: p.PayoutContractAddress, BeneficiaryAddress: f.wallet,
			ProceedPayoutID: p.ID, BalanceUnits: "0", AmountUnits: f.amount.String(), AmountToReceive: human(f.amount, p.PayoutDecimals),
			Status: store.ItemPending, Reason: f.reason, Kind: f.kind,
		})
	}
	p.FeeUnits, p.VatUnits = fee.fee.String(), fee.vat.String()
	checksum := checksumOf(p, supply.String(), idx.ScannedBlock, items)
	perToken := decimal.NewFromBigInt(total, -int32(p.PayoutDecimals)).Div(decimal.NewFromBigInt(supply, -int32(p.TokenDecimals)))

	err = e.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("proceed_payout_id = ?", p.ID).Delete(&store.PayoutItem{}).Error; err != nil {
			return err
		}
		if err := tx.CreateInBatches(items, 500).Error; err != nil {
			return err
		}
		res := tx.Model(&store.ProceedPayout{}).Where("id = ? AND status = ?", p.ID, store.StatusPreparing).Updates(map[string]interface{}{
			"status": store.StatusLocked, "snapshot_block": idx.ScannedBlock, "scanned_block": idx.ScannedBlock, "supply_units": supply.String(),
			"eligible_units": eligible.String(), "payable_units": payable.String(), "retained_units": new(big.Int).Sub(total, payable).String(), // excluded holders' share and rounding dust
			"paid_units": "0", "holder_count": len(items), "excluded_count": excludedCount, "paid_count": 0, "failed_count": 0,
			"amount_per_token": perToken.Round(18).String(),
			"fee_units":        fee.fee.String(), "vat_units": fee.vat.String(), "vat_percent": fee.vatPercent, "fee_wallet": fee.feeWallet,
			"vat_wallet": fee.vatWallet, "holder_payable": total.String(),
			"schedule_checksum": checksum, "payment_schedule_ready": 1, "locked_at": now,
			"note": fmt.Sprintf("%d holders at block %d; %d excluded", len(items), idx.ScannedBlock, excludedCount),
		})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errors.New("the payout changed state meanwhile")
		}
		return nil
	})
	if err != nil {
		return err
	}
	e.activity(fmt.Sprintf("payout %d locked for approval", p.ID))
	e.publish(ctx, "locked", p.ID, fmt.Sprintf("%d holders at block %d", len(items), idx.ScannedBlock))
	return nil
}

// ScheduleChecksum recomputes a payout's checksum from its stored items, to
// detect a schedule (or fee) changed after it was approved.
func (e *Engine) ScheduleChecksum(p *store.ProceedPayout) (string, error) {
	var items []store.PayoutItem
	if err := e.DB.Where("proceed_payout_id = ?", p.ID).Find(&items).Error; err != nil {
		return "", err
	}
	return checksumOf(p, p.SupplyUnits, p.SnapshotBlock, items), nil
}

// checksumOf digests what approvers approve: the amount, fee terms, the
// snapshot and every item's balance, amount and (locked) status. Statuses
// items reach afterwards (paid, failed, excluded or marked paid by an admin)
// count as they were when locked.
func checksumOf(p *store.ProceedPayout, supply string, block uint64, items []store.PayoutItem) string {
	sorted := append([]store.PayoutItem(nil), items...)
	kind := func(it store.PayoutItem) string {
		if it.Kind == "" {
			return store.KindHolder
		}
		return it.Kind
	}
	sort.Slice(sorted, func(i, j int) bool {
		if kind(sorted[i]) != kind(sorted[j]) {
			return kind(sorted[i]) < kind(sorted[j])
		}
		return strings.ToLower(sorted[i].BeneficiaryAddress) < strings.ToLower(sorted[j].BeneficiaryAddress)
	})
	var digest strings.Builder
	fmt.Fprintf(&digest, "%d|%s|%s|%d|%s|%s|%s|%s|%s\n", p.ID, p.TotalUnits, supply, block, p.FeeType, p.FeeValue, p.FeeCap, p.FeeUnits, p.VatUnits)
	for _, it := range sorted {
		status := it.Status
		switch {
		case status == store.ItemExcluded && it.ActionBy == "", status == store.ItemSkipped:
		default:
			status = store.ItemPending
		}
		fmt.Fprintf(&digest, "%s|%s|%s|%s|%s\n", kind(it), strings.ToLower(it.BeneficiaryAddress), it.BalanceUnits, it.AmountUnits, status)
	}
	sum := sha256.Sum256([]byte(digest.String()))
	return hex.EncodeToString(sum[:])
}

// platformWallets are the addresses never paid, with why.
func (e *Engine) platformWallets(p *store.ProceedPayout) (map[string]string, error) {
	out := map[string]string{key(e.Cfg.PayoutSafe): "the payout Safe"}
	for _, a := range e.Cfg.ExcludedAddresses {
		out[key(a)] = "excluded by configuration"
	}
	var asset store.TokenizedAsset
	if err := e.DB.First(&asset, "id = ?", p.TokenizedAssetID).Error; err != nil {
		return nil, err
	}
	if asset.MarketMakingWallet != nil && common.IsHexAddress(*asset.MarketMakingWallet) {
		out[key(common.HexToAddress(*asset.MarketMakingWallet))] = "the asset's market-making wallet"
	}
	var profile store.User
	if err := e.DB.Where("username = ?", e.Cfg.IssuingProfile).First(&profile).Error; err == nil {
		var wallets []store.UserWallet
		if err := e.DB.Where("user_id = ?", profile.ID).Find(&wallets).Error; err != nil {
			return nil, err
		}
		for _, w := range wallets {
			if common.IsHexAddress(w.ID) {
				out[key(common.HexToAddress(w.ID))] = "a tokenization issuing / distribution wallet"
			}
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if e.Cfg.OfferBook != (common.Address{}) {
		out[key(e.Cfg.OfferBook)] = "the offer book (its tokens are credited to their sellers)"
	}
	return out, nil
}

// creditOfferBookSellers moves the tokens the offer book holds for open sell
// offers back to their sellers, pro rata to what each offer still sells
// (app-backend's offer book index), so tokens on sale still earn.
func (e *Engine) creditOfferBookSellers(balances map[string]*big.Int, token string, excluded map[string]string) error {
	if e.Cfg.OfferBook == (common.Address{}) {
		return nil
	}
	book := key(e.Cfg.OfferBook)
	held := balances[book]
	if held == nil || held.Sign() == 0 {
		return nil
	}
	var offers []store.OfferBookOffer
	if err := e.DB.Where("open = ? AND LOWER(sell_token) = ?", true, token).Find(&offers).Error; err != nil {
		return nil // no offer book index here: leave the book's tokens unpaid
	}
	perSeller := map[string]*big.Int{}
	total := new(big.Int)
	for _, o := range offers {
		r := bigOf(o.Remaining)
		if r.Sign() <= 0 || !common.IsHexAddress(o.Seller) {
			continue
		}
		s := key(common.HexToAddress(o.Seller))
		if perSeller[s] == nil {
			perSeller[s] = new(big.Int)
		}
		perSeller[s].Add(perSeller[s], r)
		total.Add(total, r)
	}
	if total.Sign() == 0 {
		return nil
	}
	credited := new(big.Int)
	for s, r := range perSeller {
		share := new(big.Int).Div(new(big.Int).Mul(held, r), total)
		if balances[s] == nil {
			balances[s] = new(big.Int)
		}
		balances[s].Add(balances[s], share)
		credited.Add(credited, share)
	}
	balances[book] = new(big.Int).Sub(held, credited) // rounding dust stays with the (excluded) book
	return nil
}

// usernames maps holder addresses to the Trovo users owning them.
func (e *Engine) usernames(holders []string) (map[string]string, error) {
	out := map[string]string{}
	for i := 0; i < len(holders); i += 500 {
		end := i + 500
		if end > len(holders) {
			end = len(holders)
		}
		type row struct {
			Wallet   string
			Username string
		}
		var rows []row
		if err := e.DB.Table("user_wallets").Select("LOWER(user_wallets.id) AS wallet, users.username AS username").
			Joins("JOIN users ON users.id = user_wallets.user_id").
			Where("LOWER(user_wallets.id) IN ?", holders[i:end]).Scan(&rows).Error; err != nil {
			return nil, err
		}
		for _, r := range rows {
			out[r.Wallet] = r.Username
		}
	}
	return out, nil
}

func human(units *big.Int, decimals int) float64 {
	f, _ := decimal.NewFromBigInt(units, -int32(decimals)).Float64()
	return f
}

// payoutFee is a payout's processing fee and the VAT on it.
type payoutFee struct {
	fee, vat, distributable *big.Int
	feeWallet, vatWallet    string
	vatPercent              string
}

// payoutFee computes the payout's processing fee (FIXED: FeeValue of the
// payout token; PERCENT: FeeValue % of the total, at most FeeCap when it is
// above 0) and the VAT on it at the asset country's rate (country_configs),
// and checks their wallets are configured (the PROCEED_PAYOUT_FEE service
// fee's wallet, set from TM; VAT_WALLET or the VAT service fee's wallet).
func (e *Engine) payoutFee(p *store.ProceedPayout) (*payoutFee, error) {
	total := bigOf(p.TotalUnits)
	f := &payoutFee{fee: new(big.Int), vat: new(big.Int), vatPercent: "0"}
	value, err := decimal.NewFromString(nonEmpty(p.FeeValue, "0"))
	if err != nil || value.IsNegative() {
		return nil, fmt.Errorf("invalid payout fee %q", p.FeeValue)
	}
	feeCap, err := decimal.NewFromString(nonEmpty(p.FeeCap, "0"))
	if err != nil || feeCap.IsNegative() {
		return nil, fmt.Errorf("invalid payout fee cap %q", p.FeeCap)
	}
	switch strings.ToUpper(nonEmpty(p.FeeType, store.FeeFixed)) {
	case store.FeeFixed:
		f.fee = value.Shift(int32(p.PayoutDecimals)).Truncate(0).BigInt()
	case store.FeePercent:
		f.fee = decimal.NewFromBigInt(total, 0).Mul(value).Div(decimal.NewFromInt(100)).Truncate(0).BigInt()
		if c := feeCap.Shift(int32(p.PayoutDecimals)).Truncate(0).BigInt(); c.Sign() > 0 && f.fee.Cmp(c) > 0 {
			f.fee = c
		}
	default:
		return nil, fmt.Errorf("invalid payout fee type %q", p.FeeType)
	}
	if f.fee.Sign() > 0 {
		var asset store.TokenizedAsset
		e.DB.First(&asset, "id = ?", p.TokenizedAssetID)
		country := "NG"
		if asset.AssetCountryLocation != nil && strings.TrimSpace(*asset.AssetCountryLocation) != "" {
			country = strings.ToUpper(strings.TrimSpace(*asset.AssetCountryLocation))
		}
		var cc store.CountryConfig
		e.DB.Where("country_code = ?", country).First(&cc)
		vatPct := decimal.NewFromFloat(cc.VATPercent)
		f.vatPercent = vatPct.String()
		f.vat = decimal.NewFromBigInt(f.fee, 0).Mul(vatPct).Div(decimal.NewFromInt(100)).Truncate(0).BigInt()
		if f.feeWallet, err = e.serviceFeeWallet(store.FeeServiceID, ""); err != nil {
			return nil, fmt.Errorf("the payout fee wallet: %w (set it in TM)", err)
		}
		if f.vat.Sign() > 0 {
			if f.vatWallet, err = e.serviceFeeWallet("VAT", os.Getenv("VAT_WALLET")); err != nil {
				return nil, fmt.Errorf("the VAT wallet: %w", err)
			}
			if strings.EqualFold(f.vatWallet, f.feeWallet) {
				return nil, errors.New("the payout fee wallet and the VAT wallet must differ")
			}
		}
	}
	f.distributable = new(big.Int).Sub(total, new(big.Int).Add(f.fee, f.vat))
	if f.distributable.Sign() <= 0 {
		return nil, fmt.Errorf("the fee (%s) and VAT (%s) leave nothing to distribute", units(f.fee, p.PayoutDecimals), units(f.vat, p.PayoutDecimals))
	}
	return f, nil
}

// serviceFeeWallet is the address of a service_fees row's wallet (override
// first, when set).
func (e *Engine) serviceFeeWallet(id, override string) (string, error) {
	addr := strings.TrimSpace(override)
	if addr == "" {
		var row store.ServiceFee
		if err := e.DB.Where("id = ? AND inactive = 0", id).First(&row).Error; err != nil {
			return "", errors.New("not configured")
		}
		addr = strings.TrimSpace(row.FeeWalletSecretKey)
	}
	if !common.IsHexAddress(addr) {
		return "", errors.New("not configured as an address")
	}
	return common.HexToAddress(addr).Hex(), nil
}

func nonEmpty(s, def string) string {
	if strings.TrimSpace(s) == "" {
		return def
	}
	return strings.TrimSpace(s)
}

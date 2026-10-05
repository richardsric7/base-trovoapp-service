package engine

import (
	"context"
	"fmt"
	"math/big"

	"trovo-payout-engine/internal/chain"
	"trovo-payout-engine/internal/safe"
	"trovo-payout-engine/internal/store"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

// gas a transfer adds to a batch, and a batch's own overhead (estimates for
// the funding check; batches themselves are estimated exactly)
const (
	gasPerTransfer = 45_000
	gasPerBatch    = 90_000
)

// checkFunding confirms an approved payout can be paid: its schedule is the
// one approved, the configured signers can sign for the payout Safe, the
// Safe holds what is still to be paid (on top of what other payouts of the
// same token still owe) and the executing signer has ETH for the gas. It
// then starts paying; otherwise the payout goes back to APPROVED with why.
func (e *Engine) checkFunding(ctx context.Context, p *store.ProceedPayout) error {
	fail := func(reason string) error {
		_, err := e.setStatus(p.ID, store.StatusFundingCheckRequested, store.StatusApproved, map[string]interface{}{"note": "Funding not confirmed: " + reason})
		e.publish(ctx, "funding-failed", p.ID, reason)
		return err
	}
	sum, err := e.ScheduleChecksum(p)
	if err != nil {
		return err
	}
	if sum != p.ScheduleChecksum {
		return fail("the schedule changed after it was locked; re-prepare it")
	}
	var approvals int64
	if err := e.DB.Model(&store.PayoutApproval{}).Where("proceed_payout_id = ? AND schedule_checksum = ?", p.ID, p.ScheduleChecksum).Count(&approvals).Error; err != nil {
		return err
	}
	if approvals < int64(p.ApprovalsRequired) || p.ApprovalsRequired < 1 {
		return fail(fmt.Sprintf("%d of %d approvals", approvals, p.ApprovalsRequired))
	}
	if !common.IsHexAddress(p.PayoutSafeAddress) || common.HexToAddress(p.PayoutSafeAddress) != e.Cfg.PayoutSafe {
		return fail("the schedule was prepared for another payout Safe; re-prepare it")
	}
	signers, err := e.Signers.Signers(ctx)
	if err != nil {
		return fail(err.Error())
	}
	if err := safe.CheckSigners(ctx, e.Chain, e.Cfg.PayoutSafe, signers); err != nil {
		return fail(err.Error())
	}

	token := common.HexToAddress(p.PayoutContractAddress)
	owed, pendingCount, err := e.unpaid(p.ID)
	if err != nil {
		return err
	}
	others, err := e.owedByOthers(p)
	if err != nil {
		return err
	}
	need := new(big.Int).Add(owed, others)
	have, err := chain.BalanceOf(ctx, e.Chain, token, e.Cfg.PayoutSafe, nil)
	if err != nil {
		return err
	}
	if have.Cmp(need) < 0 {
		short := new(big.Int).Sub(need, have)
		return fail(fmt.Sprintf("the payout Safe %s holds %s %s; %s %s is needed (%s for this payout, %s for others in progress) - send %s more",
			e.Cfg.PayoutSafe.Hex(), units(have, p.PayoutDecimals), p.PayoutAssetCode, units(need, p.PayoutDecimals), p.PayoutAssetCode,
			units(owed, p.PayoutDecimals), units(others, p.PayoutDecimals), units(short, p.PayoutDecimals)))
	}

	gasCost, err := e.gasCost(ctx, pendingCount)
	if err != nil {
		return err
	}
	ethHave, err := e.Chain.BalanceAt(ctx, signers[0].Address, nil)
	if err != nil {
		return err
	}
	if ethHave.Cmp(gasCost) < 0 {
		return fail(fmt.Sprintf("the executing signer %s holds %s ETH; about %s ETH is needed for gas", signers[0].Address.Hex(), units(ethHave, 18), units(gasCost, 18)))
	}
	ok, err := e.setStatus(p.ID, store.StatusFundingCheckRequested, store.StatusPaying, map[string]interface{}{
		"funding_checked_at": now(), "started_at": now(),
		"note": fmt.Sprintf("Funding confirmed: the Safe holds %s %s (needed %s); gas about %s ETH", units(have, p.PayoutDecimals), p.PayoutAssetCode, units(need, p.PayoutDecimals), units(gasCost, 18)),
	})
	if err != nil {
		return err
	}
	if ok {
		e.publish(ctx, "funded", p.ID, "funding confirmed; paying")
	}
	return nil
}

// unpaid is what a payout still has to pay (pending or queued), and how many
// transfers that is.
func (e *Engine) unpaid(id uint64) (*big.Int, int, error) {
	var items []store.PayoutItem
	if err := e.DB.Select("amount_units").Where("proceed_payout_id = ? AND status IN ?", id, []string{store.ItemPending, store.ItemQueued}).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	sum := new(big.Int)
	for _, it := range items {
		sum.Add(sum, bigOf(it.AmountUnits))
	}
	return sum, len(items), nil
}

// owedByOthers is what other payouts of the same token, being paid or
// paused, still owe from the same Safe.
func (e *Engine) owedByOthers(p *store.ProceedPayout) (*big.Int, error) {
	var others []store.ProceedPayout
	if err := e.DB.Where("id <> ? AND LOWER(payout_contract_address) = LOWER(?) AND status IN ?", p.ID, p.PayoutContractAddress,
		[]string{store.StatusPaying, store.StatusPaused}).Find(&others).Error; err != nil {
		return nil, err
	}
	sum := new(big.Int)
	for _, o := range others {
		owed, _, err := e.unpaid(o.ID)
		if err != nil {
			return nil, err
		}
		sum.Add(sum, owed)
	}
	return sum, nil
}

// gasCost estimates the ETH paying n transfers takes, with a 50% margin.
func (e *Engine) gasCost(ctx context.Context, n int) (*big.Int, error) {
	if n == 0 {
		return new(big.Int), nil
	}
	batches := (n + e.Cfg.BatchSize - 1) / e.Cfg.BatchSize
	gas := big.NewInt(int64(n*gasPerTransfer + batches*gasPerBatch))
	head, err := e.Chain.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, err
	}
	tip, err := e.Chain.SuggestGasTipCap(ctx)
	if err != nil {
		tip = big.NewInt(1_000_000)
	}
	price := new(big.Int).Set(tip)
	if head.BaseFee != nil {
		price.Add(price, new(big.Int).Mul(head.BaseFee, big.NewInt(2)))
	}
	cost := new(big.Int).Mul(gas, price)
	return cost.Add(cost, new(big.Int).Div(cost, big.NewInt(2))), nil
}

func units(v *big.Int, decimals int) string {
	if v == nil {
		return "0"
	}
	return decimal.NewFromBigInt(v, -int32(decimals)).String()
}

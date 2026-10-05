package engine

import (
	"context"
	"fmt"
	"strings"

	"trovo-payout-engine/internal/chain"
	"trovo-payout-engine/internal/safe"
	"trovo-payout-engine/internal/store"

	"github.com/ethereum/go-ethereum/common"
)

// sweep carries out an admin's request (tm-api) to move the payout Safe's
// whole balance of a token to PROCEED_PAYOUT_SWEEP_ADDRESS, e.g. what is
// left after payouts or funds sent by mistake. It is refused while a payout
// of that token is being funded or paid.
func (e *Engine) sweep(ctx context.Context) error {
	var s store.EngineState
	if err := e.DB.First(&s, 1).Error; err != nil || s.SweepToken == "" {
		return err
	}
	done := func(result string) error {
		e.publish(ctx, "sweep", 0, result)
		return e.DB.Model(&store.EngineState{}).Where("id = 1").Updates(map[string]interface{}{"sweep_token": "", "sweep_result": trim(result, 500)}).Error
	}
	if e.Cfg.SweepAddress == (common.Address{}) {
		return done("refused: PROCEED_PAYOUT_SWEEP_ADDRESS is not configured")
	}
	if !common.IsHexAddress(s.SweepToken) {
		return done("refused: not a token address")
	}
	token := common.HexToAddress(s.SweepToken)
	var busy int64
	e.DB.Model(&store.ProceedPayout{}).Where("LOWER(payout_contract_address) = ? AND status IN ?", strings.ToLower(token.Hex()),
		[]string{store.StatusFundingCheckRequested, store.StatusPaying, store.StatusPaused}).Count(&busy)
	if busy > 0 {
		return done(fmt.Sprintf("refused: %d payout(s) of this token are being funded or paid", busy))
	}
	bal, err := chain.BalanceOf(ctx, e.Chain, token, e.Cfg.PayoutSafe, nil)
	if err != nil {
		return err
	}
	if bal.Sign() == 0 {
		return done("nothing to sweep")
	}
	signers, err := e.Signers.Signers(ctx)
	if err != nil {
		return err
	}
	if err := safe.CheckSigners(ctx, e.Chain, e.Cfg.PayoutSafe, signers); err != nil {
		return err
	}
	b, err := e.build(ctx, signers, []safe.Call{{To: token, Data: chain.TransferData(e.Cfg.SweepAddress, bal)}})
	if err != nil {
		return done("failed: " + err.Error())
	}
	tx, err := safe.SignedTx(ctx, e.Chain, e.ChainID, signers[0], e.Cfg.PayoutSafe, b.execData, b.gas*13/10, e.Cfg.GasPriceMultiplier)
	if err != nil {
		return err
	}
	if err := e.Chain.SendTransaction(ctx, tx); err != nil && !alreadyKnown(err) {
		return done("failed: " + err.Error())
	}
	r, err := chain.WaitReceipt(ctx, e.Chain, tx.Hash(), e.Cfg.ReceiptTimeout)
	if err != nil {
		return done(fmt.Sprintf("sent %s; not mined yet - check it on the explorer", tx.Hash().Hex()))
	}
	if !safe.Executed(r, e.Cfg.PayoutSafe, b.safeTxHash) {
		return done(fmt.Sprintf("reverted: %s", tx.Hash().Hex()))
	}
	return done(fmt.Sprintf("swept %s (base units) to %s in %s, requested by %s", bal, e.Cfg.SweepAddress.Hex(), tx.Hash().Hex(), s.SweepRequestedBy))
}

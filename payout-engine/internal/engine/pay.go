package engine

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"trovo-payout-engine/internal/chain"
	"trovo-payout-engine/internal/config"
	"trovo-payout-engine/internal/safe"
	"trovo-payout-engine/internal/store"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func now() time.Time { return time.Now().UTC() }

// maxBatchesPerTick bounds one tick's paying so pauses, the kill switch and
// other payouts are seen between batches.
const maxBatchesPerTick = 20

// pay sends a payout's pending transfers, BatchSize at a time, each batch
// one Safe transaction through MultiSendCallOnly. A batch is simulated
// first; when it would revert, the transfers that fail are found by
// bisection (in simulation, costing no gas) and marked FAILED, and the rest
// go ahead. Every batch is recorded with its signed transaction before it
// is broadcast (see recoverBatches).
func (e *Engine) pay(ctx context.Context, p *store.ProceedPayout) error {
	for i := 0; i < maxBatchesPerTick; i++ {
		if ctx.Err() != nil {
			return nil
		}
		if halted, _ := e.halted(); halted || e.currentStatus(p.ID) != store.StatusPaying {
			return nil
		}
		var inFlight int64
		e.DB.Model(&store.PayoutBatch{}).Where("proceed_payout_id = ? AND status = ?", p.ID, store.BatchSubmitted).Count(&inFlight)
		if inFlight > 0 {
			return nil // recoverBatches settles it first
		}
		var items []store.PayoutItem
		if err := e.DB.Where("proceed_payout_id = ? AND status = ?", p.ID, store.ItemPending).Order("id").Limit(e.Cfg.BatchSize).Find(&items).Error; err != nil {
			return err
		}
		if len(items) == 0 {
			return e.finish(ctx, p)
		}
		signers, err := e.Signers.Signers(ctx)
		if err != nil {
			return err
		}
		if err := safe.CheckSigners(ctx, e.Chain, e.Cfg.PayoutSafe, signers); err != nil {
			return err
		}
		good, err := e.payable(ctx, p, signers, items)
		if err != nil {
			return err
		}
		if len(good) == 0 {
			continue // all of them failed; take the next ones
		}
		if err := e.sendBatch(ctx, p, signers, good); err != nil {
			return err
		}
	}
	return nil
}

func (e *Engine) calls(p *store.ProceedPayout, items []store.PayoutItem) []safe.Call {
	token := common.HexToAddress(p.PayoutContractAddress)
	out := make([]safe.Call, len(items))
	for i, it := range items {
		out[i] = safe.Call{To: token, Data: chain.TransferData(common.HexToAddress(it.BeneficiaryAddress), bigOf(it.AmountUnits))}
	}
	return out
}

// built is a signed Safe transaction for a set of calls.
type built struct {
	nonce      uint64
	safeTxHash []byte
	execData   []byte
	gas        uint64
}

// build signs calls for the Safe's current nonce and estimates its gas; a
// failing estimate means the batch would revert.
func (e *Engine) build(ctx context.Context, signers []config.Signer, calls []safe.Call) (*built, error) {
	to, data, op, err := safe.Target(calls, e.Cfg.MultiSendCallOnly)
	if err != nil {
		return nil, err
	}
	nonce, err := safe.Nonce(ctx, e.Chain, e.Cfg.PayoutSafe)
	if err != nil {
		return nil, err
	}
	hash, err := safe.TxHash(e.ChainID, e.Cfg.PayoutSafe, to, data, op, nonce)
	if err != nil {
		return nil, err
	}
	sigs, err := safe.Sign(hash, signers)
	if err != nil {
		return nil, err
	}
	execData, err := safe.ExecData(to, data, op, sigs)
	if err != nil {
		return nil, err
	}
	gas, err := e.Chain.EstimateGas(ctx, ethereum.CallMsg{From: signers[0].Address, To: &e.Cfg.PayoutSafe, Data: execData})
	if err != nil {
		return nil, errReverts{err}
	}
	return &built{nonce: nonce.Uint64(), safeTxHash: hash, execData: execData, gas: gas}, nil
}

type errReverts struct{ err error }

func (e errReverts) Error() string { return "the batch would revert: " + e.err.Error() }

// payable returns the items that can be paid together within MaxBatchGas,
// marking the ones whose transfer reverts FAILED.
func (e *Engine) payable(ctx context.Context, p *store.ProceedPayout, signers []config.Signer, items []store.PayoutItem) ([]store.PayoutItem, error) {
	b, err := e.build(ctx, signers, e.calls(p, items))
	if err == nil {
		if b.gas <= e.Cfg.MaxBatchGas || len(items) == 1 {
			return items, nil
		}
		return e.payable(ctx, p, signers, items[:len(items)/2]) // too big for one transaction
	}
	var rev errReverts
	if !errors.As(err, &rev) {
		return nil, err
	}
	if len(items) == 1 {
		reason := trim("its transfer reverts: "+rev.err.Error(), 300)
		e.DB.Model(&store.PayoutItem{}).Where("id = ? AND status = ?", items[0].ID, store.ItemPending).
			Updates(map[string]interface{}{"status": store.ItemFailed, "reason": reason, "cannot_receive_asset": 1, "updated_at": now()})
		return nil, nil
	}
	mid := len(items) / 2
	left, err := e.payable(ctx, p, signers, items[:mid])
	if err != nil {
		return nil, err
	}
	right, err := e.payable(ctx, p, signers, items[mid:])
	if err != nil {
		return nil, err
	}
	both := append(left, right...)
	if len(both) == len(items) {
		// each half passes but not together (e.g. the Safe's balance runs
		// out): pay one half now, the other next round
		return left, nil
	}
	return both, nil
}

// sendBatch records, broadcasts and settles one batch.
func (e *Engine) sendBatch(ctx context.Context, p *store.ProceedPayout, signers []config.Signer, items []store.PayoutItem) error {
	b, err := e.build(ctx, signers, e.calls(p, items))
	if err != nil {
		return err
	}
	tx, err := safe.SignedTx(ctx, e.Chain, e.ChainID, signers[0], e.Cfg.PayoutSafe, b.execData, b.gas*13/10, e.Cfg.GasPriceMultiplier)
	if err != nil {
		return err
	}
	raw, err := tx.MarshalBinary()
	if err != nil {
		return err
	}
	total := new(big.Int)
	ids := make([]string, len(items))
	for i, it := range items {
		total.Add(total, bigOf(it.AmountUnits))
		ids[i] = it.ID
	}
	batch := store.PayoutBatch{
		CreatedAt: now(), UpdatedAt: now(), ProceedPayoutID: p.ID, Status: store.BatchSubmitted, SafeNonce: b.nonce,
		SafeTxHash: hexutil.Encode(b.safeTxHash), TxHash: tx.Hash().Hex(), RawTx: hexutil.Encode(raw), Executor: signers[0].Address.Hex(),
		ItemCount: len(items), AmountUnits: total.String(),
	}
	err = e.DB.Transaction(func(dbtx *gorm.DB) error {
		if err := dbtx.Create(&batch).Error; err != nil {
			return err
		}
		res := dbtx.Model(&store.PayoutItem{}).Where("id IN ? AND status = ?", ids, store.ItemPending).
			Updates(map[string]interface{}{"status": store.ItemQueued, "batch_id": batch.ID, "updated_at": now()})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != int64(len(ids)) {
			return errors.New("items changed while the batch was built")
		}
		return nil
	})
	if err != nil {
		return err
	}
	if err := e.Chain.SendTransaction(ctx, tx); err != nil && !alreadyKnown(err) {
		// not broadcast: nothing can execute it, so put the items back
		e.dropBatch(batch, "broadcast failed: "+err.Error())
		return fmt.Errorf("broadcasting batch %d: %w", batch.ID, err)
	}
	e.activity(fmt.Sprintf("payout %d: batch %d (%d transfers) sent %s", p.ID, batch.ID, len(items), batch.TxHash))
	r, err := chain.WaitReceipt(ctx, e.Chain, tx.Hash(), e.Cfg.ReceiptTimeout)
	if err != nil {
		return nil // still pending: recoverBatches follows it
	}
	return e.settle(ctx, batch, r)
}

func alreadyKnown(err error) bool {
	m := strings.ToLower(err.Error())
	return strings.Contains(m, "already known") || strings.Contains(m, "known transaction")
}

// settle records a mined batch: its items PAID when the Safe executed it,
// back to PENDING when it reverted.
func (e *Engine) settle(ctx context.Context, batch store.PayoutBatch, r *types.Receipt) error {
	hash, _ := hexutil.Decode(batch.SafeTxHash)
	if safe.Executed(r, e.Cfg.PayoutSafe, hash) {
		paidAt := now()
		err := e.DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Model(&store.PayoutBatch{}).Where("id = ?", batch.ID).Updates(map[string]interface{}{"status": store.BatchMined, "gas_used": r.GasUsed, "updated_at": paidAt}).Error; err != nil {
				return err
			}
			var feeItems []store.PayoutItem
			if err := tx.Where("batch_id = ? AND status = ? AND kind IN ?", batch.ID, store.ItemQueued, []string{store.KindFee, store.KindVat}).Find(&feeItems).Error; err != nil {
				return err
			}
			if err := tx.Model(&store.PayoutItem{}).Where("batch_id = ? AND status = ?", batch.ID, store.ItemQueued).
				Updates(map[string]interface{}{"status": store.ItemPaid, "tx_hash": batch.TxHash, "paid_at": paidAt, "updated_at": paidAt}).Error; err != nil {
				return err
			}
			return e.recordFees(tx, batch, feeItems, paidAt)
		})
		if err != nil {
			return err
		}
		e.recount(batch.ProceedPayoutID)
		e.publish(ctx, "batch-mined", batch.ProceedPayoutID, fmt.Sprintf("batch %d: %d transfers in %s", batch.ID, batch.ItemCount, batch.TxHash))
		return nil
	}
	err := e.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&store.PayoutBatch{}).Where("id = ?", batch.ID).Updates(map[string]interface{}{"status": store.BatchReverted, "gas_used": r.GasUsed, "error": "reverted on-chain", "updated_at": now()}).Error; err != nil {
			return err
		}
		return tx.Model(&store.PayoutItem{}).Where("batch_id = ? AND status = ?", batch.ID, store.ItemQueued).
			Updates(map[string]interface{}{"status": store.ItemPending, "batch_id": 0, "updated_at": now()}).Error
	})
	e.publish(ctx, "batch-reverted", batch.ProceedPayoutID, fmt.Sprintf("batch %d reverted (%s); its transfers will be retried", batch.ID, batch.TxHash))
	return err
}

// dropBatch puts a batch's items back to PENDING (it never executed).
func (e *Engine) dropBatch(batch store.PayoutBatch, why string) {
	e.DB.Transaction(func(tx *gorm.DB) error {
		tx.Model(&store.PayoutBatch{}).Where("id = ?", batch.ID).Updates(map[string]interface{}{"status": store.BatchDropped, "error": trim(why, 500), "updated_at": now()})
		return tx.Model(&store.PayoutItem{}).Where("batch_id = ? AND status = ?", batch.ID, store.ItemQueued).
			Updates(map[string]interface{}{"status": store.ItemPending, "batch_id": 0, "updated_at": now()}).Error
	})
}

// recoverBatches settles batches left SUBMITTED (the engine stopped, or a
// receipt took longer than ReceiptTimeout): mined ones by their receipt;
// unmined ones are rebroadcast while the Safe nonce they use is still
// free, and dropped (their items paid again later) once it is not.
func (e *Engine) recoverBatches(ctx context.Context) error {
	var batches []store.PayoutBatch
	if err := e.DB.Where("status = ?", store.BatchSubmitted).Order("id").Find(&batches).Error; err != nil {
		return err
	}
	for _, b := range batches {
		hash := common.HexToHash(b.TxHash)
		if r, err := e.Chain.TransactionReceipt(ctx, hash); err == nil && r != nil {
			if err := e.settle(ctx, b, r); err != nil {
				return err
			}
			continue
		}
		nonce, err := safe.Nonce(ctx, e.Chain, e.Cfg.PayoutSafe)
		if err != nil {
			return err
		}
		if nonce.Uint64() > b.SafeNonce {
			// the nonce was used by another transaction, so this one can
			// never execute
			e.dropBatch(b, fmt.Sprintf("its Safe nonce %d was used by another transaction", b.SafeNonce))
			continue
		}
		raw, err := hexutil.Decode(b.RawTx)
		if err != nil {
			e.dropBatch(b, "unreadable signed transaction")
			continue
		}
		var tx types.Transaction
		if err := tx.UnmarshalBinary(raw); err != nil {
			e.dropBatch(b, "unreadable signed transaction")
			continue
		}
		if err := e.Chain.SendTransaction(ctx, &tx); err != nil && !alreadyKnown(err) {
			if strings.Contains(strings.ToLower(err.Error()), "nonce too low") {
				// the executor's nonce went to another transaction
				e.dropBatch(b, "replaced: "+err.Error())
				continue
			}
			return fmt.Errorf("rebroadcasting batch %d: %w", b.ID, err)
		}
		if time.Since(b.CreatedAt) > 15*time.Minute {
			e.note(b.ProceedPayoutID, fmt.Sprintf("batch %d (%s) is not mined after %v; check the executor's gas price", b.ID, b.TxHash, time.Since(b.CreatedAt).Round(time.Minute)))
		}
	}
	return nil
}

// recount refreshes a payout's totals from its items.
func (e *Engine) recount(id uint64) {
	type agg struct {
		Status string
		N      int
	}
	var rows []agg
	holders := e.DB.Where("proceed_payout_id = ? AND (kind = ? OR kind = '' OR kind IS NULL)", id, store.KindHolder)
	holders.Model(&store.PayoutItem{}).Select("status, COUNT(*) AS n").Group("status").Scan(&rows)
	counts := map[string]int{}
	for _, r := range rows {
		counts[r.Status] = r.N
	}
	var paid []store.PayoutItem
	e.DB.Select("amount_units").Where("proceed_payout_id = ? AND status = ? AND (kind = ? OR kind = '' OR kind IS NULL)", id, store.ItemPaid, store.KindHolder).Find(&paid)
	sum := new(big.Int)
	for _, it := range paid {
		sum.Add(sum, bigOf(it.AmountUnits))
	}
	e.DB.Model(&store.ProceedPayout{}).Where("id = ?", id).Updates(map[string]interface{}{
		"paid_count": counts[store.ItemPaid], "failed_count": counts[store.ItemFailed], "excluded_count": counts[store.ItemExcluded], "paid_units": sum.String(),
	})
}

// finish closes a payout once nothing is left to pay.
func (e *Engine) finish(ctx context.Context, p *store.ProceedPayout) error {
	var queued int64
	e.DB.Model(&store.PayoutItem{}).Where("proceed_payout_id = ? AND status = ?", p.ID, store.ItemQueued).Count(&queued)
	if queued > 0 {
		return nil
	}
	e.recount(p.ID)
	var failed int64
	e.DB.Model(&store.PayoutItem{}).Where("proceed_payout_id = ? AND status = ?", p.ID, store.ItemFailed).Count(&failed)
	to := store.StatusCompleted
	note := "All holders paid."
	if failed > 0 {
		to = store.StatusCompletedWithFailures
		note = fmt.Sprintf("%d holder(s) could not be paid; retry them or mark them paid after paying them another way.", failed)
	}
	ok, err := e.setStatus(p.ID, store.StatusPaying, to, map[string]interface{}{"completed_at": now(), "payout_completed": 1, "note": note})
	if err == nil && ok {
		e.activity(fmt.Sprintf("payout %d %s", p.ID, strings.ToLower(to)))
		e.publish(ctx, "completed", p.ID, note)
	}
	return err
}

// recordFees records the payout fee and its VAT, once paid, in
// fee_collections (TM's fee report), like the platform's other fees.
func (e *Engine) recordFees(tx *gorm.DB, batch store.PayoutBatch, items []store.PayoutItem, paidAt time.Time) error {
	if len(items) == 0 {
		return nil
	}
	var p store.ProceedPayout
	if err := tx.First(&p, batch.ProceedPayoutID).Error; err != nil {
		return err
	}
	var asset store.TokenizedAsset
	tx.Select("asset_code").First(&asset, "id = ?", p.TokenizedAssetID)
	from := p.Batch
	if asset.AssetCode != nil {
		from = *asset.AssetCode
	}
	for _, it := range items {
		feeType := store.FeeTypePayout
		if it.Kind == store.KindVat {
			feeType = store.FeeTypePayoutVat
		}
		contract, hash := p.PayoutContractAddress, batch.TxHash
		amount, _ := decimal.NewFromBigInt(bigOf(it.AmountUnits), -int32(p.PayoutDecimals)).Float64()
		rec := store.FeeCollection{
			CreatedAt: paidAt, UpdatedAt: paidAt, ID: "payout-" + it.ID, FromUsername: from, FromWalletAddress: p.PayoutSafeAddress,
			FromWalletAlias: "payout " + p.Batch, FeeType: feeType, Amount: amount, AssetCode: p.PayoutAssetCode,
			ContractAddress: &contract, DestinationWallet: it.BeneficiaryAddress, TransactionHash: &hash,
		}
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&rec).Error; err != nil {
			return err
		}
	}
	return nil
}

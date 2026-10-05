package users

import (
	"math/big"
	"strings"
	"time"

	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/shopspring/decimal"
)

// ProceedPayoutReceipt is a proceeds payout to one of the user's wallets
// (paid, or scheduled in a locked payout), for the app's dividend history.
type ProceedPayoutReceipt struct {
	ID                 string     `json:"id"`
	ProceedPayoutID    uint64     `json:"proceedPayoutId"`
	TokenizedAssetID   string     `json:"tokenizedAssetId"`
	AssetCode          string     `json:"assetCode"`
	AssetName          string     `json:"assetName"`
	PayoutAssetCode    string     `json:"payoutAssetCode"`
	PayoutContract     string     `json:"payoutContractAddress"`
	Amount             string     `json:"amount"`         // in the payout token
	TokensHeld         string     `json:"tokensHeld"`     // at the snapshot
	AmountPerToken     string     `json:"amountPerToken"` // in the payout token
	BeneficiaryAddress string     `json:"beneficiaryAddress"`
	WalletAlias        string     `json:"walletAlias"`
	Status             string     `json:"status"` // PENDING (scheduled), QUEUED, PAID, FAILED
	TxHash             string     `json:"txHash"`
	PaidAt             *time.Time `json:"paidAt"`
	ScheduledAt        *time.Time `json:"scheduledAt"` // when the schedule was locked
}

// payoutVisibleStatuses are the payouts a holder sees: once their schedule
// is locked, and not cancelled.
var payoutVisibleStatuses = []string{
	userModels.ProceedPayoutStatusLocked, userModels.ProceedPayoutStatusApproved, userModels.ProceedPayoutStatusFundingCheckRequested,
	userModels.ProceedPayoutStatusPaying, userModels.ProceedPayoutStatusPaused, userModels.ProceedPayoutStatusCompleted,
	userModels.ProceedPayoutStatusCompletedWithFailures,
}

// unitsToHuman turns base units into a decimal string.
func unitsToHuman(units string, decimals int) string {
	n, ok := new(big.Int).SetString(strings.TrimSpace(units), 10)
	if !ok {
		return "0"
	}
	return decimal.NewFromBigInt(n, -int32(decimals)).String()
}

// ListProceedPayoutReceipts lists the proceeds payouts to user's wallets,
// newest first, optionally for one tokenized asset.
func ListProceedPayoutReceipts(user *userModels.User, tokenizedAssetID string, gc *sharedconfig.GlobalConfig) ([]ProceedPayoutReceipt, error) {
	out := make([]ProceedPayoutReceipt, 0)
	var wallets []userModels.UserWallet
	if err := gc.DB.Where("user_id = ?", user.ID).Find(&wallets).Error; err != nil {
		return nil, &tErrors.ErrorTemporaryServerError{}
	}
	if len(wallets) == 0 {
		return out, nil
	}
	aliases := make(map[string]string, len(wallets))
	addrs := make([]string, 0, len(wallets))
	for _, w := range wallets {
		a := strings.ToLower(w.ID)
		aliases[a] = w.Alias
		addrs = append(addrs, a)
	}
	var payouts []userModels.ProceedPayout
	q := gc.DB.Preload("TokenizedAsset").Where("status IN ?", payoutVisibleStatuses)
	if tokenizedAssetID != "" {
		q = q.Where("tokenized_asset_id = ?", tokenizedAssetID)
	}
	if err := q.Order("id DESC").Limit(200).Find(&payouts).Error; err != nil {
		return nil, &tErrors.ErrorTemporaryServerError{}
	}
	if len(payouts) == 0 {
		return out, nil
	}
	byID := make(map[uint64]userModels.ProceedPayout, len(payouts))
	ids := make([]uint64, 0, len(payouts))
	for _, p := range payouts {
		byID[p.ID] = p
		ids = append(ids, p.ID)
	}
	var items []userModels.TokenizedAssetPayoutSchedule
	if err := gc.DB.Where("proceed_payout_id IN ? AND LOWER(beneficiary_address) IN ? AND kind = ? AND status IN ?", ids, addrs, userModels.PayoutItemKindHolder,
		[]string{userModels.PayoutItemPending, userModels.PayoutItemQueued, userModels.PayoutItemPaid, userModels.PayoutItemFailed}).
		Order("proceed_payout_id DESC").Find(&items).Error; err != nil {
		return nil, &tErrors.ErrorTemporaryServerError{}
	}
	for _, it := range items {
		p := byID[it.ProceedPayoutID]
		r := ProceedPayoutReceipt{
			ID: it.ID, ProceedPayoutID: p.ID, TokenizedAssetID: p.TokenizedAssetID, PayoutAssetCode: p.PayoutAssetCode,
			PayoutContract: p.PayoutContractAddress, Amount: unitsToHuman(it.AmountUnits, p.PayoutDecimals),
			TokensHeld: unitsToHuman(it.BalanceUnits, p.TokenDecimals), AmountPerToken: p.AmountPerToken,
			BeneficiaryAddress: it.BeneficiaryAddress, WalletAlias: aliases[strings.ToLower(it.BeneficiaryAddress)],
			Status: it.Status, TxHash: it.TxHash, PaidAt: it.PaidAt, ScheduledAt: p.LockedAt,
		}
		if p.TokenizedAsset.AssetCode != nil {
			r.AssetCode = *p.TokenizedAsset.AssetCode
		}
		if p.TokenizedAsset.AssetName != nil {
			r.AssetName = *p.TokenizedAsset.AssetName
		}
		out = append(out, r)
	}
	return out, nil
}

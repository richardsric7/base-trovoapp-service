package engine

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"trovo-payout-engine/internal/store"
)

// notifyPaid tells holders about their payouts once paid: one push
// notification per user and payout (a user's wallets together), sent once.
func (e *Engine) notifyPaid(ctx context.Context) error {
	var items []store.PayoutItem
	if err := e.DB.Where("status = ? AND notified = 0 AND (kind = ? OR kind = '' OR kind IS NULL)", store.ItemPaid, store.KindHolder).Order("proceed_payout_id").Limit(2000).Find(&items).Error; err != nil {
		return err
	}
	type group struct {
		payout   uint64
		username string
		amount   *big.Int
		ids      []string
	}
	groups := map[string]*group{}
	var silent []string // no Trovo user to tell
	for _, it := range items {
		if it.Username == "" {
			silent = append(silent, it.ID)
			continue
		}
		k := fmt.Sprintf("%d|%s", it.ProceedPayoutID, it.Username)
		if groups[k] == nil {
			groups[k] = &group{payout: it.ProceedPayoutID, username: it.Username, amount: new(big.Int)}
		}
		groups[k].amount.Add(groups[k].amount, bigOf(it.AmountUnits))
		groups[k].ids = append(groups[k].ids, it.ID)
	}
	if len(silent) > 0 {
		e.DB.Model(&store.PayoutItem{}).Where("id IN ?", silent).Update("notified", 1)
	}
	payouts := map[uint64]store.ProceedPayout{}
	for _, g := range groups {
		p, ok := payouts[g.payout]
		if !ok {
			if err := e.DB.First(&p, g.payout).Error; err != nil {
				continue
			}
			payouts[g.payout] = p
		}
		if e.Push != nil {
			var u store.User
			if e.DB.Where("username = ?", g.username).First(&u).Error == nil && u.PushNotificationToken != nil && *u.PushNotificationToken != "" {
				title := "Proceeds payout received"
				body := fmt.Sprintf("%s %s from %s has been paid to your wallet.", units(g.amount, p.PayoutDecimals), p.PayoutAssetCode, e.assetName(p))
				data := map[string]string{"route": "dividendHistory", "proceedPayoutId": fmt.Sprint(p.ID), "tokenizedAssetId": p.TokenizedAssetID}
				if err := e.Push.Send(ctx, *u.PushNotificationToken, title, body, data); err != nil {
					log.Printf("[notify] %s: %v", g.username, err) // a stale token: not retried
				}
			}
		}
		e.DB.Model(&store.PayoutItem{}).Where("id IN ?", g.ids).Update("notified", 1)
	}
	return nil
}

func (e *Engine) assetName(p store.ProceedPayout) string {
	var a store.TokenizedAsset
	if e.DB.Select("asset_name, asset_code").First(&a, "id = ?", p.TokenizedAssetID).Error == nil {
		if a.AssetName != nil && *a.AssetName != "" {
			return *a.AssetName
		}
		if a.AssetCode != nil {
			return *a.AssetCode
		}
	}
	return "your tokenized asset"
}

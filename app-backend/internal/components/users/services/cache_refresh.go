package users

import (
	"fmt"
	"strings"
	"time"

	paymentModels "trovo-wallet-api/internal/components/payments/models"
	userModels "trovo-wallet-api/internal/components/users/models"
	"trovo-wallet-api/internal/sharedconfig"
)

// RefreshCachesForNewPayments clears the cached balances, user objects and
// payment history of every wallet that sent or received a payment since
// the last call - including deposits from outside Trovo, which no request
// here sees. It reads the payment history payment-history-engine records
// (every native and token transfer touching a tracked wallet) rather than
// the chain itself.
func RefreshCachesForNewPayments(gc *sharedconfig.GlobalConfig) {
	cur := userModels.PaymentWatchCursor{ID: "cache-refresh"}
	if gc.DB.First(&cur, "id = ?", cur.ID).Error != nil {
		cur.At = time.Now().UTC().Add(-time.Minute)
		if gc.DB.Create(&cur).Error != nil {
			return
		}
	}
	type row struct {
		FromAddress     string
		ToAddress       string
		TransactionDate time.Time
	}
	var rows []row
	err := gc.DB.Model(&paymentModels.PaymentHistory{}).Select("from_address, to_address, transaction_date").
		Where("transaction_date > ?", cur.At).Order("transaction_date ASC").Limit(500).Scan(&rows).Error
	if err != nil || len(rows) == 0 {
		return
	}
	set := map[string]bool{}
	for _, r := range rows {
		set[strings.ToUpper(r.FromAddress)] = true
		set[strings.ToUpper(r.ToAddress)] = true
	}
	var addresses []string
	for a := range set {
		if a != "" {
			addresses = append(addresses, a)
		}
	}
	var wallets []userModels.UserWallet
	gc.DB.Where("UPPER(id) IN ?", addresses).Find(&wallets)
	for i := range wallets {
		w := &wallets[i]
		w.InvalidateUserCache(gc)
		gc.RedisCache.InvalidateCachedHttpResponse(fmt.Sprintf("GetBalance_%s", w.ID), fmt.Sprintf("[GET] /v1/users/payments/%v", w.ID), fmt.Sprintf("GetNFTs_%s", w.ID))
		if owner, err := w.GetWalletOwner(gc.DB, gc); err == nil {
			gc.RedisCache.InvalidateCachedHttpResponse(fmt.Sprintf("[GET] /v1/users/%v", owner.Username))
		}
	}
	gc.DB.Model(&userModels.PaymentWatchCursor{}).Where("id = ?", cur.ID).Update("at", rows[len(rows)-1].TransactionDate)
}

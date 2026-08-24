package users

import (
	"log"
	"time"

	userModels "trovo-wallet-api/internal/components/users/models"
	"trovo-wallet-api/internal/sharedconfig"
)

func GetPaymentConfigByServiceProvider(serviceProvider string, gc *sharedconfig.GlobalConfig) (config userModels.FiatPaymentConfig, err error) {
	err = gc.DB.Where("service_provider = ?", serviceProvider).First(&config).Error
	return
}

func GetFaucetConfigByUserCase(useCase string, gc *sharedconfig.GlobalConfig) (config userModels.FaucetConfig, err error) {
	err = gc.DB.Where("use_case = ?", useCase).First(&config).Error
	return
}

func SavePaymentWebhookData(provider, data string, gc *sharedconfig.GlobalConfig) error {

	t := userModels.PaymentWebhookRequest{
		ServiceProvider: provider,
		Data:            data,
	}

	return gc.DB.Save(&t).Error
}

func SaveUserPaymentData(username, provider, paymentType, txID string, amount float64, gc *sharedconfig.GlobalConfig) error {

	t := userModels.FiatPayment{
		ServiceProvider: provider,
		Username:        username,
		TransactionID:   txID,
		Amount:          amount,
		PaymentType:     paymentType,
	}

	return gc.DB.Save(&t).Error
}

func SaveUserPaymentInvoiceData(username, provider, paymentType, txID, status string, walletAlias, walletPublicKey, tokenizedAssetID, transaction, transactionSignature *string, amount float64, gc *sharedconfig.GlobalConfig) error {

	t := userModels.FiatPaymentInvoice{
		ID:                   txID,
		ServiceProvider:      provider,
		Username:             username,
		Amount:               amount,
		PaymentType:          paymentType,
		Status:               status,
		WalletAlias:          walletAlias,
		WalletPublicKey:      walletPublicKey,
		TokenizedAssetID:     tokenizedAssetID,
		Transaction:          transaction,
		TransactionSignature: transactionSignature,
	}

	return gc.DB.Save(&t).Error
}

func GetUserPaymentData(username string, gc *sharedconfig.GlobalConfig) (ps []userModels.FiatPayment) {
	ps = make([]userModels.FiatPayment, 0)

	gc.DB.Order("id DESC").Where("username = ?", username).Limit(100).Find(&ps)

	return
}
func GetUserPaymentInvoices(username string, gc *sharedconfig.GlobalConfig) (ps []userModels.FiatPaymentInvoice) {
	ps = make([]userModels.FiatPaymentInvoice, 0)

	gc.DB.Order("created_at DESC").Where("username = ?", username).Limit(100).Find(&ps)

	return
}

// ExpireStalePaymentInvoices flips any invoice that has been stuck PENDING for more than 2 days over to
// EXPIRED. This is the recovery path for a fiat asset purchase whose webhook never arrived (payment
// never completed, or the callback was lost) - it also releases any channel account still reserved for
// such an invoice's TokenizedAssetSubscription, so the account isn't leaked out of the pool forever.
func ExpireStalePaymentInvoices(gc *sharedconfig.GlobalConfig) error {
	cutoff := time.Now().Add(-48 * time.Hour)
	var stale []userModels.FiatPaymentInvoice
	if err := gc.DB.Where("status = ? AND created_at < ?", "PENDING", cutoff).Find(&stale).Error; err != nil {
		return err
	}
	for _, inv := range stale {
		if inv.TransactionSource != nil {
			gc.ReleaseInUseChannelAccount(*inv.TransactionSource)
			log.Printf("[ExpireStalePaymentInvoices] released channel account %v held by expiring invoice %v\n", *inv.TransactionSource, inv.ID)
		}
		if err := gc.DB.Model(&userModels.FiatPaymentInvoice{}).Where("id = ?", inv.ID).Update("status", "EXPIRED").Error; err != nil {
			log.Printf("[ExpireStalePaymentInvoices] error expiring invoice %v: %v\n", inv.ID, err)
		}
	}
	return nil
}

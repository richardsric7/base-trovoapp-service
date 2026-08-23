package users

import (
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

func SaveUserPaymentInvoiceData(username, provider, paymentType, txID, status string, walletAlias, walletPublicKey, tokenizedAssetID *string, amount float64, gc *sharedconfig.GlobalConfig) error {

	t := userModels.FiatPaymentInvoice{
		ID:               txID,
		ServiceProvider:  provider,
		Username:         username,
		Amount:           amount,
		PaymentType:      paymentType,
		Status:           status,
		WalletAlias:      walletAlias,
		WalletPublicKey:  walletPublicKey,
		TokenizedAssetID: tokenizedAssetID,
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

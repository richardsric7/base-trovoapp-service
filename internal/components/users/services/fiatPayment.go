package users

import (
	userModels "trovo-wallet-api/internal/components/users/models"
	"trovo-wallet-api/internal/sharedconfig"
)

func GetPaymentConfigByServiceProvider(serviceProvider string, gc *sharedconfig.GlobalConfig) (config userModels.FiatPaymentConfig, err error) {
	err = gc.DB.Where("service_provider = ?", serviceProvider).First(&config).Error
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

package users

import "time"

type FiatPaymentConfig struct {
	ID               uint64 `json:"-"`
	ServiceProvider  string `json:"serviceProvider"`
	Token            string `json:"token"`
	PublicKey        string `json:"publicKey"`
	SecretKey        string `json:"secretKey"`
	EncryptionKey    string `json:"encryptionKey"`
	VerificationHash string `json:"verificationHash"`
}

type PaymentWebhookRequest struct {
	ID              uint64
	ServiceProvider string
	Data            string
}

type FiatPayment struct {
	ID              uint64
	ServiceProvider string
	Username        string
	TransactionID   string
	Amount          float64
	PaymentType     string
}

type FlutterwaveWebhook struct {
	Event string `json:"event"`
	Data  struct {
		ID                int       `json:"id"`
		TxRef             string    `json:"tx_ref"`
		FlwRef            string    `json:"flw_ref"`
		DeviceFingerprint string    `json:"device_fingerprint"`
		Amount            int       `json:"amount"`
		Currency          string    `json:"currency"`
		ChargedAmount     int       `json:"charged_amount"`
		AppFee            int       `json:"app_fee"`
		MerchantFee       int       `json:"merchant_fee"`
		ProcessorResponse string    `json:"processor_response"`
		AuthModel         string    `json:"auth_model"`
		IP                string    `json:"ip"`
		Narration         string    `json:"narration"`
		Status            string    `json:"status"`
		PaymentType       string    `json:"payment_type"`
		CreatedAt         time.Time `json:"created_at"`
		AccountID         int       `json:"account_id"`
		Customer          struct {
			ID          int       `json:"id"`
			Name        string    `json:"name"`
			PhoneNumber string    `json:"phone_number"`
			Email       string    `json:"email"`
			CreatedAt   time.Time `json:"created_at"`
		} `json:"customer"`
	} `json:"data"`
	MetaData struct {
		CheckoutInitAddress string `json:"__CheckoutInitAddress"`
		ConsumerID          string `json:"consumer_id"`
		ConsumerMac         string `json:"consumer_mac"`
		UserID              string `json:"user_id"`
		TransactionType     string `json:"transaction_type"`
		Product             string `json:"product"`
	} `json:"meta_data"`
	EventType string `json:"event.type"`
}

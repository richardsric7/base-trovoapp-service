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
	CreatedAt       time.Time
	ID              uint64
	ServiceProvider string
	Data            string
}

type FaucetConfig struct {
	CreatedAt time.Time
	ID        uint64
	UseCase   string //ACTIVATION
	SecretKey string
}

type FiatPayment struct {
	CreatedAt       time.Time
	ID              uint64
	ServiceProvider string
	Username        string
	TransactionID   string
	Amount          float64
	PaymentType     string
}

type FiatPaymentInvoice struct {
	ID                   string    `json:"id"`
	CreatedAt            time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"createdAt"`
	ServiceProvider      string    `json:"serviceProvider"`
	Username             string    `json:"username"`
	Amount               float64   `json:"amount"`
	PaymentType          string    `json:"paymentType"`
	Status               string    `gorm:"default:'PENDING'" json:"status"`
	Refunded             int       `gorm:"default:0" json:"refunded"`
	TokenizedAssetID     *string   `gorm:"null;size:100" json:"tokenizedAssetId"`
	WalletAlias          *string   `gorm:"null;size:100" json:"walletAlias"`
	WalletPublicKey      *string   `gorm:"null;size:100" json:"walletPublicKey"`
	Transaction          *string   `json:"transaction"`
	TransactionSignature *string   `json:"transactionSignature"`
	TransactionSource    *string   `gorm:"null;size:100" json:"-"`
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

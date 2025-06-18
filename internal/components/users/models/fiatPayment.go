package users

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

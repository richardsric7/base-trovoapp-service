package models

// StablerailConfig mirrors app-backend's stablerail_configs table (one row):
// Stablerail's API key and address, and whether bank deposits and
// withdrawals are on. app-backend reads the first row on every request.
type StablerailConfig struct {
	ID               uint64 `gorm:"column:id;primaryKey"`
	ApiKey           string `gorm:"column:api_key"`
	FintechID        string `gorm:"column:fintech_id"`
	BaseUrl          string `gorm:"column:base_url"`
	EnableStablerail int    `gorm:"column:enable_stablerail;default:0"`
}

func (StablerailConfig) TableName() string { return "stablerail_configs" }

// MaskSecret shows enough of a secret to recognise it, never enough to use
// it. Empty stays empty, so "not set" is visible.
func MaskSecret(s string) string {
	switch {
	case s == "":
		return ""
	case len(s) <= 8:
		return "****"
	default:
		return s[:4] + "…" + s[len(s)-2:]
	}
}

// ProviderSecretStatus describes a stored secret without revealing it.
type ProviderSecretStatus struct {
	Set  bool   `json:"set"`
	Hint string `json:"hint" example:"sbx:…Xy"`
}

func SecretStatus(s string) ProviderSecretStatus {
	return ProviderSecretStatus{Set: s != "", Hint: MaskSecret(s)}
}

// KycProviderSettings is one KYC provider's credentials, masked.
type KycProviderSettings struct {
	ServiceProvider string               `json:"service_provider" example:"sumsub"`
	Token           ProviderSecretStatus `json:"token"`
	SecretKey       ProviderSecretStatus `json:"secret_key"`
}

// StablerailSettings is the Stablerail configuration, API key masked.
type StablerailSettings struct {
	Enabled   bool                 `json:"enabled"`
	BaseURL   string               `json:"base_url" example:"https://api.stablerail.example.com"`
	FintechID string               `json:"fintech_id" example:"ft_123"`
	ApiKey    ProviderSecretStatus `json:"api_key"`
}

// ProviderSettingsResponse is GET /provider-settings.
type ProviderSettingsResponse struct {
	Kyc        []KycProviderSettings `json:"kyc"`
	Stablerail StablerailSettings    `json:"stablerail"`
}

// KycProviderSettingsRequest updates one KYC provider. An empty field keeps
// the stored value, so a secret never has to be re-entered or shown.
type KycProviderSettingsRequest struct {
	Token     string `json:"token" example:"sbx:AbC..."`
	SecretKey string `json:"secret_key" example:"<secret>"`
}

// StablerailSettingsRequest updates Stablerail. An empty api_key keeps the
// stored key.
type StablerailSettingsRequest struct {
	Enabled   bool   `json:"enabled"`
	BaseURL   string `json:"base_url" example:"https://api.stablerail.example.com"`
	FintechID string `json:"fintech_id" example:"ft_123"`
	ApiKey    string `json:"api_key" example:"<api key>"`
}

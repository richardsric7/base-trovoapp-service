package models

// FaucetConfig represents a faucet configuration record
type FaucetConfig struct {
	ID              int64  `gorm:"column:id;primaryKey" json:"id"`
	ServiceProvider string `gorm:"column:service_provider" json:"service_provider"`
	Token           string `gorm:"column:token" json:"token"`
	SecretKey       string `gorm:"column:secret_key" json:"secret_key"`
}

// TableName overrides the default gorm table name
func (FaucetConfig) TableName() string {
	return "faucet_configs"
}

// KycConfig represents a KYC configuration record
type KycConfig struct {
	ID              int64  `gorm:"column:id;primaryKey" json:"id"`
	ServiceProvider string `gorm:"column:service_provider" json:"service_provider"`
	Token           string `gorm:"column:token" json:"token"`
	SecretKey       string `gorm:"column:secret_key" json:"secret_key"`
}

// TableName overrides the default gorm table name
func (KycConfig) TableName() string {
	return "kyc_configs"
}

// FaucetConfigRequest represents the request payload for create/update operations
type FaucetConfigRequest struct {
	Action          string `json:"action" binding:"required,oneof=create update" example:"create" enums:"create,update"`
	ID              int64  `json:"id,omitempty" example:"1"`
	ServiceProvider string `json:"service_provider" binding:"required" example:"flutterwave"`
	Token           string `json:"token" binding:"required" example:"FLW_TOKEN_123"`
	SecretKey       string `json:"secret_key" binding:"required" example:"FLW_SECRET_456"`
}

// KycConfigRequest represents the request payload for create/update operations
type KycConfigRequest struct {
	Action          string `json:"action" binding:"required,oneof=create update" example:"create" enums:"create,update"`
	ID              int64  `json:"id,omitempty" example:"1"`
	ServiceProvider string `json:"service_provider" binding:"required" example:"flutterwave"`
	Token           string `json:"token" binding:"required" example:"FLW_TOKEN_123"`
	SecretKey       string `json:"secret_key" binding:"required" example:"KYC_SECRET_789"`
}

// FaucetConfigResponse represents the response for faucet config operations
type FaucetConfigResponse struct {
	ID              int64  `json:"id" example:"1"`
	ServiceProvider string `json:"service_provider" example:"flutterwave"`
	Token           string `json:"token" example:"FLW_TOKEN_123"`
	SecretKey       string `json:"secret_key" example:"FLW_SECRET_456"`
}

// KycConfigResponse represents the response for KYC config operations
type KycConfigResponse struct {
	ID              int64  `json:"id" example:"1"`
	ServiceProvider string `json:"service_provider" example:"flutterwave"`
	Token           string `json:"token" example:"FLW_TOKEN_123"`
	SecretKey       string `json:"secret_key" example:"KYC_SECRET_789"`
}

type DojaWidget struct {
	ID        string `gorm:"column:id;primaryKey;type:text" json:"id"`
	Level     int64  `gorm:"column:level;type:bigint" json:"level"`
	Corporate int64  `gorm:"column:corporate;type:bigint" json:"corporate"`
}

// TableName explicitly sets the table name
func (DojaWidget) TableName() string {
	return "doja_widgets"
}

// KycLevel represents a KYC level record
type KycLevel struct {
	ID           string `gorm:"column:id;primaryKey" json:"id"`
	UserCategory string `gorm:"column:user_category" json:"user_category"`
}

// TableName overrides the default gorm table name
func (KycLevel) TableName() string {
	return "kyc_levels"
}

// DojaWidgetRequest represents the request payload for create/update operations
type DojaWidgetRequest struct {
	Action    string `json:"action" binding:"required,oneof=create update" example:"create" enums:"create,update"`
	ID        string `json:"id,omitempty" example:"widget_001"`
	Level     int64  `json:"level" binding:"required" example:"1"`
	Corporate int64  `json:"corporate" binding:"required" example:"0"`
}

// KycLevelRequest represents the request payload for create/update operations
type KycLevelRequest struct {
	Action       string `json:"action" binding:"required,oneof=create update" example:"create" enums:"create,update"`
	ID           string `json:"id,omitempty" example:"level_001"`
	UserCategory string `json:"user_category" binding:"required" example:"individual"`
}

// DojaWidgetResponse represents the response for doja widget operations
type DojaWidgetResponse struct {
	ID        string `json:"id" example:"widget_001"`
	Level     int64  `json:"level" example:"1"`
	Corporate int64  `json:"corporate" example:"0"`
}

// KycLevelResponse represents the response for KYC level operations
type KycLevelResponse struct {
	ID           string `json:"id" example:"level_001"`
	UserCategory string `json:"user_category" example:"individual"`
}

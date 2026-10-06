package usermetrics

import (
	"errors"
	"net/http"
	"strings"

	"admin-panel-dashboard/internal/middleware"
	"admin-panel-dashboard/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// providerSettingsPermission gates reading and changing outside providers'
// credentials (KYC providers, Stablerail).
const providerSettingsPermission = "MANAGE_SETTINGS"

// KycProviders are the KYC providers app-backend reads from kyc_configs.
var KycProviders = []string{"sumsub", "doja"}

func checkProviderSettingsAccess(c *gin.Context, adminDB *gorm.DB) bool {
	ad, err := middleware.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return false
	}
	var admin models.AdminUser
	if err := adminDB.First(&admin, "wallet_user_id = ?", ad.UserID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return false
	}
	if !models.HasPermission(adminDB, &admin, providerSettingsPermission) {
		c.JSON(http.StatusForbidden, gin.H{"error": "you need the MANAGE_SETTINGS permission"})
		return false
	}
	return true
}

// LoadProviderSettings reads the KYC providers' and Stablerail's settings
// with every secret masked.
func LoadProviderSettings(walletDB *gorm.DB) (models.ProviderSettingsResponse, error) {
	out := models.ProviderSettingsResponse{Kyc: []models.KycProviderSettings{}}
	for _, p := range KycProviders {
		var cfg models.KycConfig
		err := walletDB.Where("service_provider = ?", p).First(&cfg).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return out, err
		}
		out.Kyc = append(out.Kyc, models.KycProviderSettings{
			ServiceProvider: p,
			Token:           models.SecretStatus(cfg.Token),
			SecretKey:       models.SecretStatus(cfg.SecretKey),
		})
	}
	var sr models.StablerailConfig
	if err := walletDB.Order("id").First(&sr).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return out, err
	}
	out.Stablerail = models.StablerailSettings{
		Enabled:   sr.EnableStablerail == 1,
		BaseURL:   sr.BaseUrl,
		FintechID: sr.FintechID,
		ApiKey:    models.SecretStatus(sr.ApiKey),
	}
	return out, nil
}

// SaveKycProviderSettings creates or updates one provider's row. Empty
// fields keep the stored value.
func SaveKycProviderSettings(walletDB *gorm.DB, provider string, req models.KycProviderSettingsRequest) error {
	var cfg models.KycConfig
	err := walletDB.Where("service_provider = ?", provider).First(&cfg).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if strings.TrimSpace(req.Token) == "" || strings.TrimSpace(req.SecretKey) == "" {
			return errors.New("token and secret_key are both needed the first time")
		}
		return walletDB.Create(&models.KycConfig{ServiceProvider: provider, Token: strings.TrimSpace(req.Token), SecretKey: strings.TrimSpace(req.SecretKey)}).Error
	}
	if err != nil {
		return err
	}
	updates := map[string]interface{}{}
	if v := strings.TrimSpace(req.Token); v != "" {
		updates["token"] = v
	}
	if v := strings.TrimSpace(req.SecretKey); v != "" {
		updates["secret_key"] = v
	}
	if len(updates) == 0 {
		return nil
	}
	return walletDB.Model(&models.KycConfig{}).Where("id = ?", cfg.ID).Updates(updates).Error
}

// SaveStablerailSettings creates or updates the single stablerail_configs
// row. An empty api_key keeps the stored key.
func SaveStablerailSettings(walletDB *gorm.DB, req models.StablerailSettingsRequest) error {
	enabled := 0
	if req.Enabled {
		enabled = 1
	}
	baseURL := strings.TrimRight(strings.TrimSpace(req.BaseURL), "/")
	if req.Enabled && baseURL == "" {
		return errors.New("base_url is needed to turn Stablerail on")
	}
	var sr models.StablerailConfig
	err := walletDB.Order("id").First(&sr).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if req.Enabled && strings.TrimSpace(req.ApiKey) == "" {
			return errors.New("api_key is needed to turn Stablerail on")
		}
		return walletDB.Create(&models.StablerailConfig{ApiKey: strings.TrimSpace(req.ApiKey), FintechID: strings.TrimSpace(req.FintechID), BaseUrl: baseURL, EnableStablerail: enabled}).Error
	}
	if err != nil {
		return err
	}
	if req.Enabled && sr.ApiKey == "" && strings.TrimSpace(req.ApiKey) == "" {
		return errors.New("api_key is needed to turn Stablerail on")
	}
	updates := map[string]interface{}{
		"base_url":          baseURL,
		"fintech_id":        strings.TrimSpace(req.FintechID),
		"enable_stablerail": enabled,
	}
	if v := strings.TrimSpace(req.ApiKey); v != "" {
		updates["api_key"] = v
	}
	return walletDB.Model(&models.StablerailConfig{}).Where("id = ?", sr.ID).Updates(updates).Error
}

// GetProviderSettingsHandler godoc
// @Summary Outside providers' settings
// @Description KYC providers' (Sumsub, Doja) and Stablerail's settings, with every key masked. Needs MANAGE_SETTINGS.
// @Tags ProviderSettings
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Success 200 {object} models.ProviderSettingsResponse
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /provider-settings [get]
func GetProviderSettingsHandler(adminDB, walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !checkProviderSettingsAccess(c, adminDB) {
			return
		}
		out, err := LoadProviderSettings(walletDB)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not read provider settings"})
			return
		}
		c.JSON(http.StatusOK, out)
	}
}

// SaveKycProviderSettingsHandler godoc
// @Summary Set a KYC provider's credentials
// @Description Creates or updates the kyc_configs row app-backend reads for this provider. Empty fields keep the stored value. Needs MANAGE_SETTINGS.
// @Tags ProviderSettings
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param provider path string true "sumsub or doja"
// @Param data body models.KycProviderSettingsRequest true "Credentials"
// @Success 200 {object} models.ProviderSettingsResponse
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /provider-settings/kyc/{provider} [put]
func SaveKycProviderSettingsHandler(adminDB, walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !checkProviderSettingsAccess(c, adminDB) {
			return
		}
		provider := strings.ToLower(c.Param("provider"))
		known := false
		for _, p := range KycProviders {
			known = known || p == provider
		}
		if !known {
			c.JSON(http.StatusBadRequest, gin.H{"error": "provider must be one of: " + strings.Join(KycProviders, ", ")})
			return
		}
		var req models.KycProviderSettingsRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
			return
		}
		if err := SaveKycProviderSettings(walletDB, provider, req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		out, _ := LoadProviderSettings(walletDB)
		c.JSON(http.StatusOK, out)
	}
}

// SaveStablerailSettingsHandler godoc
// @Summary Set Stablerail's settings
// @Description Creates or updates the stablerail_configs row app-backend reads on every bank deposit or withdrawal request. An empty api_key keeps the stored key. Needs MANAGE_SETTINGS.
// @Tags ProviderSettings
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token" default(Bearer <your-token>)
// @Param data body models.StablerailSettingsRequest true "Settings"
// @Success 200 {object} models.ProviderSettingsResponse
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /provider-settings/stablerail [put]
func SaveStablerailSettingsHandler(adminDB, walletDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !checkProviderSettingsAccess(c, adminDB) {
			return
		}
		var req models.StablerailSettingsRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
			return
		}
		if err := SaveStablerailSettings(walletDB, req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		out, _ := LoadProviderSettings(walletDB)
		c.JSON(http.StatusOK, out)
	}
}

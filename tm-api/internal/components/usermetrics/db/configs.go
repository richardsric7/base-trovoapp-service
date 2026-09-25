package usermetrics

import (
	"admin-panel-dashboard/internal/models"
	"errors"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

// SaveFaucetConfig creates or updates a faucet config
func SaveFaucetConfig(db *gorm.DB, req models.FaucetConfigRequest) error {
	if req.Action == "create" {
		config := models.FaucetConfig{
			ServiceProvider: req.ServiceProvider,
			Token:           req.Token,
			SecretKey:       req.SecretKey,
		}
		return db.Create(&config).Error
	} else if req.Action == "update" {
		if req.ID == 0 {
			return errors.New("ID is required for update")
		}
		var existing models.FaucetConfig
		if err := db.First(&existing, req.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("faucet config not found")
			}
			return err
		}
		updateData := models.FaucetConfig{
			ServiceProvider: req.ServiceProvider,
			Token:           req.Token,
			SecretKey:       req.SecretKey,
		}
		return db.Model(&models.FaucetConfig{}).Where("id = ?", req.ID).Updates(updateData).Error
	}
	return errors.New("invalid action: must be 'create' or 'update'")
}

// GetAllFaucetConfigs retrieves all faucet configs
func GetAllFaucetConfigs(db *gorm.DB) ([]models.FaucetConfigResponse, error) {
	var configs []models.FaucetConfig
	if err := db.Find(&configs).Error; err != nil {
		return nil, err
	}
	responses := make([]models.FaucetConfigResponse, len(configs))
	for i, config := range configs {
		responses[i] = models.FaucetConfigResponse{
			ID:              config.ID,
			ServiceProvider: config.ServiceProvider,
			Token:           config.Token,
			SecretKey:       config.SecretKey,
		}
	}
	return responses, nil
}

// GetFaucetConfigByID retrieves a faucet config by ID
func GetFaucetConfigByID(db *gorm.DB, id int64) (*models.FaucetConfigResponse, error) {
	var config models.FaucetConfig
	if err := db.First(&config, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("faucet config not found")
		}
		return nil, err
	}
	return &models.FaucetConfigResponse{
		ID:              config.ID,
		ServiceProvider: config.ServiceProvider,
		Token:           config.Token,
		SecretKey:       config.SecretKey,
	}, nil
}

// DeleteFaucetConfig deletes a faucet config by ID
func DeleteFaucetConfig(db *gorm.DB, id int64) error {
	return db.Delete(&models.FaucetConfig{}, id).Error
}

// SaveKycConfig creates or updates a KYC config
func SaveKycConfig(db *gorm.DB, req models.KycConfigRequest) error {
	if req.Action == "create" {
		config := models.KycConfig{
			ServiceProvider: req.ServiceProvider,
			Token:           req.Token,
			SecretKey:       req.SecretKey,
		}
		return db.Create(&config).Error
	} else if req.Action == "update" {
		if req.ID == 0 {
			return errors.New("ID is required for update")
		}
		var existing models.KycConfig
		if err := db.First(&existing, req.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("kyc config not found")
			}
			return err
		}
		updateData := models.KycConfig{
			ServiceProvider: req.ServiceProvider,
			Token:           req.Token,
			SecretKey:       req.SecretKey,
		}
		return db.Model(&models.KycConfig{}).Where("id = ?", req.ID).Updates(updateData).Error
	}
	return errors.New("invalid action: must be 'create' or 'update'")
}

// GetAllKycConfigs retrieves all KYC configs
func GetAllKycConfigs(db *gorm.DB) ([]models.KycConfigResponse, error) {
	var configs []models.KycConfig
	if err := db.Find(&configs).Error; err != nil {
		return nil, err
	}
	responses := make([]models.KycConfigResponse, len(configs))
	for i, config := range configs {
		responses[i] = models.KycConfigResponse{
			ID:              config.ID,
			ServiceProvider: config.ServiceProvider,
			Token:           config.Token,
			SecretKey:       config.SecretKey,
		}
	}
	return responses, nil
}

// GetKycConfigByID retrieves a KYC config by ID
func GetKycConfigByID(db *gorm.DB, id int64) (*models.KycConfigResponse, error) {
	var config models.KycConfig
	if err := db.First(&config, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("kyc config not found")
		}
		return nil, err
	}
	return &models.KycConfigResponse{
		ID:              config.ID,
		ServiceProvider: config.ServiceProvider,
		Token:           config.Token,
		SecretKey:       config.SecretKey,
	}, nil
}

// DeleteKycConfig deletes a KYC config by ID
func DeleteKycConfig(db *gorm.DB, id int64) error {
	return db.Delete(&models.KycConfig{}, id).Error
}

// SaveDojaWidget creates or updates a doja widget
func SaveDojaWidget(db *gorm.DB, req models.DojaWidgetRequest) error {
	if req.Action == "create" {
		widget := models.DojaWidget{
			ID:        uuid.New().String(),
			Level:     req.Level,
			Corporate: req.Corporate,
		}
		return db.Create(&widget).Error
	} else if req.Action == "update" {
		if req.ID == "" {
			return errors.New("ID is required for update")
		}
		var existing models.DojaWidget
		if err := db.First(&existing, "id = ?", req.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("doja widget not found")
			}
			return err
		}
		updateData := models.DojaWidget{
			Level:     req.Level,
			Corporate: req.Corporate,
		}
		return db.Model(&models.DojaWidget{}).Where("id = ?", req.ID).Updates(updateData).Error
	}
	return errors.New("invalid action: must be 'create' or 'update'")
}

// GetAllDojaWidgets retrieves all doja widgets
func GetAllDojaWidgets(db *gorm.DB) ([]models.DojaWidgetResponse, error) {
	var widgets []models.DojaWidget
	if err := db.Find(&widgets).Error; err != nil {
		return nil, err
	}
	responses := make([]models.DojaWidgetResponse, len(widgets))
	for i, widget := range widgets {
		responses[i] = models.DojaWidgetResponse{
			ID:        widget.ID,
			Level:     widget.Level,
			Corporate: widget.Corporate,
		}
	}
	return responses, nil
}

// GetDojaWidgetByID retrieves a doja widget by ID
func GetDojaWidgetByID(db *gorm.DB, id string) (*models.DojaWidgetResponse, error) {
	var widget models.DojaWidget
	if err := db.First(&widget, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("doja widget not found")
		}
		return nil, err
	}
	return &models.DojaWidgetResponse{
		ID:        widget.ID,
		Level:     widget.Level,
		Corporate: widget.Corporate,
	}, nil
}

// DeleteDojaWidget deletes a doja widget by ID
func DeleteDojaWidget(db *gorm.DB, id string) error {
	return db.Where("id = ?", id).Delete(&models.DojaWidget{}).Error
}

// SaveKycLevel creates or updates a KYC level
func SaveKycLevel(db *gorm.DB, req models.KycLevelRequest) error {
	if req.Action == "create" {
		// Check if user_category already exists (case-insensitive)
		var existing models.KycLevel
		if err := db.Where("LOWER(user_category) = LOWER(?)", req.UserCategory).First(&existing).Error; err == nil {
			return errors.New("user category already exists")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		level := models.KycLevel{
			ID:           uuid.New().String(),
			UserCategory: req.UserCategory,
		}
		return db.Create(&level).Error
	} else if req.Action == "update" {
		if req.ID == "" {
			return errors.New("ID is required for update")
		}
		var existing models.KycLevel
		if err := db.First(&existing, "id = ?", req.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("kyc level not found")
			}
			return err
		}

		// Check if user_category already exists for a different record (case-insensitive)
		var duplicate models.KycLevel
		if err := db.Where("LOWER(user_category) = LOWER(?) AND id != ?", req.UserCategory, req.ID).First(&duplicate).Error; err == nil {
			return errors.New("user category already exists")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		updateData := models.KycLevel{
			UserCategory: req.UserCategory,
		}
		return db.Model(&models.KycLevel{}).Where("id = ?", req.ID).Updates(updateData).Error
	}
	return errors.New("invalid action: must be 'create' or 'update'")
}

// GetAllKycLevels retrieves all KYC levels
func GetAllKycLevels(db *gorm.DB) ([]models.KycLevelResponse, error) {
	var levels []models.KycLevel
	if err := db.Find(&levels).Error; err != nil {
		return nil, err
	}
	responses := make([]models.KycLevelResponse, len(levels))
	for i, level := range levels {
		responses[i] = models.KycLevelResponse{
			ID:           level.ID,
			UserCategory: level.UserCategory,
		}
	}
	return responses, nil
}

// GetKycLevelByID retrieves a KYC level by ID
func GetKycLevelByID(db *gorm.DB, id string) (*models.KycLevelResponse, error) {
	var level models.KycLevel
	if err := db.First(&level, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("kyc level not found")
		}
		return nil, err
	}
	return &models.KycLevelResponse{
		ID:           level.ID,
		UserCategory: level.UserCategory,
	}, nil
}

// DeleteKycLevel deletes a KYC level by ID
func DeleteKycLevel(db *gorm.DB, id string) error {
	return db.Where("id = ?", id).Delete(&models.KycLevel{}).Error
}

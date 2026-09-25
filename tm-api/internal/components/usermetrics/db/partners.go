package usermetrics

import (
	"admin-panel-dashboard/internal/models"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"gorm.io/gorm"
)

func SaveAssetManager(action string, req models.AssetManagerRequest, db *gorm.DB) (uint64, error) {
	if action == "create" {
		entity := models.AssetManager{
			AssetManagerName:    req.AssetManagerName,
			AssetManagerAddress: req.AssetManagerAddress,
			AssetManagerCountry: req.AssetManagerCountry,
			FeePercent:          req.FeePercent,
			FeeFixed:            req.FeeFixed,
		}
		if err := db.Create(&entity).Error; err != nil {
			return 0, err
		}
		return uint64(entity.ID), nil
	} else if action == "update" {
		if req.ID <= 0 {
			return 0, errors.New("missing id in request body for update")
		}
		var existing models.AssetManager
		if err := db.First(&existing, req.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return 0, fmt.Errorf("asset manager with ID %d not found for update", req.ID)
			}
			return 0, err
		}

		updateData := models.AssetManager{
			AssetManagerName:    req.AssetManagerName,
			AssetManagerAddress: req.AssetManagerAddress,
			AssetManagerCountry: req.AssetManagerCountry,
			FeePercent:          req.FeePercent,
			FeeFixed:            req.FeeFixed,
		}
		if err := db.Model(&models.AssetManager{}).Where("id = ?", req.ID).Updates(updateData).Error; err != nil {
			return 0, err
		}
		return uint64(req.ID), nil
	}
	return 0, fmt.Errorf("invalid action specified: %s", action)
}

func GetPartnerEntityByID(entityType string, id int, db *gorm.DB) (interface{}, error) {
	switch entityType {
	case "asset_manager":
		var entity models.AssetManager
		err := db.First(&entity, id).Error
		return entity, err
	case "asset_issuing_house":
		var entity models.AssetIssuingHouse
		err := db.First(&entity, id).Error
		return entity, err
	case "approved_asset_custodian":
		var entity models.ApprovedAssetCustodian
		err := db.First(&entity, id).Error
		return entity, err
	case "legal_and_professionals":
		var entity models.LegalAndProfessionalsPartners
		err := db.First(&entity, id).Error
		return entity, err
	case "rating_agency":
		var entity models.RatingAgency
		err := db.First(&entity, id).Error
		return entity, err
	case "trustees":
		var entity models.Trustees
		err := db.First(&entity, id).Error
		return entity, err
	case "legal_adviser":
		var entity models.LegalAdviser
		err := db.First(&entity, id).Error
		return entity, err
	case "financial_adviser":
		var entity models.FinancialAdviser
		err := db.First(&entity, id).Error
		return entity, err
	default:
		return nil, fmt.Errorf("invalid partner type")
	}
}

func DeletePartnerEntity(typeStr, idStr string, db *gorm.DB) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return errors.New("invalid id")
	}
	switch typeStr {
	case "asset_manager":
		return db.Delete(&models.AssetManager{}, id).Error
	case "asset_issuing_house":
		return db.Delete(&models.AssetIssuingHouse{}, id).Error
	case "approved_asset_custodian":
		return db.Delete(&models.ApprovedAssetCustodian{}, id).Error
	case "legal_and_professionals":
		return db.Delete(&models.LegalAndProfessionalsPartners{}, id).Error
	case "rating_agency":
		return db.Delete(&models.RatingAgency{}, id).Error
	case "trustees":
		return db.Delete(&models.Trustees{}, id).Error
	case "legal_adviser":
		return db.Delete(&models.LegalAdviser{}, id).Error
	case "financial_adviser":
		return db.Delete(&models.FinancialAdviser{}, id).Error
	default:
		return errors.New("invalid type specified")
	}
}

func SavePartnerEntity(req models.SavePartnerRequest, db *gorm.DB) error {
	switch req.Type {
	case "asset_manager":
		return handleAssetManager(req, db)
	case "asset_issuing_house":
		return handleAssetIssuingHouse(req, db)
	case "approved_asset_custodian":
		return handleAssetCustodian(req, db)
	default:
		return errors.New("invalid type specified")
	}
}

func ListPartnerEntities(typeStr string, db *gorm.DB) (interface{}, error) {
	switch typeStr {
	case "asset_manager":
		var list []models.AssetManager
		if err := db.Find(&list).Error; err != nil {
			return nil, err
		}
		return list, nil
	case "asset_issuing_house":
		var list []models.AssetIssuingHouse
		if err := db.Find(&list).Error; err != nil {
			return nil, err
		}
		return list, nil
	case "approved_asset_custodian":
		var list []models.ApprovedAssetCustodian
		if err := db.Find(&list).Error; err != nil {
			return nil, err
		}
		return list, nil
	case "legal_and_professionals":
		var list []models.LegalAndProfessionalsPartners
		if err := db.Find(&list).Error; err != nil {
			return nil, err
		}
		return list, nil
	case "rating_agency":
		var list []models.RatingAgency
		if err := db.Find(&list).Error; err != nil {
			return nil, err
		}
		return list, nil
	case "trustees":
		var list []models.Trustees
		if err := db.Find(&list).Error; err != nil {
			return nil, err
		}
		return list, nil
	case "legal_adviser":
		var list []models.LegalAdviser
		if err := db.Find(&list).Error; err != nil {
			return nil, err
		}
		return list, nil
	case "financial_adviser":
		var list []models.FinancialAdviser
		if err := db.Find(&list).Error; err != nil {
			return nil, err
		}
		return list, nil
	default:
		return nil, errors.New("invalid type specified")
	}
}

func GetPartnerByID(typeStr string, id int, db *gorm.DB) (interface{}, error) {
	switch typeStr {
	case "asset_manager":
		var m models.AssetManager
		if err := db.First(&m, id).Error; err != nil {
			return nil, err
		}
		return m, nil
	case "asset_issuing_house":
		var i models.AssetIssuingHouse
		if err := db.First(&i, id).Error; err != nil {
			return nil, err
		}
		return i, nil
	case "approved_asset_custodian":
		var a models.ApprovedAssetCustodian
		if err := db.First(&a, id).Error; err != nil {
			return nil, err
		}
		return a, nil
	case "legal_adviser":
		var l models.LegalAdviser
		if err := db.First(&l, id).Error; err != nil {
			return nil, err
		}
		return l, nil
	case "financial_adviser":
		var f models.FinancialAdviser
		if err := db.First(&f, id).Error; err != nil {
			return nil, err
		}
		return f, nil
	default:
		return nil, errors.New("invalid type specified")
	}
}

func handleAssetManager(req models.SavePartnerRequest, db *gorm.DB) error {
	var entity models.AssetManager
	payloadBytes, err := json.Marshal(req.Payload)
	if err != nil {
		return errors.New("failed to marshal payload")
	}
	if err = json.Unmarshal(payloadBytes, &entity); err != nil {
		return errors.New("invalid payload for asset_manager")
	}
	if req.Action == "create" {
		return db.Create(&entity).Error
	} else if req.Action == "update" {
		return db.Model(&models.AssetManager{}).Where("id = ?", entity.ID).Updates(entity).Error
	}
	return errors.New("unsupported action")
}

func handleAssetIssuingHouse(req models.SavePartnerRequest, db *gorm.DB) error {
	var entity models.AssetIssuingHouse
	payloadBytes, err := json.Marshal(req.Payload)
	if err != nil {
		return errors.New("failed to marshal payload")
	}
	if err = json.Unmarshal(payloadBytes, &entity); err != nil {
		return errors.New("invalid payload for asset_issuing_house")
	}
	if req.Action == "create" {
		return db.Create(&entity).Error
	} else if req.Action == "update" {
		return db.Model(&models.AssetIssuingHouse{}).Where("id = ?", entity.ID).Updates(entity).Error
	}
	return errors.New("unsupported action")
}

func handleAssetCustodian(req models.SavePartnerRequest, db *gorm.DB) error {
	var entity models.ApprovedAssetCustodian
	payloadBytes, err := json.Marshal(req.Payload)
	if err != nil {
		return errors.New("failed to marshal payload")
	}
	if err := json.Unmarshal(payloadBytes, &entity); err != nil {
		return errors.New("invalid payload for approved_asset_custodian")
	}
	if req.Action == "create" {
		return db.Create(&entity).Error
	} else if req.Action == "update" {
		return db.Model(&models.ApprovedAssetCustodian{}).Where("id = ?", entity.ID).Updates(entity).Error
	}
	return errors.New("unsupported action")
}

type Country struct {
	// country_code is the primary key. The wallet-owned `countries` table has
	// no `id` column - see 000001_baseline.up.sql, which ends with
	// `ADD CONSTRAINT countries_pkey PRIMARY KEY (country_code)`. Declaring an
	// autoIncrement `id` here made GORM select a column that does not exist,
	// so every read of this table returned 500.
	CountryCode                              string     `gorm:"column:country_code;type:varchar(2);primaryKey" json:"country_code"`
	SecTokenizationFeePercent                float64    `gorm:"column:sec_tokenization_fee_percent;type:numeric" json:"sec_tokenization_fee_percent"`
	SecTradeFeePercent                       float64    `gorm:"column:sec_trade_fee_percent;type:numeric" json:"sec_trade_fee_percent"`
	RegionName                               string     `gorm:"column:region_name;type:text" json:"region_name"`
	SecTokenizationFee                       float64    `gorm:"column:sec_tokenization_fee;type:numeric" json:"sec_tokenization_fee"`
	SecTradeFee                              float64    `gorm:"column:sec_trade_fee;type:numeric" json:"sec_trade_fee"`
	MinTokenizationFee                       float64    `gorm:"column:min_tokenization_fee;type:numeric" json:"min_tokenization_fee"`
	TokenizationApplicationFee               float64    `gorm:"column:tokenization_application_fee;type:numeric" json:"tokenization_application_fee"`
	TokenizationApplicationFeeAsset          string     `gorm:"column:tokenization_application_fee_asset;type:text" json:"tokenization_application_fee_asset"`
	CountryName                              string     `gorm:"column:country_name;type:text" json:"country_name"`
	QuoteCurrencyCode                        string     `gorm:"column:quote_currency_code;type:text" json:"quote_currency_code"`
	LegalAndProfessionalFeePercent           float64    `gorm:"column:legal_and_professional_fee_percent;type:numeric" json:"legal_and_professional_fee_percent"`
	RatingAgencyFeePercent                   float64    `gorm:"column:rating_agency_fee_percent;type:numeric" json:"rating_agency_fee_percent"`
	VATPercent                               float64    `gorm:"column:vat_percent;type:numeric" json:"vat_percent"`
	RegulatorName                            string     `gorm:"column:regulator_name;type:text" json:"regulator_name"`
	SecTokenizationFeeFixed                  float64    `gorm:"column:sec_tokenization_fee_fixed;type:numeric" json:"sec_tokenization_fee_fixed"`
	SecTradeFeeFixed                         float64    `gorm:"column:sec_trade_fee_fixed;type:numeric" json:"sec_trade_fee_fixed"`
	LegalAndProfessionalFeeFixed             float64    `gorm:"column:legal_and_professional_fee_fixed;type:numeric" json:"legal_and_professional_fee_fixed"`
	RatingAgencyFeeFixed                     float64    `gorm:"column:rating_agency_fee_fixed;type:numeric" json:"rating_agency_fee_fixed"`
	FiatLabel                                string     `gorm:"column:fiat_label;type:text" json:"fiat_label"`
	FiatGlyph                                string     `gorm:"column:fiat_glyph;type:text" json:"fiat_glyph"`
	CreatedAt                                time.Time  `gorm:"column:created_at;type:timestamp with time zone" json:"created_at"`
	UpdatedAt                                time.Time  `gorm:"column:updated_at;type:timestamp with time zone" json:"updated_at"`
	DeletedAt                                *time.Time `gorm:"column:deleted_at;type:timestamp with time zone" json:"deleted_at,omitempty"`
	MinTrovBalanceForTokenizationApplication float64    `gorm:"column:min_trov_balance_for_tokenization_application;type:numeric" json:"min_trov_balance_for_tokenization_application"`
	FiatActivationAmount                     float64    `gorm:"column:fiat_activation_amount;type:numeric" json:"fiat_activation_amount"`
	TrovTokenActivationPercent               float64    `gorm:"column:trov_token_activation_percent;type:numeric" json:"trov_token_activation_percent"`
}

// insert query to update the country table which already exist
// An update query will be executed if the country code already exisT (id,created_at)t

func SaveCountry(db *gorm.DB, req models.CountryRequest) error {
	country := Country{
		CountryCode:                              req.CountryCode,
		RegionName:                               req.RegionName,
		CountryName:                              req.CountryName,
		QuoteCurrencyCode:                        req.QuoteCurrencyCode,
		RegulatorName:                            req.RegulatorName,
		FiatLabel:                                req.FiatLabel,
		FiatGlyph:                                req.FiatGlyph,
		SecTokenizationFeePercent:                req.SecTokenizationFeePercent,
		SecTradeFeePercent:                       req.SecTradeFeePercent,
		SecTokenizationFee:                       req.SecTokenizationFee,
		SecTradeFee:                              req.SecTradeFee,
		MinTokenizationFee:                       req.MinTokenizationFee,
		TokenizationApplicationFee:               req.TokenizationApplicationFee,
		TokenizationApplicationFeeAsset:          req.TokenizationApplicationFeeAsset,
		LegalAndProfessionalFeePercent:           req.LegalAndProfessionalFeePercent,
		RatingAgencyFeePercent:                   req.RatingAgencyFeePercent,
		VATPercent:                               req.VATPercent,
		SecTokenizationFeeFixed:                  req.SecTokenizationFeeFixed,
		SecTradeFeeFixed:                         req.SecTradeFeeFixed,
		LegalAndProfessionalFeeFixed:             req.LegalAndProfessionalFeeFixed,
		RatingAgencyFeeFixed:                     req.RatingAgencyFeeFixed,
		MinTrovBalanceForTokenizationApplication: req.MinTrovBalanceForTokenizationApplication,
		FiatActivationAmount:                     req.FiatActivationAmount,
		TrovTokenActivationPercent:               req.TrovTokenActivationPercent,
	}

	if req.Action == "create" {
		return db.Create(&country).Error
	}

	// Keyed on country_code, the table's actual primary key.
	return db.Model(&Country{}).Where("country_code = ?", req.CountryCode).Updates(country).Error
}

func ListCountries(db *gorm.DB) ([]models.CountryResponse, error) {
	var countries []Country
	if err := db.Find(&countries).Error; err != nil {
		return nil, err
	}

	var result []models.CountryResponse
	for _, c := range countries {
		result = append(result, models.CountryResponse{
			CountryCode:                              c.CountryCode,
			RegionName:                               c.RegionName,
			CountryName:                              c.CountryName,
			QuoteCurrencyCode:                        c.QuoteCurrencyCode,
			RegulatorName:                            c.RegulatorName,
			FiatLabel:                                c.FiatLabel,
			FiatGlyph:                                c.FiatGlyph,
			SecTokenizationFeePercent:                c.SecTokenizationFeePercent,
			SecTradeFeePercent:                       c.SecTradeFeePercent,
			SecTokenizationFee:                       c.SecTokenizationFee,
			SecTradeFee:                              c.SecTradeFee,
			MinTokenizationFee:                       c.MinTokenizationFee,
			TokenizationApplicationFee:               c.TokenizationApplicationFee,
			TokenizationApplicationFeeAsset:          c.TokenizationApplicationFeeAsset,
			LegalAndProfessionalFeePercent:           c.LegalAndProfessionalFeePercent,
			RatingAgencyFeePercent:                   c.RatingAgencyFeePercent,
			VATPercent:                               c.VATPercent,
			SecTokenizationFeeFixed:                  c.SecTokenizationFeeFixed,
			SecTradeFeeFixed:                         c.SecTradeFeeFixed,
			LegalAndProfessionalFeeFixed:             c.LegalAndProfessionalFeeFixed,
			RatingAgencyFeeFixed:                     c.RatingAgencyFeeFixed,
			MinTrovBalanceForTokenizationApplication: c.MinTrovBalanceForTokenizationApplication,
			FiatActivationAmount:                     c.FiatActivationAmount,
			TrovTokenActivationPercent:               c.TrovTokenActivationPercent,
		})
	}
	return result, nil
}

// DeleteCountry removes a country by its country_code (the primary key). The
// parameter is still named id because it arrives as the {id} path segment.
func DeleteCountry(db *gorm.DB, id string) error {
	return db.Where("country_code = ?", id).Delete(&Country{}).Error
}

func GetCountryByID(db *gorm.DB, id string) (*models.CountryResponse, error) {
	var country Country
	if err := db.First(&country, "country_code = ?", id).Error; err != nil {
		return nil, err
	}

	response := &models.CountryResponse{
		CountryCode:                     country.CountryCode,
		RegionName:                      country.RegionName,
		CountryName:                     country.CountryName,
		QuoteCurrencyCode:               country.QuoteCurrencyCode,
		RegulatorName:                   country.RegulatorName,
		FiatLabel:                       country.FiatLabel,
		FiatGlyph:                       country.FiatGlyph,
		SecTokenizationFeePercent:       country.SecTokenizationFeePercent,
		SecTradeFeePercent:              country.SecTradeFeePercent,
		SecTokenizationFee:              country.SecTokenizationFee,
		SecTradeFee:                     country.SecTradeFee,
		MinTokenizationFee:              country.MinTokenizationFee,
		TokenizationApplicationFee:      country.TokenizationApplicationFee,
		TokenizationApplicationFeeAsset: country.TokenizationApplicationFeeAsset,
		LegalAndProfessionalFeePercent:  country.LegalAndProfessionalFeePercent,
		RatingAgencyFeePercent:          country.RatingAgencyFeePercent,
		VATPercent:                      country.VATPercent,
		SecTokenizationFeeFixed:         country.SecTokenizationFeeFixed,
		SecTradeFeeFixed:                country.SecTradeFeeFixed,
		LegalAndProfessionalFeeFixed:    country.LegalAndProfessionalFeeFixed,
		RatingAgencyFeeFixed:            country.RatingAgencyFeeFixed,
	}

	return response, nil
}

func SaveAssetIssuingHouse(action string, req models.AssetIssuingHouseRequest, db *gorm.DB) (uint64, error) {
	if action == "create" {
		entity := models.AssetIssuingHouse{
			AssetIssuingHouseName:    req.AssetIssuingHouseName,
			AssetIssuingHouseAddress: req.AssetIssuingHouseAddress,
			AssetIssuingHouseCountry: req.AssetIssuingHouseCountry,
			FeePercent:               req.FeePercent,
			FeeFixed:                 req.FeeFixed,
		}
		if err := db.Create(&entity).Error; err != nil {
			return 0, err
		}
		return uint64(entity.ID), nil
	} else if action == "update" {
		if req.ID == 0 {
			return 0, errors.New("missing id in request body for update")
		}
		var existing models.AssetIssuingHouse
		if err := db.First(&existing, req.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return 0, fmt.Errorf("asset issuing house with ID %d not found for update", req.ID)
			}
			return 0, err
		}

		updateData := models.AssetIssuingHouse{
			AssetIssuingHouseName:    req.AssetIssuingHouseName,
			AssetIssuingHouseAddress: req.AssetIssuingHouseAddress,
			AssetIssuingHouseCountry: req.AssetIssuingHouseCountry,
			FeePercent:               req.FeePercent,
			FeeFixed:                 req.FeeFixed,
		}
		if err := db.Model(&models.AssetIssuingHouse{}).Where("id = ?", req.ID).Updates(updateData).Error; err != nil {
			return 0, err
		}
		return uint64(req.ID), nil
	}
	return 0, fmt.Errorf("invalid action specified: %s", action)
}

func SaveApprovedAssetCustodian(action string, req models.ApprovedAssetCustodianRequest, db *gorm.DB) (uint64, error) {
	// Convert string to float64
	feePercent, err := strconv.ParseFloat(req.FeePercent, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid fee percent: %v", err)
	}
	feeFixed, err := strconv.ParseFloat(req.FeeFixed, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid fee fixed: %v", err)
	}
	if action == "create" {
		entity := models.ApprovedAssetCustodian{
			AssetCustodianName:    req.AssetCustodianName,
			AssetCustodianAddress: req.AssetCustodianAddress,
			AssetCustodianCountry: req.AssetCustodianCountry,
			RequirementDocument:   req.RequirementDocument,
			FeePercent:            feePercent,
			FeeFixed:              feeFixed,
		}
		if err := db.Create(&entity).Error; err != nil {
			return 0, err
		}
		return uint64(entity.ID), nil
	} else if action == "update" {
		if req.ID == 0 {
			return 0, errors.New("missing id in request body for update")
		}
		var existing models.ApprovedAssetCustodian
		if err := db.First(&existing, req.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return 0, fmt.Errorf("approved asset custodian with ID %d not found for update", req.ID)
			}
			return 0, err
		}

		updateData := models.ApprovedAssetCustodian{
			AssetCustodianName:    req.AssetCustodianName,
			AssetCustodianAddress: req.AssetCustodianAddress,
			AssetCustodianCountry: req.AssetCustodianCountry,
			RequirementDocument:   req.RequirementDocument,
			FeePercent:            feePercent,
			FeeFixed:              feeFixed,
		}
		if err := db.Model(&models.ApprovedAssetCustodian{}).Where("id = ?", req.ID).Updates(updateData).Error; err != nil {
			return 0, err
		}
		return uint64(req.ID), nil
	}
	return 0, fmt.Errorf("invalid action specified: %s", action)
}

func SaveLegalAndProfessionals(action string, req models.LegalAndProfessionalsRequest, db *gorm.DB) (uint64, error) {
	// Convert string to float64
	feePercent, err := strconv.ParseFloat(req.FeePercent, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid fee percent: %v", err)
	}
	feeFixed, err := strconv.ParseFloat(req.FeeFixed, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid fee fixed: %v", err)
	}

	if action == "create" {
		entity := models.LegalAndProfessionalsPartners{
			PartnerName:    req.PartnerName,
			PartnerAddress: req.PartnerAddress,
			PartnerCountry: req.PartnerCountry,
			FeePercent:     feePercent,
			FeeFixed:       feeFixed,
		}
		if err := db.Create(&entity).Error; err != nil {
			return 0, err
		}
		return uint64(entity.ID), nil
	} else if action == "update" {
		if req.ID <= 0 {
			return 0, errors.New("missing id in request body for update")
		}
		var existing models.LegalAndProfessionalsPartners
		if err := db.First(&existing, req.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return 0, fmt.Errorf("legal and professional with ID %d not found for update", req.ID)
			}
			return 0, err
		}

		updateData := models.LegalAndProfessionalsPartners{
			PartnerName:    req.PartnerName,
			PartnerAddress: req.PartnerAddress,
			PartnerCountry: req.PartnerCountry,
			FeePercent:     feePercent,
			FeeFixed:       feeFixed,
		}
		if err := db.Model(&models.LegalAndProfessionalsPartners{}).Where("id = ?", req.ID).Updates(updateData).Error; err != nil {
			return 0, err
		}
		return uint64(req.ID), nil
	}
	return 0, fmt.Errorf("invalid action specified: %s", action)
}

// SaveRatingAgency saves or updates a rating agency in the database
func SaveRatingAgency(action string, req models.RatingAgencyRequest, db *gorm.DB) (uint64, error) {
	// Convert string fees to float64
	feePercent, err := strconv.ParseFloat(req.FeePercent, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid fee percent: %v", err)
	}
	feeFixed, err := strconv.ParseFloat(req.FeeFixed, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid fee fixed: %v", err)
	}

	if action == "create" {
		entity := models.RatingAgency{
			AgencyName:    req.AgencyName,
			AgencyAddress: req.AgencyAddress,
			AgencyCountry: req.AgencyCountry,
			FeePercent:    feePercent,
			FeeFixed:      feeFixed,
		}
		if err := db.Table("rating_agencies").Create(&entity).Error; err != nil {
			return 0, err
		}
		return uint64(entity.ID), nil
	} else if action == "update" {
		if req.ID <= 0 {
			return 0, errors.New("missing id in request body for update")
		}
		var existing models.RatingAgency
		if err := db.Table("rating_agencies").First(&existing, req.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return 0, fmt.Errorf("rating agency with ID %d not found for update", req.ID)
			}
			return 0, err
		}

		updateData := models.RatingAgency{
			AgencyName:    req.AgencyName,
			AgencyAddress: req.AgencyAddress,
			AgencyCountry: req.AgencyCountry,
			FeePercent:    feePercent,
			FeeFixed:      feeFixed,
		}
		if err := db.Table("rating_agencies").Where("id = ?", req.ID).Updates(updateData).Error; err != nil {
			return 0, err
		}
		return uint64(req.ID), nil
	}
	return 0, fmt.Errorf("invalid action specified: %s", action)
}

func ListAllPartnerEntities(db *gorm.DB) (map[string]interface{}, error) {
	results := make(map[string]interface{})
	var err error

	var managers []models.AssetManager
	if err = db.Find(&managers).Error; err != nil {
		log.Printf("[DB_LIST_ALL] Error fetching asset managers: %v", err)
		return nil, fmt.Errorf("failed to fetch asset managers: %w", err)
	}
	results["asset_manager"] = managers

	var houses []models.AssetIssuingHouse
	if err = db.Find(&houses).Error; err != nil {
		log.Printf("[DB_LIST_ALL] Error fetching asset issuing houses: %v", err)
		return nil, fmt.Errorf("failed to fetch asset issuing houses: %w", err)
	}
	results["asset_issuing_house"] = houses

	var custodians []models.ApprovedAssetCustodian
	if err = db.Find(&custodians).Error; err != nil {
		log.Printf("[DB_LIST_ALL] Error fetching approved asset custodians: %v", err)
		return nil, fmt.Errorf("failed to fetch approved asset custodians: %w", err)
	}
	results["approved_asset_custodian"] = custodians
	var legalPartners []models.LegalAndProfessionalsPartners
	if err = db.Find(&legalPartners).Error; err != nil {
		log.Printf("[DB_LIST_ALL] Error fetching legal and professional partners: %v", err)
		return nil, fmt.Errorf("failed to fetch legal and professional partners: %w", err)
	}
	results["legal_and_professionals"] = legalPartners
	var ratingAgencies []models.RatingAgency
	if err = db.Find(&ratingAgencies).Error; err != nil {
		log.Printf("[DB_LIST_ALL] Error fetching rating agencies: %v", err)
		return nil, fmt.Errorf("failed to fetch rating agencies: %w", err)
	}
	results["rating_agency"] = ratingAgencies
	var trustees []models.Trustees
	if err = db.Find(&trustees).Error; err != nil {
		log.Printf("[DB_LIST_ALL] Error fetching trustees: %v", err)
		return nil, fmt.Errorf("failed to fetch trustees: %w", err)
	}
	results["trustees"] = trustees

	var legalAdvisers []models.LegalAdviser
	if err = db.Find(&legalAdvisers).Error; err != nil {
		log.Printf("[DB_LIST_ALL] Error fetching legal advisers: %v", err)
		return nil, fmt.Errorf("failed to fetch legal advisers: %w", err)
	}
	results["legal_adviser"] = legalAdvisers

	var financialAdvisers []models.FinancialAdviser
	if err = db.Find(&financialAdvisers).Error; err != nil {
		log.Printf("[DB_LIST_ALL] Error fetching financial advisers: %v", err)
		return nil, fmt.Errorf("failed to fetch financial advisers: %w", err)
	}
	results["financial_adviser"] = financialAdvisers

	return results, nil
}

// SaveLegalAdviser saves or updates a legal adviser in the database
func SaveLegalAdviser(action string, req models.LegalAdviserRequest, db *gorm.DB) (uint64, error) {
	feePercent, err := strconv.ParseFloat(req.FeePercent, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid fee percent: %v", err)
	}
	feeFixed, err := strconv.ParseFloat(req.FeeFixed, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid fee fixed: %v", err)
	}

	if action == "create" {
		entity := models.LegalAdviser{
			AdviserName:    req.AdviserName,
			AdviserAddress: req.AdviserAddress,
			AdviserCountry: req.AdviserCountry,
			FeePercent:     feePercent,
			FeeFixed:       feeFixed,
		}
		if err := db.Create(&entity).Error; err != nil {
			return 0, err
		}
		return uint64(entity.ID), nil
	} else if action == "update" {
		if req.ID <= 0 {
			return 0, errors.New("missing id in request body for update")
		}
		var existing models.LegalAdviser
		if err := db.First(&existing, req.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return 0, fmt.Errorf("legal adviser with ID %d not found for update", req.ID)
			}
			return 0, err
		}
		updateData := models.LegalAdviser{
			AdviserName:    req.AdviserName,
			AdviserAddress: req.AdviserAddress,
			AdviserCountry: req.AdviserCountry,
			FeePercent:     feePercent,
			FeeFixed:       feeFixed,
		}
		if err := db.Model(&models.LegalAdviser{}).Where("id = ?", req.ID).Updates(updateData).Error; err != nil {
			return 0, err
		}
		return uint64(req.ID), nil
	}
	return 0, fmt.Errorf("invalid action specified: %s", action)
}

// SaveFinancialAdviser saves or updates a financial adviser in the database
func SaveFinancialAdviser(action string, req models.FinancialAdviserRequest, db *gorm.DB) (uint64, error) {
	feePercent, err := strconv.ParseFloat(req.FeePercent, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid fee percent: %v", err)
	}
	feeFixed, err := strconv.ParseFloat(req.FeeFixed, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid fee fixed: %v", err)
	}

	if action == "create" {
		entity := models.FinancialAdviser{
			AdviserName:    req.AdviserName,
			AdviserAddress: req.AdviserAddress,
			AdviserCountry: req.AdviserCountry,
			FeePercent:     feePercent,
			FeeFixed:       feeFixed,
		}
		if err := db.Create(&entity).Error; err != nil {
			return 0, err
		}
		return uint64(entity.ID), nil
	} else if action == "update" {
		if req.ID <= 0 {
			return 0, errors.New("missing id in request body for update")
		}
		var existing models.FinancialAdviser
		if err := db.First(&existing, req.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return 0, fmt.Errorf("financial adviser with ID %d not found for update", req.ID)
			}
			return 0, err
		}
		updateData := models.FinancialAdviser{
			AdviserName:    req.AdviserName,
			AdviserAddress: req.AdviserAddress,
			AdviserCountry: req.AdviserCountry,
			FeePercent:     feePercent,
			FeeFixed:       feeFixed,
		}
		if err := db.Model(&models.FinancialAdviser{}).Where("id = ?", req.ID).Updates(updateData).Error; err != nil {
			return 0, err
		}
		return uint64(req.ID), nil
	}
	return 0, fmt.Errorf("invalid action specified: %s", action)
}

func SaveTrustees(action string, req models.TrusteesRequest, db *gorm.DB) (uint64, error) {
	feePercent, err := strconv.ParseFloat(req.FeePercent, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid fee percent: %v", err)
	}
	feeFixed, err := strconv.ParseFloat(req.FeeFixed, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid fee fixed: %v", err)
	}
	if action == "create" {
		entity := models.Trustees{
			TrusteeName:    req.TrusteeName,
			TrusteeAddress: req.TrusteeAddress,
			TrusteeCountry: req.TrusteeCountry,
			FeePercent:     feePercent,
			FeeFixed:       feeFixed,
		}
		if err := db.Table("trustees").Create(&entity).Error; err != nil {
			return 0, err
		}
		return uint64(entity.ID), nil
	} else if action == "update" {
		if req.ID == 0 {
			return 0, errors.New("missing id in request body for update")
		}
		var existing models.Trustees
		if err := db.Table("trustees").First(&existing, req.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return 0, fmt.Errorf("trustees with ID %d not found for update", req.ID)
			}
			return 0, err
		}
		updateData := models.Trustees{
			TrusteeName:    req.TrusteeName,
			TrusteeAddress: req.TrusteeAddress,
			TrusteeCountry: req.TrusteeCountry,
			FeePercent:     feePercent,
			FeeFixed:       feeFixed,
		}
		if err := db.Table("trustees").Where("id = ?", req.ID).Updates(updateData).Error; err != nil {
			return 0, err
		}
		return uint64(req.ID), nil
	}
	return 0, fmt.Errorf("invalid action specified: %s", action)
}

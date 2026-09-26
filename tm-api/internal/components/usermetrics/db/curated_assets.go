package usermetrics

import (
	"admin-panel-dashboard/internal/models"
	"errors"

	"gorm.io/gorm"
)

func boolToInt(b bool) uint64 {
	if b {
		return 1
	}
	return 0
}

// ListCuratedAssets paginates/filters the curated-asset catalog for the
// admin management table.
func ListCuratedAssets(db *gorm.DB, req models.CuratedAssetListRequest) ([]models.CuratedAsset, int64, error) {
	query := db.Model(&models.CuratedAsset{})
	if req.AssetCode != "" {
		query = query.Where("LOWER(asset_code) LIKE LOWER(?)", "%"+req.AssetCode+"%")
	}
	if req.AssetClassID != 0 {
		query = query.Where("asset_class_id = ?", req.AssetClassID)
	}
	if req.P2PEnabled != nil {
		query = query.Where("p2p_enabled = ?", *req.P2PEnabled)
	}
	if req.Inactive != nil {
		query = query.Where("inactive = ?", boolToInt(*req.Inactive))
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var assets []models.CuratedAsset
	offset := (req.Page - 1) * req.PageSize
	if err := query.Order("priority, asset_code").Offset(offset).Limit(req.PageSize).Find(&assets).Error; err != nil {
		return nil, 0, err
	}
	return assets, total, nil
}

// GetCuratedAssetByID fetches a single curated asset.
func GetCuratedAssetByID(db *gorm.DB, id uint64) (models.CuratedAsset, error) {
	var asset models.CuratedAsset
	if err := db.First(&asset, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return asset, errors.New("curated asset not found")
		}
		return asset, err
	}
	return asset, nil
}

// SaveCuratedAsset creates or updates a curated asset. The update path
// writes through a map rather than a struct: GORM's struct-based Updates
// silently skips Go zero values (false, 0, "") entirely, independent of
// any `default:` column tag - it would make turning P2PEnabled/
// Withdrawable/Inactive back off (or clearing a text field) a silent
// no-op instead of the write the admin asked for.
func SaveCuratedAsset(db *gorm.DB, req models.CuratedAssetRequest) (models.CuratedAsset, error) {
	if req.Action == "create" {
		asset := models.CuratedAsset{
			AssetCode:                   req.AssetCode,
			AssetName:                   req.AssetName,
			ContractAddress:             req.ContractAddress,
			Description:                 req.Description,
			Website:                     req.Website,
			AssetConditions:             req.AssetConditions,
			AssetLimit:                  req.AssetLimit,
			AssetRedemptionInstructions: req.AssetRedemptionInstructions,
			ContactEmail:                req.ContactEmail,
			Priority:                    req.Priority,
			AssetClassID:                req.AssetClassID,
			Organization:                req.Organization,
			Withdrawable:                boolToInt(req.Withdrawable),
			GenerateDepositAddress:      boolToInt(req.GenerateDepositAddress),
			DecimalPlaces:               req.DecimalPlaces,
			Inactive:                    boolToInt(req.Inactive),
			P2PEnabled:                  req.P2PEnabled,
		}
		if req.ImageURL != "" {
			asset.ImageURL = &req.ImageURL
		}
		if req.RealAssetImageURL != "" {
			asset.RealAssetImageURL = &req.RealAssetImageURL
		}
		if req.ClosedGroup != "" {
			asset.ClosedGroup = &req.ClosedGroup
		}
		if err := db.Create(&asset).Error; err != nil {
			return asset, err
		}
		return asset, nil
	}

	if req.Action == "update" {
		if req.ID == 0 {
			return models.CuratedAsset{}, errors.New("id is required for update")
		}
		existing, err := GetCuratedAssetByID(db, req.ID)
		if err != nil {
			return existing, err
		}
		updates := map[string]interface{}{
			"asset_name":                     req.AssetName,
			"contract_address":               req.ContractAddress,
			"description":                    req.Description,
			"image_url":                      nullableString(req.ImageURL),
			"website":                        req.Website,
			"asset_conditions":               req.AssetConditions,
			"asset_limit":                    req.AssetLimit,
			"asset_redemption_instructions":  req.AssetRedemptionInstructions,
			"contact_email":                  req.ContactEmail,
			"priority":                       req.Priority,
			"asset_class_id":                 req.AssetClassID,
			"organization":                   req.Organization,
			"withdrawable":                   boolToInt(req.Withdrawable),
			"generate_deposit_address":       boolToInt(req.GenerateDepositAddress),
			"decimal_places":                 req.DecimalPlaces,
			"real_asset_image_url":           nullableString(req.RealAssetImageURL),
			"inactive":                       boolToInt(req.Inactive),
			"closed_group":                   nullableString(req.ClosedGroup),
			"p2p_enabled":                    req.P2PEnabled,
		}
		if err := db.Model(&models.CuratedAsset{}).Where("id = ?", req.ID).Updates(updates).Error; err != nil {
			return existing, err
		}
		return GetCuratedAssetByID(db, req.ID)
	}

	return models.CuratedAsset{}, errors.New("invalid action: must be 'create' or 'update'")
}

func nullableString(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}

// SetCuratedAssetP2PEnabled is the dedicated toggle for the P2P-availability
// flag - a map-based single-column update for the same false-gets-dropped
// reason SaveCuratedAsset's update path uses one.
func SetCuratedAssetP2PEnabled(db *gorm.DB, id uint64, enabled bool) (models.CuratedAsset, error) {
	if _, err := GetCuratedAssetByID(db, id); err != nil {
		return models.CuratedAsset{}, err
	}
	if err := db.Model(&models.CuratedAsset{}).Where("id = ?", id).Update("p2p_enabled", enabled).Error; err != nil {
		return models.CuratedAsset{}, err
	}
	return GetCuratedAssetByID(db, id)
}

// SetCuratedAssetInactive is the dedicated toggle for retiring/restoring a
// curated asset without deleting its row - curated assets are referenced by
// historical wallet/offer/order data, so admins retire one by setting this
// flag rather than deleting it.
func SetCuratedAssetInactive(db *gorm.DB, id uint64, inactive bool) (models.CuratedAsset, error) {
	if _, err := GetCuratedAssetByID(db, id); err != nil {
		return models.CuratedAsset{}, err
	}
	if err := db.Model(&models.CuratedAsset{}).Where("id = ?", id).Update("inactive", boolToInt(inactive)).Error; err != nil {
		return models.CuratedAsset{}, err
	}
	return GetCuratedAssetByID(db, id)
}

// ListAssetClasses returns the reference list of asset categories, for the
// admin form's category dropdown.
func ListAssetClasses(db *gorm.DB) ([]models.AssetClass, error) {
	var classes []models.AssetClass
	if err := db.Order("asset_class").Find(&classes).Error; err != nil {
		return nil, err
	}
	return classes, nil
}

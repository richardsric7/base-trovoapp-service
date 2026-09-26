package usermetrics

import (
	"admin-panel-dashboard/internal/models"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// intFlag mirrors boolToInt (curated_assets.go, same package) but returns
// plain int - ServiceLink's permission/status columns are int in
// app-backend, unlike CuratedAsset's uint64 ones.
func intFlag(b bool) int {
	if b {
		return 1
	}
	return 0
}

// ListServiceLinks paginates/filters the service-link (white-label
// partner integration) catalog for the admin management table.
func ListServiceLinks(db *gorm.DB, req models.ServiceLinkListRequest) ([]models.ServiceLink, int64, error) {
	query := db.Model(&models.ServiceLink{})
	if req.OwnerUsername != "" {
		query = query.Where("LOWER(owner_username) LIKE LOWER(?)", "%"+req.OwnerUsername+"%")
	}
	if req.ShortName != "" {
		query = query.Where("LOWER(short_name) LIKE LOWER(?)", "%"+req.ShortName+"%")
	}
	if req.Inactive != nil {
		query = query.Where("inactive = ?", intFlag(*req.Inactive))
	}
	if req.Verified != nil {
		query = query.Where("verified = ?", intFlag(*req.Verified))
	}
	if req.Suspended != nil {
		query = query.Where("suspended = ?", intFlag(*req.Suspended))
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var links []models.ServiceLink
	offset := (req.Page - 1) * req.PageSize
	if err := query.Order("created_at desc").Offset(offset).Limit(req.PageSize).Find(&links).Error; err != nil {
		return nil, 0, err
	}
	return links, total, nil
}

// GetServiceLinkByID fetches a single service link.
func GetServiceLinkByID(db *gorm.DB, id string) (models.ServiceLink, error) {
	var link models.ServiceLink
	if err := db.First(&link, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return link, errors.New("service link not found")
		}
		return link, err
	}
	return link, nil
}

// GetServiceLinkOwnerSummary looks up the handful of user-table fields
// worth showing alongside a service link (there's no FK, so this is a
// best-effort lookup by username - a nil return means the owner username
// no longer resolves to a real user, which the caller should surface
// rather than fail on).
func GetServiceLinkOwnerSummary(db *gorm.DB, username string) (*models.ServiceLinkOwnerSummary, error) {
	var owner models.ServiceLinkOwnerSummary
	err := db.Table("users").
		Select("email, address, kyc_verified, is_merchant, merchant_online, corporate, membership_type, verified, suspended").
		Where("username = ?", username).
		Take(&owner).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &owner, nil
}

// AttachServiceLinkOwners resolves each link's owner summary in one pass.
func AttachServiceLinkOwners(db *gorm.DB, links []models.ServiceLink) ([]models.ServiceLinkWithOwner, error) {
	out := make([]models.ServiceLinkWithOwner, 0, len(links))
	for _, link := range links {
		owner, err := GetServiceLinkOwnerSummary(db, link.OwnerUsername)
		if err != nil {
			return nil, err
		}
		out = append(out, models.ServiceLinkWithOwner{ServiceLink: link, Owner: owner})
	}
	return out, nil
}

// SaveServiceLink creates or updates a service link. Create looks up the
// owner user by username (there's no DB-level FK, so this is the only
// validation that the account exists) and copies their wallet address
// onto the new row; ID and ApiKey are always server-generated GUIDs,
// never accepted from the client. Update writes through a map, not a
// struct, for the same reason SaveCuratedAsset does: GORM's struct-based
// Updates silently drops a permission flag being turned back off.
func SaveServiceLink(db *gorm.DB, req models.ServiceLinkRequest) (models.ServiceLink, error) {
	if req.Action == "create" {
		var owner struct {
			Address string
		}
		if err := db.Table("users").Select("address").Where("username = ?", req.OwnerUsername).Take(&owner).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return models.ServiceLink{}, errors.New("owner username does not exist")
			}
			return models.ServiceLink{}, err
		}

		link := models.ServiceLink{
			ID:                                    uuid.NewString(),
			ApiKey:                                uuid.NewString(),
			OwnerUsername:                         req.OwnerUsername,
			Address:                               owner.Address,
			ShortName:                             req.ShortName,
			LongName:                              req.LongName,
			LoginPermission:                       intFlag(req.LoginPermission),
			PaymentPermission:                     intFlag(req.PaymentPermission),
			TokenInfoPermission:                   intFlag(req.TokenInfoPermission),
			AuthorizationPermission:               intFlag(req.AuthorizationPermission),
			EventPermission:                       intFlag(req.EventPermission),
			AllowUserInfo:                         intFlag(req.AllowUserInfo),
			PushNotificationPermission:            intFlag(req.PushNotificationPermission),
			IncludePhoneNumbers:                   intFlag(req.IncludePhoneNumbers),
			IncludeUserBalances:                   intFlag(req.IncludeUserBalances),
			TokenizedAssetAuthorizationPermission: intFlag(req.TokenizedAssetAuthorizationPermission),
			CreateUsersPermission:                 intFlag(req.CreateUsersPermission),
			AllowReferralForRegisteredUsers:       intFlag(req.AllowReferralForRegisteredUsers),
			Verified:                              intFlag(req.Verified),
			Inactive:                              intFlag(req.Inactive),
		}
		if err := db.Create(&link).Error; err != nil {
			return link, err
		}
		return link, nil
	}

	if req.Action == "update" {
		if req.ID == "" {
			return models.ServiceLink{}, errors.New("id is required for update")
		}
		existing, err := GetServiceLinkByID(db, req.ID)
		if err != nil {
			return existing, err
		}
		updates := map[string]interface{}{
			"short_name":                               req.ShortName,
			"long_name":                                req.LongName,
			"login_permission":                         intFlag(req.LoginPermission),
			"payment_permission":                       intFlag(req.PaymentPermission),
			"token_info_permission":                    intFlag(req.TokenInfoPermission),
			"authorization_permission":                 intFlag(req.AuthorizationPermission),
			"event_permission":                         intFlag(req.EventPermission),
			"allow_user_info":                          intFlag(req.AllowUserInfo),
			"push_notification_permission":             intFlag(req.PushNotificationPermission),
			"include_phone_numbers":                    intFlag(req.IncludePhoneNumbers),
			"include_user_balances":                    intFlag(req.IncludeUserBalances),
			"tokenized_asset_authorization_permission": intFlag(req.TokenizedAssetAuthorizationPermission),
			"create_users_permission":                  intFlag(req.CreateUsersPermission),
			"allow_referral_for_registered_users":      intFlag(req.AllowReferralForRegisteredUsers),
			"verified":                                 intFlag(req.Verified),
			"inactive":                                 intFlag(req.Inactive),
		}
		if err := db.Model(&models.ServiceLink{}).Where("id = ?", req.ID).Updates(updates).Error; err != nil {
			return existing, err
		}
		return GetServiceLinkByID(db, req.ID)
	}

	return models.ServiceLink{}, errors.New("invalid action: must be 'create' or 'update'")
}

// SetServiceLinkInactive is the dedicated toggle for retiring/restoring a
// service link without touching any of its other fields or permissions.
func SetServiceLinkInactive(db *gorm.DB, id string, inactive bool) (models.ServiceLink, error) {
	if _, err := GetServiceLinkByID(db, id); err != nil {
		return models.ServiceLink{}, err
	}
	if err := db.Model(&models.ServiceLink{}).Where("id = ?", id).Update("inactive", intFlag(inactive)).Error; err != nil {
		return models.ServiceLink{}, err
	}
	return GetServiceLinkByID(db, id)
}

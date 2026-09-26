package p2p

import (
	p2pModels "trovo-wallet-api/internal/components/p2p/models"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CreateOfferInput is the client-supplied shape for creating an offer. The
// backend resolves and validates the curated asset itself (Plan Section 6) -
// the client cannot supply an arbitrary contract address.
type CreateOfferInput struct {
	OfferType          string
	Asset              string // assetCode
	PaymentMethod      p2pModels.PaymentMethod
	Country            string
	CountryCode        string
	Currency           string
	PriceType          string
	Price              string
	PriceMargin        string
	MinOrderAmount     string
	MaxOrderAmount     string
	AvailableLiquidity string
	Remark             string
}

// CreateOffer validates the input against the curated asset table and
// creates a DRAFT offer owned by merchantUsername/merchantUserID.
func CreateOffer(gc *sharedconfig.GlobalConfig, merchantUsername, merchantUserID string, in CreateOfferInput) (p2pModels.Offer, error) {
	if in.OfferType != p2pModels.OfferTypeBuy && in.OfferType != p2pModels.OfferTypeSell {
		return p2pModels.Offer{}, &tErrors.CustomError{Param: "offerType", Err: "error-invalid-offer-type", ErrMessage: "offerType must be BUY or SELL"}
	}

	curatedAsset, err := userModels.Currency(in.Asset).GetCurratedAsset(gc)
	if err != nil || curatedAsset.AssetCode == "" {
		return p2pModels.Offer{}, &tErrors.CustomError{Param: "asset", Err: "error-unsupported-asset", ErrMessage: "This asset is not supported on the P2P marketplace."}
	}

	if _, e := decimal.NewFromString(in.Price); e != nil {
		return p2pModels.Offer{}, &tErrors.CustomError{Param: "price", Err: "error-invalid-price", ErrMessage: "price must be a valid decimal number"}
	}
	minAmt, e := decimal.NewFromString(in.MinOrderAmount)
	if e != nil {
		return p2pModels.Offer{}, &tErrors.CustomError{Param: "minOrderAmount", Err: "error-invalid-min-order-amount", ErrMessage: "minOrderAmount must be a valid decimal number"}
	}
	maxAmt, e := decimal.NewFromString(in.MaxOrderAmount)
	if e != nil {
		return p2pModels.Offer{}, &tErrors.CustomError{Param: "maxOrderAmount", Err: "error-invalid-max-order-amount", ErrMessage: "maxOrderAmount must be a valid decimal number"}
	}
	if maxAmt.LessThan(minAmt) {
		return p2pModels.Offer{}, &tErrors.CustomError{Param: "maxOrderAmount", Err: "error-max-below-min", ErrMessage: "maxOrderAmount cannot be less than minOrderAmount"}
	}

	offer := p2pModels.Offer{
		ID:                 gc.GenerateUUIDString(),
		MerchantUsername:   merchantUsername,
		MerchantUserID:     merchantUserID,
		OfferType:          in.OfferType,
		Asset:              curatedAsset.AssetCode,
		ContractAddress:    curatedAsset.ContractAddress,
		PaymentMethod:      in.PaymentMethod,
		Country:            in.Country,
		CountryCode:        in.CountryCode,
		Currency:           in.Currency,
		PriceType:          orDefaultStr(in.PriceType, "FIXED"),
		Price:              in.Price,
		PriceMargin:        orDefaultStr(in.PriceMargin, "0"),
		MinOrderAmount:     in.MinOrderAmount,
		MaxOrderAmount:     in.MaxOrderAmount,
		AvailableLiquidity: orDefaultStr(in.AvailableLiquidity, "0"),
		ReservedLiquidity:  "0",
		Remark:             in.Remark,
		AvailabilityStatus: p2pModels.OfferAvailabilityOffline,
		Status:             p2pModels.OfferStatusDraft,
		Version:            1,
	}

	if err := gc.DB.Omit(clause.Associations).Create(&offer).Error; err != nil {
		return p2pModels.Offer{}, &tErrors.ErrorTemporaryServerError{}
	}
	RecordAuditEvent(gc, "", offer.ID, p2pModels.EventOfferCreated, merchantUserID, offer)
	return offer, nil
}

func orDefaultStr(v, def string) string {
	if v == "" {
		return def
	}
	return v
}

// ActivateOffer transitions a DRAFT/PAUSED offer to ACTIVE + ONLINE.
func ActivateOffer(gc *sharedconfig.GlobalConfig, offerID, merchantUserID string) (p2pModels.Offer, error) {
	offer, err := GetOfferByID(gc.DB, offerID)
	if err != nil {
		return offer, &tErrors.CustomError{Param: "offerId", Err: "error-offer-not-found", ErrMessage: "Offer not found"}
	}
	if offer.MerchantUserID != merchantUserID {
		return offer, &tErrors.CustomError{Param: "offerId", Err: "error-forbidden", ErrMessage: "You do not own this offer", Code: 403}
	}
	offer.Status = p2pModels.OfferStatusActive
	offer.AvailabilityStatus = p2pModels.OfferAvailabilityOnline
	if err := gc.DB.Save(&offer).Error; err != nil {
		return offer, &tErrors.ErrorTemporaryServerError{}
	}
	RecordAuditEvent(gc, "", offer.ID, p2pModels.EventOfferUpdated, merchantUserID, offer)
	return offer, nil
}

// PauseOffer transitions an offer to PAUSED + OFFLINE.
func PauseOffer(gc *sharedconfig.GlobalConfig, offerID, merchantUserID string) (p2pModels.Offer, error) {
	offer, err := GetOfferByID(gc.DB, offerID)
	if err != nil {
		return offer, &tErrors.CustomError{Param: "offerId", Err: "error-offer-not-found", ErrMessage: "Offer not found"}
	}
	if offer.MerchantUserID != merchantUserID {
		return offer, &tErrors.CustomError{Param: "offerId", Err: "error-forbidden", ErrMessage: "You do not own this offer", Code: 403}
	}
	offer.Status = p2pModels.OfferStatusPaused
	offer.AvailabilityStatus = p2pModels.OfferAvailabilityOffline
	if err := gc.DB.Save(&offer).Error; err != nil {
		return offer, &tErrors.ErrorTemporaryServerError{}
	}
	RecordAuditEvent(gc, "", offer.ID, p2pModels.EventOfferUpdated, merchantUserID, offer)
	return offer, nil
}

// UpdateOfferInput is the client-supplied shape for editing an existing
// offer. OfferType/Asset are intentionally not editable here - they drive
// the curated-asset validation and role mapping (Plan Section 11) done at
// creation time, so changing either is a new offer, not an edit.
type UpdateOfferInput struct {
	PaymentMethod      p2pModels.PaymentMethod
	PriceType          string
	Price              string
	PriceMargin        string
	MinOrderAmount     string
	MaxOrderAmount     string
	AvailableLiquidity string
	Remark             string
}

// UpdateOffer edits an existing offer's terms. Allowed from any
// non-terminal status (DRAFT/ACTIVE/PAUSED/OUT_OF_LIQUIDITY) - a CLOSED or
// EXPIRED offer cannot be revived by editing it, it must be recreated.
func UpdateOffer(gc *sharedconfig.GlobalConfig, offerID, merchantUserID string, in UpdateOfferInput) (p2pModels.Offer, error) {
	offer, err := GetOfferByID(gc.DB, offerID)
	if err != nil {
		return offer, &tErrors.CustomError{Param: "offerId", Err: "error-offer-not-found", ErrMessage: "Offer not found"}
	}
	if offer.MerchantUserID != merchantUserID {
		return offer, &tErrors.CustomError{Param: "offerId", Err: "error-forbidden", ErrMessage: "You do not own this offer", Code: 403}
	}
	if offer.Status == p2pModels.OfferStatusClosed || offer.Status == p2pModels.OfferStatusExpired {
		return offer, &tErrors.CustomError{Param: "offerId", Err: "error-invalid-offer-state", ErrMessage: "This offer can no longer be edited"}
	}
	if _, e := decimal.NewFromString(in.Price); e != nil {
		return offer, &tErrors.CustomError{Param: "price", Err: "error-invalid-price", ErrMessage: "price must be a valid decimal number"}
	}
	minAmt, e := decimal.NewFromString(in.MinOrderAmount)
	if e != nil {
		return offer, &tErrors.CustomError{Param: "minOrderAmount", Err: "error-invalid-min-order-amount", ErrMessage: "minOrderAmount must be a valid decimal number"}
	}
	maxAmt, e := decimal.NewFromString(in.MaxOrderAmount)
	if e != nil {
		return offer, &tErrors.CustomError{Param: "maxOrderAmount", Err: "error-invalid-max-order-amount", ErrMessage: "maxOrderAmount must be a valid decimal number"}
	}
	if maxAmt.LessThan(minAmt) {
		return offer, &tErrors.CustomError{Param: "maxOrderAmount", Err: "error-max-below-min", ErrMessage: "maxOrderAmount cannot be less than minOrderAmount"}
	}

	offer.PaymentMethod = in.PaymentMethod
	offer.PriceType = orDefaultStr(in.PriceType, offer.PriceType)
	offer.Price = in.Price
	offer.PriceMargin = orDefaultStr(in.PriceMargin, "0")
	offer.MinOrderAmount = in.MinOrderAmount
	offer.MaxOrderAmount = in.MaxOrderAmount
	// AvailableLiquidity is only raised/lowered by the merchant's own
	// top-up amount, not overwritten wholesale - a live offer may already
	// have some of its liquidity reserved by in-flight orders, and a plain
	// overwrite here could either strand a reservation or double count it.
	if in.AvailableLiquidity != "" {
		delta, e := decimal.NewFromString(in.AvailableLiquidity)
		if e != nil {
			return offer, &tErrors.CustomError{Param: "availableLiquidity", Err: "error-invalid-liquidity", ErrMessage: "availableLiquidity must be a valid decimal number"}
		}
		offer.AvailableLiquidity = decimal.RequireFromString(orDefaultStr(offer.AvailableLiquidity, "0")).Add(delta).String()
	}
	offer.Remark = in.Remark
	offer.Version++

	if err := gc.DB.Save(&offer).Error; err != nil {
		return offer, &tErrors.ErrorTemporaryServerError{}
	}
	RecordAuditEvent(gc, "", offer.ID, p2pModels.EventOfferUpdated, merchantUserID, offer)
	return offer, nil
}

// CloseOffer permanently retires an offer (Plan Section 13's implied
// merchant control over listing lifecycle) - unlike Pause, a closed offer
// cannot be reactivated. Orders already in flight reference their own
// snapshot of the offer's terms taken at creation time, so closing does
// not affect them.
func CloseOffer(gc *sharedconfig.GlobalConfig, offerID, merchantUserID string) (p2pModels.Offer, error) {
	offer, err := GetOfferByID(gc.DB, offerID)
	if err != nil {
		return offer, &tErrors.CustomError{Param: "offerId", Err: "error-offer-not-found", ErrMessage: "Offer not found"}
	}
	if offer.MerchantUserID != merchantUserID {
		return offer, &tErrors.CustomError{Param: "offerId", Err: "error-forbidden", ErrMessage: "You do not own this offer", Code: 403}
	}
	if offer.Status == p2pModels.OfferStatusClosed {
		return offer, &tErrors.CustomError{Param: "offerId", Err: "error-invalid-offer-state", ErrMessage: "This offer is already closed"}
	}
	offer.Status = p2pModels.OfferStatusClosed
	offer.AvailabilityStatus = p2pModels.OfferAvailabilityOffline
	if err := gc.DB.Save(&offer).Error; err != nil {
		return offer, &tErrors.ErrorTemporaryServerError{}
	}
	RecordAuditEvent(gc, "", offer.ID, p2pModels.EventOfferUpdated, merchantUserID, offer)
	return offer, nil
}

// GetOfferByID fetches a single offer.
func GetOfferByID(db *gorm.DB, offerID string) (p2pModels.Offer, error) {
	var offer p2pModels.Offer
	err := db.Where("id = ?", offerID).First(&offer).Error
	return offer, err
}

// MarketplaceFilter narrows the browsable offer list (Plan Section 13).
type MarketplaceFilter struct {
	OfferType   string
	Asset       string
	CountryCode string
	Currency    string
	Page        int
	PageSize    int
}

// ListMarketplaceOffers returns ACTIVE+ONLINE offers matching the filter,
// ordered by the most recent ranking snapshot (Plan Section 14) where one
// exists, falling back to newest-first.
func ListMarketplaceOffers(db *gorm.DB, f MarketplaceFilter) ([]p2pModels.Offer, int64, error) {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize < 1 || f.PageSize > 100 {
		f.PageSize = 20
	}
	q := db.Model(&p2pModels.Offer{}).Where("status = ?", p2pModels.OfferStatusActive)
	if f.OfferType != "" {
		q = q.Where("offer_type = ?", f.OfferType)
	}
	if f.Asset != "" {
		q = q.Where("asset = ?", f.Asset)
	}
	if f.CountryCode != "" {
		q = q.Where("country_code = ?", f.CountryCode)
	}
	if f.Currency != "" {
		q = q.Where("currency = ?", f.Currency)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var offers []p2pModels.Offer
	// Online offers first; within each group, a ranked offer (CalculateRanking's
	// periodic sweep, Plan Section 14) sorts by its rank ascending (1 = best),
	// and an offer with no snapshot yet falls back to newest-first.
	err := q.Joins("LEFT JOIN offer_ranking_snapshots ON offer_ranking_snapshots.offer_id = offers.id").
		Order("offers.availability_status desc, offer_ranking_snapshots.rank IS NULL, offer_ranking_snapshots.rank ASC, offers.updated_at desc").
		Offset((f.Page - 1) * f.PageSize).Limit(f.PageSize).Find(&offers).Error
	return offers, total, err
}

// ListMerchantOffers returns every offer owned by a merchant, any status.
func ListMerchantOffers(db *gorm.DB, merchantUserID string) ([]p2pModels.Offer, error) {
	var offers []p2pModels.Offer
	err := db.Where("merchant_user_id = ?", merchantUserID).Order("created_at desc").Find(&offers).Error
	return offers, err
}

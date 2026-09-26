package p2p

import (
	p2pModels "trovo-wallet-api/internal/components/p2p/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/sharedconfig"

	"gorm.io/gorm"
)

// CreatePaymentMethodInput is the client-supplied shape for saving a new
// merchant payment method.
type CreatePaymentMethodInput struct {
	PaymentChannel string
	Provider       string
	Account        string
}

// CreatePaymentMethod saves a new fiat settlement channel for merchantUserID,
// active by default so it can be selected from a SELL offer immediately -
// including from the offer-creation form's own inline "Add new payment
// method" quick-create flow.
func CreatePaymentMethod(gc *sharedconfig.GlobalConfig, merchantUserID, merchantUsername string, in CreatePaymentMethodInput) (p2pModels.MerchantPaymentMethod, error) {
	if in.PaymentChannel == "" || in.Account == "" {
		return p2pModels.MerchantPaymentMethod{}, &tErrors.CustomError{Param: "paymentChannel", Err: "error-invalid-payment-method", ErrMessage: "paymentChannel and account are required"}
	}
	pm := p2pModels.MerchantPaymentMethod{
		ID:               gc.GenerateUUIDString(),
		MerchantUserID:   merchantUserID,
		MerchantUsername: merchantUsername,
		PaymentChannel:   in.PaymentChannel,
		Provider:         in.Provider,
		Account:          in.Account,
		IsActive:         true,
	}
	if err := gc.DB.Create(&pm).Error; err != nil {
		return p2pModels.MerchantPaymentMethod{}, &tErrors.ErrorTemporaryServerError{}
	}
	return pm, nil
}

// ListMyPaymentMethods returns every payment method merchantUserID has ever
// saved, active or not - the merchant's own management screen needs to
// show disabled ones too (with a re-enable action), not just hide them.
func ListMyPaymentMethods(db *gorm.DB, merchantUserID string) ([]p2pModels.MerchantPaymentMethod, error) {
	var methods []p2pModels.MerchantPaymentMethod
	err := db.Where("merchant_user_id = ?", merchantUserID).Order("created_at desc").Find(&methods).Error
	return methods, err
}

// GetPaymentMethodByID fetches a single payment method.
func GetPaymentMethodByID(db *gorm.DB, id string) (p2pModels.MerchantPaymentMethod, error) {
	var pm p2pModels.MerchantPaymentMethod
	err := db.Where("id = ?", id).First(&pm).Error
	return pm, err
}

// resolveSellPaymentMethod validates that paymentMethodID is one of
// merchantUserID's own active payment methods - the only way a SELL offer
// may reference one (Plan: offers select a payment method by reference,
// not a bespoke per-offer value, so an edit to the payment method can
// propagate to the offer - see UpdatePaymentMethod).
func resolveSellPaymentMethod(db *gorm.DB, merchantUserID, paymentMethodID string) (p2pModels.MerchantPaymentMethod, error) {
	if paymentMethodID == "" {
		return p2pModels.MerchantPaymentMethod{}, &tErrors.CustomError{Param: "paymentMethodId", Err: "error-payment-method-required", ErrMessage: "A payment method is required for a SELL offer"}
	}
	pm, err := GetPaymentMethodByID(db, paymentMethodID)
	if err != nil {
		return pm, &tErrors.CustomError{Param: "paymentMethodId", Err: "error-payment-method-not-found", ErrMessage: "Payment method not found"}
	}
	if pm.MerchantUserID != merchantUserID {
		return pm, &tErrors.CustomError{Param: "paymentMethodId", Err: "error-forbidden", ErrMessage: "You do not own this payment method", Code: 403}
	}
	if !pm.IsActive {
		return pm, &tErrors.CustomError{Param: "paymentMethodId", Err: "error-payment-method-inactive", ErrMessage: "This payment method is disabled"}
	}
	return pm, nil
}

// UpdatePaymentMethodInput is the client-supplied shape for editing an
// existing payment method's details.
type UpdatePaymentMethodInput struct {
	PaymentChannel string
	Provider       string
	Account        string
}

// UpdatePaymentMethod edits a payment method's details and propagates the
// change to every live offer currently referencing it, so a SELL offer
// always shows its merchant's current settlement details. This never
// touches an already-created Order.PaymentMethodSnapshot (Plan: "editing
// payment method does not change the snapshot in existing orders") - that
// copy was frozen at order-creation time and is not linked back here.
func UpdatePaymentMethod(gc *sharedconfig.GlobalConfig, id, merchantUserID string, in UpdatePaymentMethodInput) (p2pModels.MerchantPaymentMethod, error) {
	pm, err := GetPaymentMethodByID(gc.DB, id)
	if err != nil {
		return pm, &tErrors.CustomError{Param: "paymentMethodId", Err: "error-payment-method-not-found", ErrMessage: "Payment method not found"}
	}
	if pm.MerchantUserID != merchantUserID {
		return pm, &tErrors.CustomError{Param: "paymentMethodId", Err: "error-forbidden", ErrMessage: "You do not own this payment method", Code: 403}
	}
	if in.PaymentChannel == "" || in.Account == "" {
		return pm, &tErrors.CustomError{Param: "paymentChannel", Err: "error-invalid-payment-method", ErrMessage: "paymentChannel and account are required"}
	}
	pm.PaymentChannel = in.PaymentChannel
	pm.Provider = in.Provider
	pm.Account = in.Account
	if err := gc.DB.Save(&pm).Error; err != nil {
		return pm, &tErrors.ErrorTemporaryServerError{}
	}
	gc.DB.Model(&p2pModels.Offer{}).Where("payment_method_id = ?", pm.ID).Updates(map[string]interface{}{
		"payment_method_payment_channel": pm.PaymentChannel,
		"payment_method_provider":        pm.Provider,
		"payment_method_account":         pm.Account,
	})
	return pm, nil
}

// SetPaymentMethodActive toggles a payment method's active flag. It is
// never hard-deleted (Plan: "you cannot delete payment method"), and
// cannot be deactivated while any non-terminal offer still references it
// (Plan: "nobody can remove payment method if it's actively in use by an
// offer... either they remove the offer or edit it to change the payment
// method before they can disable/enable" it) - CLOSED/EXPIRED offers don't
// count, since they can never trade again regardless.
func SetPaymentMethodActive(gc *sharedconfig.GlobalConfig, id, merchantUserID string, active bool) (p2pModels.MerchantPaymentMethod, error) {
	pm, err := GetPaymentMethodByID(gc.DB, id)
	if err != nil {
		return pm, &tErrors.CustomError{Param: "paymentMethodId", Err: "error-payment-method-not-found", ErrMessage: "Payment method not found"}
	}
	if pm.MerchantUserID != merchantUserID {
		return pm, &tErrors.CustomError{Param: "paymentMethodId", Err: "error-forbidden", ErrMessage: "You do not own this payment method", Code: 403}
	}
	if !active {
		var inUseCount int64
		gc.DB.Model(&p2pModels.Offer{}).
			Where("payment_method_id = ? AND status NOT IN ?", pm.ID, []string{p2pModels.OfferStatusClosed, p2pModels.OfferStatusExpired}).
			Count(&inUseCount)
		if inUseCount > 0 {
			return pm, &tErrors.CustomError{Param: "paymentMethodId", Err: "error-payment-method-in-use", ErrMessage: "This payment method is in use by an active offer. Close the offer or change its payment method first.", Code: 409}
		}
	}
	// GORM's struct-Save silently drops a false zero-value for a field
	// with a `default:` tag - Model().Update forces it to persist (the
	// same fix used for User.MerchantOnline earlier in this codebase).
	if err := gc.DB.Model(&p2pModels.MerchantPaymentMethod{}).Where("id = ?", pm.ID).Update("is_active", active).Error; err != nil {
		return pm, &tErrors.ErrorTemporaryServerError{}
	}
	pm.IsActive = active
	return pm, nil
}

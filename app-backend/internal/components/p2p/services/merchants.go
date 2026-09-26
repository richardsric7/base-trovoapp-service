package p2p

import (
	p2pModels "trovo-wallet-api/internal/components/p2p/models"
	usersDB "trovo-wallet-api/internal/components/users/db"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/sharedconfig"
)

// merchantKYCLevelRequired is the minimum completed KYC level a user needs
// before they can request P2P merchant status - User.KYCVerified is the
// single canonical field for a user's completed KYC level across both
// providers (SumSub and Dojah both write it), not the per-provider
// UserKYCProgress/UserDojaKYCProgress rows, which only track an
// in-flight verification's own initiated/done steps.
const merchantKYCLevelRequired = 2

// RequestMerchantStatus enables userID as a P2P merchant once they've
// completed KYC level 2, per your instruction. Already being a merchant is
// a no-op success (not an error) - simpler for the client than having to
// special-case a double-request. A user who isn't at KYC level 2 yet gets
// a clear, actionable error instead of silently failing.
func RequestMerchantStatus(gc *sharedconfig.GlobalConfig, userID string) (userModels.User, error) {
	user, err := usersDB.GetUser(userID, gc.DB, gc)
	if err != nil {
		return user, err
	}
	if user.IsMerchant {
		return user, nil
	}
	if user.KYCVerified < merchantKYCLevelRequired {
		return user, &tErrors.CustomError{
			Param:      "kycLevel",
			Err:        "error-kyc-level-2-required",
			ErrMessage: "You need to complete KYC level 2 before you can become a P2P merchant.",
			Code:       403,
		}
	}
	user.IsMerchant = true
	user.MerchantOnline = true
	if err := gc.DB.Model(&userModels.User{}).Where("id = ?", user.ID).
		Updates(map[string]interface{}{"is_merchant": true, "merchant_online": true}).Error; err != nil {
		return user, &tErrors.ErrorTemporaryServerError{}
	}
	RecordAuditEvent(gc, "", "", p2pModels.EventMerchantEnabled, user.ID, nil)
	return user, nil
}

// requireMerchantOnline is CreateOrder/QuoteOrderFees's half of the
// merchant-online gate: ListMarketplaceOffers keeps an offline merchant's
// offers out of search results, but a customer who already has an
// offerID (a bookmark, a stale client cache, a shared link) could still
// try to order against it directly - this closes that gap by rejecting
// the order/quote outright rather than relying on search filtering alone.
func requireMerchantOnline(gc *sharedconfig.GlobalConfig, merchantUserID string) error {
	merchant, err := usersDB.GetUser(merchantUserID, gc.DB, gc)
	if err != nil {
		return err
	}
	if !merchant.IsMerchant || !merchant.MerchantOnline {
		return &tErrors.CustomError{
			Param:      "offerId",
			Err:        "error-merchant-offline",
			ErrMessage: "This merchant is currently offline and not accepting new orders.",
			Code:       403,
		}
	}
	return nil
}

// SetMerchantOnlineStatus toggles userID's merchant online/offline status.
// Going offline hides every one of the merchant's offers from marketplace
// search (ListMarketplaceOffers) and blocks new orders against them
// (CreateOrder/QuoteOrderFees) until they toggle back online - see those
// functions' own doc comments for why this is enforced by a join against
// User.MerchantOnline rather than by bulk-flipping every affected offer's
// own AvailabilityStatus.
func SetMerchantOnlineStatus(gc *sharedconfig.GlobalConfig, userID string, online bool) (userModels.User, error) {
	user, err := usersDB.GetUser(userID, gc.DB, gc)
	if err != nil {
		return user, err
	}
	if !user.IsMerchant {
		return user, &tErrors.CustomError{
			Param:      "userId",
			Err:        "error-not-a-merchant",
			ErrMessage: "Only merchants can go online/offline. Request merchant status first.",
			Code:       403,
		}
	}
	wasOnline := user.MerchantOnline
	user.MerchantOnline = online
	if err := gc.DB.Model(&userModels.User{}).Where("id = ?", user.ID).
		Update("merchant_online", online).Error; err != nil {
		return user, &tErrors.ErrorTemporaryServerError{}
	}
	eventName := p2pModels.EventMerchantOffline
	if online {
		eventName = p2pModels.EventMerchantOnline
	}
	RecordAuditEvent(gc, "", "", eventName, user.ID, nil)

	// Push-notify only on an actual online->offline transition, not a
	// redundant re-toggle-off call, so the merchant doesn't get repeat
	// pings if the client retries or calls this twice - they might
	// otherwise forget their offers have gone dark and lose sales without
	// realizing why.
	if wasOnline && !online {
		NotifyUsername(gc, user.Username, "You're offline on P2P",
			"Your P2P offers won't appear in marketplace search until you go back online.",
			map[string]string{"type": "P2P_MERCHANT_OFFLINE"})
	}
	return user, nil
}

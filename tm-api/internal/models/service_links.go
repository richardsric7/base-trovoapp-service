package models

import "time"

// ServiceLink is a writable mirror of app-backend's white-label partner
// integration record (internal/components/servicelinks/models.ServiceLink
// in the app-backend repo, table `service_links`). Like CuratedAsset,
// tm-api writes to this table directly: app-backend has never exposed any
// endpoint to create or edit a ServiceLink row - every existing row was
// provisioned by a direct DB insert - so managing which partner
// integrations exist, and what each one is permitted to do, is a pure
// admin/config concern with no app-backend business logic of its own.
//
// ID and ApiKey are plain string columns in app-backend (no DB-level UUID
// type), but every consumer treats them as GUIDs; tm-api generates both
// with uuid.NewString() on create and never accepts them from the client.
type ServiceLink struct {
	ID string `gorm:"column:id;primaryKey" json:"id"`
	// ApiKey is the credential the partner integration authenticates with
	// against app-backend's /v1/servicelinks and /v1/trovo-api routes.
	ApiKey    string    `gorm:"column:api_key" json:"apiKey"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	// OwnerUsername must be an existing row in the shared users table -
	// there's no DB-level foreign key (app-backend never declared one), so
	// tm-api validates it against `users.username` itself on create.
	OwnerUsername string `gorm:"column:owner_username" json:"ownerUsername"`
	// Address mirrors the owner's own wallet address (users.address) at
	// creation time - copied automatically rather than typed by the admin,
	// since it must always agree with the owner account's own public key.
	Address   string `gorm:"column:address" json:"address"`
	ShortName string `gorm:"column:short_name" json:"shortName"`
	LongName  string `gorm:"column:long_name" json:"longName"`

	// Permission flags - each gates a specific capability in app-backend's
	// servicelinks module (see handlers_impl.go's per-endpoint checks).
	LoginPermission                       int `gorm:"column:login_permission" json:"loginPermission"`
	PaymentPermission                     int `gorm:"column:payment_permission" json:"paymentPermission"`
	TokenInfoPermission                   int `gorm:"column:token_info_permission" json:"tokenInfoPermission"`
	AuthorizationPermission               int `gorm:"column:authorization_permission" json:"authorizationPermission"`
	EventPermission                       int `gorm:"column:event_permission" json:"eventPermission"`
	AllowUserInfo                         int `gorm:"column:allow_user_info" json:"allowUserInfo"`
	PushNotificationPermission            int `gorm:"column:push_notification_permission" json:"pushNotificationPermission"`
	IncludePhoneNumbers                   int `gorm:"column:include_phone_numbers" json:"includePhoneNumbers"`
	IncludeUserBalances                   int `gorm:"column:include_user_balances" json:"includeUserBalances"`
	TokenizedAssetAuthorizationPermission int `gorm:"column:tokenized_asset_authorization_permission" json:"tokenizedAssetAuthorizationPermission"`
	CreateUsersPermission                 int `gorm:"column:create_users_permission" json:"createUsersPermission"`
	AllowReferralForRegisteredUsers       int `gorm:"column:allow_referral_for_registered_users" json:"allowReferralForRegisteredUsers"`
	// Verified gates app-backend's requireActiveServiceLink alongside
	// Inactive/Suspended - it acts as an overall "fully vetted" toggle
	// rather than a single-endpoint permission, but is managed the same
	// way as the flags above.
	Verified int `gorm:"column:verified" json:"verified"`

	// Inactive is the routine admin on/off switch (0 = active, 1 =
	// inactive - the same convention as CuratedAsset.Inactive). Suspended
	// is a separate, harder block used by fraud/violation processes
	// elsewhere in app-backend; tm-api surfaces it read-only here rather
	// than managing it from this page.
	Inactive         int     `gorm:"column:inactive" json:"inactive"`
	Suspended        int     `gorm:"column:suspended" json:"suspended"`
	SuspensionReason *string `gorm:"column:suspension_reason" json:"suspensionReason"`
}

func (ServiceLink) TableName() string { return "service_links" }

// ServiceLinkOwnerSummary is a best-effort enrichment of a ServiceLink
// list/detail response with a few fields from the owner's own user
// record (email, KYC/merchant/account status) - looked up by username,
// since there's no FK to join on at the DB level.
type ServiceLinkOwnerSummary struct {
	Email          string `json:"email"`
	Address        string `json:"address"`
	KYCVerified    int    `json:"kycVerified"`
	IsMerchant     bool   `json:"isMerchant"`
	MerchantOnline bool   `json:"merchantOnline"`
	Corporate      int    `json:"corporate"`
	MembershipType int    `json:"membershipType"`
	Verified       int    `json:"verified"`
	Suspended      int    `json:"suspended"`
}

// ServiceLinkWithOwner is what the list/get endpoints actually return -
// the ServiceLink row plus its owner summary (nil if the owner username
// no longer resolves to a real user).
type ServiceLinkWithOwner struct {
	ServiceLink
	Owner *ServiceLinkOwnerSummary `json:"owner"`
}

// ServiceLinkRequest is the create/update payload. Only OwnerUsername
// (create only - immutable after), ShortName, LongName and the
// permission/Inactive flags are admin-editable; ID and ApiKey are always
// server-generated, and Address is always copied from the owner's own
// user record rather than accepted from the client.
type ServiceLinkRequest struct {
	Action                                string `json:"action" binding:"required,oneof=create update" enums:"create,update"`
	ID                                    string `json:"id,omitempty"`
	OwnerUsername                         string `json:"ownerUsername" binding:"required"`
	ShortName                             string `json:"shortName" binding:"required"`
	LongName                              string `json:"longName"`
	LoginPermission                       bool   `json:"loginPermission"`
	PaymentPermission                     bool   `json:"paymentPermission"`
	TokenInfoPermission                   bool   `json:"tokenInfoPermission"`
	AuthorizationPermission               bool   `json:"authorizationPermission"`
	EventPermission                       bool   `json:"eventPermission"`
	AllowUserInfo                         bool   `json:"allowUserInfo"`
	PushNotificationPermission            bool   `json:"pushNotificationPermission"`
	IncludePhoneNumbers                   bool   `json:"includePhoneNumbers"`
	IncludeUserBalances                   bool   `json:"includeUserBalances"`
	TokenizedAssetAuthorizationPermission bool   `json:"tokenizedAssetAuthorizationPermission"`
	CreateUsersPermission                 bool   `json:"createUsersPermission"`
	AllowReferralForRegisteredUsers       bool   `json:"allowReferralForRegisteredUsers"`
	Verified                              bool   `json:"verified"`
	Inactive                              bool   `json:"inactive"`
}

// ServiceLinkListRequest filters the paginated service-links admin table.
type ServiceLinkListRequest struct {
	Page          int
	PageSize      int
	OwnerUsername string
	ShortName     string
	Inactive      *bool
	Verified      *bool
	Suspended     *bool
}

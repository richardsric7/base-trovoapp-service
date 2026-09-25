package models

import (
	"time"
)

// Organization and Invite Status Constants
const (
	// Organization Statuses
	OrganizationStatusPending   = "PENDING"
	OrganizationStatusActive    = "ACTIVE"
	OrganizationStatusInactive  = "INACTIVE"
	OrganizationStatusSuspended = "SUSPENDED"

	// Invite Statuses
	InviteStatusPending       = "PENDING"
	EndInviteStatusValidEmail = "VALID_EMAIL" // Status after email verification

	InviteStatusAccepted   = "ACCEPTED"
	InviteStatusInProgress = "IN_PROGRESS"
	InviteStatusRejected   = "REJECTED"
	InviteStatusExpired    = "EXPIRED"
	InviteStatusVerified   = "VERIFIED"

	// Organization Types
	OrganizationTypeAssetManager       = "ASSET_MANAGER"
	OrganizationTypeAssetCustodian     = "ASSET_CUSTODIAN"
	OrganizationTypeRatingAgency       = "RATING_AGENCY"
	OrganizationTypeRegulator          = "REGULATOR"
	OrganizationTypeLegalAgency        = "LEGAL_AGENCY"
	OrganizationTypeProfessionalAgency = "PROFESSIONAL_AGENCY"

	// Stakeholder Types (must match partner types from existing /partners endpoints)
	StakeholderTypeCustodian    = "approved_asset_custodian"
	StakeholderTypeAssetManager = "asset_manager"
	StakeholderTypeIssuingHouse = "asset_issuing_house"
	StakeholderTypeLegal        = "legal_and_professionals"
	StakeholderTypeRatingAgency = "rating_agency"
	StakeholderTypeTrustee      = "trustees"
	// Legal Adviser and Financial Adviser are distinct A5 (due diligence & transaction
	// structuring) roles, separate from legal_and_professionals. They can be created,
	// linked to an org, and assigned to an asset, but are intentionally NOT granted
	// Stakeholder Portal access (they 403 from the portal — see PRD §2.3 / OI-14).
	StakeholderTypeLegalAdviser     = "legal_adviser"
	StakeholderTypeFinancialAdviser = "financial_adviser"
)

// Organization represents an organization in the system
type Organization struct {
	ID                string    `gorm:"primaryKey" json:"id"`
	Name              string    `gorm:"not null" json:"name"`
	Email             string    `gorm:"not null;unique" json:"email"`
	Type              string    `gorm:"not null" json:"type"`                     // CORPORATE, INDIVIDUAL
	Status            string    `gorm:"not null;default:'PENDING'" json:"status"` // PENDING, VERIFIED, ACTIVE, SUSPENDED
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	CreatedBy         string    `gorm:"not null" json:"created_by"`             // ID of the user who created this organization
	OrganizationCount int64     `gorm:"default:0" json:"team_member_count"`     // Count of members in the organization
	StakeholderID     *uint64   `gorm:"null" json:"stakeholder_id,omitempty"`   // Links to stakeholder in Trovo DB
	StakeholderType   *string   `gorm:"null" json:"stakeholder_type,omitempty"` // Type of linked stakeholder
	Level             *int      `gorm:"null" json:"level,omitempty"`            // Compliance requirement tier (e.g. Basic/Enhanced/Institutional); nil == unset
}

type OrganizationResponse struct {
	ID                string         `json:"id"`
	Name              string         `json:"name"`
	Email             string         `json:"email"`
	Type              string         `json:"type"`
	Status            string         `json:"status"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	TrovoAdminInfo    TrovoAdminInfo `json:"created_by_admin,omitempty"` // Admin who created this organization
	OrganizationCount int64          `json:"team_member_count"`
	StakeholderID     *uint64        `json:"stakeholder_id,omitempty"`   // Stakeholder linkage
	StakeholderType   *string        `json:"stakeholder_type,omitempty"` // Stakeholder type
	Level             *int           `json:"level,omitempty"`
	Address           string         `json:"address,omitempty"`
	Country           string         `json:"country,omitempty"`
	FeeFixed          string         `json:"fee_fixed,omitempty"`
	FeePercent        string         `json:"fee_percent,omitempty"`
}
type TrovoAdminInfo struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Username  string `json:"username"`
}

// OrganizationMember represents a member of an organization
type OrganizationMember struct {
	ID                     string       `gorm:"primaryKey" json:"id"`
	OrganizationID         string       `gorm:"not null;index:idx_org_member" json:"organization_id"`
	Email                  string       `gorm:"not null;unique" json:"email"`
	Password               string       `gorm:"null" json:"-"`                            // Hashed password
	Role                   string       `gorm:"not null;default:'MEMBER'" json:"role"`    // ADMIN, MEMBER
	Status                 string       `gorm:"not null;default:'PENDING'" json:"status"` // PENDING, ACTIVE, SUSPENDED
	CreatedAt              time.Time    `json:"created_at"`
	UpdatedAt              time.Time    `json:"updated_at"`
	CreatedBy              string       `gorm:"not null" json:"created_by"` // ID of the user who invited this member
	Organization           Organization `gorm:"foreignKey:OrganizationID" json:"organization,omitempty"`
	VerificationOTP        string       `gorm:"null" json:"-"`          // OTP for email verification
	OTPExpiresAt           time.Time    `gorm:"null" json:"-"`          // OTP expiration time
	PasswordSetupToken     string       `gorm:"null" json:"-"`          // Token for password setup
	PasswordSetupExpiresAt time.Time    `gorm:"null" json:"-"`          // Password setup token expiration time
	PasswordResetOTP       string       `gorm:"null" json:"-"`          // OTP for password reset
	PasswordResetExpiresAt time.Time    `gorm:"null" json:"-"`          // Expiration time for password reset OTP
	FirstName              string       `gorm:"null" json:"first_name"` // First name of the member
	LastName               string       `gorm:"null" json:"last_name"`  // Last name of the member
	InviteID               string       `gorm:"null" json:"invite_id"`  // ID of the invite that created this member
	KYCStatus              int          `gorm:"null" json:"kyc_status"` // KYC status of the member
	TrovoWalletUsername    *string      `gorm:"null" json:"trovo_wallet_username,omitempty"`
	IsWalletLinked         bool         `gorm:"default:false" json:"is_wallet_linked"`
	SessionVersion         uint64       `gorm:"not null;default:0" json:"-"`
}

// OrganizationInvite represents an invitation to join an organization
type OrganizationInvite struct {
	ID                    string    `gorm:"primaryKey" json:"id"`
	OrganizationID        string    `gorm:"not null" json:"organization_id"`
	Email                 string    `gorm:"not null;uniqueIndex:idx_org_invite_email" json:"email"` // Email of the invited member
	EmailValidationStatus string    `gorm:"null" json:"email_validation_status"`
	EmailValidatedAt      time.Time `gorm:"null" json:"email_validated_at"`
	Role                  string    `gorm:"not null;default:'MEMBER'" json:"role"`    // ADMIN, MEMBER
	Status                string    `gorm:"not null;default:'PENDING'" json:"status"` // PENDING, ACCEPTED, REJECTED, EXPIRED
	ExpiresAt             time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
	InvitedBy             string    `gorm:"not null" json:"invited_by,omitempty"`
	VerificationOTP       string    `gorm:"null" json:"-"` // OTP for email verification
	OTPExpiresAt          time.Time `gorm:"null" json:"-"` // OTP expiration time
}

// OrganizationCreateRequest represents the request to create a new organization
type OrganizationCreateRequest struct {
	// Organization details
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Type  string `json:"type" binding:"required"`

	// Root user (organization admin) details
	RootUserFirstName string `json:"root_user_first_name" binding:"required"`
	RootUserLastName  string `json:"root_user_last_name" binding:"required"`
	RootUserEmail     string `json:"root_user_email" binding:"required,email"`
	RootUserPassword  string `json:"root_user_password" binding:"required,min=8"`
}

// RootUserValidatesEmail represents the request to validate a root user's email
type RootUserValidatesEmail struct {
	OTP string `json:"otp" binding:"required"`
}

// OrganizationVerifyRequest represents the request to verify an organization's email
type OrganizationVerifyRequest struct {
	OrganizationID string `json:"organization_id" binding:"required"`
	OTP            string `json:"otp" binding:"required"`
}

// OrganizationSetPasswordRequest represents the request to set an organization's password
type OrganizationSetPasswordRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

// LinkWalletRequest is the request body for linking a Trovo Wallet to the authenticated org member (onboarded user). Only username is required.
type LinkWalletRequest struct {
	Username string `json:"username" binding:"required"`
}

// WalletLinkAuthorizationRequest starts wallet linking for an organization member.
type WalletLinkAuthorizationRequest struct {
	WalletUsername string `json:"wallet_username" binding:"required"`
}

// OrganizationInviteRequest represents the request to invite a member to an organization
type OrganizationInviteRequest struct {
	Email          string `json:"email" binding:"required,email"`
	FirstName      string `json:"first_name" binding:"required"`
	LastName       string `json:"last_name" binding:"required"`
	OrganizationID string `json:"organization_id"`
}

// OrganizationMemberResponse represents the response for organization members
type OrganizationMemberResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// ResendPasswordSetupRequest represents the request to resend the password setup token
type ResendPasswordSetupRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// OrganizationMemberLoginRequest represents the request for member login
type OrganizationMemberLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// OrganizationMemberPasswordResetRequest represents the request to reset a member's password
type OrganizationMemberPasswordResetRequest struct {
	OrganizationID string `json:"organization_id" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
}

// OrganizationMemberSetPasswordRequest represents the request to set a new password after reset
type OrganizationMemberSetPasswordRequest struct {
	OrganizationID string `json:"organization_id" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	OTP            string `json:"otp" binding:"required"`
	NewPassword    string `json:"new_password" binding:"required,min=8"`
}

// RootUserSetupProfileRequest represents the request to set up a root user's profile
type RootUserSetupProfileRequest struct {
	Email           string `json:"email" binding:"required,email"`
	Token           string `json:"token" binding:"required"`
	FirstName       string `json:"first_name" binding:"required"`
	LastName        string `json:"last_name" binding:"required"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
	// Organization details
	OrganizationName  string `json:"organization_name" binding:"required"`
	OrganizationType  string `json:"organization_type" binding:"required,oneof=CORPORATE INDIVIDUAL"`
	OrganizationEmail string `json:"organization_email" binding:"required,email"`
}

// ResendOTPRequest represents the request to resend OTP for root user
type ResendOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	Token string `json:"token" binding:"required"`
}

// TeamMemberVerifyOTPRequest represents the request to verify a team member's OTP
type TeamMemberVerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	Token string `json:"token" binding:"required"`
	OTP   string `json:"otp" binding:"required"`
}

// TeamMemberSetupProfileRequest represents the request to set up a team member's profile
type TeamMemberSetupProfileRequest struct {
	Email           string `json:"email" binding:"required,email"`
	Token           string `json:"token" binding:"required"`
	FirstName       string `json:"first_name" binding:"required"`
	LastName        string `json:"last_name" binding:"required"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

// TeamMemberResendOTPRequest represents the request to resend OTP for team member
type TeamMemberResendOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	Token string `json:"token" binding:"required"`
}

// TokenResponse represents a response containing a JWT token
type TokenResponse struct {
	Token string `json:"token"`
}

// SuccessResponse represents a generic success response
type SuccessResponse struct {
	Message string `json:"message"`
}

// SendEmailVerificationOTPRequest represents the request to send email verification OTP
type SendEmailVerificationOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	Token string `json:"token" binding:"required"`
}

// VerifyEmailOTPRequest represents the request to verify email OTP
type VerifyEmailOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required"`
}

// ValidateMemberInviteRequest represents the request to validate a member's invite
type ValidateMemberInviteRequest struct {
	Otp string `json:"otp" binding:"required"`
}

// SetupRootUserPasswordRequest represents the request to set up a root user's password
type SetupUserPasswordRequest struct {
	Password string `json:"password" binding:"required,min=8"`
	OTP      string `json:"otp" binding:"required"`
}

type InviteOrgMemberPayload struct {
	OrganizationName string
	OrganizationType string
	AdminFirstName   string
	AdminLastName    string
	AdminEmail       string
	AdminID          string
}

type InviteUserFlowStruct struct {
	OrganizationName string
	OrganizationType string
	InviteID         string
	OrganizationID   string
}

package models

import "time"

// AdminAccessLog is the security / access audit trail for Trovo staff. It records
// who signed in, and every privileged action an admin takes across the panel that
// has no other trail today. It deliberately does NOT duplicate the stakeholder
// governance log (fund_release / distribution / due_diligence, in stakeholder_audit_logs)
// or the vault-signer on-chain log (vault_signer_audit_logs) -- this log owns only
// the untracked surface.
//
// Identity is snapshotted onto the row (Username / FullName / Email) rather than
// stored purely as a foreign key, so the record stays truthful if the admin is later
// renamed or deleted, and so the list endpoint needs no per-row join.
type AdminAccessLog struct {
	ID     string `gorm:"primaryKey" json:"id"`
	Event  string `gorm:"index" json:"event"`    // login.success, admin.suspend, secret.write, …
	Status string `gorm:"index" json:"status"`   // successful | failed | pending
	Action string `json:"action"`                // humanised label rendered in the UI "Action" column

	// Category drives retention. Mutations are kept long-term; a future "read"
	// category (sensitive reads, not wired yet) will be purged on a schedule.
	Category string `gorm:"index" json:"category"` // session | account | secret | org | asset | config

	// Actor -- the admin who performed the action. ActorAdminID is nullable because
	// a failed login has no resolved admin id (only a submitted username).
	ActorAdminID *uint  `gorm:"index" json:"actor_admin_id,omitempty"`
	Username     string `gorm:"index" json:"username"`
	FullName     string `json:"fullname"`
	Email        string `json:"email,omitempty"`
	Role         string `json:"role,omitempty"`

	// Target of the action, when there is one distinct from the actor (e.g. the
	// admin who was suspended, the org that was deactivated, the secret slot written).
	Target string `json:"target,omitempty"`

	// Request context.
	IPAddress   string `json:"ip_address,omitempty"`
	Location    string `json:"location,omitempty"`     // "City, Country" via the existing IPAPI lookup
	UserAgent   string `json:"user_agent,omitempty"`
	Method      string `json:"method,omitempty"`       // HTTP method
	Path        string `json:"path,omitempty"`         // request path
	LoginMethod string `json:"login_method,omitempty"` // push_approval | qr -- admin login is delegated, never a password
	Detail      string `json:"detail,omitempty"`       // failure reason, or extra context

	OccurredAt time.Time `gorm:"index" json:"occurred_at"`
}

// TableName pins the table name so it is stable regardless of GORM pluralisation rules.
func (AdminAccessLog) TableName() string { return "admin_access_logs" }

// Access-log event names. Grouped by the category they carry.
const (
	// session
	EventLoginSuccess   = "login.success"
	EventLoginFailure   = "login.failure"
	EventLoginSuspended = "login.suspended"
	EventLogout         = "logout"
	EventTokenRefresh   = "token.refresh"

	// account
	EventAdminInvite       = "admin.invite"
	EventAdminSuspend      = "admin.suspend"
	EventAdminUnsuspend    = "admin.unsuspend"
	EventAdminStatusChange = "admin.status_change"
	EventAdminRemove       = "admin.remove"
	EventUserSuspend       = "user.suspend"
	EventUserLiftSuspend   = "user.lift_suspension"
	EventUserKycChange     = "user.kyc_level_change"
	EventUserToggle        = "user.toggle"

	// secret
	EventSecretWrite            = "secret.write"
	EventManagedSecretRegister  = "managed_secret.register"
	EventSecretAssignmentGrant  = "secret_assignment.grant"
	EventSecretAssignmentEdit   = "secret_assignment.edit"
	EventSecretAssignmentRevoke = "secret_assignment.revoke"
	EventPersonalEnvCreate      = "personal_env.create"
	EventPersonalEnvDelete      = "personal_env.delete"

	// config
	EventConfigBulkUpdate = "config.bulk_update"
	EventConfigDelete     = "config.delete"
	EventFeeConfigChange  = "fee_config.change"
	EventFeeConfigDelete  = "fee_config.delete"
	EventKycConfigChange  = "kyc_config.change"
	EventKycLevelChange   = "kyc_level.change"

	// org
	EventOrgCreate       = "organization.create"
	EventOrgUpdate       = "organization.update"
	EventOrgDeactivate   = "organization.deactivate"
	EventOrgMemberInvite = "org_member.invite"
	EventWalletLink      = "wallet.link"

	// asset
	EventAssetAssignment     = "asset_assignment.create"
	EventAssetCurationChange = "asset_curation.change"
	EventTokenizationMint    = "tokenization.mint"
	EventTokenizationVet     = "tokenization.vet"
	EventTokenizationFail    = "tokenization.fail"
	EventMintingUserGrant    = "minting_user.grant"
	EventMintingUserRevoke   = "minting_user.revoke"
	EventComplianceCreate    = "compliance.create"
)

// Access-log status values (canonical lowercase; the UI filter's title-case is
// reconciled to these).
const (
	AccessStatusSuccessful = "successful"
	AccessStatusFailed     = "failed"
	AccessStatusPending    = "pending"
)

// Access-log categories.
const (
	AccessCategorySession = "session"
	AccessCategoryAccount = "account"
	AccessCategorySecret  = "secret"
	AccessCategoryOrg     = "org"
	AccessCategoryAsset   = "asset"
	AccessCategoryConfig  = "config"
)

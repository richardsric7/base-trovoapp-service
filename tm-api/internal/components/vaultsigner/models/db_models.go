package models

import "time"

// VaultSignerManagedSecret describes one shared, admin-registered signer-slot
// secret: a CSV of Base keys stored in Vault, controlling a Safe wallet
// on-chain. CSV capacity is not tracked here — it's whatever the live CSV in
// Vault actually holds, with a minimum of 4 entries enforced at write time.
type VaultSignerManagedSecret struct {
	ID                 string    `gorm:"primaryKey" json:"id"`
	Label              string    `gorm:"not null" json:"label"`
	VaultMount         string    `gorm:"not null" json:"vault_mount"`
	VaultPath          string    `gorm:"not null" json:"vault_path"`
	VaultField         string    `gorm:"not null" json:"vault_field"`
	WalletAddress      string    `gorm:"not null;index" json:"wallet_address"`
	ActiveSigningCount int       `gorm:"not null" json:"active_signing_count"`
	CreatedAt          time.Time `json:"created_at"`
}

func (VaultSignerManagedSecret) TableName() string {
	return "vault_signer_managed_secrets"
}

// VaultSignerAssignment maps an owner to a VaultSignerManagedSecret.
// Position is nullable and starts unset: an assignment claims its 1-based
// Signer Position only the first time its owner submits a real signer value
// (PUT /me/vault-signer/secrets/:secretId), not at assignment-creation time.
// Position IS NULL means "no one has ever set a real key for this
// assignment yet" — whether it's freshly created or was reassigned from a
// previous owner who also never set one. Reassigning an assignment to a
// different owner (PATCH) always retains whatever Position currently holds,
// null or not — the new owner simply claims/uses that same position the
// next time they submit their own value.
//
// Postgres treats multiple NULLs in a unique index as distinct from each
// other, so any number of not-yet-claimed assignments can coexist for the
// same managed secret without tripping the (ManagedSecretID, Position)
// uniqueness constraint below — that constraint only ever engages once a
// position is actually claimed.
//
// OwnerRefID stores the owner's real, stable primary key (AdminUser.ID for
// trovo_admin, OrganizationMember.ID for org_member) rather than a
// Username/Email — those fields are mutable, and a stored copy of a mutable
// field would silently go stale if it ever changed. OwnerType disambiguates
// OwnerRefID, since the two ID namespaces aren't unique across each other.
//
// (ManagedSecretID, OwnerRefID) is unique: the same owner can never hold
// more than one assignment row on the same managed secret. This is enforced
// both here (DB-level, via idx_vault_signer_assignment_secret_owner — a real
// constraint the database itself checks under concurrent writes, not just an
// app-level convention) and in services.CreateAssignment/EditAssignment
// (app-level, for a friendly 409 instead of a raw DB error surfacing to the
// caller). Without this, self-service lookups that key off
// (managed_secret_id, owner_ref_id, owner_type) alone (handlers.loadOwnedSecret)
// would have no way to pick between two rows for the same owner on the same
// secret — GORM's First() would pick one arbitrarily and the other would
// become permanently unreachable through PUT/GET /me/vault-signer/secrets/:secretId,
// even though it could still hold a claimed, live position.
//
// This also keeps Signer Position well-defined per (ManagedSecretID,
// OwnerRefID): since an owner can now only ever hold one assignment on a
// given managed secret, that pair maps to at most one Position. Position
// numbering itself is unchanged — ClaimPosition (services/signer_service.go)
// already scopes the lowest-unclaimed lookup to one managed secret at a
// time, so positions run 1, 2, 3, ... independently per managed secret, not
// globally — this constraint doesn't touch that, it just guarantees no owner
// can occupy more than one of those numbers on the same secret.
type VaultSignerAssignment struct {
	ID              string    `gorm:"primaryKey" json:"id"`
	ManagedSecretID string    `gorm:"not null;uniqueIndex:idx_vault_signer_assignment_secret_position;uniqueIndex:idx_vault_signer_assignment_secret_owner" json:"managed_secret_id"`
	Position        *int      `gorm:"uniqueIndex:idx_vault_signer_assignment_secret_position" json:"position"`
	OwnerRefID      string    `gorm:"not null;index:idx_vault_signer_assignment_owner;uniqueIndex:idx_vault_signer_assignment_secret_owner" json:"owner_ref_id"`
	OwnerType       string    `gorm:"not null;index:idx_vault_signer_assignment_owner" json:"owner_type"`
	AssignedAt      time.Time `json:"assigned_at"`
}

func (VaultSignerAssignment) TableName() string {
	return "vault_signer_assignments"
}

// VaultSignerAuditLog is a write-once historical record for every flow this
// feature has: signer-slot value writes, on-chain swaps, assignment
// lifecycle events, and personal-secret create/delete. It is written after
// the fact for traceability only — nothing polls it or reads it back to
// decide behavior.
type VaultSignerAuditLog struct {
	ID                 string    `gorm:"primaryKey" json:"id"`
	Kind               string    `gorm:"not null;index" json:"kind"`
	ManagedSecretID    *string   `gorm:"index" json:"managed_secret_id,omitempty"`
	Position           *int      `json:"position,omitempty"`
	VaultKey           *string   `json:"vault_key,omitempty"`
	ActorID            string    `gorm:"not null;index" json:"actor_id"`
	ActorType          string    `gorm:"not null" json:"actor_type"`
	VaultVersionBefore *int      `json:"vault_version_before,omitempty"`
	VaultVersionAfter  *int      `json:"vault_version_after,omitempty"`
	BaseTxHash         *string   `json:"base_tx_hash,omitempty"`
	BaseTxStatus       *string   `json:"base_tx_status,omitempty"`
	ChangedAt          time.Time `json:"changed_at"`
	IPAddress          *string   `json:"ip_address,omitempty"`
}

func (VaultSignerAuditLog) TableName() string {
	return "vault_signer_audit_logs"
}

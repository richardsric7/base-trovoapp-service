package services

import (
	"context"
	"fmt"
	"os"
	"strings"

	"admin-panel-dashboard/internal/components/vaultsigner/vaultclient"
	"admin-panel-dashboard/internal/middleware"
	coreModels "admin-panel-dashboard/internal/models"

	"github.com/gin-gonic/gin"
	vaultapi "github.com/hashicorp/vault/api"
	"gorm.io/gorm"
)

// PersonalEnvPrefix is one PERSONAL_ENVS entry, parsed once at startup.
type PersonalEnvPrefix struct {
	Prefix        string
	FriendlyLabel string
}

// ParsePersonalEnvs parses the PERSONAL_ENVS env var — a CSV of
// PREFIX:FriendlyLabel entries (Section 2) — into an ordered slice. There is
// no VaultSignerPersonalEnvDefinition table; this env var is the single
// source of truth for both the prefix and its display name. Fails fast if
// any entry doesn't contain the ":" separator, rather than silently falling
// back to an auto-derived label.
func ParsePersonalEnvs(raw string) ([]PersonalEnvPrefix, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	entries := strings.Split(raw, ",")
	prefixes := make([]PersonalEnvPrefix, 0, len(entries))
	for _, entry := range entries {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		parts := strings.SplitN(entry, ":", 2)
		if len(parts) != 2 || strings.TrimSpace(parts[0]) == "" || strings.TrimSpace(parts[1]) == "" {
			return nil, fmt.Errorf("PERSONAL_ENVS entry %q is not in PREFIX:FriendlyLabel format", entry)
		}
		prefixes = append(prefixes, PersonalEnvPrefix{
			Prefix:        strings.TrimSpace(parts[0]),
			FriendlyLabel: strings.TrimSpace(parts[1]),
		})
	}
	return prefixes, nil
}

// personalEnvMountAndPathPrefix reads the Vault mount/path-prefix config for
// personal secrets, applying the documented defaults (Section 2).
func personalEnvMountAndPathPrefix() (mount, pathPrefix string) {
	mount = os.Getenv("PERSONAL_ENV_VAULT_MOUNT")
	if mount == "" {
		mount = "secret"
	}
	pathPrefix = os.Getenv("PERSONAL_ENV_VAULT_PATH_PREFIX")
	if pathPrefix == "" {
		pathPrefix = "personal-envs"
	}
	return mount, pathPrefix
}

// ResolveTrovoUsername looks up the caller's actual Trovo app username —
// AdminUser.Username for a Trovo admin, OrganizationMember.TrovoWalletUsername
// for an org member — using current_user_id/user_type from
// AllowOrgOrTrovoAdminNormalized. This is the identity used to derive the
// Vault key for PERSONAL SECRETS ONLY. It is deliberately NOT the identity
// used for managed secrets (VaultSignerAssignment.OwnerRefID, which stores
// OrganizationMember.ID, resolved from Email only once at
// assignment-creation time — never TrovoWalletUsername) — the two resolvers
// exist because these two flows use two different, unrelated identity
// fields for an org member, and conflating them is a correctness bug: an
// org member with no linked wallet can still own a signer slot, but cannot
// use personal secrets.
func ResolveTrovoUsername(c *gin.Context, db *gorm.DB) (username string, ok bool) {
	userType := c.GetString("user_type")
	currentUserID := c.GetString("current_user_id")

	if userType == "trovo_admin" {
		var admin middleware.AdminUser
		if err := db.First(&admin, "id = ?", currentUserID).Error; err != nil {
			return "", false
		}
		return admin.Username, true // always non-empty — not-null in schema
	}

	var member coreModels.OrganizationMember
	if err := db.First(&member, "id = ?", currentUserID).Error; err != nil {
		return "", false
	}
	if member.TrovoWalletUsername == nil || *member.TrovoWalletUsername == "" {
		return "", false // not linked — caller must show the friendly message below
	}
	return *member.TrovoWalletUsername, true
}

// FriendlyMessageNotLinked is the exact caller-facing copy for the
// "wallet not linked" gate (Section 5b/6) — always "personal secrets" in any
// user-facing text, never "personal env(s)."
const FriendlyMessageNotLinked = "Connect your profile to the Trovo app to enable personal secrets."

// DeriveVaultKey builds the full personal-secret Vault key name:
// {PREFIX}_{USERNAME}, uppercased — e.g. AUTO_APPROVE + tunde -> AUTO_APPROVE_TUNDE.
// username here is always AdminUser.Username or OrganizationMember.TrovoWalletUsername
// (ResolveTrovoUsername's result) — never owner_id, never Email.
func DeriveVaultKey(prefix, username string) string {
	return strings.ToUpper(prefix) + "_" + strings.ToUpper(username)
}

// PersonalEnvListItem is one row of GET /me/vault-signer/personal-envs.
type PersonalEnvListItem struct {
	Prefix        string `json:"prefix"`
	FriendlyLabel string `json:"friendlyLabel"`
	VaultKey      string `json:"vaultKey"`
	Exists        bool   `json:"exists"`
}

// ListPersonalEnvs answers GET /me/vault-signer/personal-envs — never
// touches a value, only Vault metadata existence checks (Section 5b).
func ListPersonalEnvs(ctx context.Context, vc *vaultapi.Client, prefixes []PersonalEnvPrefix, username string) ([]PersonalEnvListItem, error) {
	mount, pathPrefix := personalEnvMountAndPathPrefix()
	items := make([]PersonalEnvListItem, 0, len(prefixes))
	for _, p := range prefixes {
		vaultKey := DeriveVaultKey(p.Prefix, username)
		path := fmt.Sprintf("%s/%s", pathPrefix, vaultKey)
		exists, err := vaultclient.Exists(ctx, vc, mount, path)
		if err != nil {
			return nil, err
		}
		items = append(items, PersonalEnvListItem{
			Prefix:        p.Prefix,
			FriendlyLabel: p.FriendlyLabel,
			VaultKey:      vaultKey,
			Exists:        exists,
		})
	}
	return items, nil
}

// CreatePersonalEnv validates newValue through the Section 5c pipeline, then
// creates the personal secret (409 if it already exists — Section 5b).
func CreatePersonalEnv(ctx context.Context, vc *vaultapi.Client, prefix, username, newValueRaw string) (vaultKey string, version int, err error) {
	if _, err := ValidateSignerValue(newValueRaw); err != nil {
		return "", 0, err
	}
	mount, pathPrefix := personalEnvMountAndPathPrefix()
	vaultKey = DeriveVaultKey(prefix, username)
	path := fmt.Sprintf("%s/%s", pathPrefix, vaultKey)
	version, err = vaultclient.Create(ctx, vc, mount, path, strings.TrimSpace(newValueRaw))
	if err != nil {
		return "", 0, err
	}
	return vaultKey, version, nil
}

// DeletePersonalEnv fully purges the personal secret (404 if it doesn't
// exist — Section 5b). Returns the vault key for audit logging.
func DeletePersonalEnv(ctx context.Context, vc *vaultapi.Client, prefix, username string) (vaultKey string, err error) {
	mount, pathPrefix := personalEnvMountAndPathPrefix()
	vaultKey = DeriveVaultKey(prefix, username)
	path := fmt.Sprintf("%s/%s", pathPrefix, vaultKey)
	exists, err := vaultclient.Exists(ctx, vc, mount, path)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", vaultclient.ErrNotFound
	}
	if err := vaultclient.Delete(ctx, vc, mount, path); err != nil {
		return "", err
	}
	return vaultKey, nil
}

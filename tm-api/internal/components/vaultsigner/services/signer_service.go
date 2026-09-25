package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	vaultsignermodels "admin-panel-dashboard/internal/components/vaultsigner/models"
	"admin-panel-dashboard/internal/components/vaultsigner/vaultclient"
	tErrors "admin-panel-dashboard/internal/errors"
	"admin-panel-dashboard/internal/evmkeypair"
	"admin-panel-dashboard/internal/gnosissafe"
	"admin-panel-dashboard/internal/middleware"
	coreModels "admin-panel-dashboard/internal/models"
	"admin-panel-dashboard/internal/network"

	"github.com/ecnepsnai/discord"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	vaultapi "github.com/hashicorp/vault/api"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// csvIndex converts a 1-based Signer Position to a 0-based Vault CSV index —
// the single translation point (Section 4) every Vault/on-chain touchpoint
// goes through, rather than re-deriving it locally.
func csvIndex(position int) int { return position - 1 }

// ResolveOwnerRefID resolves the human-friendly identifier an admin submits
// when creating/editing an assignment (username for trovo_admin, email for
// org_member — Section 6) to the real, stable OwnerRefID that gets stored
// (AdminUser.ID / OrganizationMember.ID — Section 4). This lookup happens
// once, at write time; it is never repeated for authorization checks, which
// compare current_user_id/user_type directly against the stored OwnerRefID.
func ResolveOwnerRefID(db *gorm.DB, ownerIdentifier, ownerType string) (ownerRefID string, err error) {
	switch ownerType {
	case "trovo_admin":
		var admin middleware.AdminUser
		if err := db.First(&admin, "username = ?", ownerIdentifier).Error; err != nil {
			return "", fmt.Errorf("no trovo admin with username %q: %w", ownerIdentifier, err)
		}
		return fmt.Sprint(admin.ID), nil
	case "org_member":
		var member coreModels.OrganizationMember
		if err := db.First(&member, "email = ?", ownerIdentifier).Error; err != nil {
			return "", fmt.Errorf("no organization member with email %q: %w", ownerIdentifier, err)
		}
		return member.ID, nil
	default:
		return "", fmt.Errorf("unknown owner type %q", ownerType)
	}
}

// ResolveOwnerLabel is the reverse lookup, resolved fresh at read time
// (Section 6) so a listing never shows a stale copy of a mutable field.
func ResolveOwnerLabel(db *gorm.DB, ownerRefID, ownerType string) (label string, err error) {
	switch ownerType {
	case "trovo_admin":
		var admin middleware.AdminUser
		if err := db.First(&admin, "id = ?", ownerRefID).Error; err != nil {
			return "", err
		}
		return admin.Username, nil
	case "org_member":
		var member coreModels.OrganizationMember
		if err := db.First(&member, "id = ?", ownerRefID).Error; err != nil {
			return "", err
		}
		return member.Email, nil
	default:
		return "", fmt.Errorf("unknown owner type %q", ownerType)
	}
}

// ErrOwnerAlreadyAssigned means this owner already has an assignment row on
// this managed secret — (ManagedSecretID, OwnerRefID) is unique (Section 4:
// VaultSignerAssignment doc comment). Surfaced as 409 by the handler.
var ErrOwnerAlreadyAssigned = errors.New("this owner is already assigned to this managed secret")

// isUniqueViolation reports whether err is a Postgres unique-constraint
// violation (SQLSTATE 23505) — the DB-level half of the
// (ManagedSecretID, OwnerRefID) guard. The app-level pre-checks in
// CreateAssignment/EditAssignment below are a friendly fast path, not the
// real guarantee: only the DB constraint is race-safe under concurrent
// requests, so both paths translate to the same ErrOwnerAlreadyAssigned
// rather than letting a raw DB error leak to the caller on the rare race.
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

// CreateAssignment creates an assignment with no Position — Position stays
// NULL until the owner submits their own signer value for the first time
// (ClaimPosition, called from the PUT handler). Creation never touches Vault
// or the chain at all now: it's a plain row insert, with no locking needed
// since there's no position to race over yet.
//
// (ManagedSecretID, OwnerRefID) is unique (Section 4) — the same owner can
// never hold two assignment rows on the same managed secret. Checked here
// first for a friendly error, and enforced underneath by the DB's own
// unique index (idx_vault_signer_assignment_secret_owner) in case two
// requests for the same owner race past the pre-check.
func CreateAssignment(db *gorm.DB, managedSecretID, ownerRefID, ownerType string) (*vaultsignermodels.VaultSignerAssignment, error) {
	var existing vaultsignermodels.VaultSignerAssignment
	err := db.Where("managed_secret_id = ? AND owner_ref_id = ?", managedSecretID, ownerRefID).
		First(&existing).Error
	if err == nil {
		return nil, ErrOwnerAlreadyAssigned
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	assignment := &vaultsignermodels.VaultSignerAssignment{
		ID:              uuid.NewString(),
		ManagedSecretID: managedSecretID,
		Position:        nil,
		OwnerRefID:      ownerRefID,
		OwnerType:       ownerType,
		AssignedAt:      time.Now(),
	}
	if err := db.Create(assignment).Error; err != nil {
		if isUniqueViolation(err) {
			return nil, ErrOwnerAlreadyAssigned
		}
		return nil, err
	}
	return assignment, nil
}

// ErrNoPositionAvailable means every entry in the managed secret's CSV is
// already claimed by another assignment — there is nothing left for this
// assignment to claim.
var ErrNoPositionAvailable = errors.New("no signer position is available to claim on this managed secret")

// ClaimPosition assigns the lowest unclaimed 1-based position to assignment,
// the first time its owner submits a real signer value. A no-op — returns
// the existing value — if the assignment already has one, which covers both
// "already claimed by this same owner" and "reassigned, but the previous
// owner had already claimed it" (Position is retained across reassignment,
// per EditAssignment).
//
// Every CSV entry is assumed to already be a real, pre-existing on-chain
// signer (Section 0) — claiming a position never grows the CSV or writes a
// placeholder into it; it only decides *which already-existing* entry this
// assignment now controls. If every entry is already claimed by some other
// assignment, there is nothing left to claim (ErrNoPositionAvailable).
//
// SELECT ... FOR UPDATE on the managed secret row serializes concurrent
// first-time claims on the same managed secret, so two different owners'
// first PUT can't race for the same lowest-available position.
func ClaimPosition(ctx context.Context, db *gorm.DB, vc *vaultapi.Client, assignmentID string) (position int, err error) {
	err = db.Transaction(func(tx *gorm.DB) error {
		var assignment vaultsignermodels.VaultSignerAssignment
		if err := tx.First(&assignment, "id = ?", assignmentID).Error; err != nil {
			return err
		}
		if assignment.Position != nil {
			position = *assignment.Position
			return nil
		}

		var secret vaultsignermodels.VaultSignerManagedSecret
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&secret, "id = ?", assignment.ManagedSecretID).Error; err != nil {
			return err
		}

		csvValue, _, err := vaultclient.ReadCSV(ctx, vc, secret)
		if err != nil {
			return err
		}
		totalPositions := len(strings.Split(csvValue, ","))

		var claimedPositions []int
		if err := tx.Model(&vaultsignermodels.VaultSignerAssignment{}).
			Where("managed_secret_id = ? AND position IS NOT NULL", assignment.ManagedSecretID).
			Pluck("position", &claimedPositions).Error; err != nil {
			return err
		}
		claimed := make(map[int]bool, len(claimedPositions))
		for _, p := range claimedPositions {
			claimed[p] = true
		}

		for p := 1; p <= totalPositions; p++ {
			if !claimed[p] {
				position = p
				break
			}
		}
		if position == 0 {
			return ErrNoPositionAvailable
		}

		return tx.Model(&vaultsignermodels.VaultSignerAssignment{}).
			Where("id = ?", assignmentID).UpdateColumn("position", position).Error
	})
	if err != nil {
		return 0, err
	}
	return position, nil
}

// EditAssignment reassigns the owner of an existing assignment, looked up by
// its own ID (not by position — position may still be null, or may belong
// to an owner who never logged in to claim it). Position is always retained
// exactly as-is, whatever it currently holds. This never touches Vault or
// the chain — whatever Stellar key currently sits at this position (if any
// has been claimed) doesn't change just because the DB row's owner changed;
// the new owner picks it up the normal way, through
// PUT /me/vault-signer/secrets/:secretId, which already runs the full
// validation/swap pipeline (Section 6) and, if this is their first time,
// ClaimPosition first.
//
// (ManagedSecretID, OwnerRefID) is unique (Section 4): reassigning to an
// owner who already holds a different assignment on this same managed
// secret is rejected the same way a duplicate create is — checked here
// first (excluding this assignment's own row, since keeping the same owner
// is always a no-op save, not a conflict with itself), and backstopped by
// the DB's unique index for the concurrent case.
func EditAssignment(db *gorm.DB, assignmentID string, newOwnerRefID, newOwnerType string) (*vaultsignermodels.VaultSignerAssignment, error) {
	var assignment vaultsignermodels.VaultSignerAssignment
	if err := db.First(&assignment, "id = ?", assignmentID).Error; err != nil {
		return nil, err
	}

	if newOwnerRefID != assignment.OwnerRefID {
		var existing vaultsignermodels.VaultSignerAssignment
		err := db.Where("managed_secret_id = ? AND owner_ref_id = ? AND id <> ?",
			assignment.ManagedSecretID, newOwnerRefID, assignmentID).First(&existing).Error
		if err == nil {
			return nil, ErrOwnerAlreadyAssigned
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	assignment.OwnerRefID = newOwnerRefID
	assignment.OwnerType = newOwnerType
	if err := db.Save(&assignment).Error; err != nil {
		if isUniqueViolation(err) {
			return nil, ErrOwnerAlreadyAssigned
		}
		return nil, err
	}
	return &assignment, nil
}

// DeleteAssignmentResult is what the handler needs to build the DELETE
// response (Section 6).
type DeleteAssignmentResult struct {
	DeletedPosition  int
	RenumberedCount  int
	VaultCollapsed   bool
	OnChainAttempted bool
	BaseTxHash       string
	BaseTxStatus     string // "success" | "failed", only meaningful if OnChainAttempted
}

// ErrCSVFloorViolation means the deletion would bring the CSV below the
// minimum entry count (Section 4).
var ErrCSVFloorViolation = errors.New("deleting this assignment would bring the CSV below the minimum entry count")

// ErrVaultUnreachableReverted is returned when the on-chain removal
// succeeded but the follow-up Vault write couldn't be reached — the removal
// was reversed on-chain and every change was rolled back (Section 5e).
var ErrVaultUnreachableReverted = &tErrors.CustomError{
	Err:        "vault-unreachable-changes-reverted",
	ErrMessage: "Vault service could not be reached; all changes have been reverted.",
	Code:       http.StatusBadGateway,
}

// DeleteAssignment is Section 5e's full delete flow, looked up by the
// assignment's own ID (not by position — position may be null). An
// assignment that never had its position claimed deletes with nothing more
// than a row delete: no CSV entry was ever assigned to it, so there is
// nothing to renumber, collapse, or float a floor check against. Everything
// below the early return only applies once a position has actually been
// claimed.
//
// For a claimed position: a single CSV read to decide light vs. full path,
// DB delete + renumber, and — for a real key with enough redundancy — an
// on-chain removal attempted *before* the Vault write, with an on-chain
// reversal if the Vault write then fails.
//
// A refinement made while implementing this against Section 4/5d's existing
// active/spare distinction: the on-chain trigger additionally requires the
// deleted position itself to be within the wallet's active range
// (index < ActiveSigningCount) — a spare position's key, however valid, was
// never registered on-chain to begin with (Section 4a/5d), so there is
// nothing there to remove regardless of how many valid keys exist elsewhere.
// validKeypairCount itself still counts across the whole CSV, matching the
// request's literal wording — it's a capacity/redundancy check, not a tally
// of live on-chain signers; restricting it to the active range would cap it
// at ActiveSigningCount and make "> ActiveSigningCount" unsatisfiable.
func DeleteAssignment(ctx context.Context, db *gorm.DB, vc *vaultapi.Client, assignmentID string) (*DeleteAssignmentResult, error) {
	var assignment vaultsignermodels.VaultSignerAssignment
	if err := db.First(&assignment, "id = ?", assignmentID).Error; err != nil {
		return nil, err
	}

	if assignment.Position == nil {
		// Never claimed a position — nothing was ever written to Vault or
		// the chain for this assignment. A plain delete, no renumbering, no
		// floor check (there's no CSV entry to collapse against).
		if err := db.Delete(&assignment).Error; err != nil {
			return nil, err
		}
		return &DeleteAssignmentResult{}, nil
	}
	position := *assignment.Position
	managedSecretID := assignment.ManagedSecretID

	var secret vaultsignermodels.VaultSignerManagedSecret
	if err := db.First(&secret, "id = ?", managedSecretID).Error; err != nil {
		return nil, err
	}

	// Step 1/2: read the CSV once — reused for the floor check, the
	// light/full-path decision, the collapse, and (if needed) the on-chain
	// reversal snapshot. The floor check is against the CSV's own live
	// length, not the assignment count — the two are no longer 1:1 now that
	// positions are claimed rather than auto-assigned at creation, so an
	// unclaimed assignment must never count toward "how many CSV entries
	// would remain."
	originalCSV, originalVersion, err := vaultclient.ReadCSV(ctx, vc, secret)
	if err != nil {
		return nil, err
	}
	parts := strings.Split(originalCSV, ",")
	if len(parts)-1 < vaultclient.MinCSVEntries {
		return nil, ErrCSVFloorViolation
	}
	index := csvIndex(position)
	if index < 0 || index >= len(parts) {
		return nil, vaultclient.ErrIndexOutOfRange
	}
	targetIsValidKey := IsValidKeypair(parts[index])

	// Counted across the whole CSV, not just the active range: this is a
	// capacity/redundancy check ("is it safe to shrink the active set by
	// one"), not a literal tally of live on-chain signers — restricting it
	// to the active range would cap the count at ActiveSigningCount and make
	// "> ActiveSigningCount" impossible to satisfy, silently disabling
	// on-chain removal entirely. The position-must-be-active requirement
	// below is a separate, independent gate.
	validKeypairCount := vaultclient.CountValidKeypairs(parts, IsValidKeypair)

	result := &DeleteAssignmentResult{DeletedPosition: position}

	if !targetIsValidKey {
		// Light path: defensive fallback only. Every CSV entry is assumed to
		// already be a real, pre-existing signer (Section 0), so a claimed
		// position's entry should always be a valid keypair by construction
		// — this branch exists in case that assumption is ever violated
		// (e.g. Vault drift outside this service's control), not because
		// it's an expected outcome of the normal claim/swap flow. No
		// on-chain call regardless of validKeypairCount.
		err := db.Transaction(func(tx *gorm.DB) error {
			if err := deleteAndRenumber(tx, managedSecretID, position, result); err != nil {
				return err
			}
			if _, err := vaultclient.CollapseCSV(ctx, vc, secret, originalCSV, originalVersion, index); err != nil {
				return err
			}
			result.VaultCollapsed = true
			return nil
		})
		return result, err
	}

	// Full path: the deleted position holds a real key.
	onChainNeeded := index < secret.ActiveSigningCount && validKeypairCount > secret.ActiveSigningCount
	removedKP, parseErr := evmkeypair.ParseFull(strings.TrimSpace(parts[index]))
	if parseErr != nil {
		return nil, fmt.Errorf("re-parsing target key: %w", parseErr) // can't happen — targetIsValidKey already confirmed this
	}
	removedAddress := common.HexToAddress(removedKP.Address())

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := deleteAndRenumber(tx, managedSecretID, position, result); err != nil {
			return err
		}

		client := network.GetBlockchainClient()
		var capturedThreshold *big.Int
		onChainRemovalHappened := false

		if onChainNeeded {
			result.OnChainAttempted = true
			safe, err := gnosissafe.New(client, secret.WalletAddress)
			if err != nil {
				return fmt.Errorf("resolving safe address: %w", err)
			}
			owners, err := safe.Owners(ctx)
			if err != nil {
				return fmt.Errorf("fetching safe owners: %w", err)
			}
			prevOwner, err := gnosissafe.PrevOwner(owners, removedAddress)
			if err != nil {
				return fmt.Errorf("signer %s is not currently registered on safe %s: %w", removedAddress.Hex(), secret.WalletAddress, err)
			}
			// The threshold is preserved as-is across the removal (not
			// lowered) — the Safe contract itself will revert if this is
			// invalid, matching this codebase's established "submit and let
			// the network reject invalid state" pattern.
			threshold, err := safe.Threshold(ctx)
			if err != nil {
				return fmt.Errorf("fetching safe threshold: %w", err)
			}
			capturedThreshold = threshold

			signers, err := GatherActiveSigners(ctx, vc, secret, index)
			if err != nil {
				return err
			}
			innerCalldata, err := gnosissafe.RemoveOwnerCalldata(prevOwner, removedAddress, threshold)
			if err != nil {
				return fmt.Errorf("encoding removeOwner call: %w", err)
			}
			execResult, submitErr := submitSafeOwnerChange(ctx, safe, innerCalldata, signers)
			if submitErr != nil {
				result.BaseTxStatus = "failed"
				return submitErr // nothing has touched Vault yet — a clean rollback
			}
			result.BaseTxHash = execResult.TxHash
			if !execResult.Success {
				result.BaseTxStatus = "failed"
				return fmt.Errorf("execTransaction for removeOwner reverted on-chain (tx %s)", execResult.TxHash)
			}
			result.BaseTxStatus = "success"
			onChainRemovalHappened = true
		}

		// Vault write — after on-chain removal succeeded, or after
		// determining it wasn't needed.
		if _, collapseErr := vaultclient.CollapseCSV(ctx, vc, secret, originalCSV, originalVersion, index); collapseErr == nil {
			result.VaultCollapsed = true
			return nil // commit
		} else if !onChainRemovalHappened {
			// Vault failed on its own, nothing on-chain to reverse.
			return collapseErr
		} else {
			return reverseOnChainRemoval(ctx, client, secret, index, removedAddress, capturedThreshold, collapseErr)
		}
	})

	return result, err
}

// reverseOnChainRemoval re-adds removedAddress to the safe at
// capturedThreshold — the compensating action for a successful on-chain
// removal followed by a Vault write that couldn't be reached (Section 5e).
// Safe has no "undo removeOwner" primitive, so reversal is a distinct
// addOwnerWithThreshold call, not a replay of the removal.
func reverseOnChainRemoval(ctx context.Context, client *ethclient.Client, secret vaultsignermodels.VaultSignerManagedSecret, index int, removedAddress common.Address, capturedThreshold *big.Int, collapseErr error) error {
	safe, safeErr := gnosissafe.New(client, secret.WalletAddress)
	if safeErr != nil {
		alertDoubleFailure(collapseErr, safeErr)
		return fmt.Errorf("vault write failed and could not resolve safe to reverse on-chain removal: %w", safeErr)
	}

	vc, vcErr := vaultclient.NewClient()
	if vcErr != nil {
		alertDoubleFailure(collapseErr, vcErr)
		return fmt.Errorf("vault write failed and could not build a vault client to gather signers for reversal: %w", vcErr)
	}
	signers, sigErr := GatherActiveSigners(ctx, vc, secret, index)
	if sigErr != nil {
		alertDoubleFailure(collapseErr, sigErr)
		return sigErr
	}

	innerCalldata, encErr := gnosissafe.AddOwnerCalldata(removedAddress, capturedThreshold)
	if encErr != nil {
		alertDoubleFailure(collapseErr, encErr)
		return fmt.Errorf("encoding addOwnerWithThreshold call to reverse removal: %w", encErr)
	}
	execResult, submitErr := submitSafeOwnerChange(ctx, safe, innerCalldata, signers)
	if submitErr != nil {
		// The genuine residual risk (Section 5e): on-chain removal
		// succeeded, Vault failed, and the reversal also failed. Bounded
		// retries belong here in production; for now this alerts loudly and
		// leaves the assignment uncommitted rather than guessing which
		// system to trust.
		alertDoubleFailure(collapseErr, submitErr)
		return fmt.Errorf("vault write failed and on-chain reversal also failed — manual reconciliation required: %w", submitErr)
	}
	if !execResult.Success {
		revertErr := fmt.Errorf("execTransaction for addOwnerWithThreshold (reversal) reverted on-chain (tx %s)", execResult.TxHash)
		alertDoubleFailure(collapseErr, revertErr)
		return fmt.Errorf("vault write failed and on-chain reversal also failed — manual reconciliation required: %w", revertErr)
	}

	// Reversal succeeded — the wallet is back to where it started. Returning
	// a non-nil error rolls back the DB transaction so it matches.
	return ErrVaultUnreachableReverted
}

// alertDoubleFailure surfaces the true residual-risk case loudly, following
// this codebase's established discord.Say pattern (internal/network/main.go)
// rather than letting it show up only as an audit-log row.
func alertDoubleFailure(vaultErr, secondErr error) {
	discord.WebhookURL = "https://discord.com/api/webhooks/824381163367170058/OXSX51RHd9DyLFbFipjdW3yXmyYC8SWwqd6HiXl6UtDzu75RxS1LzWA800hWereJJumw"
	if len(os.Getenv("EXPANSION_NETWORK_ERROR_WEBHOOK")) > 50 {
		discord.WebhookURL = os.Getenv("EXPANSION_NETWORK_ERROR_WEBHOOK")
	}
	msg := fmt.Sprintf("[vaultsigner.DeleteAssignment] NEEDS MANUAL RECONCILIATION: vault write failed (%v) and on-chain reversal also failed (%v)", vaultErr, secondErr)
	if err := discord.Say(msg); err != nil {
		log.Println("[vaultsigner.DeleteAssignment] discord alert itself failed:", err)
	}
}

func deleteAndRenumber(tx *gorm.DB, managedSecretID string, position int, result *DeleteAssignmentResult) error {
	if err := tx.Where("managed_secret_id = ? AND position = ?", managedSecretID, position).
		Delete(&vaultsignermodels.VaultSignerAssignment{}).Error; err != nil {
		return err
	}
	res := tx.Model(&vaultsignermodels.VaultSignerAssignment{}).
		Where("managed_secret_id = ? AND position > ?", managedSecretID, position).
		UpdateColumn("position", gorm.Expr("position - 1"))
	if res.Error != nil {
		return res.Error
	}
	result.RenumberedCount = int(res.RowsAffected)
	return nil
}

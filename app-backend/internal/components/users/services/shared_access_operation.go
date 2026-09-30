package users

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"trovo-wallet-api/internal/aa"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/evmkeypair"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"
)

// OperationSharedAccess is the wallet operation kind for shared access
// changes (the Safe's owners and threshold).
const OperationSharedAccess = "SHARED ACCESS"

// sharedAccessContext is what a shared access operation saves: the Safe
// state it moves the wallet to, checked again when it is signed.
type sharedAccessContext struct {
	OwnersHash string `json:"ownersHash"`
}

// On-chain, shared access is the wallet Safe's owners and threshold: the
// owner's key plus every APPROVER's key, with the number of approvals
// needed as the threshold (1, the owner alone, when there are no
// approvers). VIEW-ONLY and INITIATOR permissions are enforced by the
// backend only. An issuing wallet's linked distribution wallet mirrors it
// in the same operation, through the module the issuing wallet has on it.
func sharedAccessTarget(owner *userModels.User, approvers []*userModels.User, approvalsNeeded int) ([]common.Address, int64) {
	seen := map[common.Address]bool{}
	var owners []common.Address
	add := func(a string) {
		if !common.IsHexAddress(a) {
			return
		}
		addr := common.HexToAddress(a)
		if !seen[addr] {
			seen[addr] = true
			owners = append(owners, addr)
		}
	}
	add(owner.PrimarySigner)
	for _, u := range approvers {
		add(u.PrimarySigner)
	}
	threshold := int64(approvalsNeeded)
	if len(owners) == 1 || threshold < 1 {
		threshold = 1
	}
	return owners, threshold
}

// approversOf loads the users holding APPROVER permission on walletID, as
// seen through db (a transaction holding staged changes sees them).
func approversOf(walletID string, db *gorm.DB, gc *sharedconfig.GlobalConfig) ([]*userModels.User, error) {
	var perms []userModels.WalletPermission
	if err := db.Where("wallet_address = ? AND permission = ?", walletID, "APPROVER").Find(&perms).Error; err != nil {
		return nil, err
	}
	var out []*userModels.User
	for _, p := range perms {
		u, err := userModels.Username(p.TargetUsername).GetSimpleUser(gc.DB, gc)
		if err != nil {
			return nil, err
		}
		out = append(out, &u)
	}
	return out, nil
}

// sharedAccessCalls returns the calls moving wallet (and linked, if any) to
// the target owners and threshold; none when the chain already matches.
func sharedAccessCalls(ctx context.Context, wallet, linked *userModels.UserWallet, target []common.Address, threshold int64, gc *sharedconfig.GlobalConfig) ([]aa.Call, error) {
	aw, err := walletOwners(ctx, wallet, gc)
	if err != nil {
		return nil, &tErrors.ErrorTemporaryServerError{}
	}
	calls, err := aa.PlanOwnerChange(aw.Address, aw.Owners, aw.Threshold, target, threshold)
	if err != nil {
		return nil, &tErrors.CustomError{Param: "numberOfApprovalsNeeded", Err: "error-approvers-not-enough", ErrMessage: err.Error(), Code: http.StatusBadRequest}
	}
	if linked == nil {
		return calls, nil
	}
	managed := false
	for _, m := range linked.InitialModuleList() {
		if strings.EqualFold(common.HexToAddress(m).Hex(), aw.Address.Hex()) {
			managed = true
		}
	}
	if !managed {
		return nil, &tErrors.CustomError{Param: "linkedWalletAddress", Err: "error-linked-wallet-not-managed", ErrMessage: fmt.Sprintf("The linked wallet %v is not managed by %v, so its access cannot be mirrored.", linked.Alias, wallet.Alias), Code: http.StatusConflict}
	}
	lw, err := walletOwners(ctx, linked, gc)
	if err != nil {
		return nil, &tErrors.ErrorTemporaryServerError{}
	}
	if deployed, err := aa.Deployed(ctx, gc.BantuExpansionClient, lw.Address); err != nil || !deployed {
		return nil, &tErrors.CustomError{Param: "linkedWalletAddress", Err: "error-linked-wallet-not-activated", ErrMessage: fmt.Sprintf("The linked wallet %v is still being created. Please try again once it is active.", linked.Alias), Code: http.StatusConflict}
	}
	linkedCalls, err := aa.PlanOwnerChange(lw.Address, lw.Owners, lw.Threshold, target, threshold)
	if err != nil {
		return nil, &tErrors.ErrorTemporaryServerError{}
	}
	for _, c := range linkedCalls {
		calls = append(calls, aa.ViaModule(lw.Address, c))
	}
	return calls, nil
}

// prepareSharedAccessOperation builds the operation applying a shared access
// change. validity 0 is the default (signed by the owner right away); a
// wallet that already has approvers passes the longer approval validity.
func prepareSharedAccessOperation(ctx context.Context, initiator, owner *userModels.User, wallet *userModels.UserWallet, calls []aa.Call, target []common.Address, threshold int64, validity time.Duration, gc *sharedconfig.GlobalConfig) (*PreparedWalletOperation, error) {
	return PrepareWalletOperation(ctx, OperationSharedAccess, initiator, owner, wallet, calls, validity, sharedAccessContext{OwnersHash: aa.OwnersHash(target, threshold)}, gc)
}

// loadSharedAccessOperation loads the operation the app was given and
// checks it applies exactly the requested change.
func loadSharedAccessOperation(transaction string, wallet *userModels.UserWallet, target []common.Address, threshold int64, gc *sharedconfig.GlobalConfig) (*userModels.WalletOperation, *aa.Prepared, error) {
	rec, p, err := LoadWalletOperation(transaction, wallet.ID, OperationSharedAccess, gc)
	if err != nil {
		return nil, nil, err
	}
	var c sharedAccessContext
	if rec.Context == nil || json.Unmarshal([]byte(*rec.Context), &c) != nil || c.OwnersHash != aa.OwnersHash(target, threshold) {
		return nil, nil, &tErrors.CustomError{Param: "transaction", Err: "transaction mismatch", ErrMessage: "The transaction does not match this change. Please start again.", Code: http.StatusBadRequest}
	}
	return rec, p, nil
}

// verifyStatementSignature checks an app signature (base64) over a base64
// message - a shared access change that needs no on-chain change is
// authorized by the owner's signature over its statement.
func verifyStatementSignature(signer, messageB64, sigB64 string) error {
	msg, err := base64.StdEncoding.DecodeString(messageB64)
	if err != nil {
		return &tErrors.CustomError{Param: "transaction", Err: "error-invalid-transaction", ErrMessage: "The transaction is not valid."}
	}
	sig, err := DecodeAppSignature(sigB64)
	if err != nil {
		return err
	}
	if evmkeypair.VerifyPersonal(common.HexToAddress(signer), msg, sig) != nil {
		return &tErrors.CustomError{Param: "transactionSignature", Err: "error-invalid-signature", ErrMessage: "The signature does not authorize this change.", Code: http.StatusForbidden}
	}
	return nil
}

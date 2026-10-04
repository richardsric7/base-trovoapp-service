package users

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"trovo-wallet-api/internal/aa"
	userModels "trovo-wallet-api/internal/components/users/models"
	tErrors "trovo-wallet-api/internal/errors"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"
)

// Approval requests on Safe wallets carry the base64 hash of what the
// approvers sign in TransactionXdr: a pending wallet operation (payments,
// withdrawals, shared access changes that change the Safe's owners), or a
// statement (shared access changes that only touch VIEW-ONLY/INITIATOR
// permissions, which live in the database). Each approver signs it with
// their key, which must be one of the Safe's owners.

// approvalOperation returns the pending wallet operation an approval
// request is for, or nil when it is not operation-backed.
func approvalOperation(p *userModels.PendingAuth, gc *sharedconfig.GlobalConfig) (*userModels.WalletOperation, *aa.Prepared) {
	hash, err := base64.StdEncoding.DecodeString(strings.TrimSpace(p.TransactionXdr))
	if err != nil || len(hash) != 32 {
		return nil, nil
	}
	var rec userModels.WalletOperation
	if err := gc.DB.Where("id = ? AND wallet_address = ?", common.BytesToHash(hash).Hex(), p.WalletAddress).First(&rec).Error; err != nil {
		return nil, nil
	}
	prep, err := aa.UnmarshalPrepared([]byte(rec.Prepared))
	if err != nil {
		return nil, nil
	}
	return &rec, prep
}

// statementApproval reports whether an approval request is authorized by
// signatures over a statement only (nothing to submit on-chain).
func statementApproval(p *userModels.PendingAuth) bool {
	return p.TransactionType == "MODIFY SHARED ACCESS" || p.TransactionType == "DISABLE SHARED ACCESS"
}

// checkApprover refuses an approver whose key is not an owner of the Safe
// as the operation was built (e.g. an approver whose account was recovered
// onto a new key that has not been re-added yet).
func checkApprover(prep *aa.Prepared, signer string) error {
	for _, o := range prep.Owners {
		if strings.EqualFold(o.Hex(), signer) {
			return nil
		}
	}
	return &tErrors.CustomError{Param: "id", Err: "error-signer-not-co-signer", ErrMessage: "Your key is not a co-signer of this wallet, so your approval cannot authorize this transaction. Please contact the wallet owner.", Code: http.StatusForbidden}
}

// approvalSignatures collects the approvers' signatures on an approval
// request (through db, so a signature saved in the same transaction counts).
func approvalSignatures(pendingAuthID string, db *gorm.DB) ([]aa.OwnerSignature, error) {
	var rows []userModels.PendingTransactionSignature
	if err := db.Where("pending_auth_id = ?", pendingAuthID).Find(&rows).Error; err != nil {
		return nil, err
	}
	var sigs []aa.OwnerSignature
	for _, r := range rows {
		sig, err := DecodeAppSignature(r.TransactionWithSignature)
		if err != nil {
			continue
		}
		sigs = append(sigs, aa.OwnerSignature{Owner: common.HexToAddress(r.ApproverSignerAddress), Signature: sig})
	}
	return sigs, nil
}

// notifyApprovalCompleted tells everyone with access to the wallet that an
// approval request was completed and submitted.
func notifyApprovalCompleted(signerUser *userModels.User, p *userModels.PendingAuth, wallet *userModels.UserWallet, gc *sharedconfig.GlobalConfig) {
	notified := map[string]bool{}
	dataPayload := map[string]string{"route": "pendingApproval"}
	title := fmt.Sprintf("%v completed the %v approval on wallet %v!", signerUser.Username, p.TransactionType, wallet.Alias)
	body := fmt.Sprintf("%v completed the %v request:\n%v", signerUser.Username, p.TransactionType, p.Description)
	for _, v := range wallet.GetPermissionList(gc.DB) {
		u, err := userModels.Username(v.TargetUsername).GetSimpleUser(gc.DB, gc)
		if err != nil || u.PushNotificationToken == nil || notified[*u.PushNotificationToken] {
			continue
		}
		notified[*u.PushNotificationToken] = true
		u.SendPushMessage(title, body, "", dataPayload, gc)
		u.InvalidateUserCache(gc)
	}
	wallet.InvalidateUserCache(gc)
}

package users

import (
	"fmt"
	"strings"
	"time"
	userModels "trovo-wallet-api/internal/components/users/models"
	"trovo-wallet-api/internal/sharedconfig"
)

// StableRailInitiateOfframp asks Stablerail to pay cNGN held in the user's
// Stablerail wallet to their bank account (/cngnofframp).
func StableRailInitiateOfframp(reqData userModels.StablerailOfframpRequest, gc *sharedconfig.GlobalConfig) (*userModels.StableRailOfframpResponse, error) {
	var offrampResp userModels.StableRailOfframpResponse
	if err := stablerailPost("/cngnofframp", reqData, &offrampResp, gc); err != nil {
		return nil, err
	}
	return &offrampResp, nil
}

// offrampFinalStatuses are the final offramp statuses (lower case);
// DEPOSIT_FAILED is final too (nothing left the wallet).
var offrampFinalStatuses = []string{"completed", "complete", "success", "successful", "failed", "failure", "cancelled", "canceled", "rejected", "reversed", "deposit_failed"}

// offrampFinal reports whether an offramp status is final.
func offrampFinal(status string) bool {
	s := strings.ToLower(strings.TrimSpace(status))
	for _, f := range offrampFinalStatuses {
		if s == f {
			return true
		}
	}
	return false
}

// UpdateOfframpStatuses moves bank withdrawals along: transfers to the
// Stablerail wallet that failed are marked so, deposited ones are paid out,
// and Stablerail is asked for the status of the rest until they finish, when
// the user is told.
func UpdateOfframpStatuses(gc *sharedconfig.GlobalConfig) {
	var pending []userModels.StablerailOfframp
	// finished ones are left out here, so they cannot crowd out the rest
	gc.DB.Where("created_at > ? AND LOWER(status) NOT IN ?", time.Now().Add(-30*24*time.Hour), offrampFinalStatuses).Order("created_at").Limit(200).Find(&pending)
	for _, o := range pending {
		switch o.Status {
		case userModels.OfframpDepositing:
			checkOfframpDeposit(o, gc)
			continue
		case userModels.OfframpDeposited, userModels.OfframpRequestFailed:
			if o.Attempts < maxOfframpRequestAttempts {
				requestOfframpPayout(o.ID, gc)
			}
			continue
		case userModels.OfframpRequesting:
			// an instance stopped while asking; ask again (Stablerail will not
			// pay the same deposit twice)
			if time.Since(o.UpdatedAt) > 5*time.Minute {
				gc.DB.Model(&userModels.StablerailOfframp{}).Where("id = ? AND status = ?", o.ID, userModels.OfframpRequesting).Update("status", userModels.OfframpRequestFailed)
			}
			continue
		case userModels.OfframpDepositFailed:
			continue
		}
		if offrampFinal(o.Status) {
			continue
		}
		requestID := nonEmpty(o.RequestID, o.ID)
		res, err := GetOfframpStatus(requestID, gc)
		if err != nil || res == nil || res.Data.Status == "" {
			continue
		}
		status := res.Data.Status
		if strings.EqualFold(status, o.Status) {
			continue
		}
		gc.DB.Model(&userModels.StablerailOfframp{}).Where("id = ?", o.ID).Update("status", status)
		if !offrampFinal(status) {
			continue
		}
		u, err := userModels.Username(o.TrovoUsername).GetSimpleUser(gc.DB, gc)
		if err != nil {
			continue
		}
		title, body := "Bank withdrawal completed", fmt.Sprintf("Your withdrawal of %v %v to your bank account has been paid.", o.BaseAmount, o.Ticker)
		if f := strings.ToUpper(status); f != "COMPLETED" && f != "COMPLETE" && f != "SUCCESS" && f != "SUCCESSFUL" {
			title, body = "Bank withdrawal failed", fmt.Sprintf("Your withdrawal of %v %v to your bank account could not be completed (%v). Please contact support.", o.BaseAmount, o.Ticker, status)
			gc.LogDiscordFailedRequest(fmt.Sprintf("[UpdateOfframpStatuses] offramp %v (%v) of %v ended %v", o.ID, requestID, o.TrovoUsername, status))
		}
		u.SendPushMessage(title, body, "", map[string]string{"route": "basicTransactionHistory"}, gc)
	}
}

// checkOfframpDeposit marks a withdrawal whose transfer to the Stablerail
// wallet failed (a mined transfer is recorded by recordBankWithdrawalDeposit).
func checkOfframpDeposit(o userModels.StablerailOfframp, gc *sharedconfig.GlobalConfig) {
	if o.UserOpHash == nil {
		return
	}
	var op userModels.WalletOperation
	if gc.DB.Where("user_op_hash = ?", *o.UserOpHash).First(&op).Error != nil {
		return
	}
	mined := op.Status == userModels.WalletOperationIncluded && op.Success != nil
	if mined && *op.Success {
		// mined, but the hook did not record it (e.g. the instance stopped)
		gc.DB.Model(&userModels.StablerailOfframp{}).Where("id = ? AND status = ?", o.ID, userModels.OfframpDepositing).
			Updates(map[string]interface{}{"status": userModels.OfframpDeposited, "tx_hash": derefString(op.TxHash)})
		return
	}
	failed := op.Status == userModels.WalletOperationFailed || mined
	if !failed {
		return
	}
	gc.DB.Model(&userModels.StablerailOfframp{}).Where("id = ? AND status = ?", o.ID, userModels.OfframpDepositing).
		Updates(map[string]interface{}{"status": userModels.OfframpDepositFailed, "error": derefString(op.Error)})
	if u, err := userModels.Username(o.TrovoUsername).GetSimpleUser(gc.DB, gc); err == nil {
		u.SendPushMessage("Bank withdrawal failed", fmt.Sprintf("Your withdrawal of NGN %v did not go through and nothing left your wallet.", o.BaseAmount), "", map[string]string{"route": "basicTransactionHistory"}, gc)
	}
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

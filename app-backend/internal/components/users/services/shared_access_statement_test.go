package users

import (
	"encoding/base64"
	"testing"

	userModels "trovo-wallet-api/internal/components/users/models"
	"trovo-wallet-api/internal/evmkeypair"

	"github.com/ethereum/go-ethereum/common"
)

func TestSharedAccessAuthorizationPayloadBindsTheChange(t *testing.T) {
	perms := []userModels.WalletPermissionInfo{
		{TargetUsername: "Bob", Permission: "approver"},
		{TargetUsername: "alice", Permission: "INITIATOR"},
	}
	reordered := []userModels.WalletPermissionInfo{perms[1], perms[0]}

	a := sharedAccessAuthorizationPayload("0xWallet", 2, perms)
	if a != sharedAccessAuthorizationPayload("0xWallet", 2, reordered) {
		t.Fatal("statement must not depend on permission order")
	}
	for name, other := range map[string]string{
		"threshold":   sharedAccessAuthorizationPayload("0xWallet", 3, perms),
		"wallet":      sharedAccessAuthorizationPayload("0xOther", 2, perms),
		"permissions": sharedAccessAuthorizationPayload("0xWallet", 2, perms[:1]),
	} {
		if other == a {
			t.Errorf("statement must change when the %s changes", name)
		}
	}

	// the apps decode the payload and personal_sign it; the backend verifies
	signer, _ := evmkeypair.Random()
	statement, err := base64.StdEncoding.DecodeString(a)
	if err != nil {
		t.Fatal(err)
	}
	sig, _ := signer.SignBase64(statement)
	raw, _ := base64.StdEncoding.DecodeString(sig)
	if err := evmkeypair.VerifyPersonal(common.HexToAddress(signer.Address()), statement, raw); err != nil {
		t.Fatalf("signature over the statement must verify: %v", err)
	}
	other, _ := base64.StdEncoding.DecodeString(sharedAccessAuthorizationPayload("0xWallet", 3, perms))
	if evmkeypair.VerifyPersonal(common.HexToAddress(signer.Address()), other, raw) == nil {
		t.Fatal("a signature must not authorize a different change")
	}
}

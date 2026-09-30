package users

import (
	"encoding/base64"
	"strings"
	"testing"

	"trovo-wallet-api/internal/aa"
	userModels "trovo-wallet-api/internal/components/users/models"
	"trovo-wallet-api/internal/evmkeypair"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/ethereum/go-ethereum/common"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestSharedAccessTarget(t *testing.T) {
	owner := &userModels.User{PrimarySigner: "0x000000000000000000000000000000000000000a"}
	b := &userModels.User{PrimarySigner: "0x000000000000000000000000000000000000000b"}
	c := &userModels.User{PrimarySigner: "0x000000000000000000000000000000000000000c"}

	owners, th := sharedAccessTarget(owner, []*userModels.User{b, c}, 2)
	if len(owners) != 3 || th != 2 || owners[0] != common.HexToAddress(owner.PrimarySigner) {
		t.Fatalf("approvers: %v/%d", owners, th)
	}
	// the owner approving too is not a second copy of their key
	owners, th = sharedAccessTarget(owner, []*userModels.User{owner, b}, 1)
	if len(owners) != 2 || th != 1 {
		t.Fatalf("owner as approver: %v/%d", owners, th)
	}
	// no approvers: the owner alone
	owners, th = sharedAccessTarget(owner, nil, 3)
	if len(owners) != 1 || th != 1 {
		t.Fatalf("disabled: %v/%d", owners, th)
	}
}

func TestVerifyStatementSignature(t *testing.T) {
	k, _ := evmkeypair.Random()
	other, _ := evmkeypair.Random()
	statement := sharedAccessAuthorizationPayload("0xwallet", 1, []userModels.WalletPermissionInfo{{TargetUsername: "Bob", Permission: "view-only"}})
	msg, _ := base64.StdEncoding.DecodeString(statement)
	sig, _ := evmkeypair.SignPersonal(k.PrivateKey(), msg)
	sigB64 := base64.StdEncoding.EncodeToString(sig)
	if err := verifyStatementSignature(k.Address(), statement, sigB64); err != nil {
		t.Fatal(err)
	}
	if err := verifyStatementSignature(other.Address(), statement, sigB64); err == nil {
		t.Fatal("a signature by another key must be refused")
	}
}

// An approval request carries the operation's hash; each approver's
// signature is collected for it, and only the Safe's owners may approve.
func TestApprovalOperationAndSignatures(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&userModels.WalletOperation{}, &userModels.PendingTransactionSignature{}); err != nil {
		t.Fatal(err)
	}
	gc := &sharedconfig.GlobalConfig{DB: db}

	ap1, _ := evmkeypair.Random()
	ap2, _ := evmkeypair.Random()
	outsider, _ := evmkeypair.Random()
	hash := common.HexToHash("0x" + strings.Repeat("5a", 32))
	prep := &aa.Prepared{SafeOpHash: hash, Owners: []common.Address{common.HexToAddress(ap1.Address()), common.HexToAddress(ap2.Address())}, Threshold: 2}
	raw, err := prep.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	wallet := "0x00000000000000000000000000000000000000f0"
	db.Create(&userModels.WalletOperation{ID: hash.Hex(), WalletAddress: wallet, Kind: OperationPayment, Initiator: "alice", Status: userModels.WalletOperationPending, Prepared: string(raw)})

	p := &userModels.PendingAuth{ID: "auth-1", WalletAddress: wallet, TransactionType: "PAYMENT", TransactionXdr: base64.StdEncoding.EncodeToString(hash.Bytes())}
	rec, got := approvalOperation(p, gc)
	if rec == nil || got.Threshold != 2 {
		t.Fatal("the operation behind the request must be found")
	}
	if err := checkApprover(got, outsider.Address()); err == nil {
		t.Fatal("a non-owner must not approve")
	}
	if err := checkApprover(got, strings.ToLower(ap2.Address())); err != nil {
		t.Fatal(err)
	}
	// another wallet's request cannot point at this operation
	if rec, _ := approvalOperation(&userModels.PendingAuth{WalletAddress: "0xother", TransactionXdr: p.TransactionXdr}, gc); rec != nil {
		t.Fatal("operation of another wallet")
	}

	for i, k := range []*evmkeypair.Full{ap1, ap2} {
		sig, _ := evmkeypair.SignPersonal(k.PrivateKey(), hash.Bytes())
		sigB64 := base64.StdEncoding.EncodeToString(sig)
		if err := verifyStatementSignature(k.Address(), p.TransactionXdr, sigB64); err != nil {
			t.Fatal(err)
		}
		db.Create(&userModels.PendingTransactionSignature{ID: string(rune('a' + i)), PendingAuthID: p.ID, Approver: "ap" + string(rune('1'+i)), ApproverSignerAddress: k.Address(), TransactionWithSignature: sigB64})
	}
	sigs, err := approvalSignatures(p.ID, db)
	if err != nil || len(sigs) != 2 {
		t.Fatalf("signatures %v %v", sigs, err)
	}
	// they satisfy the Safe's threshold
	if _, err := aa.AssembleSignature(hash, 0, 0, got.Owners, int(got.Threshold), sigs); err != nil {
		t.Fatal(err)
	}
	if !statementApproval(&userModels.PendingAuth{TransactionType: "DISABLE SHARED ACCESS"}) || statementApproval(p) {
		t.Fatal("statement approvals are shared access changes only")
	}
}

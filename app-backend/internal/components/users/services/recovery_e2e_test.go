package users

import (
	"context"
	"encoding/json"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"trovo-wallet-api/internal/aa"
	"trovo-wallet-api/internal/cache"
	userModels "trovo-wallet-api/internal/components/users/models"
	"trovo-wallet-api/internal/db"
	"trovo-wallet-api/internal/evmkeypair"
	"trovo-wallet-api/internal/gnosissafe"
	"trovo-wallet-api/internal/sharedconfig"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestAccountRecoveryOnLocalChain runs account recovery against real
// contracts: the local stack (paymaster/contracts/scripts/local-stack.js,
// with its dev bundler) and the recovery module on the same node
// (recovery/contracts/scripts/local-recovery.js, 60 second period).
// Skipped unless
//
//	AA_LOCAL_STACK=$PWD/../../../../../paymaster/contracts/deployments/local-stack.json
//	AA_LOCAL_RECOVERY=$PWD/../../../../../recovery/contracts/deployments/local-recovery.json
//
// (absolute paths) are set. A user turns recovery on for their primary
// wallet and a sub-wallet (in each wallet's first operation), shares the
// sub-wallet with a co-signer (which removes the guardian from it), then
// recovers: the platform's guardian Safe starts the recovery, the user
// cancels it with their key, it is started again and finalized after the
// period, the account moves to the new key, and the shared wallet gets a
// REPLACE SIGNER request the co-signer approves.
func TestAccountRecoveryOnLocalChain(t *testing.T) {
	stackFile, recoveryFile := os.Getenv("AA_LOCAL_STACK"), os.Getenv("AA_LOCAL_RECOVERY")
	if stackFile == "" || recoveryFile == "" {
		t.Skip("AA_LOCAL_STACK / AA_LOCAL_RECOVERY not set")
	}
	read := func(f string) map[string]string {
		raw, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		m := map[string]string{}
		if err := json.Unmarshal(raw, &m); err != nil {
			t.Fatal(err)
		}
		return m
	}
	stack, rec := read(stackFile), read(recoveryFile)
	for _, k := range []string{"ENTRYPOINT_ADDRESS", "SAFE_PROXY_FACTORY_ADDRESS", "SAFE_SINGLETON_ADDRESS", "SAFE_MODULE_SETUP_ADDRESS", "SAFE_4337_MODULE_ADDRESS", "SAFE_MULTISEND_CALL_ONLY_ADDRESS", "BUNDLER_URL"} {
		t.Setenv(k, stack[k])
	}
	t.Setenv("BASE_CHAIN_ID", stack["BASE_CHAIN_ID"])
	t.Setenv("RECOVERY_MODULE_ADDRESS", rec["RECOVERY_MODULE_ADDRESS"])
	t.Setenv("RECOVERY_PERIOD_SECONDS", rec["RECOVERY_PERIOD_SECONDS"])
	t.Setenv("WALLET_OPERATION_GAS_BUFFER_PERCENT", "300")
	module := common.HexToAddress(rec["RECOVERY_MODULE_ADDRESS"])

	ctx := context.Background()
	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		t.Fatal(err)
	}
	chainID, _ := client.ChainID(ctx)
	gdb, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	db.MigrateDB(gdb)
	gc := &sharedconfig.GlobalConfig{DB: gdb, BantuExpansionClient: client, RedisCache: &cache.RedisCache{}}
	cfg := aa.ConfigFromEnv(chainID)

	admin, _ := evmkeypair.ParseFull("0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80") // hardhat #0
	fund := func(to common.Address) {
		t.Helper()
		h, err := sendValue(ctx, client, chainID, admin, to, big.NewInt(1e17))
		if err != nil {
			t.Fatal(err)
		}
		if err := gnosissafe.WaitSuccess(ctx, client, h); err != nil {
			t.Fatal(err)
		}
	}

	// the platform's guardian: a Safe owned by the recovery key
	guardianKey, _ := evmkeypair.Random()
	fund(common.HexToAddress(guardianKey.Address()))
	guardianSafe, err := gnosissafe.DeploySafe(ctx, client, chainID, admin, gnosissafe.DeployConfig{
		ProxyFactory: cfg.ProxyFactory, Singleton: cfg.Singleton,
		FallbackHandler: cfg.Safe4337Module, // any contract; the local stack has no CompatibilityFallbackHandler
	}, []common.Address{common.HexToAddress(guardianKey.Address())}, 1, big.NewInt(time.Now().UnixNano()))
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("ACCOUNT_RECOVERY_GUARDIAN_SIGNERS", "0x"+guardianKey.Seed())
	t.Setenv("ACCOUNT_RECOVERY_GUARDIAN_SAFE", guardianSafe.Hex())
	guardian, err := loadRecoveryGuardian(ctx, gc)
	if err != nil || guardian.Address != guardianSafe {
		t.Fatalf("guardian %+v: %v", guardian, err)
	}

	// --- the user: primary wallet and a sub-wallet, and a co-signer
	oldKey, _ := evmkeypair.Random()
	newKey, _ := evmkeypair.Random()
	bobKey, _ := evmkeypair.Random()
	user := userModels.User{ID: "u-alice", Username: "alice", PrimarySigner: strings.ToUpper(oldKey.Address()), AccountRecoveryEnabled: 1, HasSecurityQuestions: 1}
	user.Address = userModels.PrimarySafeDeployment(user.PrimarySigner).Address
	user.BuildPrimaryWallet()
	subD, _ := userModels.NewSubWalletSafeDeployment(user.PrimarySigner)
	sub := userModels.UserWallet{ID: subD.Address, Alias: "alice-savings", Signer: user.PrimarySigner, UserID: user.ID, InitialOwners: user.PrimarySigner, InitialThreshold: 1, SafeSaltNonce: subD.SaltNonce, SafeVersion: aa.SafeVersion}
	bob := userModels.User{ID: "u-bob", Username: "bob", PrimarySigner: strings.ToUpper(bobKey.Address())}
	bob.Address = userModels.PrimarySafeDeployment(bob.PrimarySigner).Address
	bob.BuildPrimaryWallet()
	for _, v := range []interface{}{&user, &bob, &user.UserWallets[0], &bob.UserWallets[0], &sub} {
		if err := gdb.Omit("UserWallets", "WalletsSharedWithUser").Create(v).Error; err != nil {
			t.Fatal(err)
		}
	}
	primary := user.UserWallets[0]
	fund(common.HexToAddress(primary.ID))
	fund(common.HexToAddress(sub.ID))

	// signs and submits the operation the backend prepared (as the app does)
	sign := func(label string, kind string, tx string, wallet string, keys ...*evmkeypair.Full) {
		t.Helper()
		r, p, err := LoadWalletOperation(tx, wallet, kind, gc)
		if err != nil {
			t.Fatalf("%s: %v", label, err)
		}
		var sigs []aa.OwnerSignature
		for _, k := range keys {
			sig, _ := evmkeypair.SignPersonal(k.PrivateKey(), p.SafeOpHash.Bytes())
			sigs = append(sigs, aa.OwnerSignature{Owner: common.HexToAddress(k.Address()), Signature: sig})
		}
		h, err := SubmitWalletOperation(ctx, r, p, sigs, gc)
		if err != nil {
			t.Fatalf("%s: submit: %v", label, err)
		}
		wctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
		rcpt, err := operationBuilder(gc).Bundler.(*aa.Bundler).WaitReceipt(wctx, common.HexToHash(h))
		if err != nil || !rcpt.Success {
			t.Fatalf("%s: %v (success=%v)", label, err, rcpt != nil && rcpt.Success)
		}
	}
	guarded := func(w string) bool {
		t.Helper()
		ok, err := aa.IsGuardian(ctx, client, module, common.HexToAddress(w), guardianSafe)
		if err != nil {
			t.Fatal(err)
		}
		return ok
	}

	// --- turn recovery on: one operation per wallet (both still to be
	// deployed by it)
	for _, w := range []userModels.UserWallet{primary, sub} {
		aw, err := walletOwners(ctx, &w, gc)
		if err != nil {
			t.Fatal(err)
		}
		op, err := PrepareWalletOperation(ctx, OperationRecoveryEnable, &user, &user, &w, aa.EnableRecoveryCalls(aw.Address, module, guardianSafe, false), 0, nil, gc)
		if err != nil {
			t.Fatal(err)
		}
		sign("enable "+w.Alias, OperationRecoveryEnable, op.Transaction, w.ID, oldKey)
		if !guarded(w.ID) {
			t.Fatalf("%s not covered", w.Alias)
		}
	}
	full, err := userModels.Username("alice").GetFullUser(gdb, gc)
	if err != nil {
		t.Fatal(err)
	}
	all, err := ownedRecoveryWallets(ctx, &full, common.HexToAddress(user.PrimarySigner), module, guardianSafe, gc)
	if err != nil || len(all) != 2 || !all[0].Guarded || !all[1].Guarded {
		t.Fatalf("covered wallets %+v: %v", all, err)
	}

	// --- share the sub-wallet with bob (1 approval): the same operation
	// removes the guardian, so the platform never replaces keys on it
	target := []common.Address{common.HexToAddress(user.PrimarySigner), common.HexToAddress(bob.PrimarySigner)}
	calls, err := sharedAccessCalls(ctx, &sub, nil, target, 1, gc)
	if err != nil {
		t.Fatal(err)
	}
	op, err := prepareSharedAccessOperation(ctx, &user, &user, &sub, calls, target, 1, 0, gc)
	if err != nil {
		t.Fatal(err)
	}
	sign("share", OperationSharedAccess, op.Transaction, sub.ID, oldKey)
	if guarded(sub.ID) {
		t.Fatal("a wallet with a co-signer must not stay covered")
	}
	gdb.Model(&userModels.UserWallet{}).Where("id = ?", sub.ID).Updates(map[string]interface{}{"shared_access_enabled": 1, "number_of_approvals_needed": 1})
	gdb.Create(&userModels.WalletPermission{ID: "perm-bob", WalletAddress: sub.ID, TargetUsername: "bob", Permission: "APPROVER"})

	if full, err = userModels.Username("alice").GetFullUser(gdb, gc); err != nil {
		t.Fatal(err)
	}
	all, err = ownedRecoveryWallets(ctx, &full, common.HexToAddress(user.PrimarySigner), module, guardianSafe, gc)
	if err != nil || len(all) != 1 || all[0].Wallet.ID != primary.ID {
		t.Fatalf("covered wallets after sharing %+v: %v", all, err)
	}

	// --- recovery started by the guardian, canceled by the user's key
	old, nw := common.HexToAddress(user.PrimarySigner), common.HexToAddress(newKey.Address())
	arl, err := startAccountRecovery(ctx, &full, old, nw, module, guardian, all, gc)
	if err != nil {
		t.Fatal(err)
	}
	if st := GetAccountRecoveryStatus("alice", newKey.Address(), gc); st.Status != userModels.AccountRecoveryPending || st.ExecuteAfter == nil {
		t.Fatalf("status %+v", st)
	}
	cancelPayload := &userModels.UserAccountRecoveryPayload{}
	if err := CancelAccountRecovery(&full, cancelPayload, gc); err != nil || len(cancelPayload.Transactions) != 1 {
		t.Fatalf("cancel prepare: %v (%d)", err, len(cancelPayload.Transactions))
	}
	sign("cancel", OperationRecoveryCancel, cancelPayload.Transactions[0], primary.ID, oldKey)
	ProcessAccountRecoveries(ctx, gc)
	gdb.First(arl, "id = ?", arl.ID)
	if arl.Status != userModels.AccountRecoveryCanceled {
		t.Fatalf("after cancel: %v", arl.Status)
	}

	// --- started again; finalized after the period; the account moves
	arl, err = startAccountRecovery(ctx, &full, old, nw, module, guardian, all, gc)
	if err != nil {
		t.Fatal(err)
	}
	ProcessAccountRecoveries(ctx, gc) // too early: nothing happens
	if owners, _, _ := aa.OnchainOwners(ctx, client, common.HexToAddress(primary.ID)); !containsAddress(owners, old) {
		t.Fatal("finalized before the recovery period")
	}
	var period int64 = 61
	client.Client().CallContext(ctx, nil, "evm_increaseTime", period)
	client.Client().CallContext(ctx, nil, "evm_mine")
	ProcessAccountRecoveries(ctx, gc) // finalizes
	owners, _, err := aa.OnchainOwners(ctx, client, common.HexToAddress(primary.ID))
	if err != nil || len(owners) != 1 || owners[0] != nw {
		t.Fatalf("primary owners %v (want %v, old %v): %v", owners, nw, old, err)
	}
	if !guarded(primary.ID) {
		t.Fatal("the primary wallet should stay covered")
	}
	ProcessAccountRecoveries(ctx, gc) // completes
	gdb.First(arl, "id = ?", arl.ID)
	if arl.Status != userModels.AccountRecoveryCompleted {
		t.Fatalf("recovery %v", arl.Status)
	}
	var after userModels.User
	gdb.First(&after, "id = ?", user.ID)
	if !strings.EqualFold(after.PrimarySigner, newKey.Address()) || after.Address != user.Address {
		t.Fatalf("user after recovery: signer %v address %v", after.PrimarySigner, after.Address)
	}

	// --- the new key moves funds from the primary wallet; the old one cannot
	var p userModels.UserWallet
	gdb.First(&p, "id = ?", primary.ID)
	payee, _ := evmkeypair.Random()
	pay := []aa.Call{{To: common.HexToAddress(payee.Address()), Value: big.NewInt(1000)}}
	op, err = PrepareWalletOperation(ctx, "PAYMENT", &after, &after, &p, pay, 0, nil, gc)
	if err != nil {
		t.Fatal(err)
	}
	sign("pay with the new key", "PAYMENT", op.Transaction, p.ID, newKey)
	if bal, _ := client.BalanceAt(ctx, common.HexToAddress(payee.Address()), nil); bal.Cmp(big.NewInt(1000)) != 0 {
		t.Fatalf("payee got %v", bal)
	}

	// --- the shared wallet: bob approves the REPLACE SIGNER request
	var pa userModels.PendingAuth
	if err := gdb.First(&pa, "wallet_address = ? AND transaction_type = ?", sub.ID, OperationReplaceSigner).Error; err != nil {
		t.Fatalf("no replace signer request: %v", err)
	}
	if strings.Contains(pa.Description, "Warning") {
		t.Fatalf("unexpected warning: %v", pa.Description)
	}
	opRec, opPrep := approvalOperation(&pa, gc)
	if opRec == nil {
		t.Fatal("request is not operation-backed")
	}
	if err := checkApprover(opPrep, bob.PrimarySigner); err != nil {
		t.Fatal(err)
	}
	sig, _ := evmkeypair.SignPersonal(bobKey.PrivateKey(), opPrep.SafeOpHash.Bytes())
	h, err := SubmitWalletOperation(ctx, opRec, opPrep, []aa.OwnerSignature{{Owner: common.HexToAddress(bob.PrimarySigner), Signature: sig}}, gc)
	if err != nil {
		t.Fatal(err)
	}
	wctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if r, err := operationBuilder(gc).Bundler.(*aa.Bundler).WaitReceipt(wctx, common.HexToHash(h)); err != nil || !r.Success {
		t.Fatalf("replace signer: %v", err)
	}
	owners, _, _ = aa.OnchainOwners(ctx, client, common.HexToAddress(sub.ID))
	if len(owners) != 2 || !containsAddress(owners, nw) || !containsAddress(owners, common.HexToAddress(bob.PrimarySigner)) || containsAddress(owners, old) {
		t.Fatalf("shared wallet owners %v", owners)
	}

	// --- the watcher: the backend's own recoveries raise no alert; one
	// started directly with the guardian's key (a misused key) does, and
	// the owner can still cancel it
	var alerts []string
	unexpectedRecoveryAlert = func(_ *sharedconfig.GlobalConfig, msg string) { alerts = append(alerts, msg) }
	defer func() { unexpectedRecoveryAlert = func(gc *sharedconfig.GlobalConfig, msg string) { gc.LogDiscordFailedRequest(msg) } }()
	ProcessAccountRecoveries(ctx, gc)
	if len(alerts) != 0 {
		t.Fatalf("alerts on the backend's recoveries: %v", alerts)
	}
	thief, _ := evmkeypair.Random()
	nonce, _ := aa.RecoveryNonce(ctx, client, module, common.HexToAddress(primary.ID))
	if _, err := guardian.exec(ctx, []aa.Call{aa.ConfirmRecoveryCall(module, common.HexToAddress(primary.ID), []common.Address{common.HexToAddress(thief.Address())}, 1, nonce)}, gc); err != nil {
		t.Fatal(err)
	}
	ProcessAccountRecoveries(ctx, gc)
	if len(alerts) != 1 || !strings.Contains(strings.ToLower(alerts[0]), strings.ToLower(primary.ID)[2:]) {
		t.Fatalf("alerts %v", alerts)
	}
	gdb.First(&p, "id = ?", primary.ID)
	op, err = PrepareWalletOperation(ctx, OperationRecoveryCancel, &after, &after, &p, []aa.Call{aa.CancelRecoveryCall(module)}, 0, nil, gc)
	if err != nil {
		t.Fatal(err)
	}
	sign("owner cancels the unexpected recovery", OperationRecoveryCancel, op.Transaction, p.ID, newKey)
	if req, _ := aa.ReadRecoveryRequest(ctx, client, module, common.HexToAddress(primary.ID)); req.ExecutableAt != 0 {
		t.Fatal("unexpected recovery still pending")
	}
}

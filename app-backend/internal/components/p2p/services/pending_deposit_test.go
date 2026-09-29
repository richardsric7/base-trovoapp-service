package p2p

import (
	"context"
	"strings"
	"testing"

	p2pModels "trovo-wallet-api/internal/components/p2p/models"
	userModels "trovo-wallet-api/internal/components/users/models"
	"trovo-wallet-api/internal/sharedconfig"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// A deposit made through the API is recorded under its userOpHash and
// applied to the order only once its wallet operation is mined, under the
// mined transaction's hash.
func TestPendingDepositResolution(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&p2pModels.Order{}, &p2pModels.BlockchainDeposit{}, &p2pModels.Refund{}, &p2pModels.P2PAuditEvent{}, &userModels.WalletOperation{}); err != nil {
		t.Fatal(err)
	}
	gc := &sharedconfig.GlobalConfig{DB: db}

	opHash, minedTx := "0x"+strings.Repeat("ab", 32), "0x"+strings.Repeat("cd", 32)
	op := userModels.WalletOperation{ID: "0x" + strings.Repeat("01", 32), WalletAddress: "0xw", Kind: "PAYMENT", Initiator: "alice", Status: userModels.WalletOperationSubmitted, Prepared: "{}", UserOpHash: &opHash}
	order := p2pModels.Order{ID: "order-1", OrderStatus: p2pModels.OrderStatusAwaitingEscrowDeposit, SellerEscrowAssetAmount: "10", ExpectedEscrowAmount: "10", DepositedEscrowAmount: "0"}
	pending := p2pModels.BlockchainDeposit{ID: "dep-1", OrderID: order.ID, Sender: "0xw", Token: "USDC", Amount: "10", TransactionHash: opHash, DetectionSource: detectionPendingOperation}
	for _, row := range []interface{}{&op, &order, &pending} {
		if err := db.Create(row).Error; err != nil {
			t.Fatal(err)
		}
	}

	// not mined yet (no bundler configured, still SUBMITTED): nothing happens
	if err := resolvePendingDeposit(context.Background(), gc, pending); err != nil {
		t.Fatal(err)
	}
	if o, _ := GetOrderByID(db, order.ID); o.OrderStatus != p2pModels.OrderStatusAwaitingEscrowDeposit {
		t.Fatalf("order moved on before its deposit was mined: %s", o.OrderStatus)
	}

	// mined: the sweep applies it under the transaction hash
	success := true
	db.Model(&userModels.WalletOperation{}).Where("id = ?", op.ID).Updates(map[string]interface{}{"status": userModels.WalletOperationIncluded, "tx_hash": minedTx, "success": &success})
	resolvePendingDeposits(gc)
	var d p2pModels.BlockchainDeposit
	db.First(&d, "id = ?", pending.ID)
	if d.TransactionHash != minedTx || d.DetectionSource != "P2P_API" || !d.IsCanonical {
		t.Fatalf("deposit not applied under the mined hash: %+v", d)
	}
	o, _ := GetOrderByID(db, order.ID)
	if o.OrderStatus != p2pModels.OrderStatusAwaitingPayment || o.EscrowDepositTransactionHash != minedTx || o.DepositedEscrowAmount != "10" {
		t.Fatalf("order not advanced: %+v", o)
	}

	// resolving again is a no-op (no double counting)
	resolvePendingDeposits(gc)
	if err := resolvePendingDeposit(context.Background(), gc, pending); err != nil {
		t.Fatal(err)
	}
	if o, _ := GetOrderByID(db, order.ID); o.DepositedEscrowAmount != "10" {
		t.Fatalf("deposit counted twice: %s", o.DepositedEscrowAmount)
	}

	// a failed operation's pending deposit is dropped
	failedHash := "0x" + strings.Repeat("ef", 32)
	failedOp := userModels.WalletOperation{ID: "0x" + strings.Repeat("02", 32), WalletAddress: "0xw", Kind: "PAYMENT", Initiator: "alice", Status: userModels.WalletOperationFailed, Prepared: "{}", UserOpHash: &failedHash}
	failedDep := p2pModels.BlockchainDeposit{ID: "dep-2", OrderID: order.ID, Sender: "0xw", Token: "USDC", Amount: "10", TransactionHash: failedHash, DetectionSource: detectionPendingOperation}
	db.Create(&failedOp)
	db.Create(&failedDep)
	if err := resolvePendingDeposit(context.Background(), gc, failedDep); err == nil {
		t.Fatal("a failed deposit must be reported")
	}
	var n int64
	db.Model(&p2pModels.BlockchainDeposit{}).Where("id = ?", failedDep.ID).Count(&n)
	if n != 0 {
		t.Fatal("the failed deposit must be dropped")
	}
}

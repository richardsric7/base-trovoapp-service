package users

import (
	"testing"
	"time"

	"trovo-wallet-api/internal/cache"
	userModels "trovo-wallet-api/internal/components/users/models"
	appdb "trovo-wallet-api/internal/db"
	"trovo-wallet-api/internal/sharedconfig"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestListProceedPayoutReceipts(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:"+t.TempDir()+"/payouts.db"), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	appdb.MigrateDB(db)
	gc := &sharedconfig.GlobalConfig{DB: db, RedisCache: &cache.RedisCache{}}

	// wallet ids are stored upper-case; schedule lines hold checksummed addresses
	const alice = "0X5AAEB6053F3E94C9B9A09F33669435E7EF1BEAED"
	const aliceChecksummed = "0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed"
	const bob = "0xfB6916095ca1df60bB79Ce92cE3Ea74c37c5d359"
	db.Exec("INSERT INTO user_wallets (id, alias, signer, user_id, primary_wallet) VALUES (?, 'alice', 'S', 7, 1)", alice)
	code, name, other, ng := "FARM", "Farm Fund", "BOND", "NG"
	for _, a := range []userModels.TokenizedAsset{{ID: "42", AssetCode: &code, AssetName: &name, AssetCountryLocation: &ng}, {ID: "43", AssetCode: &other, AssetCountryLocation: &ng}} {
		if err := db.Omit("TokenizedAssetType").Create(&a).Error; err != nil {
			t.Fatal(err)
		}
	}

	locked := time.Now().UTC()
	payouts := []userModels.ProceedPayout{
		{ID: 1, TokenizedAssetID: "42", Batch: "FARM-d1", Status: userModels.ProceedPayoutStatusCompleted, PayoutAssetCode: "CNGN", PayoutDecimals: 6, TokenDecimals: 0, AmountPerToken: "0.95", LockedAt: &locked},
		{ID: 2, TokenizedAssetID: "42", Batch: "FARM-d2", Status: userModels.ProceedPayoutStatusCancelled, PayoutAssetCode: "CNGN", PayoutDecimals: 6},
		{ID: 3, TokenizedAssetID: "43", Batch: "BOND-d3", Status: userModels.ProceedPayoutStatusLocked, PayoutAssetCode: "CNGN", PayoutDecimals: 6, LockedAt: &locked},
		{ID: 4, TokenizedAssetID: "43", Batch: "BOND-d4", Status: userModels.ProceedPayoutStatusPreparing, PayoutAssetCode: "CNGN", PayoutDecimals: 6},
	}
	for i := range payouts {
		if err := db.Create(&payouts[i]).Error; err != nil {
			t.Fatal(err)
		}
	}
	item := func(id string, payout uint64, addr, status, kind string) userModels.TokenizedAssetPayoutSchedule {
		return userModels.TokenizedAssetPayoutSchedule{ID: id, ProceedPayoutID: payout, BeneficiaryAddress: addr, BalanceUnits: "600",
			AmountUnits: "570000000", Status: status, Kind: kind, Batch: "b", TokenizedAssetID: "42", PayoutAssetCode: "CNGN", PayoutContractAddress: "0xc"}
	}
	items := []userModels.TokenizedAssetPayoutSchedule{
		item("paid", 1, aliceChecksummed, userModels.PayoutItemPaid, userModels.PayoutItemKindHolder),
		item("bob", 1, bob, userModels.PayoutItemPaid, userModels.PayoutItemKindHolder),
		item("cancelled", 2, aliceChecksummed, userModels.PayoutItemPending, userModels.PayoutItemKindHolder),
		item("scheduled", 3, aliceChecksummed, userModels.PayoutItemPending, userModels.PayoutItemKindHolder),
		item("preparing", 4, aliceChecksummed, userModels.PayoutItemPending, userModels.PayoutItemKindHolder),
		item("excluded", 3, "0x0000000000000000000000000000000000000001", userModels.PayoutItemExcluded, userModels.PayoutItemKindHolder),
	}
	for i := range items {
		if err := db.Create(&items[i]).Error; err != nil {
			t.Fatal(err)
		}
	}

	user := userModels.User{ID: "7"}
	got, err := ListProceedPayoutReceipts(&user, "", gc)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0].ID != "scheduled" || got[1].ID != "paid" {
		t.Fatalf("receipts: %+v", got)
	}
	r := got[1]
	if r.Amount != "570" || r.TokensHeld != "600" || r.AssetCode != "FARM" || r.AssetName != "Farm Fund" || r.WalletAlias != "alice" || r.Status != userModels.PayoutItemPaid {
		t.Fatalf("paid receipt: %+v", r)
	}
	if only, _ := ListProceedPayoutReceipts(&user, "42", gc); len(only) != 1 || only[0].ID != "paid" {
		t.Fatalf("filtered by asset: %+v", only)
	}
	if none, _ := ListProceedPayoutReceipts(&userModels.User{ID: "8"}, "", gc); len(none) != 0 {
		t.Fatalf("a user without wallets: %+v", none)
	}
}

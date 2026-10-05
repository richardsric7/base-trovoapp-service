package users

import (
	"testing"

	"trovo-wallet-api/internal/sharedconfig"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestUSDAmountIn(t *testing.T) {
	t.Setenv("DOLLAR_ASSET", "USDB:0x0000000000000000000000000000000000000001")
	t.Setenv("NAIRA_ASSET", "")
	for _, c := range []struct {
		code, contract string
		want           string
	}{
		{"USDB", "0x0000000000000000000000000000000000000001", "2.5"},
		{"usdc", "0x0000000000000000000000000000000000000002", "2.5"},
		{"USDT", "0x0000000000000000000000000000000000000003", "2.5"},
	} {
		got, err := usdAmountIn(2.5, c.code, c.contract, nil)
		if err != nil || got.String() != c.want {
			t.Fatalf("%s: got %v, %v; want %s", c.code, got, err, c.want)
		}
	}
	// a $2.50 fee must never become 1/2.5 of a token (the pre-port inversion)
	if got, _ := usdAmountIn(2.5, "USDC", "0x02", nil); got.String() == "0.4" {
		t.Fatal("inverted conversion")
	}
	if _, err := usdAmountIn(2.5, "TROV", "0x0000000000000000000000000000000000000004", nil); err == nil {
		t.Fatal("a non-stablecoin cannot be priced without a DEX")
	}
}

func TestFeeExemptProfile(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&sharedconfig.FeeExemptUser{}); err != nil {
		t.Fatal(err)
	}
	gc := &sharedconfig.GlobalConfig{DB: db}
	t.Setenv("TOKENIZATION_ISSUING_PROFILE", "")
	if !feeExemptProfile("atprofile", gc) || feeExemptProfile("alice", gc) {
		t.Fatal("default exempt profile is atprofile")
	}
	t.Setenv("TOKENIZATION_ISSUING_PROFILE", "issuer")
	if !feeExemptProfile("ISSUER", gc) || feeExemptProfile("atprofile", gc) {
		t.Fatal("the configured issuing profile is exempt")
	}
	// the admin-managed list (tm-api), case-insensitively
	db.Create(&sharedconfig.FeeExemptUser{Username: "Treasury", Reason: "platform market maker", AddedBy: "admin@trovo"})
	if !feeExemptProfile("treasury", gc) || feeExemptProfile("alice", gc) {
		t.Fatal("listed accounts are exempt, others are not")
	}
}

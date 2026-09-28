package network

import (
	"testing"

	"trovo-wallet-api/internal/basetxn"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// countryConfigRow mirrors just the columns isInternalBalanceAsset's raw
// query reads from internal/components/users/models.CountryConfig - kept
// separate (rather than importing that package) for the same import-cycle
// reason documented on isInternalBalanceAsset itself.
type countryConfigRow struct {
	CountryCode               string `gorm:"primaryKey;size:2"`
	InternalBalanceTokenCode  string
	InternalTokenIssuer       string
}

func (countryConfigRow) TableName() string { return "country_configs" }

// tokenizedAssetRow mirrors just the columns isTokenizedAsset's raw query
// reads from Tokenized_Assets.
type tokenizedAssetRow struct {
	ID                     string `gorm:"primaryKey"`
	AssetTokenizationStatus int
	AssetCode              string
}

func (tokenizedAssetRow) TableName() string { return "Tokenized_Assets" }

// setupAuthTestDB wires an in-memory SQLite DB directly into authDB,
// bypassing SetDB (which also dials a blockchain RPC client via
// basetxn.SetDefaultBuilder - unwanted/unavailable in a unit test).
func setupAuthTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	// A named, shared-cache in-memory DB persists across connections opened
	// from the same DSN within the process - naming it per-test keeps each
	// test's DB isolated from the others instead of colliding on one shared
	// "file::memory:" instance.
	dsn := "file:" + t.Name() + "?mode=memory&cache=shared"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory sqlite db: %v", err)
	}
	if err := db.AutoMigrate(&WalletAssetAuthorization{}, &countryConfigRow{}, &tokenizedAssetRow{}); err != nil {
		t.Fatalf("failed to automigrate: %v", err)
	}
	authDB = db
	t.Cleanup(func() { authDB = nil })
	return db
}

func TestRequiresWalletAuthorization(t *testing.T) {
	db := setupAuthTestDB(t)

	if err := db.Create(&tokenizedAssetRow{ID: "ta-1", AssetCode: "REIT1", AssetTokenizationStatus: 5}).Error; err != nil {
		t.Fatalf("seeding tokenized asset: %v", err)
	}
	if err := db.Create(&countryConfigRow{CountryCode: "NG", InternalBalanceTokenCode: "NGNX", InternalTokenIssuer: "0xIssuer"}).Error; err != nil {
		t.Fatalf("seeding country config: %v", err)
	}

	if !RequiresWalletAuthorization("REIT1") {
		t.Error("expected a market-ready tokenized asset (status 5) to require wallet authorization")
	}
	if !RequiresWalletAuthorization("ngnx") {
		t.Error("expected the internal balance token to require wallet authorization (case-insensitively)")
	}
	if RequiresWalletAuthorization("RANDOMB20") {
		t.Error("expected a plain, non-regulated B20 asset to not require wallet authorization")
	}
}

func TestSetAndIsWalletAuthorizedForAsset_InternalBalanceAsset(t *testing.T) {
	db := setupAuthTestDB(t)
	if err := db.Create(&countryConfigRow{CountryCode: "NG", InternalBalanceTokenCode: "NGNX", InternalTokenIssuer: "0xissuer"}).Error; err != nil {
		t.Fatalf("seeding country config: %v", err)
	}

	wallet := "0xWALLET"
	asset := basetxn.CreditAsset{Code: "NGNX", Issuer: "0xissuer"}

	if IsWalletAuthorizedForAsset(wallet, asset) {
		t.Fatal("expected wallet to NOT be authorized for the internal balance asset before any grant")
	}

	if err := SetWalletAssetAuthorization(wallet, asset, true, "0xIssuerAuthority", "smoke test grant"); err != nil {
		t.Fatalf("SetWalletAssetAuthorization(grant) failed: %v", err)
	}
	if !IsWalletAuthorizedForAsset(wallet, asset) {
		t.Fatal("expected wallet to be authorized for the internal balance asset after granting")
	}

	// Revoking (authorized=false) must actually persist - GORM's
	// struct-to-Assign conversion silently drops false/"" fields, which is
	// exactly the bug SetWalletAssetAuthorization's own map-based Assign
	// call guards against (see its doc comment).
	if err := SetWalletAssetAuthorization(wallet, asset, false, "0xIssuerAuthority", "smoke test revoke"); err != nil {
		t.Fatalf("SetWalletAssetAuthorization(revoke) failed: %v", err)
	}
	if IsWalletAuthorizedForAsset(wallet, asset) {
		t.Fatal("expected wallet to NOT be authorized after revoking")
	}
}

func TestIsWalletAuthorizedForAsset_UnregulatedAssetAlwaysAuthorized(t *testing.T) {
	setupAuthTestDB(t)

	wallet := "0xWALLET"
	asset := basetxn.CreditAsset{Code: "RANDOMB20", Issuer: "0xissuer"}

	if !IsWalletAuthorizedForAsset(wallet, asset) {
		t.Error("expected a non-regulated B20 asset to always be authorized, with no authorization row")
	}
}

func TestIsWalletAuthorizedForAsset_NativeAssetAlwaysAuthorized(t *testing.T) {
	setupAuthTestDB(t)

	if !IsWalletAuthorizedForAsset("0xWALLET", basetxn.NativeAsset{}) {
		t.Error("expected the native asset to always be authorized")
	}
}

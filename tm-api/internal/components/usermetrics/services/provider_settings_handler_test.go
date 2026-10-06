package usermetrics

import (
	"strings"
	"testing"

	"admin-panel-dashboard/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func providerSettingsDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&models.KycConfig{}, &models.StablerailConfig{}); err != nil {
		t.Fatal(err)
	}
	return db
}

func TestKycProviderSettingsKeepSecretsAndNeverShowThem(t *testing.T) {
	db := providerSettingsDB(t)
	if err := SaveKycProviderSettings(db, "sumsub", models.KycProviderSettingsRequest{Token: "sbx:token-value-123"}); err == nil {
		t.Fatal("the first save must need both values")
	}
	if err := SaveKycProviderSettings(db, "sumsub", models.KycProviderSettingsRequest{Token: "sbx:token-value-123", SecretKey: "secret-value-456"}); err != nil {
		t.Fatal(err)
	}
	// rotating only the token keeps the stored secret
	if err := SaveKycProviderSettings(db, "sumsub", models.KycProviderSettingsRequest{Token: "sbx:token-value-789"}); err != nil {
		t.Fatal(err)
	}
	var row models.KycConfig
	db.Where("service_provider = ?", "sumsub").First(&row)
	if row.Token != "sbx:token-value-789" || row.SecretKey != "secret-value-456" {
		t.Fatalf("stored %+v", row)
	}
	out, err := LoadProviderSettings(db)
	if err != nil {
		t.Fatal(err)
	}
	if len(out.Kyc) != len(KycProviders) {
		t.Fatalf("want every provider listed, got %+v", out.Kyc)
	}
	for _, k := range out.Kyc {
		if strings.Contains(k.Token.Hint, "789") && len(k.Token.Hint) > 8 || strings.Contains(k.SecretKey.Hint, "value") {
			t.Fatalf("secret leaked: %+v", k)
		}
		if k.ServiceProvider == "sumsub" && (!k.Token.Set || !k.SecretKey.Set) {
			t.Fatalf("sumsub should be set: %+v", k)
		}
		if k.ServiceProvider == "doja" && (k.Token.Set || k.SecretKey.Set) {
			t.Fatalf("doja should be unset: %+v", k)
		}
	}
}

func TestStablerailSettings(t *testing.T) {
	db := providerSettingsDB(t)
	if err := SaveStablerailSettings(db, models.StablerailSettingsRequest{Enabled: true, BaseURL: "https://sr.example.com"}); err == nil {
		t.Fatal("turning it on without an API key must fail")
	}
	if err := SaveStablerailSettings(db, models.StablerailSettingsRequest{Enabled: true, BaseURL: "https://sr.example.com/", FintechID: "ft_1", ApiKey: "sr-api-key-000111"}); err != nil {
		t.Fatal(err)
	}
	// changing the address keeps the key; turning it off is allowed
	if err := SaveStablerailSettings(db, models.StablerailSettingsRequest{Enabled: false, BaseURL: "https://sr2.example.com", FintechID: "ft_1"}); err != nil {
		t.Fatal(err)
	}
	var rows []models.StablerailConfig
	db.Find(&rows)
	if len(rows) != 1 || rows[0].ApiKey != "sr-api-key-000111" || rows[0].BaseUrl != "https://sr2.example.com" || rows[0].EnableStablerail != 0 {
		t.Fatalf("stored %+v", rows)
	}
	out, _ := LoadProviderSettings(db)
	if out.Stablerail.Enabled || !out.Stablerail.ApiKey.Set || strings.Contains(out.Stablerail.ApiKey.Hint, "000111") {
		t.Fatalf("got %+v", out.Stablerail)
	}
}

package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	coreModels "admin-panel-dashboard/internal/models"
	serverModels "admin-panel-dashboard/internal/server/models"
	"admin-panel-dashboard/internal/trovosdk"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestRequestWalletLinkAuthorizationAllowsUnlinkedMember(t *testing.T) {
	server, member, walletAPI := walletLinkTestServer(t)
	defer walletAPI.Close()
	t.Setenv("LOGIN_CALLBACK_URL", "https://manager.example/callback")

	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Set("member_id", member.ID)
	response, err := SendWalletLinkAuthorizationRequest(server, ctx, "alice@trovo")
	if err != nil {
		t.Fatalf("request authorization: %v", err)
	}
	if response.AuthID != "wallet-auth-1" {
		t.Fatalf("auth id = %q, want wallet-auth-1", response.AuthID)
	}
	if err := server.AdminDB.First(&member, "id = ?", member.ID).Error; err != nil {
		t.Fatal(err)
	}
	if member.IsWalletLinked || member.TrovoWalletUsername != nil {
		t.Fatalf("wallet was persisted before verification: %+v", member)
	}
}

func TestRequestWalletLinkAuthorizationRejectsUnknownWallet(t *testing.T) {
	server, member, walletAPI := walletLinkTestServer(t)
	defer walletAPI.Close()
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Set("member_id", member.ID)

	if _, err := SendWalletLinkAuthorizationRequest(server, ctx, "missing@trovo"); err == nil {
		t.Fatal("expected unknown wallet to be rejected")
	}
}

func TestVerifyWalletLinkPersistsApprovedWallet(t *testing.T) {
	server, member, walletAPI := walletLinkTestServer(t)
	defer walletAPI.Close()

	if err := VerifyWalletLink(server, member.ID, "alice@trovo", "wallet-auth-1"); err != nil {
		t.Fatalf("verify wallet link: %v", err)
	}
	if err := server.AdminDB.First(&member, "id = ?", member.ID).Error; err != nil {
		t.Fatal(err)
	}
	if !member.IsWalletLinked || member.TrovoWalletUsername == nil || *member.TrovoWalletUsername != "alice@trovo" {
		t.Fatalf("wallet link not persisted: %+v", member)
	}
}

func TestVerifyWalletLinkCannotReplaceExistingWallet(t *testing.T) {
	server, member, walletAPI := walletLinkTestServer(t)
	defer walletAPI.Close()
	linked := "existing@trovo"
	if err := server.AdminDB.Model(&member).Updates(map[string]interface{}{
		"is_wallet_linked": true, "trovo_wallet_username": linked,
	}).Error; err != nil {
		t.Fatal(err)
	}

	if err := VerifyWalletLink(server, member.ID, "alice@trovo", "wallet-auth-1"); err == nil {
		t.Fatal("expected replacement attempt to be rejected")
	}
	if err := server.AdminDB.First(&member, "id = ?", member.ID).Error; err != nil {
		t.Fatal(err)
	}
	if member.TrovoWalletUsername == nil || *member.TrovoWalletUsername != linked {
		t.Fatalf("existing wallet changed: %+v", member)
	}
}

func TestVerifyWalletLinkDoesNotPersistFailedAuthorization(t *testing.T) {
	server, member, walletAPI := walletLinkTestServer(t)
	defer walletAPI.Close()

	if err := VerifyWalletLink(server, member.ID, "alice@trovo", "invalid-auth"); err == nil {
		t.Fatal("expected invalid authorization to fail")
	}
	if err := server.AdminDB.First(&member, "id = ?", member.ID).Error; err != nil {
		t.Fatal(err)
	}
	if member.IsWalletLinked || member.TrovoWalletUsername != nil {
		t.Fatalf("failed authorization persisted wallet link: %+v", member)
	}
}

func walletLinkTestServer(t *testing.T) (*serverModels.Server, coreModels.OrganizationMember, *httptest.Server) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	adminDB := walletLinkSQLiteDB(t, "admin")
	walletDB := walletLinkSQLiteDB(t, "wallet")
	if err := adminDB.AutoMigrate(&coreModels.Organization{}, &coreModels.OrganizationMember{}); err != nil {
		t.Fatal(err)
	}
	if err := walletDB.AutoMigrate(&coreModels.User{}); err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	organization := coreModels.Organization{
		ID: uuid.NewString(), Name: "Trustee", Email: uuid.NewString() + "@example.com",
		Type: coreModels.OrganizationTypeAssetCustodian, Status: coreModels.OrganizationStatusActive,
		CreatedBy: "test", CreatedAt: now, UpdatedAt: now,
	}
	member := coreModels.OrganizationMember{
		ID: uuid.NewString(), OrganizationID: organization.ID, Email: uuid.NewString() + "@example.com",
		Role: string(coreModels.OrganizationSuperAdmin), Status: coreModels.OrganizationStatusActive,
		CreatedBy: "test", CreatedAt: now, UpdatedAt: now,
	}
	if err := adminDB.Create(&organization).Error; err != nil {
		t.Fatal(err)
	}
	if err := adminDB.Create(&member).Error; err != nil {
		t.Fatal(err)
	}
	if err := walletDB.Create(&coreModels.User{ID: uuid.NewString(), Username: "alice", Email: "alice@example.com"}).Error; err != nil {
		t.Fatal(err)
	}

	walletAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/v1/servicelinks/authorize/request/alice":
			_ = json.NewEncoder(w).Encode(trovosdk.TrovoWalletAuthorizationData{AuthID: "wallet-auth-1", DynamicLink: "https://wallet.example/approve", QRCode: "qr"})
		case "/v1/servicelinks/authorize/verify/manager/alice/wallet-auth-1":
			_ = json.NewEncoder(w).Encode(trovosdk.TrovoWalletAuthorizationRespomseData{Message: "verified"})
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(func() { _ = os.Unsetenv("LOGIN_CALLBACK_URL") })
	return &serverModels.Server{
		AdminDB: adminDB, TrovoWalletDB: walletDB,
		GC: &coreModels.GlobalConfig{ServiceLink: &trovosdk.ServiceLink{ApiBaseUrl: walletAPI.URL, ServiceUsername: "manager", ApiKey: "test-key"}},
	}, member, walletAPI
}

func walletLinkSQLiteDB(t *testing.T, name string) *gorm.DB {
	t.Helper()
	database, err := gorm.Open(sqlite.Open("file:"+name+"-"+uuid.NewString()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	return database
}

package services

import (
	"context"
	"testing"
	"time"

	stakeholderDB "admin-panel-dashboard/internal/components/stakeholder/db"
	"admin-panel-dashboard/internal/components/stakeholder/models"
	coreModels "admin-panel-dashboard/internal/models"

	"github.com/shopspring/decimal"
)

func TestCustodianDashboardProvidesDesignDataWithoutMixingCurrencies(t *testing.T) {
	database := testDB(t,
		&models.StakeholderAssetAssignment{},
		&models.SegregatedAccount{},
		&models.FundReleaseRequest{},
		&models.ComplianceItem{},
		&models.StakeholderNotification{},
		&models.StakeholderAuditLog{},
		&coreModels.Organization{},
		&coreModels.OrganizationMember{},
	)
	custodianStakeholderID := uint64(17)
	custodianType := models.StakeholderTypeAssetCustodian
	organizations := []coreModels.Organization{
		{ID: "dashboard-custodian", Name: "Custodian Bank", Email: "dashboard-custodian@example.com", Type: "CORPORATE", Status: coreModels.OrganizationStatusActive, CreatedBy: "test", StakeholderID: &custodianStakeholderID, StakeholderType: &custodianType},
		{ID: "dashboard-manager", Name: "Asset Manager", Email: "dashboard-manager@example.com", Type: "CORPORATE", Status: coreModels.OrganizationStatusActive, CreatedBy: "test"},
	}
	if err := database.Create(&organizations).Error; err != nil {
		t.Fatalf("seed organizations: %v", err)
	}
	if err := database.Create(&coreModels.OrganizationMember{
		ID: "dashboard-member", OrganizationID: "dashboard-custodian", Email: "dashboard-member@example.com",
		Role: "ADMIN", Status: "ACTIVE", CreatedBy: "test", FirstName: "Nancy", LastName: "Custodian",
	}).Error; err != nil {
		t.Fatalf("seed member: %v", err)
	}

	assets := []stakeholderDB.TokenizedAsset{
		{ID: "dashboard-asset-1", AssetCode: "AGR1", AssetSector: "Agriculture", AssetQuoteCurrency: "CNGN", AssetCurrentValue: 100, CustodianFeeValue: 10, ApprovedAssetCustodianID: custodianStakeholderID, CreatedAt: time.Now().Add(-time.Hour)},
		{ID: "dashboard-asset-2", AssetCode: "REAL1", AssetSector: "Real Estate", AssetQuoteCurrency: "USD", AssetCurrentValue: 200, CustodianFeeValue: 5, ApprovedAssetCustodianID: custodianStakeholderID, CreatedAt: time.Now().Add(-2 * time.Hour)},
		{ID: "dashboard-asset-3", AssetCode: "AGR2", AssetSector: "Agriculture", AssetQuoteCurrency: "CNGN", AssetCurrentValue: 50, CustodianFeeValue: 2, ApprovedAssetCustodianID: custodianStakeholderID, CreatedAt: time.Now().Add(-3 * time.Hour)},
	}
	if err := database.Create(&[]models.SegregatedAccount{
		{ID: "dashboard-account-1", AssetID: assets[0].ID, CustodianOrgID: "dashboard-custodian", AccountType: "custody", AccountName: "CNGN account", Balance: decimal.NewFromInt(1000), Currency: "CNGN"},
		{ID: "dashboard-account-2", AssetID: assets[1].ID, CustodianOrgID: "dashboard-custodian", AccountType: "custody", AccountName: "USD account", Balance: decimal.NewFromInt(50), Currency: "USD"},
	}).Error; err != nil {
		t.Fatalf("seed accounts: %v", err)
	}
	if err := database.Create(&models.FundReleaseRequest{
		ID: "dashboard-release", AssetID: assets[0].ID, RequesterOrgID: "dashboard-manager", RequesterMemberID: "manager-member",
		TrusteeOrgID: "dashboard-trustee", CustodianOrgID: "dashboard-custodian", Amount: decimal.NewFromInt(20), Currency: "CNGN", Purpose: "works", Status: models.FundReleaseStatusExecutionPending,
	}).Error; err != nil {
		t.Fatalf("seed fund release: %v", err)
	}
	if err := database.Create(&models.ComplianceItem{ID: "dashboard-compliance", OrgID: "dashboard-custodian", Category: "KYC", Requirement: "Review", Status: models.ComplianceStatusPending}).Error; err != nil {
		t.Fatalf("seed compliance: %v", err)
	}
	now := time.Now()
	managerOrgID := "dashboard-manager"
	if err := database.Create(&models.StakeholderNotification{
		ID: "dashboard-notification", RecipientOrgID: "dashboard-custodian", SenderOrgID: &managerOrgID,
		Type: "fund_release.ready", Title: "Fund release ready", Message: "Execute release", CreatedAt: now,
	}).Error; err != nil {
		t.Fatalf("seed notification: %v", err)
	}
	if err := database.Create(&models.StakeholderAuditLog{
		ID: "dashboard-audit", ActorMemberID: "dashboard-member", ActorOrgID: "dashboard-custodian",
		Action: models.ActionComplianceUpdate, EntityType: models.EntityComplianceItem, EntityID: "dashboard-compliance",
		AfterState: models.JSONMap{"status": models.ComplianceStatusPending}, CreatedAt: now.Add(-time.Minute),
	}).Error; err != nil {
		t.Fatalf("seed audit: %v", err)
	}

	assetService := NewAssetService(database, fakeTokenizationReadClient{assets: assets}, nil)
	dashboardService := NewDashboardService(database, assetService)
	response, err := dashboardService.Dashboard(context.Background(), AuthContext{
		OrganizationID: "dashboard-custodian", StakeholderID: custodianStakeholderID, DashboardRole: models.DashboardRoleAssetCustodian,
	})
	if err != nil {
		t.Fatalf("custodian dashboard: %v", err)
	}
	if response.PortfolioSummary == nil || response.PortfolioSummary.AssetCount != 3 {
		t.Fatalf("portfolio = %+v, want three assets", response.PortfolioSummary)
	}
	assertCurrencyTotals(t, response.PortfolioSummary.Values, map[string]string{"CNGN": "150", "USD": "200"})
	assertCurrencyTotals(t, response.Summary["fees_generated"].([]models.CurrencyTotal), map[string]string{"CNGN": "12", "USD": "5"})
	assertCurrencyTotals(t, response.Summary["account_balances"].([]models.CurrencyTotal), map[string]string{"CNGN": "1000", "USD": "50"})
	if len(response.RecentAssets) != 3 {
		t.Fatalf("recent assets = %d, want 3", len(response.RecentAssets))
	}
	if len(response.RecentActivity) != 2 || response.RecentActivity[0].ID != "dashboard-notification" {
		t.Fatalf("recent activity = %+v, want notification then audit", response.RecentActivity)
	}
	if response.RecentActivity[0].ActorOrganization != "Asset Manager" || response.RecentActivity[1].ActorName != "Nancy Custodian" {
		t.Fatalf("activity actors = %+v, want resolved organization and member names", response.RecentActivity)
	}
}

func assertCurrencyTotals(t *testing.T, totals []models.CurrencyTotal, expected map[string]string) {
	t.Helper()
	if len(totals) != len(expected) {
		t.Fatalf("currency totals = %+v, want %v", totals, expected)
	}
	for _, total := range totals {
		if expected[total.Currency] != total.Amount {
			t.Fatalf("currency total %s = %s, want %s", total.Currency, total.Amount, expected[total.Currency])
		}
	}
}

func TestTrusteeAndManagerDashboardsProvideDesignSections(t *testing.T) {
	database := testDB(t,
		&models.StakeholderAssetAssignment{},
		&models.DueDiligenceChecklist{},
		&models.FundReleaseRequest{},
		&models.Distribution{},
		&models.RevenueRecord{},
		&models.StakeholderNotification{},
		&models.StakeholderAuditLog{},
		&coreModels.Organization{},
		&coreModels.OrganizationMember{},
	)
	trusteeID, managerID := uint64(21), uint64(22)
	assets := []stakeholderDB.TokenizedAsset{
		{ID: "role-dashboard-asset-1", AssetCode: "ROLE1", AssetSector: "Agriculture", AssetQuoteCurrency: "CNGN", AssetCurrentValue: 100, TrusteeID: trusteeID, AssetManagerID: managerID, CreatedAt: time.Now()},
		{ID: "role-dashboard-asset-2", AssetCode: "ROLE2", AssetSector: "Real Estate", AssetQuoteCurrency: "USD", AssetCurrentValue: 250, TrusteeID: trusteeID, AssetManagerID: managerID, CreatedAt: time.Now().Add(-time.Hour)},
	}
	assetService := NewAssetService(database, fakeTokenizationReadClient{assets: assets}, nil)
	dashboardService := NewDashboardService(database, assetService)

	for _, test := range []struct {
		name string
		auth AuthContext
	}{
		{name: "trustee", auth: AuthContext{OrganizationID: "role-dashboard-trustee", StakeholderID: trusteeID, DashboardRole: models.DashboardRoleTrustee}},
		{name: "asset manager", auth: AuthContext{OrganizationID: "role-dashboard-manager", StakeholderID: managerID, DashboardRole: models.DashboardRoleAssetManager}},
	} {
		t.Run(test.name, func(t *testing.T) {
			response, err := dashboardService.Dashboard(context.Background(), test.auth)
			if err != nil {
				t.Fatalf("dashboard: %v", err)
			}
			if len(response.RecentAssets) != 2 {
				t.Fatalf("recent assets = %d, want 2", len(response.RecentAssets))
			}
			if response.PortfolioSummary == nil || response.PortfolioSummary.AssetCount != 2 {
				t.Fatalf("portfolio = %+v, want two assets", response.PortfolioSummary)
			}
			assertCurrencyTotals(t, response.PortfolioSummary.Values, map[string]string{"CNGN": "100", "USD": "250"})
			if response.RecentActivity == nil {
				t.Fatal("recent activity must be an empty array, not nil")
			}
		})
	}
}

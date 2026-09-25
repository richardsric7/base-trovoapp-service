package services

import (
	"context"
	"testing"

	stakeholderDB "admin-panel-dashboard/internal/components/stakeholder/db"
	"admin-panel-dashboard/internal/components/stakeholder/models"
	coreModels "admin-panel-dashboard/internal/models"

	"github.com/shopspring/decimal"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type fundSummarySubscriptionFixture struct {
	ID               string `gorm:"primaryKey"`
	TokenizedAssetID string
	WalletPublicKey  string
	Amount           float64
}

func (fundSummarySubscriptionFixture) TableName() string { return "tokenized_asset_subscriptions" }

func TestTrusteeFundManagementSummaryUsesScopedCanonicalSources(t *testing.T) {
	adminDB := testDB(t,
		&models.StakeholderAssetAssignment{}, &models.StakeholderDocument{}, &coreModels.Organization{},
		&models.FundReleaseRequest{}, &models.Distribution{}, &models.RevenueRecord{}, &models.SegregatedAccount{},
	)
	walletDB, err := gorm.Open(sqlite.Open("file:"+t.Name()+"-wallet?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open wallet database: %v", err)
	}
	if err := walletDB.AutoMigrate(&stakeholderDB.TokenizedAsset{}, &fundSummarySubscriptionFixture{}); err != nil {
		t.Fatalf("migrate wallet fixtures: %v", err)
	}

	trusteeOrgID := "trustee-org"
	trusteeStakeholderID := uint64(30)
	otherTrusteeID := uint64(31)
	if err := adminDB.Create(&models.StakeholderAssetAssignment{
		ID: "assignment-1", AssetID: "asset-1", AssetCode: "AST1", TrusteeOrgID: &trusteeOrgID,
		TrusteeStakeholderID: &trusteeStakeholderID, Status: models.AssignmentStatusActive,
	}).Error; err != nil {
		t.Fatalf("seed assignment: %v", err)
	}
	assets := []stakeholderDB.TokenizedAsset{
		{
			ID: "asset-1", AssetCode: "AST1", AssetQuoteCurrency: "cngn", TrusteeID: trusteeStakeholderID,
			FeeInFiat: 10, SECTokenizationFeeValue: 2, CustodianFeeValue: 3, AssetManagerFeeValue: 4,
			IssuingHouseFeeValue: 5, LegalAndProfessionalFeeValue: 6, RatingAgencyFeeValue: 7,
			TrusteeFeeValue: 8, VATValue: 9,
		},
		{ID: "asset-hidden", AssetCode: "HIDDEN", AssetQuoteCurrency: "CNGN", TrusteeID: otherTrusteeID, FeeInFiat: 999},
	}
	if err := walletDB.Create(&assets).Error; err != nil {
		t.Fatalf("seed wallet assets: %v", err)
	}
	if err := walletDB.Create(&[]fundSummarySubscriptionFixture{
		{ID: "subscription-1", TokenizedAssetID: "asset-1", WalletPublicKey: "wallet-1", Amount: 1000},
		{ID: "subscription-hidden", TokenizedAssetID: "asset-hidden", WalletPublicKey: "wallet-hidden", Amount: 9000},
	}).Error; err != nil {
		t.Fatalf("seed primary sales: %v", err)
	}

	seedFundManagementRows(t, adminDB, "asset-1", "CNGN", trusteeOrgID)
	seedFundManagementRows(t, adminDB, "asset-hidden", "CNGN", "other-trustee")
	if err := adminDB.Create(&models.FundReleaseRequest{
		ID: "asset-1-pending-release", AssetID: "asset-1", RequesterOrgID: "manager", RequesterMemberID: "member",
		TrusteeOrgID: trusteeOrgID, Amount: decimal.NewFromInt(999), Currency: "CNGN", Purpose: "pending", Status: models.FundReleaseStatusSubmitted,
	}).Error; err != nil {
		t.Fatalf("seed pending release: %v", err)
	}
	if err := adminDB.Create(&models.SegregatedAccount{
		ID: "asset-1-inactive-account", AssetID: "asset-1", CustodianOrgID: "custodian",
		AccountType: models.AccountTypeDevelopmentFund, AccountName: "Inactive", Balance: decimal.NewFromInt(999), Currency: "CNGN", Status: models.AccountStatusInactive,
	}).Error; err != nil {
		t.Fatalf("seed inactive milestone account: %v", err)
	}

	assetService := NewAssetService(adminDB, stakeholderDB.NewGormTokenizationReadClient(walletDB), nil)
	service := NewFundManagementService(adminDB, walletDB, assetService)
	auth := AuthContext{OrganizationID: trusteeOrgID, StakeholderID: trusteeStakeholderID, DashboardRole: models.DashboardRoleTrustee}
	summary, err := service.TrusteeSummary(context.Background(), auth)
	if err != nil {
		t.Fatalf("get trustee fund summary: %v", err)
	}
	if summary.AssetCount != 1 {
		t.Fatalf("asset count = %d, want 1", summary.AssetCount)
	}
	assertCurrencyAmount(t, summary.TotalAmountFromPrimarySales, "CNGN", "1000")
	assertCurrencyAmount(t, summary.TotalIncomeFromAssets, "CNGN", "300")
	assertCurrencyAmount(t, summary.TotalAmountPaidToIssuers, "CNGN", "200")
	assertCurrencyAmount(t, summary.TotalPaidToInvestors, "CNGN", "50")
	assertCurrencyAmount(t, summary.TotalAmountProcessed, "CNGN", "250")
	assertCurrencyAmount(t, summary.MilestonePaymentBalance, "CNGN", "400")
	assertCurrencyAmount(t, summary.TotalFeesGenerated, "CNGN", "54")
}

func TestTrusteeFundManagementSummaryRejectsOtherRoles(t *testing.T) {
	service := NewFundManagementService(nil, nil, nil)
	_, err := service.TrusteeSummary(context.Background(), AuthContext{DashboardRole: models.DashboardRoleAssetManager})
	if ErrorStatus(err) != 403 {
		t.Fatalf("status = %d, want 403", ErrorStatus(err))
	}
}

func seedFundManagementRows(t *testing.T, database *gorm.DB, assetID, currency, trusteeOrgID string) {
	t.Helper()
	rows := []interface{}{
		&models.RevenueRecord{ID: assetID + "-revenue", AssetID: assetID, ManagerOrgID: "manager", Amount: decimal.NewFromInt(300), Currency: currency, Source: "operations", Status: models.RevenueStatusRecorded, CreatedByMemberID: "member"},
		&models.FundReleaseRequest{ID: assetID + "-release", AssetID: assetID, RequesterOrgID: "manager", RequesterMemberID: "member", TrusteeOrgID: trusteeOrgID, Amount: decimal.NewFromInt(200), Currency: currency, Purpose: "works", Status: models.FundReleaseStatusCompleted},
		&models.Distribution{ID: assetID + "-distribution", AssetID: assetID, ProposedByOrgID: "manager", ProposedByMemberID: "member", TrusteeOrgID: trusteeOrgID, Amount: decimal.NewFromInt(50), Currency: currency, Source: "income", Status: models.DistributionStatusCompleted},
		&models.SegregatedAccount{ID: assetID + "-account", AssetID: assetID, CustodianOrgID: "custodian", AccountType: models.AccountTypeDevelopmentFund, AccountName: "Development", Balance: decimal.NewFromInt(400), Currency: currency, Status: models.AccountStatusActive},
	}
	for _, row := range rows {
		if err := database.Create(row).Error; err != nil {
			t.Fatalf("seed %T: %v", row, err)
		}
	}
}

func assertCurrencyAmount(t *testing.T, totals []models.CurrencyTotal, currency, amount string) {
	t.Helper()
	for _, total := range totals {
		if total.Currency == currency {
			if total.Amount != amount {
				t.Fatalf("%s amount = %s, want %s", currency, total.Amount, amount)
			}
			return
		}
	}
	t.Fatalf("currency %s missing from totals %+v", currency, totals)
}

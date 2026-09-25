package db

import (
	"context"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type tokenizedAssetSubscriptionFixture struct {
	ID               string `gorm:"primaryKey"`
	TokenizedAssetID string
	WalletAddress    string
}

func (tokenizedAssetSubscriptionFixture) TableName() string { return "tokenized_asset_subscriptions" }

func TestListAssetsFiltersWalletStakeholderAssignments(t *testing.T) {
	database, err := gorm.Open(sqlite.Open("file:tokenization-read-client?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := database.AutoMigrate(
		&TokenizedAsset{}, &tokenizedAssetSubscriptionFixture{}, &AssetTokenizationDocument{},
		&TokenizationFeeProofOfPayment{}, &WalletBank{}, &WalletTokenizationFee{}, &WalletTokenizationStatus{},
	); err != nil {
		t.Fatalf("migrate tokenized assets: %v", err)
	}
	bankID, feeID := uint64(4), uint64(8)
	assets := []TokenizedAsset{
		{
			ID: "asset-1", AssetCode: "ONE", AssetManagerID: 10, ApprovedAssetCustodianID: 20, TrusteeID: 30,
			FeeInFiat: 1, SECTokenizationFeeValue: 2, CustodianFeeValue: 3, AssetManagerFeeValue: 4,
			IssuingHouseFeeValue: 5, LegalAndProfessionalFeeValue: 6, RatingAgencyFeeValue: 7,
			TrusteeFeeValue: 8, VATValue: 9, NumberOfTokenToBeIssued: 1000,
			AssetDescription: "Income property", AssetLogo: "https://example.com/logo.png", OfferingType: "PRIVATE",
			BankID: &bankID, TokenizationFeeID: &feeID, NumberOfTokenToBeSold: 600, PricePerToken: 2,
			RiskSharingMechanismCompletionGuarantees: 1, RiskSharingMechanismPPPs: 1,
			RiskSharingMechanismHedgeInstruments: 1, IndependentMonitoringList: "Independent monitor",
			ESGSafeguardsSustainabilityCerts: 1, ESGSafeguardsCommunityEngagementPlan: 1,
			SecurityMeasuresSurveillanceSystems: 1, SecurityMeasuresOnSitePersonnel: 1,
			SecurityMeasuresPerimeterSecurity: 1, SecurityMeasuresCriticalInfrastructure: 1,
			LegalAdvisor: "Legal adviser", FinancialAdvisor: "Financial adviser",
		},
		{ID: "asset-2", AssetCode: "TWO", AssetManagerID: 11, ApprovedAssetCustodianID: 21, TrusteeID: 31},
	}
	if err := database.Create(&assets).Error; err != nil {
		t.Fatalf("seed tokenized assets: %v", err)
	}
	if err := database.Create(&WalletBank{ID: bankID, BankName: "Trovo Bank", CountryCode: "NG"}).Error; err != nil {
		t.Fatalf("seed bank: %v", err)
	}
	if err := database.Create(&WalletTokenizationFee{ID: feeID, FeeDescription: "Standard", FeeFiatPercentage: 2}).Error; err != nil {
		t.Fatalf("seed fee: %v", err)
	}
	if err := database.Create(&[]WalletTokenizationStatus{{ID: 0, Description: "Draft"}, {ID: 5, Description: "Primary sale"}}).Error; err != nil {
		t.Fatalf("seed statuses: %v", err)
	}
	if err := database.Create(&AssetTokenizationDocument{ID: 11, CreatedAt: time.Now(), TokenizedAssetID: "asset-1", DocumentType: "title", DocumentTitle: "Title deed", DocumentURL: "https://example.com/title.pdf"}).Error; err != nil {
		t.Fatalf("seed asset document: %v", err)
	}
	if err := database.Create(&TokenizationFeeProofOfPayment{ID: 12, CreatedAt: time.Now(), TokenizedAssetID: "asset-1", TokenizationFeePaymentMethodID: "FIAT", TransactionReference: "PAY-1", DocumentURL: "https://example.com/payment.pdf"}).Error; err != nil {
		t.Fatalf("seed payment proof: %v", err)
	}
	if err := database.Create(&[]tokenizedAssetSubscriptionFixture{
		{ID: "sub-1", TokenizedAssetID: "asset-1", WalletAddress: "wallet-a"},
		{ID: "sub-2", TokenizedAssetID: "asset-1", WalletAddress: "wallet-a"},
		{ID: "sub-3", TokenizedAssetID: "asset-1", WalletAddress: "wallet-b"},
	}).Error; err != nil {
		t.Fatalf("seed subscriptions: %v", err)
	}

	client := NewGormTokenizationReadClient(database)
	tests := []struct {
		name    string
		filters AssetFilters
	}{
		{name: "manager", filters: AssetFilters{ManagerStakeholderID: uint64Pointer(10)}},
		{name: "custodian", filters: AssetFilters{CustodianStakeholderID: uint64Pointer(20)}},
		{name: "trustee", filters: AssetFilters{TrusteeStakeholderID: uint64Pointer(30)}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			page, err := client.ListAssets(context.Background(), test.filters)
			if err != nil {
				t.Fatalf("list assets: %v", err)
			}
			if page.Total != 1 || len(page.Records) != 1 || page.Records[0].ID != "asset-1" {
				t.Fatalf("page = %+v, want only asset-1", page)
			}
			if page.Records[0].TrusteeID != 30 {
				t.Fatalf("trustee id = %d, want 30", page.Records[0].TrusteeID)
			}
			if page.Records[0].TokenHolderCount != 2 {
				t.Fatalf("token holder count = %d, want 2", page.Records[0].TokenHolderCount)
			}
		})
	}
	allAssets, err := client.ListAllAssets(context.Background(), AssetFilters{CustodianStakeholderID: uint64Pointer(20)})
	if err != nil {
		t.Fatalf("list all assets: %v", err)
	}
	if len(allAssets) != 1 || allAssets[0].CustodianFeeValue != 3 {
		t.Fatalf("all assets = %+v, want custodian fee value from WalletDB", allAssets)
	}
	if allAssets[0].ConfiguredFeeValue() != 45 {
		t.Fatalf("configured fee value = %v, want 45", allAssets[0].ConfiguredFeeValue())
	}
	detail, err := client.GetAsset(context.Background(), "asset-1")
	if err != nil {
		t.Fatalf("get asset: %v", err)
	}
	if detail.TokenHolderCount != 2 {
		t.Fatalf("detail token holder count = %d, want 2", detail.TokenHolderCount)
	}
	if detail.NumberOfTokenToBeIssued != 1000 {
		t.Fatalf("issued token count = %v, want 1000", detail.NumberOfTokenToBeIssued)
	}
	if detail.RiskSharingMechanismCompletionGuarantees != 1 || detail.IndependentMonitoringList != "Independent monitor" || detail.LegalAdvisor != "Legal adviser" {
		t.Fatalf("asset protection fields were not loaded from WalletDB: %+v", detail)
	}
	related, err := client.GetAssetRelatedData(context.Background(), *detail)
	if err != nil {
		t.Fatalf("get related asset data: %v", err)
	}
	if len(related.Documents) != 1 || related.Documents[0].DocumentTitle != "Title deed" {
		t.Fatalf("related documents = %+v", related.Documents)
	}
	if len(related.PaymentProofs) != 1 || related.PaymentProofs[0].TransactionReference != "PAY-1" {
		t.Fatalf("related payment proofs = %+v", related.PaymentProofs)
	}
	if related.Bank == nil || related.Bank.BankName != "Trovo Bank" || related.PreferredFee == nil || related.PreferredFee.ID != feeID {
		t.Fatalf("related bank/fee = %+v / %+v", related.Bank, related.PreferredFee)
	}
	if len(related.StatusCatalog) != 2 {
		t.Fatalf("status catalog = %+v", related.StatusCatalog)
	}
}

func uint64Pointer(value uint64) *uint64 {
	return &value
}

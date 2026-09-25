package services

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	stakeholderDB "admin-panel-dashboard/internal/components/stakeholder/db"
	"admin-panel-dashboard/internal/components/stakeholder/models"
	coreModels "admin-panel-dashboard/internal/models"

	"github.com/shopspring/decimal"
)

func TestListDistributionsFiltersByRoleAndStatus(t *testing.T) {
	db := testDB(t, &models.Distribution{}, &coreModels.Organization{}, &coreModels.OrganizationMember{})
	now := time.Now()
	trusteeOrgID := "trustee-org"
	managerOrg := coreModels.Organization{ID: "manager-org", Name: "Manager Org", Email: "manager-org@example.com", Type: "CORPORATE", Status: coreModels.OrganizationStatusActive, CreatedBy: "test"}
	managerMember := coreModels.OrganizationMember{ID: "manager-member", OrganizationID: managerOrg.ID, Email: "manager@example.com", Role: "ADMIN", Status: "ACTIVE", CreatedBy: "test", FirstName: "Ada", LastName: "Manager"}
	if err := db.Create(&managerOrg).Error; err != nil {
		t.Fatalf("seed manager organization: %v", err)
	}
	if err := db.Create(&managerMember).Error; err != nil {
		t.Fatalf("seed manager member: %v", err)
	}

	distributions := []models.Distribution{
		{
			ID:                 "dist-proposed",
			AssetID:            "asset-1",
			AssetCode:          "AST1",
			ProposedByOrgID:    "manager-org",
			ProposedByMemberID: "manager-member",
			TrusteeOrgID:       trusteeOrgID,
			Amount:             decimal.NewFromInt(100),
			Currency:           "CNGN",
			Source:             "wallet",
			Status:             models.DistributionStatusProposed,
			CreatedAt:          now.Add(-time.Minute),
			UpdatedAt:          now.Add(-time.Minute),
		},
		{
			ID:                 "dist-authorized",
			AssetID:            "asset-2",
			AssetCode:          "AST2",
			ProposedByOrgID:    "manager-org",
			ProposedByMemberID: "manager-member",
			TrusteeOrgID:       trusteeOrgID,
			Amount:             decimal.NewFromInt(200),
			Currency:           "CNGN",
			Source:             "wallet",
			Status:             models.DistributionStatusAuthorized,
			CreatedAt:          now.Add(-2 * time.Minute),
			UpdatedAt:          now.Add(-2 * time.Minute),
		},
		{
			ID:                 "dist-completed",
			AssetID:            "asset-3",
			AssetCode:          "AST3",
			ProposedByOrgID:    "manager-org",
			ProposedByMemberID: "manager-member",
			TrusteeOrgID:       trusteeOrgID,
			Amount:             decimal.NewFromInt(300),
			Currency:           "CNGN",
			Source:             "wallet",
			Status:             models.DistributionStatusCompleted,
			CreatedAt:          now.Add(-3 * time.Minute),
			UpdatedAt:          now.Add(-3 * time.Minute),
		},
		{
			ID:                 "dist-other-org",
			AssetID:            "asset-4",
			AssetCode:          "AST4",
			ProposedByOrgID:    "manager-org",
			ProposedByMemberID: "manager-member",
			TrusteeOrgID:       "other-trustee-org",
			Amount:             decimal.NewFromInt(400),
			Currency:           "CNGN",
			Source:             "wallet",
			Status:             models.DistributionStatusProposed,
			CreatedAt:          now.Add(-4 * time.Minute),
			UpdatedAt:          now.Add(-4 * time.Minute),
		},
	}

	for _, distribution := range distributions {
		if err := db.Create(&distribution).Error; err != nil {
			t.Fatalf("seed distribution: %v", err)
		}
	}

	service := NewFinancialService(db, nil, nil, nil, nil, nil)
	auth := AuthContext{
		MemberID:       "trustee-member",
		OrganizationID: trusteeOrgID,
		DashboardRole:  models.DashboardRoleTrustee,
	}

	pendingRecords, pendingTotal, err := service.ListDistributions(context.Background(), auth, 1, 20, false)
	if err != nil {
		t.Fatalf("list pending distributions: %v", err)
	}
	if pendingTotal != 1 {
		t.Fatalf("pending total = %d, want 1", pendingTotal)
	}
	if len(pendingRecords) != 1 {
		t.Fatalf("pending records = %d, want 1", len(pendingRecords))
	}
	if pendingRecords[0].Status != models.DistributionStatusProposed {
		t.Fatalf("pending record status = %s, want %s", pendingRecords[0].Status, models.DistributionStatusProposed)
	}
	if pendingRecords[0].Requester.Email != managerMember.Email || pendingRecords[0].Requester.OrganizationName != managerOrg.Name {
		t.Fatalf("pending requester = %+v", pendingRecords[0].Requester)
	}

	historyRecords, historyTotal, err := service.ListDistributions(context.Background(), auth, 1, 20, true)
	if err != nil {
		t.Fatalf("list history distributions: %v", err)
	}
	if historyTotal != 2 {
		t.Fatalf("history total = %d, want 2", historyTotal)
	}
	if len(historyRecords) != 2 {
		t.Fatalf("history records = %d, want 2", len(historyRecords))
	}
}

func TestGetDistributionReturnsCanonicalTokenBreakdown(t *testing.T) {
	database := testDB(t, &models.Distribution{}, &models.StakeholderAssetAssignment{}, &models.StakeholderDocument{}, &coreModels.Organization{}, &coreModels.OrganizationMember{})
	trusteeOrgID := "trustee-org"
	trusteeStakeholderID := uint64(30)
	distribution := models.Distribution{
		ID: "dist-detail", AssetID: "asset-1", AssetCode: "AST1", ProposedByOrgID: "manager-org",
		ProposedByMemberID: "manager-member", TrusteeOrgID: trusteeOrgID, Amount: decimal.NewFromInt(1250),
		Currency: "CNGN", Source: "operating income", Status: models.DistributionStatusProposed,
	}
	if err := database.Create(&distribution).Error; err != nil {
		t.Fatalf("seed distribution: %v", err)
	}
	managerOrg := coreModels.Organization{ID: distribution.ProposedByOrgID, Name: "Manager Org", Email: "detail-manager-org@example.com", Type: "CORPORATE", Status: coreModels.OrganizationStatusActive, CreatedBy: "test"}
	managerMember := coreModels.OrganizationMember{ID: distribution.ProposedByMemberID, OrganizationID: managerOrg.ID, Email: "detail-manager@example.com", Role: "ADMIN", Status: "ACTIVE", CreatedBy: "test", FirstName: "Chidi", LastName: "Manager"}
	if err := database.Create(&managerOrg).Error; err != nil {
		t.Fatalf("seed manager organization: %v", err)
	}
	if err := database.Create(&managerMember).Error; err != nil {
		t.Fatalf("seed manager member: %v", err)
	}
	if err := database.Create(&models.StakeholderAssetAssignment{
		ID: "assignment-1", AssetID: distribution.AssetID, AssetCode: distribution.AssetCode,
		TrusteeOrgID: &trusteeOrgID, TrusteeStakeholderID: &trusteeStakeholderID, Status: models.AssignmentStatusActive,
	}).Error; err != nil {
		t.Fatalf("seed assignment: %v", err)
	}
	assets := NewAssetService(database, fakeTokenizationReadClient{assets: []stakeholderDB.TokenizedAsset{
		{ID: distribution.AssetID, AssetCode: distribution.AssetCode, TrusteeID: trusteeStakeholderID, NumberOfTokenToBeIssued: 500, TokenHolderCount: 12},
	}}, nil)
	service := NewFinancialService(database, assets, nil, nil, nil, nil)
	auth := AuthContext{OrganizationID: trusteeOrgID, StakeholderID: trusteeStakeholderID, DashboardRole: models.DashboardRoleTrustee}

	detail, err := service.GetDistribution(context.Background(), auth, distribution.ID)
	if err != nil {
		t.Fatalf("get distribution: %v", err)
	}
	if detail.Breakdown.TotalTokens != "500" || detail.Breakdown.TotalTokenHolders != 12 {
		t.Fatalf("token breakdown = %+v", detail.Breakdown)
	}
	if detail.Breakdown.NetIncome.Amount != "1250" || detail.Breakdown.DistributionAmountPerToken.Amount != "2.5" {
		t.Fatalf("amount breakdown = %+v", detail.Breakdown)
	}
	if detail.Requester.Email != managerMember.Email || detail.Requester.OrganizationName != managerOrg.Name {
		t.Fatalf("requester = %+v", detail.Requester)
	}
	if detail.Payout.Status != "not_registered" || detail.Payout.Payouts == nil {
		t.Fatalf("default payout = %+v", detail.Payout)
	}
	service.payouts = &fakeDistributionPayoutClient{payout: &models.DistributionPayoutResponse{
		DistributionID: distribution.ID, TokenizedAssetID: distribution.AssetID, Amount: "1250", Currency: "CNGN",
		Payouts: []models.DistributionPayoutRecordResponse{{
			ID: "payout-1", BeneficiaryAddress: "GBENEFICIARY", ConfirmedTokenBalance: "10",
			Amount: "25", Currency: "CNGN", Paid: true, CreatedAt: time.Date(2026, time.September, 3, 6, 30, 0, 0, time.UTC),
		}},
	}}
	csvContent, err := service.DistributionPayoutCSV(context.Background(), auth, distribution.ID)
	if err != nil {
		t.Fatalf("export payout CSV: %v", err)
	}
	csvText := string(csvContent)
	if !strings.Contains(csvText, "payout_id,beneficiary_address,confirmed_token_balance,amount,currency,cannot_receive_asset,paid,created_at") ||
		!strings.Contains(csvText, "payout-1,GBENEFICIARY,10,25,CNGN,false,true,2026-09-03T06:30:00Z") {
		t.Fatalf("unexpected payout CSV:\n%s", csvText)
	}

	auth.OrganizationID = "other-trustee"
	if _, err := service.GetDistribution(context.Background(), auth, distribution.ID); ErrorStatus(err) != 404 {
		t.Fatalf("cross-organization detail status = %d, want 404", ErrorStatus(err))
	}
}

type fakeDistributionPayoutClient struct {
	registered *DistributionPayoutRegistration
	payout     *models.DistributionPayoutResponse
	err        error
}

func (f *fakeDistributionPayoutClient) Register(_ context.Context, registration DistributionPayoutRegistration) (*models.DistributionPayoutResponse, error) {
	f.registered = &registration
	return f.payout, f.err
}

func (f *fakeDistributionPayoutClient) Get(_ context.Context, _ string) (*models.DistributionPayoutResponse, error) {
	return f.payout, f.err
}

func TestAuthorizeDistributionRollsBackWhenWalletRegistrationFails(t *testing.T) {
	database := testDB(t, &models.Distribution{}, &models.StakeholderAuthorizationChallenge{})
	now := time.Now()
	distribution := models.Distribution{
		ID: "b1ec38f7-6bfd-46fb-b693-5dad05a905be", AssetID: "asset-1", AssetCode: "AST1",
		ProposedByOrgID: "manager-org", ProposedByMemberID: "manager-member", TrusteeOrgID: "trustee-org",
		Amount: decimal.NewFromInt(1250), Currency: "CNGN", Source: "income",
		Status: models.DistributionStatusProposed, CreatedAt: now, UpdatedAt: now,
	}
	if err := database.Create(&distribution).Error; err != nil {
		t.Fatal(err)
	}
	challenge := models.StakeholderAuthorizationChallenge{
		ID: "challenge-1", MemberID: "trustee-member", OrganizationID: distribution.TrusteeOrgID,
		StakeholderType: models.StakeholderTypeTrustee, DashboardRole: models.DashboardRoleTrustee,
		Action: models.ActionDistributionAuthorize, EntityType: models.EntityDistribution, EntityID: distribution.ID,
		AuthID: "wallet-auth-1", Status: models.AuthorizationStatusVerified, ExpiresAt: now.Add(time.Hour),
		CreatedAt: now, UpdatedAt: now,
	}
	if err := database.Create(&challenge).Error; err != nil {
		t.Fatal(err)
	}
	auth := AuthContext{MemberID: challenge.MemberID, OrganizationID: challenge.OrganizationID, DashboardRole: models.DashboardRoleTrustee}
	payouts := &fakeDistributionPayoutClient{err: errors.New("Wallet unavailable")}
	service := NewFinancialService(database, nil, nil, nil, NewAuthorizationService(database, nil), nil, payouts)

	if _, err := service.AuthorizeDistribution(context.Background(), auth, distribution.ID, models.ChallengeActionRequest{ChallengeID: challenge.ID}); err == nil {
		t.Fatal("authorize distribution succeeded while Wallet registration failed")
	}
	if err := database.First(&distribution, "id = ?", distribution.ID).Error; err != nil {
		t.Fatal(err)
	}
	if distribution.Status != models.DistributionStatusProposed {
		t.Fatalf("distribution status = %q after rollback", distribution.Status)
	}
	if err := database.First(&challenge, "id = ?", challenge.ID).Error; err != nil {
		t.Fatal(err)
	}
	if challenge.Status != models.AuthorizationStatusVerified {
		t.Fatalf("challenge status = %q after rollback", challenge.Status)
	}

	payouts.err = nil
	payouts.payout = &models.DistributionPayoutResponse{DistributionID: distribution.ID, Status: "registered", Payouts: []models.DistributionPayoutRecordResponse{}}
	updated, err := service.AuthorizeDistribution(context.Background(), auth, distribution.ID, models.ChallengeActionRequest{ChallengeID: challenge.ID})
	if err != nil {
		t.Fatalf("authorize retry: %v", err)
	}
	if updated.Status != models.DistributionStatusAuthorized || payouts.registered == nil || payouts.registered.DistributionID != distribution.ID {
		t.Fatalf("authorization result=%+v registration=%+v", updated, payouts.registered)
	}
	if err := database.First(&challenge, "id = ?", challenge.ID).Error; err != nil {
		t.Fatal(err)
	}
	if challenge.Status != models.AuthorizationStatusUsed {
		t.Fatalf("challenge status = %q, want used", challenge.Status)
	}
}

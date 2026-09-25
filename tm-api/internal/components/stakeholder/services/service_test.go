package services

import (
	"context"
	"errors"
	"net/url"
	"testing"
	"time"

	stakeholderDB "admin-panel-dashboard/internal/components/stakeholder/db"
	"admin-panel-dashboard/internal/components/stakeholder/models"
	coreModels "admin-panel-dashboard/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func testDB(t *testing.T, migrate ...interface{}) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+url.QueryEscape(t.Name())+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if len(migrate) > 0 {
		if err := db.AutoMigrate(migrate...); err != nil {
			t.Fatalf("automigrate: %v", err)
		}
	}
	return db
}

func TestAuditRecordRollsBackWithTransaction(t *testing.T) {
	db := testDB(t, &models.StakeholderAuditLog{})
	audit := NewAuditService(db)
	expectedErr := errors.New("rollback")

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := audit.Record(context.Background(), tx, AuditEvent{
			ActorMemberID: "member-1",
			ActorOrgID:    "org-1",
			ActorRole:     models.DashboardRoleAssetManager,
			Action:        models.ActionFundReleaseCreate,
			EntityType:    models.EntityFundReleaseRequest,
			EntityID:      "fr-1",
		}); err != nil {
			t.Fatalf("record audit: %v", err)
		}
		return expectedErr
	})
	if !errors.Is(err, expectedErr) {
		t.Fatalf("transaction err = %v, want %v", err, expectedErr)
	}

	var count int64
	if err := db.Model(&models.StakeholderAuditLog{}).Count(&count).Error; err != nil {
		t.Fatalf("count audit logs: %v", err)
	}
	if count != 0 {
		t.Fatalf("audit rows = %d, want 0 after rollback", count)
	}
}

func TestAuditRecordNormalizesNilStateMaps(t *testing.T) {
	db := testDB(t, &models.StakeholderAuditLog{})
	audit := NewAuditService(db)

	if err := audit.Record(context.Background(), db, AuditEvent{
		ActorMemberID: "member-1",
		ActorOrgID:    "org-1",
		ActorRole:     models.DashboardRoleAssetManager,
		Action:        models.ActionProfileUpdate,
		EntityType:    models.EntityProfile,
		EntityID:      "member-1",
	}); err != nil {
		t.Fatalf("record audit with nil maps: %v", err)
	}

	var log models.StakeholderAuditLog
	if err := db.First(&log).Error; err != nil {
		t.Fatalf("load audit log: %v", err)
	}
	if log.BeforeState == nil || log.AfterState == nil || log.Metadata == nil {
		t.Fatalf("audit maps should be normalized, got before=%v after=%v metadata=%v", log.BeforeState, log.AfterState, log.Metadata)
	}
}

func TestAuthorizationChallengeConsumeIsOneTimeAndBound(t *testing.T) {
	db := testDB(t, &models.StakeholderAuthorizationChallenge{})
	service := NewAuthorizationService(db, nil)
	auth := AuthContext{
		MemberID:       "member-1",
		OrganizationID: "org-1",
		DashboardRole:  models.DashboardRoleTrustee,
	}
	challenge := models.StakeholderAuthorizationChallenge{
		ID:              "challenge-1",
		MemberID:        auth.MemberID,
		OrganizationID:  auth.OrganizationID,
		StakeholderType: models.StakeholderTypeTrustee,
		DashboardRole:   auth.DashboardRole,
		Action:          models.ActionFundReleaseApprove,
		EntityType:      models.EntityFundReleaseRequest,
		EntityID:        "fr-1",
		AuthID:          "auth-1",
		Status:          models.AuthorizationStatusVerified,
		ExpiresAt:       time.Now().Add(time.Minute),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	if err := db.Create(&challenge).Error; err != nil {
		t.Fatalf("seed challenge: %v", err)
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		return service.ConsumeVerifiedChallenge(context.Background(), tx, auth, models.ActionFundReleaseApprove, models.EntityFundReleaseRequest, "fr-1", challenge.ID)
	})
	if err != nil {
		t.Fatalf("consume verified challenge: %v", err)
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		return service.ConsumeVerifiedChallenge(context.Background(), tx, auth, models.ActionFundReleaseApprove, models.EntityFundReleaseRequest, "fr-1", challenge.ID)
	})
	if err == nil {
		t.Fatal("expected second consume to fail")
	}
	if ErrorStatus(err) != 403 {
		t.Fatalf("second consume status = %d, want 403", ErrorStatus(err))
	}
}

func TestAuthorizationChallengeExpiredCannotBeConsumed(t *testing.T) {
	db := testDB(t, &models.StakeholderAuthorizationChallenge{})
	service := NewAuthorizationService(db, nil)
	auth := AuthContext{MemberID: "member-1", OrganizationID: "org-1", DashboardRole: models.DashboardRoleTrustee}
	challenge := models.StakeholderAuthorizationChallenge{
		ID:              "challenge-expired",
		MemberID:        auth.MemberID,
		OrganizationID:  auth.OrganizationID,
		StakeholderType: models.StakeholderTypeTrustee,
		DashboardRole:   auth.DashboardRole,
		Action:          models.ActionDistributionAuthorize,
		EntityType:      models.EntityDistribution,
		EntityID:        "dist-1",
		AuthID:          "auth-1",
		Status:          models.AuthorizationStatusVerified,
		ExpiresAt:       time.Now().Add(-time.Minute),
		CreatedAt:       time.Now().Add(-2 * time.Minute),
		UpdatedAt:       time.Now().Add(-2 * time.Minute),
	}
	if err := db.Create(&challenge).Error; err != nil {
		t.Fatalf("seed challenge: %v", err)
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		return service.ConsumeVerifiedChallenge(context.Background(), tx, auth, models.ActionDistributionAuthorize, models.EntityDistribution, "dist-1", challenge.ID)
	})
	if err == nil {
		t.Fatal("expected expired challenge to fail")
	}
	if ErrorStatus(err) != 403 {
		t.Fatalf("expired consume status = %d, want 403", ErrorStatus(err))
	}
}

func TestAuthorizationChallengeWrongActionCannotBeConsumed(t *testing.T) {
	db := testDB(t, &models.StakeholderAuthorizationChallenge{})
	service := NewAuthorizationService(db, nil)
	auth := AuthContext{MemberID: "member-1", OrganizationID: "org-1", DashboardRole: models.DashboardRoleTrustee}
	challenge := models.StakeholderAuthorizationChallenge{
		ID:              "challenge-wrong-action",
		MemberID:        auth.MemberID,
		OrganizationID:  auth.OrganizationID,
		StakeholderType: models.StakeholderTypeTrustee,
		DashboardRole:   auth.DashboardRole,
		Action:          models.ActionFundReleaseApprove,
		EntityType:      models.EntityFundReleaseRequest,
		EntityID:        "fr-1",
		AuthID:          "auth-1",
		Status:          models.AuthorizationStatusVerified,
		ExpiresAt:       time.Now().Add(time.Minute),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	if err := db.Create(&challenge).Error; err != nil {
		t.Fatalf("seed challenge: %v", err)
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		return service.ConsumeVerifiedChallenge(context.Background(), tx, auth, models.ActionDistributionAuthorize, models.EntityDistribution, "dist-1", challenge.ID)
	})
	if err == nil {
		t.Fatal("expected wrong-action challenge consumption to fail")
	}
	if ErrorStatus(err) != 403 {
		t.Fatalf("wrong-action consume status = %d, want 403", ErrorStatus(err))
	}

	var stored models.StakeholderAuthorizationChallenge
	if err := db.First(&stored, "id = ?", challenge.ID).Error; err != nil {
		t.Fatalf("load challenge: %v", err)
	}
	if stored.Status != models.AuthorizationStatusVerified || stored.UsedAt != nil {
		t.Fatalf("challenge status=%q used_at=%v, want verified and unused", stored.Status, stored.UsedAt)
	}
}

type fakeTokenizationReadClient struct {
	assets []stakeholderDB.TokenizedAsset
}

type fakeTokenizationDetailReadClient struct {
	fakeTokenizationReadClient
	related map[string]stakeholderDB.TokenizedAssetRelatedData
}

func (f fakeTokenizationDetailReadClient) GetAssetRelatedData(_ context.Context, asset stakeholderDB.TokenizedAsset) (stakeholderDB.TokenizedAssetRelatedData, error) {
	if related, ok := f.related[asset.ID]; ok {
		return related, nil
	}
	return stakeholderDB.TokenizedAssetRelatedData{
		Documents: []stakeholderDB.AssetTokenizationDocument{}, PaymentProofs: []stakeholderDB.TokenizationFeeProofOfPayment{}, StatusCatalog: []stakeholderDB.WalletTokenizationStatus{},
	}, nil
}

func (f fakeTokenizationReadClient) ListAssets(ctx context.Context, filters stakeholderDB.AssetFilters) (stakeholderDB.AssetPage, error) {
	return stakeholderDB.AssetPage{Records: f.assets, Total: int64(len(f.assets)), Page: 1, Limit: 20}, nil
}

func (f fakeTokenizationReadClient) ListAllAssets(ctx context.Context, filters stakeholderDB.AssetFilters) ([]stakeholderDB.TokenizedAsset, error) {
	return f.assets, nil
}

func (f fakeTokenizationReadClient) GetAsset(ctx context.Context, idOrCode string) (*stakeholderDB.TokenizedAsset, error) {
	for _, asset := range f.assets {
		if asset.ID == idOrCode || asset.AssetCode == idOrCode {
			return &asset, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func TestAssetServiceFiltersVisibilityByRole(t *testing.T) {
	db := testDB(t, &models.StakeholderAssetAssignment{}, &coreModels.Organization{})
	managerOrgID := "manager-org"
	trusteeOrgID := "trustee-org"
	managerStakeholderID := uint64(11)
	trusteeStakeholderID := uint64(22)
	if err := db.Create(&models.StakeholderAssetAssignment{
		ID:                        "assignment-1",
		AssetID:                   "asset-1",
		AssetCode:                 "AST1",
		TrusteeOrgID:              &trusteeOrgID,
		TrusteeStakeholderID:      &trusteeStakeholderID,
		AssetManagerOrgID:         &managerOrgID,
		AssetManagerStakeholderID: &managerStakeholderID,
		Status:                    models.AssignmentStatusActive,
	}).Error; err != nil {
		t.Fatalf("seed assignment: %v", err)
	}
	service := NewAssetService(db, fakeTokenizationReadClient{assets: []stakeholderDB.TokenizedAsset{
		{ID: "asset-1", AssetCode: "AST1", AssetName: "Visible", AssetManagerID: managerStakeholderID},
		{ID: "asset-2", AssetCode: "AST2", AssetName: "Hidden", AssetManagerID: 999},
	}}, nil)

	managerResp, err := service.ListAssets(context.Background(), AuthContext{
		OrganizationID: managerOrgID,
		StakeholderID:  managerStakeholderID,
		DashboardRole:  models.DashboardRoleAssetManager,
	}, models.AssetListFilters{Page: 1, Limit: 20})
	if err != nil {
		t.Fatalf("manager list assets: %v", err)
	}
	managerRecords := managerResp.Records.([]models.TokenizedAssetResponse)
	if len(managerRecords) != 1 || managerRecords[0].ID != "asset-1" {
		t.Fatalf("manager records = %+v, want only asset-1", managerRecords)
	}
	if managerResp.Summary == nil || managerResp.Summary.Total != 1 {
		t.Fatalf("manager summary = %+v, want total 1", managerResp.Summary)
	}

	trusteeResp, err := service.ListAssets(context.Background(), AuthContext{
		OrganizationID: trusteeOrgID,
		StakeholderID:  trusteeStakeholderID,
		DashboardRole:  models.DashboardRoleTrustee,
	}, models.AssetListFilters{Page: 1, Limit: 20})
	if err != nil {
		t.Fatalf("trustee list assets: %v", err)
	}
	trusteeRecords := trusteeResp.Records.([]models.TokenizedAssetResponse)
	if len(trusteeRecords) != 1 || trusteeRecords[0].ID != "asset-1" {
		t.Fatalf("trustee records = %+v, want only explicitly assigned asset-1", trusteeRecords)
	}

	unassignedTrusteeResp, err := service.ListAssets(context.Background(), AuthContext{
		OrganizationID: "other-trustee",
		StakeholderID:  333,
		DashboardRole:  models.DashboardRoleTrustee,
	}, models.AssetListFilters{Page: 1, Limit: 20})
	if err != nil {
		t.Fatalf("unassigned trustee list assets: %v", err)
	}
	unassignedRecords := unassignedTrusteeResp.Records.([]models.TokenizedAssetResponse)
	if len(unassignedRecords) != 0 {
		t.Fatalf("unassigned trustee records = %+v, want none", unassignedRecords)
	}
}

func TestAssetServiceReturnsLifecycleStatisticsAndRichDetails(t *testing.T) {
	db := testDB(t, &models.StakeholderAssetAssignment{}, &models.StakeholderDocument{}, &models.StakeholderAssetOperation{}, &models.DueDiligenceChecklist{}, &coreModels.Organization{})
	managerID, custodianID, trusteeID := uint64(71), uint64(72), uint64(73)
	managerOrgID := "rich-manager-org"
	organizations := []coreModels.Organization{
		{ID: managerOrgID, Name: "Rich Manager", Email: "manager@example.com", Type: "CORPORATE", Status: coreModels.OrganizationStatusActive, CreatedBy: "test", StakeholderID: &managerID, StakeholderType: stringPointer(models.StakeholderTypeAssetManager)},
		{ID: "rich-custodian-org", Name: "Rich Custodian", Email: "custodian@example.com", Type: "CORPORATE", Status: coreModels.OrganizationStatusActive, CreatedBy: "test", StakeholderID: &custodianID, StakeholderType: stringPointer(models.StakeholderTypeAssetCustodian)},
		{ID: "rich-trustee-org", Name: "Rich Trustee", Email: "trustee@example.com", Type: "CORPORATE", Status: coreModels.OrganizationStatusActive, CreatedBy: "test", StakeholderID: &trusteeID, StakeholderType: stringPointer(models.StakeholderTypeTrustee)},
	}
	if err := db.Create(&organizations).Error; err != nil {
		t.Fatalf("seed organizations: %v", err)
	}
	asset := stakeholderDB.TokenizedAsset{
		ID: "rich-asset", AssetCode: "RICH", AssetName: "Rich Asset",
		AssetManagerID: managerID, ApprovedAssetCustodianID: custodianID, TrusteeID: trusteeID,
		AssetTokenizationStatus: models.AssetTokenizationStatusPrimarySale,
		VettingStatus:           1, TokenHolderCount: 4, InitiatorUsername: "tokenizer",
		IssuingWalletAddress: "GISSUING", IssuingWalletAlias: "issuer_wallet",
		ExemptedCountries: "US, CA", AssetDescription: "Rich details", AssetLogo: "https://example.com/rich.png",
		OfferingType: "PRIVATE", AssetCountryLocation: "NG", AssetPhysicalAddress: "Lagos",
		OwnershipType: "DIRECT", OwnershipKind: "CORPORATE", AssetOwnerName: "Rich Owner",
		NumberOfTokenToBeSold: 60, TotalTokenHeldByManager: 40, PricePerToken: 10,
		FeeInFiat: 100, FeeInAsset: 5, TokenizationApplicationFee: 50, TokenizationApplicationFeeAsset: "CNGN",
		RiskSharingMechanismCompletionGuarantees: 1, RiskSharingMechanismPPPs: 1,
		RiskSharingMechanismHedgeInstruments: 1, IndependentMonitoringList: "Deloitte",
		ESGSafeguardsSustainabilityCerts: 1, ESGSafeguardsCommunityEngagementPlan: 1,
		SecurityMeasuresSurveillanceSystems: 1, SecurityMeasuresOnSitePersonnel: 1,
		SecurityMeasuresPerimeterSecurity: 1, SecurityMeasuresCriticalInfrastructure: 1,
		LegalAdvisor: "Legal Counsel", FinancialAdvisor: "Financial Counsel",
	}
	liquidated := stakeholderDB.TokenizedAsset{
		ID: "liquidated-asset", AssetCode: "LIQ", AssetManagerID: managerID,
		AssetTokenizationStatus: models.AssetTokenizationStatusLiquidated, DueDiligenceFail: 1,
	}
	if err := db.Create(&models.StakeholderDocument{
		ID: "rich-doc", AssetID: asset.ID, UploadedByOrgID: managerOrgID, UploadedByMemberID: "member",
		Category: models.DocumentCategoryGeneral, Title: "Asset document", FileURL: "https://example.com/doc.pdf",
		AccessRoles: models.JSONStringArray{models.DashboardRoleAssetManager}, Status: models.DocumentStatusActive,
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}).Error; err != nil {
		t.Fatalf("seed document: %v", err)
	}
	auditedAt := time.Now().Add(-time.Hour).UTC().Truncate(time.Second)
	if err := db.Create(&models.DueDiligenceChecklist{
		ID: "rich-audit", AssetID: asset.ID, AssetCode: asset.AssetCode, TrusteeOrgID: "rich-trustee-org",
		Status: models.DueDiligenceStatusApproved, ApprovedAt: &auditedAt, CreatedAt: auditedAt, UpdatedAt: auditedAt,
	}).Error; err != nil {
		t.Fatalf("seed completed due-diligence audit: %v", err)
	}
	service := NewAssetService(db, fakeTokenizationDetailReadClient{
		fakeTokenizationReadClient: fakeTokenizationReadClient{assets: []stakeholderDB.TokenizedAsset{asset, liquidated}},
		related: map[string]stakeholderDB.TokenizedAssetRelatedData{asset.ID: {
			Documents:     []stakeholderDB.AssetTokenizationDocument{{ID: 50, TokenizedAssetID: asset.ID, DocumentTitle: "Wallet title deed", DocumentURL: "https://wallet.example/title.pdf"}},
			PaymentProofs: []stakeholderDB.TokenizationFeeProofOfPayment{{ID: 51, TokenizedAssetID: asset.ID, TransactionReference: "PAY-51"}},
			Bank:          &stakeholderDB.WalletBank{ID: 2, BankName: "Rich Bank", CountryCode: "NG"},
			PreferredFee:  &stakeholderDB.WalletTokenizationFee{ID: 3, FeeDescription: "Preferred"},
			StatusCatalog: []stakeholderDB.WalletTokenizationStatus{{ID: 1, Description: "Submitted"}, {ID: 5, Description: "Primary sale"}, {ID: 7, Description: "Liquidated"}},
		}},
	}, nil)
	auth := AuthContext{OrganizationID: managerOrgID, StakeholderID: managerID, DashboardRole: models.DashboardRoleAssetManager}

	page, err := service.ListAssets(context.Background(), auth, models.AssetListFilters{})
	if err != nil {
		t.Fatalf("list assets: %v", err)
	}
	if page.Summary == nil || page.Summary.Total != 2 || page.Summary.Active != 1 || page.Summary.Liquidated != 1 {
		t.Fatalf("summary = %+v, want total=2 active=1 liquidated=1", page.Summary)
	}
	records := page.Records.([]models.TokenizedAssetResponse)
	if records[0].ComplianceStatus != models.AssetComplianceStatusVerified || records[0].CustodyStatus != models.AssetCustodyStatusActive || records[0].TokenHolderCount != 4 {
		t.Fatalf("asset lifecycle response = %+v", records[0])
	}

	detail, err := service.GetAssetDetails(context.Background(), auth, asset.ID)
	if err != nil {
		t.Fatalf("get asset details: %v", err)
	}
	if detail.TokenizerUsername != "tokenizer" || detail.Wallets.IssuingWalletAddress != "GISSUING" || len(detail.ExemptedCountries) != 2 {
		t.Fatalf("detail = %+v", detail)
	}
	if len(detail.AssignedStakeholders) != 3 || len(detail.Documents) != 1 || detail.Documents[0].ID != "rich-doc" {
		t.Fatalf("detail relationships = %+v", detail)
	}
	if detail.AssetProfile.LogoURL == "" || detail.Ownership.OwnerName != "Rich Owner" || detail.Offering.AmountToBeRaised != 600 {
		t.Fatalf("rich asset groups = %+v", detail)
	}
	if detail.RiskAndCompliance.RiskAssessmentScore != nil || detail.RiskAndCompliance.MaximumRiskScore != 100 || detail.RiskAndCompliance.RiskLevel != nil || detail.RiskAndCompliance.ComplianceStatus != models.AssetComplianceStatusVerified || detail.RiskAndCompliance.LastAuditDate == nil || !detail.RiskAndCompliance.LastAuditDate.Equal(auditedAt) {
		t.Fatalf("risk and compliance details = %+v", detail.RiskAndCompliance)
	}
	if !detail.Protection.CompletionGuarantees || !detail.Protection.PublicPrivatePartnerships || !detail.Protection.HedgeInstruments || detail.Protection.IndependentMonitoringList != "Deloitte" || !detail.Protection.ESGSustainabilityCertifications || !detail.Protection.ESGCommunityEngagementPlan || !detail.Protection.SurveillanceSystems || !detail.Protection.OnSiteSecurityPersonnel || !detail.Protection.PerimeterSecurity || !detail.Protection.CriticalInfrastructureProtection || detail.Protection.LegalAdviser != "Legal Counsel" || detail.Protection.FinancialAdviser != "Financial Counsel" {
		t.Fatalf("expanded asset protection = %+v", detail.Protection)
	}
	if detail.ReceivingAccount.BankName != "Rich Bank" || len(detail.ApplicationFee.Payments) != 1 || len(detail.TokenizationDocuments) != 1 {
		t.Fatalf("rich asset relations = %+v", detail)
	}
	if len(detail.Timeline.Stages) != 3 || !detail.Timeline.Stages[1].Current {
		t.Fatalf("timeline = %+v", detail.Timeline)
	}
}

func TestAssetServiceUsesWalletVettingAssignmentsForEveryRole(t *testing.T) {
	db := testDB(t, &models.StakeholderAssetAssignment{}, &coreModels.Organization{})
	managerID, custodianID, trusteeID := uint64(20), uint64(17), uint64(31)
	organizations := []coreModels.Organization{
		{ID: "manager-org-wallet", Name: "Manager", Email: "manager-wallet@example.com", Type: "CORPORATE", Status: coreModels.OrganizationStatusActive, CreatedBy: "test", StakeholderID: &managerID, StakeholderType: stringPointer(models.StakeholderTypeAssetManager)},
		{ID: "custodian-org-wallet", Name: "Custodian", Email: "custodian-wallet@example.com", Type: "CORPORATE", Status: coreModels.OrganizationStatusActive, CreatedBy: "test", StakeholderID: &custodianID, StakeholderType: stringPointer(models.StakeholderTypeAssetCustodian)},
		{ID: "trustee-org-wallet", Name: "Trustee", Email: "trustee-wallet@example.com", Type: "CORPORATE", Status: coreModels.OrganizationStatusActive, CreatedBy: "test", StakeholderID: &trusteeID, StakeholderType: stringPointer(models.StakeholderTypeTrustee)},
	}
	if err := db.Create(&organizations).Error; err != nil {
		t.Fatalf("seed organizations: %v", err)
	}

	asset := stakeholderDB.TokenizedAsset{
		ID:                       "wallet-assigned-asset",
		AssetCode:                "WALLET1",
		AssetManagerID:           managerID,
		ApprovedAssetCustodianID: custodianID,
		TrusteeID:                trusteeID,
	}
	service := NewAssetService(db, fakeTokenizationReadClient{assets: []stakeholderDB.TokenizedAsset{asset}}, nil)

	roles := []AuthContext{
		{OrganizationID: "manager-org-wallet", StakeholderID: managerID, DashboardRole: models.DashboardRoleAssetManager},
		{OrganizationID: "custodian-org-wallet", StakeholderID: custodianID, DashboardRole: models.DashboardRoleAssetCustodian},
		{OrganizationID: "trustee-org-wallet", StakeholderID: trusteeID, DashboardRole: models.DashboardRoleTrustee},
	}
	for _, auth := range roles {
		response, err := service.ListAssets(context.Background(), auth, models.AssetListFilters{})
		if err != nil {
			t.Fatalf("%s list assets: %v", auth.DashboardRole, err)
		}
		records := response.Records.([]models.TokenizedAssetResponse)
		if len(records) != 1 || records[0].ID != asset.ID {
			t.Fatalf("%s records = %+v, want wallet-assigned asset", auth.DashboardRole, records)
		}
	}

	assignment, _, err := service.RequireAssignmentForAsset(context.Background(), roles[0], asset.ID)
	if err != nil {
		t.Fatalf("require wallet assignment: %v", err)
	}
	if assignment.AssetManagerOrgID == nil || *assignment.AssetManagerOrgID != "manager-org-wallet" ||
		assignment.CustodianOrgID == nil || *assignment.CustodianOrgID != "custodian-org-wallet" ||
		assignment.TrusteeOrgID == nil || *assignment.TrusteeOrgID != "trustee-org-wallet" {
		t.Fatalf("resolved assignment = %+v, want all linked portal organizations", assignment)
	}
}

func TestAssetServiceInactivePortalAssignmentOverridesWalletFallback(t *testing.T) {
	db := testDB(t, &models.StakeholderAssetAssignment{}, &coreModels.Organization{})
	trusteeID := uint64(41)
	if err := db.Create(&models.StakeholderAssetAssignment{
		ID: "inactive-assignment", AssetID: "asset-inactive", Status: models.AssignmentStatusInactive,
	}).Error; err != nil {
		t.Fatalf("seed inactive assignment: %v", err)
	}
	service := NewAssetService(db, fakeTokenizationReadClient{assets: []stakeholderDB.TokenizedAsset{
		{ID: "asset-inactive", TrusteeID: trusteeID},
	}}, nil)

	response, err := service.ListAssets(context.Background(), AuthContext{
		OrganizationID: "trustee-org", StakeholderID: trusteeID, DashboardRole: models.DashboardRoleTrustee,
	}, models.AssetListFilters{})
	if err != nil {
		t.Fatalf("list assets: %v", err)
	}
	if records := response.Records.([]models.TokenizedAssetResponse); len(records) != 0 {
		t.Fatalf("records = %+v, want inactive explicit assignment to revoke visibility", records)
	}
}

func stringPointer(value string) *string {
	return &value
}

func TestDocumentServiceUnscopedDocumentsAreScopedToUploaderOrg(t *testing.T) {
	db := testDB(t, &models.StakeholderDocument{})
	doc := models.StakeholderDocument{
		ID:                 "doc-unscoped",
		UploadedByOrgID:    "manager-org-a",
		UploadedByMemberID: "manager-member-a",
		Category:           "internal",
		Title:              "Unscoped document",
		FileURL:            "https://example.com/doc.pdf",
		AccessRoles:        models.JSONStringArray{models.DashboardRoleAssetManager},
		Status:             models.DocumentStatusActive,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
	if err := db.Create(&doc).Error; err != nil {
		t.Fatalf("seed document: %v", err)
	}

	service := NewDocumentService(db, nil, nil, nil)
	ownDocs, total, err := service.List(context.Background(), AuthContext{
		OrganizationID: "manager-org-a",
		DashboardRole:  models.DashboardRoleAssetManager,
	}, 1, 20, "")
	if err != nil {
		t.Fatalf("list own unscoped docs: %v", err)
	}
	if total != 1 || len(ownDocs) != 1 || ownDocs[0].ID != doc.ID {
		t.Fatalf("own docs = %+v, total=%d, want uploaded document", ownDocs, total)
	}

	otherDocs, total, err := service.List(context.Background(), AuthContext{
		OrganizationID: "manager-org-b",
		DashboardRole:  models.DashboardRoleAssetManager,
	}, 1, 20, "")
	if err != nil {
		t.Fatalf("list other unscoped docs: %v", err)
	}
	if total != 0 || len(otherDocs) != 0 {
		t.Fatalf("other org docs = %+v, total=%d, want none", otherDocs, total)
	}

	_, err = service.Get(context.Background(), AuthContext{
		OrganizationID: "manager-org-b",
		DashboardRole:  models.DashboardRoleAssetManager,
	}, doc.ID)
	if err == nil {
		t.Fatal("expected other org document detail lookup to fail")
	}
	if ErrorStatus(err) != 404 {
		t.Fatalf("document detail status = %d, want 404", ErrorStatus(err))
	}
}

func TestAssetServiceFKOnlyRolesSeeAssignedAssets(t *testing.T) {
	db := testDB(t, &models.StakeholderAssetAssignment{}, &coreModels.Organization{})

	legalID, financialID, issuingID, ratingID := uint64(51), uint64(52), uint64(53), uint64(54)
	service := NewAssetService(db, fakeTokenizationReadClient{assets: []stakeholderDB.TokenizedAsset{
		{ID: "asset-mine", AssetCode: "MINE", AssetName: "Mine",
			LegalAdviserID: legalID, FinancialAdviserID: financialID,
			AssetIssuingHouseID: issuingID, RatingAgencyID: ratingID},
		{ID: "asset-other", AssetCode: "OTHER", AssetName: "Other",
			LegalAdviserID: 999, FinancialAdviserID: 999,
			AssetIssuingHouseID: 999, RatingAgencyID: 999},
	}}, nil)

	cases := []struct {
		role          string
		stakeholderID uint64
	}{
		{models.DashboardRoleLegalAdviser, legalID},
		{models.DashboardRoleFinancialAdviser, financialID},
		{models.DashboardRoleIssuingHouse, issuingID},
		{models.DashboardRoleRatingAgency, ratingID},
	}
	for _, tc := range cases {
		t.Run(tc.role, func(t *testing.T) {
			resp, err := service.ListAssets(context.Background(), AuthContext{
				OrganizationID: "org-" + tc.role,
				StakeholderID:  tc.stakeholderID,
				DashboardRole:  tc.role,
			}, models.AssetListFilters{Page: 1, Limit: 20})
			if err != nil {
				t.Fatalf("%s list assets: %v", tc.role, err)
			}
			records := resp.Records.([]models.TokenizedAssetResponse)
			if len(records) != 1 || records[0].ID != "asset-mine" {
				t.Fatalf("%s records = %+v, want only asset-mine", tc.role, records)
			}

			// A stakeholder id that matches no asset FK sees nothing.
			empty, err := service.ListAssets(context.Background(), AuthContext{
				OrganizationID: "org-none",
				StakeholderID:  8888,
				DashboardRole:  tc.role,
			}, models.AssetListFilters{Page: 1, Limit: 20})
			if err != nil {
				t.Fatalf("%s empty list: %v", tc.role, err)
			}
			if len(empty.Records.([]models.TokenizedAssetResponse)) != 0 {
				t.Fatalf("%s unassigned should see no assets", tc.role)
			}
		})
	}
}

func TestStructuringConfirmComplete(t *testing.T) {
	db := testDB(t,
		&models.StakeholderAssetAssignment{}, &coreModels.Organization{},
		&models.StakeholderStructuringStatus{}, &models.StakeholderAuditLog{},
		&models.StakeholderDocument{},
	)

	legalID := uint64(61)
	asset := stakeholderDB.TokenizedAsset{ID: "asset-la", AssetCode: "LA1", LegalAdviserID: legalID}
	assets := NewAssetService(db, fakeTokenizationReadClient{assets: []stakeholderDB.TokenizedAsset{asset}}, nil)
	audit := NewAuditService(db)
	documents := NewDocumentService(db, assets, audit, nil)
	svc := NewStructuringService(db, assets, documents, audit)

	legalAuth := AuthContext{OrganizationID: "la-org", StakeholderID: legalID, MemberID: "m1", DashboardRole: models.DashboardRoleLegalAdviser}

	// Happy path.
	status, err := svc.ConfirmComplete(context.Background(), legalAuth, "asset-la", models.ConfirmStructuringRequest{Notes: "done"})
	if err != nil {
		t.Fatalf("confirm complete: %v", err)
	}
	if status.Status != models.StructuringStatusComplete || status.Workstream != models.StructuringWorkstreamLegal {
		t.Fatalf("status = %+v, want complete/legal", status)
	}

	// Re-confirm → 409.
	if _, err := svc.ConfirmComplete(context.Background(), legalAuth, "asset-la", models.ConfirmStructuringRequest{}); ErrorStatus(err) != 409 {
		t.Fatalf("re-confirm status = %d, want 409", ErrorStatus(err))
	}

	// Non-adviser role → 403.
	trusteeAuth := AuthContext{OrganizationID: "t-org", StakeholderID: 1, MemberID: "m2", DashboardRole: models.DashboardRoleTrustee}
	if _, err := svc.ConfirmComplete(context.Background(), trusteeAuth, "asset-la", models.ConfirmStructuringRequest{}); ErrorStatus(err) != 403 {
		t.Fatalf("trustee status = %d, want 403", ErrorStatus(err))
	}

	// Adviser who does not own the asset → 404 (asset not visible).
	otherAuth := AuthContext{OrganizationID: "o-org", StakeholderID: 9999, MemberID: "m3", DashboardRole: models.DashboardRoleLegalAdviser}
	if _, err := svc.ConfirmComplete(context.Background(), otherAuth, "asset-la", models.ConfirmStructuringRequest{}); ErrorStatus(err) != 404 {
		t.Fatalf("non-owner status = %d, want 404", ErrorStatus(err))
	}
}

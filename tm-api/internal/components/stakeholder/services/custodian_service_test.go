package services

import (
	"context"
	"testing"
	"time"

	stakeholderDB "admin-panel-dashboard/internal/components/stakeholder/db"
	"admin-panel-dashboard/internal/components/stakeholder/models"
	coreModels "admin-panel-dashboard/internal/models"

	"github.com/google/uuid"
)

func TestCreateComplianceCreatesPendingItemAndAudit(t *testing.T) {
	database := testDB(t, &coreModels.Organization{}, &models.ComplianceItem{}, &models.StakeholderAuditLog{})
	stakeholderType := models.StakeholderTypeAssetCustodian
	stakeholderID := uint64(17)
	organization := coreModels.Organization{
		ID: "custodian-org", Name: "Custodian", Email: "custodian@example.com",
		Type: coreModels.OrganizationTypeAssetCustodian, Status: coreModels.OrganizationStatusActive,
		CreatedBy: "admin", StakeholderID: &stakeholderID, StakeholderType: &stakeholderType,
	}
	if err := database.Create(&organization).Error; err != nil {
		t.Fatalf("seed custodian organization: %v", err)
	}

	service := NewCustodianService(database, nil, nil, NewAuditService(database))
	actor := AuthContext{MemberID: "admin-1", DashboardRole: "trovo_admin"}
	request := models.CreateComplianceItemRequest{
		CustodianOrgID: organization.ID,
		Category:       " Regulatory ",
		Requirement:    " Monthly custody report ",
		DueDate:        "2026-09-30",
	}

	item, err := service.CreateCompliance(context.Background(), actor, request)
	if err != nil {
		t.Fatalf("create compliance item: %v", err)
	}
	if item.OrgID != organization.ID || item.Category != "Regulatory" || item.Requirement != "Monthly custody report" {
		t.Fatalf("unexpected compliance item: %+v", item)
	}
	if item.Status != models.ComplianceStatusPending || item.DueDate == nil {
		t.Fatalf("status=%q due_date=%v, want pending with due date", item.Status, item.DueDate)
	}

	var audit models.StakeholderAuditLog
	if err := database.First(&audit, "entity_id = ?", item.ID).Error; err != nil {
		t.Fatalf("load compliance audit: %v", err)
	}
	if audit.Action != models.ActionComplianceCreate || audit.ActorMemberID != actor.MemberID {
		t.Fatalf("unexpected audit: %+v", audit)
	}

	_, err = service.CreateCompliance(context.Background(), actor, request)
	if err == nil || ErrorStatus(err) != 409 {
		t.Fatalf("duplicate error=%v status=%d, want 409", err, ErrorStatus(err))
	}
}

func TestCreateComplianceRejectsNonCustodianOrganization(t *testing.T) {
	database := testDB(t, &coreModels.Organization{}, &models.ComplianceItem{})
	stakeholderType := models.StakeholderTypeAssetManager
	stakeholderID := uint64(20)
	organization := coreModels.Organization{
		ID: "manager-org", Name: "Manager", Email: "manager@example.com",
		Type: coreModels.OrganizationTypeAssetManager, Status: coreModels.OrganizationStatusActive,
		CreatedBy: "admin", StakeholderID: &stakeholderID, StakeholderType: &stakeholderType,
	}
	if err := database.Create(&organization).Error; err != nil {
		t.Fatalf("seed manager organization: %v", err)
	}

	service := NewCustodianService(database, nil, nil, nil)
	_, err := service.CreateCompliance(context.Background(), AuthContext{}, models.CreateComplianceItemRequest{
		CustodianOrgID: organization.ID, Category: "regulatory", Requirement: "Monthly report",
	})
	if err == nil || ErrorStatus(err) != 422 {
		t.Fatalf("error=%v status=%d, want 422", err, ErrorStatus(err))
	}
}

func TestListComplianceForAdminReturnsAllAndAppliesFilters(t *testing.T) {
	database := testDB(t, &models.ComplianceItem{}, &models.StakeholderDocument{}, &models.ComplianceItemDocument{})
	now := time.Now()
	items := []models.ComplianceItem{
		{
			ID: "compliance-a", OrgID: "custodian-a", AssetID: "asset-1", Category: "regulatory",
			Requirement: "Monthly report", Status: models.ComplianceStatusPending, CreatedAt: now.Add(-3 * time.Hour), UpdatedAt: now,
		},
		{
			ID: "compliance-b", OrgID: "custodian-a", AssetID: "asset-2", Category: "custody",
			Requirement: "Proof of custody", Status: models.ComplianceStatusComplete, CreatedAt: now.Add(-2 * time.Hour), UpdatedAt: now,
		},
		{
			ID: "compliance-c", OrgID: "custodian-b", AssetID: "asset-1", Category: "regulatory",
			Requirement: "Annual filing", Status: models.ComplianceStatusPending, CreatedAt: now.Add(-time.Hour), UpdatedAt: now,
		},
	}
	if err := database.Create(&items).Error; err != nil {
		t.Fatalf("seed compliance items: %v", err)
	}

	document := models.StakeholderDocument{
		ID: uuid.NewString(), AssetID: "asset-1", UploadedByOrgID: "custodian-a", UploadedByMemberID: "member-1",
		Category: models.DocumentCategoryCustodyCompliance, Title: "Monthly report evidence", StorageObjectID: "stakeholder/report.pdf",
		AccessRoles: models.JSONStringArray{models.DashboardRoleAssetCustodian}, Status: models.DocumentStatusActive,
	}
	if err := database.Create(&document).Error; err != nil {
		t.Fatalf("seed compliance document: %v", err)
	}
	if err := database.Create(&models.ComplianceItemDocument{
		ComplianceItemID: items[0].ID, DocumentID: document.ID, CreatedAt: now,
	}).Error; err != nil {
		t.Fatalf("link compliance document: %v", err)
	}

	service := NewCustodianService(database, nil, nil, nil)
	page, total, err := service.ListComplianceForAdmin(context.Background(), 1, 2, AdminComplianceListFilters{})
	if err != nil {
		t.Fatalf("list all compliance items: %v", err)
	}
	if total != 3 || len(page) != 2 {
		t.Fatalf("total=%d records=%d, want total=3 records=2", total, len(page))
	}

	filtered, total, err := service.ListComplianceForAdmin(context.Background(), 1, 20, AdminComplianceListFilters{
		CustodianOrgID: " custodian-a ", AssetID: " asset-1 ", Status: " PENDING ",
	})
	if err != nil {
		t.Fatalf("list filtered compliance items: %v", err)
	}
	if total != 1 || len(filtered) != 1 || filtered[0].ID != items[0].ID {
		t.Fatalf("filtered items=%+v total=%d, want compliance-a only", filtered, total)
	}
	if len(filtered[0].Documents) != 1 || filtered[0].Documents[0].ID != document.ID {
		t.Fatalf("admin compliance documents=%+v, want %s", filtered[0].Documents, document.ID)
	}

	_, _, err = service.ListComplianceForAdmin(context.Background(), 1, 20, AdminComplianceListFilters{Status: "unknown"})
	if err == nil || ErrorStatus(err) != 400 {
		t.Fatalf("invalid status error=%v status=%d, want 400", err, ErrorStatus(err))
	}
}

func TestUpdateComplianceClearsCompletionMetadataWhenReopened(t *testing.T) {
	database := testDB(t, &models.ComplianceItem{}, &models.StakeholderDocument{}, &models.ComplianceItemDocument{}, &models.StakeholderAuditLog{})
	completedAt := time.Now().Add(-time.Hour)
	completedBy := "custodian-member"
	item := models.ComplianceItem{
		ID: "compliance-1", OrgID: "custodian-org", Category: "regulatory",
		Requirement: "Monthly report", Status: models.ComplianceStatusComplete,
		CompletedAt: &completedAt, CompletedByMemberID: &completedBy,
	}
	if err := database.Create(&item).Error; err != nil {
		t.Fatalf("seed compliance item: %v", err)
	}

	service := NewCustodianService(database, nil, NewDocumentService(database, nil, nil, nil), NewAuditService(database))
	updated, err := service.UpdateCompliance(context.Background(), AuthContext{
		MemberID: completedBy, OrganizationID: item.OrgID, DashboardRole: models.DashboardRoleAssetCustodian,
	}, item.ID, models.UpdateComplianceItemRequest{Status: models.ComplianceStatusPending})
	if err != nil {
		t.Fatalf("reopen compliance item: %v", err)
	}
	if updated.CompletedAt != nil || updated.CompletedByMemberID != nil {
		t.Fatalf("completion metadata should be cleared: %+v", updated)
	}
}

func TestUpdateComplianceLinksManagedDocuments(t *testing.T) {
	database := testDB(t, &models.ComplianceItem{}, &models.StakeholderDocument{}, &models.ComplianceItemDocument{})
	item := models.ComplianceItem{
		ID: "compliance-with-evidence", OrgID: "custodian-org", AssetID: "asset-1",
		Category: "regulatory", Requirement: "Monthly report", Status: models.ComplianceStatusPending,
	}
	if err := database.Create(&item).Error; err != nil {
		t.Fatalf("seed compliance item: %v", err)
	}
	document := models.StakeholderDocument{
		ID: uuid.NewString(), AssetID: item.AssetID, UploadedByOrgID: item.OrgID, UploadedByMemberID: "member-1",
		Category: models.DocumentCategoryCustodyCompliance, Title: "Custody report", StorageObjectID: "stakeholder/doc.pdf",
		AccessRoles: models.JSONStringArray{models.DashboardRoleAssetCustodian, models.DashboardRoleTrustee}, Status: models.DocumentStatusActive,
	}
	if err := database.Create(&document).Error; err != nil {
		t.Fatalf("seed document: %v", err)
	}
	documentIDs := []string{document.ID}
	auth := AuthContext{MemberID: "member-1", OrganizationID: item.OrgID, DashboardRole: models.DashboardRoleAssetCustodian}
	service := NewCustodianService(database, nil, NewDocumentService(database, nil, nil, nil), nil)

	updated, err := service.UpdateCompliance(context.Background(), auth, item.ID, models.UpdateComplianceItemRequest{
		Status: models.ComplianceStatusComplete, DocumentIDs: &documentIDs,
	})
	if err != nil {
		t.Fatalf("update compliance evidence: %v", err)
	}
	if len(updated.Documents) != 1 || updated.Documents[0].ID != document.ID {
		t.Fatalf("updated documents = %+v, want %s", updated.Documents, document.ID)
	}
	records, _, err := service.ListCompliance(context.Background(), auth, 1, 20)
	if err != nil {
		t.Fatalf("list compliance: %v", err)
	}
	if len(records) != 1 || len(records[0].Documents) != 1 || records[0].Documents[0].ID != document.ID {
		t.Fatalf("listed compliance evidence = %+v", records)
	}
}

func TestCustodianCreatesAndUpdatesAssignedAssetAccount(t *testing.T) {
	database := testDB(t, &models.StakeholderAssetAssignment{}, &models.SegregatedAccount{}, &models.StakeholderAuditLog{}, &models.StakeholderDocument{}, &coreModels.Organization{})
	custodianOrgID := "custodian-org"
	custodianStakeholderID := uint64(17)
	if err := database.Create(&models.StakeholderAssetAssignment{
		ID: "assignment-1", AssetID: "asset-1", AssetCode: "ASSET1", CustodianOrgID: &custodianOrgID,
		CustodianStakeholderID: &custodianStakeholderID, Status: models.AssignmentStatusActive,
	}).Error; err != nil {
		t.Fatalf("seed assignment: %v", err)
	}
	assets := NewAssetService(database, fakeTokenizationReadClient{assets: []stakeholderDB.TokenizedAsset{
		{ID: "asset-1", AssetCode: "ASSET1", ApprovedAssetCustodianID: custodianStakeholderID},
	}}, nil)
	service := NewCustodianService(database, assets, nil, NewAuditService(database))
	auth := AuthContext{
		MemberID: "custodian-member", OrganizationID: custodianOrgID,
		StakeholderID: custodianStakeholderID, DashboardRole: models.DashboardRoleAssetCustodian,
	}

	account, err := service.CreateAccount(context.Background(), auth, models.CreateSegregatedAccountRequest{
		AssetID: "asset-1", AccountType: models.AccountTypeRevenue, AccountName: "Revenue account",
		Balance: "0", Currency: "ngn", BankDetails: models.JSONMap{"bank": "Trovo Bank"},
	})
	if err != nil {
		t.Fatalf("create account: %v", err)
	}
	if account.AssetCode != "ASSET1" || account.Currency != "NGN" || !account.Balance.IsZero() || account.Status != models.AccountStatusActive {
		t.Fatalf("created account = %+v", account)
	}

	updated, err := service.UpdateAccount(context.Background(), auth, account.ID, models.UpdateSegregatedAccountRequest{
		Balance: "1250.50", Status: models.AccountStatusInactive,
	})
	if err != nil {
		t.Fatalf("update account: %v", err)
	}
	if updated.Balance.String() != "1250.5" || updated.Status != models.AccountStatusInactive {
		t.Fatalf("updated account = %+v", updated)
	}
	if _, err := service.CreateAccount(context.Background(), auth, models.CreateSegregatedAccountRequest{
		AssetID: "unassigned", AccountType: models.AccountTypeRevenue, AccountName: "Invalid", Currency: "NGN",
	}); ErrorStatus(err) != 404 {
		t.Fatalf("unassigned asset status = %d, want 404 (err=%v)", ErrorStatus(err), err)
	}
}

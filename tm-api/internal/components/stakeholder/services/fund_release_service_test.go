package services

import (
	"context"
	"testing"
	"time"

	"admin-panel-dashboard/internal/components/stakeholder/models"
	coreModels "admin-panel-dashboard/internal/models"

	"github.com/shopspring/decimal"
)

func currencyAmount(t *testing.T, totals []models.CurrencyTotal, currency string) string {
	t.Helper()
	for _, total := range totals {
		if total.Currency == currency {
			return total.Amount
		}
	}
	t.Fatalf("currency %s not found in %+v", currency, totals)
	return ""
}

func TestCanTransitionFundRelease(t *testing.T) {
	valid := [][2]string{
		{models.FundReleaseStatusDraft, models.FundReleaseStatusSubmitted},
		{models.FundReleaseStatusSubmitted, models.FundReleaseStatusExecutionPending},
		{models.FundReleaseStatusExecutionPending, models.FundReleaseStatusProcessing},
		{models.FundReleaseStatusProcessing, models.FundReleaseStatusCompleted},
	}
	for _, transition := range valid {
		if !CanTransitionFundRelease(transition[0], transition[1]) {
			t.Fatalf("expected %s -> %s to be valid", transition[0], transition[1])
		}
	}

	invalid := [][2]string{
		{models.FundReleaseStatusCompleted, models.FundReleaseStatusProcessing},
		{models.FundReleaseStatusTrusteeRejected, models.FundReleaseStatusExecutionPending},
		{models.FundReleaseStatusSubmitted, models.FundReleaseStatusCompleted},
	}
	for _, transition := range invalid {
		if CanTransitionFundRelease(transition[0], transition[1]) {
			t.Fatalf("expected %s -> %s to be invalid", transition[0], transition[1])
		}
	}
}

func TestUpdateExecutionStatusRequiresExecuteStepUpBeforeLeavingExecutionPending(t *testing.T) {
	db := testDB(t, &models.FundReleaseRequest{})
	now := time.Now()
	request := models.FundReleaseRequest{
		ID:                "fr-execution-pending",
		AssetID:           "asset-1",
		RequesterOrgID:    "manager-org",
		RequesterMemberID: "manager-member",
		TrusteeOrgID:      "trustee-org",
		CustodianOrgID:    "custodian-org",
		Amount:            decimal.NewFromInt(1000),
		Currency:          "NGN",
		Purpose:           "regression test",
		Status:            models.FundReleaseStatusExecutionPending,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	if err := db.Create(&request).Error; err != nil {
		t.Fatalf("seed fund release: %v", err)
	}

	service := NewFundReleaseService(db, nil, nil, nil, nil, nil)
	_, err := service.UpdateExecutionStatus(context.Background(), AuthContext{
		MemberID:       "custodian-member",
		OrganizationID: "custodian-org",
		DashboardRole:  models.DashboardRoleAssetCustodian,
	}, request.ID, models.UpdateFundReleaseStatusRequest{
		Status:             models.FundReleaseStatusCompleted,
		ExecutionReference: "bypass-test",
	})
	if err == nil {
		t.Fatal("expected status endpoint to reject execution without step-up")
	}
	if ErrorStatus(err) != 403 {
		t.Fatalf("status = %d, want 403", ErrorStatus(err))
	}

	var stored models.FundReleaseRequest
	if err := db.First(&stored, "id = ?", request.ID).Error; err != nil {
		t.Fatalf("load fund release: %v", err)
	}
	if stored.Status != models.FundReleaseStatusExecutionPending {
		t.Fatalf("status = %q, want %q", stored.Status, models.FundReleaseStatusExecutionPending)
	}
	if stored.ExecutedAt != nil || stored.ExecutedByMemberID != nil {
		t.Fatal("execution fields should not be set by rejected status update")
	}
}

func TestUpdateExecutionStatusAllowsPostExecutionProgression(t *testing.T) {
	db := testDB(t, &models.FundReleaseRequest{})
	now := time.Now()
	executedBy := "custodian-member"
	request := models.FundReleaseRequest{
		ID:                 "fr-processing",
		AssetID:            "asset-1",
		RequesterOrgID:     "manager-org",
		RequesterMemberID:  "manager-member",
		TrusteeOrgID:       "trustee-org",
		CustodianOrgID:     "custodian-org",
		Amount:             decimal.NewFromInt(1000),
		Currency:           "NGN",
		Purpose:            "regression test",
		Status:             models.FundReleaseStatusProcessing,
		ExecutedByMemberID: &executedBy,
		ExecutedAt:         &now,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	if err := db.Create(&request).Error; err != nil {
		t.Fatalf("seed fund release: %v", err)
	}

	service := NewFundReleaseService(db, nil, nil, nil, nil, nil)
	updated, err := service.UpdateExecutionStatus(context.Background(), AuthContext{
		MemberID:       "custodian-member",
		OrganizationID: "custodian-org",
		DashboardRole:  models.DashboardRoleAssetCustodian,
	}, request.ID, models.UpdateFundReleaseStatusRequest{
		Status:             models.FundReleaseStatusCompleted,
		ExecutionReference: "staging-exec-001",
	})
	if err != nil {
		t.Fatalf("update execution status: %v", err)
	}
	if updated.Status != models.FundReleaseStatusCompleted {
		t.Fatalf("status = %q, want %q", updated.Status, models.FundReleaseStatusCompleted)
	}
}

func TestFundReleaseDetailsIncludeRequesterDocumentsAndReceivingAccount(t *testing.T) {
	db := testDB(t, &models.FundReleaseRequest{}, &models.StakeholderDocument{}, &coreModels.Organization{}, &coreModels.OrganizationMember{})
	now := time.Now()
	if err := db.Create(&coreModels.Organization{
		ID: "manager-org", Name: "Green Asset Managers", Email: "org@example.com", Type: "CORPORATE",
		Status: coreModels.OrganizationStatusActive, CreatedBy: "test",
	}).Error; err != nil {
		t.Fatalf("seed organization: %v", err)
	}
	if err := db.Create(&coreModels.OrganizationMember{
		ID: "manager-member", OrganizationID: "manager-org", Email: "requester@example.com", Role: "MEMBER",
		Status: "ACTIVE", CreatedBy: "test", FirstName: "Ada", LastName: "Okafor",
	}).Error; err != nil {
		t.Fatalf("seed member: %v", err)
	}
	document := models.StakeholderDocument{
		ID: "fund-doc", AssetID: "asset-1", UploadedByOrgID: "manager-org", UploadedByMemberID: "manager-member",
		Category: models.DocumentCategoryFundReleaseSupporting, Title: "Invoice", StorageObjectID: "object-1",
		OriginalFilename: "invoice.pdf", MimeType: "application/pdf", SizeBytes: 512,
		AccessRoles: models.JSONStringArray{models.DashboardRoleAssetManager, models.DashboardRoleTrustee},
		Status:      models.DocumentStatusActive, CreatedAt: now, UpdatedAt: now,
	}
	if err := db.Create(&document).Error; err != nil {
		t.Fatalf("seed document: %v", err)
	}
	request := models.FundReleaseRequest{
		ID: "fr-details", AssetID: "asset-1", RequesterOrgID: "manager-org", RequesterMemberID: "manager-member",
		TrusteeOrgID: "trustee-org", Amount: decimal.NewFromInt(100), Currency: "NGN", Purpose: "Equipment",
		Status: models.FundReleaseStatusSubmitted, SupportingDocumentIDs: models.JSONStringArray{document.ID},
		ReceivingBank: "Access Bank", ReceivingAccountName: "Green Harvest Ltd", ReceivingAccountNumber: "0123456789",
		CreatedAt: now, UpdatedAt: now,
	}
	if err := db.Create(&request).Error; err != nil {
		t.Fatalf("seed fund release: %v", err)
	}

	service := NewFundReleaseService(db, nil, nil, nil, nil, nil)
	detail, err := service.GetDetails(context.Background(), AuthContext{
		OrganizationID: "manager-org", DashboardRole: models.DashboardRoleAssetManager,
	}, request.ID)
	if err != nil {
		t.Fatalf("get details: %v", err)
	}
	if detail.Requester.FirstName != "Ada" || detail.Requester.OrganizationName != "Green Asset Managers" {
		t.Fatalf("requester = %+v", detail.Requester)
	}
	if len(detail.SupportingDocuments) != 1 || detail.SupportingDocuments[0].DownloadPath != "/api/v1/stakeholder/shared/documents/fund-doc/download" {
		t.Fatalf("documents = %+v", detail.SupportingDocuments)
	}
	if detail.ReceivingAccount.Bank != "Access Bank" || detail.ReceivingAccount.AccountNumber != "0123456789" {
		t.Fatalf("receiving account = %+v", detail.ReceivingAccount)
	}
}

func TestFundReleaseListIncludesRequesterAndReviewerNames(t *testing.T) {
	db := testDB(t, &models.FundReleaseRequest{}, &coreModels.Organization{}, &coreModels.OrganizationMember{})
	manager := coreModels.Organization{
		ID: "manager-org", Name: "Green Asset Managers", Email: "manager@example.com", Type: "CORPORATE",
		Status: coreModels.OrganizationStatusActive, CreatedBy: "test",
	}
	trustee := coreModels.Organization{
		ID: "trustee-org", Name: "Trusted Trustees", Email: "trustee@example.com", Type: "CORPORATE",
		Status: coreModels.OrganizationStatusActive, CreatedBy: "test",
	}
	requester := coreModels.OrganizationMember{
		ID: "requester-member", OrganizationID: manager.ID, Email: "ada@example.com", FirstName: "Ada", LastName: "Okafor",
		Role: string(coreModels.OrganizationMemberRole), Status: coreModels.OrganizationStatusActive, CreatedBy: "test",
	}
	reviewer := coreModels.OrganizationMember{
		ID: "reviewer-member", OrganizationID: trustee.ID, Email: "nancy@example.com", FirstName: "Nancy", LastName: "Rike",
		Role: string(coreModels.OrganizationMemberRole), Status: coreModels.OrganizationStatusActive, CreatedBy: "test",
	}
	for _, record := range []interface{}{&manager, &trustee, &requester, &reviewer} {
		if err := db.Create(record).Error; err != nil {
			t.Fatalf("seed %T: %v", record, err)
		}
	}
	reviewerID := reviewer.ID
	request := models.FundReleaseRequest{
		ID: "fund-list-1", AssetID: "asset-1", RequesterOrgID: manager.ID, RequesterMemberID: requester.ID,
		TrusteeOrgID: trustee.ID, Amount: decimal.NewFromInt(100), Currency: "CNGN", Purpose: "Milestone",
		Status: models.FundReleaseStatusExecutionPending, ReviewedByMemberID: &reviewerID,
	}
	if err := db.Create(&request).Error; err != nil {
		t.Fatal(err)
	}

	service := NewFundReleaseService(db, nil, nil, nil, nil, nil)
	records, total, err := service.List(context.Background(), AuthContext{OrganizationID: trustee.ID, DashboardRole: models.DashboardRoleTrustee}, 1, 20, "", "")
	if err != nil {
		t.Fatalf("list fund releases: %v", err)
	}
	if total != 1 || len(records) != 1 {
		t.Fatalf("records=%d total=%d, want 1", len(records), total)
	}
	if records[0].Requester.Name != "Ada Okafor" || records[0].Requester.OrganizationName != manager.Name {
		t.Fatalf("requester=%+v", records[0].Requester)
	}
	if records[0].Reviewer == nil || records[0].Reviewer.Name != "Nancy Rike" || records[0].Reviewer.OrganizationName != trustee.Name {
		t.Fatalf("reviewer=%+v", records[0].Reviewer)
	}
	if records[0].RequesterMemberID != requester.ID || records[0].ReviewedByMemberID == nil || *records[0].ReviewedByMemberID != reviewer.ID {
		t.Fatalf("original identifiers were not preserved: %+v", records[0])
	}
}

func TestAssetManagerFundSummaryUsesExplicitStatusBuckets(t *testing.T) {
	db := testDB(t, &models.FundReleaseRequest{})
	now := time.Now()
	requests := []models.FundReleaseRequest{
		{ID: "fr-submitted", AssetID: "asset-1", RequesterOrgID: "manager-org", RequesterMemberID: "m", TrusteeOrgID: "t", Amount: decimal.NewFromInt(10), Currency: "NGN", Purpose: "one", Status: models.FundReleaseStatusSubmitted, CreatedAt: now, UpdatedAt: now},
		{ID: "fr-approved", AssetID: "asset-1", RequesterOrgID: "manager-org", RequesterMemberID: "m", TrusteeOrgID: "t", Amount: decimal.NewFromInt(50), Currency: "NGN", Purpose: "two", Status: models.FundReleaseStatusExecutionPending, CreatedAt: now, UpdatedAt: now},
		{ID: "fr-completed", AssetID: "asset-1", RequesterOrgID: "manager-org", RequesterMemberID: "m", TrusteeOrgID: "t", Amount: decimal.NewFromInt(20), Currency: "NGN", Purpose: "three", Status: models.FundReleaseStatusCompleted, CreatedAt: now, UpdatedAt: now},
		{ID: "fr-rejected", AssetID: "asset-1", RequesterOrgID: "manager-org", RequesterMemberID: "m", TrusteeOrgID: "t", Amount: decimal.NewFromInt(30), Currency: "NGN", Purpose: "four", Status: models.FundReleaseStatusTrusteeRejected, CreatedAt: now, UpdatedAt: now},
		{ID: "fr-other-org", AssetID: "asset-1", RequesterOrgID: "other-org", RequesterMemberID: "m", TrusteeOrgID: "t", Amount: decimal.NewFromInt(999), Currency: "NGN", Purpose: "hidden", Status: models.FundReleaseStatusCompleted, CreatedAt: now, UpdatedAt: now},
	}
	if err := db.Create(&requests).Error; err != nil {
		t.Fatalf("seed fund releases: %v", err)
	}
	service := NewFundReleaseService(db, nil, nil, nil, nil, nil)
	summary, err := service.AssetManagerSummary(context.Background(), AuthContext{
		OrganizationID: "manager-org", DashboardRole: models.DashboardRoleAssetManager,
	}, "")
	if err != nil {
		t.Fatalf("summary: %v", err)
	}
	if summary.Requested.Count != 4 || currencyAmount(t, summary.Requested.Amounts, "NGN") != "110" {
		t.Fatalf("requested = %+v", summary.Requested)
	}
	if summary.Approved.Count != 2 || currencyAmount(t, summary.Approved.Amounts, "NGN") != "70" {
		t.Fatalf("approved = %+v", summary.Approved)
	}
	if summary.Pending.Count != 2 || currencyAmount(t, summary.Pending.Amounts, "NGN") != "60" {
		t.Fatalf("pending = %+v", summary.Pending)
	}
	if summary.Rejected.Count != 1 || currencyAmount(t, summary.Rejected.Amounts, "NGN") != "30" {
		t.Fatalf("rejected = %+v", summary.Rejected)
	}
	if currencyAmount(t, summary.Remaining, "NGN") != "50" {
		t.Fatalf("remaining = %+v", summary.Remaining)
	}
}

func TestReceivingAccountFieldsMustBeComplete(t *testing.T) {
	_, _, _, err := validateReceivingAccount(models.CreateFundReleaseRequest{ReceivingBank: "Bank only"})
	if ErrorStatus(err) != 400 {
		t.Fatalf("status = %d, want 400", ErrorStatus(err))
	}
	bank, name, number, err := validateReceivingAccount(models.CreateFundReleaseRequest{
		ReceivingBank: " Bank ", ReceivingAccountName: " Account Name ", ReceivingAccountNumber: " 1234 ",
	})
	if err != nil || bank != "Bank" || name != "Account Name" || number != "1234" {
		t.Fatalf("validated account = %q %q %q err=%v", bank, name, number, err)
	}
}

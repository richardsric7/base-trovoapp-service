package services

import (
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"admin-panel-dashboard/internal/components/stakeholder/models"
	coreModels "admin-panel-dashboard/internal/models"

	"github.com/google/uuid"
)

func boolPointer(value bool) *bool { return &value }

func complianceRequirementTestDB(t *testing.T) *ComplianceRequirementService {
	t.Helper()
	database := testDB(t,
		&coreModels.Organization{},
		&models.StakeholderDocument{},
		&models.StakeholderAuditLog{},
		&models.ComplianceRequirementTemplate{},
		&models.ComplianceRequirementTemplateItem{},
		&models.ComplianceRequirementInstance{},
	)
	documents := NewDocumentService(database, nil, nil, nil)
	return NewComplianceRequirementService(database, NewAuditService(database), documents)
}

func TestComplianceTemplatePreservesExplicitOptionalAndDefaultsRequired(t *testing.T) {
	service := complianceRequirementTestDB(t)
	ctx := context.Background()

	optionalTemplate, err := service.CreateTemplate(ctx, AuthContext{MemberID: "admin-1"}, models.CreateComplianceTemplateRequest{
		OrgType: coreModels.StakeholderTypeAssetManager,
		Level:   1,
		Items: []models.ComplianceTemplateItemInput{{
			Name: "Optional note", Category: "kyb", InputType: models.ComplianceInputTypeText, Required: boolPointer(false),
		}},
	})
	if err != nil {
		t.Fatalf("create optional template: %v", err)
	}
	if optionalTemplate.Items[0].Required == nil || *optionalTemplate.Items[0].Required {
		t.Fatal("explicit required=false was changed to true")
	}
	var storedOptional models.ComplianceRequirementTemplateItem
	if err := service.db.First(&storedOptional, "id = ?", optionalTemplate.Items[0].ID).Error; err != nil {
		t.Fatalf("load optional item: %v", err)
	}
	if storedOptional.Required == nil || *storedOptional.Required {
		t.Fatal("stored optional item is required")
	}

	defaultTemplate, err := service.CreateTemplate(ctx, AuthContext{MemberID: "admin-1"}, models.CreateComplianceTemplateRequest{
		OrgType: coreModels.StakeholderTypeAssetManager,
		Level:   2,
		Items: []models.ComplianceTemplateItemInput{{
			Name: "Registration certificate", Category: "kyb", InputType: models.ComplianceInputTypeDocumentUpload,
		}},
	})
	if err != nil {
		t.Fatalf("create default-required template: %v", err)
	}
	if defaultTemplate.Items[0].Required == nil || !*defaultTemplate.Items[0].Required {
		t.Fatal("omitted required should default to true")
	}
}

func TestAutoAssignComplianceRequirementsIsIdempotentAndSnapshotsDefinition(t *testing.T) {
	service := complianceRequirementTestDB(t)
	ctx := context.Background()
	template, err := service.CreateTemplate(ctx, AuthContext{MemberID: "admin-1"}, models.CreateComplianceTemplateRequest{
		OrgType: coreModels.StakeholderTypeTrustee,
		Level:   1,
		Items: []models.ComplianceTemplateItemInput{{
			Name: "Trust deed", Category: "kyb", Description: "Upload the executed deed",
			InputType: models.ComplianceInputTypeDocumentUpload, Required: boolPointer(false),
		}},
	})
	if err != nil {
		t.Fatalf("create template: %v", err)
	}
	organization := coreModels.Organization{
		ID: "org-1", Name: "Trustee", Email: "trustee@example.com", Type: coreModels.StakeholderTypeTrustee,
		Status: coreModels.OrganizationStatusActive, CreatedBy: "admin",
	}
	if err := service.db.Create(&organization).Error; err != nil {
		t.Fatalf("create organization: %v", err)
	}

	for i := 0; i < 2; i++ {
		if err := AutoAssignComplianceRequirements(ctx, service.db, "org-1", coreModels.StakeholderTypeTrustee, 1); err != nil {
			t.Fatalf("auto assign pass %d: %v", i+1, err)
		}
	}
	var instances []models.ComplianceRequirementInstance
	if err := service.db.Where("org_id = ?", "org-1").Find(&instances).Error; err != nil {
		t.Fatalf("list instances: %v", err)
	}
	if len(instances) != 1 {
		t.Fatalf("instances=%d, want 1", len(instances))
	}
	instance := instances[0]
	if instance.TemplateItemID == nil || *instance.TemplateItemID != template.Items[0].ID ||
		instance.Description != "Upload the executed deed" || instance.Required == nil || *instance.Required {
		t.Fatalf("definition snapshot is incomplete: %+v", instance)
	}

	duplicate := instance
	duplicate.ID = uuid.NewString()
	if err := service.db.Create(&duplicate).Error; err == nil {
		t.Fatal("database accepted a duplicate org/template-item assignment")
	}
}

func TestUpdateTemplatePreservesItemIdentityAndDoesNotReassignCompletedRequirement(t *testing.T) {
	service := complianceRequirementTestDB(t)
	ctx := context.Background()
	template, err := service.CreateTemplate(ctx, AuthContext{MemberID: "admin-1"}, models.CreateComplianceTemplateRequest{
		OrgType: coreModels.StakeholderTypeTrustee,
		Level:   1,
		Items: []models.ComplianceTemplateItemInput{{
			Name: "Trust deed", Category: "kyb", InputType: models.ComplianceInputTypeDocumentUpload,
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	organization := coreModels.Organization{
		ID: "org-versioning", Name: "Trustee", Email: "versioning@example.com", Type: coreModels.StakeholderTypeTrustee,
		Status: coreModels.OrganizationStatusActive, CreatedBy: "admin",
	}
	if err := service.db.Create(&organization).Error; err != nil {
		t.Fatal(err)
	}
	if err := AutoAssignComplianceRequirements(ctx, service.db, organization.ID, coreModels.StakeholderTypeTrustee, 1); err != nil {
		t.Fatal(err)
	}

	originalItemID := template.Items[0].ID
	updated, err := service.UpdateTemplate(ctx, AuthContext{MemberID: "admin-1"}, template.ID, models.UpdateComplianceTemplateRequest{
		Items: []models.ComplianceTemplateItemInput{
			{ID: originalItemID, Name: "Executed trust deed", Category: "kyb", InputType: models.ComplianceInputTypeDocumentUpload},
			{Name: "Trustee licence", Category: "regulatory", InputType: models.ComplianceInputTypeDocumentUpload},
		},
	})
	if err != nil {
		t.Fatalf("update template: %v", err)
	}
	if updated.Items[0].ID != originalItemID {
		t.Fatalf("updated item id=%q, want stable id %q", updated.Items[0].ID, originalItemID)
	}
	if err := AutoAssignComplianceRequirements(ctx, service.db, organization.ID, coreModels.StakeholderTypeTrustee, 1); err != nil {
		t.Fatal(err)
	}
	var instances []models.ComplianceRequirementInstance
	if err := service.db.Where("org_id = ?", organization.ID).Order("created_at asc").Find(&instances).Error; err != nil {
		t.Fatal(err)
	}
	if len(instances) != 2 {
		t.Fatalf("instances=%d, want original plus one genuinely new requirement", len(instances))
	}
	if instances[0].Requirement != "Trust deed" || instances[0].TemplateVersion == nil || *instances[0].TemplateVersion != 1 {
		t.Fatalf("existing assignment snapshot changed: %+v", instances[0])
	}
}

func TestRequirementDetailIsReadOnlyAndReviewTransitionIsExplicit(t *testing.T) {
	service := complianceRequirementTestDB(t)
	ctx := context.Background()
	organization := coreModels.Organization{
		ID: "org-1", Name: "Org", Email: "review-org@example.com", Type: coreModels.StakeholderTypeAssetManager,
		Status: coreModels.OrganizationStatusActive, CreatedBy: "admin",
	}
	if err := service.db.Create(&organization).Error; err != nil {
		t.Fatal(err)
	}
	submitter := "member-1"
	instance := models.ComplianceRequirementInstance{
		ID: uuid.NewString(), OrgID: "org-1", Category: "kyb", Requirement: "Certificate",
		Required: boolPointer(true), InputType: models.ComplianceInputTypeText,
		Status: models.ComplianceRequirementStatusSubmitted, SubmittedByMemberID: &submitter,
		SubmissionData: models.JSONMap{"text": "submitted"}, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	if err := service.db.Create(&instance).Error; err != nil {
		t.Fatalf("seed instance: %v", err)
	}

	detail, err := service.GetRequirementDetail(ctx, instance.ID)
	if err != nil {
		t.Fatalf("get detail: %v", err)
	}
	if detail.Status != models.ComplianceRequirementStatusSubmitted {
		t.Fatalf("GET changed status to %q", detail.Status)
	}

	actor := AuthContext{MemberID: "admin-1", DashboardRole: "ROOT_SUPER_ADMIN"}
	started, err := service.StartReview(ctx, actor, instance.ID)
	if err != nil || started.Status != models.ComplianceRequirementStatusUnderReview {
		t.Fatalf("start review: instance=%+v err=%v", started, err)
	}
	if _, err := service.StartReview(ctx, actor, instance.ID); err != nil {
		t.Fatalf("repeat start review should be idempotent: %v", err)
	}

	approved, err := service.ReviewRequirement(ctx, actor, instance.ID, models.ReviewComplianceRequirementRequest{Decision: models.ComplianceRequirementStatusApproved})
	if err != nil {
		t.Fatalf("approve requirement: %v", err)
	}
	if approved.CompletedByMemberID == nil || *approved.CompletedByMemberID != submitter {
		t.Fatalf("approved requirement lost submitter: %+v", approved)
	}
}

func TestSubmitComplianceRequirementValidatesInputTypeAndDocumentOwnership(t *testing.T) {
	service := complianceRequirementTestDB(t)
	ctx := context.Background()
	auth := AuthContext{OrganizationID: "org-1", MemberID: "member-1", MemberRole: "ROOT_SUPER_ADMIN"}

	textRequirement := models.ComplianceRequirementInstance{
		ID: uuid.NewString(), OrgID: auth.OrganizationID, Category: "kyb", Requirement: "Explain ownership",
		Required: boolPointer(true), InputType: models.ComplianceInputTypeText, Status: models.ComplianceRequirementStatusPending,
	}
	if err := service.db.Create(&textRequirement).Error; err != nil {
		t.Fatal(err)
	}
	if _, err := service.SubmitOrgRequirement(ctx, auth, textRequirement.ID, models.SubmitComplianceRequirementRequest{SubmissionData: map[string]interface{}{"value": "wrong key"}}); ErrorStatus(err) != http.StatusBadRequest {
		t.Fatalf("invalid text submission status=%d err=%v", ErrorStatus(err), err)
	}
	if _, err := service.SubmitOrgRequirement(ctx, auth, textRequirement.ID, models.SubmitComplianceRequirementRequest{SubmissionData: map[string]interface{}{"text": "Beneficial owners listed"}}); err != nil {
		t.Fatalf("valid text submission: %v", err)
	}
	var submitAudit models.StakeholderAuditLog
	if err := service.db.Where("entity_id = ? AND action = ?", textRequirement.ID, models.ActionComplianceRequirementSubmit).First(&submitAudit).Error; err != nil {
		t.Fatalf("load submission audit: %v", err)
	}
	if submitAudit.ActorRole != auth.MemberRole {
		t.Fatalf("submission audit role=%q, want %q", submitAudit.ActorRole, auth.MemberRole)
	}

	document := models.StakeholderDocument{
		ID: uuid.NewString(), UploadedByOrgID: auth.OrganizationID, UploadedByMemberID: auth.MemberID,
		Category: models.DocumentCategoryComplianceRequirement, Title: "Certificate", StorageObjectID: "private/certificate.pdf",
		Status: models.DocumentStatusActive,
	}
	if err := service.db.Create(&document).Error; err != nil {
		t.Fatal(err)
	}
	documentRequirement := models.ComplianceRequirementInstance{
		ID: uuid.NewString(), OrgID: auth.OrganizationID, Category: "kyb", Requirement: "Upload certificate",
		Required: boolPointer(true), InputType: models.ComplianceInputTypeDocumentUpload, Status: models.ComplianceRequirementStatusPending,
	}
	if err := service.db.Create(&documentRequirement).Error; err != nil {
		t.Fatal(err)
	}
	if _, err := service.SubmitOrgRequirement(ctx, AuthContext{OrganizationID: "other-org", MemberID: "member-2"}, documentRequirement.ID, models.SubmitComplianceRequirementRequest{SubmissionData: map[string]interface{}{"document_id": document.ID}}); ErrorStatus(err) != http.StatusNotFound {
		t.Fatalf("cross-org requirement access status=%d err=%v", ErrorStatus(err), err)
	}
	if _, err := service.SubmitOrgRequirement(ctx, auth, documentRequirement.ID, models.SubmitComplianceRequirementRequest{SubmissionData: map[string]interface{}{"document_id": document.ID}}); err != nil {
		t.Fatalf("valid document submission: %v", err)
	}
}

func TestOrganizationComplianceDocumentUploadAndOwnerDownload(t *testing.T) {
	service := complianceRequirementTestDB(t)
	storage := &fakeAdminDocumentStorage{}
	documents := NewDocumentService(service.db, nil, NewAuditService(service.db), storage)
	auth := AuthContext{OrganizationID: "legal-org", MemberID: "legal-member"}
	content := []byte("%PDF-1.4\nprivate compliance evidence")
	document, err := documents.Upload(context.Background(), auth, models.UploadDocumentRequest{
		Category: models.DocumentCategoryComplianceRequirement,
		Title:    "Incorporation certificate",
		Filename: "certificate.pdf",
		Content:  content,
	})
	if err != nil {
		t.Fatalf("upload organization compliance document: %v", err)
	}
	if document.UploadedByOrgID != auth.OrganizationID || document.AssetID != "" || document.Category != models.DocumentCategoryComplianceRequirement {
		t.Fatalf("unexpected document metadata: %+v", document)
	}
	if err := documents.ValidateForUse(context.Background(), auth, "", models.DocumentCategoryComplianceRequirement, []string{document.ID}); err != nil {
		t.Fatalf("validate uploaded compliance document: %v", err)
	}
	download, err := documents.DownloadOwnedByOrganization(context.Background(), auth.OrganizationID, document.ID, models.DocumentCategoryComplianceRequirement)
	if err != nil {
		t.Fatalf("download uploaded document: %v", err)
	}
	defer download.Body.Close()
	downloaded, err := io.ReadAll(download.Body)
	if err != nil {
		t.Fatalf("read downloaded document: %v", err)
	}
	if string(downloaded) != string(content) {
		t.Fatalf("downloaded content=%q, want %q", downloaded, content)
	}
	if _, err := documents.DownloadOwnedByOrganization(context.Background(), "other-org", document.ID, models.DocumentCategoryComplianceRequirement); ErrorStatus(err) != http.StatusNotFound {
		t.Fatalf("cross-org download status=%d err=%v", ErrorStatus(err), err)
	}
}

func TestCreateAdHocRequirementIsSeparateAndRejectsDuplicates(t *testing.T) {
	service := complianceRequirementTestDB(t)
	org := coreModels.Organization{
		ID: "org-1", Name: "Org", Email: "org@example.com", Type: coreModels.StakeholderTypeLegal,
		Status: coreModels.OrganizationStatusActive, CreatedBy: "admin",
	}
	if err := service.db.Create(&org).Error; err != nil {
		t.Fatal(err)
	}
	req := models.CreateComplianceRequirementRequest{
		OrgID: org.ID, Category: "kyb", Requirement: "Ownership declaration",
		InputType: models.ComplianceInputTypeStructuredForm, DueDate: "2026-10-01",
	}
	created, err := service.CreateAdHocRequirement(context.Background(), AuthContext{MemberID: "admin-1"}, req)
	if err != nil {
		t.Fatalf("create ad-hoc requirement: %v", err)
	}
	if created.InputType != models.ComplianceInputTypeStructuredForm || created.Required == nil || !*created.Required {
		t.Fatalf("unexpected ad-hoc requirement: %+v", created)
	}
	if _, err := service.CreateAdHocRequirement(context.Background(), AuthContext{MemberID: "admin-1"}, req); ErrorStatus(err) != http.StatusConflict {
		t.Fatalf("duplicate status=%d err=%v", ErrorStatus(err), err)
	}
}

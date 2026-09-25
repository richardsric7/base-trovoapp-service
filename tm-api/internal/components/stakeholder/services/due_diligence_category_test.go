package services

import (
	"context"
	"testing"

	stakeholderDB "admin-panel-dashboard/internal/components/stakeholder/db"
	"admin-panel-dashboard/internal/components/stakeholder/models"
	coreModels "admin-panel-dashboard/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TestUpdateDueDiligenceCategoryUpdatesOnlyRequestedCategory(t *testing.T) {
	service, database, auth := categoryDueDiligenceService(t)
	response, err := service.UpdateCategory(context.Background(), auth, "category-asset", models.UpdateDueDiligenceCategoryRequest{
		Category: " Legal ", Status: "verified", Notes: "Documents confirmed",
	})
	if err != nil {
		t.Fatalf("update category: %v", err)
	}
	if response.Category != "legal" || response.Status != models.DueDiligenceItemStatusComplete || len(response.Items) != 2 {
		t.Fatalf("response = %+v, want two completed legal items", response)
	}
	for _, item := range response.Items {
		if item.Status != models.DueDiligenceItemStatusComplete || item.Notes != "Documents confirmed" || item.VerifiedByMemberID == nil || *item.VerifiedByMemberID != auth.MemberID || item.VerifiedAt == nil {
			t.Fatalf("updated item = %+v, want verified legal item", item)
		}
	}
	var untouched int64
	if err := database.Model(&models.DueDiligenceItem{}).
		Where("checklist_id = ? AND category != ? AND status = ?", response.Items[0].ChecklistID, "legal", models.DueDiligenceItemStatusPending).
		Count(&untouched).Error; err != nil {
		t.Fatalf("count untouched items: %v", err)
	}
	if untouched != 2 {
		t.Fatalf("untouched items = %d, want 2", untouched)
	}
	var auditLog models.StakeholderAuditLog
	if err := database.Where("action = ?", models.ActionDueDiligenceCategoryUpdate).First(&auditLog).Error; err != nil {
		t.Fatalf("load category audit: %v", err)
	}
	if auditLog.EntityType != models.EntityDueDiligenceCategory {
		t.Fatalf("audit entity type = %q, want category", auditLog.EntityType)
	}
}

func TestDueDiligenceChecklistIncludesSelectedWalletAssetDocuments(t *testing.T) {
	service, _, auth := categoryDueDiligenceService(t)
	checklist, err := service.GetChecklist(context.Background(), auth, "category-asset")
	if err != nil {
		t.Fatalf("get checklist: %v", err)
	}
	if len(checklist.AvailableAssetDocuments) != 1 || checklist.AvailableAssetDocuments[0].Title != "Wallet title deed" {
		t.Fatalf("available asset documents = %+v", checklist.AvailableAssetDocuments)
	}
}

func TestUpdateDueDiligenceCategoryResetClearsVerification(t *testing.T) {
	service, _, auth := categoryDueDiligenceService(t)
	if _, err := service.UpdateCategory(context.Background(), auth, "category-asset", models.UpdateDueDiligenceCategoryRequest{Category: "legal", Status: "complete"}); err != nil {
		t.Fatalf("complete category: %v", err)
	}
	response, err := service.UpdateCategory(context.Background(), auth, "category-asset", models.UpdateDueDiligenceCategoryRequest{Category: "legal", Status: "pending"})
	if err != nil {
		t.Fatalf("reset category: %v", err)
	}
	for _, item := range response.Items {
		if item.Status != models.DueDiligenceItemStatusPending || item.VerifiedByMemberID != nil || item.VerifiedAt != nil {
			t.Fatalf("reset item = %+v, want pending without verifier", item)
		}
	}
}

func TestUpdateDueDiligenceCategoryValidatesCategoryAndChecklistState(t *testing.T) {
	service, database, auth := categoryDueDiligenceService(t)
	_, err := service.UpdateCategory(context.Background(), auth, "category-asset", models.UpdateDueDiligenceCategoryRequest{Category: "unknown", Status: "complete"})
	if ErrorStatus(err) != 404 {
		t.Fatalf("unknown category status = %d, want 404", ErrorStatus(err))
	}
	checklist, err := service.GetChecklist(context.Background(), auth, "category-asset")
	if err != nil {
		t.Fatalf("get checklist: %v", err)
	}
	if err := database.Model(&models.DueDiligenceChecklist{}).Where("id = ?", checklist.ID).Update("status", models.DueDiligenceStatusApproved).Error; err != nil {
		t.Fatalf("finalize checklist: %v", err)
	}
	_, err = service.UpdateCategory(context.Background(), auth, "category-asset", models.UpdateDueDiligenceCategoryRequest{Category: "legal", Status: "complete"})
	if ErrorStatus(err) != 409 {
		t.Fatalf("finalized checklist status = %d, want 409", ErrorStatus(err))
	}
}

func TestUpdateDueDiligenceCategoryRejectsInvalidStatusAndNonTrustee(t *testing.T) {
	service, _, auth := categoryDueDiligenceService(t)
	_, err := service.UpdateCategory(context.Background(), auth, "category-asset", models.UpdateDueDiligenceCategoryRequest{Category: "legal", Status: "unknown"})
	if ErrorStatus(err) != 400 {
		t.Fatalf("invalid status = %d, want 400", ErrorStatus(err))
	}
	auth.DashboardRole = models.DashboardRoleAssetManager
	_, err = service.UpdateCategory(context.Background(), auth, "category-asset", models.UpdateDueDiligenceCategoryRequest{Category: "legal", Status: "complete"})
	if ErrorStatus(err) != 403 {
		t.Fatalf("non-trustee status = %d, want 403", ErrorStatus(err))
	}
}

func TestTrusteeDueDiligenceItemReturnsLinkedManagerEvidence(t *testing.T) {
	service, database, auth := categoryDueDiligenceService(t)
	checklist, err := service.GetChecklist(context.Background(), auth, "category-asset")
	if err != nil {
		t.Fatalf("get checklist: %v", err)
	}
	document := models.StakeholderDocument{
		ID: uuid.NewString(), AssetID: checklist.AssetID, UploadedByOrgID: "manager-org", UploadedByMemberID: "manager-member",
		Category: models.DocumentCategoryDueDiligence, Title: "Ownership evidence", StorageObjectID: "stakeholder/evidence.pdf",
		AccessRoles: models.JSONStringArray{models.DashboardRoleAssetManager, models.DashboardRoleTrustee}, Status: models.DocumentStatusActive,
	}
	if err := database.Create(&document).Error; err != nil {
		t.Fatalf("seed due diligence document: %v", err)
	}
	documentIDs := []string{document.ID}
	updated, err := service.UpdateItem(context.Background(), auth, checklist.AssetID, checklist.Items[0].ID, models.UpdateDueDiligenceItemRequest{
		Status: models.DueDiligenceItemStatusComplete, Notes: "Reviewed", DocumentIDs: &documentIDs,
	})
	if err != nil {
		t.Fatalf("link due diligence evidence: %v", err)
	}
	if len(updated.Documents) != 1 || updated.Documents[0].ID != document.ID {
		t.Fatalf("updated documents = %+v, want %s", updated.Documents, document.ID)
	}
	reloaded, err := service.GetChecklist(context.Background(), auth, checklist.AssetID)
	if err != nil {
		t.Fatalf("reload checklist: %v", err)
	}
	var found bool
	for _, item := range reloaded.Items {
		if item.ID == updated.ID {
			found = len(item.Documents) == 1 && item.Documents[0].ID == document.ID
		}
	}
	if !found {
		t.Fatalf("linked evidence missing from checklist: %+v", reloaded.Items)
	}
}

func categoryDueDiligenceService(t *testing.T) (*DueDiligenceService, *gorm.DB, AuthContext) {
	t.Helper()
	database := testDB(t,
		&models.StakeholderAssetAssignment{},
		&models.DueDiligenceChecklist{},
		&models.DueDiligenceItem{},
		&models.StakeholderDocument{},
		&models.DueDiligenceItemDocument{},
		&models.StakeholderAuditLog{},
		&coreModels.Organization{},
	)
	trusteeOrgID := "category-trustee-org"
	trusteeStakeholderID := uint64(72)
	if err := database.Create(&models.StakeholderAssetAssignment{
		ID: "category-assignment", AssetID: "category-asset", AssetCode: "CATEGORY",
		TrusteeOrgID: &trusteeOrgID, TrusteeStakeholderID: &trusteeStakeholderID, Status: models.AssignmentStatusActive,
	}).Error; err != nil {
		t.Fatalf("seed assignment: %v", err)
	}
	auth := AuthContext{
		MemberID: "category-trustee-member", OrganizationID: trusteeOrgID,
		StakeholderID: trusteeStakeholderID, DashboardRole: models.DashboardRoleTrustee,
	}
	asset := stakeholderDB.TokenizedAsset{ID: "category-asset", AssetCode: "CATEGORY", TrusteeID: trusteeStakeholderID}
	assets := NewAssetService(database, fakeTokenizationDetailReadClient{
		fakeTokenizationReadClient: fakeTokenizationReadClient{assets: []stakeholderDB.TokenizedAsset{asset}},
		related: map[string]stakeholderDB.TokenizedAssetRelatedData{asset.ID: {
			Documents: []stakeholderDB.AssetTokenizationDocument{{ID: 91, TokenizedAssetID: asset.ID, DocumentTitle: "Wallet title deed", DocumentURL: "https://wallet.example/title.pdf"}},
		}},
	}, nil)
	audit := NewAuditService(database)
	documents := NewDocumentService(database, assets, audit, nil)
	return NewDueDiligenceService(database, assets, documents, audit, nil), database, auth
}

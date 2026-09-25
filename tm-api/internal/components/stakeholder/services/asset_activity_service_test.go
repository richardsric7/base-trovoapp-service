package services

import (
	"context"
	"testing"
	"time"

	stakeholderDB "admin-panel-dashboard/internal/components/stakeholder/db"
	"admin-panel-dashboard/internal/components/stakeholder/models"
	coreModels "admin-panel-dashboard/internal/models"
)

func TestListAssetActivitiesReturnsOnlyVisibleAssetEventsNewestFirst(t *testing.T) {
	database := testDB(t,
		&models.StakeholderAssetAssignment{},
		&models.FundReleaseRequest{},
		&models.RevenueRecord{},
		&models.Distribution{},
		&models.AssetValuation{},
		&models.SegregatedAccount{},
		&models.ComplianceItem{},
		&models.StakeholderDocument{},
		&models.StakeholderAssetOperation{},
		&models.StakeholderReport{},
		&models.StakeholderStructuringStatus{},
		&models.DueDiligenceChecklist{},
		&models.DueDiligenceItem{},
		&models.StakeholderAuditLog{},
		&coreModels.Organization{},
		&coreModels.OrganizationMember{},
	)
	asset := stakeholderDB.TokenizedAsset{ID: "activity-asset", AssetCode: "ACT", AssetManagerID: 71}
	service := NewAssetService(database, fakeTokenizationReadClient{assets: []stakeholderDB.TokenizedAsset{asset}}, nil)
	actorOrganization := coreModels.Organization{
		ID: "trustee-org", Name: "Trusted Trustees", Email: "trustee@example.com", Type: "CORPORATE",
		Status: coreModels.OrganizationStatusActive, CreatedBy: "test",
	}
	actorMember := coreModels.OrganizationMember{
		ID: "trustee-member", OrganizationID: actorOrganization.ID, Email: "nancy@example.com", FirstName: "Nancy", LastName: "Rike",
		Role: string(coreModels.OrganizationMemberRole), Status: coreModels.OrganizationStatusActive, CreatedBy: "test",
	}
	if err := database.Create(&actorOrganization).Error; err != nil {
		t.Fatalf("seed actor organization: %v", err)
	}
	if err := database.Create(&actorMember).Error; err != nil {
		t.Fatalf("seed actor member: %v", err)
	}

	distributions := []models.Distribution{
		{ID: "asset-distribution", AssetID: asset.ID, AssetCode: asset.AssetCode, ProposedByOrgID: "manager-org", ProposedByMemberID: "manager-member", TrusteeOrgID: "trustee-org", Currency: "CNGN", Source: "income", Status: models.DistributionStatusProposed},
		{ID: "other-distribution", AssetID: "other-asset", AssetCode: "OTHER", ProposedByOrgID: "other-org", ProposedByMemberID: "other-member", TrusteeOrgID: "other-trustee", Currency: "CNGN", Source: "income", Status: models.DistributionStatusProposed},
	}
	if err := database.Create(&distributions).Error; err != nil {
		t.Fatalf("seed distributions: %v", err)
	}
	checklist := models.DueDiligenceChecklist{ID: "asset-checklist", AssetID: asset.ID, AssetCode: asset.AssetCode, TrusteeOrgID: "trustee-org", Status: models.DueDiligenceStatusInReview}
	if err := database.Create(&checklist).Error; err != nil {
		t.Fatalf("seed due-diligence checklist: %v", err)
	}

	now := time.Now()
	logs := []models.StakeholderAuditLog{
		{ID: "distribution-log", Action: models.ActionDistributionSubmit, EntityType: models.EntityDistribution, EntityID: distributions[0].ID, CreatedAt: now.Add(-time.Minute)},
		{ID: "category-log", ActorMemberID: actorMember.ID, ActorOrgID: actorOrganization.ID, Action: models.ActionDueDiligenceCategoryUpdate, EntityType: models.EntityDueDiligenceCategory, EntityID: checklist.ID + ":legal", CreatedAt: now},
		{ID: "unrelated-log", Action: models.ActionDistributionSubmit, EntityType: models.EntityDistribution, EntityID: distributions[1].ID, CreatedAt: now.Add(time.Minute)},
	}
	if err := database.Create(&logs).Error; err != nil {
		t.Fatalf("seed audit logs: %v", err)
	}

	auth := AuthContext{OrganizationID: "manager-org", StakeholderID: 71, DashboardRole: models.DashboardRoleAssetManager}
	activities, total, err := service.ListAssetActivities(context.Background(), auth, asset.ID, 1, 1)
	if err != nil {
		t.Fatalf("list asset activities: %v", err)
	}
	if total != 2 || len(activities) != 1 || activities[0].ID != "category-log" {
		t.Fatalf("activities=%+v total=%d, want newest of two asset events", activities, total)
	}
	if activities[0].ActorName != "Nancy Rike" || activities[0].ActorOrganization != actorOrganization.Name {
		t.Fatalf("activity actor=%+v, want resolved member and organization names", activities[0])
	}

	_, _, err = service.ListAssetActivities(context.Background(), AuthContext{
		OrganizationID: "other-manager", StakeholderID: 999, DashboardRole: models.DashboardRoleAssetManager,
	}, asset.ID, 1, 20)
	if err == nil || ErrorStatus(err) != 404 {
		t.Fatalf("hidden asset error=%v status=%d, want 404", err, ErrorStatus(err))
	}
}

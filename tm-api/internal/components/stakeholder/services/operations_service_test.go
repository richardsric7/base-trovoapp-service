package services

import (
	"context"
	"testing"
	"time"

	stakeholderDB "admin-panel-dashboard/internal/components/stakeholder/db"
	"admin-panel-dashboard/internal/components/stakeholder/models"
	coreModels "admin-panel-dashboard/internal/models"
)

func TestOperationsServiceValidatesAndPersistsOperationalUpdate(t *testing.T) {
	database := testDB(t, &models.StakeholderAssetAssignment{}, &models.StakeholderAssetOperation{}, &coreModels.Organization{})
	managerID := uint64(21)
	managerOrgID := "manager-org"
	asset := stakeholderDB.TokenizedAsset{ID: "asset-1", AssetCode: "AST", AssetManagerID: managerID}
	stakeholderType := models.StakeholderTypeAssetManager
	if err := database.Create(&coreModels.Organization{ID: managerOrgID, Name: "Manager", Email: "operations-manager@example.com", Type: "CORPORATE", Status: coreModels.OrganizationStatusActive, CreatedBy: "test", StakeholderID: &managerID, StakeholderType: &stakeholderType}).Error; err != nil {
		t.Fatalf("seed manager organization: %v", err)
	}
	assets := NewAssetService(database, fakeTokenizationReadClient{assets: []stakeholderDB.TokenizedAsset{asset}}, nil)
	service := NewOperationsService(database, assets, nil, nil, nil)
	auth := AuthContext{MemberID: "member-1", OrganizationID: managerOrgID, StakeholderID: managerID, DashboardRole: models.DashboardRoleAssetManager}

	if _, err := service.UpdateAssetOperation(context.Background(), auth, asset.ID, models.UpdateAssetOperationRequest{OperationalStatus: "unknown"}); ErrorStatus(err) != 400 {
		t.Fatalf("invalid status code = %d, want 400", ErrorStatus(err))
	}
	if _, err := service.UpdateAssetOperation(context.Background(), auth, asset.ID, models.UpdateAssetOperationRequest{
		OperationalStatus: models.AssetOperationalStatusActive,
		Metadata:          models.AssetOperationMetadata{EffectiveAt: "not-a-time"},
	}); ErrorStatus(err) != 400 {
		t.Fatalf("invalid effective_at code = %d, want 400", ErrorStatus(err))
	}

	effectiveAt := time.Now().UTC().Truncate(time.Second).Format(time.RFC3339)
	operation, err := service.UpdateAssetOperation(context.Background(), auth, asset.ID, models.UpdateAssetOperationRequest{
		OperationalStatus: " ACTIVE ", Notes: "Inspected",
		Metadata: models.AssetOperationMetadata{EffectiveAt: effectiveAt, Source: "manual"},
	})
	if err != nil {
		t.Fatalf("update operation: %v", err)
	}
	if operation.OperationalStatus != models.AssetOperationalStatusActive || operation.Metadata["effective_at"] != effectiveAt || operation.Metadata["source"] != "manual" {
		t.Fatalf("operation = %+v", operation)
	}
}

func TestOperationsServiceCreatesManagedDocumentReport(t *testing.T) {
	database := testDB(t, &models.StakeholderAssetAssignment{}, &models.StakeholderDocument{}, &models.StakeholderReport{}, &coreModels.Organization{})
	managerID := uint64(22)
	managerOrgID := "report-manager-org"
	asset := stakeholderDB.TokenizedAsset{ID: "asset-report", AssetCode: "RPT", AssetManagerID: managerID}
	stakeholderType := models.StakeholderTypeAssetManager
	if err := database.Create(&coreModels.Organization{ID: managerOrgID, Name: "Report Manager", Email: "report-manager@example.com", Type: "CORPORATE", Status: coreModels.OrganizationStatusActive, CreatedBy: "test", StakeholderID: &managerID, StakeholderType: &stakeholderType}).Error; err != nil {
		t.Fatalf("seed manager organization: %v", err)
	}
	assets := NewAssetService(database, fakeTokenizationReadClient{assets: []stakeholderDB.TokenizedAsset{asset}}, nil)
	documents := NewDocumentService(database, assets, nil, nil)
	service := NewOperationsService(database, assets, documents, nil, nil)
	auth := AuthContext{MemberID: "report-member", OrganizationID: managerOrgID, StakeholderID: managerID, DashboardRole: models.DashboardRoleAssetManager}
	document := models.StakeholderDocument{
		ID: "83907271-4e92-4143-92dd-1737286a9bab", AssetID: asset.ID, AssetCode: asset.AssetCode,
		UploadedByOrgID: managerOrgID, UploadedByMemberID: auth.MemberID,
		Category: models.DocumentCategoryAssetReport, Title: "Income report",
		FileURL: "wallet://report-document", StorageObjectID: "object-report-document",
		Status: models.DocumentStatusActive, AccessRoles: models.JSONStringArray{models.DashboardRoleAssetManager, models.DashboardRoleTrustee},
	}
	if err := database.Create(&document).Error; err != nil {
		t.Fatalf("seed report document: %v", err)
	}

	if _, err := service.CreateReport(context.Background(), auth, models.CreateReportRequest{
		AssetID: asset.ID, ReportType: "unsupported", Title: "Invalid", DocumentID: document.ID,
	}); ErrorStatus(err) != 400 {
		t.Fatalf("invalid report type code = %d, want 400", ErrorStatus(err))
	}
	report, err := service.CreateReport(context.Background(), auth, models.CreateReportRequest{
		AssetID: asset.ID, ReportType: models.ReportTypeIncome, Title: "August income", DocumentID: document.ID,
	})
	if err != nil {
		t.Fatalf("create report: %v", err)
	}
	if report.DocumentID == nil || *report.DocumentID != document.ID || report.FileURL != "" {
		t.Fatalf("report document linkage = %+v", report)
	}
	if types, err := service.ReportTypes(auth); err != nil || len(types) != 2 {
		t.Fatalf("report types = %+v, err=%v", types, err)
	}
}

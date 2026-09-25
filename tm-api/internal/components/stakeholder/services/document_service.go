package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"admin-panel-dashboard/internal/components/stakeholder/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DocumentService struct {
	db      *gorm.DB
	assets  *AssetService
	audit   *AuditService
	storage DocumentStorageClient
}

func NewDocumentService(db *gorm.DB, assets *AssetService, audit *AuditService, storage DocumentStorageClient) *DocumentService {
	return &DocumentService{db: db, assets: assets, audit: audit, storage: storage}
}

type DocumentCategory struct {
	Value       string   `json:"value"`
	Label       string   `json:"label"`
	AccessRoles []string `json:"access_roles"`
}

var documentCategories = []DocumentCategory{
	{models.DocumentCategoryValuationReport, "Valuation report", []string{models.DashboardRoleAssetManager, models.DashboardRoleTrustee}},
	{models.DocumentCategoryAssetReport, "Asset report", []string{models.DashboardRoleAssetManager, models.DashboardRoleTrustee}},
	{models.DocumentCategoryFundReleaseSupporting, "Fund release supporting document", []string{models.DashboardRoleAssetManager, models.DashboardRoleTrustee, models.DashboardRoleAssetCustodian}},
	{models.DocumentCategoryDueDiligence, "Due diligence", []string{models.DashboardRoleAssetManager, models.DashboardRoleTrustee}},
	{models.DocumentCategoryCustodyCompliance, "Custody compliance", []string{models.DashboardRoleAssetCustodian, models.DashboardRoleTrustee}},
	{models.DocumentCategoryComplianceRequirement, "Organization compliance requirement", []string{models.DashboardRoleAssetManager, models.DashboardRoleTrustee, models.DashboardRoleAssetCustodian, models.DashboardRoleLegalAdviser, models.DashboardRoleFinancialAdviser, models.DashboardRoleIssuingHouse, models.DashboardRoleRatingAgency}},
	{models.DocumentCategoryLegal, "Legal", []string{models.DashboardRoleLegalAdviser, models.DashboardRoleAssetManager, models.DashboardRoleTrustee}},
	{models.DocumentCategoryFinancial, "Financial", []string{models.DashboardRoleFinancialAdviser, models.DashboardRoleAssetManager, models.DashboardRoleTrustee}},
	{models.DocumentCategoryStructuring, "Structuring", []string{models.DashboardRoleLegalAdviser, models.DashboardRoleFinancialAdviser, models.DashboardRoleAssetManager, models.DashboardRoleTrustee}},
	{models.DocumentCategoryGeneral, "General", []string{models.DashboardRoleTrustee, models.DashboardRoleAssetCustodian, models.DashboardRoleAssetManager, models.DashboardRoleLegalAdviser, models.DashboardRoleFinancialAdviser, models.DashboardRoleIssuingHouse, models.DashboardRoleRatingAgency}},
}

func (s *DocumentService) Categories() []DocumentCategory {
	result := make([]DocumentCategory, len(documentCategories))
	copy(result, documentCategories)
	return result
}

func (s *DocumentService) Upload(ctx context.Context, auth AuthContext, req models.UploadDocumentRequest) (*models.StakeholderDocument, error) {
	if s.storage == nil {
		return nil, NewHTTPError(http.StatusServiceUnavailable, "document storage is not configured")
	}
	category, accessRoles, err := validatedDocumentCategory(req.Category, req.AccessRoles)
	if err != nil {
		return nil, err
	}
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, NewHTTPError(http.StatusBadRequest, "title is required")
	}
	if len(req.Content) == 0 {
		return nil, NewHTTPError(http.StatusBadRequest, "document_file is required")
	}
	if len(req.Content) > 10<<20 {
		return nil, NewHTTPError(http.StatusRequestEntityTooLarge, "document exceeds 10MB")
	}
	mimeType := http.DetectContentType(req.Content)
	if mimeType != "application/pdf" && mimeType != "image/jpeg" && mimeType != "image/png" {
		return nil, NewHTTPError(http.StatusBadRequest, "unsupported document type; use PDF, JPEG, or PNG")
	}
	assetID := strings.TrimSpace(req.AssetID)
	assetCode := ""
	if assetID == "" && category != models.DocumentCategoryGeneral && category != models.DocumentCategoryCustodyCompliance && category != models.DocumentCategoryComplianceRequirement {
		return nil, NewHTTPError(http.StatusBadRequest, "asset_id is required for this document category")
	}
	if assetID != "" {
		if s.assets == nil {
			return nil, NewHTTPError(http.StatusServiceUnavailable, "asset service is not configured")
		}
		asset, err := s.assets.GetAsset(ctx, auth, assetID)
		if err != nil {
			return nil, err
		}
		assetID, assetCode = asset.ID, asset.AssetCode
	}
	filename := filepath.Base(strings.ReplaceAll(strings.TrimSpace(req.Filename), `\`, "/"))
	filename = strings.NewReplacer(`"`, "", "\r", "", "\n", "").Replace(filename)
	if filename == "." || filename == "" {
		filename = "document"
	}
	stored, err := s.storage.Upload(ctx, filename, mimeType, req.Content)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	document := models.StakeholderDocument{
		ID: uuid.NewString(), AssetID: assetID, AssetCode: assetCode,
		UploadedByOrgID: auth.OrganizationID, UploadedByMemberID: auth.MemberID,
		Category: category, Title: title, FileURL: "", StorageObjectID: stored.ObjectID,
		OriginalFilename: stored.OriginalFilename, MimeType: stored.MimeType,
		SizeBytes: stored.SizeBytes, SHA256: stored.SHA256, Version: 1,
		AccessRoles: models.JSONStringArray(accessRoles), Status: models.DocumentStatusActive,
		CreatedAt: now, UpdatedAt: now,
	}
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&document).Error; err != nil {
			return err
		}
		if s.audit != nil {
			event := AuditEventFromAuth(auth, models.ActionDocumentCreate, models.EntityDocument, document.ID)
			event.AfterState = models.JSONMap{"asset_id": document.AssetID, "category": document.Category, "title": document.Title, "access_roles": accessRoles}
			return s.audit.Record(ctx, tx, event)
		}
		return nil
	})
	if err != nil {
		if cleanupErr := s.storage.Delete(ctx, stored.ObjectID); cleanupErr != nil {
			return nil, fmt.Errorf("save document metadata: %w (storage cleanup failed: %v)", err, cleanupErr)
		}
		return nil, err
	}
	return &document, nil
}

func (s *DocumentService) Create(ctx context.Context, auth AuthContext, req models.CreateDocumentRequest) (*models.StakeholderDocument, error) {
	if strings.TrimSpace(req.AssetID) != "" {
		if _, err := s.assets.GetAsset(ctx, auth, req.AssetID); err != nil {
			return nil, err
		}
	}
	category := strings.TrimSpace(req.Category)
	accessRoles := make([]string, 0, len(req.AccessRoles))
	for _, role := range req.AccessRoles {
		role = models.Normalize(role)
		if !models.IsSupportedDashboardRole(role) {
			return nil, NewHTTPError(http.StatusBadRequest, "access_roles contains unsupported dashboard role")
		}
		accessRoles = append(accessRoles, role)
	}
	version := req.Version
	if version < 1 {
		version = 1
	}
	now := time.Now()
	document := models.StakeholderDocument{
		ID:                 uuid.NewString(),
		AssetID:            strings.TrimSpace(req.AssetID),
		AssetCode:          strings.TrimSpace(req.AssetCode),
		UploadedByOrgID:    auth.OrganizationID,
		UploadedByMemberID: auth.MemberID,
		Category:           category,
		Title:              strings.TrimSpace(req.Title),
		FileURL:            strings.TrimSpace(req.FileURL),
		MimeType:           strings.TrimSpace(req.MimeType),
		Version:            version,
		AccessRoles:        models.JSONStringArray(accessRoles),
		Status:             models.DocumentStatusActive,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	if document.Category == "" || document.Title == "" || document.FileURL == "" {
		return nil, NewHTTPError(http.StatusBadRequest, "category, title, and file_url are required")
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&document).Error; err != nil {
			return err
		}
		if s.audit != nil {
			event := AuditEventFromAuth(auth, models.ActionDocumentCreate, models.EntityDocument, document.ID)
			event.AfterState = models.JSONMap{"asset_id": document.AssetID, "category": document.Category, "title": document.Title, "access_roles": accessRoles}
			if err := s.audit.Record(ctx, tx, event); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &document, nil
}

func (s *DocumentService) List(ctx context.Context, auth AuthContext, page, limit int, assetID string) ([]models.StakeholderDocument, int64, error) {
	query := s.db.WithContext(ctx).Where("status = ?", models.DocumentStatusActive)
	if strings.TrimSpace(assetID) != "" {
		if _, err := s.assets.GetAsset(ctx, auth, assetID); err != nil {
			return nil, 0, err
		}
		query = query.Where("asset_id = ?", assetID)
	}
	var docs []models.StakeholderDocument
	if err := query.Order("created_at desc").Find(&docs).Error; err != nil {
		return nil, 0, err
	}
	filtered := make([]models.StakeholderDocument, 0, len(docs))
	for _, doc := range docs {
		if !models.RoleCanAccessDocument(auth.DashboardRole, []string(doc.AccessRoles)) {
			continue
		}
		if doc.AssetID != "" {
			if _, err := s.assets.GetAsset(ctx, auth, doc.AssetID); err != nil {
				continue
			}
		} else if doc.UploadedByOrgID != auth.OrganizationID {
			continue
		}
		filtered = append(filtered, doc)
	}
	total := int64(len(filtered))
	start := (page - 1) * limit
	if start >= len(filtered) {
		return []models.StakeholderDocument{}, total, nil
	}
	end := start + limit
	if end > len(filtered) {
		end = len(filtered)
	}
	return filtered[start:end], total, nil
}

func (s *DocumentService) Get(ctx context.Context, auth AuthContext, id string) (*models.StakeholderDocument, error) {
	var doc models.StakeholderDocument
	if err := s.db.WithContext(ctx).Where("id = ? AND status = ?", id, models.DocumentStatusActive).First(&doc).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewHTTPError(http.StatusNotFound, "document not found")
		}
		return nil, err
	}
	if !models.RoleCanAccessDocument(auth.DashboardRole, []string(doc.AccessRoles)) {
		return nil, NewHTTPError(http.StatusNotFound, "document not found")
	}
	if doc.AssetID != "" {
		if _, err := s.assets.GetAsset(ctx, auth, doc.AssetID); err != nil {
			return nil, err
		}
	} else if doc.UploadedByOrgID != auth.OrganizationID {
		return nil, NewHTTPError(http.StatusNotFound, "document not found")
	}
	return &doc, nil
}

func (s *DocumentService) Download(ctx context.Context, auth AuthContext, id string) (*DocumentDownload, error) {
	doc, err := s.Get(ctx, auth, id)
	if err != nil {
		return nil, err
	}
	return s.downloadStored(ctx, doc)
}

func (s *DocumentService) DownloadOwnedByOrganization(ctx context.Context, organizationID, id, category string) (*DocumentDownload, error) {
	var doc models.StakeholderDocument
	query := s.db.WithContext(ctx).Where("id = ? AND uploaded_by_org_id = ? AND status = ?", id, organizationID, models.DocumentStatusActive)
	if strings.TrimSpace(category) != "" {
		query = query.Where("category = ?", category)
	}
	if err := query.First(&doc).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewHTTPError(http.StatusNotFound, "document not found")
		}
		return nil, err
	}
	return s.downloadStored(ctx, &doc)
}

func (s *DocumentService) downloadStored(ctx context.Context, doc *models.StakeholderDocument) (*DocumentDownload, error) {
	if strings.TrimSpace(doc.StorageObjectID) == "" {
		return nil, NewHTTPError(http.StatusUnprocessableEntity, "this legacy document has no managed file; use file_url")
	}
	if s.storage == nil {
		return nil, NewHTTPError(http.StatusServiceUnavailable, "document storage is not configured")
	}
	return s.storage.Download(ctx, doc.StorageObjectID)
}

func (s *DocumentService) ValidateForUse(ctx context.Context, auth AuthContext, assetID, category string, ids []string) error {
	if len(ids) == 0 {
		return NewHTTPError(http.StatusUnprocessableEntity, "at least one document_id is required")
	}
	seen := make(map[string]struct{}, len(ids))
	for _, rawID := range ids {
		id := strings.TrimSpace(rawID)
		if _, err := uuid.Parse(id); err != nil {
			return NewHTTPError(http.StatusUnprocessableEntity, "document_id is invalid")
		}
		if _, exists := seen[id]; exists {
			return NewHTTPError(http.StatusUnprocessableEntity, "document_ids must be unique")
		}
		seen[id] = struct{}{}
		var doc models.StakeholderDocument
		err := s.db.WithContext(ctx).Where("id = ? AND status = ?", id, models.DocumentStatusActive).First(&doc).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return NewHTTPError(http.StatusUnprocessableEntity, "document_id does not reference an active document")
		}
		if err != nil {
			return err
		}
		if doc.UploadedByOrgID != auth.OrganizationID || doc.AssetID != assetID || doc.Category != category || strings.TrimSpace(doc.StorageObjectID) == "" {
			return NewHTTPError(http.StatusUnprocessableEntity, "document_id does not belong to this organization, asset, and document category")
		}
	}
	return nil
}

// ValidateAccessibleForUse validates evidence selected from an assigned asset.
// Unlike ValidateForUse, the uploader may be another organization in the same
// stakeholder workflow, but the caller must have document-role access.
func (s *DocumentService) ValidateAccessibleForUse(ctx context.Context, auth AuthContext, assetID, category string, ids []string) error {
	seen := make(map[string]struct{}, len(ids))
	for _, rawID := range ids {
		id := strings.TrimSpace(rawID)
		if _, err := uuid.Parse(id); err != nil {
			return NewHTTPError(http.StatusUnprocessableEntity, "document_id is invalid")
		}
		if _, exists := seen[id]; exists {
			return NewHTTPError(http.StatusUnprocessableEntity, "document_ids must be unique")
		}
		seen[id] = struct{}{}
		var document models.StakeholderDocument
		err := s.db.WithContext(ctx).Where("id = ? AND status = ?", id, models.DocumentStatusActive).First(&document).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return NewHTTPError(http.StatusUnprocessableEntity, "document_id does not reference an active document")
		}
		if err != nil {
			return err
		}
		if document.AssetID != assetID || document.Category != category || strings.TrimSpace(document.StorageObjectID) == "" ||
			!models.RoleCanAccessDocument(auth.DashboardRole, []string(document.AccessRoles)) {
			return NewHTTPError(http.StatusUnprocessableEntity, "document_id is not accessible for this asset and document category")
		}
	}
	return nil
}

func validatedDocumentCategory(value string, requestedRoles []string) (string, []string, error) {
	category := models.Normalize(value)
	var configured *DocumentCategory
	for i := range documentCategories {
		if documentCategories[i].Value == category {
			configured = &documentCategories[i]
			break
		}
	}
	if configured == nil {
		return "", nil, NewHTTPError(http.StatusBadRequest, "unsupported document category")
	}
	roles := append([]string{}, configured.AccessRoles...)
	seen := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		seen[role] = struct{}{}
	}
	for _, role := range requestedRoles {
		role = models.Normalize(role)
		if !models.IsSupportedDashboardRole(role) {
			return "", nil, NewHTTPError(http.StatusBadRequest, "access_roles contains unsupported dashboard role")
		}
		if _, exists := seen[role]; !exists {
			return "", nil, NewHTTPError(http.StatusBadRequest, "access_roles contains a role not allowed for this category")
		}
	}
	return category, roles, nil
}

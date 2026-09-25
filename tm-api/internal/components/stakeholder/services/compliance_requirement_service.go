package services

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"admin-panel-dashboard/internal/components/stakeholder/models"
	coreModels "admin-panel-dashboard/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ComplianceRequirementService implements the templated KYC/KYB compliance
// system: template CRUD, the admin review queue, the organisation submission
// surface, and auto-assignment of requirement instances from templates. See
// PRD-kyc-compliance-templates.md.
type ComplianceRequirementService struct {
	db        *gorm.DB
	audit     *AuditService
	documents *DocumentService
}

func NewComplianceRequirementService(db *gorm.DB, audit *AuditService, documents *DocumentService) *ComplianceRequirementService {
	return &ComplianceRequirementService{db: db, audit: audit, documents: documents}
}

var complianceInputTypes = []string{
	models.ComplianceInputTypeDocumentUpload,
	models.ComplianceInputTypeText,
	models.ComplianceInputTypeStructuredForm,
}

var complianceRequirementStatuses = []string{
	models.ComplianceRequirementStatusDraft,
	models.ComplianceRequirementStatusPending,
	models.ComplianceRequirementStatusSubmitted,
	models.ComplianceRequirementStatusUnderReview,
	models.ComplianceRequirementStatusApproved,
	models.ComplianceRequirementStatusRejected,
	models.ComplianceRequirementStatusOverdue,
	models.ComplianceRequirementStatusWaived,
}

func normalizeComplianceStatusFilter(value string) (string, error) {
	status := strings.ToLower(strings.TrimSpace(value))
	if status != "" && !ContainsString(complianceRequirementStatuses, status) {
		return "", NewHTTPError(http.StatusBadRequest, "invalid compliance requirement status")
	}
	return status, nil
}

var complianceOrganizationTypes = []string{
	coreModels.StakeholderTypeAssetManager,
	coreModels.StakeholderTypeIssuingHouse,
	coreModels.StakeholderTypeCustodian,
	coreModels.StakeholderTypeLegal,
	coreModels.StakeholderTypeRatingAgency,
	coreModels.StakeholderTypeTrustee,
	coreModels.StakeholderTypeLegalAdviser,
	coreModels.StakeholderTypeFinancialAdviser,
}

func normalizeComplianceOrganizationType(value string) (string, error) {
	orgType := strings.ToLower(strings.TrimSpace(value))
	if !ContainsString(complianceOrganizationTypes, orgType) {
		return "", NewHTTPError(http.StatusBadRequest, "unsupported org_type")
	}
	return orgType, nil
}

func complianceItemRequired(value *bool) *bool {
	if value == nil {
		defaultRequired := true
		return &defaultRequired
	}
	required := *value
	return &required
}

// ---- Templates ----

func (s *ComplianceRequirementService) ListTemplates(ctx context.Context, page, limit int, orgType string, level *int) ([]models.ComplianceRequirementTemplate, int64, error) {
	query := s.db.WithContext(ctx).Model(&models.ComplianceRequirementTemplate{})
	if orgType != "" {
		var err error
		orgType, err = normalizeComplianceOrganizationType(orgType)
		if err != nil {
			return nil, 0, err
		}
		query = query.Where("org_type = ?", orgType)
	}
	if level != nil {
		if *level <= 0 {
			return nil, 0, NewHTTPError(http.StatusBadRequest, "level must be a positive integer")
		}
		query = query.Where("level = ?", *level)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var templates []models.ComplianceRequirementTemplate
	if err := query.Preload("Items").Order("created_at desc").Offset((page - 1) * limit).Limit(limit).Find(&templates).Error; err != nil {
		return nil, 0, err
	}
	return templates, total, nil
}

func validateTemplateItems(items []models.ComplianceTemplateItemInput) error {
	if len(items) == 0 {
		return NewHTTPError(http.StatusBadRequest, "at least one requirement item is required")
	}
	seenIDs := make(map[string]struct{}, len(items))
	seenDefinitions := make(map[string]struct{}, len(items))
	for _, item := range items {
		name := strings.TrimSpace(item.Name)
		category := strings.TrimSpace(item.Category)
		if name == "" || category == "" {
			return NewHTTPError(http.StatusBadRequest, "each item requires a name and category")
		}
		if !ContainsString(complianceInputTypes, strings.TrimSpace(item.InputType)) {
			return NewHTTPError(http.StatusBadRequest, "invalid input_type for requirement item")
		}
		if id := strings.TrimSpace(item.ID); id != "" {
			if _, err := uuid.Parse(id); err != nil {
				return NewHTTPError(http.StatusBadRequest, "template item id is invalid")
			}
			if _, exists := seenIDs[id]; exists {
				return NewHTTPError(http.StatusBadRequest, "template item ids must be unique")
			}
			seenIDs[id] = struct{}{}
		}
		definitionKey := strings.ToLower(category) + "\x00" + strings.ToLower(name)
		if _, exists := seenDefinitions[definitionKey]; exists {
			return NewHTTPError(http.StatusBadRequest, "template item names must be unique within a category")
		}
		seenDefinitions[definitionKey] = struct{}{}
	}
	return nil
}

func (s *ComplianceRequirementService) CreateTemplate(ctx context.Context, actor AuthContext, req models.CreateComplianceTemplateRequest) (*models.ComplianceRequirementTemplate, error) {
	orgType, err := normalizeComplianceOrganizationType(req.OrgType)
	if err != nil {
		return nil, err
	}
	if req.Level <= 0 {
		return nil, NewHTTPError(http.StatusBadRequest, "level must be a positive integer")
	}
	if err := validateTemplateItems(req.Items); err != nil {
		return nil, err
	}

	now := time.Now()
	template := models.ComplianceRequirementTemplate{
		ID:        uuid.NewString(),
		OrgType:   orgType,
		Level:     req.Level,
		Version:   1,
		CreatedAt: now,
		UpdatedAt: now,
	}
	for _, item := range req.Items {
		template.Items = append(template.Items, models.ComplianceRequirementTemplateItem{
			ID:          uuid.NewString(),
			Name:        strings.TrimSpace(item.Name),
			Category:    strings.TrimSpace(item.Category),
			Description: strings.TrimSpace(item.Description),
			InputType:   strings.TrimSpace(item.InputType),
			Required:    complianceItemRequired(item.Required),
			CreatedAt:   now,
		})
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existingCount int64
		if err := tx.Model(&models.ComplianceRequirementTemplate{}).
			Where("org_type = ? AND level = ?", orgType, req.Level).
			Count(&existingCount).Error; err != nil {
			return err
		}
		if existingCount > 0 {
			return NewHTTPError(http.StatusConflict, "a template already exists for this organisation type and level")
		}
		if err := tx.Create(&template).Error; err != nil {
			return err
		}
		if s.audit != nil {
			event := AuditEventFromAuth(actor, models.ActionComplianceTemplateCreate, models.EntityComplianceTemplate, template.ID)
			event.AfterState = models.JSONMap{"org_type": template.OrgType, "level": template.Level, "version": template.Version, "item_count": len(template.Items)}
			if err := s.audit.Record(ctx, tx, event); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (s *ComplianceRequirementService) UpdateTemplate(ctx context.Context, actor AuthContext, id string, req models.UpdateComplianceTemplateRequest) (*models.ComplianceRequirementTemplate, error) {
	if err := validateTemplateItems(req.Items); err != nil {
		return nil, err
	}
	if req.Level != nil && *req.Level <= 0 {
		return nil, NewHTTPError(http.StatusBadRequest, "level must be a positive integer")
	}
	var normalizedOrgType *string
	if req.OrgType != nil {
		value, err := normalizeComplianceOrganizationType(*req.OrgType)
		if err != nil {
			return nil, err
		}
		normalizedOrgType = &value
	}

	var updated models.ComplianceRequirementTemplate
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var template models.ComplianceRequirementTemplate
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id).First(&template).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return NewHTTPError(http.StatusNotFound, "compliance template not found")
			}
			return err
		}
		before := models.JSONMap{"org_type": template.OrgType, "level": template.Level, "version": template.Version}

		if normalizedOrgType != nil {
			template.OrgType = *normalizedOrgType
		}
		if req.Level != nil {
			template.Level = *req.Level
		}
		var conflictingCount int64
		if err := tx.Model(&models.ComplianceRequirementTemplate{}).
			Where("id <> ? AND org_type = ? AND level = ?", template.ID, template.OrgType, template.Level).
			Count(&conflictingCount).Error; err != nil {
			return err
		}
		if conflictingCount > 0 {
			return NewHTTPError(http.StatusConflict, "a template already exists for this organisation type and level")
		}
		template.Version++
		template.UpdatedAt = time.Now()

		var existingItems []models.ComplianceRequirementTemplateItem
		if err := tx.Where("template_id = ?", template.ID).Find(&existingItems).Error; err != nil {
			return err
		}
		existingByID := make(map[string]models.ComplianceRequirementTemplateItem, len(existingItems))
		existingByDefinition := make(map[string]models.ComplianceRequirementTemplateItem, len(existingItems))
		for _, item := range existingItems {
			existingByID[item.ID] = item
			key := strings.ToLower(strings.TrimSpace(item.Category)) + "\x00" + strings.ToLower(strings.TrimSpace(item.Name))
			existingByDefinition[key] = item
		}

		now := time.Now()
		items := make([]models.ComplianceRequirementTemplateItem, 0, len(req.Items))
		retainedIDs := make([]string, 0, len(req.Items))
		selectedIDs := make(map[string]struct{}, len(req.Items))
		for _, input := range req.Items {
			itemID := strings.TrimSpace(input.ID)
			createdAt := now
			if itemID != "" {
				existing, exists := existingByID[itemID]
				if !exists {
					return NewHTTPError(http.StatusBadRequest, "template item id does not belong to this template")
				}
				createdAt = existing.CreatedAt
			} else {
				key := strings.ToLower(strings.TrimSpace(input.Category)) + "\x00" + strings.ToLower(strings.TrimSpace(input.Name))
				if existing, exists := existingByDefinition[key]; exists {
					if _, alreadySelected := selectedIDs[existing.ID]; alreadySelected {
						itemID = uuid.NewString()
					} else {
						itemID = existing.ID
						createdAt = existing.CreatedAt
					}
				} else {
					itemID = uuid.NewString()
				}
			}
			items = append(items, models.ComplianceRequirementTemplateItem{
				ID:          itemID,
				TemplateID:  template.ID,
				Name:        strings.TrimSpace(input.Name),
				Category:    strings.TrimSpace(input.Category),
				Description: strings.TrimSpace(input.Description),
				InputType:   strings.TrimSpace(input.InputType),
				Required:    complianceItemRequired(input.Required),
				CreatedAt:   createdAt,
			})
			retainedIDs = append(retainedIDs, itemID)
			selectedIDs[itemID] = struct{}{}
		}
		if err := tx.Save(&template).Error; err != nil {
			return err
		}
		for i := range items {
			if err := tx.Save(&items[i]).Error; err != nil {
				return err
			}
		}
		if err := tx.Where("template_id = ? AND id NOT IN ?", template.ID, retainedIDs).
			Delete(&models.ComplianceRequirementTemplateItem{}).Error; err != nil {
			return err
		}
		template.Items = items
		updated = template

		if s.audit != nil {
			event := AuditEventFromAuth(actor, models.ActionComplianceTemplateUpdate, models.EntityComplianceTemplate, template.ID)
			event.BeforeState = before
			event.AfterState = models.JSONMap{"org_type": template.OrgType, "level": template.Level, "version": template.Version, "item_count": len(items)}
			if err := s.audit.Record(ctx, tx, event); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &updated, nil
}

// ---- Admin review queue ----

func (s *ComplianceRequirementService) ListRequirements(ctx context.Context, page, limit int, orgID, category string, level *int, status string) ([]models.ComplianceRequirementInstance, int64, error) {
	status, err := normalizeComplianceStatusFilter(status)
	if err != nil {
		return nil, 0, err
	}
	query := s.db.WithContext(ctx).Model(&models.ComplianceRequirementInstance{})
	if orgID != "" {
		query = query.Where("org_id = ?", orgID)
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if level != nil {
		query = query.Where("org_id IN (?)", s.db.WithContext(ctx).Model(&coreModels.Organization{}).Select("id").Where("level = ?", *level))
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var instances []models.ComplianceRequirementInstance
	if err := query.Order("created_at desc").Offset((page - 1) * limit).Limit(limit).Find(&instances).Error; err != nil {
		return nil, 0, err
	}
	return instances, total, nil
}

// GetRequirementDetail is deliberately read-only. Workflow transitions must be
// explicit mutations so browser retries and prefetching cannot change state.
func (s *ComplianceRequirementService) GetRequirementDetail(ctx context.Context, id string) (*models.ComplianceRequirementInstance, error) {
	var instance models.ComplianceRequirementInstance
	if err := s.db.WithContext(ctx).Where("id = ?", id).First(&instance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewHTTPError(http.StatusNotFound, "compliance requirement not found")
		}
		return nil, err
	}
	return &instance, nil
}

func (s *ComplianceRequirementService) DownloadSubmissionDocument(ctx context.Context, requirementID, documentID string) (*DocumentDownload, error) {
	if s.documents == nil {
		return nil, NewHTTPError(http.StatusServiceUnavailable, "document service is not configured")
	}
	instance, err := s.GetRequirementDetail(ctx, requirementID)
	if err != nil {
		return nil, err
	}
	linkedID, ok := instance.SubmissionData["document_id"].(string)
	if !ok || strings.TrimSpace(linkedID) == "" || strings.TrimSpace(linkedID) != strings.TrimSpace(documentID) {
		return nil, NewHTTPError(http.StatusNotFound, "document is not attached to this compliance requirement")
	}
	return s.documents.DownloadOwnedByOrganization(ctx, instance.OrgID, linkedID, models.DocumentCategoryComplianceRequirement)
}

// StartReview explicitly and idempotently moves a submitted requirement into
// review while preserving the admin actor in the audit trail.
func (s *ComplianceRequirementService) StartReview(ctx context.Context, actor AuthContext, id string) (*models.ComplianceRequirementInstance, error) {
	var updated models.ComplianceRequirementInstance
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var instance models.ComplianceRequirementInstance
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id).First(&instance).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return NewHTTPError(http.StatusNotFound, "compliance requirement not found")
			}
			return err
		}
		if instance.Status == models.ComplianceRequirementStatusUnderReview {
			updated = instance
			return nil
		}
		if instance.Status != models.ComplianceRequirementStatusSubmitted {
			return NewHTTPError(http.StatusConflict, "only a submitted requirement can enter review")
		}
		before := instance.Status
		instance.Status = models.ComplianceRequirementStatusUnderReview
		instance.UpdatedAt = time.Now()
		if err := tx.Save(&instance).Error; err != nil {
			return err
		}
		if s.audit != nil {
			event := AuditEventFromAuth(actor, models.ActionComplianceRequirementReviewStart, models.EntityComplianceRequirementInstance, instance.ID)
			event.BeforeState = models.JSONMap{"status": before}
			event.AfterState = models.JSONMap{"status": instance.Status}
			if err := s.audit.Record(ctx, tx, event); err != nil {
				return err
			}
		}
		updated = instance
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &updated, nil
}

func (s *ComplianceRequirementService) ReviewRequirement(ctx context.Context, actor AuthContext, id string, req models.ReviewComplianceRequirementRequest) (*models.ComplianceRequirementInstance, error) {
	decision := strings.TrimSpace(req.Decision)
	if decision != models.ComplianceRequirementStatusApproved && decision != models.ComplianceRequirementStatusRejected {
		return nil, NewHTTPError(http.StatusBadRequest, "decision must be 'approved' or 'rejected'")
	}
	rejectionReason := strings.TrimSpace(req.RejectionReason)
	if decision == models.ComplianceRequirementStatusRejected && rejectionReason == "" {
		return nil, NewHTTPError(http.StatusBadRequest, "rejection_reason is required when rejecting")
	}

	var updated models.ComplianceRequirementInstance
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var instance models.ComplianceRequirementInstance
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id).First(&instance).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return NewHTTPError(http.StatusNotFound, "compliance requirement not found")
			}
			return err
		}
		if instance.Status != models.ComplianceRequirementStatusSubmitted && instance.Status != models.ComplianceRequirementStatusUnderReview {
			return NewHTTPError(http.StatusConflict, "requirement is not awaiting review")
		}
		before := instance.Status
		now := time.Now()
		adminID := actor.MemberID
		instance.ReviewedByAdminID = &adminID
		instance.ReviewedAt = &now
		instance.UpdatedAt = now

		if decision == models.ComplianceRequirementStatusApproved {
			instance.Status = models.ComplianceRequirementStatusApproved
			instance.CompletedAt = &now
			instance.CompletedByMemberID = instance.SubmittedByMemberID
			instance.RejectionReason = nil
			var org coreModels.Organization
			if err := tx.Select("level").Where("id = ?", instance.OrgID).First(&org).Error; err != nil {
				return err
			}
			instance.AssignedLevelAtApproval = org.Level
		} else {
			instance.Status = models.ComplianceRequirementStatusRejected
			instance.RejectionReason = &rejectionReason
			instance.CompletedAt = nil
			instance.CompletedByMemberID = nil
		}

		if err := tx.Save(&instance).Error; err != nil {
			return err
		}
		updated = instance

		if s.audit != nil {
			event := AuditEventFromAuth(actor, models.ActionComplianceRequirementReview, models.EntityComplianceRequirementInstance, instance.ID)
			event.BeforeState = models.JSONMap{"status": before}
			event.AfterState = models.JSONMap{"status": instance.Status, "rejection_reason": instance.RejectionReason}
			if err := s.audit.Record(ctx, tx, event); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &updated, nil
}

// CreateAdHocRequirement assigns a one-off organization-level requirement.
// The legacy asset/custodian compliance API remains backed by ComplianceItem;
// keeping these contracts separate avoids losing asset association and evidence.
func (s *ComplianceRequirementService) CreateAdHocRequirement(ctx context.Context, actor AuthContext, req models.CreateComplianceRequirementRequest) (*models.ComplianceRequirementInstance, error) {
	orgID := strings.TrimSpace(req.OrgID)
	category := strings.TrimSpace(req.Category)
	requirement := strings.TrimSpace(req.Requirement)
	if orgID == "" || category == "" || requirement == "" {
		return nil, NewHTTPError(http.StatusBadRequest, "org_id, category, and requirement are required")
	}
	inputType := strings.ToLower(strings.TrimSpace(req.InputType))
	if !ContainsString(complianceInputTypes, inputType) {
		return nil, NewHTTPError(http.StatusBadRequest, "invalid input_type for requirement")
	}

	dueDate, err := ParseDate(req.DueDate)
	if err != nil {
		return nil, err
	}

	var organization coreModels.Organization
	if err := s.db.WithContext(ctx).First(&organization, "id = ?", orgID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewHTTPError(http.StatusNotFound, "organization not found")
		}
		return nil, err
	}
	if organization.Status != coreModels.OrganizationStatusActive {
		return nil, NewHTTPError(http.StatusUnprocessableEntity, "organization is not active")
	}

	now := time.Now()
	instance := models.ComplianceRequirementInstance{
		ID:          uuid.NewString(),
		OrgID:       organization.ID,
		Category:    category,
		Requirement: requirement,
		Status:      models.ComplianceRequirementStatusPending,
		Description: strings.TrimSpace(req.Description),
		Required:    complianceItemRequired(req.Required),
		InputType:   inputType,
		DueDate:     dueDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		duplicateQuery := tx.Model(&models.ComplianceRequirementInstance{}).
			Where("org_id = ? AND template_item_id IS NULL AND LOWER(category) = ? AND LOWER(requirement) = ?", orgID, strings.ToLower(category), strings.ToLower(requirement))
		if dueDate == nil {
			duplicateQuery = duplicateQuery.Where("due_date IS NULL")
		} else {
			duplicateQuery = duplicateQuery.Where("due_date = ?", *dueDate)
		}
		var duplicateCount int64
		if err := duplicateQuery.Count(&duplicateCount).Error; err != nil {
			return err
		}
		if duplicateCount > 0 {
			return NewHTTPError(http.StatusConflict, "compliance requirement already exists for this organization and due date")
		}
		if err := tx.Create(&instance).Error; err != nil {
			return err
		}
		if s.audit != nil {
			event := AuditEventFromAuth(actor, models.ActionComplianceRequirementCreate, models.EntityComplianceRequirementInstance, instance.ID)
			event.AfterState = models.JSONMap{"org_id": organization.ID, "category": instance.Category, "requirement": instance.Requirement, "status": instance.Status}
			if err := s.audit.Record(ctx, tx, event); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &instance, nil
}

// ---- Organisation-side ----

func (s *ComplianceRequirementService) ListOrgRequirements(ctx context.Context, auth AuthContext, page, limit int, status string) ([]models.ComplianceRequirementInstance, int64, error) {
	status, err := normalizeComplianceStatusFilter(status)
	if err != nil {
		return nil, 0, err
	}
	query := s.db.WithContext(ctx).Model(&models.ComplianceRequirementInstance{}).Where("org_id = ?", auth.OrganizationID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var instances []models.ComplianceRequirementInstance
	if err := query.Order("created_at desc").Offset((page - 1) * limit).Limit(limit).Find(&instances).Error; err != nil {
		return nil, 0, err
	}
	return instances, total, nil
}

var submittableComplianceStatuses = []string{
	models.ComplianceRequirementStatusDraft,
	models.ComplianceRequirementStatusPending,
	models.ComplianceRequirementStatusRejected,
}

func (s *ComplianceRequirementService) validateSubmissionData(ctx context.Context, auth AuthContext, instance *models.ComplianceRequirementInstance, data map[string]interface{}) error {
	switch instance.InputType {
	case models.ComplianceInputTypeDocumentUpload:
		documentID, ok := data["document_id"].(string)
		if !ok || strings.TrimSpace(documentID) == "" {
			return NewHTTPError(http.StatusBadRequest, "submission_data.document_id is required for document_upload requirements")
		}
		if s.documents == nil {
			return NewHTTPError(http.StatusServiceUnavailable, "document service is not configured")
		}
		return s.documents.ValidateForUse(ctx, auth, "", models.DocumentCategoryComplianceRequirement, []string{documentID})
	case models.ComplianceInputTypeText:
		value, ok := data["text"].(string)
		if !ok || strings.TrimSpace(value) == "" {
			return NewHTTPError(http.StatusBadRequest, "submission_data.text is required for text requirements")
		}
		return nil
	case models.ComplianceInputTypeStructuredForm:
		fields, ok := data["fields"].(map[string]interface{})
		if !ok || len(fields) == 0 {
			return NewHTTPError(http.StatusBadRequest, "submission_data.fields must be a non-empty object for structured_form requirements")
		}
		return nil
	default:
		return NewHTTPError(http.StatusUnprocessableEntity, "requirement has an unsupported input_type")
	}
}

func (s *ComplianceRequirementService) SubmitOrgRequirement(ctx context.Context, auth AuthContext, itemID string, req models.SubmitComplianceRequirementRequest) (*models.ComplianceRequirementInstance, error) {
	if len(req.SubmissionData) == 0 {
		return nil, NewHTTPError(http.StatusBadRequest, "submission_data is required")
	}
	var updated models.ComplianceRequirementInstance
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var instance models.ComplianceRequirementInstance
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND org_id = ?", itemID, auth.OrganizationID).First(&instance).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return NewHTTPError(http.StatusNotFound, "compliance requirement not found")
			}
			return err
		}
		if !ContainsString(submittableComplianceStatuses, instance.Status) {
			return NewHTTPError(http.StatusConflict, "requirement cannot be submitted from its current status")
		}
		if err := s.validateSubmissionData(ctx, auth, &instance, req.SubmissionData); err != nil {
			return err
		}
		now := time.Now()
		memberID := auth.MemberID
		instance.Status = models.ComplianceRequirementStatusSubmitted
		instance.SubmissionData = models.JSONMap(req.SubmissionData)
		instance.SubmittedByMemberID = &memberID
		instance.SubmittedAt = &now
		instance.RejectionReason = nil
		instance.ReviewedByAdminID = nil
		instance.ReviewedAt = nil
		instance.CompletedAt = nil
		instance.CompletedByMemberID = nil
		instance.AssignedLevelAtApproval = nil
		instance.UpdatedAt = now
		if err := tx.Save(&instance).Error; err != nil {
			return err
		}
		updated = instance
		if s.audit != nil {
			event := AuditEventFromAuth(auth, models.ActionComplianceRequirementSubmit, models.EntityComplianceRequirementInstance, instance.ID)
			event.AfterState = models.JSONMap{"status": instance.Status}
			if err := s.audit.Record(ctx, tx, event); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &updated, nil
}

// UpdateOrgRequirementStatus is intentionally narrow: it only allows an
// organisation to mark a pending item as "draft" (in progress). It cannot be
// used to bypass submit/review — those go through SubmitOrgRequirement and
// ReviewRequirement respectively.
func (s *ComplianceRequirementService) UpdateOrgRequirementStatus(ctx context.Context, auth AuthContext, itemID string, req models.UpdateOrgComplianceRequirementRequest) (*models.ComplianceRequirementInstance, error) {
	status := strings.TrimSpace(req.Status)
	if status != models.ComplianceRequirementStatusDraft {
		return nil, NewHTTPError(http.StatusBadRequest, "status can only be set to 'draft' via this endpoint")
	}
	var updated models.ComplianceRequirementInstance
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var instance models.ComplianceRequirementInstance
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND org_id = ?", itemID, auth.OrganizationID).First(&instance).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return NewHTTPError(http.StatusNotFound, "compliance requirement not found")
			}
			return err
		}
		if instance.Status != models.ComplianceRequirementStatusPending {
			return NewHTTPError(http.StatusConflict, "only a pending requirement can be marked as draft")
		}
		before := instance.Status
		instance.Status = status
		instance.UpdatedAt = time.Now()
		if err := tx.Save(&instance).Error; err != nil {
			return err
		}
		if s.audit != nil {
			event := AuditEventFromAuth(auth, models.ActionComplianceRequirementDraft, models.EntityComplianceRequirementInstance, instance.ID)
			event.BeforeState = models.JSONMap{"status": before}
			event.AfterState = models.JSONMap{"status": instance.Status}
			if err := s.audit.Record(ctx, tx, event); err != nil {
				return err
			}
		}
		updated = instance
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &updated, nil
}

// ---- Auto-assignment ----

// AutoAssignComplianceRequirements resolves the template for (orgType, level)
// and creates a ComplianceRequirementInstance for each of its items that the
// org doesn't already have an instance for. It never modifies or deletes
// existing instances, which is what keeps already-approved items untouched
// across a level change (PRD §6.2). A missing template is a no-op, not an
// error — not every (org_type, level) combination need be templated yet.
//
// Exported as a package-level function (not a method) so the `organizations`
// component can call it directly from the org-creation and org-update
// handlers without needing a fully-wired service instance.
func AutoAssignComplianceRequirements(ctx context.Context, db *gorm.DB, orgID, orgType string, level int) error {
	orgID = strings.TrimSpace(orgID)
	if orgID == "" || strings.TrimSpace(orgType) == "" || level <= 0 {
		return nil
	}
	normalizedOrgType, err := normalizeComplianceOrganizationType(orgType)
	if err != nil {
		return err
	}

	var template models.ComplianceRequirementTemplate
	err = db.WithContext(ctx).Where("org_type = ? AND level = ?", normalizedOrgType, level).
		Preload("Items").First(&template).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	if len(template.Items) == 0 {
		return nil
	}

	itemIDs := make([]string, 0, len(template.Items))
	for _, item := range template.Items {
		itemIDs = append(itemIDs, item.ID)
	}
	var existing []models.ComplianceRequirementInstance
	if err := db.WithContext(ctx).Model(&models.ComplianceRequirementInstance{}).
		Where("org_id = ? AND template_item_id IN ?", orgID, itemIDs).
		Find(&existing).Error; err != nil {
		return err
	}
	assigned := make(map[string]bool, len(existing))
	for _, instance := range existing {
		if instance.TemplateItemID != nil {
			assigned[*instance.TemplateItemID] = true
		}
	}

	now := time.Now()
	version := template.Version
	toCreate := make([]models.ComplianceRequirementInstance, 0, len(template.Items))
	for _, item := range template.Items {
		if assigned[item.ID] {
			continue
		}
		itemID := item.ID
		toCreate = append(toCreate, models.ComplianceRequirementInstance{
			ID:              uuid.NewString(),
			OrgID:           orgID,
			Category:        item.Category,
			Requirement:     item.Name,
			Description:     item.Description,
			Required:        item.Required,
			TemplateItemID:  &itemID,
			TemplateVersion: &version,
			Status:          models.ComplianceRequirementStatusPending,
			InputType:       item.InputType,
			CreatedAt:       now,
			UpdatedAt:       now,
		})
	}
	if len(toCreate) == 0 {
		return nil
	}
	return db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "org_id"}, {Name: "template_item_id"}},
		DoNothing: true,
	}).Create(&toCreate).Error
}

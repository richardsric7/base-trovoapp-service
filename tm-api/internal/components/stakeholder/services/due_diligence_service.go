package services

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"admin-panel-dashboard/internal/components/stakeholder/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DueDiligenceService struct {
	db            *gorm.DB
	assets        *AssetService
	documents     *DocumentService
	audit         *AuditService
	notifications *NotificationService
}

func NewDueDiligenceService(db *gorm.DB, assets *AssetService, documents *DocumentService, audit *AuditService, notifications *NotificationService) *DueDiligenceService {
	return &DueDiligenceService{db: db, assets: assets, documents: documents, audit: audit, notifications: notifications}
}

func (s *DueDiligenceService) GetChecklist(ctx context.Context, auth AuthContext, assetID string) (*models.DueDiligenceChecklist, error) {
	if auth.DashboardRole != models.DashboardRoleTrustee {
		return nil, NewHTTPError(http.StatusForbidden, "trustee role is required")
	}
	assignment, asset, err := s.assets.RequireAssignmentForAsset(ctx, auth, assetID)
	if err != nil {
		return nil, err
	}
	if assignment.TrusteeOrgID == nil || *assignment.TrusteeOrgID != auth.OrganizationID {
		return nil, NewHTTPError(http.StatusNotFound, "asset not found")
	}
	checklist, err := s.ensureChecklist(ctx, auth, asset.ID, asset.AssetCode)
	if err != nil {
		return nil, err
	}
	if err := s.attachItemDocuments(ctx, auth, checklist.Items); err != nil {
		return nil, err
	}
	assetDocuments, err := s.assets.WalletAssetDocuments(ctx, auth, asset.ID)
	if err != nil {
		return nil, err
	}
	checklist.AvailableAssetDocuments = assetDocuments
	return checklist, nil
}

func (s *DueDiligenceService) UpdateItem(ctx context.Context, auth AuthContext, assetID, itemID string, req models.UpdateDueDiligenceItemRequest) (*models.DueDiligenceItem, error) {
	checklist, err := s.GetChecklist(ctx, auth, assetID)
	if err != nil {
		return nil, err
	}
	status := strings.TrimSpace(req.Status)
	if !ContainsString([]string{models.DueDiligenceItemStatusPending, models.DueDiligenceItemStatusComplete, models.DueDiligenceItemStatusFailed, models.DueDiligenceItemStatusNotApplicable}, status) {
		return nil, NewHTTPError(http.StatusBadRequest, "invalid due diligence item status")
	}
	if req.DocumentIDs != nil {
		if s.documents == nil {
			return nil, NewHTTPError(http.StatusServiceUnavailable, "document service is not configured")
		}
		if err := s.documents.ValidateAccessibleForUse(ctx, auth, checklist.AssetID, models.DocumentCategoryDueDiligence, *req.DocumentIDs); err != nil {
			return nil, err
		}
	}
	var updated models.DueDiligenceItem
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var item models.DueDiligenceItem
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND checklist_id = ?", itemID, checklist.ID).First(&item).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return NewHTTPError(http.StatusNotFound, "due diligence item not found")
			}
			return err
		}
		before := item.Status
		now := time.Now()
		item.Status = status
		item.Notes = strings.TrimSpace(req.Notes)
		if status == models.DueDiligenceItemStatusComplete || status == models.DueDiligenceItemStatusFailed || status == models.DueDiligenceItemStatusNotApplicable {
			item.VerifiedByMemberID = &auth.MemberID
			item.VerifiedAt = &now
		}
		item.UpdatedAt = now
		if err := tx.Save(&item).Error; err != nil {
			return err
		}
		if req.DocumentIDs != nil {
			if err := replaceDueDiligenceItemDocuments(tx, item.ID, *req.DocumentIDs, now); err != nil {
				return err
			}
		}
		updated = item
		if s.audit != nil {
			event := AuditEventFromAuth(auth, models.ActionDueDiligenceItemUpdate, models.EntityDueDiligenceItem, item.ID)
			event.BeforeState = models.JSONMap{"status": before}
			event.AfterState = models.JSONMap{"status": item.Status, "notes": item.Notes, "document_ids": req.DocumentIDs}
			if err := s.audit.Record(ctx, tx, event); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	items := []models.DueDiligenceItem{updated}
	if err := s.attachItemDocuments(ctx, auth, items); err != nil {
		return nil, err
	}
	updated = items[0]
	return &updated, nil
}

func (s *DueDiligenceService) UpdateCategory(ctx context.Context, auth AuthContext, assetID string, req models.UpdateDueDiligenceCategoryRequest) (*models.DueDiligenceCategoryResponse, error) {
	category := strings.ToLower(strings.TrimSpace(req.Category))
	if category == "" {
		return nil, NewHTTPError(http.StatusBadRequest, "category is required")
	}
	status := strings.ToLower(strings.TrimSpace(req.Status))
	if status == "verified" {
		status = models.DueDiligenceItemStatusComplete
	}
	if !ContainsString([]string{models.DueDiligenceItemStatusPending, models.DueDiligenceItemStatusComplete, models.DueDiligenceItemStatusFailed, models.DueDiligenceItemStatusNotApplicable}, status) {
		return nil, NewHTTPError(http.StatusBadRequest, "invalid due diligence category status")
	}
	checklist, err := s.GetChecklist(ctx, auth, assetID)
	if err != nil {
		return nil, err
	}
	notes := strings.TrimSpace(req.Notes)
	var response models.DueDiligenceCategoryResponse
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var lockedChecklist models.DueDiligenceChecklist
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND trustee_org_id = ?", checklist.ID, auth.OrganizationID).
			First(&lockedChecklist).Error; err != nil {
			return err
		}
		if lockedChecklist.Status == models.DueDiligenceStatusApproved || lockedChecklist.Status == models.DueDiligenceStatusRejected {
			return NewHTTPError(http.StatusConflict, "due diligence checklist is already finalized")
		}

		var items []models.DueDiligenceItem
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("checklist_id = ? AND LOWER(category) = ?", lockedChecklist.ID, category).
			Order("created_at asc").Find(&items).Error; err != nil {
			return err
		}
		if len(items) == 0 {
			return NewHTTPError(http.StatusNotFound, "due diligence category not found")
		}

		itemIDs := make([]string, 0, len(items))
		beforeStatuses := make(map[string]string, len(items))
		for _, item := range items {
			itemIDs = append(itemIDs, item.ID)
			beforeStatuses[item.ID] = item.Status
		}
		now := time.Now()
		updates := map[string]interface{}{
			"status": status, "notes": notes, "updated_at": now,
			"verified_by_member_id": nil, "verified_at": nil,
		}
		if status != models.DueDiligenceItemStatusPending {
			updates["verified_by_member_id"] = auth.MemberID
			updates["verified_at"] = now
		}
		if err := tx.Model(&models.DueDiligenceItem{}).Where("id IN ?", itemIDs).Updates(updates).Error; err != nil {
			return err
		}
		if err := tx.Where("id IN ?", itemIDs).Order("created_at asc").Find(&items).Error; err != nil {
			return err
		}

		response = models.DueDiligenceCategoryResponse{Category: category, Status: status, Notes: notes, Items: items}
		if status != models.DueDiligenceItemStatusPending {
			memberID := auth.MemberID
			verifiedAt := now
			response.VerifiedByMemberID = &memberID
			response.VerifiedAt = &verifiedAt
		}
		if s.audit != nil {
			event := AuditEventFromAuth(auth, models.ActionDueDiligenceCategoryUpdate, models.EntityDueDiligenceCategory, lockedChecklist.ID+":"+category)
			event.BeforeState = models.JSONMap{"category": category, "item_statuses": beforeStatuses}
			event.AfterState = models.JSONMap{"category": category, "status": status, "notes": notes, "item_count": len(items)}
			if err := s.audit.Record(ctx, tx, event); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if err := s.attachItemDocuments(ctx, auth, response.Items); err != nil {
		return nil, err
	}
	return &response, nil
}

func replaceDueDiligenceItemDocuments(tx *gorm.DB, itemID string, documentIDs []string, now time.Time) error {
	if err := tx.Where("due_diligence_item_id = ?", itemID).Delete(&models.DueDiligenceItemDocument{}).Error; err != nil {
		return err
	}
	links := make([]models.DueDiligenceItemDocument, 0, len(documentIDs))
	for _, documentID := range documentIDs {
		links = append(links, models.DueDiligenceItemDocument{DueDiligenceItemID: itemID, DocumentID: strings.TrimSpace(documentID), CreatedAt: now})
	}
	if len(links) > 0 {
		return tx.Create(&links).Error
	}
	return nil
}

func (s *DueDiligenceService) attachItemDocuments(ctx context.Context, auth AuthContext, items []models.DueDiligenceItem) error {
	if len(items) == 0 {
		return nil
	}
	itemIDs := make([]string, 0, len(items))
	for i := range items {
		items[i].Documents = []models.StakeholderDocument{}
		itemIDs = append(itemIDs, items[i].ID)
	}
	var links []models.DueDiligenceItemDocument
	if err := s.db.WithContext(ctx).Where("due_diligence_item_id IN ?", itemIDs).Order("created_at asc").Find(&links).Error; err != nil {
		return err
	}
	if len(links) == 0 {
		return nil
	}
	documentIDs := make([]string, 0, len(links))
	for _, link := range links {
		documentIDs = append(documentIDs, link.DocumentID)
	}
	var documents []models.StakeholderDocument
	if err := s.db.WithContext(ctx).Where("id IN ? AND status = ?", documentIDs, models.DocumentStatusActive).Find(&documents).Error; err != nil {
		return err
	}
	documentByID := make(map[string]models.StakeholderDocument, len(documents))
	for _, document := range documents {
		if models.RoleCanAccessDocument(auth.DashboardRole, []string(document.AccessRoles)) {
			documentByID[document.ID] = document
		}
	}
	itemIndex := make(map[string]int, len(items))
	for i := range items {
		itemIndex[items[i].ID] = i
	}
	for _, link := range links {
		if document, ok := documentByID[link.DocumentID]; ok {
			if index, exists := itemIndex[link.DueDiligenceItemID]; exists {
				items[index].Documents = append(items[index].Documents, document)
			}
		}
	}
	return nil
}

func (s *DueDiligenceService) Approve(ctx context.Context, auth AuthContext, assetID string) (*models.DueDiligenceChecklist, error) {
	return s.finalize(ctx, auth, assetID, models.DueDiligenceStatusApproved, "")
}

func (s *DueDiligenceService) Reject(ctx context.Context, auth AuthContext, assetID string, reason string) (*models.DueDiligenceChecklist, error) {
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return nil, NewHTTPError(http.StatusBadRequest, "rejection reason is required")
	}
	return s.finalize(ctx, auth, assetID, models.DueDiligenceStatusRejected, reason)
}

func (s *DueDiligenceService) finalize(ctx context.Context, auth AuthContext, assetID, status, reason string) (*models.DueDiligenceChecklist, error) {
	checklist, err := s.GetChecklist(ctx, auth, assetID)
	if err != nil {
		return nil, err
	}
	var updated models.DueDiligenceChecklist
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var locked models.DueDiligenceChecklist
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND trustee_org_id = ?", checklist.ID, auth.OrganizationID).First(&locked).Error; err != nil {
			return err
		}
		if locked.Status == models.DueDiligenceStatusApproved || locked.Status == models.DueDiligenceStatusRejected {
			return NewHTTPError(http.StatusConflict, "due diligence checklist is already finalized")
		}
		before := locked.Status
		now := time.Now()
		locked.Status = status
		locked.UpdatedAt = now
		if status == models.DueDiligenceStatusApproved {
			locked.ApprovedByMemberID = &auth.MemberID
			locked.ApprovedAt = &now
		} else {
			locked.RejectedByMemberID = &auth.MemberID
			locked.RejectedAt = &now
			locked.RejectionReason = reason
		}
		if err := tx.Save(&locked).Error; err != nil {
			return err
		}
		updated = locked
		notificationType := "due_diligence.approved"
		title := "Due diligence approved"
		if status == models.DueDiligenceStatusRejected {
			notificationType = "due_diligence.rejected"
			title = "Due diligence rejected"
		}
		if s.notifications != nil {
			if assignment, err := s.assets.assignmentForAsset(ctx, locked.AssetID); err == nil {
				recipients := []string{}
				if assignment.AssetManagerOrgID != nil {
					recipients = append(recipients, *assignment.AssetManagerOrgID)
				}
				if assignment.CustodianOrgID != nil {
					recipients = append(recipients, *assignment.CustodianOrgID)
				}
				for _, recipient := range recipients {
					if recipient == "" || recipient == auth.OrganizationID {
						continue
					}
					if err := s.notifications.Create(ctx, tx, NotificationInput{
						RecipientOrgID:    recipient,
						SenderOrgID:       &auth.OrganizationID,
						Type:              notificationType,
						Title:             title,
						Message:           "Trustee due diligence status changed.",
						RelatedEntityType: models.EntityDueDiligenceChecklist,
						RelatedEntityID:   locked.ID,
					}); err != nil {
						return err
					}
				}
			}
		}
		if s.audit != nil {
			action := models.ActionDueDiligenceApprove
			if status == models.DueDiligenceStatusRejected {
				action = models.ActionDueDiligenceReject
			}
			event := AuditEventFromAuth(auth, action, models.EntityDueDiligenceChecklist, locked.ID)
			event.BeforeState = models.JSONMap{"status": before}
			event.AfterState = models.JSONMap{"status": locked.Status, "reason": reason}
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

func (s *DueDiligenceService) ensureChecklist(ctx context.Context, auth AuthContext, assetID, assetCode string) (*models.DueDiligenceChecklist, error) {
	var checklist models.DueDiligenceChecklist
	err := s.db.WithContext(ctx).Preload("Items").Where("asset_id = ? AND trustee_org_id = ?", assetID, auth.OrganizationID).First(&checklist).Error
	if err == nil {
		return &checklist, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	now := time.Now()
	checklist = models.DueDiligenceChecklist{
		ID:           uuid.NewString(),
		AssetID:      assetID,
		AssetCode:    assetCode,
		TrusteeOrgID: auth.OrganizationID,
		Status:       models.DueDiligenceStatusInReview,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	defaults := defaultDueDiligenceItems(checklist.ID, now)
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&checklist).Error; err != nil {
			return err
		}
		if len(defaults) > 0 {
			if err := tx.Create(&defaults).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	checklist.Items = defaults
	return &checklist, nil
}

func defaultDueDiligenceItems(checklistID string, now time.Time) []models.DueDiligenceItem {
	items := []struct {
		category string
		item     string
	}{
		{"legal", "Verify asset ownership documentation"},
		{"legal", "Verify no lien or encumbrance evidence"},
		{"financial", "Review asset valuation evidence"},
		{"operational", "Review asset protection and insurance evidence"},
	}
	result := make([]models.DueDiligenceItem, 0, len(items))
	for _, item := range items {
		result = append(result, models.DueDiligenceItem{
			ID:          uuid.NewString(),
			ChecklistID: checklistID,
			Category:    item.category,
			Item:        item.item,
			Status:      models.DueDiligenceItemStatusPending,
			CreatedAt:   now,
			UpdatedAt:   now,
		})
	}
	return result
}

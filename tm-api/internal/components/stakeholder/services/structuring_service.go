package services

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	models "admin-panel-dashboard/internal/components/stakeholder/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// StructuringService handles the Tier-1 A5 "confirm workstream complete" action
// for Legal and Financial Advisers (PRD §5A.3 / OI-14).
type StructuringService struct {
	db        *gorm.DB
	assets    *AssetService
	documents *DocumentService
	audit     *AuditService
}

func NewStructuringService(db *gorm.DB, assets *AssetService, documents *DocumentService, audit *AuditService) *StructuringService {
	return &StructuringService{db: db, assets: assets, documents: documents, audit: audit}
}

// workstreamForRole maps an adviser dashboard role to its A5 workstream.
func workstreamForRole(role string) (string, bool) {
	switch role {
	case models.DashboardRoleLegalAdviser:
		return models.StructuringWorkstreamLegal, true
	case models.DashboardRoleFinancialAdviser:
		return models.StructuringWorkstreamFinancial, true
	default:
		return "", false
	}
}

// ConfirmComplete records that the calling adviser's A5 workstream is complete
// for the given asset. Only Legal/Financial Advisers may call it, and only for
// assets assigned to them. Re-confirming an already-complete workstream is a 409.
func (s *StructuringService) ConfirmComplete(ctx context.Context, auth AuthContext, assetIDOrCode string, req models.ConfirmStructuringRequest) (*models.StakeholderStructuringStatus, error) {
	workstream, ok := workstreamForRole(auth.DashboardRole)
	if !ok {
		return nil, NewHTTPError(http.StatusForbidden, "legal or financial adviser role is required")
	}

	// Ownership: GetAsset enforces that the adviser is assigned to this asset
	// (via the FK-only visibility path). Returns 404 otherwise.
	asset, err := s.assets.GetAsset(ctx, auth, assetIDOrCode)
	if err != nil {
		return nil, err
	}

	// Optionally attach an evidence document.
	documentID := req.DocumentID
	if documentID == nil && strings.TrimSpace(req.FileURL) != "" {
		title := strings.TrimSpace(req.Title)
		if title == "" {
			title = "A5 structuring evidence"
		}
		doc, docErr := s.documents.Create(ctx, auth, models.CreateDocumentRequest{
			AssetID:   asset.ID,
			AssetCode: asset.AssetCode,
			Category:  "structuring",
			Title:     title,
			FileURL:   req.FileURL,
		})
		if docErr != nil {
			return nil, docErr
		}
		documentID = &doc.ID
	}

	var status models.StakeholderStructuringStatus
	now := time.Now()
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		findErr := tx.Where("asset_id = ? AND workstream = ?", asset.ID, workstream).First(&status).Error
		switch {
		case findErr == nil:
			if status.Status == models.StructuringStatusComplete {
				return NewHTTPError(http.StatusConflict, "structuring workstream already confirmed complete")
			}
		case errors.Is(findErr, gorm.ErrRecordNotFound):
			status = models.StakeholderStructuringStatus{
				ID:         uuid.NewString(),
				AssetID:    asset.ID,
				AssetCode:  asset.AssetCode,
				Workstream: workstream,
				OrgID:      auth.OrganizationID,
				CreatedAt:  now,
			}
		default:
			return findErr
		}

		before := status.Status
		status.Status = models.StructuringStatusComplete
		status.OrgID = auth.OrganizationID
		status.StakeholderID = auth.StakeholderID
		status.DocumentID = documentID
		status.Notes = req.Notes
		status.ConfirmedByMemberID = &auth.MemberID
		status.ConfirmedAt = &now
		status.UpdatedAt = now
		if err := tx.Save(&status).Error; err != nil {
			return err
		}

		if s.audit != nil {
			event := AuditEventFromAuth(auth, models.ActionStructuringComplete, models.EntityStructuringStatus, status.ID)
			event.BeforeState = models.JSONMap{"status": before}
			event.AfterState = models.JSONMap{"status": status.Status, "workstream": workstream}
			if err := s.audit.Record(ctx, tx, event); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &status, nil
}

package services

import (
	"context"
	"strings"

	"admin-panel-dashboard/internal/components/stakeholder/models"

	"gorm.io/gorm"
)

type assetActivityScope struct {
	entityType string
	entityIDs  []string
}

// ListAssetActivities returns every recorded portal action tied to an asset.
// Visibility is checked against the canonical WalletDB assignment before any
// AdminDB audit data is returned.
func (s *AssetService) ListAssetActivities(ctx context.Context, auth AuthContext, idOrCode string, page, limit int) ([]models.StakeholderAuditLog, int64, error) {
	asset, _, err := s.visibleAsset(ctx, auth, idOrCode)
	if err != nil {
		return nil, 0, err
	}

	scopes, categoryPrefixes, err := s.assetActivityScopes(ctx, asset.ID)
	if err != nil {
		return nil, 0, err
	}
	query := s.adminDB.WithContext(ctx).Model(&models.StakeholderAuditLog{}).Where("1 = 0")
	for _, scope := range scopes {
		if len(scope.entityIDs) > 0 {
			query = query.Or("(entity_type = ? AND entity_id IN ?)", scope.entityType, scope.entityIDs)
		}
	}
	for _, checklistID := range categoryPrefixes {
		query = query.Or("(entity_type = ? AND entity_id LIKE ?)", models.EntityDueDiligenceCategory, checklistID+":%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	activities := make([]models.StakeholderAuditLog, 0)
	if err := query.Order("created_at desc, id desc").Offset((page - 1) * limit).Limit(limit).Find(&activities).Error; err != nil {
		return nil, 0, err
	}
	if err := s.enrichAssetActivityActors(ctx, activities); err != nil {
		return nil, 0, err
	}
	return activities, total, nil
}

func (s *AssetService) enrichAssetActivityActors(ctx context.Context, activities []models.StakeholderAuditLog) error {
	memberIDs := make([]string, 0, len(activities))
	organizationIDs := make([]string, 0, len(activities))
	for _, activity := range activities {
		memberIDs = append(memberIDs, activity.ActorMemberID)
		organizationIDs = append(organizationIDs, activity.ActorOrgID)
	}
	membersByID, organizationsByID, err := loadOrganizationActors(ctx, s.adminDB, memberIDs, organizationIDs)
	if err != nil {
		return err
	}
	for index := range activities {
		activity := &activities[index]
		if member, ok := membersByID[activity.ActorMemberID]; ok && (activity.ActorOrgID == "" || member.OrganizationID == activity.ActorOrgID) {
			activity.ActorName = organizationMemberDisplayName(member)
		}
		if organization, ok := organizationsByID[activity.ActorOrgID]; ok {
			activity.ActorOrganization = organization.Name
		}
	}
	return nil
}

func (s *AssetService) assetActivityScopes(ctx context.Context, assetID string) ([]assetActivityScope, []string, error) {
	assetID = strings.TrimSpace(assetID)
	scopes := []assetActivityScope{
		{entityType: "tokenized_asset", entityIDs: []string{assetID}},
		{entityType: "asset", entityIDs: []string{assetID}},
	}

	modelsByEntity := []struct {
		entityType string
		model      interface{}
	}{
		{models.EntityAssetAssignment, &models.StakeholderAssetAssignment{}},
		{models.EntityFundReleaseRequest, &models.FundReleaseRequest{}},
		{models.EntityRevenueRecord, &models.RevenueRecord{}},
		{models.EntityDistribution, &models.Distribution{}},
		{models.EntityAssetValuation, &models.AssetValuation{}},
		{models.EntitySegregatedAccount, &models.SegregatedAccount{}},
		{models.EntityComplianceItem, &models.ComplianceItem{}},
		{models.EntityDocument, &models.StakeholderDocument{}},
		{models.EntityAssetOperation, &models.StakeholderAssetOperation{}},
		{models.EntityReport, &models.StakeholderReport{}},
		{models.EntityStructuringStatus, &models.StakeholderStructuringStatus{}},
	}
	for _, candidate := range modelsByEntity {
		ids, err := pluckAssetEntityIDs(s.adminDB.WithContext(ctx), candidate.model, assetID)
		if err != nil {
			return nil, nil, err
		}
		scopes = append(scopes, assetActivityScope{entityType: candidate.entityType, entityIDs: ids})
	}

	checklistIDs, err := pluckAssetEntityIDs(s.adminDB.WithContext(ctx), &models.DueDiligenceChecklist{}, assetID)
	if err != nil {
		return nil, nil, err
	}
	scopes = append(scopes, assetActivityScope{entityType: models.EntityDueDiligenceChecklist, entityIDs: checklistIDs})
	if len(checklistIDs) > 0 {
		var itemIDs []string
		if err := s.adminDB.WithContext(ctx).Model(&models.DueDiligenceItem{}).
			Where("checklist_id IN ?", checklistIDs).Pluck("id", &itemIDs).Error; err != nil {
			return nil, nil, err
		}
		scopes = append(scopes, assetActivityScope{entityType: models.EntityDueDiligenceItem, entityIDs: itemIDs})
	}
	return scopes, checklistIDs, nil
}

func pluckAssetEntityIDs(db *gorm.DB, model interface{}, assetID string) ([]string, error) {
	var ids []string
	if err := db.Model(model).Where("asset_id = ?", assetID).Pluck("id", &ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

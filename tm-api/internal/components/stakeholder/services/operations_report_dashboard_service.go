package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"admin-panel-dashboard/internal/components/stakeholder/models"
	coreModels "admin-panel-dashboard/internal/models"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type OperationsService struct {
	db            *gorm.DB
	assets        *AssetService
	documents     *DocumentService
	audit         *AuditService
	notifications *NotificationService
}

func NewOperationsService(db *gorm.DB, assets *AssetService, documents *DocumentService, audit *AuditService, notifications *NotificationService) *OperationsService {
	return &OperationsService{db: db, assets: assets, documents: documents, audit: audit, notifications: notifications}
}

func (s *OperationsService) UpdateAssetOperation(ctx context.Context, auth AuthContext, assetID string, req models.UpdateAssetOperationRequest) (*models.StakeholderAssetOperation, error) {
	if auth.DashboardRole != models.DashboardRoleAssetManager {
		return nil, NewHTTPError(http.StatusForbidden, "asset manager role is required")
	}
	_, asset, err := s.assets.RequireAssignmentForAsset(ctx, auth, assetID)
	if err != nil {
		return nil, err
	}
	status := models.Normalize(req.OperationalStatus)
	if !ContainsString(models.AllowedAssetOperationalStatuses(), status) {
		return nil, NewHTTPError(http.StatusBadRequest, "invalid operational_status; allowed values are active, inactive, paused, maintenance")
	}
	metadata := models.JSONMap{}
	if effectiveAt := strings.TrimSpace(req.Metadata.EffectiveAt); effectiveAt != "" {
		if _, err := time.Parse(time.RFC3339, effectiveAt); err != nil {
			return nil, NewHTTPError(http.StatusBadRequest, "metadata.effective_at must be RFC3339")
		}
		metadata["effective_at"] = effectiveAt
	}
	if source := strings.TrimSpace(req.Metadata.Source); source != "" {
		if len(source) > 64 {
			return nil, NewHTTPError(http.StatusBadRequest, "metadata.source exceeds 64 characters")
		}
		metadata["source"] = source
	}
	var operation models.StakeholderAssetOperation
	now := time.Now()
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Where("asset_id = ? AND manager_org_id = ?", asset.ID, auth.OrganizationID).First(&operation).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			operation = models.StakeholderAssetOperation{
				ID:           uuid.NewString(),
				AssetID:      asset.ID,
				AssetCode:    asset.AssetCode,
				ManagerOrgID: auth.OrganizationID,
				CreatedAt:    now,
			}
		}
		before := operation.OperationalStatus
		operation.OperationalStatus = status
		operation.Notes = strings.TrimSpace(req.Notes)
		operation.Metadata = metadata
		operation.UpdatedByMemberID = auth.MemberID
		operation.UpdatedAt = now
		if err := tx.Save(&operation).Error; err != nil {
			return err
		}
		if s.audit != nil {
			event := AuditEventFromAuth(auth, models.ActionAssetOperationUpdate, models.EntityAssetOperation, operation.ID)
			event.BeforeState = models.JSONMap{"operational_status": before}
			event.AfterState = models.JSONMap{"operational_status": operation.OperationalStatus, "notes": operation.Notes}
			if err := s.audit.Record(ctx, tx, event); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &operation, nil
}

func (s *OperationsService) CreateReport(ctx context.Context, auth AuthContext, req models.CreateReportRequest) (*models.StakeholderReport, error) {
	if auth.DashboardRole != models.DashboardRoleAssetManager {
		return nil, NewHTTPError(http.StatusForbidden, "asset manager role is required")
	}
	reportType := models.Normalize(req.ReportType)
	if !ContainsString(models.AllowedReportTypes(), reportType) {
		return nil, NewHTTPError(http.StatusBadRequest, "invalid report_type; allowed values are income, operational")
	}
	var trusteeOrgID string
	assetID := strings.TrimSpace(req.AssetID)
	assetCode := strings.TrimSpace(req.AssetCode)
	if assetID != "" {
		assignment, asset, err := s.assets.RequireAssignmentForAsset(ctx, auth, assetID)
		if err != nil {
			return nil, err
		}
		assetID = asset.ID
		assetCode = asset.AssetCode
		if assignment.TrusteeOrgID != nil {
			trusteeOrgID = *assignment.TrusteeOrgID
		}
	}
	documentID := strings.TrimSpace(req.DocumentID)
	fileURL := strings.TrimSpace(req.FileURL)
	if documentID == "" && fileURL == "" {
		return nil, NewHTTPError(http.StatusBadRequest, "document_id is required")
	}
	if documentID != "" {
		if s.documents == nil {
			return nil, NewHTTPError(http.StatusServiceUnavailable, "document service is not configured")
		}
		if err := s.documents.ValidateForUse(ctx, auth, assetID, models.DocumentCategoryAssetReport, []string{documentID}); err != nil {
			return nil, err
		}
		fileURL = ""
	}
	now := time.Now()
	report := models.StakeholderReport{
		ID:                  uuid.NewString(),
		AssetID:             assetID,
		AssetCode:           assetCode,
		ManagerOrgID:        auth.OrganizationID,
		TrusteeOrgID:        trusteeOrgID,
		ReportType:          reportType,
		Title:               strings.TrimSpace(req.Title),
		FileURL:             fileURL,
		Status:              models.ReportStatusGenerated,
		GeneratedByMemberID: auth.MemberID,
		CreatedAt:           now,
		UpdatedAt:           now,
	}
	if documentID != "" {
		report.DocumentID = &documentID
	}
	if report.ReportType == "" || report.Title == "" {
		return nil, NewHTTPError(http.StatusBadRequest, "report_type and title are required")
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&report).Error; err != nil {
			return err
		}
		if s.audit != nil {
			event := AuditEventFromAuth(auth, models.ActionReportGenerate, models.EntityReport, report.ID)
			event.AfterState = models.JSONMap{"status": report.Status, "report_type": report.ReportType}
			if err := s.audit.Record(ctx, tx, event); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (s *OperationsService) ReportTypes(auth AuthContext) ([]models.ReportTypeResponse, error) {
	if auth.DashboardRole != models.DashboardRoleAssetManager {
		return nil, NewHTTPError(http.StatusForbidden, "asset manager role is required")
	}
	return []models.ReportTypeResponse{
		{Value: models.ReportTypeIncome, Label: "Income"},
		{Value: models.ReportTypeOperational, Label: "Operational"},
	}, nil
}

func (s *OperationsService) ListReports(ctx context.Context, auth AuthContext, page, limit int) ([]models.StakeholderReport, int64, error) {
	if auth.DashboardRole != models.DashboardRoleAssetManager {
		return nil, 0, NewHTTPError(http.StatusForbidden, "asset manager role is required")
	}
	query := s.db.WithContext(ctx).Model(&models.StakeholderReport{}).Where("manager_org_id = ?", auth.OrganizationID)
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var reports []models.StakeholderReport
	if err := query.Order("created_at desc").Offset((page - 1) * limit).Limit(limit).Find(&reports).Error; err != nil {
		return nil, 0, err
	}
	return reports, total, nil
}

func (s *OperationsService) SubmitReport(ctx context.Context, auth AuthContext, id string) (*models.StakeholderReport, error) {
	if auth.DashboardRole != models.DashboardRoleAssetManager {
		return nil, NewHTTPError(http.StatusForbidden, "asset manager role is required")
	}
	var report models.StakeholderReport
	now := time.Now()
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND manager_org_id = ?", id, auth.OrganizationID).First(&report).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return NewHTTPError(http.StatusNotFound, "report not found")
			}
			return err
		}
		if report.Status == models.ReportStatusSubmitted {
			return NewHTTPError(http.StatusConflict, "report already submitted")
		}
		before := report.Status
		report.Status = models.ReportStatusSubmitted
		report.SubmittedAt = &now
		report.SubmittedByMemberID = &auth.MemberID
		report.UpdatedAt = now
		if err := tx.Save(&report).Error; err != nil {
			return err
		}
		if s.notifications != nil && report.TrusteeOrgID != "" {
			if err := s.notifications.Create(ctx, tx, NotificationInput{
				RecipientOrgID:    report.TrusteeOrgID,
				SenderOrgID:       &auth.OrganizationID,
				Type:              "report.submitted",
				Title:             "Report submitted",
				Message:           "An asset manager report was submitted to the trustee.",
				RelatedEntityType: models.EntityReport,
				RelatedEntityID:   report.ID,
			}); err != nil {
				return err
			}
		}
		if s.audit != nil {
			event := AuditEventFromAuth(auth, models.ActionReportSubmit, models.EntityReport, report.ID)
			event.BeforeState = models.JSONMap{"status": before}
			event.AfterState = models.JSONMap{"status": report.Status}
			if err := s.audit.Record(ctx, tx, event); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &report, nil
}

type DashboardService struct {
	db     *gorm.DB
	assets *AssetService
}

func NewDashboardService(db *gorm.DB, assets *AssetService) *DashboardService {
	return &DashboardService{db: db, assets: assets}
}

func (s *DashboardService) Dashboard(ctx context.Context, auth AuthContext) (*models.DashboardResponse, error) {
	switch auth.DashboardRole {
	case models.DashboardRoleTrustee:
		return s.trusteeDashboard(ctx, auth)
	case models.DashboardRoleAssetCustodian:
		return s.custodianDashboard(ctx, auth)
	case models.DashboardRoleAssetManager:
		return s.managerDashboard(ctx, auth)
	default:
		// View-only roles (legal/financial adviser, issuing house, rating agency)
		// share a generic dashboard: their assigned assets + recent activity.
		if models.IsSupportedDashboardRole(auth.DashboardRole) {
			return s.viewOnlyDashboard(ctx, auth)
		}
		return nil, NewHTTPError(http.StatusForbidden, "unsupported stakeholder role")
	}
}

// viewOnlyDashboard is the generic dashboard for portal roles without
// role-specific aggregates. It surfaces the caller's assigned assets, a
// portfolio summary, recent assets and recent activity — all role-agnostic.
func (s *DashboardService) viewOnlyDashboard(ctx context.Context, auth AuthContext) (*models.DashboardResponse, error) {
	assets, err := s.assets.ListAllAssets(ctx, auth)
	if err != nil {
		return nil, err
	}
	portfolio, recentAssets, recentActivity, err := s.dashboardPresentation(ctx, auth, assets)
	if err != nil {
		return nil, err
	}
	return &models.DashboardResponse{Role: auth.DashboardRole, Summary: map[string]interface{}{
		"assigned_assets": len(assets),
	}, RecentAssets: recentAssets, PortfolioSummary: &portfolio, RecentActivity: recentActivity}, nil
}

func (s *DashboardService) trusteeDashboard(ctx context.Context, auth AuthContext) (*models.DashboardResponse, error) {
	assets, err := s.assets.ListAllAssets(ctx, auth)
	if err != nil {
		return nil, err
	}
	portfolio, recentAssets, recentActivity, err := s.dashboardPresentation(ctx, auth, assets)
	if err != nil {
		return nil, err
	}
	pendingDueDiligence := countWhere(s.db.WithContext(ctx).Model(&models.DueDiligenceChecklist{}), "trustee_org_id = ? AND status IN ?", auth.OrganizationID, []string{models.DueDiligenceStatusDraft, models.DueDiligenceStatusInReview})
	pendingFundReleases := countWhere(s.db.WithContext(ctx).Model(&models.FundReleaseRequest{}), "trustee_org_id = ? AND status = ?", auth.OrganizationID, models.FundReleaseStatusSubmitted)
	pendingDistributions := countWhere(s.db.WithContext(ctx).Model(&models.Distribution{}), "trustee_org_id = ? AND status = ?", auth.OrganizationID, models.DistributionStatusProposed)
	return &models.DashboardResponse{Role: auth.DashboardRole, Summary: map[string]interface{}{
		"assets_under_trust":                  len(assets),
		"pending_due_diligence":               pendingDueDiligence,
		"pending_fund_release_approvals":      pendingFundReleases,
		"pending_distribution_authorizations": pendingDistributions,
	}, RecentAssets: recentAssets, PortfolioSummary: &portfolio, RecentActivity: recentActivity}, nil
}

func (s *DashboardService) custodianDashboard(ctx context.Context, auth AuthContext) (*models.DashboardResponse, error) {
	assets, err := s.assets.ListAllAssets(ctx, auth)
	if err != nil {
		return nil, err
	}
	portfolio, fees := custodianPortfolio(assets)

	accountBalances := make(map[string]decimal.Decimal)
	var accountCount int64
	if err := s.db.WithContext(ctx).Model(&models.SegregatedAccount{}).Where("custodian_org_id = ?", auth.OrganizationID).Count(&accountCount).Error; err != nil {
		return nil, err
	}
	var accounts []models.SegregatedAccount
	if err := s.db.WithContext(ctx).Where("custodian_org_id = ?", auth.OrganizationID).Find(&accounts).Error; err != nil {
		return nil, err
	}
	for _, account := range accounts {
		currency := normalizedDashboardCurrency(account.Currency)
		accountBalances[currency] = accountBalances[currency].Add(account.Balance)
	}
	pendingReleases := countWhere(s.db.WithContext(ctx).Model(&models.FundReleaseRequest{}), "custodian_org_id = ? AND status IN ?", auth.OrganizationID, []string{models.FundReleaseStatusExecutionPending, models.FundReleaseStatusProcessing})
	pendingCompliance := countWhere(s.db.WithContext(ctx).Model(&models.ComplianceItem{}), "org_id = ? AND status IN ?", auth.OrganizationID, []string{models.ComplianceStatusPending, models.ComplianceStatusOverdue})
	recentActivity, err := s.dashboardRecentActivity(ctx, auth, 10)
	if err != nil {
		return nil, err
	}
	recentAssets := assets
	if len(recentAssets) > 7 {
		recentAssets = recentAssets[:7]
	}
	legacyBalance := decimal.Zero
	for _, amount := range accountBalances {
		legacyBalance = legacyBalance.Add(amount)
	}
	return &models.DashboardResponse{Role: auth.DashboardRole, Summary: map[string]interface{}{
		"segregated_accounts":      accountCount,
		"total_account_balance":    legacyBalance.String(),
		"account_balances":         currencyTotals(accountBalances),
		"pending_fund_releases":    pendingReleases,
		"pending_compliance_items": pendingCompliance,
		"total_asset_count":        portfolio.AssetCount,
		"assets_under_custody":     portfolio.Values,
		"fees_generated":           fees,
	}, RecentAssets: recentAssets, PortfolioSummary: &portfolio, RecentActivity: recentActivity}, nil
}

func custodianPortfolio(assets []models.TokenizedAssetResponse) (models.CustodianPortfolioSummary, []models.CurrencyTotal) {
	totals := make(map[string]decimal.Decimal)
	feeTotals := make(map[string]decimal.Decimal)
	type categoryAccumulator struct {
		count  int
		values map[string]decimal.Decimal
	}
	categories := make(map[string]*categoryAccumulator)
	for _, asset := range assets {
		currency := normalizedDashboardCurrency(asset.AssetQuoteCurrency)
		value := decimal.NewFromFloat(asset.AssetCurrentValue)
		totals[currency] = totals[currency].Add(value)
		feeTotals[currency] = feeTotals[currency].Add(decimal.NewFromFloat(asset.CustodianFeeValue))
		category := strings.TrimSpace(asset.AssetSector)
		if category == "" {
			category = "Uncategorized"
		}
		accumulator, ok := categories[category]
		if !ok {
			accumulator = &categoryAccumulator{values: make(map[string]decimal.Decimal)}
			categories[category] = accumulator
		}
		accumulator.count++
		accumulator.values[currency] = accumulator.values[currency].Add(value)
	}

	breakdown := make([]models.CustodianPortfolioCategory, 0, len(categories))
	for category, accumulator := range categories {
		breakdown = append(breakdown, models.CustodianPortfolioCategory{
			Category: category, AssetCount: accumulator.count, Values: currencyTotals(accumulator.values),
		})
	}
	sort.Slice(breakdown, func(i, j int) bool {
		if breakdown[i].AssetCount == breakdown[j].AssetCount {
			return breakdown[i].Category < breakdown[j].Category
		}
		return breakdown[i].AssetCount > breakdown[j].AssetCount
	})
	return models.CustodianPortfolioSummary{AssetCount: len(assets), Values: currencyTotals(totals), Categories: breakdown}, currencyTotals(feeTotals)
}

func currencyTotals(totals map[string]decimal.Decimal) []models.CurrencyTotal {
	currencies := make([]string, 0, len(totals))
	for currency := range totals {
		currencies = append(currencies, currency)
	}
	sort.Strings(currencies)
	result := make([]models.CurrencyTotal, 0, len(currencies))
	for _, currency := range currencies {
		result = append(result, models.CurrencyTotal{Currency: currency, Amount: totals[currency].String()})
	}
	return result
}

func normalizedDashboardCurrency(currency string) string {
	currency = strings.ToUpper(strings.TrimSpace(currency))
	if currency == "" {
		return "UNKNOWN"
	}
	return currency
}

func (s *DashboardService) dashboardRecentActivity(ctx context.Context, auth AuthContext, limit int) ([]models.DashboardActivity, error) {
	if limit < 1 {
		return []models.DashboardActivity{}, nil
	}
	var notifications []models.StakeholderNotification
	if err := s.db.WithContext(ctx).Where("recipient_org_id = ?", auth.OrganizationID).
		Order("created_at desc").Limit(limit).Find(&notifications).Error; err != nil {
		return nil, err
	}
	var auditLogs []models.StakeholderAuditLog
	if err := s.db.WithContext(ctx).Where("actor_org_id = ?", auth.OrganizationID).
		Order("created_at desc").Limit(limit).Find(&auditLogs).Error; err != nil {
		return nil, err
	}

	organizationIDs := make([]string, 0, len(notifications)+1)
	organizationIDs = append(organizationIDs, auth.OrganizationID)
	for _, notification := range notifications {
		if notification.SenderOrgID != nil && *notification.SenderOrgID != "" {
			organizationIDs = append(organizationIDs, *notification.SenderOrgID)
		}
	}
	var organizations []coreModels.Organization
	if err := s.db.WithContext(ctx).Where("id IN ?", organizationIDs).Find(&organizations).Error; err != nil {
		return nil, err
	}
	organizationNames := make(map[string]string, len(organizations))
	for _, organization := range organizations {
		organizationNames[organization.ID] = organization.Name
	}

	memberIDs := make([]string, 0, len(auditLogs))
	for _, auditLog := range auditLogs {
		if auditLog.ActorMemberID != "" {
			memberIDs = append(memberIDs, auditLog.ActorMemberID)
		}
	}
	var members []coreModels.OrganizationMember
	if len(memberIDs) > 0 {
		if err := s.db.WithContext(ctx).Where("id IN ?", memberIDs).Find(&members).Error; err != nil {
			return nil, err
		}
	}
	memberNames := make(map[string]string, len(members))
	for _, member := range members {
		name := strings.TrimSpace(strings.TrimSpace(member.FirstName) + " " + strings.TrimSpace(member.LastName))
		if name == "" {
			name = member.Email
		}
		memberNames[member.ID] = name
	}

	activity := make([]models.DashboardActivity, 0, len(notifications)+len(auditLogs))
	for _, notification := range notifications {
		actorOrganization := ""
		if notification.SenderOrgID != nil {
			actorOrganization = organizationNames[*notification.SenderOrgID]
		}
		activity = append(activity, models.DashboardActivity{
			ID: notification.ID, Source: "notification", Type: notification.Type, Title: notification.Title,
			Message: notification.Message, ActorOrganization: actorOrganization,
			EntityType: notification.RelatedEntityType, EntityID: notification.RelatedEntityID, CreatedAt: notification.CreatedAt,
		})
	}
	for _, auditLog := range auditLogs {
		status := ""
		if value, ok := auditLog.AfterState["status"]; ok && value != nil {
			status = fmt.Sprint(value)
		}
		activity = append(activity, models.DashboardActivity{
			ID: auditLog.ID, Source: "audit", Type: auditLog.Action, Title: auditLog.Action, Status: status,
			ActorOrganization: organizationNames[auditLog.ActorOrgID], ActorName: memberNames[auditLog.ActorMemberID],
			EntityType: auditLog.EntityType, EntityID: auditLog.EntityID, CreatedAt: auditLog.CreatedAt,
		})
	}
	sort.SliceStable(activity, func(i, j int) bool { return activity[i].CreatedAt.After(activity[j].CreatedAt) })
	if len(activity) > limit {
		activity = activity[:limit]
	}
	return activity, nil
}

func (s *DashboardService) managerDashboard(ctx context.Context, auth AuthContext) (*models.DashboardResponse, error) {
	assets, err := s.assets.ListAllAssets(ctx, auth)
	if err != nil {
		return nil, err
	}
	portfolio, recentAssets, recentActivity, err := s.dashboardPresentation(ctx, auth, assets)
	if err != nil {
		return nil, err
	}
	var ytdRevenue decimal.Decimal
	startOfYear := time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	var records []models.RevenueRecord
	if err := s.db.WithContext(ctx).Where("manager_org_id = ? AND created_at >= ?", auth.OrganizationID, startOfYear).Find(&records).Error; err != nil {
		return nil, err
	}
	for _, record := range records {
		ytdRevenue = ytdRevenue.Add(record.Amount)
	}
	pendingApprovals := countWhere(s.db.WithContext(ctx).Model(&models.FundReleaseRequest{}), "requester_org_id = ? AND status = ?", auth.OrganizationID, models.FundReleaseStatusSubmitted)
	return &models.DashboardResponse{Role: auth.DashboardRole, Summary: map[string]interface{}{
		"assets_managed":            len(assets),
		"ytd_revenue":               ytdRevenue.String(),
		"pending_trustee_approvals": pendingApprovals,
	}, RecentAssets: recentAssets, PortfolioSummary: &portfolio, RecentActivity: recentActivity}, nil
}

func (s *DashboardService) dashboardPresentation(ctx context.Context, auth AuthContext, assets []models.TokenizedAssetResponse) (models.CustodianPortfolioSummary, []models.TokenizedAssetResponse, []models.DashboardActivity, error) {
	portfolio, _ := custodianPortfolio(assets)
	recentAssets := assets
	if len(recentAssets) > 7 {
		recentAssets = recentAssets[:7]
	}
	recentActivity, err := s.dashboardRecentActivity(ctx, auth, 10)
	if err != nil {
		return models.CustodianPortfolioSummary{}, nil, nil, err
	}
	return portfolio, recentAssets, recentActivity, nil
}

func countWhere(query *gorm.DB, condition string, args ...interface{}) int64 {
	var count int64
	_ = query.Where(condition, args...).Count(&count).Error
	return count
}

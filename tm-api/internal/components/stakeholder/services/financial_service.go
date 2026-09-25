package services

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"admin-panel-dashboard/internal/components/stakeholder/models"
	coreModels "admin-panel-dashboard/internal/models"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FinancialService struct {
	db             *gorm.DB
	assets         *AssetService
	audit          *AuditService
	notifications  *NotificationService
	authorizations *AuthorizationService
	documents      *DocumentService
	payouts        DistributionPayoutClient
}

func NewFinancialService(db *gorm.DB, assets *AssetService, audit *AuditService, notifications *NotificationService, authorizations *AuthorizationService, documents *DocumentService, payoutClients ...DistributionPayoutClient) *FinancialService {
	service := &FinancialService{db: db, assets: assets, audit: audit, notifications: notifications, authorizations: authorizations, documents: documents}
	if len(payoutClients) > 0 {
		service.payouts = payoutClients[0]
	}
	return service
}

func (s *FinancialService) ListRevenue(ctx context.Context, auth AuthContext, page, limit int) ([]models.RevenueRecord, int64, error) {
	if auth.DashboardRole != models.DashboardRoleAssetManager {
		return nil, 0, NewHTTPError(http.StatusForbidden, "asset manager role is required")
	}
	query := s.db.WithContext(ctx).Model(&models.RevenueRecord{}).Where("manager_org_id = ?", auth.OrganizationID)
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var records []models.RevenueRecord
	if err := query.Order("created_at desc").Offset((page - 1) * limit).Limit(limit).Find(&records).Error; err != nil {
		return nil, 0, err
	}
	return records, total, nil
}

func (s *FinancialService) CreateRevenue(ctx context.Context, auth AuthContext, req models.CreateRevenueRequest) (*models.RevenueRecord, error) {
	if auth.DashboardRole != models.DashboardRoleAssetManager {
		return nil, NewHTTPError(http.StatusForbidden, "asset manager role is required")
	}
	amount, err := ParseRequiredDecimal(req.Amount, "amount")
	if err != nil {
		return nil, err
	}
	currency, err := NormalizeCurrency(req.Currency)
	if err != nil {
		return nil, err
	}
	_, asset, err := s.assets.RequireAssignmentForAsset(ctx, auth, req.AssetID)
	if err != nil {
		return nil, err
	}
	periodStart, err := ParseDate(req.PeriodStart)
	if err != nil {
		return nil, err
	}
	periodEnd, err := ParseDate(req.PeriodEnd)
	if err != nil {
		return nil, err
	}
	collectedAt, err := ParseDate(req.CollectedAt)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	record := models.RevenueRecord{
		ID:                uuid.NewString(),
		AssetID:           asset.ID,
		AssetCode:         asset.AssetCode,
		ManagerOrgID:      auth.OrganizationID,
		PeriodStart:       periodStart,
		PeriodEnd:         periodEnd,
		Source:            strings.TrimSpace(req.Source),
		Amount:            amount,
		Currency:          currency,
		Status:            models.RevenueStatusRecorded,
		CollectedAt:       collectedAt,
		CreatedByMemberID: auth.MemberID,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	if record.Source == "" {
		return nil, NewHTTPError(http.StatusBadRequest, "source is required")
	}
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&record).Error; err != nil {
			return err
		}
		if s.audit != nil {
			event := AuditEventFromAuth(auth, models.ActionRevenueCreate, models.EntityRevenueRecord, record.ID)
			event.AfterState = models.JSONMap{"status": record.Status, "amount": record.Amount.String(), "currency": record.Currency}
			if err := s.audit.Record(ctx, tx, event); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (s *FinancialService) SubmitDistribution(ctx context.Context, auth AuthContext, req models.SubmitDistributionRequest) (*models.Distribution, error) {
	if auth.DashboardRole != models.DashboardRoleAssetManager {
		return nil, NewHTTPError(http.StatusForbidden, "asset manager role is required")
	}
	amount, err := ParseRequiredDecimal(req.Amount, "amount")
	if err != nil {
		return nil, err
	}
	currency, err := NormalizeCurrency(req.Currency)
	if err != nil {
		return nil, err
	}
	assignment, asset, err := s.assets.RequireAssignmentForAsset(ctx, auth, req.AssetID)
	if err != nil {
		return nil, err
	}
	if assignment.TrusteeOrgID == nil || *assignment.TrusteeOrgID == "" {
		return nil, NewHTTPError(http.StatusUnprocessableEntity, "asset has no trustee assignment")
	}
	scheduledDate, err := ParseDate(req.ScheduledDate)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	distribution := models.Distribution{
		ID:                 uuid.NewString(),
		AssetID:            asset.ID,
		AssetCode:          asset.AssetCode,
		ProposedByOrgID:    auth.OrganizationID,
		ProposedByMemberID: auth.MemberID,
		TrusteeOrgID:       *assignment.TrusteeOrgID,
		Amount:             amount,
		Currency:           currency,
		Source:             strings.TrimSpace(req.Source),
		ScheduledDate:      scheduledDate,
		Status:             models.DistributionStatusProposed,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	if distribution.Source == "" {
		return nil, NewHTTPError(http.StatusBadRequest, "source is required")
	}
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&distribution).Error; err != nil {
			return err
		}
		if len(req.RevenueRecordIDs) > 0 {
			if err := tx.Model(&models.RevenueRecord{}).
				Where("id IN ? AND manager_org_id = ?", req.RevenueRecordIDs, auth.OrganizationID).
				Updates(map[string]interface{}{"status": models.RevenueStatusSubmittedForDistribution, "updated_at": now}).Error; err != nil {
				return err
			}
		}
		if s.notifications != nil {
			if err := s.notifications.Create(ctx, tx, NotificationInput{
				RecipientOrgID:    distribution.TrusteeOrgID,
				SenderOrgID:       &auth.OrganizationID,
				Type:              "distribution.proposed",
				Title:             "Distribution proposed",
				Message:           fmt.Sprintf("A distribution proposal for %s %s is pending trustee authorization.", distribution.Currency, distribution.Amount.String()),
				RelatedEntityType: models.EntityDistribution,
				RelatedEntityID:   distribution.ID,
			}); err != nil {
				return err
			}
		}
		if s.audit != nil {
			event := AuditEventFromAuth(auth, models.ActionDistributionSubmit, models.EntityDistribution, distribution.ID)
			event.AfterState = models.JSONMap{"status": distribution.Status, "amount": distribution.Amount.String(), "currency": distribution.Currency}
			if err := s.audit.Record(ctx, tx, event); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &distribution, nil
}

func (s *FinancialService) ListDistributions(ctx context.Context, auth AuthContext, page, limit int, history bool) ([]models.DistributionListItemResponse, int64, error) {
	if auth.DashboardRole != models.DashboardRoleTrustee {
		return nil, 0, NewHTTPError(http.StatusForbidden, "trustee role is required")
	}
	query := s.db.WithContext(ctx).Model(&models.Distribution{}).Where("trustee_org_id = ?", auth.OrganizationID)
	if history {
		query = query.Where("status IN ?", []string{models.DistributionStatusAuthorized, models.DistributionStatusRejected, models.DistributionStatusCompleted, models.DistributionStatusFailed})
	} else {
		query = query.Where("status = ?", models.DistributionStatusProposed)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var records []models.Distribution
	if err := query.Order("created_at desc").Offset((page - 1) * limit).Limit(limit).Find(&records).Error; err != nil {
		return nil, 0, err
	}
	enriched, err := s.enrichDistributions(ctx, records)
	if err != nil {
		return nil, 0, err
	}
	return enriched, total, nil
}

func (s *FinancialService) GetDistribution(ctx context.Context, auth AuthContext, id string) (*models.DistributionDetailResponse, error) {
	if auth.DashboardRole != models.DashboardRoleTrustee {
		return nil, NewHTTPError(http.StatusForbidden, "trustee role is required")
	}
	var distribution models.Distribution
	if err := s.db.WithContext(ctx).Where("id = ? AND trustee_org_id = ?", strings.TrimSpace(id), auth.OrganizationID).First(&distribution).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewHTTPError(http.StatusNotFound, "distribution not found")
		}
		return nil, err
	}
	_, asset, err := s.assets.RequireAssignmentForAsset(ctx, auth, distribution.AssetID)
	if err != nil {
		return nil, err
	}
	totalTokens := decimal.NewFromFloat(asset.NumberOfTokensIssued)
	amountPerToken := decimal.Zero
	if totalTokens.GreaterThan(decimal.Zero) {
		amountPerToken = distribution.Amount.Div(totalTokens).Round(8)
	}
	amount := models.DistributionAmountResponse{Amount: distribution.Amount.String(), Currency: distribution.Currency}
	enriched, err := s.enrichDistributions(ctx, []models.Distribution{distribution})
	if err != nil {
		return nil, err
	}
	payout := models.DistributionPayoutResponse{
		DistributionID: distribution.ID, TokenizedAssetID: distribution.AssetID,
		Amount: distribution.Amount.String(), Currency: distribution.Currency,
		AmountPerToken: amountPerToken.String(), Status: "not_registered",
		Payouts: make([]models.DistributionPayoutRecordResponse, 0),
	}
	if s.payouts != nil {
		walletPayout, err := s.payouts.Get(ctx, distribution.ID)
		if err != nil {
			return nil, err
		}
		if walletPayout != nil {
			payout = *walletPayout
			if payout.Payouts == nil {
				payout.Payouts = make([]models.DistributionPayoutRecordResponse, 0)
			}
		}
	}
	return &models.DistributionDetailResponse{
		Distribution: distribution,
		Requester:    enriched[0].Requester,
		Breakdown: models.DistributionBreakdownResponse{
			TotalTokens: totalTokens.String(), TotalTokenHolders: asset.TokenHolderCount,
			NetIncome:                  amount,
			DistributionAmountPerToken: models.DistributionAmountResponse{Amount: amountPerToken.String(), Currency: distribution.Currency},
			Calculation:                "net_income / total_tokens",
		},
		Payout: payout,
	}, nil
}

func (s *FinancialService) DistributionPayoutCSV(ctx context.Context, auth AuthContext, id string) ([]byte, error) {
	detail, err := s.GetDistribution(ctx, auth, id)
	if err != nil {
		return nil, err
	}

	var output bytes.Buffer
	writer := csv.NewWriter(&output)
	if err := writer.Write([]string{
		"payout_id", "beneficiary_address", "confirmed_token_balance", "amount", "currency",
		"cannot_receive_asset", "paid", "created_at",
	}); err != nil {
		return nil, fmt.Errorf("write payout CSV header: %w", err)
	}
	for _, payout := range detail.Payout.Payouts {
		if err := writer.Write([]string{
			payout.ID,
			payout.BeneficiaryAddress,
			payout.ConfirmedTokenBalance,
			payout.Amount,
			payout.Currency,
			strconv.FormatBool(payout.CannotReceiveAsset),
			strconv.FormatBool(payout.Paid),
			payout.CreatedAt.UTC().Format(time.RFC3339),
		}); err != nil {
			return nil, fmt.Errorf("write payout CSV row: %w", err)
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("flush payout CSV: %w", err)
	}
	return output.Bytes(), nil
}

func (s *FinancialService) enrichDistributions(ctx context.Context, records []models.Distribution) ([]models.DistributionListItemResponse, error) {
	result := make([]models.DistributionListItemResponse, len(records))
	if len(records) == 0 {
		return result, nil
	}
	memberIDs := make([]string, 0, len(records))
	organizationIDs := make([]string, 0, len(records))
	for _, record := range records {
		memberIDs = append(memberIDs, record.ProposedByMemberID)
		organizationIDs = append(organizationIDs, record.ProposedByOrgID)
	}
	var members []coreModels.OrganizationMember
	if err := s.db.WithContext(ctx).Where("id IN ?", memberIDs).Find(&members).Error; err != nil {
		return nil, err
	}
	var organizations []coreModels.Organization
	if err := s.db.WithContext(ctx).Where("id IN ?", organizationIDs).Find(&organizations).Error; err != nil {
		return nil, err
	}
	membersByID := make(map[string]coreModels.OrganizationMember, len(members))
	for _, member := range members {
		membersByID[member.ID] = member
	}
	organizationsByID := make(map[string]coreModels.Organization, len(organizations))
	for _, organization := range organizations {
		organizationsByID[organization.ID] = organization
	}
	for i, record := range records {
		requester := models.FundReleaseRequesterResponse{
			MemberID:       record.ProposedByMemberID,
			OrganizationID: record.ProposedByOrgID,
		}
		if member, ok := membersByID[record.ProposedByMemberID]; ok && member.OrganizationID == record.ProposedByOrgID {
			requester.FirstName = member.FirstName
			requester.LastName = member.LastName
			requester.Email = member.Email
		}
		if organization, ok := organizationsByID[record.ProposedByOrgID]; ok {
			requester.OrganizationName = organization.Name
		}
		result[i] = models.DistributionListItemResponse{Distribution: record, Requester: requester}
	}
	return result, nil
}

func (s *FinancialService) AuthorizeDistribution(ctx context.Context, auth AuthContext, id string, req models.ChallengeActionRequest) (*models.Distribution, error) {
	if auth.DashboardRole != models.DashboardRoleTrustee {
		return nil, NewHTTPError(http.StatusForbidden, "trustee role is required")
	}
	var updated models.Distribution
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		distribution, err := s.lockDistribution(ctx, tx, id, auth.OrganizationID)
		if err != nil {
			return err
		}
		if distribution.Status != models.DistributionStatusProposed {
			return NewHTTPError(http.StatusConflict, "distribution cannot be authorized from current status")
		}
		if s.authorizations != nil {
			if err := s.authorizations.ConsumeVerifiedChallenge(ctx, tx, auth, models.ActionDistributionAuthorize, models.EntityDistribution, distribution.ID, req.ChallengeID); err != nil {
				return err
			}
		}
		now := time.Now()
		before := distribution.Status
		distribution.Status = models.DistributionStatusAuthorized
		distribution.AuthorizedByMemberID = &auth.MemberID
		distribution.AuthorizedAt = &now
		distribution.UpdatedAt = now
		if err := tx.Save(distribution).Error; err != nil {
			return err
		}
		if s.payouts != nil {
			if _, err := s.payouts.Register(ctx, DistributionPayoutRegistration{
				DistributionID: distribution.ID, TokenizedAssetID: distribution.AssetID,
				Amount: distribution.Amount.String(), Currency: distribution.Currency,
			}); err != nil {
				return err
			}
		}
		updated = *distribution
		if err := s.notifyDistributionDecision(ctx, tx, auth.OrganizationID, distribution, "distribution.authorized", "Distribution authorized", "A distribution proposal was authorized by the trustee."); err != nil {
			return err
		}
		if s.audit != nil {
			event := AuditEventFromAuth(auth, models.ActionDistributionAuthorize, models.EntityDistribution, distribution.ID)
			event.BeforeState = models.JSONMap{"status": before}
			event.AfterState = models.JSONMap{"status": distribution.Status}
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

func (s *FinancialService) RejectDistribution(ctx context.Context, auth AuthContext, id string, req models.RejectRequest) (*models.Distribution, error) {
	if auth.DashboardRole != models.DashboardRoleTrustee {
		return nil, NewHTTPError(http.StatusForbidden, "trustee role is required")
	}
	reason := strings.TrimSpace(req.Reason)
	if reason == "" {
		return nil, NewHTTPError(http.StatusBadRequest, "rejection reason is required")
	}
	var updated models.Distribution
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		distribution, err := s.lockDistribution(ctx, tx, id, auth.OrganizationID)
		if err != nil {
			return err
		}
		if distribution.Status != models.DistributionStatusProposed {
			return NewHTTPError(http.StatusConflict, "distribution cannot be rejected from current status")
		}
		before := distribution.Status
		distribution.Status = models.DistributionStatusRejected
		distribution.RejectionReason = reason
		distribution.UpdatedAt = time.Now()
		if err := tx.Save(distribution).Error; err != nil {
			return err
		}
		updated = *distribution
		if err := s.notifyDistributionDecision(ctx, tx, auth.OrganizationID, distribution, "distribution.rejected", "Distribution rejected", "A distribution proposal was rejected by the trustee."); err != nil {
			return err
		}
		if s.audit != nil {
			event := AuditEventFromAuth(auth, models.ActionDistributionReject, models.EntityDistribution, distribution.ID)
			event.BeforeState = models.JSONMap{"status": before}
			event.AfterState = models.JSONMap{"status": distribution.Status, "reason": reason}
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

func (s *FinancialService) ListValuations(ctx context.Context, auth AuthContext, assetID string, page, limit int) ([]models.AssetValuation, int64, error) {
	if auth.DashboardRole != models.DashboardRoleAssetManager {
		return nil, 0, NewHTTPError(http.StatusForbidden, "asset manager role is required")
	}
	if _, _, err := s.assets.RequireAssignmentForAsset(ctx, auth, assetID); err != nil {
		return nil, 0, err
	}
	query := s.db.WithContext(ctx).Model(&models.AssetValuation{}).Where("asset_id = ? AND manager_org_id = ?", assetID, auth.OrganizationID)
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var records []models.AssetValuation
	if err := query.Order("valuation_date desc").Offset((page - 1) * limit).Limit(limit).Find(&records).Error; err != nil {
		return nil, 0, err
	}
	return records, total, nil
}

func (s *FinancialService) CreateValuation(ctx context.Context, auth AuthContext, req models.CreateValuationRequest) (*models.AssetValuation, error) {
	if auth.DashboardRole != models.DashboardRoleAssetManager {
		return nil, NewHTTPError(http.StatusForbidden, "asset manager role is required")
	}
	valuationAmount, err := ParseRequiredDecimal(req.Valuation, "valuation")
	if err != nil {
		return nil, err
	}
	currency, err := NormalizeCurrency(req.Currency)
	if err != nil {
		return nil, err
	}
	_, asset, err := s.assets.RequireAssignmentForAsset(ctx, auth, req.AssetID)
	if err != nil {
		return nil, err
	}
	if s.documents == nil {
		return nil, NewHTTPError(http.StatusServiceUnavailable, "document service is not configured")
	}
	if err := s.documents.ValidateForUse(ctx, auth, asset.ID, models.DocumentCategoryValuationReport, []string{req.ReportDocumentID}); err != nil {
		return nil, err
	}
	valuationDate, err := ParseDate(req.ValuationDate)
	if err != nil {
		return nil, err
	}
	if valuationDate == nil {
		return nil, NewHTTPError(http.StatusBadRequest, "valuation_date is required")
	}
	now := time.Now()
	valuation := models.AssetValuation{
		ID:                uuid.NewString(),
		AssetID:           asset.ID,
		AssetCode:         asset.AssetCode,
		ManagerOrgID:      auth.OrganizationID,
		Valuation:         valuationAmount,
		Currency:          currency,
		Methodology:       strings.TrimSpace(req.Methodology),
		ValuationDate:     *valuationDate,
		ReportDocumentID:  &req.ReportDocumentID,
		Status:            models.ValuationStatusSubmitted,
		CreatedByMemberID: auth.MemberID,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	if valuation.Methodology == "" {
		return nil, NewHTTPError(http.StatusBadRequest, "methodology is required")
	}
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&valuation).Error; err != nil {
			return err
		}
		if s.audit != nil {
			event := AuditEventFromAuth(auth, models.ActionValuationCreate, models.EntityAssetValuation, valuation.ID)
			event.AfterState = models.JSONMap{"status": valuation.Status, "valuation": valuation.Valuation.String(), "currency": valuation.Currency}
			if err := s.audit.Record(ctx, tx, event); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &valuation, nil
}

func (s *FinancialService) RequestIndependentValuation(ctx context.Context, auth AuthContext, id string) (*models.AssetValuation, error) {
	if auth.DashboardRole != models.DashboardRoleAssetManager {
		return nil, NewHTTPError(http.StatusForbidden, "asset manager role is required")
	}
	var updated models.AssetValuation
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var valuation models.AssetValuation
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND manager_org_id = ?", id, auth.OrganizationID).First(&valuation).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return NewHTTPError(http.StatusNotFound, "valuation not found")
			}
			return err
		}
		before := valuation.Status
		valuation.Status = models.ValuationStatusIndependentRequested
		valuation.UpdatedAt = time.Now()
		if err := tx.Save(&valuation).Error; err != nil {
			return err
		}
		updated = valuation
		if s.audit != nil {
			event := AuditEventFromAuth(auth, models.ActionValuationRequestIndependent, models.EntityAssetValuation, valuation.ID)
			event.BeforeState = models.JSONMap{"status": before}
			event.AfterState = models.JSONMap{"status": valuation.Status}
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

func (s *FinancialService) lockDistribution(ctx context.Context, tx *gorm.DB, id, trusteeOrgID string) (*models.Distribution, error) {
	var distribution models.Distribution
	if err := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ? AND trustee_org_id = ?", id, trusteeOrgID).
		First(&distribution).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewHTTPError(http.StatusNotFound, "distribution not found")
		}
		return nil, err
	}
	return &distribution, nil
}

func (s *FinancialService) notifyDistributionDecision(ctx context.Context, tx *gorm.DB, senderOrgID string, distribution *models.Distribution, notificationType, title, message string) error {
	if s.notifications == nil {
		return nil
	}
	if distribution.ProposedByOrgID == "" || distribution.ProposedByOrgID == senderOrgID {
		return nil
	}
	return s.notifications.Create(ctx, tx, NotificationInput{
		RecipientOrgID:    distribution.ProposedByOrgID,
		SenderOrgID:       &senderOrgID,
		Type:              notificationType,
		Title:             title,
		Message:           message,
		RelatedEntityType: models.EntityDistribution,
		RelatedEntityID:   distribution.ID,
	})
}

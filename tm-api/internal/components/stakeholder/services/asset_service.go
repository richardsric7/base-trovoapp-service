package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	stakeholderDB "admin-panel-dashboard/internal/components/stakeholder/db"
	"admin-panel-dashboard/internal/components/stakeholder/models"
	coreModels "admin-panel-dashboard/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AssetService struct {
	adminDB *gorm.DB
	client  stakeholderDB.TokenizationReadClient
	audit   *AuditService
}

const assetRiskScoreMaximum = 100

func NewAssetService(adminDB *gorm.DB, client stakeholderDB.TokenizationReadClient, audit *AuditService) *AssetService {
	return &AssetService{adminDB: adminDB, client: client, audit: audit}
}

func (s *AssetService) CreateAssignment(ctx context.Context, actor AuthContext, assignedBy string, req models.CreateAssetAssignmentRequest) (*models.StakeholderAssetAssignment, error) {
	if strings.TrimSpace(req.Status) == "" {
		req.Status = models.AssignmentStatusActive
	}
	if req.Status != models.AssignmentStatusActive && req.Status != models.AssignmentStatusInactive {
		return nil, NewHTTPError(http.StatusBadRequest, "invalid assignment status")
	}
	trustee, err := s.requirePortalOrganization(req.TrusteeOrgID, models.DashboardRoleTrustee)
	if err != nil {
		return nil, err
	}
	var manager *coreModels.Organization
	if req.AssetManagerOrgID != nil && strings.TrimSpace(*req.AssetManagerOrgID) != "" {
		manager, err = s.requirePortalOrganization(*req.AssetManagerOrgID, models.DashboardRoleAssetManager)
		if err != nil {
			return nil, err
		}
	}
	var custodian *coreModels.Organization
	if req.CustodianOrgID != nil && strings.TrimSpace(*req.CustodianOrgID) != "" {
		custodian, err = s.requirePortalOrganization(*req.CustodianOrgID, models.DashboardRoleAssetCustodian)
		if err != nil {
			return nil, err
		}
	}
	now := time.Now()
	assignment := models.StakeholderAssetAssignment{
		ID:                   uuid.NewString(),
		AssetID:              strings.TrimSpace(req.AssetID),
		AssetCode:            strings.TrimSpace(req.AssetCode),
		TrusteeOrgID:         &trustee.ID,
		TrusteeStakeholderID: trustee.StakeholderID,
		Status:               req.Status,
		AssignedBy:           assignedBy,
		CreatedAt:            now,
		UpdatedAt:            now,
	}
	if manager != nil {
		assignment.AssetManagerOrgID = &manager.ID
		assignment.AssetManagerStakeholderID = manager.StakeholderID
	}
	if custodian != nil {
		assignment.CustodianOrgID = &custodian.ID
		assignment.CustodianStakeholderID = custodian.StakeholderID
	}

	err = s.adminDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing models.StakeholderAssetAssignment
		err := tx.Where("asset_id = ?", assignment.AssetID).First(&existing).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if err == nil {
			assignment.ID = existing.ID
			assignment.CreatedAt = existing.CreatedAt
			if err := tx.Model(&existing).Updates(map[string]interface{}{
				"asset_code":                   assignment.AssetCode,
				"trustee_org_id":               assignment.TrusteeOrgID,
				"trustee_stakeholder_id":       assignment.TrusteeStakeholderID,
				"asset_manager_org_id":         assignment.AssetManagerOrgID,
				"asset_manager_stakeholder_id": assignment.AssetManagerStakeholderID,
				"custodian_org_id":             assignment.CustodianOrgID,
				"custodian_stakeholder_id":     assignment.CustodianStakeholderID,
				"status":                       assignment.Status,
				"assigned_by":                  assignment.AssignedBy,
				"updated_at":                   now,
			}).Error; err != nil {
				return err
			}
		} else if err := tx.Create(&assignment).Error; err != nil {
			return err
		}
		if s.audit != nil {
			event := AuditEventFromAuth(actor, models.ActionAssetAssignmentCreate, models.EntityAssetAssignment, assignment.ID)
			event.AfterState = models.JSONMap{
				"asset_id":             assignment.AssetID,
				"trustee_org_id":       assignment.TrusteeOrgID,
				"asset_manager_org_id": assignment.AssetManagerOrgID,
				"custodian_org_id":     assignment.CustodianOrgID,
				"status":               assignment.Status,
			}
			if actor.MemberID == "" {
				event.ActorMemberID = assignedBy
				event.ActorRole = "trovo_admin"
			}
			if err := s.audit.Record(ctx, tx, event); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &assignment, nil
}

func (s *AssetService) ListAssets(ctx context.Context, auth AuthContext, filters models.AssetListFilters) (*models.PaginatedResponse, error) {
	assignments, err := s.assignmentsForRole(ctx, auth)
	if err != nil {
		return nil, err
	}
	assetFilters, err := s.assetFiltersForRole(auth, assignments, filters)
	if err != nil {
		return nil, err
	}
	page, err := s.client.ListAssets(ctx, assetFilters)
	if err != nil {
		return nil, err
	}
	pageAssignments, err := s.assignmentsForAssets(ctx, page.Records)
	if err != nil {
		return nil, err
	}
	assignmentByAsset, err := s.resolveAssignments(ctx, page.Records, pageAssignments)
	if err != nil {
		return nil, err
	}
	responses := make([]models.TokenizedAssetResponse, 0, len(page.Records))
	for _, asset := range page.Records {
		if s.assetVisibleToAuth(auth, asset, assignmentByAsset) {
			responses = append(responses, s.assetResponse(auth, asset, assignmentByAsset))
		}
	}
	summary, err := s.assetStatistics(ctx, auth, assetFilters)
	if err != nil {
		return nil, err
	}
	return &models.PaginatedResponse{
		Records: responses,
		Meta:    PaginationMeta(page.Page, page.Limit, int64(summary.Total)),
		Summary: &summary,
	}, nil
}

func (s *AssetService) GetAsset(ctx context.Context, auth AuthContext, idOrCode string) (*models.TokenizedAssetResponse, error) {
	asset, assignmentByAsset, err := s.visibleAsset(ctx, auth, idOrCode)
	if err != nil {
		return nil, err
	}
	resp := s.assetResponse(auth, *asset, assignmentByAsset)
	return &resp, nil
}

func (s *AssetService) GetAssetDetails(ctx context.Context, auth AuthContext, idOrCode string) (*models.TokenizedAssetDetailResponse, error) {
	asset, assignmentByAsset, err := s.visibleAsset(ctx, auth, idOrCode)
	if err != nil {
		return nil, err
	}
	stakeholders, err := s.assignedStakeholders(ctx, *asset)
	if err != nil {
		return nil, err
	}
	documents, err := s.assetDocuments(ctx, auth, asset.ID)
	if err != nil {
		return nil, err
	}
	operation, err := s.assetOperation(ctx, auth, *asset, assignmentByAsset)
	if err != nil {
		return nil, err
	}
	related, err := s.assetRelatedData(ctx, *asset)
	if err != nil {
		return nil, err
	}
	lastAuditDate, err := s.lastCompletedAssetAuditDate(ctx, asset.ID)
	if err != nil {
		return nil, err
	}
	walletDocuments := walletAssetDocumentResponses(related.Documents)
	base := s.assetResponse(auth, *asset, assignmentByAsset)
	// Prefer the resolved human asset-type label (from the wallet lookup) over the
	// raw numeric asset_type id, so the org portal no longer needs the admin-only
	// tokenization params proxy to translate it.
	if related.AssetTypeName != "" {
		base.AssetType = related.AssetTypeName
	}
	return &models.TokenizedAssetDetailResponse{
		TokenizedAssetResponse: base,
		TokenizerUsername:      asset.InitiatorUsername,
		AssignedStakeholders:   stakeholders,
		Wallets: models.AssetWalletDetails{
			IssuingWalletPublicKey:             asset.IssuingWalletPublicKey,
			IssuingWalletAlias:                 asset.IssuingWalletAlias,
			MarketMakingWallet:                 asset.MarketMakingWallet,
			InitialOwnerPreferredWalletAddress: asset.InitialOwnerPreferredWalletAddress,
			WalletToHoldAssetsNotForSale:       asset.WalletToHoldAssetsNotForSale,
		},
		ExemptedCountries: splitCountries(asset.ExemptedCountries),
		Documents:         documents,
		AssetProfile: models.AssetProfileDetails{
			Description: asset.AssetDescription, LogoURL: asset.AssetLogo, Website: asset.AssetWebsite,
			OfferingType: asset.OfferingType, Country: asset.AssetCountryLocation, PhysicalAddress: asset.AssetPhysicalAddress,
			Longitude: asset.AssetLongitude, Latitude: asset.AssetLatitude,
		},
		Ownership: models.AssetOwnershipDetails{
			Type: asset.OwnershipType, Kind: asset.OwnershipKind, OwnerName: asset.AssetOwnerName, OwnerAddress: asset.AssetOwnerAddress,
			RetainedOrContributedValue: asset.AssetOwnerRetainedOrContributedValue, CostOutsideValuation: asset.AssetMiscCostOutsideValuation,
		},
		RiskAndCompliance: models.AssetRiskAndComplianceDetails{
			MaximumRiskScore: assetRiskScoreMaximum,
			ComplianceStatus: assetComplianceStatus(*asset),
			LastAuditDate:    lastAuditDate,
		},
		Protection: models.AssetProtectionDetails{
			Methods: splitCountries(asset.ProtectionMethods), InsuranceCompanyName: asset.InsuranceCompanyName,
			InsurancePolicyNumber: asset.InsurancePolicyNumber, InsurancePolicyHolder: asset.InsurancePolicyHolder,
			InsuranceCoveragePercentage:  asset.PercentageValueOfInsurance,
			FreeFromLiensAndEncumbrances: asset.IsFreeFromLiensAndEncumbrances != 0,
			TransferTitleToCustodian:     asset.AgreeTransferTitleToCustodian != 0, RevenueGuarantees: asset.ContractualProtectionRevGuarantees != 0,
			PerformanceBond: asset.ContractualProtectionPerfBond != 0, ServiceLevelAgreement: asset.ContractualProtectionSLA != 0,
			CompletionGuarantees:             asset.RiskSharingMechanismCompletionGuarantees != 0,
			PublicPrivatePartnerships:        asset.RiskSharingMechanismPPPs != 0,
			HedgeInstruments:                 asset.RiskSharingMechanismHedgeInstruments != 0,
			IndependentMonitoringList:        asset.IndependentMonitoringList,
			ESGSustainabilityCertifications:  asset.ESGSafeguardsSustainabilityCerts != 0,
			ESGCommunityEngagementPlan:       asset.ESGSafeguardsCommunityEngagementPlan != 0,
			SurveillanceSystems:              asset.SecurityMeasuresSurveillanceSystems != 0,
			OnSiteSecurityPersonnel:          asset.SecurityMeasuresOnSitePersonnel != 0,
			PerimeterSecurity:                asset.SecurityMeasuresPerimeterSecurity != 0,
			CriticalInfrastructureProtection: asset.SecurityMeasuresCriticalInfrastructure != 0,
			LegalAdviser:                     asset.LegalAdvisor,
			FinancialAdviser:                 asset.FinancialAdvisor,
			Other:                            asset.OtherAssetProtection,
		},
		Offering: models.AssetOfferingDetails{
			TokensIssued: asset.NumberOfTokenToBeIssued, TokensForSale: asset.NumberOfTokenToBeSold,
			TokensNotForSale: asset.TotalTokenHeldByManager, MaximumTokensForSale: asset.MaxNumberOfTokenAvailableForSale,
			PricePerToken: asset.PricePerToken, AmountToBeRaised: asset.NumberOfTokenToBeSold * asset.PricePerToken,
			SalesStart: asset.SalesStart, SalesEnd: asset.SalesEnd, CapQuantity: asset.CapQuantity,
			CapAmount: asset.CapAmountInFiat, CapDurationDays: asset.CapDurationInDays,
			PayoutCycle: asset.ProceedCycle, PayoutCurrency: asset.ProceedPayoutCurrency, PayoutType: asset.ProceedPayoutType,
		},
		Fees: models.AssetTokenizationFeeDetails{
			FiatAmount: asset.FeeInFiat, AssetAmount: asset.FeeInAsset, AssetPercentage: asset.FeeInAssetPercent,
			Preferred: preferredFeeResponse(related.PreferredFee),
		},
		ReceivingAccount: models.AssetReceivingAccountDetails{
			BankID: asset.BankID, BankName: walletBankName(related.Bank), BankCountry: walletBankCountry(related.Bank),
			AccountName: asset.BeneficiaryName, AccountNumber: asset.AccountNumber,
		},
		ApplicationFee: models.AssetApplicationFeeDetails{
			Amount: asset.TokenizationApplicationFee, Asset: asset.TokenizationApplicationFeeAsset,
			Payments: applicationFeePaymentResponses(related.PaymentProofs),
		},
		TokenizationDocuments: walletDocuments,
		Timeline:              tokenizationTimeline(*asset, related.StatusCatalog),
		OperationalUpdate:     operation,
	}, nil
}

func (s *AssetService) lastCompletedAssetAuditDate(ctx context.Context, assetID string) (*time.Time, error) {
	var checklist models.DueDiligenceChecklist
	err := s.adminDB.WithContext(ctx).
		Where("asset_id = ? AND status IN ?", assetID, []string{models.DueDiligenceStatusApproved, models.DueDiligenceStatusRejected}).
		Order("updated_at desc").
		First(&checklist).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if checklist.ApprovedAt != nil {
		return checklist.ApprovedAt, nil
	}
	if checklist.RejectedAt != nil {
		return checklist.RejectedAt, nil
	}
	return &checklist.UpdatedAt, nil
}

func (s *AssetService) assetOperation(ctx context.Context, auth AuthContext, asset stakeholderDB.TokenizedAsset, assignmentByAsset map[string]models.StakeholderAssetAssignment) (*models.StakeholderAssetOperation, error) {
	managerOrgID := ""
	if auth.DashboardRole == models.DashboardRoleAssetManager {
		managerOrgID = auth.OrganizationID
	} else if assignment, ok := findAssignmentForAsset(asset, assignmentByAsset); ok && assignment.AssetManagerOrgID != nil {
		managerOrgID = strings.TrimSpace(*assignment.AssetManagerOrgID)
	}
	if managerOrgID == "" {
		return nil, nil
	}
	var operation models.StakeholderAssetOperation
	if err := s.adminDB.WithContext(ctx).Where("asset_id = ? AND manager_org_id = ?", asset.ID, managerOrgID).First(&operation).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &operation, nil
}

func (s *AssetService) assetRelatedData(ctx context.Context, asset stakeholderDB.TokenizedAsset) (stakeholderDB.TokenizedAssetRelatedData, error) {
	client, ok := s.client.(stakeholderDB.TokenizationDetailReadClient)
	if !ok {
		return stakeholderDB.TokenizedAssetRelatedData{
			Documents: []stakeholderDB.AssetTokenizationDocument{}, PaymentProofs: []stakeholderDB.TokenizationFeeProofOfPayment{}, StatusCatalog: []stakeholderDB.WalletTokenizationStatus{},
		}, nil
	}
	return client.GetAssetRelatedData(ctx, asset)
}

func (s *AssetService) WalletAssetDocuments(ctx context.Context, auth AuthContext, idOrCode string) ([]models.WalletAssetDocumentResponse, error) {
	asset, _, err := s.visibleAsset(ctx, auth, idOrCode)
	if err != nil {
		return nil, err
	}
	related, err := s.assetRelatedData(ctx, *asset)
	if err != nil {
		return nil, err
	}
	return walletAssetDocumentResponses(related.Documents), nil
}

func walletAssetDocumentResponses(documents []stakeholderDB.AssetTokenizationDocument) []models.WalletAssetDocumentResponse {
	result := make([]models.WalletAssetDocumentResponse, 0, len(documents))
	for _, document := range documents {
		result = append(result, models.WalletAssetDocumentResponse{
			ID: document.ID, DocumentType: document.DocumentType, Title: document.DocumentTitle,
			URL: document.DocumentURL, Public: document.ShowToPublic != 0, CreatedAt: document.CreatedAt,
		})
	}
	return result
}

func applicationFeePaymentResponses(payments []stakeholderDB.TokenizationFeeProofOfPayment) []models.AssetApplicationFeePaymentResponse {
	result := make([]models.AssetApplicationFeePaymentResponse, 0, len(payments))
	for _, payment := range payments {
		result = append(result, models.AssetApplicationFeePaymentResponse{
			ID: payment.ID, PaymentMethodID: payment.TokenizationFeePaymentMethodID,
			TransactionReference: payment.TransactionReference, DocumentURL: payment.DocumentURL, CreatedAt: payment.CreatedAt,
		})
	}
	return result
}

func preferredFeeResponse(fee *stakeholderDB.WalletTokenizationFee) *models.PreferredTokenizationFee {
	if fee == nil {
		return nil
	}
	return &models.PreferredTokenizationFee{
		ID: fee.ID, Description: fee.FeeDescription, FiatPercentage: fee.FeeFiatPercentage,
		FiatCap: fee.FeeFiatCap, AssetPercentage: fee.FeeAssetPercentage,
	}
}

func walletBankName(bank *stakeholderDB.WalletBank) string {
	if bank == nil {
		return ""
	}
	return bank.BankName
}

func walletBankCountry(bank *stakeholderDB.WalletBank) string {
	if bank == nil {
		return ""
	}
	return bank.CountryCode
}

func tokenizationTimeline(asset stakeholderDB.TokenizedAsset, catalog []stakeholderDB.WalletTokenizationStatus) models.AssetTokenizationTimeline {
	current := uint64(0)
	if asset.AssetTokenizationStatus > 0 {
		current = uint64(asset.AssetTokenizationStatus)
	}
	timeline := models.AssetTokenizationTimeline{
		CurrentStatusID: current, DateSubmitted: nonZeroTime(asset.DateSubmitted), DateApproved: nonZeroTime(asset.DateOfApproval),
		MintingDate: nonZeroTime(asset.MintingDate), Stages: make([]models.AssetTokenizationTimelineStage, 0, len(catalog)),
	}
	for _, status := range catalog {
		timeline.Stages = append(timeline.Stages, models.AssetTokenizationTimelineStage{
			StatusID: status.ID, Description: status.Description, Reached: status.ID <= current, Current: status.ID == current,
		})
	}
	if len(timeline.Stages) == 0 {
		timeline.Stages = append(timeline.Stages, models.AssetTokenizationTimelineStage{StatusID: current, Reached: true, Current: true})
	}
	return timeline
}

func nonZeroTime(value time.Time) *time.Time {
	if value.IsZero() {
		return nil
	}
	copy := value
	return &copy
}

func (s *AssetService) visibleAsset(ctx context.Context, auth AuthContext, idOrCode string) (*stakeholderDB.TokenizedAsset, map[string]models.StakeholderAssetAssignment, error) {
	asset, err := s.client.GetAsset(ctx, idOrCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, NewHTTPError(http.StatusNotFound, "asset not found")
		}
		return nil, nil, err
	}
	assignments, err := s.assignmentsForAsset(ctx, asset.ID, asset.AssetCode)
	if err != nil {
		return nil, nil, err
	}
	assignmentByAsset, err := s.resolveAssignments(ctx, []stakeholderDB.TokenizedAsset{*asset}, assignments)
	if err != nil {
		return nil, nil, err
	}
	if !s.assetVisibleToAuth(auth, *asset, assignmentByAsset) {
		return nil, nil, NewHTTPError(http.StatusNotFound, "asset not found")
	}
	return asset, assignmentByAsset, nil
}

func (s *AssetService) ListAllAssets(ctx context.Context, auth AuthContext) ([]models.TokenizedAssetResponse, error) {
	assignments, err := s.assignmentsForRole(ctx, auth)
	if err != nil {
		return nil, err
	}
	filters, err := s.assetFiltersForRole(auth, assignments, models.AssetListFilters{})
	if err != nil {
		return nil, err
	}
	assets, err := s.client.ListAllAssets(ctx, filters)
	if err != nil {
		return nil, err
	}
	pageAssignments, err := s.assignmentsForAssets(ctx, assets)
	if err != nil {
		return nil, err
	}
	assignmentByAsset, err := s.resolveAssignments(ctx, assets, pageAssignments)
	if err != nil {
		return nil, err
	}
	responses := make([]models.TokenizedAssetResponse, 0, len(assets))
	for _, asset := range assets {
		if s.assetVisibleToAuth(auth, asset, assignmentByAsset) {
			responses = append(responses, s.assetResponse(auth, asset, assignmentByAsset))
		}
	}
	return responses, nil
}

func (s *AssetService) RequireAssignmentForAsset(ctx context.Context, auth AuthContext, assetID string) (*models.StakeholderAssetAssignment, *models.TokenizedAssetResponse, error) {
	asset, err := s.GetAsset(ctx, auth, assetID)
	if err != nil {
		return nil, nil, err
	}
	if asset.Assignment == nil {
		return nil, nil, NewHTTPError(http.StatusUnprocessableEntity, "asset has no stakeholder assignment")
	}
	return asset.Assignment, asset, nil
}

func (s *AssetService) requirePortalOrganization(orgID string, requiredRole string) (*coreModels.Organization, error) {
	orgID = strings.TrimSpace(orgID)
	if orgID == "" {
		return nil, NewHTTPError(http.StatusBadRequest, "organization_id is required")
	}
	var org coreModels.Organization
	if err := s.adminDB.First(&org, "id = ?", orgID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewHTTPError(http.StatusBadRequest, "organization not found")
		}
		return nil, err
	}
	if org.StakeholderID == nil || org.StakeholderType == nil {
		return nil, NewHTTPError(http.StatusBadRequest, "organization has no stakeholder linkage")
	}
	role, ok := models.DashboardRoleForStakeholderType(*org.StakeholderType)
	if !ok || role != requiredRole {
		return nil, NewHTTPError(http.StatusBadRequest, "organization has invalid stakeholder role for assignment")
	}
	return &org, nil
}

func (s *AssetService) assignmentsForRole(ctx context.Context, auth AuthContext) ([]models.StakeholderAssetAssignment, error) {
	query := s.adminDB.WithContext(ctx).Where("status = ?", models.AssignmentStatusActive)
	switch auth.DashboardRole {
	case models.DashboardRoleTrustee:
		query = query.Where("trustee_org_id = ? OR trustee_stakeholder_id = ?", auth.OrganizationID, auth.StakeholderID)
	case models.DashboardRoleAssetManager:
		query = query.Where("asset_manager_org_id = ? OR asset_manager_stakeholder_id = ?", auth.OrganizationID, auth.StakeholderID)
	case models.DashboardRoleAssetCustodian:
		query = query.Where("custodian_org_id = ? OR custodian_stakeholder_id = ?", auth.OrganizationID, auth.StakeholderID)
	default:
		// FK-only roles (legal/financial adviser, issuing house, rating agency)
		// have no columns on stakeholder_asset_assignments — they are scoped
		// purely via the tokenized_assets FK column (see assetVisibleToAuth).
		if models.IsSupportedDashboardRole(auth.DashboardRole) {
			return []models.StakeholderAssetAssignment{}, nil
		}
		return nil, NewHTTPError(http.StatusForbidden, "unsupported stakeholder role")
	}
	var assignments []models.StakeholderAssetAssignment
	if err := query.Find(&assignments).Error; err != nil {
		return nil, err
	}
	return assignments, nil
}

func (s *AssetService) assignmentsForAsset(ctx context.Context, assetID, assetCode string) ([]models.StakeholderAssetAssignment, error) {
	var assignments []models.StakeholderAssetAssignment
	query := s.adminDB.WithContext(ctx)
	if strings.TrimSpace(assetCode) != "" {
		query = query.Where("asset_id = ? OR asset_code = ?", assetID, assetCode)
	} else {
		query = query.Where("asset_id = ?", assetID)
	}
	if err := query.Find(&assignments).Error; err != nil {
		return nil, err
	}
	return assignments, nil
}

func (s *AssetService) assignmentsForAssets(ctx context.Context, assets []stakeholderDB.TokenizedAsset) ([]models.StakeholderAssetAssignment, error) {
	assetIDs := make([]string, 0, len(assets))
	assetCodes := make([]string, 0, len(assets))
	for _, asset := range assets {
		assetIDs = append(assetIDs, asset.ID)
		if strings.TrimSpace(asset.AssetCode) != "" {
			assetCodes = append(assetCodes, asset.AssetCode)
		}
	}
	if len(assetIDs) == 0 {
		return []models.StakeholderAssetAssignment{}, nil
	}
	query := s.adminDB.WithContext(ctx).Where("asset_id IN ?", assetIDs)
	if len(assetCodes) > 0 {
		query = query.Or("asset_code IN ?", assetCodes)
	}
	var assignments []models.StakeholderAssetAssignment
	if err := query.Find(&assignments).Error; err != nil {
		return nil, err
	}
	return assignments, nil
}

func (s *AssetService) assignmentForAsset(ctx context.Context, assetID string) (*models.StakeholderAssetAssignment, error) {
	var assignment models.StakeholderAssetAssignment
	if err := s.adminDB.WithContext(ctx).Where("asset_id = ? AND status = ?", assetID, models.AssignmentStatusActive).First(&assignment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewHTTPError(http.StatusUnprocessableEntity, "asset has no stakeholder assignment")
		}
		return nil, err
	}
	return &assignment, nil
}

func (s *AssetService) assetFiltersForRole(auth AuthContext, assignments []models.StakeholderAssetAssignment, filters models.AssetListFilters) (stakeholderDB.AssetFilters, error) {
	dateFrom, err := ParseDate(filters.DateFrom)
	if err != nil {
		return stakeholderDB.AssetFilters{}, err
	}
	dateTo, err := ParseDate(filters.DateTo)
	if err != nil {
		return stakeholderDB.AssetFilters{}, err
	}
	status, err := stakeholderDB.ParseOptionalStatus(filters.Status)
	if err != nil {
		return stakeholderDB.AssetFilters{}, NewHTTPError(http.StatusBadRequest, "invalid asset status")
	}
	assetIDs, assetCodes := assignmentAssetIDsAndCodes(assignments)
	result := stakeholderDB.AssetFilters{
		Page:       filters.Page,
		Limit:      filters.Limit,
		Status:     status,
		AssetType:  filters.AssetType,
		Sector:     filters.Sector,
		Search:     filters.Search,
		DateFrom:   dateFrom,
		DateTo:     dateTo,
		AssetIDs:   assetIDs,
		AssetCodes: assetCodes,
	}
	switch auth.DashboardRole {
	case models.DashboardRoleTrustee:
		result.TrusteeStakeholderID = &auth.StakeholderID
	case models.DashboardRoleAssetManager:
		result.ManagerStakeholderID = &auth.StakeholderID
	case models.DashboardRoleAssetCustodian:
		result.CustodianStakeholderID = &auth.StakeholderID
	default:
		// FK-only roles are scoped by their tokenized_assets FK column.
		switch models.AssetFKColumnForRole(auth.DashboardRole) {
		case models.AssetFKColumnLegalAdviser:
			result.LegalAdviserID = &auth.StakeholderID
		case models.AssetFKColumnFinancial:
			result.FinancialAdviserID = &auth.StakeholderID
		case models.AssetFKColumnIssuingHouse:
			result.IssuingHouseID = &auth.StakeholderID
		case models.AssetFKColumnRatingAgency:
			result.RatingAgencyID = &auth.StakeholderID
		default:
			return stakeholderDB.AssetFilters{}, NewHTTPError(http.StatusForbidden, "unsupported stakeholder role")
		}
	}
	return result, nil
}

func (s *AssetService) assetVisibleToAuth(auth AuthContext, asset stakeholderDB.TokenizedAsset, assignmentByAsset map[string]models.StakeholderAssetAssignment) bool {
	// FK-only roles have no columns on stakeholder_asset_assignments, so the
	// synthesized wallet assignment above never carries them. They are scoped
	// purely by the asset's own FK column against the caller's stakeholder id.
	if !models.RoleHasAssignmentColumns(auth.DashboardRole) && models.IsSupportedDashboardRole(auth.DashboardRole) {
		return s.assetVisibleByFK(auth, asset)
	}

	assignment, hasAssignment := findAssignmentForAsset(asset, assignmentByAsset)
	if hasAssignment {
		if assignment.Status != models.AssignmentStatusActive {
			return false
		}
		switch auth.DashboardRole {
		case models.DashboardRoleTrustee:
			return (assignment.TrusteeOrgID != nil && *assignment.TrusteeOrgID == auth.OrganizationID) ||
				(assignment.TrusteeStakeholderID != nil && *assignment.TrusteeStakeholderID == auth.StakeholderID)
		case models.DashboardRoleAssetManager:
			return (assignment.AssetManagerOrgID != nil && *assignment.AssetManagerOrgID == auth.OrganizationID) ||
				(assignment.AssetManagerStakeholderID != nil && *assignment.AssetManagerStakeholderID == auth.StakeholderID)
		case models.DashboardRoleAssetCustodian:
			return (assignment.CustodianOrgID != nil && *assignment.CustodianOrgID == auth.OrganizationID) ||
				(assignment.CustodianStakeholderID != nil && *assignment.CustodianStakeholderID == auth.StakeholderID)
		default:
			return false
		}
	}
	switch auth.DashboardRole {
	case models.DashboardRoleTrustee:
		return asset.TrusteeID == auth.StakeholderID
	case models.DashboardRoleAssetManager:
		return asset.AssetManagerID == auth.StakeholderID
	case models.DashboardRoleAssetCustodian:
		return asset.ApprovedAssetCustodianID == auth.StakeholderID
	default:
		return s.assetVisibleByFK(auth, asset)
	}
}

// assetVisibleByFK reports whether the asset's stakeholder FK column for the
// caller's role matches the caller's stakeholder id. A zero id never matches
// (default-0 unassigned column).
func (s *AssetService) assetVisibleByFK(auth AuthContext, asset stakeholderDB.TokenizedAsset) bool {
	if auth.StakeholderID == 0 {
		return false
	}
	switch models.AssetFKColumnForRole(auth.DashboardRole) {
	case models.AssetFKColumnTrustee:
		return asset.TrusteeID == auth.StakeholderID
	case models.AssetFKColumnAssetManager:
		return asset.AssetManagerID == auth.StakeholderID
	case models.AssetFKColumnCustodian:
		return asset.ApprovedAssetCustodianID == auth.StakeholderID
	case models.AssetFKColumnLegalAdviser:
		return asset.LegalAdviserID == auth.StakeholderID
	case models.AssetFKColumnFinancial:
		return asset.FinancialAdviserID == auth.StakeholderID
	case models.AssetFKColumnIssuingHouse:
		return asset.AssetIssuingHouseID == auth.StakeholderID
	case models.AssetFKColumnRatingAgency:
		return asset.RatingAgencyID == auth.StakeholderID
	default:
		return false
	}
}

func (s *AssetService) assetResponse(auth AuthContext, asset stakeholderDB.TokenizedAsset, assignmentByAsset map[string]models.StakeholderAssetAssignment) models.TokenizedAssetResponse {
	assignment, hasAssignment := findAssignmentForAsset(asset, assignmentByAsset)
	var assignmentPtr *models.StakeholderAssetAssignment
	if hasAssignment {
		assignmentCopy := assignment
		assignmentPtr = &assignmentCopy
	}
	return models.TokenizedAssetResponse{
		ID:                    asset.ID,
		AssetCode:             asset.AssetCode,
		AssetName:             asset.AssetName,
		AssetType:             asset.AssetType,
		AssetSector:           asset.AssetSector,
		AssetSubSector:        asset.AssetSubSector,
		Status:                asset.AssetTokenizationStatus,
		AssetQuoteCurrency:    asset.AssetQuoteCurrency,
		AssetCurrentValue:     asset.AssetCurrentValue,
		ValueOfTokenizedAsset: asset.ValueOfTokenizedAsset,
		NumberOfTokensIssued:  asset.NumberOfTokenToBeIssued,
		AssetManagerID:        asset.AssetManagerID,
		ApprovedCustodianID:   asset.ApprovedAssetCustodianID,
		TrusteeID:             asset.TrusteeID,
		TokenHolderCount:      asset.TokenHolderCount,
		ComplianceStatus:      assetComplianceStatus(asset),
		CustodyStatus:         assetCustodyStatus(asset.AssetTokenizationStatus),
		CustodianFeeValue:     asset.CustodianFeeValue,
		ConfiguredFeeValue:    asset.ConfiguredFeeValue(),
		CreatedAt:             asset.CreatedAt,
		UpdatedAt:             asset.UpdatedAt,
		Actions:               actionFlagsForRole(auth.DashboardRole),
		Assignment:            assignmentPtr,
	}
}

func (s *AssetService) assetStatistics(ctx context.Context, auth AuthContext, filters stakeholderDB.AssetFilters) (models.AssetStatistics, error) {
	filters.Page, filters.Limit = 0, 0
	assets, err := s.client.ListAllAssets(ctx, filters)
	if err != nil {
		return models.AssetStatistics{}, err
	}
	assignments, err := s.assignmentsForAssets(ctx, assets)
	if err != nil {
		return models.AssetStatistics{}, err
	}
	assignmentByAsset, err := s.resolveAssignments(ctx, assets, assignments)
	if err != nil {
		return models.AssetStatistics{}, err
	}
	var result models.AssetStatistics
	for _, asset := range assets {
		if !s.assetVisibleToAuth(auth, asset, assignmentByAsset) {
			continue
		}
		result.Total++
		switch assetCustodyStatus(asset.AssetTokenizationStatus) {
		case models.AssetCustodyStatusActive:
			result.Active++
		case models.AssetCustodyStatusLiquidated:
			result.Liquidated++
		}
	}
	return result, nil
}

func assetComplianceStatus(asset stakeholderDB.TokenizedAsset) string {
	if asset.DueDiligenceFail != 0 {
		return models.AssetComplianceStatusRejected
	}
	if asset.VettingStatus != 0 {
		return models.AssetComplianceStatusVerified
	}
	return models.AssetComplianceStatusPending
}

func assetCustodyStatus(status int) string {
	switch status {
	case models.AssetTokenizationStatusPrimarySale, models.AssetTokenizationStatusSecondarySale:
		return models.AssetCustodyStatusActive
	case models.AssetTokenizationStatusLiquidated:
		return models.AssetCustodyStatusLiquidated
	default:
		return models.AssetCustodyStatusInactive
	}
}

func (s *AssetService) assignedStakeholders(ctx context.Context, asset stakeholderDB.TokenizedAsset) ([]models.AssignedStakeholderResponse, error) {
	type stakeholderRef struct {
		role, stakeholderType string
		id                    uint64
	}
	refs := []stakeholderRef{
		{models.DashboardRoleAssetManager, models.StakeholderTypeAssetManager, asset.AssetManagerID},
		{models.DashboardRoleAssetCustodian, models.StakeholderTypeAssetCustodian, asset.ApprovedAssetCustodianID},
		{models.DashboardRoleTrustee, models.StakeholderTypeTrustee, asset.TrusteeID},
		{models.DashboardRoleLegalAdviser, models.StakeholderTypeLegalAdviser, asset.LegalAdviserID},
		{models.DashboardRoleFinancialAdviser, models.StakeholderTypeFinancialAdviser, asset.FinancialAdviserID},
		{models.DashboardRoleIssuingHouse, models.StakeholderTypeIssuingHouse, asset.AssetIssuingHouseID},
		{models.DashboardRoleRatingAgency, models.StakeholderTypeRatingAgency, asset.RatingAgencyID},
	}
	ids := make([]uint64, 0, len(refs))
	for _, ref := range refs {
		if ref.id > 0 {
			ids = append(ids, ref.id)
		}
	}
	organizations := make([]coreModels.Organization, 0)
	if len(ids) > 0 {
		if err := s.adminDB.WithContext(ctx).Where("stakeholder_id IN ?", ids).Find(&organizations).Error; err != nil {
			return nil, err
		}
	}
	organizationByKey := make(map[string]coreModels.Organization, len(organizations))
	for _, organization := range organizations {
		if organization.StakeholderID != nil && organization.StakeholderType != nil {
			organizationByKey[stakeholderKey(*organization.StakeholderType, *organization.StakeholderID)] = organization
		}
	}
	result := make([]models.AssignedStakeholderResponse, 0, len(refs))
	for _, ref := range refs {
		if ref.id == 0 {
			continue
		}
		item := models.AssignedStakeholderResponse{Role: ref.role, StakeholderID: ref.id, StakeholderType: ref.stakeholderType}
		if organization, ok := organizationByKey[stakeholderKey(ref.stakeholderType, ref.id)]; ok {
			organizationID := organization.ID
			item.OrganizationID = &organizationID
			item.OrganizationName = organization.Name
			item.OrganizationEmail = organization.Email
		}
		result = append(result, item)
	}
	return result, nil
}

func (s *AssetService) assetDocuments(ctx context.Context, auth AuthContext, assetID string) ([]models.StakeholderDocument, error) {
	documents := make([]models.StakeholderDocument, 0)
	var candidates []models.StakeholderDocument
	if err := s.adminDB.WithContext(ctx).
		Where("asset_id = ? AND status = ?", assetID, models.DocumentStatusActive).
		Order("created_at desc").Find(&candidates).Error; err != nil {
		return nil, err
	}
	for _, document := range candidates {
		if models.RoleCanAccessDocument(auth.DashboardRole, []string(document.AccessRoles)) {
			documents = append(documents, document)
		}
	}
	return documents, nil
}

func splitCountries(value string) []string {
	result := make([]string, 0)
	for _, country := range strings.FieldsFunc(value, func(r rune) bool { return r == ',' || r == ';' }) {
		if country = strings.TrimSpace(country); country != "" {
			result = append(result, country)
		}
	}
	return result
}

func (s *AssetService) resolveAssignments(ctx context.Context, assets []stakeholderDB.TokenizedAsset, explicit []models.StakeholderAssetAssignment) (map[string]models.StakeholderAssetAssignment, error) {
	result := assignmentMap(explicit)
	stakeholderIDs := make([]uint64, 0, len(assets)*3)
	for _, asset := range assets {
		for _, id := range []uint64{asset.AssetManagerID, asset.ApprovedAssetCustodianID, asset.TrusteeID} {
			if id > 0 {
				stakeholderIDs = append(stakeholderIDs, id)
			}
		}
	}

	organizations := make([]coreModels.Organization, 0)
	if len(stakeholderIDs) > 0 {
		if err := s.adminDB.WithContext(ctx).Where("stakeholder_id IN ?", stakeholderIDs).Find(&organizations).Error; err != nil {
			return nil, err
		}
	}
	organizationByStakeholder := make(map[string]string, len(organizations))
	for _, organization := range organizations {
		if organization.StakeholderID == nil || organization.StakeholderType == nil {
			continue
		}
		organizationByStakeholder[stakeholderKey(*organization.StakeholderType, *organization.StakeholderID)] = organization.ID
	}

	for _, asset := range assets {
		assignment, found := findAssignmentForAsset(asset, result)
		if !found {
			assignment = models.StakeholderAssetAssignment{
				ID:         "wallet:" + asset.ID,
				AssetID:    asset.ID,
				AssetCode:  asset.AssetCode,
				Status:     models.AssignmentStatusActive,
				AssignedBy: "wallet_vetting",
			}
		} else if assignment.Status != models.AssignmentStatusActive {
			continue
		}
		mergeWalletStakeholder(&assignment.AssetManagerStakeholderID, &assignment.AssetManagerOrgID, asset.AssetManagerID, models.StakeholderTypeAssetManager, organizationByStakeholder)
		mergeWalletStakeholder(&assignment.CustodianStakeholderID, &assignment.CustodianOrgID, asset.ApprovedAssetCustodianID, models.StakeholderTypeAssetCustodian, organizationByStakeholder)
		mergeWalletStakeholder(&assignment.TrusteeStakeholderID, &assignment.TrusteeOrgID, asset.TrusteeID, models.StakeholderTypeTrustee, organizationByStakeholder)
		result["id:"+asset.ID] = assignment
		if asset.AssetCode != "" {
			result["code:"+asset.AssetCode] = assignment
		}
	}
	return result, nil
}

func mergeWalletStakeholder(stakeholderID **uint64, organizationID **string, walletID uint64, stakeholderType string, organizations map[string]string) {
	if walletID == 0 {
		return
	}
	if *stakeholderID == nil {
		id := walletID
		*stakeholderID = &id
	}
	if *organizationID == nil {
		if orgID, ok := organizations[stakeholderKey(stakeholderType, walletID)]; ok {
			id := orgID
			*organizationID = &id
		}
	}
}

func stakeholderKey(stakeholderType string, stakeholderID uint64) string {
	return strings.TrimSpace(stakeholderType) + ":" + fmt.Sprint(stakeholderID)
}

func assignmentMap(assignments []models.StakeholderAssetAssignment) map[string]models.StakeholderAssetAssignment {
	result := make(map[string]models.StakeholderAssetAssignment, len(assignments)*2)
	for _, assignment := range assignments {
		if assignment.AssetID != "" {
			result["id:"+assignment.AssetID] = assignment
		}
		if assignment.AssetCode != "" {
			result["code:"+assignment.AssetCode] = assignment
		}
	}
	return result
}

func findAssignmentForAsset(asset stakeholderDB.TokenizedAsset, assignments map[string]models.StakeholderAssetAssignment) (models.StakeholderAssetAssignment, bool) {
	if assignment, ok := assignments["id:"+asset.ID]; ok {
		return assignment, true
	}
	if assignment, ok := assignments["code:"+asset.AssetCode]; ok {
		return assignment, true
	}
	return models.StakeholderAssetAssignment{}, false
}

func assignmentAssetIDsAndCodes(assignments []models.StakeholderAssetAssignment) ([]string, []string) {
	ids := make([]string, 0, len(assignments))
	codes := make([]string, 0, len(assignments))
	for _, assignment := range assignments {
		if strings.TrimSpace(assignment.AssetID) != "" {
			ids = append(ids, assignment.AssetID)
		}
		if strings.TrimSpace(assignment.AssetCode) != "" {
			codes = append(codes, assignment.AssetCode)
		}
	}
	return ids, codes
}

func actionFlagsForRole(role string) models.AssetActionFlags {
	switch role {
	case models.DashboardRoleTrustee:
		return models.AssetActionFlags{
			CanApproveDueDiligence:   true,
			CanApproveFundRelease:    true,
			CanAuthorizeDistribution: true,
		}
	case models.DashboardRoleAssetCustodian:
		return models.AssetActionFlags{
			CanExecuteFundRelease: true,
			CanUpdateAccount:      true,
		}
	case models.DashboardRoleAssetManager:
		return models.AssetActionFlags{
			CanRequestFundRelease: true,
			CanRecordRevenue:      true,
			CanSubmitValuation:    true,
		}
	default:
		return models.AssetActionFlags{}
	}
}

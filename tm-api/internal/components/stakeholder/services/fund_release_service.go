package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"admin-panel-dashboard/internal/components/stakeholder/models"
	coreModels "admin-panel-dashboard/internal/models"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FundReleaseService struct {
	db             *gorm.DB
	assets         *AssetService
	audit          *AuditService
	notifications  *NotificationService
	authorizations *AuthorizationService
	documents      *DocumentService
}

func NewFundReleaseService(db *gorm.DB, assets *AssetService, audit *AuditService, notifications *NotificationService, authorizations *AuthorizationService, documents *DocumentService) *FundReleaseService {
	return &FundReleaseService{db: db, assets: assets, audit: audit, notifications: notifications, authorizations: authorizations, documents: documents}
}

func CanTransitionFundRelease(from, to string) bool {
	allowed := map[string][]string{
		models.FundReleaseStatusDraft:            {models.FundReleaseStatusSubmitted},
		models.FundReleaseStatusSubmitted:        {models.FundReleaseStatusTrusteeApproved, models.FundReleaseStatusExecutionPending, models.FundReleaseStatusTrusteeRejected},
		models.FundReleaseStatusTrusteeApproved:  {models.FundReleaseStatusExecutionPending},
		models.FundReleaseStatusExecutionPending: {models.FundReleaseStatusProcessing, models.FundReleaseStatusCompleted, models.FundReleaseStatusFailed},
		models.FundReleaseStatusProcessing:       {models.FundReleaseStatusCompleted, models.FundReleaseStatusFailed},
	}
	for _, status := range allowed[from] {
		if status == to {
			return true
		}
	}
	return false
}

func (s *FundReleaseService) Create(ctx context.Context, auth AuthContext, req models.CreateFundReleaseRequest) (*models.FundReleaseRequest, error) {
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
	if s.documents == nil {
		return nil, NewHTTPError(http.StatusServiceUnavailable, "document service is not configured")
	}
	if err := s.documents.ValidateForUse(ctx, auth, asset.ID, models.DocumentCategoryFundReleaseSupporting, req.SupportingDocumentIDs); err != nil {
		return nil, err
	}
	receivingBank, receivingAccountName, receivingAccountNumber, err := validateReceivingAccount(req)
	if err != nil {
		return nil, err
	}
	custodianOrgID := ""
	if assignment.CustodianOrgID != nil {
		custodianOrgID = *assignment.CustodianOrgID
	}
	now := time.Now()
	request := models.FundReleaseRequest{
		ID:                     uuid.NewString(),
		AssetID:                asset.ID,
		AssetCode:              asset.AssetCode,
		RequesterOrgID:         auth.OrganizationID,
		RequesterMemberID:      auth.MemberID,
		TrusteeOrgID:           *assignment.TrusteeOrgID,
		CustodianOrgID:         custodianOrgID,
		Amount:                 amount,
		Currency:               currency,
		Purpose:                strings.TrimSpace(req.Purpose),
		Status:                 models.FundReleaseStatusSubmitted,
		SupportingDocumentIDs:  models.JSONStringArray(req.SupportingDocumentIDs),
		ReceivingBank:          receivingBank,
		ReceivingAccountName:   receivingAccountName,
		ReceivingAccountNumber: receivingAccountNumber,
		CreatedAt:              now,
		UpdatedAt:              now,
	}
	if request.Purpose == "" {
		return nil, NewHTTPError(http.StatusBadRequest, "purpose is required")
	}
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&request).Error; err != nil {
			return err
		}
		if s.notifications != nil {
			if err := s.notifications.Create(ctx, tx, NotificationInput{
				RecipientOrgID:    request.TrusteeOrgID,
				SenderOrgID:       &auth.OrganizationID,
				Type:              "fund_release.submitted",
				Title:             "Fund release request submitted",
				Message:           fmt.Sprintf("A fund release request for %s %s is pending trustee review.", request.Currency, request.Amount.String()),
				RelatedEntityType: models.EntityFundReleaseRequest,
				RelatedEntityID:   request.ID,
			}); err != nil {
				return err
			}
		}
		if s.audit != nil {
			event := AuditEventFromAuth(auth, models.ActionFundReleaseCreate, models.EntityFundReleaseRequest, request.ID)
			event.AfterState = models.JSONMap{"status": request.Status, "amount": request.Amount.String(), "currency": request.Currency}
			if err := s.audit.Record(ctx, tx, event); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &request, nil
}

func (s *FundReleaseService) List(ctx context.Context, auth AuthContext, page, limit int, status, assetID string) ([]models.FundReleaseListItemResponse, int64, error) {
	query := s.db.WithContext(ctx).Model(&models.FundReleaseRequest{})
	switch auth.DashboardRole {
	case models.DashboardRoleAssetManager:
		query = query.Where("requester_org_id = ?", auth.OrganizationID)
	case models.DashboardRoleTrustee:
		query = query.Where("trustee_org_id = ?", auth.OrganizationID)
	case models.DashboardRoleAssetCustodian:
		query = query.Where("custodian_org_id = ?", auth.OrganizationID)
		if strings.TrimSpace(status) == "" {
			query = query.Where("status IN ?", []string{models.FundReleaseStatusExecutionPending, models.FundReleaseStatusProcessing, models.FundReleaseStatusCompleted, models.FundReleaseStatusFailed})
		}
	default:
		return nil, 0, NewHTTPError(http.StatusForbidden, "unsupported stakeholder role")
	}
	if strings.TrimSpace(status) != "" {
		query = query.Where("status = ?", strings.TrimSpace(status))
	}
	if strings.TrimSpace(assetID) != "" {
		asset, err := s.assets.GetAsset(ctx, auth, assetID)
		if err != nil {
			return nil, 0, err
		}
		query = query.Where("asset_id = ?", asset.ID)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var records []models.FundReleaseRequest
	if err := query.Order("created_at desc").Offset((page - 1) * limit).Limit(limit).Find(&records).Error; err != nil {
		return nil, 0, err
	}
	enriched, err := s.enrichFundReleaseList(ctx, records)
	if err != nil {
		return nil, 0, err
	}
	return enriched, total, nil
}

func (s *FundReleaseService) GetDetails(ctx context.Context, auth AuthContext, id string) (*models.FundReleaseDetailResponse, error) {
	request, err := s.Get(ctx, auth, id)
	if err != nil {
		return nil, err
	}
	requester, err := s.requesterDetails(ctx, *request)
	if err != nil {
		return nil, err
	}
	documents, err := s.supportingDocumentDetails(ctx, auth, request.SupportingDocumentIDs)
	if err != nil {
		return nil, err
	}
	return &models.FundReleaseDetailResponse{
		FundReleaseRequest:  *request,
		Requester:           requester,
		SupportingDocuments: documents,
		ReceivingAccount: models.FundReleaseReceivingAccountResponse{
			Bank: request.ReceivingBank, AccountName: request.ReceivingAccountName, AccountNumber: request.ReceivingAccountNumber,
		},
	}, nil
}

func (s *FundReleaseService) AssetManagerSummary(ctx context.Context, auth AuthContext, assetID string) (*models.AssetManagerFundSummary, error) {
	if auth.DashboardRole != models.DashboardRoleAssetManager {
		return nil, NewHTTPError(http.StatusForbidden, "asset manager role is required")
	}
	query := s.db.WithContext(ctx).Where("requester_org_id = ?", auth.OrganizationID)
	if strings.TrimSpace(assetID) != "" {
		asset, err := s.assets.GetAsset(ctx, auth, assetID)
		if err != nil {
			return nil, err
		}
		query = query.Where("asset_id = ?", asset.ID)
	}
	var requests []models.FundReleaseRequest
	if err := query.Find(&requests).Error; err != nil {
		return nil, err
	}
	type aggregate struct {
		count   int
		amounts map[string]decimal.Decimal
	}
	newAggregate := func() aggregate { return aggregate{amounts: make(map[string]decimal.Decimal)} }
	requested, approved, released, pending, rejected := newAggregate(), newAggregate(), newAggregate(), newAggregate(), newAggregate()
	add := func(target *aggregate, request models.FundReleaseRequest) {
		target.count++
		target.amounts[request.Currency] = target.amounts[request.Currency].Add(request.Amount)
	}
	for _, request := range requests {
		if request.Status != models.FundReleaseStatusDraft {
			add(&requested, request)
		}
		switch request.Status {
		case models.FundReleaseStatusTrusteeApproved, models.FundReleaseStatusExecutionPending, models.FundReleaseStatusProcessing, models.FundReleaseStatusCompleted:
			add(&approved, request)
		}
		if request.Status == models.FundReleaseStatusCompleted {
			add(&released, request)
		}
		switch request.Status {
		case models.FundReleaseStatusSubmitted, models.FundReleaseStatusTrusteeApproved, models.FundReleaseStatusExecutionPending, models.FundReleaseStatusProcessing:
			add(&pending, request)
		case models.FundReleaseStatusTrusteeRejected:
			add(&rejected, request)
		}
	}
	remaining := make(map[string]decimal.Decimal, len(approved.amounts))
	for currency, amount := range approved.amounts {
		balance := amount.Sub(released.amounts[currency])
		if balance.IsNegative() {
			balance = decimal.Zero
		}
		remaining[currency] = balance
	}
	responseAggregate := func(value aggregate) models.FundReleaseAggregate {
		return models.FundReleaseAggregate{Count: value.count, Amounts: currencyTotals(value.amounts)}
	}
	return &models.AssetManagerFundSummary{
		Requested: responseAggregate(requested), Approved: responseAggregate(approved), Released: responseAggregate(released),
		Pending: responseAggregate(pending), Rejected: responseAggregate(rejected), Remaining: currencyTotals(remaining),
	}, nil
}

func validateReceivingAccount(req models.CreateFundReleaseRequest) (string, string, string, error) {
	bank := strings.TrimSpace(req.ReceivingBank)
	accountName := strings.TrimSpace(req.ReceivingAccountName)
	accountNumber := strings.TrimSpace(req.ReceivingAccountNumber)
	provided := 0
	for _, value := range []string{bank, accountName, accountNumber} {
		if value != "" {
			provided++
		}
	}
	if provided != 0 && provided != 3 {
		return "", "", "", NewHTTPError(http.StatusBadRequest, "receiving_bank, receiving_account_name, and receiving_account_number must be provided together")
	}
	if len(bank) > 128 || len(accountName) > 160 || len(accountNumber) > 64 {
		return "", "", "", NewHTTPError(http.StatusBadRequest, "receiving account details exceed allowed length")
	}
	return bank, accountName, accountNumber, nil
}

func (s *FundReleaseService) requesterDetails(ctx context.Context, request models.FundReleaseRequest) (models.FundReleaseRequesterResponse, error) {
	result := models.FundReleaseRequesterResponse{MemberID: request.RequesterMemberID, OrganizationID: request.RequesterOrgID}
	var member coreModels.OrganizationMember
	if err := s.db.WithContext(ctx).First(&member, "id = ? AND organization_id = ?", request.RequesterMemberID, request.RequesterOrgID).Error; err == nil {
		result.FirstName, result.LastName, result.Email = member.FirstName, member.LastName, member.Email
		result.Name = organizationMemberDisplayName(member)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return result, err
	}
	var organization coreModels.Organization
	if err := s.db.WithContext(ctx).First(&organization, "id = ?", request.RequesterOrgID).Error; err == nil {
		result.OrganizationName = organization.Name
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return result, err
	}
	return result, nil
}

func (s *FundReleaseService) enrichFundReleaseList(ctx context.Context, records []models.FundReleaseRequest) ([]models.FundReleaseListItemResponse, error) {
	result := make([]models.FundReleaseListItemResponse, len(records))
	if len(records) == 0 {
		return result, nil
	}

	memberIDs := make([]string, 0, len(records)*2)
	organizationIDs := make([]string, 0, len(records)*2)
	for _, record := range records {
		memberIDs = append(memberIDs, record.RequesterMemberID)
		organizationIDs = append(organizationIDs, record.RequesterOrgID)
		if record.ReviewedByMemberID != nil && strings.TrimSpace(*record.ReviewedByMemberID) != "" {
			memberIDs = append(memberIDs, *record.ReviewedByMemberID)
			organizationIDs = append(organizationIDs, record.TrusteeOrgID)
		}
	}
	membersByID, organizationsByID, err := loadOrganizationActors(ctx, s.db, memberIDs, organizationIDs)
	if err != nil {
		return nil, err
	}

	for index, record := range records {
		requester := fundReleaseActor(record.RequesterMemberID, record.RequesterOrgID, membersByID, organizationsByID)
		var reviewer *models.FundReleaseRequesterResponse
		if record.ReviewedByMemberID != nil && strings.TrimSpace(*record.ReviewedByMemberID) != "" {
			resolved := fundReleaseActor(*record.ReviewedByMemberID, record.TrusteeOrgID, membersByID, organizationsByID)
			reviewer = &resolved
		}
		result[index] = models.FundReleaseListItemResponse{FundReleaseRequest: record, Requester: requester, Reviewer: reviewer}
	}
	return result, nil
}

func fundReleaseActor(memberID, fallbackOrganizationID string, membersByID map[string]coreModels.OrganizationMember, organizationsByID map[string]coreModels.Organization) models.FundReleaseRequesterResponse {
	result := models.FundReleaseRequesterResponse{MemberID: memberID, OrganizationID: fallbackOrganizationID}
	if member, ok := membersByID[memberID]; ok && member.OrganizationID == fallbackOrganizationID {
		result.Name = organizationMemberDisplayName(member)
		result.FirstName = member.FirstName
		result.LastName = member.LastName
		result.Email = member.Email
	}
	if organization, ok := organizationsByID[result.OrganizationID]; ok {
		result.OrganizationName = organization.Name
	}
	return result
}

func (s *FundReleaseService) supportingDocumentDetails(ctx context.Context, auth AuthContext, ids models.JSONStringArray) ([]models.FundReleaseDocumentResponse, error) {
	result := make([]models.FundReleaseDocumentResponse, 0, len(ids))
	if len(ids) == 0 {
		return result, nil
	}
	var documents []models.StakeholderDocument
	if err := s.db.WithContext(ctx).Where("id IN ? AND status = ?", []string(ids), models.DocumentStatusActive).Find(&documents).Error; err != nil {
		return nil, err
	}
	byID := make(map[string]models.StakeholderDocument, len(documents))
	for _, document := range documents {
		byID[document.ID] = document
	}
	for _, id := range ids {
		document, ok := byID[id]
		if !ok || !models.RoleCanAccessDocument(auth.DashboardRole, []string(document.AccessRoles)) {
			continue
		}
		item := models.FundReleaseDocumentResponse{
			ID: document.ID, Title: document.Title, Category: document.Category,
			OriginalFilename: document.OriginalFilename, MimeType: document.MimeType, SizeBytes: document.SizeBytes,
			ExternalURL: document.FileURL,
		}
		if document.StorageObjectID != "" {
			item.DownloadPath = "/api/v1/stakeholder/shared/documents/" + document.ID + "/download"
		}
		result = append(result, item)
	}
	return result, nil
}

func (s *FundReleaseService) Get(ctx context.Context, auth AuthContext, id string) (*models.FundReleaseRequest, error) {
	var request models.FundReleaseRequest
	query := s.db.WithContext(ctx).Where("id = ?", id)
	switch auth.DashboardRole {
	case models.DashboardRoleAssetManager:
		query = query.Where("requester_org_id = ?", auth.OrganizationID)
	case models.DashboardRoleTrustee:
		query = query.Where("trustee_org_id = ?", auth.OrganizationID)
	case models.DashboardRoleAssetCustodian:
		query = query.Where("custodian_org_id = ?", auth.OrganizationID)
	default:
		return nil, NewHTTPError(http.StatusForbidden, "unsupported stakeholder role")
	}
	if err := query.First(&request).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewHTTPError(http.StatusNotFound, "fund release request not found")
		}
		return nil, err
	}
	return &request, nil
}

func (s *FundReleaseService) Approve(ctx context.Context, auth AuthContext, id string, req models.ChallengeActionRequest) (*models.FundReleaseRequest, error) {
	if auth.DashboardRole != models.DashboardRoleTrustee {
		return nil, NewHTTPError(http.StatusForbidden, "trustee role is required")
	}
	var updated models.FundReleaseRequest
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		request, err := s.lockRequestForOrg(ctx, tx, id, "trustee_org_id", auth.OrganizationID)
		if err != nil {
			return err
		}
		if !CanTransitionFundRelease(request.Status, models.FundReleaseStatusExecutionPending) {
			return NewHTTPError(http.StatusConflict, "fund release request cannot be approved from current status")
		}
		if s.authorizations != nil {
			if err := s.authorizations.ConsumeVerifiedChallenge(ctx, tx, auth, models.ActionFundReleaseApprove, models.EntityFundReleaseRequest, request.ID, req.ChallengeID); err != nil {
				return err
			}
		}
		now := time.Now()
		before := request.Status
		request.Status = models.FundReleaseStatusExecutionPending
		request.ReviewedByMemberID = &auth.MemberID
		request.ReviewedAt = &now
		request.UpdatedAt = now
		if err := tx.Save(request).Error; err != nil {
			return err
		}
		updated = *request
		if s.notifications != nil {
			if err := s.notifyFundReleaseTransition(ctx, tx, auth.OrganizationID, request, "fund_release.approved", "Fund release approved", "A fund release request was approved by the trustee."); err != nil {
				return err
			}
		}
		if s.audit != nil {
			event := AuditEventFromAuth(auth, models.ActionFundReleaseApprove, models.EntityFundReleaseRequest, request.ID)
			event.BeforeState = models.JSONMap{"status": before}
			event.AfterState = models.JSONMap{"status": request.Status}
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

func (s *FundReleaseService) Reject(ctx context.Context, auth AuthContext, id string, req models.RejectRequest) (*models.FundReleaseRequest, error) {
	if auth.DashboardRole != models.DashboardRoleTrustee {
		return nil, NewHTTPError(http.StatusForbidden, "trustee role is required")
	}
	reason := strings.TrimSpace(req.Reason)
	if reason == "" {
		return nil, NewHTTPError(http.StatusBadRequest, "rejection reason is required")
	}
	var updated models.FundReleaseRequest
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		request, err := s.lockRequestForOrg(ctx, tx, id, "trustee_org_id", auth.OrganizationID)
		if err != nil {
			return err
		}
		if !CanTransitionFundRelease(request.Status, models.FundReleaseStatusTrusteeRejected) {
			return NewHTTPError(http.StatusConflict, "fund release request cannot be rejected from current status")
		}
		now := time.Now()
		before := request.Status
		request.Status = models.FundReleaseStatusTrusteeRejected
		request.ReviewedByMemberID = &auth.MemberID
		request.ReviewedAt = &now
		request.RejectionReason = reason
		request.UpdatedAt = now
		if err := tx.Save(request).Error; err != nil {
			return err
		}
		updated = *request
		if s.notifications != nil {
			if err := s.notifications.Create(ctx, tx, NotificationInput{
				RecipientOrgID:    request.RequesterOrgID,
				SenderOrgID:       &auth.OrganizationID,
				Type:              "fund_release.rejected",
				Title:             "Fund release rejected",
				Message:           "A fund release request was rejected by the trustee.",
				RelatedEntityType: models.EntityFundReleaseRequest,
				RelatedEntityID:   request.ID,
			}); err != nil {
				return err
			}
		}
		if s.audit != nil {
			event := AuditEventFromAuth(auth, models.ActionFundReleaseReject, models.EntityFundReleaseRequest, request.ID)
			event.BeforeState = models.JSONMap{"status": before}
			event.AfterState = models.JSONMap{"status": request.Status, "reason": reason}
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

func (s *FundReleaseService) Execute(ctx context.Context, auth AuthContext, id string, req models.ExecuteFundReleaseRequest) (*models.FundReleaseRequest, error) {
	if auth.DashboardRole != models.DashboardRoleAssetCustodian {
		return nil, NewHTTPError(http.StatusForbidden, "asset custodian role is required")
	}
	nextStatus := strings.TrimSpace(req.Status)
	if nextStatus == "" {
		nextStatus = models.FundReleaseStatusProcessing
	}
	if nextStatus != models.FundReleaseStatusProcessing && nextStatus != models.FundReleaseStatusCompleted && nextStatus != models.FundReleaseStatusFailed {
		return nil, NewHTTPError(http.StatusBadRequest, "invalid execution status")
	}
	var updated models.FundReleaseRequest
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		request, err := s.lockRequestForOrg(ctx, tx, id, "custodian_org_id", auth.OrganizationID)
		if err != nil {
			return err
		}
		if !CanTransitionFundRelease(request.Status, nextStatus) {
			return NewHTTPError(http.StatusConflict, "fund release request cannot be executed from current status")
		}
		if s.authorizations != nil {
			if err := s.authorizations.ConsumeVerifiedChallenge(ctx, tx, auth, models.ActionFundReleaseExecute, models.EntityFundReleaseRequest, request.ID, req.ChallengeID); err != nil {
				return err
			}
		}
		before := request.Status
		now := time.Now()
		request.Status = nextStatus
		request.ExecutedByMemberID = &auth.MemberID
		request.ExecutedAt = &now
		request.ExecutionReference = strings.TrimSpace(req.ExecutionReference)
		request.FailureReason = strings.TrimSpace(req.FailureReason)
		request.UpdatedAt = now
		if err := tx.Save(request).Error; err != nil {
			return err
		}
		updated = *request
		if s.notifications != nil {
			if err := s.notifyFundReleaseTransition(ctx, tx, auth.OrganizationID, request, "fund_release.executed", "Fund release execution updated", "A fund release request execution status was updated."); err != nil {
				return err
			}
		}
		if s.audit != nil {
			event := AuditEventFromAuth(auth, models.ActionFundReleaseExecute, models.EntityFundReleaseRequest, request.ID)
			event.BeforeState = models.JSONMap{"status": before}
			event.AfterState = models.JSONMap{"status": request.Status, "execution_reference": request.ExecutionReference, "failure_reason": request.FailureReason}
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

func (s *FundReleaseService) UpdateExecutionStatus(ctx context.Context, auth AuthContext, id string, req models.UpdateFundReleaseStatusRequest) (*models.FundReleaseRequest, error) {
	if auth.DashboardRole != models.DashboardRoleAssetCustodian {
		return nil, NewHTTPError(http.StatusForbidden, "asset custodian role is required")
	}
	nextStatus := strings.TrimSpace(req.Status)
	if nextStatus != models.FundReleaseStatusProcessing && nextStatus != models.FundReleaseStatusCompleted && nextStatus != models.FundReleaseStatusFailed {
		return nil, NewHTTPError(http.StatusBadRequest, "invalid execution status")
	}
	var updated models.FundReleaseRequest
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		request, err := s.lockRequestForOrg(ctx, tx, id, "custodian_org_id", auth.OrganizationID)
		if err != nil {
			return err
		}
		if request.Status == models.FundReleaseStatusExecutionPending {
			return NewHTTPError(http.StatusForbidden, "fund release must be executed with wallet step-up before status updates")
		}
		if !CanTransitionFundRelease(request.Status, nextStatus) && request.Status != nextStatus {
			return NewHTTPError(http.StatusConflict, "fund release request cannot move to requested status")
		}
		before := request.Status
		now := time.Now()
		request.Status = nextStatus
		if req.ExecutionReference != "" {
			request.ExecutionReference = strings.TrimSpace(req.ExecutionReference)
		}
		request.FailureReason = strings.TrimSpace(req.FailureReason)
		request.UpdatedAt = now
		if nextStatus == models.FundReleaseStatusCompleted || nextStatus == models.FundReleaseStatusFailed {
			request.ExecutedByMemberID = &auth.MemberID
			request.ExecutedAt = &now
		}
		if err := tx.Save(request).Error; err != nil {
			return err
		}
		updated = *request
		if s.notifications != nil {
			if err := s.notifyFundReleaseTransition(ctx, tx, auth.OrganizationID, request, "fund_release.status_updated", "Fund release status updated", "A fund release request execution status changed."); err != nil {
				return err
			}
		}
		if s.audit != nil {
			event := AuditEventFromAuth(auth, models.ActionFundReleaseStatusUpdate, models.EntityFundReleaseRequest, request.ID)
			event.BeforeState = models.JSONMap{"status": before}
			event.AfterState = models.JSONMap{"status": request.Status, "execution_reference": request.ExecutionReference, "failure_reason": request.FailureReason}
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

func (s *FundReleaseService) lockRequestForOrg(ctx context.Context, tx *gorm.DB, id string, orgColumn string, orgID string) (*models.FundReleaseRequest, error) {
	var request models.FundReleaseRequest
	if err := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ? AND "+orgColumn+" = ?", id, orgID).
		First(&request).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewHTTPError(http.StatusNotFound, "fund release request not found")
		}
		return nil, err
	}
	return &request, nil
}

func (s *FundReleaseService) notifyFundReleaseTransition(ctx context.Context, tx *gorm.DB, senderOrgID string, request *models.FundReleaseRequest, notificationType, title, message string) error {
	recipients := []string{request.RequesterOrgID, request.TrusteeOrgID, request.CustodianOrgID}
	for _, recipient := range recipients {
		if recipient == "" || recipient == senderOrgID {
			continue
		}
		if err := s.notifications.Create(ctx, tx, NotificationInput{
			RecipientOrgID:    recipient,
			SenderOrgID:       &senderOrgID,
			Type:              notificationType,
			Title:             title,
			Message:           message,
			RelatedEntityType: models.EntityFundReleaseRequest,
			RelatedEntityID:   request.ID,
		}); err != nil {
			return err
		}
	}
	return nil
}

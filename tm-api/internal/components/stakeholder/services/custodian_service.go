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
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type CustodianService struct {
	db        *gorm.DB
	assets    *AssetService
	documents *DocumentService
	audit     *AuditService
}

type AdminComplianceListFilters struct {
	CustodianOrgID string
	AssetID        string
	Status         string
}

func NewCustodianService(db *gorm.DB, assets *AssetService, documents *DocumentService, audit *AuditService) *CustodianService {
	return &CustodianService{db: db, assets: assets, documents: documents, audit: audit}
}

func (s *CustodianService) CreateCompliance(ctx context.Context, actor AuthContext, req models.CreateComplianceItemRequest) (*models.ComplianceItem, error) {
	orgID := strings.TrimSpace(req.CustodianOrgID)
	category := strings.TrimSpace(req.Category)
	requirement := strings.TrimSpace(req.Requirement)
	if orgID == "" || category == "" || requirement == "" {
		return nil, NewHTTPError(http.StatusBadRequest, "custodian_org_id, category, and requirement are required")
	}

	dueDate, err := ParseDate(req.DueDate)
	if err != nil {
		return nil, err
	}

	var organization coreModels.Organization
	if err := s.db.WithContext(ctx).First(&organization, "id = ?", orgID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewHTTPError(http.StatusNotFound, "custodian organization not found")
		}
		return nil, err
	}
	if organization.Status != coreModels.OrganizationStatusActive {
		return nil, NewHTTPError(http.StatusUnprocessableEntity, "custodian organization is not active")
	}
	if organization.StakeholderID == nil || organization.StakeholderType == nil {
		return nil, NewHTTPError(http.StatusUnprocessableEntity, "organization is not linked to an asset custodian stakeholder")
	}
	role, supported := models.DashboardRoleForStakeholderType(*organization.StakeholderType)
	if !supported || role != models.DashboardRoleAssetCustodian {
		return nil, NewHTTPError(http.StatusUnprocessableEntity, "organization is not linked to an asset custodian stakeholder")
	}
	assetID := strings.TrimSpace(req.AssetID)
	if assetID != "" {
		if s.assets == nil {
			return nil, NewHTTPError(http.StatusServiceUnavailable, "asset service is not configured")
		}
		asset, err := s.assets.GetAsset(ctx, AuthContext{
			OrganizationID: organization.ID, StakeholderID: *organization.StakeholderID, DashboardRole: models.DashboardRoleAssetCustodian,
		}, assetID)
		if err != nil {
			return nil, err
		}
		assetID = asset.ID
	}

	now := time.Now()
	item := models.ComplianceItem{
		ID:          uuid.NewString(),
		OrgID:       organization.ID,
		AssetID:     assetID,
		Category:    category,
		Requirement: requirement,
		Status:      models.ComplianceStatusPending,
		DueDate:     dueDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		duplicateQuery := tx.Model(&models.ComplianceItem{}).
			Where("org_id = ? AND LOWER(category) = ? AND LOWER(requirement) = ?", organization.ID, strings.ToLower(category), strings.ToLower(requirement))
		if assetID == "" {
			duplicateQuery = duplicateQuery.Where("(asset_id = '' OR asset_id IS NULL)")
		} else {
			duplicateQuery = duplicateQuery.Where("asset_id = ?", assetID)
		}
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
			return NewHTTPError(http.StatusConflict, "compliance item already exists for this custodian and due date")
		}
		if err := tx.Create(&item).Error; err != nil {
			return err
		}
		if s.audit != nil {
			event := AuditEventFromAuth(actor, models.ActionComplianceCreate, models.EntityComplianceItem, item.ID)
			event.AfterState = models.JSONMap{
				"org_id": organization.ID, "category": item.Category, "requirement": item.Requirement,
				"asset_id": item.AssetID, "status": item.Status, "due_date": item.DueDate,
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
	return &item, nil
}

func (s *CustodianService) CreateAccount(ctx context.Context, auth AuthContext, req models.CreateSegregatedAccountRequest) (*models.SegregatedAccount, error) {
	if auth.DashboardRole != models.DashboardRoleAssetCustodian {
		return nil, NewHTTPError(http.StatusForbidden, "asset custodian role is required")
	}
	if s.assets == nil {
		return nil, NewHTTPError(http.StatusServiceUnavailable, "asset service is not configured")
	}
	_, asset, err := s.assets.RequireAssignmentForAsset(ctx, auth, req.AssetID)
	if err != nil {
		return nil, err
	}
	accountType, status, accountName, balance, currency, err := normalizedAccountInput(req.AccountType, req.Status, req.AccountName, req.Balance, req.Currency, true)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	memberID := auth.MemberID
	account := models.SegregatedAccount{
		ID: uuid.NewString(), AssetID: asset.ID, AssetCode: asset.AssetCode, CustodianOrgID: auth.OrganizationID,
		AccountType: accountType, AccountName: accountName, Balance: balance, Currency: currency, Status: status,
		BankDetails: req.BankDetails, CreatedByMemberID: &memberID, CreatedAt: now, UpdatedAt: now,
	}
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&account).Error; err != nil {
			return err
		}
		if s.audit != nil {
			event := AuditEventFromAuth(auth, models.ActionSegregatedAccountCreate, models.EntitySegregatedAccount, account.ID)
			event.AfterState = models.JSONMap{"asset_id": account.AssetID, "account_type": account.AccountType, "currency": account.Currency, "status": account.Status}
			return s.audit.Record(ctx, tx, event)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (s *CustodianService) UpdateAccount(ctx context.Context, auth AuthContext, id string, req models.UpdateSegregatedAccountRequest) (*models.SegregatedAccount, error) {
	if auth.DashboardRole != models.DashboardRoleAssetCustodian {
		return nil, NewHTTPError(http.StatusForbidden, "asset custodian role is required")
	}
	var updated models.SegregatedAccount
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var account models.SegregatedAccount
		if err := tx.Where("id = ? AND custodian_org_id = ?", id, auth.OrganizationID).First(&account).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return NewHTTPError(http.StatusNotFound, "segregated account not found")
			}
			return err
		}
		before := models.JSONMap{"account_type": account.AccountType, "account_name": account.AccountName, "balance": account.Balance.String(), "currency": account.Currency, "status": account.Status}
		if strings.TrimSpace(req.AccountType) != "" {
			if !ContainsString(accountTypes(), strings.TrimSpace(req.AccountType)) {
				return NewHTTPError(http.StatusBadRequest, "invalid account_type")
			}
			account.AccountType = strings.TrimSpace(req.AccountType)
		}
		if strings.TrimSpace(req.AccountName) != "" {
			account.AccountName = strings.TrimSpace(req.AccountName)
		}
		if strings.TrimSpace(req.Balance) != "" {
			balance, err := decimal.NewFromString(strings.TrimSpace(req.Balance))
			if err != nil || balance.IsNegative() {
				return NewHTTPError(http.StatusBadRequest, "balance must be a non-negative decimal")
			}
			account.Balance = balance
		}
		if strings.TrimSpace(req.Currency) != "" {
			currency, err := NormalizeCurrency(req.Currency)
			if err != nil {
				return err
			}
			account.Currency = currency
		}
		if strings.TrimSpace(req.Status) != "" {
			if !ContainsString(accountStatuses(), strings.TrimSpace(req.Status)) {
				return NewHTTPError(http.StatusBadRequest, "invalid account status")
			}
			account.Status = strings.TrimSpace(req.Status)
		}
		if req.BankDetails != nil {
			account.BankDetails = req.BankDetails
		}
		account.UpdatedAt = time.Now()
		if err := tx.Save(&account).Error; err != nil {
			return err
		}
		updated = account
		if s.audit != nil {
			event := AuditEventFromAuth(auth, models.ActionSegregatedAccountUpdate, models.EntitySegregatedAccount, account.ID)
			event.BeforeState = before
			event.AfterState = models.JSONMap{"account_type": account.AccountType, "account_name": account.AccountName, "balance": account.Balance.String(), "currency": account.Currency, "status": account.Status}
			return s.audit.Record(ctx, tx, event)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &updated, nil
}

func (s *CustodianService) ListAccounts(ctx context.Context, auth AuthContext, page, limit int) ([]models.SegregatedAccount, int64, error) {
	if auth.DashboardRole != models.DashboardRoleAssetCustodian {
		return nil, 0, NewHTTPError(http.StatusForbidden, "asset custodian role is required")
	}
	query := s.db.WithContext(ctx).Model(&models.SegregatedAccount{}).Where("custodian_org_id = ?", auth.OrganizationID)
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var accounts []models.SegregatedAccount
	if err := query.Order("created_at desc").Offset((page - 1) * limit).Limit(limit).Find(&accounts).Error; err != nil {
		return nil, 0, err
	}
	return accounts, total, nil
}

func (s *CustodianService) GetAccount(ctx context.Context, auth AuthContext, id string) (*models.SegregatedAccount, error) {
	if auth.DashboardRole != models.DashboardRoleAssetCustodian {
		return nil, NewHTTPError(http.StatusForbidden, "asset custodian role is required")
	}
	var account models.SegregatedAccount
	if err := s.db.WithContext(ctx).Where("id = ? AND custodian_org_id = ?", id, auth.OrganizationID).First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewHTTPError(http.StatusNotFound, "segregated account not found")
		}
		return nil, err
	}
	if _, _, err := s.assets.RequireAssignmentForAsset(ctx, auth, account.AssetID); err != nil {
		return nil, err
	}
	return &account, nil
}

func (s *CustodianService) ListCompliance(ctx context.Context, auth AuthContext, page, limit int) ([]models.ComplianceItem, int64, error) {
	if auth.DashboardRole != models.DashboardRoleAssetCustodian {
		return nil, 0, NewHTTPError(http.StatusForbidden, "asset custodian role is required")
	}
	query := s.db.WithContext(ctx).Model(&models.ComplianceItem{}).Where("org_id = ?", auth.OrganizationID)
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []models.ComplianceItem
	if err := query.Order("due_date asc NULLS LAST, created_at desc").Offset((page - 1) * limit).Limit(limit).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	if err := s.attachComplianceDocuments(ctx, auth, items); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (s *CustodianService) ListComplianceForAdmin(ctx context.Context, page, limit int, filters AdminComplianceListFilters) ([]models.ComplianceItem, int64, error) {
	query := s.db.WithContext(ctx).Model(&models.ComplianceItem{})
	if orgID := strings.TrimSpace(filters.CustodianOrgID); orgID != "" {
		query = query.Where("org_id = ?", orgID)
	}
	if assetID := strings.TrimSpace(filters.AssetID); assetID != "" {
		query = query.Where("asset_id = ?", assetID)
	}
	if status := strings.ToLower(strings.TrimSpace(filters.Status)); status != "" {
		if !ContainsString(complianceStatuses(), status) {
			return nil, 0, NewHTTPError(http.StatusBadRequest, "invalid compliance status")
		}
		query = query.Where("status = ?", status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []models.ComplianceItem
	if err := query.Order("due_date asc NULLS LAST, created_at desc").Offset((page - 1) * limit).Limit(limit).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	if err := s.attachComplianceDocumentsWithAccess(ctx, "", items, true); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (s *CustodianService) UpdateCompliance(ctx context.Context, auth AuthContext, id string, req models.UpdateComplianceItemRequest) (*models.ComplianceItem, error) {
	if auth.DashboardRole != models.DashboardRoleAssetCustodian {
		return nil, NewHTTPError(http.StatusForbidden, "asset custodian role is required")
	}
	status := strings.TrimSpace(req.Status)
	if !ContainsString([]string{models.ComplianceStatusPending, models.ComplianceStatusComplete, models.ComplianceStatusOverdue, models.ComplianceStatusWaived}, status) {
		return nil, NewHTTPError(http.StatusBadRequest, "invalid compliance status")
	}
	var existing models.ComplianceItem
	if err := s.db.WithContext(ctx).Where("id = ? AND org_id = ?", id, auth.OrganizationID).First(&existing).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewHTTPError(http.StatusNotFound, "compliance item not found")
		}
		return nil, err
	}
	if req.DocumentIDs != nil {
		if s.documents == nil {
			return nil, NewHTTPError(http.StatusServiceUnavailable, "document service is not configured")
		}
		if len(*req.DocumentIDs) > 0 {
			if err := s.documents.ValidateForUse(ctx, auth, existing.AssetID, models.DocumentCategoryCustodyCompliance, *req.DocumentIDs); err != nil {
				return nil, err
			}
		}
	}
	var updated models.ComplianceItem
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var item models.ComplianceItem
		if err := tx.Where("id = ? AND org_id = ?", id, auth.OrganizationID).First(&item).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return NewHTTPError(http.StatusNotFound, "compliance item not found")
			}
			return err
		}
		before := item.Status
		now := time.Now()
		item.Status = status
		item.UpdatedAt = now
		if status == models.ComplianceStatusComplete {
			if item.CompletedAt == nil {
				item.CompletedAt = &now
				item.CompletedByMemberID = &auth.MemberID
			}
		} else {
			item.CompletedAt = nil
			item.CompletedByMemberID = nil
		}
		if err := tx.Save(&item).Error; err != nil {
			return err
		}
		if req.DocumentIDs != nil {
			if err := replaceComplianceDocuments(tx, item.ID, *req.DocumentIDs, now); err != nil {
				return err
			}
		}
		updated = item
		if s.audit != nil {
			event := AuditEventFromAuth(auth, models.ActionComplianceUpdate, models.EntityComplianceItem, item.ID)
			event.BeforeState = models.JSONMap{"status": before}
			event.AfterState = models.JSONMap{"status": item.Status, "document_ids": req.DocumentIDs}
			if err := s.audit.Record(ctx, tx, event); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	items := []models.ComplianceItem{updated}
	if err := s.attachComplianceDocuments(ctx, auth, items); err != nil {
		return nil, err
	}
	updated = items[0]
	return &updated, nil
}

func normalizedAccountInput(accountType, status, accountName, balanceValue, currencyValue string, creating bool) (string, string, string, decimal.Decimal, string, error) {
	accountType = strings.TrimSpace(accountType)
	if !ContainsString(accountTypes(), accountType) {
		return "", "", "", decimal.Zero, "", NewHTTPError(http.StatusBadRequest, "invalid account_type")
	}
	status = strings.TrimSpace(status)
	if status == "" && creating {
		status = models.AccountStatusActive
	}
	if !ContainsString(accountStatuses(), status) {
		return "", "", "", decimal.Zero, "", NewHTTPError(http.StatusBadRequest, "invalid account status")
	}
	accountName = strings.TrimSpace(accountName)
	if accountName == "" {
		return "", "", "", decimal.Zero, "", NewHTTPError(http.StatusBadRequest, "account_name is required")
	}
	balance := decimal.Zero
	var err error
	if strings.TrimSpace(balanceValue) != "" {
		balance, err = decimal.NewFromString(strings.TrimSpace(balanceValue))
		if err != nil || balance.IsNegative() {
			return "", "", "", decimal.Zero, "", NewHTTPError(http.StatusBadRequest, "balance must be a non-negative decimal")
		}
	}
	currency, err := NormalizeCurrency(currencyValue)
	if err != nil {
		return "", "", "", decimal.Zero, "", err
	}
	return accountType, status, accountName, balance, currency, nil
}

func accountTypes() []string {
	return []string{models.AccountTypeSaleProceeds, models.AccountTypeDevelopmentFund, models.AccountTypeRevenue, models.AccountTypeReserve}
}

func accountStatuses() []string {
	return []string{models.AccountStatusActive, models.AccountStatusInactive, models.AccountStatusClosed}
}

func complianceStatuses() []string {
	return []string{
		models.ComplianceStatusPending,
		models.ComplianceStatusComplete,
		models.ComplianceStatusOverdue,
		models.ComplianceStatusWaived,
	}
}

func replaceComplianceDocuments(tx *gorm.DB, itemID string, documentIDs []string, now time.Time) error {
	if err := tx.Where("compliance_item_id = ?", itemID).Delete(&models.ComplianceItemDocument{}).Error; err != nil {
		return err
	}
	links := make([]models.ComplianceItemDocument, 0, len(documentIDs))
	for _, documentID := range documentIDs {
		links = append(links, models.ComplianceItemDocument{ComplianceItemID: itemID, DocumentID: strings.TrimSpace(documentID), CreatedAt: now})
	}
	if len(links) > 0 {
		return tx.Create(&links).Error
	}
	return nil
}

func (s *CustodianService) attachComplianceDocuments(ctx context.Context, auth AuthContext, items []models.ComplianceItem) error {
	return s.attachComplianceDocumentsWithAccess(ctx, auth.DashboardRole, items, false)
}

func (s *CustodianService) attachComplianceDocumentsWithAccess(ctx context.Context, dashboardRole string, items []models.ComplianceItem, allowAll bool) error {
	if len(items) == 0 {
		return nil
	}
	itemIDs := make([]string, 0, len(items))
	for i := range items {
		items[i].Documents = []models.StakeholderDocument{}
		itemIDs = append(itemIDs, items[i].ID)
	}
	var links []models.ComplianceItemDocument
	if err := s.db.WithContext(ctx).Where("compliance_item_id IN ?", itemIDs).Order("created_at asc").Find(&links).Error; err != nil {
		return err
	}
	documentIDs := make([]string, 0, len(links))
	for _, link := range links {
		documentIDs = append(documentIDs, link.DocumentID)
	}
	if len(documentIDs) == 0 {
		return nil
	}
	var documents []models.StakeholderDocument
	if err := s.db.WithContext(ctx).Where("id IN ? AND status = ?", documentIDs, models.DocumentStatusActive).Find(&documents).Error; err != nil {
		return err
	}
	documentByID := make(map[string]models.StakeholderDocument, len(documents))
	for _, document := range documents {
		if allowAll || models.RoleCanAccessDocument(dashboardRole, []string(document.AccessRoles)) {
			documentByID[document.ID] = document
		}
	}
	itemIndex := make(map[string]int, len(items))
	for i := range items {
		itemIndex[items[i].ID] = i
	}
	for _, link := range links {
		if document, ok := documentByID[link.DocumentID]; ok {
			if index, exists := itemIndex[link.ComplianceItemID]; exists {
				items[index].Documents = append(items[index].Documents, document)
			}
		}
	}
	return nil
}

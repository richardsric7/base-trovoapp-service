package models

import "strings"

const (
	DashboardRoleTrustee        = "trustee"
	DashboardRoleAssetCustodian = "asset_custodian"
	DashboardRoleAssetManager   = "asset_manager"
	// New portal roles (view-only unless noted). Tier 1 = action-capable
	// (Legal/Financial Adviser confirm their A5 workstream); Tier 2 = view-only.
	DashboardRoleLegalAdviser     = "legal_adviser"
	DashboardRoleFinancialAdviser = "financial_adviser"
	DashboardRoleIssuingHouse     = "asset_issuing_house"
	DashboardRoleRatingAgency     = "rating_agency"

	RouteSegmentTrustee          = "trustee"
	RouteSegmentCustodian        = "custodian"
	RouteSegmentAssetManager     = "asset-manager"
	RouteSegmentShared           = "shared"
	RouteSegmentLegalAdviser     = "legal-adviser"
	RouteSegmentFinancialAdviser = "financial-adviser"
	RouteSegmentIssuingHouse     = "issuing-house"
	RouteSegmentRatingAgency     = "rating-agency"

	StakeholderTypeTrustee          = "trustees"
	StakeholderTypeAssetCustodian   = "approved_asset_custodian"
	StakeholderTypeAssetManager     = "asset_manager"
	StakeholderTypeLegalAdviser     = "legal_adviser"
	StakeholderTypeFinancialAdviser = "financial_adviser"
	StakeholderTypeIssuingHouse     = "asset_issuing_house"
	StakeholderTypeRatingAgency     = "rating_agency"

	// tokenized_assets FK columns used to scope "assets assigned to me".
	AssetFKColumnTrustee      = "trustee_id"
	AssetFKColumnAssetManager = "asset_manager_id"
	AssetFKColumnCustodian    = "approved_asset_custodian_id"
	AssetFKColumnLegalAdviser = "legal_adviser_id"
	AssetFKColumnFinancial    = "financial_adviser_id"
	AssetFKColumnIssuingHouse = "asset_issuing_house_id"
	AssetFKColumnRatingAgency = "rating_agency_id"

	// Structuring workstream (Tier-1 A5 confirmation).
	StructuringWorkstreamLegal     = "legal"
	StructuringWorkstreamFinancial = "financial"
	StructuringStatusPending       = "pending"
	StructuringStatusComplete      = "complete"
	ActionStructuringComplete      = "structuring.complete"
	EntityStructuringStatus        = "stakeholder_structuring_status"
)

const (
	AssignmentStatusActive   = "active"
	AssignmentStatusInactive = "inactive"

	DueDiligenceStatusDraft    = "draft"
	DueDiligenceStatusInReview = "in_review"
	DueDiligenceStatusApproved = "approved"
	DueDiligenceStatusRejected = "rejected"

	DueDiligenceItemStatusPending       = "pending"
	DueDiligenceItemStatusComplete      = "complete"
	DueDiligenceItemStatusFailed        = "failed"
	DueDiligenceItemStatusNotApplicable = "not_applicable"

	FundReleaseStatusDraft            = "draft"
	FundReleaseStatusSubmitted        = "submitted"
	FundReleaseStatusTrusteeApproved  = "trustee_approved"
	FundReleaseStatusExecutionPending = "execution_pending"
	FundReleaseStatusProcessing       = "processing"
	FundReleaseStatusCompleted        = "completed"
	FundReleaseStatusTrusteeRejected  = "trustee_rejected"
	FundReleaseStatusFailed           = "failed"

	DistributionStatusProposed   = "proposed"
	DistributionStatusAuthorized = "authorized"
	DistributionStatusRejected   = "rejected"
	DistributionStatusProcessing = "processing"
	DistributionStatusCompleted  = "completed"
	DistributionStatusFailed     = "failed"

	RevenueStatusDraft                    = "draft"
	RevenueStatusRecorded                 = "recorded"
	RevenueStatusSubmittedForDistribution = "submitted_for_distribution"

	ValuationStatusSubmitted            = "submitted"
	ValuationStatusIndependentRequested = "independent_requested"
	ValuationStatusAccepted             = "accepted"
	ValuationStatusRejected             = "rejected"

	AccountStatusActive   = "active"
	AccountStatusInactive = "inactive"
	AccountStatusClosed   = "closed"

	AccountTypeSaleProceeds    = "sale_proceeds"
	AccountTypeDevelopmentFund = "development_funds"
	AccountTypeRevenue         = "revenue"
	AccountTypeReserve         = "reserve"

	ComplianceStatusPending  = "pending"
	ComplianceStatusComplete = "complete"
	ComplianceStatusOverdue  = "overdue"
	ComplianceStatusWaived   = "waived"

	// ComplianceRequirementStatus* is the shared status enum for
	// ComplianceRequirementInstance (templated + ad-hoc), matching the
	// frontend's ComplianceRequirementStatus enum exactly. See
	// PRD-kyc-compliance-templates.md §6.5.
	ComplianceRequirementStatusDraft       = "draft"
	ComplianceRequirementStatusSubmitted   = "submitted"
	ComplianceRequirementStatusUnderReview = "under_review"
	ComplianceRequirementStatusApproved    = "approved"
	ComplianceRequirementStatusRejected    = "rejected"
	ComplianceRequirementStatusPending     = "pending"
	ComplianceRequirementStatusOverdue     = "overdue"
	ComplianceRequirementStatusWaived      = "waived"

	ComplianceInputTypeDocumentUpload = "document_upload"
	ComplianceInputTypeText           = "text"
	ComplianceInputTypeStructuredForm = "structured_form"

	// Asset compliance is derived from WalletDB's tokenization vetting fields.
	AssetComplianceStatusPending  = "pending"
	AssetComplianceStatusVerified = "verified"
	AssetComplianceStatusRejected = "rejected"

	// Custody status is a portal-friendly projection of WalletDB's canonical
	// tokenization lifecycle: primary/secondary sale assets are active and
	// status 7 is liquidated.
	AssetCustodyStatusInactive           = "inactive"
	AssetCustodyStatusActive             = "active"
	AssetCustodyStatusLiquidated         = "liquidated"
	AssetTokenizationStatusPrimarySale   = 5
	AssetTokenizationStatusSecondarySale = 6
	AssetTokenizationStatusLiquidated    = 7

	DocumentStatusActive     = "active"
	DocumentStatusSuperseded = "superseded"
	DocumentStatusDeleted    = "deleted"

	DocumentCategoryValuationReport       = "valuation_report"
	DocumentCategoryAssetReport           = "asset_report"
	DocumentCategoryFundReleaseSupporting = "fund_release_supporting"
	DocumentCategoryDueDiligence          = "due_diligence"
	DocumentCategoryCustodyCompliance     = "custody_compliance"
	DocumentCategoryComplianceRequirement = "compliance_requirement"
	DocumentCategoryLegal                 = "legal"
	DocumentCategoryFinancial             = "financial"
	DocumentCategoryStructuring           = "structuring"
	DocumentCategoryGeneral               = "general"

	AuthorizationStatusPending  = "pending"
	AuthorizationStatusVerified = "verified"
	AuthorizationStatusExpired  = "expired"
	AuthorizationStatusFailed   = "failed"
	AuthorizationStatusUsed     = "used"

	ReportStatusGenerated = "generated"
	ReportStatusSubmitted = "submitted"

	ReportTypeIncome      = "income"
	ReportTypeOperational = "operational"

	AssetOperationalStatusActive      = "active"
	AssetOperationalStatusInactive    = "inactive"
	AssetOperationalStatusPaused      = "paused"
	AssetOperationalStatusMaintenance = "maintenance"
)

func AllowedReportTypes() []string {
	return []string{ReportTypeIncome, ReportTypeOperational}
}

func AllowedAssetOperationalStatuses() []string {
	return []string{
		AssetOperationalStatusActive,
		AssetOperationalStatusInactive,
		AssetOperationalStatusPaused,
		AssetOperationalStatusMaintenance,
	}
}

const (
	ActionAssetAssignmentCreate            = "asset_assignment.create"
	ActionProfileUpdate                    = "profile.update"
	ActionNotificationPreferencesUpdate    = "profile.notification_preferences.update"
	ActionDocumentCreate                   = "document.create"
	ActionFundReleaseCreate                = "fund_release.create"
	ActionFundReleaseApprove               = "fund_release.approve"
	ActionFundReleaseReject                = "fund_release.reject"
	ActionFundReleaseExecute               = "fund_release.execute"
	ActionFundReleaseStatusUpdate          = "fund_release.status_update"
	ActionDueDiligenceItemUpdate           = "due_diligence.item_update"
	ActionDueDiligenceCategoryUpdate       = "due_diligence.category_update"
	ActionDueDiligenceApprove              = "due_diligence.approve"
	ActionDueDiligenceReject               = "due_diligence.reject"
	ActionRevenueCreate                    = "revenue.create"
	ActionDistributionSubmit               = "distribution.submit"
	ActionDistributionAuthorize            = "distribution.authorize"
	ActionDistributionReject               = "distribution.reject"
	ActionValuationCreate                  = "valuation.create"
	ActionValuationRequestIndependent      = "valuation.request_independent"
	ActionComplianceCreate                 = "compliance.create"
	ActionComplianceUpdate                 = "compliance.update"
	ActionComplianceTemplateCreate         = "compliance_template.create"
	ActionComplianceTemplateUpdate         = "compliance_template.update"
	ActionComplianceRequirementCreate      = "compliance_requirement.create"
	ActionComplianceRequirementDraft       = "compliance_requirement.draft"
	ActionComplianceRequirementReviewStart = "compliance_requirement.review_start"
	ActionComplianceRequirementReview      = "compliance_requirement.review"
	ActionComplianceRequirementSubmit      = "compliance_requirement.submit"
	ActionSegregatedAccountCreate          = "segregated_account.create"
	ActionSegregatedAccountUpdate          = "segregated_account.update"
	ActionAssetOperationUpdate             = "asset.operation_update"
	ActionReportGenerate                   = "report.generate"
	ActionReportSubmit                     = "report.submit"
)

const (
	EntityAssetAssignment               = "stakeholder_asset_assignment"
	EntityFundReleaseRequest            = "fund_release_request"
	EntityDueDiligenceChecklist         = "due_diligence_checklist"
	EntityDueDiligenceItem              = "due_diligence_item"
	EntityDueDiligenceCategory          = "due_diligence_category"
	EntityRevenueRecord                 = "revenue_record"
	EntityDistribution                  = "distribution"
	EntityAssetValuation                = "asset_valuation"
	EntitySegregatedAccount             = "segregated_account"
	EntityComplianceItem                = "compliance_item"
	EntityComplianceTemplate            = "compliance_requirement_template"
	EntityComplianceRequirementInstance = "compliance_requirement_instance"
	EntityDocument                      = "stakeholder_document"
	EntityAuthorization                 = "stakeholder_authorization_challenge"
	EntityProfile                       = "stakeholder_profile"
	EntityAssetOperation                = "stakeholder_asset_operation"
	EntityReport                        = "stakeholder_report"
)

// RoleConfig is the single source of truth describing a portal-enabled
// stakeholder role. Adding a new dashboard type is a one-line append here;
// the gate maps and per-role helpers are all derived from this slice.
type RoleConfig struct {
	Role              string // dashboard_role
	StakeholderType   string // canonical stakeholder_type (as stored on the org)
	RouteSegment      string // URL segment under /api/v1/stakeholder/
	AssetFKColumn     string // tokenized_assets FK column used for "assigned to me"
	HasAssignmentCols bool   // true only for roles with columns on stakeholder_asset_assignments
	Tier              int    // 1 = action-capable, 2 = view-only
}

// roleConfigs registers every portal-enabled role. The first three keep their
// richer explicit-assignment wiring (HasAssignmentCols); the rest are FK-only.
var roleConfigs = []RoleConfig{
	{DashboardRoleTrustee, StakeholderTypeTrustee, RouteSegmentTrustee, AssetFKColumnTrustee, true, 1},
	{DashboardRoleAssetCustodian, StakeholderTypeAssetCustodian, RouteSegmentCustodian, AssetFKColumnCustodian, true, 1},
	{DashboardRoleAssetManager, StakeholderTypeAssetManager, RouteSegmentAssetManager, AssetFKColumnAssetManager, true, 1},
	{DashboardRoleLegalAdviser, StakeholderTypeLegalAdviser, RouteSegmentLegalAdviser, AssetFKColumnLegalAdviser, false, 1},
	{DashboardRoleFinancialAdviser, StakeholderTypeFinancialAdviser, RouteSegmentFinancialAdviser, AssetFKColumnFinancial, false, 1},
	{DashboardRoleIssuingHouse, StakeholderTypeIssuingHouse, RouteSegmentIssuingHouse, AssetFKColumnIssuingHouse, false, 2},
	{DashboardRoleRatingAgency, StakeholderTypeRatingAgency, RouteSegmentRatingAgency, AssetFKColumnRatingAgency, false, 2},
}

var (
	supportedRoles                 = map[string]struct{}{}
	stakeholderTypeToDashboardRole = map[string]string{}
	roleConfigByRole               = map[string]RoleConfig{}
)

func init() {
	for _, rc := range roleConfigs {
		supportedRoles[rc.Role] = struct{}{}
		stakeholderTypeToDashboardRole[rc.StakeholderType] = rc.Role
		roleConfigByRole[rc.Role] = rc
	}
}

func Normalize(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func DashboardRoleForStakeholderType(stakeholderType string) (string, bool) {
	role, ok := stakeholderTypeToDashboardRole[Normalize(stakeholderType)]
	return role, ok
}

func IsSupportedDashboardRole(role string) bool {
	_, ok := supportedRoles[Normalize(role)]
	return ok
}

func SupportedDashboardRoles() []string {
	roles := make([]string, 0, len(roleConfigs))
	for _, rc := range roleConfigs {
		roles = append(roles, rc.Role)
	}
	return roles
}

// RoleConfigForRole returns the registry entry for a dashboard role.
func RoleConfigForRole(role string) (RoleConfig, bool) {
	rc, ok := roleConfigByRole[Normalize(role)]
	return rc, ok
}

// AssetFKColumnForRole returns the tokenized_assets FK column that scopes
// "assets assigned to me" for the given role (empty if unknown).
func AssetFKColumnForRole(role string) string {
	if rc, ok := RoleConfigForRole(role); ok {
		return rc.AssetFKColumn
	}
	return ""
}

// RoleHasAssignmentColumns reports whether the role is backed by explicit
// columns on stakeholder_asset_assignments (vs FK-only wallet-fallback scoping).
func RoleHasAssignmentColumns(role string) bool {
	if rc, ok := RoleConfigForRole(role); ok {
		return rc.HasAssignmentCols
	}
	return false
}

func RoleCanAccessDocument(role string, accessRoles []string) bool {
	role = Normalize(role)
	if len(accessRoles) == 0 {
		return true
	}
	for _, accessRole := range accessRoles {
		if Normalize(accessRole) == role {
			return true
		}
	}
	return false
}

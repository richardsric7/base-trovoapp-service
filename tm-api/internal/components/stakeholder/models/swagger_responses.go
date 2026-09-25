package models

// StakeholderErrorResponse is the error envelope for stakeholder-portal
// @Failure responses (single "error" string). Named distinctly to avoid
// colliding with the core models.ErrorResponse in swag's global type table.
type StakeholderErrorResponse struct {
	Error string `json:"error" example:"stakeholder role is not allowed for this route"`
}

// This file defines Swagger response-wrapper types only. Handlers return their
// payload under a top-level "data" key (gin.H{"data": <model>}); these structs
// model that envelope so the generated OpenAPI schema reflects the real shape.
// List endpoints return a PaginatedResponse whose Records field is typed as
// interface{} at runtime — the *TypedPage structs below document the concrete
// element type for each list endpoint.

// --- Single-object envelopes: { "data": <model> } ---

type AssetAssignmentDataResponse struct {
	Data StakeholderAssetAssignment `json:"data"`
}

type ProfileDataResponse struct {
	Data StakeholderProfileResponse `json:"data"`
}

type NotificationPreferenceDataResponse struct {
	Data StakeholderNotificationPreference `json:"data"`
}

type AssetDataResponse struct {
	Data TokenizedAssetDetailResponse `json:"data"`
}

type AuthorizationChallengeDataResponse struct {
	Data AuthorizationChallengeResponse `json:"data"`
}

type DocumentDataResponse struct {
	Data StakeholderDocument `json:"data"`
}

type NotificationDataResponse struct {
	Data StakeholderNotification `json:"data"`
}

type DashboardDataResponse struct {
	Data DashboardResponse `json:"data"`
}

type DueDiligenceChecklistDataResponse struct {
	Data DueDiligenceChecklist `json:"data"`
}

type DueDiligenceItemDataResponse struct {
	Data DueDiligenceItem `json:"data"`
}

type DueDiligenceCategoryDataResponse struct {
	Data DueDiligenceCategoryResponse `json:"data"`
}

type FundReleaseDataResponse struct {
	Data FundReleaseRequest `json:"data"`
}

type FundReleaseDetailDataResponse struct {
	Data FundReleaseDetailResponse `json:"data"`
}

type AssetManagerFundSummaryDataResponse struct {
	Data AssetManagerFundSummary `json:"data"`
}

type TrusteeFundManagementSummaryDataResponse struct {
	Data TrusteeFundManagementSummary `json:"data"`
}

type DistributionDataResponse struct {
	Data Distribution `json:"data"`
}

type DistributionDetailDataResponse struct {
	Data DistributionDetailResponse `json:"data"`
}

type RevenueDataResponse struct {
	Data RevenueRecord `json:"data"`
}

type ValuationDataResponse struct {
	Data AssetValuation `json:"data"`
}

type SegregatedAccountDataResponse struct {
	Data SegregatedAccount `json:"data"`
}

type ComplianceItemDataResponse struct {
	Data ComplianceItem `json:"data"`
}

type ComplianceRequirementTemplateDataResponse struct {
	Data ComplianceRequirementTemplate `json:"data"`
}

type ComplianceRequirementInstanceDataResponse struct {
	Data ComplianceRequirementInstance `json:"data"`
}

type AssetOperationDataResponse struct {
	Data StakeholderAssetOperation `json:"data"`
}

type ReportDataResponse struct {
	Data StakeholderReport `json:"data"`
}

type ReportTypeListDataResponse struct {
	Data []ReportTypeResponse `json:"data"`
}

type StructuringStatusDataResponse struct {
	Data StakeholderStructuringStatus `json:"data"`
}

// --- Paginated envelopes: { "data": { "records": [<model>...], "meta": {...} } } ---
// Each *Page struct redeclares Records with a concrete element type so swag emits
// the item schema (PaginatedResponse.Records is interface{} at runtime).

type AssetPage struct {
	Records []TokenizedAssetResponse `json:"records"`
	Meta    PaginationMeta           `json:"meta"`
	Summary *AssetStatistics         `json:"summary,omitempty"`
}

type AssetPageDataResponse struct {
	Data AssetPage `json:"data"`
}

type DocumentPage struct {
	Records []StakeholderDocument `json:"records"`
	Meta    PaginationMeta        `json:"meta"`
}

type DocumentPageDataResponse struct {
	Data DocumentPage `json:"data"`
}

type ComplianceRequirementTemplatePage struct {
	Records []ComplianceRequirementTemplate `json:"records"`
	Meta    PaginationMeta                  `json:"meta"`
}

type ComplianceRequirementTemplatePageDataResponse struct {
	Data ComplianceRequirementTemplatePage `json:"data"`
}

type ComplianceRequirementInstancePage struct {
	Records []ComplianceRequirementInstance `json:"records"`
	Meta    PaginationMeta                  `json:"meta"`
}

type ComplianceRequirementInstancePageDataResponse struct {
	Data ComplianceRequirementInstancePage `json:"data"`
}

type NotificationPage struct {
	Records []StakeholderNotification `json:"records"`
	Meta    PaginationMeta            `json:"meta"`
}

type NotificationPageDataResponse struct {
	Data NotificationPage `json:"data"`
}

type AuditLogPage struct {
	Records []StakeholderAuditLog `json:"records"`
	Meta    PaginationMeta        `json:"meta"`
}

type AuditLogPageDataResponse struct {
	Data AuditLogPage `json:"data"`
}

type FundReleasePage struct {
	Records []FundReleaseListItemResponse `json:"records"`
	Meta    PaginationMeta                `json:"meta"`
}

type FundReleasePageDataResponse struct {
	Data FundReleasePage `json:"data"`
}

type DistributionPage struct {
	Records []DistributionListItemResponse `json:"records"`
	Meta    PaginationMeta                 `json:"meta"`
}

type DistributionPageDataResponse struct {
	Data DistributionPage `json:"data"`
}

type RevenuePage struct {
	Records []RevenueRecord `json:"records"`
	Meta    PaginationMeta  `json:"meta"`
}

type RevenuePageDataResponse struct {
	Data RevenuePage `json:"data"`
}

type ValuationPage struct {
	Records []AssetValuation `json:"records"`
	Meta    PaginationMeta   `json:"meta"`
}

type ValuationPageDataResponse struct {
	Data ValuationPage `json:"data"`
}

type SegregatedAccountPage struct {
	Records []SegregatedAccount `json:"records"`
	Meta    PaginationMeta      `json:"meta"`
}

type SegregatedAccountPageDataResponse struct {
	Data SegregatedAccountPage `json:"data"`
}

type ComplianceItemPage struct {
	Records []ComplianceItem `json:"records"`
	Meta    PaginationMeta   `json:"meta"`
}

type ComplianceItemPageDataResponse struct {
	Data ComplianceItemPage `json:"data"`
}

type ReportPage struct {
	Records []StakeholderReport `json:"records"`
	Meta    PaginationMeta      `json:"meta"`
}

type ReportPageDataResponse struct {
	Data ReportPage `json:"data"`
}

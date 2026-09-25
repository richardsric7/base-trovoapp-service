package models

type PaginationRequest struct {
	Page  int
	Limit int
}

// ConfirmStructuringRequest is submitted by a Legal or Financial Adviser to
// confirm their A5 workstream is complete for an asset. An optional document
// can be attached as evidence (either a pre-uploaded document_id, or a
// title+file_url to create-and-link one).
type ConfirmStructuringRequest struct {
	DocumentID *string `json:"document_id"`
	Title      string  `json:"title"`
	FileURL    string  `json:"file_url"`
	Notes      string  `json:"notes"`
}

type AssetListFilters struct {
	Page      int
	Limit     int
	Status    string
	AssetType string
	Sector    string
	Search    string
	DateFrom  string
	DateTo    string
}

type CreateAssetAssignmentRequest struct {
	AssetID           string  `json:"asset_id" binding:"required"`
	AssetCode         string  `json:"asset_code"`
	TrusteeOrgID      string  `json:"trustee_org_id" binding:"required"`
	AssetManagerOrgID *string `json:"asset_manager_org_id"`
	CustodianOrgID    *string `json:"custodian_org_id"`
	Status            string  `json:"status"`
}

type CreateAuthorizationRequest struct {
	Action     string `json:"action" binding:"required"`
	EntityType string `json:"entity_type" binding:"required"`
	EntityID   string `json:"entity_id" binding:"required"`
}

type VerifyAuthorizationRequest struct {
	AuthID string `json:"auth_id"`
}

type CreateFundReleaseRequest struct {
	AssetID                string   `json:"asset_id" binding:"required"`
	Amount                 string   `json:"amount" binding:"required"`
	Currency               string   `json:"currency" binding:"required"`
	Purpose                string   `json:"purpose" binding:"required"`
	SupportingDocumentIDs  []string `json:"supporting_document_ids" binding:"required,min=1,dive,required"`
	ReceivingBank          string   `json:"receiving_bank"`
	ReceivingAccountName   string   `json:"receiving_account_name"`
	ReceivingAccountNumber string   `json:"receiving_account_number"`
}

type ChallengeActionRequest struct {
	ChallengeID string `json:"challenge_id" binding:"required"`
}

type RejectRequest struct {
	Reason string `json:"reason" binding:"required"`
}

type ExecuteFundReleaseRequest struct {
	ChallengeID        string `json:"challenge_id" binding:"required"`
	ExecutionReference string `json:"execution_reference"`
	Status             string `json:"status"`
	FailureReason      string `json:"failure_reason"`
}

type UpdateFundReleaseStatusRequest struct {
	Status             string `json:"status" binding:"required"`
	ExecutionReference string `json:"execution_reference"`
	FailureReason      string `json:"failure_reason"`
}

type UpdateDueDiligenceItemRequest struct {
	Status      string    `json:"status" binding:"required"`
	Notes       string    `json:"notes"`
	DocumentIDs *[]string `json:"document_ids,omitempty"`
}

type UpdateDueDiligenceCategoryRequest struct {
	Category string `json:"category" binding:"required"`
	Status   string `json:"status" binding:"required"`
	Notes    string `json:"notes"`
}

type CreateRevenueRequest struct {
	AssetID     string `json:"asset_id" binding:"required"`
	PeriodStart string `json:"period_start"`
	PeriodEnd   string `json:"period_end"`
	Source      string `json:"source" binding:"required"`
	Amount      string `json:"amount" binding:"required"`
	Currency    string `json:"currency" binding:"required"`
	CollectedAt string `json:"collected_at"`
}

type SubmitDistributionRequest struct {
	AssetID          string   `json:"asset_id" binding:"required"`
	RevenueRecordIDs []string `json:"revenue_record_ids"`
	Amount           string   `json:"amount" binding:"required"`
	Currency         string   `json:"currency" binding:"required"`
	Source           string   `json:"source" binding:"required"`
	ScheduledDate    string   `json:"scheduled_date"`
}

type CreateValuationRequest struct {
	AssetID          string `json:"asset_id" binding:"required"`
	Valuation        string `json:"valuation" binding:"required"`
	Currency         string `json:"currency" binding:"required"`
	Methodology      string `json:"methodology" binding:"required"`
	ValuationDate    string `json:"valuation_date" binding:"required"`
	ReportDocumentID string `json:"report_document_id" binding:"required"`
}

type UploadDocumentRequest struct {
	AssetID     string
	Category    string
	Title       string
	AccessRoles []string
	Filename    string
	Content     []byte
}

type UpdateComplianceItemRequest struct {
	Status      string    `json:"status" binding:"required"`
	DocumentIDs *[]string `json:"document_ids,omitempty"`
}

type CreateComplianceItemRequest struct {
	CustodianOrgID string `json:"custodian_org_id" binding:"required"`
	Category       string `json:"category" binding:"required"`
	Requirement    string `json:"requirement" binding:"required"`
	DueDate        string `json:"due_date"`
	AssetID        string `json:"asset_id"`
}

type ComplianceTemplateItemInput struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name" binding:"required"`
	Category    string `json:"category" binding:"required"`
	Description string `json:"description"`
	InputType   string `json:"input_type" binding:"required"`
	Required    *bool  `json:"required"`
}

type CreateComplianceTemplateRequest struct {
	OrgType string                        `json:"org_type" binding:"required"`
	Level   int                           `json:"level" binding:"required"`
	Items   []ComplianceTemplateItemInput `json:"items" binding:"required,min=1,dive"`
}

type UpdateComplianceTemplateRequest struct {
	OrgType *string                       `json:"org_type"`
	Level   *int                          `json:"level"`
	Items   []ComplianceTemplateItemInput `json:"items" binding:"required,min=1,dive"`
}

type CreateComplianceRequirementRequest struct {
	OrgID       string `json:"org_id" binding:"required"`
	Category    string `json:"category" binding:"required"`
	Requirement string `json:"requirement" binding:"required"`
	Description string `json:"description"`
	InputType   string `json:"input_type" binding:"required"`
	Required    *bool  `json:"required"`
	DueDate     string `json:"due_date"`
}

type ReviewComplianceRequirementRequest struct {
	Decision        string `json:"decision" binding:"required"`
	RejectionReason string `json:"rejection_reason"`
}

type SubmitComplianceRequirementRequest struct {
	SubmissionData map[string]interface{} `json:"submission_data" binding:"required"`
}

type UpdateOrgComplianceRequirementRequest struct {
	Status string `json:"status" binding:"required"`
}

type CreateSegregatedAccountRequest struct {
	AssetID     string  `json:"asset_id" binding:"required"`
	AccountType string  `json:"account_type" binding:"required"`
	AccountName string  `json:"account_name" binding:"required"`
	Balance     string  `json:"balance"`
	Currency    string  `json:"currency" binding:"required"`
	Status      string  `json:"status"`
	BankDetails JSONMap `json:"bank_details"`
}

type UpdateSegregatedAccountRequest struct {
	AccountType string  `json:"account_type"`
	AccountName string  `json:"account_name"`
	Balance     string  `json:"balance"`
	Currency    string  `json:"currency"`
	Status      string  `json:"status"`
	BankDetails JSONMap `json:"bank_details"`
}

type CreateDocumentRequest struct {
	AssetID     string   `json:"asset_id"`
	AssetCode   string   `json:"asset_code"`
	Category    string   `json:"category" binding:"required"`
	Title       string   `json:"title" binding:"required"`
	FileURL     string   `json:"file_url" binding:"required"`
	MimeType    string   `json:"mime_type"`
	Version     int      `json:"version"`
	AccessRoles []string `json:"access_roles"`
}

type UpdateProfileRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type UpdateNotificationPreferencesRequest struct {
	EmailEnabled *bool `json:"email_enabled"`
	InAppEnabled *bool `json:"in_app_enabled"`
}

type UpdateAssetOperationRequest struct {
	OperationalStatus string                 `json:"operational_status" binding:"required" enums:"active,inactive,paused,maintenance"`
	Notes             string                 `json:"notes"`
	Metadata          AssetOperationMetadata `json:"metadata"`
}

// AssetOperationMetadata describes the source and effective time of a
// dashboard-local operational update. It is deliberately separate from the
// canonical WalletDB tokenization lifecycle.
type AssetOperationMetadata struct {
	EffectiveAt string `json:"effective_at,omitempty" example:"2026-08-31T10:30:00Z"`
	Source      string `json:"source,omitempty" example:"manual"`
}

type CreateReportRequest struct {
	AssetID    string `json:"asset_id" binding:"required"`
	AssetCode  string `json:"asset_code"`
	ReportType string `json:"report_type" binding:"required" enums:"income,operational"`
	Title      string `json:"title" binding:"required"`
	DocumentID string `json:"document_id"`
	// FileURL is retained temporarily for existing clients. New clients must
	// upload an asset_report document and send document_id.
	FileURL string `json:"file_url,omitempty"`
}

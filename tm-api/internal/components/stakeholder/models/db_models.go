package models

import (
	coreModels "admin-panel-dashboard/internal/models"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type JSONMap map[string]interface{}

func (m JSONMap) Value() (driver.Value, error) {
	if m == nil {
		return []byte("{}"), nil
	}
	return json.Marshal(m)
}

func (m *JSONMap) Scan(value interface{}) error {
	if m == nil {
		return nil
	}
	if value == nil {
		*m = JSONMap{}
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("unsupported JSONMap scan type %T", value)
	}
	if len(bytes) == 0 {
		*m = JSONMap{}
		return nil
	}
	return json.Unmarshal(bytes, m)
}

type JSONStringArray []string

func (a JSONStringArray) Value() (driver.Value, error) {
	if a == nil {
		return []byte("[]"), nil
	}
	return json.Marshal(a)
}

func (a *JSONStringArray) Scan(value interface{}) error {
	if a == nil {
		return nil
	}
	if value == nil {
		*a = JSONStringArray{}
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("unsupported JSONStringArray scan type %T", value)
	}
	if len(bytes) == 0 {
		*a = JSONStringArray{}
		return nil
	}
	return json.Unmarshal(bytes, a)
}

type StakeholderAssetAssignment struct {
	ID                        string    `gorm:"primaryKey" json:"id"`
	AssetID                   string    `gorm:"not null;uniqueIndex" json:"asset_id"`
	AssetCode                 string    `gorm:"index" json:"asset_code"`
	TrusteeOrgID              *string   `gorm:"index" json:"trustee_org_id,omitempty"`
	TrusteeStakeholderID      *uint64   `gorm:"index" json:"trustee_stakeholder_id,omitempty"`
	AssetManagerOrgID         *string   `gorm:"index" json:"asset_manager_org_id,omitempty"`
	AssetManagerStakeholderID *uint64   `gorm:"index" json:"asset_manager_stakeholder_id,omitempty"`
	CustodianOrgID            *string   `gorm:"index" json:"custodian_org_id,omitempty"`
	CustodianStakeholderID    *uint64   `gorm:"index" json:"custodian_stakeholder_id,omitempty"`
	Status                    string    `gorm:"not null;default:'active';index" json:"status"`
	AssignedBy                string    `json:"assigned_by"`
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
}

func (StakeholderAssetAssignment) TableName() string {
	return "stakeholder_asset_assignments"
}

type DueDiligenceChecklist struct {
	ID                      string                        `gorm:"primaryKey" json:"id"`
	AssetID                 string                        `gorm:"not null;index" json:"asset_id"`
	AssetCode               string                        `gorm:"index" json:"asset_code"`
	TrusteeOrgID            string                        `gorm:"not null;index" json:"trustee_org_id"`
	Status                  string                        `gorm:"not null;default:'in_review';index" json:"status"`
	ApprovedByMemberID      *string                       `gorm:"index" json:"approved_by_member_id,omitempty"`
	ApprovedAt              *time.Time                    `json:"approved_at,omitempty"`
	RejectedByMemberID      *string                       `gorm:"index" json:"rejected_by_member_id,omitempty"`
	RejectedAt              *time.Time                    `json:"rejected_at,omitempty"`
	RejectionReason         string                        `json:"rejection_reason,omitempty"`
	CreatedAt               time.Time                     `json:"created_at"`
	UpdatedAt               time.Time                     `json:"updated_at"`
	Items                   []DueDiligenceItem            `gorm:"foreignKey:ChecklistID" json:"items,omitempty"`
	AvailableAssetDocuments []WalletAssetDocumentResponse `gorm:"-" json:"available_asset_documents"`
}

func (DueDiligenceChecklist) TableName() string {
	return "due_diligence_checklists"
}

type DueDiligenceItem struct {
	ID                 string                `gorm:"primaryKey" json:"id"`
	ChecklistID        string                `gorm:"not null;index" json:"checklist_id"`
	Category           string                `gorm:"not null" json:"category"`
	Item               string                `gorm:"not null" json:"item"`
	Status             string                `gorm:"not null;default:'pending';index" json:"status"`
	Notes              string                `json:"notes,omitempty"`
	VerifiedByMemberID *string               `gorm:"index" json:"verified_by_member_id,omitempty"`
	VerifiedAt         *time.Time            `json:"verified_at,omitempty"`
	CreatedAt          time.Time             `json:"created_at"`
	UpdatedAt          time.Time             `json:"updated_at"`
	Documents          []StakeholderDocument `gorm:"-" json:"documents"`
}

func (DueDiligenceItem) TableName() string {
	return "due_diligence_items"
}

type FundReleaseRequest struct {
	ID                     string          `gorm:"primaryKey" json:"id"`
	AssetID                string          `gorm:"not null;index" json:"asset_id"`
	AssetCode              string          `gorm:"index" json:"asset_code"`
	RequesterOrgID         string          `gorm:"not null;index" json:"requester_org_id"`
	RequesterMemberID      string          `gorm:"not null;index" json:"requester_member_id"`
	TrusteeOrgID           string          `gorm:"not null;index" json:"trustee_org_id"`
	CustodianOrgID         string          `gorm:"index" json:"custodian_org_id"`
	Amount                 decimal.Decimal `gorm:"type:numeric(30,8);not null" json:"amount"`
	Currency               string          `gorm:"size:12;not null" json:"currency"`
	Purpose                string          `gorm:"not null" json:"purpose"`
	Status                 string          `gorm:"not null;default:'submitted';index" json:"status"`
	SupportingDocumentIDs  JSONStringArray `gorm:"type:jsonb" json:"supporting_document_ids,omitempty"`
	ReceivingBank          string          `gorm:"not null;default:''" json:"-"`
	ReceivingAccountName   string          `gorm:"not null;default:''" json:"-"`
	ReceivingAccountNumber string          `gorm:"not null;default:''" json:"-"`
	ReviewedByMemberID     *string         `gorm:"index" json:"reviewed_by_member_id,omitempty"`
	ReviewedAt             *time.Time      `json:"reviewed_at,omitempty"`
	RejectionReason        string          `json:"rejection_reason,omitempty"`
	ExecutedByMemberID     *string         `gorm:"index" json:"executed_by_member_id,omitempty"`
	ExecutedAt             *time.Time      `json:"executed_at,omitempty"`
	ExecutionReference     string          `json:"execution_reference,omitempty"`
	FailureReason          string          `json:"failure_reason,omitempty"`
	CreatedAt              time.Time       `json:"created_at"`
	UpdatedAt              time.Time       `json:"updated_at"`
}

func (FundReleaseRequest) TableName() string {
	return "fund_release_requests"
}

type Distribution struct {
	ID                   string          `gorm:"primaryKey" json:"id"`
	AssetID              string          `gorm:"not null;index" json:"asset_id"`
	AssetCode            string          `gorm:"index" json:"asset_code"`
	ProposedByOrgID      string          `gorm:"not null;index" json:"proposed_by_org_id"`
	ProposedByMemberID   string          `gorm:"not null;index" json:"proposed_by_member_id"`
	TrusteeOrgID         string          `gorm:"not null;index" json:"trustee_org_id"`
	Amount               decimal.Decimal `gorm:"type:numeric(30,8);not null" json:"amount"`
	Currency             string          `gorm:"size:12;not null" json:"currency"`
	Source               string          `gorm:"not null" json:"source"`
	ScheduledDate        *time.Time      `json:"scheduled_date,omitempty"`
	Status               string          `gorm:"not null;default:'proposed';index" json:"status"`
	AuthorizedByMemberID *string         `gorm:"index" json:"authorized_by_member_id,omitempty"`
	AuthorizedAt         *time.Time      `json:"authorized_at,omitempty"`
	RejectionReason      string          `json:"rejection_reason,omitempty"`
	CreatedAt            time.Time       `json:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at"`
}

func (Distribution) TableName() string {
	return "distributions"
}

type RevenueRecord struct {
	ID                string          `gorm:"primaryKey" json:"id"`
	AssetID           string          `gorm:"not null;index" json:"asset_id"`
	AssetCode         string          `gorm:"index" json:"asset_code"`
	ManagerOrgID      string          `gorm:"not null;index" json:"manager_org_id"`
	PeriodStart       *time.Time      `json:"period_start,omitempty"`
	PeriodEnd         *time.Time      `json:"period_end,omitempty"`
	Source            string          `gorm:"not null" json:"source"`
	Amount            decimal.Decimal `gorm:"type:numeric(30,8);not null" json:"amount"`
	Currency          string          `gorm:"size:12;not null" json:"currency"`
	Status            string          `gorm:"not null;default:'recorded';index" json:"status"`
	CollectedAt       *time.Time      `json:"collected_at,omitempty"`
	CreatedByMemberID string          `gorm:"not null;index" json:"created_by_member_id"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

func (RevenueRecord) TableName() string {
	return "revenue_records"
}

type AssetValuation struct {
	ID                string          `gorm:"primaryKey" json:"id"`
	AssetID           string          `gorm:"not null;index" json:"asset_id"`
	AssetCode         string          `gorm:"index" json:"asset_code"`
	ManagerOrgID      string          `gorm:"not null;index" json:"manager_org_id"`
	Valuation         decimal.Decimal `gorm:"type:numeric(30,8);not null" json:"valuation"`
	Currency          string          `gorm:"size:12;not null" json:"currency"`
	Methodology       string          `gorm:"not null" json:"methodology"`
	ValuationDate     time.Time       `gorm:"not null;index" json:"valuation_date"`
	ReportDocumentID  *string         `gorm:"index" json:"report_document_id,omitempty"`
	Status            string          `gorm:"not null;default:'submitted';index" json:"status"`
	CreatedByMemberID string          `gorm:"not null;index" json:"created_by_member_id"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

func (AssetValuation) TableName() string {
	return "asset_valuations"
}

type SegregatedAccount struct {
	ID                string          `gorm:"primaryKey" json:"id"`
	AssetID           string          `gorm:"not null;index" json:"asset_id"`
	AssetCode         string          `gorm:"index" json:"asset_code"`
	CustodianOrgID    string          `gorm:"not null;index" json:"custodian_org_id"`
	AccountType       string          `gorm:"not null;index" json:"account_type"`
	AccountName       string          `gorm:"not null" json:"account_name"`
	Balance           decimal.Decimal `gorm:"type:numeric(30,8);not null" json:"balance"`
	Currency          string          `gorm:"size:12;not null" json:"currency"`
	Status            string          `gorm:"not null;default:'active';index" json:"status"`
	BankDetails       JSONMap         `gorm:"type:jsonb" json:"bank_details,omitempty"`
	CreatedByMemberID *string         `gorm:"index" json:"created_by_member_id,omitempty"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

func (SegregatedAccount) TableName() string {
	return "segregated_accounts"
}

type ComplianceItem struct {
	ID                  string                `gorm:"primaryKey" json:"id"`
	OrgID               string                `gorm:"not null;index" json:"org_id"`
	AssetID             string                `gorm:"index" json:"asset_id,omitempty"`
	Category            string                `gorm:"not null" json:"category"`
	Requirement         string                `gorm:"not null" json:"requirement"`
	Status              string                `gorm:"not null;default:'pending';index" json:"status"`
	DueDate             *time.Time            `gorm:"index" json:"due_date,omitempty"`
	CompletedAt         *time.Time            `json:"completed_at,omitempty"`
	CompletedByMemberID *string               `gorm:"index" json:"completed_by_member_id,omitempty"`
	CreatedAt           time.Time             `json:"created_at"`
	UpdatedAt           time.Time             `json:"updated_at"`
	Documents           []StakeholderDocument `gorm:"-" json:"documents"`
}

func (ComplianceItem) TableName() string {
	return "compliance_items"
}

type ComplianceItemDocument struct {
	ComplianceItemID string               `gorm:"primaryKey;not null" json:"compliance_item_id"`
	ComplianceItem   *ComplianceItem      `gorm:"foreignKey:ComplianceItemID;references:ID;constraint:compliance_item_documents_compliance_item_id_fkey,OnDelete:CASCADE" json:"-"`
	DocumentID       string               `gorm:"primaryKey;not null;index:idx_compliance_item_documents_document_id" json:"document_id"`
	Document         *StakeholderDocument `gorm:"foreignKey:DocumentID;references:ID;constraint:compliance_item_documents_document_id_fkey,OnDelete:CASCADE" json:"-"`
	CreatedAt        time.Time            `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (ComplianceItemDocument) TableName() string { return "compliance_item_documents" }

type DueDiligenceItemDocument struct {
	DueDiligenceItemID string               `gorm:"primaryKey;not null" json:"due_diligence_item_id"`
	DueDiligenceItem   *DueDiligenceItem    `gorm:"foreignKey:DueDiligenceItemID;references:ID;constraint:due_diligence_item_documents_due_diligence_item_id_fkey,OnDelete:CASCADE" json:"-"`
	DocumentID         string               `gorm:"primaryKey;not null;index:idx_due_diligence_item_documents_document_id" json:"document_id"`
	Document           *StakeholderDocument `gorm:"foreignKey:DocumentID;references:ID;constraint:due_diligence_item_documents_document_id_fkey,OnDelete:CASCADE" json:"-"`
	CreatedAt          time.Time            `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (DueDiligenceItemDocument) TableName() string { return "due_diligence_item_documents" }

type StakeholderDocument struct {
	ID                 string          `gorm:"primaryKey" json:"id"`
	AssetID            string          `gorm:"index" json:"asset_id,omitempty"`
	AssetCode          string          `gorm:"index" json:"asset_code,omitempty"`
	UploadedByOrgID    string          `gorm:"not null;index" json:"uploaded_by_org_id"`
	UploadedByMemberID string          `gorm:"not null;index" json:"uploaded_by_member_id"`
	Category           string          `gorm:"not null;index" json:"category"`
	Title              string          `gorm:"not null" json:"title"`
	FileURL            string          `gorm:"not null" json:"file_url"`
	StorageObjectID    string          `gorm:"index" json:"-"`
	OriginalFilename   string          `json:"original_filename,omitempty"`
	MimeType           string          `json:"mime_type,omitempty"`
	SizeBytes          int64           `json:"size_bytes,omitempty"`
	SHA256             string          `json:"sha256,omitempty"`
	Version            int             `gorm:"not null;default:1" json:"version"`
	AccessRoles        JSONStringArray `gorm:"type:jsonb" json:"access_roles,omitempty"`
	Status             string          `gorm:"not null;default:'active';index" json:"status"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}

func (StakeholderDocument) TableName() string {
	return "stakeholder_documents"
}

type StakeholderNotification struct {
	ID                string     `gorm:"primaryKey" json:"id"`
	RecipientOrgID    string     `gorm:"not null;index" json:"recipient_org_id"`
	SenderOrgID       *string    `gorm:"index" json:"sender_org_id,omitempty"`
	Type              string     `gorm:"not null;index" json:"type"`
	Title             string     `gorm:"not null" json:"title"`
	Message           string     `gorm:"not null" json:"message"`
	RelatedEntityType string     `gorm:"index" json:"related_entity_type,omitempty"`
	RelatedEntityID   string     `gorm:"index" json:"related_entity_id,omitempty"`
	ReadAt            *time.Time `gorm:"index" json:"read_at,omitempty"`
	CreatedAt         time.Time  `gorm:"index" json:"created_at"`
}

func (StakeholderNotification) TableName() string {
	return "stakeholder_notifications"
}

type StakeholderAuditLog struct {
	ID                string    `gorm:"primaryKey" json:"id"`
	RequestID         string    `gorm:"index" json:"request_id,omitempty"`
	ActorMemberID     string    `gorm:"index" json:"actor_member_id,omitempty"`
	ActorOrgID        string    `gorm:"index" json:"actor_org_id,omitempty"`
	ActorRole         string    `gorm:"index" json:"actor_role,omitempty"`
	ActorName         string    `gorm:"-" json:"actor_name,omitempty"`
	ActorOrganization string    `gorm:"-" json:"actor_organization,omitempty"`
	Action            string    `gorm:"not null;index" json:"action"`
	EntityType        string    `gorm:"not null;index" json:"entity_type"`
	EntityID          string    `gorm:"not null;index" json:"entity_id"`
	BeforeState       JSONMap   `gorm:"type:jsonb" json:"before_state,omitempty"`
	AfterState        JSONMap   `gorm:"type:jsonb" json:"after_state,omitempty"`
	Metadata          JSONMap   `gorm:"type:jsonb" json:"metadata,omitempty"`
	CreatedAt         time.Time `gorm:"index" json:"created_at"`
}

func (StakeholderAuditLog) TableName() string {
	return "stakeholder_audit_logs"
}

type StakeholderAuthorizationChallenge struct {
	ID              string     `gorm:"primaryKey" json:"id"`
	MemberID        string     `gorm:"not null;index" json:"member_id"`
	OrganizationID  string     `gorm:"not null;index" json:"organization_id"`
	StakeholderType string     `gorm:"not null" json:"stakeholder_type"`
	DashboardRole   string     `gorm:"not null;index" json:"dashboard_role"`
	Action          string     `gorm:"not null;index" json:"action"`
	EntityType      string     `gorm:"not null;index" json:"entity_type"`
	EntityID        string     `gorm:"not null;index" json:"entity_id"`
	AuthID          string     `gorm:"not null;index" json:"auth_id"`
	Status          string     `gorm:"not null;default:'pending';index" json:"status"`
	ExpiresAt       time.Time  `gorm:"not null;index" json:"expires_at"`
	VerifiedAt      *time.Time `json:"verified_at,omitempty"`
	UsedAt          *time.Time `json:"used_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func (StakeholderAuthorizationChallenge) TableName() string {
	return "stakeholder_authorization_challenges"
}

type StakeholderNotificationPreference struct {
	ID             string    `gorm:"primaryKey" json:"id"`
	MemberID       string    `gorm:"not null;uniqueIndex" json:"member_id"`
	OrganizationID string    `gorm:"not null;index" json:"organization_id"`
	EmailEnabled   bool      `gorm:"not null;default:true" json:"email_enabled"`
	InAppEnabled   bool      `gorm:"not null;default:true" json:"in_app_enabled"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (StakeholderNotificationPreference) TableName() string {
	return "stakeholder_notification_preferences"
}

type StakeholderAssetOperation struct {
	ID                string    `gorm:"primaryKey" json:"id"`
	AssetID           string    `gorm:"not null;index:idx_asset_operation_asset_org,unique" json:"asset_id"`
	AssetCode         string    `gorm:"index" json:"asset_code"`
	ManagerOrgID      string    `gorm:"not null;index:idx_asset_operation_asset_org,unique" json:"manager_org_id"`
	OperationalStatus string    `json:"operational_status,omitempty"`
	Notes             string    `json:"notes,omitempty"`
	Metadata          JSONMap   `gorm:"type:jsonb" json:"metadata,omitempty"`
	UpdatedByMemberID string    `gorm:"not null;index" json:"updated_by_member_id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (StakeholderAssetOperation) TableName() string {
	return "stakeholder_asset_operations"
}

type StakeholderReport struct {
	ID                  string     `gorm:"primaryKey" json:"id"`
	AssetID             string     `gorm:"index" json:"asset_id,omitempty"`
	AssetCode           string     `gorm:"index" json:"asset_code,omitempty"`
	ManagerOrgID        string     `gorm:"not null;index" json:"manager_org_id"`
	TrusteeOrgID        string     `gorm:"index" json:"trustee_org_id,omitempty"`
	ReportType          string     `gorm:"not null;index" json:"report_type"`
	Title               string     `gorm:"not null" json:"title"`
	DocumentID          *string    `gorm:"index" json:"document_id,omitempty"`
	FileURL             string     `json:"file_url,omitempty"`
	Status              string     `gorm:"not null;default:'generated';index" json:"status"`
	GeneratedByMemberID string     `gorm:"not null;index" json:"generated_by_member_id"`
	SubmittedAt         *time.Time `json:"submitted_at,omitempty"`
	SubmittedByMemberID *string    `gorm:"index" json:"submitted_by_member_id,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

func (StakeholderReport) TableName() string {
	return "stakeholder_reports"
}

// StakeholderStructuringStatus records a Legal/Financial Adviser's confirmation
// that their A5 workstream (legal documentation / financial structuring) is
// complete for an asset. One row per (asset, workstream) — see PRD §5A.3 / OI-14.
type StakeholderStructuringStatus struct {
	ID                  string     `gorm:"primaryKey" json:"id"`
	AssetID             string     `gorm:"not null;uniqueIndex:idx_structuring_asset_workstream,priority:1;index" json:"asset_id"`
	AssetCode           string     `gorm:"index" json:"asset_code,omitempty"`
	Workstream          string     `gorm:"not null;uniqueIndex:idx_structuring_asset_workstream,priority:2" json:"workstream"` // legal | financial
	Status              string     `gorm:"not null;default:'pending';index" json:"status"`
	OrgID               string     `gorm:"not null;index" json:"org_id"`
	StakeholderID       uint64     `gorm:"index" json:"stakeholder_id,omitempty"`
	DocumentID          *string    `gorm:"index" json:"document_id,omitempty"`
	Notes               string     `json:"notes,omitempty"`
	ConfirmedByMemberID *string    `gorm:"index" json:"confirmed_by_member_id,omitempty"`
	ConfirmedAt         *time.Time `json:"confirmed_at,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

func (StakeholderStructuringStatus) TableName() string {
	return "stakeholder_structuring_statuses"
}

// ComplianceRequirementTemplate defines the requirement items required of
// organizations matching (org_type, level). See PRD-kyc-compliance-templates.md
// §6.1/§7. `Version` increments on every edit to `Items` (no branching/rollback
// — it exists solely so a ComplianceRequirementInstance can record which
// definition it was created against).
type ComplianceRequirementTemplate struct {
	ID        string                              `gorm:"primaryKey" json:"id"`
	OrgType   string                              `gorm:"not null;uniqueIndex:idx_compliance_template_org_type_level" json:"org_type"`
	Level     int                                 `gorm:"not null;uniqueIndex:idx_compliance_template_org_type_level" json:"level"`
	Version   int                                 `gorm:"not null;default:1" json:"version"`
	CreatedAt time.Time                           `json:"created_at"`
	UpdatedAt time.Time                           `json:"updated_at"`
	Items     []ComplianceRequirementTemplateItem `gorm:"foreignKey:TemplateID" json:"items"`
}

func (ComplianceRequirementTemplate) TableName() string {
	return "compliance_requirement_templates"
}

type ComplianceRequirementTemplateItem struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	TemplateID  string    `gorm:"not null;index;constraint:compliance_requirement_template_items_template_id_fkey,OnDelete:CASCADE" json:"-"`
	Name        string    `gorm:"not null" json:"name"`
	Category    string    `gorm:"not null" json:"category"`
	Description string    `json:"description,omitempty"`
	InputType   string    `gorm:"not null" json:"input_type"`
	Required    *bool     `gorm:"not null;default:true" json:"required"`
	CreatedAt   time.Time `json:"created_at"`
}

func (ComplianceRequirementTemplateItem) TableName() string {
	return "compliance_requirement_template_items"
}

// ComplianceRequirementInstance is a requirement assigned to a specific
// organization — either generated from a ComplianceRequirementTemplate item,
// or created ad-hoc by an admin (TemplateItemID/TemplateVersion nil). Both
// kinds share the same submit -> review -> approve/reject/resubmit lifecycle.
// Mirrors PRD-kyc-compliance-templates.md §7.
type ComplianceRequirementInstance struct {
	ID                      string                  `gorm:"primaryKey" json:"id"`
	OrgID                   string                  `gorm:"not null;index;uniqueIndex:idx_compliance_requirement_org_template_item,priority:1" json:"org_id"`
	Category                string                  `gorm:"not null" json:"category"`
	Requirement             string                  `gorm:"not null" json:"requirement"`
	Description             string                  `json:"description,omitempty"`
	Required                *bool                   `gorm:"not null;default:true" json:"required"`
	TemplateItemID          *string                 `gorm:"index;uniqueIndex:idx_compliance_requirement_org_template_item,priority:2" json:"template_item_id,omitempty"`
	TemplateVersion         *int                    `json:"template_version,omitempty"`
	Status                  string                  `gorm:"not null;default:'pending';index" json:"status"`
	InputType               string                  `gorm:"not null;default:'text'" json:"input_type"`
	SubmissionData          JSONMap                 `gorm:"type:jsonb" json:"submission_data,omitempty"`
	SubmittedByMemberID     *string                 `gorm:"index" json:"submitted_by_member_id,omitempty"`
	SubmittedAt             *time.Time              `json:"submitted_at,omitempty"`
	RejectionReason         *string                 `json:"rejection_reason,omitempty"`
	ReviewedByAdminID       *string                 `json:"reviewed_by_admin_id,omitempty"`
	ReviewedAt              *time.Time              `json:"reviewed_at,omitempty"`
	CompletedAt             *time.Time              `json:"completed_at,omitempty"`
	CompletedByMemberID     *string                 `json:"completed_by_member_id,omitempty"`
	AssignedLevelAtApproval *int                    `json:"assigned_level_at_approval,omitempty"`
	DueDate                 *time.Time              `json:"due_date,omitempty"`
	CreatedAt               time.Time               `json:"created_at"`
	UpdatedAt               time.Time               `json:"updated_at"`
	Organization            coreModels.Organization `gorm:"foreignKey:OrgID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
}

func (ComplianceRequirementInstance) TableName() string {
	return "compliance_requirement_instances"
}

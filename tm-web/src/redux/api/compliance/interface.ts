// Shared compliance requirement template/instance types, imported by both the
// admin (baseApi) and organisation (orgApi) compliance slices so both sides of
// the submit -> review -> approve/reject/resubmit lifecycle agree on one shape.
// Mirrors PRD-kyc-compliance-templates.md §7 (ComplianceTemplate /
// ComplianceTemplateItem / ComplianceRequirementInstance).

export enum ComplianceRequirementStatus {
  Draft = "draft",
  Submitted = "submitted",
  UnderReview = "under_review",
  Approved = "approved",
  Rejected = "rejected",
  Pending = "pending",
  Overdue = "overdue",
  Waived = "waived",
}

export enum ComplianceInputType {
  DocumentUpload = "document_upload",
  Text = "text",
  StructuredForm = "structured_form",
}

export interface ComplianceTemplateItem {
  id: string;
  name: string;
  category: string;
  description?: string;
  input_type: ComplianceInputType | string;
  required: boolean;
}

export interface ComplianceTemplate {
  id: string;
  org_type: string;
  level: number;
  version: number;
  items: ComplianceTemplateItem[];
  created_at?: string;
  updated_at?: string;
}

export interface ComplianceRequirementInstance {
  id: string;
  org_id: string;
  category: string;
  requirement: string;
  template_item_id?: string;
  template_version?: number;
  status: ComplianceRequirementStatus | string;
  input_type?: ComplianceInputType | string;
  submission_data?: Record<string, unknown>;
  rejection_reason?: string;
  reviewed_by_admin_id?: string;
  reviewed_at?: string;
  completed_at?: string;
  completed_by_member_id?: string;
  assigned_level_at_approval?: number;
  due_date?: string;
  created_at?: string;
  updated_at?: string;
}

export interface CompliancePaginationMeta {
  page?: number;
  limit?: number;
  total?: number;
  total_pages?: number;
}

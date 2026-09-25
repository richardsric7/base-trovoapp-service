export interface TrusteeSummary {
  assets_under_trust?: number;
  pending_distribution_authorizations?: number;
  pending_due_diligence?: number;
  pending_fund_release_approvals?: number;
}

export interface TrusteeAssetActions {
  canApproveDueDiligence?: boolean;
  canApproveFundRelease?: boolean;
  canAuthorizeDistribution?: boolean;
  canExecuteFundRelease?: boolean;
  canUpdateAccount?: boolean;
  canRequestFundRelease?: boolean;
  canRecordRevenue?: boolean;
  canSubmitValuation?: boolean;
}

export interface TrusteeAssetAssignment {
  id?: string;
  asset_id?: string;
  asset_code?: string;
  trustee_org_id?: string;
  trustee_stakeholder_id?: number;
  asset_manager_org_id?: string;
  asset_manager_stakeholder_id?: number;
  custodian_org_id?: string;
  custodian_stakeholder_id?: number;
  status?: string;
  assigned_by?: string;
  created_at?: string;
  updated_at?: string;
}

export interface TrusteeRecentAsset {
  id?: string;
  assetCode?: string;
  assetName?: string;
  assetType?: string;
  assetSector?: string;
  assetSubSector?: string;
  assetTokenizationStatus?: number;
  assetQuoteCurrency?: string;
  assetCurrentValue?: number;
  valueOfTokenizedAsset?: number;
  assetManagerId?: number;
  approvedAssetCustodianId?: number;
  trusteeId?: number;
  createdAt?: string;
  updatedAt?: string;
  actions?: TrusteeAssetActions;
  assignment?: TrusteeAssetAssignment;
}

export interface TrusteePortfolioValue {
  currency?: string;
  amount?: string;
}

export interface TrusteePortfolioCategory {
  category?: string;
  asset_count?: number;
  values?: TrusteePortfolioValue[];
}

export interface TrusteePortfolioSummary {
  asset_count?: number;
  values?: TrusteePortfolioValue[];
  categories?: TrusteePortfolioCategory[];
}

export interface TrusteeRecentActivity {
  id?: string;
  source?: string;
  type?: string;
  title?: string;
  message?: string;
  status?: string;
  actor_organization?: string;
  actor_name?: string;
  entity_type?: string;
  entity_id?: string;
  created_at?: string;
}

export interface TrusteeDashboardData {
  role?: string;
  summary?: TrusteeSummary;
  recent_assets?: TrusteeRecentAsset[];
  portfolio_summary?: TrusteePortfolioSummary;
  recent_activity?: TrusteeRecentActivity[];
}

export interface TrusteeDashboardResponse {
  message?: string;
  data: TrusteeDashboardData;
  timestamp?: string;
  status?: string;
}

export interface ICurrencyTotal {
  amount: string | number;
  currency: string;
}

export interface ITrusteeFundManagementSummary {
  asset_count: number;
  milestone_payment_balance: ICurrencyTotal[];
  total_amount_from_primary_sales: ICurrencyTotal[];
  total_amount_paid_to_issuers: ICurrencyTotal[];
  total_amount_processed: ICurrencyTotal[];
  total_fees_generated: ICurrencyTotal[];
  total_income_from_assets: ICurrencyTotal[];
  total_paid_to_investors: ICurrencyTotal[];
}

export interface ITrusteeFundManagementSummaryResponse {
  data: ITrusteeFundManagementSummary;
}

export interface IDueDiligenceItem {
  notes?: string;
  id: string;
  checklist_id: string;
  category: string;
  item: string;
  status: string;
  created_at: string;
  updated_at: string;
}

export interface IDueDiligenceChecklist {
  id: string;
  asset_id: string;
  asset_code: string;
  trustee_org_id: string;
  status: string;
  created_at: string;
  updated_at: string;
  items: IDueDiligenceItem[];
}

export interface IDueDiligenceResponse {
  data: IDueDiligenceChecklist;
}

export interface IApproveDueDiligenceResponse {
  [key: string]: unknown;
}

export interface IRejectDueDiligenceRequest {
  reason: string;
}

export interface IUpdateDueDiligenceItemRequest {
  assetId: string;
  itemId: string;
  notes?: string;
  status: string;
}

export interface IVerifyDueDiligenceCategoryPayload {
  category: string;
  notes?: string;
  status: "pending" | "complete" | "failed" | "not_applicable" | string;
}

export interface IVerifyDueDiligenceCategoryRequest {
  assetId: string;
  payload: IVerifyDueDiligenceCategoryPayload;
}

export interface IVerifyDueDiligenceCategoryResponse {
  [key: string]: unknown;
}

export interface IFundReleaseSupportingDocument {
  category: string;
  download_path: string;
  external_url: string;
  id: string;
  mime_type: string;
  original_filename: string;
  size_bytes: number;
  title: string;
}

export interface IFundReleaseActor {
  member_id: string;
  name?: string;
  first_name?: string;
  last_name?: string;
  email?: string;
  organization_id: string;
  organization_name?: string;
}

export interface IFundReleaseRecord {
  id: string;
  asset_id: string;
  asset_code: string;
  requester_org_id: string;
  requester_member_id: string;
  trustee_org_id: string;
  custodian_org_id: string;
  amount: string;
  currency: string;
  purpose: string;
  status: string;
  supporting_document_ids: string[];
  receiving_account?: {
    account_name: string;
    account_number: string;
    bank: string;
  };
  requester?: IFundReleaseActor | null;
  reviewer?: IFundReleaseActor | null;
  supporting_documents?: IFundReleaseSupportingDocument[];
  reviewed_by_member_id?: string;
  reviewed_at?: string;
  rejection_reason?: string;
  executed_by_member_id?: string;
  executed_at?: string;
  execution_reference?: string;
  failure_reason?: string;
  created_at: string;
  updated_at: string;
}

export interface IFundReleaseMeta {
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

export interface IFundReleaseListResponse {
  data: {
    records: IFundReleaseRecord[];
    meta: IFundReleaseMeta;
  };
}

export interface IFundReleaseResponse {
  data: IFundReleaseRecord;
}

export interface IApproveFundReleaseRequest {
  challenge_id: string;
}

export interface IRejectFundReleaseRequest {
  reason: string;
}

export interface IFundReleaseActionResponse {
  [key: string]: unknown;
}

export interface ITrusteeFundReleaseListParams {
  page?: number;
  limit?: number;
  status?: string;
  asset_id?: string;
}

export interface ITrusteeDistribution {
  id: string;
  asset_id: string;
  asset_code: string;
  proposed_by_org_id: string;
  proposed_by_member_id: string;
  trustee_org_id: string;
  amount: string | number;
  currency: string;
  source: string;
  scheduled_date: string;
  status: string;
  authorized_at?: string;
  authorized_by_member_id?: string;
  rejection_reason?: string;
  created_at: string;
  updated_at: string;
}

export interface IDistributionCurrencyAmount {
  amount: string | number;
  currency: string;
}

export interface ITrusteeDistributionBreakdown {
  calculation: string;
  distribution_amount_per_token: IDistributionCurrencyAmount;
  net_income: IDistributionCurrencyAmount;
  total_token_holders: number;
  total_tokens: string;
}

export interface ITrusteeDistributionPayoutRecord {
  id: string;
  beneficiary_address: string;
  confirmed_token_balance: string;
  amount: string;
  currency: string;
  cannot_receive_asset: boolean;
  paid: boolean;
  created_at: string;
}

export interface ITrusteeDistributionPayout {
  distribution_id: string;
  tokenized_asset_id: string;
  batch?: string;
  amount: string;
  currency: string;
  amount_per_token: string;
  status: string;
  payment_schedule_ready: boolean;
  payout_completed: boolean;
  payouts: ITrusteeDistributionPayoutRecord[];
  created_at?: string;
  updated_at?: string;
}

export interface ITrusteeDistributionDetailResponse {
  data: {
    breakdown: ITrusteeDistributionBreakdown;
    distribution: ITrusteeDistribution;
    payout: ITrusteeDistributionPayout;
  };
}

export interface IDistributionMeta {
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

export interface IDistributionListParams {
  page: number;
  limit: number;
}

export interface ITrusteeDistributionListResponse {
  data: {
    records: ITrusteeDistribution[];
    meta: IDistributionMeta;
  };
}

export interface IRejectDistributionRequest {
  reason: string;
}

export interface IDistributionActionResponse {
  [key: string]: unknown;
}

export interface AssetManagerDashboardMetric {
  total?: number | string;
  pending?: number | string;
  value?: number | string;
}


export interface AssetManagerDashboardSummary {
  total_assets?: number;
  assets_under_management?: number | string;
  fees_generated?: number | string;
  milestone_verifications?: AssetManagerDashboardMetric | number;
  fund_release_requests?: AssetManagerDashboardMetric | number;
  distribution_requests?: AssetManagerDashboardMetric | number;
}

export interface AssetManagerDashboardAsset {
  asset_id?: string;
  asset_code?: string;
  asset_name?: string;
  asset_type?: string;
  asset_sector?: string;
  status?: string;
}

export interface AssetManagerPortfolioItem {
  asset_sector?: string;
  label?: string;
  count?: number;
  percentage?: number;
}

export interface AssetManagerDashboardActivity {
  id?: string;
  asset_code?: string;
  asset_name?: string;
  description?: string;
  activity_type?: string;
  requested_by?: string;
  status?: string;
  created_at?: string;
}

export interface AssetManagerDashboardData {
  summary?: AssetManagerDashboardSummary;
  recent_assets?: AssetManagerDashboardAsset[];
  portfolio_summary?: AssetManagerPortfolioItem[];
  recent_activities?: AssetManagerDashboardActivity[];
  currency?: string;
}

export interface AssetManagerDashboardResponse {
  message: string;
  data: AssetManagerDashboardData;
  timestamp: string;
  status: string;
}

// Detailed types matching the Asset Manager Dashboard API response (/stakeholder/asset-manager/dashboard)
export interface AmSummary {
  assets_managed?: number;
  pending_trustee_approvals?: number;
  ytd_revenue?: number | string;
}

export interface AmAssetActions {
  canApproveDueDiligence?: boolean;
  canApproveFundRelease?: boolean;
  canAuthorizeDistribution?: boolean;
  canExecuteFundRelease?: boolean;
  canUpdateAccount?: boolean;
  canRequestFundRelease?: boolean;
  canRecordRevenue?: boolean;
  canSubmitValuation?: boolean;
}

export interface AmAssetAssignment {
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

export interface AmRecentAsset {
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
  numberOfTokenToBeIssued?: number;
  assetManagerId?: number;
  approvedAssetCustodianId?: number;
  trusteeId?: number;
  tokenHolderCount?: number;
  complianceStatus?: string;
  custodyStatus?: string;
  createdAt?: string;
  updatedAt?: string;
  actions?: AmAssetActions;
  assignment?: AmAssetAssignment;
}

export interface AmPortfolioValue {
  currency?: string;
  amount?: string | number;
}

export interface AmPortfolioCategory {
  category?: string;
  asset_count?: number;
  values?: AmPortfolioValue[];
}

export interface AmPortfolioSummary {
  asset_count?: number;
  values?: AmPortfolioValue[];
  categories?: AmPortfolioCategory[];
}

export interface AmRecentActivity {
  id?: string;
  source?: string;
  type?: string;
  title?: string;
  status?: string;
  actor_organization?: string;
  actor_name?: string;
  entity_type?: string;
  entity_id?: string;
  created_at?: string;
}

export interface AmDashboardData {
  role?: string;
  summary?: AmSummary;
  recent_assets?: AmRecentAsset[];
  portfolio_summary?: AmPortfolioSummary;
  recent_activity?: AmRecentActivity[];
}

export interface AmDashboardResponse {
  message?: string;
  data: AmDashboardData;
  timestamp?: string;
  status?: string;
}

export interface AssetManagerApiResponse<T = Record<string, unknown>> {
  message?: string;
  data: T;
  timestamp?: string;
  status?: string;
}

export interface GetAssetManagerAssetsParams {
  page?: number;
  limit?: number;
  status?: string;
  asset_type?: string;
  asset_sector?: string;
  search?: string;
}

export interface UpdateAssetOperationalInfoPayload {
  metadata?: Record<string, unknown>;
  notes?: string;
  operational_status?: string;
}

export interface UpdateAssetOperationalInfoRequest {
  assetId: string;
  payload: UpdateAssetOperationalInfoPayload;
}

export interface CreateFundReleasePayload {
  amount: string;
  asset_id: string;
  currency: string;
  purpose: string;
  receiving_account_name?: string;
  receiving_account_number?: string;
  receiving_bank?: string;
  supporting_document_ids: string[];
}

export interface GenerateAssetManagerReportPayload {
  asset_code: string;
  asset_id: string;
  document_id: string;
  report_type: string;
  title: string;
}

export interface AssetManagerReportType {
  value: string;
  label: string;
}

export type AssetManagerReportTypesResponse =
  AssetManagerApiResponse<AssetManagerReportType[]>;

export interface AssetManagerAssetsData {
  records: AmRecentAsset[];
  meta: {
    page: number;
    limit: number;
    total: number;
    total_pages: number;
  };
  summary: {
    total: number;
    active: number;
    liquidated: number;
  };
}

export type AssetManagerAssetsResponse =
  AssetManagerApiResponse<AssetManagerAssetsData>;

export interface GetAssetManagerFundReleasesParams {
  page?: number;
  limit?: number;
  status?: string;
  asset_id?: string;
}

export interface AssetManagerFundReleaseRecord {
  requester?: import("../trustees/interface").IFundReleaseActor | null;
  reviewer?: import("../trustees/interface").IFundReleaseActor | null;
  id: string;
  amount: string | number;
  asset_code: string;
  asset_id: string;
  created_at: string;
  currency: string;
  custodian_org_id: string;
  executed_at?: string | null;
  executed_by_member_id?: string | null;
  execution_reference?: string | null;
  failure_reason?: string | null;
  purpose: string;
  rejection_reason?: string | null;
  requester_member_id: string;
  requester_org_id: string;
  reviewed_at?: string | null;
  reviewed_by_member_id?: string | null;
  status: string;
  supporting_document_ids: string[];
  trustee_org_id: string;
  updated_at: string;
}

export interface AssetManagerFundReleasesData {
  meta: {
    limit: number;
    page: number;
    total: number;
    total_pages: number;
  };
  records: AssetManagerFundReleaseRecord[];
}

export type AssetManagerFundReleasesResponse =
  AssetManagerApiResponse<AssetManagerFundReleasesData>;

export type AssetManagerFundReleaseResponse =
  AssetManagerApiResponse<AssetManagerFundReleaseRecord>;

export interface AssetManagerCurrencyTotal {
  amount: string;
  currency: string;
}

export interface AssetManagerFundReleaseAggregate {
  amounts: AssetManagerCurrencyTotal[];
  count: number;
}

export interface AssetManagerFundSummary {
  requested: AssetManagerFundReleaseAggregate;
  approved: AssetManagerFundReleaseAggregate;
  released: AssetManagerFundReleaseAggregate;
  pending: AssetManagerFundReleaseAggregate;
  rejected: AssetManagerFundReleaseAggregate;
  remaining_balance: AssetManagerCurrencyTotal[];
}

export type AssetManagerFundSummaryResponse =
  AssetManagerApiResponse<AssetManagerFundSummary>;

export interface AssetManagerReport {
  id: string;
  asset_id: string;
  asset_code: string;
  manager_org_id: string;
  trustee_org_id: string;
  report_type: string;
  title: string;
  document_id?: string | null;
  file_url?: string;
  status: "generated" | "submitted" | string;
  generated_by_member_id: string;
  submitted_at?: string | null;
  submitted_by_member_id?: string | null;
  created_at: string;
  updated_at: string;
}

export interface AssetManagerReportsData {
  records: AssetManagerReport[];
  meta: {
    page: number;
    limit: number;
    total: number;
    total_pages: number;
  };
}

export type AssetManagerReportsResponse =
  AssetManagerApiResponse<AssetManagerReportsData>;

export interface GetAssetManagerReportsParams {
  page?: number;
  limit?: number;
}

export interface RecordAssetRevenuePayload {
  amount: string;
  asset_id: string;
  collected_at: string;
  currency: string;
  period_end: string;
  period_start: string;
  source: string;
}

export interface AssetManagerRevenueRecord {
  id: string;
  asset_id: string;
  asset_code?: string;
  manager_org_id?: string;
  period_start?: string;
  period_end?: string;
  source: string;
  amount: string;
  currency: string;
  status: "draft" | "recorded" | "submitted_for_distribution" | string;
  collected_at?: string;
  created_at: string;
  updated_at?: string;
}

export interface AssetManagerRevenueData {
  records: AssetManagerRevenueRecord[];
  meta: {
    page: number;
    limit: number;
    total: number;
    total_pages: number;
  };
}

export type AssetManagerRevenueResponse =
  AssetManagerApiResponse<AssetManagerRevenueData>;

export interface GetAssetManagerRevenueParams {
  page?: number;
  limit?: number;
}

export interface SubmitRevenueDistributionPayload {
  amount: string;
  asset_id: string;
  currency: string;
  revenue_record_ids: string[];
  scheduled_date: string;
  source: string;
}

export interface SubmitAssetValuationPayload {
  asset_id: string;
  currency: string;
  methodology: string;
  report_document_id: string;
  valuation: string;
  valuation_date: string;
}

export interface AssetManagerValuationRecord {
  asset_code: string;
  asset_id: string;
  created_at: string;
  created_by_member_id: string;
  currency: string;
  id: string;
  manager_org_id: string;
  methodology: string;
  report_document_id: string;
  status: string;
  updated_at: string;
  valuation: string | number;
  valuation_date: string;
}

export interface AssetManagerValuationHistoryData {
  meta: {
    limit: number;
    page: number;
    total: number;
    total_pages: number;
  };
  records: AssetManagerValuationRecord[];
}

export type AssetManagerValuationHistoryResponse =
  AssetManagerApiResponse<AssetManagerValuationHistoryData>;

export type AssetManagerValuationResponse =
  AssetManagerApiResponse<AssetManagerValuationRecord>;

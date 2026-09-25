export interface FaAssetActions {
  canApproveDueDiligence?: boolean;
  canApproveFundRelease?: boolean;
  canAuthorizeDistribution?: boolean;
  canExecuteFundRelease?: boolean;
  canUpdateAccount?: boolean;
  canRequestFundRelease?: boolean;
  canRecordRevenue?: boolean;
  canSubmitValuation?: boolean;
}

export interface FaAssetAssignment {
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

export interface FaRecentAsset {
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
  numberOfTokenToBeIssued?: number;
  tokenHolderCount?: number;
  complianceStatus?: string;
  custodyStatus?: string;
  createdAt?: string;
  updatedAt?: string;
  actions?: FaAssetActions;
  assignment?: FaAssetAssignment;
}

export interface FaPortfolioValue {
  currency?: string;
  amount?: string | number;
}

export interface FaPortfolioCategory {
  category?: string;
  asset_count?: number;
  values?: FaPortfolioValue[];
}

export interface FaPortfolioSummary {
  asset_count?: number;
  values?: FaPortfolioValue[];
  categories?: FaPortfolioCategory[];
}

export interface FaRecentActivity {
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

export interface FaDashboardData {
  role?: string;
  summary?: Record<string, unknown>;
  recent_assets?: FaRecentAsset[];
  portfolio_summary?: FaPortfolioSummary;
  recent_activity?: FaRecentActivity[];
}

export interface FaDashboardResponse {
  message?: string;
  data: FaDashboardData;
  timestamp?: string;
  status?: string;
}

export interface CompleteAssetStructuringPayload {
  assetId: string;
  document_id?: string;
  title?: string;
  file_url?: string;
  notes: string;
}

export interface StakeholderStructuringStatus {
  id?: string;
  asset_id?: string;
  asset_code?: string;
  org_id?: string;
  stakeholder_id?: number;
  workstream?: string;
  status?: string;
  document_id?: string;
  notes?: string;
  confirmed_by_member_id?: string;
  confirmed_at?: string;
  created_at?: string;
  updated_at?: string;
}

export interface FaApiResponse<T = Record<string, unknown>> {
  message?: string;
  data: T;
  timestamp?: string;
  status?: string;
}

export type CompleteAssetStructuringResponse =
  FaApiResponse<StakeholderStructuringStatus>;

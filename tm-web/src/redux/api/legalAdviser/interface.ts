export interface LaAssetActions {
  canApproveDueDiligence?: boolean;
  canApproveFundRelease?: boolean;
  canAuthorizeDistribution?: boolean;
  canExecuteFundRelease?: boolean;
  canUpdateAccount?: boolean;
  canRequestFundRelease?: boolean;
  canRecordRevenue?: boolean;
  canSubmitValuation?: boolean;
}

export interface LaAssetAssignment {
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

export interface LaRecentAsset {
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
  actions?: LaAssetActions;
  assignment?: LaAssetAssignment;
}

export interface LaPortfolioValue {
  currency?: string;
  amount?: string | number;
}

export interface LaPortfolioCategory {
  category?: string;
  asset_count?: number;
  values?: LaPortfolioValue[];
}

export interface LaPortfolioSummary {
  asset_count?: number;
  values?: LaPortfolioValue[];
  categories?: LaPortfolioCategory[];
}

export interface LaRecentActivity {
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

export interface LaDashboardData {
  role?: string;
  summary?: Record<string, unknown>;
  recent_assets?: LaRecentAsset[];
  portfolio_summary?: LaPortfolioSummary;
  recent_activity?: LaRecentActivity[];
}

export interface LaDashboardResponse {
  message?: string;
  data: LaDashboardData;
  timestamp?: string;
  status?: string;
}

export interface CompleteAssetStructuringPayload {
  assetId: string;
  document_id: string;
  file_url: string;
  notes: string;
  title: string;
}

export interface LaApiResponse<T = Record<string, unknown>> {
  message?: string;
  data: T;
  timestamp?: string;
  status?: string;
}

export type CompleteAssetStructuringResponse = LaApiResponse<
  Record<string, unknown>
>;

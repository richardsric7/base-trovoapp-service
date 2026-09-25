export interface CustodianCurrencyAmount {
  currency: string;
  amount: string;
}

export interface CustodianDashboardSummary {
  total_asset_count?: number;
  assets_under_custody?: CustodianCurrencyAmount[];
  fees_generated?: CustodianCurrencyAmount[];
  account_balances?: CustodianCurrencyAmount[];
  total_account_balance?: string;
  segregated_accounts?: number;
  pending_fund_releases?: number;
  pending_compliance_items?: number;
}

export interface CustodianPortfolioCategory {
  category: string;
  asset_count: number;
  values: CustodianCurrencyAmount[];
}

export interface CustodianPortfolioSummary {
  asset_count: number;
  values: CustodianCurrencyAmount[];
  categories: CustodianPortfolioCategory[];
}

// Mirrors StakeholderAssetAssignment (backend/internal/components/stakeholder/models/db_models.go)
export interface CustodianAssetAssignment {
  id?: string;
  asset_id?: string;
  asset_code?: string;
  status?: string;
  assigned_by?: string;
  created_at?: string;
  updated_at?: string;
}

// Mirrors TokenizedAssetResponse (backend/internal/components/stakeholder/models/responses.go)
export interface CustodianDashboardAsset {
  id?: string;
  assetCode?: string;
  assetName?: string;
  assetType?: string;
  assetSector?: string;
  assetSubSector?: string;
  assetTokenizationStatus?: number;
  assetQuoteCurrency?: string;
  assetCurrentValue?: number;
  createdAt?: string;
  updatedAt?: string;
  assignment?: CustodianAssetAssignment;
}

// Mirrors DashboardActivity (backend/internal/components/stakeholder/models/responses.go)
export interface CustodianDashboardActivity {
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

export interface CustodianDashboardData {
  role?: string;
  summary?: CustodianDashboardSummary;
  portfolio_summary?: CustodianPortfolioSummary;
  recent_assets?: CustodianDashboardAsset[];
  recent_activity?: CustodianDashboardActivity[];
}

export interface CustodianDashboardResponse {
  data: CustodianDashboardData;
}

export interface CustodianPaginationMeta {
  page?: number;
  limit?: number;
  total?: number;
  total_pages?: number;
}

export interface GetCustodianComplianceParams {
  page?: number;
  limit?: number;
}

// Mirrors ComplianceItem (backend/internal/components/stakeholder/models/db_models.go)
export interface ComplianceItem {
  id: string;
  org_id: string;
  category: string;
  requirement: string;
  status: string;
  due_date?: string;
  completed_at?: string;
  completed_by_member_id?: string;
  created_at: string;
  updated_at: string;
}

export interface CustodianComplianceResponse {
  data: {
    records: ComplianceItem[];
    meta: CustodianPaginationMeta;
  };
}

export const ComplianceStatus = {
  Pending: "pending",
  Complete: "complete",
  Overdue: "overdue",
  Waived: "waived",
} as const;

export interface UpdateCustodianComplianceItemPayload {
  status: string;
}

export interface UpdateCustodianComplianceItemRequest {
  itemId: string;
  payload: UpdateCustodianComplianceItemPayload;
}

export const toCurrencyAmounts = (
  amounts?: CustodianCurrencyAmount[],
): { currency: string; amount: number }[] =>
  Array.isArray(amounts)
    ? amounts.map((item) => ({
        currency: (item.currency ?? "").toUpperCase(),
        amount: Number(item.amount ?? 0) || 0,
      }))
    : [];

export const portfolioCategoryLabel = (category: CustodianPortfolioCategory) =>
  category.category || "Uncategorized";

export const portfolioCategoryCount = (category: CustodianPortfolioCategory) =>
  category.asset_count ?? 0;

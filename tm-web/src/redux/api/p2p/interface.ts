// These interfaces mirror tm-api's P2P admin endpoints, which now read
// from the current P2P module's schema (app-backend's internal/components/p2p)
// instead of the old, dead legacy P2P trading platform's schema. Field
// names/types below match that module directly (decimal amounts as
// strings, orderStatus as a string enum, disputes rather than appeals).

export interface IP2PStatistics {
  totalOffers: number;
  activeOffers: number;
  totalOrders: number;
  completedOrders: number;
  openDisputes: number;
  resolvedDisputes: number;
}

export interface IP2PMetrics {
  message: string;
  data: IP2PStatistics;
  status: string;
  timestamp: string;
}

export interface ITopTrader {
  merchantUsername: string;
  completedTrades: number;
}

export interface P2POrder {
  id: string;
  offerId: string;
  customerUserId: string;
  customerUsername: string;
  merchantUserId: string;
  merchantUsername: string;
  offerType: "BUY" | "SELL";
  asset: string;
  currency: string;
  price: string;
  specifiedAssetAmount: string;
  paymentAmount: string;
  orderStatus: string;
  isDisputed: boolean;
  createdAt: string;
}

export interface ITradeStatistics {
  totalTrades: number;
  completedTrades: number;
  cancelledTrades: number;
  expiredTrades: number;
  tradeSuccessRate: number;
  openDisputes: number;
  topTraders: ITopTrader[];
  recentTrades: P2POrder[];
}

export interface ITradeStatisticsResponse {
  message: string;
  data: ITradeStatistics;
  status: string;
  timestamp: string;
}

export interface TradeListQueryParams {
  page?: number;
  pageSize?: number;
  offerType?: string;
  status?: string;
  username?: string;
  createdAt?: string;
}

export interface TradeListResponse {
  message: string;
  data: {
    data: P2POrder[];
    total: number;
    page: number;
    pageSize: number;
  };
  status: string;
  timestamp: string;
}

export interface P2PUser {
  id: string;
  username: string;
  email: string;
  phone: string;
  registrationDate: string;
  suspended: boolean;
  isMerchant: boolean;
  merchantCompletedTrades: number;
  merchantCompletionRate: string;
  customerCompletedTrades: number;
  customerCompletionRate: string;
}

export interface P2PUserListQueryParams {
  page?: number;
  pageSize?: number;
  username?: string;
  email?: string;
  phone?: string;
  search?: string;
}

export interface P2PUserListResponse {
  message: string;
  data: {
    data: P2PUser[];
    total: number;
    page: number;
    pageSize: number;
  };
  status: string;
  timestamp: string;
}

// --- P2P Market Reports ---
// These mirror tm-api's /p2p/reports/* endpoints (usermetrics/db/reports.go),
// gated server-side on the ACCESS_REPORTS admin permission.

export type ReportRange = "7d" | "30d" | "90d" | "1y";

export interface ReportRangeParams {
  range?: ReportRange;
}

export interface IVolumeReportPoint {
  date: string;
  totalOrders: number;
  completedOrders: number;
  completedVolume: string;
}

export interface IVolumeReport {
  series: IVolumeReportPoint[];
  totalOrders: number;
  completedOrders: number;
  completedVolume: string;
}

export interface ICountItem {
  key: string;
  count: number;
}

export interface IVolumeItem {
  key: string;
  count: number;
  volume: string;
}

export interface IDistributionReport {
  byAsset: IVolumeItem[];
  byCurrency: IVolumeItem[];
  byCountry: IVolumeItem[];
  byStatus: ICountItem[];
}

export interface IDisputeReportPoint {
  date: string;
  opened: number;
  resolved: number;
}

export interface IDisputeReport {
  series: IDisputeReportPoint[];
  totalOpened: number;
  totalResolved: number;
  stillOpen: number;
  resolutionRate: number;
  averageResolutionTimeSeconds: number;
  bySubject: ICountItem[];
  byResolution: ICountItem[];
}

export interface IRevenueReportPoint {
  date: string;
  platformFee: string;
  regulatoryFee: string;
  vat: string;
  total: string;
}

export interface IRevenueReport {
  series: IRevenueReportPoint[];
  totalPlatformFee: string;
  totalRegulatoryFee: string;
  totalVat: string;
  totalRevenue: string;
}

export interface IGrowthReportPoint {
  date: string;
  newOffers: number;
  newMerchants: number;
}

export interface IGrowthReport {
  series: IGrowthReportPoint[];
  totalNewOffers: number;
  totalNewMerchants: number;
}

export interface IMerchantLeaderboardEntry {
  merchantId: string;
  merchantUsername: string;
  completedTrades: number;
  completedTradeVolume: string;
  completionRate: string;
  averageOrderCompletionTime: number;
  cancelledOrders: number;
  expiredOrders: number;
  disputesOpened: number;
  disputesResolvedAgainstMerchant: number;
  lastActivityAt: string | null;
}

export interface MerchantLeaderboardQueryParams {
  page?: number;
  pageSize?: number;
  sortBy?: "completedTrades" | "completedVolume" | "completionRate" | "disputesOpened";
}

export interface MerchantLeaderboardResponse {
  message: string;
  data: {
    data: IMerchantLeaderboardEntry[];
    total: number;
    page: number;
    pageSize: number;
  };
  status: string;
  timestamp: string;
}

export interface IHourlyActivityPoint {
  hour: number;
  count: number;
}

export interface IHourlyActivityReport {
  series: IHourlyActivityPoint[];
  peakHour: number;
}

// A generic wrapper for the single-payload report responses (volume,
// distribution, disputes, revenue, growth, hourly activity) - all share
// tm-api's serverResponse.JSON envelope shape.
export interface ReportResponse<T> {
  message: string;
  data: T;
  status: string;
  timestamp: string;
}

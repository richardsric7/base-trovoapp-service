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

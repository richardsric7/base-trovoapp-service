export interface IP2PMetrics {
  data: {
    average_time_per_session: number;
    daily_active_and_new_users: {
      active_users: number;
      new_users: number;
    };
    monthly_active_and_new_users: {
      active_users: number;
      new_users: number;
    };
    total_active_users: number;
    total_sessions: number;
    total_users: number;
    weekly_active_and_new_users: {
      active_users: number;
      new_users: number;
    };
    weekly_percentage_change: {
      new_users_percentage_change: number;
      active_users_percentage_change: number;
      total_sessions_percentage_change: number;
      average_time_percentage_change: number;
    };
  };
  message: string;
  status: string;
  timestamp: string;
}

export interface ITopTrader {
  offer_maker: string;
  order_count: number;
}

export interface IOrderStatus {
  status: string;
}

export interface IRecentTrade {
  id: string;
  orderType: string;
  CreatedAt: string;
  UpdatedAt: string;
  ExpiresAt: string;
  AcceptedAt: string;
  cancelAfter: string;
  offerId: string;
  offerType: string;
  offerMaker: string;
  offerMakerPhone: string;
  offerMakerCountryCode: string | null;
  offerMaxTimePerTransaction: number;
  offerTaker: string;
  offerTakerPhone: string;
  offerPaymentMethodId: string | null;
  offerPaymentChannelId: string | null;
  offerCurrencyPaymentMethodId: string;
  offerCurrencyPaymentChannelId: string;
  offerCurrencyPaymentMethodName: string;
  offerCurrencyPaymentMethodDestinationAccount: string;
  offerPaymentCurrencyMethodBankName: string;
  offerCurrencyPaymentMethodAccountOpeningBranch: string;
  offerCurrencyPaymentMethodCountryCode: string;
  offerCurrencyID: string;
  offerAssetAmount: number;
  offerAssetId: string;
  offerAssetPrice: number;
  offerMinTradeAmount: number;
  offerMaxTradeAmount: number;
  offerRemark: string;
  takerPaymentMethodId: string;
  takerPaymentChannelId: string;
  takerPaymentMethodDestinationAccount: string;
  orderEscrowAddress: string;
  orderPaymentMemo: string;
  orderAmount: number;
  orderMakerFee: number;
  orderTakerFee: number;
  orderEscrowTransactionId: string | null;
  orderAssetReleaseTransactionId: string | null;
  orderStatusId: number;
  orderStatus: IOrderStatus;
  dynamicLink: string | null;
  qrCode: string | null;
  FiatDepositTransactionID: string | null;
}

export interface ITradeStatisticsResponse {
  total_trades: number;
  completed_trades: number;
  cancelled_trades: number;
  total_volume: number;
  average_trade_size: number;
  trade_success_rate: number;
  trades_on_appeal: number;
  top_traders: ITopTrader[];
  recent_trades: IRecentTrade[];
}

export interface P2PUser {
  id: string;
  username: string;
  email: string;
  phone: string;
  registration_date: string;
  account_status: number;
  merchant_status: number;
  violations: string | null;
  city: string;
}

export interface P2PUserListResponse {
  message: string;
  data: {
    data: P2PUser[];
    total: number;
    page: number;
    pageSize: number;
  };
  timestamp: string;
  status: string;
}

export interface P2PUserListQueryParams {
  page?: number;
  pageSize?: number;
  username?: string;
  email?: string;
  phone?: string;
  first_name?: string;
  last_name?: string;
  city?: string;
  address?: string;
  suspended?: number;
  kyc_level?: number;
  admin_level?: number;
  registration_date_from?: string;
  registration_date_to?: string;
  search?: string;
}
export interface TradeStatus {
  status: string;
}

export interface TradeOrder {
  id: string;
  orderType: "BUY" | "SELL";
  CreatedAt: string;
  ExpiresAt: string;
  AcceptedAt: string;
  cancelAfter: string;

  offerId: string;
  offerType: "BUY" | "SELL";
  offerMaker: string;
  offerMakerPhone: string;
  offerMakerCountryCode: string;
  offerMaxTimePerTransaction: number;

  offerTaker: string;
  offerTakerPhone: string;

  offerPaymentMethodId: string;
  offerPaymentChannelId: string;
  offerPaymentMethodName: string | null;
  offerPaymentMethodDestinationAccount: string;
  offerPaymentMethodMemo: string | null;
  offerPaymentMethodBankName: string | null;
  offerPaymentMethodAccountOpeningBranch: string | null;

  offerCurrencyPaymentMethodId: string;
  offerCurrencyPaymentChannelId: string;
  offerCurrencyPaymentMethodName: string;
  offerCurrencyPaymentMethodDestinationAccount: string;
  offerCurrencyPaymentMethodMemo: string | null;
  offerPaymentCurrencyMethodBankName: string;
  offerCurrencyPaymentMethodAccountOpeningBranch: string | null;
  offerCurrencyPaymentMethodCountryCode: string;
  offerCurrencyPaymentMethodCurrencyId: string | null;

  offerCurrencyID: string;
  offerAssetAmount: number;
  offerAssetId: string;
  offerAssetPrice: number;
  offerMinTradeAmount: number;
  offerMaxTradeAmount: number;
  offerRemark: string;

  takerPaymentMethodId: string;
  takerPaymentChannelId: string;
  takerPaymentMethodName: string;
  takerPaymentMethodDestinationAccount: string;
  takerPaymentMethodMemo: string | null;
  takerPaymentMethodBankName: string;
  takerPaymentMethodAccountOpeningBranch: string;
  takerPaymentMethodCountryCode: string;
  takerPaymentMethodCurrencyId: string;

  orderEscrowAddress: string;
  orderPaymentMemo: string;
  orderAmount: number;
  orderMakerFee: number;
  orderTakerFee: number;
  orderEscrowTransactionId: string | null;
  orderAssetReleaseTransactionId: string | null;

  orderStatusId: number;
  orderStatus: TradeStatus;

  dynamicLink: string | null;
  qrCode: string | null;
  FiatDepositTransactionID: string | null;
}

export interface TradeListQueryParams {
  page?: string | number;
  pageSize?: string | number;
  orderType?: string; // or restrict to "BUY" | "SELL"
  offerMaker?: string;
  offerTaker?: string;
  status?: string;
  [key: string]: string | number | undefined;
}
export interface TradeListResponse {
  message: string;
  data: {
    data: TradeOrder[];
    page: number;
    page_size: number;
    total: number;
  };
}

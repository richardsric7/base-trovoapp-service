import { baseApi } from "@/redux/baseApi";
import {
  IP2PMetrics,
  ITradeStatisticsResponse,
  P2PUserListQueryParams,
  P2PUserListResponse,
  TradeListQueryParams,
  TradeListResponse,
} from "./interface";

export const p2pApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    p2pMetrics: builder.query<IP2PMetrics, void>({
      query: () => ({
        url: "/p2p/statistics",
        method: "GET",
      }),
    }),

    trades: builder.query<TradeListResponse, TradeListQueryParams>({
      query: (params) => ({
        url: "/orders/trade",
        method: "GET",
        params,
      }),
    }),

    tradeStatistics: builder.query<ITradeStatisticsResponse, void>({
      query: () => ({
        url: "/trades/statistics",
        method: "GET",
      }),
    }),

    p2pUsers: builder.query<P2PUserListResponse, P2PUserListQueryParams>({
      query: (params) => ({
        url: "/p2p/users",
        method: "GET",
        params,
      }),
    }),
  }),
});

export const {
  useP2pMetricsQuery,
  useTradesQuery,
  useTradeStatisticsQuery,
  useP2pUsersQuery,
} = p2pApi;

import { baseApi } from "@/redux/baseApi";
import {
  IP2PMetrics,
  ITradeStatisticsResponse,
  P2PUserListQueryParams,
  P2PUserListResponse,
  TradeListQueryParams,
  TradeListResponse,
  ReportRangeParams,
  ReportResponse,
  IVolumeReport,
  IDistributionReport,
  IDisputeReport,
  IRevenueReport,
  IGrowthReport,
  IHourlyActivityReport,
  MerchantLeaderboardQueryParams,
  MerchantLeaderboardResponse,
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

    // --- P2P Market Reports ---
    p2pVolumeReport: builder.query<ReportResponse<IVolumeReport>, ReportRangeParams>({
      query: (params) => ({
        url: "/p2p/reports/volume",
        method: "GET",
        params,
      }),
    }),

    p2pDistributionReport: builder.query<ReportResponse<IDistributionReport>, ReportRangeParams>({
      query: (params) => ({
        url: "/p2p/reports/distribution",
        method: "GET",
        params,
      }),
    }),

    p2pDisputeReport: builder.query<ReportResponse<IDisputeReport>, ReportRangeParams>({
      query: (params) => ({
        url: "/p2p/reports/disputes",
        method: "GET",
        params,
      }),
    }),

    p2pRevenueReport: builder.query<ReportResponse<IRevenueReport>, ReportRangeParams>({
      query: (params) => ({
        url: "/p2p/reports/revenue",
        method: "GET",
        params,
      }),
    }),

    p2pGrowthReport: builder.query<ReportResponse<IGrowthReport>, ReportRangeParams>({
      query: (params) => ({
        url: "/p2p/reports/growth",
        method: "GET",
        params,
      }),
    }),

    p2pHourlyActivityReport: builder.query<ReportResponse<IHourlyActivityReport>, ReportRangeParams>({
      query: (params) => ({
        url: "/p2p/reports/hourly-activity",
        method: "GET",
        params,
      }),
    }),

    p2pMerchantLeaderboard: builder.query<MerchantLeaderboardResponse, MerchantLeaderboardQueryParams>({
      query: (params) => ({
        url: "/p2p/reports/merchants",
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
  useP2pVolumeReportQuery,
  useP2pDistributionReportQuery,
  useP2pDisputeReportQuery,
  useP2pRevenueReportQuery,
  useP2pGrowthReportQuery,
  useP2pHourlyActivityReportQuery,
  useP2pMerchantLeaderboardQuery,
} = p2pApi;

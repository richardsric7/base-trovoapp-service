import { baseApi } from "@/redux/baseApi";
import {
  ICountriesStats,
  IReferrerCountResponse,
  IUsersMetrics,
  IWalletDistributionResponse,
  RecentRegistration,
  RecentRegistrationsResponse,
} from "./interface";

export const metricsApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    usersMetrics: builder.query<IUsersMetrics, void>({
      query: () => ({
        url: "/users/metrics",
        method: "GET",
      }),
    }),

    walletDist: builder.query<IWalletDistributionResponse, void>({
      query: () => ({
        url: "/wallet/distribution",
        method: "GET",
      }),
    }),

    referrerCount: builder.query<IReferrerCountResponse, void>({
      query: () => ({
        url: "/referrer/count",
        method: "GET",
      }),
    }),

    countiresStatistics: builder.query<ICountriesStats, void>({
      query: () => ({
        url: "/country/statistics",
        method: "GET",
      }),
    }),
    recentReg: builder.query<
      RecentRegistrationsResponse,
      { page?: number; pageSize?: number }
    >({
      query: (params) => ({
        url: "/recent/registrations",
        method: "GET",
        params: params,
      }),
    }),
  }),
});

export const {
  useUsersMetricsQuery,
  useWalletDistQuery,
  useReferrerCountQuery,
  useCountiresStatisticsQuery,
  useRecentRegQuery,
} = metricsApi;

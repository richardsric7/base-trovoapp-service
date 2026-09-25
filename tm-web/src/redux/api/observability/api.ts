import { baseApi } from "@/redux/baseApi";
import { tagTypes } from "@/redux/baseApi/tagTypes";
import {
  ICapabilities,
  IEnvelope,
  IHealthResponse,
  IIssuesResponse,
  IRecentErrorsResponse,
  ITraceResponse,
  ITrendsResponse,
} from "./interface";

export const observabilityApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getObservabilityCapabilities: builder.query<IEnvelope<ICapabilities>, void>({
      query: () => ({ url: "/observability/capabilities", method: "GET" }),
      providesTags: [tagTypes.OBSERVABILITY],
    }),

    getSystemHealth: builder.query<IEnvelope<IHealthResponse>, void>({
      query: () => ({ url: "/observability/health", method: "GET" }),
      providesTags: [tagTypes.OBSERVABILITY],
    }),

    getSystemIssues: builder.query<
      IEnvelope<IIssuesResponse>,
      { app?: string; limit?: number } | void
    >({
      query: (params) => ({
        url: "/observability/issues",
        method: "GET",
        params: params || undefined,
      }),
      providesTags: [tagTypes.OBSERVABILITY],
    }),

    getSystemTrends: builder.query<IEnvelope<ITrendsResponse>, { window: string }>({
      query: (params) => ({
        url: "/observability/trends",
        method: "GET",
        params,
      }),
      providesTags: [tagTypes.OBSERVABILITY],
    }),

    getRequestTrace: builder.query<IEnvelope<ITraceResponse>, { requestId: string }>({
      query: ({ requestId }) => ({
        url: `/observability/trace/${encodeURIComponent(requestId)}`,
        method: "GET",
      }),
      providesTags: [tagTypes.OBSERVABILITY],
    }),

    getRecentErrors: builder.query<IEnvelope<IRecentErrorsResponse>, { limit?: number } | void>({
      query: (params) => ({
        url: "/observability/recent-errors",
        method: "GET",
        params: params || undefined,
      }),
      providesTags: [tagTypes.OBSERVABILITY],
    }),
  }),
  overrideExisting: false,
});

export const {
  useGetObservabilityCapabilitiesQuery,
  useGetSystemHealthQuery,
  useGetSystemIssuesQuery,
  useGetSystemTrendsQuery,
  useGetRequestTraceQuery,
  useLazyGetRequestTraceQuery,
  useGetRecentErrorsQuery,
} = observabilityApi;

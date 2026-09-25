import { orgApi } from "@/redux/baseApi/orgApi";
import type { IFundReleaseResponse } from "../trustees/interface";
import {
  AmDashboardResponse,
  AssetManagerApiResponse,
  AssetManagerAssetsResponse,
  AssetManagerDashboardResponse,
  CreateFundReleasePayload,
  GetAssetManagerFundReleasesParams,
  AssetManagerFundReleaseResponse,
  AssetManagerFundReleasesResponse,
  AssetManagerFundSummaryResponse,
  GenerateAssetManagerReportPayload,
  GetAssetManagerAssetsParams,
  GetAssetManagerReportsParams,
  AssetManagerReportsResponse,
  AssetManagerReportTypesResponse,
  AssetManagerRevenueResponse,
  GetAssetManagerRevenueParams,
  RecordAssetRevenuePayload,
  SubmitAssetValuationPayload,
  SubmitRevenueDistributionPayload,
  UpdateAssetOperationalInfoRequest,
  AssetManagerValuationHistoryResponse,
  AssetManagerValuationResponse,
} from "./interface";

type AssetManagerResponse = AssetManagerApiResponse<Record<string, unknown>>;

export const assetManagerApi = orgApi.injectEndpoints({
  endpoints: (builder) => ({
    getAssetManagerFundRelease: builder.query<IFundReleaseResponse, string>({
      query: (requestId) => ({
        url: `/stakeholder/asset-manager/fund-releases/${encodeURIComponent(requestId)}`,
        method: "GET",
      }),
      providesTags: ["stakeholders"],
    }),
    getAssetManagerDashboard: builder.query<
      AmDashboardResponse,
      void
    >({
      query: () => ({
        url: "/stakeholder/asset-manager/dashboard",
        method: "GET",
      }),
    }),
    getAssetManagerAssets: builder.query<
      AssetManagerAssetsResponse,
      GetAssetManagerAssetsParams | void
    >({
      query: (params) => ({
        url: "/stakeholder/asset-manager/assets",
        method: "GET",
        params: params || undefined,
      }),
      providesTags: ["stakeholders"],
    }),
    updateAssetManagerAsset: builder.mutation<
      AssetManagerResponse,
      UpdateAssetOperationalInfoRequest
    >({
      query: ({ assetId, payload }) => ({
        url: `/stakeholder/asset-manager/assets/${assetId}`,
        method: "PUT",
        data: payload,
      }),
      invalidatesTags: ["stakeholders"],
    }),
    getAssetManagerFundReleases: builder.query<
      AssetManagerFundReleasesResponse,
      GetAssetManagerFundReleasesParams | void
    >({
      query: (params) => ({
        url: "/stakeholder/asset-manager/fund-releases",
        method: "GET",
        params: params || undefined,
      }),
      providesTags: ["stakeholders"],
    }),
    createAssetManagerFundRelease: builder.mutation<
      AssetManagerFundReleaseResponse,
      CreateFundReleasePayload
    >({
      query: (payload) => ({
        url: "/stakeholder/asset-manager/fund-releases",
        method: "POST",
        data: payload,
      }),
      invalidatesTags: ["stakeholders"],
    }),
    getAssetManagerFundSummary: builder.query<
      AssetManagerFundSummaryResponse,
      { asset_id?: string } | void
    >({
      query: (params) => ({
        url: "/stakeholder/asset-manager/fund-releases/summary",
        method: "GET",
        params: params || undefined,
      }),
      providesTags: ["stakeholders"],
    }),
    getAssetManagerReports: builder.query<
      AssetManagerReportsResponse,
      GetAssetManagerReportsParams | void
    >({
      query: (params) => ({
        url: "/stakeholder/asset-manager/reports",
        method: "GET",
        params: params || undefined,
      }),
      providesTags: ["stakeholders"],
    }),
    getAssetManagerReportTypes: builder.query<
      AssetManagerReportTypesResponse,
      void
    >({
      query: () => ({
        url: "/stakeholder/asset-manager/report-types",
        method: "GET",
      }),
    }),
    generateAssetManagerReport: builder.mutation<
      AssetManagerResponse,
      GenerateAssetManagerReportPayload
    >({
      query: (payload) => ({
        url: "/stakeholder/asset-manager/reports",
        method: "POST",
        data: payload,
      }),
      invalidatesTags: ["stakeholders"],
    }),
    submitAssetManagerReport: builder.mutation<
      AssetManagerResponse,
      string
    >({
      query: (reportId) => ({
        url: `/stakeholder/asset-manager/reports/${reportId}/submit`,
        method: "POST",
      }),
      invalidatesTags: ["stakeholders"],
    }),
    getAssetManagerRevenue: builder.query<
      AssetManagerRevenueResponse,
      GetAssetManagerRevenueParams | void
    >({
      query: (params) => ({
        url: "/stakeholder/asset-manager/revenue",
        method: "GET",
        params: params || undefined,
      }),
      providesTags: ["stakeholders"],
    }),
    recordAssetManagerRevenue: builder.mutation<
      AssetManagerResponse,
      RecordAssetRevenuePayload
    >({
      query: (payload) => ({
        url: "/stakeholder/asset-manager/revenue",
        method: "POST",
        data: payload,
      }),
      invalidatesTags: ["stakeholders"],
    }),
    submitAssetManagerRevenueDistribution: builder.mutation<
      AssetManagerResponse,
      SubmitRevenueDistributionPayload
    >({
      query: (payload) => ({
        url: "/stakeholder/asset-manager/revenue/submit-distribution",
        method: "POST",
        data: payload,
      }),
      invalidatesTags: ["stakeholders"],
    }),
    submitAssetManagerValuation: builder.mutation<
      AssetManagerValuationResponse,
      SubmitAssetValuationPayload
    >({
      query: (payload) => ({
        url: "/stakeholder/asset-manager/valuations",
        method: "POST",
        data: payload,
      }),
      invalidatesTags: ["stakeholders"],
    }),
    getAssetManagerValuationHistory: builder.query<
      AssetManagerValuationHistoryResponse,
      string
    >({
      query: (assetId) => ({
        url: `/stakeholder/asset-manager/valuations/${assetId}`,
        method: "GET",
      }),
      providesTags: ["stakeholders"],
    }),
    requestIndependentAssetValuation: builder.mutation<
      AssetManagerValuationResponse,
      string
    >({
      query: (valuationId) => ({
        url: `/stakeholder/asset-manager/valuations/${valuationId}/request-independent`,
        method: "POST",
      }),
      invalidatesTags: ["stakeholders"],
    }),
  }),
  overrideExisting: false,
});

export const {
  useGetAssetManagerFundReleaseQuery,
  useGetAssetManagerDashboardQuery,
  useGetAssetManagerAssetsQuery,
  useUpdateAssetManagerAssetMutation,
  useGetAssetManagerFundReleasesQuery,
  useCreateAssetManagerFundReleaseMutation,
  useGetAssetManagerFundSummaryQuery,
  useGetAssetManagerReportsQuery,
  useGetAssetManagerReportTypesQuery,
  useGenerateAssetManagerReportMutation,
  useSubmitAssetManagerReportMutation,
  useGetAssetManagerRevenueQuery,
  useRecordAssetManagerRevenueMutation,
  useSubmitAssetManagerRevenueDistributionMutation,
  useSubmitAssetManagerValuationMutation,
  useGetAssetManagerValuationHistoryQuery,
  useRequestIndependentAssetValuationMutation,
} = assetManagerApi;

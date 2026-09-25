import { orgApi } from "@/redux/baseApi/orgApi";
import {
  IApproveDueDiligenceResponse,
  IApproveFundReleaseRequest,
  IDueDiligenceResponse,
  IFundReleaseActionResponse,
  IFundReleaseListResponse,
  IFundReleaseResponse,
  IDistributionActionResponse,
  IDistributionListParams,
  IRejectDistributionRequest,
  ITrusteeDistributionListResponse,
  ITrusteeDistributionDetailResponse,
  ITrusteeFundManagementSummaryResponse,
  ITrusteeFundReleaseListParams,
  IRejectDueDiligenceRequest,
  IRejectFundReleaseRequest,
  IUpdateDueDiligenceItemRequest,
  IVerifyDueDiligenceCategoryRequest,
  IVerifyDueDiligenceCategoryResponse,
  TrusteeDashboardResponse,
} from "./interface";
import { tagTypes } from "@/redux/baseApi/tagTypes";

export const trusteeApi = orgApi.injectEndpoints({
  endpoints: (builder) => ({
    getTrusteeDashboard: builder.query<TrusteeDashboardResponse, void>({
      query: () => ({
        url: "/stakeholder/trustee/dashboard",
        method: "GET",
      }),
    }),
    getTrusteeFundManagementSummary: builder.query<
      ITrusteeFundManagementSummaryResponse,
      void
    >({
      query: () => ({
        url: "/stakeholder/trustee/fund-management/summary",
        method: "GET",
      }),
    }),
    getTrusteeDueDiligence: builder.query<IDueDiligenceResponse, string>({
      query: (assetId) => ({
        url: `/stakeholder/trustee/due-diligence/${assetId}`,
        method: "GET",
      }),
    }),
    approveTrusteeDueDiligence: builder.mutation<
      IApproveDueDiligenceResponse,
      string
    >({
      query: (assetId) => ({
        url: `/stakeholder/trustee/due-diligence/${assetId}/approve`,
        method: "POST",
      }),
    }),
    rejectTrusteeDueDiligence: builder.mutation<
      IApproveDueDiligenceResponse,
      { assetId: string; payload: IRejectDueDiligenceRequest }
    >({
      query: ({ assetId, payload }) => ({
        url: `/stakeholder/trustee/due-diligence/${assetId}/reject`,
        method: "POST",
        body: payload,
      }),
    }),
    updateTrusteeDueDiligenceItem: builder.mutation<
      IApproveDueDiligenceResponse,
      IUpdateDueDiligenceItemRequest
    >({
      query: ({ assetId, itemId, ...body }) => ({
        url: `/stakeholder/trustee/due-diligence/${assetId}/items/${itemId}`,
        method: "PUT",
        body,
      }),
    }),
    verifyTrusteeDueDiligenceCategory: builder.mutation<
      IVerifyDueDiligenceCategoryResponse,
      IVerifyDueDiligenceCategoryRequest
    >({
      query: ({ assetId, payload }) => ({
        url: `/stakeholder/trustee/due-diligence/${assetId}/categories`,
        method: "PUT",
        body: payload,
        data: payload,
      }),
    }),

    // Fund Releases
    getTrusteeFundReleases: builder.query<
      IFundReleaseListResponse,
      ITrusteeFundReleaseListParams | void
    >({
      query: (params) => ({
        url: "/stakeholder/trustee/fund-releases",
        method: "GET",
        params: params || undefined,
      }),
      providesTags: [tagTypes.FUND_RELEASES],
    }),
    getTrusteeFundRelease: builder.query<IFundReleaseResponse, string>({
      query: (requestId) => ({
        url: `/stakeholder/trustee/fund-releases/${requestId}`,
        method: "GET",
      }),
      providesTags: [tagTypes.FUND_RELEASES],
    }),
    approveTrusteeFundRelease: builder.mutation<
      IFundReleaseActionResponse,
      { requestId: string; payload: IApproveFundReleaseRequest }
    >({
      query: ({ requestId, payload }) => ({
        url: `/stakeholder/trustee/fund-releases/${requestId}/approve`,
        method: "POST",
        body: payload,
      }),
      invalidatesTags: [tagTypes.FUND_RELEASES],
    }),
    rejectTrusteeFundRelease: builder.mutation<
      IFundReleaseActionResponse,
      { requestId: string; payload: IRejectFundReleaseRequest }
    >({
      query: ({ requestId, payload }) => ({
        url: `/stakeholder/trustee/fund-releases/${requestId}/reject`,
        method: "POST",
        body: payload,
      }),
      invalidatesTags: [tagTypes.FUND_RELEASES],
    }),

    // Distributions
    getTrusteeDistributions: builder.query<
      ITrusteeDistributionListResponse,
      IDistributionListParams
    >({
      query: (params) => ({
        url: "/stakeholder/trustee/distributions",
        method: "GET",
        params,
      }),
      providesTags: ["distributions"],
    }),
    getTrusteeDistributionHistory: builder.query<
      ITrusteeDistributionListResponse,
      IDistributionListParams
    >({
      query: (params) => ({
        url: "/stakeholder/trustee/distributions/history",
        method: "GET",
        params,
      }),
      providesTags: ["distributions"],
    }),
    getTrusteeDistribution: builder.query<
      ITrusteeDistributionDetailResponse,
      string
    >({
      query: (distributionId) => ({
        url: `/stakeholder/trustee/distributions/${distributionId}`,
        method: "GET",
      }),
      providesTags: ["distributions"],
    }),
    downloadTrusteeDistributionPayouts: builder.mutation<Blob, string>({
      query: (distributionId) => ({
        url: `/stakeholder/trustee/distributions/${encodeURIComponent(distributionId)}/payouts/download`,
        method: "GET",
        responseType: "blob",
      }),
    }),
    authorizeTrusteeDistribution: builder.mutation<
      IDistributionActionResponse,
      { distributionId: string; challengeId: string }
    >({
      query: ({ distributionId, challengeId }) => ({
        url: `/stakeholder/trustee/distributions/${distributionId}/authorize`,
        method: "POST",
        body: { challenge_id: challengeId },
      }),
      invalidatesTags: ["distributions"],
    }),
    rejectTrusteeDistribution: builder.mutation<
      IDistributionActionResponse,
      { distributionId: string; payload: IRejectDistributionRequest }
    >({
      query: ({ distributionId, payload }) => ({
        url: `/stakeholder/trustee/distributions/${distributionId}/reject`,
        method: "POST",
        body: payload,
      }),
      invalidatesTags: ["distributions"],
    }),
  }),

  overrideExisting: false,
});

export const {
  useGetTrusteeDashboardQuery,
  useGetTrusteeFundManagementSummaryQuery,
  useGetTrusteeDueDiligenceQuery,
  useLazyGetTrusteeDueDiligenceQuery,
  useApproveTrusteeDueDiligenceMutation,
  useRejectTrusteeDueDiligenceMutation,
  useUpdateTrusteeDueDiligenceItemMutation,
  useVerifyTrusteeDueDiligenceCategoryMutation,
  useGetTrusteeFundReleasesQuery,
  useGetTrusteeFundReleaseQuery,
  useLazyGetTrusteeFundReleaseQuery,
  useApproveTrusteeFundReleaseMutation,
  useRejectTrusteeFundReleaseMutation,
  useGetTrusteeDistributionsQuery,
  useGetTrusteeDistributionHistoryQuery,
  useGetTrusteeDistributionQuery,
  useDownloadTrusteeDistributionPayoutsMutation,
  useAuthorizeTrusteeDistributionMutation,
  useRejectTrusteeDistributionMutation,
} = trusteeApi;

import { orgApi } from "@/redux/baseApi/orgApi";
import {
  IStakeholderAssetDetailResponse,
  IStakeholderProfileResponse,
  IUpdateStakeholderProfileRequest,
  IUpdateNotificationPreferencesRequest,
  IStakeholderNotificationPreferenceResponse,
  IStakeholderAssetsResponse,
  IStakeholderAssetsQueryParams,
  IStakeholderAuditTrailResponse,
  IStakeholderNotificationsResponse,
  IMarkStakeholderNotificationReadResponse,
  IAuthorizationChallengeResponse,
  ICreateAuthorizationRequest,
  IVerifyAuthorizationRequest,
  IStakeholderDashboardResponse,
  IVerifyWalletLinkRequest,
  IVerifyWalletLinkResponse,
  IWalletLinkAuthorizationResponse,
  IStakeholderDocumentCategory,
  IStakeholderDocumentListParams,
  IStakeholderDocumentResponse,
  IStakeholderDocumentsResponse,
  ICreateStakeholderDocumentRequest,
  IStakeholderAssetActivitiesQueryParams,
  IStakeholderAssetActivitiesResponse,
} from "./interface";
import { tagTypes } from "@/redux/baseApi/tagTypes";

export const sharedStakeholderApi = orgApi.injectEndpoints({
  endpoints: (builder) => ({
    getRatingAgencyDashboard: builder.query<
      IStakeholderDashboardResponse,
      void
    >({
      query: () => ({
        url: "/stakeholder/rating-agency/dashboard",
        method: "GET",
      }),
    }),
    getIssuingHouseDashboard: builder.query<
      IStakeholderDashboardResponse,
      void
    >({
      query: () => ({
        url: "/stakeholder/issuing-house/dashboard",
        method: "GET",
      }),
    }),
    getStakeholderProfile: builder.query<IStakeholderProfileResponse, void>({
      query: () => ({
        url: `/stakeholder/shared/profile`,
        method: "GET",
      }),
      providesTags: [tagTypes.STAKEHOLDERS],
    }),
    updateStakeholderProfile: builder.mutation<
      IStakeholderProfileResponse,
      IUpdateStakeholderProfileRequest
    >({
      query: (payload) => ({
        url: "/stakeholder/shared/profile",
        method: "PUT",
        body: payload,
      }),
      invalidatesTags: [tagTypes.STAKEHOLDERS],
    }),
    updateStakeholderNotificationPreferences: builder.mutation<
      IStakeholderNotificationPreferenceResponse,
      IUpdateNotificationPreferencesRequest
    >({
      query: (payload) => ({
        url: "/stakeholder/shared/profile/notifications",
        method: "PUT",
        body: payload,
      }),
    }),
    getStakeholderAssets: builder.query<
      IStakeholderAssetsResponse,
      IStakeholderAssetsQueryParams | void
    >({
      query: (params) => ({
        url: `/stakeholder/shared/assets`,
        method: "GET",
        params: params || undefined,
      }),
    }),
    getStakeholderAssetDetail: builder.query<
      IStakeholderAssetDetailResponse,
      string
    >({
      query: (assetId) => ({
        url: `/stakeholder/shared/assets/${assetId}`,
        method: "GET",
      }),
    }),
    getStakeholderAssetActivities: builder.query<
      IStakeholderAssetActivitiesResponse,
      IStakeholderAssetActivitiesQueryParams
    >({
      query: ({ assetId, ...params }) => ({
        url: `/stakeholder/shared/assets/${assetId}/activities`,
        method: "GET",
        params,
      }),
    }),
    getStakeholderDocumentCategories: builder.query<
      { data: IStakeholderDocumentCategory[] },
      void
    >({
      query: () => ({
        url: "/stakeholder/shared/document-categories",
        method: "GET",
      }),
    }),
    getStakeholderDocuments: builder.query<
      IStakeholderDocumentsResponse,
      IStakeholderDocumentListParams | void
    >({
      query: (params) => ({
        url: "/stakeholder/shared/documents",
        method: "GET",
        params: params || undefined,
      }),
      providesTags: [tagTypes.STAKEHOLDERS],
    }),
    getStakeholderDocument: builder.query<IStakeholderDocumentResponse, string>({
      query: (documentId) => ({
        url: `/stakeholder/shared/documents/${documentId}`,
        method: "GET",
      }),
    }),
    createStakeholderDocument: builder.mutation<
      IStakeholderDocumentResponse,
      ICreateStakeholderDocumentRequest
    >({
      query: (payload) => ({
        url: "/stakeholder/shared/documents",
        method: "POST",
        data: payload,
      }),
      invalidatesTags: [tagTypes.STAKEHOLDERS],
    }),
    uploadStakeholderDocument: builder.mutation<
      IStakeholderDocumentResponse,
      FormData
    >({
      query: (payload) => ({
        url: "/stakeholder/shared/documents/upload",
        method: "POST",
        data: payload,
      }),
      invalidatesTags: [tagTypes.STAKEHOLDERS],
    }),
    downloadStakeholderDocument: builder.mutation<Blob, string>({
      query: (documentId) => ({
        url: `/stakeholder/shared/documents/${documentId}/download`,
        method: "GET",
        responseType: "blob",
      }),
    }),
    requestWalletLinkAuthorization: builder.mutation<
      IWalletLinkAuthorizationResponse,
      Pick<IVerifyWalletLinkRequest, "wallet_username">
    >({
      query: (payload) => ({
        url: "/organizations/members/link-wallet/request-auth",
        method: "POST",
        body: payload,
      }),
    }),
    verifyWalletLink: builder.mutation<
      IVerifyWalletLinkResponse,
      IVerifyWalletLinkRequest
    >({
      query: (payload) => ({
        url: "/organizations/members/link-wallet/verify",
        method: "POST",
        body: payload,
      }),
      invalidatesTags: [tagTypes.STAKEHOLDERS],
    }),
    getStakeholderAuditTrail: builder.query<
      IStakeholderAuditTrailResponse,
      void
    >({
      query: () => ({
        url: `/stakeholder/shared/audit-trail`,
        method: "GET",
      }),
    }),
    getStakeholderNotifications: builder.query<
      IStakeholderNotificationsResponse,
      void
    >({
      query: () => ({
        url: `/stakeholder/shared/notifications`,
        method: "GET",
      }),
      providesTags: [tagTypes.STAKEHOLDERS],
    }),
    markStakeholderNotificationRead: builder.mutation<
      IMarkStakeholderNotificationReadResponse,
      string
    >({
      query: (id) => ({
        url: `/stakeholder/shared/notifications/${id}/read`,
        method: "PUT",
      }),
      invalidatesTags: [tagTypes.STAKEHOLDERS],
    }),
    createStakeholderAuthorization: builder.mutation<
      IAuthorizationChallengeResponse,
      ICreateAuthorizationRequest
    >({
      query: (payload) => ({
        url: "/stakeholder/shared/authorizations",
        method: "POST",
        body: payload,
      }),
    }),
    verifyStakeholderAuthorization: builder.mutation<
      IAuthorizationChallengeResponse,
      IVerifyAuthorizationRequest
    >({
      query: ({ challengeId, ...payload }) => ({
        url: `/stakeholder/shared/authorizations/${challengeId}/verify`,
        method: "POST",
        body: payload,
      }),
    }),
  }),
  overrideExisting: false,
});

export const {
  useGetRatingAgencyDashboardQuery,
  useGetIssuingHouseDashboardQuery,
  useGetStakeholderProfileQuery,
  useLazyGetStakeholderProfileQuery,
  useUpdateStakeholderProfileMutation,
  useUpdateStakeholderNotificationPreferencesMutation,
  useRequestWalletLinkAuthorizationMutation,
  useVerifyWalletLinkMutation,
  useGetStakeholderAssetsQuery,
  useLazyGetStakeholderAssetsQuery,
  useGetStakeholderAssetDetailQuery,
  useLazyGetStakeholderAssetDetailQuery,
  useGetStakeholderAssetActivitiesQuery,
  useGetStakeholderDocumentCategoriesQuery,
  useGetStakeholderDocumentsQuery,
  useLazyGetStakeholderDocumentQuery,
  useCreateStakeholderDocumentMutation,
  useUploadStakeholderDocumentMutation,
  useDownloadStakeholderDocumentMutation,
  useGetStakeholderAuditTrailQuery,
  useLazyGetStakeholderAuditTrailQuery,
  useGetStakeholderNotificationsQuery,
  useLazyGetStakeholderNotificationsQuery,
  useMarkStakeholderNotificationReadMutation,
  useCreateStakeholderAuthorizationMutation,
  useVerifyStakeholderAuthorizationMutation,
} = sharedStakeholderApi;

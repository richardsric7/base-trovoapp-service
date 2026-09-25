import { baseApi } from "@/redux/baseApi";
import {
  VetTokenizedAssetRequest,
  TokenizationQueryParams,
  TokenizationResponse,
  TokenizationDetailResponse,
  DueDiligenceFailReason,
  UpdateTokenizationPayload,
  TokenizationParamsResponse,
  TokenizedAssetStatisticsResponse,
} from "./interface";

export const assetTokenization = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getTokenizationList: builder.query<
      TokenizationResponse,
      TokenizationQueryParams
    >({
      query: (params) => ({
        url: "/tokenization/list",
        method: "GET",
        params,
      }),
    }),
    getTokenizationParams: builder.query<TokenizationParamsResponse, void>({
      query: () => ({
        url: "/tokenization",
        method: "GET",
      }),
    }),

    getTokenizationDetail: builder.query<TokenizationDetailResponse, string>({
      query: (tokenizedAssetID) => ({
        url: `/tokenization/detail/${tokenizedAssetID}`,
        method: "GET",
      }),
      providesTags: ["tokenization"],
    }),

    completeVetting: builder.mutation<
      TokenizationResponse,
      VetTokenizedAssetRequest
    >({
      query: ({ tokenizedAssetID, vettingPayload }) => ({
        url: `/tokenization/vet/${tokenizedAssetID}`,
        method: "PUT",
        data: vettingPayload,
      }),
    }),

    confirmPayment: builder.mutation({
      query: (tokenizedAssetID) => ({
        url: `/tokenization/fee/${tokenizedAssetID}`,
        method: "POST",
      }),
    }),

    approveTokenizationAndMint: builder.mutation({
      query: (tokenizedAssetID) => ({
        url: `/tokenization/mint/${tokenizedAssetID}`,
        method: "POST",
      }),
    }),

    FailDueDiligence: builder.mutation<
      TokenizationResponse,
      DueDiligenceFailReason
    >({
      query: ({ tokenizedAssetID, reason }) => ({
        url: `/tokenization/failed/${tokenizedAssetID}`,
        method: "POST",
        data: { reason },
      }),
    }),

    // updateTokenizationInfo: builder.mutation<
    //   any,
    //   { tokenizedAssetID: string; payload: Partial<UpdateTokenizationPayload> }
    // >({
    //   query: ({ tokenizedAssetID, payload }) => ({
    //     url: `/tokenization/update/${tokenizedAssetID}`,
    //     method: "PUT",
    //     data: payload,
    //   }),
    //   invalidatesTags: ["tokenization"],
    // }),

    // api/assettokenization.ts

    updateTokenizationInfo: builder.mutation<
      any,
      { tokenizedAssetID: string; data: UpdateTokenizationPayload }
    >({
      query: ({ tokenizedAssetID, data }) => ({
        url: `/tokenization/update/${tokenizedAssetID}`,
        method: "PUT",
        data: data,
      }),
      invalidatesTags: ["tokenization"],
    }),

    deleteTokenization: builder.mutation<any, string>({
      query: (tokenizationID) => ({
        url: `/tokenization/${tokenizationID}`,
        method: "DELETE",
      }),
      // Optionally, invalidate cache/tag if your UI needs to refetch list/details after deletion
      invalidatesTags: ["tokenization"],
    }),

    updateTokenizationSalesDates: builder.mutation<
      void,
      { tokenizedAssetID: string; salesStart?: string; salesEnd?: string }
    >({
      query: ({ tokenizedAssetID, salesStart, salesEnd }) => ({
        url: `/tokenization/salesdate/${tokenizedAssetID}`,
        method: "PUT",
        data: {
          salesStart: salesStart || "",
          salesEnd: salesEnd || "",
        },
      }),
      invalidatesTags: ["tokenization"],
    }),
    updateLogo: builder.mutation<
      void,
      { tokenizedAssetID: string; documentFile: File }
    >({
      query: ({ tokenizedAssetID, documentFile }) => {
        const formData = new FormData();
        formData.append("documentFile", documentFile);

        return {
          url: `/tokenization/logo/${tokenizedAssetID}`,
          method: "PUT",
          data: formData,
          transformRequest: [(data: any) => data],
        };
      },
      invalidatesTags: ["tokenization"],
    }),

    uploadTokenizationDocument: builder.mutation<
      any,
      {
        tokenizedAssetID: string;
        documentType: string;
        documentTitle: string;
        documentFile: File;
      }
    >({
      query: ({
        tokenizedAssetID,
        documentType,
        documentTitle,
        documentFile,
      }) => {
        const formData = new FormData();
        formData.append("tokenizedAssetId", tokenizedAssetID);
        formData.append("documentType", documentType);
        formData.append("documentTitle", documentTitle);
        formData.append("documentFile", documentFile);

        return {
          url: `/tokenization/document`,
          method: "PUT",
          data: formData,

          transformRequest: [(data: any) => data],
        };
      },
      invalidatesTags: ["tokenization"],
    }),

    getTokenizationStatistics: builder.query<
      TokenizedAssetStatisticsResponse,
      void
    >({
      query: () => ({
        url: `/tokenization/statistics`,
        method: "GET",
      }),
      providesTags: ["tokenization"],
    }),
  }),
});

export const {
  useGetTokenizationListQuery,
  useCompleteVettingMutation,
  useConfirmPaymentMutation,
  useApproveTokenizationAndMintMutation,
  useGetTokenizationDetailQuery,
  useFailDueDiligenceMutation,
  useUpdateTokenizationInfoMutation,
  useUpdateTokenizationSalesDatesMutation,
  useUpdateLogoMutation,
  useUploadTokenizationDocumentMutation,
  useGetTokenizationParamsQuery,
  useDeleteTokenizationMutation,
  useGetTokenizationStatisticsQuery,
} = assetTokenization;

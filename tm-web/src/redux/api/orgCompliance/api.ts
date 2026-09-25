import { orgApi } from "@/redux/baseApi/orgApi";
import { tagTypes } from "@/redux/baseApi/tagTypes";
import { IStakeholderDocumentResponse } from "@/redux/api/sharedstakeholders";
import {
  IOrgComplianceListParams,
  IOrgComplianceListResponse,
  ISubmitComplianceItemRequest,
  IUpdateOrgComplianceItemRequest,
} from "./interface";

export const orgComplianceApi = orgApi.injectEndpoints({
  endpoints: (builder) => ({
    getOrgCompliance: builder.query<
      IOrgComplianceListResponse,
      IOrgComplianceListParams | void
    >({
      query: (params) => ({
        url: "/stakeholder/org/compliance",
        method: "GET",
        params: params || undefined,
      }),
      providesTags: [tagTypes.COMPLIANCE_REQUIREMENTS],
    }),
    submitOrgComplianceItem: builder.mutation<
      Record<string, unknown>,
      ISubmitComplianceItemRequest
    >({
      query: ({ itemId, payload }) => ({
        url: `/stakeholder/org/compliance/${itemId}/submit`,
        method: "POST",
        data: payload,
      }),
      invalidatesTags: [tagTypes.COMPLIANCE_REQUIREMENTS],
    }),
    updateOrgComplianceItem: builder.mutation<
      Record<string, unknown>,
      IUpdateOrgComplianceItemRequest
    >({
      query: ({ itemId, payload }) => ({
        url: `/stakeholder/org/compliance/${itemId}`,
        method: "PUT",
        data: payload,
      }),
      invalidatesTags: [tagTypes.COMPLIANCE_REQUIREMENTS],
    }),
    // Dedicated compliance-document upload (category is forced server-side),
    // distinct from the generic shared/documents/upload endpoint.
    uploadComplianceDocument: builder.mutation<
      IStakeholderDocumentResponse,
      FormData
    >({
      query: (payload) => ({
        url: "/stakeholder/org/compliance/documents/upload",
        method: "POST",
        data: payload,
      }),
    }),
    downloadComplianceDocument: builder.mutation<Blob, string>({
      query: (documentId) => ({
        url: `/stakeholder/org/compliance/documents/${encodeURIComponent(documentId)}/download`,
        method: "GET",
        responseType: "blob",
      }),
    }),
  }),
  overrideExisting: false,
});

export const {
  useGetOrgComplianceQuery,
  useSubmitOrgComplianceItemMutation,
  useUpdateOrgComplianceItemMutation,
  useUploadComplianceDocumentMutation,
  useDownloadComplianceDocumentMutation,
} = orgComplianceApi;

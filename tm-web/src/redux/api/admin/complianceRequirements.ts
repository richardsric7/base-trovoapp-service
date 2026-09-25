import { baseApi } from "@/redux/baseApi";
import { tagTypes } from "@/redux/baseApi/tagTypes";
import {
  ComplianceRequirementInstance,
  CompliancePaginationMeta,
} from "@/redux/api/compliance/interface";

export interface IComplianceRequirementListParams {
  page?: number;
  limit?: number;
  org_id?: string;
  category?: string;
  level?: number;
  status?: string;
}

export interface IComplianceRequirementListResponse {
  data: {
    records: ComplianceRequirementInstance[];
    meta: CompliancePaginationMeta;
  };
}

export interface IComplianceRequirementDetailResponse {
  data: ComplianceRequirementInstance;
}

export interface IReviewComplianceRequirementRequest {
  id: string;
  payload: {
    decision: "approved" | "rejected";
    rejection_reason?: string;
  };
}

export interface ICreateComplianceRequirementRequest {
  org_id: string;
  category: string;
  requirement: string;
  description?: string;
  input_type: string;
  required?: boolean;
  due_date?: string;
}

export const complianceRequirementsApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getComplianceRequirements: builder.query<
      IComplianceRequirementListResponse,
      IComplianceRequirementListParams | void
    >({
      query: (params) => ({
        url: "/stakeholder/admin/compliance-requirements",
        method: "GET",
        params: params || undefined,
      }),
      providesTags: [tagTypes.COMPLIANCE_REQUIREMENTS],
    }),
    getComplianceRequirementDetail: builder.query<
      IComplianceRequirementDetailResponse,
      string
    >({
      query: (id) => ({
        url: `/stakeholder/admin/compliance-requirements/${id}`,
        method: "GET",
      }),
      providesTags: [tagTypes.COMPLIANCE_REQUIREMENTS],
    }),
    reviewComplianceRequirement: builder.mutation<
      IComplianceRequirementDetailResponse,
      IReviewComplianceRequirementRequest
    >({
      query: ({ id, payload }) => ({
        url: `/stakeholder/admin/compliance-requirements/${id}/review`,
        method: "PUT",
        data: payload,
      }),
      invalidatesTags: [tagTypes.COMPLIANCE_REQUIREMENTS],
    }),
    // Transitions a submitted requirement into under_review. The backend no
    // longer does this as a side effect of GET (reads must stay idempotent),
    // so the review queue calls this explicitly when opening a submitted item.
    startComplianceRequirementReview: builder.mutation<
      IComplianceRequirementDetailResponse,
      string
    >({
      query: (id) => ({
        url: `/stakeholder/admin/compliance-requirements/${id}/review/start`,
        method: "PUT",
      }),
      invalidatesTags: [tagTypes.COMPLIANCE_REQUIREMENTS],
    }),
    // Ad-hoc assignment to any organisation (not custodian-only), sharing the
    // same instance lifecycle as templated requirements.
    createComplianceRequirement: builder.mutation<
      IComplianceRequirementDetailResponse,
      ICreateComplianceRequirementRequest
    >({
      query: (payload) => ({
        url: "/stakeholder/admin/compliance-requirements",
        method: "POST",
        data: payload,
      }),
      invalidatesTags: [tagTypes.COMPLIANCE_REQUIREMENTS],
    }),
    downloadComplianceRequirementDocument: builder.mutation<
      Blob,
      { requirementId: string; documentId: string }
    >({
      query: ({ requirementId, documentId }) => ({
        url: `/stakeholder/admin/compliance-requirements/${encodeURIComponent(requirementId)}/documents/${encodeURIComponent(documentId)}/download`,
        method: "GET",
        responseType: "blob",
      }),
    }),
  }),
  overrideExisting: false,
});

export const {
  useGetComplianceRequirementsQuery,
  useGetComplianceRequirementDetailQuery,
  useReviewComplianceRequirementMutation,
  useStartComplianceRequirementReviewMutation,
  useCreateComplianceRequirementMutation,
  useDownloadComplianceRequirementDocumentMutation,
} = complianceRequirementsApi;

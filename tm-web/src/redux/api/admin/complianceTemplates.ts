import { baseApi } from "@/redux/baseApi";
import { tagTypes } from "@/redux/baseApi/tagTypes";
import {
  ComplianceTemplate,
  ComplianceTemplateItem,
  CompliancePaginationMeta,
} from "@/redux/api/compliance/interface";

export interface IComplianceTemplateListParams {
  page?: number;
  limit?: number;
  org_type?: string;
  level?: number;
}

export interface IComplianceTemplateListResponse {
  data: {
    records: ComplianceTemplate[];
    meta: CompliancePaginationMeta;
  };
}

export interface IComplianceTemplateResponse {
  data: ComplianceTemplate;
}

export interface ICreateComplianceTemplateRequest {
  org_type: string;
  level: number;
  items: Omit<ComplianceTemplateItem, "id">[];
}

export interface IUpdateComplianceTemplateRequest {
  id: string;
  payload: {
    org_type?: string;
    level?: number;
    items: Omit<ComplianceTemplateItem, "id">[];
  };
}

export const complianceTemplatesApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getComplianceTemplates: builder.query<
      IComplianceTemplateListResponse,
      IComplianceTemplateListParams | void
    >({
      query: (params) => ({
        url: "/stakeholder/admin/compliance-templates",
        method: "GET",
        params: params || undefined,
      }),
      providesTags: [tagTypes.COMPLIANCE_TEMPLATES],
    }),
    createComplianceTemplate: builder.mutation<
      IComplianceTemplateResponse,
      ICreateComplianceTemplateRequest
    >({
      query: (payload) => ({
        url: "/stakeholder/admin/compliance-templates",
        method: "POST",
        data: payload,
      }),
      invalidatesTags: [tagTypes.COMPLIANCE_TEMPLATES],
    }),
    // Bumps the template's `version` server-side. Existing
    // ComplianceRequirementInstance rows already assigned keep the
    // template_version they were created under and are NOT retroactively
    // changed by this — only requirement instances assigned from this point
    // forward use the new version.
    updateComplianceTemplate: builder.mutation<
      IComplianceTemplateResponse,
      IUpdateComplianceTemplateRequest
    >({
      query: ({ id, payload }) => ({
        url: `/stakeholder/admin/compliance-templates/${id}`,
        method: "PUT",
        data: payload,
      }),
      invalidatesTags: [tagTypes.COMPLIANCE_TEMPLATES],
    }),
  }),
  overrideExisting: false,
});

export const {
  useGetComplianceTemplatesQuery,
  useCreateComplianceTemplateMutation,
  useUpdateComplianceTemplateMutation,
} = complianceTemplatesApi;

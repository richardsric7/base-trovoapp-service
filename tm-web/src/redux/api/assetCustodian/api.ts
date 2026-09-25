import { orgApi } from "@/redux/baseApi/orgApi";
import {
  CustodianComplianceResponse,
  CustodianDashboardResponse,
  GetCustodianComplianceParams,
  UpdateCustodianComplianceItemRequest,
} from "./interface";

export const assetCustodianApi = orgApi.injectEndpoints({
  endpoints: (builder) => ({
    getCustodianDashboard: builder.query<CustodianDashboardResponse, void>({
      query: () => ({
        url: "/stakeholder/custodian/dashboard",
        method: "GET",
      }),
      providesTags: ["stakeholders"],
    }),
    getCustodianCompliance: builder.query<
      CustodianComplianceResponse,
      GetCustodianComplianceParams | void
    >({
      query: (params) => ({
        url: "/stakeholder/custodian/compliance",
        method: "GET",
        params: params || undefined,
      }),
      providesTags: ["stakeholders"],
    }),
    updateCustodianComplianceItem: builder.mutation<
      Record<string, unknown>,
      UpdateCustodianComplianceItemRequest
    >({
      query: ({ itemId, payload }) => ({
        url: `/stakeholder/custodian/compliance/${itemId}`,
        method: "PUT",
        data: payload,
      }),
      invalidatesTags: ["stakeholders"],
    }),
  }),
  overrideExisting: false,
});

export const {
  useGetCustodianDashboardQuery,
  useGetCustodianComplianceQuery,
  useUpdateCustodianComplianceItemMutation,
} = assetCustodianApi;

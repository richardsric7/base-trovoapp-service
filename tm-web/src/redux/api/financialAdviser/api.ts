import { orgApi } from "@/redux/baseApi/orgApi";
import {
  CompleteAssetStructuringPayload,
  CompleteAssetStructuringResponse,
  FaDashboardResponse,
} from "./interface";

export const financialAdviserApi = orgApi.injectEndpoints({
  endpoints: (builder) => ({
    getFinancialAdviserDashboard: builder.query<FaDashboardResponse, void>({
      query: () => ({
        url: "/stakeholder/financial-adviser/dashboard",
        method: "GET",
      }),
      providesTags: ["stakeholders"],
    }),
    completeAssetStructuring: builder.mutation<
      CompleteAssetStructuringResponse,
      CompleteAssetStructuringPayload
    >({
      query: ({ assetId, ...payload }) => ({
        url: `/stakeholder/financial-adviser/structuring/${assetId}/complete`,
        method: "POST",
        data: payload,
      }),
      invalidatesTags: ["stakeholders"],
    }),
  }),
  overrideExisting: false,
});

export const {
  useGetFinancialAdviserDashboardQuery,
  useCompleteAssetStructuringMutation,
} = financialAdviserApi;

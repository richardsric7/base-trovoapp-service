import { orgApi } from "@/redux/baseApi/orgApi";
import {
  CompleteAssetStructuringPayload,
  CompleteAssetStructuringResponse,
  LaDashboardResponse,
} from "./interface";

export const legalAdviserApi = orgApi.injectEndpoints({
  endpoints: (builder) => ({
    getLegalAdviserDashboard: builder.query<LaDashboardResponse, void>({
      query: () => ({
        url: "/stakeholder/legal-adviser/dashboard",
        method: "GET",
      }),
      providesTags: ["stakeholders"],
    }),
    completeAssetStructuring: builder.mutation<
      CompleteAssetStructuringResponse,
      CompleteAssetStructuringPayload
    >({
      query: ({ assetId, ...payload }) => ({
        url: `/stakeholder/legal-adviser/structuring/${assetId}/complete`,
        method: "POST",
        data: payload,
      }),
      invalidatesTags: ["stakeholders"],
    }),
  }),
  overrideExisting: false,
});

export const {
  useGetLegalAdviserDashboardQuery,
  useCompleteAssetStructuringMutation,
} = legalAdviserApi;

import { baseApi } from "@/redux/baseApi";
import { tagTypes } from "@/redux/baseApi/tagTypes";

import {
  ICreateCustodianComplianceRequest,
  ICreateCustodianComplianceResponse,
} from "./interface";

export const custodianComplianceApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    createCustodianCompliance: builder.mutation<
      ICreateCustodianComplianceResponse,
      ICreateCustodianComplianceRequest
    >({
      query: (payload) => ({
        url: "/stakeholder/admin/custodian-compliance",
        method: "POST",
        data: payload,
      }),
      // Ad-hoc requirement assignment feeds into the same review queue as
      // templated requirements (PRD §6.2), so invalidate the shared tag.
      invalidatesTags: [tagTypes.COMPLIANCE_REQUIREMENTS],
    }),
  }),
  overrideExisting: false,
});

export const { useCreateCustodianComplianceMutation } = custodianComplianceApi;

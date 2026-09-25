import { baseApi } from "@/redux/baseApi";
import { tagTypes } from "@/redux/baseApi/tagTypes";
import {
  IAuditRecord,
  IAuditRecordResponse,
  IAuditTrailParams,
} from "./interface";

export const auditTrailApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getAuditTrail: builder.query<
      IPaginatedResponse<IAuditRecord>,
      IAuditTrailParams
    >({
      query: (params) => ({
        url: "/audit-trail",
        method: "GET",
        params,
      }),
      providesTags: [tagTypes.AUDIT_TRAIL],
    }),

    getAuditActivity: builder.query<IAuditRecordResponse, { id: string }>({
      query: ({ id }) => ({
        url: `/audit-trail/${id}`,
        method: "GET",
      }),
      providesTags: [tagTypes.AUDIT_TRAIL],
    }),
  }),
  overrideExisting: false,
});

export const { useGetAuditTrailQuery, useGetAuditActivityQuery } = auditTrailApi;

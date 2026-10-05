import { baseApi } from "@/redux/baseApi";
import { tagTypes } from "@/redux/baseApi/tagTypes";
import { AddFeeExemptUserRequest, FeeExemptUsersResponse } from "./interface";

export const feeExemptionsApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getFeeExemptUsers: builder.query<FeeExemptUsersResponse, void>({
      query: () => ({ url: "/fee/exempt-users", method: "GET" }),
      providesTags: [tagTypes.FEE_EXEMPTIONS],
    }),

    addFeeExemptUser: builder.mutation<unknown, AddFeeExemptUserRequest>({
      query: (data) => ({ url: "/fee/exempt-users", method: "POST", data }),
      invalidatesTags: [tagTypes.FEE_EXEMPTIONS],
    }),

    removeFeeExemptUser: builder.mutation<unknown, string>({
      query: (username) => ({ url: `/fee/exempt-users/${encodeURIComponent(username)}`, method: "DELETE" }),
      invalidatesTags: [tagTypes.FEE_EXEMPTIONS],
    }),
  }),
});

export const { useGetFeeExemptUsersQuery, useAddFeeExemptUserMutation, useRemoveFeeExemptUserMutation } =
  feeExemptionsApi;

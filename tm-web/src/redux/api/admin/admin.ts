import { baseApi } from "@/redux/baseApi";

import {
  IInviteAdmin,
  IAdminResponse,
  ISuspendAdmin,
  IAdminListResponse,
  IRevokeAdminSuspension,
  IUpdateAdminRolePayload,
  IGetSuspensionHistory,
  IViolationOption,
  IAdminListRequest,
} from "./interface";

export const adminApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    inviteAdmin: builder.mutation<IAdminResponse, IInviteAdmin>({
      query: (payload) => ({
        url: "/admin/invite",
        method: "POST",
        data: payload,
      }),
      invalidatesTags: ["Admins"],
    }),
    listAdmins: builder.query<IAdminListResponse, IAdminListRequest>({
      query: (params) => ({
        url: "/admin/list",
        method: "GET",
        params,
      }),
      providesTags: ["Admins"],
    }),
    suspendAdmin: builder.mutation<IAdminResponse, ISuspendAdmin>({
      query: (payload) => ({
        url: "/admin/suspend",
        method: "PATCH",
        data: payload,
      }),
      invalidatesTags: ["Admins"],
    }),

    updateAdminRole: builder.mutation<IAdminResponse, IUpdateAdminRolePayload>({
      query: (payload) => ({
        url: "/admin/modify/status",
        method: "PATCH",
        data: payload,
      }),
      invalidatesTags: ["Admins"],
    }),

    unsuspendAdmin: builder.mutation<IAdminResponse, IRevokeAdminSuspension>({
      query: (payload) => ({
        url: "/admin/unsuspend",
        method: "PATCH",
        data: payload,
      }),
      invalidatesTags: ["Admins"],
    }),
    getSuspensionHistory: builder.query<
      IGetSuspensionHistory[],
      { email: string }
    >({
      query: ({ email }) => ({
        url: `/suspension/history/${email}`, // Use path parameter
        method: "GET",
      }),
    }),

    getVoilations: builder.query<IViolationOption[], void>({
      query: () => ({
        url: "/configurations/all",
        method: "GET",
      }),
    }),
  }),
});

export const {
  useInviteAdminMutation,
  useListAdminsQuery,
  useSuspendAdminMutation,
  useUpdateAdminRoleMutation,
  useUnsuspendAdminMutation,
  useGetSuspensionHistoryQuery,
  useGetVoilationsQuery,
} = adminApi;

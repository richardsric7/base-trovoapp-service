// organizationEndpoints.ts

import { IInviteMember, IMembersResponse } from "../organizations";
import {
  FundReleasesResponse,
  GetFundReleasesParams,
  OrganizationDetailsLogin,
} from "./interface";

import { orgApi } from "@/redux/baseApi/orgApi";
import { tagTypes } from "@/redux/baseApi/tagTypes";
import { IFundReleaseResponse } from "../trustees/interface";
import { FundReleaseExecutionPayload } from "./interface";

export const orgOnlyEndpoints = orgApi.injectEndpoints({
  endpoints: (builder) => ({
    getOrganizationDetailsData: builder.query<OrganizationDetailsLogin, void>({
      query: () => ({
        url: `/organizations/details`,
        method: "GET",
      }),
    }),

    logoutOrganizationMember: builder.mutation<void, void>({
      query: () => ({
        url: "/organizations/member/logout",
        method: "POST",
      }),
    }),

    // get members
    getOrganizationMembers: builder.query<
      IMembersResponse,
      {
        page?: number;
        pageSize?: number;
      }
    >({
      query: (params) => ({
        url: `/organizations/members`,
        method: "GET",
        params,
      }),
    }),
    // invite team member
    inviteMember: builder.mutation<IInviteMember, IInviteMember>({
      query: (payload) => ({
        url: `/organizations/members`,
        method: "POST",
        data: payload,
      }),
    }),

    // get custodian fund release requests
    getCustodianFundRelease: builder.query<IFundReleaseResponse, string>({
      query: (id) => ({ url: `/stakeholder/custodian/fund-releases/${encodeURIComponent(id)}`, method: "GET" }),
      providesTags: [tagTypes.FUND_RELEASES],
    }),
    executeCustodianFundRelease: builder.mutation<IFundReleaseResponse, { requestId: string; payload: FundReleaseExecutionPayload & { challenge_id: string } }>({
      query: ({ requestId, payload }) => ({ url: `/stakeholder/custodian/fund-releases/${encodeURIComponent(requestId)}/execute`, method: "POST", data: payload }),
      invalidatesTags: [tagTypes.FUND_RELEASES],
    }),
    updateCustodianFundReleaseStatus: builder.mutation<IFundReleaseResponse, { requestId: string; payload: FundReleaseExecutionPayload }>({
      query: ({ requestId, payload }) => ({ url: `/stakeholder/custodian/fund-releases/${encodeURIComponent(requestId)}/status`, method: "PUT", data: payload }),
      invalidatesTags: [tagTypes.FUND_RELEASES],
    }),
    getFundReleases: builder.query<
      FundReleasesResponse,
      GetFundReleasesParams | void
    >({
      query: (params) => ({
        url: `/stakeholder/custodian/fund-releases`,
        method: "GET",
        params: params ?? undefined,
      }),
      providesTags: [tagTypes.FUND_RELEASES],
    }),
  }),
  overrideExisting: false,
});

export const {
  useGetOrganizationDetailsDataQuery,
  useLazyGetOrganizationDetailsDataQuery,
  useLogoutOrganizationMemberMutation,
  useInviteMemberMutation,
  useGetOrganizationMembersQuery,
  useGetFundReleasesQuery,
  useGetCustodianFundReleaseQuery,
  useExecuteCustodianFundReleaseMutation,
  useUpdateCustodianFundReleaseStatusMutation,
} = orgOnlyEndpoints;

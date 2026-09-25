import { baseApi } from "@/redux/baseApi";
import {
  GetOrganizationsQueryParams,
  GetOrganizationsResponse,
  IDeactivateOrganizationResponse,
  IInviteMember,
  IInviteOrganizationRequest,
  IInviteOrganizationResponse,
  IMemberLoginPayload,
  IMemberLoginResponse,
  IMembersResponse,
  IResendOtpPayload,
  ISetupRootUserPasswordRequest,
  IValidateEmailWithOtpRequest,
  OrganizationDetailsResponse,
} from "./interface";

export const organizationsList = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    // Validate root user invite
    validateRootUserInvite: builder.mutation<any, { inviteId: string }>({
      query: ({ inviteId }) => ({
        url: `/organizations/invite/validate/${inviteId}`,
        method: "POST",
      }),
    }),

    validateRootUserEmail: builder.mutation<any, IValidateEmailWithOtpRequest>({
      query: ({ inviteId, otp }) => ({
        url: `/organizations/invite/validate-email/${inviteId}`,
        method: "POST",
        data: { otp },
      }),
    }),
    // Setup root user password
    setupRootUserPassword: builder.mutation<any, ISetupRootUserPasswordRequest>(
      {
        query: ({ inviteId, password, otp }) => ({
          url: `/organizations/setup-password/${inviteId}`,
          method: "POST",
          data: { password, otp },
        }),
      },
    ),

    //login
    loginMember: builder.mutation<IMemberLoginResponse, IMemberLoginPayload>({
      query: (payload) => ({
        url: `/organizations/member/login`,
        method: "POST",
        data: payload,
      }),
    }),

    // Resend OTP for team member
    resendTeamMemberOtp: builder.mutation<any, IResendOtpPayload>({
      query: (payload) => ({
        url: `/organizations/team-member/resend-otp`,
        method: "POST",
        data: payload,
      }),
    }),

    // Get organization list
    getOrganizationsList: builder.query<
      GetOrganizationsResponse,
      GetOrganizationsQueryParams
    >({
      query: (params) => ({
        url: "/organizations",
        method: "GET",
        params,
      }),
      // providesTags: ["Organizations"],
    }),

    // organizationDetails
    // getOrganizationDetails: builder.query<OrganizationDetailsResponse, string>({
    //   query: (orgId) => ({
    //     url: "/organizations/details",
    //     method: "GET",
    //     params: { organization_id: orgId },
    //   }),
    //   // providesTags: ["Organizations"],
    // }),

    getOrganizationDetails: builder.query<OrganizationDetailsResponse, string>({
      query: (orgId) => {
        return {
          url: "/organizations/details",
          method: "GET",
          params: { organization_id: orgId },
        };
      },
    }),

    inviteOragnization: builder.mutation<
      IInviteOrganizationResponse,
      IInviteOrganizationRequest
    >({
      query: ({ stakeholder_type, ...body }) => ({
        url: "/organizations/invite/root-user",
        method: "POST",
        data: body,
        params: {
          stakeholder_type,
        },
      }),
      invalidatesTags: ["Organizations"],
    }),
    // inviteOragnization: builder.mutation<
    //   IInviteOrganizationResponse,
    //   IInviteOrganizationRequest
    // >({
    //   query: ({ stakeholder_type, ...body }) => ({
    //     url: "/organizations/invite/root-user",
    //     method: "POST",
    //     data: body,
    //     params: stakeholder_type ? { stakeholder_type } : undefined,
    //   }),
    //   invalidatesTags: ["Organizations"],
    // }),

    // invite team member
    inviteMember: builder.mutation<IInviteMember, IInviteMember>({
      query: (payload) => ({
        url: `/organizations/members`,
        method: "POST",
        data: payload,
      }),
    }),
    // get members
    getOrganizationMembers: builder.query<
      IMembersResponse,
      {
        organization_id: string;
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

    // Reset password
    resetMemberPassword: builder.mutation<
      any,
      { email: string; organization_id: string }
    >({
      query: (data) => ({
        url: "/organizations/members/reset-password",
        method: "POST",
        data,
      }),
    }),

    // update organizationn

    // In @/redux/api/organizations
    updateOrganizationInfo: builder.mutation({
      query: ({ organizationId, stakeholderType, body }) => ({
        url: `/organizations/update/${organizationId}`,
        method: "PUT",
        params: {
          stakeholder_type: stakeholderType, // Query parameter
        },
        data: body, // Request body
      }),
      invalidatesTags: ["Organizations"],
    }),

    deactivateOrganization: builder.mutation<
      IDeactivateOrganizationResponse,
      string
    >({
      query: (organization_id) => ({
        url: `/organizations/deactivate/${organization_id}`,
        method: "PUT",
      }),
      invalidatesTags: ["Organizations"],
    }),
  }),
});

export const {
  useGetOrganizationsListQuery,
  useInviteOragnizationMutation,
  useValidateRootUserInviteMutation,
  useValidateRootUserEmailMutation,
  useInviteMemberMutation,
  useResendTeamMemberOtpMutation,
  useLoginMemberMutation,
  useSetupRootUserPasswordMutation,
  useGetOrganizationDetailsQuery,
  useResetMemberPasswordMutation,
  useGetOrganizationMembersQuery,
  useUpdateOrganizationInfoMutation,
  useDeactivateOrganizationMutation,
} = organizationsList;

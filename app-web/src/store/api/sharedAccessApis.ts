import { Payload } from '../../types/payload';
import { baseApi } from './baseapi';

export const sharedAccessApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getSharedAccessWallets: builder.query({
      query: (payload: Payload) => ({
        url: '/v1/shared-access/users/account',
        method: 'GET',
        data: {
          creds: {
            signer: payload.signer,
            address: payload.address,
            secretKey: payload.secretKey,
          },
        },
      }),
    }),
    getApprovals: builder.query({
      query: (payload: Payload) => ({
        url: `/v1/shared-access/approvals?&excludeUserApproved=${payload.body.excludeUserApproved ?? 1}&limit=${payload.body.limit ?? 50}${payload.body.query ?? ''}&page=${payload.body.page ?? 1}`,
        method: 'GET',
        data: {
          creds: {
            signer: payload.signer,
            address: payload.address,
            secretKey: payload.secretKey,
          },
        },
      }),
    }),
    addSharedAccess: builder.mutation({
      query: (payload: Payload) => ({
        url: '/v1/shared-access/users/account',
        method: 'POST',
        data: {
          payload: payload.body,
          creds: {
            signer: payload.signer,
            address: payload.address,
            secretKey: payload.secretKey,
          },
        },
      }),
    }),
    updateSharedAccess: builder.mutation({
      query: (payload: Payload) => ({
        url: '/v1/shared-access/users/account',
        method: 'PUT',
        data: {
          payload: payload.body,
          creds: {
            signer: payload.signer,
            address: payload.address,
            secretKey: payload.secretKey,
          },
        },
      }),
    }),
    disableSharedAccess: builder.mutation({
      query: (payload: Payload) => ({
        url: '/v1/shared-access/users/account',
        method: 'DELETE',
        data: {
          payload: payload.body,
          creds: {
            signer: payload.signer,
            address: payload.address,
            secretKey: payload.secretKey,
          },
        },
      }),
    }),
    approveSharedAccess: builder.mutation({
      query: (payload: Payload) => ({
        url: `/v1/shared-access/approval/${payload.body.id}`,
        method: 'POST',
        data: {
          payload: payload.body,
          creds: {
            signer: payload.signer,
            address: payload.address,
            secretKey: payload.secretKey,
          },
        },
      }),
    }),
    rejectSharedAccess: builder.mutation({
      query: (payload: Payload) => ({
        url: `/v1/shared-access/approval/${payload.body.id}`,
        method: 'DELETE',
        data: {
          payload: payload.body,
          creds: {
            signer: payload.signer,
            address: payload.address,
            secretKey: payload.secretKey,
          },
        },
      }),
    }),
    getWalletBalance: builder.query({
      query: (payload: Payload) => ({
        url: '/v1/shared-access/wallet-balances',
        method: 'GET',
        data: {
          creds: {
            signer: payload.signer,
            address: payload.address,
            secretKey: payload.secretKey,
          },
        },
      }),
    }),
    checkUsername: builder.query({
      query: (payload: Payload) => ({
        url: `/v1/users/${payload.body.username}`,
        method: 'GET',
        data: {
          creds: {
            signer: payload.signer,
            address: payload.address,
            secretKey: payload.secretKey,
          },
        },
      }),
    }),
  }),
});

export const {
  useGetSharedAccessWalletsQuery,
  useLazyGetSharedAccessWalletsQuery,
  useGetApprovalsQuery,
  useLazyGetApprovalsQuery,
  useAddSharedAccessMutation,
  useUpdateSharedAccessMutation,
  useDisableSharedAccessMutation,
  useApproveSharedAccessMutation,
  useRejectSharedAccessMutation,
  useGetWalletBalanceQuery,
  useLazyGetWalletBalanceQuery,
  useCheckUsernameQuery,
  useLazyCheckUsernameQuery,
} = sharedAccessApi;

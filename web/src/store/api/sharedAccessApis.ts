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
            publicKey: payload.publicKey,
            secretKey: payload.secretKey,
          },
        },
      }),
    }),
    getApprovals: builder.query({
      query: (payload: Payload) => ({
        url: '/v1/shared-access/approvals',
        method: 'GET',
        data: {
          creds: {
            signer: payload.signer,
            publicKey: payload.publicKey,
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
            publicKey: payload.publicKey,
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
            publicKey: payload.publicKey,
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
            publicKey: payload.publicKey,
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
            publicKey: payload.publicKey,
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
            publicKey: payload.publicKey,
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
            publicKey: payload.publicKey,
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
  useApproveSharedAccessMutation,
  useRejectSharedAccessMutation,
  useGetWalletBalanceQuery,
  useLazyGetWalletBalanceQuery,
  useCheckUsernameQuery,
  useLazyCheckUsernameQuery,
} = sharedAccessApi;

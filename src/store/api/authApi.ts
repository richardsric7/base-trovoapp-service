import { Payload } from '../../types/payload';
import { baseApi } from './baseapi';

export const authApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    register: builder.mutation({
      query: (payload: Payload) => ({
        url: '/v1/users',
        method: 'POST',        
        data: {
          payload: payload.body,
          creds: {
            signer: payload.signer,
            publicKey: payload.publicKey,
            secretKey: payload.secretKey,
          }
        }, 
      }),
    }),
    verifyLogin: builder.query({
      query: (payload) => ({
        url: `/v1/users/verify/${payload.username}/${payload.loginId}`,
        method: 'GET',
      }),
    }),
    getUserInfo: builder.query({
      query: (payload) => ({
        url: `/v1/users/detail/${payload}`,
        method: 'GET',
      }),
      providesTags: ['user'],
    }),
    toggleUserAvailability: builder.mutation({
      query: (payload) => ({
        url: '/v1/users/toggle',
        method: 'PUT',
        data: payload,
      }),
      invalidatesTags: ['user'],
    }),
  }),
});

export const {
  useRegisterMutation,
  useGetUserInfoQuery,
  useVerifyLoginQuery,
  useToggleUserAvailabilityMutation,
} = authApi;

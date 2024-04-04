import { baseApi } from './baseapi';

export const authApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    login: builder.mutation({
      query: (payload) => ({
        url: '/v1/login',
        method: 'POST',
        data: payload,
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
  useLoginMutation,
  useGetUserInfoQuery,
  useVerifyLoginQuery,
  useToggleUserAvailabilityMutation,
} = authApi;

import { baseApi } from "@/redux/baseApi";
import { ILogin, ILoginRes, IVerifyLogin, IWalletConnectResponse } from "./interface";

export const authApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    login: builder.mutation<ILoginRes, ILogin>({
      query: (payload) => ({
        url: '/login',
        method: 'POST',
        data: payload,
      })
    }),
    logout: builder.mutation<ISuccessResponse<null>, void>({
      query: () => ({
        url: '/logout',
        method: 'POST',
      })
    }),
    verifyLogin: builder.query<ISuccessResponse<IWalletConnectResponse>, IVerifyLogin>({
      query: (payload) => ({
        url: `/users/verify/${payload.targetUser}/${payload.loginID}`,
        method: 'GET',
      })
    }),
  })
});

export const {
  useLoginMutation,
  useVerifyLoginQuery,
  useLogoutMutation
  
} = authApi;

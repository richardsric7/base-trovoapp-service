import { baseApi } from "@/redux/baseApi";
import { tagTypes } from "@/redux/baseApi/tagTypes";
import {
  IFetchUserResponse,
  ISuspendUsers,
  iUsers,
  IUsersParams,
  PaymentHistoryParams,
  PaymentHistoryResponse,
  WalletBalancesResponse,
} from "./interface";

export const usersApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getUsers: builder.query<IPaginatedResponse<iUsers>, IUsersParams>({
      query: (params) => ({
        url: "/users",
        method: "GET",
        params: params,
      }),
      providesTags: [tagTypes.USERS],
    }),

    fetchUserByCriteria: builder.query<
      IFetchUserResponse,
      { id?: string; email?: string; username?: string }
    >({
      query: ({ id, email, username }) => ({
        url: "/users/profile",
        method: "GET",
        params: { id, email, username },
      }),
    }),
    suspendUser: builder.mutation<iUsers, ISuspendUsers>({
      query: (payload) => ({
        url: "/user/suspend-or-reactivate",
        method: "PATCH",
        data: payload,
      }),
      invalidatesTags: [tagTypes.USERS],
    }),
    getPaymentHistory: builder.query<
      PaymentHistoryResponse,
      PaymentHistoryParams
    >({
      query: (params) => ({
        url: "/payment/history",
        method: "GET",
        params,
      }),
    }),

    getWalletBalances: builder.query<WalletBalancesResponse, string>({
      query: (walletPublicKey) => ({
        url: `/wallet-balances/${walletPublicKey}`,
        method: "GET",
      }),
    }),
  }),
});

export const {
  useGetUsersQuery,
  useFetchUserByCriteriaQuery,
  useSuspendUserMutation,
  useGetPaymentHistoryQuery,
  useGetWalletBalancesQuery,
} = usersApi;

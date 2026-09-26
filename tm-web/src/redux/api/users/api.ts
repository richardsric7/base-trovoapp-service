import { baseApi } from "@/redux/baseApi";
import { tagTypes } from "@/redux/baseApi/tagTypes";
import {
  IFetchUserResponse,
  ISuspendOrLiftUserPayload,
  iUsers,
  IUsersParams,
  IUserSuspensionHistoryEntry,
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
      providesTags: [tagTypes.USERS],
    }),
    suspendUser: builder.mutation<{ message: string }, ISuspendOrLiftUserPayload>({
      query: (payload) => ({
        url: "/admin/users/suspend",
        method: "PATCH",
        data: payload,
      }),
      invalidatesTags: [tagTypes.USERS],
    }),
    liftUserSuspension: builder.mutation<{ message: string }, ISuspendOrLiftUserPayload>({
      query: (payload) => ({
        url: "/admin/users/lift-suspension",
        method: "PATCH",
        data: payload,
      }),
      invalidatesTags: [tagTypes.USERS],
    }),
    getUserSuspensionHistory: builder.query<IUserSuspensionHistoryEntry[], string>({
      query: (email) => ({
        url: `/users/suspension-history/${email}`,
        method: "GET",
      }),
      providesTags: [tagTypes.USERS],
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
      query: (walletAddress) => ({
        url: `/wallet-balances/${walletAddress}`,
        method: "GET",
      }),
    }),
  }),
});

export const {
  useGetUsersQuery,
  useFetchUserByCriteriaQuery,
  useSuspendUserMutation,
  useLiftUserSuspensionMutation,
  useGetUserSuspensionHistoryQuery,
  useGetPaymentHistoryQuery,
  useGetWalletBalancesQuery,
} = usersApi;

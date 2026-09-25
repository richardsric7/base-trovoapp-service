import { baseApi } from "@/redux/baseApi";
import {
  MintingUser,
  MintingUserApiResponse,
  MintingUserPayload,
  MintingUserQueryParams,
  MintingUserSaveResponse,
  SearchMintingUsersParams,
} from "./interface";

export const mintingUserApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getMintingUsers: builder.query<
      MintingUserApiResponse,
      MintingUserQueryParams
    >({
      query: (params) => ({
        url: "/tokenization/minting/users",
        method: "GET",
        params: params,
      }),
      providesTags: ["mintingUsers"],
    }),

    addMintingUser: builder.mutation<
      MintingUserSaveResponse,
      Partial<MintingUser>
    >({
      query: (data) => ({
        url: "/tokenization/minting/users",
        method: "POST",
        data,
      }),
      invalidatesTags: ["mintingUsers"],
    }),

    deleteMintingUser: builder.mutation({
      query: ({ id, role }: { id: number; role: string }) => ({
        url: `/tokenization/minting/users`,
        method: "DELETE",
        data: { id, role },
      }),
      invalidatesTags: ["mintingUsers"],
    }),

    searchMintingUsers: builder.query<
      MintingUserPayload,
      SearchMintingUsersParams
    >({
      query: (params) => ({
        url: `/tokenization/minting/users/search`,
        method: "GET",
        params,
      }),
    }),
  }),
});

export const {
  useGetMintingUsersQuery,
  useAddMintingUserMutation,
  useDeleteMintingUserMutation,
  useSearchMintingUsersQuery,
} = mintingUserApi;

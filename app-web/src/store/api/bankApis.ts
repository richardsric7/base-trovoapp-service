import { baseApi } from './baseapi';
import { tagTypes } from './baseapi/tagTypes';
import { P2PCreds } from './p2pApis';
import { BankWithdrawal, StablerailBank, StablerailProfile, VirtualAccount } from '../../types/bank';

type WithCreds<T> = T & { creds: P2PCreds };

export type BankWithdrawalBody = {
  amount: string;
  accountNumber: string;
  bankCode: string;
  transaction?: string;
  transactionSignature?: string;
};

// bankApi: Naira deposits (onramp) and withdrawals to a Nigerian bank account
// (offramp) through Stablerail. Users are onboarded with Stablerail by the
// backend when KYC level 1 (BVN) completes, so there is no BVN call here.
// creds.address is the wallet deposits go to and withdrawals leave from.
export const bankApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getStablerailProfile: builder.query<StablerailProfile, WithCreds<{}>>({
      query: ({ creds }) => ({ url: '/v1/users/stablerail/profile', method: 'GET', data: { creds } }),
    }),
    getStablerailBanks: builder.query<StablerailBank[], WithCreds<{}>>({
      query: ({ creds }) => ({ url: '/v1/users/stablerail/banks', method: 'GET', data: { creds } }),
    }),
    getBankWithdrawals: builder.query<{ withdrawals: BankWithdrawal[] }, WithCreds<{}>>({
      query: ({ creds }) => ({ url: '/v1/users/stablerail/withdrawals', method: 'GET', data: { creds } }),
      providesTags: [tagTypes.bankWithdrawal],
    }),
    // depositNaira returns the virtual account to pay into
    depositNaira: builder.mutation<{ data: { virtualAccount: VirtualAccount } }, WithCreds<{ amount: string }>>({
      query: ({ creds, amount }) => ({
        url: `/v1/users/stablerail/onrampcngn/${encodeURIComponent(amount)}`,
        method: 'POST',
        data: { payload: {}, creds },
      }),
    }),
    // withdrawToBank: without transactionSignature it returns the transaction
    // to sign (and messages to show); with it, it submits the withdrawal
    withdrawToBank: builder.mutation<
      { transaction: string; messages: string[]; transactionId: string; withdrawalId: string },
      WithCreds<{ body: BankWithdrawalBody }>
    >({
      query: ({ creds, body }) => ({ url: '/v1/users/stablerail/withdraw', method: 'POST', data: { payload: body, creds } }),
      invalidatesTags: (_r, _e, arg) => (arg.body.transactionSignature ? [tagTypes.bankWithdrawal] : []),
    }),
  }),
});

export const {
  useGetStablerailProfileQuery,
  useGetStablerailBanksQuery,
  useGetBankWithdrawalsQuery,
  useDepositNairaMutation,
  useWithdrawToBankMutation,
} = bankApi;

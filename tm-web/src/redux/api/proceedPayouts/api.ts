import { baseApi } from "@/redux/baseApi";
import { tagTypes } from "@/redux/baseApi/tagTypes";
import {
  Envelope,
  IEngineState,
  IFeeConfig,
  IFeesReport,
  IPayout,
  IPayoutDetail,
  IPayoutItem,
  IPayoutsReport,
  ItemFilters,
  PayoutFilters,
} from "./interface";

const clean = (params: Record<string, unknown>) =>
  Object.fromEntries(Object.entries(params).filter(([, v]) => v !== undefined && v !== ""));

type PayoutAction = "prepare" | "approve" | "reject" | "confirm-funding" | "pause" | "resume" | "cancel" | "retry-failed";
type ItemAction = "exclude" | "include" | "mark-paid";

export const proceedPayoutsApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getPayouts: builder.query<Envelope<{ payouts: IPayout[]; total: number }>, PayoutFilters>({
      query: (params) => ({ url: "/proceed-payouts", method: "GET", params: clean({ ...params }) }),
      providesTags: [tagTypes.PROCEED_PAYOUTS],
    }),
    getPayout: builder.query<Envelope<IPayoutDetail>, number>({
      query: (id) => ({ url: `/proceed-payouts/${id}`, method: "GET" }),
      providesTags: [tagTypes.PROCEED_PAYOUTS],
    }),
    getPayoutItems: builder.query<Envelope<{ items: IPayoutItem[]; total: number }>, ItemFilters>({
      query: ({ id, ...params }) => ({ url: `/proceed-payouts/${id}/items`, method: "GET", params: clean({ ...params }) }),
      providesTags: [tagTypes.PROCEED_PAYOUTS],
    }),
    payoutAction: builder.mutation<Envelope<IPayout>, { id: number; action: PayoutAction; reason?: string }>({
      query: ({ id, action, reason }) => ({
        url: `/proceed-payouts/${id}/${action}`,
        method: "POST",
        data: reason !== undefined ? { reason } : undefined,
      }),
      invalidatesTags: [tagTypes.PROCEED_PAYOUTS],
    }),
    setPayoutFee: builder.mutation<Envelope<IPayout>, { id: number; feeType: string; feeValue: string; feeCap?: string }>({
      query: ({ id, ...data }) => ({ url: `/proceed-payouts/${id}/fee`, method: "PUT", data }),
      invalidatesTags: [tagTypes.PROCEED_PAYOUTS],
    }),
    payoutItemAction: builder.mutation<
      Envelope<IPayoutItem>,
      { id: number; itemId: string; action: ItemAction; reason?: string; reference?: string }
    >({
      query: ({ id, itemId, action, reason, reference }) => ({
        url: `/proceed-payouts/${id}/items/${itemId}/${action}`,
        method: "POST",
        data: action === "mark-paid" ? { reference } : reason !== undefined ? { reason } : undefined,
      }),
      invalidatesTags: [tagTypes.PROCEED_PAYOUTS],
    }),
    getPayoutEngine: builder.query<Envelope<IEngineState>, void>({
      query: () => ({ url: "/proceed-payouts/engine", method: "GET" }),
      providesTags: [tagTypes.PROCEED_PAYOUTS],
    }),
    payoutEngineAction: builder.mutation<Envelope<IEngineState>, { action: "halt" | "unhalt" | "sweep"; reason?: string; token?: string }>({
      query: ({ action, reason, token }) => ({
        url: `/proceed-payouts/engine/${action}`,
        method: "POST",
        data: action === "sweep" ? { token } : action === "halt" ? { reason } : undefined,
      }),
      invalidatesTags: [tagTypes.PROCEED_PAYOUTS],
    }),
    getPayoutFeeConfig: builder.query<Envelope<IFeeConfig>, void>({
      query: () => ({ url: "/proceed-payouts/fee-config", method: "GET" }),
      providesTags: [tagTypes.PROCEED_PAYOUTS],
    }),
    setPayoutFeeConfig: builder.mutation<Envelope<IFeeConfig>, { feeWallet: string; feeType: string; feeValue: string; feeCap?: string }>({
      query: (data) => ({ url: "/proceed-payouts/fee-config", method: "PUT", data }),
      invalidatesTags: [tagTypes.PROCEED_PAYOUTS],
    }),
    getPayoutsReport: builder.query<Envelope<IPayoutsReport>, PayoutFilters>({
      query: (params) => ({ url: "/proceed-payouts/reports/payouts", method: "GET", params: clean({ ...params }) }),
      providesTags: [tagTypes.PROCEED_PAYOUTS],
    }),
    getPayoutFeesReport: builder.query<Envelope<IFeesReport>, PayoutFilters>({
      query: (params) => ({ url: "/proceed-payouts/reports/fees", method: "GET", params: clean({ ...params }) }),
      providesTags: [tagTypes.PROCEED_PAYOUTS],
    }),
  }),
});

export const {
  useGetPayoutsQuery,
  useGetPayoutQuery,
  useGetPayoutItemsQuery,
  usePayoutActionMutation,
  useSetPayoutFeeMutation,
  usePayoutItemActionMutation,
  useGetPayoutEngineQuery,
  usePayoutEngineActionMutation,
  useGetPayoutFeeConfigQuery,
  useSetPayoutFeeConfigMutation,
  useGetPayoutsReportQuery,
  useGetPayoutFeesReportQuery,
} = proceedPayoutsApi;

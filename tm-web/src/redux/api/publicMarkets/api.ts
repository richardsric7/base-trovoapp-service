import { baseApi } from "@/redux/baseApi";
import { tagTypes } from "@/redux/baseApi/tagTypes";
import {
  AssetInput,
  Envelope,
  ExchangeInput,
  IActionDetail,
  IAssetDetail,
  IAssetRow,
  IBatchView,
  IConfirmationRow,
  ICorporateAction,
  ICustodianRow,
  IDealingMember,
  IExchangeDetail,
  IExchangeRow,
  IExecution,
  IHealth,
  IInstruction,
  IJobRequest,
  IOrderDetail,
  IOverview,
  IPartnerEvent,
  IPartnerWallet,
  IPMAsset,
  IPMOrder,
  IPMPermissions,
  IPMSettings,
  IPriceRow,
  IPriceSnapshot,
  IReconciliation,
  IReconRun,
  IWalletStats,
  Paged,
  PartnerInput,
} from "./interface";

const clean = (params: Record<string, unknown>) =>
  Object.fromEntries(Object.entries(params).filter(([, v]) => v !== undefined && v !== ""));

const T = tagTypes.PUBLIC_MARKETS;
const BASE = "/public-markets";

type AssetAction = "go-live" | "halt" | "resume";
type BatchAction = "approve" | "reject" | "roll";
type InstructionAction = "retry" | "handled";
type ExchangeAction = "rotate-secret" | "suspend" | "activate" | "replay-dead-letters";

export const publicMarketsApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getPMPermissions: builder.query<Envelope<IPMPermissions>, void>({
      query: () => ({ url: `${BASE}/me`, method: "GET" }),
      providesTags: [T],
    }),
    getPMOverview: builder.query<Envelope<IOverview>, void>({
      query: () => ({ url: `${BASE}/overview`, method: "GET" }),
      providesTags: [T],
    }),
    getPMHealth: builder.query<Envelope<IHealth>, void>({
      query: () => ({ url: `${BASE}/health`, method: "GET" }),
      providesTags: [T],
    }),
    getPMPartnerEvents: builder.query<Envelope<{ events: IPartnerEvent[]; total: number }>, { source?: string; kind?: string } & Paged>({
      query: (params) => ({ url: `${BASE}/partner-events`, method: "GET", params: clean({ ...params }) }),
      providesTags: [T],
    }),

    // assets
    getPMAssets: builder.query<Envelope<IAssetRow[]>, { market?: string; type?: string; status?: string; search?: string; custodianId?: number }>({
      query: (params) => ({ url: `${BASE}/assets`, method: "GET", params: clean({ ...params }) }),
      providesTags: [T],
    }),
    getPMAsset: builder.query<Envelope<IAssetDetail>, string>({
      query: (id) => ({ url: `${BASE}/assets/${encodeURIComponent(id)}`, method: "GET" }),
      providesTags: [T],
    }),
    createPMAsset: builder.mutation<Envelope<IPMAsset>, AssetInput>({
      query: (data) => ({ url: `${BASE}/assets`, method: "POST", data }),
      invalidatesTags: [T],
    }),
    updatePMAsset: builder.mutation<Envelope<IPMAsset>, { id: string } & AssetInput>({
      query: ({ id, ...data }) => ({ url: `${BASE}/assets/${encodeURIComponent(id)}`, method: "PUT", data }),
      invalidatesTags: [T],
    }),
    registerPMContract: builder.mutation<Envelope<IPMAsset>, { id: string; contractAddress: string; issuingSafeAddress: string }>({
      query: ({ id, ...data }) => ({ url: `${BASE}/assets/${encodeURIComponent(id)}/contract`, method: "POST", data }),
      invalidatesTags: [T],
    }),
    pmAssetAction: builder.mutation<Envelope<IPMAsset | IJobRequest>, { id: string; action: AssetAction; reason?: string }>({
      query: ({ id, action, reason }) => ({
        url: `${BASE}/assets/${encodeURIComponent(id)}/${action}`,
        method: "POST",
        data: reason !== undefined ? { reason } : undefined,
      }),
      invalidatesTags: [T],
    }),
    setPMPrice: builder.mutation<Envelope<IPMAsset>, { id: string; price: string }>({
      query: ({ id, price }) => ({ url: `${BASE}/assets/${encodeURIComponent(id)}/price`, method: "POST", data: { price } }),
      invalidatesTags: [T],
    }),
    recordPMPosition: builder.mutation<Envelope<unknown>, { id: string; unitsHeld: string; asOf: string; reference: string }>({
      query: ({ id, ...data }) => ({ url: `${BASE}/assets/${encodeURIComponent(id)}/position`, method: "POST", data }),
      invalidatesTags: [T],
    }),
    getPMPriceHistory: builder.query<Envelope<IPriceSnapshot[]>, { id: string; limit?: number }>({
      query: ({ id, limit }) => ({ url: `${BASE}/assets/${encodeURIComponent(id)}/prices`, method: "GET", params: clean({ limit }) }),
      providesTags: [T],
    }),
    getPMPrices: builder.query<Envelope<{ prices: IPriceRow[]; executions: IExecution[] }>, void>({
      query: () => ({ url: `${BASE}/prices`, method: "GET" }),
      providesTags: [T],
    }),

    // orders, batches, instructions
    getPMOrders: builder.query<
      Envelope<{ orders: IPMOrder[]; total: number }>,
      { type?: string; state?: string; channel?: string; path?: string; asset?: string; batch?: string; exchange?: string; search?: string; from?: string; to?: string } & Paged
    >({
      query: (params) => ({ url: `${BASE}/orders`, method: "GET", params: clean({ ...params }) }),
      providesTags: [T],
    }),
    getPMOrder: builder.query<Envelope<IOrderDetail>, string>({
      query: (id) => ({ url: `${BASE}/orders/${encodeURIComponent(id)}`, method: "GET" }),
      providesTags: [T],
    }),
    getPMBatches: builder.query<Envelope<{ batches: IBatchView[]; total: number }>, { status?: string; asset?: string } & Paged>({
      query: (params) => ({ url: `${BASE}/batches`, method: "GET", params: clean({ ...params }) }),
      providesTags: [T],
    }),
    pmBatchAction: builder.mutation<Envelope<IBatchView>, { id: string; action: BatchAction; reason?: string }>({
      query: ({ id, action, reason }) => ({
        url: `${BASE}/batches/${encodeURIComponent(id)}/${action}`,
        method: "POST",
        data: reason !== undefined ? { reason } : undefined,
      }),
      invalidatesTags: [T],
    }),
    getPMInstructions: builder.query<Envelope<{ instructions: IInstruction[]; total: number }>, { status?: string; asset?: string; batch?: string } & Paged>({
      query: (params) => ({ url: `${BASE}/instructions`, method: "GET", params: clean({ ...params }) }),
      providesTags: [T],
    }),
    pmInstructionAction: builder.mutation<Envelope<IInstruction>, { id: string; action: InstructionAction; reason?: string }>({
      query: ({ id, action, reason }) => ({
        url: `${BASE}/instructions/${encodeURIComponent(id)}/${action}`,
        method: "POST",
        data: reason !== undefined ? { reason } : undefined,
      }),
      invalidatesTags: [T],
    }),
    recordPMOutcome: builder.mutation<
      Envelope<IPartnerEvent>,
      { id: string; status: string; executedQuantity?: string; executedPrice?: string; settledQuantity?: string; settlementDate?: string; custodianReference?: string }
    >({
      query: ({ id, ...data }) => ({ url: `${BASE}/instructions/${encodeURIComponent(id)}/outcome`, method: "POST", data }),
      invalidatesTags: [T],
    }),

    // reconciliation & jobs
    getPMReconciliation: builder.query<Envelope<IReconciliation>, void>({
      query: () => ({ url: `${BASE}/reconciliation`, method: "GET" }),
      providesTags: [T],
    }),
    getPMReconciliationRuns: builder.query<Envelope<{ runs: IReconRun[]; total: number }>, { asset?: string } & Paged>({
      query: (params) => ({ url: `${BASE}/reconciliation/runs`, method: "GET", params: clean({ ...params }) }),
      providesTags: [T],
    }),
    runPMReconciliation: builder.mutation<Envelope<IJobRequest>, { target?: string }>({
      query: (data) => ({ url: `${BASE}/reconciliation/run`, method: "POST", data }),
      invalidatesTags: [T],
    }),
    requestPMPositionFeed: builder.mutation<Envelope<IJobRequest>, { target: string }>({
      query: (data) => ({ url: `${BASE}/position-feed`, method: "POST", data }),
      invalidatesTags: [T],
    }),
    getPMJobs: builder.query<Envelope<IJobRequest[]>, void>({
      query: () => ({ url: `${BASE}/jobs`, method: "GET" }),
      providesTags: [T],
    }),

    // corporate actions
    getPMActions: builder.query<Envelope<{ actions: ICorporateAction[]; total: number }>, { status?: string; asset?: string } & Paged>({
      query: (params) => ({ url: `${BASE}/corporate-actions`, method: "GET", params: clean({ ...params }) }),
      providesTags: [T],
    }),
    getPMAction: builder.query<Envelope<IActionDetail>, { id: string; status?: string } & Paged>({
      query: ({ id, ...params }) => ({ url: `${BASE}/corporate-actions/${encodeURIComponent(id)}`, method: "GET", params: clean({ ...params }) }),
      providesTags: [T],
    }),
    declarePMAction: builder.mutation<
      Envelope<IPartnerEvent>,
      { assetCode: string; eventType: string; recordDate: string; payDate?: string; amountPerUnit?: string; description?: string; reference?: string }
    >({
      query: (data) => ({ url: `${BASE}/corporate-actions`, method: "POST", data }),
      invalidatesTags: [T],
    }),
    pmActionAction: builder.mutation<Envelope<IActionDetail | ICorporateAction>, { id: string; action: "approve" | "cancel"; reason?: string }>({
      query: ({ id, action, reason }) => ({
        url: `${BASE}/corporate-actions/${encodeURIComponent(id)}/${action}`,
        method: "POST",
        data: reason !== undefined ? { reason } : undefined,
      }),
      invalidatesTags: [T],
    }),
    getPMConfirmations: builder.query<Envelope<IConfirmationRow[]>, void>({
      query: () => ({ url: `${BASE}/confirmations`, method: "GET" }),
      providesTags: [T],
    }),

    // exchanges & wallets
    getPMExchanges: builder.query<
      Envelope<{ exchanges: IExchangeRow[]; candidates: { id: string; shortName: string; longName: string }[] }>,
      { status?: string; search?: string }
    >({
      query: (params) => ({ url: `${BASE}/exchanges`, method: "GET", params: clean({ ...params }) }),
      providesTags: [T],
    }),
    getPMExchange: builder.query<Envelope<IExchangeDetail>, { id: string; deliveries?: string }>({
      query: ({ id, deliveries }) => ({ url: `${BASE}/exchanges/${encodeURIComponent(id)}`, method: "GET", params: clean({ deliveries }) }),
      providesTags: [T],
    }),
    onboardPMExchange: builder.mutation<Envelope<{ exchange: IExchangeRow; signingSecret: string }>, ExchangeInput>({
      query: (data) => ({ url: `${BASE}/exchanges`, method: "POST", data }),
      invalidatesTags: [T],
    }),
    updatePMExchange: builder.mutation<Envelope<IExchangeRow>, { id: string } & ExchangeInput>({
      query: ({ id, ...data }) => ({ url: `${BASE}/exchanges/${encodeURIComponent(id)}`, method: "PUT", data }),
      invalidatesTags: [T],
    }),
    pmExchangeAction: builder.mutation<
      Envelope<{ signingSecret?: string; previousValidUntil?: string; replayed?: number } & Partial<IExchangeRow>>,
      { id: string; action: ExchangeAction }
    >({
      query: ({ id, action }) => ({ url: `${BASE}/exchanges/${encodeURIComponent(id)}/${action}`, method: "POST" }),
      invalidatesTags: [T],
    }),
    requestPMWithdrawal: builder.mutation<Envelope<IJobRequest>, { id: string; amount: string }>({
      query: ({ id, amount }) => ({ url: `${BASE}/exchanges/${encodeURIComponent(id)}/withdrawals`, method: "POST", data: { amount } }),
      invalidatesTags: [T],
    }),
    replayPMWebhook: builder.mutation<Envelope<unknown>, string>({
      query: (id) => ({ url: `${BASE}/webhooks/${encodeURIComponent(id)}/replay`, method: "POST" }),
      invalidatesTags: [T],
    }),
    getPMWallets: builder.query<
      Envelope<{ wallets: IPartnerWallet[]; total: number; stats: IWalletStats; personalDataVisible: boolean }>,
      { exchange?: string; status?: string; search?: string } & Paged
    >({
      query: (params) => ({ url: `${BASE}/wallets`, method: "GET", params: clean({ ...params }) }),
      providesTags: [T],
    }),

    // settings & partners
    getPMSettings: builder.query<Envelope<IPMSettings>, void>({
      query: () => ({ url: `${BASE}/settings`, method: "GET" }),
      providesTags: [T],
    }),
    updatePMSettings: builder.mutation<Envelope<IPMSettings>, Partial<IPMSettings>>({
      query: (data) => ({ url: `${BASE}/settings`, method: "PUT", data }),
      invalidatesTags: [T],
    }),
    getPMDealingMembers: builder.query<Envelope<IDealingMember[]>, void>({
      query: () => ({ url: `${BASE}/dealing-members`, method: "GET" }),
      providesTags: [T],
    }),
    savePMDealingMember: builder.mutation<Envelope<IDealingMember>, { id?: number } & PartnerInput>({
      query: ({ id, ...data }) => ({ url: id ? `${BASE}/dealing-members/${id}` : `${BASE}/dealing-members`, method: id ? "PUT" : "POST", data }),
      invalidatesTags: [T],
    }),
    getPMCustodians: builder.query<Envelope<ICustodianRow[]>, void>({
      query: () => ({ url: `${BASE}/custodians`, method: "GET" }),
      providesTags: [T],
    }),
    configurePMCustodian: builder.mutation<Envelope<unknown>, { id: number } & PartnerInput>({
      query: ({ id, ...data }) => ({ url: `${BASE}/custodians/${id}`, method: "PUT", data }),
      invalidatesTags: [T],
    }),
  }),
});

export const {
  useGetPMPermissionsQuery,
  useGetPMOverviewQuery,
  useGetPMHealthQuery,
  useGetPMPartnerEventsQuery,
  useGetPMAssetsQuery,
  useGetPMAssetQuery,
  useCreatePMAssetMutation,
  useUpdatePMAssetMutation,
  useRegisterPMContractMutation,
  usePmAssetActionMutation,
  useSetPMPriceMutation,
  useRecordPMPositionMutation,
  useGetPMPriceHistoryQuery,
  useGetPMPricesQuery,
  useGetPMOrdersQuery,
  useGetPMOrderQuery,
  useGetPMBatchesQuery,
  usePmBatchActionMutation,
  useGetPMInstructionsQuery,
  usePmInstructionActionMutation,
  useRecordPMOutcomeMutation,
  useGetPMReconciliationQuery,
  useGetPMReconciliationRunsQuery,
  useRunPMReconciliationMutation,
  useRequestPMPositionFeedMutation,
  useGetPMJobsQuery,
  useGetPMActionsQuery,
  useGetPMActionQuery,
  useDeclarePMActionMutation,
  usePmActionActionMutation,
  useGetPMConfirmationsQuery,
  useGetPMExchangesQuery,
  useGetPMExchangeQuery,
  useOnboardPMExchangeMutation,
  useUpdatePMExchangeMutation,
  usePmExchangeActionMutation,
  useRequestPMWithdrawalMutation,
  useReplayPMWebhookMutation,
  useGetPMWalletsQuery,
  useGetPMSettingsQuery,
  useUpdatePMSettingsMutation,
  useGetPMDealingMembersQuery,
  useSavePMDealingMemberMutation,
  useGetPMCustodiansQuery,
  useConfigurePMCustodianMutation,
} = publicMarketsApi;

import { baseApi } from './baseapi';
import { tagTypes } from './baseapi/tagTypes';
import { P2POffer, P2POrder, P2PDispute, P2POrderFeeQuote, P2PPaymentMethod } from '../../types/p2p';

export type P2PCreds = {
  signer: string;
  address: string;
  secretKey: string;
};

type WithCreds<T> = T & { creds: P2PCreds };

export type CreateOfferBody = {
  offerType: 'BUY' | 'SELL';
  asset: string;
  paymentMethod: P2PPaymentMethod;
  country: string;
  countryCode: string;
  currency: string;
  priceType: string;
  price: string;
  priceMargin?: string;
  minOrderAmount: string;
  maxOrderAmount: string;
  availableLiquidity?: string;
  remark?: string;
};

export const p2pApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    // ---- Offers / marketplace ----
    listP2PMarketplaceOffers: builder.query<
      { data: P2POffer[]; total: number },
      WithCreds<{ offerType?: string; asset?: string; countryCode?: string; currency?: string; page?: number; pageSize?: number }>
    >({
      query: ({ creds, ...filters }) => {
        const qp = new URLSearchParams();
        Object.entries(filters).forEach(([k, v]) => {
          if (v !== undefined && v !== '') qp.set(k, String(v));
        });
        return {
          url: `/v1/p2p/offers?${qp.toString()}`,
          method: 'GET',
          data: { creds },
        };
      },
      providesTags: [tagTypes.p2pOffer],
    }),
    getP2POffer: builder.query<P2POffer, WithCreds<{ offerId: string }>>({
      query: ({ creds, offerId }) => ({
        url: `/v1/p2p/offers/${offerId}`,
        method: 'GET',
        data: { creds },
      }),
      providesTags: [tagTypes.p2pOffer],
    }),
    quoteP2POrderFees: builder.query<P2POrderFeeQuote, WithCreds<{ offerId: string; amount: string }>>({
      query: ({ creds, offerId, amount }) => ({
        url: `/v1/p2p/offers/${offerId}/quote?amount=${encodeURIComponent(amount)}`,
        method: 'GET',
        data: { creds },
      }),
    }),
    createP2POffer: builder.mutation<P2POffer, WithCreds<{ body: CreateOfferBody }>>({
      query: ({ creds, body }) => ({
        url: '/v1/p2p/offers',
        method: 'POST',
        data: { payload: body, creds },
      }),
      invalidatesTags: [tagTypes.p2pOffer],
    }),
    activateP2POffer: builder.mutation<P2POffer, WithCreds<{ offerId: string }>>({
      query: ({ creds, offerId }) => ({
        url: `/v1/p2p/offers/${offerId}/activate`,
        method: 'POST',
        data: { creds },
      }),
      invalidatesTags: [tagTypes.p2pOffer],
    }),
    pauseP2POffer: builder.mutation<P2POffer, WithCreds<{ offerId: string }>>({
      query: ({ creds, offerId }) => ({
        url: `/v1/p2p/offers/${offerId}/pause`,
        method: 'POST',
        data: { creds },
      }),
      invalidatesTags: [tagTypes.p2pOffer],
    }),
    listMyP2POffers: builder.query<{ data: P2POffer[] }, WithCreds<{}>>({
      query: ({ creds }) => ({
        url: '/v1/p2p/my-offers',
        method: 'GET',
        data: { creds },
      }),
      providesTags: [tagTypes.p2pOffer],
    }),

    // ---- Orders ----
    createP2POrder: builder.mutation<P2POrder, WithCreds<{ offerId: string; specifiedAssetAmount: string }>>({
      query: ({ creds, offerId, specifiedAssetAmount }) => ({
        url: '/v1/p2p/orders',
        method: 'POST',
        data: { payload: { offerId, specifiedAssetAmount }, creds },
      }),
      invalidatesTags: [tagTypes.p2pOrder, tagTypes.p2pOffer],
    }),
    listMyP2POrders: builder.query<
      { data: P2POrder[]; total: number },
      WithCreds<{ role?: string; status?: string; page?: number; pageSize?: number }>
    >({
      query: ({ creds, ...filters }) => {
        const qp = new URLSearchParams();
        Object.entries(filters).forEach(([k, v]) => {
          if (v !== undefined && v !== '') qp.set(k, String(v));
        });
        return {
          url: `/v1/p2p/orders?${qp.toString()}`,
          method: 'GET',
          data: { creds },
        };
      },
      providesTags: [tagTypes.p2pOrder],
    }),
    getP2POrder: builder.query<P2POrder, WithCreds<{ orderId: string }>>({
      query: ({ creds, orderId }) => ({
        url: `/v1/p2p/orders/${orderId}`,
        method: 'GET',
        data: { creds },
      }),
      providesTags: [tagTypes.p2pOrder],
    }),
    acceptP2POrder: builder.mutation<P2POrder, WithCreds<{ orderId: string }>>({
      query: ({ creds, orderId }) => ({
        url: `/v1/p2p/orders/${orderId}/accept`,
        method: 'POST',
        data: { creds },
      }),
      invalidatesTags: [tagTypes.p2pOrder],
    }),
    rejectP2POrder: builder.mutation<P2POrder, WithCreds<{ orderId: string }>>({
      query: ({ creds, orderId }) => ({
        url: `/v1/p2p/orders/${orderId}/reject`,
        method: 'POST',
        data: { creds },
      }),
      invalidatesTags: [tagTypes.p2pOrder],
    }),
    cancelP2POrder: builder.mutation<P2POrder, WithCreds<{ orderId: string }>>({
      query: ({ creds, orderId }) => ({
        url: `/v1/p2p/orders/${orderId}/cancel`,
        method: 'POST',
        data: { creds },
      }),
      invalidatesTags: [tagTypes.p2pOrder],
    }),

    // ---- Escrow deposit (two-phase build -> sign -> commit, same
    // contract as /v1/users/payment) ----
    escrowDepositBuild: builder.mutation<any, WithCreds<{ orderId: string }>>({
      query: ({ creds, orderId }) => ({
        url: `/v1/p2p/orders/${orderId}/escrow-deposit`,
        method: 'POST',
        data: { payload: {}, creds },
      }),
    }),
    escrowDepositCommit: builder.mutation<
      any,
      WithCreds<{ orderId: string; transaction: string; transactionSignature: string }>
    >({
      query: ({ creds, orderId, transaction, transactionSignature }) => ({
        url: `/v1/p2p/orders/${orderId}/escrow-deposit`,
        method: 'POST',
        data: { payload: { transaction, transactionSignature, commit: 1 }, creds },
      }),
      invalidatesTags: [tagTypes.p2pOrder],
    }),

    // ---- Fiat payment stage ----
    markP2PPaymentSent: builder.mutation<P2POrder, WithCreds<{ orderId: string }>>({
      query: ({ creds, orderId }) => ({
        url: `/v1/p2p/orders/${orderId}/payment-sent`,
        method: 'POST',
        data: { creds },
      }),
      invalidatesTags: [tagTypes.p2pOrder],
    }),
    confirmP2PPaymentReceived: builder.mutation<P2POrder, WithCreds<{ orderId: string }>>({
      query: ({ creds, orderId }) => ({
        url: `/v1/p2p/orders/${orderId}/payment-confirmed`,
        method: 'POST',
        data: { creds },
      }),
      invalidatesTags: [tagTypes.p2pOrder],
    }),

    // ---- Disputes ----
    getOpenP2PDisputeForOrder: builder.query<P2PDispute, WithCreds<{ orderId: string }>>({
      query: ({ creds, orderId }) => ({
        url: `/v1/p2p/orders/${orderId}/dispute`,
        method: 'GET',
        data: { creds },
      }),
    }),
    openP2PDispute: builder.mutation<
      P2PDispute,
      WithCreds<{ orderId: string; subject: string; description: string; evidence?: string[] }>
    >({
      query: ({ creds, orderId, subject, description, evidence }) => ({
        url: `/v1/p2p/orders/${orderId}/disputes`,
        method: 'POST',
        data: { payload: { subject, description, evidence: evidence ?? [] }, creds },
      }),
      invalidatesTags: [tagTypes.p2pOrder],
    }),
    merchantConfirmsP2PPayment: builder.mutation<P2POrder, WithCreds<{ disputeId: string }>>({
      query: ({ creds, disputeId }) => ({
        url: `/v1/p2p/disputes/${disputeId}/merchant-confirms-payment`,
        method: 'POST',
        data: { creds },
      }),
      invalidatesTags: [tagTypes.p2pOrder],
    }),
    buyerConfirmsP2PNotPaid: builder.mutation<P2POrder, WithCreds<{ disputeId: string }>>({
      query: ({ creds, disputeId }) => ({
        url: `/v1/p2p/disputes/${disputeId}/buyer-confirms-not-paid`,
        method: 'POST',
        data: { creds },
      }),
      invalidatesTags: [tagTypes.p2pOrder],
    }),
  }),
});

export const {
  useListP2PMarketplaceOffersQuery,
  useGetP2POfferQuery,
  useQuoteP2POrderFeesQuery,
  useCreateP2POfferMutation,
  useActivateP2POfferMutation,
  usePauseP2POfferMutation,
  useListMyP2POffersQuery,
  useCreateP2POrderMutation,
  useListMyP2POrdersQuery,
  useGetP2POrderQuery,
  useAcceptP2POrderMutation,
  useRejectP2POrderMutation,
  useCancelP2POrderMutation,
  useEscrowDepositBuildMutation,
  useEscrowDepositCommitMutation,
  useMarkP2PPaymentSentMutation,
  useConfirmP2PPaymentReceivedMutation,
  useGetOpenP2PDisputeForOrderQuery,
  useOpenP2PDisputeMutation,
  useMerchantConfirmsP2PPaymentMutation,
  useBuyerConfirmsP2PNotPaidMutation,
} = p2pApi;

import { Payload } from '../../types/payload';
import { baseApi } from './baseapi';

export const authApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({    
    sendAsset: builder.mutation({
      query: (payload: Payload) => ({        
        url: payload.body.isSharedWallet ? '/v1/shared-access/payment' : '/v1/users/payment',
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
    receiveAsset: builder.query({
      query: (payload: Payload) => {
        return ({
          url: `/v1/users/payment/generate/${payload.body.alias}?paymentDestination=${payload.body.publicKey}&assetCode=${payload.body.assetCode}&assetIssuer=${payload.body.assetIssuer}&amount=${payload.body.amount}&memo=${payload.body.memo != null ? encodeURIComponent(payload.body.memo.toString()) : ''}`,
          method: 'GET',        
          data: {
            creds: {
              signer: payload.signer,
              publicKey: payload.publicKey,
              secretKey: payload.secretKey,
            }
          }, 
        })
      },
    }),
    swapAsset: builder.mutation({
      query: (payload: Payload) => ({
        url: payload.body.isSharedWallet ? '/v1/users/swap' : '/v1/users/swap',
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
    fetchFiatPayments: builder.query({
      query: (payload: Payload) => ({
        url: `/v1/users/payments/${payload.publicKey}?limit=${payload.body.limit ?? 50}${payload.body.query ?? ''}`,
        method: 'GET',
        data: {
          creds: {
            signer: payload.signer,
            publicKey: payload.publicKey,
            secretKey: payload.secretKey,
          }
        },
      }),
    }),
    fetchFiatAmountForActivation: builder.query({
      query: (payload: Payload) => ({
        url: '/v1/users/activate/fiat',
        method: 'GET',
        data: {
          creds: {
            signer: payload.signer,
            publicKey: payload.publicKey,
            secretKey: payload.secretKey,
          }
        },
      }),
    }),
    fetchTokenizedAssets: builder.query({
      query: (payload: Payload) => ({
        url: `/v1/tokenization/list?onlyWithUserPermission=0&salesList=${payload.body.status}&${payload.body.filters}`,
        method: 'GET',
        data: {
          creds: {
            signer: payload.signer,
            publicKey: payload.publicKey,
            secretKey: payload.secretKey,
          }
        },
      }),
    }),
    fetchExpressedInterests: builder.query({
      query: (payload: Payload) => ({
        url: `/v1/tokenization/expressed-interests`,
        method: 'GET',
        data: {
          creds: {
            signer: payload.signer,
            publicKey: payload.publicKey,
            secretKey: payload.secretKey,
          }
        },
      }),
    }),
    fetchSubscriptions: builder.query({
      query: (payload: Payload) => ({
        url: `/v1/tokenization/subscriptions`,
        method: 'GET',
        data: {
          creds: {
            signer: payload.signer,
            publicKey: payload.publicKey,
            secretKey: payload.secretKey,
          }
        },
      }),
    }),
    fetchTokenizationData: builder.query({
      query: (payload: Payload) => ({
        url: `/v1/tokenization`,
        method: 'GET',
        data: {
          creds: {
            signer: payload.signer,
            publicKey: payload.publicKey,
            secretKey: payload.secretKey,
          }
        },
      }),
    }),
    buyTokenizedAssets: builder.mutation({      
      query: (payload: Payload) => {
        console.log('from mutation ==> ', payload);
        return ({
        url: payload.body.isSharedWallet
            ? `/v1/shared-access/tokenization/subscriptions/${payload.body.assetId}`
            : `/v1/tokenization/subscriptions/${payload.body.assetId}`,
        method: 'POST',
        data: {
          payload: payload.body.data,
          creds: {
            signer: payload.signer,
            publicKey: payload.publicKey,
            secretKey: payload.secretKey,
          }
        },
      })
      },
    }),
    subscribeTokenizedAssets: builder.mutation({
      query: (payload: Payload) => ({
        url:
             `/v1/tokenization/expressed-interests/${payload.body.assetId}`,
        method: 'POST',
        data: {
          payload: {'amount': payload.body.amount},
          creds: {
            signer: payload.signer,
            publicKey: payload.publicKey,
            secretKey: payload.secretKey,
          }
        },
      }),
    }),
    addSubwallet: builder.mutation({
      query: (payload: Payload) => ({
        url:
             `/v1/users/subwallet`,
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
    addAsset: builder.mutation({
      query: (payload: Payload) => ({        
        url: payload.body.isSharedWallet
            ? '/v1/shared-access/users/asset/opt-in'
            : '/v1/users/asset/opt-in',
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
    removeAsset: builder.mutation({
      query: (payload: Payload) => ({        
        url: payload.body.isSharedWallet
            ? '/v1/shared-access/users/asset/opt-out'
            : '/v1/users/asset/opt-out',
        method: 'DELETE',
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
  }),
});

export const {
  useSendAssetMutation,
  useAddSubwalletMutation,
  useSwapAssetMutation,
  useBuyTokenizedAssetsMutation,
  useAddAssetMutation,
  useRemoveAssetMutation,
  useSubscribeTokenizedAssetsMutation,
  useLazyReceiveAssetQuery,
  useFetchFiatPaymentsQuery,
  useFetchTokenizationDataQuery,
  useFetchTokenizedAssetsQuery,
  useFetchExpressedInterestsQuery,
  useFetchSubscriptionsQuery,
  useLazyFetchFiatAmountForActivationQuery,
} = authApi;

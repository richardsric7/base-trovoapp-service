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
  }),
});

export const {
  useSendAssetMutation,
  useSwapAssetMutation,
  useLazyReceiveAssetQuery,
} = authApi;

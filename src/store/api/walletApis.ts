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
  }),
});

export const {
  useSendAssetMutation,
} = authApi;

import { Payload } from '../../types/payload';
import { baseApi } from './baseapi';

export const authApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    register: builder.mutation({
      query: (payload: Payload) => ({
        url: '/v1/users',
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
    getUser: builder.query({
      query: (payload: Payload) => {
        let url = `/v1/users/${payload.body.userId}`;
        
        if(payload.body.import){
          url = `${url}?type=import`;
        }

        return ({
          url: url,
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
  }),
});

export const {
  useRegisterMutation,
  useLazyGetUserQuery,
} = authApi;

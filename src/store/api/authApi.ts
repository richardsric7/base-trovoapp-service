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
    recoveryOtp: builder.mutation({
      query: (payload: Payload) => ({
        url: `/v1/account/recovery/request-email-otp/${payload.body.username}`,
        method: 'POST',        
        data: {
          payload: "",
          creds: {
            signer: payload.signer,
            publicKey: payload.publicKey,
            secretKey: payload.secretKey,
          }
        }, 
      }),
    }),
    verifyEmailOtp: builder.mutation({
      query: (payload: Payload) => ({
        url: `/v1/account/recovery/verify-email-otp/${payload.body.username}/${payload.body.otp}`,
        method: 'POST',        
        data: {
          payload: "",
          creds: {
            signer: payload.signer,
            publicKey: payload.publicKey,
            secretKey: payload.secretKey,
          }
        }, 
      }),
    }),
    fetchSecurityQuestions: builder.query({
      query: (payload: Payload) => {
        let url = `/v1/security-questions/${payload.body.username}`;

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
    submitSecurityAnswers: builder.mutation({
      query: (payload: Payload) => ({
        url: `/v1/verify-answers/${payload.body.username}`,
        method: 'POST',        
        data: {
          payload: payload.body.answers,
          creds: {
            signer: payload.signer,
            publicKey: payload.publicKey,
            secretKey: payload.secretKey,
          }
        }, 
      }),
    }),
    requestAccountRecovery: builder.mutation({
      query: (payload: Payload) => ({
        url: `/v1/users/account/recover`,
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
  useRecoveryOtpMutation,
  useVerifyEmailOtpMutation,
  useFetchSecurityQuestionsQuery,
  useSubmitSecurityAnswersMutation,
  useRequestAccountRecoveryMutation,
  useLazyGetUserQuery,
} = authApi;

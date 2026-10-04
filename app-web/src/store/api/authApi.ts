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
            address: payload.address,
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
            address: payload.address,
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
            address: payload.address,
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
              address: payload.address,
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
            address: payload.address,
            secretKey: payload.secretKey,
          }
        }, 
      }),
    }),
    restoreInactiveAccount: builder.mutation({
      query: (payload: Payload) => ({
        url: '/v1/users/inactive-account/recover',
        method: 'POST',        
        data: {
          payload: payload.body,
          creds: {
            signer: payload.signer,
            address: payload.address,
            secretKey: payload.secretKey,
          }
        }, 
      }),
    }),
    SetupSecurityQuestions: builder.mutation({
      query: (payload: Payload) => ({
        url: '/v1/security-questions',
        method: 'POST',        
        data: {
          payload: payload.body,
          creds: {
            signer: payload.signer,
            address: payload.address,
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
            address: payload.address,
            secretKey: payload.secretKey,
          }
        }, 
      }),
    }),
    // state of the recovery onto the requesting (new) key: PENDING until the
    // recovery period is over, then COMPLETED (or CANCELED by the owner)
    getRecoveryStatus: builder.query({
      query: (payload: Payload) => ({
        url: `/v1/account/recovery/status/${payload.body.username}`,
        method: 'GET',
        data: {
          creds: {
            signer: payload.signer,
            address: payload.address,
            secretKey: payload.secretKey,
          }
        },
      }),
    }),
    getUser: builder.query({
      query: (payload: Payload) => {
        let url = `/v1/users/${payload.body.userId}`;
        
        if(payload.body.import){
          url = `${url}?type=import&type=refresh`;
        }

        return ({
          url: url,
          method: 'GET',        
          data: {
            creds: {
              signer: payload.signer,
              address: payload.address,
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
  useRestoreInactiveAccountMutation,
  useSetupSecurityQuestionsMutation,
  useFetchSecurityQuestionsQuery,
  useSubmitSecurityAnswersMutation,
  useRequestAccountRecoveryMutation,
  useLazyGetRecoveryStatusQuery,
  useLazyGetUserQuery,
} = authApi;

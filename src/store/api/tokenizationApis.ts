import { Payload } from '../../types/payload';
import { baseApi } from './baseapi';

type DocumentUploadPayload = {
  signer: string;
  publicKey: string;
  secretKey: string;
  tokenizedAssetId: string;
  documentTitle: string;
  documentType: string;
  file: File;
};

type DocumentDeletePayload = {
  signer: string;
  publicKey: string;
  secretKey: string;
  documentId: string;
};

type UploadAssetLogoPayload = {
  signer: string;
  publicKey: string;
  secretKey: string;
  tokenizedAssetId: string;
  file: File;
};

type TokenizationDetailPayload = {
  signer: string;
  publicKey: string;
  secretKey: string;
  tokenizedAssetId: string;
};

type ConfirmTokenizationPayload = {
  signer: string;
  publicKey: string;
  secretKey: string;
  tokenizedAssetId: string;
  body: string;
};

type UploadTokenizationFeeProofPayload = {
  signer: string;
  publicKey: string;
  secretKey: string;
  tokenizedAssetId: string;
  tokenizationFeePaymentMethodID: string;
  transactionReference: string;
  file?: File;
};

type DeleteTokenizationFeeProofPayload = {
  signer: string;
  publicKey: string;
  secretKey: string;
  documentId: string;
};

type ConfirmTokenizationFeePayload = {
  signer: string;
  publicKey: string;
  secretKey: string;
  tokenizedAssetId: string;
  body: string;
};

export const tokenizationApi = baseApi.injectEndpoints({
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
    fetchFormFields: builder.query({
      query: (payload: Payload) => {
        return ({
          url: `/v1/forms/${payload.body.formId}`,
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
    fetchTokenizationAssets: builder.query({
      query: (payload: Payload) => ({
        url: `/v1/tokenization/list`,
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
    uploadDocument: builder.mutation({
      query: (payload: DocumentUploadPayload) => {
        const formData = new FormData();
        formData.append('documentFile', payload.file);
        formData.append('tokenizedAssetId', payload.tokenizedAssetId);
        formData.append('documentTitle', payload.documentTitle);
        formData.append('documentType', payload.documentType);

        return {
          url: '/v1/tokenization/document',
          method: 'PUT',
          data: {
            payload: formData,
            creds: {
              signer: payload.signer,
              publicKey: payload.publicKey,
              secretKey: payload.secretKey,
            },
            isFormData: true,
          },
        };
      },
    }),
uploadAssetLogo: builder.mutation({
      query: (payload: UploadAssetLogoPayload) => {
        const formData = new FormData();
        formData.append('documentFile', payload.file);
        formData.append('tokenizedAssetID', payload.tokenizedAssetId);
        formData.append('documentType', '');
        formData.append('documentTitle', '');

        return {
          url: '/v1/tokenization/logo',
          method: 'PUT',
          data: {
            payload: formData,
            creds: {
              signer: payload.signer,
              publicKey: payload.publicKey,
              secretKey: payload.secretKey,
            },
            isFormData: true,
          },
        };
      },
    }),
    deleteDocument: builder.mutation({
      query: (payload: DocumentDeletePayload) => ({
        url: `/v1/tokenization/document/${payload.documentId}`,
        method: 'DELETE',
        data: {
          payload: {},
          creds: {
            signer: payload.signer,
            publicKey: payload.publicKey,
            secretKey: payload.secretKey,
          },
        },
      }),
    }),
    getTokenizationDetail: builder.query({
      query: (payload: TokenizationDetailPayload) => ({
        url: `/v1/tokenization/detail/${payload.tokenizedAssetId}`,
        method: 'GET',
        data: {
          creds: {
            signer: payload.signer,
            publicKey: payload.publicKey,
            secretKey: payload.secretKey,
          },
        },
      }),
    }),
    uploadTokenizationFeeProof: builder.mutation({
      query: (payload: UploadTokenizationFeeProofPayload) => {
        const formData = new FormData();

        if (payload.file) {
          formData.append('documentFile', payload.file);
        }

        formData.append('transactionReference', payload.transactionReference);
        formData.append(
          'tokenizationFeePaymentMethodID',
          payload.tokenizationFeePaymentMethodID,
        );

        return {
          url: `/v1/tokenization/fee/${payload.tokenizedAssetId}`,
          method: 'PUT',
          data: {
            payload: formData,
            creds: {
              signer: payload.signer,
              publicKey: payload.publicKey,
              secretKey: payload.secretKey,
            },
            isFormData: true,
          },
        };
      },
    }),
    deleteTokenizationFeeProof: builder.mutation({
      query: (payload: DeleteTokenizationFeeProofPayload) => ({
        url: `/v1/tokenization/fee/${payload.documentId}`,
        method: 'DELETE',
        data: {
          payload: {},
          creds: {
            signer: payload.signer,
            publicKey: payload.publicKey,
            secretKey: payload.secretKey,
          },
        },
      }),
    }),
    confirmTokenizationFee: builder.mutation({
      query: (payload: ConfirmTokenizationFeePayload) => ({
        url: `/v1/tokenization/fee/${payload.tokenizedAssetId}`,
        method: 'POST',
        data: {
          payload: payload.body,
          creds: {
            signer: payload.signer,
            publicKey: payload.publicKey,
            secretKey: payload.secretKey,
          },
        },
      }),
    }),
    confirmTokenization: builder.mutation({
      query: (payload: ConfirmTokenizationPayload) => ({
        url: `/v1/tokenization/confirm/${payload.tokenizedAssetId}`,
        method: 'PUT',
        data: {
          payload: payload.body,
          creds: {
            signer: payload.signer,
            publicKey: payload.publicKey,
            secretKey: payload.secretKey,
          },
        },
      }),
    }),
    submitTokenization: builder.mutation({
      query: (payload: Payload) => ({
        url: '/v1/tokenization',
        method: 'POST',
        data: {
          payload: payload.body,
          creds: {
            signer: payload.signer,
            publicKey: payload.publicKey,
            secretKey: payload.secretKey,
          },
        },
      }),
    }),
  }),
});

export const {
  useSendAssetMutation,  
  useLazyFetchFormFieldsQuery,
  useFetchTokenizationAssetsQuery,
useUploadDocumentMutation,
  useUploadAssetLogoMutation,
  useDeleteDocumentMutation,
  useLazyGetTokenizationDetailQuery,
  useUploadTokenizationFeeProofMutation,
  useDeleteTokenizationFeeProofMutation,
  useConfirmTokenizationFeeMutation,
  useConfirmTokenizationMutation,
  useSubmitTokenizationMutation,
} = tokenizationApi;

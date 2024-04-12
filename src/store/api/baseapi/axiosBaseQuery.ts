import axios, { AxiosError, AxiosRequestConfig } from 'axios';
import { getStorage } from '../../../utils/storage';
import { TOKEN } from '../../constants';
import { signHTTP } from '../../../utils/trovoSDK';
import { SerializedError } from '@reduxjs/toolkit';


axios.interceptors.request.use(
  async (config: any) => {
    const access_token = getStorage(TOKEN)?.accessToken ?? '';
      config.headers = {
      ...config.headers,
      Authorization: `Bearer ${access_token}`
    };
    return config;
  },
  (error: any) => {
    return Promise.reject(error);
  }
);

type Credentials = {
  signer: string,
  publicKey: string,
  secretKey: string,
}

export interface SuccessResponse {
  data: any;
}

export interface ErrorResponse {
  error: {
      status: any;
      data: any;
  };
}

export const axiosBaseQuery =
  ({
    baseUrl = '',
    baseHeaders = {}
  }) =>
  async ({ url, method, data, headers = {} }: AxiosRequestConfig<any>) => {
    const creds = data.creds as Credentials;
    const payload = data.payload;
    try {
      const result = await axios({
        url: baseUrl + url,
        method,
        data: payload,
        headers: { ...baseHeaders, ...headers, ...getRequestHeaders(url!, creds.signer, creds.publicKey, creds.secretKey) },
      });
      return { data: result.data };
    } catch (axiosError: any) {
      let err = axiosError;
      return {
        error: {
          status: err.response?.status,
          data: {...(typeof err.response?.data === 'object' ? err.response.data : {})},
        },
      };
    }
  };

  const getRequestHeaders = (uri: string, signer: string, publicKey: string, secretKey: string) => {
    const deviceId = navigator.userAgent;
    const ms = Date.now();
    const serverTs = Math.round(ms / 1000).toString();
    const toSign = uri + signer + serverTs;
    console.log('toSign: ', toSign); 
    const signHttp = signHTTP(toSign, secretKey);
    console.log('pubkey: ', publicKey);  
    console.log('uri: ', uri);

    return {
      "X-TW-SIGNATURE": signHttp,
      "X-TW-PUBLIC-KEY": publicKey,
      "X-TW-SIGNER": signer,
      "X-TW-DEVICE-ID": deviceId,
      // "X-TW-APP-VERSION": "1",
      "X-TW-TIMESTAMP": serverTs
    }    
  }

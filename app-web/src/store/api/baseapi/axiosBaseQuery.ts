import axios, { AxiosError, AxiosRequestConfig } from 'axios';
import { getStorage } from '../../../utils/storage';
import { TOKEN } from '../../constants';
import { signHTTP } from '../../../utils/trovoSDK';


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
  address: string,
  secretKey: string,
}

export interface SuccessResponse {
  data: any;
}

export interface ErrorResponse {  
      status: any;
      data: any;
}

export const axiosBaseQuery =
  ({
    baseUrl = '',
    baseHeaders = {}
  }) =>
  async ({ url, method, data, headers = {} }: AxiosRequestConfig<any>) => {
    const creds = data?.creds as Credentials;
    const payload = data?.payload;
    const isFormData = data?.isFormData as boolean | undefined;
    try {
      const requestHeaders: any = creds
        ? { ...baseHeaders, ...headers, ...getRequestHeaders(url!, creds.signer, creds.address, creds.secretKey) }
        : { ...baseHeaders, ...headers };

      // For FormData uploads, let axios set Content-Type automatically
      if (isFormData) {
        delete requestHeaders['Content-Type'];
      }

      const result = await axios({
        url: baseUrl + url,
        method,
        data: payload,
        headers: requestHeaders,
      });
      return { data: result.data };
    } catch (axiosError: any) {
      let err = axiosError as AxiosError;
      return {
        error: {
          status: err.response?.status,
          data: {...(typeof err.response?.data === 'object' ? err.response.data : {})},
        },
      };
    }
  };

  const getRequestHeaders = (uri: string, signer: string, address: string, secretKey: string) => {
    const deviceId = navigator.userAgent;
    const ms = Date.now();
    const serverTs = Math.round(ms / 1000).toString();
    const toSign = uri + signer + serverTs;
    console.log('toSign: ', toSign); 
    const signature = signHTTP(toSign, secretKey);
    console.log('pubkey: ', address);  
    console.log('uri: ', uri);
    console.log('signature: ', signature);
    console.log('timestamp: ', serverTs);

    return {
      "X-TW-SIGNATURE": signature,
      "X-TW-PUBLIC-KEY": address,
      "X-TW-SIGNER": signer,
      "X-TW-DEVICE-ID": deviceId,
      // "X-TW-APP-VERSION": "1",
      "X-TW-TIMESTAMP": serverTs
    }    
  }

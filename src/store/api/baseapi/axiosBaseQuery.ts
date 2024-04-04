import axios, { AxiosError, AxiosRequestConfig } from 'axios';
import { getStorage } from '../../../utils/storage';
import { TOKEN } from '../../constants';
import { BaseQueryFn } from '@reduxjs/toolkit/query';


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

export const axiosBaseQuery =
  ({
    baseUrl = '',
    baseHeaders = {}
  }) =>
  async ({ url, method, data, headers = {} }: AxiosRequestConfig<any>) => {
    const result = await axios({
      url: baseUrl + url,
      method,
      data,
      headers: { ...baseHeaders, ...headers }
    });
    return { data: result.data };
  };

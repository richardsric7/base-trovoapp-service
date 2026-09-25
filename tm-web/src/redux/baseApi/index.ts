import { BASE_URL } from '@/config';
import { createApi } from '@reduxjs/toolkit/query/react';
import { axiosBaseQuery } from './axiosBasedQuery';
import { tagTypes } from './tagTypes';

export const baseApi = createApi({
  baseQuery: axiosBaseQuery({
    baseUrl: BASE_URL,
    baseHeaders: {
      'Content-Type': 'application/json',
    },
  }),
  tagTypes: Object.values(tagTypes),
  endpoints: () => ({}),
});

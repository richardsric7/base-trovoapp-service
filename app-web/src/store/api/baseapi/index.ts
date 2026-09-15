import { createApi } from '@reduxjs/toolkit/query/react'
import { BASE_URL } from '../../config';
import { tagTypes } from './tagTypes';
import { axiosBaseQuery} from './axiosBaseQuery'

// Define a service using a base URL and expected endpoints
export const baseApi = createApi({
  baseQuery: axiosBaseQuery({ baseUrl: BASE_URL}),
  tagTypes: Object.values(tagTypes),
  endpoints: () => ({})
})
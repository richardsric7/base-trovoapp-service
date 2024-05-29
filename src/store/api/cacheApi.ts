import { baseApi } from './baseapi';

export const authApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({    
    fetchFiatRates: builder.query({
      query: () => {
        let url = `/v1/rates`;

        return ({
          url: url,
          method: 'GET',
        })
      },
    }),  
        
    fetchVersionInfo: builder.query({
      query: () => {
        let url = `/v1/app-version`;

        return ({
          url: url,
          method: 'GET',
        })
      },
    }),  
    
    fetchAnnouncements: builder.query({
      query: () => {
        let url = `/v1/announcements`;
        
        return ({
          url: url,
          method: 'GET', 
        })
      },
    }),    
  }),
});

export const {
  useFetchAnnouncementsQuery,
  useFetchFiatRatesQuery,
  useFetchVersionInfoQuery
} = authApi;

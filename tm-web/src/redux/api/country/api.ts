import { baseApi } from "@/redux/baseApi";
import {
  Country,
  CountryPayload,
  ICountryPayload,
  ICountrySaveResponse,
} from "./interface";

export const countryList = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getCountryList: builder.query<CountryPayload, void>({
      query: () => ({
        url: "/country/list",
        method: "GET",
      }),
      providesTags: ["country"],
    }),

    getCountryById: builder.query<Country, number>({
      query: (id) => ({
        url: `/country/get/${id}`,
        method: "GET",
      }),
    }),

    addOrUpdateCountry: builder.mutation<ICountrySaveResponse, ICountryPayload>(
      {
        query: (data) => ({
          url: "/country/save",
          method: "POST",
          data,
        }),
        invalidatesTags: ["country"],
      }
    ),

    deleteCountry: builder.mutation<ICountrySaveResponse, number>({
      query: (id) => ({
        url: `/country/delete/${id}`,
        method: "DELETE",
      }),
      invalidatesTags: ["country"],
    }),
  }),
});

export const {
  useGetCountryListQuery,
  useGetCountryByIdQuery,
  useAddOrUpdateCountryMutation,
  useDeleteCountryMutation,
} = countryList;

import { baseApi } from "@/redux/baseApi";
import { tagTypes } from "@/redux/baseApi/tagTypes";
import {
  ServiceLinkListQueryParams,
  ServiceLinkListResponse,
  ServiceLinkRequest,
  ServiceLinkResponse,
} from "./interface";

export const serviceLinksApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getServiceLinks: builder.query<ServiceLinkListResponse, ServiceLinkListQueryParams>({
      query: (params) => ({
        url: "/service-links",
        method: "GET",
        params,
      }),
      providesTags: [tagTypes.SERVICE_LINKS],
    }),

    getServiceLinkById: builder.query<ServiceLinkResponse, string>({
      query: (id) => ({
        url: `/service-links/${id}`,
        method: "GET",
      }),
      providesTags: [tagTypes.SERVICE_LINKS],
    }),

    saveServiceLink: builder.mutation<ServiceLinkResponse, ServiceLinkRequest>({
      query: (data) => ({
        url: "/service-links",
        method: "POST",
        data,
      }),
      invalidatesTags: [tagTypes.SERVICE_LINKS],
    }),

    setServiceLinkInactive: builder.mutation<ServiceLinkResponse, { id: string; inactive: boolean }>({
      query: ({ id, inactive }) => ({
        url: `/service-links/${id}/inactive`,
        method: "PUT",
        data: { inactive },
      }),
      invalidatesTags: [tagTypes.SERVICE_LINKS],
    }),
  }),
});

export const {
  useGetServiceLinksQuery,
  useGetServiceLinkByIdQuery,
  useSaveServiceLinkMutation,
  useSetServiceLinkInactiveMutation,
} = serviceLinksApi;

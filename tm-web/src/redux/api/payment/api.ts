import { baseApi } from "@/redux/baseApi";
import { FiatPaymentRecordsResponse, GetFiatPaymentsParams } from "./interface";

export const paymentApi = baseApi.injectEndpoints({
  endpoints: (build) => ({
    getFiatPayments: build.query<
      FiatPaymentRecordsResponse,
      GetFiatPaymentsParams
    >({
      query: (params) => ({
        url: "/fiat/payments",
        params: params,
        method: "GET",
      }),
    }),
  }),
  overrideExisting: false,
});

export const { useGetFiatPaymentsQuery } = paymentApi;

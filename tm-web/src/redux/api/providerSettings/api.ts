import { baseApi } from "@/redux/baseApi";
import { tagTypes } from "@/redux/baseApi/tagTypes";
import { ProviderSettings, SaveKycProviderRequest, SaveStablerailRequest } from "./interface";

export const providerSettingsApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getProviderSettings: builder.query<ProviderSettings, void>({
      query: () => ({ url: "/provider-settings", method: "GET" }),
      providesTags: [tagTypes.PROVIDER_SETTINGS],
    }),

    saveKycProvider: builder.mutation<ProviderSettings, SaveKycProviderRequest>({
      query: ({ provider, ...data }) => ({
        url: `/provider-settings/kyc/${encodeURIComponent(provider)}`,
        method: "PUT",
        data,
      }),
      invalidatesTags: [tagTypes.PROVIDER_SETTINGS],
    }),

    saveStablerail: builder.mutation<ProviderSettings, SaveStablerailRequest>({
      query: (data) => ({ url: "/provider-settings/stablerail", method: "PUT", data }),
      invalidatesTags: [tagTypes.PROVIDER_SETTINGS],
    }),
  }),
});

export const { useGetProviderSettingsQuery, useSaveKycProviderMutation, useSaveStablerailMutation } =
  providerSettingsApi;

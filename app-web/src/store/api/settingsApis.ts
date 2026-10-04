import { baseApi } from './baseapi';
import { P2PCreds } from './p2pApis';

// A curated asset users can pay network fees in (see app-backend's
// GET /v1/users/settings/gas-fee-assets).
export type GasFeeAssetOption = {
  assetCode: string;
  assetName: string;
  contractAddress: string;
  imageUrl?: string | null;
};

export type GasFeeAssetsResponse = {
  gasFeeAsset: string | null;
  assets: GasFeeAssetOption[];
};

export const settingsApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getGasFeeAssets: builder.query<GasFeeAssetsResponse, { creds: P2PCreds }>({
      query: ({ creds }) => ({
        url: '/v1/users/settings/gas-fee-assets',
        method: 'GET',
        data: { creds },
      }),
    }),
    setGasFeeAsset: builder.mutation<
      { gasFeeAsset: string | null },
      { creds: P2PCreds; assetCode: string }
    >({
      query: ({ creds, assetCode }) => ({
        url: '/v1/users/settings/gas-fee-asset',
        method: 'PUT',
        data: { payload: { assetCode }, creds },
      }),
    }),
  }),
});

export const { useGetGasFeeAssetsQuery, useSetGasFeeAssetMutation } =
  settingsApi;

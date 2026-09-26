import { baseApi } from "@/redux/baseApi";
import {
  AssetClassListResponse,
  CuratedAssetListQueryParams,
  CuratedAssetListResponse,
  CuratedAssetRequest,
  CuratedAssetResponse,
} from "./interface";

export const curatedAssetsApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getCuratedAssets: builder.query<CuratedAssetListResponse, CuratedAssetListQueryParams>({
      query: (params) => ({
        url: "/assets/curated",
        method: "GET",
        params,
      }),
      providesTags: ["curatedAssets"],
    }),

    getCuratedAssetById: builder.query<CuratedAssetResponse, number>({
      query: (id) => ({
        url: `/assets/curated/${id}`,
        method: "GET",
      }),
      providesTags: ["curatedAssets"],
    }),

    saveCuratedAsset: builder.mutation<CuratedAssetResponse, CuratedAssetRequest>({
      query: (data) => ({
        url: "/assets/curated",
        method: "POST",
        data,
      }),
      invalidatesTags: ["curatedAssets"],
    }),

    setCuratedAssetP2PEnabled: builder.mutation<CuratedAssetResponse, { id: number; p2pEnabled: boolean }>({
      query: ({ id, p2pEnabled }) => ({
        url: `/assets/curated/${id}/p2p-enabled`,
        method: "PUT",
        data: { p2pEnabled },
      }),
      invalidatesTags: ["curatedAssets"],
    }),

    getAssetClasses: builder.query<AssetClassListResponse, void>({
      query: () => ({
        url: "/asset-classes",
        method: "GET",
      }),
    }),
  }),
});

export const {
  useGetCuratedAssetsQuery,
  useGetCuratedAssetByIdQuery,
  useSaveCuratedAssetMutation,
  useSetCuratedAssetP2PEnabledMutation,
  useGetAssetClassesQuery,
} = curatedAssetsApi;

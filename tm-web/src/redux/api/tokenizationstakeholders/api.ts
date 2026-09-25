import { baseApi } from "@/redux/baseApi";
import {
  PartnerListResponse,
  AssetManager,
  AssetIssuingHouse,
  AssetCustodian,
  PartnerType,
  LegalAndProfessional,
  RatingAgency,
  Trustee,
  LegalAdviser,
  FinancialAdviser,
  StakeholderData,
} from "./interface";

// Map type to result type
type PartnerMap = {
  asset_manager: AssetManager[];
  asset_issuing_house: AssetIssuingHouse[];
  approved_asset_custodian: AssetCustodian[];
  legal_and_professionals: LegalAndProfessional[];
  rating_agency: RatingAgency[];
  trustees: Trustee[];
  legal_adviser: LegalAdviser[];
  financial_adviser: FinancialAdviser[];
};

type PartnerQueryArg = { type?: PartnerType };

export const assetTokenization = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getPartnersList: builder.query<
      PartnerListResponse | PartnerMap[PartnerType],
      PartnerQueryArg
    >({
      query: (params) => ({
        url: "/partners/list",
        method: "GET",
        params,
      }),
      providesTags: ["stakeholders"],
    }),

    getSinglePartner: builder.query<
      StakeholderData,
      { id: number; type: PartnerType }
    >({
      query: ({ id, type }) => ({
        url: "/partners/one",
        method: "GET",
        params: { id, type },
      }),
    }),

    createPartner: builder.mutation({
      query: ({ type, action, data }) => ({
        url: `/partners/save?type=${type}&action=${action}`,
        method: "POST",
        data,
      }),
      invalidatesTags: ["stakeholders"],
    }),

    deletePartner: builder.mutation<void, { type: PartnerType; id: number }>({
      query: ({ type, id }) => ({
        url: `/partners/delete`,
        method: "DELETE",
        params: { type, id },
      }),
      invalidatesTags: ["stakeholders"],
    }),
  }),
});

export const {
  useGetPartnersListQuery,
  useGetSinglePartnerQuery,
  useCreatePartnerMutation,
  useDeletePartnerMutation,
} = assetTokenization;

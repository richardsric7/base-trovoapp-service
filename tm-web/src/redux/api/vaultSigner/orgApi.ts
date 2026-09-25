import { orgApi } from "@/redux/baseApi/orgApi";
import { tagTypes } from "@/redux/baseApi/tagTypes";
import {
  ICreatePersonalEnvRequest,
  ICreatePersonalEnvResponse,
  IMySecretListResponse,
  IMySecretValueResponse,
  IPersonalEnvListResponse,
  IPutSecretRequest,
  IPutSecretResponse,
} from "./interface";

// Organisation-portal counterpart of api.ts's self-service endpoints, authenticated via
// orgApi's raw-org-token axios interceptor instead of baseApi's Bearer JWT. The backend
// endpoint (AllowOrgOrTrovoAdminNormalized) is identical either way — only the axios
// instance/auth header differs, matching how getOrganizationMembers is dual-injected into
// both baseApi (redux/api/organizations) and orgApi (redux/api/org) today. Hooks are
// suffixed "Org" so they can be re-exported alongside api.ts's hooks from ./index without
// name collisions.
export const vaultSignerOrgApi = orgApi.injectEndpoints({
  endpoints: (builder) => ({
    listMySecretsOrg: builder.query<IMySecretListResponse, void>({
      query: () => ({
        url: "/me/vault-signer/secrets",
        method: "GET",
      }),
      providesTags: [tagTypes.VAULT_SIGNER_SECRETS],
    }),
    getMySecretValueOrg: builder.query<IMySecretValueResponse, string>({
      query: (secretId) => ({
        url: `/me/vault-signer/secrets/${secretId}`,
        method: "GET",
      }),
    }),
    putMySecretValueOrg: builder.mutation<
      IPutSecretResponse,
      { secretId: string; payload: IPutSecretRequest }
    >({
      query: ({ secretId, payload }) => ({
        url: `/me/vault-signer/secrets/${secretId}`,
        method: "PUT",
        body: payload,
      }),
      invalidatesTags: [tagTypes.VAULT_SIGNER_SECRETS],
    }),
    listPersonalEnvsOrg: builder.query<IPersonalEnvListResponse, void>({
      query: () => ({
        url: "/me/vault-signer/personal-envs",
        method: "GET",
      }),
      providesTags: [tagTypes.VAULT_SIGNER_PERSONAL_ENVS],
    }),
    createPersonalEnvOrg: builder.mutation<
      ICreatePersonalEnvResponse,
      { prefix: string; payload: ICreatePersonalEnvRequest }
    >({
      query: ({ prefix, payload }) => ({
        url: `/me/vault-signer/personal-envs/${prefix}`,
        method: "POST",
        body: payload,
      }),
      invalidatesTags: [tagTypes.VAULT_SIGNER_PERSONAL_ENVS],
    }),
    deletePersonalEnvOrg: builder.mutation<void, string>({
      query: (prefix) => ({
        url: `/me/vault-signer/personal-envs/${prefix}`,
        method: "DELETE",
      }),
      invalidatesTags: [tagTypes.VAULT_SIGNER_PERSONAL_ENVS],
    }),
  }),
  overrideExisting: false,
});

export const {
  useListMySecretsOrgQuery,
  useGetMySecretValueOrgQuery,
  usePutMySecretValueOrgMutation,
  useListPersonalEnvsOrgQuery,
  useCreatePersonalEnvOrgMutation,
  useDeletePersonalEnvOrgMutation,
} = vaultSignerOrgApi;

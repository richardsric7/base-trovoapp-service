import { baseApi } from "@/redux/baseApi";
import { tagTypes } from "@/redux/baseApi/tagTypes";
import {
  ICreatePersonalEnvRequest,
  ICreatePersonalEnvResponse,
  IDeleteAssignmentRequest,
  IDeleteAssignmentResponse,
  IEditAssignmentRequest,
  IListAssignmentsResponse,
  IListAuditLogResponse,
  IListManagedSecretsResponse,
  IMySecretListResponse,
  IMySecretValueResponse,
  IPersonalEnvListResponse,
  IPutSecretRequest,
  IPutSecretResponse,
  IRegisterManagedSecretRequest,
  IRegisterManagedSecretResponse,
  ICreateAssignmentRequest,
  IAssignmentResponse,
} from "./interface";

// This file mixes two endpoint groups with two different backend gates — do not assume one
// middleware covers both:
//
// - Self-service (/me/vault-signer/*) — backend: middleware.AllowOrgOrTrovoAdminNormalized.
//   Reachable by EITHER a Trovo Admin or an Organization member. This injection is the
//   dashboard side, authenticated via baseApi's Bearer-JWT axios interceptor. The same six
//   endpoints are also injected into orgApi (see ./orgApi.ts) for the organisation portal,
//   exactly like getOrganizationMembers already is for both portals in redux/api/organizations
//   and redux/api/org.
// - Admin (/admin/vault-signer/*) — backend: middleware.AuthenticateSuperAdmin. Trovo
//   SuperAdmin ONLY — this middleware has no OrganizationAuth branch at all, so an org member's
//   token is rejected outright (401) before the role check even runs; it is never reachable
//   from the organisation portal, which is why these endpoints are only ever injected here into
//   baseApi and have no ./orgApi.ts counterpart. Covers all three admin surfaces uniformly:
//   managed-secrets*, managed-secrets/:id/assignments*, and audit-log — none of them fall back
//   to the looser self-service gate.
export const vaultSignerApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    // --- Self-service (AllowOrgOrTrovoAdminNormalized): signer slots ---
    listMySecrets: builder.query<IMySecretListResponse, void>({
      query: () => ({
        url: "/me/vault-signer/secrets",
        method: "GET",
      }),
      providesTags: [tagTypes.VAULT_SIGNER_SECRETS],
    }),
    getMySecretValue: builder.query<IMySecretValueResponse, string>({
      query: (secretId) => ({
        url: `/me/vault-signer/secrets/${secretId}`,
        method: "GET",
      }),
    }),
    putMySecretValue: builder.mutation<
      IPutSecretResponse,
      { secretId: string; payload: IPutSecretRequest }
    >({
      query: ({ secretId, payload }) => ({
        url: `/me/vault-signer/secrets/${secretId}`,
        method: "PUT",
        data: payload,
      }),
      invalidatesTags: [tagTypes.VAULT_SIGNER_SECRETS],
    }),

    // --- Self-service (AllowOrgOrTrovoAdminNormalized): personal secrets ---
    listPersonalEnvs: builder.query<IPersonalEnvListResponse, void>({
      query: () => ({
        url: "/me/vault-signer/personal-envs",
        method: "GET",
      }),
      providesTags: [tagTypes.VAULT_SIGNER_PERSONAL_ENVS],
    }),
    createPersonalEnv: builder.mutation<
      ICreatePersonalEnvResponse,
      { prefix: string; payload: ICreatePersonalEnvRequest }
    >({
      query: ({ prefix, payload }) => ({
        url: `/me/vault-signer/personal-envs/${prefix}`,
        method: "POST",
        data: payload,
      }),
      invalidatesTags: [tagTypes.VAULT_SIGNER_PERSONAL_ENVS],
    }),
    deletePersonalEnv: builder.mutation<void, string>({
      query: (prefix) => ({
        url: `/me/vault-signer/personal-envs/${prefix}`,
        method: "DELETE",
      }),
      invalidatesTags: [tagTypes.VAULT_SIGNER_PERSONAL_ENVS],
    }),

    // --- Admin (AuthenticateSuperAdmin — Trovo SuperAdmin only): managed secrets ---
    listManagedSecrets: builder.query<IListManagedSecretsResponse, void>({
      query: () => ({
        url: "/admin/vault-signer/managed-secrets",
        method: "GET",
      }),
      providesTags: [tagTypes.VAULT_SIGNER_MANAGED_SECRETS],
    }),
    registerManagedSecret: builder.mutation<
      IRegisterManagedSecretResponse,
      IRegisterManagedSecretRequest
    >({
      query: (payload) => ({
        url: "/admin/vault-signer/managed-secrets",
        method: "POST",
        data: payload,
      }),
      invalidatesTags: [tagTypes.VAULT_SIGNER_MANAGED_SECRETS],
    }),

    // --- Admin (AuthenticateSuperAdmin — Trovo SuperAdmin only): assignments ---
    listAssignments: builder.query<IListAssignmentsResponse, string>({
      query: (managedSecretId) => ({
        url: `/admin/vault-signer/managed-secrets/${managedSecretId}/assignments`,
        method: "GET",
      }),
      providesTags: [tagTypes.VAULT_SIGNER_ASSIGNMENTS],
    }),
    createAssignment: builder.mutation<IAssignmentResponse, ICreateAssignmentRequest>({
      query: ({ managedSecretId, ...payload }) => ({
        url: `/admin/vault-signer/managed-secrets/${managedSecretId}/assignments`,
        method: "POST",
        data: payload,
      }),
      invalidatesTags: [
        tagTypes.VAULT_SIGNER_MANAGED_SECRETS,
        tagTypes.VAULT_SIGNER_ASSIGNMENTS,
      ],
    }),
    editAssignment: builder.mutation<IAssignmentResponse, IEditAssignmentRequest>({
      query: ({ managedSecretId, assignmentId, ...payload }) => ({
        url: `/admin/vault-signer/managed-secrets/${managedSecretId}/assignments/${assignmentId}`,
        method: "PATCH",
        data: payload,
      }),
      invalidatesTags: [
        tagTypes.VAULT_SIGNER_MANAGED_SECRETS,
        tagTypes.VAULT_SIGNER_ASSIGNMENTS,
      ],
    }),
    deleteAssignment: builder.mutation<IDeleteAssignmentResponse, IDeleteAssignmentRequest>({
      query: ({ managedSecretId, assignmentId }) => ({
        url: `/admin/vault-signer/managed-secrets/${managedSecretId}/assignments/${assignmentId}`,
        method: "DELETE",
      }),
      invalidatesTags: [
        tagTypes.VAULT_SIGNER_MANAGED_SECRETS,
        tagTypes.VAULT_SIGNER_ASSIGNMENTS,
      ],
    }),

    // --- Admin (AuthenticateSuperAdmin — Trovo SuperAdmin only): audit log ---
    listAuditLog: builder.query<IListAuditLogResponse, void>({
      query: () => ({
        url: "/admin/vault-signer/audit-log",
        method: "GET",
      }),
      providesTags: [tagTypes.VAULT_SIGNER_AUDIT_LOG],
    }),
  }),
});

export const {
  useListMySecretsQuery,
  useGetMySecretValueQuery,
  usePutMySecretValueMutation,
  useListPersonalEnvsQuery,
  useCreatePersonalEnvMutation,
  useDeletePersonalEnvMutation,
  useListManagedSecretsQuery,
  useRegisterManagedSecretMutation,
  useListAssignmentsQuery,
  useCreateAssignmentMutation,
  useEditAssignmentMutation,
  useDeleteAssignmentMutation,
  useListAuditLogQuery,
} = vaultSignerApi;

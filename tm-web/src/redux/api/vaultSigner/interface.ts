// Types are matched directly against the backend handler code
// (internal/components/vaultsigner/handlers/*.go, models/db_models.go) rather than the
// generic `additionalProperties: true` swagger stubs, which don't carry field names.

// ---- Self-service: signer slots (/me/vault-signer/secrets*) ----

export interface IMySecretListItem {
  managedSecretId: string;
  label: string;
  position: number | null;
}

export interface IMySecretListResponse {
  data: IMySecretListItem[];
}

export interface IMySecretValue {
  position: number | null;
  value?: string;
}

export interface IMySecretValueResponse {
  data: IMySecretValue;
}

export interface IPutSecretRequest {
  newValue: string;
}

export interface IPutSecretResponseData {
  vaultVersion: number;
  stellarTxHash?: string;
  stellarTxStatus?: string;
  error?: string;
}

export interface IPutSecretResponse {
  data: IPutSecretResponseData;
}

// ---- Self-service: personal secrets (/me/vault-signer/personal-envs*) ----

export interface IPersonalEnvListItem {
  prefix: string;
  friendlyLabel: string;
  vaultKey: string;
  exists: boolean;
}

export interface IPersonalEnvListResponse {
  data: IPersonalEnvListItem[];
}

export interface ICreatePersonalEnvRequest {
  newValue: string;
}

export interface ICreatePersonalEnvResponseData {
  prefix: string;
  vaultKey: string;
  vaultVersion: number;
}

export interface ICreatePersonalEnvResponse {
  data: ICreatePersonalEnvResponseData;
}

// ---- Admin: managed secrets (/admin/vault-signer/managed-secrets*) ----

export type OwnerType = "trovo_admin" | "org_member";

export interface IManagedSecret {
  id: string;
  label: string;
  vault_mount: string;
  vault_path: string;
  vault_field: string;
  wallet_public_key: string;
  active_signing_count: number;
  created_at: string;
}

export interface IAssignment {
  id: string;
  managedSecretId: string;
  position: number | null;
  ownerRefId: string;
  ownerType: OwnerType;
  ownerLabel: string;
  assignedAt: string;
}

export interface IManagedSecretWithAssignments {
  managedSecret: IManagedSecret;
  assignments: IAssignment[];
}

export interface IListManagedSecretsResponse {
  data: IManagedSecretWithAssignments[];
}

export interface IRegisterManagedSecretRequest {
  label: string;
  vaultMount: string;
  vaultPath: string;
  vaultField: string;
  walletPublicKey: string;
  activeSigningCount: number;
}

export interface IRegisterManagedSecretResponse {
  data: IManagedSecret;
}

export interface IListAssignmentsResponse {
  data: IAssignment[];
}

export interface IOwnerRequest {
  ownerIdentifier: string;
  ownerType: OwnerType;
}

export interface IAssignmentResponse {
  data: IAssignment;
}

export interface IOnChainRemoval {
  attempted: boolean;
  stellarTxHash?: string;
  stellarTxStatus?: string;
}

export interface IDeleteAssignmentResponseData {
  claimed: boolean;
  deletedPosition?: number;
  renumberedCount?: number;
  vaultCollapsed?: boolean;
  onChainRemoval?: IOnChainRemoval | null;
}

export interface IDeleteAssignmentResponse {
  data: IDeleteAssignmentResponseData;
}

export interface IAuditLogRow {
  id: string;
  kind: string;
  managed_secret_id?: string;
  position?: number;
  vault_key?: string;
  actor_id: string;
  actor_type: string;
  vault_version_before?: number;
  vault_version_after?: number;
  stellar_tx_hash?: string;
  stellar_tx_status?: string;
  changed_at: string;
  ip_address?: string;
}

export interface IListAuditLogResponse {
  data: IAuditLogRow[];
}

export interface ICreateAssignmentRequest extends IOwnerRequest {
  managedSecretId: string;
}

export interface IEditAssignmentRequest extends IOwnerRequest {
  managedSecretId: string;
  assignmentId: string;
}

export interface IDeleteAssignmentRequest {
  managedSecretId: string;
  assignmentId: string;
}

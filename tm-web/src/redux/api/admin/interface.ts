export interface IInviteAdmin {
  admin_role: string;
  email_or_username: string;
}

// Define User with all required properties of a single admin
export interface User {
  ID: number;
  Email: string;
  Username: string;
  FirstName: string;
  LastName: string;
  Role: string;
  Status: string;
  CreatedAt: string;
  UpdatedAt: string;
  IsAdmin: boolean;
  WalletUserID: string;
}
export interface IViolationOption {
  id: number;
  label: string;
}

// Define IAdminListResponse for a list of admins
export interface IAdminListResponse {
  admins: User[];

  pagination: {
    page: number;
    pageSize: number;
    total: number;
    totalPages: number;
  };
}

export interface IAdminListRequest {
  page?: number;
  pageSize?: number;
}

export type IAdminResponse = User;

export interface ISuspendAdmin {
  email: string;
  suspension_note: string;
  suspension_reason_id: number;
}

export interface IUpdateAdminRolePayload {
  email: string;
  admin_role: string;
}

export interface IRevokeAdminSuspension {
  email: string;
}

export interface ICreateCustodianComplianceRequest {
  custodian_org_id: string;
  category: string;
  requirement: string;
  due_date?: string;
  asset_id?: string;
}

// Mirrors ComplianceItem (backend/internal/components/stakeholder/models/db_models.go)
export interface ICustodianComplianceItem {
  id: string;
  org_id: string;
  asset_id?: string;
  category: string;
  requirement: string;
  status: string;
  due_date?: string;
  completed_at?: string;
  completed_by_member_id?: string;
  created_at: string;
  updated_at: string;
}

export interface ICreateCustodianComplianceResponse {
  data: ICustodianComplianceItem;
}

export interface IGetSuspensionHistory {
  ActionPerformedBy: string;
  actionType: string;
  email: string;
  id: 0;
  Reason: string;
  SuspensionDateTime: string;
  SuspensionNote: string;
  username: string;
}

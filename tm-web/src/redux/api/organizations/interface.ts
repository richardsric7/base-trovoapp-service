export interface GetOrganizationsQueryParams {
  page?: number; // Page number (default: 1)
  pageSize?: number; // Page size (default: 10)
  name?: string; // Partial match on organization name
  email?: string; // Partial match on organization email
  type?: string; // Exact match on organization type
  status?: string; // Exact match on organization status
  createdBy?: string; // Partial match on createdBy
  search?: string;
}

export interface Organization {
  id: string;
  name: string;
  email: string;
  type: string;
  status: string;
  createdBy: string;
  created_at: string;
  updated_at: string;
  stakeholder_type: string;
  stakeholder_id: number;
  created_by_admin: AdminUser;
  team_member_count: number;
  address?: string;
  country?: string;
  fee_fixed?: string;
  fee_percent?: string;
}

export interface OrganizationDetailsData {
  organization: Organization;
  organization_admin: OrganizationAdmin;
  stakeholder: Stakeholder;
}

export interface OrganizationDetailsResponse {
  message: string;
  status: "OK" | string;
  timestamp?: string;
  data: OrganizationDetailsData;
}
export interface OrganizationAdmin {
  first_name: string;
  last_name: string;
  email: string;
  role: "ROOT_SUPER_ADMIN" | string;
  status: "ACTIVE" | "INACTIVE" | string;
}

export interface Stakeholder {
  id: number;
  agency_name: string;
  agency_address: string;
  agency_country: string;
  fee_percent: number;
  fee_fixed: number;
}
export interface OrganizationDetailsData {
  organization: Organization;
  organization_admin: OrganizationAdmin;
  stakeholder: Stakeholder;
}

export interface AdminUser {
  first_name: string;
  last_name: string;
  email: string;
  username: string;
}

export interface PaginationInfo {
  page: number;
  pageSize: number;
  total: number;
  totalPages: number;
}

export interface GetOrganizationsResponse {
  organizations: Organization[];
  pagination: PaginationInfo;
}

export interface IInviteOragnization {
  organization_name: string;
  address: string;
  country: string;
  admin_first_name: string;
  admin_last_name: string;
  admin_email: string;
  fee_fixed: string;
  fee_percent: string;
}

export interface IInviteMember {
  email: string;
  first_name: string;
  last_name: string;
  organization_id?: string;
}

export interface IOrganizationMember {
  id: string;
  email: string;
  role: string;
  status: string;
  created_at: string;
}

export interface IValidateInviteResponse {
  valid: boolean;
  message: string;
}

export interface IValidateEmailWithOtpRequest {
  inviteId: string;
  otp: string;
}

export interface ISetupRootUserPasswordRequest {
  inviteId: string;
  otp: string;
  password: string;
}

export interface IMemberLoginPayload {
  email: string;
  password: string;
}
export interface IMemberLoginResponse {
  message: string;
  data: {
    token: string;
    // userInfo?: any;  // if the API ever returns userInfo in the future
  };
  timestamp: string;
  status: string;
}

export interface IResendOtpPayload {
  email: string;
  token: string;
}

export interface IResendPasswordSetupPayload {
  email: string;
}
export interface IVerifyOtpPayload {
  email: string;
  otp: string;
  token: string;
}

export interface IResendPasswordSetupResponse {
  [key: string]: string;
}

export interface IOrganizationMember {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  role: string;
  status: string;
  created_at: string;
}

export interface IMembersResponse {
  message: string;
  data: {
    members: IOrganizationMember[];
    pagination: {
      page: number;
      pageSize: number;
      total: number;
      totalPages: number;
    };
  };
  timestamp: string;
  status: string;
}
export interface IInviteOrganizationResponse {
  invite_id: string;
  message: string;
  organization_id: string;
}

export interface IInviteOrganizationRequest {
  organization_name: string;
  address: string;
  country: string;
  admin_first_name: string;
  admin_last_name: string;
  admin_email: string;
  fee_fixed: string;
  fee_percent: string;
  stakeholder_type?: string;
}

export interface UpdateOrganizationPayload {
  organizationId: string;
  stakeholderType: string;
  data: {
    organization_name: string;
    organization_type: string;
    address: string;
    country: string;
    admin_first_name: string;
    admin_last_name: string;
    admin_email: string;
    fee_fixed: string;
    fee_percent: string;
  };
}

// deactivate
export interface IDeactivateOrganizationRequest {
  organization_id: string;
}

export interface IDeactivateOrganizationResponse {
  message: string;
  organization_id: string;
  organization_name: string;
  status: string;
}

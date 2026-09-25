export interface OrganizationDetailsLogin {
  message: string;
  status: "OK" | "ERROR";
  timestamp: string;
  data: OrganizationDetailsData;
}
export interface OrganizationDetailsData {
  organization: Organization;
  organization_admin: OrganizationAdmin;
  stakeholder: Stakeholder;
}

export interface Organization {
  id: string;
  name: string;
  email: string;
  type: string;
  status: "PENDING" | "ACTIVE" | "INACTIVE" | string;
  stakeholder_id: number;
  stakeholder_type: string;
  team_member_count: number;
  created_at: string;
  updated_at: string;
  created_by: string;
}

export interface OrganizationAdmin {
  first_name: string;
  last_name: string;
  email: string;
  role: string;
  status: "ACTIVE" | "INACTIVE" | string;
}

export interface Stakeholder {
  id: number;
  asset_issuing_house_name: string;
  asset_issuing_house_address: string;
  asset_issuing_house_country: string;
  fee_fixed: number;
  fee_percent: number;
}

export type FundReleaseItem = import("../trustees/interface").IFundReleaseRecord;
export type FundReleasesResponse = import("../trustees/interface").IFundReleaseListResponse;

export interface GetFundReleasesParams {
  status?: string;
  asset_id?: string;
  page?: number;
  limit?: number;
}

export interface FundReleaseExecutionPayload {
  status: "processing" | "completed" | "failed";
  execution_reference?: string;
  failure_reason?: string;
}

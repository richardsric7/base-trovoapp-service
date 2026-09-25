import { ComplianceRequirementInstance } from "@/redux/api/compliance/interface";

export interface CompliancePaginationMeta {
  page?: number;
  limit?: number;
  total?: number;
  total_pages?: number;
}

export interface IOrgComplianceListParams {
  page?: number;
  limit?: number;
  status?: string;
}

export interface IOrgComplianceListResponse {
  data: {
    records: ComplianceRequirementInstance[];
    meta: CompliancePaginationMeta;
  };
}

export interface ISubmitComplianceItemRequest {
  itemId: string;
  payload: {
    submission_data: Record<string, unknown>;
  };
}

export interface IUpdateOrgComplianceItemRequest {
  itemId: string;
  payload: {
    status: string;
  };
}

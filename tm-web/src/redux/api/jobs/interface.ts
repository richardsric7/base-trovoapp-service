export type WorkType = "fulltime" | "parttime" | "contract";
export type LocationType = "remote" | "onsite" | "hybrid";

export type SectionKey =
  | "About The Role"
  | "What We're Looking For"
  | "What You'll Own"
  | "What Success Looks Like";

export type JobPostPayload = {
  heading: string;
  company_overview: string;
  location: LocationType;
  years_of_experience: string;
  work_type: WorkType;
  sections: Record<SectionKey, string[]>;
  application: {
    apply: string;
    location: string;
    subject: string;
  };
};

export interface IJob {
  id: string;
  role: string;
  heading: string;
  company_overview: string;
  location: LocationType;
  years_of_experience: string;
  work_type: WorkType;
  sections: Record<SectionKey, string[]>;
  application: {
    apply: string;
    location: string;
    subject: string;
  };
  status: string;
  created_at: string;
  updated_at: string;
}

export interface JobResponse {
  message: string;
  data: IJob;
  timestamp: string;
  status: "OK" | "ERROR";
}

export interface JobsPagination {
  page: number;
  pageSize: number;
  total: number;
  totalPages: number;
}

export interface GetJobsResponse {
  message: string;
  data: IJob[]; // ✅ array directly
  pagination: JobsPagination; // optional (backend may add later)
  timestamp: string;
  status: "OK";
}

export interface GetJobsParams {
  status: "published" | "draft" | "disabled";
  search?: string;
  location?: string;
  work_type?: string;
  years_of_experience?: string;
  page?: number;
  pageSize?: number;
}

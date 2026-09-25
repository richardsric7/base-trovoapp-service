// One row of the admin access / audit trail — the shape returned by the backend
// GET /api/v1/audit-trail. Mirrors models.AdminAccessLog plus the pre-formatted
// `date` display string the table prints verbatim.
export interface IAuditRecord {
  id: string;
  event: string;
  status: "successful" | "failed" | "pending";
  action: string; // humanised label for the "Action" column
  category: string;

  actor_admin_id?: number;
  username: string;
  fullname: string;
  email?: string;
  role?: string;

  target?: string;

  ip_address?: string;
  location?: string;
  user_agent?: string;
  method?: string;
  path?: string;
  login_method?: string;
  detail?: string;

  occurred_at: string; // ISO
  date: string; // "DD Mon, YYYY hh:mm AM/PM" — printed as-is by the table
}

export interface IAuditTrailParams {
  page?: number;
  pageSize?: number;
  search?: string;
  status?: "successful" | "failed" | "pending" | string;
  date?: string; // YYYY-MM-DD
}

// The detail endpoint returns a single record under `data`.
export interface IAuditRecordResponse {
  status: number;
  message: string;
  data: IAuditRecord;
}

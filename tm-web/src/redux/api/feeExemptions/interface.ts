// These mirror tm-api's /fee/exempt-users endpoints, which write
// app-backend's fee_exempt_users table: accounts that pay no platform
// service fees (the platform's own trading or operations accounts).
// app-backend reads the list whenever it charges a fee.

export interface IFeeExemptUser {
  username: string;
  reason: string;
  addedBy: string;
  createdAt: string;
}

export interface FeeExemptUsersResponse {
  message: string;
  data: IFeeExemptUser[];
  status: string;
}

export interface AddFeeExemptUserRequest {
  username: string;
  reason: string;
}

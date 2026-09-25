export interface MintingUser {
  id: number;
  username: string;
  role: string;
  created_at: string;
  updated_at: string;
}

// In your mintingUserApi.ts
export interface MintingUserQueryParams {
  role?: string;
  page?: number;
  page_size?: number;
}

export interface MintingUserListMeta {
  data: MintingUser[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface MintingUserApiResponse {
  message: string;
  data: MintingUserListMeta;
  timestamp: string;
  status: string;
}
export interface MintingUserPayload {
  users: MintingUser[];
  total: number;
}

export interface MintingUserSaveResponse {
  message: string;
  data?: MintingUser[];
  success: boolean;
}

export interface SearchMintingUsersParams {
  query: string;
  role?: string;
  page?: number;
  page_size?: number;
}

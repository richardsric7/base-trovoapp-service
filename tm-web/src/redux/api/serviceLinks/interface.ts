// These mirror tm-api's /service-links endpoints
// (usermetrics/db/service_links.go), which write directly to the same
// service_links table app-backend's servicelinks module reads from.
// app-backend has never exposed a create/edit endpoint for this model -
// every existing row was provisioned by a direct DB insert - so this is
// the only place these accounts can be managed from.

export interface IServiceLinkOwner {
  email: string;
  address: string;
  kycVerified: number;
  isMerchant: boolean;
  merchantOnline: boolean;
  corporate: number;
  membershipType: number;
  verified: number;
  suspended: number;
}

export interface IServiceLink {
  id: string;
  apiKey: string;
  createdAt: string;
  updatedAt: string;
  ownerUsername: string;
  address: string;
  shortName: string;
  longName: string;
  loginPermission: number;
  paymentPermission: number;
  tokenInfoPermission: number;
  authorizationPermission: number;
  eventPermission: number;
  allowUserInfo: number;
  pushNotificationPermission: number;
  includePhoneNumbers: number;
  includeUserBalances: number;
  tokenizedAssetAuthorizationPermission: number;
  createUsersPermission: number;
  allowReferralForRegisteredUsers: number;
  verified: number;
  inactive: number;
  suspended: number;
  suspensionReason: string | null;
  owner: IServiceLinkOwner | null;
}

export interface ServiceLinkListQueryParams {
  page?: number;
  pageSize?: number;
  ownerUsername?: string;
  shortName?: string;
  inactive?: "true" | "false";
  verified?: "true" | "false";
  suspended?: "true" | "false";
}

export interface ServiceLinkListResponse {
  message: string;
  data: {
    data: IServiceLink[];
    total: number;
    page: number;
    pageSize: number;
  };
  status: string;
  timestamp: string;
}

export interface ServiceLinkResponse {
  message: string;
  data: IServiceLink;
  status: string;
  timestamp: string;
}

// The permission flags are booleans in/out at this layer - tm-api stores
// them as 0/1 int columns (matching app-backend's own convention), the
// same int<->bool mapping CuratedAssetRequest uses for Withdrawable/
// Inactive.
export interface ServiceLinkRequest {
  action: "create" | "update";
  id?: string;
  ownerUsername: string;
  shortName: string;
  longName?: string;
  loginPermission?: boolean;
  paymentPermission?: boolean;
  tokenInfoPermission?: boolean;
  authorizationPermission?: boolean;
  eventPermission?: boolean;
  allowUserInfo?: boolean;
  pushNotificationPermission?: boolean;
  includePhoneNumbers?: boolean;
  includeUserBalances?: boolean;
  tokenizedAssetAuthorizationPermission?: boolean;
  createUsersPermission?: boolean;
  allowReferralForRegisteredUsers?: boolean;
  verified?: boolean;
  inactive?: boolean;
}

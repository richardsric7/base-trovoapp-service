export interface iUsers {
  id: string;
  user_name: string;
  created_at: number;
  updated_at: number;
  telegram_connected_at: number;
  last_name: string;
  first_name: string;
  middle_name: string;
  mobile: number;
  contact_phone: string;
  telegram: string;
  email: string;
  image_thumbnail: string;
  country_code: string;
  latitude: number;
  longitude: number;
  city: string;
  region: string;
  region_name: string;
  time_zone: string;
  isp: string;
  public_ip: number;
  address: string;
  bantu_talk: string;
  suspended: number;
  kyc_level: number;
  max_asset_per_offer: number;
  max_asset_per_order: number;
  suspension_reason: string;
  offline: number;
  admin_level: number;
  telegram_notifications: number;
}

export interface IUsersParams {
  page?: number;
  pageSize?: number;
  username?: string;
  email?: string;
  phone?: string;
  first_name?: string;
  last_name?: string;
  city?: string;
  search?: string;
}

export interface IUserInfo {
  account_recovery_enabled: boolean;
  account_recovery_expires_on: string | null;

  corporate: boolean;

  created_at: string;
  email: string;
  first_name: string;
  has_security_questions: boolean;
  id: string;

  kyc_verified: boolean;
  last_name: string;
  LastRecoveredAccountOn: string;
  LastUpdatedMobileOn: string;
  Latitude: number;
  Longitude: number;
  MembershipExpiry: string | null;
  MembershipType: string;
  mobile: string;
  mobile_verified: boolean;
  public_iP: string;
  address: string;
  primary_signer: string;
  push_notification_token: string;
  referral_link: string;
  referral_qr_code: string;
  referrer: string | null;

  suspended: number;
  suspension_reason: string | null;

  UpdatedAt: string;
  username: string;
  verified: boolean;
}

export interface IWallet {
  address: string;
  alias: string;
  signer: string;
  userId: string;
  primaryWallet: number;
}

export interface IFetchUserResponse {
  message: string;
  data: {
    user_info: IUserInfo;
    wallets: IWallet[];
  };
  status: string;
  timestamp: string;
}

export interface ISuspendUsers {
  email: string;
  suspension_note: string;
  suspension_reason_id: number;
}

export interface PaymentHistoryParams {
  page?: number;
  pageSize?: number;
  search?: string;
  from?: string;
}

export interface PaymentHistoryQueryParams {
  page: number;
  pageSize: number;
  search: string;
  from: string;
}

export interface Transaction {
  transactionDate: string;
  transactionType: string;
  from: string;
  fromAddress: string;
  to: string;
  toAddress: string;
  memo: string;
  contractAddress: string;
  assetCode: string;
  amount: string;
  transactionId: string;
}

export interface PaymentHistoryResponse {
  message: string;
  data: {
    data: Transaction[];
    total?: number;
  };
}
// Define the type for each deposit address object
export interface CryptoWalletDepositAddress {
  trovoWalletAddress: string;
  createdAt: string;
  currency: string;
  depositAddress: string;
  id: string;
  network: string;
  qrCode: string;
}

// Define the type for the inTrade field
export interface InTrade {
  buyingLiabilities: string;
  sellingLiabilities: string;
}

// Define the type for each claimed wallet balance entry
export interface WalletClaim {
  amount: string;
  assetCode: string;
  contractAddress: string;
  closedGroup: string;
  cryptoWalletDepositAddresses: CryptoWalletDepositAddress[];
  imageUrl: string;
  inTrade: InTrade;
  nativePrice: string;
  qrCode: string;
  usdPrice: string;
}

// Define the overall API response shape
export interface WalletBalancesResponse {
  claimed: WalletClaim[];
}

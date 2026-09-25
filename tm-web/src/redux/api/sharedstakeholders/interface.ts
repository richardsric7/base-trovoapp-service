export interface IStakeholderMember {
  id: string;
  email: string;
  role: string;
  status: string;
  first_name: string;
  last_name: string;
}

export interface IStakeholderOrganization {
  id: string;
  name: string;
  email: string;
  type: string;
  status: string;
}

export interface IStakeholderInfo {
  stakeholder_id: number;
  stakeholder_type: string;
  dashboard_role: string;
}

export interface IStakeholderWallet {
  is_wallet_linked: boolean;
  trovo_wallet_username?: string;
}

export interface IStakeholderProfileData {
  member: IStakeholderMember;
  organization: IStakeholderOrganization;
  stakeholder: IStakeholderInfo;
  wallet: IStakeholderWallet;
}

export interface IStakeholderProfileResponse {
  data: IStakeholderProfileData;
}

export interface IUpdateStakeholderProfileRequest {
  first_name: string;
  last_name: string;
}

export interface IUpdateNotificationPreferencesRequest {
  email_enabled: boolean;
  in_app_enabled: boolean;
}

export interface IStakeholderNotificationPreference {
  id: string;
  member_id: string;
  organization_id: string;
  email_enabled: boolean;
  in_app_enabled: boolean;
  created_at: string;
  updated_at: string;
}

export interface IStakeholderNotificationPreferenceResponse {
  data: IStakeholderNotificationPreference;
}

export interface IWalletLinkAuthorizationResponse {
  authId: string;
  dynamicLink?: string;
  qrCode?: string;
}

export interface IVerifyWalletLinkRequest {
  auth_id: string;
  wallet_username: string;
}

export interface IVerifyWalletLinkResponse {
  message: string;
}

export interface IStakeholderAssetActions {
  canApproveDueDiligence: boolean;
  canApproveFundRelease: boolean;
  canAuthorizeDistribution: boolean;
  canExecuteFundRelease: boolean;
  canUpdateAccount: boolean;
  canRequestFundRelease: boolean;
  canRecordRevenue: boolean;
  canSubmitValuation: boolean;
}

export interface IStakeholderAssetAssignment {
  id?: string;
  asset_id?: string;
  asset_code?: string;
  trustee_org_id?: string;
  trustee_stakeholder_id?: number;
  asset_manager_org_id?: string;
  asset_manager_stakeholder_id?: number;
  custodian_org_id?: string;
  custodian_stakeholder_id?: number;
  status?: string;
  assigned_by?: string;
  created_at?: string;
  updated_at?: string;
}

export interface IStakeholderAssetAssignedStakeholder {
  role: string;
  stakeholder_id: number;
  stakeholder_type: string;
  organization_id: string;
  organization_name: string;
  organization_email: string;
}

export interface IStakeholderAssetWallets {
  issuing_wallet_address?: string;
  issuing_wallet_alias?: string;
  market_making_wallet?: string;
  wallet_to_hold_assets_not_for_sale?: string;
}

export interface IStakeholderAssetDocument extends IStakeholderDocument {
  [key: string]: unknown;
  uploaded_by_org_id?: string;
  uploaded_by_member_id?: string;
  sha256?: string;
}

export interface IStakeholderAssetProfile {
  description: string;
  logo_url: string;
  website: string;
  offering_type: string;
  country: string;
  physical_address: string;
  longitude: string;
  latitude: string;
}

export interface IStakeholderAssetOwnership {
  type: string;
  kind: string;
  retained_or_contributed_value: number;
  cost_outside_valuation: number;
}

export interface IStakeholderAssetProtection {
  methods: string[];
  insurance_company_name: string;
  insurance_policy_number: string;
  insurance_policy_holder: string;
  insurance_coverage_percentage: number;
  free_from_liens_and_encumbrances: boolean;
  transfer_title_to_custodian: boolean;
  revenue_guarantees: boolean;
  performance_bond: boolean;
  service_level_agreement: boolean;
  other: string;
}

export interface IStakeholderAssetOffering {
  tokens_issued: number;
  tokens_for_sale: number;
  tokens_not_for_sale: number;
  maximum_tokens_for_sale: number;
  price_per_token: number;
  amount_to_be_raised: number;
  sales_start: string;
  sales_end: string;
  cap_quantity: number;
  cap_amount: number;
  cap_duration_days: number;
  payout_cycle: string;
  payout_currency: string;
  payout_type: number;
}

export interface IStakeholderAssetPreferredFee {
  id: number;
  description: string;
  fiat_percentage: number;
  fiat_cap: number;
  asset_percentage: number;
}

export interface IStakeholderAssetTokenizationFees {
  fiat_amount: number;
  asset_amount: number;
  asset_percentage: number;
  preferred: IStakeholderAssetPreferredFee;
}

export interface IStakeholderAssetReceivingAccount {
  bank_id: number;
  bank_name: string;
  bank_country: string;
  account_name: string;
  account_number: string;
}

export interface IStakeholderAssetApplicationPayment {
  id: number;
  payment_method_id: string;
  document_url: string;
  created_at: string;
}

export interface IStakeholderAssetApplicationFee {
  amount: number;
  asset: string;
  payments: IStakeholderAssetApplicationPayment[];
}

export interface IStakeholderAssetTokenizationDocument {
  id: number;
  document_type: string;
  title: string;
  url: string;
  public: boolean;
  created_at: string;
}

export interface IStakeholderAssetTimelineStage {
  status_id: number;
  description: string;
  reached: boolean;
  current: boolean;
}

export interface IStakeholderAssetTokenizationTimeline {
  current_status_id: number;
  date_submitted: string;
  date_approved?: string;
  stages: IStakeholderAssetTimelineStage[];
}

export interface IStakeholderAssetDetailData {
  risk_and_compliance?: {
    risk_assessment_score: number | null;
    maximum_risk_score: number;
    risk_level: string | null;
    compliance_status: string;
    last_audit_date: string | null;
  };
  id: string;
  assetCode: string;
  assetName: string;
  assetType: string;
  assetSector: string;
  assetSubSector: string;
  assetTokenizationStatus: number;
  assetQuoteCurrency: string;
  assetCurrentValue: number;
  valueOfTokenizedAsset: number;
  numberOfTokenToBeIssued: number;
  assetManagerId: number;
  approvedAssetCustodianId: number;
  trusteeId?: number;
  tokenHolderCount: number;
  complianceStatus: string;
  custodyStatus: string;
  createdAt: string;
  updatedAt: string;
  actions: IStakeholderAssetActions;
  assignment?: IStakeholderAssetAssignment;
  tokenizer_username?: string;
  assigned_stakeholders?: IStakeholderAssetAssignedStakeholder[];
  wallets?: IStakeholderAssetWallets;
  exempted_countries?: string[];
  documents?: IStakeholderAssetDocument[];
  asset_profile?: IStakeholderAssetProfile;
  ownership?: IStakeholderAssetOwnership;
  asset_protection?: IStakeholderAssetProtection;
  offering?: IStakeholderAssetOffering;
  tokenization_fees?: IStakeholderAssetTokenizationFees;
  receiving_account?: IStakeholderAssetReceivingAccount;
  application_fee?: IStakeholderAssetApplicationFee;
  tokenization_documents?: IStakeholderAssetTokenizationDocument[];
  tokenization_timeline?: IStakeholderAssetTokenizationTimeline;
}

export interface IStakeholderAssetDetailResponse {
  data: IStakeholderAssetDetailData;
}

export interface IStakeholderAssetsMeta {
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

export interface IStakeholderAssetsData {
  records: IStakeholderAssetDetailData[];
  meta: IStakeholderAssetsMeta;
  summary: IStakeholderAssetsSummary;
}

export interface IStakeholderAssetsSummary {
  total: number;
  active: number;
  liquidated: number;
}

export interface IStakeholderAssetsResponse {
  data: IStakeholderAssetsData;
}

export interface IStakeholderAssetsQueryParams {
  page?: number;
  limit?: number;
  search?: string;
  status?: string;
  asset_type?: string;
  asset_sector?: string;
}

export type StakeholderAuthorizationAction =
  | "fund_release.approve"
  | "fund_release.execute"
  | "distribution.authorize";

export type StakeholderAuthorizationEntity =
  | "fund_release_request"
  | "distribution";

export interface ICreateAuthorizationRequest {
  action: StakeholderAuthorizationAction;
  entity_type: StakeholderAuthorizationEntity;
  entity_id: string;
}

export interface IAuthorizationChallenge {
  id: string;
  auth_id: string;
  dynamic_link?: string;
  qr_code?: string;
  status: "pending" | "verified" | "expired" | "failed" | "used" | string;
  expires_at: string;
}

export interface IAuthorizationChallengeResponse {
  data: IAuthorizationChallenge;
}

export interface IVerifyAuthorizationRequest {
  challengeId: string;
  auth_id: string;
}

export type IAuditTrailStateValue =
  | string
  | number
  | boolean
  | null
  | IAuditTrailStateValue[]
  | { [key: string]: IAuditTrailStateValue };

export type IAuditTrailState = Record<string, IAuditTrailStateValue>;

export interface IStakeholderAuditTrailRecord {
  id: string;
  actor_member_id: string;
  actor_org_id: string;
  actor_role: string;
  actor_name?: string;
  actor_organization?: string;
  action: string;
  entity_type: string;
  entity_id: string;
  before_state?: IAuditTrailState;
  after_state?: IAuditTrailState;
  metadata?: IAuditTrailState;
  created_at: string;
}

export interface IStakeholderAssetActivitiesQueryParams {
  assetId: string;
  page?: number;
  limit?: number;
}

export interface IStakeholderAssetActivitiesResponse {
  data: {
    records: IStakeholderAuditTrailRecord[];
    meta: IStakeholderAuditTrailMeta;
  };
}

export interface IStakeholderAuditTrailMeta {
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

export interface IStakeholderAuditTrailResponse {
  data: {
    records: IStakeholderAuditTrailRecord[];
    meta: IStakeholderAuditTrailMeta;
  };
}

export interface IStakeholderNotification {
  id: string;
  recipient_org_id: string;
  sender_org_id?: string;
  type: string;
  title: string;
  message: string;
  related_entity_type?: string;
  related_entity_id?: string;
  read_at?: string | null;
  created_at: string;
}

export interface IStakeholderNotificationsResponse {
  data: {
    records: IStakeholderNotification[];
    meta: IStakeholderAuditTrailMeta;
  };
}

export interface IMarkStakeholderNotificationReadResponse {
  data: IStakeholderNotification;
}

export interface IStakeholderDocumentCategory {
  value: string;
  label: string;
  access_roles: string[];
}

export interface IStakeholderDocument {
  id: string;
  asset_id?: string;
  asset_code?: string;
  category: string;
  title: string;
  file_url?: string;
  original_filename?: string;
  mime_type?: string;
  size_bytes?: number;
  version: number;
  access_roles?: string[];
  status: string;
  created_at: string;
  updated_at: string;
}

export interface IStakeholderDocumentsResponse {
  data: {
    records: IStakeholderDocument[];
    meta: IStakeholderAssetsMeta;
  };
}

export interface IStakeholderDocumentResponse {
  data: IStakeholderDocument;
}

export interface IStakeholderDocumentListParams {
  page?: number;
  limit?: number;
  asset_id?: string;
}

export interface ICreateStakeholderDocumentRequest {
  asset_id?: string;
  asset_code?: string;
  category: string;
  title: string;
  file_url: string;
  mime_type?: string;
  version?: number;
  access_roles?: string[];
}

export interface IStakeholderDashboardValue {
  amount?: string;
  currency?: string;
}

export interface IStakeholderDashboardCategory {
  asset_count?: number;
  category?: string;
  values?: IStakeholderDashboardValue[];
}

export interface IStakeholderDashboardActivity {
  actor_name?: string;
  actor_organization?: string;
  created_at?: string;
  entity_id?: string;
  entity_type?: string;
  id?: string;
  message?: string;
  source?: string;
  status?: string;
  title?: string;
  type?: string;
}

export interface IStakeholderDashboardData {
  portfolio_summary?: {
    asset_count?: number;
    categories?: IStakeholderDashboardCategory[];
    values?: IStakeholderDashboardValue[];
  };
  recent_activity?: IStakeholderDashboardActivity[];
  recent_assets?: IStakeholderAssetDetailData[];
  role?: string;
  summary?: Record<string, unknown>;
}

export interface IStakeholderDashboardResponse {
  data: IStakeholderDashboardData;
}

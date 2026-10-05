// Public Markets (tm-api /public-markets): tokenized NGX equities and FMDQ
// bonds. Amounts and quantities are decimal strings.

export interface Envelope<T> {
  message: string;
  data: T;
}

export type AssetStatus = "SETUP" | "LIVE" | "HALTED";

export interface IPMAsset {
  id: string;
  createdAt: string;
  assetCode: string;
  ticker: string;
  market: "NGX" | "FMDQ";
  assetType: "EQUITY" | "BOND";
  isin: string;
  instrumentName: string;
  shortName: string;
  sector: string;
  description: string;
  unitDescription: string;
  logoBackground: string;
  logoForeground: string;
  logoInitials: string;
  logoUrl: string;
  marketCap: string;
  peRatio: string;
  dividendYield: string;
  coupon: string;
  maturityDate?: string | null;
  contractAddress: string;
  issuingSafeAddress: string;
  tokenDecimals: number;
  custodianId: number;
  dealingMemberId: number;
  omnibusReference: string;
  minimumBuy: string;
  feePercent: string;
  inventoryTargetUnits: string;
  status: AssetStatus;
  haltReason: string;
  haltedAt?: string | null;
  haltedBy: string;
  lastPrice: string;
  previousClose: string;
  priceSource: string;
  priceAt?: string | null;
  createdBy: string;
}

export interface IAssetRow extends IPMAsset {
  supply: string;
  owners: number;
  position: string;
  positionAsOf?: string | null;
  positionSource: string;
  custodianName: string;
  dealingMemberName: string;
  marketValue: string;
  lastReconResult: string;
  lastReconAt?: string | null;
  priceFresh: boolean;
  marketOpen: boolean;
}

export interface IHolder {
  walletAddress: string;
  channel: string;
  serviceLinkId: string;
  exchangeName: string;
  balance: string;
  percentOfSupply: string;
  substantial: boolean;
}

export interface IReconRun {
  id: number;
  createdAt: string;
  assetId: string;
  assetCode: string;
  tokenSupply: string;
  custodianPosition: string;
  positionAsOf?: string | null;
  positionSource: string;
  ledgerTotal: string;
  inventory: string;
  delta: string;
  ledgerDelta: string;
  priceStale: boolean;
  result: "MATCHED" | "DRIFT" | "LEDGER_DRIFT" | "NO_POSITION" | "ERROR";
  halted: boolean;
  detail: string;
}

export interface IPriceSnapshot {
  id: number;
  assetId: string;
  price: string;
  source: string;
  marketHours: boolean;
  asOf: string;
  capturedAt: string;
}

export interface ICorporateAction {
  id: string;
  createdAt: string;
  assetId: string;
  assetCode: string;
  eventType: string;
  description: string;
  recordDate: string;
  payDate: string;
  amountPerUnit: string;
  currency: string;
  source: string;
  sourceReference: string;
  status: string;
  declaredBy: string;
  recordBlock: number;
  snapshotAt?: string | null;
  snapshotChecksum: string;
  holderCount: number;
  eligibleUnits: string;
  retainedUnits: string;
  grossAmount: string;
  whtAmount: string;
  netAmount: string;
  approvalsRequired: number;
  approvedAt?: string | null;
  fundedAt?: string | null;
  whtTxHash: string;
  completedAt?: string | null;
  note: string;
}

export interface ICustodianIntegration {
  custodianId: number;
  code: string;
  nomineeName: string;
  mode: "MOCK" | "REST" | "MANUAL";
  transport: string;
  authScheme: string;
  baseUrl: string;
  credentialsRef: string;
  feePercent: string;
  active: boolean;
}

export interface IDealingMember {
  id: number;
  dealingMemberName: string;
  dealingMemberAddress: string;
  dealingMemberCountry: string;
  cscsMemberCode: string;
  requirementDocument: string;
  feePercent: number;
  feeFixed: number;
  code: string;
  mode: "MOCK" | "REST" | "MANUAL";
  authScheme: string;
  baseUrl: string;
  credentialsRef: string;
  active: boolean;
  assets?: number;
}

export interface ISetupStep {
  label: string;
  detail: string;
  done: boolean;
}

export interface IAssetDetail extends IAssetRow {
  custodian?: ICustodianIntegration | null;
  dealingMember?: IDealingMember | null;
  holders: IHolder[];
  reconciliationRuns: IReconRun[];
  prices: IPriceSnapshot[];
  corporateActions: ICorporateAction[];
  setupSteps: ISetupStep[];
  openOrders: number;
}

export interface IJobRun {
  job: string;
  lastRunAt: string;
  status: "up" | "degraded" | "down";
  detail: string;
}

export interface IJobRequest {
  id: number;
  createdAt: string;
  job: string;
  target: string;
  requestedBy: string;
  doneAt?: string | null;
  result: string;
}

export interface INetBatch {
  id: string;
  createdAt: string;
  updatedAt: string;
  assetId: string;
  assetCode: string;
  sessionDate: string;
  side: "BUY" | "SELL" | "NONE";
  quantity: string;
  referencePrice: string;
  value: string;
  creationOrders: number;
  redemptionOrders: number;
  status: string;
  approvalsRequired: number;
  executedQuantity: string;
  executedPrice: string;
  executedAt?: string | null;
  settledAt?: string | null;
  note: string;
}

export interface IBatchView extends INetBatch {
  approvals: { at: string; batchId: string; approver: string }[];
  threshold: string;
}

export interface IAttentionItem {
  kind: string;
  text: string;
  where: string;
  target: string;
}

export interface IOverview {
  assets: number;
  liveAssets: number;
  haltedAssets: number;
  marketValue: string;
  settledToday: number;
  settledTodayValue: string;
  inProgress: number;
  inProgressValue: string;
  creationsToday: number;
  redemptionsToday: number;
  fastPercent: string;
  nettedPercent: string;
  halted: IPMAsset[];
  attention: IAttentionItem[];
  jobs: IJobRun[];
  batchesToday: INetBatch[];
}

export interface IPMOrder {
  id: string;
  createdAt: string;
  type: "CREATION" | "REDEMPTION";
  assetId: string;
  assetCode: string;
  channel: "TROVO_APP" | "EXCHANGE";
  serviceLinkId: string;
  externalOrderRef: string;
  username: string;
  walletAddress: string;
  partnerWalletId: string;
  fundingAssetCode: string;
  amount: string;
  fee: string;
  netAmount: string;
  quantity: string;
  referencePrice: string;
  executedPrice: string;
  path: "FAST" | "SLOW" | "NETTED" | "";
  state: string;
  note: string;
  batchId: string;
  paymentTxHash: string;
  tokenTxHash: string;
  payoutTxHash: string;
  settledAt?: string | null;
  completedAt?: string | null;
}

export interface IInstruction {
  id: string;
  createdAt: string;
  updatedAt: string;
  kind: "CUSTODIAN_CREATION" | "CUSTODIAN_REDEMPTION" | "DM_ORDER";
  partnerType: "CUSTODIAN" | "DEALING_MEMBER";
  partnerId: number;
  partnerName: string;
  assetCode: string;
  batchId: string;
  side: string;
  quantity: string;
  status: string;
  attempts: number;
  maxAttempts: number;
  nextAttemptAt?: string | null;
  lastError: string;
  lastStatusCode: number;
  partnerRef: string;
  handledBy: string;
}

export interface IOrderDetail {
  order: IPMOrder;
  lastError: string;
  attempts: number;
  events: { at: string; state: string; note: string }[];
  batch?: INetBatch | null;
  instructions: IInstruction[];
  exchangeName: string;
}

export interface IPartnerEvent {
  id: number;
  createdAt: string;
  source: string;
  kind: string;
  reference: string;
  payload: string;
  processedAt?: string | null;
  result: string;
}

export interface IReconRow extends IReconRun {
  custodianName: string;
  assetStatus: string;
  ordersSince: IPMOrder[];
}

export interface IReconciliation {
  checked: number;
  matched: number;
  drift: number;
  stalePrices: number;
  lastRunAt?: string | null;
  nextRunAt: string;
  rows: IReconRow[];
}

export interface IEntitlement {
  id: string;
  walletAddress: string;
  channel: string;
  serviceLinkId: string;
  partnerWalletId: string;
  username: string;
  units: string;
  grossAmount: string;
  whtPercent: string;
  whtAmount: string;
  netAmount: string;
  taxResidency: string;
  status: string;
  txHash: string;
  paidAt?: string | null;
  note: string;
}

export interface IActionDetail extends ICorporateAction {
  approvals: { at: string; approver: string; checksum: string }[];
  approvers: string[];
  byChannel: { channel: string; recipients: number; gross: string; wht: string; net: string; paid: number; failed: number }[];
  entitlements: IEntitlement[];
  totalEntitlements: number;
}

export interface IConfirmationRow {
  serviceLinkId: string;
  exchangeName: string;
  assetCode: string;
  walletsPaid: number;
  confirmed: number;
  outstanding: number;
  escalated: number;
  oldestDueAt?: string | null;
}

export interface IExchangeRow {
  serviceLinkId: string;
  createdAt: string;
  status: "active" | "suspended";
  environment: "sandbox" | "production";
  callbackUrl: string;
  rateLimitTier: string;
  revenueShareTier: string;
  confirmationSlaHours: number;
  fundingAddress: string;
  balance: string;
  techContact: string;
  name: string;
  shortName: string;
  rateLimitPerMinute: number;
  wallets: number;
  orders30d: number;
  deadLetters: number;
  pendingWebhooks: number;
  previousSecretValidUntil?: string | null;
}

export interface IWebhookDelivery {
  id: string;
  createdAt: string;
  event: string;
  reference: string;
  walletId: string;
  assetCode: string;
  payload: string;
  status: "PENDING" | "DELIVERED" | "DEAD_LETTER";
  attempts: number;
  lastResponseCode: number;
  lastError: string;
  deliveredAt?: string | null;
  needsConfirmation: boolean;
  confirmedAt?: string | null;
  escalatedAt?: string | null;
}

export interface IExchangeLedgerEntry {
  id: number;
  createdAt: string;
  kind: string;
  amount: string;
  balanceAfter: string;
  reference: string;
  note: string;
  createdBy: string;
}

export interface IExchangeDetail extends IExchangeRow {
  deliveries: IWebhookDelivery[];
  ledger: IExchangeLedgerEntry[];
  recentOrders: IPMOrder[];
  withdrawals: IJobRequest[];
}

export interface IPartnerWallet {
  walletId: string;
  createdAt: string;
  serviceLinkId: string;
  externalUserRef: string;
  legalName: string;
  taxIdentifier: string;
  residencyCountry: string;
  nationality: string;
  ndpaConsent: boolean;
  walletAddress: string;
  deployed: boolean;
  status: "active" | "rejected" | "suspended";
  rejectionReason: string;
  exchangeName: string;
}

export interface IWalletStats {
  provisioned: number;
  consentPercent: string;
  taxDataPercent: string;
  rejected30d: number;
  deployedWallets: number;
}

export interface IPriceRow {
  assetId: string;
  assetCode: string;
  market: string;
  lastPrice: string;
  previousClose: string;
  source: string;
  capturedAt?: string | null;
  marketOpen: boolean;
  fresh: boolean;
  ageMinutes: number;
}

export interface IExecution {
  instructionId: string;
  dealingMember: string;
  assetCode: string;
  side: string;
  quantity: string;
  executedPrice: string;
  referencePrice: string;
  deviationPercent: string;
  executedAt?: string | null;
}

export interface IPMSettings {
  updatedAt: string;
  updatedBy: string;
  tradeFeePercent: string;
  feeWallet: string;
  fundingAssetCode: string;
  netCreationThreshold: string;
  approvalsRequired: number;
  netCreationApprovers: string;
  dividendApprovers: string;
  priceStaleMinutes: number;
  instructionMaxAttempts: number;
  settlementSlaHours: number;
  confirmationSlaHours: number;
  webhookMaxAttempts: number;
  batchIntervalMinutes: number;
  reconciliationHour: number;
  whtResidentPercent: string;
  whtNonResidentPercent: string;
  whtMissingTaxIdPercent: string;
  whtWallet: string;
  substantialHoldingPercent: string;
  rateLimitTiers: string;
  aumFeePercent: string;
  fxSpreadPercent: string;
  revenueShareTiers: string;
  ngxOpen: string;
  ngxClose: string;
  fmdqOpen: string;
  fmdqClose: string;
  marketHolidays: string;
}

export interface ICustodianRow {
  id: number;
  name: string;
  country: string;
  feePercent: number;
  integration?: ICustodianIntegration | null;
  assets: number;
}

export interface IHealth {
  jobs: IJobRun[];
  requests: IJobRequest[];
  pendingEvents: number;
  pendingMockEvents: number;
  pendingWebhooks: number;
}

export interface IPMPermissions {
  email: string;
  manage: boolean;
  settings: boolean;
  personalData: boolean;
  netCreationApprover: boolean;
  dividendApprover: boolean;
}

export type AssetInput = Partial<{
  assetCode: string;
  ticker: string;
  market: string;
  assetType: string;
  isin: string;
  instrumentName: string;
  shortName: string;
  sector: string;
  description: string;
  unitDescription: string;
  custodianId: number;
  dealingMemberId: number;
  omnibusReference: string;
  minimumBuy: string;
  feePercent: string;
  inventoryTargetUnits: string;
  logoBackground: string;
  logoForeground: string;
  logoInitials: string;
  marketCap: string;
  peRatio: string;
  dividendYield: string;
  coupon: string;
  maturityDate: string;
}>;

export type ExchangeInput = Partial<{
  serviceLinkId: string;
  callbackUrl: string;
  environment: string;
  rateLimitTier: string;
  revenueShareTier: string;
  fundingAddress: string;
  techContact: string;
  confirmationSlaHours: number;
}>;

export type PartnerInput = Partial<{
  name: string;
  address: string;
  country: string;
  cscsMemberCode: string;
  requirementDocument: string;
  feePercent: number;
  feeFixed: number;
  code: string;
  nomineeName: string;
  mode: string;
  transport: string;
  authScheme: string;
  baseUrl: string;
  credentialsRef: string;
  active: boolean;
}>;

export type Paged = { page?: number; limit?: number };

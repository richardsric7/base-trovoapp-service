// These mirror tm-api's /proceed-payouts endpoints: Trovo admins drive
// payout-engine's stages for each authorized stakeholder distribution
// (prepare -> approve -> confirm funding -> paying -> completed) and read the
// payout, fee and VAT reports. Amounts ending in "Units" are token base units;
// the others are human amounts of the payout token.

export type PayoutStatus =
  | "REGISTERED"
  | "PREPARE_REQUESTED"
  | "PREPARING"
  | "LOCKED"
  | "APPROVED"
  | "FUNDING_CHECK_REQUESTED"
  | "PAYING"
  | "PAUSED"
  | "COMPLETED"
  | "COMPLETED_WITH_FAILURES"
  | "CANCELLED";

export type PayoutItemStatus = "PENDING" | "QUEUED" | "PAID" | "FAILED" | "EXCLUDED" | "SKIPPED";

export interface IPayout {
  id: number;
  createdAt: string;
  updatedAt: string;
  tokenizedAssetId: string;
  assetCode: string;
  assetName: string;
  batch: string;
  distributionId: string;
  status: PayoutStatus;
  payoutAssetCode: string;
  payoutContractAddress: string;
  payoutSafeAddress: string;
  tokenContractAddress: string;
  totalAmount: string;
  amountPerToken: string;
  holderPayable: string;
  fee: string;
  vat: string;
  paid: string;
  payable: string;
  retained: string;
  snapshotBlock: number;
  scannedBlock: number;
  holderCount: number;
  paidCount: number;
  failedCount: number;
  excludedCount: number;
  scheduleChecksum: string;
  approvalsRequired: number;
  approvals: number;
  preparedBy: string;
  lockedAt?: string;
  approvedAt?: string;
  fundingRequestedBy: string;
  startedAt?: string;
  completedAt?: string;
  note: string;
  feeType: "FIXED" | "PERCENT";
  feeValue: string;
  feeCap: string;
  feeSetBy: string;
  vatPercent: string;
  feeWallet: string;
  vatWallet: string;
}

export interface IPayoutApproval {
  id: number;
  createdAt: string;
  adminEmail: string;
  scheduleChecksum: string;
}

export interface IPayoutBatch {
  id: number;
  createdAt: string;
  status: "SUBMITTED" | "MINED" | "REVERTED" | "DROPPED";
  safeNonce: number;
  txHash: string;
  executor: string;
  itemCount: number;
  amountUnits: string;
  gasUsed: number;
  error: string;
}

export interface IEngineState {
  halted: boolean;
  haltReason: string;
  haltedBy: string;
  haltedAt?: string;
  heartbeatAt?: string;
  instance: string;
  version: string;
  activity: string;
  lastError: string;
  sweepToken: string;
  sweepRequestedBy: string;
  sweepResult: string;
  online: boolean;
}

export interface IPayoutDetail extends IPayout {
  approvalList: IPayoutApproval[];
  batches: IPayoutBatch[];
  itemCounts: Partial<Record<PayoutItemStatus, number>>;
  engine?: IEngineState;
}

export interface IPayoutItem {
  id: string;
  beneficiaryAddress: string;
  username: string;
  kind: "HOLDER" | "FEE" | "VAT";
  confirmedTokenizedAssetBalance: number;
  amountToReceive: number;
  status: PayoutItemStatus;
  reason: string;
  txHash: string;
  paidAt?: string;
  actionBy: string;
}

export interface IFeeConfig {
  feeWallet: string;
  feeType: "FIXED" | "PERCENT";
  feeValue: string;
  feeCap: string;
  vatWallet: string;
  lastUpdatedBy: string;
}

export interface ICurrencyTotals {
  currency: string;
  payouts: number;
  total: string;
  holderPayable: string;
  paidToHolders: string;
  fees: string;
  vat: string;
  holdersPaid: number;
  holdersFailed: number;
}

export interface IPayoutsReport {
  payouts: IPayout[];
  totals: ICurrencyTotals[] | null;
  byStatus: Record<string, number>;
}

export interface IFeeCollection {
  id: string;
  createdAt: string;
  assetCode: string;
  payoutBatch: string;
  feeType: "PROCEED_PAYOUT_FEE" | "PROCEED_PAYOUT_FEE_VAT";
  amount: number;
  payoutAssetCode: string;
  destinationWallet: string;
  transactionHash?: string;
}

export interface IFeesReport {
  collections: IFeeCollection[] | null;
  totals: { assetCode: string; payoutAssetCode: string; fees: number; vat: number; count: number }[] | null;
  totalFees: Record<string, number>;
  totalVat: Record<string, number>;
}

export interface Envelope<T> {
  message: string;
  data: T;
}

export interface PayoutFilters {
  status?: string;
  asset?: string;
  from?: string;
  to?: string;
  page?: number;
  limit?: number;
}

export interface ItemFilters {
  id: number;
  status?: string;
  kind?: string;
  search?: string;
  page?: number;
  limit?: number;
}

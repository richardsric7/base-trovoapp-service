export type P2PPaymentMethod = {
  paymentChannel: string;
  provider: string;
  account: string;
};

export type P2POffer = {
  id: string;
  merchantUsername: string;
  merchantUserId: string;
  offerType: 'BUY' | 'SELL';
  asset: string;
  contractAddress: string;
  paymentMethod: P2PPaymentMethod;
  country: string;
  countryCode: string;
  currency: string;
  priceType: string;
  price: string;
  priceMargin: string;
  minOrderAmount: string;
  maxOrderAmount: string;
  availableLiquidity: string;
  reservedLiquidity: string;
  remark: string;
  availabilityStatus: 'ONLINE' | 'OFFLINE';
  status: string;
  version: number;
  createdAt: string;
  updatedAt: string;
};

export type P2POrder = {
  id: string;
  offerId: string;
  customerUserId: string;
  customerUsername: string;
  merchantUserId: string;
  merchantUsername: string;
  offerType: 'BUY' | 'SELL';
  asset: string;
  assetContractAddress: string;
  paymentMethodSnapshot: P2PPaymentMethod;
  country: string;
  countryCode: string;
  currency: string;
  price: string;
  specifiedAssetAmount: string;
  paymentAmount: string;

  buyerTotalFees: string;
  buyerTotalVat: string;
  buyerTotalCharges: string;
  buyerNetAssetAmount: string;

  sellerTotalFees: string;
  sellerTotalVat: string;
  sellerTotalCharges: string;
  sellerEscrowAssetAmount: string;

  combinedPlatformFee: string;
  combinedRegulatoryFee: string;
  combinedVat: string;

  escrowDepositShortlink: string;
  escrowDepositQrCode: string;
  escrowDepositTransactionHash: string;
  assetReleaseTransactionHash: string;

  escrowDepositStatus: string;
  expectedEscrowAmount: string;
  depositedEscrowAmount: string;
  refundableAmount: string;

  assetDepositor: string;
  assetRecipient: string;
  fiatPayer: string;
  fiatRecipient: string;

  orderStatus:
    | 'AWAITING_APPROVAL'
    | 'AWAITING_ESCROW_DEPOSIT'
    | 'AWAITING_PAYMENT'
    | 'AWAITING_PAYMENT_CONFIRMATION'
    | 'COMPLETED'
    | 'REJECTED'
    | 'CANCELLED'
    | 'EXPIRED';
  isDisputed: boolean;
  expiresAt: string | null;
  completedAt: string | null;
  paymentConfirmedAt: string | null;
  assetReleasedAt: string | null;
  createdAt: string;
  updatedAt: string;
};

export type P2PDispute = {
  id: string;
  orderId: string;
  openedBy: string;
  openedAt: string;
  subject: string;
  description: string;
  evidence: string[];
  status: 'OPEN' | 'RESOLVED';
  resolution: string;
  resolvedAt: string | null;
  resolvedBy: string;
};

export type P2PRefund = {
  id: string;
  orderId: string;
  sender: string;
  token: string;
  contractAddress: string;
  amount: string;
  reason: string;
  claimed: boolean;
  claimedAt: string | null;
  transactionHash: string;
  createdAt: string;
};

export type P2PMerchantPerformance = {
  merchantId: string;
  merchantUsername: string;
  completedTrades: number;
  completedTradeVolume: string;
  completionRate: string;
  averageOrderCompletionTime: number;
  cancelledOrders: number;
  expiredOrders: number;
  disputesOpened: number;
  disputesResolved: number;
  disputesResolvedAgainstMerchant: number;
  disputesResolvedInFavorOfMerchant: number;
  lastActivityAt: string | null;
};

export type P2PCustomerPerformance = {
  customerId: string;
  customerUsername: string;
  completedTrades: number;
  completedTradeVolume: string;
  completionRate: string;
  cancelledOrders: number;
  cancelledAfterAcceptance: number;
  expiredOrders: number;
  disputesOpened: number;
  disputesResolved: number;
  lastActivityAt: string | null;
};

export type P2POrderFeeQuote = {
  offerId: string;
  specifiedAssetAmount: string;
  paymentAmount: string;
  buyerTotalFees: string;
  buyerTotalVat: string;
  buyerTotalCharges: string;
  buyerNetAssetAmount: string;
  sellerTotalFees: string;
  sellerTotalVat: string;
  sellerTotalCharges: string;
  sellerEscrowAssetAmount: string;
};

// The role-check helpers below mirror the backend's fixed offerType-based
// role mapping (Plan Section 11), the same way app-mobile's P2POrder does -
// the client only knows its own username, not the opaque
// customerUserId/merchantUserId/assetDepositor/etc. ids the API returns.
export function isCustomer(order: P2POrder, myUsername: string): boolean {
  return order.customerUsername === myUsername;
}

export function isMerchant(order: P2POrder, myUsername: string): boolean {
  return order.merchantUsername === myUsername;
}

export function isAssetDepositor(order: P2POrder, myUsername: string): boolean {
  const iAmCustomer = isCustomer(order, myUsername);
  return order.offerType === 'SELL' ? !iAmCustomer : iAmCustomer;
}

export function isFiatPayer(order: P2POrder, myUsername: string): boolean {
  const iAmCustomer = isCustomer(order, myUsername);
  return order.offerType === 'SELL' ? iAmCustomer : !iAmCustomer;
}

export function isFiatRecipient(order: P2POrder, myUsername: string): boolean {
  return !isFiatPayer(order, myUsername);
}

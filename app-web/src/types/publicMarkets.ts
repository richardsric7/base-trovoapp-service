// Public Markets: tokenized NGX equities and FMDQ bonds (app-backend
// /v1/public-markets). Amounts, prices and quantities are decimal strings.

export type PMSession = {
  open: boolean;
  market: string;
  opensAt?: string;
  closesAt?: string;
  minutesRemaining?: number;
  minutesTotal?: number;
};

export type PMCustody = {
  custodian?: string;
  nominee?: string;
  dealingMember?: string;
  depository?: string;
  settlement?: string;
  unitsHeld?: string;
  positionAsOf?: string;
  lastReconciliation?: { at?: string; result?: string };
};

export type PMCorporateAction = {
  id: string;
  assetCode: string;
  eventType: string;
  description: string;
  recordDate: string;
  payDate: string;
  amountPerUnit: string;
  currency: string;
  status: string;
};

export type PMAsset = {
  id: string;
  assetCode: string;
  ticker: string;
  name: string;
  shortName: string;
  market: string;
  assetType: string; // EQUITY | BOND | ETF
  sector: string;
  isin: string;
  description: string;
  unitDescription: string;
  logo: { background?: string; foreground?: string; initials?: string; url?: string };
  price: string;
  priceSource: string;
  priceAt?: string;
  priceLive: boolean;
  dayChangePercent: string;
  previousClose: string;
  dayHigh: string;
  dayLow: string;
  dayVolume: string;
  marketCap: string;
  peRatio: string;
  dividendYield: string;
  coupon: string;
  maturityDate?: string;
  status: 'open' | 'halted' | 'coming-soon';
  contractAddress: string;
  tokenDecimals: number;
  minimumBuy: string;
  feePercent: string;
  fundingAsset: string;
  currency: string;
  session: PMSession;
  custody?: PMCustody;
  corporateActions?: PMCorporateAction[];
  tokensInCirculation?: string;
};

export type PMPriceRange = '1D' | '1W' | '1M' | '3M' | '1Y' | 'All';

export type PMPrices = { points: { at: string; price: string }[]; previousClose: string };

export type PMQuote = {
  assetCode: string;
  side: 'buy' | 'sell';
  amount: string;
  fee: string;
  feePercent: string;
  netAmount: string;
  quantity: string;
  price: string;
  priceSource: string;
  currency: string;
  fundingAsset: string;
  marketOpen: boolean;
  path: string; // FAST | SLOW | NETTED
  note: string;
  nextSessionAt?: string;
  minimumBuy: string;
  custodianName: string;
  settlementNote: string;
};

export type PMOrder = {
  id: string;
  type: 'CREATION' | 'REDEMPTION';
  assetCode: string;
  state: string;
  path: string;
  amount: string;
  fee: string;
  netAmount: string;
  quantity: string;
  referencePrice: string;
  executedPrice: string;
  fundingAssetCode: string;
  walletAddress: string;
  walletAlias?: string;
  note: string;
  createdAt: string;
  tokenTxHash?: string;
  payoutTxHash?: string;
  paymentTxHash?: string;
  events?: { at: string; state: string; note: string }[];
};

export type PMHolding = {
  asset: PMAsset;
  quantity: string;
  marketValue: string;
  averageCost: string;
  costBasis: string;
  totalReturn: string;
  returnPercent: string;
  todayChange: string;
  incomeReceived: string;
  wallets: Record<string, string>; // wallet address -> quantity
};

export type PMActivity = {
  kind: string;
  title: string;
  detail: string;
  received: string;
  spent: string;
  reference: string;
  at: string;
};

export type PMPortfolio = {
  value: string;
  costBasis: string;
  totalReturn: string;
  returnPercent: string;
  todayChange: string;
  income: string;
  holdings: PMHolding[];
  openOrders: PMOrder[];
  activity: PMActivity[];
};

export type PMDividend = {
  id: string;
  assetCode: string;
  eventType: string;
  description: string;
  recordDate: string;
  payDate: string;
  amountPerUnit: string;
  custodianName: string;
  units: string;
  grossAmount: string;
  whtPercent: string;
  whtAmount: string;
  netAmount: string;
  status: string;
  walletAddress: string;
  txHash?: string;
  paidAt?: string;
};

export const PM_FINAL_STATES = ['complete', 'rejected', 'failed', 'cancelled'];

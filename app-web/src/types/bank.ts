// Bank (NGN) deposits and withdrawals through Stablerail - see app-backend's
// /v1/users/stablerail/... routes.
export type StablerailProfile = {
  enabled: boolean;
  onboarded: boolean;
  onboardingStatus: string;
  minimumWithdrawal: string;
};

export type StablerailBank = {
  bank_code: string;
  bank_name: string;
};

export type BankWithdrawal = {
  id: string;
  requestId: string;
  createdAt: string;
  bankCode: string;
  bankName: string;
  accountNumber: string;
  amount: number;
  status: string;
  walletAddress: string;
  txHash: string;
  error?: string;
};

export type VirtualAccount = {
  accountNumber: string;
  bankName: string;
  accountName: string;
  amount: number;
};

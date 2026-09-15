type HistoryRow = {
  id: string;
  type: string;
  assetCode: string;
  amountValue: number;
  price: string;
  description: string;
  dateLabel: string;
  dateValue: Date | null;
  walletKey: string;
  walletLabel: string;
  username: string;
  fromUsername?: string;
  toUsername?: string;
  fromPublicKey: string;
  toPublicKey: string;
  memo: string;
};

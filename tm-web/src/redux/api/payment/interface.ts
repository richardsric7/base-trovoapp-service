export interface FiatPaymentRecordsResponse {
  message: string;
  data: FiatPaymentRecordsPayload;
  timestamp: string;
  status: "OK" | "ERROR";
}

export interface GetFiatPaymentsParams {
  record_type?: "payment" | "invoice";
  service_provider?: string;
  username?: string;
  payment_type?: string;
  status?: "PENDING" | "COMPLETED";
  page?: number;
  page_size?: number;
  created_at?: string;
}

export interface FiatPaymentRecordsPayload {
  data: FiatPaymentRecord[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export type FiatPaymentRecord =
  | FiatPaymentRecordPayment
  | FiatPaymentRecordInvoice;

export interface FiatPaymentRecordPayment {
  id: string;
  record_type: "payment";
  service_provider: string;
  username: string;
  transaction_id: string;
  amount: string;
  payment_type: "ACTIVATION" | string;
  created_at: string;
}

export interface FiatPaymentRecordInvoice {
  id: string;
  record_type: "invoice";
  service_provider: "flutterwave" | string;
  username: string;
  amount: string;
  payment_type: "ACTIVATION" | string;
  status: "PENDING" | "COMPLETED" | string;
  refunded: number;
  created_at: string;
}

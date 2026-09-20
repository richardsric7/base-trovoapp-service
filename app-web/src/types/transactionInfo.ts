export type TransactionInfo= {
    transactionDate: Date,
    transactionType?: string,
    transactionDirection?: TransactionDirection,
    from: string,
    fromAddress: string,
    to: string,
    toAddress: string,
    memo: string,
    assetCode: string,
    assetIssuer: string,
    amount: string,
    transactionId: string,
}

export enum TransactionDirection {
    Send,
    Receive,
    Swap,
    Deposit,
    Withdraw,
}
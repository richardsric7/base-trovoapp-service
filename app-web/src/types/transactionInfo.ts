export type TransactionInfo= {
    transactionDate: Date,
    transactionType?: string,
    transactionDirection?: TransactionDirection,
    from: string,
    fromPublicKey: string,
    to: string,
    toPublicKey: string,
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
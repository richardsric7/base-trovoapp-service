export type Asset = {
    assetIssuer: string,
    assetCode: string,
    amount: number,
    inTrade: {
        sellingLiabilities: number,
        buyingLiabilities: number,
    },
    qrCode: string,
    imageUrl: string,
    usdPrice: number,
    nativePrice: number,
    cryptoWalletDepositAddresses: string,
    closedGroup: string,
}
  
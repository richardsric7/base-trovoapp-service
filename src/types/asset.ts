export type Asset = {
    assetIssuer: string,
    assetCode: string,
    amount: string,
    inTrade: {
        sellingLiabilities: string,
        buyingLiabilities: string,
    },
    qrCode: string,
    imageUrl: string,
    usdPrice: string,
    nativePrice: string,
    cryptoWalletDepositAddresses: string,
    closedGroup: string,
}
  
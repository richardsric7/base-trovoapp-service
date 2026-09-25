export type Asset = {
    contractAddress: string,
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
    tokenizedAsset: boolean,
    assetClassId : number,
    fundingStructure : number,
    exitWithFiat : number,
}
  
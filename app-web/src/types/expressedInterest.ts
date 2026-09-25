import { TokenizedAsset } from "./tokenizedAsset";

export interface ExpressedInterest {
    createdAt:          Date;
    updatedAt:          Date;
    tokenizedAssetId:   string;
    tokenizedAssetInfo: TokenizedAsset;
    assetCode:          string;
    contractAddress:        string;
    amount:             number;
    price:              number;
    subscriberUsername: string;
}
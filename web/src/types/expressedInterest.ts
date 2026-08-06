import { TokenizedAsset } from "./tokenizedAsset";

export interface ExpressedInterest {
    createdAt:          Date;
    updatedAt:          Date;
    tokenizedAssetId:   string;
    tokenizedAssetInfo: TokenizedAsset;
    assetCode:          string;
    assetIssuer:        string;
    amount:             number;
    price:              number;
    subscriberUsername: string;
}
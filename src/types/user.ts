import { Asset } from "./asset";
import { CuratedAsset } from "./curatedAsset";
import { DefaultAsset } from "./defaultAsset";
import { PatronMembership } from "./patronMembership";
import { Wallet } from "./wallet";

export type User = {
    username: string,
    email: string,
    imageThumbnailURL: string,
    firstName: string,
    lastName: string,
    mobile: string,
    publicKey: string,
    primarySigner: string,
    referrer: string,
    referralLink: string,
    referralQrCode: string,
    pushNotificationToken: string,
    corporate: number,
    mobileVerified: boolean,
    countryCode: string,
    membershipType: number,
    membershipExpiry: Date,
    kycVerified: boolean,
    accountRecoveryEnabled: boolean,
    userWallets: Wallet[],
    verified: boolean,
    suspended: boolean,
    hasSecurityQuestions: boolean,
    curatedSwapList: CuratedAsset[],
    patronMembership?: PatronMembership,
    downlines: {
        level1: string,
        level2: string,
        level3: string,
    },
    uplines: {
        level1: string,
        level2: string,
        level3: string,
    },
    defaultAssets: DefaultAsset[],
    tokenizedAssets: Asset[],
    secretKeys: string[],
    isLoggedIn: boolean,
    password: string,
    currency: string,
}

export enum FieldState {
    error,
    pristine,
    ok,
}
  
export interface FormFieldGuide {
    info: string;
    fieldState: FieldState;
}
  
import { Asset } from "./asset";
import { Permission } from "./permission";

export type Wallet = {
    createdAt: Date,
    publicKey: string,
    tag: string,
    description: string,
    alias: string,
    signer: string,
    userId: string,
    sharedAccessEnabled: boolean, 
    primaryWallet: boolean,
    walletType: number,
    walletThreshold: number,
    numberOfApprovalsNeeded: number, 
    permissions: Permission[],
    sharedAccessCreatedAt: Date,
    sharedAccessUpdatedAt: Date,
    nfts: any[],
    claimedAssets: Asset[],
    unclaimedAssets: Asset[],
}
  
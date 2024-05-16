import { Asset } from "./asset";
import { Permission } from "./permission";

export type SharedWallet = {
    owner: string,
    permission: string,
    walletAlias: string,
    walletDescription: string,
    walletPublicKey: string,
    walletSettings: WalletSettings,    
    assetBalances: {
        claimed: Asset[],
        unclaimed: Asset[],
    }
}

export type WalletSettings = {
    numberOfApprovalsNeeded: number,
    walletThreshold: number,
    walletType: number,
    permissions: Permission[],
}
  
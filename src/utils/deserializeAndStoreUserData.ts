import { DefaultAsset } from "../types/defaultAsset";
import { SharedWallet } from "../types/sharedWallet";
import { User } from "../types/user";
import { Wallet } from "../types/wallet";

export const deserializeUserData = (data: any): User => {
    const userData = data.userData as unknown as User;
    const walletsMap = new Map<string, Wallet>();
    const walletsSharedWithUser: SharedWallet[] = [];
    const defaultAssets: DefaultAsset[] = [];

    for (var wallet of data.userData.userWallets) {
        const w = {
        ...wallet,
        unclaimedAssets: [],
        claimedAssets: [],
        nfts: [],
        };
        walletsMap.set(wallet.publicKey, w);
    }

    for (var assetKey in data.assetBalances) {
        const assetBalance = data.assetBalances[assetKey];
        walletsMap
        .get(assetKey)
        ?.claimedAssets.push(...assetBalance.claimed);
        walletsMap
        .get(assetKey)
        ?.unclaimedAssets.push(...assetBalance.unclaimed);
    }

    for (var assetKey in data.nfts) {
        const assetBalance = data.nfts[assetKey] as Array<any>;
        walletsMap.get(assetKey)?.nfts.push(...assetBalance);
    }

    for (var wallet of data.walletsSharedWithUser) {
        walletsSharedWithUser.push(wallet as SharedWallet);
    }

    for (var asset of data.defaultAssets) {
        defaultAssets.push(asset as DefaultAsset);
    }
    
    return {
        ...userData,
        userWallets: [...walletsMap.values()],
        walletsSharedWithUser,
        defaultAssets,              
    }; 
}
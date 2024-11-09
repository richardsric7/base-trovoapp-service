// import { Asset } from "@stellar/stellar-base";
import { DefaultAsset } from "../types/defaultAsset";
import { Permission } from "../types/permission";
import { SharedWallet } from "../types/sharedWallet";
import { User } from "../types/user";
import { Wallet } from "../types/wallet";

export const deserializeUserData = (data: any): User => {
    const userData = data.userData as unknown as User;
    const walletsMap = new Map<string, Wallet>();
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
        const claimed = assetBalance.claimed;
        const unclaimed = assetBalance.unclaimed;
        if(assetBalance){
            claimed.map((a: any) => {
                walletsMap
                .get(assetKey)
                ?.claimedAssets.push({
                    ...a,
                    amount: Number(a.amount),
                    usdPrice: Number(a.usdPrice),
                    inTrade: {
                        sellingLiabilities: Number(a.inTrade.sellingLiabilities),
                        buyingLiabilities: Number(a.inTrade.buyingLiabilities),
                    },
                    nativePrice: a.nativePrice,
                });
            });

            unclaimed.map((a: any) => {
                walletsMap
                .get(assetKey)
                ?.unclaimedAssets.push({
                    ...a,
                    amount: Number(a.amount),
                    usdPrice: Number(a.usdPrice),
                    inTrade: {
                        sellingLiabilities: Number(a.inTrade.sellingLiabilities),
                        buyingLiabilities: Number(a.inTrade.buyingLiabilities),
                    },
                    nativePrice: a.nativePrice,
                });
            });            
        }
    }

    for (var assetKey in data.nfts) {
        const assetBalance = data.nfts[assetKey] as Array<any>;
        if(assetBalance){
            walletsMap.get(assetKey)?.nfts.push(...assetBalance);
        }
    }

    for (var wallet of data.walletsSharedWithUser) {
        // walletsSharedWithUser.push(wallet as SharedWallet);
        const d = wallet as SharedWallet;
        const w = {
            owner: d.owner,
            permission: d.permission,
            permissions: d.walletSettings ? d.walletSettings?.permissions as Permission[] : [],
            alias: d.walletAlias,
            signer: userData.primarySigner,
            sharedAccessEnabled: true, 
            description: d.walletDescription,
            publicKey: d.walletPublicKey,
            numberOfApprovalsNeeded: d.walletSettings ? d.walletSettings.numberOfApprovalsNeeded : null,
            walletThreshold: d.walletSettings ? d.walletSettings.walletThreshold : null,
            walletType: d.walletSettings ? d.walletSettings.walletType : null,
            claimedAssets: d.assetBalances.claimed.map((a: any) => {return {
                ...a,
                amount: Number(a.amount),
                usdPrice: Number(a.usdPrice),
                inTrade: {
                    sellingLiabilities: Number(a.inTrade.sellingLiabilities),
                    buyingLiabilities: Number(a.inTrade.buyingLiabilities),
                },
                nativePrice: a.nativePrice,
            }}),
            unclaimedAssets: d.assetBalances.unclaimed.map((a: any) => {return {
                ...a,
                amount: Number(a.amount),
                usdPrice: Number(a.usdPrice),
                inTrade: {
                    sellingLiabilities: Number(a.inTrade.sellingLiabilities),
                    buyingLiabilities: Number(a.inTrade.buyingLiabilities),
                },
                nativePrice: a.nativePrice,
            }}),
            nfts: [],
        };
        walletsMap.set(d.walletPublicKey, w as unknown as Wallet);
    }

    for (var asset of data.defaultAssets) {
        defaultAssets.push(asset as DefaultAsset);
    }
    
    return {
        ...userData,
        userWallets: [...walletsMap.values()],
        defaultAssets,              
    }; 
}
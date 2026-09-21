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
            claimedAssets: [],
            nfts: [],
        };
        w.isInitiator = getAccesses(w.permissions, userData.username).includes('INITIATOR');
        w.isApprover = getAccesses(w.permissions, userData.username).includes('APPROVER');
        w.isSharedWallet = w.sharedAccessEnabled;
        w.isSharedWalletAndCanInitiate = w.sharedAccessEnabled && w.walletThreshold == 2 && getAccesses(w.permissions, userData.username).includes('INITIATOR');
        w.canInitiate = !w.sharedAccessEnabled || (w.sharedAccessEnabled && w.walletThreshold == 2 && getAccesses(w.permissions, userData.username).includes('INITIATOR')) && w.walletType == 0;
        w.isPrimaryWallet = w.primaryWallet == 1;

        if(w.isSharedWallet && !w.canInitiate) console.log('========> cannot initiate', w);
        walletsMap.set(wallet.address, w);
    }    

    for (var assetKey in data.assetBalances) {
        const assetBalance = data.assetBalances[assetKey];
        const claimed = assetBalance.claimed;
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
        const permissions = d.walletSettings ? d.walletSettings?.permissions as Permission[] : [];
        const w = {
            owner: d.owner,
            permission: d.permission,
            permissions: permissions,
            alias: d.walletAlias,
            signer: userData.primarySigner,
            sharedAccessEnabled: true, 
            isInitiator: getAccesses(permissions, userData.username).includes('INITIATOR') ,
            isApprover: getAccesses(permissions, userData.username).includes('APPROVER'),
            isSharedWallet: true,
            isSharedWalletAndCanInitiate: d.walletSettings?.walletThreshold == 2 && getAccesses(permissions, userData.username).includes('INITIATOR'),
            canInitiate: (d.walletSettings?.walletThreshold == 2 && getAccesses(permissions, userData.username).includes('INITIATOR')) && d.walletSettings.walletType == 0,
            isPrimaryWallet: false,
            description: d.walletDescription,
            address: d.walletAddress,
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
                tokenizedAsset: a.tokenizedAsset === 1,
            }}),
            nfts: [],
        };
        walletsMap.set(d.walletAddress, w as unknown as Wallet);
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

function getAccesses(permissions: Permission[], username: string): string[] {
        return [
            ...new Set(
                (permissions ?? [])
                    .filter(p => p.targetUsername === username)
                    .map(p => p.permission)
            ),
        ];
    }
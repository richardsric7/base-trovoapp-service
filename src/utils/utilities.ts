import { BANTUBLOCKCHAINEXPLORERBASEURL, BANTUBLOCKCHAINEXPLORERTESTNETBASEURL } from "../store/constants";
import { Wallet } from "../types/wallet";

const calculateFiatValue = (assetBalance: number, fiatRate: number, usdPrice: number) =>
    getFiatRate(usdPrice, fiatRate) * assetBalance;

const getFiatRate = (usdPrice: number, fiatRate: number): number => {
  usdPrice = !usdPrice ? 0 : usdPrice;
  return (fiatRate * usdPrice);  
}

const formatToDecimal = (amount: number) => `
${amount.toLocaleString('en-NG', {
  style: 'decimal',  
  minimumFractionDigits: 7,
  maximumFractionDigits: 7,
})}`;

const shareWallets = (wallets: Wallet[]) => {
    return wallets.filter((w) => w.sharedAccessEnabled);
}

const isSharedWalletAndCanInitiate = (wallet: Wallet) => {
    if(wallet.sharedAccessEnabled && wallet.walletThreshold == 2 && wallet.permission === 'INITIATOR'){
        return true;
    }

    return false;
}

const canInitiate = (wallet: Wallet) => {
    if(!wallet.sharedAccessEnabled || (wallet.sharedAccessEnabled && wallet.walletThreshold == 2 && wallet.permission === 'INITIATOR')){
        return true;
    }

    return false;
}

const getAssetCode = (assetCode: string) => {
    // assign XBN to the asset which has an
    // empty assetCode value.
    // native token of the bantu blockchain
    // has empty values as assetCode and
    // assetIssuer
    return !assetCode ? 'XBN' : assetCode.toString();
  }

const totalAccountBalanceInCurrency = (allWallets: Wallet[], fiatRate: number) => {
    let balance = 0;
    allWallets.map((w) => {
        w.claimedAssets.map((a) => {
            balance += calculateFiatValue(a.amount, fiatRate, a.usdPrice);
        });
    });

    return balance;
}

const totalWalletBalanceInCurrency = (wallet: Wallet, fiatRate: number) => {
    let balance = 0;
    wallet.claimedAssets.map((a) => {
        balance += calculateFiatValue(a.amount, fiatRate, a.usdPrice);
    });

    return balance;
}

const getExplorerBaseUrl = (walletMode: string) => {
    return walletMode == 'Mainnet'
        ? BANTUBLOCKCHAINEXPLORERBASEURL
        : BANTUBLOCKCHAINEXPLORERTESTNETBASEURL;
}

const getBytesLength = (text: string) => new TextEncoder().encode(text).length;

export {
    calculateFiatValue,
    getBytesLength,
    shareWallets,
    getAssetCode,
    totalAccountBalanceInCurrency, 
    formatToDecimal,
    getExplorerBaseUrl,
    totalWalletBalanceInCurrency,
    isSharedWalletAndCanInitiate,
    canInitiate,
}


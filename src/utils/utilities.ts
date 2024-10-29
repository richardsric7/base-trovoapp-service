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

export {
    calculateFiatValue,
    shareWallets,
    totalAccountBalanceInCurrency, 
    formatToDecimal,
    totalWalletBalanceInCurrency,
}


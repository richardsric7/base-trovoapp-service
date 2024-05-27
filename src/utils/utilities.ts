// const calculateFiatValue = (assetBalance: string, usdPrice: string, currency: string) =>
// formatNumber(double.parse(getFiatRate(usdPrice, currency, getUnFormatted: true)) *
//         double.parse(assetBalance))
//     .toString();

// const getFiatRate = (usdPrice: number, currency: string, fiatRate: number, getUnFormatted: boolean = false) => {
//   usdPrice = !usdPrice ? 0 : usdPrice;
//   if (getUnFormatted)
//     return (fiatRate * usdPrice).toString();

//   return formatToCurrency("#,##0.00000", "en_US")
//       .format(appState.fiatRate[currency] * double.parse(usdPrice))
//       .toString();
// }

// const formatToCurrency = (amount: number, currencySymbol: string) => `
// ${currencySymbol}${amount.toLocaleString('en-NG', {
//   // style: 'currency',
//   // currency: 'GBP',
//   minimumFractionDigits: 2,
//   maximumFractionDigits: 2,
// })}`;

// export {calculateFiatValue}


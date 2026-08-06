export const formatNumber = (number: number ): string => {
    var val = Math.round(Number(number) *10000000) / 10000000;
    var parts = val.toString().split(".");
     return parts[0].replace(/\B(?=(\d{3})+(?!\d))/g, ",") + (parts[1] ? "." + parts[1] : "");
  };

  export const getFiatValue = (amount: number): string => {
  if (amount < 99000000000) {
    return formatNumberShort(amount);
  }

  return formatHistoryNumber(amount, 99000000000)
    .toString()
    .replace(/,/g, '');
};

export const formatHistoryNumber = (
  number: number,
  trimNum: number,
  isShort: boolean = false
): string => {
  // if number is greater than or equal to trimNum return compact format (e.g. 1.2M)
  if (number >= trimNum) {
    return new Intl.NumberFormat('en-US', {
      notation: 'compact',
      maximumFractionDigits: 2
    }).format(number);
  }

  return isShort ? formatNumberShort(number) : formatNumber(number);
};

export const formatNumberShort = (number: number): string => {
  const formattedString = new Intl.NumberFormat('en-US', {
    minimumFractionDigits: 0,
    maximumFractionDigits: 2
  }).format(number);

  return formattedString;
};
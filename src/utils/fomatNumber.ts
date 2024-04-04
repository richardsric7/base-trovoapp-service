export const formatNumber = (number: number ): string => {
    var val = Math.round(Number(number) *10000000) / 10000000;
    var parts = val.toString().split(".");
     return parts[0].replace(/\B(?=(\d{3})+(?!\d))/g, ",") + (parts[1] ? "." + parts[1] : "");
  };
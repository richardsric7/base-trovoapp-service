export const truncateValues = (value: number, decimalPlaces:number) => {
    const output = value.toString();
    return output.slice(0, (output.indexOf("."))+decimalPlaces)
  }
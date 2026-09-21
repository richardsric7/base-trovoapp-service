export const truncateValues = (value: number, decimalPlaces:number) => {
    const output = value.toString();
    return output.slice(0, (output.indexOf("."))+decimalPlaces)
  }

  export const truncateAddress = (address: string) => {
    if (address == null) return "Enter address";
    if (address.length <= 7) return address;
    return truncate(address,  7) +
        address.substring(address.length - 7);
  }
  
  export const truncate = (text: string,  length: number = 7, omission: string  = '....'): string => {  
    if (length >= text.length) {
      return text;
    }
  
    return text.slice(0, length) + omission;
  }
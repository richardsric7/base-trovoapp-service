export const truncateValues = (value: number, decimalPlaces:number) => {
    const output = value.toString();
    return output.slice(0, (output.indexOf("."))+decimalPlaces)
  }

  export const truncatePublicKey = (publicKey: string) => {
    if (publicKey == null) return "Enter public key";
    if (publicKey.length <= 7) return publicKey;
    return truncate(publicKey,  7) +
        publicKey.substring(publicKey.length - 7);
  }
  
  export const truncate = (text: string,  length: number = 7, omission: string  = '....'): string => {  
    if (length >= text.length) {
      return text;
    }
  
    return text.slice(0, length) + omission;
  }
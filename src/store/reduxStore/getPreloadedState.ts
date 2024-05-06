export const getPreloadedState = ()  => {
   const defalutValue = {
    auth: {
      user: undefined,
      tempData: {
        usePassphrase: false,
        passphrase: '',
        importExistingWallet: false,
        secretKey: '',
        password: '',
        agreesToTerms: false
      }      
    },
  };
  return defalutValue;
};

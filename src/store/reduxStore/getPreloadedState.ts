export const getPreloadedState = ()  => {
   const defalutValue = {
    auth: {
      user: undefined,
      regFormInfo: {
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

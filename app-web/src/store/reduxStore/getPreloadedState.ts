export const getPreloadedState = ()  => {
   const defalutValue = {
    auth: {
      user: undefined,
      tempData: {
        usePassphrase: false,
        passphrase: '',
        importExistingWallet: false,
        username: '',
        secretKey: '',
        emailOtp: '',
        address: '',
        password: '',
        agreesToTerms: false
      }      
    },
  };
  return defalutValue;
};

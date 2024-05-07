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
        publicKey: '',
        password: '',
        agreesToTerms: false
      }      
    },
  };
  return defalutValue;
};

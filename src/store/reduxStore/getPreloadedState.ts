import { User } from '../../types/user';
import { getStorage } from '../../utils/storage';
import { USER_DETAILS } from '../constants';


export const getPreloadedState = ()  => {
  const userDetails = getStorage(USER_DETAILS) as User ?? null;
  const defalutValue = {
    auth: {
      user: userDetails,
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

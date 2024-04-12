import { createSlice } from '@reduxjs/toolkit';
import { USER_DETAILS } from './constants';
import { User } from '../types/user';

type AuthState = {
  user: User,
  regFormInfo: {
    usePassphrase: boolean,
    importExistingWallet: boolean,
    secretKey: string,
    agreesToTerms: boolean
  }
}

const initialState: AuthState = {
  user: {
    username: '',
    password: '',
    firstName: '',
    lastName: '',
    email: '',
    mobileCountryCode: 'NG',
    mobile: '',
    referrer: '',
    isCorporateUser: false,
    publicKey: '',
  },
  regFormInfo: {
    usePassphrase: false,
    importExistingWallet: false,
    secretKey: '',
    agreesToTerms: false
  }
};
export const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    setUser: (state, action) => {
      const storage = localStorage;
      if (action.payload) {
        state.user = action.payload;
        storage.setItem(USER_DETAILS, JSON.stringify(state.user));
      } else {
        state = initialState;
        storage.clear();
      }
      return state;
    },
    setTempUser: (state, action) => {
      if (action.payload) {
        return {...state, user: action.payload};
      } 
      return state;
    },
    setFormState: (state, action) => {
      if (action.payload) {
        return {...state, regFormInfo: action.payload};
      } 
      return state;
    },
  },
});

export const { setUser, setTempUser, setFormState } = authSlice.actions;

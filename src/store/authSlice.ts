import { createSlice } from '@reduxjs/toolkit';
import { USER_DETAILS } from './constants';
import { User } from '../types/user';

export type AuthState = {
  user?: User,
  regFormInfo: {
    usePassphrase: boolean,
    passphrase: string,
    importExistingWallet: boolean,
    secretKey: string,
    password: string,
    agreesToTerms: boolean
  }
}

const initialState: AuthState = {
  user: undefined,
  regFormInfo: {
    usePassphrase: false,
    importExistingWallet: true,
    secretKey: 'SAZA4CU34762KCGWDAGGARTMWUBLYEXFGYQ34CK2LIXXLA5264AFLNJ5',
    password: '',
    passphrase: '',
    agreesToTerms: true
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

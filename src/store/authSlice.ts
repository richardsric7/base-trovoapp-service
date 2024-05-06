import { createSlice } from '@reduxjs/toolkit';
import { USER_DETAILS } from './constants';
import { User } from '../types/user';
import { setStorage } from '../utils/storage';

export type AuthState = {
  user?: User,
  tempData: {
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
  tempData: {
    usePassphrase: false,
    importExistingWallet: true,
    secretKey: '',
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
        state.user = action.payload.user;
        setStorage(USER_DETAILS, {__slw31H408: action.payload.encryptedUser, __39deR7sx4: action.payload.key, __i34dcY9Mn: state.user!.publicKey});               
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
        return {...state, tempData: action.payload};
      } 
      return state;
    },
  },
});

export const { setUser, setTempUser, setFormState } = authSlice.actions;

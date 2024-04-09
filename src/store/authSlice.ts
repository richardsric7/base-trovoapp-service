import { createSlice } from '@reduxjs/toolkit';
import { USER_DETAILS } from './constants';
import { authApi } from './api/authApi';
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
    username: 'Kent',
    password: 'K@nt2cky',
    firstName: 'Kennis',
    lastName: 'Maduka',
    email: 'madukakennis@gmail.com',
    phoneNumber: '+2347065027384',
    referrer: 'kenmaddy',
    isCorporateUser: false,
  },
  regFormInfo: {
    usePassphrase: false,
    importExistingWallet: false,
    secretKey: '',
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
  extraReducers: (builder) => {
    builder.addMatcher(
      authApi.endpoints.toggleUserAvailability.matchFulfilled,
      (state, { payload }) => {
        localStorage.setItem(USER_DETAILS, JSON.stringify(payload));
        return { ...state, user: payload };
      },
    );    
  },
});

export const { setUser, setTempUser, setFormState } = authSlice.actions;

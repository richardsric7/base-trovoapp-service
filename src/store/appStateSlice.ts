import { createSlice } from '@reduxjs/toolkit';
import { HIDEBALANCES, WALLETMODE } from './constants';
import { setStorage } from '../utils/storage';

export type AppStateSlice = {
  hideBalances: number,
  walletMode: string,
}

const initialState: AppStateSlice = {
  hideBalances: 0,
  walletMode: 'TESTNET',
};
export const appStateSlice = createSlice({
  name: 'appState',
  initialState,
  reducers: {
    setHideBalances: (state, action) => {
      if (action.payload) {
        state.hideBalances = action.payload;
        setStorage(HIDEBALANCES, {hideBalances: state.hideBalances});               
      } 
      return state;
    },
    setWalletMode: (state, action) => {
      if (action.payload) {
        state.walletMode = action.payload;
        setStorage(WALLETMODE, {walletMode: state.walletMode});               
      } 
      return state;
    },
  },
});

export const { setHideBalances } = appStateSlice.actions;

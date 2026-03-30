import { createSlice } from '@reduxjs/toolkit';
import { ACTIVEWALLET, HIDEBALANCES, WALLETMODE } from './constants';
import { setStorage } from '../utils/storage';
import { Wallet } from '../types/wallet';

export type AppStateSlice = {
  hideBalances: number,
  walletMode: string,
  activeWallet?: Wallet,
  activeAsset?: string,
}

const initialState: AppStateSlice = {
  hideBalances: 0,
  walletMode: 'TESTNET',
  activeWallet: undefined,
  activeAsset: undefined
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
    setActiveWallet: (state, action) => {
      if (action.payload) {
        state.activeWallet = action.payload;
        setStorage(ACTIVEWALLET, {activeWallet: state.activeWallet});               
      } 
      return state;
    },
    setActiveAsset: (state, action) => {
      if (action.payload) {
        state.activeAsset = action.payload;        
      } 
      return state;
    },
  },
});

export const { setHideBalances, setActiveWallet, setActiveAsset } = appStateSlice.actions;

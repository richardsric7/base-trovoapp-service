import { createSlice } from '@reduxjs/toolkit';
import { ACTIVEWALLET, HIDEBALANCES, WALLETMODE } from './constants';
import { setStorage } from '../utils/storage';
import { Wallet } from '../types/wallet';
import { TokenizedAsset } from '../types/tokenizedAsset';
import { TokenizationData } from '../types/tokenizationData';

export type AppStateSlice = {
  hideBalances: number,
  walletMode: string,
  activeWallet?: Wallet,
  activeAsset?: string,
  activeTokenizedAsset?: TokenizedAsset,
  tokenizationData?: TokenizationData,
  availableFinancialAssetTypes: string[],
}

export const AVAILABLE_FINANCIAL_ASSET_TYPES = [
  '1114',
  '1115',
  '1121',
  '1122',
  '1123',
  '1124',
  '1125',
  '1174',
  '1180',
  '1182',
];

const initialState: AppStateSlice = {
  hideBalances: 0,
  walletMode: 'TESTNET',
  activeWallet: undefined,
  activeAsset: undefined,
  activeTokenizedAsset: undefined,
  availableFinancialAssetTypes: AVAILABLE_FINANCIAL_ASSET_TYPES,
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
    setActiveTokenizedAsset: (state, action) => {
      if (action.payload) {
        state.activeTokenizedAsset = action.payload;        
      } 
      return state;
    },
    setTokenizationData: (state, action) => {
      if (action.payload) {
        state.tokenizationData = action.payload;        
      } 
      return state;
    },
  },
});

export const { setHideBalances, setActiveWallet, setActiveAsset, setActiveTokenizedAsset, setTokenizationData } = appStateSlice.actions;

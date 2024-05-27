import { createSlice } from '@reduxjs/toolkit';
import { ANNOUNCEMENTS, APP_VERSION, FIAT_RATES } from './constants';
import { setStorage } from '../utils/storage';

export type AppVersion = {
  ID: number,
  androidUrl: string,
  forceUpdate: number,
  iosUrl: string,
  minVersion: string,
  version: string, 
};

export type CacheState = {
  fiatRates: {},
  appVersion: AppVersion | undefined,
  announcements: [],
}

const initialState: CacheState = {
  fiatRates: {},
  appVersion: undefined,
  announcements: [],
};
export const cacheSlice = createSlice({
  name: 'cache',
  initialState,
  reducers: {
    setFiatRates: (state, action) => {
      if (action.payload) {
        state.fiatRates = action.payload;
        setStorage(FIAT_RATES, state.fiatRates);               
      } 
      return state;
    },
    setAppVersion: (state, action) => {
      if (action.payload) {
        state.fiatRates = action.payload;
        setStorage(APP_VERSION, state.fiatRates);               
      }
      return state;
    },
    setAnnouncements: (state, action) => {
      if (action.payload) {
        state.announcements = action.payload;
        setStorage(ANNOUNCEMENTS, state.announcements);               
      }
      return state;
    },
  },
});

export const { setAnnouncements, setAppVersion, setFiatRates } = cacheSlice.actions;

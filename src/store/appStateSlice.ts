import { createSlice } from '@reduxjs/toolkit';
import { HIDEBALANCES } from './constants';
import { setStorage } from '../utils/storage';

export type AppStateSlice = {
  hideBalances: number,
}

const initialState: AppStateSlice = {
  hideBalances: 0,
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
  },
});

export const { setHideBalances } = appStateSlice.actions;

import { createSlice } from '@reduxjs/toolkit';

type ToasterInfo = {
  type: 'success' | 'error' | 'info' | undefined,
  message: string,
  delay?: number,
}

type GlobalUIState = {
  showSidebar: boolean,
  showLoader: boolean,
  showToaster: ToasterInfo,
}

const initialState: GlobalUIState = {  
    showSidebar: false,
    showLoader: false,
    showToaster: {type: undefined, message: '', delay: 2000},  
};

export const sidebarSlice = createSlice({
  name: 'global',
  initialState,
  reducers: {
    toggleSidebar: (state) => {
      // eslint-disable-next-line
      return {...state, showSidebar: !state.showSidebar}
    },
    showHideLoader: (state) => {
      // eslint-disable-next-line
      return {...state, showLoader: !state.showLoader}
    },
    showLoader: (state) => {
      // eslint-disable-next-line
      return {...state, showLoader: true}
    },
    hideLoader: (state) => {
      // eslint-disable-next-line
      return {...state, showLoader: false}
    },
    showToaster: (state, action) => {
      // eslint-disable-next-line
      return {...state, showToaster: {...state.showToaster, type: action.payload.type, message: action.payload.message, delay: action.payload.delay}}
    },
    hideToaster: (state) => {
      // eslint-disable-next-line
      return {...state, showToaster: {...state.showToaster, type: undefined}}
    },
  },
});

// Action creators are generated for each case reducer function
export const { toggleSidebar, showHideLoader, showLoader, hideLoader, showToaster, hideToaster} = sidebarSlice.actions;

export default sidebarSlice.reducer;

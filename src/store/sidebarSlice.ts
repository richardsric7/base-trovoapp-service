import { createSlice } from '@reduxjs/toolkit';

export const sidebarSlice = createSlice({
  name: 'sidebar',
  initialState: {
    value: false,
  },
  reducers: {
    toggle: (state) => {
      // eslint-disable-next-line
      state.value = !state.value;
    },
  },
});

// Action creators are generated for each case reducer function
export const { toggle } = sidebarSlice.actions;

export default sidebarSlice.reducer;

import { createSlice } from '@reduxjs/toolkit';
import { TOKEN, USER_DETAILS } from './constants';
import { authApi } from './api/authApi';

const initialState = {
  user: null,
  token: { accessToken: null, refreshToken: null },
};
export const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    setToken: (state, action) => {
      state.token.accessToken = action.payload;
      const storage = localStorage;
      storage.setItem(TOKEN, JSON.stringify(state.token.accessToken));
      return state;
    },
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

export const { setUser, setToken } = authSlice.actions;

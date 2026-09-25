import { createSlice } from "@reduxjs/toolkit";
import { TOKEN, USER_DETAILS } from "../constants";
import { IWalletConnectResponse } from "../api";

interface AuthState {
  user: IWalletConnectResponse | null;
  token: {
    accessToken: string | null;
    refreshToken: string | null;
    expiresAt?: number;
  };
}

const initialState: AuthState = {
  user: null,
  token: { accessToken: null, refreshToken: null },
};
export const authSlice = createSlice({
  name: "auth",
  initialState,
  reducers: {
    setToken: (state, action) => {
      state.token = action.payload;
      const storage = localStorage;
      storage.setItem(TOKEN, JSON.stringify(state.token));

      return state;
    },

    setUser: (state, action) => {
      const storage = localStorage;
      if (action.payload) {
        state.user = action.payload;
        storage.setItem(USER_DETAILS, JSON.stringify(state.user));
      } else {
        state = initialState;
        // Only this flow's own keys — the organisation portal keeps its own separate session
        // in the same localStorage, and must not be logged out as a side effect of this.
        storage.removeItem(TOKEN);
        storage.removeItem(USER_DETAILS);
      }
      return state;
    },
  },
  // extraReducers: (builder) => {
  //   builder.addMatcher(
  //     userApi.endpoints.updateContact.matchFulfilled,
  //     (state, { payload }) => {
  //       localStorage.setItem(USER_DETAILS, JSON.stringify(payload));
  //       return { ...state, user: payload };
  //     }
  //   );
  // }
});

export const { setUser, setToken } = authSlice.actions;

import { createSlice } from "@reduxjs/toolkit";
import { TOKEN_ORG, USER_DETAILS_ORG } from "../constants";

export interface OrganizationAuthState {
  //   user: YourOrgUserType | null;
  token: { accessToken: string | null; refreshToken: string | null };
}

const initialState: OrganizationAuthState = {
  //   user: null,
  token: { accessToken: null, refreshToken: null },
};

export const organizationAuthSlice = createSlice({
  name: "organizationAuth",
  initialState,
  reducers: {
    setOrgToken: (state, action) => {
      state.token = action.payload;
      localStorage.setItem(TOKEN_ORG, JSON.stringify(state.token));
    },
    clearOrgAuth: (state) => {
      state.token = { accessToken: null, refreshToken: null };
      localStorage.removeItem(TOKEN_ORG);
      localStorage.removeItem(USER_DETAILS_ORG);
    },
    // setOrgUser: (state, action) => {
    //   if (action.payload) {
    //     state.user = action.payload;
    //     localStorage.setItem(USER_DETAILS_ORG, JSON.stringify(state.user));
    //   } else {
    //     state.user = null;
    //     state.token = { accessToken: null, refreshToken: null };
    //     localStorage.removeItem(TOKEN_ORG);
    //     localStorage.removeItem(USER_DETAILS_ORG);
    //   }
    // },
  },
});

export const { setOrgToken, clearOrgAuth } = organizationAuthSlice.actions;
export default organizationAuthSlice.actions;

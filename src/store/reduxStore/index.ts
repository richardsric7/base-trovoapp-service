import { configureStore } from "@reduxjs/toolkit";
import { setupListeners } from "@reduxjs/toolkit/query";
import axios from "axios";
import { baseApi } from "../api/baseapi";
import { authSlice } from "../authSlice";
import { refreshAccessToken } from "../refreshAccessToken";
import { sidebarSlice } from '../sidebarSlice';
import { getPreloadedState } from "./getPreloadedState";

export const store = configureStore({
    reducer:{
        api: baseApi.reducer,
        auth: authSlice.reducer,
        sidebarSlice: sidebarSlice.reducer,
    },
    middleware:(getDefaultMiddleware) => 
    getDefaultMiddleware().concat(baseApi.middleware),
    preloadedState: getPreloadedState(),
});
setupListeners(store.dispatch);

  axios.interceptors.response.use(
    (response) => {
      return response;
    },
    async function (error) {
      const originalRequest = error.config;
      if (
          error?.response?.status === 401 &&
          !originalRequest._retry
        ) {
          originalRequest._retry = true;
          const access_token = await refreshAccessToken();
          axios.defaults.headers.common['Authorization'] = `${access_token}`;
          return await axios(originalRequest);
        }
        // if (error.response.status === 422) {
        //     store.dispatch(setToken({accessToken: null, refreshToken: null}))
        //     store.dispatch(setUser(null));
        //     clearStorage()
        // }
        return Promise.reject(error);
    },
  );

  export type RootState = ReturnType<typeof store.getState>
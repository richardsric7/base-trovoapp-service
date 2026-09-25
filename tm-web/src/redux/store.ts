import { configureStore } from "@reduxjs/toolkit";
// import { baseApi } from "./baseApi";
import { setupListeners } from "@reduxjs/toolkit/query";
import { authSlice } from "./slices/authSlice";
import { organizationAuthSlice } from "./slices/orgSlice";

import { baseApi } from "./baseApi";
import { orgApi } from "./baseApi/orgApi";

export const makeStore = () => {
  const store = configureStore({
    reducer: {
      api: baseApi.reducer,
      auth: authSlice.reducer,
      organizationAuth: organizationAuthSlice.reducer,

      orgApi: orgApi.reducer,
    },
    middleware: (getDefaultMiddleware) =>
      getDefaultMiddleware()
        .concat(baseApi.middleware)
        .concat(orgApi.middleware),
    // devTools: process.env.NODE_ENV !== 'production',
  });

  // Wires up the browser visibilitychange/focus/online listeners RTK Query needs for
  // refetchOnFocus/refetchOnReconnect to do anything at all. Opt-in per query (see
  // useVerifyLoginQuery in qr-scan/page.tsx) rather than a global default, so this doesn't
  // change refetch behavior for every other table/list in the app.
  setupListeners(store.dispatch);

  return store;
};

export type AppStore = ReturnType<typeof makeStore>;
export type AppDispatch = AppStore["dispatch"];
export type RootState = ReturnType<AppStore["getState"]>;

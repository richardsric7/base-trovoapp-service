// import axios from "axios";
// import { setStorage } from "../utils/storage";
// import { TOKEN } from "./constants";
// import { BASE_URL } from "@/config";

// export const refreshAccessToken = async () => {
//   // const newAxios = axios.create();
//   const response = await axios.post(`${BASE_URL}/v1/login/token/refresh`);
//   setStorage(TOKEN, {
//     accessToken: response.data.accessToken,
//     refeshToken: response.data.refreshToken,
//   });
//   return response.data.accessToken;
// };

import axios from "axios";
import { setStorage, getStorage } from "../utils/storage";
import { TOKEN } from "./constants";
import { BASE_URL } from "@/config";

// Error types
export enum AuthErrorType {
  REFRESH_FAILED = "REFRESH_FAILED",
}

interface AuthError extends Error {
  type: AuthErrorType;
}

interface TokenData {
  accessToken: string;
  refreshToken: string;
}

// A bare axios instance: the refresh call must NOT go through the global
// interceptors. The request interceptor used to call getPreloadedState() (which
// itself triggered a refresh), and the response interceptor retries 401s by
// refreshing - either way a refresh could spawn more refreshes without bound.
const refreshClient = axios.create();

// Single-flight: every caller during one expiry window shares the same request.
let inFlight: Promise<string> | null = null;

export const refreshAccessToken = (): Promise<string> => {
  if (inFlight) return inFlight;
  inFlight = (async () => {
    try {
      const currentToken = getStorage<TokenData>(TOKEN);
      if (!currentToken?.refreshToken) {
        throw new Error("No refresh token available");
      }

      // BASE_URL already ends in /api/v1 - the old `${BASE_URL}/v1/...` built
      // /api/v1/v1/login/token/refresh, which does not exist (404), so every
      // refresh failed and the expiry window never advanced.
      const response = await refreshClient.post<TokenData>(
        `${BASE_URL}/login/token/refresh`,
        { refreshToken: currentToken.refreshToken }
      );

      setStorage(TOKEN, {
        accessToken: response.data.accessToken,
        refreshToken: response.data.refreshToken,
        // No JWT decoding client-side; assumes the same 15-minute validity
        // scheduleTokenRefresh already assumes.
        expiresAt: Date.now() + 15 * 60 * 1000,
      });

      return response.data.accessToken;
    } catch (error) {
      const authError = new Error("Failed to refresh token") as AuthError;
      authError.type = AuthErrorType.REFRESH_FAILED;
      throw authError;
    } finally {
      inFlight = null;
    }
  })();
  return inFlight;
};

import { BaseQueryFn } from "@reduxjs/toolkit/query";
import axios, {
  AxiosError,
  AxiosRequestConfig,
  AxiosRequestHeaders,
  Method,
} from "axios";
import { getPreloadedState } from "../getPreloadedState";
import { setStorage } from "../../utils/storage";
import { TOKEN } from "../constants";
import { BASE_URL } from "@/config";
import { refreshAccessToken } from "../refreshAccessToken";
import {
  REQUEST_ID_HEADER,
  setLastRequestId,
  reportError,
} from "@/observability";

// Error types
export enum AuthErrorType {
  TOKEN_EXPIRED = "TOKEN_EXPIRED",
  REFRESH_FAILED = "REFRESH_FAILED",
  NETWORK_ERROR = "NETWORK_ERROR",
  UNAUTHORIZED = "UNAUTHORIZED",
}

interface AuthError extends Error {
  type: AuthErrorType;
  status?: number;
}

// Navigation handler - Redirects to sign-in.
// Guard: the admin API client's 401/refresh-failure handlers call this to bounce
// the user to the admin sign-in. The organisation (stakeholder) portal lives under
// /organisation and authenticates with its own token via the org API client; a 401
// there (e.g. an org screen that still calls an admin-only endpoint, or a bad org
// login) must NOT eject the org user to the admin sign-in. So suppress the admin
// redirect while on an /organisation route.
export const handleNavigation = (path: string) => {
  if (typeof window === "undefined") return;
  if (window.location.pathname.startsWith("/organisation")) return;
  // Assigning location.href to the page we are already on triggers a full reload,
  // and if the trigger (e.g. an expired token in storage) is still there on the
  // next boot, the page reloads forever and looks frozen. Never self-navigate.
  const target = new URL(path, window.location.origin);
  if (target.pathname === window.location.pathname) return;
  window.location.href = path;
};

let isRefreshing = false;
let failedQueue: Array<{
  resolve: (token: string) => void;
  reject: (error: any) => void;
}> = [];

const processQueue = (error: any, token: string | null = null) => {
  failedQueue.forEach((prom) => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve(token!);
    }
  });
  failedQueue = [];
};

// Axios request interceptor to attach token
// axios.interceptors.request.use(
//   async (config) => {
//     const token = getPreloadedState().auth.token;

//     if (token?.accessToken) {
//       config.headers = {
//         ...config.headers,
//         Authorization: `Bearer ${token.accessToken}`,
//       } as AxiosRequestHeaders;
//     }

//     return config;
//   },
//   (error) => Promise.reject(error)
// );

axios.interceptors.request.use(
  async (config) => {
    const token = getPreloadedState().auth.token;

    if (token?.accessToken) {
      // Check if the request is for special endpoints that need different auth format
      const isTokenizationEndpoint =
        config.url && config.url.includes("/tokenization");
      const isWalletBalanceEndpoint =
        config.url && config.url.includes("/wallet-balances/");
      const isOrganizations =
        config.url && config.url.includes("/organizations/");
      // Apply the appropriate authorization header format
      config.headers = {
        ...config.headers,
        Authorization:
          isTokenizationEndpoint || isOrganizations || isWalletBalanceEndpoint
            ? token.accessToken
            : `Bearer ${token.accessToken}`,
      } as AxiosRequestHeaders;
    }

    return config;
  },
  (error) => Promise.reject(error),
);

// Axios response interceptor to handle 401 errors
axios.interceptors.response.use(
  (response) => {
    // The backend returns a correlation id on every response (see the observe
    // package). Recording it means a later browser error can be traced to the
    // exact backend request behind it.
    setLastRequestId(response.headers?.[REQUEST_ID_HEADER]);
    return response;
  },
  async (error) => {
    const originalRequest = error.config;
    setLastRequestId(error.response?.headers?.[REQUEST_ID_HEADER]);

    // A 5xx means the backend broke, which is worth investigating even though
    // the UI handles it gracefully. 4xx responses are deliberately not
    // reported: they are expected outcomes (validation, permissions, an
    // expired session) and would drown the real failures.
    const status = error.response?.status;
    if (typeof status === "number" && status >= 500) {
      reportError(error, {
        url: originalRequest?.url,
        method: originalRequest?.method,
        status,
      });
    }

    const isRefreshCall = String(originalRequest?.url || "").includes("/login/token/refresh");
    if (error.response?.status === 401 && !originalRequest._retry && !isRefreshCall) {
      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        })
          .then((token) => {
            originalRequest.headers.Authorization = `Bearer ${token}`;
            return axios(originalRequest);
          })
          .catch((err) => Promise.reject(err));
      }

      originalRequest._retry = true;
      isRefreshing = true;

      try {
        const newAccessToken = await refreshAccessToken();
        isRefreshing = false;
        processQueue(null, newAccessToken);

        originalRequest.headers.Authorization = `Bearer ${newAccessToken}`;
        return axios(originalRequest);
      } catch (refreshError) {
        isRefreshing = false;
        processQueue(refreshError, null);
        await setStorage(TOKEN, null);
        handleNavigation("/sign-in?reason=unauthorized");
        return Promise.reject(refreshError);
      }
    }

    return Promise.reject(error);
  },
);

// Axios base query function for Redux Toolkit

export const axiosBaseQuery =
  ({
    baseUrl = "",
    baseHeaders = {},
  }: {
    baseUrl: string;
    baseHeaders?: AxiosRequestConfig["headers"];
  }): BaseQueryFn<
    {
      url: string;
      method: Method;
      data?: AxiosRequestConfig["data"];
      params?: AxiosRequestConfig["params"];
      headers?: AxiosRequestConfig["headers"];
      transformRequest?: AxiosRequestConfig["transformRequest"]; // ✅ ADD THIS
      responseType?: AxiosRequestConfig["responseType"];
    },
    any,
    unknown
  > =>
  async ({
    url,
    method,
    data,
    params,
    headers = {},
    transformRequest,
    responseType,
  }) => {
    try {
      // If data is FormData, avoid including JSON Content-Type from baseHeaders
      // We use a robust check for FormData (checking if .append exists)
      const isFormData =
        data instanceof FormData ||
        (data && typeof data === "object" && typeof data.append === "function");

      const finalHeaders = {
        ...(isFormData ? {} : baseHeaders),
        ...headers,
      };

      // If data is FormData, we MUST NOT set Content-Type: application/json
      if (isFormData) {
        console.log("📤 Sending FormData - Headers:", finalHeaders);
        delete (finalHeaders as any)["Content-Type"];
        delete (finalHeaders as any)["content-type"];
        delete (finalHeaders as any)["Content-type"];
      } else if (
        !finalHeaders["Content-Type"] &&
        !finalHeaders["content-type"]
      ) {
        (finalHeaders as any)["Content-Type"] = "application/json";
      }

      const result = await axios({
        url: baseUrl + url,
        method,
        params,
        data,
        transformRequest,
        responseType,
        headers: finalHeaders,
      });
      return { data: result.data };
    } catch (axiosError) {
      const err = axiosError as AxiosError;

      if (err.response?.status === 401) {
        const authError = new Error("Unauthorized") as AuthError;
        authError.type = AuthErrorType.UNAUTHORIZED;
        await handleAuthError(authError);
      }

      return {
        error: {
          status: err.response?.status,
          error: err.message,
          ...(typeof err.response?.data === "object" ? err.response.data : {}),
        },
      };
    }
  };

const handleAuthError = async (error: AuthError) => {
  switch (error.type) {
    case AuthErrorType.TOKEN_EXPIRED:
      return;

    case AuthErrorType.REFRESH_FAILED:
    case AuthErrorType.UNAUTHORIZED:
      await setStorage(TOKEN, null);
      handleNavigation("/sign-in?reason=unauthorized");
      break;

    case AuthErrorType.NETWORK_ERROR:
      console.error("Network error occurred:", error);
      break;

    default:
      console.error("Unhandled auth error:", error);
  }
};

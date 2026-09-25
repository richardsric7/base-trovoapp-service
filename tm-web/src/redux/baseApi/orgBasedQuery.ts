// import axios, {
//   AxiosError,
//   AxiosRequestConfig,
//   AxiosRequestHeaders,
//   Method,
// } from "axios";
// import { BaseQueryFn } from "@reduxjs/toolkit/query";
// import { BASE_URL } from "@/config";
// import { getPreloadedOrgState } from "../getPreLoadedStateOrg";

// // orgAxios.ts
// export const orgAxios = axios.create({
//   baseURL: BASE_URL, // just base API URL, no /organization suffix
// });

// // Axios request interceptor to add Authorization header with org token
// orgAxios.interceptors.request.use(
//   (config: any) => {
//     const orgToken = getPreloadedOrgState().organizationAuth.token.accessToken;

//     if (orgToken) {
//       config.headers = config.headers || {};
//       config.headers["Authorization"] = orgToken;
//     }
//     return config;
//   },
//   (error) => Promise.reject(error)
// );

// // orgAxios.interceptors.request.use((config) => {
// //   const raw = localStorage.getItem("org_auth_token");
// //   if (raw) {
// //     const token = JSON.parse(raw).accessToken;
// //     if (token) {
// //       config.headers = config.headers || {};
// //       config.headers["Authorization"] = token;
// //     }
// //   }
// //   return config;
// // });

// const SIGNUP_URL = "/organizations/login"; // Change this to your actual signup route

// export const orgAxiosBaseQuery =
//   ({
//     baseUrl = "",
//     baseHeaders = {},
//   }: {
//     baseUrl?: string;
//     baseHeaders?: AxiosRequestConfig["headers"];
//   }): BaseQueryFn<
//     {
//       url: string;
//       method: Method;
//       data?: AxiosRequestConfig["data"];
//       params?: AxiosRequestConfig["params"];
//       headers?: AxiosRequestConfig["headers"];
//       transformRequest?: AxiosRequestConfig["transformRequest"];
//     },
//     any,
//     unknown
//   > =>
//   async ({ url, method, data, params, headers = {}, transformRequest }) => {
//     // 1. Check for org token before making the request
//     let authHeaders = {};
//     let token: string | null = null;
//     try {
//       const raw = localStorage.getItem("org_auth_token");
//       if (raw) {
//         token = JSON.parse(raw).accessToken;
//         if (token) {
//           authHeaders = { Authorization: `Bearer ${token}` };
//         }
//       }
//     } catch {}

//     if (!token) {
//       // No token: redirect and return error
//       window.location.href = SIGNUP_URL;
//       return {
//         error: {
//           status: 401,
//           error: "Unauthorized: No org token found.",
//         },
//       };
//     }

//     try {
//       const result = await orgAxios({
//         url: baseUrl + url,
//         method,
//         params,
//         data,
//         transformRequest,
//         headers: {
//           ...baseHeaders,
//           ...headers,
//           ...authHeaders,
//           ...(data instanceof FormData
//             ? {}
//             : { "Content-Type": "application/json" }),
//         },
//       });
//       return { data: result.data };
//     } catch (axiosError) {
//       const err = axiosError as AxiosError;

//       // 2. If server responds with 401, redirect to sign up
//       if (err.response?.status === 401) {
//         window.location.href = SIGNUP_URL;
//       }

//       return {
//         error: {
//           status: err.response?.status,
//           error: err.message,
//           ...(typeof err.response?.data === "object" ? err.response.data : {}),
//         },
//       };
//     }
//   };

// RTK Query baseQuery function using orgAxios instance
// export const orgAxiosBaseQuery: BaseQueryFn<any, unknown, unknown> = async ({
//   url,
//   method,
//   data,
//   params,
//   headers = {},
// }) => {
//   let authHeaders = {};
//   try {
//     const raw = localStorage.getItem("org_auth_token");
//     if (raw) {
//       const token = JSON.parse(raw).accessToken;
//       if (token) {
//         authHeaders = {
//           Authorization: `Bearer ${token}`,
//         };
//       }
//     }
//   } catch {}

//   try {
//     const result = await orgAxios({
//       url,
//       method,
//       data,
//       params,
//       headers: {
//         ...headers,
//         ...authHeaders, // merge in fresh org token headers
//       },
//     });
//     return { data: result.data };
//   } catch (axiosError: any) {
//     return {
//       error: {
//         status: axiosError.response?.status,
//         data: axiosError.response?.data || axiosError.message,
//       },
//     };
//   }
// };

import axios, { AxiosError, AxiosRequestConfig, Method } from "axios";
import { BaseQueryFn } from "@reduxjs/toolkit/query";
import { BASE_URL } from "@/config";
import {
  REQUEST_ID_HEADER,
  setLastRequestId,
  reportError,
} from "@/observability";

const SIGNUP_URL = "/organizations/login";

export const orgAxios = axios.create({
  baseURL: BASE_URL,
});

// Same correlation-id capture as the admin client: record the backend request
// id from every response so a browser error can be traced back to it, and
// report 5xx responses (a broken backend) while leaving 4xx alone - those are
// expected outcomes and would drown the real failures.
orgAxios.interceptors.response.use(
  (response) => {
    setLastRequestId(response.headers?.[REQUEST_ID_HEADER]);
    return response;
  },
  (error) => {
    setLastRequestId(error.response?.headers?.[REQUEST_ID_HEADER]);
    const status = error.response?.status;
    if (typeof status === "number" && status >= 500) {
      reportError(error, {
        url: error.config?.url,
        method: error.config?.method,
        status,
        portal: "organisation",
      });
    }
    return Promise.reject(error);
  }
);

type OrgAxiosBaseQueryArgs = {
  url: string;
  method: Method;
  data?: AxiosRequestConfig["data"];
  body?: AxiosRequestConfig["data"];
  params?: AxiosRequestConfig["params"];
  headers?: AxiosRequestConfig["headers"];
  transformRequest?: AxiosRequestConfig["transformRequest"];
  responseType?: AxiosRequestConfig["responseType"];
};

const getOrgToken = (): string | null => {
  try {
    const raw = localStorage.getItem("org_auth_token");

    if (!raw) return null;

    const parsed = JSON.parse(raw);

    return parsed?.accessToken ?? null;
  } catch {
    return null;
  }
};

export const orgAxiosBaseQuery =
  ({
    baseUrl = "",
    baseHeaders = {},
  }: {
    baseUrl?: string;
    baseHeaders?: AxiosRequestConfig["headers"];
  }): BaseQueryFn<OrgAxiosBaseQueryArgs, unknown, unknown> =>
  async ({
    url,
    method,
    data,
    body,
    params,
    headers = {},
    transformRequest,
    responseType,
  }) => {
    const token = getOrgToken();

    if (!token) {
      window.location.href = SIGNUP_URL;
      return {
        error: {
          status: 401,
          error: "Unauthorized: No organization token found.",
        },
      };
    }

    try {
      const requestData = data !== undefined ? data : body;

      const isFormData =
        requestData instanceof FormData ||
        (requestData &&
          typeof requestData === "object" &&
          typeof (requestData as FormData).append === "function");

      const finalHeaders = {
        ...(isFormData ? {} : baseHeaders),
        ...headers,

        // IMPORTANT:
        // Organization endpoints use raw token based on your Swagger test.
        Authorization: token,
      };

      if (isFormData) {
        delete (finalHeaders as Record<string, unknown>)["Content-Type"];
        delete (finalHeaders as Record<string, unknown>)["content-type"];
        delete (finalHeaders as Record<string, unknown>)["Content-type"];
      } else {
        (finalHeaders as Record<string, string>)["Content-Type"] =
          "application/json";
      }

      const result = await orgAxios({
        url: baseUrl + url,
        method,
        params,
        data: requestData,
        transformRequest,
        responseType,
        headers: finalHeaders,
      });

      return {
        data: result.data,
      };

    } catch (axiosError) {
      const err = axiosError as AxiosError;

      // Only a 401 means the org session itself is invalid/expired — that warrants
      // bouncing to the org login. A 403 means the member is authenticated but not
      // permitted for THIS specific resource (e.g. not assigned to this asset); that
      // must not eject them from the portal — surface it as an error and stay put.
      if (err.response?.status === 401) {
        window.location.href = SIGNUP_URL;
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

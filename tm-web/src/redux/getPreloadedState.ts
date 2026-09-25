import { getStorage } from "@/utils";
import { TOKEN, USER_DETAILS } from "@/redux/constants";
interface IToken {
  accessToken: string;
  refreshToken: string;
  expiresAt: number;
}
export const getPreloadedState = () => {
  const userDetails = getStorage(USER_DETAILS);
  // const token = getStorage(TOKEN) ?? {};
  const token = getStorage<IToken>(TOKEN) || {
    accessToken: "",
    refreshToken: "",
    expiresAt: 0,
  };

  // Deliberately NO refresh side effect here. This function is called from the
  // axios request interceptor on every request; a refresh triggered from here
  // spawned further refreshes (the refresh call is itself a request) without
  // bound. Proactive refresh lives only in scheduleTokenRefresh (timer-driven).

  const defaultValue = {
    auth: {
      user: userDetails,
      token: token as IToken,
    },
  };
  return defaultValue;
};

// import { getStorage } from "@/utils";
// import { TOKEN, USER_DETAILS } from "@/redux/constants";

// export const getPreloadedState = () => {
//   // Get token from cookies
//   const getCookie = (name: string) => {
//     const value = `; ${document.cookie}`;
//     const parts = value.split(`; ${name}=`);
//     if (parts.length === 2) return parts.pop()?.split(";").shift();
//     return "";
//   };

//   const token = getCookie(TOKEN) ?? "";
//   const userDetails = getStorage(USER_DETAILS);

//   return {
//     auth: {
//       user: userDetails,
//       token: token,
//     },
//   };
// };

import { refreshAccessToken } from "./refreshAccessToken";
import { getPreloadedState } from "./getPreloadedState";
import { setStorage, clearAdminSession } from "../utils/storage";
import { TOKEN } from "./constants";
import { handleNavigation } from "./baseApi/axiosBasedQuery";

let refreshTimeout: NodeJS.Timeout | null = null;

export const scheduleTokenRefresh = () => {
  const token = getPreloadedState().auth.token;
  if (!token || !token.expiresAt) return;

  const expiresIn = token.expiresAt - Date.now();
  if (expiresIn <= 0) {
    // Drop the expired session first, otherwise the sign-in page boots with the
    // same dead token and redirects again — a loop the user cannot escape.
    clearAdminSession();
    handleNavigation("/sign-in?reason=session_expired");
    return;
  }

  // Schedule token refresh 2 minutes before expiry
  const refreshTime = expiresIn - 2 * 60 * 1000;

  if (refreshTimeout) clearTimeout(refreshTimeout);

  refreshTimeout = setTimeout(async () => {
    try {
      const newAccessToken = await refreshAccessToken();
      if (newAccessToken) {
        setStorage(TOKEN, {
          accessToken: newAccessToken,
          expiresAt: Date.now() + 15 * 60 * 1000, // Assuming 15 min validity
        });
        scheduleTokenRefresh(); // Re-schedule for the new expiry
      } else {
        clearAdminSession();
        handleNavigation("/sign-in?reason=session_expired");
      }
    } catch (error) {
      clearAdminSession();
      handleNavigation("/sign-in?reason=session_expired");
    }
  }, refreshTime);
};

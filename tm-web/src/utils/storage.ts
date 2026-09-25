import { TOKEN, USER_DETAILS } from "@/redux/constants";
import { deleteCookie } from "cookies-next";

export const getStorage = <T>(itemKey: string) => {
  try {
    let value = localStorage.getItem(itemKey);
    if (value) return JSON.parse(value) as T;
    return null;
  } catch (error) {
    return null;
  }
};
export const setStorage = <T>(itemKey: string, data: T) => {
  localStorage.setItem(itemKey, JSON.stringify(data));
};
// Clears only this dashboard flow's own keys, never the whole of localStorage — the
// organisation portal keeps a separate session (org_auth_token/org_user_details) in the same
// localStorage, and logging out of one flow must not silently log the other out too.
// Wipes the admin session from BOTH localStorage and the auth cookie. Used when a
// stored token turns out to be expired/unrefreshable, so the next page load does
// not trip over the same dead token again (that was a boot-time redirect loop that
// only a manual cache wipe could break).
export const clearAdminSession = () => {
  clearStorage();
  deleteCookie(TOKEN);
};

export const clearStorage = () => {
  localStorage.removeItem(TOKEN);
  localStorage.removeItem(USER_DETAILS);
};

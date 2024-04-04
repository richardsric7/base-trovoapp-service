import { getStorage } from "../../utils/storage";
import { TOKEN, USER_DETAILS } from "../constants";

export const getPreloadedState = () => {
  const userDetails = getStorage(USER_DETAILS);
  const token = getStorage(TOKEN)?.accessToken ?? "";
  const defalutValue = {
    auth: {
      user: userDetails,
      token: token,
    },
  };
  return defalutValue;
};

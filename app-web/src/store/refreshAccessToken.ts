import axios from "axios";
import { setStorage } from "../utils/storage";
import { BASE_URL } from "./config";
import { TOKEN } from "./constants";

export const refreshAccessToken = async () => {  
    // const newAxios = axios.create();
    const response = await axios.post(
      `${BASE_URL}/v1/login/token/refresh`,
    );
    setStorage(TOKEN, { accessToken: response.data.accessToken, refeshToken: response.data.refreshToken });
    return response.data.accessToken;
  };
  
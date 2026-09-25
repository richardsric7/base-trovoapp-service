// src/getPreloadedOrgState.ts
import { OrganizationAuthState } from "./slices/orgSlice";
import { TOKEN_ORG } from "./constants";

export function getPreloadedOrgState(): {
  organizationAuth: OrganizationAuthState;
} {
  let token = null;

  try {
    const storedToken = localStorage.getItem(TOKEN_ORG);
    if (storedToken) {
      token = JSON.parse(storedToken);
    }
  } catch (error) {
    console.error("Failed to parse org token from storage", error);
  }

  return {
    organizationAuth: {
      token: {
        accessToken: token?.accessToken || null,
        refreshToken: token?.refreshToken || null, // including for structure, even if unused now
      },
    },
  };
}

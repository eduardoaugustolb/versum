import { cache } from "react";
import { getApiV1UsersMe } from "@/dal/orval/fetch/users/users";
import { getSessionCookieName } from "@/lib/auth";

export const getCachedUserAuth = cache(
  async (sessionToken?: string, baseUrl?: string) => {
    try {
      return await getApiV1UsersMe({
        baseURL: baseUrl,
        headers: sessionToken
          ? { Cookie: `${getSessionCookieName()}=${encodeURIComponent(sessionToken)}` }
          : undefined,
      } as RequestInit & { baseURL?: string });
    } catch (err) {
      console.error("Failed to fetch user session", err);
      return null;
    }
  },
);

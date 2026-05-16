import { cache } from "react";
import { getApiV1UsersMe } from "@/dal/orval/fetch/users/users";
import { getSessionCookieName } from "@/lib/auth";

export const getCachedUserAuth = cache(async (sessionToken?: string) => {
  try {
    return await getApiV1UsersMe({
      headers: sessionToken
        ? {
            Cookie: `${getSessionCookieName()}=${encodeURIComponent(sessionToken)}`,
          }
        : undefined,
    });
  } catch (err) {
    console.error("Failed to fetch user session", err);
    return null;
  }
});

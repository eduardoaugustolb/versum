import { getApiV1UsersMe } from "@/lib/kubb/gen";
import { cache } from "react";

export const getCachedUserAuth = cache(getApiV1UsersMe.bind(getApiV1UsersMe, {
  baseURL: process.env.NEXT_PUBLIC_API_URL,
}));

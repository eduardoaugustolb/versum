import { cache } from "react";
import { getApiV1UsersMe } from "@/lib/kubb/gen";

export const getCachedUserAuth = cache(
  getApiV1UsersMe.bind(getApiV1UsersMe),
);

import type { ApiVersion } from "../../api-version.ts";
import { AuthControllerV1 } from "../controllers/auth.v1.controller.ts";
import { createAuthRoutesV1 } from "./auth.v1.route.ts";

export const createAuthRoutes = (version: ApiVersion) => {
  switch (version) {
    case "v1":
      return createAuthRoutesV1(new AuthControllerV1());
    default: {
      const exhaustiveCheck: never = version;
      throw new Error(`Unsupported auth API version: ${exhaustiveCheck}`);
    }
  }
};

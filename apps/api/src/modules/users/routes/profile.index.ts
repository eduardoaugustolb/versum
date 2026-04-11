import type { ApiVersion } from "../../api-version.ts";
import { ProfileControllerV1 } from "../controllers/profile.v1.controller.ts";
import { createProfileRoutesV1 } from "./profile.v1.routes.ts";

export const createProfileRoutes = (version: ApiVersion) => {
  switch (version) {
    case "v1":
      return createProfileRoutesV1(new ProfileControllerV1());
    default: {
      const exhaustiveCheck: never = version;
      throw new Error(`Unsupported profiles API version: ${exhaustiveCheck}`);
    }
  }
};

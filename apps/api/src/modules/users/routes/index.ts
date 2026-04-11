import type { ApiVersion } from "../../api-version.ts";
import { UsersControllerV1 } from "../controllers/users.v1.controller.ts";
import { createUsersRoutesV1 } from "./users.v1.route.ts";

export const createUsersRoutes = (version: ApiVersion) => {
  switch (version) {
    case "v1":
      return createUsersRoutesV1(new UsersControllerV1());
    default: {
      const exhaustiveCheck: never = version;
      throw new Error(`Unsupported users API version: ${exhaustiveCheck}`);
    }
  }
};

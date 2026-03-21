import type { ApiVersion } from "../../api-version.ts";
import { BibleControllerV1 } from "../controllers/bible.v1.controller.ts";
import { createBibleRoutesV1 } from "./bible.v1.route.ts";

export const createBibleRoutes = (version: ApiVersion) => {
  switch (version) {
    case "v1":
      return createBibleRoutesV1(new BibleControllerV1());
    default: {
      const exhaustiveCheck: never = version;
      throw new Error(`Unsupported bible API version: ${exhaustiveCheck}`);
    }
  }
};

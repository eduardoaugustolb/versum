import type { ApiVersion } from "../../api-version.ts";
import { ConsentLogControllerV1 } from "../controllers/consent-log.v1.controller.ts";
import { createConsentLogRoutesV1 } from "./consent-log.v1.route.ts";

export const createConsentLogRoutes = (version: ApiVersion) => {
  switch (version) {
    case "v1":
      return createConsentLogRoutesV1(new ConsentLogControllerV1());
    default: {
      const exhaustiveCheck: never = version;
      throw new Error(`Unsupported consent logs API version: ${exhaustiveCheck}`);
    }
  }
};

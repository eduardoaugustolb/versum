import { OpenAPIHono } from "@hono/zod-openapi";
import { createDiscoveryRoutesV1 } from "./discovery/routes/index.ts";
import { createJourneyRoutesV1 } from "./journey/routes/index.ts";

export const createReadingRoutes = () => {
  const router = new OpenAPIHono();

  router.route("/journey", createJourneyRoutesV1());
  router.route("/discovery", createDiscoveryRoutesV1());

  return router;
};

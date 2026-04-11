import { OpenAPIHono } from "@hono/zod-openapi";
import { createDiscoveryRoutes } from "./discovery/routes/index.ts";
import { createJourneyRoutes } from "./journey/routes/index.ts";

export const createReadingRoutes = () => {
  const router = new OpenAPIHono();

  router.route("/journey", createJourneyRoutes());
  router.route("/discovery", createDiscoveryRoutes());

  return router;
};

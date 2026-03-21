import { OpenAPIHono } from "@hono/zod-openapi";
import type { ApiVersion } from "./api-version.ts";
import { createAuthRoutes } from "./auth/routes/index.ts";
import { createBibleRoutes } from "./bible/routes/index.ts";
import { createUsersRoutes } from "./users/routes/index.ts";

export const createModulesRoutes = (version: ApiVersion) => {
  const router = new OpenAPIHono();

  router.route("/auth", createAuthRoutes(version));
  router.route("/users", createUsersRoutes(version));
  router.route("/public/bible", createBibleRoutes(version));

  return router;
};

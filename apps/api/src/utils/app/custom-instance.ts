import { OpenAPIHono } from "@hono/zod-openapi";
import type { Hono } from "hono";
import { logger } from "../logger/index.ts";
import { validationErrorHook } from "./errors/validation.hook.ts";
import { SetupListeners } from "./setups/setup-listeners.ts";
import { SetupMiddlewares } from "./setups/setup-middlewares.ts";
import { SetupPlugins } from "./setups/setup-plugins.ts";
import { SetupRoutes } from "./setups/setup-routes.ts";

export class App {
  public readonly hono = new OpenAPIHono({
    defaultHook: validationErrorHook,
  });

  constructor() {
    this.setup();
  }

  setup() {
    new SetupPlugins({ app: this.hono });
    new SetupListeners({ app: this.hono });
    new SetupMiddlewares({ app: this.hono });
    new SetupRoutes({ app: this.hono });
  }

  route(path: string, app: Hono) {
    this.hono.route(path, app);
    return this;
  }

  get fetch() {
    return this.hono.fetch;
  }

  showRoutes() {
    const seen = new Set<string>();

    const routes = this.hono.routes.filter((route) => {
      if (route.method === "ALL" || route.path === "/*") return false;

      const key = `${route.method}:${route.path}`;
      if (seen.has(key)) return false;

      seen.add(key);
      return true;
    });

    logger({ level: "info", icon: "🪧" }, "Routes:");
    routes.forEach((route) => {
      logger({ level: "info", icon: "  " }, `[${route.method}]`, route.path);
    });
  }

  start() {
    try {
      this.showRoutes();
    } catch (error) {
      logger("error", error);
    }
  }
}

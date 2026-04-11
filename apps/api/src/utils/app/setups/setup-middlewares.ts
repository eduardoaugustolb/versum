import type { OpenAPIHono } from "@hono/zod-openapi";
import { GlobalRateLimiter } from "@/middlewares/rate-limiter/middleware.ts";
import { debugRequestsMiddleware } from "../../../middlewares/debug-requests.middleware.ts";

export class SetupMiddlewares {
  private readonly app: OpenAPIHono;
  private readonly globalRateLimiter: GlobalRateLimiter;

  constructor({ app }: { app: OpenAPIHono }) {
    this.app = app;
    this.globalRateLimiter = new GlobalRateLimiter();
    this.setupGlobalMiddlewares();
  }

  setupGlobalMiddlewares() {
    this.app.use(this.globalRateLimiter.middleware);
    this.app.use(debugRequestsMiddleware);
  }
}

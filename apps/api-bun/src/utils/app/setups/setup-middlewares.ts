import { OpenAPIHono } from "@hono/zod-openapi";
import { debugRequestsMiddleware } from "../../../middlewares/debug-requests.middleware.ts";

export class SetupMiddlewares {
  private readonly app: OpenAPIHono;

  constructor({ app }: { app: OpenAPIHono }) {
    this.app = app;
    this.setupGlobalMiddlewares();
  }

  setupGlobalMiddlewares() {
    this.app.use(debugRequestsMiddleware);
  }
}

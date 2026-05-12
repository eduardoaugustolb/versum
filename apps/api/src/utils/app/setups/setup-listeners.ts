import type { OpenAPIHono } from "@hono/zod-openapi";
import { redis } from "../../../infrastructure/redis";
import { logger } from "@versum/logger";
import { ErrorHandler } from "../errors/index.ts";

export class SetupListeners {
  private readonly app: OpenAPIHono;

  constructor({ app }: { app: OpenAPIHono }) {
    this.app = app;
    this.setupErrorListener();
    this.setupRedisListener();
  }

  private setupRedisListener() {
    redis.connect().catch((err) => {
      logger({ level: "error" }, "[REDIS]", `Connection failed: ${err}`);
    });
  }

  private setupErrorListener() {
    this.app.onError((error, ctx) => {
      logger(
        {
          level: "error",
        },
        `Server error: ${error.message}`,
      );

      return new ErrorHandler({ ctx }).handle(error);
    });
  }
}

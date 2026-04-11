import type { OpenAPIHono } from "@hono/zod-openapi";
import { CustomRedisClient } from "../../../infrastructure/redis";
import { logger } from "../../logger/index.ts";
import { ErrorHandler } from "../errors/index.ts";

export class SetupListeners {
  private readonly app: OpenAPIHono;

  constructor({ app }: { app: OpenAPIHono }) {
    this.app = app;
    this.setupErrorListener();
    this.setupRedisListener();
  }

  private setupRedisListener() {
    CustomRedisClient.connectAll();
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

import { ErrorHandler } from "../errors/index.ts";
import { redis } from "../../../infrastructure/redis/index.ts";
import { logger } from "../../logger/index.ts";
import { OpenAPIHono } from "@hono/zod-openapi";

export class SetupListeners {
  private readonly app: OpenAPIHono;

  constructor({ app }: { app: OpenAPIHono }) {
    this.app = app;
    this.setupErrorListener();
    this.setupRedisListener();
  }

  private setupRedisListener() {
    redis.on("connect", () => {
      logger({ level: "success" }, "[REDIS]", "Connected");
    });

    redis.on("error", (error) => {
      logger({ level: "error" }, "[REDIS]", `Error: ${error.message}`);
    });

    redis.on("reconnecting", () => {
      logger({ level: "info" }, "[REDIS]", "Reconnecting...");
    });

    redis.on("close", () => {
      logger({ level: "info" }, "[REDIS]", "Disconnected");
    });

    redis.on("ready", () => {
      logger({ level: "success" }, "[REDIS]", "Ready");
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

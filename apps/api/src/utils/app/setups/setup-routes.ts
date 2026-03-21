import { OpenAPIHono } from "@hono/zod-openapi";
import { Scalar } from "@scalar/hono-api-reference";
import { DEFAULT_API_VERSION } from "../../../modules/api-version.ts";
import { createModulesRoutes } from "../../../modules/routes.ts";
import { logger } from "../../logger/index.ts";
import { getCookie } from "hono/cookie";

const scalarCss = await Bun.file("src/assets/css/scalar.css").text();

export class SetupRoutes {
  private readonly app: OpenAPIHono;

  constructor({ app }: { app: OpenAPIHono }) {
    this.app = app;
    this.setupRoutes();
  }

  setupRoutes() {
    this.app.get("/", (c) => {
      logger("info", getCookie(c, "__Host-session"));
      return c.text("Hello World");
    });
    this.app.route(
      `/api/${DEFAULT_API_VERSION}`,
      createModulesRoutes(DEFAULT_API_VERSION),
    );

    this.app.doc("/openapi.json", {
      openapi: "3.0.0",
      info: {
        title: "Versum API",
        version: "1.0.0",
      },
    });

    this.app.get(
      "/docs",
      Scalar({
        title: "Versum API",
        url: "/openapi.json",
        customCss: scalarCss,
      }),
    );
  }
}

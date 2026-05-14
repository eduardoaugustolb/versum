import type { OpenAPIHono } from "@hono/zod-openapi";
import { cors } from "hono/cors";
export class SetupPlugins {
  private readonly app: OpenAPIHono;

  constructor({ app }: { app: OpenAPIHono }) {
    this.app = app;
    this.setupPlugins();
  }

  setupOpenApiAuthCookie() {
    const isSecure = Bun.env.COOKIE_SECURE === "true";
    const cookieName = isSecure ? "__Host-session" : "session";

    this.app.openAPIRegistry.registerComponent("securitySchemes", "cookieAuth", {
      type: "apiKey",
      in: "cookie",
      name: cookieName,
      description: "Session token obtido via magic link. Enviado automaticamente pelo browser como cookie HttpOnly.",
    });
  }

  setupPlugins() {
    this.setupOpenApiAuthCookie();

    this.app.use(
      cors({
        origin: [
          "http://localhost:3000",
          "https://app.versum.eduardoaugusto.work",
        ],
        credentials: true,
        allowHeaders: ["Content-Type", "Authorization"],
        allowMethods: ["GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"],
        exposeHeaders: ["set-cookie"],
      }),
    );
  }
}

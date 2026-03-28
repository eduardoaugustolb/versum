import { OpenAPIHono } from "@hono/zod-openapi";
import { cors } from "hono/cors";
export class SetupPlugins {
  private readonly app: OpenAPIHono;

  constructor({ app }: { app: OpenAPIHono }) {
    this.app = app;
    this.setupPlugins();
  }

  setupPlugins() {
    this.app.use(cors({
      origin: [
        "http://localhost:3000",
        "https://app.versum.eduardoaugusto.work"
      ],
      credentials: true,
      allowHeaders: ["Content-Type", "Authorization"],
      allowMethods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"]
    }))
  }
}

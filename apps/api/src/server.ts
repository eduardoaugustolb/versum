import { Hono } from "hono";
import { App } from "@/utils/app/index.ts";
import { env } from "@/utils/env/parser.ts";

const helloWorld = new Hono().get("/", async (ctx) => {
  await new Promise((resolve) => setTimeout(resolve, 1000));
  return ctx.text("Hello world!");
});
const app = new App();

app.route("/", helloWorld);

app.start();

export default {
  port: env.PORT,
  fetch: app.fetch,
};

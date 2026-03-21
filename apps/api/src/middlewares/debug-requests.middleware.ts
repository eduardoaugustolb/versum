import type { Context, Next } from "hono";
import { logger } from "../utils/logger/index.ts";

class DebugRequestsMiddleware {
  async middleware(ctx: Context, next: Next) {
    const reqUrl = new URL(ctx.req.url);
    logger("info", "Request received:", `[${ctx.req.method}]`, reqUrl.pathname);
    logger("debug", "Request started:", `[${ctx.req.method}]`, reqUrl.pathname);
    const start = Date.now();
    const res = await next();
    const duration = Date.now() - start;
    logger(
      "debug",
      "Request completed:",
      `[${ctx.req.method}]`,
      reqUrl.pathname,
      `${duration}ms`,
    );
    return res;
  }
}

export const debugRequestsMiddleware = new DebugRequestsMiddleware().middleware;

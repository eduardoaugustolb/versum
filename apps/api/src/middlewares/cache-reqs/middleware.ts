import type { Context, Next } from "hono";
import { cacheReqsRedisConn as redis } from "../../infrastructure/redis";
import { logger } from "../../utils/logger/index.ts";

export class CacheMiddleware {
  private readonly config: {
    ttl: number;
    cacheKeyPrefix: string;
  };

  constructor({
    config,
  }: {
    config: {
      ttl: number;
      cacheKeyPrefix: string;
    };
  }) {
    this.config = config;
  }

  middleware = async (ctx: Context, next: Next) => {
    const { ttl, cacheKeyPrefix } = this.config;
    const url = new URL(ctx.req.url);
    const key = cacheKeyPrefix
      .replaceAll("{req}", ctx.req.url)
      .replaceAll("{method}", ctx.req.method)
      .replaceAll("{path}", url.pathname)
      .replaceAll("{query}", url.searchParams.toString());
    const ignoreCache =
      (url.hostname === "localhost" &&
        url.searchParams.get("ignoreCache") === "true") ||
      ctx.req.method !== "GET";

    if (ignoreCache) {
      logger(
        "info",
        "[CACHE MIDDLEWARE]",
        `[${ctx.req.method}] ${url.pathname}`,
      );
      return await next();
    }

    const cachedResponse = await redis.get(key);
    if (cachedResponse) {
      const acctualyTtl = await redis.ttl(key);
      const cacheExpireAt = new Date(
        Date.now() + acctualyTtl * 1000,
      ).toISOString();
      return ctx.json({
        ...JSON.parse(cachedResponse),
        cache: true,
        cacheExpireAt,
      });
    }

    await next();

    if (ctx.res.status === 200) {
      const body = await ctx.res.clone().json();

      await redis.set(key, JSON.stringify(body), "EX", ttl);
    }
  };
}

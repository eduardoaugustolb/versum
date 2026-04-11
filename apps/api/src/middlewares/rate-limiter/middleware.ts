import type { Context, Next } from "hono";
import { rateLimitRedisConn as redis } from "../../infrastructure/redis";

interface iConfig {
  windowMs?: number;
  limit?: number;
  keyGenerator?: (c: Context) => string;
}

export class RateLimiterMiddleware {
  private readonly windowMs: number;
  private readonly limit: number;
  private readonly getKey: (c: Context) => string;

  constructor(config: iConfig = {}) {
    this.windowMs = config.windowMs ?? 60000;
    this.limit = config.limit ?? 100;
    this.getKey = config.keyGenerator ?? this.defaultKeyGenerator;
    this.middleware = this.middleware.bind(this);
  }

  private defaultKeyGenerator(c: Context): string {
    const forwarded = c.req.header("x-forwarded-for");
    if (forwarded) {
      const firstIp = forwarded.split(",")[0];
      if (firstIp) {
        return firstIp.trim();
      }
    }

    const realIp = c.req.header("x-real-ip");
    if (realIp) {
      return realIp.trim();
    }

    return "unknown-ip";
  }

  async middleware(c: Context, next: Next) {
    try {
      const identifier = this.getKey(c);
      const windowStart = Math.floor(Date.now() / this.windowMs);
      const redisKey = `ratelimit:${identifier}:${c.req.path}:${windowStart}`;

      const currentCount = await redis.incr(redisKey);

      if (currentCount === 1) {
        await redis.expire(redisKey, this.windowMs / 1000);
      }

      c.header("X-RateLimit-Limit", String(this.limit));
      c.header(
        "X-RateLimit-Remaining",
        String(Math.max(0, this.limit - currentCount)),
      );

      if (currentCount > this.limit) {
        const ttl = await redis.ttl(redisKey);
        c.header("Retry-After", String(ttl));

        return c.json(
          {
            success: false,
            message: "Muitas requisições. Tente novamente em breve.",
            retryAfter: ttl,
          },
          429,
        );
      }

      await next();
    } catch (error) {
      console.error("Erro no rate limiter:", error);
      await next();
    }
  }
}

export class GlobalRateLimiter extends RateLimiterMiddleware {
  constructor() {
    super({
      windowMs: 60_000,
      limit: 100,
    });
  }
}

export class MagicLinkRateLimiter extends RateLimiterMiddleware {
  constructor() {
    super({
      windowMs: 60_000,
      limit: 5,
      keyGenerator: (c) => {
        const forwarded = c.req.header("x-forwarded-for");
        if (forwarded) {
          const firstIp = forwarded.split(",")[0];
          if (firstIp) {
            return firstIp.trim();
          }
        }
        const realIp = c.req.header("x-real-ip");
        if (realIp) {
          return realIp.trim();
        }
        return "unknown-ip";
      },
    });
  }
}

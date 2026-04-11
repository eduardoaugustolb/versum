import { env } from "@/utils/env";
import { CustomRedisClient } from "../client";

const url = new URL(env.REDIS_DATABASE_URL);
url.pathname = `rate_limit_${env.BUN_ENV}`;

export const rateLimitRedisConn = new CustomRedisClient(url.toString(), {
  maxRetries: 1,
});

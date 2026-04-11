import { env } from "../../../utils/env/index.ts";
import { CustomRedisClient } from "../index.ts";

const url = new URL(env.REDIS_DATABASE_URL);
url.pathname = `main_${env.BUN_ENV}`;

export const mainRedisConn = new CustomRedisClient(url.toString(), {
  maxRetries: 1,
});

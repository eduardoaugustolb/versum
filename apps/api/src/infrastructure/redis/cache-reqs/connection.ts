import { env } from "../../../utils/env";
import { CustomRedisClient } from "..";

const url = new URL(env.REDIS_DATABASE_URL);
url.pathname = `cache_reqs_${env.BUN_ENV}`;

export const cacheReqsRedisConn = new CustomRedisClient(url.toString(), {
  maxRetries: 1,
});

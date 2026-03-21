import { env } from "../../utils/env/index.ts";
import { CustomRedisClient } from "./index.ts";

export const redis = new CustomRedisClient(env.REDIS_DATABASE_URL, {
  maxRetriesPerRequest: 1,
  enableReadyCheck: true,
});

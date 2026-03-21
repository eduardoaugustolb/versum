import { Redis, type RedisOptions } from "ioredis";
import { env } from "../../utils/env/index.ts";

const cert = await Bun.file(".certs/redis-certificate.pem").text();

export class CustomRedisClient extends Redis {
  constructor(connectionUri: string, options: RedisOptions) {
    super(connectionUri, getCustomRedisClientOptions(options));
  }
}

function getCustomRedisClientOptions(options: RedisOptions): RedisOptions {
  if (!env.REDIS_DATABASE_CERT_PATH || options.tls) return options;
  if (!cert) return options;
  return {
    ...options,
    tls: {
      key: cert,
      cert,
      ca: cert,
      rejectUnauthorized: false,
    },
  };
}

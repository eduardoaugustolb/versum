import { RedisClient, type RedisOptions } from "bun";
import { logger } from "@/utils/logger/logger.ts";
import { env } from "../../utils/env/index.ts";

const clientsMap = new Map<string, CustomRedisClient>();

const cert = await Bun.file(".certs/redis-certificate.pem").text();

export class CustomRedisClient extends RedisClient {
  options: RedisOptions;
  connUri: string;
  constructor(connectionUri: string, options: RedisOptions) {
    super(connectionUri, getCustomRedisClientOptions(options));
    this.options = options;
    this.connUri = connectionUri;
    this.setupListeners = this.setupListeners.bind(this);
    this.setupListeners();
    clientsMap.set(connectionUri, this);
  }

  private setupListeners = () => {
    this.onconnect = () => {
      const connUrl = new URL(this.connUri);
      logger(
        { level: "success" },
        "[REDIS]",
        `Connected! Db: "${connUrl.pathname.replaceAll("/", "")}"`,
      );
    };

    this.onclose = (error) => {
      logger({ level: "error" }, "[REDIS]", `Error: ${error.message}`);
    };
  };

  static async connectAll() {
    await Promise.all(Array.from(clientsMap.values()).map((c) => c.connect()));
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

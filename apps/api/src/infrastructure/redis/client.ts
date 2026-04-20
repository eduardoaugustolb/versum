import { RedisClient, type RedisOptions } from "bun";
import { logger } from "@versum/logger";
import { env } from "../../utils/env/index.ts";

let clientsInitialized = false;

const cert = await Bun.file(".certs/redis-certificate.pem").text();

const keepAlive = (connUrl: URL, client: RedisClient) => {
  return setInterval(async () => {
    const pingResult = await client.ping();
    if (pingResult !== "PONG") {
      logger(
        { level: "error" },
        "[REDIS]",
        `Ping failed! Db: "${connUrl.pathname.replaceAll("/", "")}"`,
      );

      const errorEmbed = {
        title: "Redis Ping Failed",
        description: `Db: "${connUrl.pathname.replaceAll("/", "")}"`,
        color: 16711680,
      };

      const response = await fetch(`${env.DISCORD_WEBHOOK_URL}?wait=true`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ embeds: [errorEmbed] }),
      });

      if (!response.ok) {
        throw new Error(`Failed to send ping error message: ${response.statusText}`);
      }
    }
  }, 10000);
}

export class CustomRedisClient extends RedisClient {
  options: RedisOptions;
  connUrl: URL;
  static instances: (CustomRedisClient & {connUrl: URL})[] = [];

  constructor(connectionUri: string, options: RedisOptions) {
    // const uri = new URL(connectionUri).pathname.replaceAll("/", "");
    super(connectionUri, getCustomRedisClientOptions(options));
    this.options = options;
    this.connUrl = new URL(connectionUri);
    this.setupListeners = this.setupListeners.bind(this);
    this.setupListeners();

    const existingIndex = CustomRedisClient.instances.findIndex(
      (c) => c.connUrl.href === connectionUri,
    );
    if (existingIndex >= 0) {
      CustomRedisClient.instances[existingIndex] = this;
    } else {
      CustomRedisClient.instances.push(this);
    }
  }

  private setupListeners = () => {
    let keepAliveInterval: NodeJS.Timeout | null = null;
    this.onconnect = () => {

      if (!CustomRedisClient.instances.some(c => c.connected === false)) {
        logger("success", "[REDIS]", `All connections established. Databases: ${CustomRedisClient.instances.map(c => c.connUrl.pathname.replaceAll("/", "")).join(", ")}`)
      }

      keepAliveInterval = keepAlive(this.connUrl, this);
    };

    this.onclose = async (error) => {
      const dbName = this.connUrl.pathname.replaceAll("/", "");

      logger({ level: "error" }, "[REDIS]", `Connection failed: ${error.message}. Db: "${dbName}"`);
      if (keepAliveInterval) clearInterval(keepAliveInterval);
    };
  };

  static async connectAll() {
    if (clientsInitialized) return;
    clientsInitialized = true;

    await Promise.all(
      CustomRedisClient.instances.map((c) => c.connect()),
    );
  }

  static resetForHotReload() {
    clientsInitialized = false;
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

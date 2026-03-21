import { drizzle } from "drizzle-orm/node-postgres";
import { env } from "../../utils/env/index.ts";
import { Pool } from "pg";
import * as schema from "./schema.ts";

const cert = await Bun.file(".certs/postgre-certificate.pem").text();

const dbUrl = new URL(env.DATABASE_URL!);

const pgPool = new Pool({
  connectionString: env.DATABASE_URL,
  host: dbUrl.hostname,
  port: parseInt(dbUrl.port),
  user: dbUrl.username,
  password: dbUrl.password,
  database: dbUrl.pathname.replace("/", ""),
  ssl: {
    ca: cert,
    cert,
    key: cert,
    rejectUnauthorized: false,
  },
});

export const db = drizzle(pgPool, {
  schema,
  logger: true,
});

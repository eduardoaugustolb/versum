import { defineConfig } from "drizzle-kit";
import { readFileSync } from "fs";

const cert = readFileSync(".certs/postgre-certificate.pem", "utf-8");

const dbUrl = new URL(process.env.DATABASE_URL!);

export default defineConfig({
  out: "./drizzle",
  schema: "./src/infrastructure/db/schema.ts",
  dialect: "postgresql",
  dbCredentials: {
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
  },
});

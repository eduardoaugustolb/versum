import { defineConfig } from "drizzle-kit";

const cert = Deno.readTextFileSync(".certs/postgre-certificate.pem");

const dbUrl = new URL(Deno.env.get("DATABASE_URL")!);

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

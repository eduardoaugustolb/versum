import { envSchema } from "./schema.ts";
export const mockEnv = {
  DATABASE_URL: "postgresql://test",
  DATABASE_CERT_PATH: "/test",
  REDIS_DATABASE_URL: "redis://test",
  REDIS_DATABASE_CERT_PATH: "/test",
  PORT: "3000",
  BUN_ENV: "test",
  APP_URL: "http://test",
  ENCRYPT_SECRET: "test-secret-key-min-32-chars-long!!",
  RESEND_API_KEY: "test",
  WEB_CLIENT_APP_URL: "http://localhost:3000",
};

export const env = envSchema.parse(
  Bun.env.BUN_ENV === "test" ? mockEnv : Bun.env,
);

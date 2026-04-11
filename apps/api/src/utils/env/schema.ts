import { z } from "zod";

export const envSchema = z.object({
  DATABASE_URL: z.string(),
  DATABASE_CERT_PATH: z.string(),
  REDIS_DATABASE_URL: z.string(),
  REDIS_DATABASE_CERT_PATH: z.string(),
  PORT: z.string(),
  BUN_ENV: z.string(),
  ENCRYPT_SECRET: z.string(),
  RESEND_API_KEY: z.string(),
  DEBUG: z.string().optional(),
  WEB_CLIENT_APP_URL: z.url(),
});

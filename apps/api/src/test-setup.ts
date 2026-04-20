import { vi } from "vitest";

const env = {
  DATABASE_URL: "postgresql://test:test@localhost:5432/test",
  DATABASE_CERT_PATH: "/test",
  REDIS_DATABASE_URL: "redis://localhost:6379",
  REDIS_DATABASE_CERT_PATH: "/test",
  PORT: "3000",
  BUN_ENV: "test",
  DISCORD_WEBHOOK_URL: "https://discord.com/api/webhooks/webhook_id/webhook_token",
  ENCRYPT_SECRET: "test-secret-key-min-32-chars-long!!",
  RESEND_API_KEY: "test",
  WEB_CLIENT_APP_URL: "http://localhost:3000",
};

Object.defineProperty(globalThis, "Bun", {
  value: {
    env,
    file: vi.fn().mockReturnValue({
      text: vi.fn().mockResolvedValue("mock-cert-content"),
    }),
  },
  writable: true,
});

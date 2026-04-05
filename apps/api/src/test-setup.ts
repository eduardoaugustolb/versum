import { vi } from "vitest";

Object.defineProperty(globalThis, "Bun", {
  value: {
    env: {
      DATABASE_URL: "postgresql://test:test@localhost:5432/test",
      DATABASE_CERT_PATH: "/test",
      REDIS_DATABASE_URL: "redis://localhost:6379",
      REDIS_DATABASE_CERT_PATH: "/test",
      PORT: "3000",
      BUN_ENV: "test",
      APP_URL: "http://localhost:3000",
      ENCRYPT_SECRET: "test-secret-key-min-32-chars-long!!",
      RESEND_API_KEY: "test",
      WEB_CLIENT_APP_URL: "http://localhost:3000",
    },
    file: vi.fn().mockReturnValue({
      text: vi.fn().mockResolvedValue("mock-cert-content"),
    }),
  },
  writable: true,
  configurable: true,
});

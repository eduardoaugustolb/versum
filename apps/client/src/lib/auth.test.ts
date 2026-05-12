import { describe, it, expect, afterEach } from "vitest";
import { getSessionCookieName } from "./auth";

describe("getSessionCookieName", () => {
  const ORIGINAL = process.env.NEXT_PUBLIC_API_URL;

  afterEach(() => {
    process.env.NEXT_PUBLIC_API_URL = ORIGINAL;
  });

  it("retorna __Host-session para HTTPS", () => {
    process.env.NEXT_PUBLIC_API_URL = "https://api.versum.work";
    expect(getSessionCookieName()).toBe("__Host-session");
  });

  it("retorna session para HTTP", () => {
    process.env.NEXT_PUBLIC_API_URL = "http://localhost:4002";
    expect(getSessionCookieName()).toBe("session");
  });

  it("retorna session se NEXT_PUBLIC_API_URL for undefined", () => {
    delete process.env.NEXT_PUBLIC_API_URL;
    expect(getSessionCookieName()).toBe("session");
  });
});

import { describe, it, expect } from "vitest";
import { guardRoutes, DEFAULT_GUARD } from "./config";

describe("guardRoutes", () => {
  it("possui 3 rotas", () => {
    expect(guardRoutes).toHaveLength(3);
  });

  it("login usa startWith e kind guest", () => {
    const login = guardRoutes.find((r) => "startWith" in r && r.startWith === "/login");
    expect(login).toBeDefined();
    expect(login!.kind).toBe("guest");
    expect(login!.redirectTo).toBe("/");
  });

  it("magic-link usa path exato e kind guest", () => {
    const ml = guardRoutes.find((r) => "path" in r && r.path === "/auth/magic-link");
    expect(ml).toBeDefined();
    expect(ml!.kind).toBe("guest");
    expect(ml!.redirectTo).toBe("/");
  });

  it("onboarding usa path exato e kind onboarding", () => {
    const onboarding = guardRoutes.find((r) => "path" in r && r.path === "/onboarding");
    expect(onboarding).toBeDefined();
    expect(onboarding!.kind).toBe("onboarding");
    expect(onboarding!.redirectTo).toBe("/");
  });
});

describe("DEFAULT_GUARD", () => {
  it("authRedirectTo aponta para /login", () => {
    expect(DEFAULT_GUARD.authRedirectTo).toBe("/login");
  });

  it("onboardingRedirectTo aponta para /onboarding", () => {
    expect(DEFAULT_GUARD.onboardingRedirectTo).toBe("/onboarding");
  });
});

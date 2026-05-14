import { describe, it, expect } from "vitest";
import { matchGuard, evaluateGuard } from "./guards";

describe("matchGuard", () => {
  it("retorna a guard pelo path exato", () => {
    const guard = matchGuard("/onboarding");
    expect(guard).toBeDefined();
    expect(guard!.kind).toBe("onboarding");
  });

  it("retorna a guard pelo startWith", () => {
    const guard = matchGuard("/login/test");
    expect(guard).toBeDefined();
    expect(guard!.kind).toBe("guest");
  });

  it("retorna undefined para rota desconhecida", () => {
    expect(matchGuard("/desconhecido")).toBeUndefined();
  });
});

describe("evaluateGuard", () => {
  describe("kind: guest", () => {
    const guard = { kind: "guest" as const, path: "/login", redirectTo: "/" };

    it("redirect se autenticado", () => {
      expect(evaluateGuard(guard, { isAuthenticated: true, onboardingComplete: false })).toBe("/");
    });

    it("null se não autenticado", () => {
      expect(evaluateGuard(guard, { isAuthenticated: false, onboardingComplete: false })).toBeNull();
    });
  });

  describe("kind: auth", () => {
    const guard = { kind: "auth" as const, path: "/dashboard", redirectTo: "/login" };

    it("null se autenticado", () => {
      expect(evaluateGuard(guard, { isAuthenticated: true, onboardingComplete: true })).toBeNull();
    });

    it("redirect se não autenticado", () => {
      expect(evaluateGuard(guard, { isAuthenticated: false, onboardingComplete: false })).toBe("/login");
    });
  });

  describe("kind: onboarding", () => {
    const guard = { kind: "onboarding" as const, path: "/onboarding", redirectTo: "/" };

    it("redirect se autenticado e onboarding completo", () => {
      expect(evaluateGuard(guard, { isAuthenticated: true, onboardingComplete: true })).toBe("/");
    });

    it("redirect se não autenticado", () => {
      expect(evaluateGuard(guard, { isAuthenticated: false, onboardingComplete: false })).toBe("/");
    });

    it("null se autenticado e onboarding incompleto (pode acessar)", () => {
      expect(evaluateGuard(guard, { isAuthenticated: true, onboardingComplete: false })).toBeNull();
    });
  });

  describe("kind: onboarding-complete", () => {
    const guard = { kind: "onboarding-complete" as const, path: "/app", redirectTo: "/onboarding" };

    it("null se autenticado e onboarding completo", () => {
      expect(evaluateGuard(guard, { isAuthenticated: true, onboardingComplete: true })).toBeNull();
    });

    it("redirect se não autenticado", () => {
      expect(evaluateGuard(guard, { isAuthenticated: false, onboardingComplete: false })).toBe("/onboarding");
    });

    it("redirect se autenticado mas onboarding incompleto", () => {
      expect(evaluateGuard(guard, { isAuthenticated: true, onboardingComplete: false })).toBe("/onboarding");
    });
  });
});

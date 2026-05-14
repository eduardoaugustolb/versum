import type { GuardRoute, EvaluateGuardContext } from "./types";
import { guardRoutes } from "./config";

export function matchGuard(pathname: string): GuardRoute | undefined {
  return guardRoutes.find((r) =>
    "startWith" in r
      ? pathname.startsWith(r.startWith)
      : r.path === pathname,
  );
}

export function evaluateGuard(
  guard: GuardRoute,
  ctx: EvaluateGuardContext,
): string | null {
  switch (guard.kind) {
    case "guest":
      return ctx.isAuthenticated ? guard.redirectTo : null;
    case "auth":
      return ctx.isAuthenticated ? null : guard.redirectTo;
    case "onboarding":
      if (!ctx.isAuthenticated) return guard.redirectTo;
      if (ctx.onboardingComplete) return guard.redirectTo;
      return null;
    case "onboarding-complete":
      if (!ctx.isAuthenticated) return guard.redirectTo;
      if (!ctx.onboardingComplete) return guard.redirectTo;
      return null;
  }
}

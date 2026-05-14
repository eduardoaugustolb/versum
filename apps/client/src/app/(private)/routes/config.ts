import type { GuardRoute } from "./types";

export const guardRoutes: GuardRoute[] = [
  { kind: "guest", startWith: "/login", redirectTo: "/" },
  { kind: "guest", path: "/auth/magic-link", redirectTo: "/" },
  { kind: "onboarding", path: "/onboarding", redirectTo: "/" },
];

export const DEFAULT_GUARD = {
  authRedirectTo: "/login",
  onboardingRedirectTo: "/onboarding",
} as const;

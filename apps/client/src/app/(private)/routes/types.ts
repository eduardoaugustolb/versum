export type RouteMatch =
  | { path: string }
  | { startWith: string };

export type GuardRoute = RouteMatch & {
  kind: "guest" | "auth" | "onboarding" | "onboarding-complete";
  redirectTo: string;
};

export type EvaluateGuardContext = {
  isAuthenticated: boolean;
  onboardingComplete: boolean;
};

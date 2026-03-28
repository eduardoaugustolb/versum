type Route = {
  startWith: string;
  redirectTo: string;
} & (
  {  path: string;
} | { startWith: string }
);

type ProtectedRoute = Route & ({
  requiresAuth: boolean;
} | {
  requiresGhest: boolean;
});

export const protectedRoutes: ProtectedRoute[] = [
  {
    startWith: "/login",
    redirectTo: "/",
    requiresGhest: true,
  },
  {
    startWith: "/auth/magic-link",
    redirectTo: "/",
    requiresGhest: true,
  },
];

export function findProtectedRoute({ pathname }: { pathname: string }) {
  return protectedRoutes.find((r) => r.startWith ? pathname.startsWith(r.startWith) : "path" in r ? r.path === pathname : false)
}

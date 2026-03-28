import { cookies, headers } from "next/headers";
import { redirect } from "next/navigation";
import { getCachedUserAuth } from "@/dal/auth/get-cached-session";
import { findProtectedRoute } from "./protected-routes";

export default async function Layout({
  children,
}: {
  children: React.ReactNode;
}) {
  const pathname = (await headers()).get("x-invoke-path");
  const sessionCookie = (await cookies()).get("__Host-session");
  console.debug("pathname", pathname, "sessionCookie", sessionCookie);
  if (!pathname) return null;

  const route = findProtectedRoute({ pathname });

  if (!route) return null;

  if (sessionCookie && "requiresGhest" in route && route.requiresGhest) {
    redirect(route.redirectTo);
  }
  if (sessionCookie && "requiresAuth" in route && route.requiresAuth) {
    try {
      const userAuth = await getCachedUserAuth();
      if (sessionCookie && !userAuth) {
        redirect(route.redirectTo);
      }
    } catch {
      redirect(route.redirectTo);
    }
  }

  return children;
}

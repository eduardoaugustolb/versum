import { cookies, headers } from "next/headers";
import { redirect } from "next/navigation";
import { matchGuard, evaluateGuard, DEFAULT_GUARD } from "./routes";
import { getSessionCookieName } from "@/lib/auth";
import { getCachedUserAuth } from "@/dal/auth/get-cached-session";

export default async function Layout({
  children,
}: {
  children: React.ReactNode;
}) {
  const sessionCookieName = getSessionCookieName();
  const pathname = (await headers()).get("x-invoke-path") ?? "/";
  const sessionCookie = (await cookies()).get(sessionCookieName);

  const isAuthenticated = !!sessionCookie;

  let onboardingComplete = false;
  if (isAuthenticated) {
    const userData = await getCachedUserAuth(sessionCookie?.value);

    if (!userData) {
      redirect("/auth/clear-session");
    }

    onboardingComplete = userData?.data?.onboardingIsCompleted ?? false;
  }

  const guard = matchGuard(pathname);

  if (guard) {
    const redirectTo = evaluateGuard(guard, { isAuthenticated, onboardingComplete });
    if (redirectTo) redirect(redirectTo);
  } else {
    if (!isAuthenticated) redirect(DEFAULT_GUARD.authRedirectTo);
    if (!onboardingComplete) redirect(DEFAULT_GUARD.onboardingRedirectTo);
  }

  return children;
}

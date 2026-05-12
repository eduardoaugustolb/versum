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
  const pathname = (await headers()).get("next-url") ?? "/";
  const sessionCookie = (await cookies()).get(sessionCookieName);

  const isAuthenticated = !!sessionCookie;

  let onboardingComplete = false;
  if (isAuthenticated) {
    try {
      const userData = await getCachedUserAuth(
        sessionCookie?.value,
        process.env.NEXT_PUBLIC_API_URL,
      );
      onboardingComplete = userData?.onboardingIsCompleted ?? false;
    } catch (error) {
      console.error("Erro ao verificar onboarding:", error);
    }
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

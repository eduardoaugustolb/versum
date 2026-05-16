import { type NextRequest, NextResponse } from "next/server";
import { getSessionCookieName } from "@/lib/auth";

export async function GET(req: NextRequest) {
  const sessionCookieName = getSessionCookieName();
  const authCookie = req.cookies.get(sessionCookieName);

  if (authCookie) {
    return NextResponse.redirect("/");
  }

  const url = req.nextUrl.clone();
  const params = url.searchParams;
  const token = params.get("token") ? String(params.get("token")) : undefined;

  if (!token) {
    return NextResponse.redirect("/login");
  }

  const routeRedirect = url.clone();
  routeRedirect.searchParams.delete("token");
  routeRedirect.pathname = "/";

  try {
    const forwardedFor = req.headers.get("x-forwarded-for");
    const realIp = req.headers.get("x-real-ip");
    const userAgent = req.headers.get("user-agent") ?? "unknown";
    const clientIp = forwardedFor?.split(",")[0]?.trim() ?? realIp ?? "unknown";

    const apiUrl = new URL("/api/v1/auth/magic-link", process.env.NEXT_PUBLIC_API_URL);
    apiUrl.searchParams.set("token", token);

    const res = await fetch(apiUrl.toString(), {
      method: "GET",
      credentials: "include",
      headers: {
        "x-forwarded-for": clientIp,
        "user-agent": userAgent,
      },
    });

    if (!res.ok) {
      throw new Error(`HTTP ${res.status}`);
    }

    const redirectRes = NextResponse.redirect(routeRedirect.toString());

    const setCookie = res.headers.get("set-cookie");
    if (setCookie) {
      redirectRes.headers.set("set-cookie", setCookie);
    }

    return redirectRes;
  } catch {
    routeRedirect.searchParams.set("error", "invalid_or_expired_token");
    return NextResponse.redirect(routeRedirect.toString());
  }
}

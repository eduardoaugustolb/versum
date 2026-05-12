import { type NextRequest, NextResponse } from "next/server";
import { matchGuard } from "./app/(private)/routes";
import { getSessionCookieName } from "./lib/auth";

class Response {
  res?: NextResponse;

  next({ url }: { url: URL }) {
    console.debug("🛡️ Proxy executando para:", `"${url.pathname}"`);
    this.res = NextResponse.next();

    this.setupHeaders({ url, res: this.res });

    return this.res;
  }

  redirect(to: URL, { url }: { url: URL }) {
    console.debug("🛡️ Redirecionando para:", to.toString());
    this.res = NextResponse.redirect(to);
    this.setupHeaders({ url, res: this.res });

    return this.res;
  }

  private setupHeaders({ url, res }: { url: URL; res: NextResponse }) {
    res.headers.set("x-invoke-path", url.pathname);
    res.headers.set("X-Frame-Options", "DENY");
    res.headers.set("X-Content-Type-Options", "nosniff");
    res.headers.set("Referrer-Policy", "strict-origin-when-cross-origin");
    res.headers.set(
      "Strict-Transport-Security",
      "max-age=31536000; includeSubDomains; preload",
    );
  }
}

export default function proxy(req: NextRequest) {
  const response = new Response();

  const isAuthenticated = !!req.cookies.get(getSessionCookieName());
  const url = req.nextUrl.clone();
  console.debug(
    "🛡️ Proxy executando para:",
    `"${url.pathname}"`,
    "| Autenticado:",
    isAuthenticated,
  );

  if (!url.pathname) {
    console.debug("🛡️ Pathname vazio, redirecionando para:", "/");
    return response.redirect(new URL("/", url), {
      url,
    });
  }

  const guard = matchGuard(url.pathname);

  if (guard?.kind === "guest" && isAuthenticated) {
    const routeRedirect = url.clone();
    routeRedirect.pathname = guard.redirectTo;
    console.debug("🛡️ Redirecionando para (guest):", routeRedirect);
    return response.redirect(routeRedirect, { url });
  }

  return response.next({ url });
}

export const config = {
  matcher: [
    // Ignore Next.js internals, static assets, favicon and API routes
    "/((?!api|_next|.*\\.[\\w]+$).*)",
  ],
};

import type {
  Client,
  RequestConfig,
  ResponseConfig,
} from "@kubb/plugin-client/clients/fetch";
import { type NextRequest, NextResponse } from "next/server";
import { getApiV1AuthMagicLink } from "@/lib/kubb/gen";

export async function GET(req: NextRequest) {
  const authCookie = req.cookies.get("__Host-session");

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
    let rawResponse: Response | undefined;

    const forwardedFor = req.headers.get("x-forwarded-for");
    const realIp = req.headers.get("x-real-ip");
    const userAgent = req.headers.get("user-agent") ?? "unknown";
    const clientIp = forwardedFor?.split(",")[0]?.trim() ?? realIp ?? "unknown";

    const customClient: Client = async <T>(
      config: RequestConfig,
    ): Promise<ResponseConfig<T>> => {
      const normalizedParams = new URLSearchParams();
      Object.entries(config.params || {}).forEach(([key, value]) => {
        if (value !== undefined)
          normalizedParams.append(key, value === null ? "null" : String(value));
      });

      let targetUrl = [config.baseURL, config.url].filter(Boolean).join("");
      if (config.params) targetUrl += `?${normalizedParams}`;

      const headers = new Headers(config.headers);
      headers.set("x-forwarded-for", clientIp);
      headers.set("user-agent", userAgent);

      rawResponse = await fetch(targetUrl, {
        method: config.method?.toUpperCase(),
        signal: config.signal,
        headers,
        credentials: "include",
      });

      return {
        data:
          [204, 205, 304].includes(rawResponse.status) || !rawResponse.body
            ? ({} as T)
            : await rawResponse.json(),
        status: rawResponse.status,
        statusText: rawResponse.statusText,
        headers: rawResponse.headers,
      };
    };

    await getApiV1AuthMagicLink(
      { token },
      {
        client: customClient,
      },
    );

    const res = NextResponse.redirect(routeRedirect.toString());

    const setCookie = rawResponse?.headers?.get("set-cookie");
    if (setCookie) {
      res.headers.set("set-cookie", setCookie);
    }

    return res;
  } catch {
    return NextResponse.redirect(routeRedirect.toString());
  }
}

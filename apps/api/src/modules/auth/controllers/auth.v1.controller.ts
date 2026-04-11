import type { Context } from "hono";
import { getCookie, setCookie } from "hono/cookie";
import { BadRequestError } from "../../../utils/app/errors/index.ts";
import { AuthServiceV1 } from "../services/auth.v1.service.ts";

const isSecure = Bun.env.COOKIE_SECURE === "true";
const cookieOptions = {
  httpOnly: true,
  secure: isSecure,
  sameSite: (isSecure ? "strict" : "lax") as "strict" | "lax",
  path: "/",
};

export class AuthControllerV1 {
  private readonly service: AuthServiceV1;

  constructor({ service }: { service?: AuthServiceV1 } = {}) {
    this.service = service ?? new AuthServiceV1();
  }

  authenticateWithMagicLink = async (c: Context) => {
    const token = c.req.query("token") ?? "";

    const sessionToken = await this.service.createSessionWithMagicLink({
      bruteMagicLinkToken: token,
      ip: c.req.header("x-forwarded-for") ?? "unknown",
      userAgent: c.req.header("user-agent") ?? "unknown",
    });

    const cookieName = isSecure ? "__Host-session" : "session";
    setCookie(c, cookieName, sessionToken, {
      ...cookieOptions,
      expires: new Date(Date.now() + 1000 * 60 * 60 * 24 * 30),
    });

    return c.json({ message: "Logged in successfully!" }, 200);
  };

  createAndSendMagicLink = async (c: Context) => {
    const { email } = (await c.req.json()) as { email: string };
    const magicLink = await this.service.createMagicLink({
      email,
    });

    await this.service.sendMagicLink({
      email,
      magicLink,
    });

    return c.json({ message: "Magic link sent!" }, 200);
  };

  logout = async (c: Context) => {
    const cookieName = isSecure ? "__Host-session" : "session";
    const sessionCookie = getCookie(c, cookieName);
    if (!sessionCookie) throw new BadRequestError("Session cookie not found");
    const sessionPublicId = sessionCookie.split(".")[0];
    if (!sessionPublicId) throw new BadRequestError("Invalid session cookie");
    await this.service.revokeSession({ sessionPublicId });
    setCookie(c, cookieName, "", {
      ...cookieOptions,
      expires: new Date(0),
    });
    return c.json({ message: "Logged out successfully!" }, 200);
  };
}

import type { Context } from "hono";
import { AuthServiceV1 } from "../services/auth.v1.service.ts";
import { setCookie, getCookie } from "hono/cookie";
import { BadRequestError } from "../../../utils/app/errors/index.ts";

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

    setCookie(c, "__Host-session", sessionToken, {
      httpOnly: true,
      secure: true,
      sameSite: "strict",
      expires: new Date(Date.now() + 1000 * 60 * 60 * 24 * 30), // One month
      path: "/",
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
    const sessionCookie = getCookie(c, "__Host-session");
    if (!sessionCookie) throw new BadRequestError("Session cookie not found");
    const sessionPublicId = sessionCookie.split(".")[0];
    await this.service.revokeSession({ sessionPublicId });
    setCookie(c, "__Host-session", "", {
      httpOnly: true,
      secure: true,
      sameSite: "strict",
      expires: new Date(0),
      path: "/",
    });
    return c.json({ message: "Logged out successfully!" }, 200);
  };
}

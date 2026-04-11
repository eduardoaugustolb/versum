import type { Context, Next } from "hono";
import { getCookie, setCookie } from "hono/cookie";
import { ValidateSession } from "../modules/auth/helpers/validate-session.ts";
import { AuthRepository } from "../modules/auth/repositories/auth.repository.ts";
import { AuthServiceV1 } from "../modules/auth/services/auth.v1.service.ts";
import { UnauthorizedError } from "../utils/app/errors/index.ts";
import { logger } from "../utils/logger/index.ts";

const isSecure = Bun.env.COOKIE_SECURE === "true";
const cookieName = isSecure ? "__Host-session" : "session";

export class AuthMiddleware {
  private readonly service: AuthServiceV1;
  private readonly repo: AuthRepository;
  constructor({
    repository,
    service,
  }: {
    service?: AuthServiceV1;
    repository?: AuthRepository;
  } = {}) {
    this.service = service ?? new AuthServiceV1();
    this.repo = repository ?? new AuthRepository();
    this.validateSession = this.validateSession.bind(this);
  }

  async validateSession(ctx: Context, next: Next) {
    const sessionCookie = getCookie(ctx, cookieName);
    logger(
      "debug",
      "[Auth] cookieName:",
      cookieName,
      "sessionCookie:",
      sessionCookie,
    );
    if (!sessionCookie) throw new UnauthorizedError("Session cookie not found");
    const [sessionCookiePublicId] = sessionCookie.split(".");

    if (!sessionCookiePublicId)
      throw new UnauthorizedError("Invalid session cookie");

    const session = await this.repo.getSessionByPublicId({
      publicId: sessionCookiePublicId,
    });
    logger("debug", "[Auth] session from publicId:", session?.id);

    const validatedSession = new ValidateSession({
      session,
      forwardedFor: ctx.req.header("x-forwarded-for") ?? "unknown",
    });
    logger("debug", "[Auth] session validated, calling refreshSession");

    const refreshSessionResult = await this.service.refreshSession({
      sessionId: validatedSession.session.id,
    });

    logger(
      "debug",
      "[Auth] refreshSession done, rotated:",
      refreshSessionResult.rotated,
    );

    if (refreshSessionResult.rotated) {
      const { session, token } = refreshSessionResult;
      setCookie(ctx, cookieName, token, {
        path: "/",
        httpOnly: true,
        secure: isSecure,
        sameSite: isSecure ? "Strict" : "Lax",
        maxAge: 60 * 60 * 24 * 30, // 30 days
      });

      ctx.set("session", session);
    } else {
      ctx.set("session", validatedSession.session);
    }

    await next();
  }
}

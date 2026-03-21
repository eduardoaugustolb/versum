import type { Context } from "hono";
import { AuthServiceV1 } from "../modules/auth/services/auth.v1.service.ts";
import { getCookie } from "hono/cookie";
import { UnauthorizedError } from "../utils/app/errors/index.ts";
import { AuthRepository } from "../modules/auth/repositories/auth.repository.ts";
import { ValidateSession } from "../modules/auth/helpers/validate-session.ts";
import { setCookie } from "hono/cookie";

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
  }

  async validateSession(ctx: Context) {
    const sessionCookie = getCookie(ctx, "__Host-session");
    if (!sessionCookie) throw new UnauthorizedError("Session cookie not found");
    const [sessionCookiePublicId] = sessionCookie.split(".");
    const session = await this.repo.getSessionByPublicId({
      publicId: sessionCookiePublicId,
    });
    const validatedSession = new ValidateSession({
      session,
      forwardedFor: ctx.req.header("x-forwarded-for") ?? "unknown",
    });

    const refreshSessionResult = await this.service.refreshSession({
      sessionId: validatedSession.session.id,
    });

    if (refreshSessionResult.rotated) {
      const { session, token } = refreshSessionResult;
      setCookie(ctx, "__Host-session", token, {
        path: "/",
        httpOnly: true,
        secure: true,
        sameSite: "Strict",
        maxAge: 60 * 60 * 24 * 30, // 30 days
      });

      ctx.set("session", session);
      return session;
    }

    ctx.set("session", validatedSession.session);

    return validatedSession.session;
  }
}

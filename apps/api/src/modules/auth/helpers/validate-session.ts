import {
  NotFoundError,
  UnauthorizedError,
} from "../../../utils/app/errors/index.ts";
import type { Session } from "../repositories/auth.types.repository.ts";

export class ValidateSession {
  session: Session;

  constructor({
    session,
    forwardedFor,
  }: {
    session?: Session | null;
    forwardedFor?: string;
  }) {
    if (!session) throw new NotFoundError("Session not found");
    if (session.expiresAt.getTime() < Date.now())
      throw new UnauthorizedError("Session expired");
    if (session.revokedAt && session.revokedAt.getTime() < Date.now())
      throw new UnauthorizedError("Session revoked");
    // if (session.ip !== forwardedFor)
    //   throw new UnauthorizedError("Session IP mismatch");

    this.session = session;
  }
}

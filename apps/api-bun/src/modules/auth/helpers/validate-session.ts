import {
  NotFoundError,
  UnauthorizedError,
} from "../../../utils/app/errors/index.ts";
import { Session } from "../repositories/auth.types.repository.ts";

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
    if (session.expires_at.getTime() < Date.now())
      throw new UnauthorizedError("Session expired");
    if (session.revoked_at && session.revoked_at.getTime() < Date.now())
      throw new UnauthorizedError("Session revoked");
    if (session.ip !== forwardedFor)
      throw new UnauthorizedError("Session IP mismatch");

    this.session = session;
  }
}

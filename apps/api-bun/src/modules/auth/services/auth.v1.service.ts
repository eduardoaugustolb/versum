import { EmailProvider } from "../../../infrastructure/resend/email-provider.ts";
import {
  NotFoundError,
  UnauthorizedError,
} from "../../../utils/app/errors/index.ts";
import { env } from "../../../utils/env/index.ts";
import { UserRepository } from "../../users/repositoties/user.repository.ts";
import { AuthRepository } from "../repositories/auth.repository.ts";
import argon2 from "argon2";
import { Session } from "../repositories/auth.types.repository.ts";
import { ValidateSession } from "../helpers/validate-session.ts";

export class AuthServiceV1 {
  private readonly repository: AuthRepository;
  private readonly userRepository: UserRepository;
  private readonly emailProvider: EmailProvider;

  constructor({
    repository,
    emailProvider,
    userRepository,
  }: {
    repository?: AuthRepository;
    emailProvider?: EmailProvider;
    userRepository?: UserRepository;
  } = {}) {
    this.repository = repository ?? new AuthRepository();
    this.userRepository = userRepository ?? new UserRepository();
    this.emailProvider = emailProvider ?? new EmailProvider();
  }

  async createMagicLink({ email }: { email: string }) {
    const token = this.generateRandomToken();

    const tokenHash = await argon2.hash(token);

    const savedToken = await this.repository.createMagicLink({
      email,
      tokenHash,
      expiresAt: new Date(Date.now() + 1000 * 60 * 10), // 10 minutes
    });

    const magicLinkTokenFormated = savedToken?.public_id + "." + token;
    const magicLinkUrl = new URL(env.WEB_CLIENT_APP_URL);
    magicLinkUrl.pathname = "/auth/magic-link/";
    magicLinkUrl.searchParams.set("token", magicLinkTokenFormated);

    return magicLinkUrl.toString();
  }

  async sendMagicLink({
    email,
    magicLink,
  }: {
    email: string;
    magicLink: string;
  }) {
    const html = (
      await Bun.file("src/assets/html/magic-link.html").text()
    ).replaceAll("%{magicLink}%", magicLink);

    await this.emailProvider.sendEmail({
      to: email,
      subject: "Seu link para entrar no Versum",
      html,
    });
  }

  async createSessionWithMagicLink({
    bruteMagicLinkToken,
    userAgent,
    ip,
  }: {
    bruteMagicLinkToken: string;
    userAgent: string;
    ip: string;
  }) {
    const [magicLinkPublicId, magicLinkToken] = bruteMagicLinkToken.split(".");

    const magicLink = await this.repository.getMagicLinkByPublicId({
      publicId: magicLinkPublicId,
    });

    if (!magicLink) throw new UnauthorizedError("Invalid magic link token!");

    if (magicLink.expires_at.getTime() < new Date().getTime())
      throw new UnauthorizedError("Expired magic link token!");
    if (
      magicLink.invalidated_at &&
      magicLink.invalidated_at.getTime() < new Date().getTime()
    )
      throw new UnauthorizedError("Invalidated magic link token!");
    if (magicLink.used_at && magicLink.used_at.getTime() < new Date().getTime())
      throw new UnauthorizedError("Used magic link token!");

    const tokenIsValid = await argon2.verify(
      magicLink?.token_hash,
      magicLinkToken,
    );

    if (!tokenIsValid) throw new UnauthorizedError("Magic links don't match!");

    await this.repository.updateMagicLink({
      publicId: magicLinkPublicId,
      invalidated_at: new Date(),
    });

    let user = await this.userRepository.findByEmail({
      email: magicLink.email,
    });

    if (!user) {
      user = await this.userRepository.create({
        email: magicLink.email,
        name: magicLink.email.split("@")[0],
        username: magicLink.email.split("@")[0].toLowerCase(),
      });
    }

    const sessionToken = this.generateRandomToken();

    const sessionTokenHash = await argon2.hash(sessionToken);

    const session = await this.repository.createSession({
      userId: user.id,
      userAgent,
      ip,
      tokenHash: sessionTokenHash,
    });

    return session.public_id + "." + sessionToken;
  }

  async refreshSession({
    sessionId,
  }: {
    sessionId: string;
  }): Promise<
    | { session: Session; rotated: false }
    | { session: Session; rotated: true; token: string }
  > {
    const session = await this.repository.getSessionById({
      id: sessionId,
    });

    const validatedSession = new ValidateSession({
      session,
      forwardedFor: session?.ip,
    }).session;

    const refreshWindowMs = 1000 * 60 * 60 * 24 * 10; // 10 days before expiration
    const timeUntilExpirationMs =
      validatedSession.expires_at.getTime() - Date.now();

    if (timeUntilExpirationMs > refreshWindowMs) {
      return {
        session: validatedSession,
        rotated: false,
      };
    }

    const sessionToken = this.generateRandomToken();

    await this.repository.rotateSession({
      id: validatedSession.id,
      publicId: validatedSession.public_id,
      expires_at: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000),
      token_hash: await argon2.hash(sessionToken),
    });

    const updatedSession = await this.repository.getSessionById({
      id: validatedSession.id,
    });

    if (!updatedSession) throw new NotFoundError("Session not found");

    return {
      token: updatedSession.public_id + "." + sessionToken,
      session: updatedSession,
      rotated: true,
    };
  }

  async revokeSession({ sessionPublicId }: { sessionPublicId: string }) {
    await this.repository.revokeSessionByPublicId({
      publicId: sessionPublicId,
    });
  }

  private generateRandomToken() {
    const tokenBytes = crypto.getRandomValues(new Uint8Array(32));
    const token: string = Array.from(tokenBytes)
      .map((b) => b.toString(32).padStart(2, "0"))
      .join("")
      .replaceAll(".", "");

    return token;
  }
}

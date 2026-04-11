import { and, eq, isNull } from "drizzle-orm";
import { db as drizzle } from "../../../infrastructure/db/index.ts";
import { magicLinks } from "../db/magic-links.table.ts";
import { sessions } from "../db/sessions.table.ts";
import type {
  CreateMagicLinkParams,
  CreateSessionParams,
  iAuthRepository,
  MagicLink,
  Session,
} from "./auth.types.repository.ts";

export class AuthRepository implements iAuthRepository {
  private readonly db: typeof drizzle;

  constructor({ db }: { db?: typeof drizzle } = {}) {
    this.db = db ?? drizzle;
  }
  async createMagicLink({
    email,
    expiresAt,
    tokenHash,
  }: CreateMagicLinkParams): Promise<MagicLink | undefined> {
    const [magicLink] = await this.db
      .insert(magicLinks)
      .values({
        email,
        expiresAt,
        tokenHash,
      })
      .returning();
    return magicLink;
  }

  async createSession({
    ip,
    userAgent,
    userId,
    tokenHash,
  }: CreateSessionParams): Promise<Session | undefined> {
    const [session] = await this.db
      .insert(sessions)
      .values({
        userId: userId,
        ip,
        userAgent,
        tokenHash,
        expiresAt: new Date(Date.now() + 1000 * 60 * 60 * 24 * 30), // One month
      })
      .returning();
    return session;
  }

  async getSessionByToken({
    tokenHash,
  }: {
    tokenHash: string;
  }): Promise<Session | null> {
    const session = await this.db.query.sessions.findFirst({
      where: (sessions, { eq }) => eq(sessions.tokenHash, tokenHash),
    });
    return session ?? null;
  }

  async getSessionById({ id }: { id: string }): Promise<Session | null> {
    const session = await this.db.query.sessions.findFirst({
      where: (sessions, { eq }) => eq(sessions.id, id),
    });
    return session ?? null;
  }

  async getSessionByPublicId({
    publicId,
  }: {
    publicId: string;
  }): Promise<Session | null> {
    const session = await this.db.query.sessions.findFirst({
      where: (sessions, { eq }) => eq(sessions.publicId, publicId),
    });
    return session ?? null;
  }

  async getSessionsByUserId({
    userId,
  }: {
    userId: string;
  }): Promise<Session[] | null> {
    const userSessions = await this.db.query.sessions.findMany({
      where: (sessions, { eq }) => eq(sessions.userId, userId),
    });
    return userSessions.length > 0 ? userSessions : null;
  }

  async getMagicLinkByPublicId({
    publicId,
  }: {
    publicId: string;
  }): Promise<MagicLink | null> {
    const magicLink = await this.db.query.magicLinks.findFirst({
      where: (magicLinks, { eq }) => eq(magicLinks.publicId, publicId),
    });
    return magicLink ?? null;
  }

  async getMagicLinkByEmail({
    email,
  }: {
    email: string;
  }): Promise<MagicLink | null> {
    const magicLink = await this.db.query.magicLinks.findFirst({
      where: (magicLinks, { eq }) => eq(magicLinks.email, email),
    });
    return magicLink ?? null;
  }

  async updateMagicLink({
    publicId,
    invalidatedAt,
    usedAt,
  }: {
    publicId: string;
    invalidatedAt?: Date;
    usedAt?: Date;
  }): Promise<MagicLink | null> {
    const magicLink = await this.db.query.magicLinks.findFirst({
      where: (magicLinks, { eq }) =>
        and(
          eq(magicLinks.publicId, publicId),
          isNull(magicLinks.usedAt),
          isNull(magicLinks.invalidatedAt),
        ),
    });
    if (!magicLink) return null;

    await this.db
      .update(magicLinks)
      .set({
        invalidatedAt,
        usedAt,
      })
      .where(
        and(
          eq(magicLinks.publicId, magicLink.publicId),
          isNull(magicLinks.usedAt),
          isNull(magicLinks.invalidatedAt),
        ),
      );

    const updatedMagicLink = await this.db.query.magicLinks.findFirst({
      where: (magicLinks, { eq }) => eq(magicLinks.publicId, publicId),
    });

    return updatedMagicLink ?? null;
  }

  async rotateSession({
    id,
    publicId,
    expiresAt,
    tokenHash,
  }: {
    id: string;
    publicId: string;
    expiresAt: Date;
    tokenHash: string;
  }) {
    await this.db
      .update(sessions)
      .set({
        tokenHash,
        expiresAt,
      })
      .where(
        and(
          eq(sessions.publicId, publicId),
          eq(sessions.id, id),
          isNull(sessions.revokedAt),
        ),
      );
  }

  async revokeSessionByPublicId({
    publicId,
  }: {
    publicId: string;
  }): Promise<void> {
    await this.db
      .update(sessions)
      .set({
        revokedAt: new Date(),
      })
      .where(and(eq(sessions.publicId, publicId), isNull(sessions.revokedAt)));
  }
}

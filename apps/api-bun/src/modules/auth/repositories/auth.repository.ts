import { and, eq, isNull } from "drizzle-orm";
import { db as drizzle } from "../../../infrastructure/db/index.ts";
import { magicLinks } from "../db/magic-links.table.ts";
import { sessions } from "../db/sessions.table.ts";
import {
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
    expiresAt: expires_at,
    tokenHash: token_hash,
  }: CreateMagicLinkParams): Promise<MagicLink> {
    const [magicLink] = await this.db
      .insert(magicLinks)
      .values({
        email,
        expires_at,
        token_hash,
      })
      .returning();
    return magicLink;
  }

  async createSession({
    ip,
    userAgent,
    userId,
    tokenHash: token_hash,
  }: CreateSessionParams): Promise<Session> {
    const [session] = await this.db
      .insert(sessions)
      .values({
        user_id: userId,
        ip,
        user_agent: userAgent,
        token_hash,
        expires_at: new Date(Date.now() + 1000 * 60 * 60 * 24 * 30), // One month
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
      where: (sessions, { eq }) => eq(sessions.token_hash, tokenHash),
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
      where: (sessions, { eq }) => eq(sessions.public_id, publicId),
    });
    return session ?? null;
  }

  async getSessionsByUserId({
    userId,
  }: {
    userId: string;
  }): Promise<Session[] | null> {
    const userSessions = await this.db.query.sessions.findMany({
      where: (sessions, { eq }) => eq(sessions.user_id, userId),
    });
    return userSessions.length > 0 ? userSessions : null;
  }

  async getMagicLinkByPublicId({
    publicId,
  }: {
    publicId: string;
  }): Promise<MagicLink | null> {
    const magicLink = await this.db.query.magicLinks.findFirst({
      where: (magicLinks, { eq }) => eq(magicLinks.public_id, publicId),
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
    invalidated_at,
    used_at,
  }: {
    publicId: string;
    invalidated_at?: Date;
    used_at?: Date;
  }): Promise<MagicLink | null> {
    const magicLink = await this.db.query.magicLinks.findFirst({
      where: (magicLinks, { eq }) =>
        and(
          eq(magicLinks.public_id, publicId),
          isNull(magicLinks.used_at),
          isNull(magicLinks.invalidated_at),
        ),
    });
    if (!magicLink) return null;

    await this.db
      .update(magicLinks)
      .set({
        invalidated_at,
        used_at,
      })
      .where(
        and(
          eq(magicLinks.public_id, magicLink.public_id),
          isNull(magicLinks.used_at),
          isNull(magicLinks.invalidated_at),
        ),
      );

    const updatedMagicLink = await this.db.query.magicLinks.findFirst({
      where: (magicLinks, { eq }) => eq(magicLinks.public_id, publicId),
    });

    return updatedMagicLink ?? null;
  }

  async rotateSession({
    id,
    publicId,
    expires_at,
    token_hash,
  }: {
    id: string;
    publicId: string;
    expires_at: Date;
    token_hash: string;
  }) {
    await this.db
      .update(sessions)
      .set({
        token_hash,
        expires_at,
      })
      .where(
        and(
          eq(sessions.public_id, publicId),
          eq(sessions.id, id),
          isNull(sessions.revoked_at),
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
        revoked_at: new Date(),
      })
      .where(
        and(eq(sessions.public_id, publicId), isNull(sessions.revoked_at)),
      );
  }
}

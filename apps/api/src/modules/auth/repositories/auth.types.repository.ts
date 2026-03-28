import type { InferSelectModel } from "drizzle-orm";
import type { magicLinks, sessions } from "../../../infrastructure/db/schema.ts";

export type MagicLink = InferSelectModel<typeof magicLinks>;
export type Session = InferSelectModel<typeof sessions>;

export type CreateMagicLinkParams = {
  email: string;
  expiresAt: Date;
  tokenHash: string;
};

export type CreateSessionParams = {
  userId: string;
  ip: string;
  userAgent: string;
  tokenHash: string;
};

export interface iAuthRepository {
  createMagicLink(params: CreateMagicLinkParams): Promise<MagicLink | undefined>;
  createSession(params: CreateSessionParams): Promise<Session | undefined>;

  getSessionByToken(params: { tokenHash: string }): Promise<Session | null>;
  getSessionById(params: { id: string }): Promise<Session | null>;
  getSessionsByUserId(params: { userId: string }): Promise<Session[] | null>;

  getMagicLinkByPublicId(params: {
    publicId: string;
  }): Promise<MagicLink | null>;
  getMagicLinkByEmail(params: { email: string }): Promise<MagicLink | null>;

  rotateSession(params: {
    id: string;
    publicId: string;
    expiresAt: Date;
    tokenHash: string;
  }): Promise<void>;

  revokeSessionByPublicId(params: { publicId: string }): Promise<void>;
}

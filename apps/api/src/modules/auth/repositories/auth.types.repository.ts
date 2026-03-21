import { InferSelectModel } from "drizzle-orm";
import { magicLinks, sessions } from "../../../infrastructure/db/schema.ts";

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
  createMagicLink(params: CreateMagicLinkParams): Promise<MagicLink>;
  createSession(params: CreateSessionParams): Promise<Session>;

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
    expires_at: Date;
    token_hash: string;
  }): Promise<void>;

  revokeSessionByPublicId(params: { publicId: string }): Promise<void>;
}

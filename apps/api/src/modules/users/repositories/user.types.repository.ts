import type { InferInsertModel, InferSelectModel } from "drizzle-orm";
import type { users } from "../../../infrastructure/db/schema";

export type User = InferSelectModel<typeof users>;
export type CreateUserParams = InferInsertModel<typeof users>;
export type UpdateUserParams = Partial<Omit<User, "id" | "createdAt">>;

import type { Profile } from "./profile.types.repository";

export interface UserExportData {
  user: User;
  profile: Profile | null;
  readings: unknown[];
  discoveryReadings: unknown[];
  likes: unknown[];
  marks: unknown[];
  sessions: unknown[];
  consentLogs: unknown[];
}

export interface iUserRepository {
  create(params: CreateUserParams): Promise<User>;
  findById(params: { id: string }): Promise<User | null>;
  findByIdWithProfile(params: { id: string }): Promise<(User & { profile: Profile | undefined }) | null>;
  findByIdWithAllData(params: { id: string }): Promise<UserExportData | null>;
  findByEmail(params: { email: string }): Promise<User | null>;
  updateUser(params: UpdateUserParams & { id: string }): Promise<User>;
  deleteUser(params: { id: string }): Promise<void>;
}

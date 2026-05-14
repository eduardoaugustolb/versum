import type { InferInsertModel, InferSelectModel } from "drizzle-orm";
import type { users } from "../../../infrastructure/db/schema";

export type User = InferSelectModel<typeof users>;
export type CreateUserParams = InferInsertModel<typeof users>;
export type UpdateUserParams = Partial<Omit<User, "id" | "createdAt">>;

import type { Profile } from "./profile.types.repository";

export interface iUserRepository {
  create(params: CreateUserParams): Promise<User>;
  findById(params: { id: string }): Promise<User | null>;
  findByIdWithProfile(params: { id: string }): Promise<(User & { profile: Profile | undefined }) | null>;
  findByEmail(params: { email: string }): Promise<User | null>;
  updateUser(params: UpdateUserParams & { id: string }): Promise<User>;
}

import type { InferInsertModel, InferSelectModel } from "drizzle-orm";
import type { users } from "../../../infrastructure/db/schema";

export type User = InferSelectModel<typeof users>;
export type CreateUserParams = InferInsertModel<typeof users>;
export type UpdateUserParams = Partial<Omit<User, "id" | "createdAt">>;

export interface iUserRepository {
  create(params: CreateUserParams): Promise<User>;
  findById(params: { id: string }): Promise<User | null>;
  findByEmail(params: { email: string }): Promise<User | null>;
  updateUser(params: UpdateUserParams & { id: string }): Promise<User>;
}

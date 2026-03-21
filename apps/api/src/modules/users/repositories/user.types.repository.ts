import { InferInsertModel, InferSelectModel } from "drizzle-orm";
import { users } from "../../../infrastructure/db/schema.ts";

export type User = InferSelectModel<typeof users>;
export type CreateUserParams = InferInsertModel<typeof users>;

export interface iUserRepository {
  create(params: CreateUserParams): Promise<User>;
  findById(params: { id: string }): Promise<User | null>;
  findByEmail(params: { email: string }): Promise<User | null>;
  findByUsername(params: { username: string }): Promise<User | null>;
}

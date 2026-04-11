import type { InferInsertModel, InferSelectModel } from "drizzle-orm";
import type { profiles } from "../../../infrastructure/db/schema";

export type Profile = InferSelectModel<typeof profiles>;
export type CreateProfileParams = InferInsertModel<typeof profiles>;
export type UpdateProfileParams = Partial<
  Omit<Profile, "id" | "userId" | "createdAt" | "updatedAt">
>;

export interface iProfileRepository {
  create(params: CreateProfileParams): Promise<Profile>;
  findById(params: { id: string }): Promise<Profile | null>;
  findByUserId(params: { userId: string }): Promise<Profile | null>;
  findByUsername(params: { username: string }): Promise<Profile | null>;
  update(params: UpdateProfileParams & { id: string }): Promise<Profile>;
  existsByUsername(params: { username: string }): Promise<boolean>;
}

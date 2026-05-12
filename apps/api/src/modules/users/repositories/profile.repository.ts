import { eq } from "drizzle-orm";
import { db as drizzle } from "../../../infrastructure/db/index";
import { profiles } from "../../../infrastructure/db/schema";
import type {
  CreateProfileParams,
  iProfileRepository,
  Profile,
  UpdateProfileParams,
} from "./profile.types.repository";
import { InternalServerError } from "@/utils/app/errors";

export class ProfileRepository implements iProfileRepository {
  private readonly db: typeof drizzle;

  constructor({ db }: { db?: typeof drizzle } = {}) {
    this.db = db ?? drizzle;
  }

  async create(params: CreateProfileParams): Promise<Profile> {
    const [profile] = await this.db.insert(profiles).values(params).returning();

    if (!profile) {
      throw new InternalServerError("Create failed");
    }

    return profile;
  }

  async findById({ id }: { id: string }): Promise<Profile | null> {
    const profile = await this.db.query.profiles.findFirst({
      where: (profiles, { eq }) => eq(profiles.id, id),
    });
    return profile ?? null;
  }

  async findByUserId({ userId }: { userId: string }): Promise<Profile | null> {
    const profile = await this.db.query.profiles.findFirst({
      where: (profiles, { eq }) => eq(profiles.userId, userId),
    });
    return profile ?? null;
  }

  async findByUsername({
    username,
  }: {
    username: string;
  }): Promise<Profile | null> {
    const profile = await this.db.query.profiles.findFirst({
      where: (profiles, { eq }) =>
        eq(profiles.username, username.toLowerCase()),
    });
    return profile ?? null;
  }

  async update(params: UpdateProfileParams & { id: string }): Promise<Profile> {
    const existing = await this.findById({ id: params.id });
    if (!existing) {
      throw new Error("Profile not found");
    }

    const [updated] = await this.db
      .update(profiles)
      .set({
        ...params,
        updatedAt: new Date(),
      })
      .where(eq(profiles.id, params.id))
      .returning();

    if (!updated) {
      throw new InternalServerError("Update failed");
    }

    return updated;
  }

  async existsByUsername({ username }: { username: string }): Promise<{ exists: boolean, profileId: string } | { exists: false }> {
    const profile = await this.db.query.profiles.findFirst({
      where: (profiles, { eq }) =>
        eq(profiles.username, username.toLowerCase()),
      columns: { id: true },
    });
    return profile ? { exists: true, profileId: profile.id } : { exists: false };
  }
}

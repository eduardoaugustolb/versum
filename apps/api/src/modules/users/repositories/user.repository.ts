import { eq } from "drizzle-orm";
import { db as drizzle } from "../../../infrastructure/db/index";
import { users } from "../../../infrastructure/db/schema";
import type {
  CreateUserParams,
  iUserRepository,
  UpdateUserParams,
  User,
} from "./user.types.repository";
import { InternalServerError } from "@/utils/app/errors";
import type { Profile } from "./profile.types.repository";

export class UserRepository implements iUserRepository {
  private readonly db: typeof drizzle;

  constructor({ db }: { db?: typeof drizzle } = {}) {
    this.db = db ?? drizzle;
  }

  async create(params: CreateUserParams): Promise<User> {
    const [user] = await this.db.insert(users).values(params).returning();
    if (!user) {
      throw new InternalServerError("Failed to create user");
    }
    return user;
  }

  async findById({ id }: { id: string }): Promise<User | null> {
    const user = await this.db.query.users.findFirst({
      where: (users, { eq }) => eq(users.id, id),
    });

    return user ?? null;
  }

  async findByIdWithProfile({ id }: { id: string }): Promise<(User & { profile: Profile | undefined }) | null> {
    const user = await this.db.query.users.findFirst({
      where: (users, { eq }) => eq(users.id, id),
      with: {
        profile: true
      }
    });

    return user ?? null;
  }

  async findByEmail({ email }: { email: string }): Promise<User | null> {
    const user = await this.db.query.users.findFirst({
      where: (users, { eq }) => eq(users.email, email.toLowerCase()),
    });

    return user ?? null;
  }

  async updateUser(params: UpdateUserParams & { id: string }): Promise<User> {
    const [updated] = await this.db
      .update(users)
      .set(params)
      .where(eq(users.id, params.id))
      .returning();

    if (!updated) {
      throw new InternalServerError("Failed to update user");
    }

    return updated;
  }
}

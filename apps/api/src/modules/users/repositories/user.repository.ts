import { eq } from "drizzle-orm";
import { db as drizzle } from "../../../infrastructure/db/index.ts";
import { users } from "../../../infrastructure/db/schema.ts";
import type {
  CreateUserParams,
  iUserRepository,
  User,
} from "./user.types.repository.ts";

export class UserRepository implements iUserRepository {
  private readonly db: typeof drizzle;

  constructor({ db }: { db?: typeof drizzle } = {}) {
    this.db = db ?? drizzle;
  }

  async create(params: CreateUserParams): Promise<User> {
    const [user] = await this.db.insert(users).values(params).returning();
    return user;
  }

  async findById({ id }: { id: string }): Promise<User | null> {
    const user = await this.db.query.users.findFirst({
      where: (users, { eq }) => eq(users.id, id),
    });

    return user ?? null;
  }

  async findByEmail({ email }: { email: string }): Promise<User | null> {
    const user = await this.db.query.users.findFirst({
      where: (users, { eq }) => eq(users.email, email),
    });

    return user ?? null;
  }

  async findByUsername({
    username,
  }: {
    username: string;
  }): Promise<User | null> {
    const user = await this.db.query.users.findFirst({
      where: (users, { eq }) => eq(users.username, username),
    });

    return user ?? null;
  }

  async updateUser(params: Partial<User> & { id: string }): Promise<User> {
    const user = await this.db.query.users.findFirst({
      where: (users, { eq }) => eq(users.id, params.id),
    });

    if (!user) {
      throw new Error("User not found");
    }

    const updatedUser = await this.db
      .update(users)
      .set({
        ...params,
      })
      .where(eq(users.id, user.id))
      .returning();

    return updatedUser[0];
  }
}

import { UserRepository } from "../repositories/user.repository.ts";
import type {
  CreateUserParams,
  User,
} from "../repositories/user.types.repository.ts";

export class UserServiceV1 {
  private readonly repository: UserRepository;

  constructor({ repository }: { repository?: UserRepository } = {}) {
    this.repository = repository ?? new UserRepository();
  }

  async createUser(params: CreateUserParams): Promise<User> {
    return await this.repository.create(params);
  }

  async getUserById({ id }: { id: string }): Promise<User> {
    const user = await this.repository.findById({ id });

    if (!user) {
      throw new Error("User not found");
    }

    return user;
  }

  async updateUser(params: Partial<User> & { id: string }): Promise<User> {
    const user = await this.repository.findById({ id: params.id });

    if (!user) {
      throw new Error("User not found");
    }

    const updatedUser = await this.repository.updateUser({
      ...params,
      id: user.id,
    });

    return updatedUser;
  }

  async getUserByEmail({ email }: { email: string }): Promise<User | null> {
    return await this.repository.findByEmail({ email });
  }

  async getUserByUsername({
    username,
  }: {
    username: string;
  }): Promise<User | null> {
    return await this.repository.findByUsername({ username });
  }
}

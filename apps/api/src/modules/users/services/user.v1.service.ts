import { UserRepository } from "../repositories/user.repository";
import type {
  CreateUserParams,
  User,
} from "../repositories/user.types.repository";

const MAX_EMAIL_LENGTH = 255;

export class UserServiceV1 {
  private readonly repository: UserRepository;

  constructor({ repository }: { repository?: UserRepository } = {}) {
    this.repository = repository ?? new UserRepository();
  }

  private validateEmail(email: string): void {
    const trimmed = email.trim().toLowerCase();

    if (trimmed.length === 0) {
      throw new Error("Email is required");
    }

    if (trimmed.length > MAX_EMAIL_LENGTH) {
      throw new Error(`Email must not exceed ${MAX_EMAIL_LENGTH} characters`);
    }

    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(trimmed)) {
      throw new Error("Invalid email format");
    }
  }

  async createUser(params: CreateUserParams): Promise<User> {
    const normalizedEmail = params.email.trim().toLowerCase();
    this.validateEmail(normalizedEmail);

    return await this.repository.create({
      ...params,
      email: normalizedEmail,
    });
  }

  async getUserById({ id }: { id: string }): Promise<User> {
    const user = await this.repository.findById({ id });

    if (!user) {
      throw new Error("User not found");
    }

    return user;
  }

  async getUserByEmail({ email }: { email: string }): Promise<User | null> {
    return await this.repository.findByEmail({ email });
  }

  async updateUserEmail(params: { id: string; email: string }): Promise<User> {
    const user = await this.repository.findById({ id: params.id });

    if (!user) {
      throw new Error("User not found");
    }

    this.validateEmail(params.email);

    const normalizedEmail = params.email.trim().toLowerCase();

    const existingWithEmail = await this.repository.findByEmail({
      email: normalizedEmail,
    });

    if (existingWithEmail && existingWithEmail.id !== user.id) {
      throw new Error("Email already in use");
    }

    return await this.repository.updateUser({
      id: user.id,
      email: normalizedEmail,
    });
  }
}

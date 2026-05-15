import { AuthRepository } from "../../auth/repositories/auth.repository";
import type { Profile } from "../repositories/profile.types.repository";
import { ProfileRepository } from "../repositories/profile.repository";
import { UserRepository } from "../repositories/user.repository";
import type {
  CreateUserParams,
  User,
  UserExportData,
} from "../repositories/user.types.repository";
import type { ExportUserDataResponse } from "../schemas/v1/users.v1.common.schema";

const MAX_EMAIL_LENGTH = 255;

export class UserServiceV1 {
  private readonly repository: UserRepository;
  private readonly authRepository: AuthRepository;
  private readonly profileRepository: ProfileRepository;

  constructor({
    repository,
    authRepository,
    profileRepository,
  }: {
    repository?: UserRepository;
    authRepository?: AuthRepository;
    profileRepository?: ProfileRepository;
  } = {}) {
    this.repository = repository ?? new UserRepository();
    this.authRepository = authRepository ?? new AuthRepository();
    this.profileRepository = profileRepository ?? new ProfileRepository();
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

  async getUserByIdWithProfile({ id }: { id: string }): Promise<User & { profile: Profile | undefined }> {
    const user = await this.repository.findByIdWithProfile({ id });

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

  async deleteUser({ id }: { id: string }): Promise<void> {
    const user = await this.repository.findById({ id });

    if (!user) {
      throw new Error("User not found");
    }

    await this.authRepository.deleteSessionsByUserId({ userId: id });
    await this.profileRepository.deleteByUserId({ userId: id });
    await this.repository.deleteUser({ id });
  }

  async exportUserData({ id }: { id: string }): Promise<ExportUserDataResponse> {
    const data = await this.repository.findByIdWithAllData({ id });

    if (!data) {
      throw new Error("User not found");
    }

    const sessions = (data.sessions ?? []).map((s: Record<string, unknown>) => ({
      createdAt: (s.createdAt as Date).toISOString(),
      ip: s.ip as string,
      userAgent: s.userAgent as string,
      expiresAt: (s.expiresAt as Date).toISOString(),
    }));

    return {
      exportedAt: new Date().toISOString(),
      user: {
        email: data.user.email,
        createdAt: data.user.createdAt.toISOString(),
      },
      profile: data.profile
        ? {
            username: data.profile.username,
            name: data.profile.name,
            bio: data.profile.bio ?? null,
            pictureUrl: data.profile.pictureUrl ?? null,
          }
        : null,
      sessions,
      readingHistory: {
        journey: (data.readings ?? []).map((r: Record<string, unknown>) => ({
          chapterId: r.chapterId as string,
          readAt: (r.readAt as Date).toISOString(),
        })),
        discovery: (data.discoveryReadings ?? []).map((r: Record<string, unknown>) => ({
          verseId: r.verseId as string,
          readAt: (r.readAt as Date).toISOString(),
        })),
      },
      annotations: (data.marks ?? []).map((m: Record<string, unknown>) => ({
        verseId: m.verseId as string,
        selectedVerseId: m.selectedVerseId as string,
        annotation: (m.annotation as string | null) ?? null,
        isPublic: m.isPublic as boolean,
        createdAt: (m.createdAt as Date).toISOString(),
      })),
      likes: (data.likes ?? []).map((l: Record<string, unknown>) => ({
        verseId: l.verseId as string,
        createdAt: (l.createdAt as Date).toISOString(),
      })),
      consentLogs: (data.consentLogs ?? []).map((c: Record<string, unknown>) => ({
        id: c.id as string,
        userId: c.userId as string,
        purpose: c.purpose as string,
        granted: c.granted as boolean,
        ip: c.ip as string,
        userAgent: c.userAgent as string,
        createdAt: (c.createdAt as Date).toISOString(),
      })),
    };
  }
}

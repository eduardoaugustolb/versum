import { beforeEach, describe, expect, it, vi } from "vitest";
import type { User } from "../repositories/user.types.repository";
import { UserServiceV1 } from "./user.v1.service";

describe("UserServiceV1", () => {
  let service: UserServiceV1;

  const mockUser: User = {
    id: "123e4567-e89b-12d3-a456-426614174000",
    email: "john@example.com",
    createdAt: new Date("2024-01-01T00:00:00Z"),
  };

  const createMockRepository = () => ({
    db: {},
    create: vi.fn<() => Promise<User>>(),
    findById: vi.fn<() => Promise<User | null>>(),
    findByIdWithProfile: vi.fn<() => Promise<User | null>>(),
    findByEmail: vi.fn<() => Promise<User | null>>(),
    deleteUser: vi.fn<() => Promise<void>>(),
  });

  const createMockAuthRepository = () => ({
    db: {},
    deleteSessionsByUserId: vi.fn<() => Promise<void>>(),
  });

  const createMockProfileRepository = () => ({
    db: {},
    deleteByUserId: vi.fn<() => Promise<void>>(),
  });

  const createService = ({
    mockRepository,
    mockAuthRepository,
    mockProfileRepository,
  }: {
    mockRepository: ReturnType<typeof createMockRepository>;
    mockAuthRepository: ReturnType<typeof createMockAuthRepository>;
    mockProfileRepository: ReturnType<typeof createMockProfileRepository>;
  }) =>
    new UserServiceV1({
      repository:
        mockRepository as unknown as import("../repositories/user.repository").UserRepository,
      authRepository:
        mockAuthRepository as unknown as import("../../auth/repositories/auth.repository").AuthRepository,
      profileRepository:
        mockProfileRepository as unknown as import("../repositories/profile.repository").ProfileRepository,
    });

  beforeEach(() => {
    const mockRepository = createMockRepository();
    const mockAuthRepository = createMockAuthRepository();
    const mockProfileRepository = createMockProfileRepository();
    service = createService({ mockRepository, mockAuthRepository, mockProfileRepository });
    vi.clearAllMocks();
  });

  describe("createUser", () => {
    it("should create a new user with normalized email", async () => {
      const mockRepository = createMockRepository();
      const mockAuthRepository = createMockAuthRepository();
      const mockProfileRepository = createMockProfileRepository();
      mockRepository.create.mockResolvedValue(mockUser);
      service = createService({ mockRepository, mockAuthRepository, mockProfileRepository });

      const createParams = {
        email: "John@Example.COM",
      };

      const result = await service.createUser(createParams);

      expect(result).toEqual(mockUser);
      expect(mockRepository.create).toHaveBeenCalledWith({
        email: "john@example.com",
      });
    });

    it("should throw error for invalid email format", async () => {
      const mockRepository = createMockRepository();
      const mockAuthRepository = createMockAuthRepository();
      const mockProfileRepository = createMockProfileRepository();
      service = createService({ mockRepository, mockAuthRepository, mockProfileRepository });

      await expect(
        service.createUser({ email: "invalid-email" }),
      ).rejects.toThrow("Invalid email format");
    });

    it("should throw error for empty email", async () => {
      const mockRepository = createMockRepository();
      const mockAuthRepository = createMockAuthRepository();
      const mockProfileRepository = createMockProfileRepository();
      service = createService({ mockRepository, mockAuthRepository, mockProfileRepository });

      await expect(service.createUser({ email: "   " })).rejects.toThrow(
        "Email is required",
      );
    });

    it("should throw error for email exceeding max length", async () => {
      const mockRepository = createMockRepository();
      const mockAuthRepository = createMockAuthRepository();
      const mockProfileRepository = createMockProfileRepository();
      service = createService({ mockRepository, mockAuthRepository, mockProfileRepository });

      const longEmail = `${"a".repeat(250)}@example.com`;
      await expect(service.createUser({ email: longEmail })).rejects.toThrow(
        "Email must not exceed 255 characters",
      );
    });
  });

  describe("getUserById", () => {
    it("should return user when found", async () => {
      const mockRepository = createMockRepository();
      const mockAuthRepository = createMockAuthRepository();
      const mockProfileRepository = createMockProfileRepository();
      mockRepository.findById.mockResolvedValue(mockUser);
      service = createService({ mockRepository, mockAuthRepository, mockProfileRepository });

      const result = await service.getUserById({ id: mockUser.id });

      expect(result).toEqual(mockUser);
      expect(mockRepository.findById).toHaveBeenCalledWith({ id: mockUser.id });
    });

    it("should throw error when user not found", async () => {
      const mockRepository = createMockRepository();
      const mockAuthRepository = createMockAuthRepository();
      const mockProfileRepository = createMockProfileRepository();
      mockRepository.findById.mockResolvedValue(null);
      service = createService({ mockRepository, mockAuthRepository, mockProfileRepository });

      await expect(
        service.getUserById({ id: "nonexistent-id" }),
      ).rejects.toThrow("User not found");
    });
  });

  describe("getUserByEmail", () => {
    it("should return user when found", async () => {
      const mockRepository = createMockRepository();
      const mockAuthRepository = createMockAuthRepository();
      const mockProfileRepository = createMockProfileRepository();
      mockRepository.findByEmail.mockResolvedValue(mockUser);
      service = createService({ mockRepository, mockAuthRepository, mockProfileRepository });

      const result = await service.getUserByEmail({ email: mockUser.email });

      expect(result).toEqual(mockUser);
      expect(mockRepository.findByEmail).toHaveBeenCalledWith({
        email: mockUser.email,
      });
    });

    it("should return null when user not found", async () => {
      const mockRepository = createMockRepository();
      const mockAuthRepository = createMockAuthRepository();
      const mockProfileRepository = createMockProfileRepository();
      mockRepository.findByEmail.mockResolvedValue(null);
      service = createService({ mockRepository, mockAuthRepository, mockProfileRepository });

      const result = await service.getUserByEmail({
        email: "nonexistent@example.com",
      });

      expect(result).toBeNull();
    });
  });

  describe("deleteUser", () => {
    it("should delete user and related data", async () => {
      const mockRepository = createMockRepository();
      const mockAuthRepository = createMockAuthRepository();
      const mockProfileRepository = createMockProfileRepository();
      mockRepository.findById.mockResolvedValue(mockUser);
      mockAuthRepository.deleteSessionsByUserId.mockResolvedValue(undefined);
      mockProfileRepository.deleteByUserId.mockResolvedValue(undefined);
      mockRepository.deleteUser.mockResolvedValue(undefined);
      service = createService({ mockRepository, mockAuthRepository, mockProfileRepository });

      await service.deleteUser({ id: mockUser.id });

      expect(mockRepository.findById).toHaveBeenCalledWith({ id: mockUser.id });
      expect(mockAuthRepository.deleteSessionsByUserId).toHaveBeenCalledWith({ userId: mockUser.id });
      expect(mockProfileRepository.deleteByUserId).toHaveBeenCalledWith({ userId: mockUser.id });
      expect(mockRepository.deleteUser).toHaveBeenCalledWith({ id: mockUser.id });
    });

    it("should delete sessions before profile before user", async () => {
      const mockRepository = createMockRepository();
      const mockAuthRepository = createMockAuthRepository();
      const mockProfileRepository = createMockProfileRepository();
      mockRepository.findById.mockResolvedValue(mockUser);
      service = createService({ mockRepository, mockAuthRepository, mockProfileRepository });

      await service.deleteUser({ id: mockUser.id });

      const sessionOrder = mockAuthRepository.deleteSessionsByUserId.mock.invocationCallOrder[0];
      const profileOrder = mockProfileRepository.deleteByUserId.mock.invocationCallOrder[0];
      const userOrder = mockRepository.deleteUser.mock.invocationCallOrder[0];

      expect(sessionOrder).toBeLessThan(profileOrder!);
      expect(profileOrder!).toBeLessThan(userOrder!);
    });

    it("should throw error when user not found", async () => {
      const mockRepository = createMockRepository();
      const mockAuthRepository = createMockAuthRepository();
      const mockProfileRepository = createMockProfileRepository();
      mockRepository.findById.mockResolvedValue(null);
      service = createService({ mockRepository, mockAuthRepository, mockProfileRepository });

      await expect(
        service.deleteUser({ id: "nonexistent-id" }),
      ).rejects.toThrow("User not found");

      expect(mockAuthRepository.deleteSessionsByUserId).not.toHaveBeenCalled();
      expect(mockProfileRepository.deleteByUserId).not.toHaveBeenCalled();
      expect(mockRepository.deleteUser).not.toHaveBeenCalled();
    });
  });
});

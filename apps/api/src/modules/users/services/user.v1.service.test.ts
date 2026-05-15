import { beforeEach, describe, expect, it, vi } from "vitest";
import type { User, UserExportData } from "../repositories/user.types.repository";
import { UserServiceV1 } from "./user.v1.service";

describe("UserServiceV1", () => {
  let service: UserServiceV1;

  const mockUser: User = {
    id: "123e4567-e89b-12d3-a456-426614174000",
    email: "john@example.com",
    createdAt: new Date("2024-01-01T00:00:00Z"),
  };

  const mockUserExportData: UserExportData = {
    user: mockUser,
    profile: {
      id: "profile-id",
      userId: mockUser.id,
      username: "john",
      name: "John Doe",
      bio: "Bio text",
      pictureUrl: "https://example.com/avatar.jpg",
      createdAt: new Date("2024-01-01T00:00:00Z"),
      updatedAt: new Date("2024-01-01T00:00:00Z"),
    },
    readings: [
      {
        id: "reading-1",
        userId: mockUser.id,
        chapterId: "chapter-1",
        readAt: new Date("2024-06-01T00:00:00Z"),
      },
    ],
    discoveryReadings: [
      {
        id: "dreading-1",
        userId: mockUser.id,
        verseId: "verse-1",
        readAt: new Date("2024-06-15T00:00:00Z"),
      },
    ],
    likes: [
      {
        id: "like-1",
        userId: mockUser.id,
        verseId: "verse-2",
        createdAt: new Date("2024-06-10T00:00:00Z"),
      },
    ],
    marks: [
      {
        id: "mark-1",
        userId: mockUser.id,
        verseId: "verse-3",
        selectedVerseId: "verse-3",
        annotation: "Great verse",
        isPublic: false,
        createdAt: new Date("2024-06-05T00:00:00Z"),
      },
    ],
    sessions: [
      {
        id: "session-1",
        publicId: "pub-session-1",
        userId: mockUser.id,
        ip: "127.0.0.1",
        userAgent: "Mozilla/5.0",
        tokenHash: "hash123",
        createdAt: new Date("2024-01-01T00:00:00Z"),
        updatedAt: new Date("2024-01-01T00:00:00Z"),
        expiresAt: new Date("2025-01-01T00:00:00Z"),
        revokedAt: null,
      },
    ],
    consentLogs: [
      {
        id: "consent-1",
        userId: mockUser.id,
        purpose: "profile_content",
        granted: true,
        ip: "127.0.0.1",
        userAgent: "Mozilla/5.0",
        createdAt: new Date("2024-01-01T00:00:00Z"),
      },
    ],
  };

  const createMockRepository = () => ({
    db: {},
    create: vi.fn<() => Promise<User>>(),
    findById: vi.fn<() => Promise<User | null>>(),
    findByIdWithProfile: vi.fn<() => Promise<User | null>>(),
    findByIdWithAllData: vi.fn<() => Promise<UserExportData | null>>(),
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

  describe("exportUserData", () => {
    it("should return all user data in export format", async () => {
      const mockRepository = createMockRepository();
      const mockAuthRepository = createMockAuthRepository();
      const mockProfileRepository = createMockProfileRepository();
      mockRepository.findByIdWithAllData.mockResolvedValue(mockUserExportData);
      service = createService({ mockRepository, mockAuthRepository, mockProfileRepository });

      const result = await service.exportUserData({ id: mockUser.id });

      expect(result.user.email).toBe("john@example.com");
      expect(result.profile).not.toBeNull();
      expect(result.profile!.username).toBe("john");
      expect(result.profile!.name).toBe("John Doe");
      expect(result.sessions).toHaveLength(1);
      expect(result.sessions[0].ip).toBe("127.0.0.1");
      expect(result.sessions[0]).not.toHaveProperty("tokenHash");
      expect(result.readingHistory.journey).toHaveLength(1);
      expect(result.readingHistory.discovery).toHaveLength(1);
      expect(result.annotations).toHaveLength(1);
      expect(result.annotations[0].annotation).toBe("Great verse");
      expect(result.likes).toHaveLength(1);
      expect(result.consentLogs).toHaveLength(1);
      expect(result.exportedAt).toBeDefined();
      expect(mockRepository.findByIdWithAllData).toHaveBeenCalledWith({ id: mockUser.id });
    });

    it("should throw error when user not found", async () => {
      const mockRepository = createMockRepository();
      const mockAuthRepository = createMockAuthRepository();
      const mockProfileRepository = createMockProfileRepository();
      mockRepository.findByIdWithAllData.mockResolvedValue(null);
      service = createService({ mockRepository, mockAuthRepository, mockProfileRepository });

      await expect(
        service.exportUserData({ id: "nonexistent-id" }),
      ).rejects.toThrow("User not found");
    });

    it("should return partial data when profile/sessions are empty", async () => {
      const mockRepository = createMockRepository();
      const mockAuthRepository = createMockAuthRepository();
      const mockProfileRepository = createMockProfileRepository();
      const partialData: UserExportData = {
        user: mockUser,
        profile: null,
        readings: [],
        discoveryReadings: [],
        likes: [],
        marks: [],
        sessions: [],
        consentLogs: [],
      };
      mockRepository.findByIdWithAllData.mockResolvedValue(partialData);
      service = createService({ mockRepository, mockAuthRepository, mockProfileRepository });

      const result = await service.exportUserData({ id: mockUser.id });

      expect(result.profile).toBeNull();
      expect(result.sessions).toHaveLength(0);
      expect(result.readingHistory.journey).toHaveLength(0);
      expect(result.readingHistory.discovery).toHaveLength(0);
      expect(result.annotations).toHaveLength(0);
      expect(result.likes).toHaveLength(0);
      expect(result.consentLogs).toHaveLength(0);
      expect(result.user.email).toBe("john@example.com");
    });
  });
});

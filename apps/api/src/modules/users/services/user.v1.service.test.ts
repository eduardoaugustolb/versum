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
    findByEmail: vi.fn<() => Promise<User | null>>(),
  });

  beforeEach(() => {
    const mockRepository = createMockRepository();
    service = new UserServiceV1({
      repository:
        mockRepository as unknown as import("../repositories/user.repository").UserRepository,
    });
    vi.clearAllMocks();
  });

  describe("createUser", () => {
    it("should create a new user with normalized email", async () => {
      const mockRepository = createMockRepository();
      mockRepository.create.mockResolvedValue(mockUser);
      service = new UserServiceV1({
        repository:
          mockRepository as unknown as import("../repositories/user.repository").UserRepository,
      });

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
      service = new UserServiceV1({
        repository:
          mockRepository as unknown as import("../repositories/user.repository").UserRepository,
      });

      await expect(
        service.createUser({ email: "invalid-email" }),
      ).rejects.toThrow("Invalid email format");
    });

    it("should throw error for empty email", async () => {
      const mockRepository = createMockRepository();
      service = new UserServiceV1({
        repository:
          mockRepository as unknown as import("../repositories/user.repository").UserRepository,
      });

      await expect(service.createUser({ email: "   " })).rejects.toThrow(
        "Email is required",
      );
    });

    it("should throw error for email exceeding max length", async () => {
      const mockRepository = createMockRepository();
      service = new UserServiceV1({
        repository:
          mockRepository as unknown as import("../repositories/user.repository").UserRepository,
      });

      const longEmail = `${"a".repeat(250)}@example.com`;
      await expect(service.createUser({ email: longEmail })).rejects.toThrow(
        "Email must not exceed 255 characters",
      );
    });
  });

  describe("getUserById", () => {
    it("should return user when found", async () => {
      const mockRepository = createMockRepository();
      mockRepository.findById.mockResolvedValue(mockUser);
      service = new UserServiceV1({
        repository:
          mockRepository as unknown as import("../repositories/user.repository").UserRepository,
      });

      const result = await service.getUserById({ id: mockUser.id });

      expect(result).toEqual(mockUser);
      expect(mockRepository.findById).toHaveBeenCalledWith({ id: mockUser.id });
    });

    it("should throw error when user not found", async () => {
      const mockRepository = createMockRepository();
      mockRepository.findById.mockResolvedValue(null);
      service = new UserServiceV1({
        repository:
          mockRepository as unknown as import("../repositories/user.repository").UserRepository,
      });

      await expect(
        service.getUserById({ id: "nonexistent-id" }),
      ).rejects.toThrow("User not found");
    });
  });

  describe("getUserByEmail", () => {
    it("should return user when found", async () => {
      const mockRepository = createMockRepository();
      mockRepository.findByEmail.mockResolvedValue(mockUser);
      service = new UserServiceV1({
        repository:
          mockRepository as unknown as import("../repositories/user.repository").UserRepository,
      });

      const result = await service.getUserByEmail({ email: mockUser.email });

      expect(result).toEqual(mockUser);
      expect(mockRepository.findByEmail).toHaveBeenCalledWith({
        email: mockUser.email,
      });
    });

    it("should return null when user not found", async () => {
      const mockRepository = createMockRepository();
      mockRepository.findByEmail.mockResolvedValue(null);
      service = new UserServiceV1({
        repository:
          mockRepository as unknown as import("../repositories/user.repository").UserRepository,
      });

      const result = await service.getUserByEmail({
        email: "nonexistent@example.com",
      });

      expect(result).toBeNull();
    });
  });
});

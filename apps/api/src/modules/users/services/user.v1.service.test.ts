import { describe, it, expect, vi, beforeEach } from "vitest";
import { UserServiceV1 } from "./user.v1.service";
import type { User } from "../repositories/user.types.repository";

describe("UserServiceV1", () => {
  let service: UserServiceV1;

  const mockUser: User = {
    id: "123e4567-e89b-12d3-a456-426614174000",
    username: "johndoe",
    name: "John Doe",
    email: "john@example.com",
    bio: "Software developer",
    pictureUrl: "https://example.com/avatar.jpg",
    createdAt: new Date("2024-01-01T00:00:00Z"),
  };

  const createMockRepository = () => ({
    db: {},
    create: vi.fn<() => Promise<User>>(),
    findById: vi.fn<() => Promise<User | null>>(),
    findByEmail: vi.fn<() => Promise<User | null>>(),
    findByUsername: vi.fn<() => Promise<User | null>>(),
    updateUser: vi.fn<() => Promise<User>>(),
  });

  beforeEach(() => {
    const mockRepository = createMockRepository();
    service = new UserServiceV1({ repository: mockRepository as unknown as import("../repositories/user.repository").UserRepository });
    vi.clearAllMocks();
  });

  describe("createUser", () => {
    it("should create a new user", async () => {
      const mockRepository = createMockRepository();
      mockRepository.create.mockResolvedValue(mockUser);
      service = new UserServiceV1({ repository: mockRepository as unknown as import("../repositories/user.repository").UserRepository });

      const createParams = {
        username: "janedoe",
        name: "Jane Doe",
        email: "jane@example.com",
        bio: null,
        pictureUrl: null,
      };

      const result = await service.createUser(createParams);

      expect(result).toEqual(mockUser);
      expect(mockRepository.create).toHaveBeenCalledWith(createParams);
    });
  });

  describe("getUserById", () => {
    it("should return user when found", async () => {
      const mockRepository = createMockRepository();
      mockRepository.findById.mockResolvedValue(mockUser);
      service = new UserServiceV1({ repository: mockRepository as unknown as import("../repositories/user.repository").UserRepository });

      const result = await service.getUserById({ id: mockUser.id });

      expect(result).toEqual(mockUser);
      expect(mockRepository.findById).toHaveBeenCalledWith({ id: mockUser.id });
    });

    it("should throw error when user not found", async () => {
      const mockRepository = createMockRepository();
      mockRepository.findById.mockResolvedValue(null);
      service = new UserServiceV1({ repository: mockRepository as unknown as import("../repositories/user.repository").UserRepository });

      await expect(
        service.getUserById({ id: "nonexistent-id" })
      ).rejects.toThrow("User not found");
    });
  });

  describe("updateUser", () => {
    it("should update user when found", async () => {
      const mockRepository = createMockRepository();
      const updatedUser = { ...mockUser, name: "John Updated" };
      mockRepository.findById.mockResolvedValue(mockUser);
      mockRepository.updateUser.mockResolvedValue(updatedUser);
      service = new UserServiceV1({ repository: mockRepository as unknown as import("../repositories/user.repository").UserRepository });

      const result = await service.updateUser({
        id: mockUser.id,
        name: "John Updated",
      });

      expect(result.name).toBe("John Updated");
      expect(mockRepository.updateUser).toHaveBeenCalledWith({
        id: mockUser.id,
        name: "John Updated",
      });
    });

    it("should throw error when user to update not found", async () => {
      const mockRepository = createMockRepository();
      mockRepository.findById.mockResolvedValue(null);
      service = new UserServiceV1({ repository: mockRepository as unknown as import("../repositories/user.repository").UserRepository });

      await expect(
        service.updateUser({
          id: "nonexistent-id",
          name: "New Name",
        })
      ).rejects.toThrow("User not found");
    });

    it("should preserve user id when updating", async () => {
      const mockRepository = createMockRepository();
      const updatedUser = { ...mockUser, bio: "New bio" };
      mockRepository.findById.mockResolvedValue(mockUser);
      mockRepository.updateUser.mockResolvedValue(updatedUser);
      service = new UserServiceV1({ repository: mockRepository as unknown as import("../repositories/user.repository").UserRepository });

      await service.updateUser({
        id: mockUser.id,
        bio: "New bio",
      });

      expect(mockRepository.updateUser).toHaveBeenCalledWith({
        id: mockUser.id,
        bio: "New bio",
      });
    });
  });

  describe("getUserByEmail", () => {
    it("should return user when found", async () => {
      const mockRepository = createMockRepository();
      mockRepository.findByEmail.mockResolvedValue(mockUser);
      service = new UserServiceV1({ repository: mockRepository as unknown as import("../repositories/user.repository").UserRepository });

      const result = await service.getUserByEmail({ email: mockUser.email });

      expect(result).toEqual(mockUser);
      expect(mockRepository.findByEmail).toHaveBeenCalledWith({ email: mockUser.email });
    });

    it("should return null when user not found", async () => {
      const mockRepository = createMockRepository();
      mockRepository.findByEmail.mockResolvedValue(null);
      service = new UserServiceV1({ repository: mockRepository as unknown as import("../repositories/user.repository").UserRepository });

      const result = await service.getUserByEmail({ email: "nonexistent@example.com" });

      expect(result).toBeNull();
    });
  });

  describe("getUserByUsername", () => {
    it("should return user when found", async () => {
      const mockRepository = createMockRepository();
      mockRepository.findByUsername.mockResolvedValue(mockUser);
      service = new UserServiceV1({ repository: mockRepository as unknown as import("../repositories/user.repository").UserRepository });

      const result = await service.getUserByUsername({ username: mockUser.username });

      expect(result).toEqual(mockUser);
      expect(mockRepository.findByUsername).toHaveBeenCalledWith({ username: mockUser.username });
    });

    it("should return null when user not found", async () => {
      const mockRepository = createMockRepository();
      mockRepository.findByUsername.mockResolvedValue(null);
      service = new UserServiceV1({ repository: mockRepository as unknown as import("../repositories/user.repository").UserRepository });

      const result = await service.getUserByUsername({ username: "nonexistent" });

      expect(result).toBeNull();
    });
  });
});

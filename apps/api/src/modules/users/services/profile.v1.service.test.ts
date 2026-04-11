import { beforeEach, describe, expect, it, vi } from "vitest";
import type { Profile } from "../repositories/profile.types.repository";
import { ProfileServiceV1 } from "./profile.v1.service";

describe("ProfileServiceV1", () => {
  let service: ProfileServiceV1;

  const mockProfile: Profile = {
    id: "123e4567-e89b-12d3-a456-426614174001",
    userId: "123e4567-e89b-12d3-a456-426614174000",
    username: "johndoe",
    name: "John Doe",
    bio: "Software developer",
    pictureUrl: "https://example.com/avatar.jpg",
    createdAt: new Date("2024-01-01T00:00:00Z"),
    updatedAt: new Date("2024-01-01T00:00:00Z"),
  };

  const createMockRepository = () => ({
    db: {},
    create: vi.fn<() => Promise<Profile>>(),
    findById: vi.fn<() => Promise<Profile | null>>(),
    findByUserId: vi.fn<() => Promise<Profile | null>>(),
    findByUsername: vi.fn<() => Promise<Profile | null>>(),
    update: vi.fn<() => Promise<Profile>>(),
    existsByUsername: vi.fn<() => Promise<boolean>>(),
  });

  beforeEach(() => {
    const mockRepository = createMockRepository();
    service = new ProfileServiceV1({
      repository:
        mockRepository as unknown as import("../repositories/profile.repository").ProfileRepository,
    });
    vi.clearAllMocks();
  });

  describe("getProfileByUserId", () => {
    it("should return profile when found", async () => {
      const mockRepository = createMockRepository();
      mockRepository.findByUserId.mockResolvedValue(mockProfile);
      service = new ProfileServiceV1({
        repository:
          mockRepository as unknown as import("../repositories/profile.repository").ProfileRepository,
      });

      const result = await service.getProfileByUserId({
        userId: mockProfile.userId,
      });

      expect(result).toEqual(mockProfile);
      expect(mockRepository.findByUserId).toHaveBeenCalledWith({
        userId: mockProfile.userId,
      });
    });

    it("should return null when profile not found", async () => {
      const mockRepository = createMockRepository();
      mockRepository.findByUserId.mockResolvedValue(null);
      service = new ProfileServiceV1({
        repository:
          mockRepository as unknown as import("../repositories/profile.repository").ProfileRepository,
      });

      const result = await service.getProfileByUserId({
        userId: "nonexistent",
      });

      expect(result).toBeNull();
    });
  });

  describe("getProfileByUsername", () => {
    it("should return profile when found", async () => {
      const mockRepository = createMockRepository();
      mockRepository.findByUsername.mockResolvedValue(mockProfile);
      service = new ProfileServiceV1({
        repository:
          mockRepository as unknown as import("../repositories/profile.repository").ProfileRepository,
      });

      const result = await service.getProfileByUsername({
        username: mockProfile.username,
      });

      expect(result).toEqual(mockProfile);
      expect(mockRepository.findByUsername).toHaveBeenCalledWith({
        username: mockProfile.username,
      });
    });

    it("should return null when profile not found", async () => {
      const mockRepository = createMockRepository();
      mockRepository.findByUsername.mockResolvedValue(null);
      service = new ProfileServiceV1({
        repository:
          mockRepository as unknown as import("../repositories/profile.repository").ProfileRepository,
      });

      const result = await service.getProfileByUsername({
        username: "nonexistent",
      });

      expect(result).toBeNull();
    });
  });

  describe("createProfile", () => {
    it("should create a new profile with sanitized data", async () => {
      const mockRepository = createMockRepository();
      mockRepository.existsByUsername.mockResolvedValue(false);
      mockRepository.create.mockResolvedValue(mockProfile);
      service = new ProfileServiceV1({
        repository:
          mockRepository as unknown as import("../repositories/profile.repository").ProfileRepository,
      });

      const createParams = {
        userId: mockProfile.userId,
        username: "JohnDoe",
        name: "<script>alert('xss')</script>John Doe",
        bio: null,
        pictureUrl: null,
      };

      const result = await service.createProfile(createParams);

      expect(result).toEqual(mockProfile);
      expect(mockRepository.create).toHaveBeenCalledWith(
        expect.objectContaining({
          username: "johndoe",
          name: "&lt;script&gt;alert(&#x27;xss&#x27;)&lt;&#x2F;script&gt;John Doe",
        }),
      );
    });

    it("should throw error for invalid username format", async () => {
      const mockRepository = createMockRepository();
      service = new ProfileServiceV1({
        repository:
          mockRepository as unknown as import("../repositories/profile.repository").ProfileRepository,
      });

      await expect(
        service.createProfile({
          userId: mockProfile.userId,
          username: "john doe",
          name: "John Doe",
        }),
      ).rejects.toThrow(
        "Username can only contain letters, numbers, and underscores",
      );
    });

    it("should throw error for username too short", async () => {
      const mockRepository = createMockRepository();
      service = new ProfileServiceV1({
        repository:
          mockRepository as unknown as import("../repositories/profile.repository").ProfileRepository,
      });

      await expect(
        service.createProfile({
          userId: mockProfile.userId,
          username: "ab",
          name: "John Doe",
        }),
      ).rejects.toThrow("Username must be at least 3 characters");
    });

    it("should throw error for username already in use", async () => {
      const mockRepository = createMockRepository();
      mockRepository.existsByUsername.mockResolvedValue(true);
      service = new ProfileServiceV1({
        repository:
          mockRepository as unknown as import("../repositories/profile.repository").ProfileRepository,
      });

      await expect(
        service.createProfile({
          userId: mockProfile.userId,
          username: "johndoe",
          name: "John Doe",
        }),
      ).rejects.toThrow("Username already in use");
    });

    it("should throw error for empty name", async () => {
      const mockRepository = createMockRepository();
      service = new ProfileServiceV1({
        repository:
          mockRepository as unknown as import("../repositories/profile.repository").ProfileRepository,
      });

      await expect(
        service.createProfile({
          userId: mockProfile.userId,
          username: "johndoe",
          name: "",
        }),
      ).rejects.toThrow("Name is required");
    });

    it("should throw error for bio exceeding max length", async () => {
      const mockRepository = createMockRepository();
      service = new ProfileServiceV1({
        repository:
          mockRepository as unknown as import("../repositories/profile.repository").ProfileRepository,
      });

      const longBio = "a".repeat(501);
      await expect(
        service.createProfile({
          userId: mockProfile.userId,
          username: "johndoe",
          name: "John Doe",
          bio: longBio,
        }),
      ).rejects.toThrow("Bio must not exceed 500 characters");
    });

    it("should throw error for invalid picture URL", async () => {
      const mockRepository = createMockRepository();
      service = new ProfileServiceV1({
        repository:
          mockRepository as unknown as import("../repositories/profile.repository").ProfileRepository,
      });

      await expect(
        service.createProfile({
          userId: mockProfile.userId,
          username: "johndoe",
          name: "John Doe",
          pictureUrl: "not-a-valid-url",
        }),
      ).rejects.toThrow("Picture URL must be a valid URL");
    });

    it("should throw error for non-HTTPS picture URL", async () => {
      const mockRepository = createMockRepository();
      service = new ProfileServiceV1({
        repository:
          mockRepository as unknown as import("../repositories/profile.repository").ProfileRepository,
      });

      await expect(
        service.createProfile({
          userId: mockProfile.userId,
          username: "johndoe",
          name: "John Doe",
          pictureUrl: "http://example.com/avatar.jpg",
        }),
      ).rejects.toThrow("Picture URL must use HTTPS");
    });
  });

  describe("updateProfile", () => {
    it("should update profile with sanitized data", async () => {
      const mockRepository = createMockRepository();
      mockRepository.findByUserId.mockResolvedValue(mockProfile);
      mockRepository.existsByUsername.mockResolvedValue(false);
      mockRepository.update.mockResolvedValue({
        ...mockProfile,
        name: "Jane Doe",
      });
      service = new ProfileServiceV1({
        repository:
          mockRepository as unknown as import("../repositories/profile.repository").ProfileRepository,
      });

      const result = await service.updateProfile({
        userId: mockProfile.userId,
        name: "<b>Jane Doe</b>",
      });

      expect(result.name).toBe("Jane Doe");
      expect(mockRepository.update).toHaveBeenCalledWith(
        expect.objectContaining({
          id: mockProfile.id,
          name: "&lt;b&gt;Jane Doe&lt;&#x2F;b&gt;",
        }),
      );
    });

    it("should throw error when profile not found", async () => {
      const mockRepository = createMockRepository();
      mockRepository.findByUserId.mockResolvedValue(null);
      service = new ProfileServiceV1({
        repository:
          mockRepository as unknown as import("../repositories/profile.repository").ProfileRepository,
      });

      await expect(
        service.updateProfile({
          userId: "nonexistent",
          name: "John Updated",
        }),
      ).rejects.toThrow("Profile not found");
    });

    it("should throw error for duplicate username", async () => {
      const mockRepository = createMockRepository();
      mockRepository.findByUserId.mockResolvedValue(mockProfile);
      mockRepository.existsByUsername.mockResolvedValue(true);
      service = new ProfileServiceV1({
        repository:
          mockRepository as unknown as import("../repositories/profile.repository").ProfileRepository,
      });

      await expect(
        service.updateProfile({
          userId: mockProfile.userId,
          username: "existinguser",
        }),
      ).rejects.toThrow("Username already in use");
    });
  });

  describe("validation", () => {
    it("should normalize username to lowercase", async () => {
      const mockRepository = createMockRepository();
      mockRepository.existsByUsername.mockResolvedValue(false);
      mockRepository.create.mockResolvedValue(mockProfile);
      service = new ProfileServiceV1({
        repository:
          mockRepository as unknown as import("../repositories/profile.repository").ProfileRepository,
      });

      await service.createProfile({
        userId: mockProfile.userId,
        username: "JOHNDOE",
        name: "John Doe",
      });

      expect(mockRepository.create).toHaveBeenCalledWith(
        expect.objectContaining({
          username: "johndoe",
        }),
      );
    });

    it("should sanitize HTML from name", async () => {
      const mockRepository = createMockRepository();
      mockRepository.existsByUsername.mockResolvedValue(false);
      mockRepository.create.mockResolvedValue(mockProfile);
      service = new ProfileServiceV1({
        repository:
          mockRepository as unknown as import("../repositories/profile.repository").ProfileRepository,
      });

      await service.createProfile({
        userId: mockProfile.userId,
        username: "johndoe",
        name: '<img src="x" onerror="alert(1)">John',
      });

      expect(mockRepository.create).toHaveBeenCalledWith(
        expect.objectContaining({
          name: "&lt;img src=&quot;x&quot; onerror=&quot;alert(1)&quot;&gt;John",
        }),
      );
    });
  });
});

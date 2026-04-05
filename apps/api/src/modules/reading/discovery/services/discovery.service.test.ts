import { describe, it, expect, vi, beforeEach } from "vitest";
import { DiscoveryService } from "./discovery.service";
import type { DiscoveryRepository } from "../repositories/discovery.repository";

describe("DiscoveryService", () => {
  let service: DiscoveryService;

  const mockChapter = {
    chapter: {
      id: "223e4567-e89b-12d3-a456-426614174001",
      number: 1,
    },
    book: {
      id: "123e4567-e89b-12d3-a456-426614174000",
      name: "Genesis",
      slug: "genesis",
    },
  };

  const mockVerses = [
    {
      id: "323e4567-e89b-12d3-a456-426614174002",
      number: 1,
      text: "In the beginning God created the heavens and the earth.",
    },
    {
      id: "423e4567-e89b-12d3-a456-426614174003",
      number: 2,
      text: "And the earth was without form, and void.",
    },
  ];

  const createMockRepository = () => ({
    db: {},
    getRandomVersesForChapter: vi.fn<() => Promise<typeof mockVerses>>(),
    findChapterById: vi.fn<() => Promise<typeof mockChapter | null>>(),
    getVersesByIds: vi.fn<() => Promise<unknown[]>>(),
    getReadVersesCount: vi.fn<() => Promise<number>>(),
  });

  beforeEach(() => {
    const mockRepository = createMockRepository();
    service = new DiscoveryService({ repository: mockRepository as unknown as DiscoveryRepository });
    vi.clearAllMocks();
  });

  describe("getNextVerses", () => {
    it("should return verses with context when chapter exists", async () => {
      const mockRepository = createMockRepository();
      mockRepository.findChapterById.mockResolvedValue(mockChapter);
      mockRepository.getRandomVersesForChapter.mockResolvedValue(mockVerses);
      service = new DiscoveryService({ repository: mockRepository as unknown as DiscoveryRepository });

      const result = await service.getNextVerses(mockChapter.chapter.id);

      expect(result).toEqual([
        {
          id: mockVerses[0]?.id,
          number: mockVerses[0]?.number,
          text: mockVerses[0]?.text,
          chapter: {
            id: mockChapter.chapter.id,
            number: mockChapter.chapter.number,
          },
          book: {
            id: mockChapter.book.id,
            name: mockChapter.book.name,
            slug: mockChapter.book.slug,
          },
        },
        {
          id: mockVerses[1]?.id,
          number: mockVerses[1]?.number,
          text: mockVerses[1]?.text,
          chapter: {
            id: mockChapter.chapter.id,
            number: mockChapter.chapter.number,
          },
          book: {
            id: mockChapter.book.id,
            name: mockChapter.book.name,
            slug: mockChapter.book.slug,
          },
        },
      ]);
    });

    it("should throw error when chapter not found", async () => {
      const mockRepository = createMockRepository();
      mockRepository.findChapterById.mockResolvedValue(null);
      service = new DiscoveryService({ repository: mockRepository as unknown as DiscoveryRepository });

      await expect(
        service.getNextVerses("nonexistent-chapter-id")
      ).rejects.toThrow("Chapter not found");
    });

    it("should call repository with correct chapter id", async () => {
      const mockRepository = createMockRepository();
      const chapterId = "test-chapter-id";
      mockRepository.findChapterById.mockResolvedValue(mockChapter);
      mockRepository.getRandomVersesForChapter.mockResolvedValue(mockVerses);
      service = new DiscoveryService({ repository: mockRepository as unknown as DiscoveryRepository });

      await service.getNextVerses(chapterId);

      expect(mockRepository.findChapterById).toHaveBeenCalledWith(chapterId);
      expect(mockRepository.getRandomVersesForChapter).toHaveBeenCalledWith(chapterId);
    });

    it("should return empty array when no verses found", async () => {
      const mockRepository = createMockRepository();
      mockRepository.findChapterById.mockResolvedValue(mockChapter);
      mockRepository.getRandomVersesForChapter.mockResolvedValue([]);
      service = new DiscoveryService({ repository: mockRepository as unknown as DiscoveryRepository });

      const result = await service.getNextVerses(mockChapter.chapter.id);

      expect(result).toEqual([]);
    });
  });

  describe("markVersesAsRead", () => {
    it("should call repository with correct params", async () => {
      const mockRepository = createMockRepository();
      mockRepository.getVersesByIds.mockResolvedValue([]);
      service = new DiscoveryService({ repository: mockRepository as unknown as DiscoveryRepository });

      const userId = "user-123";
      const verseIds = ["verse-1", "verse-2"];

      await service.markVersesAsRead({ userId, verseIds });

      expect(mockRepository.getVersesByIds).toHaveBeenCalledWith(verseIds);
    });
  });

  describe("getStats", () => {
    it("should return stats with verses read count", async () => {
      const mockRepository = createMockRepository();
      mockRepository.getReadVersesCount.mockResolvedValue(42);
      service = new DiscoveryService({ repository: mockRepository as unknown as DiscoveryRepository });

      const result = await service.getStats("user-123");

      expect(result).toEqual({ versesRead: 42 });
      expect(mockRepository.getReadVersesCount).toHaveBeenCalledWith("user-123");
    });

    it("should return 0 verses read when user has no readings", async () => {
      const mockRepository = createMockRepository();
      mockRepository.getReadVersesCount.mockResolvedValue(0);
      service = new DiscoveryService({ repository: mockRepository as unknown as DiscoveryRepository });

      const result = await service.getStats("user-with-no-readings");

      expect(result.versesRead).toBe(0);
    });
  });
});

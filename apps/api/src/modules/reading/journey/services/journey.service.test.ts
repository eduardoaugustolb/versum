import { describe, it, expect, vi, beforeEach } from "vitest";
import { JourneyService } from "./journey.service";
import type { JourneyRepository } from "../repositories/journey.repository";

describe("JourneyService", () => {
  let service: JourneyService;

  const mockChapterData = {
    chapter: {
      id: "223e4567-e89b-12d3-a456-426614174001",
      number: 1,
      totalVerses: 31,
    },
    book: {
      id: "123e4567-e89b-12d3-a456-426614174000",
      order: 1,
      name: "Genesis",
      slug: "genesis",
    },
    verses: [
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
    ],
  };

  const createMockRepository = () => ({
    db: {},
    findCurrentChapter: vi.fn<() => Promise<typeof mockChapterData | null>>(),
    findNextChapter: vi.fn<() => Promise<typeof mockChapterData | null>>(),
    markChapterAsRead: vi.fn<() => Promise<void>>(),
    getReadChaptersCount: vi.fn<() => Promise<number>>(),
    hasAnyReadings: vi.fn<() => Promise<boolean>>(),
  });

  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe("getCurrentChapter", () => {
    it("should return current chapter when user has readings", async () => {
      const mockRepository = createMockRepository();
      mockRepository.findCurrentChapter.mockResolvedValue(mockChapterData);
      service = new JourneyService({ repository: mockRepository as unknown as JourneyRepository });

      const result = await service.getCurrentChapter("user-123");

      expect(result).toEqual({
        chapter: {
          id: mockChapterData.chapter.id,
          number: mockChapterData.chapter.number,
          totalVerses: mockChapterData.chapter.totalVerses,
        },
        book: {
          id: mockChapterData.book.id,
          order: mockChapterData.book.order,
          name: mockChapterData.book.name,
          slug: mockChapterData.book.slug,
        },
        verses: mockChapterData.verses,
      });
    });

    it("should return first chapter when user has no readings", async () => {
      const mockRepository = createMockRepository();
      mockRepository.findCurrentChapter.mockResolvedValue(mockChapterData);
      service = new JourneyService({ repository: mockRepository as unknown as JourneyRepository });

      const result = await service.getCurrentChapter("new-user");

      expect(result).toBeDefined();
      expect(result?.chapter.id).toBe(mockChapterData.chapter.id);
    });

    it("should return undefined when no chapter found", async () => {
      const mockRepository = createMockRepository();
      mockRepository.findCurrentChapter.mockResolvedValue(null);
      service = new JourneyService({ repository: mockRepository as unknown as JourneyRepository });

      const result = await service.getCurrentChapter("user-without-chapters");

      expect(result).toBeUndefined();
    });
  });

  describe("getNextChapter", () => {
    it("should return first chapter when user has no readings", async () => {
      const mockRepository = createMockRepository();
      mockRepository.hasAnyReadings.mockResolvedValue(false);
      mockRepository.findCurrentChapter.mockResolvedValue(mockChapterData);
      service = new JourneyService({ repository: mockRepository as unknown as JourneyRepository });

      const result = await service.getNextChapter("new-user");

      expect(result).toBeDefined();
      expect(result?.chapter.id).toBe(mockChapterData.chapter.id);
      expect(mockRepository.hasAnyReadings).toHaveBeenCalledWith("new-user");
    });

    it("should return next chapter and mark it as read when user has readings", async () => {
      const mockRepository = createMockRepository();
      const nextChapter = { ...mockChapterData, chapter: { ...mockChapterData.chapter, id: "next-chapter-id", number: 2 } };
      mockRepository.hasAnyReadings.mockResolvedValue(true);
      mockRepository.findCurrentChapter.mockResolvedValue(mockChapterData);
      mockRepository.findNextChapter.mockResolvedValue(nextChapter);
      service = new JourneyService({ repository: mockRepository as unknown as JourneyRepository });

      const result = await service.getNextChapter("user-123");

      expect(result).toBeDefined();
      expect(result?.chapter.id).toBe("next-chapter-id");
      expect(mockRepository.markChapterAsRead).toHaveBeenCalledWith({
        userId: "user-123",
        chapterId: "next-chapter-id",
      });
    });

    it("should return undefined when no next chapter exists", async () => {
      const mockRepository = createMockRepository();
      mockRepository.hasAnyReadings.mockResolvedValue(true);
      mockRepository.findCurrentChapter.mockResolvedValue(mockChapterData);
      mockRepository.findNextChapter.mockResolvedValue(null);
      service = new JourneyService({ repository: mockRepository as unknown as JourneyRepository });

      const result = await service.getNextChapter("user-at-end");

      expect(result).toBeUndefined();
      expect(mockRepository.markChapterAsRead).not.toHaveBeenCalled();
    });

    it("should not mark chapter as read when no next chapter found", async () => {
      const mockRepository = createMockRepository();
      mockRepository.hasAnyReadings.mockResolvedValue(true);
      mockRepository.findCurrentChapter.mockResolvedValue(mockChapterData);
      mockRepository.findNextChapter.mockResolvedValue(null);
      service = new JourneyService({ repository: mockRepository as unknown as JourneyRepository });

      await service.getNextChapter("user-123");

      expect(mockRepository.markChapterAsRead).not.toHaveBeenCalled();
    });
  });

  describe("getStats", () => {
    it("should return stats with chapters read count", async () => {
      const mockRepository = createMockRepository();
      mockRepository.getReadChaptersCount.mockResolvedValue(42);
      service = new JourneyService({ repository: mockRepository as unknown as JourneyRepository });

      const result = await service.getStats("user-123");

      expect(result).toEqual({ chaptersRead: 42 });
      expect(mockRepository.getReadChaptersCount).toHaveBeenCalledWith("user-123");
    });

    it("should return 0 chapters read when user has no readings", async () => {
      const mockRepository = createMockRepository();
      mockRepository.getReadChaptersCount.mockResolvedValue(0);
      service = new JourneyService({ repository: mockRepository as unknown as JourneyRepository });

      const result = await service.getStats("new-user");

      expect(result.chaptersRead).toBe(0);
    });
  });
});

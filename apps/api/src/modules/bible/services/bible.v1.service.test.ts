import { describe, it, expect, vi, beforeEach } from "vitest";
import { BibleServiceV1 } from "./bible.v1.service";
import type { BibleRepository } from "../repositories/bible.repository";
import type { Book, Chapter, Verse, PaginatedResult } from "../repositories/bible.types.repository";

describe("BibleServiceV1", () => {
  let service: BibleServiceV1;

  const mockBook: Book = {
    id: "123e4567-e89b-12d3-a456-426614174000",
    order: 1,
    name: "Genesis",
    slug: "genesis",
    niceName: "Gênesis",
    testament: "OLD",
    totalChapters: 50,
  };

  const mockChapter: Chapter = {
    id: "223e4567-e89b-12d3-a456-426614174001",
    bookId: "123e4567-e89b-12d3-a456-426614174000",
    number: 1,
    totalVerses: 31,
  };

  const mockVerse: Verse = {
    id: "323e4567-e89b-12d3-a456-426614174002",
    chapterId: "223e4567-e89b-12d3-a456-426614174001",
    number: 1,
    text: "In the beginning God created the heavens and the earth.",
  };

  const createMockRepository = () => ({
    db: {},
    findBooksPaginated: vi.fn<() => Promise<PaginatedResult<Book>>>(),
    findBookByDynamicId: vi.fn<() => Promise<Book | null>>(),
    findChaptersPaginated: vi.fn<() => Promise<PaginatedResult<Chapter>>>(),
    findChapterByNumberAndDynamicId: vi.fn<() => Promise<Chapter | null>>(),
    findVersesPaginated: vi.fn<() => Promise<PaginatedResult<Verse>>>(),
    findVerse: vi.fn<() => Promise<Verse | null>>(),
  });

  beforeEach(() => {
    const mockRepository = createMockRepository();
    service = new BibleServiceV1({ repository: mockRepository as unknown as BibleRepository });
    vi.clearAllMocks();
  });

  describe("getBooksPaginated", () => {
    it("should return paginated books", async () => {
      const mockRepository = createMockRepository();
      const paginatedResult: PaginatedResult<Book> = {
        data: [mockBook],
        total: 1,
      };
      mockRepository.findBooksPaginated.mockResolvedValue(paginatedResult);
      service = new BibleServiceV1({ repository: mockRepository as unknown as BibleRepository });

      const result = await service.getBooksPaginated({ page: 1, limit: 10 });

      expect(result).toEqual(paginatedResult);
      expect(mockRepository.findBooksPaginated).toHaveBeenCalledWith({ page: 1, limit: 10 });
    });

    it("should return empty result when no books", async () => {
      const mockRepository = createMockRepository();
      const paginatedResult: PaginatedResult<Book> = {
        data: [],
        total: 0,
      };
      mockRepository.findBooksPaginated.mockResolvedValue(paginatedResult);
      service = new BibleServiceV1({ repository: mockRepository as unknown as BibleRepository });

      const result = await service.getBooksPaginated({ page: 1, limit: 10 });

      expect(result.data).toHaveLength(0);
      expect(result.total).toBe(0);
    });
  });

  describe("getBookByDynamicId", () => {
    it("should return book when found", async () => {
      const mockRepository = createMockRepository();
      mockRepository.findBookByDynamicId.mockResolvedValue(mockBook);
      service = new BibleServiceV1({ repository: mockRepository as unknown as BibleRepository });

      const result = await service.getBookByDynamicId({ dynamicId: "genesis" });

      expect(result).toEqual(mockBook);
      expect(mockRepository.findBookByDynamicId).toHaveBeenCalledWith({ dynamicId: "genesis" });
    });

    it("should throw error when book not found", async () => {
      const mockRepository = createMockRepository();
      mockRepository.findBookByDynamicId.mockResolvedValue(null);
      service = new BibleServiceV1({ repository: mockRepository as unknown as BibleRepository });

      await expect(
        service.getBookByDynamicId({ dynamicId: "nonexistent" })
      ).rejects.toThrow("Book not found");
    });
  });

  describe("getChaptersPaginated", () => {
    it("should return paginated chapters", async () => {
      const mockRepository = createMockRepository();
      const paginatedResult: PaginatedResult<Chapter> = {
        data: [mockChapter],
        total: 1,
      };
      mockRepository.findChaptersPaginated.mockResolvedValue(paginatedResult);
      service = new BibleServiceV1({ repository: mockRepository as unknown as BibleRepository });

      const result = await service.getChaptersPaginated({
        dynamicId: "genesis",
        page: 1,
        limit: 10,
      });

      expect(result).toEqual(paginatedResult);
    });

    it("should throw error when book not found or no chapters", async () => {
      const mockRepository = createMockRepository();
      const paginatedResult: PaginatedResult<Chapter> = {
        data: [],
        total: 0,
      };
      mockRepository.findChaptersPaginated.mockResolvedValue(paginatedResult);
      service = new BibleServiceV1({ repository: mockRepository as unknown as BibleRepository });

      await expect(
        service.getChaptersPaginated({
          dynamicId: "nonexistent",
          page: 1,
          limit: 10,
        })
      ).rejects.toThrow("Book not found or no chapters available");
    });
  });

  describe("getChapter", () => {
    it("should return chapter when found", async () => {
      const mockRepository = createMockRepository();
      mockRepository.findChapterByNumberAndDynamicId.mockResolvedValue(mockChapter);
      service = new BibleServiceV1({ repository: mockRepository as unknown as BibleRepository });

      const result = await service.getChapter({
        dynamicId: "genesis",
        chapterNumber: 1,
      });

      expect(result).toEqual(mockChapter);
    });

    it("should throw error when chapter not found", async () => {
      const mockRepository = createMockRepository();
      mockRepository.findChapterByNumberAndDynamicId.mockResolvedValue(null);
      service = new BibleServiceV1({ repository: mockRepository as unknown as BibleRepository });

      await expect(
        service.getChapter({
          dynamicId: "genesis",
          chapterNumber: 999,
        })
      ).rejects.toThrow("Chapter not found");
    });
  });

  describe("getVersesPaginated", () => {
    it("should return paginated verses", async () => {
      const mockRepository = createMockRepository();
      const paginatedResult: PaginatedResult<Verse> = {
        data: [mockVerse],
        total: 1,
      };
      mockRepository.findVersesPaginated.mockResolvedValue(paginatedResult);
      service = new BibleServiceV1({ repository: mockRepository as unknown as BibleRepository });

      const result = await service.getVersesPaginated({
        dynamicId: "genesis",
        chapterNumber: 1,
        page: 1,
        limit: 10,
      });

      expect(result).toEqual(paginatedResult);
    });

    it("should throw error when chapter not found or no verses", async () => {
      const mockRepository = createMockRepository();
      const paginatedResult: PaginatedResult<Verse> = {
        data: [],
        total: 0,
      };
      mockRepository.findVersesPaginated.mockResolvedValue(paginatedResult);
      service = new BibleServiceV1({ repository: mockRepository as unknown as BibleRepository });

      await expect(
        service.getVersesPaginated({
          dynamicId: "genesis",
          chapterNumber: 999,
          page: 1,
          limit: 10,
        })
      ).rejects.toThrow("Chapter not found or no verses available");
    });
  });

  describe("getVerse", () => {
    it("should return verse when found", async () => {
      const mockRepository = createMockRepository();
      mockRepository.findVerse.mockResolvedValue(mockVerse);
      service = new BibleServiceV1({ repository: mockRepository as unknown as BibleRepository });

      const result = await service.getVerse({
        dynamicId: "genesis",
        chapterNumber: 1,
        verseNumber: 1,
      });

      expect(result).toEqual(mockVerse);
    });

    it("should throw error when verse not found", async () => {
      const mockRepository = createMockRepository();
      mockRepository.findVerse.mockResolvedValue(null);
      service = new BibleServiceV1({ repository: mockRepository as unknown as BibleRepository });

      await expect(
        service.getVerse({
          dynamicId: "genesis",
          chapterNumber: 1,
          verseNumber: 999,
        })
      ).rejects.toThrow("Verse not found");
    });
  });
});

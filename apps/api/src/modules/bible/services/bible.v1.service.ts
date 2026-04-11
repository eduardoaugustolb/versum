import { BibleRepository } from "../repositories/bible.repository.ts";
import type {
  Book,
  Chapter,
  PaginatedResult,
  PaginationParams,
  Verse,
  VerseParams,
} from "../repositories/bible.types.repository.ts";

export class BibleServiceV1 {
  private readonly repository: BibleRepository;

  constructor({ repository }: { repository?: BibleRepository } = {}) {
    this.repository = repository ?? new BibleRepository();
  }

  async getBooksPaginated(
    params: PaginationParams,
  ): Promise<PaginatedResult<Book>> {
    return this.repository.findBooksPaginated(params);
  }

  async getBookByDynamicId({ dynamicId }: { dynamicId: string }) {
    const book = await this.repository.findBookByDynamicId({ dynamicId });

    if (!book) {
      throw new Error("Book not found");
    }

    return book;
  }

  async getChaptersPaginated(
    params: PaginationParams & { dynamicId: string },
  ): Promise<PaginatedResult<Chapter>> {
    const result = await this.repository.findChaptersPaginated(params);

    if (result.total === 0) {
      throw new Error("Book not found or no chapters available");
    }

    return result;
  }

  async getChapter({
    dynamicId,
    chapterNumber,
  }: {
    dynamicId: string;
    chapterNumber: number;
  }) {
    const chapter = await this.repository.findChapterByNumberAndDynamicId({
      dynamicId,
      chapterNumber,
    });

    if (!chapter) {
      throw new Error("Chapter not found");
    }

    return chapter;
  }

  async getVersesPaginated(
    params: PaginationParams & Omit<VerseParams, "verseNumber">,
  ): Promise<PaginatedResult<Verse>> {
    const result = await this.repository.findVersesPaginated(params);

    if (result.total === 0) {
      throw new Error("Chapter not found or no verses available");
    }

    return result;
  }

  async getVerse(params: VerseParams) {
    const verse = await this.repository.findVerse(params);

    if (!verse) {
      throw new Error("Verse not found");
    }

    return verse;
  }
}

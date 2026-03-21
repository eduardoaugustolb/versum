import { BibleRepository } from "../repositories/bible.repository.ts";
import type {
  Book,
  Chapter,
  Verse,
} from "../repositories/bible.types.repository.ts";

interface PaginatedResult<T> {
  data: T[];
  total: number;
}

interface PaginationParams {
  page: number;
  limit: number;
}

interface ChaptersPaginationParams extends PaginationParams {
  bookOrder: number;
}

interface VersesPaginationParams extends PaginationParams {
  bookOrder: number;
  chapterNumber: number;
}

export class BibleServiceV1 {
  private readonly repository: BibleRepository;

  constructor({ repository }: { repository?: BibleRepository } = {}) {
    this.repository = repository ?? new BibleRepository();
  }
  // ------------------------
  // Books
  // ------------------------

  async getBooksPaginated(
    params: PaginationParams,
  ): Promise<PaginatedResult<Book>> {
    const books = await this.repository.findAllBooks();

    const total = books.length;

    const start = (params.page - 1) * params.limit;
    const end = start + params.limit;

    const data = books.slice(start, end);

    return { data, total };
  }

  async getBookByOrder({ order }: { order: number }) {
    const book = await this.repository.findBookByOrder({
      order,
    });

    if (!book) {
      throw new Error("Book not found");
    }

    return book;
  }

  // ------------------------
  // Chapters
  // ------------------------

  async getChaptersPaginated(
    params: ChaptersPaginationParams,
  ): Promise<PaginatedResult<Chapter>> {
    const chapters = await this.repository.findChaptersByBookOrder(params);

    if (!chapters.length) {
      throw new Error("Book not found or no chapters available");
    }

    const total = chapters.length;

    const start = (params.page - 1) * params.limit;
    const end = start + params.limit;

    const data = chapters.slice(start, end);

    return { data, total };
  }

  async getChapter({
    bookOrder,
    chapterNumber,
  }: {
    bookOrder: number;
    chapterNumber: number;
  }) {
    const chapter = await this.repository.findChapterByNumberAndBookOrder({
      bookOrder,
      chapterNumber,
    });

    if (!chapter) {
      throw new Error("Chapter not found");
    }

    return chapter;
  }

  // ------------------------
  // Verses
  // ------------------------

  async getVersesPaginated(
    params: VersesPaginationParams,
  ): Promise<PaginatedResult<Verse>> {
    const verses = await this.repository.findVerses(params);

    if (!verses.length) {
      throw new Error("Chapter not found or no verses available");
    }

    const total = verses.length;

    const start = (params.page - 1) * params.limit;
    const end = start + params.limit;

    const data = verses.slice(start, end);

    return { data, total };
  }

  async getVerse(params: {
    bookOrder: number;
    chapterNumber: number;
    verseNumber: number;
  }) {
    const verse = await this.repository.findVerse(params);

    if (!verse) {
      throw new Error("Verse not found");
    }

    return verse;
  }
}

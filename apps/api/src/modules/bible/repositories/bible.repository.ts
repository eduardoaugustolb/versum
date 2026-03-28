import { and, eq, or, sql, count } from "drizzle-orm";
import {
  bibleBooks as books,
  bibleChapters as chapters,
  bibleVerses as verses,
} from "../../../infrastructure/db/schema.ts";
import type {
  iBibleRepository,
  PaginationParams,
  PaginatedResult,
  BookDynamicIdParams,
  ChapterParams,
  VerseParams,
  Book,
  Chapter,
  Verse,
} from "./bible.types.repository.ts";
import { db as drizzle } from "../../../infrastructure/db/index.ts";

export class BibleRepository implements iBibleRepository {
  private readonly db: typeof drizzle;

  constructor({ db }: { db?: typeof drizzle } = {}) {
    this.db = db ?? drizzle;
  }

  async findBooksPaginated({
    page,
    limit,
  }: PaginationParams): Promise<PaginatedResult<Book>> {
    const offset = (page - 1) * limit;

    const [data, countResult] = await Promise.all([
      this.db
        .select()
        .from(books)
        .orderBy(books.name)
        .limit(limit)
        .offset(offset),
      this.db.select({ total: count() }).from(books),
    ]);

    const total = countResult[0]?.total ?? 0;

    return { data, total: Number(total) };
  }

  async findBookByDynamicId({
    dynamicId,
  }: BookDynamicIdParams): Promise<Book | null> {
    const [book] = await this.db
      .select()
      .from(books)
      .where(or(eq(books.slug, dynamicId), eq(books.name, dynamicId)));

    return book ?? null;
  }

  async findChaptersPaginated({
    dynamicId,
    page,
    limit,
  }: PaginationParams & Omit<ChapterParams, "chapterNumber">): Promise<PaginatedResult<Chapter>> {
    const offset = (page - 1) * limit;

    const bookResult = await this.db
      .select({ id: books.id })
      .from(books)
      .where(or(eq(books.slug, dynamicId), eq(books.name, dynamicId)));

    if (!bookResult.length) {
      return { data: [], total: 0 };
    }

    const book = bookResult[0];
    if (!book) {
      return { data: [], total: 0 };
    }
    const bookId = book.id;

    const [data, countResult] = await Promise.all([
      this.db
        .select()
        .from(chapters)
        .where(eq(chapters.bookId, bookId))
        .orderBy(chapters.number)
        .limit(limit)
        .offset(offset),
      this.db
        .select({ total: count() })
        .from(chapters)
        .where(eq(chapters.bookId, bookId)),
    ]);

    const total = countResult[0]?.total ?? 0;

    return { data, total: Number(total) };
  }

  async findChapterByNumberAndDynamicId({
    dynamicId,
    chapterNumber,
  }: ChapterParams): Promise<Chapter | null> {
    const bookResult = await this.db
      .select({ id: books.id })
      .from(books)
      .where(or(eq(books.slug, dynamicId), eq(books.name, dynamicId)));

    if (!bookResult.length) {
      return null;
    }

    const book = bookResult[0];
    if (!book) {
      return null;
    }
    const bookId = book.id;

    const [chapter] = await this.db
      .select()
      .from(chapters)
      .where(
        and(eq(chapters.bookId, bookId), eq(chapters.number, chapterNumber)),
      );

    return chapter ?? null;
  }

  async findVersesPaginated({
    dynamicId,
    chapterNumber,
    page,
    limit,
  }: PaginationParams & Omit<VerseParams, "verseNumber">): Promise<PaginatedResult<Verse>> {
    const offset = (page - 1) * limit;

    const chapterResult = await this.db
      .select({ id: chapters.id })
      .from(chapters)
      .innerJoin(books, eq(chapters.bookId, books.id))
      .where(
        and(
          or(eq(books.slug, dynamicId), eq(books.name, dynamicId)),
          eq(chapters.number, chapterNumber),
        ),
      );

    if (!chapterResult.length) {
      return { data: [], total: 0 };
    }

    const chapterResultItem = chapterResult[0];
    if (!chapterResultItem) {
      return { data: [], total: 0 };
    }
    const chapterId = chapterResultItem.id;

    const [data, countResult] = await Promise.all([
      this.db
        .select()
        .from(verses)
        .where(eq(verses.chapterId, chapterId))
        .orderBy(verses.number)
        .limit(limit)
        .offset(offset),
      this.db
        .select({ total: count() })
        .from(verses)
        .where(eq(verses.chapterId, chapterId)),
    ]);

    const total = countResult[0]?.total ?? 0;

    return { data, total: Number(total) };
  }

  async findVerse({
    dynamicId,
    chapterNumber,
    verseNumber,
  }: VerseParams): Promise<Verse | null> {
    const [result] = await this.db
      .select({
        verse: verses,
      })
      .from(verses)
      .innerJoin(chapters, eq(verses.chapterId, chapters.id))
      .innerJoin(books, eq(chapters.bookId, books.id))
      .where(
        and(
          or(eq(books.slug, dynamicId), eq(books.name, dynamicId)),
          eq(chapters.number, chapterNumber),
          eq(verses.number, verseNumber),
        ),
      );

    return result?.verse ?? null;
  }
}

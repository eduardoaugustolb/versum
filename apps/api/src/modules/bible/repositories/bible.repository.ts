import { and, eq } from "drizzle-orm";
import {
  bibleBooks as books,
  bibleChapters as chapters,
  bibleVerses as verses,
} from "../../../infrastructure/db/schema.ts";
import type {
  iBibleRepository,
  ChapterParams,
  VerseParams,
} from "./bible.types.repository.ts";
import { db as drizzle } from "../../../infrastructure/db/index.ts";

export class BibleRepository implements iBibleRepository {
  private readonly db: typeof drizzle;

  constructor({ db }: { db?: typeof drizzle } = {}) {
    this.db = db ?? drizzle;
  }

  async findAllBooks() {
    return await this.db.select().from(books);
  }

  async findBookByOrder({ order }: { order: number }) {
    const [book] = await this.db
      .select()
      .from(books)
      .where(eq(books.order, order));

    return book ?? null;
  }

  async findChaptersByBookOrder({
    bookOrder,
  }: Omit<ChapterParams, "chapterNumber">) {
    return await this.db
      .select({
        chapter: chapters,
      })
      .from(chapters)
      .innerJoin(books, eq(chapters.book_id, books.id))
      .where(eq(books.order, bookOrder))
      .then((r) => r.map((row) => row.chapter));
  }

  async findChapterByNumberAndBookOrder({
    bookOrder,
    chapterNumber,
  }: ChapterParams) {
    const [result] = await this.db
      .select({
        chapter: chapters,
      })
      .from(chapters)
      .innerJoin(books, eq(chapters.book_id, books.id))
      .where(
        and(eq(books.order, bookOrder), eq(chapters.number, chapterNumber)),
      );

    return result?.chapter ?? null;
  }

  async findVerses({
    bookOrder,
    chapterNumber,
  }: Omit<VerseParams, "verseNumber">) {
    return await this.db
      .select({
        verse: verses,
      })
      .from(verses)
      .innerJoin(chapters, eq(verses.chapter_id, chapters.id))
      .innerJoin(books, eq(chapters.book_id, books.id))
      .where(
        and(eq(books.order, bookOrder), eq(chapters.number, chapterNumber)),
      )
      .then((r) => r.map((row) => row.verse));
  }

  async findVerse({ bookOrder, chapterNumber, verseNumber }: VerseParams) {
    const [result] = await this.db
      .select({
        verse: verses,
      })
      .from(verses)
      .innerJoin(chapters, eq(verses.chapter_id, chapters.id))
      .innerJoin(books, eq(chapters.book_id, books.id))
      .where(
        and(
          eq(books.order, bookOrder),
          eq(chapters.number, chapterNumber),
          eq(verses.number, verseNumber),
        ),
      );

    return result?.verse ?? null;
  }
}

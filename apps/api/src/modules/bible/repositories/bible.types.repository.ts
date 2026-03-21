import type { InferSelectModel } from "drizzle-orm";
import type {
  bibleBooks as books,
  bibleChapters as chapters,
  bibleVerses as verses,
} from "../../../infrastructure/db/schema.ts";

export type Book = InferSelectModel<typeof books>;
export type Chapter = InferSelectModel<typeof chapters>;
export type Verse = InferSelectModel<typeof verses>;

export type ChapterParams = {
  bookOrder: number;
  chapterNumber: number;
};

export type VerseParams = {
  bookOrder: number;
  chapterNumber: number;
  verseNumber: number;
};

export interface iBibleRepository {
  findAllBooks(): Promise<Book[]>;
  findBookByOrder(p: { order: number }): Promise<Book | null>;

  findChaptersByBookOrder(
    p: Omit<ChapterParams, "chapterNumber">,
  ): Promise<Chapter[]>;
  findChapterByNumberAndBookOrder(p: ChapterParams): Promise<Chapter | null>;

  findVerses(p: Omit<VerseParams, "verseNumber">): Promise<Verse[]>;

  findVerse({
    bookOrder,
    chapterNumber,
    verseNumber,
  }: VerseParams): Promise<Verse | null>;
}

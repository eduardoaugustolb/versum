import type { InferSelectModel } from "drizzle-orm";
import type {
  bibleBooks as books,
  bibleChapters as chapters,
  bibleVerses as verses,
} from "../../../infrastructure/db/schema.ts";

export type Book = InferSelectModel<typeof books>;
export type Chapter = InferSelectModel<typeof chapters>;
export type Verse = InferSelectModel<typeof verses>;

export interface PaginatedResult<T> {
  data: T[];
  total: number;
}

export interface PaginationParams {
  page: number;
  limit: number;
}

export type DynamicIdParams = {
  dynamicId: string;
};

export type BookDynamicIdParams = DynamicIdParams;

export type ChapterParams = {
  dynamicId: string;
  chapterNumber: number;
};

export type VerseParams = {
  dynamicId: string;
  chapterNumber: number;
  verseNumber: number;
};

export interface iBibleRepository {
  findBooksPaginated(p: PaginationParams): Promise<PaginatedResult<Book>>;
  findBookByDynamicId(p: BookDynamicIdParams): Promise<Book | null>;

  findChaptersPaginated(
    p: PaginationParams & Omit<ChapterParams, "chapterNumber">,
  ): Promise<PaginatedResult<Chapter>>;
  findChapterByNumberAndDynamicId(p: ChapterParams): Promise<Chapter | null>;

  findVersesPaginated(
    p: PaginationParams & Omit<VerseParams, "verseNumber">,
  ): Promise<PaginatedResult<Verse>>;
  findVerse(p: VerseParams): Promise<Verse | null>;
}

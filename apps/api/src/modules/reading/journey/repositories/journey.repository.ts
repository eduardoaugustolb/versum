import { and, count, desc, eq, gt, or } from "drizzle-orm";
import { db as drizzle } from "../../../../infrastructure/db/index.ts";
import {
  bibleBooks,
  bibleChapters,
  bibleVerses,
} from "../../../../infrastructure/db/schema.ts";
import { journeyReadings } from "../../../interactions/db/readings.table.ts";

export type ChapterWithContent = {
  chapter: typeof bibleChapters.$inferSelect;
  book: typeof bibleBooks.$inferSelect;
  verses: (typeof bibleVerses.$inferSelect)[];
};

export class JourneyRepository {
  private readonly db: typeof drizzle;

  constructor({ db }: { db?: typeof drizzle } = {}) {
    this.db = db ?? drizzle;
  }

  async findNextChapterToRead(
    userId: string,
  ): Promise<ChapterWithContent | null> {
    const lastRead = await this.db
      .select({ chapterId: journeyReadings.chapterId })
      .from(journeyReadings)
      .where(eq(journeyReadings.userId, userId))
      .orderBy(desc(journeyReadings.readAt))
      .limit(1);

    if (lastRead.length === 0 || !lastRead[0]?.chapterId) {
      const [firstChapter] = await this.db
        .select({
          chapter: bibleChapters,
          book: bibleBooks,
        })
        .from(bibleChapters)
        .innerJoin(bibleBooks, eq(bibleChapters.bookId, bibleBooks.id))
        .orderBy(bibleBooks.order, bibleChapters.number)
        .limit(1);

      if (!firstChapter) return null;
      return this.buildChapterWithVerses(
        firstChapter.chapter,
        firstChapter.book,
      );
    }

    const nextChapter = await this.findNextChapterSequential(
      lastRead[0].chapterId,
    );
    return nextChapter;
  }

  private async findNextChapterSequential(
    currentChapterId: string,
  ): Promise<ChapterWithContent | null> {
    const current = await this.db
      .select({
        chapter: bibleChapters,
        book: bibleBooks,
      })
      .from(bibleChapters)
      .innerJoin(bibleBooks, eq(bibleChapters.bookId, bibleBooks.id))
      .where(eq(bibleChapters.id, currentChapterId))
      .limit(1);

    if (!current[0]) return null;

    const { chapter: currentChapter, book: currentBook } = current[0];

    const [nextChapter] = await this.db
      .select({
        chapter: bibleChapters,
        book: bibleBooks,
      })
      .from(bibleChapters)
      .innerJoin(bibleBooks, eq(bibleChapters.bookId, bibleBooks.id))
      .where(
        or(
          gt(bibleBooks.order, currentBook.order),
          and(
            eq(bibleBooks.order, currentBook.order),
            gt(bibleChapters.number, currentChapter.number),
          ),
        ),
      )
      .orderBy(bibleBooks.order, bibleChapters.number)
      .limit(1);

    if (!nextChapter) return null;
    return this.buildChapterWithVerses(nextChapter.chapter, nextChapter.book);
  }

  async findChaptersAfter(
    currentChapterId: string,
    count: number,
  ): Promise<ChapterWithContent[]> {
    const current = await this.db
      .select({
        chapter: bibleChapters,
        book: bibleBooks,
      })
      .from(bibleChapters)
      .innerJoin(bibleBooks, eq(bibleChapters.bookId, bibleBooks.id))
      .where(eq(bibleChapters.id, currentChapterId))
      .limit(1);

    if (!current[0]) return [];

    const { chapter: currentChapter, book: currentBook } = current[0];

    const chapters = await this.db
      .select({
        chapter: bibleChapters,
        book: bibleBooks,
      })
      .from(bibleChapters)
      .innerJoin(bibleBooks, eq(bibleChapters.bookId, bibleBooks.id))
      .where(
        or(
          gt(bibleBooks.order, currentBook.order),
          and(
            eq(bibleBooks.order, currentBook.order),
            gt(bibleChapters.number, currentChapter.number),
          ),
        ),
      )
      .orderBy(bibleBooks.order, bibleChapters.number)
      .limit(count);

    return Promise.all(
      chapters.map((c) => this.buildChapterWithVerses(c.chapter, c.book)),
    );
  }

  private async buildChapterWithVerses(
    chapter: typeof bibleChapters.$inferSelect,
    book: typeof bibleBooks.$inferSelect,
  ): Promise<ChapterWithContent> {
    const verses = await this.db
      .select()
      .from(bibleVerses)
      .where(eq(bibleVerses.chapterId, chapter.id))
      .orderBy(bibleVerses.number);

    return { chapter, book, verses };
  }

  async markChapterAsRead({
    userId,
    chapterId,
  }: {
    userId: string;
    chapterId: string;
  }): Promise<void> {
    await this.db
      .insert(journeyReadings)
      .values({
        userId,
        chapterId,
        readAt: new Date(),
      })
      .onConflictDoUpdate({
        target: [journeyReadings.userId, journeyReadings.chapterId],
        set: {
          readAt: new Date(),
        },
      });
  }

  async getReadChaptersCount(userId: string): Promise<number> {
    const [result] = await this.db
      .select({ count: count() })
      .from(journeyReadings)
      .where(eq(journeyReadings.userId, userId));

    return result?.count ?? 0;
  }

  async getTotalChapters(): Promise<number> {
    const [result] = await this.db
      .select({ count: count() })
      .from(bibleChapters);

    return result?.count ?? 0;
  }

  async isChapterRead(userId: string, chapterId: string): Promise<boolean> {
    const [result] = await this.db
      .select({ count: count() })
      .from(journeyReadings)
      .where(
        and(
          eq(journeyReadings.userId, userId),
          eq(journeyReadings.chapterId, chapterId),
        ),
      );

    return (result?.count ?? 0) > 0;
  }
}

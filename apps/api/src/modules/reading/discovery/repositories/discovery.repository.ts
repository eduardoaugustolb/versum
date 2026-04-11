import { eq, sql } from "drizzle-orm";
import { db as drizzle } from "../../../../infrastructure/db/index.ts";
import {
  bibleBooks,
  bibleChapters,
  bibleVerses,
} from "../../../../infrastructure/db/schema.ts";

export class DiscoveryRepository {
  private readonly db: typeof drizzle;

  constructor({ db }: { db?: typeof drizzle } = {}) {
    this.db = db ?? drizzle;
  }

  async getRandomVersesForChapter(
    chapterId: string,
  ): Promise<(typeof bibleVerses.$inferSelect)[]> {
    const verses = await this.db
      .select()
      .from(bibleVerses)
      .where(eq(bibleVerses.chapterId, chapterId))
      .orderBy(bibleVerses.number);

    return verses;
  }

  async findChapterById(chapterId: string) {
    const [result] = await this.db
      .select({
        chapter: bibleChapters,
        book: bibleBooks,
      })
      .from(bibleChapters)
      .innerJoin(bibleBooks, eq(bibleChapters.bookId, bibleBooks.id))
      .where(eq(bibleChapters.id, chapterId));

    return result ?? null;
  }

  async getVersesByIds(verseIds: string[]) {
    return this.db
      .select()
      .from(bibleVerses)
      .innerJoin(bibleChapters, eq(bibleVerses.chapterId, bibleChapters.id))
      .innerJoin(bibleBooks, eq(bibleChapters.bookId, bibleBooks.id))
      .where(
        sql`${bibleVerses.id} IN (${sql.join(
          verseIds.map((id) => sql`${id}`),
          sql`, `,
        )})`,
      );
  }

  async getReadVersesCount(_userId: string): Promise<number> {
    return 0;
  }
}

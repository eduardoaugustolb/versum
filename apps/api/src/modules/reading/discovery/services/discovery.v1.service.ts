import { DiscoveryRepositoryV1 } from "../repositories/discovery.v1.repository.ts";

export interface VerseWithContext {
  id: string;
  number: number;
  text: string;
  chapter: {
    id: string;
    number: number;
  };
  book: {
    id: string;
    name: string;
    slug: string;
  };
}

export interface DiscoveryStats {
  versesRead: number;
}

export class DiscoveryServiceV1 {
  private readonly repository: DiscoveryRepositoryV1;

  constructor({ repository }: { repository?: DiscoveryRepositoryV1 } = {}) {
    this.repository = repository ?? new DiscoveryRepositoryV1();
  }

  async getNextVerses(chapterId: string): Promise<VerseWithContext[]> {
    const chapter = await this.repository.findChapterById(chapterId);
    if (!chapter) {
      throw new Error("Chapter not found");
    }

    const verses = await this.repository.getRandomVersesForChapter(chapterId);

    return verses.map((verse) => ({
      id: verse.id,
      number: verse.number,
      text: verse.text,
      chapter: {
        id: chapter.chapter.id,
        number: chapter.chapter.number,
      },
      book: {
        id: chapter.book.id,
        name: chapter.book.name,
        slug: chapter.book.slug,
      },
    }));
  }

  async markVersesAsRead({
    userId,
    verseIds,
  }: {
    userId: string;
    verseIds: string[];
  }): Promise<void> {
    const _verses = await this.repository.getVersesByIds(verseIds);
  }

  async getStats(userId: string): Promise<DiscoveryStats> {
    const versesRead = await this.repository.getReadVersesCount(userId);
    return { versesRead };
  }
}

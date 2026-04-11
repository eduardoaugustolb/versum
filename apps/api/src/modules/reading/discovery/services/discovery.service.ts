import { DiscoveryRepository } from "../repositories/discovery.repository.ts";

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

export class DiscoveryService {
  private readonly repository: DiscoveryRepository;

  constructor({ repository }: { repository?: DiscoveryRepository } = {}) {
    this.repository = repository ?? new DiscoveryRepository();
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

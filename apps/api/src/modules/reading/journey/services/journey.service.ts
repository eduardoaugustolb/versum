import { logger } from "@/utils/logger";
import {
  type ChapterWithContent,
  JourneyRepository,
} from "../repositories/journey.repository.ts";

export interface ChapterWithContentResponse {
  chapter: {
    id: string;
    number: number;
    totalVerses: number;
  };
  book: {
    id: string;
    order: number;
    name: string;
    slug: string;
  };
  verses: {
    id: string;
    number: number;
    text: string;
  }[];
}

export interface JourneyProgress {
  chaptersRead: number;
  chaptersRemaining: number;
  totalChapters: number;
  percentComplete: number;
  isAtEnd: boolean;
}

export interface FeedResponse {
  current: ChapterWithContentResponse | null;
  nextItems: ChapterWithContentResponse[];
  progress: JourneyProgress;
}

export interface StatusResponse {
  chaptersRead: number;
  chaptersRemaining: number;
  totalChapters: number;
  percentComplete: number;
  isAtEnd: boolean;
}

export class JourneyService {
  private readonly repository: JourneyRepository;

  constructor({ repository }: { repository?: JourneyRepository } = {}) {
    this.repository = repository ?? new JourneyRepository();
  }

  private mapChapter(result: ChapterWithContent): ChapterWithContentResponse {
    return {
      chapter: {
        id: result.chapter.id,
        number: result.chapter.number,
        totalVerses: result.chapter.totalVerses,
      },
      book: {
        id: result.book.id,
        order: result.book.order,
        name: result.book.name,
        slug: result.book.slug,
      },
      verses: result.verses.map((v) => ({
        id: v.id,
        number: v.number,
        text: v.text,
      })),
    };
  }

  private async calculateProgress(userId: string): Promise<JourneyProgress> {
    const [chaptersRead, totalChapters] = await Promise.all([
      this.repository.getReadChaptersCount(userId),
      this.repository.getTotalChapters(),
    ]);

    const chaptersRemaining = Math.max(0, totalChapters - chaptersRead);
    const percentComplete =
      totalChapters > 0 ? Math.round((chaptersRead / totalChapters) * 100) : 0;

    return {
      chaptersRead,
      chaptersRemaining,
      totalChapters,
      percentComplete,
      isAtEnd: chaptersRemaining === 0,
    };
  }

  async getFeed(userId: string, bufferSize: number = 4): Promise<FeedResponse> {
    const [current, totalChapters] = await Promise.all([
      this.repository.findNextChapterToRead(userId),
      this.repository.getTotalChapters(),
    ]);

    let nextItems: ChapterWithContent[] = [];

    if (current && bufferSize > 0) {
      nextItems = await this.repository.findChaptersAfter(
        current.chapter.id,
        bufferSize,
      );
    }

    const chaptersRead = await this.repository.getReadChaptersCount(userId);
    const chaptersRemaining = Math.max(0, totalChapters - chaptersRead);
    const isAtEnd = chaptersRemaining === 0;

    const progress: JourneyProgress = {
      chaptersRead,
      chaptersRemaining,
      totalChapters,
      percentComplete:
        totalChapters > 0
          ? Math.round((chaptersRead / totalChapters) * 100)
          : 0,
      isAtEnd,
    };

    return {
      current: current ? this.mapChapter(current) : null,
      nextItems: nextItems.map((c) => this.mapChapter(c)),
      progress,
    };
  }

  async markCurrentAsRead(userId: string): Promise<{ success: boolean }> {
    const currentChapter = await this.repository.findNextChapterToRead(userId);

    if (!currentChapter) {
      logger("debug", "[Journey] No chapter to mark as read");
      return { success: true };
    }

    logger(
      "debug",
      `[Journey] Marking chapter as read: ${currentChapter.chapter.id} (${currentChapter.book.name} ${currentChapter.chapter.number})`,
    );

    await this.repository.markChapterAsRead({
      userId,
      chapterId: currentChapter.chapter.id,
    });

    return { success: true };
  }

  async getStatus(userId: string): Promise<StatusResponse> {
    const progress = await this.calculateProgress(userId);

    return {
      chaptersRead: progress.chaptersRead,
      chaptersRemaining: progress.chaptersRemaining,
      totalChapters: progress.totalChapters,
      percentComplete: progress.percentComplete,
      isAtEnd: progress.isAtEnd,
    };
  }
}

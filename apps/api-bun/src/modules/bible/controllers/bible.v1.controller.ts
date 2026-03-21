import type { Context } from "hono";
import { BibleServiceV1 } from "../services/bible.v1.service.ts";
import { PaginationViewModel } from "../../../view-models/default/pagination.view-model.ts";
import { parsePagination } from "../../../utils/pagination/index.ts";
import { SuccessViewModel } from "../../../view-models/default/success.view-model.ts";
import { BadRequestError } from "../../../utils/app/errors/index.ts";

export class BibleControllerV1 {
  private readonly service: BibleServiceV1;

  constructor({ service }: { service?: BibleServiceV1 } = {}) {
    this.service = service ?? new BibleServiceV1();
  }

  // ------------------------
  // Books (paginated)
  // ------------------------

  getBooks = async (c: Context) => {
    const { page, limit } = parsePagination(c.req.query());

    const result = await this.service.getBooksPaginated({
      page,
      limit,
    });

    return c.json(
      SuccessViewModel.create(
        result.data,
        PaginationViewModel.create({
          page,
          limit,
          total_items: result.total,
        }),
      ),
      200,
    );
  };

  getBookByOrder = async (c: Context) => {
    const order = this.parsePositiveInt(
      c.req.param("order"),
      "Book order",
      150,
    );

    const book = await this.service.getBookByOrder({
      order,
    });

    return c.json(
      SuccessViewModel.create(
        book,
        PaginationViewModel.create({
          page: 1,
          limit: 1,
          total_items: 1,
        }),
      ),
      200,
    );
  };

  // ------------------------
  // Chapters
  // ------------------------

  getChapters = async (c: Context) => {
    const order = this.parsePositiveInt(
      c.req.param("order"),
      "Book order",
      150,
    );
    const { page, limit } = parsePagination(c.req.query());

    const result = await this.service.getChaptersPaginated({
      bookOrder: order,
      page,
      limit,
    });

    return c.json(
      SuccessViewModel.create(
        result.data,
        PaginationViewModel.create({
          page,
          limit,
          total_items: result.total,
        }),
      ),
      200,
    );
  };

  getChapter = async (c: Context) => {
    const order = this.parsePositiveInt(
      c.req.param("order"),
      "Book order",
      150,
    );
    const chapter_number = this.parsePositiveInt(
      c.req.param("number"),
      "Chapter number",
      150,
    );

    const chapter = await this.service.getChapter({
      bookOrder: order,
      chapterNumber: chapter_number,
    });

    return c.json(SuccessViewModel.create(chapter), 200);
  };

  // ------------------------
  // Verses
  // ------------------------

  getVerses = async (c: Context) => {
    const order = this.parsePositiveInt(
      c.req.param("order"),
      "Book order",
      150,
    );
    const chapter_number = this.parsePositiveInt(
      c.req.param("number"),
      "Chapter number",
      150,
    );
    const { page, limit } = parsePagination(c.req.query());

    const result = await this.service.getVersesPaginated({
      bookOrder: order,
      chapterNumber: chapter_number,
      page,
      limit,
    });

    return c.json(
      SuccessViewModel.create(
        result.data,
        PaginationViewModel.create({
          page,
          limit,
          total_items: result.total,
        }),
      ),
      200,
    );
  };

  getVerse = async (c: Context) => {
    const order = this.parsePositiveInt(
      c.req.param("order"),
      "Book order",
      150,
    );
    const chapter_number = this.parsePositiveInt(
      c.req.param("number"),
      "Chapter number",
      150,
    );
    const verse_number = this.parsePositiveInt(
      c.req.param("verse"),
      "Verse number",
      150,
    );

    const verse = await this.service.getVerse({
      bookOrder: order,
      chapterNumber: chapter_number,
      verseNumber: verse_number,
    });

    return c.json(SuccessViewModel.create(verse), 200);
  };

  // ------------------------
  // Utils
  // ------------------------

  private parsePositiveInt(value: string, field: string, max: number): number {
    if (!/^\d+$/.test(value)) {
      throw new BadRequestError(`${field} must be a valid positive integer`);
    }

    const num = Number(value);

    if (!Number.isSafeInteger(num)) {
      throw new BadRequestError(`${field} is too large`);
    }

    if (num < 1 || num > max) {
      throw new BadRequestError(`${field} is out of allowed range`);
    }

    return num;
  }
}

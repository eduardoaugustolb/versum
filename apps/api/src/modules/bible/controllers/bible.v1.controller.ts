import type { Context } from "hono";
import { BadRequestError } from "../../../utils/app/errors/index.ts";
import { parsePagination } from "../../../utils/pagination/index.ts";
import { PaginationViewModel } from "../../../view-models/default/pagination.view-model.ts";
import { SuccessViewModel } from "../../../view-models/default/success.view-model.ts";
import { BibleServiceV1 } from "../services/bible.v1.service.ts";

export class BibleControllerV1 {
  private readonly service: BibleServiceV1;

  constructor({ service }: { service?: BibleServiceV1 } = {}) {
    this.service = service ?? new BibleServiceV1();
  }

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
          totalItems: result.total,
        }),
      ),
      200,
    );
  };

  getBookByDynamicId = async (c: Context) => {
    const dynamicId = c.req.param("dynamicId");

    if (!dynamicId) {
      throw new BadRequestError("Dynamic ID is required");
    }

    const book = await this.service.getBookByDynamicId({
      dynamicId,
    });

    return c.json(
      SuccessViewModel.create(
        book,
        PaginationViewModel.create({
          page: 1,
          limit: 1,
          totalItems: 1,
        }),
      ),
      200,
    );
  };

  getChapters = async (c: Context) => {
    const dynamicId = c.req.param("dynamicId");
    const { page, limit } = parsePagination(c.req.query());

    if (!dynamicId) {
      throw new BadRequestError("Dynamic ID is required");
    }

    const result = await this.service.getChaptersPaginated({
      dynamicId,
      page,
      limit,
    });

    return c.json(
      SuccessViewModel.create(
        result.data,
        PaginationViewModel.create({
          page,
          limit,
          totalItems: result.total,
        }),
      ),
      200,
    );
  };

  getChapter = async (c: Context) => {
    const dynamicId = c.req.param("dynamicId");
    const chapterNumber = this.parsePositiveInt(
      c.req.param("number")!,
      "Chapter number",
    );

    if (!dynamicId) {
      throw new BadRequestError("Dynamic ID is required");
    }

    const chapter = await this.service.getChapter({
      dynamicId,
      chapterNumber,
    });

    return c.json(SuccessViewModel.create(chapter), 200);
  };

  getVerses = async (c: Context) => {
    const dynamicId = c.req.param("dynamicId");
    const chapterNumber = this.parsePositiveInt(
      c.req.param("number")!,
      "Chapter number",
    );
    const { page, limit } = parsePagination(c.req.query());

    if (!dynamicId) {
      throw new BadRequestError("Dynamic ID is required");
    }

    const result = await this.service.getVersesPaginated({
      dynamicId,
      chapterNumber,
      page,
      limit,
    });

    return c.json(
      SuccessViewModel.create(
        result.data,
        PaginationViewModel.create({
          page,
          limit,
          totalItems: result.total,
        }),
      ),
      200,
    );
  };

  getVerse = async (c: Context) => {
    const dynamicId = c.req.param("dynamicId");
    const chapterNumber = this.parsePositiveInt(
      c.req.param("number")!,
      "Chapter number",
    );
    const verseNumber = this.parsePositiveInt(
      c.req.param("verse")!,
      "Verse number",
    );

    if (!dynamicId) {
      throw new BadRequestError("Dynamic ID is required");
    }

    const verse = await this.service.getVerse({
      dynamicId,
      chapterNumber,
      verseNumber,
    });

    return c.json(SuccessViewModel.create(verse), 200);
  };

  private parsePositiveInt(value: string, field: string): number {
    if (!/^\d+$/.test(value)) {
      throw new BadRequestError(`${field} must be a valid positive integer`);
    }

    const num = Number(value);

    if (!Number.isSafeInteger(num)) {
      throw new BadRequestError(`${field} is too large`);
    }

    if (num < 1) {
      throw new BadRequestError(`${field} must be a positive integer`);
    }

    return num;
  }
}

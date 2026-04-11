import { createRoute, OpenAPIHono } from "@hono/zod-openapi";
import { CacheMiddleware } from "../../../middlewares/cache-reqs/middleware.ts";
import { validationErrorHook } from "../../../utils/app/errors/validation.hook.ts";
import type { BibleControllerV1 } from "../controllers/bible.v1.controller.ts";
import {
  getBookByDynamicIdResponses,
  getBooksResponses,
} from "../schemas/bible.v1.books.schema.ts";
import {
  getChapterResponses,
  getChaptersResponses,
} from "../schemas/bible.v1.chapters.schema.ts";
import {
  chapterNumberParamSchema,
  dynamicIdParamSchema,
  paginationQuerySchema,
  verseNumberParamSchema,
} from "../schemas/bible.v1.common.schema.ts";
import {
  getVerseResponses,
  getVersesResponses,
} from "../schemas/bible.v1.verses.schema.ts";

export const createBibleRoutesV1 = (controller: BibleControllerV1) => {
  const router = new OpenAPIHono({
    defaultHook: validationErrorHook,
  });

  router.use(
    "/*",
    new CacheMiddleware({
      config: {
        ttl: 300,
        cacheKeyPrefix: "cache:{method}:{path}:{query}",
      },
    }).middleware,
  );

  const getBooksRoute = createRoute({
    method: "get",
    path: "/books",
    tags: ["Bíblia"],
    summary: "Listar livros da Bíblia",
    description:
      "Retorna uma lista paginada de todos os livros da Bíblia, com informações como nome, slug, testamento e número de capítulos.",
    request: {
      query: paginationQuerySchema,
    },
    responses: getBooksResponses,
  });

  router.openapi(getBooksRoute, controller.getBooks);

  const getBookByDynamicIdRoute = createRoute({
    method: "get",
    path: "/books/{dynamicId}",
    tags: ["Bíblia"],
    summary: "Obter livro por slug ou nome",
    description:
      "Retorna um livro específico da Bíblia baseado em seu slug (ex: 'genesis') ou nome (ex: 'Gênesis').",
    request: {
      params: dynamicIdParamSchema,
    },
    responses: getBookByDynamicIdResponses,
  });

  router.openapi(getBookByDynamicIdRoute, controller.getBookByDynamicId);

  const getChaptersRoute = createRoute({
    method: "get",
    path: "/books/{dynamicId}/chapters",
    tags: ["Bíblia"],
    summary: "Listar capítulos de um livro",
    description:
      "Retorna uma lista paginada de todos os capítulos de um livro específico da Bíblia.",
    request: {
      params: dynamicIdParamSchema,
      query: paginationQuerySchema,
    },
    responses: getChaptersResponses,
  });

  router.openapi(getChaptersRoute, controller.getChapters);

  const getChapterRoute = createRoute({
    method: "get",
    path: "/books/{dynamicId}/chapters/{number}",
    tags: ["Bíblia"],
    summary: "Obter capítulo por número",
    description:
      "Retorna um capítulo específico de um livro da Bíblia, identificado pelo slug/nome do livro e número do capítulo.",
    request: {
      params: dynamicIdParamSchema.merge(chapterNumberParamSchema),
    },
    responses: getChapterResponses,
  });

  router.openapi(getChapterRoute, controller.getChapter);

  const getVersesRoute = createRoute({
    method: "get",
    path: "/books/{dynamicId}/chapters/{number}/verses",
    tags: ["Bíblia"],
    summary: "Listar versículos de um capítulo",
    description:
      "Retorna uma lista paginada de todos os versículos de um capítulo específico da Bíblia.",
    request: {
      params: dynamicIdParamSchema.merge(chapterNumberParamSchema),
      query: paginationQuerySchema,
    },
    responses: getVersesResponses,
  });

  router.openapi(getVersesRoute, controller.getVerses);

  const getVerseRoute = createRoute({
    method: "get",
    path: "/books/{dynamicId}/chapters/{number}/verses/{verse}",
    tags: ["Bíblia"],
    summary: "Obter versículo por número",
    description:
      "Retorna um versículo específico de um capítulo da Bíblia, identificado pelo slug/nome do livro, número do capítulo e número do versículo.",
    request: {
      params: dynamicIdParamSchema
        .merge(chapterNumberParamSchema)
        .merge(verseNumberParamSchema),
    },
    responses: getVerseResponses,
  });

  router.openapi(getVerseRoute, controller.getVerse);

  return router;
};

export const createBibleRoutes = createBibleRoutesV1;

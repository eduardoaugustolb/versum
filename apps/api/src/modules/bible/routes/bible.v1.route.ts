import { BibleControllerV1 } from "../controllers/bible.v1.controller.ts";
import { CacheMiddleware } from "../../../middlewares/cache.middleware.ts";
import { OpenAPIHono, createRoute } from "@hono/zod-openapi";
import {
  bookOrderParamSchema,
  chapterNumberParamSchema,
  verseNumberParamSchema,
  paginationQuerySchema,
} from "../schemas/bible.v1.common.schema.ts";
import {
  getBooksResponses,
  getBookByOrderResponses,
} from "../schemas/bible.v1.books.schema.ts";
import {
  getChaptersResponses,
  getChapterResponses,
} from "../schemas/bible.v1.chapters.schema.ts";
import {
  getVersesResponses,
  getVerseResponses,
} from "../schemas/bible.v1.verses.schema.ts";
import { validationErrorHook } from "../../../utils/app/errors/validation.hook.ts";

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

  // GET /books - List all books (paginated)
  const getBooksRoute = createRoute({
    method: "get",
    path: "/books",
    tags: ["Bíblia"],
    summary: "Listar livros da Bíblia",
    description:
      "Retorna uma lista paginada de todos os livros da Bíblia, com informações como nome, ordem, testamento e número de capítulos.",
    request: {
      query: paginationQuerySchema,
    },
    responses: getBooksResponses,
  });

  router.openapi(getBooksRoute, controller.getBooks);

  // GET /books/:order - Get book by order
  const getBookByOrderRoute = createRoute({
    method: "get",
    path: "/books/{order}",
    tags: ["Bíblia"],
    summary: "Obter livro por ordem",
    description:
      "Retorna um livro específico da Bíblia baseado em sua ordem numérica (1-150).",
    request: {
      params: bookOrderParamSchema,
    },
    responses: getBookByOrderResponses,
  });

  router.openapi(getBookByOrderRoute, controller.getBookByOrder);

  // GET /books/:order/chapters - List chapters of a book
  const getChaptersRoute = createRoute({
    method: "get",
    path: "/books/{order}/chapters",
    tags: ["Bíblia"],
    summary: "Listar capítulos de um livro",
    description:
      "Retorna uma lista paginada de todos os capítulos de um livro específico da Bíblia.",
    request: {
      params: bookOrderParamSchema,
      query: paginationQuerySchema,
    },
    responses: getChaptersResponses,
  });

  router.openapi(getChaptersRoute, controller.getChapters);
  // GET /books/:order/chapters/:number - Get chapter by number
  const getChapterRoute = createRoute({
    method: "get",
    path: "/books/{order}/chapters/{number}",
    tags: ["Bíblia"],
    summary: "Obter capítulo por número",
    description:
      "Retorna um capítulo específico de um livro da Bíblia, identificado pela ordem do livro e número do capítulo.",
    request: {
      params: bookOrderParamSchema.merge(chapterNumberParamSchema),
    },
    responses: getChapterResponses,
  });

  router.openapi(getChapterRoute, controller.getChapter);

  // GET /books/:order/chapters/:number/verses - List verses of a chapter
  const getVersesRoute = createRoute({
    method: "get",
    path: "/books/{order}/chapters/{number}/verses",
    tags: ["Bíblia"],
    summary: "Listar versículos de um capítulo",
    description:
      "Retorna uma lista paginada de todos os versículos de um capítulo específico da Bíblia.",
    request: {
      params: bookOrderParamSchema.merge(chapterNumberParamSchema),
      query: paginationQuerySchema,
    },
    responses: getVersesResponses,
  });

  router.openapi(getVersesRoute, controller.getVerses);

  // GET /books/:order/chapters/:number/verses/:verse - Get verse by number
  const getVerseRoute = createRoute({
    method: "get",
    path: "/books/{order}/chapters/{number}/verses/{verse}",
    tags: ["Bíblia"],
    summary: "Obter versículo por número",
    description:
      "Retorna um versículo específico de um capítulo da Bíblia, identificado pela ordem do livro, número do capítulo e número do versículo.",
    request: {
      params: bookOrderParamSchema
        .merge(chapterNumberParamSchema)
        .merge(verseNumberParamSchema),
    },
    responses: getVerseResponses,
  });

  router.openapi(getVerseRoute, controller.getVerse);

  return router;
};

export const createBibleRoutes = createBibleRoutesV1;

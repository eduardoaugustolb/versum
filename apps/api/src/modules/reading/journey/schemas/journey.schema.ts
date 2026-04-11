import { z } from "@hono/zod-openapi";
import { createSuccessResponseSchema } from "../../../bible/schemas/bible.v1.common.schema.ts";

export const chapterResponseSchema = z.object({
  id: z.string().uuid().openapi({ description: "ID do capítulo" }),
  number: z
    .number()
    .int()
    .positive()
    .openapi({ description: "Número do capítulo" }),
  totalVerses: z
    .number()
    .int()
    .positive()
    .openapi({ description: "Total de versículos" }),
});

export const bookResponseSchema = z.object({
  id: z.string().uuid().openapi({ description: "ID do livro" }),
  order: z
    .number()
    .int()
    .positive()
    .openapi({ description: "Ordem canônica do livro (1-73)" }),
  name: z.string().openapi({ description: "Nome do livro" }),
  slug: z.string().openapi({ description: "Slug do livro" }),
});

export const verseResponseSchema = z.object({
  id: z.string().uuid().openapi({ description: "ID do versículo" }),
  number: z
    .number()
    .int()
    .positive()
    .openapi({ description: "Número do versículo" }),
  text: z.string().openapi({ description: "Texto do versículo" }),
});

export const chapterWithContentSchema = z.object({
  chapter: chapterResponseSchema,
  book: bookResponseSchema,
  verses: z.array(verseResponseSchema),
});

export const progressResponseSchema = z.object({
  chaptersRead: z
    .number()
    .int()
    .nonnegative()
    .openapi({ description: "Capítulos lidos" }),
  chaptersRemaining: z
    .number()
    .int()
    .nonnegative()
    .openapi({ description: "Capítulos restantes" }),
  totalChapters: z
    .number()
    .int()
    .positive()
    .openapi({ description: "Total de capítulos" }),
  percentComplete: z
    .number()
    .min(0)
    .max(100)
    .openapi({ description: "Porcentagem concluída" }),
  isAtEnd: z.boolean().openapi({ description: "Se chegou ao fim da Bíblia" }),
});

export const feedResponseDataSchema = z.object({
  current: chapterWithContentSchema
    .nullable()
    .openapi({ description: "Próximo capítulo a ler" }),
  nextItems: z
    .array(chapterWithContentSchema)
    .openapi({ description: "Capítulos para pre-fetch" }),
  progress: progressResponseSchema,
});

export const statusResponseDataSchema = z.object({
  chaptersRead: z
    .number()
    .int()
    .nonnegative()
    .openapi({ description: "Capítulos lidos" }),
  chaptersRemaining: z
    .number()
    .int()
    .nonnegative()
    .openapi({ description: "Capítulos restantes" }),
  totalChapters: z
    .number()
    .int()
    .positive()
    .openapi({ description: "Total de capítulos" }),
  percentComplete: z
    .number()
    .min(0)
    .max(100)
    .openapi({ description: "Porcentagem concluída" }),
  isAtEnd: z.boolean().openapi({ description: "Se chegou ao fim da Bíblia" }),
});

export const feedQuerySchema = z.object({
  "buffer-size": z.coerce
    .number()
    .int()
    .min(0)
    .max(4)
    .default(4)
    .openapi({
      param: { name: "buffer-size", in: "query", required: false },
      description: "Quantidade de itens no pre-fetch (0-4, padrão: 4)",
      example: 4,
    }),
});

export const feedResponseSchema = createSuccessResponseSchema(
  "JourneyFeedResponse",
  feedResponseDataSchema,
);

export const nextProgressResponseSchema = createSuccessResponseSchema(
  "JourneyNextProgressResponse",
  z.object({ success: z.boolean() }),
);

export const statusResponseSchema = createSuccessResponseSchema(
  "JourneyStatusResponse",
  statusResponseDataSchema,
);

export const journeyStatsResponseSchema = createSuccessResponseSchema(
  "JourneyStatsResponse",
  z.object({
    chaptersRead: z
      .number()
      .int()
      .nonnegative()
      .openapi({ description: "Total de capítulos lidos" }),
  }),
);

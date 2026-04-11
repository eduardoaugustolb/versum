import { z } from "@hono/zod-openapi";
import { createSuccessResponseSchema } from "../../../bible/schemas/bible.v1.common.schema.ts";

export const markVersesRequestSchema = z.object({
  verseIds: z.array(z.string().uuid()).openapi({
    example: ["550e8400-e29b-41d4-a716-446655440000"],
    description: "IDs dos versículos lidos",
  }),
});

export const verseWithContextResponseSchema = z.object({
  id: z.string().uuid().openapi({ description: "ID do versículo" }),
  number: z
    .number()
    .int()
    .positive()
    .openapi({ description: "Número do versículo" }),
  text: z.string().openapi({ description: "Texto do versículo" }),
  chapter: z.object({
    id: z.string().uuid().openapi({ description: "ID do capítulo" }),
    number: z
      .number()
      .int()
      .positive()
      .openapi({ description: "Número do capítulo" }),
  }),
  book: z.object({
    id: z.string().uuid().openapi({ description: "ID do livro" }),
    order: z
      .number()
      .int()
      .positive()
      .openapi({ description: "Ordem canônica do livro (1-73)" }),
    name: z.string().openapi({ description: "Nome do livro" }),
    slug: z.string().openapi({ description: "Slug do livro" }),
  }),
});

export const nextVersesResponseDataSchema = z.array(
  verseWithContextResponseSchema,
);

export const discoveryStatsResponseDataSchema = z.object({
  versesRead: z.number().int().nonnegative().openapi({
    description: "Total de versículos lidos",
  }),
});

export const nextVersesResponseSchema = createSuccessResponseSchema(
  "DiscoveryNextVersesResponse",
  nextVersesResponseDataSchema,
);

export const discoveryStatsResponseSchema = createSuccessResponseSchema(
  "DiscoveryStatsResponse",
  discoveryStatsResponseDataSchema,
);

export const markVersesResponseSchema = createSuccessResponseSchema(
  "DiscoveryMarkVersesResponse",
  z.object({ success: z.boolean() }),
);

export type MarkVersesRequest = z.infer<typeof markVersesRequestSchema>;

import { z } from "@hono/zod-openapi";

export class BibleCommonSchemasV1 {
  static readonly testamentEnum = z
    .enum(["OLD", "NEW"])
    .describe("Testamento bíblico");

  static readonly bookSchema = z
    .object({
      id: z.uuid().describe("ID único do livro"),
      name: z.string().max(100).describe("Nome do livro"),
      order: z
        .number()
        .int()
        .min(1)
        .max(150)
        .describe("Ordem do livro na Bíblia"),
      testament: BibleCommonSchemasV1.testamentEnum.describe(
        "Testamento ao qual o livro pertence",
      ),
      total_chapters: z
        .number()
        .int()
        .positive()
        .describe("Número total de capítulos"),
    })
    .openapi("Book", {
      description: "Livro da Bíblia",
    });

  static readonly chapterSchema = z
    .object({
      id: z.uuid().describe("ID único do capítulo"),
      book_id: z.uuid().describe("ID do livro ao qual o capítulo pertence"),
      number: z.number().int().positive().describe("Número do capítulo"),
      total_verses: z
        .number()
        .int()
        .positive()
        .describe("Número total de versículos"),
    })
    .openapi("Chapter", {
      description: "Capítulo de um livro da Bíblia",
    });

  static readonly verseSchema = z
    .object({
      id: z.uuid().describe("ID único do versículo"),
      chapter_id: z
        .uuid()
        .describe("ID do capítulo ao qual o versículo pertence"),
      number: z.number().int().positive().describe("Número do versículo"),
      text: z.string().describe("Texto do versículo"),
    })
    .openapi("Verse", {
      description: "Versículo de um capítulo da Bíblia",
    });

  static readonly paginationViewModelSchema = z
    .object({
      current_page: z.number().int().positive().describe("Página atual"),
      total_pages: z
        .number()
        .int()
        .positive()
        .describe("Número total de páginas"),
      total_items: z
        .number()
        .int()
        .nonnegative()
        .describe("Número total de itens"),
      items_per_page: z
        .number()
        .int()
        .positive()
        .describe("Número de itens por página"),
      has_next_page: z.boolean().describe("Indica se existe próxima página"),
      has_prev_page: z.boolean().describe("Indica se existe página anterior"),
    })
    .openapi("PaginationViewModel", {
      description: "Informações de paginação",
    });

  static createSuccessResponseSchema<T extends z.ZodType>(
    name: string,
    dataSchema: T,
    includePagination = false,
  ) {
    const schema = z.object({
      success: z
        .boolean()
        .default(true)
        .describe("Indica se a requisição foi bem-sucedida"),
      message: z.string().optional().describe("Mensagem opcional de contexto"),
      data: dataSchema.optional().describe("Dados da resposta"),
      pagination: includePagination
        ? BibleCommonSchemasV1.paginationViewModelSchema
            .optional()
            .describe("Informações de paginação")
        : BibleCommonSchemasV1.paginationViewModelSchema
            .optional()
            .describe("Informações de paginação"),
    });

    return schema.openapi(name, {
      description: `Resposta de sucesso para ${name}`,
    });
  }

  static readonly bookOrderParamSchema = z.object({
    order: z
      .string()
      .regex(/^\d+$/, "Ordem do livro deve ser um inteiro positivo válido")
      .openapi({
        param: {
          name: "order",
          in: "path",
          required: true,
        },
        example: "1",
        description: "Ordem do livro na Bíblia (1-150)",
      }),
  });

  static readonly chapterNumberParamSchema = z.object({
    number: z
      .string()
      .regex(/^\d+$/, "Número do capítulo deve ser um inteiro positivo válido")
      .openapi({
        param: {
          name: "number",
          in: "path",
          required: true,
        },
        example: "1",
        description: "Número do capítulo",
      }),
  });

  static readonly verseNumberParamSchema = z.object({
    verse: z
      .string()
      .regex(/^\d+$/, "Número do versículo deve ser um inteiro positivo válido")
      .openapi({
        param: {
          name: "verse",
          in: "path",
          required: true,
        },
        example: "1",
        description: "Número do versículo",
      }),
  });

  static readonly paginationQuerySchema = z.object({
    page: z
      .string()
      .optional()
      .default("1")
      .openapi({
        param: {
          name: "page",
          in: "query",
          required: false,
        },
        example: "1",
        description: "Número da página (padrão: 1)",
      }),
    limit: z
      .string()
      .optional()
      .default("10")
      .openapi({
        param: {
          name: "limit",
          in: "query",
          required: false,
        },
        example: "10",
        description: "Limite de itens por página (padrão: 10)",
      }),
  });
}

export const testamentEnum = BibleCommonSchemasV1.testamentEnum;
export const bookSchema = BibleCommonSchemasV1.bookSchema;
export const chapterSchema = BibleCommonSchemasV1.chapterSchema;
export const verseSchema = BibleCommonSchemasV1.verseSchema;
export const paginationViewModelSchema =
  BibleCommonSchemasV1.paginationViewModelSchema;
export const createSuccessResponseSchema =
  BibleCommonSchemasV1.createSuccessResponseSchema;
export const bookOrderParamSchema = BibleCommonSchemasV1.bookOrderParamSchema;
export const chapterNumberParamSchema =
  BibleCommonSchemasV1.chapterNumberParamSchema;
export const verseNumberParamSchema =
  BibleCommonSchemasV1.verseNumberParamSchema;
export const paginationQuerySchema = BibleCommonSchemasV1.paginationQuerySchema;

export type Book = z.infer<typeof BibleCommonSchemasV1.bookSchema>;
export type Chapter = z.infer<typeof BibleCommonSchemasV1.chapterSchema>;
export type Verse = z.infer<typeof BibleCommonSchemasV1.verseSchema>;
export type Testament = z.infer<typeof BibleCommonSchemasV1.testamentEnum>;

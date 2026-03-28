import { z } from "@hono/zod-openapi";

export class BibleCommonSchemasV1 {
  static readonly testamentEnum = z
    .enum(["OLD", "NEW"])
    .describe("Testamento bíblico");

  static readonly bookSchema = z
    .object({
      id: z.uuid().describe("ID único do livro"),
      name: z.string().max(100).describe("Nome do livro"),
      slug: z.string().max(10).describe("Slug do livro"),
      niceName: z.string().max(100).describe("Nome amigável do livro"),
      testament: BibleCommonSchemasV1.testamentEnum.describe(
        "Testamento ao qual o livro pertence",
      ),
      totalChapters: z
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
      bookId: z.uuid().describe("ID do livro ao qual o capítulo pertence"),
      number: z.number().int().positive().describe("Número do capítulo"),
      totalVerses: z
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
      chapterId: z
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
      currentPage: z.number().int().positive().describe("Página atual"),
      totalPages: z
        .number()
        .int()
        .positive()
        .describe("Número total de páginas"),
      totalItems: z
        .number()
        .int()
        .nonnegative()
        .describe("Número total de itens"),
      itemsPerPage: z
        .number()
        .int()
        .positive()
        .describe("Número de itens por página"),
      hasNextPage: z.boolean().describe("Indica se existe próxima página"),
      hasPrevPage: z.boolean().describe("Indica se existe página anterior"),
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

  static readonly dynamicIdParamSchema = z.object({
    dynamicId: z
      .string()
      .min(1)
      .openapi({
        param: {
          name: "dynamicId",
          in: "path",
          required: true,
        },
        example: "genesis",
        description: "Slug ou nome do livro (ex: 'genesis' ou 'Gênesis')",
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
export const dynamicIdParamSchema = BibleCommonSchemasV1.dynamicIdParamSchema;
export const chapterNumberParamSchema =
  BibleCommonSchemasV1.chapterNumberParamSchema;
export const verseNumberParamSchema =
  BibleCommonSchemasV1.verseNumberParamSchema;
export const paginationQuerySchema = BibleCommonSchemasV1.paginationQuerySchema;

export type Book = z.infer<typeof BibleCommonSchemasV1.bookSchema>;
export type Chapter = z.infer<typeof BibleCommonSchemasV1.chapterSchema>;
export type Verse = z.infer<typeof BibleCommonSchemasV1.verseSchema>;
export type Testament = z.infer<typeof BibleCommonSchemasV1.testamentEnum>;

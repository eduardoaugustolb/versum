import { z } from "@hono/zod-openapi";
import { createErrorResponses } from "../../../utils/app/errors/openapi.ts";
import { BibleCommonSchemasV1 } from "./bible.v1.common.schema.ts";

export class BooksSchemasV1 {
  static readonly getBooksResponseSchema =
    BibleCommonSchemasV1.createSuccessResponseSchema(
      "GetBooksResponse",
      z.array(BibleCommonSchemasV1.bookSchema),
      true,
    );

  static readonly getBooksResponses = {
    200: {
      content: {
        "application/json": {
          schema: BooksSchemasV1.getBooksResponseSchema,
        },
      },
      description: "Lista de livros da Bíblia retornada com sucesso",
    },
    ...createErrorResponses([400, 429, 500]),
  };

  static readonly getBookByDynamicIdResponseSchema =
    BibleCommonSchemasV1.createSuccessResponseSchema(
      "GetBookByDynamicIdResponse",
      BibleCommonSchemasV1.bookSchema,
      true,
    );

  static readonly getBookByDynamicIdResponses = {
    200: {
      content: {
        "application/json": {
          schema: BooksSchemasV1.getBookByDynamicIdResponseSchema,
        },
      },
      description: "Livro encontrado e retornado com sucesso",
    },
    ...createErrorResponses([400, 404, 429, 500]),
  };
}

export const getBooksResponseSchema = BooksSchemasV1.getBooksResponseSchema;
export const getBooksResponses = BooksSchemasV1.getBooksResponses;
export const getBookByDynamicIdResponseSchema =
  BooksSchemasV1.getBookByDynamicIdResponseSchema;
export const getBookByDynamicIdResponses =
  BooksSchemasV1.getBookByDynamicIdResponses;

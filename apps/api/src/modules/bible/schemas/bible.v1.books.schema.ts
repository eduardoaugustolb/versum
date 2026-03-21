import { z } from "@hono/zod-openapi";
import { BibleCommonSchemasV1 } from "./bible.v1.common.schema.ts";
import { createErrorResponses } from "../../../utils/app/errors/openapi.ts";

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
    ...createErrorResponses([400, 500]),
  };

  static readonly getBookByOrderResponseSchema =
    BibleCommonSchemasV1.createSuccessResponseSchema(
      "GetBookByOrderResponse",
      BibleCommonSchemasV1.bookSchema,
      true,
    );

  static readonly getBookByOrderResponses = {
    200: {
      content: {
        "application/json": {
          schema: BooksSchemasV1.getBookByOrderResponseSchema,
        },
      },
      description: "Livro encontrado e retornado com sucesso",
    },
    ...createErrorResponses([400, 404, 500]),
  };
}

export const getBooksResponseSchema = BooksSchemasV1.getBooksResponseSchema;
export const getBooksResponses = BooksSchemasV1.getBooksResponses;
export const getBookByOrderResponseSchema =
  BooksSchemasV1.getBookByOrderResponseSchema;
export const getBookByOrderResponses = BooksSchemasV1.getBookByOrderResponses;

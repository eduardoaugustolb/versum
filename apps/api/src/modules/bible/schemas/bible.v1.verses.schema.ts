import { z } from "@hono/zod-openapi";
import { BibleCommonSchemasV1 } from "./bible.v1.common.schema.ts";
import { createErrorResponses } from "../../../utils/app/errors/openapi.ts";

export class VersesSchemasV1 {
  static readonly getVersesResponseSchema =
    BibleCommonSchemasV1.createSuccessResponseSchema(
      "GetVersesResponse",
      z.array(BibleCommonSchemasV1.verseSchema),
      true,
    );

  static readonly getVersesResponses = {
    200: {
      content: {
        "application/json": {
          schema: VersesSchemasV1.getVersesResponseSchema,
        },
      },
      description: "Lista de versículos do capítulo retornada com sucesso",
    },
    ...createErrorResponses([400, 404, 500]),
  };

  static readonly getVerseResponseSchema =
    BibleCommonSchemasV1.createSuccessResponseSchema(
      "GetVerseResponse",
      BibleCommonSchemasV1.verseSchema,
      false,
    );

  static readonly getVerseResponses = {
    200: {
      content: {
        "application/json": {
          schema: VersesSchemasV1.getVerseResponseSchema,
        },
      },
      description: "Versículo encontrado e retornado com sucesso",
    },
    ...createErrorResponses([400, 404, 500]),
  };
}

export const getVersesResponseSchema = VersesSchemasV1.getVersesResponseSchema;
export const getVersesResponses = VersesSchemasV1.getVersesResponses;
export const getVerseResponseSchema = VersesSchemasV1.getVerseResponseSchema;
export const getVerseResponses = VersesSchemasV1.getVerseResponses;

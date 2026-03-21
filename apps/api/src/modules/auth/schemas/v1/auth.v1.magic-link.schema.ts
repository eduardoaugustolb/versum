import { AuthCommonSchemasV1 } from "./auth.v1.common.schema.ts";
import { createErrorResponses } from "../../../../utils/app/errors/openapi.ts";

export class MagicLinkSchemasV1 {
  static readonly createAndSendMagicLinkResponseSchema =
    AuthCommonSchemasV1.messageSchema.openapi("SendMagicLinkResponse", {
      description: "Solicitação de magic link aceita",
    });

  static readonly createAndSendMagicLinkResponses = {
    200: {
      content: {
        "application/json": {
          schema: MagicLinkSchemasV1.createAndSendMagicLinkResponseSchema,
        },
      },
      description: "Magic link enviado",
    },
    ...createErrorResponses([400, 500]),
  };

  static readonly authenticateWithMagicLinkResponseSchema =
    AuthCommonSchemasV1.messageSchema.openapi(
      "AuthenticateWithMagicLinkResponse",
      {
        description: "Sessão criada",
      },
    );

  static readonly authenticateWithMagicLinkResponses = {
    200: {
      content: {
        "application/json": {
          schema: MagicLinkSchemasV1.authenticateWithMagicLinkResponseSchema,
        },
      },
      description: "Login realizado com sucesso",
    },
    ...createErrorResponses([400, 401, 500]),
  };
}

export const createAndSendMagicLinkResponseSchema =
  MagicLinkSchemasV1.createAndSendMagicLinkResponseSchema;
export const createAndSendMagicLinkResponses =
  MagicLinkSchemasV1.createAndSendMagicLinkResponses;
export const authenticateWithMagicLinkResponseSchema =
  MagicLinkSchemasV1.authenticateWithMagicLinkResponseSchema;
export const authenticateWithMagicLinkResponses =
  MagicLinkSchemasV1.authenticateWithMagicLinkResponses;

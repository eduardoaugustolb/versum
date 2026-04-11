import { createErrorResponses } from "../../../../utils/app/errors/openapi.ts";
import { AuthCommonSchemasV1 } from "./auth.v1.common.schema.ts";

export class LogoutSchemasV1 {
  static readonly logoutResponseSchema =
    AuthCommonSchemasV1.messageSchema.openapi("LogoutResponse", {
      description: "Sessão encerrada com sucesso",
    });

  static readonly logoutResponses = {
    200: {
      content: {
        "application/json": {
          schema: LogoutSchemasV1.logoutResponseSchema,
        },
      },
      description: "Logout realizado com sucesso",
    },
    ...createErrorResponses([400, 429, 500]),
  };
}

export const logoutResponseSchema = LogoutSchemasV1.logoutResponseSchema;
export const logoutResponses = LogoutSchemasV1.logoutResponses;

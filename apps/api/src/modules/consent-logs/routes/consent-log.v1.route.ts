import { createRoute, OpenAPIHono } from "@hono/zod-openapi";
import { AuthMiddleware } from "@/middlewares/auth.middleware.ts";
import { createErrorResponses } from "../../../utils/app/errors/openapi.ts";
import { validationErrorHook } from "../../../utils/app/errors/validation.hook.ts";
import type { ConsentLogControllerV1 } from "../controllers/consent-log.v1.controller.ts";
import {
  consentHistoryResponseSchema,
  recordConsentBodySchema,
} from "../schemas/v1/consent-log.v1.common.schema.ts";

export const createConsentLogRoutesV1 = (controller: ConsentLogControllerV1) => {
  const router = new OpenAPIHono({
    defaultHook: validationErrorHook,
  });

  const authMiddleware = new AuthMiddleware();

  const recordConsentRoute = createRoute({
    method: "post",
    path: "/",
    tags: ["Consent Logs"],
    summary: "Registrar consentimentos",
    description: "Registra um ou mais consentimentos do usuário autenticado.",
    security: [{ cookieAuth: [] }],
    request: {
      body: {
        content: {
          "application/json": {
            schema: recordConsentBodySchema,
          },
        },
      },
    },
    responses: {
      201: {
        content: {
          "application/json": {
            schema: consentHistoryResponseSchema,
          },
        },
        description: "Consentimentos registrados com sucesso",
      },
      ...createErrorResponses([400, 401, 429, 500]),
    },
  });

  const getConsentHistoryRoute = createRoute({
    method: "get",
    path: "/",
    tags: ["Consent Logs"],
    summary: "Obter histórico de consentimentos",
    description: "Retorna todos os registros de consentimento do usuário autenticado.",
    security: [{ cookieAuth: [] }],
    responses: {
      200: {
        content: {
          "application/json": {
            schema: consentHistoryResponseSchema,
          },
        },
        description: "Histórico de consentimentos retornado com sucesso",
      },
      ...createErrorResponses([401, 429, 500]),
    },
  });

  router.use("/*", authMiddleware.validateSession);

  router.openapi(recordConsentRoute, controller.recordConsent);
  router.openapi(getConsentHistoryRoute, controller.getConsentHistory);

  return router;
};

import { createRoute, OpenAPIHono } from "@hono/zod-openapi";
import { AuthMiddleware } from "@/middlewares/auth.middleware.ts";
import { createErrorResponses } from "../../../utils/app/errors/openapi.ts";
import { validationErrorHook } from "../../../utils/app/errors/validation.hook.ts";
import type { UsersControllerV1 } from "../controllers/users.v1.controller.ts";
import {
  deleteAuthenticatedUserResponseSchema,
  exportUserDataResponseSchema,
  getAuthenticatedUserResponseSchema,
  updateAuthenticatedUserBodySchema,
  updateAuthenticatedUserResponseSchema,
} from "../schemas/v1/users.v1.common.schema.ts";

export const createUsersRoutesV1 = (controller: UsersControllerV1) => {
  const router = new OpenAPIHono({
    defaultHook: validationErrorHook,
  });

  const authMiddleware = new AuthMiddleware();

  const getMeRoute = createRoute({
    method: "get",
    path: "/@me",
    tags: ["Users"],
    summary: "Obter usuário autenticado",
      description: "Retorna os dados do usuário autenticado.",
      security: [{ cookieAuth: [] }],
    responses: {
      200: {
        content: {
          "application/json": {
            schema: getAuthenticatedUserResponseSchema,
          },
        },
        description: "Usuário autenticado retornado com sucesso",
      },
      ...createErrorResponses([401, 404, 429, 500]),
    },
  });

  const updateMeRoute = createRoute({
    method: "patch",
    path: "/@me",
    tags: ["Users"],
    summary: "Atualizar usuário autenticado",
      description: "Atualiza os dados do usuário autenticado.",
      security: [{ cookieAuth: [] }],
    request: {
      body: {
        content: {
          "application/json": {
            schema: updateAuthenticatedUserBodySchema,
          },
        },
      },
    },
    responses: {
      200: {
        content: {
          "application/json": {
            schema: updateAuthenticatedUserResponseSchema,
          },
        },
        description: "Usuário autenticado atualizado com sucesso",
      },
      ...createErrorResponses([400, 401, 404, 409, 429, 500]),
    },
  });

  const deleteMeRoute = createRoute({
    method: "delete",
    path: "/@me",
    tags: ["Users"],
    summary: "Deletar usuário autenticado",
    description: "Remove a conta do usuário e todos os dados associados (direito de eliminação LGPD Art. 18, VI).",
    security: [{ cookieAuth: [] }],
    responses: {
      204: {
        description: "Usuário deletado com sucesso",
      },
      ...createErrorResponses([401, 404, 429, 500]),
    },
  });

  const exportMeRoute = createRoute({
    method: "get",
    path: "/@me/export",
    tags: ["Users"],
    summary: "Exportar dados do usuário autenticado",
    description: "Exporta todos os dados pessoais do usuário (direito de portabilidade LGPD Art. 18, II e V).",
    security: [{ cookieAuth: [] }],
    responses: {
      200: {
        content: {
          "application/json": {
            schema: exportUserDataResponseSchema,
          },
        },
        description: "Dados do usuário exportados com sucesso",
      },
      ...createErrorResponses([401, 404, 429, 500]),
    },
  });

  router.use("/@me", authMiddleware.validateSession);

  router.openapi(getMeRoute, controller.getAuthenticatedUser);
  router.openapi(updateMeRoute, controller.updateAuthenticatedUser);
  router.openapi(deleteMeRoute, controller.deleteAuthenticatedUser);
  router.openapi(exportMeRoute, controller.exportUserData);

  return router;
};

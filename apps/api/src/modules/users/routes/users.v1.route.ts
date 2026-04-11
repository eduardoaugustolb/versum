import { createRoute, OpenAPIHono } from "@hono/zod-openapi";
import { AuthMiddleware } from "@/middlewares/auth.middleware.ts";
import { createErrorResponses } from "../../../utils/app/errors/openapi.ts";
import { validationErrorHook } from "../../../utils/app/errors/validation.hook.ts";
import type { UsersControllerV1 } from "../controllers/users.v1.controller.ts";
import {
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

  router.use("/@me", authMiddleware.validateSession);

  router.openapi(getMeRoute, controller.getAutheticatedUser);
  router.openapi(updateMeRoute, controller.updateAuthenticatedUser);

  return router;
};

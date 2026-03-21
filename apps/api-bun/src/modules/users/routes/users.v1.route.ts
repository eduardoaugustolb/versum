import { createRoute, OpenAPIHono } from "@hono/zod-openapi";
import { UsersControllerV1 } from "../controllers/users.v1.controller.ts";
import {
  getAuthenticatedUserResponseSchema,
  updateAuthenticatedUserBodySchema,
  updateAuthenticatedUserResponseSchema,
  usernameParamSchema,
} from "../schemas/v1/users.v1.common.schema.ts";
import { validationErrorHook } from "../../../utils/app/errors/validation.hook.ts";
import { createErrorResponses } from "../../../utils/app/errors/openapi.ts";
import { AuthMiddleware } from "../../../middlewares/auth.middleware.ts";

export class UsersV1Routes {
  router = new OpenAPIHono({
    defaultHook: validationErrorHook,
  });

  constructor({ controller }: { controller?: UsersControllerV1 } = {}) {
    this.setupRoutes(controller ?? new UsersControllerV1());
  }

  private setupRoutes(controller: UsersControllerV1) {
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
        ...createErrorResponses([401, 404, 500]),
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
        ...createErrorResponses([400, 401, 404, 500]),
      },
    });

    const getByUsernameRoute = this.getByUsernameRoute();

    this.router.use("/@me", async (c, next) => {
      await authMiddleware.validateSession(c);
      await next();
    });

    this.router.openapi(getMeRoute, controller.getAutheticatedUser);
    this.router.openapi(updateMeRoute, controller.updateAuthenticatedUser);
    this.router.openapi(getByUsernameRoute, controller.getUserByUsername);
    this.router.get("/:username", controller.getUserByUsername);
  }

  private getByUsernameRoute() {
    return createRoute({
      method: "get",
      path: "/@{username}",
      tags: ["Users"],
      summary: "Obter usuário por username",
      description: "Retorna os dados públicos de um usuário pelo username.",
      request: {
        params: usernameParamSchema,
      },
      responses: {
        200: {
          content: {
            "application/json": {
              schema: getAuthenticatedUserResponseSchema,
            },
          },
          description: "Usuário retornado com sucesso",
        },
        ...createErrorResponses([404, 500]),
      },
    });
  }
}

export const createUsersRoutesV1 = (controller?: UsersControllerV1) =>
  new UsersV1Routes({ controller }).router;

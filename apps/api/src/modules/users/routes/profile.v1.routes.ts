import { createRoute, OpenAPIHono } from "@hono/zod-openapi";
import { AuthMiddleware } from "@/middlewares/auth.middleware.ts";
import { createErrorResponses } from "../../../utils/app/errors/openapi.ts";
import { validationErrorHook } from "../../../utils/app/errors/validation.hook.ts";
import type { ProfileControllerV1 } from "../controllers/profile.v1.controller.ts";
import {
  createProfileBodySchema,
  createProfileResponseSchema,
  getAuthenticatedProfileResponseSchema,
  profileSchema,
  updateAuthenticatedProfileBodySchema,
  updateAuthenticatedProfileResponseSchema,
  usernameParamSchema,
} from "../schemas/v1/profiles.v1.common.schema.ts";

export const createProfileRoutesV1 = (controller: ProfileControllerV1) => {
  const router = new OpenAPIHono({
    defaultHook: validationErrorHook,
  });

  const authMiddleware = new AuthMiddleware();

  const createProfileRoute = createRoute({
    method: "post",
    path: "/",
    tags: ["Profiles"],
    summary: "Criar perfil do usuário autenticado",
    description: "Cria um novo perfil para o usuário autenticado.",
    request: {
      body: {
        content: {
          "application/json": {
            schema: createProfileBodySchema,
          },
        },
      },
    },
    responses: {
      201: {
        content: {
          "application/json": {
            schema: createProfileResponseSchema,
          },
        },
        description: "Perfil criado com sucesso",
      },
      ...createErrorResponses([400, 401, 409, 429, 500]),
    },
  });

  const getMeRoute = createRoute({
    method: "get",
    path: "/@me",
    tags: ["Profiles"],
    summary: "Obter perfil do usuário autenticado",
    description: "Retorna os dados do perfil do usuário autenticado.",
    responses: {
      200: {
        content: {
          "application/json": {
            schema: getAuthenticatedProfileResponseSchema,
          },
        },
        description: "Perfil retornado com sucesso",
      },
      ...createErrorResponses([401, 404, 429, 500]),
    },
  });

  const updateMeRoute = createRoute({
    method: "patch",
    path: "/@me",
    tags: ["Profiles"],
    summary: "Atualizar perfil do usuário autenticado",
    description: "Atualiza os dados do perfil do usuário autenticado.",
    request: {
      body: {
        content: {
          "application/json": {
            schema: updateAuthenticatedProfileBodySchema,
          },
        },
      },
    },
    responses: {
      200: {
        content: {
          "application/json": {
            schema: updateAuthenticatedProfileResponseSchema,
          },
        },
        description: "Perfil atualizado com sucesso",
      },
      ...createErrorResponses([400, 401, 404, 409, 429, 500]),
    },
  });

  const getByUsernameRoute = createRoute({
    method: "get",
    path: "/{username}",
    tags: ["Profiles"],
    summary: "Obter perfil por username",
    description: "Retorna os dados públicos de um perfil pelo username.",
    request: {
      params: usernameParamSchema,
    },
    responses: {
      200: {
        content: {
          "application/json": {
            schema: profileSchema,
          },
        },
        description: "Perfil retornado com sucesso",
      },
      ...createErrorResponses([404, 429, 500]),
    },
  });

  router.use("/@me", authMiddleware.validateSession);

  router.openapi(createProfileRoute, controller.createProfile);
  router.openapi(getMeRoute, controller.getAuthenticatedProfile);
  router.openapi(updateMeRoute, controller.updateAuthenticatedProfile);
  router.openapi(getByUsernameRoute, controller.getProfileByUsername);

  return router;
};

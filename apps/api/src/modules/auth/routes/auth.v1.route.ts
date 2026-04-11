import { createRoute, OpenAPIHono } from "@hono/zod-openapi";
import { MagicLinkRateLimiter } from "@/middlewares/rate-limiter/middleware.ts";
import { validationErrorHook } from "../../../utils/app/errors/validation.hook.ts";
import type { AuthControllerV1 } from "../controllers/auth.v1.controller.ts";
import {
  magicLinkTokenQuerySchema,
  sendMagicLinkBodySchema,
} from "../schemas/v1/auth.v1.common.schema.ts";
import { logoutResponses } from "../schemas/v1/auth.v1.logout.schema.ts";
import {
  authenticateWithMagicLinkResponses,
  createAndSendMagicLinkResponses,
} from "../schemas/v1/auth.v1.magic-link.schema.ts";

export const createAuthRoutesV1 = (controller: AuthControllerV1) => {
  const router = new OpenAPIHono({
    defaultHook: validationErrorHook,
  });

  const magicLinkRateLimiter = new MagicLinkRateLimiter();
  router.use("/magic-link", magicLinkRateLimiter.middleware);

  const sendMagicLinkRouteV1 = createRoute({
    method: "post",
    path: "/magic-link",
    tags: ["Auth"],
    summary: "Enviar magic link",
    description: "Envia um link mágico para autenticação via e-mail.",
    request: {
      body: {
        content: {
          "application/json": {
            schema: sendMagicLinkBodySchema,
          },
        },
      },
    },
    responses: createAndSendMagicLinkResponses,
  });

  router.openapi(sendMagicLinkRouteV1, controller.createAndSendMagicLink);

  const authenticateWithMagicLinkRouteV1 = createRoute({
    method: "get",
    path: "/magic-link",
    tags: ["Auth"],
    summary: "Autenticar com magic link",
    description: "Cria uma sessão a partir de um token de magic link.",
    request: {
      query: magicLinkTokenQuerySchema,
    },
    responses: authenticateWithMagicLinkResponses,
  });

  router.openapi(
    authenticateWithMagicLinkRouteV1,
    controller.authenticateWithMagicLink,
  );

  const logoutRouteV1 = createRoute({
    method: "post",
    path: "/logout",
    tags: ["Auth"],
    summary: "Encerrar sessão",
    description: "Encerra a sessão atual baseada no cookie de sessão.",
    responses: logoutResponses,
  });

  router.openapi(logoutRouteV1, controller.logout);

  return router;
};

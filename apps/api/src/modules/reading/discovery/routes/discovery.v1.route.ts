import { createRoute, OpenAPIHono, z } from "@hono/zod-openapi";
import { AuthMiddleware } from "@/middlewares/auth.middleware.ts";
import { createErrorResponses } from "@/utils/app/errors/openapi.ts";
import { validationErrorHook } from "@/utils/app/errors/validation.hook.ts";
import { logger } from "@versum/logger";
import { DiscoveryControllerV1 } from "../controllers/discovery.v1.controller.ts";
import {
  discoveryStatsResponseSchema,
  markVersesRequestSchema,
  markVersesResponseSchema,
  nextVersesResponseSchema,
} from "../schemas/discovery.v1.schema.ts";

export class DiscoveryRoutesV1 {
  router = new OpenAPIHono({
    defaultHook: validationErrorHook,
  });

  constructor({ controller }: { controller?: DiscoveryControllerV1 } = {}) {
    this.setupRoutes(controller ?? new DiscoveryControllerV1());
  }

  private setupRoutes(controller: DiscoveryControllerV1) {
    const authMiddleware = new AuthMiddleware();
    logger("debug", "[Discovery] routes setup");

    this.router.use("/*", authMiddleware.validateSession.bind(authMiddleware));

    const getNextVersesRoute = createRoute({
      method: "get",
      path: "/next",
      tags: ["Discovery"],
      summary: "Versículos do capítulo",
      description:
        "Retorna todos os versículos de um capítulo para leitura no modo Discovery.",
      security: [{ cookieAuth: [] }],
      request: {
        query: z.object({
          chapterId: z.string().uuid().optional().openapi({
            example: "550e8400-e29b-41d4-a716-446655440000",
            description: "ID do capítulo (opcional)",
          }),
        }),
      },
      responses: {
        200: {
          content: {
            "application/json": {
              schema: nextVersesResponseSchema,
            },
          },
          description: "Versículos retornados",
        },
        ...createErrorResponses([401, 500]),
      },
    });

    this.router.openapi(getNextVersesRoute, controller.getNextVerses);

    const markVersesRoute = createRoute({
      method: "post",
      path: "/",
      tags: ["Discovery"],
      summary: "Marcar versículos como lidos",
      description: "Registra que o usuário leu versículos específicos.",
      security: [{ cookieAuth: [] }],
      request: {
        body: {
          content: {
            "application/json": {
              schema: markVersesRequestSchema,
            },
          },
        },
      },
      responses: {
        200: {
          content: {
            "application/json": {
              schema: markVersesResponseSchema,
            },
          },
          description: "Versículos marcados como lidos",
        },
        ...createErrorResponses([400, 401, 500]),
      },
    });

    this.router.openapi(markVersesRoute, controller.markVersesAsRead);

    const getStatsRoute = createRoute({
      method: "get",
      path: "/stats",
      tags: ["Discovery"],
      summary: "Estatísticas do Discovery",
      description: "Retorna estatísticas de leitura no modo Discovery.",
      security: [{ cookieAuth: [] }],
      responses: {
        200: {
          content: {
            "application/json": {
              schema: discoveryStatsResponseSchema,
            },
          },
          description: "Estatísticas retornadas",
        },
        ...createErrorResponses([401, 500]),
      },
    });

    this.router.openapi(getStatsRoute, controller.getStats);
  }
}

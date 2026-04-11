import { createRoute, OpenAPIHono } from "@hono/zod-openapi";
import { AuthMiddleware } from "../../../../middlewares/auth.middleware.ts";
import { createErrorResponses } from "../../../../utils/app/errors/openapi.ts";
import { validationErrorHook } from "../../../../utils/app/errors/validation.hook.ts";
import { JourneyController } from "../controllers/journey.controller.ts";
import {
  feedQuerySchema,
  feedResponseSchema,
  nextProgressResponseSchema,
  statusResponseSchema,
} from "../schemas/journey.schema.ts";

export class JourneyRoutes {
  router = new OpenAPIHono({
    defaultHook: validationErrorHook,
  });

  constructor({ controller }: { controller?: JourneyController } = {}) {
    this.setupRoutes(controller ?? new JourneyController());
  }

  private setupRoutes(controller: JourneyController) {
    const authMiddleware = new AuthMiddleware();

    this.router.use("/*", authMiddleware.validateSession.bind(authMiddleware));

    const feedRoute = createRoute({
      method: "get",
      path: "/feed",
      tags: ["Journey"],
      summary: "Feed de leitura",
      description:
        "Retorna o capítulo atual e os próximos itens para pre-fetch.",
      request: {
        query: feedQuerySchema,
      },
      responses: {
        200: {
          content: {
            "application/json": {
              schema: feedResponseSchema,
            },
          },
          description: "Feed retornado com sucesso",
        },
        ...createErrorResponses([401, 500]),
      },
    });

    this.router.openapi(feedRoute, controller.getFeed);

    const nextRoute = createRoute({
      method: "post",
      path: "/next",
      tags: ["Journey"],
      summary: "Avançar progresso",
      description: "Salva o capítulo atual como lido e avança para o próximo.",
      responses: {
        200: {
          content: {
            "application/json": {
              schema: nextProgressResponseSchema,
            },
          },
          description: "Progresso salvo com sucesso",
        },
        ...createErrorResponses([401, 500]),
      },
    });

    this.router.openapi(nextRoute, controller.markCurrentAsRead);

    const statusRoute = createRoute({
      method: "get",
      path: "/status",
      tags: ["Journey"],
      summary: "Status da jornada",
      description: "Retorna o status atual do progresso de leitura.",
      responses: {
        200: {
          content: {
            "application/json": {
              schema: statusResponseSchema,
            },
          },
          description: "Status retornado com sucesso",
        },
        ...createErrorResponses([401, 500]),
      },
    });

    this.router.openapi(statusRoute, controller.getStatus);
  }
}

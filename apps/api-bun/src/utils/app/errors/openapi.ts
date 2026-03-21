import { z } from "@hono/zod-openapi";

export const apiErrorResponseSchema = z
  .object({
    success: z.literal(false).describe("Indica que ocorreu um erro"),
    error: z.string().describe("Mensagem de erro"),
    code: z.string().optional().describe("Código do erro"),
  })
  .openapi("ApiErrorResponse", {
    description: "Resposta padrão de erro da API",
  });

type ErrorStatusCode = 400 | 401 | 403 | 404 | 409 | 422 | 429 | 500;

type ErrorResponseObject = {
  content: {
    "application/json": {
      schema: typeof apiErrorResponseSchema;
    };
  };
  description: string;
};

export const createErrorResponse = (
  description: string,
): ErrorResponseObject => ({
  content: {
    "application/json": {
      schema: apiErrorResponseSchema,
    },
  },
  description,
});

export const DEFAULT_ERROR_RESPONSE_DESCRIPTIONS: Record<
  ErrorStatusCode,
  string
> = {
  400: "Requisição inválida",
  401: "Não autorizado",
  403: "Acesso negado",
  404: "Recurso não encontrado",
  409: "Conflito de estado",
  422: "Entidade não processável",
  429: "Muitas requisições",
  500: "Erro interno do servidor",
};

export const createErrorResponses = (
  statusesOrDescriptions:
    | ErrorStatusCode[]
    | Partial<Record<ErrorStatusCode, string>>,
  overrides: Partial<Record<ErrorStatusCode, string>> = {},
): Record<number, ErrorResponseObject> => {
  const descriptions = Array.isArray(statusesOrDescriptions)
    ? Object.fromEntries(
        statusesOrDescriptions.map((status) => [
          status,
          DEFAULT_ERROR_RESPONSE_DESCRIPTIONS[status],
        ]),
      )
    : statusesOrDescriptions;

  const merged = {
    ...descriptions,
    ...overrides,
  };

  const entries = Object.entries(merged).map(([status, description]) => [
    Number(status),
    createErrorResponse(description ?? "Erro"),
  ]);

  return Object.fromEntries(entries);
};

import { z } from "@hono/zod-openapi";

const paginationViewModelSchema = z
  .object({
    currentPage: z.number().int().positive().describe("Página atual"),
    totalPages: z.number().int().positive().describe("Número total de páginas"),
    totalItems: z.number().int().nonnegative().describe("Número total de itens"),
    itemsPerPage: z.number().int().positive().describe("Número de itens por página"),
    hasNextPage: z.boolean().describe("Indica se existe próxima página"),
    hasPrevPage: z.boolean().describe("Indica se existe página anterior"),
  })
  .openapi("PaginationViewModel", {
    description: "Informações de paginação",
  });

export function createSuccessResponseSchema<T extends z.ZodType>(
  name: string,
  dataSchema: T,
  includePagination = false,
) {
  const schema: Record<string, z.ZodType> = {
    success: z.boolean().default(true).describe("Indica se a requisição foi bem-sucedida"),
    message: z.string().optional().describe("Mensagem opcional de contexto"),
    data: dataSchema.optional().describe("Dados da resposta"),
  };

  if (includePagination) {
    schema.pagination = paginationViewModelSchema
      .optional()
      .describe("Informações de paginação");
  }

  return z.object(schema).openapi(name, {
    description: `Resposta de sucesso para ${name}`,
  });
}

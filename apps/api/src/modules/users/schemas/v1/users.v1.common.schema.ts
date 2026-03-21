import { z } from "@hono/zod-openapi";

export class UsersCommonSchemasV1 {
  static readonly userSchema = z
    .object({
      id: z.uuid().describe("ID único do usuário"),
      username: z.string().min(1).max(50).describe("Nome de usuário único"),
      name: z.string().min(1).max(100).describe("Nome de exibição do usuário"),
      email: z.email().describe("E-mail do usuário"),
      bio: z.string().max(500).nullable().describe("Biografia do usuário"),
      picture_url: z
        .string()
        .max(500)
        .nullable()
        .describe("URL da foto de perfil"),
      created_at: z.iso.datetime().describe("Data/hora de criação da conta"),
    })
    .openapi("User", {
      description: "Representação pública do usuário",
    });

  static readonly getAuthenticatedUserResponseSchema = z
    .object({
      user: UsersCommonSchemasV1.userSchema,
    })
    .openapi("GetAuthenticatedUserResponse", {
      description: "Usuário autenticado retornado com sucesso",
    });

  static readonly updateAuthenticatedUserBodySchema = z
    .object({
      username: z
        .string()
        .min(1)
        .max(50)
        .optional()
        .describe("Novo nome de usuário"),
      name: z
        .string()
        .min(1)
        .max(100)
        .optional()
        .describe("Novo nome de exibição"),
      bio: z.string().max(500).nullable().optional().describe("Nova biografia"),
      picture_url: z
        .string()
        .max(500)
        .nullable()
        .optional()
        .describe("Nova URL da foto de perfil"),
    })
    .openapi("UpdateAuthenticatedUserBody", {
      description: "Payload para atualizar os dados do usuário autenticado",
    });

  static readonly updateAuthenticatedUserResponseSchema = z
    .object({
      user: UsersCommonSchemasV1.userSchema,
    })
    .openapi("UpdateAuthenticatedUserResponse", {
      description: "Usuário autenticado atualizado com sucesso",
    });

  static readonly usernameParamSchema = z.object({
    username: z
      .string()
      .min(1)
      .max(50)
      .openapi({
        param: {
          name: "username",
          in: "path",
          required: true,
        },
        description: "Nome de usuário",
        example: "joao",
      }),
  });
}

export const userSchema = UsersCommonSchemasV1.userSchema;
export const getAuthenticatedUserResponseSchema =
  UsersCommonSchemasV1.getAuthenticatedUserResponseSchema;
export const updateAuthenticatedUserBodySchema =
  UsersCommonSchemasV1.updateAuthenticatedUserBodySchema;
export const updateAuthenticatedUserResponseSchema =
  UsersCommonSchemasV1.updateAuthenticatedUserResponseSchema;
export const usernameParamSchema = UsersCommonSchemasV1.usernameParamSchema;

export type UserSchema = z.infer<typeof UsersCommonSchemasV1.userSchema>;
export type UpdateAuthenticatedUserBody = z.infer<
  typeof UsersCommonSchemasV1.updateAuthenticatedUserBodySchema
>;
export type UsernameParam = z.infer<
  typeof UsersCommonSchemasV1.usernameParamSchema
>;

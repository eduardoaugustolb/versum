import { z } from "@hono/zod-openapi";

export class UsersCommonSchemasV1 {
  static readonly emailSchema = z
    .string()
    .email("Invalid email format")
    .max(255, "Email must not exceed 255 characters")
    .toLowerCase()
    .openapi({
      description: "E-mail do usuário",
      example: "john@example.com",
    });

  static readonly userSchema = z
    .object({
      email: this.emailSchema,
    })
    .openapi("User", {
      description: "Dados privados do usuário",
    });

  static readonly getAuthenticatedUserResponseSchema = z
    .object({
      user: this.userSchema,
      onboardingIsCompleted: z.boolean(),
    })
    .openapi("GetAuthenticatedUserResponse", {
      description: "Usuário autenticado retornado com sucesso",
    });

  static readonly updateAuthenticatedUserBodySchema = z
    .object({
      email: this.emailSchema,
    })
    .openapi("UpdateAuthenticatedUserBody", {
      description: "Payload para atualizar os dados do usuário autenticado",
    });

  static readonly updateAuthenticatedUserResponseSchema = z
    .object({
      user: this.userSchema,
    })
    .openapi("UpdateAuthenticatedUserResponse", {
      description: "Usuário autenticado atualizado com sucesso",
    });

  static readonly deleteAuthenticatedUserResponseSchema = z
    .null()
    .openapi("DeleteAuthenticatedUserResponse", {
      description: "Usuário deletado com sucesso (sem conteúdo)",
    });
}

export const userSchema = UsersCommonSchemasV1.userSchema;
export const getAuthenticatedUserResponseSchema =
  UsersCommonSchemasV1.getAuthenticatedUserResponseSchema;
export const updateAuthenticatedUserBodySchema =
  UsersCommonSchemasV1.updateAuthenticatedUserBodySchema;
export const updateAuthenticatedUserResponseSchema =
  UsersCommonSchemasV1.updateAuthenticatedUserResponseSchema;
export const deleteAuthenticatedUserResponseSchema =
  UsersCommonSchemasV1.deleteAuthenticatedUserResponseSchema;

export type UserSchema = z.infer<typeof UsersCommonSchemasV1.userSchema>;
export type UpdateAuthenticatedUserBody = z.infer<
  typeof UsersCommonSchemasV1.updateAuthenticatedUserBodySchema
>;

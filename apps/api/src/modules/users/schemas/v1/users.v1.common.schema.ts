import { z } from "@hono/zod-openapi";
import { createSuccessResponseSchema } from "@/utils/app/schemas/success-response.ts";

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

  static readonly authenticatedUserDataSchema = z
    .object({
      user: this.userSchema,
      onboardingIsCompleted: z.boolean(),
    })
    .openapi("AuthenticatedUserData", {
      description: "Dados do usuário autenticado",
    });

  static readonly getAuthenticatedUserResponseSchema =
    createSuccessResponseSchema(
      "GetAuthenticatedUserResponse",
      this.authenticatedUserDataSchema,
    );

  static readonly updateAuthenticatedUserBodySchema = z
    .object({
      email: this.emailSchema,
    })
    .openapi("UpdateAuthenticatedUserBody", {
      description: "Payload para atualizar os dados do usuário autenticado",
    });

  static readonly updatedUserDataSchema = z
    .object({
      user: this.userSchema,
    })
    .openapi("UpdatedUserData", {
      description: "Dados do usuário atualizado",
    });

  static readonly updateAuthenticatedUserResponseSchema =
    createSuccessResponseSchema(
      "UpdateAuthenticatedUserResponse",
      this.updatedUserDataSchema,
    );

  static readonly deleteAuthenticatedUserResponseSchema = z
    .null()
    .openapi("DeleteAuthenticatedUserResponse", {
      description: "Usuário deletado com sucesso (sem conteúdo)",
    });

  static readonly consentLogSchema = z
    .object({
      id: z.string().openapi({ description: "ID do registro de consentimento" }),
      userId: z.string().openapi({ description: "ID do usuário" }),
      purpose: z.string().openapi({ description: "Finalidade do consentimento" }),
      granted: z.boolean().openapi({ description: "Se o consentimento foi concedido" }),
      ip: z.string().openapi({ description: "Endereço IP" }),
      userAgent: z.string().openapi({ description: "User-Agent" }),
      createdAt: z.string().datetime().openapi({ description: "Data do registro" }),
    })
    .openapi("ExportConsentLog", {
      description: "Registro de consentimento",
    });

  static readonly exportDataObjectSchema = z
    .object({
      exportedAt: z.string().datetime().openapi({
        description: "Data e hora da exportação",
      }),
      user: z.object({
        email: z.string().email(),
        createdAt: z.string().datetime(),
      }).openapi("ExportUser", {
        description: "Dados do usuário",
      }),
      profile: z.object({
        username: z.string(),
        name: z.string(),
        bio: z.string().nullable(),
        pictureUrl: z.string().nullable(),
      }).nullable().openapi("ExportProfile", {
        description: "Perfil do usuário",
      }),
      sessions: z.array(z.object({
        createdAt: z.string().datetime(),
        ip: z.string(),
        userAgent: z.string(),
        expiresAt: z.string().datetime(),
      })).openapi({
        description: "Sessões do usuário",
      }),
      readingHistory: z.object({
        journey: z.array(z.object({
          chapterId: z.string(),
          readAt: z.string().datetime(),
        })),
        discovery: z.array(z.object({
          verseId: z.string(),
          readAt: z.string().datetime(),
        })),
      }).openapi("ExportReadingHistory", {
        description: "Histórico de leitura",
      }),
      annotations: z.array(z.object({
        verseId: z.string(),
        selectedVerseId: z.string(),
        annotation: z.string().nullable(),
        isPublic: z.boolean(),
        createdAt: z.string().datetime(),
      })).openapi({
        description: "Anotações do usuário",
      }),
      likes: z.array(z.object({
        verseId: z.string(),
        createdAt: z.string().datetime(),
      })).openapi({
        description: "Favoritos do usuário",
      }),
      consentLogs: z.array(this.consentLogSchema).openapi({
        description: "Histórico de consentimento",
      }),
    })
    .openapi("ExportUserData", {
      description: "Dados exportados do usuário",
    });

  static readonly exportUserDataResponseSchema =
    createSuccessResponseSchema(
      "ExportUserDataResponse",
      this.exportDataObjectSchema,
    );
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
export const exportUserDataResponseSchema =
  UsersCommonSchemasV1.exportUserDataResponseSchema;

export type UserSchema = z.infer<typeof UsersCommonSchemasV1.userSchema>;
export type UpdateAuthenticatedUserBody = z.infer<
  typeof UsersCommonSchemasV1.updateAuthenticatedUserBodySchema
>;
export type ExportUserDataResponse = z.infer<
  typeof UsersCommonSchemasV1.exportUserDataResponseSchema
>;

import { z } from "@hono/zod-openapi";

const USERNAME_REGEX = /^[a-zA-Z0-9_]+$/;

export class ProfilesCommonSchemasV1 {
  static readonly usernameSchema = z
    .string()
    .min(3, "Username must be at least 3 characters")
    .max(50, "Username must not exceed 50 characters")
    .regex(
      USERNAME_REGEX,
      "Username can only contain letters, numbers, and underscores",
    )
    .toLowerCase()
    .openapi({
      description: "Username único",
      example: "johndoe",
    });

  static readonly nameSchema = z
    .string()
    .min(1, "Name is required")
    .max(100, "Name must not exceed 100 characters")
    .openapi({
      description: "Nome de exibição",
      example: "John Doe",
    });

  static readonly bioSchema = z
    .string()
    .max(500, "Bio must not exceed 500 characters")
    .nullable()
    .optional()
    .openapi({
      description: "Biografia do usuário",
      example: "Software developer and open source enthusiast",
    });

  static readonly pictureUrlSchema = z
    .string()
    .max(500, "Picture URL must not exceed 500 characters")
    .url("Picture URL must be a valid URL")
    .startsWith("https://", "Picture URL must use HTTPS")
    .nullable()
    .optional()
    .openapi({
      description: "URL da foto de perfil",
      example: "https://example.com/avatar.jpg",
    });

  static readonly profileSchema = z
    .object({
      username: this.usernameSchema,
      name: this.nameSchema,
      bio: this.bioSchema,
      pictureUrl: this.pictureUrlSchema,
    })
    .openapi("Profile", {
      description: "Representação pública do perfil",
    });

  static readonly createProfileBodySchema = z
    .object({
      username: this.usernameSchema,
      name: this.nameSchema,
      bio: this.bioSchema,
      pictureUrl: this.pictureUrlSchema,
    })
    .openapi("CreateProfileBody", {
      description: "Payload para criar o perfil do usuário",
    });

  static readonly createProfileResponseSchema = z
    .object({
      profile: this.profileSchema,
    })
    .openapi("CreateProfileResponse", {
      description: "Perfil criado com sucesso",
    });

  static readonly getAuthenticatedProfileResponseSchema = z
    .object({
      profile: this.profileSchema,
    })
    .openapi("GetAuthenticatedProfileResponse", {
      description: "Perfil do usuário autenticado",
    });

  static readonly updateAuthenticatedProfileBodySchema = z
    .object({
      username: this.usernameSchema.optional(),
      name: this.nameSchema.optional(),
      bio: this.bioSchema,
      pictureUrl: this.pictureUrlSchema,
    })
    .openapi("UpdateAuthenticatedProfileBody", {
      description: "Payload para atualizar o perfil do usuário autenticado",
    });

  static readonly updateAuthenticatedProfileResponseSchema = z
    .object({
      profile: this.profileSchema,
    })
    .openapi("UpdateAuthenticatedProfileResponse", {
      description: "Perfil atualizado com sucesso",
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
        description: "Username",
        example: "johndoe",
      }),
  });
}

export const profileSchema = ProfilesCommonSchemasV1.profileSchema;
export const createProfileBodySchema =
  ProfilesCommonSchemasV1.createProfileBodySchema;
export const createProfileResponseSchema =
  ProfilesCommonSchemasV1.createProfileResponseSchema;
export const getAuthenticatedProfileResponseSchema =
  ProfilesCommonSchemasV1.getAuthenticatedProfileResponseSchema;
export const updateAuthenticatedProfileBodySchema =
  ProfilesCommonSchemasV1.updateAuthenticatedProfileBodySchema;
export const updateAuthenticatedProfileResponseSchema =
  ProfilesCommonSchemasV1.updateAuthenticatedProfileResponseSchema;
export const usernameParamSchema = ProfilesCommonSchemasV1.usernameParamSchema;

export type ProfileSchema = z.infer<
  typeof ProfilesCommonSchemasV1.profileSchema
>;
export type CreateProfileBody = z.infer<
  typeof ProfilesCommonSchemasV1.createProfileBodySchema
>;
export type UpdateAuthenticatedProfileBody = z.infer<
  typeof ProfilesCommonSchemasV1.updateAuthenticatedProfileBodySchema
>;
export type UsernameParam = z.infer<
  typeof ProfilesCommonSchemasV1.usernameParamSchema
>;

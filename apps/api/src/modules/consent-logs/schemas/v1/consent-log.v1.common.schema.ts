import { z } from "@hono/zod-openapi";

export const CONSENT_PURPOSES = [
  "profile_content",
  "annotations",
  "likes",
  "terms",
] as const;

export type ConsentPurpose = (typeof CONSENT_PURPOSES)[number];

export class ConsentLogCommonSchemasV1 {
  static readonly consentPurposeSchema = z
    .enum(CONSENT_PURPOSES)
    .openapi({
      description: "Finalidade do consentimento",
      example: "profile_content",
    });

  static readonly consentItemSchema = z
    .object({
      purpose: this.consentPurposeSchema,
      granted: z.boolean().openapi({
        description: "true = concedido, false = revogado",
        example: true,
      }),
    })
    .openapi("ConsentItem", {
      description: "Item de consentimento",
    });

  static readonly recordConsentBodySchema = z
    .object({
      consents: z
        .array(this.consentItemSchema)
        .min(1, "Pelo menos um consentimento deve ser informado")
        .openapi({
          description: "Lista de consentimentos",
        }),
    })
    .openapi("RecordConsentBody", {
      description: "Payload para registrar consentimentos",
    });

  static readonly consentLogSchema = z
    .object({
      id: z.string().uuid(),
      userId: z.string().uuid(),
      purpose: z.string(),
      granted: z.boolean(),
      ip: z.string(),
      userAgent: z.string(),
      createdAt: z.date(),
    })
    .openapi("ConsentLog", {
      description: "Registro de consentimento",
    });

  static readonly consentHistoryResponseSchema = z
    .object({
      consents: z.array(this.consentLogSchema),
    })
    .openapi("ConsentHistoryResponse", {
      description: "Histórico de consentimentos do usuário",
    });
}

export const consentPurposeSchema = ConsentLogCommonSchemasV1.consentPurposeSchema;
export const consentItemSchema = ConsentLogCommonSchemasV1.consentItemSchema;
export const recordConsentBodySchema = ConsentLogCommonSchemasV1.recordConsentBodySchema;
export const consentLogSchema = ConsentLogCommonSchemasV1.consentLogSchema;
export const consentHistoryResponseSchema = ConsentLogCommonSchemasV1.consentHistoryResponseSchema;

export type ConsentItem = z.infer<typeof ConsentLogCommonSchemasV1.consentItemSchema>;
export type RecordConsentBody = z.infer<typeof ConsentLogCommonSchemasV1.recordConsentBodySchema>;

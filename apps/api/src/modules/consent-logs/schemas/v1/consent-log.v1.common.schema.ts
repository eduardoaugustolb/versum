import { z } from "@hono/zod-openapi";
import { createSuccessResponseSchema } from "@/utils/app/schemas/success-response.ts";

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

  static readonly consentHistoryDataSchema = z
    .object({
      consents: z.array(this.consentLogSchema),
    })
    .openapi("ConsentHistoryData", {
      description: "Dados do histórico de consentimentos",
    });

  static readonly consentHistoryResponseSchema =
    createSuccessResponseSchema(
      "ConsentHistoryResponse",
      this.consentHistoryDataSchema,
    );
}

export const consentPurposeSchema = ConsentLogCommonSchemasV1.consentPurposeSchema;
export const consentItemSchema = ConsentLogCommonSchemasV1.consentItemSchema;
export const recordConsentBodySchema = ConsentLogCommonSchemasV1.recordConsentBodySchema;
export const consentLogSchema = ConsentLogCommonSchemasV1.consentLogSchema;
export const consentHistoryDataSchema = ConsentLogCommonSchemasV1.consentHistoryDataSchema;
export const consentHistoryResponseSchema = ConsentLogCommonSchemasV1.consentHistoryResponseSchema;

export type ConsentItem = z.infer<typeof ConsentLogCommonSchemasV1.consentItemSchema>;
export type RecordConsentBody = z.infer<typeof ConsentLogCommonSchemasV1.recordConsentBodySchema>;

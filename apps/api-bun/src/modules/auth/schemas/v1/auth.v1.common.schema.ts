import { z } from "@hono/zod-openapi";

export class AuthCommonSchemasV1 {
  static readonly messageSchema = z
    .object({
      message: z.string().describe("Mensagem legível para humanos"),
    })
    .openapi("AuthMessageResponse", {
      description: "Resposta simples com mensagem",
    });

  static readonly sendMagicLinkBodySchema = z
    .object({
      email: z.email().describe("E-mail do usuário"),
    })
    .openapi("SendMagicLinkBody", {
      description: "Payload para solicitar um magic link",
    });

  static readonly magicLinkTokenQuerySchema = z.object({
    token: z
      .string()
      .min(1)
      .openapi({
        param: {
          name: "token",
          in: "query",
          required: true,
        },
        description: "Token do magic link (public_id.token)",
        example: "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee.xxxxxxxxxxxxxx",
      }),
  });
}

export const messageSchema = AuthCommonSchemasV1.messageSchema;
export const sendMagicLinkBodySchema =
  AuthCommonSchemasV1.sendMagicLinkBodySchema;
export const magicLinkTokenQuerySchema =
  AuthCommonSchemasV1.magicLinkTokenQuerySchema;

export type SendMagicLinkBody = z.infer<
  typeof AuthCommonSchemasV1.sendMagicLinkBodySchema
>;
export type MagicLinkTokenQuery = z.infer<
  typeof AuthCommonSchemasV1.magicLinkTokenQuerySchema
>;

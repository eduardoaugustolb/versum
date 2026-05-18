import { z } from "zod";

export const loginFormSchema = z.object({
  email: z.string().email("Digite um email válido"),
});

export type LoginFormValues = z.infer<typeof loginFormSchema>;

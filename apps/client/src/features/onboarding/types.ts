import z from "zod";

// SCHEMA

export const onboardingFormSchema = z.object({
  name: z
    .string()
    .min(2, "Nome deve ter pelo menos 2 caracteres")
    .max(60, "Nome muito longo"),
  username: z
    .string()
    .min(3, "Username deve ter pelo menos 3 caracteres")
    .max(30, "Username muito longo")
    .regex(
      /^[a-z0-9_]+$/,
      "Apenas letras minúsculas, números e underscores"
    ),
  bio: z.string().min(1, { message: "A bio é obrigatória" }),
});

// STEP DEFINITIONS

export type OnboardingValues = z.infer<typeof onboardingFormSchema>;

export type StepKind = "in" | "form" | "out" | "error";

export interface BaseStep {
  id: string;
  kind: StepKind;
}

export interface InStep extends BaseStep {
  kind: "in";
  label: string;
}

export interface FormStep extends BaseStep {
  kind: "form";
  field: keyof OnboardingValues;
  title: string;
  description: string;
  placeholder: string;
  inputType: "text";
}

export interface OutStep extends BaseStep {
  kind: "out";
  getLabel: (values: OnboardingValues) => string;
}

export interface ErrorStep extends BaseStep {
  kind: "error";
  title: string;
}

export type OnboardingStep = InStep | FormStep | OutStep | ErrorStep;

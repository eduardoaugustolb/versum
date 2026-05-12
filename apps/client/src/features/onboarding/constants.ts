import { FormStep, OnboardingStep } from "./types";

export const STEPS: OnboardingStep[] = [
  {
    id: "in",
    kind: "in",
    label: "Estamos quase lá...",
  },
  {
      id: "name",
      kind: "form",
      field: "name",
      title: "Qual é o seu nome?",
      description: "Como você gostaria de ser chamado?",
      placeholder: "Seu nome completo",
      inputType: "text",
    },
    {
      id: "username",
      kind: "form",
      field: "username",
      title: "Escolha um username",
      description: "Letras minúsculas, números e underscores.",
      placeholder: "ex: joao_silva",
      inputType: "text",
  },
  {
    id: "bio",
    kind: "form",
    field: "bio",
    title: "Sobre você",
    description: "Quem você é e o que você faz.",
    placeholder: "Digite uma breve descrição",
    inputType: "text",
  },
  {
    id: "out",
    kind: "out",
    getLabel: (v) => `Seja bem-vindo, ${v.name || "Fulano"}!`,
  },
] as const;


export const TOTAL_STEPS = STEPS.length;
export const FORM_STEPS = STEPS.filter((s): s is FormStep => s.kind === "form");

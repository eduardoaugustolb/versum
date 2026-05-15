import type { ConsentOption, ConsentStep, FormStep, OnboardingStep } from "./types";

export const CONSENT_OPTIONS: ConsentOption[] = [
  {
    purpose: "profile_content",
    label: "Conteúdo do perfil",
    description: "Permitir armazenar nome, username, bio e foto do perfil.",
    required: true,
  },
  {
    purpose: "annotations",
    label: "Anotações",
    description: "Permitir salvar minhas anotações e marcadores nos versículos.",
    required: false,
  },
  {
    purpose: "likes",
    label: "Favoritos",
    description: "Permitir salvar meus versículos favoritos.",
    required: false,
  },
  {
    purpose: "terms",
    label: "Termos e Privacidade",
    description: "Aceito os Termos de Uso e a Política de Privacidade.",
    required: true,
  },
];

export const STEPS: OnboardingStep[] = [
  {
    id: "in",
    kind: "in",
    label: "Estamos quase lá...",
  },
  {
    id: "consent",
    kind: "consent",
    options: CONSENT_OPTIONS,
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
];

export const TOTAL_STEPS = STEPS.length;
export const FORM_STEPS = STEPS.filter((s): s is FormStep => s.kind === "form");

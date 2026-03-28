import prompts from "prompts";
import { TRUNCATABLE_TABLES, TRUNCATE_PRESETS } from "./truncate.constants.ts";

export const truncateModeMenu = async () =>
  await prompts({
    type: "select",
    name: "mode",
    message: "Como deseja escolher o TRUNCATE?",
    choices: [
      {
        title: "Preset (grupos prontos)",
        value: "preset",
      },
      {
        title: "Tabelas individuais (multiseleção)",
        value: "custom",
      },
      { title: "Voltar", value: "back" },
    ],
  });

export const truncatePresetMenu = async () =>
  await prompts({
    type: "select",
    name: "preset",
    message: "Escolha o preset:",
    choices: [
      ...Object.entries(TRUNCATE_PRESETS).map(([key, v]) => ({
        title: v.title,
        description: v.description,
        value: key,
      })),
      { title: "Voltar", value: "back" },
    ],
  });

export const truncateTablesMultiselect = async () =>
  await prompts({
    type: "multiselect",
    name: "tables",
    message: "Selecione as tabelas (confirme com espaço, Enter para enviar):",
    choices: TRUNCATABLE_TABLES.map((t) => ({
      title: `${t.name} — ${t.label}`,
      value: t.name,
    })),
    hint: "- Espaço para marcar. Enter confirma.",
    validate: (value: string[]) =>
      value.length > 0 ? true : "Selecione ao menos uma tabela.",
  });

export const truncateConfirmPrompt = async () =>
  await prompts({
    type: "text",
    name: "token",
    message:
      'Para confirmar o TRUNCATE irreversível, digite exatamente: TRUNCAR (ou "voltar" para cancelar)',
    validate: (value) => {
      const v = value.trim().toLowerCase();
      if (v === "voltar" || v === "back") return true;
      if (value === "TRUNCAR") return true;
      return 'Digite TRUNCAR para confirmar, ou "voltar".';
    },
  });

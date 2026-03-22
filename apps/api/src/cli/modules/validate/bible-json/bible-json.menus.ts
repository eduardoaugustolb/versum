import { createConfirmActionPrompt } from "@/cli/menus/cofirm-action.menu";
import prompts from "prompts";

export const confirmValidationPromptMenu = async () => await createConfirmActionPrompt("Deseja prosseguir com a validação?");
export const saveResultJsonFilesMenu = async ({ errors}: { errors?: string[] }) => await prompts({
  type: "multiselect",
  name: "save_json_files",
  message: "Selecione os arquivos JSON que deseja salvar:",
  choices: [
    {
      value: "parsed_bible_data",
      title: "Dados parseados",
      selected: true
    },
    errors && errors.length > 0 ? {
      value: "parse_bible_errors",
      title: "Erros"
    } : undefined
  ].filter(t => t != undefined)
})

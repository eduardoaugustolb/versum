import prompts from "prompts";
import { findCliOutputFiles } from "./helpers/find-cli-output-files";

export const seedBibleJsonPathPrompt = async () => {
  const cwd = process.cwd();
  const files = await findCliOutputFiles(cwd);

  return await prompts({
    type: "autocomplete",
    name: "bible_json_path",
    message:
      'Selecione ou digite o caminho do arquivo JSON (escolha "Sair" para sair, "Listar arquivos" para listar arquivos ou "Voltar" para voltar):',
    validate: async (value) => {
      const lower = value.toLowerCase();
      if (lower === "exit" || lower === "back" || lower === "list-files") {
        return true;
      }
      if (!value.trim()) {
        return "Por favor, digite um caminho válido.";
      }
      if (!value.endsWith(".json")) {
        return "O arquivo precisa ser um .json";
      }
      const file = Bun.file(value);
      const exists = await file.exists();
      if (!exists) {
        return "Esse arquivo não existe.";
      }
      return true;
    },
    choices: [
      ...files.map((file: string) => ({
        title: file,
        value: file,
      })),
      { title: "Sair", value: "exit" },
      { title: "Voltar", value: "back" },
      { title: "Listar arquivos", value: "list-files" },
    ],
    suggest: async (input, choices) => {
      const normalized = input.toLowerCase();
      return choices.filter((choice) =>
        choice.value.toLowerCase().includes(normalized),
      );
    },
  });
};

export const confirmSeedPromptMenu = async (path: string) =>
  await prompts({
    type: "confirm",
    name: "confirm",
    message: `Deseja prosseguir com o seed da bíblia usando o arquivo ${path}?`,
    initial: true,
  });

export const seedOptionsPromptMenu = async () =>
  await prompts({
    type: "multiselect",
    name: "options",
    message: "Selecione as opções de seed:",
    choices: [
      { title: "Inserir livros", value: "books", selected: true },
      { title: "Inserir capítulos", value: "chapters", selected: true },
      { title: "Inserir versículos", value: "verses", selected: true },
    ],
  });

import prompts from "prompts";

export const databaseMenu = async () =>
  await prompts({
    type: "select",
    name: "choice",
    message: "Database — escolha uma opção:",
    choices: [
      { title: "Truncar tabelas", value: "truncate" },
      { title: "Voltar", value: "back" },
    ],
  });

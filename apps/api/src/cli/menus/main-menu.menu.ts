#!/usr/bin/env bun
import prompts from "prompts";

export const mainMenu = async () =>
  await prompts({
    type: "select",
    name: "main",
    message: "Escolha uma opção:",
    choices: [
      { title: "Bíblia", value: "bible" },
      { title: "Database", value: "database" },
      { title: "Gerar Docs OpenAPI", value: "openapi-doc" },
      { title: "Build", value: "build" },
      { title: "Deploy (Square Cloud)", value: "deploy" },
      { title: "Sair", value: "exit" },
    ],
  });

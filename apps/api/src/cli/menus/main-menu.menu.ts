#!/usr/bin/env bun
import prompts from "prompts";

export const mainMenu = async () => await prompts({
  type: "select",
  name: "main",
  message: "Escolha uma opção:",
  choices: [{title: "Seed", value: "seed"}, { title: "Validar arquivo json da bíblia", value: "validate-bible" }, { title: "Sair", value: "exit" }],
});

import prompts from "prompts";
import { logger } from "@/utils/logger";
import { buildProject } from "./build.service";
import { initCli } from "../../index";

export async function buildAction() {
  logger("info", "Iniciando build...");

  const response = await prompts({
    type: "confirm",
    name: "value",
    message: "Deseja compactar o projeto?",
    initial: true,
  });

  if (!response.value) {
    console.clear();
    await initCli(false);
    return;
  }

  try {
    const zipPath = await buildProject();
    logger("success", `Build criado com sucesso: ${zipPath}`);
  } catch (error) {
    const message = error instanceof Error ? error.message : String(error);
    logger("error", `Erro ao criar build: ${message}`);
  }

  logger("info", "Pressione Enter para continuar...");
  await prompts({ type: "text", name: "continue", message: "" });

  console.clear();
  await initCli(false);
}

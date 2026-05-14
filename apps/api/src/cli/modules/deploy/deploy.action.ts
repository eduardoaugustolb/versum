import prompts from "prompts";
import { logger } from "@versum/logger";
import { initCli } from "../../index";
import { deployToSquareCloud } from "./deploy.service";

export async function deployAction() {
  logger("info", "Iniciando deploy para Square Cloud...");

  const buildResponse = await prompts({
    type: "confirm",
    name: "value",
    message: "Deseja compilar o projeto antes do deploy?",
    initial: true,
  });

  if (!buildResponse.value) {
    console.clear();
    await initCli(false);
    return;
  }

  const apiKeyResponse = await prompts({
    type: "password",
    name: "value",
    message:
      "Informe sua Square Cloud API Key (ou defina SQUARE_API_KEY no .env):",
  });

  const apiKey =
    apiKeyResponse.value || process.env.SQUARE_API_KEY;

  if (!apiKey) {
    logger(
      "error",
      "API Key não informada. Defina SQUARE_API_KEY no .env ou informe agora.",
    );
    console.clear();
    await initCli(false);
    return;
  }

  try {
    const result = await deployToSquareCloud(apiKey);
    logger("success", `Deploy realizado com sucesso! ID: ${result.id}`);
    logger("info", `Nome: ${result.name}`);
  } catch (error) {
    const message = error instanceof Error ? error.message : String(error);
    logger("error", `Erro ao fazer deploy: ${message}`);
  }

  logger("info", "Pressione Enter para continuar...");
  await prompts({ type: "text", name: "continue", message: "" });

  console.clear();
  await initCli(false);
}

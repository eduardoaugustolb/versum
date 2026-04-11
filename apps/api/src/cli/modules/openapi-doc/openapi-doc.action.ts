import prompts from "prompts";
import { logger } from "@/utils/logger";
import { initCli } from "../../index";
import { generateOpenApiDocs } from "./openapi-doc.service";

export async function openapiDocAction() {
  logger("info", "Gerando documentação OpenAPI...");

  try {
    await generateOpenApiDocs();
    logger("success", "Documentação gerada com sucesso!");
  } catch (error) {
    logger("error", `Erro ao gerar documentação: ${error}`);
  }

  logger("info", "Pressione Enter para continuar...");
  await prompts({ type: "text", name: "continue", message: "" });

  console.clear();
  await initCli(false);
}

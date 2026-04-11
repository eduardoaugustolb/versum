import prompts from "prompts";
import { logger } from "@/utils/logger";
import { initCli } from "../../index";
import { type SeedBibleOptions, seedBibleFromJson } from "./seed/seed.action";
import {
  confirmSeedPromptMenu,
  seedBibleJsonPathPrompt,
  seedOptionsPromptMenu,
} from "./seed/seed.menus";

export const bibleMenu = async () =>
  await prompts({
    type: "select",
    name: "bible",
    message: "Escolha uma opção:",
    choices: [
      { title: "Seed (inserir no banco)", value: "seed" },
      { title: "Voltar", value: "back" },
    ],
  });

export async function bibleAction() {
  const menuResult = await bibleMenu();

  switch (menuResult.bible) {
    case "seed": {
      const { bible_json_path } = await seedBibleJsonPathPrompt();

      if (!bible_json_path) {
        console.clear();
        await bibleAction();
        return;
      }

      const lower = bible_json_path.toLowerCase();

      if (lower === "exit") {
        logger({ color: "red", icon: "", level: "info" }, "Saindo...");
        process.exit(0);
      }

      if (lower === "list-files") {
        logger("info", "Voltando para listar arquivos...");
        return await bibleAction();
      }

      if (lower === "back") {
        console.clear();
        return await bibleAction();
      }

      const { confirm } = await confirmSeedPromptMenu(bible_json_path);

      if (!confirm) {
        logger("info", "Seed cancelado.");
        return await bibleAction();
      }

      const { options } = await seedOptionsPromptMenu();

      const seedOptions: SeedBibleOptions = {
        insertBooks: options.includes("books"),
        insertChapters: options.includes("chapters"),
        insertVerses: options.includes("verses"),
      };

      await seedBibleFromJson(bible_json_path, seedOptions);

      logger("info", "Pressione Enter para continuar...");
      await prompts({ type: "text", name: "continue", message: "" });

      return await bibleAction();
    }

    case "back":
      console.clear();
      logger(
        { color: "blue", icon: "", level: "info" },
        "Voltando para o menu...",
      );
      await initCli(false);
      return;
  }
}

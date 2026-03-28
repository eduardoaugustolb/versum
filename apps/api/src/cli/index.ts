import { logger } from "@/utils/logger";
import { mainMenu } from "./menus/main-menu.menu";
import { bibleAction } from "./modules/bible/bible.action";
import { databaseAction } from "./modules/db/database.action";

export async function initCli(clear = true) {
  if (clear) console.clear();
  logger({ level: "info", color: "blue", icon: "" }, "Versum API CLI");
  const mainMenuResult = await mainMenu();

  switch (mainMenuResult.main) {
    case "bible":
      await bibleAction();
      break;
    case "database":
      await databaseAction();
      break;
    case "exit":
      process.exit(0);
      break;
  }
}

initCli();

import { logger } from "@/utils/logger";
import { mainMenu } from "./menus/main-menu.menu";
import { bibleAction } from "./modules/bible/bible.action";
import { buildAction } from "./modules/build/build.action";
import { databaseAction } from "./modules/db/database.action";
import { openapiDocAction } from "./modules/openapi-doc/openapi-doc.action";

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
    case "openapi-doc":
      await openapiDocAction();
      break;
    case "build":
      await buildAction();
      break;
    case "exit":
      process.exit(0);
      break;
  }
}

initCli();

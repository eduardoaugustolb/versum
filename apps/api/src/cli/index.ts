import { logger } from "@/utils/logger";
import { mainMenu } from "./menus/main-menu.menu";
import { validateBibleJsonAction } from "./modules/validate/bible-json/bible-json.action";

async function main() {
  console.clear();
  logger({ level: "info", color: "blue", icon: "" }, "Versum API CLI");
  const mainMenuResult = await mainMenu();
  switch (mainMenuResult.main) {
    case "seed":

      break;
    case "validate-bible":
      await validateBibleJsonAction();
      break;
    case "exit":
      process.exit(0);
      break;
  }
}

main();

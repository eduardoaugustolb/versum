import { logger } from "@/utils/logger";
import { initCli } from "../../index.ts";
import { truncateDatabaseAction } from "./truncate.action.ts";
import { databaseMenu } from "./database.menus.ts";

export async function databaseAction() {
  const result = await databaseMenu();

  if (!result?.choice) {
    return;
  }

  switch (result.choice) {
    case "truncate":
      await truncateDatabaseAction();
      return await databaseAction();

    case "back":
      console.clear();
      logger({ color: "blue", icon: "", level: "info" }, "Voltando para o menu...");
      await initCli(false);
      return;

    default:
      return await databaseAction();
  }
}

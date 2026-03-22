import { logger } from "@/utils/logger";
import { confirmValidationPromptMenu, saveResultJsonFilesMenu } from "./bible-json.menus";
import { executeValidation, writeEnrichedBible, writeMissingVerses } from "./bible-json.script";

export async function validateBibleJsonAction() {
  const confirmValidation = await confirmValidationPromptMenu();
  logger("debug", "confirmValidation", confirmValidation);
  if (confirmValidation) {

    const { errors, bible } = executeValidation();
    const saveResultJson = await saveResultJsonFilesMenu({ errors });
    if (saveResultJson.save_json_files.length > 0) {
      const jsonFiles = saveResultJson.save_json_files;
        const timestamp = new Date();
      for (const file of jsonFiles) {
        if (file === "parsed_bible_data") {
          await writeEnrichedBible(bible, timestamp);
        } else {
          await writeMissingVerses(errors, timestamp);
        }
      }
    }

  }
}

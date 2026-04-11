import { sql } from "drizzle-orm";
import prompts from "prompts";
import { db } from "@/infrastructure/db/config.ts";
import { logger } from "@/utils/logger";
import {
  assertTruncatableTableNames,
  TRUNCATE_PRESETS,
  type TruncatableTableName,
} from "./truncate.constants.ts";
import {
  truncateConfirmPrompt,
  truncateModeMenu,
  truncatePresetMenu,
  truncateTablesMultiselect,
} from "./truncate.menus.ts";

function quoteIdent(name: string): string {
  if (!/^[a-z_][a-z0-9_]*$/.test(name)) {
    throw new Error(`Identificador de tabela inválido: ${name}`);
  }
  return `"${name}"`;
}

export async function executeTruncate(tables: readonly TruncatableTableName[]) {
  assertTruncatableTableNames([...tables]);
  const unique = [...new Set(tables)];
  const list = unique.map(quoteIdent).join(", ");
  const query = `TRUNCATE TABLE ${list} RESTART IDENTITY CASCADE`;
  await db.execute(sql.raw(query));
}

export async function truncateDatabaseAction() {
  const mode = await truncateModeMenu();

  if (!mode.mode || mode.mode === "back") {
    console.clear();
    return;
  }

  let tables: TruncatableTableName[] = [];

  if (mode.mode === "preset") {
    const preset = await truncatePresetMenu();
    if (!preset.preset || preset.preset === "back") {
      return await truncateDatabaseAction();
    }
    const def = TRUNCATE_PRESETS[preset.preset];
    if (!def) {
      logger("error", "Preset inválido.");
      return await truncateDatabaseAction();
    }
    tables = [...def.tables];
    logger(
      "info",
      `${def.title}\n${def.description}\nTabelas: ${tables.join(", ")}`,
    );
  } else {
    const sel = await truncateTablesMultiselect();
    if (!sel?.tables?.length) {
      logger("info", "Nenhuma tabela selecionada.");
      return await truncateDatabaseAction();
    }
    tables = sel.tables as TruncatableTableName[];
    logger("info", `Tabelas: ${tables.join(", ")}`);
  }

  const confirm = await truncateConfirmPrompt();
  const token = confirm.token?.trim() ?? "";
  const lower = token.toLowerCase();

  if (lower === "voltar" || lower === "back") {
    logger("info", "Operação cancelada.");
    return;
  }

  if (token !== "TRUNCAR") {
    logger("error", "Confirmação incorreta. Abortado.");
    return;
  }

  try {
    await executeTruncate(tables);
    logger("success", `TRUNCATE concluído em: ${tables.join(", ")}`);
  } catch (err) {
    logger("error", `Falha no TRUNCATE: ${String(err)}`);
    console.error(err);
  }

  logger("info", "Pressione Enter para continuar...");
  await prompts({ type: "text", name: "continue", message: "" });
}

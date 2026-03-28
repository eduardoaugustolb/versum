import { existsSync, readdirSync } from "node:fs";
import { join } from "node:path";
import { cliOutputPath } from "@/cli/constants";

export async function findCliOutputFiles(cwd: string): Promise<string[]> {
  const files: string[] = [];
  const defaultBible = join(cwd, "src/assets/json/bible.json");
  if (existsSync(defaultBible)) {
    files.push(defaultBible);
  }

  const outputDir = `${cwd}/${cliOutputPath}`;

  try {
    const entries = readdirSync(outputDir, { recursive: true });

    for (const entry of entries) {
      if (typeof entry === "string" && entry.endsWith(".json") && !entry.endsWith("missing-verses.json")) {
        const full = `${outputDir}/${entry}`;
        if (!files.includes(full)) {
          files.push(full);
        }
      }
    }
  } catch {
    return files.sort().reverse();
  }

  return files.sort().reverse();
}

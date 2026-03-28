import JSZip from "jszip";
import { resolve, relative } from "node:path";
import {
  readdirSync,
  statSync,
  existsSync,
  mkdirSync,
  writeFileSync,
} from "node:fs";

const IGNORED_DIRS = new Set([
  "node_modules",
  ".git",
  ".vscode",
  ".idea",
  "dist",
  "build",
  ".cache",
  ".tmp",
  ".temp",
  ".output",
  ".build",
]);

const IGNORED_FILES = new Set([
  ".DS_Store",
  "Thumbs.db",
]);

async function addToZip(
  zip: JSZip,
  sourcePath: string,
  basePath: string,
): Promise<void> {
  const stat = statSync(sourcePath);
  const parts = sourcePath.split(/[/\\]/);
  const name = parts[parts.length - 1] ?? "";

  if (IGNORED_DIRS.has(name)) return;

  if (stat.isDirectory()) {
    const entries = readdirSync(sourcePath);
    for (const entry of entries) {
      const fullPath = resolve(sourcePath, entry);
      await addToZip(zip, fullPath, basePath);
    }
  } else {
    if (IGNORED_FILES.has(name)) return;
    if (name.startsWith(".env")) return;
    if (name === "bun.lock" || name.endsWith(".lock")) return;

    const archivePath = relative(basePath, sourcePath).replace(/\\/g, "/");
    const content = await Bun.file(sourcePath).arrayBuffer();
    zip.file(archivePath, new Uint8Array(content));
  }
}

export async function buildProject(): Promise<string> {
  const rootDir = resolve(import.meta.dir, "..", "..", "..", "..");
  const projectName = "versum-api";
  const outputDir = resolve(rootDir, ".build");
  const outputPath = resolve(outputDir, `${projectName}.zip`);

  if (!existsSync(outputDir)) {
    mkdirSync(outputDir, { recursive: true });
  }

  const zip = new JSZip();
  
  await addToZip(zip, rootDir, rootDir);
  
  const zipBuffer = await zip.generateAsync({ type: "nodebuffer", compression: "DEFLATE", compressionOptions: { level: 6 } });
  
  writeFileSync(outputPath, zipBuffer);

  return outputPath;
}

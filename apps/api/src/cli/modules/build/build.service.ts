import {
  readdirSync,
  statSync,
  readFileSync,
  existsSync,
} from "node:fs";
import { readdir, mkdir, rm, writeFile, stat } from "node:fs/promises";
import { resolve, relative } from "node:path";
import JSZip from "jszip";

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

const IGNORED_FILES = new Set([".DS_Store", "Thumbs.db"]);

const API_ROOT = resolve(import.meta.dir, "..", "..", "..", "..");
const WORKSPACE_ROOT = resolve(API_ROOT, "..", "..");
const LOGGER_ROOT = resolve(WORKSPACE_ROOT, "packages", "logger");

function shouldIgnore(name: string): boolean {
  if (IGNORED_DIRS.has(name)) return true;
  if (IGNORED_FILES.has(name)) return true;
  if (name.startsWith(".env")) return true;
  if (name.endsWith(".lock") && name !== "bun.lock") return true;
  return false;
}

export async function buildProject(): Promise<string> {
  const projectName = "versum-api";
  const buildDir = resolve(API_ROOT, ".build", projectName);
  const outputDir = resolve(API_ROOT, ".build");
  const outputPath = resolve(outputDir, `${projectName}.zip`);

  if (existsSync(buildDir)) {
    await rm(buildDir, { recursive: true, force: true });
  }
  await mkdir(buildDir, { recursive: true });

  // 1. Build server entry (compiles TS→JS, bundles internal code)
  const serverBuild = Bun.spawnSync([
    "bun",
    "build",
    "./src/server.ts",
    "--outdir",
    buildDir,
    "--target",
    "bun",
  ], {
    cwd: API_ROOT,
  });

  if (!serverBuild.success) {
    throw new Error(
      `Server build failed:\n${serverBuild.stderr.toString()}`,
    );
  }

  // 2. Build logger package
  const loggerBuildDir = resolve(buildDir, "packages", "logger");
  await mkdir(loggerBuildDir, { recursive: true });

  const loggerBuild = Bun.spawnSync([
    "bun",
    "build",
    "./index.ts",
    "--outdir",
    loggerBuildDir,
    "--target",
    "bun",
  ], {
    cwd: LOGGER_ROOT,
  });

  if (!loggerBuild.success) {
    throw new Error(
      `Logger build failed:\n${loggerBuild.stderr.toString()}`,
    );
  }

  // 3. Copy logger package.json with JS paths
  const loggerPkg = JSON.parse(
    readFileSync(resolve(LOGGER_ROOT, "package.json"), "utf-8"),
  );
  loggerPkg.module = "index.js";
  loggerPkg.main = "index.js";
  loggerPkg.types = "index.js";
  loggerPkg.exports = { ".": "./index.js" };
  delete loggerPkg.devDependencies;
  delete loggerPkg.peerDependencies;
  delete loggerPkg.scripts;

  await writeFile(
    resolve(loggerBuildDir, "package.json"),
    JSON.stringify(loggerPkg, null, 2),
  );

  // 4. Generate deploy package.json
  const apiPkg = JSON.parse(
    readFileSync(resolve(API_ROOT, "package.json"), "utf-8"),
  );

  const deployPkg: Record<string, unknown> = {
    name: projectName,
    version: apiPkg.version,
    module: "server.js",
    main: "server.js",
    type: "module",
    scripts: {
      start: "bun run server.js",
    },
  };

  const dependencies: Record<string, string> = {};
  for (const [name, version] of Object.entries(
    apiPkg.dependencies as Record<string, string>,
  )) {
    if (version === "workspace:*") {
      dependencies[name] = "file:./packages/logger";
    } else {
      dependencies[name] = version;
    }
  }
  deployPkg.dependencies = dependencies;

  await writeFile(
    resolve(buildDir, "package.json"),
    JSON.stringify(deployPkg, null, 2),
  );

  // 5. Generate lockfile by running bun install
  const installResult = Bun.spawnSync(["bun", "install"], {
    cwd: buildDir,
  });

  if (!installResult.success) {
    throw new Error(
      `Dependency installation failed:\n${installResult.stderr.toString()}`,
    );
  }

  // 6. Remove node_modules (keep zip small, SquareCloud will reinstall)
  await rm(resolve(buildDir, "node_modules"), {
    recursive: true,
    force: true,
  });

  // 7. Zip the build directory
  const zip = new JSZip();
  await addDirToZip(zip, buildDir, buildDir);

  const zipBuffer = await zip.generateAsync({
    type: "nodebuffer",
    compression: "DEFLATE",
    compressionOptions: { level: 6 },
  });

  // 6. Clean up build dir and write final zip
  await rm(buildDir, { recursive: true, force: true });
  await writeFile(outputPath, new Uint8Array(zipBuffer));

  return outputPath;
}

async function addDirToZip(
  zip: JSZip,
  dirPath: string,
  basePath: string,
): Promise<void> {
  const entries = await readdir(dirPath);

  for (const entry of entries) {
    if (shouldIgnore(entry)) continue;

    const fullPath = resolve(dirPath, entry);
    let isDir: boolean;

    try {
      isDir = (await stat(fullPath)).isDirectory();
    } catch {
      continue;
    }

    if (isDir) {
      await addDirToZip(zip, fullPath, basePath);
    } else {
      const archivePath = relative(basePath, fullPath);
      const content = await Bun.file(fullPath).arrayBuffer();
      zip.file(archivePath, new Uint8Array(content));
    }
  }
}

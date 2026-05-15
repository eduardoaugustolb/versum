import { readFileSync, existsSync } from "node:fs";
import { readdir, mkdir, rm, writeFile, stat } from "node:fs/promises";
import { resolve, relative } from "node:path";
import JSZip from "jszip";

const API_ROOT = resolve(import.meta.dir, "..", "..", "..", "..");

async function addDirToZip(
  zip: JSZip,
  dirPath: string,
  basePath: string,
): Promise<void> {
  const entries = await readdir(dirPath);

  for (const entry of entries) {
    if (entry === "node_modules" || entry === ".build") continue;
    if (entry.startsWith(".env")) continue;

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

export async function buildProject(): Promise<string> {
  const buildDir = resolve(API_ROOT, ".build", "versum-api");
  const outputPath = resolve(API_ROOT, ".build", "versum-api.zip");

  if (existsSync(buildDir)) {
    await rm(buildDir, { recursive: true, force: true });
  }
  await mkdir(buildDir, { recursive: true });

  // Read package.json to know which deps are workspace vs npm
  const apiPkg = JSON.parse(
    readFileSync(resolve(API_ROOT, "package.json"), "utf-8"),
  );

  // Externalize all npm deps (only workspace deps are bundled inline)
  const npmDeps = Object.entries(
    apiPkg.dependencies as Record<string, string>,
  )
    .filter(([, version]) => version !== "workspace:*")
    .map(([name]) => name);

  // 1. Build with Bun.build — bundles source + workspace deps,
  //    externalizes npm packages so they load from node_modules at runtime
  const result = await Bun.build({
    entrypoints: [resolve(API_ROOT, "src/server.ts")],
    outdir: resolve(buildDir, "src"),
    target: "bun",
    external: npmDeps,
  });

  if (!result.success) {
    throw new Error(
      `Server build failed:\n${result.logs
        .filter((l) => l.level === "error")
        .map((l) => l.message)
        .join("\n")}`,
    );
  }

  // 2. Generate deploy package.json (only npm deps, workspace ones are bundled)
  const dependencies: Record<string, string> = {};
  for (const [name, version] of Object.entries(
    apiPkg.dependencies as Record<string, string>,
  )) {
    if (version !== "workspace:*") {
      dependencies[name] = version;
    }
  }

  const deployPkg = {
    name: "versum-api",
    version: apiPkg.version,
    module: "src/server.js",
    main: "src/server.js",
    type: "module",
    scripts: {
      start: "bun run src/server.js",
    },
    dependencies,
  };

  await writeFile(
    resolve(buildDir, "package.json"),
    JSON.stringify(deployPkg, null, 2),
  );

  // 3. Generate lockfile
  const install = Bun.spawnSync(["bun", "install"], { cwd: buildDir });

  if (!install.success) {
    throw new Error(
      `bun install failed:\n${install.stderr.toString()}`,
    );
  }

  // 4. Remove node_modules (SquareCloud runs bun install on deploy)
  await rm(resolve(buildDir, "node_modules"), {
    recursive: true,
    force: true,
  });

  // 5. Zip
  const zip = new JSZip();
  await addDirToZip(zip, buildDir, buildDir);

  const zipBuffer = await zip.generateAsync({
    type: "nodebuffer",
    compression: "DEFLATE",
    compressionOptions: { level: 6 },
  });

  await rm(buildDir, { recursive: true, force: true });
  await writeFile(outputPath, new Uint8Array(zipBuffer));

  return outputPath;
}

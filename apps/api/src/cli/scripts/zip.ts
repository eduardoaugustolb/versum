import { buildProject } from "../modules/build/build.service";

const zipPath = await buildProject();
console.log(zipPath);

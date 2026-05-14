import { existsSync } from "node:fs";
import { buildProject } from "../build/build.service";

const SQUARE_API_URL = "https://api.squarecloud.app/v2/apps/upload";

interface SquareCloudUploadResponse {
  status: string;
  code: string;
  response: {
    id: string;
    name: string;
    description: string;
    domain: string | null;
    ram: number;
    cpu: number;
    language: {
      name: string;
      version: string;
    };
  };
}

export async function deployToSquareCloud(apiKey: string) {
  const zipPath = await buildProject();

  if (!existsSync(zipPath)) {
    throw new Error(`Arquivo zip não encontrado: ${zipPath}`);
  }

  const file = Bun.file(zipPath);
  const formData = new FormData();
  formData.append("file", new Blob([await file.arrayBuffer()]), "versum-api.zip");

  const response = await fetch(SQUARE_API_URL, {
    method: "POST",
    headers: {
      Authorization: apiKey,
    },
    body: formData,
  });

  if (!response.ok) {
    const text = await response.text();
    throw new Error(`Square Cloud API retornou ${response.status}: ${text}`);
  }

  const data = (await response.json()) as SquareCloudUploadResponse;

  if (data.status !== "success") {
    throw new Error(`Square Cloud API retornou erro: ${data.code}`);
  }

  return data.response;
}

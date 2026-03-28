import * as fs from "fs";
import * as path from "path";
import * as yaml from "yaml";
import { cliOutputPath } from "../../constants";

interface OpenApiSchema {
  type?: string;
  properties?: Record<string, OpenApiSchema>;
  items?: OpenApiSchema;
  $ref?: string;
  allOf?: OpenApiSchema[];
  anyOf?: OpenApiSchema[];
  format?: string;
  enum?: string[];
  description?: string;
  required?: string[];
  default?: string;
  example?: string;
  nullable?: boolean;
  minLength?: number;
  maxLength?: number;
  minimum?: number;
  maximum?: number;
  pattern?: string;
}

interface OpenApiParameter {
  name: string;
  in: string;
  required?: boolean;
  description?: string;
  schema?: OpenApiSchema;
}

interface OpenApiResponse {
  description?: string;
  content?: Record<string, { schema?: OpenApiSchema }>;
}

interface OpenApiOperation {
  tags?: string[];
  summary?: string;
  description?: string;
  parameters?: OpenApiParameter[];
  requestBody?: {
    content?: Record<string, { schema?: OpenApiSchema }>;
  };
  responses?: Record<string, OpenApiResponse>;
}

interface OpenApiPath {
  get?: OpenApiOperation;
  post?: OpenApiOperation;
  patch?: OpenApiOperation;
  put?: OpenApiOperation;
  delete?: OpenApiOperation;
}

interface OpenApiSpec {
  openapi?: string;
  info?: {
    title?: string;
    version?: string;
    description?: string;
  };
  paths?: Record<string, OpenApiPath>;
  components?: {
    schemas?: Record<string, OpenApiSchema>;
  };
}

function resolveRef(ref: string, spec: OpenApiSpec): OpenApiSchema | null {
  if (!ref.startsWith("#/")) return null;
  const parts = ref.slice(2).split("/");
  let current: unknown = spec;
  for (const part of parts) {
    current = (current as Record<string, unknown>)?.[part];
  }
  return current as OpenApiSchema | null;
}

function resolveSchema(schema: OpenApiSchema | undefined, spec: OpenApiSpec): OpenApiSchema {
  if (!schema) return {};

  if (schema.$ref) {
    return resolveSchema(resolveRef(schema.$ref, spec) ?? undefined, spec);
  }

  if (schema.allOf) {
    const merged: OpenApiSchema = { ...schema };
    delete merged.allOf;
    for (const subSchema of schema.allOf) {
      const resolved = resolveSchema(subSchema, spec);
      if (resolved.properties) {
        merged.properties = { ...merged.properties, ...resolved.properties };
      }
      if (resolved.required) {
        merged.required = [...(merged.required || []), ...resolved.required];
      }
      if (resolved.description) {
        merged.description = resolved.description;
      }
    }
    return merged;
  }

  if (schema.anyOf) {
    return resolveSchema(schema.anyOf[0], spec);
  }

  return schema;
}

function getSchemaType(schema: OpenApiSchema): string {
  if (schema.type) return schema.type;
  if (schema.enum) return "enum";
  if (schema.properties) return "object";
  if (schema.items) return "array";
  return "unknown";
}

function getSchemaRefName(schema: OpenApiSchema): string | null {
  if (schema.$ref) {
    const parts = schema.$ref.split("/");
    return parts[parts.length - 1] || null;
  }
  if (schema.allOf) {
    for (const sub of schema.allOf) {
      const name = getSchemaRefName(sub);
      if (name) return name;
    }
  }
  return null;
}

function formatSchemaType(schema: OpenApiSchema | undefined, spec: OpenApiSpec): string {
  if (!schema) return "`any`";

  if (schema.$ref) {
    const name = getSchemaRefName(schema);
    if (name) {
      return `[${name}](#${toAnchorId(name)})`;
    }
  }
  if (schema.allOf) {
    const name = getSchemaRefName(schema);
    if (name) {
      return `[${name}](#${toAnchorId(name)})`;
    }
  }
  const type = getSchemaType(schema);
  return `\`${schema.type || type}\``;
}

function toAnchorId(name: string): string {
  return name.toLowerCase().replace(/([A-Z])/g, "-$1").replace(/^-/, "");
}

function formatProperty(schema: OpenApiSchema, spec: OpenApiSpec, indent: number): string {
  const spaces = "  ".repeat(indent);
  const type = getSchemaType(schema);
  let result = `\`${schema.type || type}\``;

  if (schema.format) result += ` (${schema.format})`;
  if (schema.nullable) result += " | null";
  if (schema.enum) result += ` = ${schema.enum.map((e) => `\`${e}\``).join(" | ")}`;
  if (schema.default !== undefined) result += ` = ${schema.default}`;
  if (schema.description) result += `\n${spaces}  ${schema.description}`;

  return result;
}

function generateSchemaDocs(schemaName: string, schema: OpenApiSchema, spec: OpenApiSpec): string {
  const lines: string[] = [];
  const anchorId = toAnchorId(schemaName);
  lines.push(`### \`${schemaName}\` {#${anchorId}}`);

  if (schema.description) {
    lines.push(schema.description);
    lines.push("");
  }

  if (schema.type === "object" && schema.properties) {
    lines.push("| Campo | Tipo | Descrição |");
    lines.push("|-------|------|-----------|");

    const required = schema.required || [];
    for (const [propName, propSchema] of Object.entries(schema.properties)) {
      const typeStr = formatSchemaType(propSchema, spec);
      const resolved = resolveSchema(propSchema, spec);
      const isRequired = required.includes(propName) ? "**" : "";
      const desc = resolved.description || "";
      lines.push(`| ${isRequired}${propName}${isRequired} | ${typeStr} | ${desc} |`);
    }
    lines.push("");
  } else if (schema.type === "array" && schema.items) {
    const typeStr = formatSchemaType(schema.items, spec);
    const resolved = resolveSchema(schema.items, spec);
    lines.push(`Array of ${typeStr}`);
    if (resolved.description) lines.push(resolved.description);
    lines.push("");
  } else if (schema.enum) {
    lines.push(schema.enum.map((e) => `- \`${e}\``).join("\n"));
    lines.push("");
  }

  return lines.join("\n");
}

function generateEndpointDocs(path: string, method: string, operation: OpenApiOperation, spec: OpenApiSpec): string {
  const lines: string[] = [];
  const methodUpper = method.toUpperCase();
  const methodBadge = `**\`${methodUpper}\`**`;

  lines.push(`#### ${methodBadge} \`${path}\``);
  if (operation.summary) lines.push(`**${operation.summary}**`);
  if (operation.description) lines.push(operation.description);
  lines.push("");

  if (operation.parameters && operation.parameters.length > 0) {
    lines.push("**Parâmetros:**");
    lines.push("");
    lines.push("| Nome | Local | Tipo | Obrigatório | Descrição |");
    lines.push("|------|-------|------|-------------|-----------|");

    for (const param of operation.parameters) {
      const type = formatSchemaType(param.schema, spec);
      const required = param.required ? "Sim" : "Não";
      lines.push(
        `| \`${param.name}\` | ${param.in} | ${type} | ${required} | ${param.description || ""} |`
      );
    }
    lines.push("");
  }

  if (operation.requestBody?.content?.["application/json"]?.schema) {
    lines.push("**Body (JSON):**");
    const bodySchema = resolveSchema(operation.requestBody.content["application/json"].schema, spec);
    lines.push("");
    if (bodySchema.properties) {
      lines.push("| Campo | Tipo | Obrigatório | Descrição |");
      lines.push("|-------|------|-------------|-----------|");
      const required = bodySchema.required || [];
      for (const [propName, propSchema] of Object.entries(bodySchema.properties)) {
        const resolved = resolveSchema(propSchema, spec);
        const typeStr = formatSchemaType(resolved, spec);
        const isRequired = required.includes(propName) ? "Sim" : "Não";
        lines.push(
          `| \`${propName}\` | ${typeStr} | ${isRequired} | ${resolved.description || ""} |`
        );
      }
    }
    lines.push("");
  }

  if (operation.responses) {
    lines.push("**Respostas:**");
    lines.push("");
    for (const [code, response] of Object.entries(operation.responses)) {
      lines.push(`- \`${code}\`: ${response.description || ""}`);
    }
    lines.push("");
  }

  return lines.join("\n");
}

export async function generateOpenApiDocs(): Promise<void> {
  const openApiPath = path.resolve(process.cwd(), "..", "client", "openapi.yaml");

  if (!fs.existsSync(openApiPath)) {
    throw new Error(`OpenAPI file not found: ${openApiPath}`);
  }

  const openApiContent = fs.readFileSync(openApiPath, "utf-8");
  const spec: OpenApiSpec = yaml.parse(openApiContent);

  const lines: string[] = [];
  lines.push(`# ${spec.info?.title || "API Documentation"}`);
  lines.push("");
  lines.push(`**Version:** ${spec.info?.version || "1.0.0"}`);
  if (spec.openapi) lines.push(`**OpenAPI:** ${spec.openapi}`);
  lines.push("");

  if (spec.info?.description) {
    lines.push(spec.info.description);
    lines.push("");
  }

  if (spec.components?.schemas) {
    lines.push("## Schemas");
    lines.push("");

    for (const [schemaName, schema] of Object.entries(spec.components.schemas)) {
      lines.push(generateSchemaDocs(schemaName, resolveSchema(schema, spec), spec));
    }
  }

  if (spec.paths) {
    const pathsByTag: Record<string, Array<{ path: string; method: string; operation: OpenApiOperation }>> = {};

    for (const [pathStr, pathItem] of Object.entries(spec.paths)) {
      for (const [method, operation] of Object.entries(pathItem)) {
        if (!operation) continue;
        const tags = operation.tags || ["Other"];
        for (const tag of tags) {
          if (!pathsByTag[tag]) pathsByTag[tag] = [];
          pathsByTag[tag].push({ path: pathStr, method, operation });
        }
      }
    }

    lines.push("## Endpoints");
    lines.push("");

    for (const [tag, endpoints] of Object.entries(pathsByTag)) {
      lines.push(`### ${tag}`);
      lines.push("");
      for (const { path: pathStr, method, operation } of endpoints) {
        lines.push(generateEndpointDocs(pathStr, method, operation, spec));
      }
    }
  }

  const outputDir = path.resolve(process.cwd(), cliOutputPath);
  if (!fs.existsSync(outputDir)) {
    fs.mkdirSync(outputDir, { recursive: true });
  }

  const outputFile = path.join(outputDir, "api-documentation.md");
  fs.writeFileSync(outputFile, lines.join("\n"), "utf-8");
}

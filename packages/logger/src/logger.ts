import { validateLevel, validateColor, validateIcon, validateArgs } from "./validator.ts";
import { getIsDev, getIsDebugEnabled } from "./env.ts";
import { processSensitiveValues } from "./sensitive.ts";
import { getColor, getIcon } from "./formatter.ts";
import type { LoggerConfig } from "./modules/config/types/index.ts";

function log(level: string, ...args: unknown[]) {
  const validLevel = validateLevel(level);

  if (validLevel === "debug" && getIsDebugEnabled()) return;

  const sanitizedArgs = validateArgs(args);
  const processedArgs = processSensitiveValues(sanitizedArgs);

  const color = getColor(validLevel);
  const icon = getIcon(validLevel);

  console.log(color, icon, ...processedArgs);
}

export function logger(config: LoggerConfig, ...args: unknown[]) {
  const level = typeof config === "string" ? config : config.level;
  const validLevel = validateLevel(level);
  const customColor = typeof config === "object" ? validateColor(config.color) : undefined;
  const customIcon = typeof config === "object" ? validateIcon(config.icon) : undefined;

  if (validLevel === "debug" && getIsDebugEnabled()) return;

  const sanitizedArgs = validateArgs(args);
  const processedArgs = processSensitiveValues(sanitizedArgs);

  const color = getColor(validLevel, customColor);
  const icon = getIcon(validLevel, customIcon);

  console.log(color, icon, ...processedArgs);
}

export const info = (...args: unknown[]) => log("info", ...args);
export const warn = (...args: unknown[]) => log("warn", ...args);
export const error = (...args: unknown[]) => log("error", ...args);
export const debug = (...args: unknown[]) => log("debug", ...args);
export const success = (...args: unknown[]) => log("success", ...args);
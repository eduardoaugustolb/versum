import { env } from "../env/parser.ts";
import { colors, defaultColor, icons, presets } from "./constants.ts";
import { LoggerConfig, LoggerConfigObject } from "./types.ts";

export function logger(paramConfig: LoggerConfig, ...args: unknown[]) {
  let config: LoggerConfigObject;

  if (typeof paramConfig === "string") {
    config = presets[paramConfig];
  } else {
    config = paramConfig;
  }

  if (config.level === "debug" && !env.DEBUG) return;
  const color =
    colors[config.color as keyof typeof colors] || defaultColor[config.level];

  const icon = config.icon || icons[config.level as keyof typeof icons];

  console.log(color, icon, ...args);
}

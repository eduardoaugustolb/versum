import type { LoggerConfigLevel } from "../config/types";

export const COLORS = {
  reset: "\x1b[0m",
  red: "\x1b[31m",
  green: "\x1b[32m",
  yellow: "\x1b[33m",
  blue: "\x1b[34m",
  gray: "\x1b[90m",
};

export const PRESET_COLORS:Record<LoggerConfigLevel, string> = {
  debug: COLORS.yellow,
  info: COLORS.blue,
  warn: COLORS.yellow,
  error: COLORS.red,
  success: COLORS.green,
  log: COLORS.gray,
};

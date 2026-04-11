import type { LoggerConfigLevel, LoggerConfigObject } from "./types.ts";

export const colors = {
  reset: "\x1b[0m",
  red: "\x1b[31m",
  green: "\x1b[32m",
  yellow: "\x1b[33m",
  blue: "\x1b[34m",
  gray: "\x1b[90m",
};

export const defaultColor = {
  error: colors.red,
  warn: colors.yellow,
  info: colors.blue,
  debug: colors.gray,
  success: colors.green,
};

export const icons = {
  success: "🟢",
  info: "🔵",
  warn: "🟡",
  error: "🔴",
  debug: "🟣",
};

export const presets: Record<LoggerConfigLevel, LoggerConfigObject> = {
  debug: {
    level: "debug",
    color: "yellow",
  },
  info: {
    level: "info",
    color: "gray",
  },
  warn: {
    level: "warn",
    color: "gray",
  },
  error: {
    level: "error",
    color: "gray",
  },
  success: {
    level: "success",
    color: "green",
  },
};

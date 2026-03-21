import { colors, icons } from "./constants.ts";

export type LoggerConfigLevel = "debug" | "info" | "warn" | "error" | "success";

export type LoggerConfigObject = {
  level: LoggerConfigLevel;
  color?: keyof typeof colors;
  icon?: keyof typeof icons | string | null;
};

export type LoggerConfig = LoggerConfigLevel | LoggerConfigObject;

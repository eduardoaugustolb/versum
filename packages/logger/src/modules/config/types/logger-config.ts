import type { COLORS } from "../../color";
import type { ICONS } from "../../icon";

export type LoggerConfigLevel = "debug" | "info" | "warn" | "error" | "success" | "log";

export type LoggerConfigObject = {
  level: LoggerConfigLevel;
  color?: keyof typeof COLORS;
  icon?: keyof typeof ICONS | string | null;
};

export type LoggerConfig = LoggerConfigLevel | LoggerConfigObject;

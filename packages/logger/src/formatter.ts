import { PRESET_COLORS, ICONS } from "./modules";
import type { LoggerConfigLevel } from "./modules";

const DEFAULT_COLOR = "\x1b[90m"; // gray

export function getColor(level: string, customColor?: string): string {
  return customColor || (level in PRESET_COLORS ? PRESET_COLORS[level as LoggerConfigLevel] : DEFAULT_COLOR);
}

export function getIcon(level: string, customIcon?: string): string {
  return customIcon || (level in ICONS ? ICONS[level as keyof typeof ICONS] : "");
}

export function formatLog(level: string, args: unknown[]): { style: string; icon: string; message: unknown[] } {
  const style = getColor(level as LoggerConfigLevel);
  const icon = getIcon(level as LoggerConfigLevel);
  return { style, icon, message: args };
}
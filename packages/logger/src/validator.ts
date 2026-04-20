const VALID_LEVELS = ["info", "warn", "error", "debug", "log", "success"] as const;
const DEFAULT_LEVEL = "log";

export type LogLevel = (typeof VALID_LEVELS)[number];

export { VALID_LEVELS, DEFAULT_LEVEL };

export function validateLevel(level: string): LogLevel {
  if (!VALID_LEVELS.includes(level as LogLevel)) {
    console.warn(
      `[logger] Unknown level "${level}". Using "${DEFAULT_LEVEL}" as fallback.`
    );
    return DEFAULT_LEVEL;
  }
  return level as LogLevel;
}

export function validateColor(color: unknown): string | undefined {
  if (typeof color !== "string" && color !== undefined) {
    console.warn(`[logger] Invalid color type. Expected string, got ${typeof color}.`);
    return undefined;
  }
  return color;
}

export function validateIcon(icon: unknown): string | undefined {
  if (typeof icon !== "string" && icon !== undefined) {
    console.warn(`[logger] Invalid icon type. Expected string, got ${typeof icon}.`);
    return undefined;
  }
  return icon;
}

export function validateArgs(args: unknown[]): unknown[] {
  return args.map((arg) => {
    if (arg === null) return "null";
    if (arg === undefined) return "undefined";
    if (typeof arg === "function") return "[Function]";
    if (typeof arg === "symbol") return arg.toString();
    return arg;
  });
}

import { getIsDev, getIsDebugEnabled } from "./env.ts";

const SENSITIVE_MARKER = Symbol("SENSITIVE");

class SensitiveWrapper {
  constructor(
    public value: unknown,
    public shouldMask = true
  ) {}
}

export function sensitive<T>(value: T, shouldMask = true): SensitiveWrapper {
  return new SensitiveWrapper(value, shouldMask);
}

export function isSensitive(value: unknown): value is SensitiveWrapper {
  return value instanceof SensitiveWrapper;
}

function maskString(str: string): string {
  if (str.length <= 4) return "****";
  return str.slice(0, 2) + "****" + str.slice(-2);
}

function maskObject(obj: object): object {
  const masked: Record<string, unknown> = {};
  for (const [key, value] of Object.entries(obj)) {
    masked[key] = processValue(value);
  }
  return masked;
}

function maskArray(arr: unknown[]): unknown[] {
  return arr.map(processValue);
}

function processValue(value: unknown): unknown {
  const isDev = getIsDev();
  const isDebug = getIsDebugEnabled();

  if (isSensitive(value)) {
    if (value.shouldMask && !(isDev && isDebug)) {
      if (typeof value.value === "string") return maskString(value.value);
      if (Array.isArray(value.value)) return maskArray(value.value);
      if (typeof value.value === "object" && value.value !== null) {
        return maskObject(value.value as object);
      }
      return "****";
    }
    return value.value;
  }

  if (typeof value === "string") {
    return value;
  }

  if (Array.isArray(value)) {
    return value.map(processValue);
  }

  if (typeof value === "object" && value !== null) {
    return maskObject(value);
  }

  return value;
}

export function processSensitiveValues(args: unknown[]): unknown[] {
  return args.map(processValue);
}

export { SensitiveWrapper, SENSITIVE_MARKER };

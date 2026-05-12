import { ContentType } from "../types/enums"

export function maskSensitiveString(value: string): string {
  return '*'.repeat(value.length);
}

export function sanitizeSensitiveString(value: string): string {
  return value.toString().replace(/\s+/g, ' ');
}

export function splitSensitiveString(value: string): string[] {
  return value.toString().split(ContentType.SENSITIVE_STRING);
}

export function sensitiveString(value: string): string {
  const sanitized = sanitizeSensitiveString(value);
  return ContentType.SENSITIVE_STRING + sanitized + ContentType.SENSITIVE_STRING;
}

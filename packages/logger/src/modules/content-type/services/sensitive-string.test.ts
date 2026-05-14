import { describe, expect, it } from "vitest";
import {
  maskSensitiveString,
  sanitizeSensitiveString,
  splitSensitiveString,
  sensitiveString,
} from "./sensitive-string.service.ts";

describe("maskSensitiveString", () => {
  it("should mask string with asterisks", () => {
    expect(maskSensitiveString("secret")).toBe("******");
  });

  it("should handle empty string", () => {
    expect(maskSensitiveString("")).toBe("");
  });

  it("should preserve length of original string", () => {
    expect(maskSensitiveString("token_abc123").length).toBe(12);
  });
});

describe("sanitizeSensitiveString", () => {
  it("should replace multiple spaces with single space", () => {
    expect(sanitizeSensitiveString("hello    world")).toBe("hello world");
  });

  it("should not trim leading and trailing spaces", () => {
    expect(sanitizeSensitiveString("  hello world  ")).toBe(" hello world ");
  });

  it("should handle string without extra spaces", () => {
    expect(sanitizeSensitiveString("hello world")).toBe("hello world");
  });
});

describe("splitSensitiveString", () => {
  it("should split string by sensitive marker", () => {
    const input = sensitiveString("my_secret");
    expect(splitSensitiveString(input)).toEqual(["", "my_secret", ""]);
  });

  it("should return array with original string if no marker", () => {
    expect(splitSensitiveString("no_marker")).toEqual(["no_marker"]);
  });
});

describe("sensitiveString", () => {
  it("should wrap value with marker", () => {
    const result = sensitiveString("my_token");
    expect(result).toContain("my_token");
    expect(result).toContain("_:_SENSITIVE_STRING_:_");
  });

  it("should wrap value on both sides", () => {
    const result = sensitiveString("abc");
    expect(result.startsWith("_:_SENSITIVE_STRING_:_")).toBe(true);
    expect(result.endsWith("_:_SENSITIVE_STRING_:_")).toBe(true);
  });
});
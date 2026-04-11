import { describe, expect, it } from "vitest";
import { extractErrorMessage, isValidationError } from "./utils";

describe("error utils", () => {
  describe("isValidationError", () => {
    it("should return true for ZodError", () => {
      const zodError = new Error("Validation failed");
      zodError.name = "ZodError";

      expect(isValidationError(zodError)).toBe(true);
    });

    it("should return true for error with issues array", () => {
      const error = new Error("Validation failed");
      (error as Error & { issues: unknown }).issues = [
        { message: "Invalid input" },
      ];

      expect(isValidationError(error)).toBe(true);
    });

    it("should return true for ZodError with issues array", () => {
      const error = new Error("Validation failed");
      error.name = "ZodError";
      (error as Error & { issues: unknown }).issues = [
        { message: "Invalid input" },
      ];

      expect(isValidationError(error)).toBe(true);
    });

    it("should return false for regular Error", () => {
      const error = new Error("Something went wrong");

      expect(isValidationError(error)).toBe(false);
    });

    it("should return false for null", () => {
      expect(isValidationError(null)).toBe(false);
    });

    it("should return false for undefined", () => {
      expect(isValidationError(undefined)).toBe(false);
    });

    it("should return false for primitive values", () => {
      expect(isValidationError("string")).toBe(false);
      expect(isValidationError(123)).toBe(false);
      expect(isValidationError(true)).toBe(false);
    });

    it("should return false for object without name and issues", () => {
      expect(isValidationError({})).toBe(false);
    });

    it("should return true for error with empty issues array", () => {
      const error = new Error("Error");
      (error as Error & { issues: unknown }).issues = [];

      expect(isValidationError(error)).toBe(true);
    });
  });

  describe("extractErrorMessage", () => {
    it("should extract message from string error", () => {
      expect(extractErrorMessage("Simple error message")).toBe(
        "Simple error message",
      );
    });

    it("should extract message from empty string error with fallback", () => {
      expect(extractErrorMessage("", "Custom fallback")).toBe(
        "Custom fallback",
      );
    });

    it("should extract message from empty string without fallback", () => {
      expect(extractErrorMessage("")).toBe("Internal Server Error");
    });

    it("should extract message from Error object", () => {
      const error = new Error("Error message");
      expect(extractErrorMessage(error)).toBe("Error message");
    });

    it("should extract message from Error with issues array", () => {
      const error = new Error("Error");
      (error as Error & { issues: Array<{ message?: unknown }> }).issues = [
        { message: "First issue message" },
        { message: "Second issue message" },
      ];

      expect(extractErrorMessage(error)).toBe("First issue message");
    });

    it("should use fallback for Error with empty message", () => {
      const error = new Error("");
      expect(extractErrorMessage(error, "Custom fallback")).toBe(
        "Custom fallback",
      );
    });

    it("should use default fallback for Error with empty message", () => {
      const error = new Error("");
      expect(extractErrorMessage(error)).toBe("Internal Server Error");
    });

    it("should handle Error without message property", () => {
      const error = { name: "Error", message: "" };
      expect(extractErrorMessage(error)).toBe("Internal Server Error");
    });

    it("should handle Error with ZodError name and serialized message", () => {
      const error = new Error('[{"message":"Field is required"}]');
      error.name = "ZodError";

      expect(extractErrorMessage(error)).toBe("Field is required");
    });

    it("should handle Error with ZodError name and invalid serialized message", () => {
      const error = new Error('[{"invalid":"Field is required"}]');
      error.name = "ZodError";

      expect(extractErrorMessage(error)).toBe(
        '[{"invalid":"Field is required"}]',
      );
    });

    it("should handle Error with ZodError name and empty array message", () => {
      const error = new Error("[]");
      error.name = "ZodError";

      expect(extractErrorMessage(error)).toBe("[]");
    });

    it("should handle Error with ZodError name and empty object message", () => {
      const error = new Error("[{}]");
      error.name = "ZodError";

      expect(extractErrorMessage(error)).toBe("[{}]");
    });

    it("should handle null issues array", () => {
      const error = new Error("Error message");
      (error as Error & { issues: null }).issues = null;

      expect(extractErrorMessage(error)).toBe("Error message");
    });

    it("should handle issues array with empty string message", () => {
      const error = new Error("Error");
      (error as Error & { issues: Array<{ message?: unknown }> }).issues = [
        { message: "" },
      ];

      expect(extractErrorMessage(error)).toBe("Error");
    });

    it("should handle issues array with non-string message", () => {
      const error = new Error("Error");
      (error as Error & { issues: Array<{ message?: unknown }> }).issues = [
        { message: 123 },
        { message: null },
        { message: undefined },
      ];

      expect(extractErrorMessage(error)).toBe("Error");
    });

    it("should return fallback for null input", () => {
      expect(extractErrorMessage(null, "Custom fallback")).toBe(
        "Custom fallback",
      );
    });

    it("should return fallback for undefined input", () => {
      expect(extractErrorMessage(undefined, "Custom fallback")).toBe(
        "Custom fallback",
      );
    });

    it("should return default fallback for null input", () => {
      expect(extractErrorMessage(null)).toBe("Internal Server Error");
    });

    it("should return default fallback for undefined input", () => {
      expect(extractErrorMessage(undefined)).toBe("Internal Server Error");
    });

    it("should return fallback for number input", () => {
      expect(extractErrorMessage(0, "Fallback for number")).toBe(
        "Fallback for number",
      );
    });

    it("should return fallback for object input", () => {
      expect(extractErrorMessage({}, "Fallback for object")).toBe(
        "Fallback for object",
      );
    });

    it("should return fallback for array input", () => {
      expect(extractErrorMessage([], "Fallback for array")).toBe(
        "Fallback for array",
      );
    });

    it("should extract message from very long error", () => {
      const longMessage = "A".repeat(10000);
      const error = new Error(longMessage);

      expect(extractErrorMessage(error)).toBe(longMessage);
    });

    it("should handle Error with unicode characters", () => {
      const error = new Error("Erro com acentos: café, múltiplos");
      expect(extractErrorMessage(error)).toBe(
        "Erro com acentos: café, múltiplos",
      );
    });

    it("should handle Error with special characters", () => {
      const error = new Error("Special chars: <>&\"'@#$%");
      expect(extractErrorMessage(error)).toBe("Special chars: <>&\"'@#$%");
    });
  });
});

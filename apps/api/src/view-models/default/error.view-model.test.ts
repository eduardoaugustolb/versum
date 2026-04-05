import { describe, it, expect } from "vitest";
import { ApiErrorViewModel } from "./error.view-model";

describe("ApiErrorViewModel", () => {
  describe("constructor", () => {
    it("should create instance with error message only", () => {
      const vm = new ApiErrorViewModel("Something went wrong");

      expect(vm.success).toBe(false);
      expect(vm.error).toBe("Something went wrong");
      expect(vm.code).toBeUndefined();
    });

    it("should create instance with error message and code", () => {
      const vm = new ApiErrorViewModel("Not found", "NOT_FOUND");

      expect(vm.success).toBe(false);
      expect(vm.error).toBe("Not found");
      expect(vm.code).toBe("NOT_FOUND");
    });

    it("should create instance with empty string error", () => {
      const vm = new ApiErrorViewModel("");

      expect(vm.success).toBe(false);
      expect(vm.error).toBe("");
      expect(vm.code).toBeUndefined();
    });
  });



  describe("toJSON", () => {
    it("should return correct JSON with message only", () => {
      const vm = new ApiErrorViewModel("Something went wrong");
      const json = vm.toJSON();

      expect(json).toEqual({
        success: false,
        error: "Something went wrong",
      });
    });

    it("should return correct JSON with message and code", () => {
      const vm = new ApiErrorViewModel("Resource not found", "NOT_FOUND");
      const json = vm.toJSON();

      expect(json).toEqual({
        success: false,
        error: "Resource not found",
        code: "NOT_FOUND",
      });
    });

    it("should not include code in JSON when undefined", () => {
      const vm = new ApiErrorViewModel("Error");
      const json = vm.toJSON();

      expect(json).not.toHaveProperty("code");
    });

    it("should include code in JSON when provided", () => {
      const vm = new ApiErrorViewModel("Error", "ERROR_CODE");
      const json = vm.toJSON();

      expect(json).toHaveProperty("code", "ERROR_CODE");
    });

    it("should handle empty error message", () => {
      const vm = new ApiErrorViewModel("");
      const json = vm.toJSON();

      expect(json).toEqual({
        success: false,
        error: "",
      });
    });

    it("should handle special characters in error message", () => {
      const vm = new ApiErrorViewModel("Error with special chars: <>&\"'");
      const json = vm.toJSON();

      expect(json).toEqual({
        success: false,
        error: "Error with special chars: <>&\"'",
      });
    });

    it("should handle long error message", () => {
      const longMessage = "A".repeat(1000);
      const vm = new ApiErrorViewModel(longMessage);
      const json = vm.toJSON();

      expect(json.error).toBe(longMessage);
    });

    it("should handle various error codes", () => {
      const codes = ["VALIDATION_ERROR", "UNAUTHORIZED", "FORBIDDEN", "INTERNAL_ERROR", "404"];

      codes.forEach((code) => {
        const vm = new ApiErrorViewModel("Error", code);
        const json = vm.toJSON();

        expect(json.code).toBe(code);
      });
    });
  });
});

import { describe, it, expect } from "vitest";
import { parsePagination } from "./index";

describe("parsePagination", () => {
  describe("with valid inputs", () => {
    it("should parse valid page and limit", () => {
      const result = parsePagination({ page: "2", limit: "20" });

      expect(result).toEqual({ page: 2, limit: 20 });
    });

    it("should parse first page correctly", () => {
      const result = parsePagination({ page: "1", limit: "10" });

      expect(result).toEqual({ page: 1, limit: 10 });
    });

    it("should parse large values", () => {
      const result = parsePagination({ page: "100", limit: "100" });

      expect(result).toEqual({ page: 100, limit: 50 });
    });
  });

  describe("with default values", () => {
    it("should return defaults when query is empty", () => {
      const result = parsePagination({});

      expect(result).toEqual({ page: 1, limit: 10 });
    });

    it("should return defaults for undefined values", () => {
      const result = parsePagination({ page: undefined, limit: undefined });

      expect(result).toEqual({ page: 1, limit: 10 });
    });

    it("should return defaults for null values", () => {
      const result = parsePagination({ page: "", limit: "" });

      expect(result).toEqual({ page: 1, limit: 10 });
    });
  });

  describe("with invalid page values", () => {
    it("should use default for negative page", () => {
      const result = parsePagination({ page: "-1", limit: "10" });

      expect(result.page).toBe(1);
    });

    it("should use default for zero page", () => {
      const result = parsePagination({ page: "0", limit: "10" });

      expect(result.page).toBe(1);
    });

    it("should use default for non-numeric page", () => {
      const result = parsePagination({ page: "abc", limit: "10" });

      expect(result.page).toBe(1);
    });

    it("should use default for float page", () => {
      const result = parsePagination({ page: "1.5", limit: "10" });

      expect(result.page).toBe(1);
    });

    it("should use default for NaN page", () => {
      const result = parsePagination({ page: "NaN", limit: "10" });

      expect(result.page).toBe(1);
    });

    it("should use default for infinity page", () => {
      const result = parsePagination({ page: "Infinity", limit: "10" });

      expect(result.page).toBe(1);
    });
  });

  describe("with invalid limit values", () => {
    it("should use default for negative limit", () => {
      const result = parsePagination({ page: "1", limit: "-5" });

      expect(result.limit).toBe(10);
    });

    it("should use default for zero limit", () => {
      const result = parsePagination({ page: "1", limit: "0" });

      expect(result.limit).toBe(10);
    });

    it("should use default for non-numeric limit", () => {
      const result = parsePagination({ page: "1", limit: "abc" });

      expect(result.limit).toBe(10);
    });

    it("should use default for float limit", () => {
      const result = parsePagination({ page: "1", limit: "10.5" });

      expect(result.limit).toBe(10);
    });

    it("should cap limit at MAX_LIMIT (50)", () => {
      const result = parsePagination({ page: "1", limit: "100" });

      expect(result.limit).toBe(50);
    });

    it("should cap limit at MAX_LIMIT for values just over", () => {
      const result = parsePagination({ page: "1", limit: "51" });

      expect(result.limit).toBe(50);
    });

    it("should allow limit at exactly MAX_LIMIT", () => {
      const result = parsePagination({ page: "1", limit: "50" });

      expect(result.limit).toBe(50);
    });
  });

  describe("edge cases", () => {
    it("should handle very large numbers", () => {
      const result = parsePagination({ page: "999999", limit: "999999" });

      expect(result.page).toBe(999999);
      expect(result.limit).toBe(50);
    });

    it("should handle decimal strings that parse to integers", () => {
      const result = parsePagination({ page: "5.0", limit: "10.0" });

      expect(result.page).toBe(5);
      expect(result.limit).toBe(10);
    });

    it("should handle whitespace strings", () => {
      const result = parsePagination({ page: "  1  ", limit: "  10  " });

      expect(result.page).toBe(1);
      expect(result.limit).toBe(10);
    });

    it("should handle scientific notation", () => {
      const result = parsePagination({ page: "1e1", limit: "1e1" });

      expect(result.page).toBe(10);
      expect(result.limit).toBe(10);
    });
  });
});

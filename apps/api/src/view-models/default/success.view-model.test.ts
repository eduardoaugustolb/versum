import { describe, expect, it } from "vitest";
import { PaginationViewModel } from "./pagination.view-model";
import { SuccessViewModel } from "./success.view-model";

describe("SuccessViewModel", () => {
  describe("create", () => {
    it("should create instance with data only", () => {
      const vm = SuccessViewModel.create({ id: "123", name: "Test" });

      expect(vm.success).toBe(true);
      expect(vm.data).toEqual({ id: "123", name: "Test" });
      expect(vm.pagination).toBeUndefined();
      expect(vm.message).toBeUndefined();
    });

    it("should create instance with data and pagination", () => {
      const pagination = PaginationViewModel.create({
        page: 1,
        limit: 10,
        totalItems: 100,
      });

      const vm = SuccessViewModel.create({ items: ["a", "b"] }, pagination);

      expect(vm.success).toBe(true);
      expect(vm.data).toEqual({ items: ["a", "b"] });
      expect(vm.pagination).toBe(pagination);
      expect(vm.message).toBeUndefined();
    });

    it("should create instance with data and message", () => {
      const vm = SuccessViewModel.create(
        { id: "123" },
        undefined,
        "Operation successful",
      );

      expect(vm.success).toBe(true);
      expect(vm.data).toEqual({ id: "123" });
      expect(vm.pagination).toBeUndefined();
      expect(vm.message).toBe("Operation successful");
    });

    it("should create instance with all parameters", () => {
      const pagination = PaginationViewModel.create({
        page: 1,
        limit: 10,
        totalItems: 100,
      });

      const vm = SuccessViewModel.create(
        { items: ["a", "b"] },
        pagination,
        "Data retrieved successfully",
      );

      expect(vm.success).toBe(true);
      expect(vm.data).toEqual({ items: ["a", "b"] });
      expect(vm.pagination).toBe(pagination);
      expect(vm.message).toBe("Data retrieved successfully");
    });

    it("should create instance without data", () => {
      const vm = SuccessViewModel.create();

      expect(vm.success).toBe(true);
      expect(vm.data).toBeUndefined();
      expect(vm.pagination).toBeUndefined();
      expect(vm.message).toBeUndefined();
    });

    it("should handle null data", () => {
      const vm = SuccessViewModel.create(null);

      expect(vm.success).toBe(true);
      expect(vm.data).toBeNull();
    });

    it("should handle array data", () => {
      const vm = SuccessViewModel.create([1, 2, 3]);

      expect(vm.success).toBe(true);
      expect(vm.data).toEqual([1, 2, 3]);
    });

    it("should handle empty array data", () => {
      const vm = SuccessViewModel.create([]);

      expect(vm.success).toBe(true);
      expect(vm.data).toEqual([]);
    });

    it("should handle string data", () => {
      const vm = SuccessViewModel.create("hello");

      expect(vm.success).toBe(true);
      expect(vm.data).toBe("hello");
    });

    it("should handle number data", () => {
      const vm = SuccessViewModel.create(42);

      expect(vm.success).toBe(true);
      expect(vm.data).toBe(42);
    });

    it("should handle boolean data", () => {
      const vm = SuccessViewModel.create(true);

      expect(vm.success).toBe(true);
      expect(vm.data).toBe(true);
    });
  });

  describe("toJSON", () => {
    it("should return correct JSON with data only", () => {
      const vm = SuccessViewModel.create({ id: "123" });
      const json = vm.toJSON();

      expect(json).toEqual({
        success: true,
        data: { id: "123" },
      });
    });

    it("should return correct JSON with data and message", () => {
      const vm = SuccessViewModel.create({ id: "123" }, undefined, "Success!");
      const json = vm.toJSON();

      expect(json).toEqual({
        success: true,
        data: { id: "123" },
        message: "Success!",
      });
    });

    it("should return correct JSON with data and pagination", () => {
      const pagination = PaginationViewModel.create({
        page: 1,
        limit: 10,
        totalItems: 100,
      });

      const vm = SuccessViewModel.create({ items: ["a"] }, pagination);
      const json = vm.toJSON();

      expect(json).toEqual({
        success: true,
        data: { items: ["a"] },
        pagination: {
          currentPage: 1,
          totalPages: 10,
          totalItems: 100,
          itemsPerPage: 10,
          hasNextPage: true,
          hasPrevPage: false,
        },
      });
    });

    it("should return correct JSON with message only (no data)", () => {
      const vm = SuccessViewModel.create(
        undefined,
        undefined,
        "Operation completed",
      );
      const json = vm.toJSON();

      expect(json).toEqual({
        success: true,
        message: "Operation completed",
      });
    });

    it("should return correct JSON without optional fields", () => {
      const vm = SuccessViewModel.create();
      const json = vm.toJSON();

      expect(json).toEqual({
        success: true,
      });
    });

    it("should return correct JSON with all fields", () => {
      const pagination = PaginationViewModel.create({
        page: 2,
        limit: 10,
        totalItems: 50,
      });

      const vm = SuccessViewModel.create(
        { users: ["john", "jane"] },
        pagination,
        "Users retrieved",
      );
      const json = vm.toJSON();

      expect(json).toEqual({
        success: true,
        data: { users: ["john", "jane"] },
        pagination: {
          currentPage: 2,
          totalPages: 5,
          totalItems: 50,
          itemsPerPage: 10,
          hasNextPage: true,
          hasPrevPage: true,
        },
        message: "Users retrieved",
      });
    });
  });
});

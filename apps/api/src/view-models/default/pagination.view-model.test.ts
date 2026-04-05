import { describe, it, expect } from "vitest";
import { PaginationViewModel } from "./pagination.view-model";

describe("PaginationViewModel", () => {
  describe("create", () => {
    it("should create pagination with correct values", () => {
      const vm = PaginationViewModel.create({
        page: 1,
        limit: 10,
        totalItems: 100,
      });

      expect(vm.currentPage).toBe(1);
      expect(vm.itemsPerPage).toBe(10);
      expect(vm.totalItems).toBe(100);
      expect(vm.totalPages).toBe(10);
    });

    it("should calculate totalPages correctly for exact division", () => {
      const vm = PaginationViewModel.create({
        page: 1,
        limit: 10,
        totalItems: 100,
      });

      expect(vm.totalPages).toBe(10);
    });

    it("should calculate totalPages correctly for non-exact division", () => {
      const vm = PaginationViewModel.create({
        page: 1,
        limit: 10,
        totalItems: 25,
      });

      expect(vm.totalPages).toBe(3);
    });

    it("should return 1 page for 0 items", () => {
      const vm = PaginationViewModel.create({
        page: 1,
        limit: 10,
        totalItems: 0,
      });

      expect(vm.totalPages).toBe(1);
    });

    it("should return 1 page when items less than limit", () => {
      const vm = PaginationViewModel.create({
        page: 1,
        limit: 10,
        totalItems: 5,
      });

      expect(vm.totalPages).toBe(1);
    });

    it("should handle 1 item correctly", () => {
      const vm = PaginationViewModel.create({
        page: 1,
        limit: 10,
        totalItems: 1,
      });

      expect(vm.totalPages).toBe(1);
    });
  });

  describe("hasNextPage", () => {
    it("should return true when page is less than totalPages", () => {
      const vm = PaginationViewModel.create({
        page: 1,
        limit: 10,
        totalItems: 25,
      });

      expect(vm.hasNextPage).toBe(true);
    });

    it("should return false when page equals totalPages", () => {
      const vm = PaginationViewModel.create({
        page: 3,
        limit: 10,
        totalItems: 25,
      });

      expect(vm.hasNextPage).toBe(false);
    });

    it("should return false when there is only one page", () => {
      const vm = PaginationViewModel.create({
        page: 1,
        limit: 10,
        totalItems: 5,
      });

      expect(vm.hasNextPage).toBe(false);
    });

    it("should return false when totalItems is 0", () => {
      const vm = PaginationViewModel.create({
        page: 1,
        limit: 10,
        totalItems: 0,
      });

      expect(vm.hasNextPage).toBe(false);
    });
  });

  describe("hasPrevPage", () => {
    it("should return false when page is 1", () => {
      const vm = PaginationViewModel.create({
        page: 1,
        limit: 10,
        totalItems: 100,
      });

      expect(vm.hasPrevPage).toBe(false);
    });

    it("should return true when page is greater than 1", () => {
      const vm = PaginationViewModel.create({
        page: 2,
        limit: 10,
        totalItems: 100,
      });

      expect(vm.hasPrevPage).toBe(true);
    });

    it("should return true when page is last page", () => {
      const vm = PaginationViewModel.create({
        page: 10,
        limit: 10,
        totalItems: 100,
      });

      expect(vm.hasPrevPage).toBe(true);
    });
  });

  describe("toJSON", () => {
    it("should return correct JSON structure", () => {
      const vm = PaginationViewModel.create({
        page: 2,
        limit: 10,
        totalItems: 50,
      });

      const json = vm.toJSON();

      expect(json).toEqual({
        currentPage: 2,
        totalPages: 5,
        totalItems: 50,
        itemsPerPage: 10,
        hasNextPage: true,
        hasPrevPage: true,
      });
    });

    it("should return correct JSON for first page", () => {
      const vm = PaginationViewModel.create({
        page: 1,
        limit: 10,
        totalItems: 50,
      });

      const json = vm.toJSON();

      expect(json).toEqual({
        currentPage: 1,
        totalPages: 5,
        totalItems: 50,
        itemsPerPage: 10,
        hasNextPage: true,
        hasPrevPage: false,
      });
    });

    it("should return correct JSON for last page", () => {
      const vm = PaginationViewModel.create({
        page: 5,
        limit: 10,
        totalItems: 50,
      });

      const json = vm.toJSON();

      expect(json).toEqual({
        currentPage: 5,
        totalPages: 5,
        totalItems: 50,
        itemsPerPage: 10,
        hasNextPage: false,
        hasPrevPage: true,
      });
    });
  });
});

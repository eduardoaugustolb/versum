export interface PaginationInput {
  page: number;
  limit: number;
  totalItems: number;
}

export class PaginationViewModel {
  readonly currentPage: number;
  readonly totalPages: number;
  readonly totalItems: number;
  readonly itemsPerPage: number;
  readonly hasNextPage: boolean;
  readonly hasPrevPage: boolean;

  private constructor(input: PaginationInput) {
    const { page, limit, totalItems } = input;

    const totalPages = Math.max(1, Math.ceil(totalItems / limit));

    this.currentPage = page;
    this.totalItems = totalItems;
    this.itemsPerPage = limit;
    this.totalPages = totalPages;
    this.hasNextPage = page < totalPages;
    this.hasPrevPage = page > 1;
  }

  static create(input: PaginationInput) {
    return new PaginationViewModel(input);
  }

  toJSON() {
    return {
      currentPage: this.currentPage,
      totalPages: this.totalPages,
      totalItems: this.totalItems,
      itemsPerPage: this.itemsPerPage,
      hasNextPage: this.hasNextPage,
      hasPrevPage: this.hasPrevPage,
    };
  }
}

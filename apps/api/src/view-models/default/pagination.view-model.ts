export interface PaginationInput {
  page: number;
  limit: number;
  total_items: number;
}

export class PaginationViewModel {
  readonly current_page: number;
  readonly total_pages: number;
  readonly total_items: number;
  readonly items_per_page: number;
  readonly has_next_page: boolean;
  readonly has_prev_page: boolean;

  private constructor(input: PaginationInput) {
    const { page, limit, total_items } = input;

    const total_pages = Math.max(1, Math.ceil(total_items / limit));

    this.current_page = page;
    this.total_items = total_items;
    this.items_per_page = limit;
    this.total_pages = total_pages;
    this.has_next_page = page < total_pages;
    this.has_prev_page = page > 1;
  }

  static create(input: PaginationInput) {
    return new PaginationViewModel(input);
  }

  toJSON() {
    return {
      current_page: this.current_page,
      total_pages: this.total_pages,
      total_items: this.total_items,
      items_per_page: this.items_per_page,
      has_next_page: this.has_next_page,
      has_prev_page: this.has_prev_page,
    };
  }
}

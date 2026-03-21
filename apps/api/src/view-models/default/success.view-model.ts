import { PaginationViewModel } from "./pagination.view-model.ts";

export class SuccessViewModel<T> {
  success: boolean;
  data?: T;
  pagination?: PaginationViewModel;
  message?: string;

  constructor(data?: T, pagination?: PaginationViewModel, message?: string) {
    this.success = true;
    this.data = data;
    this.pagination = pagination;
    this.message = message;
  }

  static create<T>(
    data?: T,
    pagination?: PaginationViewModel,
    message?: string,
  ) {
    return new SuccessViewModel(data, pagination, message);
  }

  toJSON() {
    const response: {
      success: boolean;
      message?: string;
      data?: T;
      pagination?: ReturnType<PaginationViewModel["toJSON"]>;
    } = {
      success: this.success,
    };

    if (this.message) response.message = this.message;
    if (this.data) response.data = this.data;
    if (this.pagination) response.pagination = this.pagination.toJSON();

    return response;
  }
}

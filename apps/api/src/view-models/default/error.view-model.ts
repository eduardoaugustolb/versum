export class ApiErrorViewModel {
  success: boolean;
  error: string;
  code?: string;

  constructor(error: string, code?: string) {
    this.success = false;
    this.error = error;
    this.code = code;
  }

  toJSON() {
    const response: { success: boolean; error: string; code?: string } = {
      success: this.success,
      error: this.error,
    };

    if (this.code) response.code = this.code;

    return response;
  }
}

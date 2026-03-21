import { BaseError } from "./base.error-model.ts";

export class RateLimitError extends BaseError {
  constructor(message: string) {
    super(message);
    this.name = "RateLimitError";
  }
}

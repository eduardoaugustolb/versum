import { BaseError } from "./base.error-model.ts";

export class UnauthorizedError extends BaseError {
  constructor(message: string) {
    super(message);
    this.name = "UnauthorizedError";
  }
}

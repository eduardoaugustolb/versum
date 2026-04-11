import { BaseError } from "./base.error-model.ts";

export class ConflictError extends BaseError {
  constructor(message: string) {
    super(message);
    this.name = "ConflictError";
  }
}

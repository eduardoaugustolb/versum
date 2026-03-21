import { BaseError } from "./base.error-model.ts";

export class NotFoundError extends BaseError {
  constructor(message: string) {
    super(message);
    this.name = "NotFoundError";
  }
}

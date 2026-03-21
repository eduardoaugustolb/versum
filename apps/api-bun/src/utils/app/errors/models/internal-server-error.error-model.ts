import { BaseError } from "./base.error-model.ts";

export class InternalServerError extends BaseError {
  constructor(message: string) {
    super(message);
    this.name = "InternalServerError";
  }
}

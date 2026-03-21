import { BaseError } from "./base.error-model.ts";

export class BadRequestError extends BaseError {
  constructor(message: string) {
    super(message);
    this.name = "BadRequestError";
  }
}

export class ForbiddenError extends BaseError {
  constructor(message: string) {
    super(message);
    this.name = "ForbiddenError";
  }
}

import type { Context } from "hono";
import { logger } from "../../logger/index.ts";
import { ApiErrorViewModel } from "../../../view-models/default/error.view-model.ts";
import {
  BadRequestError,
  BaseError,
  ForbiddenError,
  InternalServerError,
  NotFoundError,
  UnauthorizedError,
} from "./index.ts";
import { extractErrorMessage, isValidationError } from "./utils.ts";

export class ErrorHandler {
  private readonly ctx: Context;

  constructor({ ctx }: { ctx: Context }) {
    this.ctx = ctx;
  }

  handle(error: unknown, details?: string) {
    if (error instanceof BaseError) {
      const statusCode = this.getStatusCodeForError(error);
      this.ctx.status(statusCode);
      return this.ctx.json(new ApiErrorViewModel(error.message, error.name));
    }

    if (isValidationError(error)) {
      this.ctx.status(400);
      return this.ctx.json(
        new ApiErrorViewModel(
          extractErrorMessage(error, "Validation error"),
          "VALIDATION_ERROR",
        ),
      );
    }

    logger(
      "error",
      details ? `Erro no ${details}:` : "Erro não tratado:",
      error,
    );

    this.ctx.status(500);
    return this.ctx.json(
      new ApiErrorViewModel(
        extractErrorMessage(error, "Internal Server Error"),
        "INTERNAL_SERVER_ERROR",
      ),
    );
  }
  private getStatusCodeForError(error: unknown) {
    if (error instanceof BadRequestError) return 400;
    if (error instanceof NotFoundError) return 404;
    if (error instanceof UnauthorizedError) return 401;
    if (error instanceof ForbiddenError) return 403;
    if (error instanceof InternalServerError) return 500;

    return 400;
  }
}

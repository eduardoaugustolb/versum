import type { Context } from "hono";
import { ApiErrorViewModel } from "../../../view-models/default/error.view-model.ts";
import { logger } from "@versum/logger";
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

  handle(err: unknown, details?: string) {
    if (err instanceof BaseError) {
      const statusCode = this.getStatusCodeForError(err);
      this.ctx.status(statusCode);
      return this.ctx.json(new ApiErrorViewModel(err.message, err.name));
    }

    if (isValidationError(err)) {
      this.ctx.status(400);
      return this.ctx.json(
        new ApiErrorViewModel(
          extractErrorMessage(err, "Validation error"),
          "VALIDATION_ERROR",
        ),
      );
    }

    logger("error",
      details ? `Erro no ${details}:` : "Erro não tratado:",
      String(err),
    );

    this.ctx.status(500);
    return this.ctx.json(
      new ApiErrorViewModel(
        extractErrorMessage(err, "Internal Server Error"),
        "INTERNAL_SERVER_ERROR",
      ),
    );
  }
  private getStatusCodeForError(err: unknown) {
    if (err instanceof BadRequestError) return 400;
    if (err instanceof NotFoundError) return 404;
    if (err instanceof UnauthorizedError) return 401;
    if (err instanceof ForbiddenError) return 403;
    if (err instanceof InternalServerError) return 500;

    return 400;
  }
}

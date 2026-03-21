import type { Context } from "hono";
import { ApiErrorViewModel } from "../../../view-models/default/error.view-model.ts";
import { extractErrorMessage } from "./utils.ts";

export const validationErrorHook = (
  result: { success: boolean; error?: unknown },
  c: Context,
) => {
  if (result.success) return;

  const message = extractErrorMessage(result.error, "Validation error");
  return c.json(new ApiErrorViewModel(message, "VALIDATION_ERROR"), 400);
};

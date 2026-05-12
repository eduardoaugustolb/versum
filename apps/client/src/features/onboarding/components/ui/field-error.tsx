import type { AnyFieldApi } from "@tanstack/react-form";

export function FieldError({ field }: { field: AnyFieldApi }) {
  return field.state.meta.isTouched && field.state.meta.errors.length ? (
    <p className="mt-1.5 text-md font-instrument-sans text-red-500">
      {field.state.meta.errors.map(e => e.message).join(", ")}
    </p>
  ) : null;
}

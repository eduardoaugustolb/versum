"use client";

import { memo } from "react";
import { cn } from "@/lib/utils";

interface ActionButtonProps {
  label: string;
  onClick?: () => void;
  type?: "button" | "submit";
  disabled?: boolean;
}

export const ActionButton = memo(function ActionButton({
  label,
  onClick,
  type = "button",
  disabled = false,
}: ActionButtonProps) {
  return (
    <button
      type={type}
      onClick={onClick}
      disabled={disabled}
      className={cn(
        "mt-8 w-full rounded-xl bg-neutral-900 px-6 py-3.5 text-sm font-medium text-white",
        "transition-colors hover:bg-neutral-700 active:scale-[.98]",
        "disabled:cursor-not-allowed disabled:opacity-40",
        "dark:bg-white dark:text-neutral-900 dark:hover:bg-neutral-200",
      )}
    >
      {label}
    </button>
  );
});

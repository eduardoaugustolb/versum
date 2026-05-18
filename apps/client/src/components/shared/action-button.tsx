"use client";

import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";

interface ActionButtonProps {
  label: string;
  onClick?: () => void;
  type?: "button" | "submit";
  disabled?: boolean;
}

export function ActionButton({
  label,
  onClick,
  type = "button",
  disabled = false,
}: ActionButtonProps) {
  return (
    <Button
      type={type}
      onClick={onClick}
      disabled={disabled}
      className={cn(
        "mt-8 w-full rounded-xl bg-neutral-900 px-6 py-3.5 text-sm font-medium text-white",
        "hover:bg-neutral-700 active:scale-[.98] h-12",
        "dark:bg-white dark:text-neutral-900 dark:hover:bg-neutral-200",
      )}
    >
      {label}
    </Button>
  );
}

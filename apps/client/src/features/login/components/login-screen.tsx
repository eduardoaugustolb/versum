"use client";

import { LoginForm } from "./login-form";
import { Toaster } from "@/components/ui/sonner";

export function LoginScreen() {
  return (
    <main className="relative mx-auto flex w-9/10 md:w-[720px] h-svh flex-col items-center justify-center overflow-hidden bg-background">
      <LoginForm />
      <Toaster />
    </main>
  );
}

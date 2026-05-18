import { useForm } from "@tanstack/react-form";
import { toast } from "sonner";
import { loginFormSchema } from "../types";
import { postApiV1AuthMagicLink } from "@/dal/orval/fetch/auth/auth";

interface UseLoginFormOptions {
  onSuccess?: () => void;
}

export function useLoginForm({ onSuccess }: UseLoginFormOptions = {}) {
  const form = useForm({
    defaultValues: { email: "" },
    validators: {
      onSubmit: loginFormSchema,
      onChange: loginFormSchema,
    },
    onSubmitInvalid() {
      toast.error("Por favor, corrija os erros no formulário.");
    },
    async onSubmit({ value }) {
      const toastId = toast.loading("Enviando magic link...");

      try {
        const res = await postApiV1AuthMagicLink({ email: value.email });
        toast.success(res.message, { id: toastId });
        form.reset();
        onSuccess?.();
      } catch (error) {
        let message = "Ocorreu um erro desconhecido.";

        if (error && typeof error === "object" && "response" in error) {
          const err = error as { response?: { status?: number } };
          if (err.response?.status === 429) {
            message =
              "Muitas tentativas. Por favor, aguarde 1 minuto antes de tentar novamente.";
          }
        } else if (error instanceof Error) {
          message = error.message;
        }

        toast.error(message, { id: toastId });
      }
    },
  });

  return form;
}

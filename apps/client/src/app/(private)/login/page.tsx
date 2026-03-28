"use client";

import { SignInIcon } from "@phosphor-icons/react";
import { useForm } from "@tanstack/react-form";
import { toast } from "sonner";
import { z } from "zod";
import { Button } from "@/components/ui/button";
import { Card, CardTitle } from "@/components/ui/card";
import {
  Field,
  FieldDescription,
  FieldGroup,
  FieldLabel,
  FieldLegend,
  FieldSet,
} from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { postApiV1AuthMagicLink } from "@/lib/kubb/gen";

const formSchema = z.object({
  email: z.email("Digite um email válido"),
});

export default function LoginPage() {
  const form = useForm({
    defaultValues: { email: "" },
    validators: {
      onSubmit: formSchema,
      onChange: formSchema,
    },
    onSubmitInvalid() {
      toast.error("Por favor, corrija os erros no formulário.");
    },
    async onSubmit({ value }) {
      const toastId = toast.loading("Enviando magic link...");

      try {
        const res = await postApiV1AuthMagicLink(
          { email: value.email },
          {
            baseURL: process.env.NEXT_PUBLIC_API_URL,
          },
        );
        toast.success(res.message, { id: toastId });
        form.reset();
      } catch (error) {
        const message =
          error instanceof Error
            ? error.message
            : "Ocorreu um erro desconhecido.";
        toast.error(message, { id: toastId });
      }
    },
  });

  return (
    <div className="min-h-screen flex items-center justify-center p-4 bg-background">
      <Card className="w-full max-w-md p-8 shadow-lg">
        <div className="text-center mb-2">
          <CardTitle className="text-2xl font-semibold">
            Bem-vindo de volta
          </CardTitle>
          <p className="text-sm text-muted-foreground">
            Enviaremos um link mágico para seu email
          </p>
        </div>

        <form
          onSubmit={(e) => {
            e.preventDefault();
            form.handleSubmit();
          }}
          className="space-y-6"
        >
          <FieldGroup>
            <FieldSet>
              <FieldLegend className="sr-only">
                Informações de login
              </FieldLegend>

              <form.Field name="email">
                {(field) => (
                  <Field>
                    <FieldLabel htmlFor={field.name}>Email</FieldLabel>
                    <Input
                      id={field.name}
                      name={field.name}
                      type="email"
                      placeholder="seu@email.com"
                      value={field.state.value}
                      onChange={(e) => field.handleChange(e.target.value)}
                      className="h-12"
                      autoComplete="email"
                      autoFocus
                    />
                    {field.state.meta.errors.length > 0 && (
                      <FieldDescription className="text-red-600">
                        {String(field.state.meta.errors[0])}
                      </FieldDescription>
                    )}
                  </Field>
                )}
              </form.Field>
            </FieldSet>
          </FieldGroup>

          <Button
            type="submit"
            variant="default"
            size="lg"
            className="w-full h-12 text-base cursor-pointer"
            disabled={form.state.isSubmitting}
          >
            <SignInIcon size={20} weight="bold" />
            {form.state.isSubmitting ? "Enviando..." : "Enviar magic link"}
          </Button>
        </form>

        <p className="text-xs text-center text-muted-foreground mt-6">
          Ao continuar, você concorda com nossos termos de uso e política de
          privacidade
        </p>
      </Card>
    </div>
  );
}
